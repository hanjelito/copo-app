package router

import (
	"copo/auth/internal/handler"
	authmw "copo/auth/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth)
		r.Get("/auth/me", authHandler.Me)
		r.Get("/users/me", userHandler.GetMe)
		r.Put("/users/me", userHandler.UpdateMe)
	})

	return r

}
