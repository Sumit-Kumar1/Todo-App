package userstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	createSession      = "INSERT INTO sessions (id, user_id, token, expiry) VALUES (?, ?, ?,?);"
	deleteSessionByID  = "DELETE FROM sessions WHERE id=?"
	getUser            = "SELECT * FROM users WHERE email = ?;"
	getSessionByUserID = "SELECT * FROM sessions WHERE user_id = ?;"
	getSessionByToken  = "SELECT * FROM sessions where token=?"
	registerQuery      = "INSERT INTO users(user_id, name, email, password) VALUES (?,?,?,?);"
	updateSession      = "UPDATE sessions SET token = ?,  expiry = ? WHERE id = ?;"
)

type Store struct {
	DB  *sql.DB
	Log *slog.Logger
}

func New(db *sql.DB, logger *slog.Logger) *Store {
	return &Store{
		DB:  db,
		Log: logger,
	}
}

func (s *Store) RegisterUser(ctx context.Context, data *models.UserData) error {
	vals := []any{data.ID, data.Name, data.Email, data.Password}
	sr, err := s.DB.ExecContext(ctx, registerQuery, vals...)
	if err != nil {
		return err
	}

	if _, err := sr.LastInsertId(); err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateSession(ctx context.Context, session *models.UserSession) error {
	vals := []any{session.ID, session.UserID, session.Token, session.Expiry}

	rs, err := s.DB.ExecContext(ctx, createSession, vals...)
	if err != nil {
		return err
	}

	if _, err := rs.RowsAffected(); err != nil {
		return errors.New("no rows were affected")
	}

	return nil
}

func (s *Store) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error) {
	var session models.UserSession

	row := s.DB.QueryRowContext(ctx, getSessionByUserID, *userID)
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound("user")
		}

		return nil, err
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.UserSession) error {
	var vals = []any{newSession.Token, newSession.Expiry, newSession.ID}

	sr, err := s.DB.ExecContext(ctx, updateSession, vals...)
	if err != nil {
		return err
	}

	if _, err := sr.RowsAffected(); err != nil {
		return err
	}

	s.Log.Info("session is refreshed", "user", newSession.UserID)

	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	var user models.UserData

	row := s.DB.QueryRowContext(ctx, getUser, email)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound("user")
		}

		return nil, err
	}

	return &user, nil
}

func (s *Store) Logout(ctx context.Context, token *uuid.UUID) error {
	var session models.UserSession

	row := s.DB.QueryRowContext(ctx, getSessionByToken, *token)
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound("session with current user")
		}
	}

	rs, err := s.DB.ExecContext(ctx, deleteSessionByID, session.ID)
	if err != nil {
		return err
	}

	if _, err := rs.RowsAffected(); err != nil {
		return errors.New("no rows were affected")
	}

	return nil
}
