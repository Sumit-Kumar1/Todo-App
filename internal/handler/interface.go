package handler

import (
	"context"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

type Servicer interface {
	TodoService
	UserService
}

type UserService interface {
	Register(ctx context.Context, req *models.RegisterReq) (*models.UserSession, error)
	Login(ctx context.Context, req *models.LoginReq) (*models.UserSession, error)
	Logout(ctx context.Context, token string) error
}

type TodoService interface {
	GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error)
	AddTask(ctx context.Context, task string, userID *uuid.UUID) (*models.Task, error)
	DeleteTask(ctx context.Context, id string, userID *uuid.UUID) error
	UpdateTask(ctx context.Context, id, title, isDone string, userID *uuid.UUID) (*models.Task, error)
	MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
