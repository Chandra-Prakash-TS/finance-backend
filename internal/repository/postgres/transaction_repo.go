package postgres

import (
	"context"
	"database/sql"
	"errors"
	"finance-backend/internal/domain"
	"fmt"
	"time"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, txn *domain.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, amount, type, category, date, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		txn.UserID, txn.Amount, txn.Type, txn.Category, txn.Date, txn.Notes,
	).Scan(&txn.ID, &txn.CreatedAt, &txn.UpdatedAt)
}

func (r *TransactionRepo) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, type, category, date, notes, created_at, updated_at
		FROM transactions
		WHERE id = $1 AND deleted_at IS NULL`

	txn := &domain.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&txn.ID, &txn.UserID, &txn.Amount, &txn.Type,
		&txn.Category, &txn.Date, &txn.Notes, &txn.CreatedAt, &txn.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return txn, err
}

func (r *TransactionRepo) List(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	where := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if filter.Type != "" {
		where += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Category != "" {
		where += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.DateFrom != "" {
		where += fmt.Sprintf(" AND date >= $%d", argIdx)
		args = append(args, filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != "" {
		where += fmt.Sprintf(" AND date <= $%d", argIdx)
		args = append(args, filter.DateTo)
		argIdx++
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) FROM transactions " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Validate sort fields
	sortBy := "date"
	switch filter.SortBy {
	case "amount", "created_at", "date", "category":
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := fmt.Sprintf(
		`SELECT id, user_id, amount, type, category, date, notes, created_at, updated_at
		FROM transactions %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortBy, sortOrder, argIdx, argIdx+1,
	)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Amount, &t.Type,
			&t.Category, &t.Date, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}
	return transactions, total, rows.Err()
}

func (r *TransactionRepo) Update(ctx context.Context, txn *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, type = $2, category = $3, date = $4, notes = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL`

	txn.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		txn.Amount, txn.Type, txn.Category, txn.Date, txn.Notes, txn.UpdatedAt, txn.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TransactionRepo) SoftDelete(ctx context.Context, id string) error {
	query := `UPDATE transactions SET deleted_at = $1, updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
