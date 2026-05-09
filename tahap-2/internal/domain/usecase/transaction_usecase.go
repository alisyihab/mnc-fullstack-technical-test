package usecase

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"

	"github.com/google/uuid"
)

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type PaymentRequest struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Remarks string  `json:"remarks"`
}

type TransferRequest struct {
	TargetUser uuid.UUID `json:"target_user" binding:"required"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Remarks    string    `json:"remarks"`
}

type TransactionUsecase interface {
	TopUp(ctx context.Context, userID uuid.UUID, req *TopUpRequest) (*models.Transaction, error)
	Payment(ctx context.Context, userID uuid.UUID, req *PaymentRequest) (*models.Transaction, error)
	Transfer(ctx context.Context, userID uuid.UUID, req *TransferRequest) (*models.Transaction, error)
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
}
