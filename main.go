package main

import (
	"context"
	"embed"
	"log"
	"todo/backend/db"
	"todo/backend/server"
	"todo/backend/service"

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
	srv := server.NewServer("8081")
	srv.Start()

	// Start Notification Scheduler
	notifCtx, stopNotif := context.WithCancel(context.Background())
	service.StartNotificationScheduler(notifCtx)

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options.
	// Note: wails.Run blocks until the window is closed, so `defer srv.Stop`
	// here would never fire. Use the OnShutdown hook instead so the HTTP server
	// actually shuts down when the app exits.
	err := wails.Run(&options.App{
		Title:  "Todo App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: app.startup,
		OnShutdown: func(ctx context.Context) {
			stopNotif()
			if err := srv.Stop(ctx); err != nil {
				log.Printf("error shutting down HTTP server: %v", err)
			}
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
