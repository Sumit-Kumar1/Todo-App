package usersvc

import (
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen --destination=mock_interface.go --package=usersvc
type AuthClient interface {
	SignUp(ctx *gin.Context, email, password string) error
	SignIn(ctx *gin.Context, email, password string) (*models.AuthUserResp, error)
	Refresh(ctx *gin.Context, token string) (*string, error)
	Revoke(ctx *gin.Context, token string) error
}
