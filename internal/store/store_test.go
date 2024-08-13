package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"todoapp/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

var (
	id    = "abcde"
	td    = time.Now()
	title = "Dog to walk"
	ctx   = context.Background()
	errDB = errors.New("some db error")
	cols  = []string{"task_title", "done_status", "added_at", "modified_at"}
)

func initMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("not able to init mocks:%v", err.Error())

		return nil, nil, err
	}

	return db, mock, nil
}

func TestStore_GetAll(t *testing.T) {
	db, mock, err := initMock(t)
	if err != nil {
		return
	}

	s := Store{
		DB: db,
	}

	col := []string{"id", "task_title", "done_status", "added_at", "modified_at"}

	tests := []struct {
		name    string
		sqlMock any
		want    []models.Task
		wantErr error
	}{
		{
			name:    "valid query - done task",
			sqlMock: mock.ExpectQuery(getAll).WillReturnRows(sqlmock.NewRows(col).AddRow(id, title, 1, td, td)),
			want:    []models.Task{{ID: id, Title: title, IsDone: true, AddedAt: td, ModifiedAt: &td}},
		},
		{
			name:    "valid query",
			sqlMock: mock.ExpectQuery(getAll).WillReturnRows(sqlmock.NewRows(col).AddRow(id, title, 0, td, td)),
			want:    []models.Task{{ID: id, Title: title, IsDone: false, AddedAt: td, ModifiedAt: &td}},
		},
		{
			name:    "no row err",
			sqlMock: mock.ExpectQuery(getAll).WillReturnRows(sqlmock.NewRows(col).CloseError(errDB)),
			want:    []models.Task{}, wantErr: nil,
		},
		{
			name:    "query err",
			sqlMock: mock.ExpectQuery(getAll).WillReturnError(errDB),
			want:    nil, wantErr: errDB,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetAll(ctx)

			assert.Equalf(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestStore_Create(t *testing.T) {
	db, mock, err := initMock(t)
	if err != nil {
		return
	}

	s := Store{
		DB: db,
	}

	tests := []struct {
		name    string
		id      string
		title   string
		mock    any
		want    *models.Task
		wantErr error
	}{
		{
			name:  "valid create",
			id:    id,
			title: title,
			mock:  mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlmock.NewResult(1, 1)),
			want:  &models.Task{ID: id, Title: title, IsDone: false, AddedAt: td},
		},
		{
			name:  "time err",
			id:    id,
			title: title,
			mock:  mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlmock.NewResult(1, 1)),
			want:  &models.Task{ID: id, Title: title, IsDone: false, AddedAt: td},
		},
		{
			name:    "exec error",
			id:      id,
			title:   title,
			mock:    mock.ExpectExec("INSERT INTO tasks").WillReturnError(errDB),
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(context.Background(), tt.id, tt.title)

			assert.Equalf(t, tt.wantErr, err, "Test[%s] Failed", tt.name)

			if got != nil {
				got.AddedAt = td
				assert.Equalf(t, tt.want, got, "Test[%s] Failed", tt.name)
			}
		})
	}
}

func TestStore_Update(t *testing.T) {
	db, mock, err := initMock(t)
	if err != nil {
		return
	}

	s := Store{
		DB: db,
	}

	updateQuery := "UPDATE tasks"

	tests := []struct {
		name    string
		id      string
		title   string
		mock    any
		want    *models.Task
		wantErr error
	}{
		{
			name:  "valid update",
			id:    id,
			title: title,
			mock:  mock.ExpectExec(updateQuery).WillReturnResult(sqlmock.NewResult(1, 1)),
			want:  &models.Task{ID: id, Title: title, IsDone: false, ModifiedAt: &td},
		},
		{
			name:    "db error",
			id:      id,
			title:   title,
			mock:    mock.ExpectExec(updateQuery).WillReturnError(errDB),
			wantErr: errDB,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Update(ctx, tt.id, tt.title)

			assert.Equalf(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)

			if got != nil {
				got.ModifiedAt = &td
				assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
			}
		})
	}
}

func TestStore_Delete(t *testing.T) {
	db, mock, err := initMock(t)
	if err != nil {
		return
	}

	s := Store{
		DB: db,
	}

	delQuery := "DELETE FROM tasks WHERE task_id=?"

	tests := []struct {
		name    string
		id      string
		mock    any
		wantErr error
	}{
		{
			name: "valid delete",
			id:   id,
			mock: mock.ExpectExec(delQuery).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			name:    "db error",
			id:      id,
			mock:    mock.ExpectExec(delQuery).WillReturnError(errDB),
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("DELETE FROM tasks").WillReturnResult(sqlmock.NewResult(1, 1))
			err := s.Delete(context.Background(), tt.id)

			assert.Equalf(t, tt.wantErr, err, "Test[%s] Failed", tt.name)
		})
	}
}

func TestStore_MarkDone(t *testing.T) {
	db, mock, err := initMock(t)
	if err != nil {
		return
	}

	s := Store{
		DB: db,
	}

	selectQuery := "SELECT task_title, done_status, added_at, modified_at FROM tasks"

	tests := []struct {
		name    string
		id      string
		mocks   []any
		want    *models.Task
		wantErr error
	}{
		{
			name: "exec error",
			id:   id,
			mocks: []any{
				mock.ExpectExec("UPDATE tasks").WillReturnError(errDB),
			},
			want:    nil,
			wantErr: errDB,
		},
		{
			name: "rows affected error",
			id:   id,
			mocks: []any{
				mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlmock.NewErrorResult(errDB)),
			},
			want:    nil,
			wantErr: errDB,
		},
		{
			name: "unexpected updates",
			id:   id,
			mocks: []any{
				mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlmock.NewResult(0, 3)),
			},
			want:    nil,
			wantErr: errors.New("not able to update"),
		},
		{
			name: "row scan err",
			id:   id,
			mocks: []any{mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.ExpectQuery(selectQuery).WillReturnRows(sqlmock.NewRows(cols).RowError(-1, errDB))},
			wantErr: sql.ErrNoRows,
		},
		{
			name: "valid mark done",
			id:   id,
			mocks: []any{mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlmock.NewResult(0, 1)),
				mock.ExpectQuery(selectQuery).WillReturnRows(sqlmock.NewRows(cols).AddRow(title, 1, td, td))},
			want: &models.Task{ID: id, Title: title, IsDone: true, AddedAt: td, ModifiedAt: &td},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.MarkDone(ctx, tt.id)

			assert.Equalf(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}
