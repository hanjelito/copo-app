package service

import (
	"context"
	"encoding/json"
	"retry/internal/model"
	"retry/internal/repository"

	kafka "github.com/segmentio/kafka-go"
)

type RetryService struct {
	repo  *repository.BookingRepository
	kafka *kafka.Writer
}

func NewRetryService(repo *repository.BookingRepository, kafka *kafka.Writer) *RetryService {
	return &RetryService{repo: repo, kafka: kafka}
}

func (s *RetryService) Insert(ctx context.Context, req *model.Booking) (*model.Booking, error) {
	if err := s.repo.RetryInsert(ctx, req); err != nil {
		return nil, err
	}
	publishCreated(ctx, s.kafka, req)
	return req, nil
}

func publishCreated(ctx context.Context, kw *kafka.Writer, event *model.Booking) {
	data, _ := json.Marshal(event)
	kw.WriteMessages(ctx, kafka.Message{
		Topic: "booking.created",
		Value: data,
	})
}
