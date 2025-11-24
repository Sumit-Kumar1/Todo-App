// Package server init server and configure server with all required hard deps
package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Configs struct {
	Name            string
	Env             string
	Host            string
	Port            string
	MigrationMethod string
}

type Health struct {
	Msg     string `json:"Message"`
	DB      bool   `json:"DB"`
	Service bool   `json:"Service"`
}

type rateLimiter struct {
	attempts    map[string]*limiterAttempt
	timeWindow  time.Duration
	maxAttempts int
	mu          sync.Mutex
}

type limiterAttempt struct {
	firstTime time.Time
	count     int
}

type Server struct {
	DB            *sql.DB
	Logger        *slog.Logger
	ShutDownFxn   func(context.Context) error
	Mux           *http.ServeMux
	Health        *Health
	loginLimiter  *rateLimiter
	globalLimiter *rateLimiter
	logFile       *os.File
	Configs       Configs
	*http.Server
}

type ServerBuilder struct {
	server *Server
}

func NewServerBuilder() *ServerBuilder {
	s := defaultServer()

	return s
}

func (sb *ServerBuilder) WithHostPort() *ServerBuilder {
	host := GetEnvOrDefault("HOST_HOST", "localhost")
	port := GetEnvOrDefault("HTTP_PORT", "9003")

	httpServer := &http.Server{
		Addr: net.JoinHostPort(host, port),
	}

	sb.server.Server = httpServer

	return sb
}

func (sb *ServerBuilder) WithLogger() *ServerBuilder {
	sb.server.Logger, sb.server.logFile = newLogger()
	return sb
}

func (sb *ServerBuilder) WithConfig() *ServerBuilder {
	sb.server.Configs.Env = GetEnvOrDefault("ENV", "dev")
	sb.server.Configs.Name = GetEnvOrDefault("APP_NAME", "todo-app")
	sb.server.Configs.MigrationMethod = GetEnvOrDefault("MIGRATION_METHOD", "UP")
	return sb
}

func (sb *ServerBuilder) WithDB() *ServerBuilder {
	sb.server.DB, _ = newDB()
	return sb
}

func (sb *ServerBuilder) WithServerTimeouts() *ServerBuilder {
	readTimeout := GetEnvOrDefault("READ_TIMEOUT", 2)
	writeTimeout := GetEnvOrDefault("WRITE_TIMEOUT", 3)
	idleTimeout := GetEnvOrDefault("IDLE_TIMEOUT", 5)

	sb.server.ReadTimeout = time.Duration(readTimeout * int(time.Second))
	sb.server.WriteTimeout = time.Duration(writeTimeout * int(time.Second))
	sb.server.IdleTimeout = time.Duration(idleTimeout * int(time.Second))

	return sb
}

func (sb *ServerBuilder) Build() *Server {
	sb.server.ShutDownFxn = sb.server.createShutdownFunction()
	return sb.server
}

func defaultServer() *ServerBuilder {
	s := &Server{
		Configs: Configs{
			Name:            "todoApp",
			Env:             "dev",
			MigrationMethod: "UP",
		},
		Mux: http.NewServeMux(),
		Health: &Health{
			DB:      false,
			Service: false,
			Msg:     "INIT HEALTH",
		},
		globalLimiter: &rateLimiter{
			attempts:    make(map[string]*limiterAttempt),
			maxAttempts: GetEnvOrDefault("GLOBAL_ATTEMPTS", 300),
			timeWindow:  time.Second * time.Duration(GetEnvOrDefault("GLOBAL_TIME_WINDOW", 60)),
		},
		loginLimiter: &rateLimiter{
			attempts:    make(map[string]*limiterAttempt),
			maxAttempts: GetEnvOrDefault("LOGIN_ATTEMPTS", 10),
			timeWindow:  time.Second * time.Duration(GetEnvOrDefault("LOGIN_TIME_WINDOW", 60)),
		},
	}

	return &ServerBuilder{server: s}
}

func (s *Server) createShutdownFunction() func(context.Context) error {
	return func(ctx context.Context) error {
		// Close the log file
		if s.logFile != nil {
			if err := s.logFile.Close(); err != nil {
				s.Logger.Error("failed to close log file", "error", err)
			} else {
				s.Logger.Info("log file closed successfully")
			}
		}

		// Close database connection
		if s.DB != nil {
			if err := s.DB.Close(); err != nil {
				s.Logger.Error("failed to close database connection", "error", err)
			} else {
				s.Logger.Info("database connection closed successfully")
			}
		}

		s.Logger.Info("server shutdown completed")
		return nil
	}
}

func GetEnvOrDefault[T string | int | bool](key string, defaultValue T) T {
	env := os.Getenv(key)
	if strings.TrimSpace(env) == "" {
		return defaultValue
	}

	var zero T
	switch any(zero).(type) {
	case string:
		return any(env).(T)
	case int:
		if iVal, err := strconv.Atoi(env); err == nil {
			return any(iVal).(T)
		}
		return defaultValue
	case bool:
		if bVal, err := strconv.ParseBool(env); err == nil {
			return any(bVal).(T)
		}
		return defaultValue
	default:
		return defaultValue
	}
}
