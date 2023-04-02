package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/windbnb/user-service/handler"
)

func HandleRequests(handler *handler.Handler) {
	router := mux.NewRouter()
	router.HandleFunc("/api/users/login", handler.Login).Methods("POST")
	router.HandleFunc("/api/users/register", handler.Register).Methods("POST")

	router.HandleFunc("/api/users/authorize/guest", handler.AuthoriseGuest).Methods("POST")
	router.HandleFunc("/api/users/authorize/host", handler.AuthoriseHost).Methods("POST")
	
	router.HandleFunc("/api/users/{id}", handler.FindUser).Methods("GET")
	router.HandleFunc("/api/users/{id}", handler.EditUser).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8081", router))
}
