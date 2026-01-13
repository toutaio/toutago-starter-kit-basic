package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
)

func TestPostRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()
	now := time.Now()

	post := &domain.Post{
		Title:      "Test Post",
		Slug:       "test-post",
		Content:    "This is test content",
		AuthorID:   1,
		Status:     domain.PostStatusDraft,
		MetaTitle:  "Test Meta Title",
		MetaDesc:   "Test meta description",
		IsFeatured: false,
	}

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(post.Title, post.Slug, post.Content, post.AuthorID, post.Status, post.MetaTitle, post.MetaDesc, post.IsFeatured, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	err = repo.Create(ctx, post)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), post.ID)
	assert.NotZero(t, post.CreatedAt)
	assert.NotZero(t, post.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "author_id", "status",
		"meta_title", "meta_desc", "is_featured", "published_at", "created_at", "updated_at",
	}).AddRow(
		1, "Test Post", "test-post", "Content", 1, domain.PostStatusPublished,
		"Meta Title", "Meta Desc", true, now, now, now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM posts WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	post, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, int64(1), post.ID)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "test-post", post.Slug)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "author_id", "status",
		"meta_title", "meta_desc", "is_featured", "published_at", "created_at", "updated_at",
	}).AddRow(
		1, "Test Post", "test-post", "Content", 1, domain.PostStatusPublished,
		"Meta Title", "Meta Desc", true, now, now, now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM posts WHERE slug = \$1`).
		WithArgs("test-post").
		WillReturnRows(rows)

	post, err := repo.GetBySlug(ctx, "test-post")
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "test-post", post.Slug)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()

	post := &domain.Post{
		ID:         1,
		Title:      "Updated Post",
		Slug:       "updated-post",
		Content:    "Updated content",
		AuthorID:   1,
		Status:     domain.PostStatusPublished,
		MetaTitle:  "Updated Meta",
		MetaDesc:   "Updated desc",
		IsFeatured: true,
	}

	mock.ExpectExec(`UPDATE posts SET`).
		WithArgs(post.Title, post.Slug, post.Content, post.Status, post.MetaTitle, post.MetaDesc, post.IsFeatured, sqlmock.AnyArg(), post.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(ctx, post)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "author_id", "status",
		"meta_title", "meta_desc", "is_featured", "published_at", "created_at", "updated_at",
	}).
		AddRow(1, "Post 1", "post-1", "Content 1", 1, domain.PostStatusPublished, "Meta 1", "Desc 1", true, now, now, now).
		AddRow(2, "Post 2", "post-2", "Content 2", 1, domain.PostStatusPublished, "Meta 2", "Desc 2", false, now, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	posts, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, "Post 2", posts[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_ListByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "author_id", "status",
		"meta_title", "meta_desc", "is_featured", "published_at", "created_at", "updated_at",
	}).AddRow(1, "Post 1", "post-1", "Content 1", 1, domain.PostStatusPublished, "Meta 1", "Desc 1", true, now, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM posts WHERE status = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(domain.PostStatusPublished, 10, 0).
		WillReturnRows(rows)

	posts, err := repo.ListByStatus(ctx, domain.PostStatusPublished, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, domain.PostStatusPublished, posts[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_ListByAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "author_id", "status",
		"meta_title", "meta_desc", "is_featured", "published_at", "created_at", "updated_at",
	}).AddRow(1, "Post 1", "post-1", "Content 1", 1, domain.PostStatusPublished, "Meta 1", "Desc 1", true, now, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM posts WHERE author_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(int64(1), 10, 0).
		WillReturnRows(rows)

	posts, err := repo.ListByAuthor(ctx, 1, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, int64(1), posts[0].AuthorID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT (.+) FROM posts WHERE id = \$1`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	post, err := repo.GetByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
