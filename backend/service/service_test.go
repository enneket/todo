package service

import (
	"database/sql"
	"testing"
	"time"
	"todo/backend/db"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) {
	// Use a shared in-memory SQLite database. The `cache=shared` pragma makes
	// every connection in the pool see the same DB; without it each connection
	// in modernc.org/sqlite gets its own private :memory: database, and
	// CREATE TABLE / SELECT land on different connections at random.
	var err error
	db.DB, err = sql.Open("sqlite", "file::memory:?cache=shared&_pragma=foreign_keys(ON)")
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
		notified_at DATETIME,
		repeat TEXT DEFAULT '',
		tags TEXT DEFAULT '[]',
		project_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
	);`
	if _, err := db.DB.Exec(createTableSQL); err != nil {
		t.Fatalf("Failed to create todos table: %v", err)
	}

	createProjectsTableSQL := `CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		color TEXT DEFAULT '#64748B',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
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
	id, err := CreateTodo(CreateTodoParams{
		Title:       "Buy Milk",
		Description: "Groceries",
		Priority:    "high",
		DueDate:     &now,
		Tags:        []string{"shopping"},
		ProjectID:   &projIDInt,
	})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Test GetTodos
	todos, err := GetTodos("")
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
	todos, _ = GetTodos("")
	if !todos[0].Completed {
		t.Error("Expected todo to be completed")
	}

	// Test UpdateTodoDetails (full payload via UpdateTodoParams)
	newTitle := "Buy Almond Milk"
	newDesc := "Updated desc"
	newPriority := "low"
	newTags := []string{"food"}
	err = UpdateTodoDetails(int(id), UpdateTodoParams{
		Title:       &newTitle,
		Description: &newDesc,
		Priority:    &newPriority,
		Tags:        &newTags,
	})
	if err != nil {
		t.Fatalf("UpdateTodoDetails failed: %v", err)
	}
	todos, _ = GetTodos("")
	if todos[0].Title != "Buy Almond Milk" {
		t.Errorf("Expected updated title 'Buy Almond Milk', got '%s'", todos[0].Title)
	}
	if todos[0].Priority != "low" {
		t.Errorf("Expected priority 'low', got '%s'", todos[0].Priority)
	}
	if len(todos[0].Tags) != 1 || todos[0].Tags[0] != "food" {
		t.Errorf("Expected tags [food], got %v", todos[0].Tags)
	}

	// Test DeleteTodo
	err = DeleteTodo(int(id))
	if err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}
	todos, _ = GetTodos("")
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos after delete, got %d", len(todos))
	}
}

