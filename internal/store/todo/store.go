package todostore

import (
	"database/sql"
	"errors"
	"time"
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

type Store struct{}

func New() *Store {
	return &Store{}
}

func (s *Store) GetAll(ctx server.Context, userID *uuid.UUID) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := ctx.DB.QueryContext(ctx, getAll, *userID)
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

	return res, nil
}

func (s *Store) Create(ctx server.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	addTS := time.Now().UTC()

	query, values := genInsertQuery(id, title, *userID, addTS)

	ctx.Logger.Debug("generated query", "sql-query", query, "values", values)

	_, err := ctx.DB.ExecContext(ctx, query, values...)
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

func (s *Store) Update(ctx server.Context, id, title string, userID *uuid.UUID) (*models.Task, error) {
	modifiedTS := time.Now()

	query, values := genUpdateQuery(id, title, *userID, modifiedTS)

	ctx.Logger.Debug("generated query", "sql-query", query, "values", values)

	_, err := ctx.DB.ExecContext(ctx, query, values...)
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

func (s *Store) Delete(ctx server.Context, id string, userID *uuid.UUID) error {
	_, err := ctx.DB.ExecContext(ctx, "DELETE FROM tasks WHERE task_id=? AND user_id=?", id, *userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) MarkDone(ctx server.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	var (
		task = models.Task{ID: id}
		done int
	)

	_, err := ctx.DB.ExecContext(ctx, "UPDATE tasks SET done_status=? WHERE task_id=? AND user_id=?", 1, id, *userID)
	if err != nil {
		return nil, err
	}

	row := ctx.DB.QueryRowContext(ctx, `SELECT task_title, done_status, added_at, modified_at FROM tasks 
		WHERE task_id=? AND user_id=?`, id, *userID)
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
