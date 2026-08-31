package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo/backend/db"
	"todo/backend/service"
)

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := service.GetTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todos == nil {
		todos = []db.Todo{}
	}
	json.NewEncoder(w).Encode(todos)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := service.CreateTodo(req.Title, req.Description, req.Priority, req.DueDate, req.RemindAt, req.Repeat, req.Tags, req.ProjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Completed != nil {
		if err := service.UpdateTodoStatus(id, *req.Completed); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	if err := service.DeleteTodo(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
