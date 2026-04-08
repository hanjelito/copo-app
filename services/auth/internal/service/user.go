package service

import (
	"context"
	"copo/auth/internal/model"
	"copo/auth/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetMe(ctx context.Context, email string) (*model.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) UpdateMe(ctx context.Context, email string, req *model.UPdateRequest) (*model.User, error) {
	return s.repo.UpdateName(ctx, email, req.Name)
}
