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

	st, err := store.New(logger)
	if err != nil {
		slog.Error("DB Creation err", "error", err)
		return
	}

	if err = st.DB.Ping(); err != nil {
		slog.Error("database not reachable", "error", err)
		return
	}

	defer func() {
		if err := st.DB.Close(); err != nil {
			logger.Error("DB Close", "error", err)
		}
	}()

	slog.Info("db connection", "success", st.DB.Stats())

	svc := service.New(st, logger)
	h := handler.New(svc, logger)

	http.HandleFunc("/", h.IndexPage)
	http.HandleFunc("/add", h.AddTask)
	http.HandleFunc("/delete/{id}", h.DeleteTask)
	http.HandleFunc("/update/{id}", h.Update)
	http.HandleFunc("/done/{id}", h.Done)
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

func healthStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(json.RawMessage(`{"status":"OK"}`))
}
