package todohttp

import (
	"html/template"
	"log/slog"
	"net/http"
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

const invalidReqMethod = "method not allowed"

type Handler struct {
	Service  TodoServicer
	Log      *slog.Logger
	template *template.Template
}

func New(todoSvc TodoServicer, log *slog.Logger) *Handler {
	tmpl := models.NewTemplate()

	return &Handler{template: tmpl, Log: log, Service: todoSvc}
}

// Root rendering endpoints
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	var tempName string

	vals := r.URL.Query()
	switch vals.Get("page") {
	case "register":
		tempName = "user-register"
	default:
		tempName = "user-login"
	}

	if err := h.template.ExecuteTemplate(w, tempName, nil); err != nil {
		h.Log.Error(err.Error(), "template-render", "index")
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
	ctx := r.Context()
	userID, ok := ctx.Value(server.CtxKey).(uuid.UUID)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	h.Log.Info("Task Done -> ", "ID", id)

	resp, err := h.Service.MarkDone(ctx, id, &userID)
	if err != nil {
		switch {
		case err.Error() == models.ErrNotFound("task").Error():
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		h.Log.Error("error while marking task done", "error", err)

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := h.template.ExecuteTemplate(w, "add", *resp); err != nil {
		h.Log.Error("template done render", "ID", id, "error", err.Error())
	}
}

func (h *Handler) addTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(server.CtxKey).(uuid.UUID)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	task := r.PostFormValue("task")

	t, err := h.Service.AddTask(ctx, task, &userID)
	if err != nil {
		h.Log.Error("error while adding task", "error", err)
		w.WriteHeader(http.StatusBadRequest)

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	h.Log.Info("Task is Added", "ID", t.ID)

	if err := h.template.ExecuteTemplate(w, "add", *t); err != nil {
		h.Log.Error("error while executing template:", "error", err.Error())
		return
	}
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(server.CtxKey).(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	tasks, err := h.Service.GetAll(r.Context(), &userID)
	if err != nil {
		if models.ErrNotFound("user").Error() == err.Error() {
			w.Header().Add("HX-Redirect", "/?page=register")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.Log.Error(err.Error(), "request", "service-getAll")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	if err := h.template.ExecuteTemplate(w, "index", tasks); err != nil {
		h.Log.Error(err.Error(), "request", "service-getAll")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(server.CtxKey).(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	h.Log.Info("Delete Request->", "ID", id)

	err := h.Service.DeleteTask(ctx, id, &userID)
	if err != nil {
		switch {
		case models.ErrNotFound("user").Error() == err.Error():
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		h.Log.Error(err.Error(), "request", "handler-delete")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(server.CtxKey).(uuid.UUID)
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

		h.Log.Error(err.Error(), "request", "handler-update")
		return
	}

	if resp == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if err := h.template.ExecuteTemplate(w, "add", *resp); err != nil {
		h.Log.Error("/updated template render", "ID", id, "error", err.Error())
		return
	}

	h.Log.Info("Task Updated", "ID", id)
}
