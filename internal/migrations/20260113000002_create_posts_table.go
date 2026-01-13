package migrations

import (
	"context"
	"github.com/toutaio/toutago-sil-migrator/pkg/sil"
)

func init() {
	sil.RegisterMigration(&Migration_20260113000002_CreatePostsTable{})
}

// Migration_20260113000002_CreatePostsTable creates the posts table
type Migration_20260113000002_CreatePostsTable struct {
	sil.BaseMigration
}

// Version returns the migration version.
func (m *Migration_20260113000002_CreatePostsTable) Version() string {
	return "20260113000002"
}

// Description returns the migration description.
func (m *Migration_20260113000002_CreatePostsTable) Description() string {
	return "create posts table"
}

// Up applies the migration.
func (m *Migration_20260113000002_CreatePostsTable) Up(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	return adapter.Exec(ctx, `
		CREATE TABLE posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			content TEXT NOT NULL,
			excerpt TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			published_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX idx_posts_slug ON posts(slug);
		CREATE INDEX idx_posts_status ON posts(status);
		CREATE INDEX idx_posts_author_id ON posts(author_id);
		CREATE INDEX idx_posts_published_at ON posts(published_at);
	`)
}

// Down reverts the migration.
func (m *Migration_20260113000002_CreatePostsTable) Down(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	return adapter.Exec(ctx, `DROP TABLE IF EXISTS posts CASCADE`)
}
