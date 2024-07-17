package handler

import (
	"html/template"
	"log"
	"net/http"

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
		log.Printf("%v on /add : not a valid request method", r.Method)
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

	if err := h.template.ExecuteTemplate(w, "add", *t); err != nil {
		log.Printf("Error while excuting template: %v", err.Error())
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		log.Print("not a htmx request")
		w.WriteHeader(http.StatusForbidden)

		return
	}

	if r.Method != http.MethodDelete {
		log.Printf("%v on /delete : not a valid request method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.PathValue("id")

	log.Printf("Delete Request for ID: %v", id)

	err := h.Service.DeleteTask(id)
	if err != nil {
		switch err.Error() {
		case "not found":
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

	log.Printf("Task Done for ID: %v", id)

	err := h.Service.MarkDone(id)
	if err != nil {
		switch err.Error() {
		case "not found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}
