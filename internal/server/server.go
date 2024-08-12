package server

import (
	"net/http"
	"time"
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

func defaultServer() *Server {
	port := "9001"
	name := "todoApp"
	env := "dev"
	host := ""
	timeout := 10 * time.Second

	return &Server{
		Server: &http.Server{
			Addr:         host + ":" + port,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
			IdleTimeout:  2 * timeout,
		},
		Configs: &Configs{
			Name: name,
			Env:  env,
		},
	}
}
