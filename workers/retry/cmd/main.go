package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"retry/internal/consumer"
	"retry/internal/repository"
	"retry/internal/service"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
)

func main() {
	godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to connect to Postgres: %v", err)
	}
	defer db.Close()

	kw := &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_URL")),
		Balancer: &kafka.LeastBytes{},
	}
	defer kw.Close()

	retryRepo := repository.NewRepository(db)
	retrySvc := service.NewRetryService(retryRepo, kw)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Printf("Shutting down retry worker...")
		cancel()
	}()

	consumer.Start(ctx, retrySvc)
}
