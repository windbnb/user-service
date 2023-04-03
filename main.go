package main

import (
	handler "github.com/windbnb/user-service/handler"
	repository "github.com/windbnb/user-service/repository"
	router "github.com/windbnb/user-service/router"
	service "github.com/windbnb/user-service/service"
	util "github.com/windbnb/user-service/util"
)

func main() {
	db := util.ConnectToDatabase()
	router.HandleRequests(&handler.Handler{Service: &service.UserService{ Repo: &repository.Repository{Db: db}}})

	defer db.Close()
}