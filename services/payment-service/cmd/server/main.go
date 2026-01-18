// Package main is the entry point for the payment-service
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the service configuration
type Config struct {
	ServerPort  string
	DatabaseURL string
	RedisURL    string
	LogLevel    string
	// Payment provider configs
	PaystackKey    string
	FlutterwaveKey string
	StripeKey      string
}

func main() {
	cfg := loadConfig()
	logger := initLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting payment-service",
		zap.String("port", cfg.ServerPort),
	)

	ctx := context.Background()
	dbPool, err := initDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer dbPool.Close()

	router := gin.New()
	router.Use(gin.Recovery())

	// Health endpoints
	router.GET("/health", healthHandler)
	router.GET("/ready", readinessHandler(dbPool))

	// API routes
	api := router.Group("/api/v1")
	{
		// Payments
		api.POST("/payments", createPaymentHandler)
		api.GET("/payments/:id", getPaymentHandler)
		api.GET("/payments", listPaymentsHandler)
		api.POST("/payments/:id/refund", refundPaymentHandler)

		// Webhooks
		api.POST("/webhooks/paystack", paystackWebhookHandler)
		api.POST("/webhooks/flutterwave", flutterwaveWebhookHandler)

		// Wallets
		api.GET("/wallets/:customer_id", getWalletHandler)
		api.POST("/wallets/:id/credit", creditWalletHandler)
		api.POST("/wallets/:id/debit", debitWalletHandler)
		api.GET("/wallets/:id/transactions", getWalletTransactionsHandler)

		// Credit accounts
		api.GET("/credit-accounts/:customer_id", getCreditAccountHandler)
		api.POST("/credit-accounts", createCreditAccountHandler)
		api.POST("/credit-accounts/:id/use", useCreditLineHandler)
		api.POST("/credit-accounts/:id/repay", repayCreditLineHandler)
		api.GET("/credit-accounts/:id/statement", getCreditStatementHandler)

		// Credit scoring
		api.GET("/credit-scores/:customer_id", getCreditScoreHandler)
		api.POST("/credit-scores/:customer_id/calculate", calculateCreditScoreHandler)

		// Disbursements
		api.POST("/disbursements", createDisbursementHandler)
		api.GET("/disbursements/:id", getDisbursementHandler)
		api.GET("/disbursements", listDisbursementsHandler)
		api.POST("/disbursements/batch", createDisbursementBatchHandler)

		// Bank accounts
		api.GET("/bank-accounts/:customer_id", listBankAccountsHandler)
		api.POST("/bank-accounts", addBankAccountHandler)
		api.DELETE("/bank-accounts/:id", deleteBankAccountHandler)
		api.POST("/bank-accounts/:id/verify", verifyBankAccountHandler)

		// Virtual accounts
		api.GET("/virtual-accounts/:customer_id", getVirtualAccountHandler)
		api.POST("/virtual-accounts", createVirtualAccountHandler)

		// Reconciliation
		api.GET("/reconciliation/report", getReconciliationReportHandler)
		api.POST("/reconciliation/flag/:payment_id", flagPaymentHandler)
	}

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server listening", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}

func loadConfig() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8084"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		PaystackKey:    getEnv("PAYSTACK_SECRET_KEY", ""),
		FlutterwaveKey: getEnv("FLUTTERWAVE_SECRET_KEY", ""),
		StripeKey:      getEnv("STRIPE_SECRET_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initLogger(level string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, _ := config.Build()
	return logger
}

func initDatabase(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	config.MaxConns = 25
	config.MinConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "payment-service"})
}

func readinessHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

// Placeholder handlers
func createPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listPaymentsHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"payments": []interface{}{}}) }
func refundPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func paystackWebhookHandler(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func flutterwaveWebhookHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func getWalletHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func creditWalletHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func debitWalletHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getWalletTransactionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"transactions": []interface{}{}})
}
func getCreditAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createCreditAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func useCreditLineHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func repayCreditLineHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getCreditStatementHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getCreditScoreHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func calculateCreditScoreHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createDisbursementHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getDisbursementHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listDisbursementsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"disbursements": []interface{}{}})
}
func createDisbursementBatchHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listBankAccountsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"bank_accounts": []interface{}{}})
}
func addBankAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func deleteBankAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func verifyBankAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getVirtualAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createVirtualAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getReconciliationReportHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func flagPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
