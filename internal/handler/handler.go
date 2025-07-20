package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"todoapp/internal/models"
)

const (
	queryPage    = "page"
	userRegister = "user-register"
	userLogin    = "user-login"
	swagger      = "swagger"
	api          = "api"
	register     = "register"
)

type UIHandler struct {
	templ *template.Template
}

func New(templ *template.Template) *UIHandler {
	templl := templ

	return &UIHandler{
		templ: templl,
	}
}

// Root rendering endpoints
func (h *UIHandler) Root(w http.ResponseWriter, r *http.Request) {
	var (
		tempName string
		ctx      = r.Context()
		logger   = models.GetLoggerFromCtx(ctx)
		vals     = r.URL.Query()
	)

	switch vals.Get(queryPage) {
	case register:
		tempName = userRegister
	case api:
		tempName = swagger
	default:
		tempName = userLogin
	}

	if err := h.templ.ExecuteTemplate(w, tempName, nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, err.Error(),
			slog.String("template-render", tempName),
		)

		models.HandleHTTPError(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *UIHandler) Swagger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	if err := h.templ.ExecuteTemplate(w, swagger, nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError,
			err.Error(), slog.String("template-render", swagger),
		)

		models.HandleHTTPError(w, err, http.StatusInternalServerError)

		return
	}
}
