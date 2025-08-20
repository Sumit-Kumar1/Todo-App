package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"todoapp/internal/handler"
	todohttp "todoapp/internal/handler/todo"
	userhttp "todoapp/internal/handler/user"
	"todoapp/internal/models"
	"todoapp/internal/service/todosvc"
	usersvc "todoapp/internal/service/user"
	sessionstore "todoapp/internal/store/session"
	todostore "todoapp/internal/store/todo"
	userstore "todoapp/internal/store/user"
)

func SetupRoutes(ctx context.Context, app *Server) {
	setupPublicRoutes(app)
	setupUserRoutes(app)
	setupTasksRoutes(ctx, app)
}

func setupTasksRoutes(ctx context.Context, app *Server) {
	templ := models.NewTemplate()
	todoStore := todostore.New(app.DB)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(templ, todoSvc)

	app.Mux.HandleFunc("/task",
		chain(todoHTTP.TaskPage, methodWithCORS(http.MethodGet), app.authMiddleware(ctx)))
	app.Mux.HandleFunc("/tasks",
		chain(todoHTTP.HandleTasks, isHTMX(), app.authMiddleware(ctx)))
	app.Mux.HandleFunc("/tasks/{id}",
		chain(todoHTTP.Update, methodWithCORS(http.MethodPut), isHTMX(), app.authMiddleware(ctx)))
	app.Mux.HandleFunc("/tasks/{id}/delete",
		chain(todoHTTP.DeleteTask, methodWithCORS(http.MethodDelete), isHTMX(), app.authMiddleware(ctx)))
	app.Mux.HandleFunc("/tasks/{id}/done",
		chain(todoHTTP.Done, methodWithCORS(http.MethodPut), isHTMX(), app.authMiddleware(ctx)))
}

func setupUserRoutes(app *Server) {
	templ := models.NewTemplate()
	usrSt := userstore.New(app.DB)
	sessionSt := sessionstore.New(app.DB)
	userSvc := usersvc.New(usrSt, sessionSt)
	usrHTTP := userhttp.New(templ, userSvc)

	app.Mux.HandleFunc("/register", chain(usrHTTP.Register, methodWithCORS(http.MethodPost)))
	app.Mux.HandleFunc("/login", chain(usrHTTP.Login, methodWithCORS(http.MethodPost), app.rateLimiterLogin()))
	app.Mux.HandleFunc("/logout", chain(usrHTTP.Logout, methodWithCORS(http.MethodPost)))
}

func setupPublicRoutes(app *Server) {
	templ := models.NewTemplate()
	h := handler.New(templ)

	public := http.FileServer(http.Dir("public"))
	openapi := http.FileServer(http.Dir("openapi"))

	app.Mux.HandleFunc("/", chain(h.Root, methodWithCORS(http.MethodGet)))
	app.Mux.Handle("/public/", http.StripPrefix("/public/", public))
	app.Mux.Handle("/openapi/", http.StripPrefix("/openapi/", openapi))
	app.Mux.Handle("/api/", http.StripPrefix("/api", chain(h.Swagger, methodWithCORS(http.MethodGet))))
	app.Mux.Handle("/healthz", chain(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		app.Health = &Health{DB: false, Service: false, Msg: "started health"}

		if app.DB == nil {
			app.Health.DB = false
			app.Health.Msg = "DB is nil"

			app.Logger.LogAttrs(r.Context(), slog.LevelError, "application DB is nil")

			http.Error(w, "db is nil", http.StatusInternalServerError)

			return
		}

		if err := app.DB.Ping(); err != nil {
			app.Health.DB = false
			app.Health.Msg = err.Error()
		}

		if isServiceHealthy(r.Context(), app.Port) {
			app.Health.Service = false
			app.Health.Msg = "Service Down"
		}

		app.Health.Service = true
		app.Health.DB = true
		app.Health.Msg = "OK"

		data, err := json.Marshal(app.Health)
		if err != nil {
			app.Logger.LogAttrs(r.Context(), slog.LevelError, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)

		endTime := time.Since(t)
		app.Logger.LogAttrs(r.Context(), slog.LevelInfo, "GET /healthz",
			slog.Any("status", app.Health),
			slog.Int64("time taken(ms)", endTime.Milliseconds()),
		)
	}, methodWithCORS(http.MethodGet)))
}

func isServiceHealthy(ctx context.Context, port string) bool {
	var url = "http://localhost:" + port

	client := http.Client{Timeout: 100 * time.Millisecond}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return false
	}
	// nolint:bodyclose // body is already closed in defer statement
	resp, err := client.Do(r)
	if err != nil {
		return false
	}

	defer func(ctx context.Context, body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
			return
		}
	}(ctx, resp.Body)

	return resp.StatusCode == http.StatusOK
}
