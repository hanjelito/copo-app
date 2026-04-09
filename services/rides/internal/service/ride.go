package service

import (
	"context"
	"copo/rides/internal/model"
	"copo/rides/internal/repository"
)

type RideService struct {
	repo *repository.RideRepository
}

func NewRideService(repo *repository.RideRepository) *RideService {
	return &RideService{repo: repo}
}

func (s *RideService) Create(ctx context.Context, driverID string, req *model.CreateRideRequest) (*model.Ride, error) {
	ride := &model.Ride{
		DriverID:    driverID,
		Origin:      req.Origin,
		Destination: req.Destination,
		Departure:   req.Departure,
		Seats:       req.Seats,
	}
	return s.repo.Create(ctx, ride)
}

func (s *RideService) GetAll(ctx context.Context) ([]*model.Ride, error) {
	return s.repo.FindAll(ctx)
}

func (s *RideService) GetByID(ctx context.Context, id string) (*model.Ride, error) {
	return s.repo.FindByID(ctx, id)
}
