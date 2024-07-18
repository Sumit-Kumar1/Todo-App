package main

import (
	"log"
	"net/http"

	"todoapp/handler"
	"todoapp/service"
)

func main() {
	s := service.New()
	h := handler.New(s)

	http.HandleFunc("/", h.IndexPage)
	http.HandleFunc("/add", h.AddTask)
	http.HandleFunc("/delete/{id}", h.DeleteTask)
	http.HandleFunc("/update/{id}", h.Update)
	http.HandleFunc("/done/{id}", h.Done)

	log.Fatal(http.ListenAndServe(":12344", nil))
}
