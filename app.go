package main

import (
	"context"
)

// App is the Wails-bound application struct. The frontend talks to the
// backend over HTTP, so this only exists to give Wails a startup hook.
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

// startup is called by Wails when the window is ready. We keep the context
// so future Wails runtime calls (e.g. dialogs, clipboard) would have access
// to the right lifecycle.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
