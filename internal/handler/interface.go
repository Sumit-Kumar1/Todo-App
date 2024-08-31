package handler

import (
	"context"
	"todoapp/internal/models"
)

type Servicer interface {
	Register(ctx context.Context, req models.RegisterReq) (models.LoginSession, error)
	Login(ctx context.Context, req models.LoginReq) (models.LoginSession, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	AddTask(ctx context.Context, task string) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
	UpdateTask(ctx context.Context, id, title, isDone string) (*models.Task, error)
	MarkDone(ctx context.Context, id string) (*models.Task, error)
}
