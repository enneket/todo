package server

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/backend/db"
	"todo/backend/service"
)

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := service.GetProjects()
	if err != nil {
		log.Printf("GetProjects failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if projects == nil {
		projects = []db.Project{}
	}
	writeJSON(w, projects)
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := service.CreateProject(req.Name, req.Description, req.Color)
	if err != nil {
		log.Printf("CreateProject failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]int64{"id": id})
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := service.UpdateProject(id, req.Name, req.Description, req.Color); err != nil {
		log.Printf("UpdateProject(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}

	if err := service.DeleteProject(id); err != nil {
		log.Printf("DeleteProject(%d) failed: %v", id, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
