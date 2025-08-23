package todostore

import (
	"context"
	"database/sql"
	liberrors "errors"
	"log/slog"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
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

func (s *Store) Create(ctx context.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)
	query := fmt.Sprintf(insertQuery,
		task.ID,
		task.UserID,
		task.Title,
		task.Description,
		task.IsDone,
		task.DueDate.UnixMilli(),
		task.AddedAt.UnixMilli(),
	)

	if _, err := s.DB.ExecContext(ctx, insertQuery, task.ID, task.UserID,
		task.Title, task.Description, task.IsDone, task.DueDate, task.AddedAt); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task added successfully",
		slog.String("task", task.ID),
	)

	return nil
}

func (s *Store) Update(ctx context.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)
	query := fmt.Sprintf(updateQuery,
		task.Title,
		task.Description,
		task.IsDone,
		task.ModifiedAt.UnixMilli(),
		task.ID,
		task.UserID,
	)

	if _, err := s.DB.ExecContext(ctx, updateQuery, task.Title, task.Description,
		task.IsDone, task.ModifiedAt, task.ID, task.UserID); err != nil {
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
		task = &models.Task{ID: id}
		err  error
	)

	res, err := s.DB.ExecContext(ctx, setDone, 1, time.Now(), id, *userID)
	if err != nil {
		return nil, err
	}

	if _, err := res.RowsAffected(); err != nil {
		return nil, err
	}

	rows, err := s.DB.QueryContext(ctx, getTaskByID, id, *userID)
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

	return task, nil
}

func populateTaskFields(rows *sqlitecloud.Result, r uint64) (*models.Task, error) {
	var (
		task models.Task
		err  error
	)

	task.ID, err = rows.GetStringValue(r, 0) // taskID
	if err != nil {
		return nil, err
	}

	v2, err := rows.GetStringValue(r, 1) // userID
	if err != nil {
		return nil, err
	}

	task.Title, err = rows.GetStringValue(r, 2) // title
	if err != nil {
		return nil, err
	}

	task.Description, err = rows.GetStringValue(r, 3) // description
	if err != nil {
		return nil, err
	}

	v4, err := rows.GetInt64Value(r, 4) // done status
	if err != nil {
		return nil, err
	}

	dd, err := rows.GetInt64Value(r, 5) // due date
	if err != nil {
		return nil, err
	}

	v5, err := rows.GetInt64Value(r, 6) // added time
	if err != nil {
		return nil, err
	}

	v6 := rows.GetInt64Value_(r, 7) // modified time

	task.UserID = uuid.MustParse(v2)

	task.IsDone = (v4 == 1)
	task.DueDate = dateRef(dd)
	task.AddedAt = *dateRef(v5)
	task.ModifiedAt = dateRef(v6)

	return &task, nil
}

func dateRef(data int64) *time.Time {
	if data == 0 {
		return nil
	}

	t := time.UnixMilli(data)

	return &t
}
