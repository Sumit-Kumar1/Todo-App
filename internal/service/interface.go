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
	GetByEmail(ctx context.Context, email string) (*models.UserData, error)
	GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error)
	RefreshSession(ctx context.Context, newSession *models.UserSession) (*models.UserSession, error)
	RegisterUser(ctx context.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error)
}

type TodoStorer interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	Create(ctx context.Context, id, title string) (*models.Task, error)
	Update(ctx context.Context, id, title string) (*models.Task, error)
	Delete(ctx context.Context, id string) error
	MarkDone(ctx context.Context, id string) (*models.Task, error)
}
