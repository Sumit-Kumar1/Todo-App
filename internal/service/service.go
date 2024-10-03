package service

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Store Storer
	Log   *slog.Logger
}

func New(st Storer, logger *slog.Logger) *Service {
	return &Service{Store: st,
		Log: logger}
}

// User Endpoints

func (s *Service) Register(ctx context.Context, req *models.RegisterReq) (*models.UserSession, error) {
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
	if err != nil {
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

func (s *Service) Login(ctx context.Context, req *models.LoginReq) (*models.UserSession, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get the user's data
	user, err := s.Store.GetByEmail(ctx, req.Email)
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
		if models.ErrNotFound("user ID").Error() == err.Error() {
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

func (s *Service) Logout(ctx context.Context, token string) error {
	t, err := uuid.Parse(token)
	if err != nil {
		return err
	}

	return s.Store.Logout(ctx, &t)
}

// Tasks endpoints

func (s *Service) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	tasks, err := s.Store.GetAll(ctx, userID)
	if err != nil {
		s.Log.Error("error in getAll", "error", err.Error())

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(ctx context.Context, title string, userID *uuid.UUID) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, models.ErrInvalid("task title")
	}

	id := generateID()

	task, err := s.Store.Create(ctx, id, title, userID)
	if err != nil {
		s.Log.Error("error in add task", "error", err.Error())
		return nil, err
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id string, userID *uuid.UUID) error {
	if err := validateID(id); err != nil {
		s.Log.Debug("", "error", err.Error(), "ID", id)
		return models.ErrInvalid("task id")
	}

	if err := s.Store.Delete(ctx, id, userID); err != nil {
		s.Log.Error("error in delete task", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	if err := validateID(id); err != nil {
		s.Log.ErrorContext(ctx, err.Error(), "task id", id)
		return nil, err
	}

	task, err := s.Store.MarkDone(ctx, id, userID)
	if err != nil {
		s.Log.Error("error in mark-done", "error", err.Error())

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, id, title, isDone string, userID *uuid.UUID) (*models.Task, error) {
	if err := validateTask(id, title, isDone); err != nil {
		s.Log.ErrorContext(ctx, err.Error(), "ID", id)
		return nil, err
	}

	task, err := s.Store.Update(ctx, id, title, userID)
	if err != nil {
		s.Log.Error("error in updating task", "error", err.Error())

		return nil, err
	}

	s.Log.Info("\nUpdated task", "task", task)

	return task, nil
}
