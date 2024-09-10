package service

import (
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

type Storer interface {
	UserStorer
	TodoStorer
}

type UserStorer interface {
	Logout(ctx server.Context, token *uuid.UUID) error
	CreateSession(ctx server.Context, session *models.UserSession) error
	GetByEmail(ctx server.Context, email string) (*models.UserData, error)
	GetSessionByID(ctx server.Context, userID *uuid.UUID) (*models.UserSession, error)
	RefreshSession(ctx server.Context, newSession *models.UserSession) (*models.UserSession, error)
	RegisterUser(ctx server.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error)
}

type TodoStorer interface {
	GetAll(ctx server.Context, userID *uuid.UUID) ([]models.Task, error)
	Create(ctx server.Context, id, title string, userID *uuid.UUID) (*models.Task, error)
	Update(ctx server.Context, id, title string, userID *uuid.UUID) (*models.Task, error)
	Delete(ctx server.Context, id string, userID *uuid.UUID) error
	MarkDone(ctx server.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
