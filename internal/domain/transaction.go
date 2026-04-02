package domain

import (
	"context"
	"time"
)

const (
	TypeIncome  = "income"
	TypeExpense = "expense"
)

type Transaction struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Amount    float64    `json:"amount"`
	Type      string     `json:"type"`
	Category  string     `json:"category"`
	Date      string     `json:"date"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type TransactionFilter struct {
	Type      string
	Category  string
	DateFrom  string
	DateTo    string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

type TransactionRepository interface {
	Create(ctx context.Context, txn *Transaction) error
	GetByID(ctx context.Context, id string) (*Transaction, error)
	List(ctx context.Context, filter TransactionFilter) ([]Transaction, int64, error)
	Update(ctx context.Context, txn *Transaction) error
	SoftDelete(ctx context.Context, id string) error
}

type CategoryTotal struct {
	Category     string  `json:"category"`
	Type         string  `json:"type"`
	TotalAmount  float64 `json:"total_amount"`
}

type TrendPoint struct {
	Period   string  `json:"period"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
}

type DashboardRepository interface {
	GetSummary(ctx context.Context, dateFrom, dateTo string) (totalIncome, totalExpenses float64, err error)
	GetCategoryTotals(ctx context.Context, dateFrom, dateTo string) ([]CategoryTotal, error)
	GetTrends(ctx context.Context, period, dateFrom, dateTo string) ([]TrendPoint, error)
	GetRecentTransactions(ctx context.Context, limit int) ([]Transaction, error)
}
