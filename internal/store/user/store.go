package userstore

import (
	"context"
	"fmt"
	"github.com/sqlitecloud/sqlitecloud-go"
	"log/slog"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	createSession      = "INSERT INTO sessions (id, user_id, token, expiry) VALUES (%v, %v, %v,%v);"
	deleteSessionByID  = "DELETE FROM sessions WHERE id=%v"
	getUser            = "SELECT * FROM users WHERE email = %v;"
	getSessionByUserID = "SELECT * FROM sessions WHERE user_id = %v;"
	//nolint:gosec //not any hardcoded credential
	getSessionByToken = "SELECT id FROM sessions where token=%v"
	registerQuery     = "INSERT INTO users(user_id, name, email, password) VALUES (%v,%v,%v,%v);"
	updateSession     = "UPDATE sessions SET token = %v,  expiry = %v WHERE id = %v;"
)

type Store struct {
	DB  *sqlitecloud.SQCloud
	Log *slog.Logger
}

func New(db *sqlitecloud.SQCloud, logger *slog.Logger) *Store {
	return &Store{
		DB:  db,
		Log: logger,
	}
}

func (s *Store) RegisterUser(ctx context.Context, data *models.UserData) error {
	if err := s.DB.Execute(fmt.Sprintf(registerQuery, data.ID, data.Name, data.Email, data.Password)); err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateSession(ctx context.Context, session *models.UserSession) error {
	if err := s.DB.Execute(fmt.Sprintf(createSession, session.ID, session.UserID, session.Token, session.Expiry)); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error) {
	var session models.UserSession

	res, err := s.DB.Select(fmt.Sprintf(getSessionByUserID, *userID))
	if err != nil {
		return nil, err
	}

	if res.GetNumberOfRows() == uint64(0) {
		return nil, models.ErrNotFound("user")
	}

	for r := uint64(0); r < res.GetNumberOfRows(); r++ {
		c1, err := res.GetStringValue(r, 0)
		if err != nil {
			return nil, err
		}
		c2, err := res.GetStringValue(r, 1)
		if err != nil {
			return nil, err
		}
		c3, err := res.GetStringValue(r, 2)
		if err != nil {
			return nil, err
		}
		c4, err := res.GetFloat64Value(r, 3)
		if err != nil {
			return nil, err
		}

		session.ID = uuid.MustParse(c1)
		session.UserID = uuid.MustParse(c2)
		session.Token = c3
		session.Expiry = time.Unix(int64(c4), 0)
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.UserSession) error {
	if err := s.DB.Execute(fmt.Sprintf(updateSession, newSession.Token, newSession.Expiry, newSession.ID)); err != nil {
		return err
	}

	s.Log.Info("session is refreshed", "user", newSession.UserID)

	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	var user models.UserData

	res, err := s.DB.Select(fmt.Sprintf(getUser, email))
	if err != nil {
		return nil, err
	}

	if res.GetNumberOfRows() == 0 {
		return nil, models.ErrNotFound("user")
	}

	for r := uint64(0); r < res.GetNumberOfRows(); r++ {
		c1, err := res.GetStringValue(r, 0)
		if err != nil {
			return nil, err
		}
		c2, err := res.GetStringValue(r, 1)
		if err != nil {
			return nil, err
		}
		c3, err := res.GetStringValue(r, 2)
		if err != nil {
			return nil, err
		}

		c4, err := res.GetStringValue(r, 3)
		if err != nil {
			return nil, err
		}

		user.ID = uuid.MustParse(c1)
		user.Name = c2
		user.Email = c3
		user.Password = c4
	}

	return &user, nil
}

func (s *Store) Logout(ctx context.Context, token *uuid.UUID) error {
	var id uuid.UUID

	res, err := s.DB.Select(fmt.Sprintf(getSessionByToken, *token))
	if err != nil {
		return err
	}

	if res.GetNumberOfRows() == 0 {
		return models.ErrNotFound("session with current user")
	}

	for r := uint64(0); r < res.GetNumberOfRows(); r++ {
		r1, err := res.GetStringValue(r, 0)
		if err != nil {
			return err
		}

		id = uuid.MustParse(r1)
	}

	if err := s.DB.Execute(fmt.Sprintf(deleteSessionByID, id)); err != nil {
		return err
	}

	return nil
}
