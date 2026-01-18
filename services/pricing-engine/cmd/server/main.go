// Package main is the entry point for the pricing engine service
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/omniroute/pricing-engine/internal/api"
	"github.com/omniroute/pricing-engine/internal/cache"
	"github.com/omniroute/pricing-engine/internal/engine"
	"github.com/omniroute/pricing-engine/internal/repository"
)

// Config holds the service configuration
type Config struct {
	ServerPort     string
	DatabaseURL    string
	RedisURL       string
	LogLevel       string
	CacheEnabled   bool
	CacheTTL       time.Duration
}

func main() {
	// Load configuration
	cfg := loadConfig()
	
	// Initialize logger
	logger := initLogger(cfg.LogLevel)
	defer logger.Sync()
	
	logger.Info("starting pricing engine service",
		zap.String("port", cfg.ServerPort),
	)
	
	// Initialize database connection pool
	ctx := context.Background()
	dbPool, err := initDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer dbPool.Close()
	
	// Initialize Redis cache
	var priceCache engine.PriceCache
	if cfg.CacheEnabled && cfg.RedisURL != "" {
		redisClient := initRedis(cfg.RedisURL)
		priceCache = cache.NewRedisCache(redisClient, "pricing:")
		logger.Info("redis cache enabled")
	} else {
		priceCache = cache.NewNullCache()
		logger.Info("caching disabled")
	}
	
	// Initialize repositories
	productRepo := repository.NewPostgresProductRepository(dbPool)
	customerRepo := repository.NewPostgresCustomerRepository(dbPool)
	priceListRepo := repository.NewPostgresPriceListRepository(dbPool)
	volumeDiscountRepo := repository.NewPostgresVolumeDiscountRepository(dbPool)
	contractPriceRepo := repository.NewPostgresContractPriceRepository(dbPool)
	promotionRepo := repository.NewPostgresPromotionRepository(dbPool)
	taxRepo := repository.NewPostgresTaxRepository(dbPool)
	
	// Initialize pricing engine
	pricingConfig := &engine.PricingConfig{
		EnableCaching:         cfg.CacheEnabled,
		CacheTTL:              cfg.CacheTTL,
		MaxConcurrentCalcs:    100,
		EnableVolumeDiscounts: true,
		EnablePromotions:      true,
		EnableContractPricing: true,
		RoundingPrecision:     2,
		RoundingMode:          "half_up",
	}
	
	pricingEngine := engine.NewPricingEngine(
		productRepo,
		customerRepo,
		priceListRepo,
		volumeDiscountRepo,
		contractPriceRepo,
		promotionRepo,
		taxRepo,
		priceCache,
		logger,
		pricingConfig,
	)
	
	// Initialize HTTP handler
	handler := api.NewPricingHandler(pricingEngine, logger)
	
	// Setup HTTP server
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	
	// Add middleware
	wrappedMux := loggingMiddleware(logger, corsMiddleware(mux))
	
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      wrappedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in goroutine
	go func() {
		logger.Info("HTTP server listening", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("shutting down server...")
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	
	logger.Info("server stopped")
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		ServerPort:   getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/omniroute"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		CacheEnabled: getEnv("CACHE_ENABLED", "true") == "true",
		CacheTTL:     parseDuration(getEnv("CACHE_TTL", "5m")),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string with a default
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 5 * time.Minute
	}
	return d
}

// initLogger initializes the zap logger
func initLogger(level string) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	
	return logger
}

// initDatabase initializes the PostgreSQL connection pool
func initDatabase(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}
	
	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	
	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return pool, nil
}

// initRedis initializes the Redis client
func initRedis(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}
	
	return redis.NewClient(opt)
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		logger.Info("http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapped.statusCode),
			zap.Duration("duration", time.Since(start)),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID, X-Request-ID")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
