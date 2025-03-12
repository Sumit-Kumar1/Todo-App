package userhttp

import (
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
}

func New(usrSvc UserServicer) *Handler {
	return &Handler{Service: usrSvc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(r.Context())

	var user models.RegisterReq

	user.Name = r.FormValue("name")
	user.LoginReq = &models.LoginReq{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	defer ctx.Done()

	resp, err := h.Service.Register(ctx, &user)
	switch {
	case err == nil:
	case errors.Is(err, models.ErrUserAlreadyExists):
		logger.LogAttrs(ctx, slog.LevelError, "user already exists, login again", slog.String("user", user.Email))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", user.Email))
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
	logger.LogAttrs(ctx, slog.LevelDebug, "user logged in successfully!", slog.String("user", user.Email))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		user   models.LoginReq
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	session, err := h.Service.Login(ctx, &user)
	if err != nil {
		if models.ErrNotFound("user").Error() == err.Error() {
			logger.LogAttrs(ctx, slog.LevelError, "user not found - login", slog.String("user", user.Email))
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		logger.LogAttrs(ctx, slog.LevelError, "error while logging in the user",
			slog.String("error", err.Error()), slog.String("email", user.Email))
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
	logger.LogAttrs(ctx, slog.LevelDebug, "login success", slog.String("user", user.Email))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	c, err := r.Cookie(token)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("invalid session:", token))
		http.Error(w, "user not logged in", http.StatusUnauthorized)
		return
	}

	if err := h.Service.Logout(ctx, c.Value); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while logging out user", slog.String("error", err.Error()))
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

	logger.LogAttrs(ctx, slog.LevelDebug, "user logout success!", slog.String("token", c.Value))
}
