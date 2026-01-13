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

func TestPageRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()
	now := time.Now()

	page := &domain.Page{
		Title:     "Test Page",
		Slug:      "test-page",
		Content:   "This is test content",
		Status:    domain.PageStatusDraft,
		MetaTitle: "Test Meta Title",
		MetaDesc:  "Test meta description",
	}

	mock.ExpectQuery(`INSERT INTO pages`).
		WithArgs(page.Title, page.Slug, page.Content, page.Status, page.MetaTitle, page.MetaDesc, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	err = repo.Create(ctx, page)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), page.ID)
	assert.NotZero(t, page.CreatedAt)
	assert.NotZero(t, page.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "status",
		"meta_title", "meta_desc", "published_at", "created_at", "updated_at",
	}).AddRow(
		1, "Test Page", "test-page", "Content", domain.PageStatusPublished,
		"Meta Title", "Meta Desc", now, now, now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM pages WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	page, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, int64(1), page.ID)
	assert.Equal(t, "Test Page", page.Title)
	assert.Equal(t, "test-page", page.Slug)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_GetBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "status",
		"meta_title", "meta_desc", "published_at", "created_at", "updated_at",
	}).AddRow(
		1, "Test Page", "test-page", "Content", domain.PageStatusPublished,
		"Meta Title", "Meta Desc", now, now, now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM pages WHERE slug = \$1`).
		WithArgs("test-page").
		WillReturnRows(rows)

	page, err := repo.GetBySlug(ctx, "test-page")
	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "test-page", page.Slug)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()

	page := &domain.Page{
		ID:        1,
		Title:     "Updated Page",
		Slug:      "updated-page",
		Content:   "Updated content",
		Status:    domain.PageStatusPublished,
		MetaTitle: "Updated Meta",
		MetaDesc:  "Updated desc",
	}

	mock.ExpectExec(`UPDATE pages SET`).
		WithArgs(page.Title, page.Slug, page.Content, page.Status, page.MetaTitle, page.MetaDesc, sqlmock.AnyArg(), page.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(ctx, page)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM pages WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "status",
		"meta_title", "meta_desc", "published_at", "created_at", "updated_at",
	}).
		AddRow(1, "Page 1", "page-1", "Content 1", domain.PageStatusPublished, "Meta 1", "Desc 1", now, now, now).
		AddRow(2, "Page 2", "page-2", "Content 2", domain.PageStatusPublished, "Meta 2", "Desc 2", now, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM pages ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	pages, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, pages, 2)
	assert.Equal(t, "Page 1", pages[0].Title)
	assert.Equal(t, "Page 2", pages[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_ListByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "content", "status",
		"meta_title", "meta_desc", "published_at", "created_at", "updated_at",
	}).AddRow(1, "Page 1", "page-1", "Content 1", domain.PageStatusPublished, "Meta 1", "Desc 1", now, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM pages WHERE status = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(domain.PageStatusPublished, 10, 0).
		WillReturnRows(rows)

	pages, err := repo.ListByStatus(ctx, domain.PageStatusPublished, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, pages, 1)
	assert.Equal(t, domain.PageStatusPublished, pages[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPageRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPageRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT (.+) FROM pages WHERE id = \$1`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	page, err := repo.GetByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, page)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
