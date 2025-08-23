package sessionstore

import (
	"context"
	"database/sql"
	liberrors "errors"
	"log/slog"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	createSession      = "INSERT INTO sessions (id, user_id, token, expiry) VALUES ($1, $2, $3,$4);"
	deleteSessionByID  = "DELETE FROM sessions WHERE id=$1;"
	getSessionByUserID = "SELECT id, user_id, token, expiry FROM sessions WHERE user_id=$1;"
	//nolint:gosec //not any hardcoded credential
	getSessionByToken = "SELECT id FROM sessions where token=$1;"
	updateSession     = "UPDATE sessions SET token=$1, expiry=$2 WHERE id=$3;"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateSession(ctx context.Context, session *models.SessionData) error {
	res, err := s.DB.ExecContext(ctx, createSession, session.ID, session.UserID, session.Token, session.Expiry)
	if err != nil {
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

	row := s.DB.QueryRowContext(ctx, getSessionByUserID, *userID)
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
			logger.LogAttrs(ctx, slog.LevelError, "store : no user session found for userID",
				slog.String("userID", userID.String()),
			)

			return nil, errors.ErrNotFound("user ID")
		}

		return nil, err
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.SessionData) error {
	logger := models.GetLoggerFromCtx(ctx)

	_, err := s.DB.ExecContext(ctx, updateSession, newSession.Token, newSession.Expiry, newSession.ID)
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
		if liberrors.Is(err, sql.ErrNoRows) {
			return errors.ErrNotFound("session with current user")
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
