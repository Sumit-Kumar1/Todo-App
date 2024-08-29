package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"todoapp/internal/models"
)

const (
	hxRequest        = "Hx-Request"
	trueStr          = "true"
	invalidReqMethod = "%s method not allowed on %s"
	notHTMX          = "not a htmx request"
)

type Handler struct {
	template *template.Template
	Service  Servicer
	Log      *slog.Logger
}

func New(s Servicer, log *slog.Logger) *Handler {
	templateDir := "views/index.html"
	tmpl := template.Must(template.ParseFiles(templateDir))

	return &Handler{template: tmpl, Service: s, Log: log}
}

func (h *Handler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Log.Error(invalidReqMethod, "method", r.Method, "endpoint", "/")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	tasks, err := h.Service.GetAll(ctx)
	if err != nil {
		h.Log.Error("error in service GetAll", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.template.ExecuteTemplate(w, "index", map[string][]models.Task{
		"Data": tasks,
	})
	if err != nil {
		h.Log.Error("error executing template", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		h.Log.Error(notHTMX)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if r.Method != http.MethodPost {
		h.Log.Error(invalidReqMethod, "method", r.Method, "endpoint", "/add")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	task := r.PostFormValue("task")
	ctx := r.Context()

	t, err := h.Service.AddTask(ctx, task)
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

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		h.Log.Error(notHTMX)
		w.WriteHeader(http.StatusForbidden)

		return
	}

	if r.Method != http.MethodDelete {
		h.Log.Error(invalidReqMethod, "method", r.Method, "endpoint", "/delete")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")
	ctx := r.Context()

	h.Log.Info("Delete Request->", "ID", id)

	err := h.Service.DeleteTask(ctx, id)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		h.Log.Error("error while deleting task", "error", err)

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPut {
		h.Log.Error(invalidReqMethod, "method", r.Method, "endpoint", "/done")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")
	ctx := r.Context()

	h.Log.Info("Task Done -> ", "ID", id)

	resp, err := h.Service.MarkDone(ctx, id)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
			w.WriteHeader(http.StatusNotFound)
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		h.Log.Error(invalidReqMethod, "method", r.Method, "endpoint", "/update")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")
	title := r.Header.Get("HX-Prompt")
	ctx := r.Context()

	resp, err := h.Service.UpdateTask(ctx, id, title, "false")
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		h.Log.Error("error while updating task", "error", err)

		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := h.template.ExecuteTemplate(w, "add", *resp); err != nil {
		h.Log.Error("/updated template render", "ID", id, "error", err.Error())
		return
	}

	h.Log.Info("Task Updated", "ID", id)
}
