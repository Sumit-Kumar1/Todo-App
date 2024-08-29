package models

import "time"

type Task struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	IsDone     bool       `json:"isDone"`
	AddedAt    time.Time  `json:"addedAt"`
	ModifiedAt *time.Time `json:"modifiedAt"`
}

type Error struct {
	Type    string `json:"type"`
	IsError bool   `json:"isError"`
	Msg     string `json:"msg"`
}
