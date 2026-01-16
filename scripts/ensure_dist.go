package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	distDir := filepath.Join("frontend", "dist")
	indexFile := filepath.Join(distDir, "index.html")

	// Create frontend/dist directory
	if err := os.MkdirAll(distDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", distDir, err)
		os.Exit(1)
	}

	// Create index.html if it doesn't exist
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		if err := os.WriteFile(indexFile, []byte("<!DOCTYPE html><html><body>Placeholder</body></html>"), 0644); err != nil {
			fmt.Printf("Error creating file %s: %v\n", indexFile, err)
			os.Exit(1)
		}
		fmt.Printf("Created placeholder %s\n", indexFile)
	} else {
		fmt.Printf("File %s already exists\n", indexFile)
	}
}
