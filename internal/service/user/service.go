package usersvc

import (
	"context"
	"log/slog"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Store UserStorer
}

func New(st UserStorer) *Service {
	return &Service{Store: st}
}

func (s *Service) Register(ctx context.Context, req *models.RegisterReq) (*models.UserSession, error) {
	logger := models.GetLoggerFromCtx(ctx)

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
	existingUser, err := s.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() != models.ErrNotFound("user").Error() {
			logger.LogAttrs(ctx, slog.LevelError, "user not found - Service.Register", slog.String("error", err.Error()),
				slog.String("user", req.Email))
			return nil, err
		}
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
		Expiry: time.Now().Add(time.Minute * 15),
	}

	user := models.UserData{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: passwd,
	}

	err = s.Store.RegisterUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "user created successfully!!",
		slog.String("email", req.Email), slog.String("userID", userID.String()))

	err = s.Store.CreateSession(ctx, &session)
	if err != nil {
		return nil, err
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "session created successfully!!", slog.String("userID", userID.String()))

	return &session, nil
}

func (s *Service) Login(ctx context.Context, req *models.LoginReq) (*models.UserSession, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get the user's data
	user, err := s.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, models.ErrNotFound("user")
	}

	if matchErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); matchErr != nil {
		return nil, models.ErrPsswdNotMatch
	}

	session, err := s.Store.GetSessionByID(ctx, &user.ID)
	if err != nil {
		if models.ErrNotFound("user ID").Error() != err.Error() {
			return nil, err
		}

		t := time.Now().Add(time.Minute * 15)
		ss := models.UserSession{
			ID:     uuid.New(),
			UserID: user.ID,
			Token:  uuid.NewString(),
			Expiry: t,
		}

		if er := s.Store.CreateSession(ctx, &ss); er != nil {
			return nil, er
		}

		return &ss, nil
	}

	if session.Expiry.Before(time.Now().UTC()) {
		session.Expiry = time.Now().Add(time.Minute * 15).UTC()
		session.Token = uuid.NewString()

		err := s.Store.RefreshSession(ctx, session)
		if err != nil {
			return nil, err
		}
	}

	return session, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	t, err := uuid.Parse(token)
	if err != nil {
		return err
	}

	return s.Store.Logout(ctx, &t)
}
