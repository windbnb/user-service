package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/service"
	"github.com/windbnb/user-service/tracer"
)

type Handler struct {
	Service *service.UserService
	Tracer opentracing.Tracer
	Closer io.Closer
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("loginHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling login at %s\n", r.URL.Path)),
	)

	var credentials model.Credentials
	json.NewDecoder(r.Body).Decode(&credentials)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	token, err := handler.Service.Login(credentials, ctx)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	json.NewEncoder(w).Encode(model.LoginResponse{Token: token})
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("registerHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling registration at %s\n", r.URL.Path)),
	)

	var userDTO model.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&userDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	createdUser, err := handler.Service.CreateUser(userDTO.ToUser(), ctx)

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
	span := tracer.StartSpanFromRequest("editUserHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling user edit at %s\n", r.URL.Path)),
	)

	params := mux.Vars(r)
	userId, _ := strconv.ParseUint(params["id"], 10, 32)

	var userDTO model.UserDTO
	json.NewDecoder(r.Body).Decode(&userDTO)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	err := handler.Service.EditUser(userDTO, userId, ctx)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) AuthoriseGuest(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("authoriseGuestHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling guest authorisation at %s\n", r.URL.Path)),
	)

	cookie := r.Header.Values("Authorization")
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(cookie[0], " ")[1]
	
	ctx := tracer.ContextWithSpan(context.Background(), span)
	err := handler.Service.AuthoriseUser(tokenString, model.GUEST, ctx)
	
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (handler *Handler) AuthoriseHost(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("authoriseHostHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling host authorisation at %s\n", r.URL.Path)),
	)

	cookie := r.Header.Values("Authorization")
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(cookie[0], " ")[1]
	
	ctx := tracer.ContextWithSpan(context.Background(), span)
	err := handler.Service.AuthoriseUser(tokenString, model.HOST, ctx)
	
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (handler *Handler) FindUser(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("findUserHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling finding user at %s\n", r.URL.Path)),
	)

	params := mux.Vars(r)
	userId, _ := strconv.ParseUint(params["id"], 10, 32)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	user, err := handler.Service.FindUser(userId, ctx)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	json.NewEncoder(w).Encode(user.ToDTO())
}

func (handler *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("deleteUserHandler", handler.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling user deletion at %s\n", r.URL.Path)),
	)

	params := mux.Vars(r)
	userId, _ := strconv.ParseUint(params["id"], 10, 32)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	err := handler.Service.DeleteUser(userId, ctx)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unreachable") {
			status = http.StatusGatewayTimeout
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: status})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
