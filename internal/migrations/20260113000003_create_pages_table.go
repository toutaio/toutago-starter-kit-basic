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

// Version returns the migration version.
func (m *Migration_20260113000003_CreatePagesTable) Version() string {
	return "20260113000003"
}

// Description returns the migration description.
func (m *Migration_20260113000003_CreatePagesTable) Description() string {
	return "create pages table"
}

// Up applies the migration.
func (m *Migration_20260113000003_CreatePagesTable) Up(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	return adapter.Exec(ctx, `
		CREATE TABLE pages (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			content TEXT NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			published_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX idx_pages_slug ON pages(slug);
		CREATE INDEX idx_pages_status ON pages(status);
		CREATE INDEX idx_pages_author_id ON pages(author_id);
	`)
}

// Down reverts the migration.
func (m *Migration_20260113000003_CreatePagesTable) Down(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	return adapter.Exec(ctx, `DROP TABLE IF EXISTS pages CASCADE`)
}
