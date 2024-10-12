package userhttp

import (
	"context"
	"todoapp/internal/models"
)

type UserServicer interface {
	Register(ctx context.Context, req *models.RegisterReq) (*models.UserSession, error)
	Login(ctx context.Context, req *models.LoginReq) (*models.UserSession, error)
	Logout(ctx context.Context, token string) error
}
