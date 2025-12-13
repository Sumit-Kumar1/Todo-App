package usersvc

import (
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Auth AuthClient
}

func New(auth AuthClient) *Service {
	return &Service{Auth: auth}
}

func (s *Service) Register(ctx *gin.Context, req *models.LoginReq) error {
	if req == nil {
		return errors.ErrRequired("user login and password")
	}

	if err := req.Validate(); err != nil {
		return err
	}

	return s.Auth.SignUp(ctx, req.Email, req.Password)
}

func (s *Service) Login(ctx *gin.Context, req *models.LoginReq) (*models.AuthUserResp, error) {
	if req == nil {
		return nil, errors.ErrRequired("user login and password")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	authResp, err := s.Auth.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	if authResp == nil {
		return nil, errors.ErrRequired("login-svc: nil auth response")
	}

	return authResp, nil
}

func (s *Service) Logout(ctx *gin.Context, token string) error {
	if token == "" {
		return errors.ErrRequired("token")
	}

	return s.Auth.Revoke(ctx, token)
}
