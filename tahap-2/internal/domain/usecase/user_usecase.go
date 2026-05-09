package usecase

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PIN         string `json:"pin" binding:"required,len=6"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PIN         string `json:"pin" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
}

type UserUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*models.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
