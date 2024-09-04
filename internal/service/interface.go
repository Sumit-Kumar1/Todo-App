package service

import (
	"context"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

type Storer interface {
	UserStorer
	TodoStorer
}

type UserStorer interface {
	Logout(ctx context.Context, token *uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*models.UserData, error)
	GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error)
	RefreshSession(ctx context.Context, newSession *models.UserSession) (*models.UserSession, error)
	RegisterUser(ctx context.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error)
}

type TodoStorer interface {
	GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error)
	Create(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error)
	Update(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error)
	Delete(ctx context.Context, id string, userID *uuid.UUID) error
	MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
