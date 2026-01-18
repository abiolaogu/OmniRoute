// Package database provides YugabyteDB client configuration and utilities.
// YugabyteDB is a PostgreSQL-compatible distributed database used for
// horizontally scalable, multi-region deployment.
package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// YugabyteConfig holds YugabyteDB connection configuration
type YugabyteConfig struct {
	// Hosts is a list of YugabyteDB node addresses for HA
	Hosts []string
	// Port is the YSQL port (default: 5433)
	Port int
	// Database name
	Database string
	// User for authentication
	User string
	// Password for authentication
	Password string
	// SSLMode (disable, require, verify-ca, verify-full)
	SSLMode string
	// MaxConnections in the pool
	MaxConnections int32
	// MinConnections to keep alive
	MinConnections int32
	// MaxConnLifetime is the maximum time a connection can live
	MaxConnLifetime time.Duration
	// LoadBalance enables YugabyteDB smart driver load balancing
	LoadBalance bool
	// TopologyKeys for zone-aware routing (e.g., "gcp.africa-south1.africa-south1-a")
	TopologyKeys string
	// ApplicationName for connection identification
	ApplicationName string
}

// DefaultYugabyteConfig returns default configuration
func DefaultYugabyteConfig() YugabyteConfig {
	return YugabyteConfig{
		Hosts:           []string{"localhost"},
		Port:            5433,
		Database:        "omniroute",
		User:            "yugabyte",
		Password:        "yugabyte",
		SSLMode:         "disable",
		MaxConnections:  25,
		MinConnections:  5,
		MaxConnLifetime: time.Hour,
		LoadBalance:     true,
		ApplicationName: "omniroute-service",
	}
}

// ConnectionString builds the YugabyteDB connection string
func (c YugabyteConfig) ConnectionString() string {
	host := strings.Join(c.Hosts, ",")

	connStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s "+
			"pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%s "+
			"application_name=%s",
		host,
		c.Port,
		c.Database,
		c.User,
		c.Password,
		c.SSLMode,
		c.MaxConnections,
		c.MinConnections,
		c.MaxConnLifetime.String(),
		c.ApplicationName,
	)

	// Add YugabyteDB-specific options
	if c.LoadBalance {
		connStr += " load_balance=true"
	}
	if c.TopologyKeys != "" {
		connStr += fmt.Sprintf(" topology_keys=%s", c.TopologyKeys)
	}

	return connStr
}

// NewYugabytePool creates a new connection pool to YugabyteDB
func NewYugabytePool(ctx context.Context, cfg YugabyteConfig) (*pgxpool.Pool, error) {
	connStr := cfg.ConnectionString()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Set runtime parameters
	poolConfig.ConnConfig.RuntimeParams["application_name"] = cfg.ApplicationName
	poolConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}

// HealthCheck performs a health check on the database connection
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return pool.Ping(ctx)
}

// TransactionFunc is a function that runs within a transaction
type TransactionFunc func(ctx context.Context, tx pgxpool.Tx) error

// WithTransaction executes a function within a database transaction
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn TransactionFunc) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// BatchInsert performs efficient batch inserts using COPY
func BatchInsert(ctx context.Context, pool *pgxpool.Pool, tableName string, columns []string, rows [][]interface{}) (int64, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	copyCount, err := conn.Conn().CopyFrom(
		ctx,
		[]string{tableName},
		columns,
		pgxCopyFromRows(rows),
	)
	if err != nil {
		return 0, fmt.Errorf("copy from: %w", err)
	}

	return copyCount, nil
}

// pgxCopyFromRows implements pgx.CopyFromSource
type pgxCopyFromRows [][]interface{}

func (r pgxCopyFromRows) Next() bool {
	return len(r) > 0
}

func (r *pgxCopyFromRows) Values() ([]interface{}, error) {
	if len(*r) == 0 {
		return nil, nil
	}
	row := (*r)[0]
	*r = (*r)[1:]
	return row, nil
}

func (r pgxCopyFromRows) Err() error {
	return nil
}
