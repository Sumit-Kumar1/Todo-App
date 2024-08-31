package handler

import (
	"encoding/json"
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
	tmpl := models.NewTemplate()

	return &Handler{template: tmpl, Service: s, Log: log}
}

// User API Handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterReq

	err := bind(r, &user)
	if err != nil {
		h.Log.Error("error while binding the request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.Service.Register(r.Context(), user)
	if err != nil {
		h.Log.Error("error while registering the user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		h.Log.Error("error while marshaling the response body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		h.Log.Error("error while writing the response body", "error", err)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.LoginReq

	err := bind(r, &user)
	if err != nil {
		h.Log.Error("error while binding the request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.Service.Login(r.Context(), user)
	if err != nil {
		h.Log.Error("error while registering the user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		h.Log.Error("error while marshaling the response body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		h.Log.Error("error while writing the response body", "error", err)
	}
}

// TODO API Handlers
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

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		h.Log.Error(notHTMX)
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

func (h *Handler) HandleIDReq(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		h.Log.Error(notHTMX)
		w.WriteHeader(http.StatusForbidden)

		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.deleteTask(w, r)
	case http.MethodPut:
		h.update(w, r)
	default:
		h.Log.Error(invalidReqMethod, "method", r.Method, "endpoint", "/delete")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		h.Log.Error(notHTMX)
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
