package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"todoapp/internal/handler"
	"todoapp/internal/server"
	"todoapp/internal/service"
	"todoapp/internal/store"

	"github.com/google/uuid"
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

	// User API
	http.HandleFunc("/register", h.Register)
	http.HandleFunc("/login", h.Login)

	// tasks API
	http.HandleFunc("/tasks", chain(h.HandleTasks, isHTMX(), authMiddleware(st.DB)))
	http.HandleFunc("/task/{id}", chain(h.HandleIDReq, isHTMX(), authMiddleware(st.DB)))
	http.HandleFunc("/task/done/{id}", chain(h.Done, isHTMX(), method(http.MethodPost), authMiddleware(st.DB)))

	http.HandleFunc("/health", healthStatus)

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

type middleware func(http.HandlerFunc) http.HandlerFunc

func chain(f http.HandlerFunc, middlewares ...middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

func method(m string) middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

			f(w, r)
		}
	}
}

func isHTMX() middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Hx-Request") != "true" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			f(w, r)
		}
	}
}

func authMiddleware(db *sql.DB) middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("user_session")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					http.Error(w, "invalid cookie", http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var (
				uid   uuid.UUID
				token uuid.UUID
			)

			row := db.QueryRowContext(r.Context(), "select user_id, token from sessions where token = ?", cookie.Value)
			err = row.Scan(&uid, &token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "invalid cookie", http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f(w, r.WithContext(context.WithValue(r.Context(), "user_id", uid)))
		}
	}
}
