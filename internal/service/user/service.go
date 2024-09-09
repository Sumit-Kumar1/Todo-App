package userservice

import (
	"time"
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Store UserStorer
}

func New(us UserStorer) *Service {
	return &Service{
		Store: us,
	}
}

func (s *Service) Register(ctx server.Context, req *models.RegisterReq) (*models.UserSession, error) {
	if req == nil {
		return nil, nil
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	passwd, err := encryptedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// check if user already exists
	existingUser, err := s.Store.GetByEmail(ctx, req.Email)
	if err != nil && !models.ErrUserNotFound.Is(err) {
		return nil, err
	}

	if existingUser != nil {
		return nil, models.ErrUserAlreadyExists
	}

	userID := uuid.New()
	sessionID := uuid.New()

	session := models.UserSession{
		ID:     sessionID,
		UserID: userID,
		Token:  uuid.NewString(),
		Expiry: time.Now().Add(time.Minute * 15).UTC(),
	}

	user := models.UserData{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: passwd,
	}

	return s.Store.RegisterUser(ctx, &user, &session)
}

func (s *Service) Login(ctx server.Context, req *models.LoginReq) (*models.UserSession, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get the user's data
	user, err := s.Store.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, models.ErrUserNotFound
	}

	if matchErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); matchErr != nil {
		return nil, models.ErrPsswdNotMatch
	}

	session, err := s.Store.GetSessionByID(ctx, &user.ID)
	if err != nil {
		if !models.ErrNotFound.Is(err) {
			return nil, err
		}

		ss := models.UserSession{
			ID:     uuid.New(),
			UserID: user.ID,
			Token:  uuid.NewString(),
			Expiry: time.Now().Add(time.Minute * 15).UTC(),
		}

		if er := s.Store.CreateSession(ctx, &ss); er != nil {
			return nil, er
		}

		return &ss, nil
	}

	if session.Expiry.Before(time.Now().UTC()) {
		session.Expiry = time.Now().Add(time.Minute * 15).UTC()
		session.Token = uuid.NewString()

		return s.Store.RefreshSession(ctx, session)
	}

	return session, nil
}

func (s *Service) Logout(ctx server.Context, token string) error {
	t, err := uuid.Parse(token)
	if err != nil {
		return err
	}

	return s.Store.Logout(ctx, &t)
}
