package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

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

	newHTTPHandler(ctx, app)

	if err = migrations.RunMigrations(app, getEnvOrDefault("MIGRATION_METHOD", "UP")); err != nil {
		app.Logger.LogAttrs(ctx, slog.LevelError, err.Error())
		return
	}

	srvErr := make(chan error, 1)
	go func() {
		app.Logger.LogAttrs(ctx, slog.LevelInfo, "application is running on", slog.String("Address", app.Addr))
		srvErr <- app.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		app.Logger.LogAttrs(ctx, slog.LevelError, err.Error())
		return
	case <-ctx.Done():
		stop()
	}

	err = app.Shutdown(context.Background())
	app.Logger.LogAttrs(ctx, slog.LevelError, "error while shutting down application", slog.String("error", err.Error()))
}

func newHTTPHandler(ctx context.Context, app *server.Server) {
	usrStore := userstore.New(app.DB)
	userSvc := usersvc.New(usrStore)
	usrHTTP := userhttp.New(userSvc)

	todoStore := todostore.New(app.DB)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	public := http.FileServer(http.Dir("public"))
	openapi := http.FileServer(http.Dir("openapi"))

	http.Handle("/public/", http.StripPrefix("/public/", public))
	http.Handle("/openapi/", http.StripPrefix("/openapi/", openapi))

	http.HandleFunc("/", server.Chain(todoHTTP.Root, server.Method(http.MethodGet)))
	http.Handle("/api", http.StripPrefix("/api", server.Chain(todoHTTP.Swagger, server.Method(http.MethodGet))))
	http.HandleFunc("/task", server.Chain(todoHTTP.TaskPage, server.Method(http.MethodGet), server.AuthMiddleware(ctx, app.DB)))

	// User API
	http.HandleFunc("/register", server.Chain(usrHTTP.Register, server.Method(http.MethodPost)))
	http.HandleFunc("/login", server.Chain(usrHTTP.Login, server.Method(http.MethodPost)))
	http.HandleFunc("/logout", server.Chain(usrHTTP.Logout, server.Method(http.MethodPost)))

	// tasks API
	http.HandleFunc("/tasks", server.Chain(todoHTTP.HandleTasks, server.IsHTMX(), server.AuthMiddleware(ctx, app.DB)))
	http.HandleFunc("/tasks/{id}", server.Chain(todoHTTP.Update, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(ctx, app.DB)))
	http.HandleFunc("/tasks/{id}/delete", server.Chain(todoHTTP.DeleteTask, server.IsHTMX(), server.AuthMiddleware(ctx, app.DB),
		server.Method(http.MethodDelete)))
	http.HandleFunc("/tasks/{id}/done", server.Chain(todoHTTP.Done, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(ctx, app.DB)))

	http.HandleFunc("/health", server.Chain(func(w http.ResponseWriter, r *http.Request) {
		if err := app.DB.Ping(); err != nil {
			app.Health = &server.Health{
				Status:   "Down",
				DBStatus: "Down",
			}

			data, mErr := json.Marshal(app.Health)
			if mErr != nil {
				http.Error(w, "not able to marshal the health status", http.StatusInternalServerError)
				return
			}

			_, _ = w.Write(data)
			return
		}

		app.Health = &server.Health{
			Status:   "Up",
			DBStatus: "Up",
		}

		data, mErr := json.Marshal(app.Health)
		if mErr != nil {
			http.Error(w, "not able to marshal the health status", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(data)
	}, server.Method(http.MethodGet)))
}

func getEnvOrDefault(key, def string) string {
	eval := os.Getenv(key)
	if eval == "" {
		return def
	}

	return eval
}
