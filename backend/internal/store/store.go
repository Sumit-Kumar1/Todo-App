package store

import (
	"database/sql"
	liberrors "errors"
	"log/slog"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

const (
	deleteTask     = "DELETE FROM tasks WHERE id=$1 AND user_id=$2;"
	getAllByUserID = "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM tasks WHERE user_id=$1;"
	getTaskByID    = "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM " +
		"tasks WHERE id=$1 AND user_id=$2;"
	insertQuery = "INSERT INTO tasks (id, user_id, title, description, done_status, due_date, added_at) VALUES " +
		"($1, $2, $3, $4, $5, $6, $7);"
	setDone     = "UPDATE tasks SET done_status=$1, modified_at=$2 WHERE id=$3 AND user_id=$4;"
	updateQuery = "UPDATE tasks SET title=$1, description=$2, done_status=$3, modified_at=$4 WHERE id=$5 AND user_id=$6;"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (s *Store) GetAll(ctx *gofr.Context, userID *uuid.UUID) ([]models.Task, error) {
	var (
		res    = make([]models.Task, 0)
		logger = models.GetLoggerFromCtx(ctx)
	)

	rows, err := ctx.SQL.QueryContext(ctx, getAllByUserID, *userID)
	if err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
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

func (s *Store) Create(ctx *gofr.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	if _, err := ctx.SQL.ExecContext(ctx, insertQuery, task.ID, task.UserID,
		task.Title, task.Description, task.IsDone, task.DueDate, task.AddedAt); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task added successfully",
		slog.String("task", task.ID),
	)

	return nil
}

func (s *Store) Update(ctx *gofr.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	if _, err := ctx.SQL.ExecContext(ctx, updateQuery, task.Title, task.Description,
		task.IsDone, task.ModifiedAt, task.ID, task.UserID); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task updated successfully",
		slog.String("task", task.ID))

	return nil
}

func (s *Store) Delete(ctx *gofr.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	if _, err := ctx.SQL.ExecContext(ctx, deleteTask, id, *userID); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task deleted successfully", slog.String("task", id))

	return nil
}

func (s *Store) MarkDone(ctx *gofr.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	var (
		task = &models.Task{ID: id}
		err  error
	)

	res, err := ctx.SQL.ExecContext(ctx, setDone, 1, time.Now(), id, *userID)
	if err != nil {
		return nil, err
	}

	if _, err := res.RowsAffected(); err != nil {
		return nil, err
	}

	rows, err := ctx.SQL.QueryContext(ctx, getTaskByID, id, *userID)
	if err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound("task")
		}

		return nil, err
	}

	for rows.Next() {
		task, err = populateTaskFields(rows)
		if err != nil {
			return nil, err
		}
	}

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
