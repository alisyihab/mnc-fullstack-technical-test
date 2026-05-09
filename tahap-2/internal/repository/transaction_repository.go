package repository

import (
	"context"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("id = ?", transactionID).
		Update("status", status).Error
}
