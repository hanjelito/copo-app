package router

import (
	"copo/rides/internal/handler"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(rideHandler *handler.RideHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Post("/rides", rideHandler.Create)
	r.Get("/rides", rideHandler.GetAll)
	r.Get("/rides/{id}", rideHandler.GetByID)

	return r
}
