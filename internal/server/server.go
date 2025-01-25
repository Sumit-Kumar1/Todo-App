package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sqlitecloud/sqlitecloud-go"
)

type Configs struct {
	Name string
	Env  string
	Host string
	Port string
}

type Server struct {
	DB          *sqlitecloud.SQCloud
	Logger      *slog.Logger
	ShutDownFxn func(context.Context) error
	Health      *Health
	Mux         *http.ServeMux
	*Configs
}

type Opts func(s *Server)

func NewServer(opts ...Opts) (*Server, error) {
	s := defaultServer()

	s.Logger = newLogger()

	db, err := newDB(s.Logger)
	if err != nil {
		return nil, err
	}

	s.DB = db

	for _, fn := range opts {
		fn(s)
	}

	return s, nil
}

func WithAppName(name string) Opts {
	return func(s *Server) {
		s.Name = name
	}
}

func WithEnv(env string) Opts {
	return func(s *Server) {
		s.Env = env
	}
}

func ServerFromEnvs() (*Server, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("error while loading env file")

		return nil, err
	}

	opts := loadEnvVars()

	return NewServer(opts...)
}

func loadEnvVars() []Opts {
	var opts []Opts

	appName := os.Getenv("APP_NAME")
	if appName != "" {
		opts = append(opts, WithAppName(appName))
	}

	env := os.Getenv("ENV")
	if env != "" {
		opts = append(opts, WithEnv(env))
	}

	return opts
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
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	return defaultValue
}
