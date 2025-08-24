package userhttp

import (
	"todoapp/internal/models"

	"gofr.dev/pkg/gofr"
)

// const (
// 	appJSON     = "application/json"
// 	contentType = "Content-Type"
// 	token       = "token"
// 	hxRedirect  = "HX-Redirect"
// 	name        = "name"
// 	email       = "email"
// 	password    = "password"
// )

type AuthAPI interface {
	SignIn(ctx *gofr.Context, email, password string) (*models.AuthUserResp, error)
	SignUp(ctx *gofr.Context, email, password string) error
	Refresh(ctx *gofr.Context) (*string, error)
	Revoke(ctx *gofr.Context) error
}

type Handler struct {
	AuthService AuthAPI
}

func New(authSvc AuthAPI) *Handler {
	return &Handler{AuthService: authSvc}
}

func (h *Handler) Register(ctx *gofr.Context) (any, error) {
	var user models.UserData
	if err := ctx.Bind(&user); err != nil {
		return nil, err
	}

	return nil, h.AuthService.SignUp(ctx, user.Email, user.Password)
}

func (h *Handler) Login(ctx *gofr.Context) (any, error) {
	var user models.UserData
	if err := ctx.Bind(&user); err != nil {
		return nil, err
	}

	resp, err := h.AuthService.SignIn(ctx, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Logout(ctx *gofr.Context) (any, error) {
	return nil, h.AuthService.Revoke(ctx)
}
