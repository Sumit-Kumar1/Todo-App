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

func New(logger *slog.Logger) (*Store, func(), error) {
	const dbFile string = "./tasks.db"

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		logger.Error("unable to connect sqlite", "error", err.Error())
		return nil, nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("database not reachable", "error", err)
		return nil, nil, err
	}

	fn := func() {
		if err := db.Close(); err != nil {
			logger.Error("DB Close", "error", err)
		}
	}

	if err := runMigration(db); err != nil {
		return nil, fn, err
	}

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

func (s *Store) GetSessionByID(ctx context.Context, userID *uuid.UUID) (*models.UserSession, error) {
	var session models.UserSession

	if userID == nil {
		return nil, errors.New("invalid user id")
	}

	row := s.DB.QueryRowContext(ctx, getSession, *userID)
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.Expiry); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) RefreshSession(ctx context.Context, newSession *models.UserSession) (*models.UserSession, error) {
	_, err := s.DB.ExecContext(ctx, refreshSession, newSession.Token, newSession.Expiry, newSession.ID)
	if err != nil {
		return nil, err
	}

	s.Log.Info("session is refreshed", "userID", newSession.UserID)

	return newSession, nil
}

func (s *Store) GetByEmail(ctx context.Context, userID string) (*models.UserData, error) {
	var user models.UserData

	row := s.DB.QueryRowContext(ctx, getUser, userID)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (s *Store) GetAll(ctx context.Context) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}

	s.Log.Debug("completed query with no error", "sql-query", getAll)

	defer rows.Close()

	for rows.Next() {
		var (
			task models.Task
			done int
		)

		err := rows.Scan(&task.ID, &task.Title, &done, &task.AddedAt, &task.ModifiedAt)
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

	return res, nil
}

func (s *Store) Create(ctx context.Context, id, title string) (*models.Task, error) {
	addTS := time.Now().UTC()

	query, values := genInsertQuery(id, title, addTS)

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

func (s *Store) Update(ctx context.Context, id, title string) (*models.Task, error) {
	modifiedTS := time.Now()

	query, values := genUpdateQuery(id, title, modifiedTS)

	s.Log.Debug("generated query", "sql-query", query, "values", values)

	_, err := s.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	t := models.Task{
		ID:         id,
		Title:      title,
		IsDone:     false,
		ModifiedAt: &modifiedTS,
	}

	return &t, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM tasks WHERE task_id=?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) MarkDone(ctx context.Context, id string) (*models.Task, error) {
	var (
		task = models.Task{ID: id}
		done int
	)

	res, err := s.DB.ExecContext(ctx, "UPDATE tasks SET done_status=? WHERE task_id=?", 1, id)
	if err != nil {
		return nil, err
	}

	val, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if val > 1 {
		return nil, errors.New("not able to update")
	}

	row := s.DB.QueryRowContext(ctx, "SELECT task_title, done_status, added_at, modified_at FROM tasks WHERE task_id=?", id)
	if err := row.Scan(&task.Title, &done, &task.AddedAt, &task.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	if done == 1 {
		task.IsDone = true
	}

	return &task, nil
}
