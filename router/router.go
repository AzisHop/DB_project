package router

import (
	"DB_project/db"
	"DB_project/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Routing() *mux.Router {
	postgres, err := db.NewDb()

	if err != nil {
		log.Fatal(err)
	}


	router := mux.NewRouter()
	//handler := handlers.CreateHandler(postgres.GetPs())
	userHandler := handlers.CreateUserHandler(postgres.GetPs())
	forumHandler := handlers.CreateForumHandler(postgres.GetPs())
	threadHandler := handlers.CreateThreadHandler(postgres.GetPs())
	postHandler := handlers.CreatePostHandler(postgres.GetPs())
	serviceHandler := handlers.CreateServiceHandler(postgres.GetPs())



	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{nickname}/create", userHandler.CreateUser).Methods(http.MethodPost)
	user.HandleFunc("/{nickname}/profile", userHandler.GetUser).Methods(http.MethodGet)
	user.HandleFunc("/{nickname}/profile", userHandler.UpdateUser).Methods(http.MethodPost)

	forum := router.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", forumHandler.CreateForum).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/details", forumHandler.GetForum).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/create", forumHandler.CreateThreadForum).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/threads", forumHandler.GetThreads).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/users", forumHandler.GetUsersForum).Methods(http.MethodGet)

	thread := router.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", threadHandler.CreatePostThread).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/details", threadHandler.GetThread).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/details", threadHandler.UpdateThread).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/posts", threadHandler.GetThreadPosts).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/vote", threadHandler.VoiceThread).Methods(http.MethodPost)

	post := router.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", postHandler.UpdatePost).Methods(http.MethodPost)
	post.HandleFunc("/{id}/details", postHandler.GetPost).Methods(http.MethodGet)

	service := router.PathPrefix("/service").Subrouter()
	service.HandleFunc("/status", serviceHandler.ServiceStatus).Methods(http.MethodGet)
	service.HandleFunc("/clear", serviceHandler.ServiceClear).Methods(http.MethodPost)

	return router
}
