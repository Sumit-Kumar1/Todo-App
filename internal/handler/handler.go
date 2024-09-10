package handler

import (
	"net/http"
	"todoapp/internal/models"
)

const (
	appJSON          = "application/json"
	contentType      = "Content-Type"
	invalidReqMethod = "method not allowed"
	token            = "token"
	keyUserID        = "user_id"
)

type Handler struct {
}

func New(s Servicer) *Handler {
	tmpl := models.NewTemplate()

	return &Handler{template: tmpl, Service: s, Log: log}
}

// Root rendering endpoints
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	var tempName string

	vals := r.URL.Query()
	switch vals.Get("page") {
	case "register":
		tempName = "user-register"
	default:
		tempName = "user-login"
	}

	if err := h.template.ExecuteTemplate(w, tempName, nil); err != nil {
		h.Log.Error(err.Error(), "template-render", "index")
		return
	}
}
