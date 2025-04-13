package userhttp

import (
	"context"

	"todoapp/internal/models"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=userhttp
type UserServicer interface {
	Register(ctx context.Context, req *models.RegisterReq) (*models.SessionData, error)
	Login(ctx context.Context, req *models.LoginReq) (*models.SessionData, error)
	Logout(ctx context.Context, token string) error
}
