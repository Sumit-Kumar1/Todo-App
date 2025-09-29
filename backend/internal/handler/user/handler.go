package userhttp

import (
	"encoding/json"
	pkgErr "errors"
	"log/slog"
	"net/http"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/handler"
	"todoapp/internal/models"
	"todoapp/internal/todocookie"
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
		logger.LogAttrs(ctx, slog.LevelError, "register:error while closing req body",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	if err := h.svc.Register(ctx, &user); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "register:service:error while registering user",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	resp := "user created successfully"

	if err := handler.WriteResponse[string](w, http.StatusCreated, &resp); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "register:error while writing response",
			slog.String("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "login:user created successfully!", slog.String("email", user.Email))
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
		logger.LogAttrs(ctx, slog.LevelError, "login:service:error while calling service",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	authCookie, refCookie := getLoginCookie(resp)
	err = todocookie.WriteCookie(w, authCookie)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:writing cookie auth",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	err = todocookie.WriteCookie(w, refCookie)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:writing cookie refresh",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	success := "user login successfully"

	if err := handler.WriteResponse[string](w, http.StatusOK, &success); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "login:error while writing response",
			slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "user login successfully")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	authCookie, err := todocookie.ReadCookie(r, authCookieName)
	if err != nil {
		if pkgErr.Is(err, http.ErrNoCookie) {
			logger.Log(ctx, slog.LevelError, "logout:service:no cookie in the request")
		}

		logger.Log(ctx, slog.LevelError, "logout:service:error while fetching cookie")
		errors.HandleHTTPError(w, err)
		return
	}

	if err := h.svc.Logout(ctx, authCookie); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "logout:error while calling service",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	_ = todocookie.WriteCookie(w, http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})

	if err := handler.WriteResponse[any](w, int(http.StatusOK), nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "user logout successfully")
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
		SameSite: http.SameSiteLaxMode,
	}

	refresh = http.Cookie{
		Name:     refreshCookieName,
		Value:    resp.RefreshToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   int(time.Hour * 6),
		SameSite: http.SameSiteLaxMode,
	}

	return
}
