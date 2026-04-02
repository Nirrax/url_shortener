package api

import (
	"net/http"

	"github.com/nirrax/url_shortener/internal/config"
	"github.com/nirrax/url_shortener/internal/service"
)

type App struct {
	config.Config
	service.Service
}

func (app *App) Run() error {
	server := &http.Server{
		Addr:    ":" + app.Config.ServerPort,
		Handler: app.NewRouter(),
	}

	return server.ListenAndServe()
}
