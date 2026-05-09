package repository

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error
}
