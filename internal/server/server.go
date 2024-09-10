package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configs struct {
	Name string `json:"name"`
	Env  string `json:"env"`
	Port int    `json:"port"`
}

type Server struct {
	Context *Context
	*http.ServeMux
	*Configs
}

type Opts func(s *Server)

func NewServer(opts ...Opts) *Server {
	s := defaultServer()
	s.Context = NewContext()

	if s.Context == nil {
		return nil
	}

	for _, fn := range opts {
		fn(s)
	}

	return s
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

func WithPort(port int) Opts {
	return func(s *Server) {
		s.Port = port
	}
}

func ServerFromEnvs() *Server {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("error while loading env file")
		return nil
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

	port := getEnvAsInt("HTTP_PORT", 9001)

	opts = append(opts, WithPort(port))

	return opts
}

func getEnvAsInt(key string, defaultVal int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}

	if val, err := strconv.Atoi(value); err == nil {
		return val
	}

	return defaultVal
}

func defaultServer() *Server {
	return &Server{
		Configs: &Configs{
			Name: "todoApp",
			Env:  "dev",
			Host: "",
			Port: "9001",
		},
		Context: &Context{},
	}
}
