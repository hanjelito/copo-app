package service

import (
	"context"
	"copo/auth/internal/model"
	"copo/auth/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	if req.Role == "" {
		req.Role = "passenger"
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  req.Role,
	}

	id, err := s.repo.Create(ctx, user, string(hashed))
	if err != nil {
		return nil, err
	}

	return s.generateTokens(id, req.Email, req.Role)
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("credenciales invalidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales invalidas")
	}
	return s.generateTokens(user.ID, user.Email, user.Role)
}

func (s *AuthService) RefreshToken(_ context.Context, req *model.RefreshRequest) (*model.AuthResponse, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("refresh token invalido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token invalido")
	}

	id := claims["id"].(string)
	email := claims["email"].(string)
	role := claims["role"].(string)

	return s.generateTokens(id, email, role)
}

func (s *AuthService) generateTokens(id, email, role string) (*model.AuthResponse, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	access, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	refresh, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
