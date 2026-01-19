package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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
	server.StartServer("8081")

	log.Println("Server started on :8081")

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
}
