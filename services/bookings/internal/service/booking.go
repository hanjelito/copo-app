package service

import (
	"context"
	"copo/bookings/internal/model"
	"copo/bookings/internal/repository"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type BookingService struct {
	repo  *repository.BookingRepository
	redis *redis.Client
	kafka *kafka.Writer
}

func NewBookingService(repo *repository.BookingRepository, redis *redis.Client, kafka *kafka.Writer) *BookingService {
	return &BookingService{repo: repo, redis: redis, kafka: kafka}
}

func (s *BookingService) Create(ctx context.Context, userID string, req *model.CreateBookingRequest) (*model.Booking, error) {
	// Try Redis: Check if a cache is available
	cacheKey := fmt.Sprintf("ride:%s:seats", req.RideID)
	seats, err := s.redis.Get(ctx, cacheKey).Int()
	if err == nil && seats <= 0 {
		return nil, fmt.Errorf("there are no seats available")
	}

	booking := &model.Booking{
		RideID: req.RideID,
		UserID: userID,
	}
	result, err := s.repo.Create(ctx, booking)
	if err != nil {
		// PG failed — publish booking.pending to the retry worker
		s.publishEvent(ctx, "booking.pending", booking)
		return nil, err
	}
	// Refresh the seat cache in Redis
	s.redis.Decr(ctx, cacheKey)

	// Publish booking.created to Kafka
	s.publishEvent(ctx, "booking.created", result)

	return result, nil
}

func (s *BookingService) GetMyBookings(ctx context.Context, userID string) ([]*model.Booking, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *BookingService) Cancel(ctx context.Context, id, userID string) error {
	//Get the booking to find out the ride_id
	err := s.repo.Cancel(ctx, id, userID)
	if err != nil {
		return err
	}
	// return the seats to Redis
	booking, err := s.repo.FindByID(ctx, id)
	if err == nil {
		cacheKey := fmt.Sprintf("ride:%s:seats", booking.RideID)
		s.redis.Incr(ctx, cacheKey)
	}
	return nil
}

func (s *BookingService) publishEvent(ctx context.Context, topic string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	s.kafka.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: data,
	})
}
