package sessionstore

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
	deleteSessionByID  = "DELETE FROM sessions WHERE id=?;"
	getSessionByUserID = "SELECT id, user_id, token, expiry FROM sessions WHERE user_id=?;"
	//nolint:gosec //not any hardcoded credential
	getSessionByToken = "SELECT id FROM sessions where token=?;"
	updateSession     = "UPDATE sessions SET token=?, expiry=? WHERE id=?;"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateSession(ctx context.Context, session *models.SessionData) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, createSession, session.ID,
		session.UserID, session.Token, session.Expiry.UnixMilli())
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while running session create query",
			slog.String("error", err.Error()),
		)

		return err
	}

	if _, err = res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.SessionData, error) {
	logger := models.GetLoggerFromCtx(ctx)

	var session models.SessionData

	res, err := s.DB.QueryContext(ctx, getSessionByUserID, *userID)
	if err == nil {
		for res.Next() {
			if err2 := res.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err2 != nil {
				return nil, err2
			}
		}

		return &session, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		logger.LogAttrs(ctx, slog.LevelError, "store : no user session found for userID",
			slog.String("userID", userID.String()),
		)

		return nil, models.ErrNotFound("user ID")
	}

	logger.LogAttrs(ctx, slog.LevelError, "store : error while fetching session by userID",
		slog.String("error", err.Error()),
	)

	return nil, err
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.SessionData) error {
	logger := models.GetLoggerFromCtx(ctx)

	_, err := s.DB.ExecContext(ctx, updateSession, newSession.Token,
		newSession.Expiry.UnixMilli(), newSession.ID)
	if err == nil {
		return nil
	}

	logger.LogAttrs(ctx, slog.LevelError, "store : error while refreshing session",
		slog.String("error", err.Error()),
	)

	return err
}

func (s *Store) Logout(ctx context.Context, token *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	var (
		id uuid.UUID
		r1 string
	)

	res := s.DB.QueryRowContext(ctx, getSessionByToken, *token)
	if err := res.Scan(&r1); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound("session with current user")
		}

		logger.LogAttrs(ctx, slog.LevelError, "error while logging out user",
			slog.String("error", err.Error()),
		)

		return err
	}

	id = uuid.MustParse(r1)

	_, err := s.DB.ExecContext(ctx, deleteSessionByID, id)

	return err
}
