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
	templateAddTask  = "add"
	templateIndex    = "index"
	userNotFound     = "user not found"
	renderErr        = "error while rendering template"
	hxRedirect       = "HX-Redirect"
)

type Handler struct {
	Service  TodoServicer
	template *template.Template
}

func New(todoSvc TodoServicer) *Handler {
	tmpl := models.NewTemplate()

	return &Handler{template: tmpl, Service: todoSvc}
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
		http.Error(w, userNotFound, http.StatusUnauthorized)
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

	if err := h.template.ExecuteTemplate(w, templateAddTask, resp.ToTaskResp()); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, renderErr, slog.String("template", templateAddTask))
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
		http.Error(w, userNotFound, http.StatusUnauthorized)
		return
	}

	t := models.TaskReq{
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
		DueDate:     r.PostFormValue("dueDate"),
	}

	task, err := h.Service.AddTask(ctx, &t, &userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.template.ExecuteTemplate(w, templateAddTask, task.ToTaskResp()); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, renderErr, slog.String("template", templateAddTask))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// nolint:revive // this is a handler get not returning
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
			w.Header().Add(hxRedirect, "/?page=register")
			w.WriteHeader(http.StatusOK)

			return
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", userID.String()))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	trs := []models.TaskResp{}

	for i := range tasks {
		trs = append(trs, *tasks[i].ToTaskResp())
	}

	if err := h.template.ExecuteTemplate(w, templateIndex, trs); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, renderErr, slog.String("template", templateIndex))
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
			http.Error(w, userNotFound, http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error(),
			slog.String("user", userID.String()),
			slog.String("task", id),
		)

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

	t := models.TaskReq{
		ID:          r.PathValue("id"),
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
		DueDate:     r.PostFormValue("dueDate"),
	}

	resp, err := h.Service.UpdateTask(ctx, t.ID, &t, false, &userID)
	if err != nil {
		switch {
		case models.ErrNotFound("user").Error() == err.Error():
			http.Error(w, userNotFound, http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		logger.LogAttrs(ctx, slog.LevelError, err.Error(),
			slog.String("user", userID.String()),
			slog.String("task", t.ID),
		)

		return
	}

	if resp == nil {
		http.Error(w, userNotFound, http.StatusNotFound)
		return
	}

	if err := h.template.ExecuteTemplate(w, templateAddTask, *resp.ToTaskResp()); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, renderErr, slog.String("template", templateAddTask))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task update done!",
		slog.String("user", userID.String()),
		slog.String("task", t.ID),
	)
}
