package repository

import (
	"context"
	"fmt"
	"log"
	"retry/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) RetryInsert(ctx context.Context, event *model.Booking) error {
	maxRetries := 5
	wait := time.Second

	for i := range maxRetries {
		var id string
		err := r.db.QueryRow(ctx, `
			INSERT INTO bookings (ride_id, user_id)
			VALUES ($1, $2)
			RETURNING id
		`,
			event.RideID,
			event.UserID,
		).Scan(&id)

		if err == nil {
			event.ID = id
			log.Printf("booking retry ok: %s intento %d", id, i+1)
			return nil
		}

		log.Printf("retry %d/%d failed %v, waiting %s", i+1, maxRetries, err, wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
		wait *= 2
	}
	return fmt.Errorf("max retries reaches for ride %s", event.RideID)
}
