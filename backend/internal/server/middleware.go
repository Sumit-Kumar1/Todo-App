package server

import (
	"context"
	"database/sql"
	libErr "errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	cookieName = "token"
)

var errMethodNotAllowed = libErr.New(http.StatusText(http.StatusMethodNotAllowed))

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

// TODO: do we still need this ??
func MethodWithCORS(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				errors.HandleHTTPError(w, errMethodNotAllowed)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "todo.zone.id")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

			f(w, r)
		}
	}
}

// TODO: fix auth Middleware
func (s *Server) AuthMiddleware(ctx context.Context) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookieVal, err := validateCookie(ctx, s.Logger, r)
			if err != nil {
				errors.HandleHTTPError(w, err)

				s.Logger.LogAttrs(ctx, slog.LevelError, "error while validating cookie",
					slog.String("error", err.Error()))

				w.Header().Set("HX-Redirect", "/")
				return
			}

			uid, err := getSessionID(ctx, s.DB, s.Logger, cookieVal)
			if err != nil {
				s.Logger.LogAttrs(ctx, slog.LevelError, "error while validating session",
					slog.String("error", err.Error()))

				http.Error(w, err.Error(), http.StatusUnauthorized)

				w.Header().Set("HX-Redirect", "/")
				return
			}

			f(w, r.WithContext(context.WithValue(ctx, models.CtxKeyUserID, *uid)))
		}
	}
}

func AddCorrelation() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			corrID := r.Header.Get(models.HeaderCorrelation)
			if strings.TrimSpace(corrID) == "" {
				corrID = uuid.NewString()
			}

			logger := slog.With(slog.Group("request",
				slog.String(string(models.HeaderCorrelation), corrID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("host", r.Host),
				slog.String("remote-addr", r.RemoteAddr),
			))

			ctx := context.WithValue(r.Context(), models.Logger, logger)

			f(w, r.WithContext(context.WithValue(ctx, models.CorrelationID, corrID)))
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
func (s *Server) rateLimiterLogin() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			email := r.PostFormValue("email")
			if strings.TrimSpace(email) == "" {
				errors.HandleHTTPError(w, errors.ErrInvalid("email"))
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
			s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "login attempt count increased", slog.Int("count", attempt.count))

			if attempt.count > s.loginLimiter.maxAttempts {
				s.Logger.LogAttrs(r.Context(), slog.LevelDebug, "login attempt count exceeded",
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

			return nil, errors.ErrInvalidCookie
		}

		return &uid, nil
	}

	if libErr.Is(err, http.ErrNoCookie) {
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

	query := "SELECT user_id FROM sessions WHERE token=$1;"

	row := db.QueryRowContext(ctx, query, *sessionToken)
	if err := row.Scan(&userID); err != nil {
		if libErr.Is(err, sql.ErrNoRows) {
			logger.LogAttrs(ctx, slog.LevelError, "no valid session found, please login again")

			return nil, errors.ErrInvalidCookie
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
