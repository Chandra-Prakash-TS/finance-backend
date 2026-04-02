package router

import (
	"finance-backend/internal/domain"
	"finance-backend/internal/handler"
	"finance-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AuthHandler        *handler.AuthHandler
	UserHandler        *handler.UserHandler
	TransactionHandler *handler.TransactionHandler
	DashboardHandler   *handler.DashboardHandler
	JWTSecret          string
	UserRepo           domain.UserRepository
}

func Setup(cfg Config) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", cfg.AuthHandler.Register)
		auth.POST("/login", cfg.AuthHandler.Login)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret, cfg.UserRepo))

	// User profile (any authenticated user)
	protected.GET("/users/me", cfg.UserHandler.GetProfile)
	protected.PUT("/users/me", cfg.UserHandler.UpdateProfile)

	// User management (admin only)
	users := protected.Group("/users")
	users.Use(middleware.RequireRole(domain.RoleAdmin))
	{
		users.GET("", cfg.UserHandler.List)
		users.GET("/:id", cfg.UserHandler.GetByID)
		users.PUT("/:id", cfg.UserHandler.Update)
		users.DELETE("/:id", cfg.UserHandler.Delete)
	}

	// Transactions
	transactions := protected.Group("/transactions")
	{
		transactions.GET("", middleware.RequireRole(domain.RoleAnalyst, domain.RoleAdmin), cfg.TransactionHandler.List)
		transactions.GET("/:id", middleware.RequireRole(domain.RoleAnalyst, domain.RoleAdmin), cfg.TransactionHandler.GetByID)
		transactions.POST("", middleware.RequireRole(domain.RoleAdmin), cfg.TransactionHandler.Create)
		transactions.PUT("/:id", middleware.RequireRole(domain.RoleAdmin), cfg.TransactionHandler.Update)
		transactions.DELETE("/:id", middleware.RequireRole(domain.RoleAdmin), cfg.TransactionHandler.Delete)
	}

	// Dashboard (all authenticated users)
	dashboard := protected.Group("/dashboard")
	dashboard.Use(middleware.RequireRole(domain.RoleViewer, domain.RoleAnalyst, domain.RoleAdmin))
	{
		dashboard.GET("/summary", cfg.DashboardHandler.GetSummary)
		dashboard.GET("/category-totals", cfg.DashboardHandler.GetCategoryTotals)
		dashboard.GET("/trends", cfg.DashboardHandler.GetTrends)
		dashboard.GET("/recent", cfg.DashboardHandler.GetRecent)
	}

	return r
}
