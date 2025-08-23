package server

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
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
			DB:      false,
			Service: false,
			Msg:     "INIT HEALTH",
		},
		globalLimiter: &rateLimiter{
			attempts:    make(map[string]*limiterAttempt),
			maxAttempts: getEnvAsInt("GLOBAL_ATTEMPTS", 300),
			timeWindow:  time.Second * time.Duration(getEnvAsInt("GLOBAL_TIME_WINDOW", 60)),
		},
		loginLimiter: &rateLimiter{
			attempts:    make(map[string]*limiterAttempt),
			maxAttempts: getEnvAsInt("LOGIN_ATTEMPTS", 10),
			timeWindow:  time.Second * time.Duration(getEnvAsInt("LOGIN_TIME_WINDOW", 60)),
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
