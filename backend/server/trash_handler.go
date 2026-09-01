package server

import (
	"log"
	"net/http"
	"todo/backend/db"
	"todo/backend/service"
)

// TrashTodosHandler returns the soft-deleted todos (the trash list).
func TrashTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := service.GetTrashedTodos()
	if err != nil {
		log.Printf("GetTrashedTodos failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if todos == nil {
		todos = []db.Todo{}
	}
	writeJSON(w, todos)
}

// TrashProjectsHandler returns the soft-deleted projects (the trash list).
func TrashProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := service.GetTrashedProjects()
	if err != nil {
		log.Printf("GetTrashedProjects failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if projects == nil {
		projects = []db.Project{}
	}
	writeJSON(w, projects)
}

// RestoreTodoHandler moves a trashed todo back into the normal list. If its
// project is also in the trash, the whole group is restored (see service.
// RestoreTodo).
func RestoreTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := service.RestoreTodo(id); err != nil {
		log.Printf("RestoreTodo(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RestoreProjectHandler moves a trashed project back into the normal projects
// list. Its todos stay in the trash and are restored individually.
func RestoreProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := service.RestoreProject(id); err != nil {
		log.Printf("RestoreProject(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PurgeTodoHandler permanently deletes a trashed todo. Only works on rows
// already in the trash — normal todos are never deleted this way.
func PurgeTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := service.PurgeTodo(id); err != nil {
		log.Printf("PurgeTodo(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PurgeProjectHandler permanently deletes a trashed project and its todos.
func PurgeProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := service.PurgeProject(id); err != nil {
		log.Printf("PurgeProject(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}