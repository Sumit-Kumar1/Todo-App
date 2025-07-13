package todostore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	deleteTask     = "DELETE FROM tasks WHERE id=? AND user_id=?;"
	getAllByUserID = "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM tasks WHERE user_id=?;"
	getTaskByID    = "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM " +
		"tasks WHERE id=? AND user_id=?;"
	insertQuery = "INSERT INTO tasks (id, user_id, title, description, done_status, due_date, added_at) VALUES " +
		"(?, ?, ?, ?, ?, ?, ?);"
	setDone     = "UPDATE tasks SET done_status=?, modified_at=? WHERE id=? AND user_id=?;"
	updateQuery = "UPDATE tasks SET title=?, description=?, done_status=?, modified_at=? WHERE id=? AND user_id=?;"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	var (
		res    = make([]models.Task, 0)
		logger = models.GetLoggerFromCtx(ctx)
	)

	rows, err := s.DB.QueryContext(ctx, getAllByUserID, *userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, nil
		}

		return nil, err
	}

	for rows.Next() {
		task, err := populateTaskFields(rows)
		if err != nil {
			return nil, err
		}

		res = append(res, *task)
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "get all tasks", slog.String("user", userID.String()))

	return res, nil
}

func (s *Store) Create(ctx context.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	if _, err := s.DB.ExecContext(ctx, insertQuery, task.ID, task.UserID,
		task.Title, task.Description, task.IsDone, task.DueDate.UnixMilli(),
		task.AddedAt.UnixMilli()); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task added successfully",
		slog.String("task", task.ID),
	)

	return nil
}

func (s *Store) Update(ctx context.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	if _, err := s.DB.ExecContext(ctx, updateQuery, task.Title, task.Description,
		task.IsDone, task.ModifiedAt.UnixMilli(), task.ID, task.UserID); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task updated successfully",
		slog.String("task", task.ID))

	return nil
}

func (s *Store) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	if _, err := s.DB.ExecContext(ctx, deleteTask, id, *userID); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task deleted successfully", slog.String("task", id))

	return nil
}

func (s *Store) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	var (
		task   = &models.Task{ID: id}
		logger = models.GetLoggerFromCtx(ctx)
		err    error
	)

	res, err := s.DB.ExecContext(ctx, setDone, 1, time.Now().UnixMilli(), id, *userID)
	if err != nil {
		return nil, err
	}

	if _, err := res.RowsAffected(); err != nil {
		return nil, err
	}

	rows, err := s.DB.QueryContext(ctx, getTaskByID, id, *userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound("task")
		}

		return nil, err
	}

	for rows.Next() {
		task, err = populateTaskFields(rows)
		if err != nil {
			return nil, err
		}
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task marked done",
		slog.String("task", id), slog.String("user", userID.String()))

	return task, nil
}

func populateTaskFields(rows *sql.Rows) (*models.Task, error) {
	var (
		task models.Task
		err  error
	)

	err = rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.IsDone,
		&task.DueDate, &task.AddedAt, &task.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
