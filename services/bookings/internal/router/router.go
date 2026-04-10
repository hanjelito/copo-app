package router

import (
	"copo/bookings/internal/handler"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(bookingHandler *handler.BookingHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Post("/bookings", bookingHandler.Create)
	r.Get("/bookings/me", bookingHandler.GetMyBookings)
	r.Delete("/bookings/{id}", bookingHandler.Cancel)

	return r
}
