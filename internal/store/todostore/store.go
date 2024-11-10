package todostore

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/sqlitecloud/sqlitecloud-go"
)

const (
	deleteTask  = "DELETE FROM tasks WHERE task_id=%v AND user_id=%v"
	getAll      = "SELECT task_id, user_id, task_title, done_status, added_at, modified_at from tasks WHERE user_id=%v"
	getTask     = "SELECT task_title, done_status, added_at, modified_at FROM tasks WHERE task_id=%v AND user_id=%v"
	insertQuery = "INSERT INTO tasks (task_id, user_id, task_title, done_status, added_at) VALUES (%v, %v, %v, %v, %v);"
	setDone     = "UPDATE tasks SET done_status=%v WHERE task_id=%v AND user_id=%v"
	updateQuery = "UPDATE tasks SET task_title=%v, done_status=%v, modified_at=%v WHERE task_id=%v AND user_id=%v;"
)

type Store struct {
	DB  *sqlitecloud.SQCloud
	Log *slog.Logger
}

func New(db *sqlitecloud.SQCloud, logger *slog.Logger) *Store {
	return &Store{
		DB:  db,
		Log: logger,
	}
}

func (s *Store) GetAll(_ context.Context, userID *uuid.UUID) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.Select(fmt.Sprintf(getAll, *userID))
	if err != nil {
		return nil, err
	}

	numRows := rows.GetNumberOfRows()

	if numRows == 0 {
		return res, nil
	}

	for r := uint64(0); r < numRows; r++ {
		var task models.Task

		task.ID, err = rows.GetStringValue(r, 0)
		if err != nil {
			return nil, err
		}
		v2, err := rows.GetStringValue(r, 1)
		if err != nil {
			return nil, err
		}
		task.Title, err = rows.GetStringValue(r, 2)
		if err != nil {
			return nil, err
		}

		v4, err := rows.GetInt64Value(r, 3)
		if err != nil {
			return nil, err
		}

		task.AddedAt, err = rows.GetSQLDateTime(r, 4)
		if err != nil {
			return nil, err
		}

		v6, err := rows.GetSQLDateTime(r, 5)
		if err != nil {
			return nil, err
		}

		task.UserID = uuid.MustParse(v2)
		task.IsDone = (v4 == 1)
		task.ModifiedAt = &v6

		res = append(res, task)
	}

	s.Log.Info("Get all the tasks", "user", userID)

	return res, nil
}

func (s *Store) Create(_ context.Context, task *models.Task) error {
	if err := s.DB.Execute(fmt.Sprintf(insertQuery, task.ID, task.UserID, task.Title, task.IsDone, task.AddedAt)); err != nil {
		return err
	}

	return nil
}

func (s *Store) Update(_ context.Context, task *models.Task) error {
	if err := s.DB.Execute(fmt.Sprintf(updateQuery, task.Title, task.IsDone, task.ModifiedAt, task.ID, task.UserID)); err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(_ context.Context, id string, userID *uuid.UUID) error {
	if err := s.DB.Execute(fmt.Sprintf(deleteTask, id, *userID)); err != nil {
		return err
	}

	return nil
}

func (s *Store) MarkDone(_ context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	var (
		task = models.Task{ID: id}
	)

	if err := s.DB.Execute(fmt.Sprintf(setDone, true, id, *userID)); err != nil {
		return nil, err
	}

	row, err := s.DB.Select(fmt.Sprintf(getTask, id, *userID))
	if err != nil {
		return nil, err
	}

	if row.GetNumberOfRows() == 0 {
		return nil, models.ErrNotFound("task")
	}

	for i := uint64(0); i < row.GetNumberOfRows(); i++ {
		c1, err := row.GetStringValue(i, 0)
		if err != nil {
			return nil, err
		}

		c2, err := row.GetStringValue(i, 1)
		if err != nil {
			return nil, err
		}

		c3, err := row.GetSQLDateTime(i, 2)
		if err != nil {
			return nil, err
		}

		c4, err := row.GetSQLDateTime(i, 3)
		if err != nil {
			return nil, err
		}

		task.Title = c1
		if task.IsDone, err = strconv.ParseBool(c2); err != nil {
			return nil, err
		}

		task.AddedAt = c3
		task.ModifiedAt = &c4
	}

	s.Log.Info("task done", "taskID", id, "user", userID)

	return &task, nil
}
