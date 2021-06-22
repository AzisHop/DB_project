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
	user.HandleFunc("/{nickname}/profile", handler.GetUser).Methods(http.MethodGet)
	user.HandleFunc("/{nickname}/profile", handler.UpdateUser).Methods(http.MethodPost)

	forum := router.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", handler.CreateForum).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/details", handler.GetForum).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/create", handler.CreateThreadForum).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/threads", handler.GetThreads).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/users", handler.GetUsersForum).Methods(http.MethodGet)

	thread := router.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", handler.CreatePostThread).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/details", handler.GetThread).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/details", handler.UpdateThread).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/posts", handler.GetThreadPosts).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/vote", handler.VoiceThread).Methods(http.MethodPost)

	post := router.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", handler.UpdatePost).Methods(http.MethodPost)
	post.HandleFunc("/{id}/details", handler.GetPost).Methods(http.MethodGet)

	service := router.PathPrefix("/service").Subrouter()
	service.HandleFunc("/status", handler.ServiceStatus).Methods(http.MethodGet)
	service.HandleFunc("/clear", handler.ServiceClear).Methods(http.MethodPost)

	server := &http.Server{
		Handler: router,
		Addr: ":5000",
	}

	log.Println("Server starting")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}