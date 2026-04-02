package service

import (
	"context"
	"finance-backend/internal/domain"
	"testing"
	"time"
)

type mockTxnRepo struct {
	txns map[string]*domain.Transaction
}

func newMockTxnRepo() *mockTxnRepo {
	return &mockTxnRepo{txns: make(map[string]*domain.Transaction)}
}

func (m *mockTxnRepo) Create(ctx context.Context, txn *domain.Transaction) error {
	txn.ID = "txn-" + txn.Category + "-" + txn.Date
	txn.CreatedAt = time.Now()
	txn.UpdatedAt = time.Now()
	m.txns[txn.ID] = txn
	return nil
}

func (m *mockTxnRepo) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if txn, ok := m.txns[id]; ok {
		return txn, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockTxnRepo) List(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	var result []domain.Transaction
	for _, t := range m.txns {
		if filter.Type != "" && t.Type != filter.Type {
			continue
		}
		if filter.Category != "" && t.Category != filter.Category {
			continue
		}
		result = append(result, *t)
	}
	return result, int64(len(result)), nil
}

func (m *mockTxnRepo) Update(ctx context.Context, txn *domain.Transaction) error {
	if _, ok := m.txns[txn.ID]; !ok {
		return domain.ErrNotFound
	}
	m.txns[txn.ID] = txn
	return nil
}

func (m *mockTxnRepo) SoftDelete(ctx context.Context, id string) error {
	if _, ok := m.txns[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.txns, id)
	return nil
}

func TestCreateTransaction_Success(t *testing.T) {
	repo := newMockTxnRepo()
	svc := NewTransactionService(repo)

	txn, err := svc.Create(context.Background(), "user-1", CreateTransactionInput{
		Amount:   1500.50,
		Type:     "income",
		Category: "Salary",
		Date:     "2024-01-15",
		Notes:    "Monthly salary",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txn.Amount != 1500.50 {
		t.Fatalf("expected amount 1500.50, got %f", txn.Amount)
	}
	if txn.UserID != "user-1" {
		t.Fatalf("expected user_id user-1, got %s", txn.UserID)
	}
}

func TestGetTransaction_NotFound(t *testing.T) {
	repo := newMockTxnRepo()
	svc := NewTransactionService(repo)

	_, err := svc.GetByID(context.Background(), "nonexistent")
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateTransaction_Success(t *testing.T) {
	repo := newMockTxnRepo()
	svc := NewTransactionService(repo)

	txn, _ := svc.Create(context.Background(), "user-1", CreateTransactionInput{
		Amount:   100,
		Type:     "expense",
		Category: "Food",
		Date:     "2024-01-15",
	})

	newAmount := 200.0
	updated, err := svc.Update(context.Background(), txn.ID, UpdateTransactionInput{
		Amount:   &newAmount,
		Category: "Groceries",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Amount != 200.0 {
		t.Fatalf("expected amount 200, got %f", updated.Amount)
	}
	if updated.Category != "Groceries" {
		t.Fatalf("expected category Groceries, got %s", updated.Category)
	}
}

func TestDeleteTransaction_Success(t *testing.T) {
	repo := newMockTxnRepo()
	svc := NewTransactionService(repo)

	txn, _ := svc.Create(context.Background(), "user-1", CreateTransactionInput{
		Amount:   100,
		Type:     "expense",
		Category: "Food",
		Date:     "2024-01-15",
	})

	if err := svc.Delete(context.Background(), txn.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := svc.GetByID(context.Background(), txn.ID)
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	repo := newMockTxnRepo()
	svc := NewTransactionService(repo)

	err := svc.Delete(context.Background(), "nonexistent")
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListTransactions_FilterByType(t *testing.T) {
	repo := newMockTxnRepo()
	svc := NewTransactionService(repo)

	svc.Create(context.Background(), "user-1", CreateTransactionInput{
		Amount: 100, Type: "income", Category: "Salary", Date: "2024-01-15",
	})
	svc.Create(context.Background(), "user-1", CreateTransactionInput{
		Amount: 50, Type: "expense", Category: "Food", Date: "2024-01-16",
	})

	txns, total, err := svc.List(context.Background(), domain.TransactionFilter{
		Type:     "income",
		Page:     1,
		PageSize: 20,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 result, got %d", total)
	}
	if txns[0].Type != "income" {
		t.Fatalf("expected income type, got %s", txns[0].Type)
	}
}
