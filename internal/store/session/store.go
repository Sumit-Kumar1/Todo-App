package sessionstore

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/sqlitecloud/sqlitecloud-go"
)

const (
	createSession      = "INSERT INTO sessions (id, user_id, token, expiry) VALUES ('%v', '%v', '%v','%v');"
	deleteSessionByID  = "DELETE FROM sessions WHERE id='%v';"
	getSessionByUserID = "SELECT id, user_id, token, expiry FROM sessions WHERE user_id='%v';"
	//nolint:gosec //not any hardcoded credential
	getSessionByToken = "SELECT id FROM sessions where token='%v';"
	updateSession     = "UPDATE sessions SET token='%v',  expiry='%v' WHERE id='%v';"
)

type Store struct {
	DB *sqlitecloud.SQCloud
}

func New(db *sqlitecloud.SQCloud) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateSession(ctx context.Context, session *models.SessionData) error {
	logger := models.GetLoggerFromCtx(ctx)

	query := fmt.Sprintf(
		createSession,
		session.ID,
		session.UserID,
		session.Token,
		session.Expiry.UnixMilli(),
	)
	if err := s.DB.Execute(query); err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"error while running session create query",
			slog.String("error", err.Error()),
		)

		return err
	}

	return nil
}

// nolint:gocognit // cognitive of 11 is okay, didn't want to split this method
func (s *Store) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.SessionData, error) {
	logger := models.GetLoggerFromCtx(ctx)

	var session models.SessionData

	res, err := s.DB.Select(fmt.Sprintf(getSessionByUserID, *userID))
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError,
			"error while fetching session by userID",
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	if res.GetNumberOfRows() == uint64(0) {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"no user session found for userID",
			slog.String("userID", userID.String()),
		)

		return nil, models.ErrNotFound("user ID")
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

		c4, err := res.GetInt64Value(r, 3)
		if err != nil {
			return nil, err
		}

		session.ID = uuid.MustParse(c1)
		session.UserID = uuid.MustParse(c2)
		session.Token = c3

		session.Expiry = time.UnixMilli(c4)
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.SessionData) error {
	logger := models.GetLoggerFromCtx(ctx)

	query := fmt.Sprintf(
		updateSession,
		newSession.Token,
		newSession.Expiry.UnixMilli(),
		newSession.ID,
	)
	if err := s.DB.Execute(query); err != nil {
		logger.LogAttrs(ctx, slog.LevelError,
			"error in refreshing session",
			slog.String("error", err.Error()),
		)

		return err
	}

	return nil
}

func (s *Store) Logout(ctx context.Context, token *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	var id uuid.UUID

	res, err := s.DB.Select(fmt.Sprintf(getSessionByToken, *token))
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError,
			"error while logging out user",
			slog.String("error", err.Error()),
		)

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

	return s.DB.Execute(fmt.Sprintf(deleteSessionByID, id))
}
