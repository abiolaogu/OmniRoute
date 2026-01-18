// Package main is the entry point for the ATC (Authority to Collect) service
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

func main() {
	cfg := loadConfig()
	logger := initLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting atc-service",
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

	// Hasura Actions
	actions := router.Group("/actions")
	{
		// ATC Grant management
		actions.POST("/create-atc-grant", createATCGrantHandler)
		actions.POST("/activate-atc-grant", activateATCGrantHandler)
		actions.POST("/suspend-atc-grant", suspendATCGrantHandler)
		actions.POST("/revoke-atc-grant", revokeATCGrantHandler)

		// Collection operations
		actions.POST("/record-collection", recordCollectionHandler)
		actions.POST("/verify-collection-authority", verifyCollectionAuthorityHandler)

		// Settlement operations
		actions.POST("/create-settlement-batch", createSettlementBatchHandler)
		actions.POST("/process-settlement", processSettlementHandler)
	}

	// REST API
	api := router.Group("/api/v1")
	{
		// ATC Grants
		api.GET("/atc-grants", listATCGrantsHandler)
		api.GET("/atc-grants/:id", getATCGrantHandler)
		api.POST("/atc-grants", createATCGrantAPIHandler)
		api.PUT("/atc-grants/:id", updateATCGrantHandler)
		api.DELETE("/atc-grants/:id", deleteATCGrantHandler)

		// Collections
		api.GET("/collections", listCollectionsHandler)
		api.GET("/collections/:id", getCollectionHandler)
		api.GET("/atc-grants/:id/collections", getATCGrantCollectionsHandler)

		// Settlements
		api.GET("/settlement-batches", listSettlementBatchesHandler)
		api.GET("/settlement-batches/:id", getSettlementBatchHandler)
		api.POST("/settlement-batches/:id/process", processSettlementBatchHandler)

		// Reports
		api.GET("/reports/commission-summary", getCommissionSummaryHandler)
		api.GET("/reports/collection-summary", getCollectionSummaryHandler)
	}

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
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

type Config struct {
	ServerPort  string
	DatabaseURL string
	RedisURL    string
	LogLevel    string
}

func loadConfig() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8087"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/atcservice?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
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
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "atc-service"})
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
func createATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func activateATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func suspendATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func revokeATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func recordCollectionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func verifyCollectionAuthorityHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createSettlementBatchHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func processSettlementHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listATCGrantsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"atc_grants": []interface{}{}})
}
func getATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createATCGrantAPIHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func updateATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func deleteATCGrantHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listCollectionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"collections": []interface{}{}})
}
func getCollectionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getATCGrantCollectionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"collections": []interface{}{}})
}
func listSettlementBatchesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"batches": []interface{}{}})
}
func getSettlementBatchHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func processSettlementBatchHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getCommissionSummaryHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getCollectionSummaryHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
