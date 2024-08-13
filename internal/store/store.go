package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"todoapp/internal/models"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func New() (*Store, error) {
	const dbFile string = "./tasks.db"

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Printf("unable to connect sqlite: %s", err.Error())
		return nil, err
	}

	const create string = `
		CREATE TABLE IF NOT EXISTS tasks(
		task_id TEXT NOT NULL PRIMARY KEY,
		task_title TEXT NOT NULL,
		done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
		added_at DATETIME NOT NULL,
		modified_at DATETIME);`

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &Store{
		DB: db,
	}, nil
}

func (s Store) GetAll(ctx context.Context) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}

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

func (s Store) Create(ctx context.Context, id, title string) (*models.Task, error) {
	addTS := time.Now()

	query, values := genInsertQuery(id, title, addTS)

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

func (s Store) Update(ctx context.Context, id, title string) (*models.Task, error) {
	modifiedTS := time.Now()

	query, values := genUpdateQuery(id, title, modifiedTS)

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

func (s Store) Delete(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM tasks WHERE task_id=?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) MarkDone(ctx context.Context, id string) (*models.Task, error) {
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
		return nil, err
	}

	if done == 1 {
		task.IsDone = true
	}

	return &task, nil
}
