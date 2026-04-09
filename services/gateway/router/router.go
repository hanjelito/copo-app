package router

import (
	"copo/gateway/internal/middleware"
	"copo/gateway/internal/proxy"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(middleware.RelateLimit())

	authURL := os.Getenv("AUTH_SERVICE_URL")
	userURL := os.Getenv("USER_SERVICE_URL")
	ridesURL := os.Getenv("RIDES_SERVICE_URL")
	bookingsURL := os.Getenv("BOOKINGS_SERVICE_URL")

	//public routes
	// Rutas públicas
	r.Post("/auth/register", proxy.To(authURL))
	r.Post("/auth/login", proxy.To(authURL))
	r.Post("/auth/refresh", proxy.To(authURL))

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/auth/me", proxy.To(authURL))
		r.Get("/users/me", proxy.To(userURL))
		r.Put("/users/me", proxy.To(userURL))
		r.Post("/rides", proxy.To(ridesURL))
		r.Get("/rides", proxy.To(ridesURL))
		r.Get("/rides/{id}", proxy.To(ridesURL))
		r.Post("/bookings", proxy.To(bookingsURL))
		r.Get("/bookings", proxy.To(bookingsURL))
		r.Delete("/bookings/{id}", proxy.To(bookingsURL))
	})

	return r

}
