// Package database provides database connection management.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/toutaio/toutago-sil-migrator/pkg/sil"
	"github.com/toutaio/toutago-sil-migrator/pkg/sil/adapters"
	"github.com/toutaio/toutago-starter-kit-basic/internal/config"

	// Import database drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	// Import migrations to register them
	_ "github.com/toutaio/toutago-starter-kit-basic/migrations"
)

// Connect establishes a database connection based on configuration.
func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	connStr := cfg.ConnectionString()
	if connStr == "" {
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Ping checks if the database connection is alive.
func Ping(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return db.Ping()
}

// Close closes the database connection.
func Close(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return db.Close()
}

// Stats returns database statistics.
func Stats(db *sql.DB) sql.DBStats {
	if db == nil {
		return sql.DBStats{}
	}
	return db.Stats()
}

// RunMigrations runs database migrations using Sil migrator
func RunMigrations(cfg config.DatabaseConfig) error {
	ctx := context.Background()

	// Create Sil config
	silConfig := sil.DefaultConfig()
	silConfig.DatabaseURL = buildDatabaseURL(cfg)
	silConfig.MigrationsDir = "./migrations"
	silConfig.Verbose = true

	// Create appropriate adapter
	var adapter sil.DatabaseAdapter
	var err error

	switch cfg.Driver {
	case "postgres":
		adapter, err = adapters.NewPostgresAdapter(silConfig)
	case "mysql":
		adapter, err = adapters.NewMySQLAdapter(silConfig)
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return fmt.Errorf("failed to create database adapter: %w", err)
	}

	// Connect to database
	if err := adapter.Connect(ctx, silConfig); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer adapter.Close()

	// Create migrator
	migrator, err := sil.NewMigrator(silConfig, adapter)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Run migrations
	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

// buildDatabaseURL builds the database URL for Sil migrator
func buildDatabaseURL(cfg config.DatabaseConfig) string {
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)
	case "mysql":
		return fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)
	default:
		return ""
	}
}
