package models

import "github.com/google/uuid"

type Task struct {
	ID     uuid.UUID
	TagID  string
	Title  string
	Desc   string
	IsDone bool
}
