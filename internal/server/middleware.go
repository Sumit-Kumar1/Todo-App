package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/sqlitecloud/sqlitecloud-go"
)

const (
	invalidCookieMsg = "user not logged in, please login again!!"
	cookieName       = "token"
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
				http.Error(
					w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)

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

func AuthMiddleware(ctx context.Context, db *sqlitecloud.SQCloud) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var (
				temp   = template.Must(template.ParseGlob("views/*"))
				uid    uuid.UUID
				logger = models.GetLoggerFromCtx(ctx)
			)

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					_ = temp.ExecuteTemplate(w, "errorPage", map[string]any{
						"Code":    http.StatusUnauthorized,
						"Message": invalidCookieMsg,
					})

					logger.LogAttrs(
						ctx,
						slog.LevelError,
						err.Error(),
						slog.String("problem", "no cokkie found, please login again!"),
					)

					return
				}

				logger.LogAttrs(ctx, slog.LevelError, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			if _, err = uuid.Parse(cookie.Value); err != nil {
				logger.LogAttrs(ctx, slog.LevelError, "invalid cokkie found, please login again")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			row, err := db.Select(
				fmt.Sprintf("SELECT user_id FROM sessions WHERE token='%s';", cookie.Value),
			)
			if err != nil {
				logger.LogAttrs(ctx, slog.LevelError, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			if row.GetNumberOfRows() == uint64(0) { // this means no rows
				logger.LogAttrs(ctx, slog.LevelError, "no valid session found, login again")
				http.Error(w, invalidCookieMsg, http.StatusUnauthorized)

				return
			}

			for r := uint64(0); r < row.GetNumberOfRows(); r++ {
				userID, err := row.GetStringValue(r, 0)
				if err != nil {
					logger.LogAttrs(ctx, slog.LevelError, err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

				uid, err = uuid.Parse(userID)
				if err != nil {
					logger.LogAttrs(ctx, slog.LevelError, err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}
			}

			f(w, r.WithContext(context.WithValue(ctx, models.CtxKeyUserID, uid)))
		}
	}
}
