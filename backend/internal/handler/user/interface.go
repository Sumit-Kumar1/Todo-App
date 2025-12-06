package userhttp

import (
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=userhttp
type UserServicer interface {
	Register(ctx *gin.Context, req *models.LoginReq) error
	Login(ctx *gin.Context, req *models.LoginReq) (*models.AuthUserResp, error)
	Logout(ctx *gin.Context, token string) error
}
