package usersvc

import (
	"context"

	"todoapp/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=usersvc
type UserStorer interface {
	GetUserByEmail(ctx context.Context, email string) (*models.UserData, error)
	RegisterUser(ctx context.Context, data *models.UserData) error
}

type SessionStorer interface {
	Logout(ctx context.Context, token *uuid.UUID) error
	CreateSession(ctx context.Context, session *models.SessionData) error
	GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.SessionData, error)
	RefreshSession(ctx context.Context, newSession *models.SessionData) error
}
