package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sqlitecloud/sqlitecloud-go"
)

type Configs struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

type Health struct {
	DBStatus string `json:"dbStatus"`
	Status   string `json:"status"`
}

type Server struct {
	DB          *sqlitecloud.SQCloud
	Logger      *slog.Logger
	ShutDownFxn func(context.Context) error
	Health      *Health
	*http.Server
	*Configs
}

type Opts func(s *Server)

func NewServer(opts ...Opts) (*Server, error) {
	s := defaultServer()

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

func WithTimeouts(read, write, idle int) Opts {
	return func(s *Server) {
		s.ReadTimeout = time.Duration(read) * time.Second
		s.WriteTimeout = time.Duration(write) * time.Second
		s.IdleTimeout = time.Duration(idle) * time.Second
	}
}

func WithPort(port string) Opts {
	return func(s *Server) {
		s.Addr = ":" + port
	}
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

	port := os.Getenv("HTTP_PORT")
	if port != "" {
		opts = append(opts, WithPort(port))
	}

	env := os.Getenv("ENV")
	if env != "" {
		opts = append(opts, WithEnv(env))
	}

	readTimeout := getEnvAsInt("READ_TIMEOUT", 180)   // Default to 3 minutes
	writeTimeout := getEnvAsInt("WRITE_TIMEOUT", 180) // Default to 3 minutes
	idleTimeout := getEnvAsInt("IDLE_TIMEOUT", 300)   // Default to 5 minutes

	opts = append(opts, WithTimeouts(readTimeout, writeTimeout, idleTimeout))
	return opts
}

func defaultServer() *Server {
	return &Server{
		Server: &http.Server{
			Addr:         ":9001",
			ReadTimeout:  time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  20 * time.Second,
		},

		Configs: &Configs{
			Name: "todoApp",
			Env:  "dev",
		},

		Logger: newLogger(),
	}
}

func newLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	slog.SetDefault(logger)

	return logger
}

func newDB(logger *slog.Logger) (*sqlitecloud.SQCloud, error) {
	config := sqlitecloud.SQCloudConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT", 8860),
		Database: os.Getenv("DB_NAME"),
		ApiKey:   os.Getenv("DB_APIKEY"),
		MaxRows:  getEnvAsInt("DB_MAXROWS", 20),
	}

	isSecure, err := strconv.ParseBool(os.Getenv("DB_SECURE_FLAG"))
	if err != nil {
		return nil, err
	}

	config.Secure = isSecure

	sqcl := sqlitecloud.New(config)

	if err := sqcl.Connect(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if !sqcl.IsConnected() {
		return nil, fmt.Errorf("not able to connect to database")
	}

	logger.Info("DB connected successfully")

	return sqcl, nil
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
