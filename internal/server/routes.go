package server

import (
	"context"
	"net/http"

	"todoapp/internal/handler"
	"todoapp/internal/service/todosvc"
	"todoapp/internal/store/todostore"

	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	usersvc "todoapp/internal/service/user"
	userstore "todoapp/internal/store/user"
)

func SetupRoutes(ctx context.Context, app *Server) {
	setupPublicRoutes(app)
	setupUserRoutes(app)
	setupTasksRoutes(ctx, app)
}

func setupTasksRoutes(ctx context.Context, app *Server) {
	todoStore := todostore.New(app.DB)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	app.Mux.HandleFunc("/task", Chain(todoHTTP.TaskPage, Method(http.MethodGet), AuthMiddleware(ctx, app.DB)))
	app.Mux.HandleFunc("/tasks", Chain(todoHTTP.HandleTasks, IsHTMX(), AuthMiddleware(ctx, app.DB)))
	app.Mux.HandleFunc("/tasks/{id}", Chain(todoHTTP.Update, IsHTMX(), Method(http.MethodPut),
		AuthMiddleware(ctx, app.DB)))
	app.Mux.HandleFunc("/tasks/{id}/delete", Chain(todoHTTP.DeleteTask, IsHTMX(), AuthMiddleware(ctx, app.DB),
		Method(http.MethodDelete)))
	app.Mux.HandleFunc("/tasks/{id}/done", Chain(todoHTTP.Done, IsHTMX(), Method(http.MethodPut),
		AuthMiddleware(context.Background(), app.DB)))
}

func setupUserRoutes(app *Server) {
	usrSt := userstore.New(app.DB)
	userSvc := usersvc.New(usrSt)
	usrHTTP := userhttp.New(userSvc)

	app.Mux.HandleFunc("/register", Chain(usrHTTP.Register, Method(http.MethodPost)))
	app.Mux.HandleFunc("/login", Chain(usrHTTP.Login, Method(http.MethodPost)))
	app.Mux.HandleFunc("/logout", Chain(usrHTTP.Logout, Method(http.MethodPost)))
}

func setupPublicRoutes(app *Server) {
	h := handler.New()

	public := http.FileServer(http.Dir("public"))
	openapi := http.FileServer(http.Dir("openapi"))

	app.Mux.HandleFunc("/", Chain(h.Root, Method(http.MethodGet)))
	app.Mux.Handle("/public/", http.StripPrefix("/public/", public))
	app.Mux.Handle("/openapi/", http.StripPrefix("/openapi/", openapi))
	app.Mux.Handle("/api", http.StripPrefix("/api", Chain(h.Swagger, Method(http.MethodGet))))
}
