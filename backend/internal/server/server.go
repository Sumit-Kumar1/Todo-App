// Package server init server and configure server with all required hard deps
package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	serverOnce     sync.Once
	serverInstance *Server
)

type Configs struct {
	Name            string
	Env             string
	Host            string
	Port            string
	ReadTimeout     int
	WriteTimeout    int
	IdleTimeout     int
	MigrationMethod string
}

type Health struct {
	DB      bool   `json:"DB"`
	Service bool   `json:"Service"`
	Msg     string `json:"Message"`
}

type rateLimiter struct {
	mu          sync.Mutex
	attempts    map[string]*limiterAttempt
	maxAttempts int
	timeWindow  time.Duration
}

type limiterAttempt struct {
	count     int
	firstTime time.Time
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
	*Configs
}

type Opts func(s *Server)

func NewServer() (*Server, error) {
	if serverInstance != nil {
		return serverInstance, nil
	}

	s, err := configureServer()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func configureServer() (*Server, error) {
	s := defaultServer()
	s.Name = GetEnvOrDefault("APP_NAME", "todo-app")
	s.Port = GetEnvOrDefault("HTTP_PORT", "9003")
	s.Env = GetEnvOrDefault("ENV", "dev")
	s.ReadTimeout = GetEnvOrDefault("READ_TIMEOUT", 2)
	s.WriteTimeout = GetEnvOrDefault("WRITE_TIMEOUT", 3)
	s.IdleTimeout = GetEnvOrDefault("IDLE_TIMEOUT", 5)
	s.MigrationMethod = GetEnvOrDefault("MIGRATION_METHOD", "UP")

	s.Logger, s.logFile = newLogger()

	db, err := newDB(s.Logger)
	if err != nil {
		return nil, err
	}

	s.DB = db

	// Set up shutdown function
	s.ShutDownFxn = s.createShutdownFunction()

	return s, nil
}

func defaultServer() *Server {
	if serverInstance != nil {
		return serverInstance
	}

	serverOnce.Do(func() {
		serverInstance = &Server{
			Configs: &Configs{
				Name: "todoApp",
				Env:  "dev",
				Host: "localhost",
				Port: "9003",
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
	})

	return serverInstance
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
