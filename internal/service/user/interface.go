package userservice

import (
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

type UserStorer interface {
	Logout(ctx server.Context, token *uuid.UUID) error
	CreateSession(ctx server.Context, session *models.UserSession) error
	GetByEmail(ctx server.Context, email string) (*models.UserData, error)
	GetSessionByID(ctx server.Context, userID *uuid.UUID) (*models.UserSession, error)
	RefreshSession(ctx server.Context, newSession *models.UserSession) (*models.UserSession, error)
	RegisterUser(ctx server.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error)
}
