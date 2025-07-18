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
	UserStore    UserStorer
	SessionStore SessionStorer
}

func New(st UserStorer, ss SessionStorer) *Service {
	return &Service{UserStore: st, SessionStore: ss}
}

func (s *Service) Register(ctx context.Context, req *models.RegisterReq) (*models.SessionData, error) {
	if req == nil {
		return nil, nil
	}

	logger := models.GetLoggerFromCtx(ctx)

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// check if user already exists
	existingUser, err := s.UserStore.GetUserByEmail(ctx, req.Email)
	if err != nil && err.Error() != models.ErrNotFound("user").Error() {
		logger.LogAttrs(ctx, slog.LevelError, "Service.Register - user not found",
			slog.String("error", err.Error()),
			slog.String("user", req.Email),
		)

		return nil, err
	}

	if existingUser != nil {
		return nil, models.ErrUserAlreadyExists
	}

	passwd, err := encryptedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.UserData{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: passwd,
	}

	if err := s.UserStore.RegisterUser(ctx, &user); err != nil {
		return nil, err
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "user created successfully!!",
		slog.String("email", req.Email), slog.String("userID", user.ID.String()))

	session := models.SessionData{
		ID:     uuid.New(),
		UserID: user.ID,
		Token:  uuid.NewString(),
		Expiry: time.Now().Add(time.Minute * 15),
	}

	if err := s.SessionStore.CreateSession(ctx, &session); err != nil {
		return nil, err
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "Service:Register - session created successfully!!",
		slog.String("userID", user.ID.String()),
	)

	return &session, nil
}

func (s *Service) Login(ctx context.Context, req *models.LoginReq) (*models.SessionData, error) {
	if req == nil {
		return nil, models.ErrRequired("login request")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get the user's data
	user, err := s.UserStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, models.ErrUserNotFound
	}

	if matchErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); matchErr != nil {
		return nil, models.ErrPsswdNotMatch
	}

	return s.handleLoginSession(ctx, user)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	t, err := uuid.Parse(token)
	if err != nil {
		return err
	}

	return s.SessionStore.Logout(ctx, &t)
}

func (s *Service) handleLoginSession(ctx context.Context, user *models.UserData) (*models.SessionData, error) {
	session, err := s.SessionStore.GetSessionByID(ctx, &user.ID)
	if err != nil {
		if models.ErrNotFound("user ID").Error() != err.Error() {
			return nil, err
		}

		t := time.Now().Add(time.Minute * 15).UTC()
		ss := models.SessionData{
			ID:     uuid.New(),
			UserID: user.ID,
			Token:  uuid.NewString(),
			Expiry: t,
		}

		if er := s.SessionStore.CreateSession(ctx, &ss); er != nil {
			return nil, er
		}

		return &ss, nil
	}

	if session.Expiry.Before(time.Now().UTC()) {
		session.Expiry = time.Now().Add(time.Minute * 15).UTC()
		session.Token = uuid.NewString()

		if err := s.SessionStore.RefreshSession(ctx, session); err != nil {
			return nil, err
		}
	}

	return session, nil
}

func encryptedPassword(password string) (string, error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwd), nil
}
