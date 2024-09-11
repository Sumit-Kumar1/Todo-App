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
	app := server.ServerFromEnvs()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	st, fn, err := store.New(logger)
	if err != nil {
		return
	}

	defer fn()

	svc := service.New(st, logger)
	h := handler.New(svc, logger)

	public := http.FileServer(http.Dir("public"))
	swagger := http.FileServer(http.Dir("openapi"))

	http.Handle("/public/", http.StripPrefix("/public/", public))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", swagger))

	http.HandleFunc("/", server.Chain(h.Root, server.Method(http.MethodGet)))
	http.HandleFunc("/task", server.Chain(h.TaskPage, server.Method(http.MethodGet), server.AuthMiddleware(st.DB)))

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

	http.HandleFunc("/health", server.Chain(healthStatus, server.Method(http.MethodGet)))

	slog.Info("application is running on", "host:port", app.Addr)

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
