package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID         string     `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Title      string     `json:"title"`
	IsDone     bool       `json:"isDone"`
	DueDate    *time.Time `json:"dueDate"`
	AddedAt    time.Time  `json:"addedAt"`
	ModifiedAt *time.Time `json:"modifiedAt"`
}

type Error struct {
	Type    string `json:"type"`
	IsError bool   `json:"isError"`
	Msg     string `json:"msg"`
}
