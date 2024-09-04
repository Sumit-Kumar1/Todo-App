package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
	"todoapp/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	id     = "abcde"
	ts     = time.Now()
	ctx    = context.Background()
	title  = "Dog to walk"
	errSvc = errors.New("some error occurred")
)

func TestService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	st := NewMockStorer(ctrl)

	s := New(st, slog.Default())
	tz2 := time.Now().Add(1000)
	uid := uuid.New()

	t1 := models.Task{
		ID:         id,
		Title:      title,
		IsDone:     false,
		AddedAt:    ts,
		ModifiedAt: &time.Time{},
	}

	t2 := models.Task{
		ID:         id,
		Title:      title,
		IsDone:     false,
		AddedAt:    tz2,
		ModifiedAt: &time.Time{},
	}

	tests := []struct {
		name    string
		mock    *gomock.Call
		want    []models.Task
		wantErr error
	}{
		{name: "valid case", mock: st.EXPECT().GetAll(ctx, &uid).Return([]models.Task{t1, t2}, nil), want: []models.Task{t1, t2}},
		{name: "empty case", mock: st.EXPECT().GetAll(ctx, &uid).Return([]models.Task{}, nil), want: []models.Task{}},
		{name: "err case", mock: st.EXPECT().GetAll(ctx, &uid).Return(nil, errSvc), want: nil, wantErr: errSvc},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := s.GetAll(ctx)

			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, tasks, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestService_AddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	st := NewMockStorer(ctrl)
	s := New(st, slog.Default())

	task := models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts}
	uid := uuid.New()

	tests := []struct {
		name    string
		id      string
		title   string
		mock    *gomock.Call
		want    *models.Task
		wantErr error
	}{
		{name: "invalid task title", title: "", wantErr: models.ErrInvalidTitle},
		{name: "valid case", title: title, id: id, mock: st.EXPECT().Create(ctx, gomock.Any(), title, &uid).Return(&task, nil),
			want: &task, wantErr: nil},
		{name: "err case", title: title, id: id, mock: st.EXPECT().Create(ctx, gomock.Any(), title, &uid).Return(nil, errSvc),
			want: nil, wantErr: errSvc},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.AddTask(ctx, tt.title)

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
	ctrl := gomock.NewController(t)
	st := NewMockStorer(ctrl)
	s := New(st, slog.Default())
	uid := uuid.New()

	tests := []struct {
		name    string
		id      string
		mock    *gomock.Call
		wantErr error
	}{
		{name: "id not found", id: id, mock: st.EXPECT().Delete(ctx, id, uid).Return(models.ErrNotFound), wantErr: models.ErrNotFound},
		{name: "invalid id", id: "ID124cde", wantErr: models.ErrInvalidID},
		{name: "empty id", id: "", wantErr: models.ErrInvalidID},
		{name: "valid case", id: id, mock: st.EXPECT().Delete(ctx, id, uid).Return(nil), wantErr: nil},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteTask(ctx, tt.id)

			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestService_MarkDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	st := NewMockStorer(ctrl)
	s := New(st, slog.Default())
	task := models.Task{ID: id, Title: title, IsDone: true, AddedAt: ts, ModifiedAt: &ts}
	uid := uuid.New()

	tests := []struct {
		name    string
		id      string
		mock    *gomock.Call
		want    *models.Task
		wantErr error
	}{
		{name: "valid case", id: id, mock: st.EXPECT().MarkDone(ctx, id, uid).Return(&task, nil), want: &task, wantErr: nil},
		{name: "invalid id", id: "128bakdhiue", want: nil, wantErr: models.ErrInvalidID},
		{name: "not found case", id: "abcze", mock: st.EXPECT().MarkDone(ctx, "abcze", uid).Return(nil, models.ErrNotFound), wantErr: models.ErrNotFound},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.MarkDone(ctx, tt.id)

			formatTime(resp, &ts)

			assert.Equalf(t, tt.want, resp, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.wantErr, err, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	st := NewMockStorer(ctrl)
	s := New(st, slog.Default())
	task := models.Task{ID: id, Title: title, IsDone: false, AddedAt: ts, ModifiedAt: &ts}
	uid := uuid.New()
	type args struct {
		id     string
		title  string
		isDone string
	}

	tests := []struct {
		name    string
		args    args
		mock    *gomock.Call
		want    *models.Task
		wantErr error
	}{
		{name: "valid case", args: args{id: id, title: title, isDone: "false"}, mock: st.EXPECT().Update(ctx, id, title, uid).Return(&task, nil), want: &task},
		{name: "invalid task", args: args{id: id, title: "", isDone: "false"}, want: nil, wantErr: models.ErrInvalidTitle},
		{name: "invalid task", args: args{id: "abcd", title: "hello world", isDone: "false"}, mock: st.EXPECT().Update(ctx, "abcd", "hello world", uid).
			Return(nil, models.ErrNotFound), wantErr: models.ErrNotFound},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.UpdateTask(ctx, tt.args.id, tt.args.title, tt.args.isDone)

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

	task.ModifiedAt = timestamp
	task.AddedAt = *timestamp
}
