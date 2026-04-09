package repository

import (
	"context"
	"copo/rides/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RideRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *RideRepository {
	return &RideRepository{db: db}
}

func (r *RideRepository) Create(ctx context.Context, ride *model.Ride) (*model.Ride, error) {
	result := &model.Ride{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO rides (driver_id, origin, destination, departure, seats)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, driver_id, origin, destination, departure, seats, created_at
	`,
		ride.DriverID,
		ride.Origin,
		ride.Destination,
		ride.Departure,
		ride.Seats,
	).
		Scan(
			&result.ID,
			&result.DriverID,
			&result.Origin,
			&result.Destination,
			&result.Departure,
			&result.Seats,
			&result.CreatedAt,
		)
	return result, err
}

func (r *RideRepository) FindAll(ctx context.Context) ([]*model.Ride, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, driver_id, origin, destination, departure, seats, created_at
		FROM rides ORDER BY departure ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*model.Ride
	for rows.Next() {
		ride := &model.Ride{}
		if err := rows.Scan(
			&ride.ID,
			&ride.DriverID,
			&ride.Origin,
			&ride.Destination,
			&ride.Departure,
			&ride.Seats,
			&ride.CreatedAt,
		); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *RideRepository) FindByID(ctx context.Context, id string) (*model.Ride, error) {
	ride := &model.Ride{}
	err := r.db.QueryRow(ctx, `
		SELECT id, driver_id, origin, destination, departure, seats, created_at
		FROM rides WHERE id = $1
		`, id).Scan(
		&ride.ID,
		&ride.DriverID,
		&ride.Origin,
		&ride.Destination,
		&ride.Departure,
		&ride.Seats,
		&ride.CreatedAt,
	)
	return ride, err
}
