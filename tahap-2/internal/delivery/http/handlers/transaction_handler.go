package handlers

import (
	"mnc-fullstack-technical-test/tahap-2/internal/delivery/http/response"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(t usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: t}
}

// TopUp godoc
// @Summary Top up balance
// @Description Add balance to current user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.TopUpRequest true "Top up request"
// @Success 200 {object} response.JSONResponse
// @Failure 400 {object} response.JSONResponse
// @Router /topup [post]
func (h *TransactionHandler) TopUp(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req usecase.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := h.transactionUsecase.TopUp(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, topUpResponse{
		TopUpID:       tx.ID,
		AmountTopUp:   tx.Amount,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		CreatedDate:   tx.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// Payment godoc
// @Summary Make a payment
// @Description Deduct balance from current user for payment
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.PaymentRequest true "Payment request"
// @Success 200 {object} response.JSONResponse
// @Failure 400 {object} response.JSONResponse
// @Router /pay [post]
func (h *TransactionHandler) Payment(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req usecase.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := h.transactionUsecase.Payment(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, paymentResponse{
		PaymentID:     tx.ID,
		Amount:        tx.Amount,
		Remarks:       tx.Remarks,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		CreatedDate:   tx.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// Transfer godoc
// @Summary Transfer balance
// @Description Transfer balance from current user to another user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.TransferRequest true "Transfer request"
// @Success 200 {object} response.JSONResponse
// @Failure 400 {object} response.JSONResponse
// @Router /transfer [post]
func (h *TransactionHandler) Transfer(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req usecase.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := h.transactionUsecase.Transfer(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, transferResponse{
		TransferID:    tx.ID,
		Amount:        tx.Amount,
		Remarks:       tx.Remarks,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		CreatedDate:   tx.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// GetTransactions godoc
// @Summary Get transaction history
// @Description Get all transactions for current user
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.JSONResponse
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	txs, err := h.transactionUsecase.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res := make([]transactionItemResponse, 0, len(txs))
	for _, tx := range txs {
		res = append(res, toTransactionItem(&tx))
	}

	response.Success(c, http.StatusOK, res)
}

func toTransactionItem(tx *models.Transaction) transactionItemResponse {
	item := transactionItemResponse{
		Status:          tx.Status,
		UserID:          tx.UserID,
		TransactionType: string(tx.TransactionType),
		Amount:          tx.Amount,
		Remarks:         tx.Remarks,
		BalanceBefore:   tx.BalanceBefore,
		BalanceAfter:    tx.BalanceAfter,
		CreatedDate:     tx.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	id := tx.ID
	switch tx.Category {
	case models.CategoryTopUp:
		item.TopUpID = &id
	case models.CategoryPayment:
		item.PaymentID = &id
	case models.CategoryTransfer:
		item.TransferID = &id
	}

	return item
}
