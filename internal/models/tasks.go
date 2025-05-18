package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string     `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsDone      bool       `json:"isDone"`
	DueDate     *time.Time `json:"dueDate"`
	AddedAt     time.Time  `json:"addedAt"`
	ModifiedAt  *time.Time `json:"modifiedAt"`
}

type TaskResp struct {
	ID          string     `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsDone      bool       `json:"isDone"`
	DueDate     *string    `json:"dueDate"`
	AddedAt     time.Time  `json:"addedAt"`
	ModifiedAt  *time.Time `json:"modifiedAt"`
}

type TaskReq struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	IsDone      bool   `json:"isDone"`
}

type Error struct {
	Type    string `json:"type"`
	IsError bool   `json:"isError"`
	Msg     string `json:"msg"`
}

func (t *Task) ToTaskResp() *TaskResp {
	tr := TaskResp{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		IsDone:      t.IsDone,
		AddedAt:     t.AddedAt,
		ModifiedAt:  t.ModifiedAt,
	}

	dd := t.DueDate.Format(time.DateOnly)
	tr.DueDate = &dd

	return &tr
}
