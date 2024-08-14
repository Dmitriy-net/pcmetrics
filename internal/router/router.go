package router

import (
	"github.com/Dmitriy-net/pcmetrics/internal/handlers"
	"github.com/Dmitriy-net/pcmetrics/internal/repository"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(repo repository.Repository) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(repo))
	r.Get("/value/{type}/{name}", handlers.GetValueHandler(repo))
	r.Get("/", handlers.ListMetricsHandler(repo))

	return r
}
