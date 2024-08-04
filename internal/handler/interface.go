package handler

import "todoapp/internal/models"

type Servicer interface {
	GetAll() []models.Task
	AddTask(task string) (*models.Task, error)
	DeleteTask(id string) error
	UpdateTask(id, title, isDone string) (*models.Task, error)
	MarkDone(id string) (*models.Task, error)
}
