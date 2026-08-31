package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"todo/backend/db"
	"todo/backend/service"
)

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := service.GetTodos()
	if err != nil {
		log.Printf("GetTodos failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if todos == nil {
		todos = []db.Todo{}
	}
	writeJSON(w, todos)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Priority    string     `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
		RemindAt    *time.Time `json:"remind_at"`
		Repeat      string     `json:"repeat"`
		Tags        []string   `json:"tags"`
		ProjectID   *int       `json:"project_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := service.CreateTodo(service.CreateTodoParams{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		RemindAt:    req.RemindAt,
		Repeat:      req.Repeat,
		Tags:        req.Tags,
		ProjectID:   req.ProjectID,
	})
	if err != nil {
		log.Printf("CreateTodo failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// Return the full todo so the client can prepend it locally without an
	// extra GET /api/todos round-trip.
	todo, err := service.GetTodo(int(id))
	if err != nil {
		log.Printf("GetTodo(%d) after create failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, todo)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req struct {
		Completed   *bool      `json:"completed"`
		Title       *string    `json:"title"`
		Description *string    `json:"description"`
		Priority    *string    `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
		RemindAt    *time.Time `json:"remind_at"`
		Repeat      *string    `json:"repeat"`
		Tags        *[]string  `json:"tags"`
		ProjectID   *int       `json:"project_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Completed != nil {
		if err := service.UpdateTodoStatus(id, *req.Completed); err != nil {
			log.Printf("UpdateTodoStatus(%d) failed: %v", id, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	// Partial update: any non-nil field gets SET, the rest stay as they are.
	params := service.UpdateTodoParams{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		RemindAt:    req.RemindAt,
		Repeat:      req.Repeat,
		Tags:        req.Tags,
		ProjectID:   req.ProjectID,
	}
	if params.Title != nil || params.Description != nil || params.Priority != nil ||
		params.DueDate != nil || params.RemindAt != nil || params.Repeat != nil ||
		params.Tags != nil || params.ProjectID != nil {
		if err := service.UpdateTodoDetails(id, params); err != nil {
			log.Printf("UpdateTodoDetails(%d) failed: %v", id, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	// Return the updated todo in full so the client can patch it in place
	// without re-GETting the list.
	todo, err := service.GetTodo(id)
	if err != nil {
		log.Printf("GetTodo(%d) after update failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, todo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	if err := service.DeleteTodo(id); err != nil {
		log.Printf("DeleteTodo(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
