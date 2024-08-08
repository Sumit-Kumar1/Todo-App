package server

import (
	"net/http"
	"time"
)

type Opts func(*Server)

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

func NewServer(opts ...Opts) *Server {
	s := defaultState()

	for _, fn := range opts {
		fn(s)
	}

	return s
}

func WithTimeouts(read, write, idle int) Opts {
	return func(s *Server) {
		s.ReadTimeout = time.Duration(read * int(time.Second))
		s.WriteTimeout = time.Duration(write * int(time.Second))
		s.IdleTimeout = time.Duration(idle * int(time.Second))
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

func defaultState() *Server {
	port := "9001"
	name := "todoApp"
	env := "dev"
	host := ""

	return &Server{
		Server: &http.Server{
			Addr:         host + ":" + port,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  20 * time.Second,
		},
		Configs: &Configs{
			Name: name,
			Env:  env,
		},
	}
}
