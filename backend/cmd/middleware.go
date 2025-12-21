package cmd

import (
	"log/slog"
	"strings"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

const (
	cookieName     = "auth"
	allowedOrigins = "http://localhost:8081"
	allowedMethods = "POST,GET,PUT,DELETE,PATCH,OPTIONS"
	allowedHeaders = "Accept,Content-Type,Content-Length,Accept-Encoding"
	corsMaxAge     = 8640
)

type Auth interface {
	Validate(ctx *gin.Context, token string) (*uuid.UUID, error)
}

type clientMiddleware struct {
	client Auth
}

func newMiddleware(client Auth) clientMiddleware {
	return clientMiddleware{
		client: client,
	}
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
			slog.Duration("latency(ms)", time.Duration(latency.Milliseconds())),
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
func (cc clientMiddleware) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := cc.validateCookie(c)
		if err != nil {
			errors.HandleHTTPError(c, err)
			return
		}

		c.Set(string(models.CtxKeyUserID), *uid)
		c.Next()
	}
}

func (cc clientMiddleware) validateCookie(c *gin.Context) (*uuid.UUID, error) {
	val, err := c.Cookie(cookieName)
	if err != nil {
		return nil, errors.ErrInvalidCookie
	}

	userID, err := cc.client.Validate(c, val)
	if err != nil {
		return nil, errors.ErrInvalidCookie
	}

	return userID, nil
}
