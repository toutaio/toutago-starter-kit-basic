package migrations

import (
	"context"

	"github.com/toutaio/toutago-sil-migrator/pkg/sil"
)

func init() {
	sil.RegisterMigration(&Migration_20260113000001_CreateUsersTable{})
}

// Migration_20260113000001_CreateUsersTable creates the users table with authentication fields
type Migration_20260113000001_CreateUsersTable struct {
	sil.BaseMigration
}

// Version returns the migration version
func (m *Migration_20260113000001_CreateUsersTable) Version() string {
	return "20260113000001"
}

// Description returns the migration description
func (m *Migration_20260113000001_CreateUsersTable) Description() string {
	return "create users table"
}

// Up applies the migration
func (m *Migration_20260113000001_CreateUsersTable) Up(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	// Try PostgreSQL syntax first
	err := adapter.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			email_verified BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	
	if err != nil {
		// Try MySQL syntax
		err = adapter.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS users (
				id INT AUTO_INCREMENT PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				name VARCHAR(255) NOT NULL,
				role VARCHAR(50) NOT NULL DEFAULT 'user',
				email_verified BOOLEAN NOT NULL DEFAULT FALSE,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`)
	}
	
	if err != nil {
		return err
	}

	// Create indexes - compatible with both databases
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`)

	return nil
}

// Down reverts the migration
func (m *Migration_20260113000001_CreateUsersTable) Down(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()
	return adapter.Exec(ctx, `DROP TABLE IF EXISTS users`)
}
