package main

import (
	"context"
	"embed"
	"log"
	"todo/backend/db"
	"todo/backend/server"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize DB
	if err := db.InitDB("todo.db"); err != nil {
		log.Fatal(err)
	}

	// Start HTTP Server
	// Use a fixed port for now, e.g., 8081
	server.StartServer("8081")
	defer server.StopServer(context.Background())

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Todo App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
