package main

import (
	"context"
	"embed"
	"ewails/app/server"
	"ewails/app/server/config"
	"ewails/app/server/services"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	config.ReadEnv()
	app := NewApp()
	es := services.EmailService{}

	serverMiddleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "emailService", &es)
		*r = *r.WithContext(ctx)
	}

	r := server.Create(serverMiddleware)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ewails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Middleware: func(next http.Handler) http.Handler {
				r.UseHandler("/", next)
				return r.Mux
			},
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
