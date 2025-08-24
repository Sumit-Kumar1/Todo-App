package usersvc

import (
	"todoapp/internal/models"

	"gofr.dev/pkg/gofr"
)

//go:generate mockgen --destination=mock_interface.go --package=usersvc
type AuthClient interface {
	SignUp(ctx *gofr.Context, email, password string) error
	SignIn(ctx *gofr.Context, email, password string) (*models.AuthUserResp, error)
	Refresh(ctx *gofr.Context) (*string, error)
	Revoke(ctx *gofr.Context, token string) error
}