func TestTodoRepeat(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	now := time.Now()
	// Create todo with daily repeat
	id, err := CreateTodo(CreateTodoParams{
		Title:    "Repeat Task",
		Priority: "high",
		DueDate:  &now,
		Repeat:   "daily",
	})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Complete the task
	err = UpdateTodoStatus(int(id), true)
	if err != nil {
		t.Fatalf("UpdateTodoStatus failed: %v", err)
	}

	todos, err := GetTodos("")
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
	todoID, _ := CreateTodo(CreateTodoParams{Title: "Main Task"})
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

	// Test UpdateSubtask (fields updated together)
	title := "Subtask 1 Updated"
	completed := true
	err = UpdateSubtask(int(id), &title, &completed)
	if err != nil {
		t.Fatalf("UpdateSubtask failed: %v", err)
	}
	subtasks, _ = GetSubtasks(todoIDInt)
	if !subtasks[0].Completed {
		t.Error("Expected subtask to be completed")
	}
	if subtasks[0].Title != "Subtask 1 Updated" {
		t.Errorf("Expected subtask title 'Subtask 1 Updated', got '%s'", subtasks[0].Title)
	}

	// A partial update with only completed must not blank the title.
	completed = false
	err = UpdateSubtask(int(id), nil, &completed)
	if err != nil {
		t.Fatalf("UpdateSubtask (partial) failed: %v", err)
	}
	subtasks, _ = GetSubtasks(todoIDInt)
	if subtasks[0].Completed {
		t.Error("Expected subtask to be incomplete")
	}
	if subtasks[0].Title != "Subtask 1 Updated" {
		t.Errorf("Expected title to survive partial update, got '%s'", subtasks[0].Title)
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

func TestGetTodosBatchedSubtasks(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Three todos, each with its own subtasks. After GetTodos every todo must
	// end up with only its own subtasks — i.e. the batch IN-query joined them
	// correctly and didn't smear subtasks across todos.
	todo1, _ := CreateTodo(CreateTodoParams{Title: "Task 1"})
	todo2, _ := CreateTodo(CreateTodoParams{Title: "Task 2"})
	todo3, _ := CreateTodo(CreateTodoParams{Title: "Task 3"})

	CreateSubtask(int(todo1), "1.a")
	CreateSubtask(int(todo1), "1.b")
	CreateSubtask(int(todo2), "2.a")
	CreateSubtask(int(todo3), "3.a")
	CreateSubtask(int(todo3), "3.b")
	CreateSubtask(int(todo3), "3.c")

	todos, err := GetTodos("")
	if err != nil {
		t.Fatalf("GetTodos failed: %v", err)
	}
	if len(todos) != 3 {
		t.Fatalf("Expected 3 todos, got %d", len(todos))
	}

	counts := map[int]int{}
	titles := map[int][]string{}
	for _, td := range todos {
		counts[td.ID] = len(td.Subtasks)
		for _, s := range td.Subtasks {
			titles[td.ID] = append(titles[td.ID], s.Title)
		}
	}
	if counts[int(todo1)] != 2 {
		t.Errorf("todo1 expected 2 subtasks, got %d (titles=%v)", counts[int(todo1)], titles[int(todo1)])
	}
	if counts[int(todo2)] != 1 {
		t.Errorf("todo2 expected 1 subtask, got %d (titles=%v)", counts[int(todo2)], titles[int(todo2)])
	}
	if counts[int(todo3)] != 3 {
		t.Errorf("todo3 expected 3 subtasks, got %d (titles=%v)", counts[int(todo3)], titles[int(todo3)])
	}
}

// todoTitles extracts titles so search failures can print what actually came
// back instead of whole todo structs.
func todoTitles(todos []db.Todo) []string {
	titles := make([]string, len(todos))
	for i, td := range todos {
		titles[i] = td.Title
	}
	return titles
}

func TestGetTodosSearch(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	CreateTodo(CreateTodoParams{Title: "Buy Milk", Description: "groceries for breakfast", Tags: []string{"shopping"}})
	CreateTodo(CreateTodoParams{Title: "Write report", Description: "quarterly numbers", Tags: []string{"work"}})
	trashID, _ := CreateTodo(CreateTodoParams{Title: "Milk the cows"})
	if err := DeleteTodo(int(trashID)); err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}

	// Title match, ASCII case-insensitive; the trashed todo stays hidden.
	todos, err := GetTodos("MILK")
	if err != nil {
		t.Fatalf("GetTodos failed: %v", err)
	}
	if len(todos) != 1 || todos[0].Title != "Buy Milk" {
		t.Errorf("Expected only 'Buy Milk' for q=MILK, got %v", todoTitles(todos))
	}

	// Description match.
	todos, _ = GetTodos("quarterly")
	if len(todos) != 1 || todos[0].Title != "Write report" {
		t.Errorf("Expected 'Write report' via description match, got %v", todoTitles(todos))
	}

	// Tag match.
	todos, _ = GetTodos("shopping")
	if len(todos) != 1 || todos[0].Title != "Buy Milk" {
		t.Errorf("Expected 'Buy Milk' via tag match, got %v", todoTitles(todos))
	}

	// LIKE wildcards in the query must be treated as literal characters.
	CreateTodo(CreateTodoParams{Title: "Save 50%_now"})
	todos, _ = GetTodos(`50%_`)
	if len(todos) != 1 || todos[0].Title != "Save 50%_now" {
		t.Errorf("Expected literal '50%%_' match, got %v", todoTitles(todos))
	}

	// Blank / whitespace-only query behaves like no filter.
	todos, _ = GetTodos("   ")
	if len(todos) != 3 {
		t.Errorf("Expected all 3 todos for blank query, got %d (%v)", len(todos), todoTitles(todos))
	}
}

func TestTodoPartialUpdate(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Seed a todo with every field populated so we can detect any unwanted wipe.
	due := time.Now().Add(24 * time.Hour)
	id, err := CreateTodo(CreateTodoParams{
		Title:       "Original",
		Description: "Original desc",
		Priority:    "medium",
		DueDate:     &due,
		Repeat:      "weekly",
		Tags:        []string{"work", "urgent"},
	})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Update only the title — everything else must stay intact.
	newTitle := "Just the title"
	if err := UpdateTodoDetails(int(id), UpdateTodoParams{Title: &newTitle}); err != nil {
		t.Fatalf("UpdateTodoDetails failed: %v", err)
	}

	todos, _ := GetTodos("")
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
	td := todos[0]
	if td.Title != "Just the title" {
		t.Errorf("title not updated: got %q", td.Title)
	}
	if td.Description != "Original desc" {
		t.Errorf("description was wiped: got %q", td.Description)
	}
	if td.Priority != "medium" {
		t.Errorf("priority was wiped: got %q", td.Priority)
	}
	if td.Repeat != "weekly" {
		t.Errorf("repeat was wiped: got %q", td.Repeat)
	}
	if len(td.Tags) != 2 {
		t.Errorf("tags were wiped: got %v", td.Tags)
	}
	if td.DueDate == nil || !td.DueDate.Equal(due) {
		t.Errorf("due_date was wiped: got %v want %v", td.DueDate, due)
	}

	// Empty params is a valid no-op (no SET clauses emitted).
	if err := UpdateTodoDetails(int(id), UpdateTodoParams{}); err != nil {
		t.Errorf("empty UpdateTodoParams should not error, got %v", err)
	}
}

func TestNotificationDedup(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Stub out the actual desktop-notification dispatch — beeep.Notify blocks
	// on dbus in headless environments.
	origNotify := sendNotification
	sendNotification = func(title, description string) error { return nil }
	defer func() { sendNotification = origNotify }()

	// Create a todo whose remind_at falls inside the [now, now+1m) window.
	remindAt := time.Now().UTC().Add(10 * time.Second)
	id, err := CreateTodo(CreateTodoParams{
		Title:    "Reminder",
		Priority: "high",
		RemindAt: &remindAt,
	})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// First scan: should mark notified_at.
	checkReminders()
	todos, _ := GetTodos("")
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
	if todos[0].NotifiedAt == nil {
		t.Fatal("Expected notified_at to be set after first checkReminders")
	}
	firstNotified := *todos[0].NotifiedAt

	// Second scan: query filters out notified_at IS NOT NULL, so notified_at
	// should stay unchanged even if the window keeps sliding.
	time.Sleep(20 * time.Millisecond)
	checkReminders()
	todos, _ = GetTodos("")
	if todos[0].NotifiedAt == nil || !todos[0].NotifiedAt.Equal(firstNotified) {
		t.Errorf("notified_at should not be overwritten by a subsequent scan; got %v want %v", todos[0].NotifiedAt, firstNotified)
	}

	// After UpdateTodoDetails, notified_at should reset so the user gets reminded again.
	if err := UpdateTodoDetails(int(id), UpdateTodoParams{RemindAt: &remindAt}); err != nil {
		t.Fatalf("UpdateTodoDetails failed: %v", err)
	}
	todos, _ = GetTodos("")
	if todos[0].NotifiedAt != nil {
		t.Errorf("Expected notified_at to be reset to NULL after UpdateTodoDetails, got %v", todos[0].NotifiedAt)
	}
}

func TestTrashTodoLifecycle(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	id, err := CreateTodo(CreateTodoParams{Title: "Doomed"})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Soft delete moves the todo from the normal list to the trash.
	if err := DeleteTodo(int(id)); err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}
	todos, _ := GetTodos("")
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos after soft delete, got %d", len(todos))
	}
	trashed, err := GetTrashedTodos()
	if err != nil {
		t.Fatalf("GetTrashedTodos failed: %v", err)
	}
	if len(trashed) != 1 || trashed[0].ID != int(id) {
		t.Fatalf("Expected 1 trashed todo (id=%d), got %v", id, trashed)
	}
	if trashed[0].DeletedAt == nil {
		t.Error("Expected trashed todo to carry a deleted_at timestamp")
	}

	// Restore brings it back; the trash empties.
	if err := RestoreTodo(int(id)); err != nil {
		t.Fatalf("RestoreTodo failed: %v", err)
	}
	todos, _ = GetTodos("")
	if len(todos) != 1 || todos[0].ID != int(id) {
		t.Errorf("Expected the todo back in the normal list, got %v", todos)
	}
	trashed, _ = GetTrashedTodos()
	if len(trashed) != 0 {
		t.Errorf("Expected empty trash after restore, got %d", len(trashed))
	}
}

