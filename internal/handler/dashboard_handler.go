package handler

import (
	"finance-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashService *service.DashboardService
}

func NewDashboardHandler(dashService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashService: dashService}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	summary, err := h.dashService.GetSummary(c.Request.Context(), dateFrom, dateTo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get summary")
		return
	}
	respondSuccess(c, http.StatusOK, summary)
}

func (h *DashboardHandler) GetCategoryTotals(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	totals, err := h.dashService.GetCategoryTotals(c.Request.Context(), dateFrom, dateTo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get category totals")
		return
	}
	respondSuccess(c, http.StatusOK, totals)
}

func (h *DashboardHandler) GetTrends(c *gin.Context) {
	period := c.DefaultQuery("period", "monthly")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	trends, err := h.dashService.GetTrends(c.Request.Context(), period, dateFrom, dateTo)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get trends")
		return
	}
	respondSuccess(c, http.StatusOK, trends)
}

func (h *DashboardHandler) GetRecent(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	transactions, err := h.dashService.GetRecentTransactions(c.Request.Context(), limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get recent transactions")
		return
	}
	respondSuccess(c, http.StatusOK, transactions)
}
