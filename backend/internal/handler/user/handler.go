package userhttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"
)

const (
	authCookieName    = "auth"
	refreshCookieName = "refresh"
)

type Handler struct {
	svc UserServicer
}

func New(svc UserServicer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	var user models.LoginReq

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:error while closing req body",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	if err := h.svc.Register(ctx, &user); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:error while closing req body",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("user created successfully!"))
	logger.LogAttrs(ctx, slog.LevelInfo, "login:user created sucessfully!", slog.String("email", user.Email))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	var user models.LoginReq

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:error while decoding req body",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	resp, err := h.svc.Login(ctx, &user)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:error while calling service",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	authCookie, refCookie := getLoginCookie(resp)
	http.SetCookie(w, &authCookie)
	http.SetCookie(w, &refCookie)

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	if err := h.svc.Logout(ctx, ""); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "logout:error while calling service",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getLoginCookie(resp *models.AuthUserResp) (auth, refresh http.Cookie) {
	if resp == nil {
		return
	}

	auth = http.Cookie{
		Name:     authCookieName,
		Value:    resp.AccessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(time.Minute * 10),
		SameSite: http.SameSiteStrictMode,
	}

	refresh = http.Cookie{
		Name:     refreshCookieName,
		Value:    resp.RefreshToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(time.Hour * 24),
		SameSite: http.SameSiteStrictMode,
	}

	return
}
