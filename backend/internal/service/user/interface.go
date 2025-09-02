package usersvc

import (
	"context"
	"todoapp/internal/models"
)

//go:generate mockgen --destination=mock_interface.go --package=usersvc
type AuthClient interface {
	SignUp(ctx context.Context, email, password string) error
	SignIn(ctx context.Context, email, password string) (*models.AuthUserResp, error)
	Refresh(ctx context.Context, token string) (*string, error)
	Revoke(ctx context.Context, token string) error
}
