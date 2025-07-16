package server

import (
	"context"
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	invalidCookieMsg = "user not logged in, please login again!!"
	cookieName       = "token"
)

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
				http.Error(w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)

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

func (s *Server) authMiddleware(ctx context.Context) middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var temp = template.Must(template.ParseGlob("views/*"))

			cookieVal, err := validateCookie(ctx, s.Logger, r)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					_ = temp.ExecuteTemplate(w, "errorPage", map[string]any{
						"Code":    http.StatusUnauthorized,
						"Message": invalidCookieMsg,
					})

					return
				}

				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				_ = temp.ExecuteTemplate(w, "errorPage", map[string]any{
					"Code":    http.StatusUnauthorized,
					"Message": invalidCookieMsg,
				})

				return
			}

			uid, err := getSessionID(ctx, s.DB, s.Logger, cookieVal)
			if err != nil {
				s.Logger.LogAttrs(ctx, slog.LevelError, "error while validating session", slog.String("error", err.Error()))

				_ = temp.ExecuteTemplate(w, "errorPage", map[string]any{
					"Code":    http.StatusUnauthorized,
					"Message": invalidCookieMsg,
				})

				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}

			f(w, r.WithContext(context.WithValue(ctx, models.CtxKeyUserID, *uid)))
		}
	}
}

func (s *Server) GlobalRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		now := time.Now()

		s.globalLimiter.mu.Lock()

		info, exists := s.globalLimiter.attempts[ip]
		if !exists {
			info = &limiterAttempt{
				count:     1,
				firstTime: now,
			}

			s.globalLimiter.attempts[ip] = info
			s.globalLimiter.mu.Unlock()
			next.ServeHTTP(w, r)

			return
		}

		if now.Sub(info.firstTime) > s.globalLimiter.timeWindow {
			info.firstTime = now
			info.count = 1
		} else {
			info.count++
		}

		if info.count > s.globalLimiter.maxAttempts {
			w.Header().Set("Retry-After", strconv.Itoa(int(s.globalLimiter.timeWindow.Seconds())))
			http.Error(w, "Too many requests. Please try again later!!", http.StatusTooManyRequests)
			s.globalLimiter.mu.Unlock()

			return
		}

		s.globalLimiter.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// nolint:gocognit // can't divide it furthur
func (s *Server) rateLimiterLogin() middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "started login rate limiter")

			email := r.FormValue("email")
			if strings.TrimSpace(email) == "" {
				http.Error(w, "invalid email provided", http.StatusBadRequest)
				s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "invalid email in rate limiter login")

				return
			}

			s.loginLimiter.mu.Lock()

			attempt, exists := s.loginLimiter.attempts[email]
			if !exists {
				attempt = &limiterAttempt{count: 0, firstTime: time.Now()}
				s.loginLimiter.attempts[email] = attempt
			}

			if time.Since(attempt.firstTime) > s.loginLimiter.timeWindow {
				attempt.count = 0
				attempt.firstTime = time.Now()

				f(w, r)

				return
			}

			attempt.count++
			s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "attempt count increased", slog.Int("count", attempt.count))

			if attempt.count > s.loginLimiter.maxAttempts {
				s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "attempt count exceeded",
					slog.Int("count", attempt.count), slog.Int("max attempt", s.loginLimiter.maxAttempts))

				http.Error(w, "Too many login attempts. Please try again later.", http.StatusTooManyRequests)

				s.loginLimiter.mu.Unlock()

				return
			}

			s.loginLimiter.mu.Unlock()
			s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "success login limiter finished")
			f(w, r)
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

func getSessionID(ctx context.Context, db *sql.DB, logger *slog.Logger, sessionToken *uuid.UUID) (*uuid.UUID, error) {
	var (
		userID string
		uid    uuid.UUID
		err    error
	)

	query := "SELECT user_id FROM sessions WHERE token=?;"

	row := db.QueryRowContext(ctx, query, *sessionToken)
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.LogAttrs(ctx, slog.LevelError, "no valid session found, please login again")

			return nil, models.ErrInvalidCookie
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error())

		return nil, err
	}

	uid, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return &uid, nil
}

// Extract client IP address from request (trusting RemoteAddr, no proxy handling)
func clientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback to whole string
	}

	return ip
}
