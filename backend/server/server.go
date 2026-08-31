package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// pathID extracts an int path parameter and writes a 400 if it's missing or
// non-numeric. Returns the id and ok=true on success.
func pathID(w http.ResponseWriter, r *http.Request, param string) (int, bool) {
	raw := r.PathValue(param)
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// maxBodyBytes caps the size of any request body the API will read.
const maxBodyBytes = 1 << 20 // 1 MiB

// writeJSON marshals v and writes it. The encoding step is the one place
// where we can fail mid-response — typically only if the client has
// disconnected. We can't change the status code at that point, but we can
// log it so we know something went wrong.
func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("response encode failed: %v", err)
	}
}

// Server owns the http.Server instance for the app. Holding it on a struct
// (instead of a package-level var) makes shutdown lifecycle explicit and
// lets multiple servers exist in tests without global state leaks.
type Server struct {
	srv *http.Server
}

// NewServer builds the server with all routes wired up. Call Start to begin
// listening; call Stop to shut it down.
func NewServer(port string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/todos", GetTodosHandler)
	mux.HandleFunc("POST /api/todos", CreateTodoHandler)
	mux.HandleFunc("PUT /api/todos/{id}", UpdateTodoHandler)
	mux.HandleFunc("DELETE /api/todos/{id}", DeleteTodoHandler)

	mux.HandleFunc("GET /api/projects", GetProjectsHandler)
	mux.HandleFunc("POST /api/projects", CreateProjectHandler)
	mux.HandleFunc("PUT /api/projects/{id}", UpdateProjectHandler)
	mux.HandleFunc("DELETE /api/projects/{id}", DeleteProjectHandler)

	mux.HandleFunc("POST /api/todos/{id}/subtasks", CreateSubtaskHandler)
	mux.HandleFunc("PUT /api/subtasks/{id}", UpdateSubtaskHandler)
	mux.HandleFunc("DELETE /api/subtasks/{id}", DeleteSubtaskHandler)

	return &Server{
		srv: &http.Server{
			Addr:    ":" + port,
			Handler: corsMiddleware(mux),
		},
	}
}

// Start launches the listener in a goroutine. ListenAndServe errors other
// than ErrServerClosed are fatal.
func (s *Server) Start() {
	go func() {
		log.Printf("Starting HTTP server on %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// Stop gracefully shuts the server down, honoring ctx deadlines.
func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
