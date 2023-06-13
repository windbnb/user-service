package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"github.com/windbnb/user-service/cronUtil"
	handler "github.com/windbnb/user-service/handler"
	repository "github.com/windbnb/user-service/repository"
	router "github.com/windbnb/user-service/router"
	service "github.com/windbnb/user-service/service"
	"github.com/windbnb/user-service/tracer"
	util "github.com/windbnb/user-service/util"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db := util.ConnectToDatabase()

	tracer, closer := tracer.Init("user-service")
	opentracing.SetGlobalTracer(tracer)
	router := router.ConfigureRouter(&handler.Handler{
		Tracer: tracer,
		Closer: closer,
		Service: &service.UserService{
			Repo: &repository.Repository{
				Db: db}}})

	cronUtil.ConfigureCronJobs(db)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3005"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	srv := &http.Server{Addr: "0.0.0.0:8081", Handler: c.Handler(router)}
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
