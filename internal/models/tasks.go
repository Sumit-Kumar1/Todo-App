package models

import "time"

type Task struct {
	ID         string
	Title      string
	IsDone     bool
	AddedAt    time.Time
	ModifiedAt *time.Time
}
