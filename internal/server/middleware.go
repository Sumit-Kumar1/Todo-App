package server

import (
	"context"
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type ContextKey string

const (
	CtxKey           ContextKey = "user_id"
	invalidCookieMsg string     = "user not logged in, please login again!!"
)

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
			var temp = template.Must(template.ParseGlob("views/*"))

			cookie, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					_ = temp.ExecuteTemplate(w, "errorPage", map[string]any{
						"Code":    http.StatusUnauthorized,
						"Message": invalidCookieMsg,
					})

					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var (
				uid   uuid.UUID
				token uuid.UUID
			)

			row := db.QueryRowContext(r.Context(), "SELECT user_id, token FROM sessions WHERE token=?", cookie.Value)
			if err := row.Scan(&uid, &token); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, invalidCookieMsg, http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f(w, r.WithContext(context.WithValue(r.Context(), CtxKey, uid)))
		}
	}
}
