package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "DEBIT"
	TransactionTypeCredit TransactionType = "CREDIT"
)

type TransactionCategory string

const (
	CategoryTopUp    TransactionCategory = "TOPUP"
	CategoryPayment  TransactionCategory = "PAYMENT"
	CategoryTransfer TransactionCategory = "TRANSFER"
)

type Transaction struct {
	ID              uuid.UUID           `gorm:"type:uuid;primaryKey" json:"transaction_id"`
	UserID          uuid.UUID           `gorm:"type:uuid;not null" json:"user_id"`
	TargetUserID    *uuid.UUID          `gorm:"type:uuid" json:"target_user_id,omitempty"`
	Amount          float64             `gorm:"type:decimal(15,2);not null" json:"amount"`
	Remarks         string              `gorm:"type:text" json:"remarks"`
	BalanceBefore   float64             `gorm:"type:decimal(15,2);not null" json:"balance_before"`
	BalanceAfter    float64             `gorm:"type:decimal(15,2);not null" json:"balance_after"`
	Status          string              `gorm:"type:varchar(20);default:'SUCCESS'" json:"status"`
	TransactionType TransactionType     `gorm:"type:varchar(10);not null" json:"transaction_type"`
	Category        TransactionCategory `gorm:"type:varchar(20);not null" json:"-"`
	CreatedAt       time.Time           `json:"created_date"`
	UpdatedAt       time.Time           `json:"updated_date"`
}

// TransactionResponse is used for custom JSON output as per PDF
type TransactionResponse struct {
	ID              uuid.UUID       `json:"id"` // Will be mapped to transfer_id, top_up_id, etc.
	Status          string          `json:"status"`
	UserID          uuid.UUID       `json:"user_id"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          float64         `json:"amount"`
	Remarks         string          `json:"remarks"`
	BalanceBefore   float64         `json:"balance_before"`
	BalanceAfter    float64         `json:"balance_after"`
	CreatedDate     string          `json:"created_date"`
}
