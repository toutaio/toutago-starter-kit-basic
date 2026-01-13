package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
)

type MockPageRepository struct {
	mock.Mock
}

func (m *MockPageRepository) Create(ctx context.Context, page *domain.Page) error {
	args := m.Called(ctx, page)
	return args.Error(0)
}

func (m *MockPageRepository) GetByID(ctx context.Context, id int64) (*domain.Page, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Page), args.Error(1)
}

func (m *MockPageRepository) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Page), args.Error(1)
}

func (m *MockPageRepository) Update(ctx context.Context, page *domain.Page) error {
	args := m.Called(ctx, page)
	return args.Error(0)
}

func (m *MockPageRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Page, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Page), args.Error(1)
}

func (m *MockPageRepository) ListByStatus(ctx context.Context, status domain.PageStatus, limit, offset int) ([]*domain.Page, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Page), args.Error(1)
}

func TestPageService_CreatePage(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	page := &domain.Page{
		Title:   "Test Page",
		Slug:    "test-page",
		Content: "Content",
		Status:  domain.PageStatusDraft,
	}

	repo.On("Create", ctx, page).Return(nil)

	err := service.CreatePage(ctx, page)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPageService_CreatePage_ValidationError(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	tests := []struct {
		name string
		page *domain.Page
	}{
		{
			name: "empty title",
			page: &domain.Page{
				Title:   "",
				Slug:    "test",
				Content: "Content",
			},
		},
		{
			name: "empty slug",
			page: &domain.Page{
				Title:   "Test",
				Slug:    "",
				Content: "Content",
			},
		},
		{
			name: "empty content",
			page: &domain.Page{
				Title:   "Test",
				Slug:    "test",
				Content: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreatePage(ctx, tt.page)
			assert.Error(t, err)
		})
	}
}

func TestPageService_GetPageByID(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	expectedPage := &domain.Page{
		ID:      1,
		Title:   "Test Page",
		Slug:    "test-page",
		Content: "Content",
	}

	repo.On("GetByID", ctx, int64(1)).Return(expectedPage, nil)

	page, err := service.GetPageByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedPage, page)
	repo.AssertExpectations(t)
}

func TestPageService_GetPageByID_NotFound(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, int64(999)).Return(nil, sql.ErrNoRows)

	page, err := service.GetPageByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, page)
	repo.AssertExpectations(t)
}

func TestPageService_GetPageBySlug(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	expectedPage := &domain.Page{
		ID:      1,
		Title:   "Test Page",
		Slug:    "test-page",
		Content: "Content",
	}

	repo.On("GetBySlug", ctx, "test-page").Return(expectedPage, nil)

	page, err := service.GetPageBySlug(ctx, "test-page")
	assert.NoError(t, err)
	assert.Equal(t, expectedPage, page)
	repo.AssertExpectations(t)
}

func TestPageService_UpdatePage(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	page := &domain.Page{
		ID:      1,
		Title:   "Updated Page",
		Slug:    "updated-page",
		Content: "Updated content",
		Status:  domain.PageStatusPublished,
	}

	repo.On("Update", ctx, page).Return(nil)

	err := service.UpdatePage(ctx, page)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPageService_DeletePage(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.DeletePage(ctx, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPageService_ListPages(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	expectedPages := []*domain.Page{
		{ID: 1, Title: "Page 1", Slug: "page-1"},
		{ID: 2, Title: "Page 2", Slug: "page-2"},
	}

	repo.On("List", ctx, 10, 0).Return(expectedPages, nil)

	pages, err := service.ListPages(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedPages, pages)
	repo.AssertExpectations(t)
}

func TestPageService_ListPublishedPages(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	expectedPages := []*domain.Page{
		{ID: 1, Title: "Page 1", Status: domain.PageStatusPublished},
	}

	repo.On("ListByStatus", ctx, domain.PageStatusPublished, 10, 0).Return(expectedPages, nil)

	pages, err := service.ListPublishedPages(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedPages, pages)
	repo.AssertExpectations(t)
}

func TestPageService_PublishPage(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()
	now := time.Now()

	page := &domain.Page{
		ID:      1,
		Title:   "Test Page",
		Slug:    "test-page",
		Content: "Content",
		Status:  domain.PageStatusDraft,
	}

	repo.On("GetByID", ctx, int64(1)).Return(page, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(p *domain.Page) bool {
		return p.ID == 1 && p.Status == domain.PageStatusPublished && p.PublishedAt != nil
	})).Return(nil)

	err := service.PublishPage(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, domain.PageStatusPublished, page.Status)
	assert.NotNil(t, page.PublishedAt)
	assert.WithinDuration(t, now, *page.PublishedAt, time.Second)
	repo.AssertExpectations(t)
}

func TestPageService_UnpublishPage(t *testing.T) {
	repo := new(MockPageRepository)
	service := NewPageService(repo)
	ctx := context.Background()

	publishedAt := time.Now()
	page := &domain.Page{
		ID:          1,
		Title:       "Test Page",
		Slug:        "test-page",
		Content:     "Content",
		Status:      domain.PageStatusPublished,
		PublishedAt: &publishedAt,
	}

	repo.On("GetByID", ctx, int64(1)).Return(page, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(p *domain.Page) bool {
		return p.ID == 1 && p.Status == domain.PageStatusDraft && p.PublishedAt == nil
	})).Return(nil)

	err := service.UnpublishPage(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, domain.PageStatusDraft, page.Status)
	assert.Nil(t, page.PublishedAt)
	repo.AssertExpectations(t)
}
