package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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
	*http.Server
	*Configs
}

type Opts func(s *Server)

func NewServer(opts ...Opts) *Server {
	s := defaultServer()
	for _, fn := range opts {
		fn(s)
	}
	return s
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

func ServerFromEnvs() *Server {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("error while loading env file")
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

func defaultServer() *Server {
	return &Server{
		Server: &http.Server{
			Addr:         ":9001",
			ReadTimeout:  3 * time.Minute,
			WriteTimeout: 3 * time.Minute,
			IdleTimeout:  5 * time.Minute,
		},
		Configs: &Configs{
			Name: "todoApp",
			Env:  "dev",
		},
	}
}
