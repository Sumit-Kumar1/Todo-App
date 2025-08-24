package userhttp

import (
	"todoapp/internal/models"

	"gofr.dev/pkg/gofr"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=userhttp
type UserServicer interface {
	Register(ctx *gofr.Context, req *models.RegisterReq) (*models.SessionData, error)
	Login(ctx *gofr.Context, req *models.LoginReq) (*models.SessionData, error)
	Logout(ctx *gofr.Context, token string) error
}
