// Package config provides configuration management for all OmniRoute services.
// It supports environment variables and provides type-safe configuration structs.
package config

import (
	"os"
	"strconv"
	"time"
)

// BaseConfig contains configuration common to all services
type BaseConfig struct {
	// Server configuration
	ServerPort         string        `env:"SERVER_PORT"`
	ServerReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`

	// Database configuration
	DatabaseURL         string        `env:"DATABASE_URL" required:"true"`
	DatabaseMaxConns    int           `env:"DATABASE_MAX_CONNS"`
	DatabaseMinConns    int           `env:"DATABASE_MIN_CONNS"`
	DatabaseMaxConnLife time.Duration `env:"DATABASE_MAX_CONN_LIFETIME"`

	// Redis configuration
	RedisURL      string `env:"REDIS_URL"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	// Logging configuration
	LogLevel  string `env:"LOG_LEVEL"`
	LogFormat string `env:"LOG_FORMAT"`

	// Telemetry configuration
	OTELEndpoint    string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTELServiceName string `env:"OTEL_SERVICE_NAME"`

	// Environment
	Environment string `env:"ENVIRONMENT"`
}

// Load returns a BaseConfig populated from environment variables with defaults
func Load() *BaseConfig {
	return &BaseConfig{
		// Server defaults
		ServerPort:         GetEnv("SERVER_PORT", "8080"),
		ServerReadTimeout:  GetDuration("SERVER_READ_TIMEOUT", 15*time.Second),
		ServerWriteTimeout: GetDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),

		// Database defaults
		DatabaseURL:         GetEnv("DATABASE_URL", ""),
		DatabaseMaxConns:    GetInt("DATABASE_MAX_CONNS", 25),
		DatabaseMinConns:    GetInt("DATABASE_MIN_CONNS", 5),
		DatabaseMaxConnLife: GetDuration("DATABASE_MAX_CONN_LIFETIME", time.Hour),

		// Redis defaults
		RedisURL:      GetEnv("REDIS_URL", "redis://localhost:6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),
		RedisDB:       GetInt("REDIS_DB", 0),

		// Logging defaults
		LogLevel:  GetEnv("LOG_LEVEL", "info"),
		LogFormat: GetEnv("LOG_FORMAT", "json"),

		// Telemetry defaults
		OTELEndpoint:    GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OTELServiceName: GetEnv("OTEL_SERVICE_NAME", "omniroute"),

		// Environment
		Environment: GetEnv("ENVIRONMENT", "development"),
	}
}

// IsProduction returns true if running in production environment
func (c *BaseConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development environment
func (c *BaseConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetEnv returns an environment variable value or a default
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetInt returns an environment variable as int or a default
func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetBool returns an environment variable as bool or a default
func GetBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// GetDuration returns an environment variable as duration or a default
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if durVal, err := time.ParseDuration(value); err == nil {
			return durVal
		}
	}
	return defaultValue
}

// MustGetEnv returns an environment variable or panics if not set
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("required environment variable " + key + " is not set")
	}
	return value
}
