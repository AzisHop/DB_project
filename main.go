package main

import (
	"DB_project/db"
	"DB_project/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//postgres, err := db.NewDb("user=postgres dbname=postgres password=admin host=127.0.0.1 port=5432 sslmode=disable")
	postgres, err := db.NewDb()

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	handler := handlers.CreateHandler(postgres.GetPs())



	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)

	server := &http.Server{
		Handler: router,
		Addr: ":5000",
	}

	log.Println("Server starting")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}