package todostore

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/sqlitecloud/sqlitecloud-go"
)

const (
	deleteTask     = "DELETE FROM tasks WHERE id='%v' AND user_id='%v';"
	getAllByUserID = "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM tasks WHERE user_id='%v';"
	getTaskByID    = "SELECT id, user_id, title, description, done_status, due_date, added_at, modified_at FROM " +
		"tasks WHERE id='%v' AND user_id='%v';"
	insertQuery = "INSERT INTO tasks (id, user_id, title, description, done_status, due_date, added_at) VALUES " +
		"('%v', '%v', '%v', '%v', %v, '%v', '%v');"
	setDone     = "UPDATE tasks SET done_status=%v, modified_at='%v' WHERE id='%v' AND user_id='%v';"
	updateQuery = "UPDATE tasks SET title='%v', done_status=%v, modified_at='%v' WHERE id='%v' AND user_id='%v';"
)

type Store struct {
	DB *sqlitecloud.SQCloud
}

func New(db *sqlitecloud.SQCloud) *Store {
	return &Store{DB: db}
}

func (s *Store) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	var (
		res    = make([]models.Task, 0)
		logger = models.GetLoggerFromCtx(ctx)
	)

	rows, err := s.DB.Select(fmt.Sprintf(getAllByUserID, *userID))
	if err != nil {
		return nil, err
	}

	numRows := rows.GetNumberOfRows()

	if numRows == 0 {
		return res, nil
	}

	for row := uint64(0); row < numRows; row++ {
		task, err := populateTaskFields(rows, row)
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

	if err := s.DB.Execute(query); err != nil {
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
		task.IsDone,
		task.ModifiedAt.UnixMilli(),
		task.ID,
		task.UserID,
	)

	if err := s.DB.Execute(query); err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task updated successfully",
		slog.String("task", task.ID))

	return nil
}

func (s *Store) Delete(ctx context.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	if err := s.DB.Execute(fmt.Sprintf(deleteTask, id, *userID)); err != nil {
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

	if err := s.DB.Execute(fmt.Sprintf(setDone, true, time.Now().UnixMilli(), id, *userID)); err != nil {
		return nil, err
	}

	row, err := s.DB.Select(fmt.Sprintf(getTaskByID, id, *userID))
	if err != nil {
		return nil, err
	}

	if row.GetNumberOfRows() == 0 {
		return nil, models.ErrNotFound("task")
	}

	for i := uint64(0); i < row.GetNumberOfRows(); i++ {
		task, err = populateTaskFields(row, i)
		if err != nil {
			return nil, err
		}
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task marked done",
		slog.String("task", id), slog.String("user", userID.String()))

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
