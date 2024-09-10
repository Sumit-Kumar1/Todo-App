package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"todoapp/internal/server"

	todoservice "todoapp/internal/service/todo"
	userservice "todoapp/internal/service/user"
	todostore "todoapp/internal/store/todo"
	userstore "todoapp/internal/store/user"
)

func main() {
	app := server.ServerFromEnvs()

	defer app.Context.Cleanup()

	us := userstore.New()
	ts := todostore.New()

	usvc := userservice.New(us)
	tsvc := todoservice.New(ts)

	app.Router.HandleFunc("GET /", h.Root)
	app.Router.
		http.HandleFunc("/", server.Chain(h.Root, server.Method(http.MethodGet)))
	http.HandleFunc("/task", server.Chain(h.TaskPage, server.Method(http.MethodGet), server.AuthMiddleware(app.Context.DB)))

	// User API
	http.HandleFunc("/register", server.Chain(h.Register, server.Method(http.MethodPost)))
	http.HandleFunc("/login", server.Chain(h.Login, server.Method(http.MethodPost)))
	http.HandleFunc("/logout", server.Chain(h.Logout, server.Method(http.MethodPost)))

	// tasks API
	http.HandleFunc("/tasks", server.Chain(h.HandleTasks, server.IsHTMX(), server.AuthMiddleware(st.DB)))
	http.HandleFunc("/tasks/{id}", server.Chain(h.Update, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(st.DB)))
	http.HandleFunc("/tasks/{id}/delete", server.Chain(h.DeleteTask, server.IsHTMX(), server.AuthMiddleware(st.DB),
		server.Method(http.MethodDelete)))
	http.HandleFunc("/tasks/{id}/done", server.Chain(h.Done, server.IsHTMX(), server.Method(http.MethodPut),
		server.AuthMiddleware(st.DB)))

	http.HandleFunc("/health", healthStatus)

	slog.Info("application is running on", "host:port", app.Configs.)

	if err := http.ListenAndServe(); err != nil {

	}
}

func healthStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json.RawMessage(`{"status":"OK"}`))
}

func runMigration(ctx server.Context) error {
	const (
		upTaskTable = `DROP TABLE IF EXISTS tasks; CREATE TABLE IF NOT EXISTS tasks(task_id TEXT PRIMARY KEY, user_id TEXT NOT NULL,
task_title TEXT NOT NULL, done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
added_at DATETIME NOT NULL, modified_at DATETIME);`
		upUserTable = `DROP TABLE IF EXISTS users; CREATE TABLE IF NOT EXISTS users(user_id TEXT NOT NULL PRIMARY KEY, name TEXT NOT NULL,
email TEXT NOT NULL UNIQUE CHECK (email LIKE '%'), password TEXT NOT NULL);`
		upSessionTable = `DROP TABLE IF EXISTS sessions; CREATE TABLE IF NOT EXISTS sessions(id TEXT PRIMARY KEY, user_id TEXT NOT NULL UNIQUE,
token TEXT NOT NULL UNIQUE, expiry DATETIME NOT NULL);`
	)

	if _, err := ctx.DB.Exec(upTaskTable); err != nil {
		return err
	}

	if _, err := ctx.DB.Exec(upUserTable); err != nil {
		return err
	}

	if _, err := ctx.DB.Exec(upSessionTable); err != nil {
		return err
	}

	return nil
}
