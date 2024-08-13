package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	timeout := 10 * time.Second

	tests := []struct {
		name string
		opts []Opts
		want *Server
	}{
		{name: "default server", opts: nil, want: &Server{
			Server:  &http.Server{Addr: ":9001", ReadTimeout: timeout, IdleTimeout: 2 * timeout, WriteTimeout: timeout},
			Configs: &Configs{Name: "todoApp", Env: "dev"},
		}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.opts...)

			assert.Equalf(t, tt.want.Addr, got.Addr, "Test[%d] Failed Address - %s", i, tt.name)
			assert.Equalf(t, tt.want.Name, got.Name, "Test[%d] Failed AppName - %s", i, tt.name)
			assert.Equalf(t, tt.want.ReadTimeout, got.ReadTimeout, "Test[%d] Failed ReadTime - %s", i, tt.name)
			assert.Equalf(t, tt.want.WriteTimeout, got.WriteTimeout, "Test[%d] Failed WriteTime - %s", i, tt.name)
			assert.Equalf(t, tt.want.IdleTimeout, got.IdleTimeout, "Test[%d] Failed IdleTime - %s", i, tt.name)
			assert.Equalf(t, tt.want.Env, got.Env, "Test[%d] Failed Env - %s", i, tt.name)
		})
	}
}

func Test_defaultServer(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		{name: "nil case", want: &Server{
			Server: &http.Server{Addr: ":9001", ReadTimeout: 10 * time.Second,
				WriteTimeout: 10 * time.Second, IdleTimeout: 20 * time.Second},
			Configs: &Configs{Name: "todoApp", Env: "dev"},
		}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultServer()
			assert.Equalf(t, tt.want, got, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
