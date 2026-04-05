package handler

import (
	stderrors "errors"
	"net/http"

	"todoapp/internal/errors"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HandleError maps a domain or HTTP error to the appropriate status code
// and writes a JSON error response.
func HandleError(c *gin.Context, err error) {
	var httpErr *errors.HTTPError

	resp := errorResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}

	switch {
	case stderrors.As(err, &httpErr):
		resp.Code = httpErr.Code
		resp.Details = httpErr.Details
	case stderrors.Is(err, errors.ErrUserAlreadyExists):
		resp.Code = http.StatusConflict
	case stderrors.Is(err, errors.ErrUserNotFound), stderrors.Is(err, errors.ErrTaskNotFound):
		resp.Code = http.StatusNotFound
	case stderrors.Is(err, errors.ErrInvalidCookie), stderrors.Is(err, errors.ErrPsswdNotMatch):
		resp.Code = http.StatusUnauthorized
	case stderrors.Is(err, errors.ErrInvalidTaskID),
		stderrors.Is(err, errors.ErrRequired),
		stderrors.Is(err, errors.ErrInvalid),
		stderrors.Is(err, errors.ErrDueDatePast):
		resp.Code = http.StatusBadRequest
	}

	c.AbortWithStatusJSON(resp.Code, resp)
}
