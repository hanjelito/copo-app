package model

import "time"

type Ride struct {
	ID          string    `json:"id"`
	DriverID    string    `json:"driver_id"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Departure   time.Time `json:"departure"`
	Seats       int       `json:"seats"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateRideRequest struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Departure   time.Time `json:"departure"`
	Seats       int       `json:"seats"`
}
