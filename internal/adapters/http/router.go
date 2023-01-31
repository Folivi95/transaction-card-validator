package http

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/saltpay/transaction-card-validator/internal/adapters/http/handlers"
	middleware2 "github.com/saltpay/transaction-card-validator/internal/adapters/http/middleware"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware2.Prometheus)
	router.Use(middleware2.ContextMiddleware)

	monitoring(router)

	return router
}

func monitoring(router *mux.Router) {
	router.HandleFunc("/internal/health_check", handlers.HealthCheck)
	router.Handle("/internal/metrics", promhttp.Handler())
}
