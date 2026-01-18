// Package main is the entry point for the bank-gateway service
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
	// NIBSS configuration
	NIBSSEndpoint string
	NIBSSBankCode string
	// Mobile Money configuration
	MPesaEndpoint   string
	MTNMomoEndpoint string
	AirtelEndpoint  string
}

func main() {
	cfg := loadConfig()
	logger := initLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting bank-gateway service",
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

	// Hasura Actions (POST endpoints for GraphQL actions)
	actions := router.Group("/actions")
	{
		actions.POST("/initiate-payment", initiatePaymentHandler)
		actions.POST("/verify-account", verifyAccountHandler)
		actions.POST("/process-bulk-payment", processBulkPaymentHandler)
		actions.POST("/create-virtual-account", createVirtualAccountHandler)
		actions.POST("/fetch-statement", fetchStatementHandler)
		actions.POST("/get-payment-status", getPaymentStatusHandler)
	}

	// Bank Webhooks
	webhooks := router.Group("/webhooks")
	{
		webhooks.POST("/nibss/payment-notification", nibssWebhookHandler)
		webhooks.POST("/mobile-money/callback", mobileMoneyWebhookHandler)
		webhooks.POST("/paystack/callback", paystackWebhookHandler)
		webhooks.POST("/flutterwave/callback", flutterwaveWebhookHandler)
	}

	// Internal APIs
	api := router.Group("/api/v1")
	{
		// Bank connections
		api.GET("/connections", listConnectionsHandler)
		api.GET("/connections/:id", getConnectionHandler)
		api.POST("/connections/:id/health-check", healthCheckConnectionHandler)

		// Payments
		api.GET("/payments", listPaymentsHandler)
		api.GET("/payments/:id", getPaymentHandler)
		api.POST("/payments/:id/retry", retryPaymentHandler)
		api.POST("/payments/:id/reverse", reversePaymentHandler)

		// Virtual accounts
		api.GET("/virtual-accounts", listVirtualAccountsHandler)
		api.GET("/virtual-accounts/:id", getVirtualAccountHandler)
		api.GET("/virtual-accounts/:id/transactions", getVirtualAccountTransactionsHandler)

		// Batches
		api.GET("/batches", listBatchesHandler)
		api.GET("/batches/:id", getBatchHandler)
		api.GET("/batches/:id/payments", getBatchPaymentsHandler)

		// Reconciliation
		api.GET("/reconciliation/report", getReconciliationReportHandler)
		api.POST("/reconciliation/run", runReconciliationHandler)
	}

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
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
		ServerPort:      getEnv("SERVER_PORT", "8086"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/bankgateway?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		NIBSSEndpoint:   getEnv("NIBSS_ENDPOINT", ""),
		NIBSSBankCode:   getEnv("NIBSS_BANK_CODE", ""),
		MPesaEndpoint:   getEnv("MPESA_ENDPOINT", ""),
		MTNMomoEndpoint: getEnv("MTN_MOMO_ENDPOINT", ""),
		AirtelEndpoint:  getEnv("AIRTEL_MONEY_ENDPOINT", ""),
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
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "bank-gateway"})
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

// Placeholder handlers - to be implemented
func initiatePaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func verifyAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func processBulkPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createVirtualAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func fetchStatementHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getPaymentStatusHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func nibssWebhookHandler(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func mobileMoneyWebhookHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func paystackWebhookHandler(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func flutterwaveWebhookHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func listConnectionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"connections": []interface{}{}})
}
func getConnectionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func healthCheckConnectionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listPaymentsHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"payments": []interface{}{}}) }
func getPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func retryPaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func reversePaymentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listVirtualAccountsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"virtual_accounts": []interface{}{}})
}
func getVirtualAccountHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getVirtualAccountTransactionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"transactions": []interface{}{}})
}
func listBatchesHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"batches": []interface{}{}}) }
func getBatchHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getBatchPaymentsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"payments": []interface{}{}})
}
func getReconciliationReportHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func runReconciliationHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
