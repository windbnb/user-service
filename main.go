package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/windbnb/user-service/cronUtil"
	handler "github.com/windbnb/user-service/handler"
	repository "github.com/windbnb/user-service/repository"
	router "github.com/windbnb/user-service/router"
	service "github.com/windbnb/user-service/service"
	util "github.com/windbnb/user-service/util"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db := util.ConnectToDatabase()
	router := router.ConfigureRouter(&handler.Handler{Service: &service.UserService{ Repo: &repository.Repository{Db: db}}})
	cronUtil.ConfigureCronJobs(db)

	srv := &http.Server{Addr: "localhost:8081", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit
	
	defer db.Close()
	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
