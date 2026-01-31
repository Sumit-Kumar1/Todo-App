package todohttp

import (
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=todohttp
type TodoServicer interface {
	GetAll(ctx *gin.Context, userID *uuid.UUID) ([]models.Task, error)
	AddTask(ctx *gin.Context, task *models.TaskReq, userID *uuid.UUID) (*models.Task, error)
	DeleteTask(ctx *gin.Context, id string, userID *uuid.UUID) error
	UpdateTask(ctx *gin.Context, id string, task *models.TaskReq, userID *uuid.UUID) (*models.Task, error)
	MarkDone(ctx *gin.Context, id string, userID *uuid.UUID) (*models.Task, error)
	GetChildTasks(ctx *gin.Context, taskID string, userID *uuid.UUID) ([]models.Task, error)
}
