// Package database provides PostgreSQL connection pool management and
// migration runner for TeslaPay microservices.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// DB wraps a pgxpool.Pool with logging and health check capabilities.
type DB struct {
	Pool   *pgxpool.Pool
	logger *zap.Logger
}

// Config holds PostgreSQL connection configuration.
type Config struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	HealthCheckFreq time.Duration
}

// DefaultConfig returns sensible default connection pool settings
// suitable for a financial services workload.
func DefaultConfig(url string) *Config {
	return &Config{
		URL:             url,
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
		HealthCheckFreq: 30 * time.Second,
	}
}

// New creates a new database connection pool. It validates the connection
// by pinging the database before returning.
func New(ctx context.Context, cfg *Config, logger *zap.Logger) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckFreq

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Validate connection before returning.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("database connection pool established",
		zap.String("host", poolCfg.ConnConfig.Host),
		zap.Uint16("port", poolCfg.ConnConfig.Port),
		zap.String("database", poolCfg.ConnConfig.Database),
		zap.Int32("max_conns", cfg.MaxConns),
		zap.Int32("min_conns", cfg.MinConns),
	)

	return &DB{Pool: pool, logger: logger}, nil
}

// Close gracefully shuts down the connection pool.
func (db *DB) Close() {
	db.Pool.Close()
	db.logger.Info("database connection pool closed")
}

// HealthCheck verifies the database connection is alive.
func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return db.Pool.Ping(ctx)
}

// DBTX is an interface that abstracts pgxpool.Pool and pgx.Tx,
// allowing repository methods to work with both direct queries and transactions.
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// WithTransaction executes fn within a database transaction. If fn returns
// an error, the transaction is rolled back. Otherwise it is committed.
// This is the primary mechanism for ensuring ACID properties in ledger postings.
func (db *DB) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // Re-throw after rollback.
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			db.logger.Error("failed to rollback transaction",
				zap.Error(rbErr),
				zap.Error(err),
			)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
