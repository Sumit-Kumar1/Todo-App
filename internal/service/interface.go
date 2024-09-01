package service

import (
	"context"
	"todoapp/internal/models"
)

type Storer interface {
	UserStorer
	TodoStorer
}

type UserStorer interface {
	GetByEmail(ctx context.Context, email string) (*models.UserData, error)
	RegisterUser(ctx context.Context, data *models.UserData) (*models.LoginSession, error)
}

type TodoStorer interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	Create(ctx context.Context, id, title string) (*models.Task, error)
	Update(ctx context.Context, id, title string) (*models.Task, error)
	Delete(ctx context.Context, id string) error
	MarkDone(ctx context.Context, id string) (*models.Task, error)
}
