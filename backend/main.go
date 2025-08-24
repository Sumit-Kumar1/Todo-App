package main

import (
	"strconv"
	"todoapp/client"

	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	"todoapp/internal/migrations"
	todosvc "todoapp/internal/service/todo"
	usersvc "todoapp/internal/service/user"
	todostore "todoapp/internal/store"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/service"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	setupTasksRoutes(app)
	setupUserRoutes(app)

	app.Run()
}

func setupTasksRoutes(app *gofr.App) {
	todoStore := todostore.New()
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	app.POST("/task", todoHTTP.AddTask)
	app.GET("/tasks", todoHTTP.GetAllTasks)
	app.PUT("/tasks/{id}", todoHTTP.Update)
	app.PUT("/task/{id}/done", todoHTTP.Done)
	app.DELETE("/tasks/{id}/delete", todoHTTP.DeleteTask)
}

func setupUserRoutes(app *gofr.App) {
	authURL := app.Config.GetOrDefault("AUTH_URL", "http://localhost:9001")
	maxRetry, err := strconv.Atoi(app.Config.GetOrDefault("MAX_RETRY", "3"))
	if err != nil {
		maxRetry = 3
	}

	retryCfg := &service.RetryConfig{MaxRetries: maxRetry}
	healthCfg := &service.HealthConfig{HealthEndpoint: "health"}

	httpSvc := service.NewHTTPService(authURL, app.Logger(), app.Metrics(), retryCfg, healthCfg)

	authClient := client.New(httpSvc)
	userSvc := usersvc.New(authClient)
	userHTTP := userhttp.New(userSvc)

	app.POST("/register", userHTTP.Register)
	app.POST("/login", userHTTP.Login)
	app.POST("/logout", userHTTP.Logout)
}
