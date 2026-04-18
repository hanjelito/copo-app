package model

type Booking struct {
	ID     string `json:"id"`
	RideID string `json:"ride_id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}
