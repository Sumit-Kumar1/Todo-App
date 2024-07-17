package handler

import "todoapp/models"

type Servicer interface {
	AddTask(task, descr string) (*models.Task, error)
	DeleteTask(id string) error
	MarkDone(id string) error
}
