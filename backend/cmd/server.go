// Package cmd setups go server and start the API
package cmd

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
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

	"github.com/gin-gonic/gin"
)

func Run(c context.Context, _ io.Writer, _ []string) error {
	router := gin.Default()

	db, err := newDB(c)
	if err != nil {
		return err
	}

	if err := migrations.RunMigrations(c, db, "UP"); err != nil {
		slog.LogAttrs(c, slog.LevelError, "error while running migrations",
			slog.String("error", err.Error()))

		return err
	}

	// Setup health endpoints firsts
	setupHealthRoutes(router, db)
	setupUserRoutes(router, db)
	setupTasksRoutes(router, db)

	return router.Run("localhost:9003")
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

	r.POST("/task", todoHTTP.AddTask)
	r.GET("/tasks", todoHTTP.GetAllTasks)
	r.PATCH("/tasks/:id/done", todoHTTP.MarkDone)
	r.PUT("/tasks/:id", todoHTTP.Update)
	r.DELETE("/tasks/:id", todoHTTP.DeleteTask)
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
