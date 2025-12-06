package cmd

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	cookieName     = "auth"
	allowedOrigins = "http://localhost:8081"
	allowedMethods = "POST,GET,PUT,DELETE,PATCH,OPTIONS"
	allowedHeaders = "Accept,Content-Type,Content-Length,Accept-Encoding"
	corsMaxAge     = 8640
)

type Claims struct {
	Email    string `json:"email"`
	ClaimUID string `json:"claimID"`
	// registeredClaim's subject is userID from db
	jwt.RegisteredClaims
}

// slogMiddleware creates a Gin middleware that logs requests using slog
func slogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status and other metadata
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Prepare attributes
		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", time.Duration(latency.Microseconds())),
		}

		if raw != "" {
			attrs = append(attrs, slog.String("query", raw))
		}
		if errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		// Log based on status code
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		logger.LogAttrs(c.Request.Context(), level, "Incoming Request", attrs...)
	}
}

func loadCORScfg() cors.Config {
	origins := getEnvOrDefault("CORS_ORIGIN", allowedOrigins)
	methods := getEnvOrDefault("CORS_METHOD", allowedMethods)
	headers := getEnvOrDefault("CORS_HEADERS", allowedHeaders)
	maxAge := getEnvOrDefault("CORS_MAX_AGE", corsMaxAge)
	cfg := cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOrigins:     strings.Split(origins, ","),
		AllowHeaders:     strings.Split(headers, ","),
		AllowMethods:     strings.Split(methods, ","),
		MaxAge:           time.Duration(maxAge),
	}

	slog.Info("cors configuration", slog.Any("allow-credentials", cfg.AllowCredentials), slog.Any("allow-origins", cfg.AllowAllOrigins),
		slog.Any("origins", cfg.AllowOrigins), slog.Any("headers", cfg.AllowHeaders), slog.Any("methods", cfg.AllowMethods),
		slog.Any("Max age", cfg.MaxAge))

	return cfg
}

// AuthMiddleware check for valid cookie and extracts user id for the user
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := validateCookie(c)
		if err != nil {
			errors.HandleHTTPError(c, err)
			return
		}

		c.Set(string(models.CtxKeyUserID), *uid)
		c.Next()
	}
}

func validateCookie(c *gin.Context) (*uuid.UUID, error) {
	val, err := c.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(val, &Claims{}, func(token *jwt.Token) (any, error) {
		secVal := "33cea8f88c5c8ad73b1700af7d72891fe3097297e59fb6cbe5fd8b545a8316d0"
		return []byte(secVal), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.ErrInvalidCookie
	}

	if !token.Valid {
		return nil, errors.ErrInvalidCookie
	}

	return extractUserID(claims.Subject)
}

func extractUserID(claimSubject string) (*uuid.UUID, error) {
	uid, err := uuid.Parse(claimSubject)
	if err != nil {
		return nil, err
	}

	if uid == uuid.Nil {
		return nil, errors.ErrInvalidCookie
	}

	return &uid, nil
}
