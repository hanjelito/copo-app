package main

import (
	"context"
	"copo/bookings/internal/handler"
	"copo/bookings/internal/repository"
	"copo/bookings/internal/router"
	"copo/bookings/internal/service"
	"log"
	"net/http"
	"os"
	"time"

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
		if err != nil {
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
	log.Printf("Bookings service run in :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
