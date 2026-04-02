package service

import (
	"context"
	"finance-backend/internal/domain"
	"testing"
)

type mockDashRepo struct{}

func (m *mockDashRepo) GetSummary(ctx context.Context, dateFrom, dateTo string) (float64, float64, error) {
	return 5000.0, 2000.0, nil
}

func (m *mockDashRepo) GetCategoryTotals(ctx context.Context, dateFrom, dateTo string) ([]domain.CategoryTotal, error) {
	return []domain.CategoryTotal{
		{Category: "Salary", Type: "income", TotalAmount: 5000},
		{Category: "Food", Type: "expense", TotalAmount: 1000},
		{Category: "Transport", Type: "expense", TotalAmount: 1000},
	}, nil
}

func (m *mockDashRepo) GetTrends(ctx context.Context, period, dateFrom, dateTo string) ([]domain.TrendPoint, error) {
	return []domain.TrendPoint{
		{Period: "2024-01-01", Income: 5000, Expenses: 2000},
	}, nil
}

func (m *mockDashRepo) GetRecentTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	return []domain.Transaction{
		{ID: "1", Amount: 100, Type: "expense", Category: "Food"},
	}, nil
}

func TestGetSummary(t *testing.T) {
	svc := NewDashboardService(&mockDashRepo{})

	summary, err := svc.GetSummary(context.Background(), "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if summary.TotalIncome != 5000 {
		t.Fatalf("expected income 5000, got %f", summary.TotalIncome)
	}
	if summary.TotalExpenses != 2000 {
		t.Fatalf("expected expenses 2000, got %f", summary.TotalExpenses)
	}
	if summary.NetBalance != 3000 {
		t.Fatalf("expected net balance 3000, got %f", summary.NetBalance)
	}
}

func TestGetCategoryTotals(t *testing.T) {
	svc := NewDashboardService(&mockDashRepo{})

	totals, err := svc.GetCategoryTotals(context.Background(), "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(totals) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(totals))
	}
}

func TestGetTrends_DefaultMonthly(t *testing.T) {
	svc := NewDashboardService(&mockDashRepo{})

	trends, err := svc.GetTrends(context.Background(), "", "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend point, got %d", len(trends))
	}
}

func TestGetRecentTransactions_DefaultLimit(t *testing.T) {
	svc := NewDashboardService(&mockDashRepo{})

	txns, err := svc.GetRecentTransactions(context.Background(), 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(txns) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txns))
	}
}
