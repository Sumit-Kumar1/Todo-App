package todohttp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"todoapp/internal/constant"
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

type Handler struct {
	Service TodoServicer
}

func New(todoSvc TodoServicer) *Handler {
	return &Handler{Service: todoSvc}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	var taskReq models.TaskReq

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		logger.LogAttrs(ctx, slog.LevelError, "invalid user id")
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "close body error", slog.String("error", err.Error()))
		}
	}(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Add task: decoding body error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	task, err := h.Service.AddTask(ctx, &taskReq, &userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "service:add task error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	data, _ := json.Marshal(task)

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(data)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		logger.LogAttrs(ctx, slog.LevelError, "invalid user id")
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	id := r.PathValue("id")

	resp, err := h.Service.MarkDone(ctx, id, &userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "PATCH:mark done", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	data, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		logger.LogAttrs(ctx, slog.LevelError, "invalid user id")
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	resp, err := h.Service.GetAll(ctx, &userID)
	if err != nil {
		if errors.ErrUserNotFound.Is(err) {
			logger.LogAttrs(ctx, slog.LevelError, "get all tasks error", slog.String("error", err.Error()))
			errors.HandleHTTPError(w, errors.ErrUserNotFound)
			return
		}

		logger.LogAttrs(ctx, slog.LevelError, "get all tasks error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	var tasks = make([]models.TaskResp, 0)

	for i := range resp {
		tasks = append(tasks, *resp[i].ToTaskResp())
	}

	data, _ := json.Marshal(tasks)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		logger.LogAttrs(ctx, slog.LevelError, "invalid user id")
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	id := r.PathValue("id")

	if err := h.Service.DeleteTask(ctx, id, &userID); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "delete task error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	var taskReq models.TaskReq

	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		logger.LogAttrs(ctx, slog.LevelError, "invalid user id")
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	id := r.PathValue("id")

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "update task: decoding body error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	resp, err := h.Service.UpdateTask(ctx, id, &taskReq, false, &userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "update task error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	data, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
