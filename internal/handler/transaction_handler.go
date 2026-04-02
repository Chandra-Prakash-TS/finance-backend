package handler

import (
	"errors"
	"finance-backend/internal/domain"
	"finance-backend/internal/service"
	"finance-backend/pkg/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	txnService *service.TransactionService
}

func NewTransactionHandler(txnService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txnService: txnService}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var input service.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondValidationError(c, []FieldError{{Field: "body", Message: err.Error()}})
		return
	}

	userID := c.GetString("user_id")
	txn, err := h.txnService.Create(c.Request.Context(), userID, input)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create transaction")
		return
	}

	respondSuccess(c, http.StatusCreated, txn)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	txn, err := h.txnService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(c, http.StatusNotFound, "NOT_FOUND", "Transaction not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get transaction")
		return
	}
	respondSuccess(c, http.StatusOK, txn)
}

func (h *TransactionHandler) List(c *gin.Context) {
	var params pagination.Params
	if err := c.ShouldBindQuery(&params); err != nil {
		respondValidationError(c, []FieldError{{Field: "query", Message: err.Error()}})
		return
	}
	params.Defaults()

	filter := domain.TransactionFilter{
		Type:      c.Query("type"),
		Category:  c.Query("category"),
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
		Page:      params.Page,
		PageSize:  params.PageSize,
		SortBy:    c.DefaultQuery("sort_by", "date"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	transactions, total, err := h.txnService.List(c.Request.Context(), filter)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list transactions")
		return
	}

	respondPaginated(c, http.StatusOK, transactions, pagination.NewMeta(params.Page, params.PageSize, total))
}

func (h *TransactionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input service.UpdateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondValidationError(c, []FieldError{{Field: "body", Message: err.Error()}})
		return
	}

	txn, err := h.txnService.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(c, http.StatusNotFound, "NOT_FOUND", "Transaction not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update transaction")
		return
	}
	respondSuccess(c, http.StatusOK, txn)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.txnService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(c, http.StatusNotFound, "NOT_FOUND", "Transaction not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete transaction")
		return
	}
	respondSuccess(c, http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
