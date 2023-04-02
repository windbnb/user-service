package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/service"
)

type Handler struct {
	Service *service.UserService
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credentials
	json.NewDecoder(r.Body).Decode(&credentials)

	token, err := handler.Service.Login(credentials)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	json.NewEncoder(w).Encode(model.LoginResponse{Token: token})
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var userDTO model.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&userDTO)

	createdUser, err := handler.Service.CreateUser(userDTO.ToUser())

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser.ToDTO())
}

func (handler *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, _ := strconv.ParseUint(params["id"], 10, 32)

	var userDTO model.UserDTO
	json.NewDecoder(r.Body).Decode(&userDTO)

	err := handler.Service.EditUser(userDTO, userId)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) AuthoriseGuest(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Values("Authorization")
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(cookie[0], " ")[1]
	
	err := handler.Service.AuthoriseUser(tokenString, model.GUEST)
	
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (handler *Handler) AuthoriseHost(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Values("Authorization")
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(cookie[0], " ")[1]
	
	err := handler.Service.AuthoriseUser(tokenString, model.HOST)
	
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (handler *Handler) FindUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, _ := strconv.ParseUint(params["id"], 10, 32)

	user, err := handler.Service.FindUser(userId)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	json.NewEncoder(w).Encode(user.ToDTO())
}
