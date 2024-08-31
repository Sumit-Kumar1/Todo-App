package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"todoapp/internal/handler"
	"todoapp/internal/server"
	"todoapp/internal/service"
	"todoapp/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	st, fn, err := store.New(logger)
	if err != nil {
		return
	}

	if err := runMigration(st); err != nil {
		logger.Error("error while running migration function", "error", err.Error())
		return
	}

	defer fn()

	svc := service.New(st, logger)
	h := handler.New(svc, logger)

	// User API

	// Todo API
	http.HandleFunc("/todo", h.AddTask)
	http.HandleFunc("/task/{id}", h.HandleIDReq)
	http.HandleFunc("/task/done/{id}", h.Done)

	http.HandleFunc("/health", healthStatus)

	app := server.NewServer(
		server.WithAppName("todoApp"),
		server.WithEnv("development"),
		server.WithPort("9001"),
	)

	slog.Info("application is running on", "host:port", app.Addr)

	err = app.ListenAndServe()
	if err != nil {
		slog.Error("error while running server", "error", err)

		return
	}
}

func runMigration(st *store.Store) error {
	const (
		createTaskTable = `CREATE TABLE IF NOT EXISTS tasks(task_id TEXT NOT NULL PRIMARY KEY,
		task_title TEXT NOT NULL, done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
		added_at DATETIME NOT NULL, modified_at DATETIME);`
		createUserTable = `CREATE TABLE IF NOT EXISTS users(user_id TEXT NOT NULL PRIMARY KEY, name TEXT NOT NULL,
email TEXT NOT NULL UNIQUE CHECK (email LIKE '%'), password TEXT NOT NULL);`
	)

	if _, err := st.DB.Exec(createTaskTable); err != nil {
		return err
	}

	if _, err := st.DB.Exec(createUserTable); err != nil {
		return err
	}

	return nil
}

func healthStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json.RawMessage(`{"status":"OK"}`))
}
