package service

import (
	"context"
	"finance-backend/internal/domain"
)

type TransactionService struct {
	txnRepo domain.TransactionRepository
}

func NewTransactionService(txnRepo domain.TransactionRepository) *TransactionService {
	return &TransactionService{txnRepo: txnRepo}
}

type CreateTransactionInput struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Type     string  `json:"type" binding:"required,oneof=income expense"`
	Category string  `json:"category" binding:"required,min=1,max=100"`
	Date     string  `json:"date" binding:"required"`
	Notes    string  `json:"notes" binding:"max=1000"`
}

type UpdateTransactionInput struct {
	Amount   *float64 `json:"amount" binding:"omitempty,gt=0"`
	Type     string   `json:"type" binding:"omitempty,oneof=income expense"`
	Category string   `json:"category" binding:"omitempty,min=1,max=100"`
	Date     string   `json:"date"`
	Notes    string   `json:"notes" binding:"max=1000"`
}

func (s *TransactionService) Create(ctx context.Context, userID string, input CreateTransactionInput) (*domain.Transaction, error) {
	txn := &domain.Transaction{
		UserID:   userID,
		Amount:   input.Amount,
		Type:     input.Type,
		Category: input.Category,
		Date:     input.Date,
		Notes:    input.Notes,
	}

	if err := s.txnRepo.Create(ctx, txn); err != nil {
		return nil, err
	}
	return txn, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.txnRepo.GetByID(ctx, id)
}

func (s *TransactionService) List(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	return s.txnRepo.List(ctx, filter)
}

func (s *TransactionService) Update(ctx context.Context, id string, input UpdateTransactionInput) (*domain.Transaction, error) {
	txn, err := s.txnRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Amount != nil {
		txn.Amount = *input.Amount
	}
	if input.Type != "" {
		txn.Type = input.Type
	}
	if input.Category != "" {
		txn.Category = input.Category
	}
	if input.Date != "" {
		txn.Date = input.Date
	}
	if input.Notes != "" {
		txn.Notes = input.Notes
	}

	if err := s.txnRepo.Update(ctx, txn); err != nil {
		return nil, err
	}
	return txn, nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) error {
	return s.txnRepo.SoftDelete(ctx, id)
}