func TestTrashProjectCascade(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	projID, _ := CreateProject("Doomed Project", "", "")
	projIDInt := int(projID)
	todoID, _ := CreateTodo(CreateTodoParams{Title: "Under project", ProjectID: &projIDInt})

	// Deleting the project takes its todos into the trash with it.
	if err := DeleteProject(projIDInt); err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
	projects, _ := GetProjects()
	if len(projects) != 0 {
		t.Errorf("Expected 0 projects after soft delete, got %d", len(projects))
	}
	todos, _ := GetTodos("")
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos after project soft delete, got %d", len(todos))
	}
	trashedProjects, err := GetTrashedProjects()
	if err != nil {
		t.Fatalf("GetTrashedProjects failed: %v", err)
	}
	if len(trashedProjects) != 1 || trashedProjects[0].ID != projIDInt {
		t.Fatalf("Expected 1 trashed project, got %v", trashedProjects)
	}
	trashedTodos, err := GetTrashedTodos()
	if err != nil {
		t.Fatalf("GetTrashedTodos failed: %v", err)
	}
	if len(trashedTodos) != 1 || trashedTodos[0].ID != int(todoID) {
		t.Fatalf("Expected 1 trashed todo, got %v", trashedTodos)
	}
	// Shared deletion timestamp keeps the group consistent.
	if !trashedProjects[0].DeletedAt.Equal(*trashedTodos[0].DeletedAt) {
		t.Error("Expected project and todo to share the same deleted_at timestamp")
	}

	// Restoring the project leaves its todos in the trash — each item keeps
	// its own deleted_at and is restored individually (RestoreTodo).
	if err := RestoreProject(projIDInt); err != nil {
		t.Fatalf("RestoreProject failed: %v", err)
	}
	projects, _ = GetProjects()
	if len(projects) != 1 || projects[0].ID != projIDInt {
		t.Errorf("Expected the project back in the normal list, got %v", projects)
	}
	todos, _ = GetTodos("")
	if len(todos) != 0 {
		t.Errorf("Expected todos to stay in the trash after project restore, got %d todos", len(todos))
	}
	trashedProjects, _ = GetTrashedProjects()
	trashedTodos, _ = GetTrashedTodos()
	if len(trashedProjects) != 0 {
		t.Errorf("Expected empty project trash after restore, got %d", len(trashedProjects))
	}
	if len(trashedTodos) != 1 || trashedTodos[0].ID != int(todoID) {
		t.Errorf("Expected the todo still in the trash after project restore, got %v", trashedTodos)
	}

	// Restoring the todo then brings it back; the project is already active
	// so no cascade is triggered.
	if err := RestoreTodo(int(todoID)); err != nil {
		t.Fatalf("RestoreTodo failed: %v", err)
	}
	todos, _ = GetTodos("")
	if len(todos) != 1 || todos[0].ID != int(todoID) {
		t.Errorf("Expected the todo back in the normal list, got %v", todos)
	}
	trashedTodos, _ = GetTrashedTodos()
	if len(trashedTodos) != 0 {
		t.Errorf("Expected empty todo trash after restore, got %d", len(trashedTodos))
	}
}

