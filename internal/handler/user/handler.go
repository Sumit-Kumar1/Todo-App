package userhttp

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"todoapp/internal/models"
)

const (
	appJSON     = "application/json"
	contentType = "Content-Type"
	token       = "token"
	hxRedirect  = "HX-Redirect"
	name        = "name"
	email       = "email"
	password    = "password"
)

type Handler struct {
	templ   *template.Template
	Service UserServicer
}

func New(templ *template.Template, usrSvc UserServicer) *Handler {
	return &Handler{templ: templ, Service: usrSvc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(r.Context())

	var user models.RegisterReq

	user.Name = r.FormValue(name)
	user.LoginReq = &models.LoginReq{
		Email:    r.FormValue(email),
		Password: r.FormValue(password),
	}

	defer ctx.Done()

	resp, err := h.Service.Register(ctx, &user)

	switch {
	case err == nil:
	case errors.Is(err, models.ErrUserAlreadyExists):
		logger.LogAttrs(ctx, slog.LevelError, "user already exists, login again",
			slog.String("user", user.Email),
		)

		models.HandleHTTPError(w, err, http.StatusConflict)

		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", user.Email))

		models.HandleHTTPError(w, err, http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	cookie := http.Cookie{
		Name:     token,
		Value:    resp.Token,
		HttpOnly: true,
		Expires:  resp.Expiry,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	w.Header().Add(hxRedirect, "/task")
	w.WriteHeader(http.StatusOK)
	logger.LogAttrs(ctx, slog.LevelDebug, "user logged in successfully!",
		slog.String("user", user.Email),
	)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		user   models.LoginReq
		ctx    = r.Context()
		logger = models.GetLoggerFromCtx(ctx)
	)

	user.Email = r.FormValue(models.Email)
	user.Password = r.FormValue(models.Password)

	session, err := h.Service.Login(ctx, &user)
	if err != nil {
		if models.ErrNotFound("user").Error() == err.Error() {
			logger.LogAttrs(ctx, slog.LevelError, "login-service: user not found",
				slog.String("user", user.Email),
			)

			models.HandleHTTPError(w, err, http.StatusNotFound)
			//http.Error(w, "user not found", http.StatusUnauthorized)

			return
		}

		models.HandleHTTPError(w, err, http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		logger.LogAttrs(ctx, slog.LevelError, "error while logging in the user",
			slog.String("error", err.Error()), slog.String(models.Email, user.Email))

		return
	}

	cookie := http.Cookie{
		Name:     token,
		Value:    session.Token,
		HttpOnly: true,
		Expires:  session.Expiry,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	w.Header().Add(hxRedirect, "/task")
	w.WriteHeader(http.StatusOK)
	logger.LogAttrs(ctx, slog.LevelDebug, "login success", slog.String(models.User, user.Email))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	c, err := r.Cookie(token)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("invalid session:", token))

		models.HandleHTTPError(w, err, http.StatusUnauthorized)
		//http.Error(w, "user not logged in", http.StatusUnauthorized)

		return
	}

	if err2 := h.Service.Logout(ctx, c.Value); err2 != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while logging out user",
			slog.String("error", err2.Error()),
		)

		models.HandleHTTPError(w, err2, http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)

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
	w.Header().Add(hxRedirect, "/")
	w.WriteHeader(http.StatusOK)

	logger.LogAttrs(ctx, slog.LevelDebug, "user logout success!!", slog.String(token, c.Value))
}
