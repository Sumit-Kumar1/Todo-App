package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

func Method(m string) Middleware {
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

func IsHTMX() Middleware {
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

func AuthMiddleware(db *sql.DB) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
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
