package mocks

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserRepository is a mock implementation of repository.UserRepository.
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}
