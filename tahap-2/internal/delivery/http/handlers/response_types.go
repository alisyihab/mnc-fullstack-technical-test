package handlers

import "github.com/google/uuid"

type registerResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	CreatedDate string    `json:"created_date"`
}

type updateProfileResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Address     string    `json:"address"`
	UpdatedDate string    `json:"updated_date"`
}

type topUpResponse struct {
	TopUpID       uuid.UUID `json:"top_up_id"`
	AmountTopUp   float64   `json:"amount_top_up"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedDate   string    `json:"created_date"`
}

type paymentResponse struct {
	PaymentID     uuid.UUID `json:"payment_id"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedDate   string    `json:"created_date"`
}

type transferResponse struct {
	TransferID    uuid.UUID `json:"transfer_id"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedDate   string    `json:"created_date"`
}

// transactionItemResponse uses omitempty on ID fields so only the relevant
// one appears (top_up_id / payment_id / transfer_id) while preserving order.
type transactionItemResponse struct {
	TopUpID         *uuid.UUID `json:"top_up_id,omitempty"`
	PaymentID       *uuid.UUID `json:"payment_id,omitempty"`
	TransferID      *uuid.UUID `json:"transfer_id,omitempty"`
	Status          string     `json:"status"`
	UserID          uuid.UUID  `json:"user_id"`
	TransactionType string     `json:"transaction_type"`
	Amount          float64    `json:"amount"`
	Remarks         string     `json:"remarks"`
	BalanceBefore   float64    `json:"balance_before"`
	BalanceAfter    float64    `json:"balance_after"`
	CreatedDate     string     `json:"created_date"`
}
