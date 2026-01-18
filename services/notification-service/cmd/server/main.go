// Package main is the entry point for the notification-service
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
	// Provider configs
	TwilioSID   string
	TwilioToken string
	SendGridKey string
	FirebaseKey string
}

func main() {
	cfg := loadConfig()
	logger := initLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting notification-service",
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
		// Notifications
		api.POST("/notifications", sendNotificationHandler)
		api.POST("/notifications/bulk", sendBulkNotificationHandler)
		api.GET("/notifications/:id", getNotificationHandler)
		api.GET("/notifications/:id/status", getNotificationStatusHandler)

		// Templates
		api.GET("/templates", listTemplatesHandler)
		api.GET("/templates/:id", getTemplateHandler)
		api.POST("/templates", createTemplateHandler)
		api.PUT("/templates/:id", updateTemplateHandler)
		api.DELETE("/templates/:id", deleteTemplateHandler)

		// USSD
		api.POST("/ussd/callback", ussdCallbackHandler)
		api.GET("/ussd/sessions/:id", getUSSDSessionHandler)
		api.GET("/ussd/menus", listUSSDMenusHandler)
		api.POST("/ussd/menus", createUSSDMenuHandler)

		// WhatsApp
		api.POST("/whatsapp/webhook", whatsappWebhookHandler)
		api.GET("/whatsapp/templates", listWhatsAppTemplatesHandler)

		// Devices (Push notifications)
		api.POST("/devices", registerDeviceHandler)
		api.DELETE("/devices/:id", unregisterDeviceHandler)

		// Preferences
		api.GET("/preferences/:user_id", getPreferencesHandler)
		api.PUT("/preferences/:user_id", updatePreferencesHandler)
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
		ServerPort:  getEnv("SERVER_PORT", "8083"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/notifications?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		TwilioSID:   getEnv("TWILIO_SID", ""),
		TwilioToken: getEnv("TWILIO_TOKEN", ""),
		SendGridKey: getEnv("SENDGRID_API_KEY", ""),
		FirebaseKey: getEnv("FIREBASE_KEY", ""),
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
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "notification-service"})
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
func sendNotificationHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func sendBulkNotificationHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getNotificationHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getNotificationStatusHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listTemplatesHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"templates": []interface{}{}}) }
func getTemplateHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func createTemplateHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func updateTemplateHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func deleteTemplateHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func ussdCallbackHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getUSSDSessionHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func listUSSDMenusHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"menus": []interface{}{}}) }
func createUSSDMenuHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func whatsappWebhookHandler(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }
func listWhatsAppTemplatesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"templates": []interface{}{}})
}
func registerDeviceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func unregisterDeviceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func getPreferencesHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func updatePreferencesHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
