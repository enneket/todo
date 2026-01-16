package server

import (
	"context"
	"log"
	"net/http"
)

var srv *http.Server

func StartServer(port string) {
	mux := http.NewServeMux()

	// Register handlers
	// Note: We need to wrap handler registration with middleware or just wrap the whole mux
	mux.HandleFunc("GET /api/todos", GetTodosHandler)
	mux.HandleFunc("POST /api/todos", CreateTodoHandler)
	mux.HandleFunc("PUT /api/todos/{id}", UpdateTodoHandler)
	mux.HandleFunc("DELETE /api/todos/{id}", DeleteTodoHandler)

	// Apply CORS
	handler := corsMiddleware(mux)

	srv = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func StopServer(ctx context.Context) error {
	if srv != nil {
		return srv.Shutdown(ctx)
	}
	return nil
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
