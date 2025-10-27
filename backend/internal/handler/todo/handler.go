package todohttp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"todoapp/internal/errors"
	"todoapp/internal/handler"
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

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
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
		logger.LogAttrs(ctx, slog.LevelError, "add task: decoding body error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	task, err := h.Service.AddTask(ctx, &taskReq, &userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "add task:service error",
			slog.String("error", err.Error()))
		errors.HandleHTTPError(w, errors.ErrUserNotFound)
		return
	}

	if err := handler.WriteResponse(w, http.StatusCreated, task.ToTaskResp()); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "add-task:error while writing response",
			slog.String("err", err.Error()))
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "task added successfully", slog.String("task id", task.ID))
}

func (h *Handler) MarkDone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
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

	if err := handler.WriteResponse(w, http.StatusOK, resp.ToTaskResp()); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "mark done:error while writing response",
			slog.String("err", err.Error()))
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "mark done task successfully", slog.String("task id", resp.ID))
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
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

	if err := handler.WriteResponse(w, http.StatusOK, &tasks); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "get-all task:error while writing response",
			slog.String("err", err.Error()))
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "fetched all task successfully")
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	id := r.PathValue("id")

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		logger.LogAttrs(ctx, slog.LevelError, "invalid user id")
		errors.HandleHTTPError(w, errors.ErrInvalidCookie)
		return
	}

	if err := h.Service.DeleteTask(ctx, id, &userID); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "delete task error", slog.String("error", err.Error()))

		errors.HandleHTTPError(w, err)
		return
	}

	if err := handler.WriteResponse[any](w, http.StatusNoContent, nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "delete task:error while writing response",
			slog.String("err", err.Error()))
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "delete task successfully", slog.String("task id", id))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	var taskReq models.TaskReq

	userID, ok := ctx.Value(models.CtxKeyUserID).(uuid.UUID)
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

	resp, err := h.Service.UpdateTask(ctx, id, &taskReq, &userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "update task error", slog.String("error", err.Error()))
		errors.HandleHTTPError(w, err)
		return
	}

	if err := handler.WriteResponse(w, http.StatusOK, resp.ToTaskResp()); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "mark done:error while writing response",
			slog.String("err", err.Error()))
		return
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "updated task successfully", slog.String("task id", resp.ID))
}
