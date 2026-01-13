package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
)

type PageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) Create(ctx context.Context, page *domain.Page) error {
	query := `
		INSERT INTO pages (title, slug, content, status, meta_title, meta_desc, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	return r.db.QueryRowContext(
		ctx,
		query,
		page.Title,
		page.Slug,
		page.Content,
		page.Status,
		page.MetaTitle,
		page.MetaDesc,
		now,
		now,
	).Scan(&page.ID, &page.CreatedAt, &page.UpdatedAt)
}

func (r *PageRepository) GetByID(ctx context.Context, id int64) (*domain.Page, error) {
	query := `
		SELECT id, title, slug, content, status, meta_title, meta_desc, published_at, created_at, updated_at
		FROM pages
		WHERE id = $1
	`

	page := &domain.Page{}
	var publishedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&page.ID,
		&page.Title,
		&page.Slug,
		&page.Content,
		&page.Status,
		&page.MetaTitle,
		&page.MetaDesc,
		&publishedAt,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if publishedAt.Valid {
		page.PublishedAt = &publishedAt.Time
	}

	return page, nil
}

func (r *PageRepository) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	query := `
		SELECT id, title, slug, content, status, meta_title, meta_desc, published_at, created_at, updated_at
		FROM pages
		WHERE slug = $1
	`

	page := &domain.Page{}
	var publishedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&page.ID,
		&page.Title,
		&page.Slug,
		&page.Content,
		&page.Status,
		&page.MetaTitle,
		&page.MetaDesc,
		&publishedAt,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if publishedAt.Valid {
		page.PublishedAt = &publishedAt.Time
	}

	return page, nil
}

func (r *PageRepository) Update(ctx context.Context, page *domain.Page) error {
	query := `
		UPDATE pages
		SET title = $1, slug = $2, content = $3, status = $4, meta_title = $5, meta_desc = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		page.Title,
		page.Slug,
		page.Content,
		page.Status,
		page.MetaTitle,
		page.MetaDesc,
		time.Now(),
		page.ID,
	)

	return err
}

func (r *PageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pages WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Page, error) {
	query := `
		SELECT id, title, slug, content, status, meta_title, meta_desc, published_at, created_at, updated_at
		FROM pages
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPages(rows)
}

func (r *PageRepository) ListByStatus(ctx context.Context, status domain.PageStatus, limit, offset int) ([]*domain.Page, error) {
	query := `
		SELECT id, title, slug, content, status, meta_title, meta_desc, published_at, created_at, updated_at
		FROM pages
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPages(rows)
}

func (r *PageRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domain.Page, error) {
	query := `
		SELECT id, title, slug, content, status, meta_title, meta_desc, published_at, created_at, updated_at
		FROM pages
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPages(rows)
}

func (r *PageRepository) scanPages(rows *sql.Rows) ([]*domain.Page, error) {
	var pages []*domain.Page

	for rows.Next() {
		page := &domain.Page{}
		var publishedAt sql.NullTime

		err := rows.Scan(
			&page.ID,
			&page.Title,
			&page.Slug,
			&page.Content,
			&page.Status,
			&page.MetaTitle,
			&page.MetaDesc,
			&publishedAt,
			&page.CreatedAt,
			&page.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if publishedAt.Valid {
			page.PublishedAt = &publishedAt.Time
		}

		pages = append(pages, page)
	}

	return pages, rows.Err()
}
