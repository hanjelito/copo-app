package repository

import (
	"context"
	"copo/bookings/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, b *model.Booking) (*model.Booking, error) {
	result := &model.Booking{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO bookings (ride_id, user_id)
		VALUES ($1, $2)
		RETURNING id, ride_id, user_id, status, created_at
	`,
		b.RideID,
		b.UserID,
	).Scan(
		&result.ID,
		&result.RideID,
		&result.UserID,
		&result.Status,
		&result.CreatedAt,
	)
	return result, err
}

func (r *BookingRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Booking, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, ride_id, user_id, status, created_at
		FROM bookings
		WHERE user_id = $1 ORDER BY  created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		b := &model.Booking{}
		if err := rows.Scan(
			&b.ID,
			&b.RideID,
			&b.UserID,
			&b.Status,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingRepository) FindByID(ctx context.Context, id string) (*model.Booking, error) {
	b := &model.Booking{}
	err := r.db.QueryRow(ctx, `
		SELECT id, ride_id, user_id, status, created_at
		FROM bookings WHERE id = $1
	`, id).Scan(
		&b.ID,
		&b.RideID,
		&b.UserID,
		&b.Status,
		&b.CreatedAt,
	)
	return b, err
}

func (r *BookingRepository) Cancel(ctx context.Context, id, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE bookings SET status = 'cancelled'
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	return err
}
