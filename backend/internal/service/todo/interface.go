package todosvc

import (
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=todosvc
type TodoStorer interface {
	GetAll(ctx *gin.Context, userID *uuid.UUID) ([]models.Task, error)
	Create(ctx *gin.Context, task *models.Task) error
	Update(ctx *gin.Context, task *models.Task) error
	Delete(ctx *gin.Context, id string, userID *uuid.UUID) error
	MarkDone(ctx *gin.Context, id string, userID *uuid.UUID) error
	GetTaskByID(ctx *gin.Context, taskID string, userID *uuid.UUID) (*models.Task, error)
}
