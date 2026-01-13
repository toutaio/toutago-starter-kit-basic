package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts (title, slug, content, author_id, status, meta_title, meta_desc, is_featured, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	return r.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Slug,
		post.Content,
		post.AuthorID,
		post.Status,
		post.MetaTitle,
		post.MetaDesc,
		post.IsFeatured,
		now,
		now,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *PostRepository) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	query := `
		SELECT id, title, slug, content, author_id, status, meta_title, meta_desc, is_featured, published_at, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	post := &domain.Post{}
	var publishedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.AuthorID,
		&post.Status,
		&post.MetaTitle,
		&post.MetaDesc,
		&post.IsFeatured,
		&publishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if publishedAt.Valid {
		post.PublishedAt = &publishedAt.Time
	}

	return post, nil
}

func (r *PostRepository) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	query := `
		SELECT id, title, slug, content, author_id, status, meta_title, meta_desc, is_featured, published_at, created_at, updated_at
		FROM posts
		WHERE slug = $1
	`

	post := &domain.Post{}
	var publishedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.AuthorID,
		&post.Status,
		&post.MetaTitle,
		&post.MetaDesc,
		&post.IsFeatured,
		&publishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if publishedAt.Valid {
		post.PublishedAt = &publishedAt.Time
	}

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `
		UPDATE posts
		SET title = $1, slug = $2, content = $3, status = $4, meta_title = $5, meta_desc = $6, is_featured = $7, updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		post.Title,
		post.Slug,
		post.Content,
		post.Status,
		post.MetaTitle,
		post.MetaDesc,
		post.IsFeatured,
		time.Now(),
		post.ID,
	)

	return err
}

func (r *PostRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostRepository) List(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	query := `
		SELECT id, title, slug, content, author_id, status, meta_title, meta_desc, is_featured, published_at, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) ListByStatus(ctx context.Context, status domain.PostStatus, limit, offset int) ([]*domain.Post, error) {
	query := `
		SELECT id, title, slug, content, author_id, status, meta_title, meta_desc, is_featured, published_at, created_at, updated_at
		FROM posts
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domain.Post, error) {
	query := `
		SELECT id, title, slug, content, author_id, status, meta_title, meta_desc, is_featured, published_at, created_at, updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) scanPosts(rows *sql.Rows) ([]*domain.Post, error) {
	var posts []*domain.Post

	for rows.Next() {
		post := &domain.Post{}
		var publishedAt sql.NullTime

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.AuthorID,
			&post.Status,
			&post.MetaTitle,
			&post.MetaDesc,
			&post.IsFeatured,
			&publishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if publishedAt.Valid {
			post.PublishedAt = &publishedAt.Time
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}
