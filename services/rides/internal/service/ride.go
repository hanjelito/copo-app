package service

import (
	"context"
	"copo/rides/internal/model"
	"copo/rides/internal/repository"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RideService struct {
	repo  *repository.RideRepository
	redis *redis.Client
}

func NewRideService(repo *repository.RideRepository, redis *redis.Client) *RideService {
	return &RideService{repo: repo, redis: redis}
}

func (s *RideService) Create(ctx context.Context, driverID string, req *model.CreateRideRequest) (*model.Ride, error) {
	ride := &model.Ride{
		DriverID:    driverID,
		Origin:      req.Origin,
		Destination: req.Destination,
		Departure:   req.Departure,
		Seats:       req.Seats,
	}
	result, err := s.repo.Create(ctx, ride)
	if err != nil {
		return nil, err
	}
	//Save seats in redis
	cacheKey := fmt.Sprintf("ride:%s:seats", result.ID)
	s.redis.Set(ctx, cacheKey, result.Seats, 0)

	return result, nil
}

func (s *RideService) GetAll(ctx context.Context) ([]*model.Ride, error) {
	return s.repo.FindAll(ctx)
}

func (s *RideService) GetByID(ctx context.Context, id string) (*model.Ride, error) {
	return s.repo.FindByID(ctx, id)
}
