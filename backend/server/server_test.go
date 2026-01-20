package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/backend/db"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) {
	var err error
	db.DB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

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

func TestGetTodosHandler(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	req, _ := http.NewRequest("GET", "/api/todos", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTodosHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Should be empty array
	expected := "[]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateTodoHandler(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	todo := map[string]interface{}{
		"title":    "Test Todo",
		"priority": "high",
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTodoHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]int
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response["id"] == 0 {
		t.Error("Expected valid ID")
	}
}

func TestProjectHandlers(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Create
	project := map[string]string{
		"name":  "Test Project",
		"color": "#000000",
	}
	body, _ := json.Marshal(project)
	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	http.HandlerFunc(CreateProjectHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("CreateProjectHandler returned wrong status: %v", rr.Code)
	}

	// Get
	req, _ = http.NewRequest("GET", "/api/projects", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetProjectsHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GetProjectsHandler returned wrong status: %v", rr.Code)
	}
	
	var projects []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &projects)
	if len(projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects))
	}
}

func TestSubtaskHandlers(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Need a todo first
	todo := map[string]interface{}{"title": "Main Task"}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	http.HandlerFunc(CreateTodoHandler).ServeHTTP(rr, req)
	
	var todoResp map[string]int
	json.Unmarshal(rr.Body.Bytes(), &todoResp)
	todoID := todoResp["id"]

	// Create Subtask
	subtask := map[string]string{"title": "Sub 1"}
	body, _ = json.Marshal(subtask)
	// PathValue not supported in httptest nicely for default mux without routing
	// But our handler uses r.PathValue("id").
	// We need to set it manually or use the mux.
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/todos/{id}/subtasks", CreateSubtaskHandler)
	mux.HandleFunc("PUT /api/subtasks/{id}", UpdateSubtaskHandler)
	mux.HandleFunc("DELETE /api/subtasks/{id}", DeleteSubtaskHandler)

	url := "/api/todos/" + string(rune(todoID+'0')) + "/subtasks" // hacky int to string for small int
	// Better:
	// url := fmt.Sprintf("/api/todos/%d/subtasks", todoID)
	// But I don't have fmt imported. Let's use hardcoded "1" since it's fresh DB.
	url = "/api/todos/1/subtasks"

	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("CreateSubtaskHandler returned wrong status: %v", rr.Code)
	}
	
	var subResp map[string]int
	json.Unmarshal(rr.Body.Bytes(), &subResp)
	subID := subResp["id"]
	if subID == 0 {
		t.Fatal("Expected valid subtask ID")
	}

	// Update Subtask
	update := map[string]interface{}{"completed": true}
	body, _ = json.Marshal(update)
	req, _ = http.NewRequest("PUT", "/api/subtasks/1", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("UpdateSubtaskHandler returned wrong status: %v", rr.Code)
	}

	// Delete Subtask
	req, _ = http.NewRequest("DELETE", "/api/subtasks/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("DeleteSubtaskHandler returned wrong status: %v", rr.Code)
	}
}

func TestUpdateDeleteTodoHandler(t *testing.T) {
	setupTestDB(t)
	defer db.DB.Close()

	// Create Todo
	todo := map[string]interface{}{"title": "To Update"}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	http.HandlerFunc(CreateTodoHandler).ServeHTTP(rr, req)

	// Mux for path values
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/todos/{id}", UpdateTodoHandler)
	mux.HandleFunc("DELETE /api/todos/{id}", DeleteTodoHandler)

	// Update
	update := map[string]interface{}{"completed": true}
	body, _ = json.Marshal(update)
	req, _ = http.NewRequest("PUT", "/api/todos/1", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("UpdateTodoHandler returned wrong status: %v", rr.Code)
	}

	// Delete
	req, _ = http.NewRequest("DELETE", "/api/todos/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("DeleteTodoHandler returned wrong status: %v", rr.Code)
	}
}
