package userhttp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"todoapp/internal/models"
)

const (
	appJSON     = "application/json"
	contentType = "Content-Type"
	token       = "token"
)

type Handler struct {
	Service UserServicer
	Log     *slog.Logger
}

func New(usrSvc UserServicer, logger *slog.Logger) *Handler {
	return &Handler{Service: usrSvc, Log: logger}
}

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
		h.Log.Error(err.Error(), slog.String("error", "User already exist please login"))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		h.Log.Error("error while registering the user", slog.String("error", err.Error()))
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
		if models.ErrNotFound("user").Error() == err.Error() {
			h.Log.Error(err.Error())
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		h.Log.Error("error while logging in the user", slog.String("error", err.Error()), slog.String("email", user.Email))
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
}
