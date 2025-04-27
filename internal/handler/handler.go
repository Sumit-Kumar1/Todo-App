package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"todoapp/internal/models"
)

type UIHandler struct {
	templ *template.Template
}

func New() *UIHandler {
	templl := models.NewTemplate()

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

	switch vals.Get("page") {
	case "register":
		tempName = "user-register"
	case "api":
		tempName = "swagger"
	default:
		tempName = "user-login"
	}

	if err := h.templ.ExecuteTemplate(w, tempName, nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, err.Error(),
			slog.String("template-render", tempName),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *UIHandler) Swagger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := models.GetLoggerFromCtx(ctx)

	if err := h.templ.ExecuteTemplate(w, "swagger", nil); err != nil {
		logger.LogAttrs(ctx, slog.LevelError,
			"error while rendering template",
			slog.String("template", "swagger"),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