func TestRestoreTodoRestoresProject(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	projID, _ := CreateProject("Grouped", "", "")
	projIDInt := int(projID)
	todoID, _ := CreateTodo(CreateTodoParams{Title: "Member", ProjectID: &projIDInt})

	if err := DeleteProject(projIDInt); err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}

	// Restoring the todo must pull its trashed project back too, so no todo
	// ever points at a deleted project.
	if err := RestoreTodo(int(todoID)); err != nil {
		t.Fatalf("RestoreTodo failed: %v", err)
	}
	projects, _ := GetProjects()
	todos, _ := GetTodos("")
	if len(projects) != 1 || projects[0].ID != projIDInt {
		t.Errorf("Expected project restored along with its todo, got %v", projects)
	}
	if len(todos) != 1 || todos[0].ID != int(todoID) {
		t.Errorf("Expected todo restored, got %v", todos)
	}
}

func TestPurgeTodoCascadesSubtasks(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	id, _ := CreateTodo(CreateTodoParams{Title: "Purge me"})
	subtaskID, _ := CreateSubtask(int(id), "child")

	if err := DeleteTodo(int(id)); err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}
	if err := PurgeTodo(int(id)); err != nil {
		t.Fatalf("PurgeTodo failed: %v", err)
	}

	trashed, _ := GetTrashedTodos()
	if len(trashed) != 0 {
		t.Errorf("Expected empty trash after purge, got %d", len(trashed))
	}
	// The todo row and its subtasks are physically gone.
	var count int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ?", id).Scan(&count); err != nil || count != 0 {
		t.Errorf("Expected todo row physically deleted, count=%d err=%v", count, err)
	}
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM subtasks WHERE id = ?", subtaskID).Scan(&count); err != nil || count != 0 {
		t.Errorf("Expected subtask cascade-deleted, count=%d err=%v", count, err)
	}
}

