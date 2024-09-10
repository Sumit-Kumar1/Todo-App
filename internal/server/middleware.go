package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Use(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

func CustomHandler() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func IsHTMX() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Hx-Request") != "true" {
				http.Error(w, "not an htmx request", http.StatusBadRequest)
				return
			}

			f(w, r)
		}
	}
}

func AuthMiddleware(ctx Context) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var (
				uid   uuid.UUID
				token uuid.UUID
				key   contextKey = "user_id"
			)

			cookie, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					http.Error(w, "invalid cookie", http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			row := ctx.DB.QueryRowContext(r.Context(), "select user_id, token from sessions where token = ?", cookie.Value)
			err = row.Scan(&uid, &token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "invalid cookie", http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f(w, r.WithContext(context.WithValue(r.Context(), key, uid)))
		}
	}
}
