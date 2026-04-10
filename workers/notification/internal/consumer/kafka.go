package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

type BookingEvent struct {
	ID     string `json:"id"`
	RideID string `json:"ride_id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

func Start(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_URL")},
		Topic:       "booking.created",
		GroupID:     "notification-worker",
		StartOffset: kafka.FirstOffset,
	})
	defer r.Close()

	log.Println("Notification worker listen booking.created...")

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("error read message: %v", err)
			continue
		}

		var event BookingEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("error deserializing message: %v", err)
			continue
		}

		log.Printf("Notificación: reserva %s confirmada para el usuario %s en el viaje %s",
			event.ID, event.UserID, event.RideID)
	}
}
