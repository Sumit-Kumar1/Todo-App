package todostore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	deleteTask  = "DELETE FROM tasks WHERE task_id=? AND user_id=?"
	getAll      = "SELECT task_id, user_id, task_title, done_status, added_at, modified_at from tasks WHERE user_id=?"
	getTask     = "SELECT task_title, done_status, added_at, modified_at FROM tasks WHERE task_id=? AND user_id=?"
	insertQuery = "INSERT INTO tasks (task_id, user_id, task_title, done_status, added_at) VALUES (?, ?, ?, ?, ?);"
	setDone     = "UPDATE tasks SET done_status=? WHERE task_id=? AND user_id=?"
	updateQuery = "UPDATE tasks SET task_title=?, done_status=?, modified_at=? WHERE task_id=? AND user_id=?;"
)

type Store struct {
	DB  *sql.DB
	Log *slog.Logger
}

func New(db *sql.DB, logger *slog.Logger) *Store {
	return &Store{
		DB:  db,
		Log: logger,
	}
}

func (s *Store) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.QueryContext(ctx, getAll, *userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.IsDone, &task.AddedAt, &task.ModifiedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, task)
	}

	s.Log.Info("Get all the tasks", "user", userID)

	return res, nil
}

func (s *Store) Create(ctx context.Context, task *models.Task) error {
	sr, err := s.DB.ExecContext(ctx, insertQuery, task.ID, task.UserID, task.Title, task.IsDone, task.AddedAt)
	if err != nil {
		return err
	}

	if _, err := sr.LastInsertId(); err != nil {
		return errors.New("task not created")
	}

	return nil
}

func (s *Store) Update(ctx context.Context, task *models.Task) error {
	sr, err := s.DB.ExecContext(ctx, updateQuery, task.Title, task.IsDone, task.ModifiedAt, task.ID, task.UserID)
	if err != nil {
		return err
	}

	if _, err := sr.LastInsertId(); err != nil {
		return errors.New("task not updated")
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	sr, err := s.DB.ExecContext(ctx, deleteTask, id, *userID)
	if err != nil {
		return err
	}

	if _, err := sr.LastInsertId(); err != nil {
		return errors.New("task not deleted")
	}

	return nil
}

func (s *Store) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	var (
		task = models.Task{ID: id}
		vals = []any{true, id, *userID}
	)

	rs, err := s.DB.ExecContext(ctx, setDone, vals...)
	if err != nil {
		return nil, err
	}

	if _, err := rs.LastInsertId(); err != nil {
		return nil, errors.New("task not marked done")
	}

	row := s.DB.QueryRowContext(ctx, getTask, id, *userID)
	if err := row.Scan(&task.Title, &task.IsDone, &task.AddedAt, &task.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound("task")
		}

		return nil, err
	}

	s.Log.Info("task done", "taskID", id, "user", userID)

	return &task, nil
}
