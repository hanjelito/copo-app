package main

import (
	"context"
	"copo/rides/internal/handler"
	"copo/rides/internal/repository"
	"copo/rides/internal/router"
	"copo/rides/internal/service"
	"log"
	"net/http"
	"os"

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

	rideRepo := repository.NewRepository(db)
	rideSvc := service.NewRideService(rideRepo)
	rideHandler := handler.NewRideHandler(rideSvc)

	r := router.New(rideHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Rides service corrido en %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
