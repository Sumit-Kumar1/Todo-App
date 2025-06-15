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
	todoStore := todostore.New(app.DB)
	todoSvc := todosvc.New(todoStore)
	todoHTTP := todohttp.New(todoSvc)

	app.Mux.HandleFunc("/task",
		app.chain(todoHTTP.TaskPage,
			app.method(http.MethodGet),
			app.authMiddleware(ctx),
		),
	)
	app.Mux.HandleFunc("/tasks",
		app.chain(todoHTTP.HandleTasks,
			app.isHTMX(),
			app.authMiddleware(ctx)))
	app.Mux.HandleFunc("/tasks/{id}",
		app.chain(todoHTTP.Update,
			app.isHTMX(),
			app.method(http.MethodPut),
			app.authMiddleware(ctx),
		))
	app.Mux.HandleFunc("/tasks/{id}/delete",
		app.chain(todoHTTP.DeleteTask,
			app.isHTMX(),
			app.authMiddleware(ctx),
			app.method(http.MethodDelete),
		),
	)
	app.Mux.HandleFunc("/tasks/{id}/done",
		app.chain(todoHTTP.Done,
			app.isHTMX(),
			app.method(http.MethodPut),
			app.authMiddleware(context.Background()),
		))
}

func setupUserRoutes(app *Server) {
	usrSt := userstore.New(app.DB)
	sessionSt := sessionstore.New(app.DB)
	userSvc := usersvc.New(usrSt, sessionSt)
	usrHTTP := userhttp.New(userSvc)

	app.Mux.HandleFunc("/register", app.chain(usrHTTP.Register, app.method(http.MethodPost)))
	app.Mux.HandleFunc("/login", app.chain(usrHTTP.Login, app.method(http.MethodPost), app.rateLimiterLogin()))
	app.Mux.HandleFunc("/logout", app.chain(usrHTTP.Logout, app.method(http.MethodPost)))
}

func setupPublicRoutes(app *Server) {
	h := handler.New()

	public := http.FileServer(http.Dir("public"))
	openapi := http.FileServer(http.Dir("openapi"))

	app.Mux.HandleFunc("/", app.chain(h.Root, app.method(http.MethodGet)))
	app.Mux.Handle("/public/", http.StripPrefix("/public/", public))
	app.Mux.Handle("/openapi/", http.StripPrefix("/openapi/", openapi))
	app.Mux.Handle("/api", http.StripPrefix("/api", app.chain(h.Swagger, app.method(http.MethodGet))))
	app.Mux.Handle("/healthz", app.chain(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		app.Health = &Health{DBStatus: false, ServiceStatus: false, Msg: "StartedHealth"}

		if app.DB != nil && !app.DB.IsConnected() {
			app.Health.DBStatus = false
			app.Health.Msg = "DB is not connected"
		}

		if err := app.DB.Ping(); err != nil {
			app.Health.DBStatus = false
			app.Health.Msg = err.Error()
		}

		if isServiceHealthy(r.Context(), app.Port) {
			app.Health.ServiceStatus = false
			app.Health.Msg = "Service Down"
		}

		app.Health.ServiceStatus = true
		app.Health.DBStatus = true
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
		app.Logger.LogAttrs(r.Context(), slog.LevelInfo, "GET/healthz",
			slog.Any("status", app.Health),
			slog.Int64("time taken(ms)", endTime.Milliseconds()),
		)
	}, app.method(http.MethodGet)))
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
