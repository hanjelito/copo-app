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
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
		log.Printf("Attemp %d/5: Unable to connect to Postgres, retrying... %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if db == nil {
		log.Fatal("Unable to connect to Postgres after 5 attempts")
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("unable to connect to Redis: %v", err)
	}

	rideRepo := repository.NewRepository(db)
	rideSvc := service.NewRideService(rideRepo, rdb)
	rideHandler := handler.NewRideHandler(rideSvc)

	r := router.New(rideHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Ride service running on:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}

	log.Println("Rides service stopped")
}
