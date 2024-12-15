package todohttp

import (
	"html/template"
	"log/slog"
	"net/http"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	invalidReqMethod = "method not allowed"
	templAddTask     = "add"
	templIndex       = "index"
)

type Handler struct {
	Service  TodoServicer
	template *template.Template
}

func New(todoSvc TodoServicer) *Handler {
	tmpl := models.NewTemplate()

	return &Handler{template: tmpl, Service: todoSvc}
}

// Root rendering endpoints
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	var (
		template string
		ctx      = r.Context()
		logger   = models.GetLoggerFromCtx(ctx)
	)

	vals := r.URL.Query()
	switch vals.Get("page") {
	case "register":
		template = "user-register"
	case "api":
		template = "swagger"
	default:
		template = "user-login"
	}

	if err := h.template.ExecuteTemplate(w, template, nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while rendering template", slog.String("template", "index"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Swagger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	if err := h.template.ExecuteTemplate(w, "swagger", nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while rendering template", slog.String("template", "swagger"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) TaskPage(w http.ResponseWriter, r *http.Request) {
	h.getAll(w, r)
}

func (h *Handler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.addTask(w, r)
	default:
		http.Error(w, invalidReqMethod, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	resp, err := h.Service.MarkDone(ctx, id, &userID)
	if err != nil {
		switch {
		case err.Error() == models.ErrNotFound("task").Error():
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		logger.LogAttrs(ctx, slog.LevelError, "error while marking task done",
			slog.String("error", err.Error()), slog.String("task", id))

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := h.template.ExecuteTemplate(w, templAddTask, *resp); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while rendering template", slog.String("template", templAddTask))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) addTask(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	task := r.PostFormValue("task")

	t, err := h.Service.AddTask(ctx, task, &userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.template.ExecuteTemplate(w, templAddTask, *t); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while rendering template", slog.String("template", templAddTask))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		_ = h.template.ExecuteTemplate(w, "errorPage", map[string]any{
			"Code":    http.StatusUnauthorized,
			"Message": "user not authorized!!",
		})

		return
	}

	tasks, err := h.Service.GetAll(r.Context(), &userID)
	if err != nil {
		if models.ErrNotFound("user").Error() == err.Error() {
			w.Header().Add("HX-Redirect", "/?page=register")
			w.WriteHeader(http.StatusOK)
			return
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", userID.String()))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	if err := h.template.ExecuteTemplate(w, templIndex, tasks); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while rendering template", slog.String("template", templIndex))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	if err := h.Service.DeleteTask(ctx, id, &userID); err != nil {
		switch {
		case models.ErrNotFound("user").Error() == err.Error():
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", userID.String()), slog.String("task", id))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	title := r.Header.Get("HX-Prompt")

	resp, err := h.Service.UpdateTask(ctx, id, title, false, &userID)
	if err != nil {
		switch {
		case models.ErrNotFound("user").Error() == err.Error():
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", userID.String()), slog.String("task", id))
		return
	}

	if resp == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if err := h.template.ExecuteTemplate(w, templAddTask, *resp); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while rendering template", slog.String("template", templAddTask))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task update done!", slog.String("user", userID.String()), slog.String("task", id))
}
