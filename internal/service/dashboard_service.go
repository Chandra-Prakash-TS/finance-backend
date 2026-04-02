package service

import (
	"context"
	"finance-backend/internal/domain"
)

type DashboardService struct {
	dashRepo domain.DashboardRepository
}

func NewDashboardService(dashRepo domain.DashboardRepository) *DashboardService {
	return &DashboardService{dashRepo: dashRepo}
}

type SummaryResponse struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
}

func (s *DashboardService) GetSummary(ctx context.Context, dateFrom, dateTo string) (*SummaryResponse, error) {
	income, expenses, err := s.dashRepo.GetSummary(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	return &SummaryResponse{
		TotalIncome:   income,
		TotalExpenses: expenses,
		NetBalance:    income - expenses,
	}, nil
}

func (s *DashboardService) GetCategoryTotals(ctx context.Context, dateFrom, dateTo string) ([]domain.CategoryTotal, error) {
	return s.dashRepo.GetCategoryTotals(ctx, dateFrom, dateTo)
}

func (s *DashboardService) GetTrends(ctx context.Context, period, dateFrom, dateTo string) ([]domain.TrendPoint, error) {
	if period == "" {
		period = "monthly"
	}
	return s.dashRepo.GetTrends(ctx, period, dateFrom, dateTo)
}

func (s *DashboardService) GetRecentTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.dashRepo.GetRecentTransactions(ctx, limit)
}
