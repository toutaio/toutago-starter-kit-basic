package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed all:migrations
var migrationsFS embed.FS

// RunMigrations executes database migrations based on the database type
func RunMigrations(db *sql.DB, dbType string) error {
	// Create migrations table
	if err := createMigrationsTable(db, dbType); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrations, err := getMigrationFiles(dbType)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Apply migrations
	for _, migration := range migrations {
		if err := applyMigration(db, migration, dbType); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration, err)
		}
	}

	return nil
}

// RollbackMigrations rolls back the last migration
func RollbackMigrations(db *sql.DB, dbType string) error {
	// Get last applied migration
	lastMigration, err := getLastMigration(db, dbType)
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	if lastMigration == "" {
		return fmt.Errorf("no migrations to rollback")
	}

	// Execute down migration
	downFile := strings.Replace(lastMigration, ".up.sql", ".down.sql", 1)
	if err := rollbackMigration(db, downFile, dbType); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration, err)
	}

	return nil
}

// GetMigrationStatus returns list of migrations with their status
func GetMigrationStatus(db *sql.DB, dbType string) ([]MigrationStatus, error) {
	// Create migrations table if not exists
	if err := createMigrationsTable(db, dbType); err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migrations
	allMigrations, err := getMigrationFiles(dbType)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db, dbType)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedMap[m] = true
	}

	var statuses []MigrationStatus
	for _, migration := range allMigrations {
		name := filepath.Base(migration)
		name = strings.TrimSuffix(name, ".up.sql")
		statuses = append(statuses, MigrationStatus{
			Version:     name,
			Description: name,
			Applied:     appliedMap[migration],
		})
	}

	return statuses, nil
}

type MigrationStatus struct {
	Version     string
	Description string
	Applied     bool
}

func createMigrationsTable(db *sql.DB, dbType string) error {
	var query string
	switch dbType {
	case "postgres":
		query = `
			CREATE TABLE IF NOT EXISTS migrations (
				id SERIAL PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`
	case "mysql":
		query = `
			CREATE TABLE IF NOT EXISTS migrations (
				id INT AUTO_INCREMENT PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	_, err := db.Exec(query)
	return err
}

func getMigrationFiles(dbType string) ([]string, error) {
	dir := fmt.Sprintf("migrations/%s", dbType)
	entries, err := migrationsFS.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrations = append(migrations, filepath.Join(dir, entry.Name()))
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

func applyMigration(db *sql.DB, migrationFile, dbType string) error {
	// Check if already applied
	applied, err := isMigrationApplied(db, migrationFile, dbType)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	// Read migration file
	content, err := migrationsFS.ReadFile(migrationFile)
	if err != nil {
		return err
	}

	// Execute migration
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	// Record migration
	_, err = db.Exec("INSERT INTO migrations (version) VALUES ($1)", migrationFile)
	return err
}

func rollbackMigration(db *sql.DB, downFile, dbType string) error {
	// Read down migration file
	content, err := migrationsFS.ReadFile(downFile)
	if err != nil {
		return err
	}

	// Execute down migration
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	// Remove migration record
	upFile := strings.Replace(downFile, ".down.sql", ".up.sql", 1)
	_, err = db.Exec("DELETE FROM migrations WHERE version = $1", upFile)
	return err
}

func isMigrationApplied(db *sql.DB, migration, dbType string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = $1", migration).Scan(&count)
	return count > 0, err
}

func getAppliedMigrations(db *sql.DB, dbType string) ([]string, error) {
	rows, err := db.Query("SELECT version FROM migrations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, rows.Err()
}

func getLastMigration(db *sql.DB, dbType string) (string, error) {
	var migration string
	err := db.QueryRow("SELECT version FROM migrations ORDER BY id DESC LIMIT 1").Scan(&migration)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return migration, err
}

func getDatabaseURL(dbType string) string {
	switch dbType {
	case "postgres":
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	case "mysql":
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	default:
		return ""
	}
}
