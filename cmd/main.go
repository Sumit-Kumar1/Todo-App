package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"todoapp/internal/handler"
	"todoapp/internal/migrations"
	"todoapp/internal/models"
	"todoapp/internal/server"
	"todoapp/internal/service/todosvc"
	"todoapp/internal/store/todostore"

	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	usersvc "todoapp/internal/service/user"
	userstore "todoapp/internal/store/user"
)

func main() {
	// handling SIGINT gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app, err := server.ServerFromEnvs()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// add logger into context
	ctx = context.WithValue(ctx, models.Logger, app.Logger)

	setupRoutes(ctx, app)

	if err = migrations.RunMigrations(ctx, app, getEnvOrDefault("MIGRATION_METHOD", "UP")); err != nil {
		slog.Error(err.Error())
		return
	}

	srvErr := make(chan error, 1)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(app.Configs.Host, app.Configs.Port),
		Handler:      app.Mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		app.Logger.LogAttrs(ctx, slog.LevelInfo, "Server started", slog.String("Address", httpServer.Addr))
		srvErr <- httpServer.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			app.Logger.LogAttrs(ctx, slog.LevelError, "error listening and serving", slog.String("error", err.Error()))
		}

		return
	case <-ctx.Done():
		stop()
	}

	err = httpServer.Shutdown(context.Background())
	app.Logger.LogAttrs(ctx, slog.LevelError, "error while shutting down the server", slog.String("error", err.Error()))
}

func setupRoutes(ctx context.Context, app *server.Server) {
	setupPublicRoutes(app)
	setupUserRoutes(app)
	setupTasksRoutes(ctx, app)
}

func setupTasksRoutes(ctx context.Context, app *server.Server) {
	todoStore := todostore.New(app.DB)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	// tasks API
	app.Mux.HandleFunc("/task", server.Chain(todoHTTP.TaskPage, server.Method(http.MethodGet), server.AuthMiddleware(ctx, app.DB)))
	app.Mux.HandleFunc("/tasks", server.Chain(todoHTTP.HandleTasks, server.IsHTMX(), server.AuthMiddleware(ctx, app.DB)))
	app.Mux.HandleFunc("/tasks/{id}", server.Chain(todoHTTP.Update, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(ctx, app.DB)))
	app.Mux.HandleFunc("/tasks/{id}/delete", server.Chain(todoHTTP.DeleteTask, server.IsHTMX(), server.AuthMiddleware(ctx, app.DB),
		server.Method(http.MethodDelete)))
	app.Mux.HandleFunc("/tasks/{id}/done", server.Chain(todoHTTP.Done, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(context.Background(), app.DB)))
}

func setupUserRoutes(app *server.Server) {
	usrSt := userstore.New(app.DB)
	userSvc := usersvc.New(usrSt)
	usrHTTP := userhttp.New(userSvc)

	// User API
	app.Mux.HandleFunc("/register", server.Chain(usrHTTP.Register, server.Method(http.MethodPost)))
	app.Mux.HandleFunc("/login", server.Chain(usrHTTP.Login, server.Method(http.MethodPost)))
	app.Mux.HandleFunc("/logout", server.Chain(usrHTTP.Logout, server.Method(http.MethodPost)))
}

func setupPublicRoutes(app *server.Server) {
	h := handler.New()

	public := http.FileServer(http.Dir("public"))
	openapi := http.FileServer(http.Dir("openapi"))

	app.Mux.HandleFunc("/", server.Chain(h.Root, server.Method(http.MethodGet)))
	app.Mux.Handle("/public/", http.StripPrefix("/public/", public))
	app.Mux.Handle("/openapi/", http.StripPrefix("/openapi/", openapi))
	app.Mux.Handle("/api", http.StripPrefix("/api", server.Chain(h.Swagger, server.Method(http.MethodGet))))
}

func getEnvOrDefault(key, def string) string {
	eval := os.Getenv(key)
	if eval == "" {
		return def
	}

	return eval
}
