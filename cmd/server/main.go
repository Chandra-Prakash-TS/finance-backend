package main

import (
	"context"
	"database/sql"
	"finance-backend/internal/config"
	"finance-backend/internal/handler"
	"finance-backend/internal/repository/postgres"
	"finance-backend/internal/router"
	"finance-backend/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Run migrations
	if err := postgres.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied")

	// Initialize repositories
	userRepo := postgres.NewUserRepo(db)
	txnRepo := postgres.NewTransactionRepo(db)
	dashRepo := postgres.NewDashboardRepo(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTTTL)
	userService := service.NewUserService(userRepo)
	txnService := service.NewTransactionService(txnRepo)
	dashService := service.NewDashboardService(dashRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	txnHandler := handler.NewTransactionHandler(txnService)
	dashHandler := handler.NewDashboardHandler(dashService)

	// Setup router
	r := router.Setup(router.Config{
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		TransactionHandler: txnHandler,
		DashboardHandler:   dashHandler,
		JWTSecret:          cfg.JWTSecret,
		UserRepo:           userRepo,
	})

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
