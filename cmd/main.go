package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"todoapp/internal/migrations"
	"todoapp/internal/server"
	"todoapp/internal/service/todosvc"
	"todoapp/internal/store/todostore"

	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	usersvc "todoapp/internal/service/user"
	userstore "todoapp/internal/store/user"
)

func main() {
	app, err := server.ServerFromEnvs()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err = migrations.RunMigrations(app, getEnvOrDefault("MIGRATION_METHOD", "UP")); err != nil {
		slog.Error(err.Error())
		return
	}

	usrSt := userstore.New(app.DB, app.Logger)
	userSvc := usersvc.New(usrSt, app.Logger)
	usrHTTP := userhttp.New(userSvc, app.Logger)

	todoStore := todostore.New(app.DB, app.Logger)
	todoSvc := todosvc.New(todoStore, app.Logger)
	todoHTTP := todohttp.New(todoSvc, app.Logger)

	public := http.FileServer(http.Dir("public"))
	swagger := http.FileServer(http.Dir("openapi"))

	http.Handle("/public/", http.StripPrefix("/public/", public))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", swagger))

	http.HandleFunc("/", server.Chain(todoHTTP.Root, server.Method(http.MethodGet)))
	http.HandleFunc("/task", server.Chain(todoHTTP.TaskPage, server.Method(http.MethodGet), server.AuthMiddleware(app.DB)))

	// User API
	http.HandleFunc("/register", server.Chain(usrHTTP.Register, server.Method(http.MethodPost)))
	http.HandleFunc("/login", server.Chain(usrHTTP.Login, server.Method(http.MethodPost)))
	http.HandleFunc("/logout", server.Chain(usrHTTP.Logout, server.Method(http.MethodPost)))

	// tasks API
	http.HandleFunc("/tasks", server.Chain(todoHTTP.HandleTasks, server.IsHTMX(), server.AuthMiddleware(app.DB)))
	http.HandleFunc("/tasks/{id}", server.Chain(todoHTTP.Update, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(app.DB)))
	http.HandleFunc("/tasks/{id}/delete", server.Chain(todoHTTP.DeleteTask, server.IsHTMX(), server.AuthMiddleware(app.DB),
		server.Method(http.MethodDelete)))
	http.HandleFunc("/tasks/{id}/done", server.Chain(todoHTTP.Done, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(app.DB)))

	http.HandleFunc("/health", server.Chain(healthStatus, server.Method(http.MethodGet)))

	app.Logger.Info("application is running on", "Address", app.Addr)

	err = app.ListenAndServe()
	if err != nil {
		slog.Error("error while running server", "error", err)

		return
	}
}

func healthStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json.RawMessage(`{"status":"OK"}`))
}

func getEnvOrDefault(key, def string) string {
	eval := os.Getenv(key)
	if eval == "" {
		return def
	}

	return eval
}
