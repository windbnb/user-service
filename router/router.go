package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/user-service/handler"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/users/login", handler.Login).Methods("POST")
	router.HandleFunc("/api/users/register", handler.Register).Methods("POST")

	router.HandleFunc("/api/users/authorize/guest", handler.AuthoriseGuest).Methods("POST")
	router.HandleFunc("/api/users/authorize/host", handler.AuthoriseHost).Methods("POST")
	
	router.HandleFunc("/api/users/{id}", handler.FindUser).Methods("GET")
	router.HandleFunc("/api/users/{id}", handler.EditUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", handler.DeleteUser).Methods("DELETE")

	return router
}
