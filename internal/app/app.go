package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"neiro/internal/app/handler"
	"neiro/internal/app/health"
	"neiro/internal/app/services"
	"neiro/internal/app/web"
	"neiro/internal/lib/middleware"

	"github.com/gorilla/mux"
)

type App struct {
	HTTPServer *web.HTTPServer
	service    services.IService

	health.LivenessChecker
	health.ReadinessChecker
}

// NewMyService creates new instance of MyService.
func NewMyService(
	log *slog.Logger,
	port int,
	useTracing bool,
	tracingAddress string,
	ttl int,
	clearInterval time.Duration,
	serviceName string,
) (*App, error) {
	const op = "app.NewMyService"
	ctx := context.Background()

	app := &App{}
	srv, err := services.NewService(log, ttl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	telemetryMiddleware, err := addTelemetryMiddleware(ctx, useTracing, tracingAddress, serviceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	router := mux.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(telemetryMiddleware)

	router.HandleFunc("/live", health.LivenessHandler(app)).Methods("GET")
	router.HandleFunc("/ready", health.ReadinessHandler(app)).Methods("GET")

	router.HandleFunc("/api/set", handler.PostItem(srv)).Methods("POST")
	router.HandleFunc("/api/get/{id}", handler.GetItem(srv)).Methods("GET")
	router.HandleFunc("/api/delete/{id}", handler.DeleteItem(srv)).Methods("DELETE")
	server, err := web.New(log, port, router)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Run cache cleaner.
	go srv.ClearCache(clearInterval)

	app.HTTPServer = server
	app.service = srv

	return app, nil
}

// Start runs the application.
func (a *App) Start() {
	a.HTTPServer.Start()
}

// Stop stops the application.
func (a *App) Stop() {
	if a != nil && a.service != nil {
		err := a.service.Close()
		if err != nil {
			fmt.Println("An error occurred closing service" + err.Error())

			return
		}
	}
}

// addTelemetryMiddleware adds telemetry middleware to the router.
func addTelemetryMiddleware(ctx context.Context, useTracing bool, tracingAddress string, serviceName string) (mux.MiddlewareFunc, error) {
	var telemetryMiddleware mux.MiddlewareFunc
	var err error
	if useTracing {
		telemetryMiddleware, err = handler.AddTelemetryMiddleware(ctx, tracingAddress, serviceName)
		if err != nil {
			return nil, err
		}
	}

	return telemetryMiddleware, nil
}

// LivenessCheck checks if the application is live.
func (a *App) LivenessCheck() bool {
	return a.service.LivenessCheck()
}

// ReadinessCheck checks if the application is ready for work.
func (a *App) ReadinessCheck() bool {
	return a.service.ReadinessCheck()
}
