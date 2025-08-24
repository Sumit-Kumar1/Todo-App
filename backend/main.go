package main

import (
	"strconv"
	"todoapp/client"
	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	"todoapp/internal/migrations"
	"todoapp/internal/service/todosvc"
	todostore "todoapp/internal/store/todo"

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
	usrHTTP := userhttp.New(authClient)

	app.POST("/register", usrHTTP.Register)
	app.POST("/login", usrHTTP.Login)
	app.POST("/logout", usrHTTP.Logout)
}
