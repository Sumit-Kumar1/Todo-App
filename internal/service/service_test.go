package service

import (
	"testing"
	"time"
	"todoapp/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestService_GetAll(t *testing.T) {
	s := New()
	tz := time.Now()
	tz2 := time.Now().Add(1000)

	t1 := models.Task{
		ID:         "abcuie",
		Title:      "Take Dog to walk",
		IsDone:     false,
		AddedAt:    tz,
		ModifiedAt: time.Time{},
	}

	t2 := models.Task{
		ID:         "abcde",
		Title:      "Take Dog to walk",
		IsDone:     false,
		AddedAt:    tz2,
		ModifiedAt: time.Time{},
	}

	tests := []struct {
		name string
		Data map[string]*models.Task
		want []models.Task
	}{
		{name: "empty case", Data: map[string]*models.Task{"abcuie": &t1}, want: []models.Task{t1}},
		{name: "nil case", Data: map[string]*models.Task{"abcuie": nil}, want: []models.Task(nil)},
		{name: "valid case", Data: map[string]*models.Task{"abcuie": &t1, "abcde": &t2}, want: []models.Task{t1, t2}},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Data = tt.Data

			assert.Equalf(t, tt.want, s.GetAll(), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestService_AddTask(t *testing.T) {
	id := "abcde"
	title := "Dog to walk"
	ts := time.Now()

	s := New()

	s.Data[id] = &models.Task{}

	tests := []struct {
		name    string
		id      string
		title   string
		want    *models.Task
		wantErr error
	}{
		{name: "invalid task title", title: "", wantErr: models.ErrInvalidTitle},
		{name: "valid case", title: title, id: id, want: &models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts}, wantErr: nil},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.AddTask(tt.title)

			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)

			if resp != nil {
				resp.ID = tt.id
				resp.AddedAt = ts

				assert.Equalf(t, tt.want, resp, "TEST[%d] Failed - %s", i, tt.name)
			}
		})
	}
}

func TestService_DeleteTask(t *testing.T) {
	id := "abcde"
	title := "Dog to walk"
	ts := time.Now()

	s := New()

	s.Data[id] = &models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts}

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{name: "id not found", id: "IDcde", wantErr: models.ErrNotFound},
		{name: "invalid id", id: "ID124cde", wantErr: models.ErrInvalidID},
		{name: "empty id", id: "", wantErr: models.ErrInvalidID},
		{name: "valid case", id: id, wantErr: nil},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteTask(tt.id)

			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestService_MarkDone(t *testing.T) {
	id := "abcde"
	title := "Dog to walk"
	ts := time.Now()

	s := New()

	s.Data[id] = &models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts}

	tests := []struct {
		name    string
		id      string
		want    *models.Task
		wantErr error
	}{
		{name: "valid case", id: id, want: &models.Task{ID: id, Title: title, IsDone: true, AddedAt: ts, ModifiedAt: ts}, wantErr: nil},
		{name: "invalid id", id: "128bakdhiue", want: nil, wantErr: models.ErrInvalidID},
		{name: "not found case", id: "abcze", want: nil, wantErr: models.ErrNotFound},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.MarkDone(tt.id)

			formatTime(resp, &ts)

			assert.Equalf(t, tt.want, resp, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestService_UpdateTask(t *testing.T) {
	id := "abcde"
	title := "Dog to walk"
	ts := time.Now()

	s := New()

	s.Data[id] = &models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts}

	type args struct {
		id     string
		title  string
		isDone string
	}

	tests := []struct {
		name    string
		args    args
		want    *models.Task
		wantErr error
	}{
		{name: "valid case", args: args{id: id, title: title, isDone: "false"}, want: &models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts, ModifiedAt: ts}, wantErr: nil},
		{name: "invalid task", args: args{id: id, title: "", isDone: "false"}, want: nil, wantErr: models.ErrInvalidTitle},
		{name: "invalid task", args: args{id: "abcd", title: "hello world", isDone: "false"}, want: nil, wantErr: models.ErrNotFound},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.UpdateTask(tt.args.id, tt.args.title, tt.args.isDone)

			formatTime(resp, &ts)

			assert.Equalf(t, tt.want, resp, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func formatTime(task *models.Task, timestamp *time.Time) {
	if task == nil || timestamp == nil {
		return
	}

	task.ModifiedAt = *timestamp
	task.AddedAt = *timestamp
}
