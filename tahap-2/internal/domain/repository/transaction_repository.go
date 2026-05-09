package repository

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) error
}
