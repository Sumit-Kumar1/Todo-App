package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sqlitecloud/sqlitecloud-go"
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
	DBStatus      bool   `json:"dbStatus"`
	ServiceStatus bool   `json:"serviceStatus"`
	Msg           string `json:"msg"`
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
	DB            *sqlitecloud.SQCloud
	Logger        *slog.Logger
	ShutDownFxn   func(context.Context) error
	Mux           *http.ServeMux
	Health        *Health
	loginLimiter  *rateLimiter
	globalLimiter *rateLimiter
	*Configs
}

type Opts func(s *Server)

func NewServer() (*Server, error) {
	s, err := configureServer()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func configureServer() (*Server, error) {
	s := defaultServer()

	if err := godotenv.Load(".env"); err != nil {
		log.Print("error while loading env file")

		return nil, err
	}

	s.Name = getEnvOrDefault("APP_NAME", "todo-app")
	s.Port = getEnvOrDefault("HTTP_PORT", "9001")
	s.Env = getEnvOrDefault("ENV", "dev")
	s.ReadTimeout = getEnvAsInt("READ_TIMEOUT", 2)
	s.WriteTimeout = getEnvAsInt("WRITE_TIMEOUT", 3)
	s.IdleTimeout = getEnvAsInt("IDLE_TIMEOUT", 5)
	s.MigrationMethod = getEnvOrDefault("MIGRATION_METHOD", "UP")

	s.Logger = newLogger()

	db, err := newDB(s.Logger)
	if err != nil {
		return nil, err
	}

	s.DB = db

	return s, nil
}

func defaultServer() *Server {
	return &Server{
		Configs: &Configs{
			Name: "todoApp",
			Env:  "dev",
			Host: "localhost",
			Port: "9001",
		},
		Mux: http.NewServeMux(),
		Health: &Health{
			DBStatus:      false,
			ServiceStatus: false,
			Msg:           "INIT HEALTH",
		},
		globalLimiter: &rateLimiter{
			attempts:    make(map[string]*limiterAttempt),
			timeWindow:  time.Minute * 1,
			maxAttempts: 20,
		},
		loginLimiter: &rateLimiter{
			attempts:    make(map[string]*limiterAttempt),
			maxAttempts: 5,
			timeWindow:  time.Minute * 1,
		},
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}

	if iVal, err := strconv.Atoi(env); err == nil {
		return iVal
	}

	return defaultValue
}

func getEnvOrDefault(key, def string) string {
	eval := os.Getenv(key)
	if eval == "" {
		return def
	}

	return eval
}
