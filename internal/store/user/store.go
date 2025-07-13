package userstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"todoapp/internal/models"
)

const (
	getUser       = "SELECT id, name, email, password FROM users WHERE email=?;"
	registerQuery = "INSERT INTO users(id, name, email, password) VALUES (?,?,?,?);"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) RegisterUser(ctx context.Context, data *models.UserData) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, registerQuery, data.ID, data.Name, data.Email, data.Password)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while running Register query",
			slog.String("error", err.Error()),
		)

		return err
	}

	if _, err2 := res.LastInsertId(); err2 != nil {
		return err2
	}

	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	logger := models.GetLoggerFromCtx(ctx)
	user := models.UserData{}

	res := s.DB.QueryRowContext(ctx, getUser, email)

	if err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		logger.LogAttrs(ctx, slog.LevelError, "error in fetching user by email",
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	return &user, nil
}
