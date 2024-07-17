package handler

import (
	"log"
	"net/http"

	"html/template"

	"github.com/angelofallars/htmx-go"
)

type Handler struct {
	template *template.Template
	Service  Servicer
}

func New(s Servicer) *Handler {
	tmpl := template.Must(template.ParseFiles("html/index.html"))
	return &Handler{template: tmpl, Service: s}
}

func (h *Handler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("not a valid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	err := h.template.ExecuteTemplate(w, "index", nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		log.Print("not a htmx request")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if r.Method != http.MethodPost {
		log.Print("not a valid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	task := r.PostFormValue("task")
	taskDesc := r.PostFormValue("desc")

	t, err := h.Service.AddTask(task, taskDesc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	if err := h.template.ExecuteTemplate(w, "add", t); err != nil {
		log.Printf("Error while excuting template: %v", err.Error())
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if r.Method != http.MethodDelete {
		log.Print("not a valid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")

	log.Printf("\nDelete Request for ID: %v", id)

	if err := h.Service.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	htmx.NewResponse().StatusCode(http.StatusAccepted).Redirect("/").Refresh(true)
}

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPut {
		log.Print("not a valid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")

	log.Printf("\nTask Done for ID: %v", id)

	if err := h.Service.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	_ = htmx.NewResponse().StatusCode(http.StatusOK).Redirect("/").Refresh(true).Write(w)
}
