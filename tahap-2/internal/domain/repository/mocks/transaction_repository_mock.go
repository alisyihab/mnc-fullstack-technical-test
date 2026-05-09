package mocks

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// TransactionRepository is a mock implementation of repository.TransactionRepository.
type TransactionRepository struct {
	mock.Mock
}

func (m *TransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *TransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *TransactionRepository) UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) error {
	args := m.Called(ctx, transactionID, status)
	return args.Error(0)
}
