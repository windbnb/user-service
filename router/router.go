package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/user-service/handler"
	"github.com/windbnb/user-service/metrics"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/users/login", metrics.MetricProxy(handler.Login)).Methods("POST")
	router.HandleFunc("/api/users/register", metrics.MetricProxy(handler.Register)).Methods("POST")

	router.HandleFunc("/api/users/authorize/guest", metrics.MetricProxy(handler.AuthoriseGuest)).Methods("POST")
	router.HandleFunc("/api/users/authorize/host", metrics.MetricProxy(handler.AuthoriseHost)).Methods("POST")

	router.HandleFunc("/api/users/{id}", metrics.MetricProxy(handler.FindUser)).Methods("GET")
	router.HandleFunc("/api/users/{id}", metrics.MetricProxy(handler.EditUser)).Methods("PUT")
	router.HandleFunc("/api/users/change-password/{id}", metrics.MetricProxy(handler.ChangePassword)).Methods("PUT")
	router.HandleFunc("/api/users/{id}", metrics.MetricProxy(handler.DeleteUser)).Methods("DELETE")

	router.Path("/metrics").Handler(metrics.MetricsHandler())

	router.HandleFunc("/probe/liveness", handler.Healthcheck)
	router.HandleFunc("/probe/readiness", handler.Ready)

	return router
}
