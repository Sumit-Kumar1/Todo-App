package main

import (
	"log"
	"net/http"
	"time"
	"todoapp/internal/handler"
	"todoapp/internal/service"
	"todoapp/internal/store"
)

func main() {
	st, err := store.New()
	if err != nil {
		log.Printf("\nDB Creation err : %s", err.Error())
		return
	}

	if err := st.DB.Ping(); err != nil {
		log.Println("not able to ping the database: ", err.Error())
		return
	}

	defer st.DB.Close()

	log.Printf("\ndb connection success %+v", st.DB.Stats())

	svc := service.New(st)
	h := handler.New(svc)

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
