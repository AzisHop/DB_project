package main

import (
	"DB_project/router"
	"log"
	"net/http"
)

func main() {
	myRouter := router.Routing()

	server := &http.Server{
		Handler: myRouter,
		Addr: ":5000",
	}

	log.Println("Server starting")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}