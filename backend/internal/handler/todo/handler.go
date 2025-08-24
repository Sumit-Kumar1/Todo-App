package todohttp

import (
	"todoapp/internal/constant"
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

const (
	templateAddTask = "add"
	templateIndex   = "index"
	renderErr       = "error while rendering template"
	hxRedirect      = "HX-Redirect"
)

type Handler struct {
	Service TodoServicer
	// template *template.Template
}

func New(todoSvc TodoServicer) *Handler {
	return &Handler{Service: todoSvc}
}

func (h *Handler) AddTask(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	t := models.TaskReq{
		Title:       ctx.Value("title").(string),
		Description: ctx.Value("description").(string),
		DueDate:     ctx.Value("dueDate").(string),
	}

	task, err := h.Service.AddTask(ctx, &t, &userID)
	if err != nil {
		return nil, err
	}

	//if err := h.template.ExecuteTemplate(ctx, templateAddTask, task.ToTaskResp()); err != nil {
	//	ctx.Logger.Errorf(renderErr, slog.String("template", templateAddTask))
	//	return nil, err
	//}

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

	//if err := h.template.ExecuteTemplate(w, templateAddTask, resp.ToTaskResp()); err != nil {
	//	ctx.Logger.Errorf("error while rendering template: %v", err)
	//
	//	return nil, err
	//}

	return resp, nil
}

// nolint:revive // this is a handler get not returning
func (h *Handler) GetAllTasks(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		//_ = h.template.ExecuteTemplate(w, "errorPage", map[string]any{
		//	"Code":    http.StatusUnauthorized,
		//	"Message": "user not authorized!!",
		//})

		return nil, errors.ErrInvalidCookie
	}

	resp, err := h.Service.GetAll(ctx, &userID)
	if err != nil {
		if errors.ErrNotFound("user").Error() == err.Error() {
			//w.Header().Add(hxRedirect, "/?page=register")
			//w.WriteHeader(http.StatusOK)

			return nil, errors.ErrUserNotFound
		}

		//logger.LogAttrs(ctx, slog.LevelError, err.Error(), slog.String("user", userID.String()))
		//errors.HandleHTTPError(w, err, http.StatusInternalServerError)

		return nil, err
	}

	//w.WriteHeader(http.StatusOK)

	var tasks = make([]models.TaskResp, 0)

	for i := range resp {
		tasks = append(tasks, *resp[i].ToTaskResp())
	}

	//if err := h.template.ExecuteTemplate(w, templateIndex, resp); err != nil {
	//	logger.LogAttrs(ctx, slog.LevelError, renderErr, slog.String("template", templateIndex))
	//	errors.HandleHTTPError(w, err, http.StatusInternalServerError)
	//}

	return tasks, nil
}

func (h *Handler) DeleteTask(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrInvalid("user")
	}

	id := ctx.PathParam("id")

	if err := h.Service.DeleteTask(ctx, id, &userID); err != nil {
		//switch {
		//case errors.ErrNotFound("user").Error() == err.Error():
		//	errors.HandleHTTPError(w, err, http.StatusNotFound)
		//default:
		//	errors.HandleHTTPError(w, err, http.StatusBadRequest)
		//}
		return nil, err
	}

	return nil, nil
}

func (h *Handler) Update(ctx *gofr.Context) (any, error) {
	userID, ok := ctx.Value(constant.CtxKeyUserID).(uuid.UUID)
	if !ok {
		return nil, errors.ErrInvalid("user")
	}

	t := models.TaskReq{
		ID:          ctx.Value("id").(string),
		Title:       ctx.Value("title").(string),
		Description: ctx.Value("description").(string),
		DueDate:     ctx.Value("dueDate").(string),
	}

	resp, err := h.Service.UpdateTask(ctx, t.ID, &t, false, &userID)
	if err != nil {
		return nil, err
	}

	//if err := h.template.ExecuteTemplate(w, templateAddTask, *resp.ToTaskResp()); err != nil {
	//	logger.LogAttrs(ctx, slog.LevelError, renderErr, slog.String("template", templateAddTask))
	//	errors.HandleHTTPError(w, err, http.StatusInternalServerError)
	//
	//	return nil, nil
	//}

	return resp, nil
}
