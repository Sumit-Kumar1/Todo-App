package userhttp

import (
	"todoapp/internal/models"

	"gofr.dev/pkg/gofr"
)

type Handler struct {
	svc UserServicer
}

func New(svc UserServicer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(ctx *gofr.Context) (any, error) {
	var user models.RegisterReq
	if err := ctx.Bind(&user); err != nil {
		return nil, err
	}

	return nil, h.svc.Register(ctx, &user)
}

func (h *Handler) Login(ctx *gofr.Context) (any, error) {
	var user models.LoginReq
	if err := ctx.Bind(&user); err != nil {
		return nil, err
	}

	authResp, err := h.svc.Login(ctx, &user)
	if err != nil {
		return nil, err
	}

	return authResp, nil
}

func (h *Handler) Logout(ctx *gofr.Context) (any, error) {
	return nil, h.svc.Logout(ctx, "")
}
