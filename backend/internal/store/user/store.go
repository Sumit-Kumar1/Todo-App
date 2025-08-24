package userstore

import (
	"database/sql"
	liberrors "errors"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"gofr.dev/pkg/gofr"
)

const (
	getUser       = "SELECT id, name, email, password FROM users WHERE email=$1;"
	registerQuery = "INSERT INTO users(id, name, email, password) VALUES ($1,$2,$3,$4);"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (s *Store) RegisterUser(ctx *gofr.Context, data *models.UserData) error {
	res, err := ctx.SQL.ExecContext(ctx, registerQuery, data.ID, data.Name, data.Email, data.Password)
	if err != nil {
		return err
	}

	if _, err2 := res.RowsAffected(); err2 != nil {
		return err2
	}

	return nil
}

func (s *Store) GetUserByEmail(ctx *gofr.Context, email string) (*models.UserData, error) {
	user := models.UserData{}

	res := ctx.SQL.QueryRowContext(ctx, getUser, email)

	if err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
