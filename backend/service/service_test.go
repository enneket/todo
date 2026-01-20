package service

import (
	"database/sql"
	"testing"
	"time"
	"todo/backend/db"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) {
	// Use in-memory SQLite database for testing
	var err error
	db.DB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Initialize tables
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		completed BOOLEAN DEFAULT FALSE,
		priority TEXT DEFAULT 'medium',
		due_date DATETIME,
		remind_at DATETIME,
		repeat TEXT DEFAULT '',
		tags TEXT DEFAULT '[]',
		project_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.DB.Exec(createTableSQL); err != nil {
		t.Fatalf("Failed to create todos table: %v", err)
	}

	createProjectsTableSQL := `CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		color TEXT DEFAULT '#64748B',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.DB.Exec(createProjectsTableSQL); err != nil {
		t.Fatalf("Failed to create projects table: %v", err)
	}

	createSubtasksTableSQL := `CREATE TABLE IF NOT EXISTS subtasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		todo_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(todo_id) REFERENCES todos(id) ON DELETE CASCADE
	);`
	if _, err := db.DB.Exec(createSubtasksTableSQL); err != nil {
		t.Fatalf("Failed to create subtasks table: %v", err)
	}
}

func TestProjectService(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Test CreateProject
	id, err := CreateProject("Work", "Work related tasks", "#EF4444")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}
	if id == 0 {
		t.Fatal("CreateProject returned invalid ID")
	}

	// Test GetProjects
	projects, err := GetProjects()
	if err != nil {
		t.Fatalf("GetProjects failed: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "Work" {
		t.Errorf("Expected project name 'Work', got '%s'", projects[0].Name)
	}

	// Test UpdateProject
	err = UpdateProject(int(id), "Work Updated", "Updated desc", "#000000")
	if err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}

	projects, _ = GetProjects()
	if projects[0].Name != "Work Updated" {
		t.Errorf("Expected updated project name 'Work Updated', got '%s'", projects[0].Name)
	}

	// Test DeleteProject
	err = DeleteProject(int(id))
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}

	projects, _ = GetProjects()
	if len(projects) != 0 {
		t.Errorf("Expected 0 projects after delete, got %d", len(projects))
	}
}

func TestTodoService(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Create a project first
	projID, _ := CreateProject("Test Project", "", "")
	projIDInt := int(projID)

	// Test CreateTodo
	now := time.Now()
	id, err := CreateTodo("Buy Milk", "Groceries", "high", &now, nil, "", []string{"shopping"}, &projIDInt)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Test GetTodos
	todos, err := GetTodos()
	if err != nil {
		t.Fatalf("GetTodos failed: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
	if todos[0].Title != "Buy Milk" {
		t.Errorf("Expected todo title 'Buy Milk', got '%s'", todos[0].Title)
	}
	if todos[0].ProjectID == nil || *todos[0].ProjectID != projIDInt {
		t.Errorf("Expected ProjectID %d, got %v", projIDInt, todos[0].ProjectID)
	}

	// Test UpdateTodoStatus
	err = UpdateTodoStatus(int(id), true)
	if err != nil {
		t.Fatalf("UpdateTodoStatus failed: %v", err)
	}
	todos, _ = GetTodos()
	if !todos[0].Completed {
		t.Error("Expected todo to be completed")
	}

	// Test UpdateTodoDetails
	err = UpdateTodoDetails(int(id), "Buy Almond Milk", "Updated desc", "low", nil, nil, "", []string{"food"}, nil)
	if err != nil {
		t.Fatalf("UpdateTodoDetails failed: %v", err)
	}
	todos, _ = GetTodos()
	if todos[0].Title != "Buy Almond Milk" {
		t.Errorf("Expected updated title 'Buy Almond Milk', got '%s'", todos[0].Title)
	}

	// Test DeleteTodo
	err = DeleteTodo(int(id))
	if err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}
	todos, _ = GetTodos()
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos after delete, got %d", len(todos))
	}
}

func TestTodoRepeat(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	now := time.Now()
	// Create todo with daily repeat
	id, err := CreateTodo("Repeat Task", "", "high", &now, nil, "daily", nil, nil)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Complete the task
	err = UpdateTodoStatus(int(id), true)
	if err != nil {
		t.Fatalf("UpdateTodoStatus failed: %v", err)
	}

	todos, err := GetTodos()
	// Should have 2 todos now: one completed (original), one pending (new)
	if len(todos) != 2 {
		t.Fatalf("Expected 2 todos, got %d", len(todos))
	}

	var newTodo, oldTodo db.Todo
	foundOld := false
	foundNew := false

	for _, todo := range todos {
		if todo.ID == int(id) {
			oldTodo = todo
			foundOld = true
		} else {
			newTodo = todo
			foundNew = true
		}
	}

	if !foundOld || !foundNew {
		t.Fatalf("Could not distinguish old and new todos. IDs: %d, %d", todos[0].ID, todos[1].ID)
	}

	if !oldTodo.Completed {
		t.Error("Original todo should be completed")
	}
	if newTodo.Completed {
		t.Error("New todo should be pending")
	}
	if newTodo.Title != "Repeat Task" {
		t.Errorf("New todo title mismatch: %s", newTodo.Title)
	}
	// Check due date: should be +1 day
	// Note: time comparisons can be tricky with nanoseconds, so we compare year/day
	expected := now.AddDate(0, 0, 1)
	if newTodo.DueDate == nil || newTodo.DueDate.Year() != expected.Year() || newTodo.DueDate.YearDay() != expected.YearDay() {
		t.Errorf("New todo due date incorrect. Expected ~%v, got %v", expected, newTodo.DueDate)
	}
}

func TestSubtaskService(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Create a todo first
	todoID, _ := CreateTodo("Main Task", "", "medium", nil, nil, "", nil, nil)
	todoIDInt := int(todoID)

	// Test CreateSubtask
	id, err := CreateSubtask(todoIDInt, "Subtask 1")
	if err != nil {
		t.Fatalf("CreateSubtask failed: %v", err)
	}

	// Test GetSubtasks
	subtasks, err := GetSubtasks(todoIDInt)
	if err != nil {
		t.Fatalf("GetSubtasks failed: %v", err)
	}
	if len(subtasks) != 1 {
		t.Fatalf("Expected 1 subtask, got %d", len(subtasks))
	}
	if subtasks[0].Title != "Subtask 1" {
		t.Errorf("Expected subtask title 'Subtask 1', got '%s'", subtasks[0].Title)
	}

	// Test UpdateSubtask
	err = UpdateSubtask(int(id), "Subtask 1 Updated", true)
	if err != nil {
		t.Fatalf("UpdateSubtask failed: %v", err)
	}
	subtasks, _ = GetSubtasks(todoIDInt)
	if !subtasks[0].Completed {
		t.Error("Expected subtask to be completed")
	}

	// Test DeleteSubtask
	err = DeleteSubtask(int(id))
	if err != nil {
		t.Fatalf("DeleteSubtask failed: %v", err)
	}
	subtasks, _ = GetSubtasks(todoIDInt)
	if len(subtasks) != 0 {
		t.Errorf("Expected 0 subtasks after delete, got %d", len(subtasks))
	}
}
