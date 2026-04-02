package postgres

import (
	"context"
	"database/sql"
	"finance-backend/internal/domain"
	"fmt"
)

type DashboardRepo struct {
	db *sql.DB
}

func NewDashboardRepo(db *sql.DB) *DashboardRepo {
	return &DashboardRepo{db: db}
}

func (r *DashboardRepo) GetSummary(ctx context.Context, dateFrom, dateTo string) (float64, float64, error) {
	where := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if dateFrom != "" {
		where += fmt.Sprintf(" AND date >= $%d", argIdx)
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		where += fmt.Sprintf(" AND date <= $%d", argIdx)
		args = append(args, dateTo)
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expenses
		FROM transactions %s`, where)

	var totalIncome, totalExpenses float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalIncome, &totalExpenses)
	return totalIncome, totalExpenses, err
}

func (r *DashboardRepo) GetCategoryTotals(ctx context.Context, dateFrom, dateTo string) ([]domain.CategoryTotal, error) {
	where := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if dateFrom != "" {
		where += fmt.Sprintf(" AND date >= $%d", argIdx)
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		where += fmt.Sprintf(" AND date <= $%d", argIdx)
		args = append(args, dateTo)
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT category, type, SUM(amount) as total
		FROM transactions %s
		GROUP BY category, type
		ORDER BY total DESC`, where)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totals []domain.CategoryTotal
	for rows.Next() {
		var ct domain.CategoryTotal
		if err := rows.Scan(&ct.Category, &ct.Type, &ct.TotalAmount); err != nil {
			return nil, err
		}
		totals = append(totals, ct)
	}
	return totals, rows.Err()
}

func (r *DashboardRepo) GetTrends(ctx context.Context, period, dateFrom, dateTo string) ([]domain.TrendPoint, error) {
	truncPeriod := "month"
	if period == "weekly" {
		truncPeriod = "week"
	}

	where := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if dateFrom != "" {
		where += fmt.Sprintf(" AND date >= $%d", argIdx)
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		where += fmt.Sprintf(" AND date <= $%d", argIdx)
		args = append(args, dateTo)
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT
			TO_CHAR(date_trunc('%s', date), 'YYYY-MM-DD') as period,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
		FROM transactions %s
		GROUP BY date_trunc('%s', date)
		ORDER BY period`, truncPeriod, where, truncPeriod)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []domain.TrendPoint
	for rows.Next() {
		var tp domain.TrendPoint
		if err := rows.Scan(&tp.Period, &tp.Income, &tp.Expenses); err != nil {
			return nil, err
		}
		trends = append(trends, tp)
	}
	return trends, rows.Err()
}

func (r *DashboardRepo) GetRecentTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, user_id, amount, type, category, date, notes, created_at, updated_at
		FROM transactions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Amount, &t.Type,
			&t.Category, &t.Date, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}
