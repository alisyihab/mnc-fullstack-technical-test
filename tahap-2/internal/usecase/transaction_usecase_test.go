package usecase_test

import (
	"context"
	"errors"
	"testing"

	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/repository/mocks"
	domainUsecase "mnc-fullstack-technical-test/tahap-2/internal/domain/usecase"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/worker"
	"mnc-fullstack-technical-test/tahap-2/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// noopCreditFn satisfies worker.CreditFn without hitting a database.
func noopCreditFn(_ context.Context, _ uuid.UUID, _ float64, _ uuid.UUID, _ string) error {
	return nil
}

func newTestTxUsecase(txRepo *mocks.TransactionRepository, userRepo *mocks.UserRepository) domainUsecase.TransactionUsecase {
	tw := worker.NewTransferWorker(1, 10, noopCreditFn)
	return usecase.NewTransactionUsecase(txRepo, userRepo, tw)
}

// ---------------------------------------------------------------------------
// TopUp
// ---------------------------------------------------------------------------

func TestTopUp_Success(t *testing.T) {
	txRepo := new(mocks.TransactionRepository)
	userRepo := new(mocks.UserRepository)
	uc := newTestTxUsecase(txRepo, userRepo)

	userID := uuid.New()
	existingUser := &models.User{ID: userID, Balance: 100.0}

	userRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)
	txRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)
	userRepo.On("UpdateBalance", mock.Anything, userID, 50.0).Return(nil)

	req := &domainUsecase.TopUpRequest{Amount: 50.0}
	tx, err := uc.TopUp(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 100.0, tx.BalanceBefore)
	assert.Equal(t, 150.0, tx.BalanceAfter)
	assert.Equal(t, "SUCCESS", tx.Status)
	userRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Payment
// ---------------------------------------------------------------------------

func TestPayment_Success(t *testing.T) {
	txRepo := new(mocks.TransactionRepository)
	userRepo := new(mocks.UserRepository)
	uc := newTestTxUsecase(txRepo, userRepo)

	userID := uuid.New()
	existingUser := &models.User{ID: userID, Balance: 200.0}

	userRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)
	txRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)
	userRepo.On("UpdateBalance", mock.Anything, userID, -75.0).Return(nil)

	req := &domainUsecase.PaymentRequest{Amount: 75.0, Remarks: "test payment"}
	tx, err := uc.Payment(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 200.0, tx.BalanceBefore)
	assert.Equal(t, 125.0, tx.BalanceAfter)
	assert.Equal(t, "SUCCESS", tx.Status)
	userRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestPayment_InsufficientBalance(t *testing.T) {
	txRepo := new(mocks.TransactionRepository)
	userRepo := new(mocks.UserRepository)
	uc := newTestTxUsecase(txRepo, userRepo)

	userID := uuid.New()
	existingUser := &models.User{ID: userID, Balance: 10.0}

	userRepo.On("FindByID", mock.Anything, userID).Return(existingUser, nil)

	req := &domainUsecase.PaymentRequest{Amount: 500.0, Remarks: "too much"}
	tx, err := uc.Payment(context.Background(), userID, req)

	assert.Nil(t, tx)
	assert.EqualError(t, err, "Balance is not enough")
	userRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Transfer
// ---------------------------------------------------------------------------

func TestTransfer_Success(t *testing.T) {
	txRepo := new(mocks.TransactionRepository)
	userRepo := new(mocks.UserRepository)
	uc := newTestTxUsecase(txRepo, userRepo)

	senderID := uuid.New()
	receiverID := uuid.New()
	sender := &models.User{ID: senderID, Balance: 300.0}
	receiver := &models.User{ID: receiverID, Balance: 50.0}

	userRepo.On("FindByID", mock.Anything, senderID).Return(sender, nil)
	userRepo.On("FindByID", mock.Anything, receiverID).Return(receiver, nil)
	txRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)
	userRepo.On("UpdateBalance", mock.Anything, senderID, -100.0).Return(nil)

	req := &domainUsecase.TransferRequest{
		TargetUser: receiverID,
		Amount:     100.0,
		Remarks:    "test transfer",
	}
	tx, err := uc.Transfer(context.Background(), senderID, req)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "PENDING", tx.Status)
	assert.Equal(t, 300.0, tx.BalanceBefore)
	assert.Equal(t, 200.0, tx.BalanceAfter)
	userRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestTransfer_InsufficientBalance(t *testing.T) {
	txRepo := new(mocks.TransactionRepository)
	userRepo := new(mocks.UserRepository)
	uc := newTestTxUsecase(txRepo, userRepo)

	senderID := uuid.New()
	sender := &models.User{ID: senderID, Balance: 50.0}

	userRepo.On("FindByID", mock.Anything, senderID).Return(sender, nil)

	req := &domainUsecase.TransferRequest{
		TargetUser: uuid.New(),
		Amount:     100.0,
	}
	tx, err := uc.Transfer(context.Background(), senderID, req)

	assert.Nil(t, tx)
	assert.EqualError(t, err, "Balance is not enough")
	userRepo.AssertExpectations(t)
}

func TestTransfer_TargetUserNotFound(t *testing.T) {
	txRepo := new(mocks.TransactionRepository)
	userRepo := new(mocks.UserRepository)
	uc := newTestTxUsecase(txRepo, userRepo)

	senderID := uuid.New()
	receiverID := uuid.New()
	sender := &models.User{ID: senderID, Balance: 500.0}

	userRepo.On("FindByID", mock.Anything, senderID).Return(sender, nil)
	userRepo.On("FindByID", mock.Anything, receiverID).Return(nil, errors.New("record not found"))

	req := &domainUsecase.TransferRequest{
		TargetUser: receiverID,
		Amount:     100.0,
	}
	tx, err := uc.Transfer(context.Background(), senderID, req)

	assert.Nil(t, tx)
	assert.EqualError(t, err, "Target user not found")
	userRepo.AssertExpectations(t)
}
