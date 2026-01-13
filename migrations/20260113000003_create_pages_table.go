package migrations

import (
	"context"

	"github.com/toutaio/toutago-sil-migrator/pkg/sil"
)

func init() {
	sil.RegisterMigration(&Migration_20260113000003_CreatePagesTable{})
}

// Migration_20260113000003_CreatePagesTable creates the pages table
type Migration_20260113000003_CreatePagesTable struct {
	sil.BaseMigration
}

// Version returns the migration version
func (m *Migration_20260113000003_CreatePagesTable) Version() string {
	return "20260113000003"
}

// Description returns the migration description
func (m *Migration_20260113000003_CreatePagesTable) Description() string {
	return "create pages table"
}

// Up applies the migration
func (m *Migration_20260113000003_CreatePagesTable) Up(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	// Try PostgreSQL syntax first
	err := adapter.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS pages (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			content TEXT NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			author_id INT NOT NULL,
			published_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	
	if err != nil {
		// Try MySQL syntax
		err = adapter.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS pages (
				id INT AUTO_INCREMENT PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				slug VARCHAR(255) UNIQUE NOT NULL,
				content TEXT NOT NULL,
				status VARCHAR(50) NOT NULL DEFAULT 'draft',
				author_id INT NOT NULL,
				published_at TIMESTAMP NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`)
	}
	
	if err != nil {
		return err
	}

	// Create indexes
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages(slug)`)
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_pages_status ON pages(status)`)
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_pages_author_id ON pages(author_id)`)

	return nil
}

// Down reverts the migration
func (m *Migration_20260113000003_CreatePagesTable) Down(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()
	return adapter.Exec(ctx, `DROP TABLE IF EXISTS pages`)
}
