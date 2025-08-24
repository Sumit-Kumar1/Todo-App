package usersvc

import (
	"todoapp/internal/models"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=usersvc
type UserStorer interface {
	GetUserByEmail(ctx *gofr.Context, email string) (*models.UserData, error)
	RegisterUser(ctx *gofr.Context, data *models.UserData) error
}

type SessionStorer interface {
	Logout(ctx *gofr.Context, token *uuid.UUID) error
	CreateSession(ctx *gofr.Context, session *models.SessionData) error
	GetSessionByID(ctx *gofr.Context, userID *uuid.UUID) (*models.SessionData, error)
	RefreshSession(ctx *gofr.Context, newSession *models.SessionData) error
}
