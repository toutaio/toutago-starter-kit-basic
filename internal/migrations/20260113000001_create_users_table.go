package migrations

import (
	"context"
	"github.com/toutaio/toutago-sil-migrator/pkg/sil"
)

func init() {
	sil.RegisterMigration(&Migration_20260113000001_CreateUsersTable{})
}

// Migration_20260113000001_CreateUsersTable creates the users table
type Migration_20260113000001_CreateUsersTable struct {
	sil.BaseMigration
}

// Version returns the migration version.
func (m *Migration_20260113000001_CreateUsersTable) Version() string {
	return "20260113000001"
}

// Description returns the migration description.
func (m *Migration_20260113000001_CreateUsersTable) Description() string {
	return "create users table"
}

// Up applies the migration.
func (m *Migration_20260113000001_CreateUsersTable) Up(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	return adapter.Exec(ctx, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			email_verified_at TIMESTAMP,
			remember_token VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
}

// Down reverts the migration.
func (m *Migration_20260113000001_CreateUsersTable) Down(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	return adapter.Exec(ctx, `DROP TABLE IF EXISTS users CASCADE`)
}
