package handler

import (
	"html/template"
	"log"
	"net/http"
	"todoapp/internal/models"
)

const (
	hxRequest        = "Hx-Request"
	hxPromtp         = "HX-Prompt"
	trueStr          = "true"
	invalidReqMethod = "%s method not allowed on %s"
	notHTMX          = "not a htmx request"
)

type Handler struct {
	template *template.Template
	Service  Servicer
}

func New(s Servicer) *Handler {
	templateDir := "views/index.html"

	tmpl := template.Must(template.ParseFiles(templateDir))
	return &Handler{template: tmpl, Service: s}
}

func (h *Handler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf(invalidReqMethod, r.Method, "/")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	tasks, err := h.Service.GetAll(ctx)
	if err != nil {
		return
	}

	err = h.template.ExecuteTemplate(w, "index", map[string][]models.Task{
		"Data": tasks,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		log.Print(notHTMX)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if r.Method != http.MethodPost {
		log.Printf(invalidReqMethod, r.Method, "/add")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	task := r.PostFormValue("task")
	ctx := r.Context()

	t, err := h.Service.AddTask(ctx, task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	log.Printf("Task is Added ID : %s", t.ID)

	if err := h.template.ExecuteTemplate(w, "add", *t); err != nil {
		log.Printf("Error while excuting template: %v", err.Error())
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		log.Print(notHTMX)
		w.WriteHeader(http.StatusForbidden)

		return
	}

	if r.Method != http.MethodDelete {
		log.Printf(invalidReqMethod, r.Method, "/delete")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")
	ctx := r.Context()

	log.Printf("Delete Request for ID: %v", id)

	err := h.Service.DeleteTask(ctx, id)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPut {
		log.Printf(invalidReqMethod, r.Method, "/done")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")
	ctx := r.Context()

	log.Printf("Task Done for ID: %v", id)

	resp, err := h.Service.MarkDone(ctx, id)
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if err := h.template.ExecuteTemplate(w, "add", *resp); err != nil {
		log.Printf("error in /done/%s\n\tError:%s", id, err.Error())
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(hxRequest) != trueStr {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if r.Method != http.MethodPost {
		log.Printf(invalidReqMethod, r.Method, "/update")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")
	title := r.Header.Get("HX-Prompt")
	ctx := r.Context()

	resp, err := h.Service.UpdateTask(ctx, id, title, "false")
	if err != nil {
		switch {
		case models.ErrNotFound.Is(err):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if err := h.template.ExecuteTemplate(w, "add", *resp); err != nil {
		log.Printf("error in /update/%s\n\tError:%s", id, err.Error())
		return
	}

	log.Printf("Task Updated : %s", id)
}
