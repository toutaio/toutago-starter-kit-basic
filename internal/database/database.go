// Package database provides database connection management.
package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/config"

	// Import database drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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
