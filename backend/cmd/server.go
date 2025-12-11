// Package cmd setups go server and start the API
package cmd

import (
	"context"
	"database/sql"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"todoapp/client"
	"todoapp/internal/handler"
	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	"todoapp/internal/migrations"
	todosvc "todoapp/internal/service/todo"
	usersvc "todoapp/internal/service/user"
	"todoapp/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const ()

func Run(c context.Context, w io.Writer) error {
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	router := gin.New()

	router.Use(slogMiddleware(logger))
	router.Use(gin.Recovery())
	router.Use(cors.New(loadCORScfg()))

	db, err := newDB(c)
	if err != nil {
		return err
	}

	if err := migrations.RunMigrations(c, db, getEnvOrDefault("MIGRATION_METHOD", "UP")); err != nil {
		log.Printf("error while running migrations: %s", err.Error())

		return err
	}

	setupHealthRoutes(router, db)
	setupUserRoutes(router, db)
	setupTasksRoutes(router, db)

	host := getEnvOrDefault("HTTP_HOST", "localhost")
	port := getEnvOrDefault("HTTP_PORT", "9003")

	slog.Info("starting server", slog.String("host", host), slog.String("port", port))

	return router.Run(net.JoinHostPort(host, port))
}

// setupHealthRoutes configures health check endpoints
func setupHealthRoutes(r *gin.Engine, db *sql.DB) {
	hh := handler.NewHealthHandler(db)

	r.GET("/health", hh.HealthHandler)
	r.GET("/ready", hh.ReadyHandler)
	r.GET("/live", hh.LivenessHandler)
}

func setupTasksRoutes(r *gin.Engine, db *sql.DB) {
	todoStore := store.New(db)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	r.POST("/task", authMiddleware(), todoHTTP.AddTask)
	r.GET("/tasks", authMiddleware(), todoHTTP.GetAllTasks)
	r.PUT("/tasks/:id", authMiddleware(), todoHTTP.Update)
	r.DELETE("/tasks/:id", authMiddleware(), todoHTTP.DeleteTask)
	r.PATCH("/tasks/:id/done", authMiddleware(), todoHTTP.MarkDone)
}

func setupUserRoutes(r *gin.Engine, app *sql.DB) {
	authURL := getEnvOrDefault("AUTH_URL", "http://localhost:9001")
	authClient := client.New(authURL)
	userSvc := usersvc.New(authClient)
	userHTTP := userhttp.New(userSvc)

	r.POST("/register", userHTTP.Register)
	r.POST("/login", userHTTP.Login)
	r.POST("/logout", userHTTP.Logout)
}

func getEnvOrDefault[T string | int | bool](key string, defaultValue T) T {
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
