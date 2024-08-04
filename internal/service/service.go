package service

import (
	"context"
	"log"
	"strings"
	"todoapp/internal/models"
)

type Service struct {
	Store Storer
}

func New(st Storer) *Service {
	return &Service{Store: st}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.Store.GetAll(ctx)
	if err != nil {
		log.Println("error in getAll : ", err.Error())

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(ctx context.Context, title string) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, models.ErrInvalidTitle
	}

	id := generateID()

	task, err := s.Store.Create(ctx, id, title)
	if err != nil {
		log.Println("error in add task : ", err.Error())

		return nil, err
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	if err := validateID(id); err != nil {
		return models.ErrInvalidID
	}

	if err := s.Store.Delete(ctx, id); err != nil {
		log.Println("error in delete task : ", err.Error())

		return err
	}

	return nil
}

func (s *Service) MarkDone(ctx context.Context, id string) (*models.Task, error) {
	if err := validateID(id); err != nil {
		return nil, models.ErrInvalidID
	}

	task, err := s.Store.MarkDone(ctx, id)
	if err != nil {
		log.Println("error in markdone : ", err.Error())

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, id, title, isDone string) (*models.Task, error) {
	if err := validateTask(id, title, isDone); err != nil {
		return nil, err
	}

	task, err := s.Store.Update(ctx, id, title, isDone)
	if err != nil {
		log.Println("error in updating task : ", err.Error())

		return nil, err
	}

	log.Printf("\nUpdated task: %+v", task)

	return task, nil
}
