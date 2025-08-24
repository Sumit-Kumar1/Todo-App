package todohttp

import (
	"todoapp/internal/constant"
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	Service TodoServicer
}

func New(todoSvc TodoServicer) *Handler {
	return &Handler{Service: todoSvc}
}

func (h *Handler) AddTask(ctx *gofr.Context) (any, error) {
	var taskReq models.TaskReq

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	if err := ctx.Bind(&taskReq); err != nil {
		return nil, err
	}

	task, err := h.Service.AddTask(ctx, &taskReq, &userID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (h *Handler) Done(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	id := ctx.PathParam("id")

	resp, err := h.Service.MarkDone(ctx, id, &userID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetAllTasks(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrInvalidCookie
	}

	resp, err := h.Service.GetAll(ctx, &userID)
	if err != nil {
		if errors.ErrUserNotFound.Is(err) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	var tasks = make([]models.TaskResp, 0)

	for i := range resp {
		tasks = append(tasks, *resp[i].ToTaskResp())
	}

	return tasks, nil
}

func (h *Handler) DeleteTask(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrInvalid("user")
	}

	id := ctx.PathParam("id")

	if err := h.Service.DeleteTask(ctx, id, &userID); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) Update(ctx *gofr.Context) (any, error) {
	var taskReq models.TaskReq

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrInvalid("user")
	}

	id := ctx.PathParam("id")

	if err := ctx.Bind(&taskReq); err != nil {
		return nil, err
	}

	resp, err := h.Service.UpdateTask(ctx, id, &taskReq, false, &userID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
