package store

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Store struct {
	DB  *sql.DB
	Log *slog.Logger
}

func New(logger *slog.Logger, dbFile string) (*Store, func(), error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		logger.Error("unable to connect sqlite", "error", err.Error())
		return nil, nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("database not reachable", "error", err)
		return nil, nil, err
	}

	logger.Info("database connection established successfully", "stats", db.Stats())

	fn := func() {
		if err := db.Close(); err != nil {
			logger.Error("DB Close", "error", err)
		}
	}

	if err := runMigration(db); err != nil {
		return nil, fn, err
	}

	logger.Info("migrations run successfully")

	return &Store{Log: logger, DB: db}, fn, nil
}

func (s *Store) RegisterUser(ctx context.Context, data *models.UserData, session *models.UserSession) (*models.UserSession, error) {
	var (
		err  error
		opts = sql.TxOptions{Isolation: sql.LevelReadCommitted}
	)

	tx, err := s.DB.BeginTx(ctx, &opts)
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

func (s *Store) CreateSession(ctx context.Context, session *models.UserSession) error {
	_, err := s.DB.ExecContext(ctx, createSession, session.ID, session.UserID, session.Token, session.Expiry)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error) {
	var session models.UserSession

	if userID == nil {
		return nil, models.ErrInvalid("user ID")
	}

	row := s.DB.QueryRowContext(ctx, getSession, *userID)
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound("user")
		}

		return nil, err
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.UserSession) (*models.UserSession, error) {
	_, err := s.DB.ExecContext(ctx, refreshSession, newSession.Token, newSession.Expiry, newSession.ID)
	if err != nil {
		return nil, err
	}

	s.Log.Info("session is refreshed", "user", newSession.UserID)

	return newSession, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*models.UserData, error) {
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
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	var session models.UserSession

	row := tx.QueryRowContext(ctx, "SELECT * FROM sessions where token=?", *token)
	if err = row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rollback(tx, models.ErrNotFound("session with current user"))
		}
	}

	if _, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE id=?", session.ID); err != nil {
		return rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		s.Log.Error(err.Error(), "incident", "logout commit")
		return err
	}

	return nil
}

func (s *Store) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.QueryContext(ctx, getAll, *userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			task models.Task
			done int
		)

		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &done, &task.AddedAt, &task.ModifiedAt)
		if err != nil {
			return nil, err
		}

		switch done {
		case 0:
			task.IsDone = false
		case 1:
			task.IsDone = true
		}

		res = append(res, task)
	}

	s.Log.Info("Get all the tasks", "user", userID)

	return res, nil
}

func (s *Store) Create(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	addTS := time.Now().UTC()

	query, values := genInsertQuery(id, title, *userID, addTS)

	s.Log.Debug("generated query", "sql-query", query, "values", values)

	_, err := s.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	var task = models.Task{
		ID:      id,
		Title:   title,
		IsDone:  false,
		AddedAt: addTS,
	}

	return &task, nil
}

func (s *Store) Update(ctx context.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	modifiedTS := time.Now()

	query, values := genUpdateQuery(id, title, *userID, modifiedTS)

	s.Log.DebugContext(ctx, "generated query", "sql-query", query, "values", values, "user", userID)

	_, err := s.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	t := models.Task{
		ID:         id,
		UserID:     *userID,
		Title:      title,
		IsDone:     false,
		ModifiedAt: &modifiedTS,
	}

	return &t, nil
}

func (s *Store) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM tasks WHERE task_id=? AND user_id=?", id, *userID)
	if err != nil {
		return err
	}

	s.Log.DebugContext(ctx, "generated query", "sql-query", "DELETE FROM tasks WHERE task_id=? AND user_id=?", "task", id, "user", userID)

	return nil
}

func (s *Store) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	var (
		task = models.Task{ID: id}
		done int
	)

	_, err := s.DB.ExecContext(ctx, "UPDATE tasks SET done_status=? WHERE task_id=? AND user_id=?", 1, id, *userID)
	if err != nil {
		return nil, err
	}

	row := s.DB.QueryRowContext(ctx, `SELECT task_title, done_status, added_at, modified_at FROM tasks 
		WHERE task_id=? AND user_id=?`, id, *userID)
	if err := row.Scan(&task.Title, &done, &task.AddedAt, &task.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound("task")
		}

		return nil, err
	}

	if done == 1 {
		task.IsDone = true
	}

	s.Log.Info("task done", "taskID", id, "user", userID)

	return &task, nil
}