func TestPurgeProjectCascadesTodos(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	projID, _ := CreateProject("Purge project", "", "")
	projIDInt := int(projID)
	todoID, _ := CreateTodo(CreateTodoParams{Title: "Inside", ProjectID: &projIDInt})

	if err := DeleteProject(projIDInt); err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
	if err := PurgeProject(projIDInt); err != nil {
		t.Fatalf("PurgeProject failed: %v", err)
	}

	trashedProjects, _ := GetTrashedProjects()
	trashedTodos, _ := GetTrashedTodos()
	if len(trashedProjects) != 0 || len(trashedTodos) != 0 {
		t.Errorf("Expected empty trash after purge, got %d projects, %d todos", len(trashedProjects), len(trashedTodos))
	}
	var count int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM projects WHERE id = ?", projID).Scan(&count); err != nil || count != 0 {
		t.Errorf("Expected project row physically deleted, count=%d err=%v", count, err)
	}
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ?", todoID).Scan(&count); err != nil || count != 0 {
		t.Errorf("Expected todo row physically deleted, count=%d err=%v", count, err)
	}
}

func TestNotificationSkipsTrashed(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Stub out the actual desktop-notification dispatch (same as
	// TestNotificationDedup) so the scan does not block on dbus.
	origNotify := sendNotification
	sendNotification = func(title, description string) error { return nil }
	defer func() { sendNotification = origNotify }()

	remindAt := time.Now().UTC().Add(10 * time.Second)
	id, err := CreateTodo(CreateTodoParams{
		Title:    "Ghost reminder",
		Priority: "high",
		RemindAt: &remindAt,
	})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// A trashed todo must never fire a reminder.
	if err := DeleteTodo(int(id)); err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}
	checkReminders()

	var notifiedAt *time.Time
	if err := db.DB.QueryRow("SELECT notified_at FROM todos WHERE id = ?", id).Scan(&notifiedAt); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if notifiedAt != nil {
		t.Error("Expected trashed todo to be skipped by the reminder scan, got notified_at set")
	}
}
