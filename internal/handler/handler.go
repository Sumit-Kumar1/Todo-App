package handler

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	appJSON          = "application/json"
	contentType      = "Content-Type"
	invalidReqMethod = "method not allowed"
	token            = "token"
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

// User API Handlers

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterReq

	user.Name = r.FormValue("name")
	user.LoginReq = &models.LoginReq{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	ctx := context.Background()
	r.Clone(ctx)

	defer ctx.Done()

	resp, err := h.Service.Register(ctx, &user)
	switch {
	case err == nil:
	case errors.Is(err, models.ErrUserAlreadyExists):
		h.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		h.Log.Error("error while registering the user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     token,
		Value:    resp.Token,
		HttpOnly: true,
		Expires:  resp.Expiry,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/task")

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.LoginReq

	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	session, err := h.Service.Login(r.Context(), &user)
	if err != nil {
		if models.ErrNotFound.Is(err) {
			h.Log.Error(err.Error())
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		h.Log.Error("error while logging in the user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     token,
		Value:    session.Token,
		HttpOnly: true,
		Expires:  session.Expiry,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/task")

	w.WriteHeader(http.StatusOK)
	h.Log.Info("login success!!")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(token)
	if err != nil {
		h.Log.Error(err.Error(), "request", "logout")
		http.Error(w, "user not logged in", http.StatusUnauthorized)
		return
	}

	if err := h.Service.Logout(r.Context(), c.Value); err != nil {
		h.Log.Error(err.Error(), "request", "logout")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)

	w.Header().Set(contentType, appJSON)
	w.Header().Add("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)

	h.Log.Info("logout success!!")

	if _, err := w.Write([]byte("logout success !!")); err != nil {
		h.Log.Error("error while writing the response body", "error", err)
	}
}

// TODO API Handlers

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
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	h.Log.Info("Task Done -> ", "ID", id)

	resp, err := h.Service.MarkDone(ctx, id, &userID)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
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
	userID, ok := ctx.Value("user_id").(uuid.UUID)
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

	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	tasks, err := h.Service.GetAll(r.Context(), &userID)
	if err != nil {
		if models.ErrUserNotFound.Is(err) || models.ErrNotFound.Is(err) {
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
	id := r.PathValue("id")
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	h.Log.Info("Delete Request->", "ID", id)

	err := h.Service.DeleteTask(ctx, id, &userID)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
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
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	title := r.Header.Get("HX-Prompt")

	resp, err := h.Service.UpdateTask(ctx, id, title, "false", &userID)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
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
