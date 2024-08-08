package main

import (
	"log"
	"net/http"
	"todoapp/internal/handler"
	"todoapp/internal/server"
	"todoapp/internal/service"
	"todoapp/internal/store"
)

func main() {
	st, err := store.New()
	if err != nil {
		log.Printf("\nDB Creation err : %s", err.Error())
		return
	}

	if err = st.DB.Ping(); err != nil {
		log.Println("database not reachable: ", err.Error())
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
	http.HandleFunc("/health", healthStatus)

	app := server.NewServer(
		server.WithAppName("todoApp"),
		server.WithEnv("development"),
		server.WithPort("12344"),
	)

	log.Printf("Server created with configs:App-Name: %s, Port: %s, env: %s", app.Name, app.Addr, app.Env)
	log.Printf("\nApplication %v server is started on port:%v", app.Name, app.Addr)

	err = app.ListenAndServe()
	if err != nil {
		log.Println("error while running server : ", err.Error())

		return
	}
}

func healthStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
