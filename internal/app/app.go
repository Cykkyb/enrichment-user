package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	HttpServer *http.Server
	log        *slog.Logger
}

func NewApp(log *slog.Logger) *App {
	return &App{
		HttpServer: nil,
		log:        log,
	}
}

func (a *App) Run(handler http.Handler, port string) error {
	a.log.Info("Starting server")

	a.HttpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	err := a.HttpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to start server", err)
	}

	return nil
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.HttpServer.Shutdown(ctx); err != nil {
		a.log.Error("Graceful shutdown failed:", err)
		return err
	}
	a.log.Info("Server gracefully stopped")

	return nil
}
