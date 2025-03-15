package userstore

import (
	"context"
	"fmt"
	"log/slog"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/sqlitecloud/sqlitecloud-go"
)

const (
	getUser       = "SELECT user_id, name, email, password FROM users WHERE email='%s';"
	registerQuery = "INSERT INTO users(user_id, name, email, password) VALUES ('%v','%v','%v','%v');"
)

type Store struct {
	DB *sqlitecloud.SQCloud
}

func New(db *sqlitecloud.SQCloud) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) RegisterUser(ctx context.Context, data *models.UserData) error {
	logger := models.GetLoggerFromCtx(ctx)

	query := fmt.Sprintf(registerQuery, data.ID, data.Name, data.Email, data.Password)
	if err := s.DB.Execute(query); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while running Register query", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	logger := models.GetLoggerFromCtx(ctx)

	var user models.UserData

	res, err := s.DB.Select(fmt.Sprintf(getUser, email))
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error in fetching user by email", slog.String("error", err.Error()))
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
