package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"copo/auth/internal/handler"
	"copo/auth/internal/repository"
	"copo/auth/internal/router"
	"copo/auth/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var err error
	var db *pgxpool.Pool
	for i := range 5 {
		db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err == nil {
			break
		}
		log.Printf("Attempt %d/5: Unable to connect to Postgres, retrying... %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if db == nil {
		log.Fatal("Unable to connect to Postgress after 5 attempts")
	}
	defer db.Close()

	//Repository
	userRepo := repository.NewRepository(db)

	//Service
	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)

	//Handler
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
