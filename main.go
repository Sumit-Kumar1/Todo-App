package main

import (
	"log"
	"net/http"
	"time"
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

	server := http.Server{
		Addr:         ":12344",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
