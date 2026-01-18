// Package main is the entry point for the SCE (Service Creation Environment) service
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
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the service configuration
type Config struct {
	ServerPort   string
	DatabaseURL  string
	RedisURL     string
	TemporalHost string
	LogLevel     string
}

func main() {
	cfg := loadConfig()
	logger := initLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting sce-service",
		zap.String("port", cfg.ServerPort),
	)

	ctx := context.Background()
	dbPool, err := initDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer dbPool.Close()

	// Initialize Temporal client (optional, may not be running locally)
	var temporalClient client.Client
	if cfg.TemporalHost != "" {
		temporalClient, err = client.Dial(client.Options{
			HostPort: cfg.TemporalHost,
		})
		if err != nil {
			logger.Warn("failed to connect to Temporal, workflows disabled", zap.Error(err))
		} else {
			defer temporalClient.Close()
			logger.Info("connected to Temporal", zap.String("host", cfg.TemporalHost))
		}
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Health endpoints
	router.GET("/health", healthHandler)
	router.GET("/ready", readinessHandler(dbPool))

	// API routes
	api := router.Group("/api/v1")
	{
		// Service definitions
		api.GET("/services", listServicesHandler)
		api.GET("/services/:id", getServiceHandler)
		api.POST("/services", createServiceHandler)
		api.PUT("/services/:id", updateServiceHandler)
		api.DELETE("/services/:id", deleteServiceHandler)

		// Service lifecycle
		api.POST("/services/:id/publish", publishServiceHandler)
		api.POST("/services/:id/deprecate", deprecateServiceHandler)
		api.POST("/services/:id/archive", archiveServiceHandler)

		// Service versions
		api.GET("/services/:id/versions", listVersionsHandler)
		api.POST("/services/:id/versions", createVersionHandler)
		api.PUT("/services/:id/versions/:version_id/activate", activateVersionHandler)

		// Workflow management
		api.GET("/services/:id/workflow", getWorkflowHandler)
		api.PUT("/services/:id/workflow", updateWorkflowHandler)
		api.POST("/services/:id/workflow/validate", validateWorkflowHandler)

		// Service execution
		api.POST("/services/:id/execute", executeServiceHandler)
		api.GET("/executions/:execution_id", getExecutionHandler)
		api.GET("/executions/:execution_id/status", getExecutionStatusHandler)
		api.POST("/executions/:execution_id/signal", signalExecutionHandler)
		api.POST("/executions/:execution_id/cancel", cancelExecutionHandler)

		// n8n integration
		api.GET("/n8n/workflows", listN8NWorkflowsHandler)
		api.POST("/n8n/workflows/:id/import", importN8NWorkflowHandler)

		// AI generation
		api.POST("/ai/generate-workflow", generateWorkflowHandler)
		api.POST("/ai/suggest-nodes", suggestNodesHandler)
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
		ServerPort:   getEnv("SERVER_PORT", "8085"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/sce?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		TemporalHost: getEnv("TEMPORAL_HOST", "localhost:7233"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
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
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "sce-service"})
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
func listServicesHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"services": []interface{}{}}) }
func getServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func updateServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func deleteServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func publishServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func deprecateServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func archiveServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listVersionsHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"versions": []interface{}{}}) }
func createVersionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func activateVersionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getWorkflowHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func updateWorkflowHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func validateWorkflowHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func executeServiceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getExecutionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getExecutionStatusHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func signalExecutionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func cancelExecutionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listN8NWorkflowsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"workflows": []interface{}{}})
}
func importN8NWorkflowHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func generateWorkflowHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func suggestNodesHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
