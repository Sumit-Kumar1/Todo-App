package handler

import (
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetContextKey(c *gin.Context) (*uuid.UUID, error) {
	ctxVal, ok := c.Get(string(models.CtxKeyUserID))
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	userID, ok := ctxVal.(uuid.UUID)
	if !ok {
		return nil, errors.Invalid("login user")
	}

	return &userID, nil
}
