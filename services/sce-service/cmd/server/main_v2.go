// OmniRoute Ecosystem - Backend API Server
// High-performance Go backend with participant-based multi-tenancy

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omniroute/backend/internal/config"
	"github.com/omniroute/backend/internal/handlers"
	"github.com/omniroute/backend/internal/middleware"
	"github.com/omniroute/backend/internal/repository"
	"github.com/omniroute/backend/internal/services"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize database connection
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize Redis for caching
	cache, err := repository.NewRedisCache(cfg.Redis)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	svc := services.NewServices(repos, cache, cfg)

	// Initialize handlers
	h := handlers.NewHandlers(svc)

	// Setup router
	router := setupRouter(h, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("starting server", "port", cfg.Server.Port, "env", cfg.Server.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exited")
}

func setupRouter(h *handlers.Handlers, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes
	// Auth routes
	mux.HandleFunc("POST /api/v1/auth/register", h.Auth.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Auth.Login)
	mux.HandleFunc("POST /api/v1/auth/verify-otp", h.Auth.VerifyOTP)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.Auth.RefreshToken)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", h.Auth.ForgotPassword)
	mux.HandleFunc("GET /api/v1/auth/me", middleware.Auth(cfg.JWT.Secret)(h.Auth.GetCurrentUser))

	// Onboarding routes
	mux.HandleFunc("GET /api/v1/onboarding/participant-types", h.Onboarding.GetParticipantTypes)
	mux.HandleFunc("POST /api/v1/onboarding/kyc", middleware.Auth(cfg.JWT.Secret)(h.Onboarding.SubmitKYC))
	mux.HandleFunc("POST /api/v1/onboarding/documents", middleware.Auth(cfg.JWT.Secret)(h.Onboarding.UploadDocument))
	mux.HandleFunc("GET /api/v1/onboarding/kyc/status", middleware.Auth(cfg.JWT.Secret)(h.Onboarding.GetKYCStatus))

	// Dashboard routes
	mux.HandleFunc("GET /api/v1/dashboard/stats", middleware.Auth(cfg.JWT.Secret)(h.Dashboard.GetStats))
	mux.HandleFunc("GET /api/v1/dashboard/notifications", middleware.Auth(cfg.JWT.Secret)(h.Dashboard.GetNotifications))
	mux.HandleFunc("GET /api/v1/dashboard/activities", middleware.Auth(cfg.JWT.Secret)(h.Dashboard.GetActivities))

	// Orders routes
	mux.HandleFunc("GET /api/v1/orders", middleware.Auth(cfg.JWT.Secret)(h.Orders.List))
	mux.HandleFunc("POST /api/v1/orders", middleware.Auth(cfg.JWT.Secret)(h.Orders.Create))
	mux.HandleFunc("GET /api/v1/orders/{id}", middleware.Auth(cfg.JWT.Secret)(h.Orders.GetByID))
	mux.HandleFunc("PUT /api/v1/orders/{id}/status", middleware.Auth(cfg.JWT.Secret)(h.Orders.UpdateStatus))

	// Products routes
	mux.HandleFunc("GET /api/v1/products", middleware.Auth(cfg.JWT.Secret)(h.Products.List))
	mux.HandleFunc("POST /api/v1/products", middleware.Auth(cfg.JWT.Secret)(h.Products.Create))
	mux.HandleFunc("GET /api/v1/products/{id}", middleware.Auth(cfg.JWT.Secret)(h.Products.GetByID))
	mux.HandleFunc("PUT /api/v1/products/{id}", middleware.Auth(cfg.JWT.Secret)(h.Products.Update))
	mux.HandleFunc("DELETE /api/v1/products/{id}", middleware.Auth(cfg.JWT.Secret)(h.Products.Delete))
	mux.HandleFunc("GET /api/v1/products/categories", h.Products.GetCategories)

	// Inventory routes
	mux.HandleFunc("GET /api/v1/inventory", middleware.Auth(cfg.JWT.Secret)(h.Inventory.List))
	mux.HandleFunc("GET /api/v1/inventory/stock-levels", middleware.Auth(cfg.JWT.Secret)(h.Inventory.GetStockLevels))
	mux.HandleFunc("GET /api/v1/inventory/alerts", middleware.Auth(cfg.JWT.Secret)(h.Inventory.GetAlerts))
	mux.HandleFunc("POST /api/v1/inventory/adjust", middleware.Auth(cfg.JWT.Secret)(h.Inventory.AdjustStock))

	// Logistics routes
	mux.HandleFunc("GET /api/v1/logistics/deliveries", middleware.Auth(cfg.JWT.Secret)(h.Logistics.ListDeliveries))
	mux.HandleFunc("GET /api/v1/logistics/deliveries/{id}", middleware.Auth(cfg.JWT.Secret)(h.Logistics.GetDelivery))
	mux.HandleFunc("GET /api/v1/logistics/tracking/{id}", middleware.Auth(cfg.JWT.Secret)(h.Logistics.TrackDelivery))
	mux.HandleFunc("GET /api/v1/logistics/routes", middleware.Auth(cfg.JWT.Secret)(h.Logistics.GetRoutes))
	mux.HandleFunc("GET /api/v1/logistics/fleet", middleware.Auth(cfg.JWT.Secret)(h.Logistics.GetFleet))

	// Finance routes
	mux.HandleFunc("GET /api/v1/finance/wallet", middleware.Auth(cfg.JWT.Secret)(h.Finance.GetWallet))
	mux.HandleFunc("GET /api/v1/finance/transactions", middleware.Auth(cfg.JWT.Secret)(h.Finance.GetTransactions))
	mux.HandleFunc("POST /api/v1/finance/withdraw", middleware.Auth(cfg.JWT.Secret)(h.Finance.Withdraw))
	mux.HandleFunc("GET /api/v1/finance/settlements", middleware.Auth(cfg.JWT.Secret)(h.Finance.GetSettlements))
	mux.HandleFunc("GET /api/v1/finance/loans", middleware.Auth(cfg.JWT.Secret)(h.Finance.GetLoans))
	mux.HandleFunc("POST /api/v1/finance/loans/apply", middleware.Auth(cfg.JWT.Secret)(h.Finance.ApplyForLoan))

	// Analytics routes
	mux.HandleFunc("GET /api/v1/analytics/sales", middleware.Auth(cfg.JWT.Secret)(h.Analytics.GetSalesAnalytics))
	mux.HandleFunc("GET /api/v1/analytics/inventory", middleware.Auth(cfg.JWT.Secret)(h.Analytics.GetInventoryAnalytics))
	mux.HandleFunc("GET /api/v1/analytics/performance", middleware.Auth(cfg.JWT.Secret)(h.Analytics.GetPerformanceMetrics))

	// Apply global middleware
	handler := middleware.Chain(
		mux,
		middleware.Logger,
		middleware.CORS(cfg.Server.AllowedOrigins),
		middleware.RateLimiter(100, time.Minute),
		middleware.RequestID,
		middleware.Recover,
	)

	return handler
}
