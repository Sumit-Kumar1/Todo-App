package store

import (
	"context"
	"database/sql"
	pkgErr "errors"
	"log/slog"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

type queries struct {
	DeleteTask     string
	GetAllByUserID string
	GetTaskByID    string
	InsertQuery    string
	SetDone        string
	UpdateQuery    string
}

var sqlQueries = queries{
	DeleteTask:     "DELETE FROM tasks WHERE id=$1 AND user_id=$2;",
	GetAllByUserID: "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM tasks WHERE user_id=$1;",
	GetTaskByID:    "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM tasks WHERE id=$1 AND user_id=$2;",
	InsertQuery:    "INSERT INTO tasks (id, user_id, title, description, done_status, due_date, added_at) VALUES ($1, $2, $3, $4, $5, $6, $7);",
	SetDone:        "UPDATE tasks SET done_status=$1, modified_at=$2 WHERE id=$3 AND user_id=$4;",
	UpdateQuery:    "UPDATE tasks SET title=$1, description=$2, done_status=$3, modified_at=$4 WHERE id=$5 AND user_id=$6;",
}

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
		task   models.Task
	)

	rows, err := s.DB.QueryContext(ctx, sqlQueries.GetAllByUserID, *userID)
	if err != nil {
		if pkgErr.Is(err, sql.ErrNoRows) {
			return res, nil
		}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.IsDone,
			&task.DueDate, &task.AddedAt, &task.ModifiedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "get all tasks", slog.String("user", userID.String()))

	return res, nil
}

func (s *Store) Create(ctx context.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, sqlQueries.InsertQuery, task.ID, task.UserID,
		task.Title, task.Description, task.IsDone, task.DueDate, task.AddedAt)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.NewConstError("not able to insert into database")
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task added successfully",
		slog.String("task", task.ID),
	)

	return nil
}

func (s *Store) Update(ctx context.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, sqlQueries.UpdateQuery, task.Title, task.Description,
		task.IsDone, task.ModifiedAt, task.ID, task.UserID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrTaskNotFound
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task updated successfully",
		slog.String("task", task.ID))

	return nil
}

func (s *Store) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, sqlQueries.DeleteTask, id, *userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrTaskNotFound
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task deleted successfully", slog.String("task", id))

	return nil
}

func (s *Store) MarkDone(ctx context.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, sqlQueries.SetDone, 1, time.Now(), id, *userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrTaskNotFound
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task marked as done", slog.String("task", id))
	return nil
}

func (s *Store) GetTaskByID(ctx context.Context, taskID string, userID *uuid.UUID) (*models.Task, error) {
	var task models.Task

	row := s.DB.QueryRowContext(ctx, sqlQueries.GetTaskByID, taskID, *userID)
	err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.IsDone,
		&task.DueDate, &task.AddedAt, &task.ModifiedAt)

	if err != nil {
		if pkgErr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrTaskNotFound
		}

		return nil, err
	}

	return &task, nil
}
