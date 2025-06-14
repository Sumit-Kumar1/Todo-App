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
				logger = models.GetLoggerFromCtx(ctx)
			)

			cookieVal, err := validateCookie(ctx, logger, r)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					_ = temp.ExecuteTemplate(w, "errorPage", map[string]any{
						"Code":    http.StatusUnauthorized,
						"Message": invalidCookieMsg,
					})

					return
				}

				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			uid, err := getSessionID(ctx, db, logger, cookieVal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}

			f(w, r.WithContext(context.WithValue(ctx, models.CtxKeyUserID, *uid)))
		}
	}
}

func validateCookie(ctx context.Context, logger *slog.Logger, r *http.Request) (*uuid.UUID, error) {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		uid, err := uuid.Parse(cookie.Value)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "invalid cookie found, please login again")

			return nil, models.ErrInvalidCookie
		}

		return &uid, nil
	}

	if errors.Is(err, http.ErrNoCookie) {
		logger.LogAttrs(ctx, slog.LevelError, err.Error(),
			slog.String("error", "no cookie found, please login again!"),
		)

		return nil, err
	}

	logger.LogAttrs(ctx, slog.LevelError, err.Error())

	return nil, err
}

func getSessionID(ctx context.Context, db *sqlitecloud.SQCloud, logger *slog.Logger, sessionToken *uuid.UUID) (*uuid.UUID, error) {
	var (
		uid uuid.UUID
		err error
	)

	row, err := db.Select(
		fmt.Sprintf("SELECT user_id FROM sessions WHERE token='%s';", *sessionToken),
	)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, err.Error())

		return nil, err
	}

	if row.GetNumberOfRows() == uint64(0) {
		logger.LogAttrs(ctx, slog.LevelError, "no valid session found, login again")

		return nil, models.ErrInvalidCookie
	}

	for r := uint64(0); r < row.GetNumberOfRows(); r++ {
		userID, err := row.GetStringValue(r, 0)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, err.Error())

			return nil, err
		}

		uid, err = uuid.Parse(userID)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, err.Error())

			return nil, err
		}
	}

	return &uid, nil
}
