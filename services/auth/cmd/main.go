package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"copo/auth/internal/handler"
	"copo/auth/internal/repository"
	"copo/auth/internal/router"
	"copo/auth/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("no se pudo conectar a postgres: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewRepository(db)
	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	r := router.New(authHandler, userHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Auth service corriendo en %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
