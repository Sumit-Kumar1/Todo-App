package userstore

import (
	"database/sql"
	"errors"
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

const (
	registerQuery  = "INSERT INTO users(user_id, name, email, password) VALUES (?,?,?,?);"
	getUser        = "SELECT * FROM users WHERE email = ?;"
	createSession  = "INSERT INTO sessions (id, user_id, token, expiry) VALUES (?, ?, ?,?);"
	getSession     = "SELECT * FROM sessions WHERE user_id = ?;"
	refreshSession = "UPDATE sessions SET token = ?,  expiry = ? WHERE id = ?;"
)

func (s *Store) RegisterUser(ctx server.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error) {
	var (
		err  error
		opts = sql.TxOptions{Isolation: sql.LevelReadCommitted}
	)

	tx, err := ctx.DB.BeginTx(ctx, &opts)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, registerQuery, data.ID, data.Name, data.Email, data.Password)
	if err != nil {
		return nil, rollback(tx, err)
	}

	_, err = tx.ExecContext(ctx, createSession, session.ID, session.UserID, session.Token, session.Expiry)
	if err != nil {
		return nil, rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, rollback(tx, err)
	}

	return session, nil
}

func (s *Store) CreateSession(ctx server.Context, session *models.UserSession) error {
	_, err := ctx.DB.ExecContext(ctx, createSession, session.ID, session.UserID, session.Token, session.Expiry)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSessionByID(ctx server.Context, userID *uuid.UUID) (*models.UserSession, error) {
	var session models.UserSession

	if userID == nil {
		return nil, errors.New("invalid user_id provided")
	}

	row := ctx.DB.QueryRowContext(ctx, getSession, *userID)
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx server.Context, newSession *models.UserSession) (*models.UserSession, error) {
	_, err := ctx.DB.ExecContext(ctx, refreshSession, newSession.Token, newSession.Expiry, newSession.ID)
	if err != nil {
		return nil, err
	}

	ctx.Logger.Info("session is refreshed", "userID", newSession.UserID)

	return newSession, nil
}

func (s *Store) GetByEmail(ctx server.Context, email string) (*models.UserData, error) {
	var user models.UserData

	row := ctx.DB.QueryRowContext(ctx, getUser, email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (s *Store) Logout(ctx server.Context, token *uuid.UUID) error {
	tx, err := ctx.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	var session models.UserSession

	row := tx.QueryRowContext(ctx, "SELECT * FROM sessions where token=?", *token)
	if err = row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rollback(tx, models.ErrNotFound)
		}
	}

	if _, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE id=?", session.ID); err != nil {
		return rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		ctx.Logger.Error(err.Error(), "incident", "logout commit")
		return err
	}

	return nil
}
