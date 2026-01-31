package todohttp

import (
	"net/http"

	"todoapp/internal/errors"
	"todoapp/internal/handler"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service TodoServicer
}

func New(todoSvc TodoServicer) *Handler {
	return &Handler{Service: todoSvc}
}

func (h *Handler) AddTask(c *gin.Context) {
	var taskReq models.TaskReq

	userID, err := handler.GetContextKey(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := c.BindJSON(&taskReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	task, err := h.Service.AddTask(c, &taskReq, userID)
	if err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, task.ToTaskResp())
}

func (h *Handler) MarkDone(c *gin.Context) {
	userID, err := handler.GetContextKey(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	id := c.Param("id")

	resp, err := h.Service.MarkDone(c, id, userID)
	if err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, resp.ToTaskResp())
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	userID, err := handler.GetContextKey(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	resp, err := h.Service.GetAll(c, userID)
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

	userID, err := handler.GetContextKey(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.Service.DeleteTask(c, id, userID); err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Update(c *gin.Context) {
	var taskReq models.TaskReq

	userID, err := handler.GetContextKey(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	id := c.Param("id")

	if err := c.BindJSON(&taskReq); err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	resp, err := h.Service.UpdateTask(c, id, &taskReq, userID)
	if err != nil {
		errors.HandleHTTPError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, resp.ToTaskResp())
}

func (h *Handler) GetChildTasks(c *gin.Context) {
	userID, err := handler.GetContextKey(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	id := c.Param("id")

	resp, err := h.Service.GetChildTasks(c, id, userID)
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
