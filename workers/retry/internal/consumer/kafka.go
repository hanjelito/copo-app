package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"retry/internal/model"
	"retry/internal/service"

	kafka "github.com/segmentio/kafka-go"
)

func Start(ctx context.Context, svc *service.RetryService) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_URL")},
		Topic:       "booking.created",
		GroupID:     "notification-worker",
		StartOffset: kafka.FirstOffset,
	})
	defer r.Close()

	log.Println("Retry worker listening booking.pending...")

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("error reading message: %v", err)
			continue
		}
		var event model.Booking
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("error deserializing message: %v", err)
			continue
		}
		if _, err := svc.Insert(ctx, &event); err != nil {
			log.Printf("booking %s failed after all retries, dropping: %v", event.ID, err)
		}
	}
}
