package handler

import "todoapp/models"

type Servicer interface {
	AddTask(task, descr string) (*models.Task, error)
	DeleteTask(id string) error
	UpdateTask(id, title, desc, isDone string) (*models.Task, error)
	MarkDone(id string) (*models.Task, error)
}
