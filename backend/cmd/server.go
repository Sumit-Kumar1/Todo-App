// Package cmd setups go server and start the API
package cmd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"todoapp/client"
	"todoapp/internal/handler"
	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	"todoapp/internal/migrations"
	"todoapp/internal/server"
	todosvc "todoapp/internal/service/todo"
	usersvc "todoapp/internal/service/user"
	"todoapp/internal/store"
)

func Run(c context.Context, _ io.Writer, _ []string) error {
	ctx, stop := signal.NotifyContext(c, os.Interrupt)
	defer stop()

	app, err := server.NewServer()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err = migrations.RunMigrations(ctx, app, app.MigrationMethod); err != nil {
		slog.LogAttrs(c, slog.LevelError, "error while running migrations",
			slog.String("error", err.Error()))

		return err
	}

	// Setup health endpoints first
	setupHealthRoutes(app)
	setupUserRoutes(app)
	setupTasksRoutes(app)

	srvErr := make(chan error, 1)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(app.Host, app.Port),
		Handler:      app.ServerWideMiddlewares(app.Mux),
		ReadTimeout:  time.Duration(app.ReadTimeout * int(time.Second)),
		WriteTimeout: time.Duration(app.WriteTimeout * int(time.Second)),
		IdleTimeout:  time.Duration(app.IdleTimeout * int(time.Second)),
	}

	go func() {
		app.Logger.LogAttrs(ctx, slog.LevelInfo, "Server started", slog.String("Address", httpServer.Addr))
		srvErr <- httpServer.ListenAndServe()
	}()

	// Wait for shutdown signal or error
	if err := handleServerError(ctx, app, srvErr); err != nil {
		return err
	}

	// Graceful shutdown
	return shutdownServer(ctx, app, httpServer)
}

// setupHealthRoutes configures health check endpoints
func setupHealthRoutes(app *server.Server) {
	hh := handler.NewHealthHandler(app)

	healthRoutes := []struct {
		methodPath string
		handler    http.HandlerFunc
		middleware []server.Middleware
	}{
		{"GET /health", hh.HealthHandler, []server.Middleware{server.AddCorrelation()}},
		{"GET /ready", hh.ReadyHandler, []server.Middleware{server.AddCorrelation()}},
		{"GET /live", hh.LivenessHandler, []server.Middleware{server.AddCorrelation()}},
	}

	for i := range healthRoutes {
		app.Mux.HandleFunc(healthRoutes[i].methodPath,
			server.Chain(healthRoutes[i].handler, healthRoutes[i].middleware...),
		)
	}
}

// handleServerError processes server errors and shutdown signals
func handleServerError(ctx context.Context, app *server.Server, srvErr chan error) error {
	select {
	case err := <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			app.Logger.LogAttrs(ctx, slog.LevelError,
				"error listening and serving",
				slog.String("error", err.Error()),
			)
			return err
		}

		app.Logger.LogAttrs(ctx, slog.LevelInfo, "server stopped due to error")

	case <-ctx.Done():
		app.Logger.LogAttrs(ctx, slog.LevelInfo, "server shutdown triggered by signal")
	}

	return nil
}

// shutdownServer gracefully shuts down the server
func shutdownServer(ctx context.Context, app *server.Server, httpServer *http.Server) error {
	if err := httpServer.Shutdown(context.Background()); err != nil {
		app.Logger.LogAttrs(ctx, slog.LevelError, "error while shutting down the server",
			slog.String("error", err.Error()),
		)

		return err
	}

	// Call the server's shutdown function to clean up resources
	if app.ShutDownFxn != nil {
		if err := app.ShutDownFxn(ctx); err != nil {
			app.Logger.LogAttrs(ctx, slog.LevelError, "error during server shutdown",
				slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}

func setupTasksRoutes(app *server.Server) {
	todoStore := store.New(app.DB)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	taskRoutes := []struct {
		methodPath string
		handler    http.HandlerFunc
		middleware []server.Middleware
	}{
		{"POST /task", todoHTTP.AddTask, []server.Middleware{server.AddCorrelation(), app.AuthMiddleware()}},
		{"GET /tasks", todoHTTP.GetAllTasks, []server.Middleware{server.AddCorrelation(), app.AuthMiddleware()}},
		{"PATCH /tasks/{id}/done", todoHTTP.MarkDone, []server.Middleware{server.AddCorrelation(), app.AuthMiddleware()}},
		{"PUT /tasks/{id}", todoHTTP.Update, []server.Middleware{server.AddCorrelation(), app.AuthMiddleware()}},
		{"DELETE /tasks/{id}", todoHTTP.DeleteTask, []server.Middleware{server.AddCorrelation(), app.AuthMiddleware()}},
	}

	for i := range taskRoutes {
		app.Mux.HandleFunc(taskRoutes[i].methodPath,
			server.Chain(taskRoutes[i].handler, taskRoutes[i].middleware...),
		)
	}
}

func setupUserRoutes(app *server.Server) {
	authURL := server.GetEnvOrDefault("AUTH_URL", "http://localhost:9001")

	authClient := client.New(authURL)
	userSvc := usersvc.New(authClient)
	userHTTP := userhttp.New(userSvc)

	taskRoutes := []struct {
		methodPath string
		handler    http.HandlerFunc
		middleware []server.Middleware
	}{
		{"POST /register", userHTTP.Register, []server.Middleware{server.AddCorrelation()}},
		{"POST /login", userHTTP.Login, []server.Middleware{server.AddCorrelation()}},
		{"POST /logout", userHTTP.Logout, []server.Middleware{server.AddCorrelation()}},
	}

	for i := range taskRoutes {
		app.Mux.HandleFunc(taskRoutes[i].methodPath,
			server.Chain(taskRoutes[i].handler, taskRoutes[i].middleware...),
		)
	}
}
