package service

import (
	"context"
	"todoapp/internal/models"
)

type Storer interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	Create(ctx context.Context, id, title string) (*models.Task, error)
	Update(ctx context.Context, id, title string) (*models.Task, error)
	Delete(ctx context.Context, id string) error
	MarkDone(ctx context.Context, id string) (*models.Task, error)
}
