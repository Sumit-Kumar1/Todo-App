package userhttp

import (
	"todoapp/internal/models"
	"todoapp/internal/server"
)

type UserServicer interface {
	Register(ctx server.Context, req *models.RegisterReq) (*models.UserSession, error)
	Login(ctx server.Context, req *models.LoginReq) (*models.UserSession, error)
	Logout(ctx server.Context, token string) error
}
