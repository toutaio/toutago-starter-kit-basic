package database

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

func TestRunMigrations_Postgres(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := sql.Open("postgres", "postgres://starter_user:starter_pass@localhost:5432/starter_db?sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test: cannot connect to postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: cannot ping postgres: %v", err)
	}

	// Clean up before test
	_, _ = db.Exec("DROP TABLE IF EXISTS sil_migrations CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS page_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS pages CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS categories CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS sessions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Run migrations
	err = RunMigrations(db, "postgres")
	if err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	// Verify migrations table exists
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'sil_migrations'
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check sil_migrations table: %v", err)
	}
	if !exists {
		t.Error("sil_migrations table does not exist")
	}

	// Verify users table exists
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'users'
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check users table: %v", err)
	}
	if !exists {
		t.Error("users table does not exist")
	}

	// Check migration status
	statuses, err := GetMigrationStatus(db, "postgres")
	if err != nil {
		t.Fatalf("GetMigrationStatus failed: %v", err)
	}
	if len(statuses) == 0 {
		t.Error("Expected migration statuses")
	}

	// Clean up after test
	_, _ = db.Exec("DROP TABLE IF EXISTS sil_migrations CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS page_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS pages CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS categories CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS sessions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS users CASCADE")
}

func TestRunMigrations_MySQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := sql.Open("mysql", "starter_user:starter_pass@tcp(localhost:3306)/starter_db?parseTime=true")
	if err != nil {
		t.Skipf("Skipping test: cannot connect to mysql: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: cannot ping mysql: %v", err)
	}

	// Clean up before test
	_, _ = db.Exec("DROP TABLE IF EXISTS page_versions")
	_, _ = db.Exec("DROP TABLE IF EXISTS pages")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_versions")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_tags")
	_, _ = db.Exec("DROP TABLE IF EXISTS posts")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags")
	_, _ = db.Exec("DROP TABLE IF EXISTS categories")
	_, _ = db.Exec("DROP TABLE IF EXISTS sessions")
	_, _ = db.Exec("DROP TABLE IF EXISTS users")
	_, _ = db.Exec("DROP TABLE IF EXISTS sil_migrations")

	// Run migrations
	err = RunMigrations(db, "mysql")
	if err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	// Verify migrations table exists
	var tableName string
	err = db.QueryRow("SHOW TABLES LIKE 'sil_migrations'").Scan(&tableName)
	if err != nil {
		t.Fatalf("Failed to check sil_migrations table: %v", err)
	}
	if tableName != "sil_migrations" {
		t.Error("sil_migrations table does not exist")
	}

	// Verify users table exists
	err = db.QueryRow("SHOW TABLES LIKE 'users'").Scan(&tableName)
	if err != nil {
		t.Fatalf("Failed to check users table: %v", err)
	}
	if tableName != "users" {
		t.Error("users table does not exist")
	}

	// Check migration status
	statuses, err := GetMigrationStatus(db, "mysql")
	if err != nil {
		t.Fatalf("GetMigrationStatus failed: %v", err)
	}
	if len(statuses) == 0 {
		t.Error("Expected migration statuses")
	}

	// Clean up after test
	_, _ = db.Exec("DROP TABLE IF EXISTS page_versions")
	_, _ = db.Exec("DROP TABLE IF EXISTS pages")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_versions")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_tags")
	_, _ = db.Exec("DROP TABLE IF EXISTS posts")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags")
	_, _ = db.Exec("DROP TABLE IF EXISTS categories")
	_, _ = db.Exec("DROP TABLE IF EXISTS sessions")
	_, _ = db.Exec("DROP TABLE IF EXISTS users")
	_, _ = db.Exec("DROP TABLE IF EXISTS sil_migrations")
}

func TestRollbackMigrations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := sql.Open("postgres", "postgres://starter_user:starter_pass@localhost:5432/starter_db?sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test: cannot connect to postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: cannot ping postgres: %v", err)
	}

	// Clean up before test
	_, _ = db.Exec("DROP TABLE IF EXISTS sil_migrations CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS page_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS pages CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS categories CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS sessions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Run migrations first
	err = RunMigrations(db, "postgres")
	if err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	statusesBefore, err := GetMigrationStatus(db, "postgres")
	if err != nil {
		t.Fatalf("GetMigrationStatus failed: %v", err)
	}

	appliedBefore := 0
	for _, s := range statusesBefore {
		if s.Applied {
			appliedBefore++
		}
	}

	// Rollback last migration
	err = RollbackMigrations(db, "postgres")
	if err != nil {
		t.Fatalf("RollbackMigrations failed: %v", err)
	}

	statusesAfter, err := GetMigrationStatus(db, "postgres")
	if err != nil {
		t.Fatalf("GetMigrationStatus failed: %v", err)
	}

	appliedAfter := 0
	for _, s := range statusesAfter {
		if s.Applied {
			appliedAfter++
		}
	}

	if appliedAfter >= appliedBefore {
		t.Errorf("Expected fewer applied migrations after rollback, got before=%d, after=%d", appliedBefore, appliedAfter)
	}

	// Clean up after test
	_, _ = db.Exec("DROP TABLE IF EXISTS sil_migrations CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS page_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS pages CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_versions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS post_tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS categories CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS sessions CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS users CASCADE")
}

func TestRunMigrations_UnsupportedDB(t *testing.T) {
	db := &sql.DB{} // Empty DB for testing error case
	err := RunMigrations(db, "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported database type")
	}
}
