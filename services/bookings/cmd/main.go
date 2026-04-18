package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"copo/bookings/internal/handler"
	"copo/bookings/internal/repository"
	"copo/bookings/internal/router"
	"copo/bookings/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	kafka "github.com/segmentio/kafka-go"
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
		log.Fatalf("Unable to connect to Postgres after 5 attempts")
	}
	defer db.Close()

	//redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("unable to connect to Redis: %v", err)
	}

	//kafka
	kw := &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_URL")),
		Balancer: &kafka.LeastBytes{},
	}
	defer kw.Close()

	bookingRepo := repository.NewRepository(db)
	bookingSvc := service.NewBookingService(bookingRepo, rdb, kw)
	bookingHandler := handler.NewBookingHandler(bookingSvc)

	r := router.New(bookingHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Booking service running on %s", port)
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

	log.Println("Booking service stopped")
}
