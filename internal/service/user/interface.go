package usersvc

import (
	"context"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=usersvc
type UserStorer interface {
	Logout(ctx context.Context, token *uuid.UUID) error
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetUserByEmail(ctx context.Context, email string) (*models.UserData, error)
	GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error)
	RefreshSession(ctx context.Context, newSession *models.UserSession) error
	RegisterUser(ctx context.Context, data *models.UserData) error
}
