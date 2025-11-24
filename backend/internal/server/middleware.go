package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"todoapp/internal/errors"
	"todoapp/internal/models"
	"todoapp/internal/todocookie"

	"github.com/google/uuid"
)

const (
	cookieName     = "auth"
	allowedOrigin  = "http://localhost:4173"
	allowedMethods = "POST, GET, PUT, DELETE, PATCH, OPTIONS"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding"
	corsMaxAge     = "8640"
)

type Claims struct {
	Email    string `json:"email"`
	ClaimUID string `json:"claimID"`
	// registeredClaim's subject is userID from db
	jwt.RegisteredClaims
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

// AuthMiddleware check for valid cookie and extracts user id for the user
// TODO: Decide wether if token is expired had to refresh and add a new auth cookie to disrupt the flow ?
func (s *Server) AuthMiddleware() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			uid, err := s.validateCookie(ctx, s.Logger, r)
			if err != nil {
				errors.HandleHTTPError(w, err)
				s.Logger.LogAttrs(ctx, slog.LevelError, "error while validating cookie",
					slog.String("error", err.Error()))

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
				slog.String(models.HeaderCorrelation, corrID),
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

// ServerWideMiddlewares CORS Part middleware
func (s *Server) ServerWideMiddlewares(next http.Handler) http.Handler {
	origin := GetEnvOrDefault("CORS_ORIGIN", allowedOrigin)
	methods := GetEnvOrDefault("CORS_METHODS", allowedMethods)
	headers := GetEnvOrDefault("CORS_HEADERS", allowedHeaders)
	maxAge := GetEnvOrDefault("CORS_MAX_AGE", corsMaxAge)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Max-Age", maxAge)

			w.WriteHeader(http.StatusOK)

			slog.Log(r.Context(), slog.LevelDebug, "responded cors req with OK")
			return
		}

		s.globalRateLimiter(next)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) globalRateLimiter(next http.Handler) http.Handler {
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

func (s *Server) validateCookie(ctx context.Context, logger *slog.Logger, r *http.Request) (*uuid.UUID, error) {
	val, err := todocookie.ReadCookie(r, cookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(val, &Claims{}, func(token *jwt.Token) (any, error) {
		secVal := GetEnvOrDefault("AUTH_SECRET", "33cea8f88c5c8ad73b1700af7d72891fe3097297e59fb6cbe5fd8b545a8316d0")
		return []byte(secVal), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.ErrInvalid("claim")
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return extractUserID(claims.Subject)
}

func extractUserID(claimSubject string) (*uuid.UUID, error) {
	uid, err := uuid.Parse(claimSubject)
	if err != nil {
		return nil, err
	}

	if uid == uuid.Nil {
		return nil, errors.ErrInvalid("nil user id in claims")
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
