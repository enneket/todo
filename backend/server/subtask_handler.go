package server

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/backend/service"
)

func CreateSubtaskHandler(w http.ResponseWriter, r *http.Request) {
	todoID, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := service.CreateSubtask(todoID, req.Title)
	if err != nil {
		log.Printf("CreateSubtask(todo=%d) failed: %v", todoID, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// Fetch the new subtask in full so the client can append it to its
	// in-memory copy without re-GETting the parent todo.
	sub, err := service.GetSubtask(int(id))
	if err != nil {
		log.Printf("GetSubtask(%d) after create failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, sub)
}

func UpdateSubtaskHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// req fields are pointers, so UpdateSubtask only updates the fields the
	// client actually sent — a partial update (e.g. toggling completion only)
	// is expected to send just that field.
	if err := service.UpdateSubtask(id, req.Title, req.Completed); err != nil {
		log.Printf("UpdateSubtask(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteSubtaskHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	if err := service.DeleteSubtask(id); err != nil {
		log.Printf("DeleteSubtask(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
