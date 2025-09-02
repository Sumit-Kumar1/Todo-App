package userhttp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"todoapp/internal/errors"
	"todoapp/internal/models"
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "Register: error while closing body",
				slog.String("error", err.Error()))
		}
	}(r.Body)

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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "login:error while closing req body",
				slog.String("error", err.Error()))
		}
	}(r.Body)

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

	data, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
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
