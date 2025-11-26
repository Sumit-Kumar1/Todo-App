package todohttp

import (
	"net/http"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service TodoServicer
}

func New(todoSvc TodoServicer) *Handler {
	return &Handler{Service: todoSvc}
}

func (h *Handler) AddTask(c *gin.Context) {
	var taskReq models.TaskReq

	if err := c.BindJSON(&taskReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//TODO:fix the userID part from getting values from cookie
	userID := uuid.New()

	task, err := h.Service.AddTask(c, &taskReq, &userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, task.ToTaskResp())
}

func (h *Handler) MarkDone(c *gin.Context) {
	userID, ok := c.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		c.AbortWithError(http.StatusNotFound, errors.ErrUserNotFound)
		return
	}

	id := c.Param("id")

	resp, err := h.Service.MarkDone(c, id, &userID)
	if err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, resp.ToTaskResp())
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	userID, ok := c.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		c.AbortWithError(http.StatusNotFound, errors.ErrUserNotFound)
		return
	}

	resp, err := h.Service.GetAll(c, &userID)
	if err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	var tasks = make([]models.TaskResp, 0)

	for i := range resp {
		tasks = append(tasks, *resp[i].ToTaskResp())
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	userID, ok := c.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		c.AbortWithError(http.StatusNotFound, errors.ErrUserNotFound)
		return
	}

	if err := h.Service.DeleteTask(c, id, &userID); err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Update(c *gin.Context) {
	var taskReq models.TaskReq

	userID, ok := c.Value(models.CtxKeyUserID).(uuid.UUID)
	if !ok {
		c.AbortWithError(http.StatusNotFound, errors.ErrUserNotFound)
		return
	}

	id := c.Param("id")

	if err := c.BindJSON(&taskReq); err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	resp, err := h.Service.UpdateTask(c, id, &taskReq, &userID)
	if err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, resp.ToTaskResp())
}
