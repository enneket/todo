package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/backend/db"
	"todo/backend/server"
)

func main() {
	log.Println("Starting headless server...")

	// Initialize DB
	if err := db.InitDB("todo.db"); err != nil {
		log.Fatal(err)
	}

	// Start HTTP Server
	srv := server.NewServer("8081")
	srv.Start()

	log.Println("Server started on :8081")

	// Wait for interrupt signal, then shut down cleanly.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
