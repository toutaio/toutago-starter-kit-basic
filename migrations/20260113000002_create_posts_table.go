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

// Version returns the migration version
func (m *Migration_20260113000002_CreatePostsTable) Version() string {
	return "20260113000002"
}

// Description returns the migration description
func (m *Migration_20260113000002_CreatePostsTable) Description() string {
	return "create posts table"
}

// Up applies the migration
func (m *Migration_20260113000002_CreatePostsTable) Up(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()

	// Try PostgreSQL syntax first
	err := adapter.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			content TEXT NOT NULL,
			excerpt TEXT,
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
			CREATE TABLE IF NOT EXISTS posts (
				id INT AUTO_INCREMENT PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				slug VARCHAR(255) UNIQUE NOT NULL,
				content TEXT NOT NULL,
				excerpt TEXT,
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
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug)`)
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status)`)
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id)`)
	adapter.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at)`)

	return nil
}

// Down reverts the migration
func (m *Migration_20260113000002_CreatePostsTable) Down(adapter sil.DatabaseAdapter) error {
	ctx := context.Background()
	return adapter.Exec(ctx, `DROP TABLE IF EXISTS posts`)
}
