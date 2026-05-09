package usecase

import (
	"context"
	"errors"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/repository"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/usecase"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/worker"

	"github.com/google/uuid"
)

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	transferWorker  *worker.TransferWorker
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	transferWorker *worker.TransferWorker,
) usecase.TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		transferWorker:  transferWorker,
	}
}

func (u *transactionUsecase) TopUp(ctx context.Context, userID uuid.UUID, req *usecase.TopUpRequest) (*models.Transaction, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore + req.Amount

	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		Amount:          req.Amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		TransactionType: models.TransactionTypeCredit,
		Category:        models.CategoryTopUp,
		Status:          "SUCCESS",
	}

	if err := u.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	if err := u.userRepo.UpdateBalance(ctx, userID, req.Amount); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *transactionUsecase) Payment(ctx context.Context, userID uuid.UUID, req *usecase.PaymentRequest) (*models.Transaction, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Balance < req.Amount {
		return nil, errors.New("Balance is not enough")
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore - req.Amount

	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		Amount:          req.Amount,
		Remarks:         req.Remarks,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		TransactionType: models.TransactionTypeDebit,
		Category:        models.CategoryPayment,
		Status:          "SUCCESS",
	}

	if err := u.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	if err := u.userRepo.UpdateBalance(ctx, userID, -req.Amount); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *transactionUsecase) Transfer(ctx context.Context, userID uuid.UUID, req *usecase.TransferRequest) (*models.Transaction, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Balance < req.Amount {
		return nil, errors.New("Balance is not enough")
	}

	_, err = u.userRepo.FindByID(ctx, req.TargetUser)
	if err != nil {
		return nil, errors.New("Target user not found")
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore - req.Amount

	// Create transaction with PENDING status — will be set to SUCCESS by worker
	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		TargetUserID:    &req.TargetUser,
		Amount:          req.Amount,
		Remarks:         req.Remarks,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		TransactionType: models.TransactionTypeDebit,
		Category:        models.CategoryTransfer,
		Status:          "PENDING",
	}

	if err := u.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Deduct sender balance synchronously
	if err := u.userRepo.UpdateBalance(ctx, userID, -req.Amount); err != nil {
		return nil, err
	}

	// Enqueue receiver credit to worker pool
	u.transferWorker.Enqueue(worker.TransferJob{
		SenderID:      userID,
		ReceiverID:    req.TargetUser,
		Amount:        req.Amount,
		TransactionID: transaction.ID,
	})

	return transaction, nil
}

func (u *transactionUsecase) GetTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	return u.transactionRepo.FindByUserID(ctx, userID)
}
