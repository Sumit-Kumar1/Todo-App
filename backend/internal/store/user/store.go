package userstore

import (
	"context"
	"database/sql"
	liberrors "errors"

	"todoapp/internal/errors"
	"todoapp/internal/models"
)

const (
	getUser       = "SELECT id, name, email, password FROM users WHERE email=$1;"
	registerQuery = "INSERT INTO users(id, name, email, password) VALUES ($1,$2,$3,$4);"
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
	res, err := s.DB.ExecContext(ctx, registerQuery, data.ID, data.Name, data.Email, data.Password)
	if err != nil {
		return err
	}

	if _, err2 := res.RowsAffected(); err2 != nil {
		return err2
	}

	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	user := models.UserData{}

	res := s.DB.QueryRowContext(ctx, getUser, email)

	if err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
