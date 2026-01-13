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

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostRepository) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) List(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Post), args.Error(1)
}

func (m *MockPostRepository) ListByStatus(ctx context.Context, status domain.PostStatus, limit, offset int) ([]*domain.Post, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Post), args.Error(1)
}

func (m *MockPostRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domain.Post, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Post), args.Error(1)
}

func TestPostService_CreatePost(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	post := &domain.Post{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "Content",
		AuthorID: 1,
		Status:   domain.PostStatusDraft,
	}

	repo.On("Create", ctx, post).Return(nil)

	err := service.CreatePost(ctx, post)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_CreatePost_ValidationError(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	tests := []struct {
		name string
		post *domain.Post
	}{
		{
			name: "empty title",
			post: &domain.Post{
				Title:    "",
				Slug:     "test",
				Content:  "Content",
				AuthorID: 1,
			},
		},
		{
			name: "empty slug",
			post: &domain.Post{
				Title:    "Test",
				Slug:     "",
				Content:  "Content",
				AuthorID: 1,
			},
		},
		{
			name: "empty content",
			post: &domain.Post{
				Title:    "Test",
				Slug:     "test",
				Content:  "",
				AuthorID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreatePost(ctx, tt.post)
			assert.Error(t, err)
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	expectedPost := &domain.Post{
		ID:      1,
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Content",
	}

	repo.On("GetByID", ctx, int64(1)).Return(expectedPost, nil)

	post, err := service.GetPostByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	repo.AssertExpectations(t)
}

func TestPostService_GetPostByID_NotFound(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, int64(999)).Return(nil, sql.ErrNoRows)

	post, err := service.GetPostByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, post)
	repo.AssertExpectations(t)
}

func TestPostService_GetPostBySlug(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	expectedPost := &domain.Post{
		ID:      1,
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Content",
	}

	repo.On("GetBySlug", ctx, "test-post").Return(expectedPost, nil)

	post, err := service.GetPostBySlug(ctx, "test-post")
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	repo.AssertExpectations(t)
}

func TestPostService_UpdatePost(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	post := &domain.Post{
		ID:       1,
		Title:    "Updated Post",
		Slug:     "updated-post",
		Content:  "Updated content",
		AuthorID: 1,
		Status:   domain.PostStatusPublished,
	}

	repo.On("Update", ctx, post).Return(nil)

	err := service.UpdatePost(ctx, post)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_DeletePost(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.DeletePost(ctx, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_ListPosts(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	expectedPosts := []*domain.Post{
		{ID: 1, Title: "Post 1", Slug: "post-1"},
		{ID: 2, Title: "Post 2", Slug: "post-2"},
	}

	repo.On("List", ctx, 10, 0).Return(expectedPosts, nil)

	posts, err := service.ListPosts(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	repo.AssertExpectations(t)
}

func TestPostService_ListPublishedPosts(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	expectedPosts := []*domain.Post{
		{ID: 1, Title: "Post 1", Status: domain.PostStatusPublished},
	}

	repo.On("ListByStatus", ctx, domain.PostStatusPublished, 10, 0).Return(expectedPosts, nil)

	posts, err := service.ListPublishedPosts(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	repo.AssertExpectations(t)
}

func TestPostService_PublishPost(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()
	now := time.Now()

	post := &domain.Post{
		ID:      1,
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Content",
		Status:  domain.PostStatusDraft,
	}

	repo.On("GetByID", ctx, int64(1)).Return(post, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(p *domain.Post) bool {
		return p.ID == 1 && p.Status == domain.PostStatusPublished && p.PublishedAt != nil
	})).Return(nil)

	err := service.PublishPost(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, domain.PostStatusPublished, post.Status)
	assert.NotNil(t, post.PublishedAt)
	assert.WithinDuration(t, now, *post.PublishedAt, time.Second)
	repo.AssertExpectations(t)
}

func TestPostService_UnpublishPost(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	publishedAt := time.Now()
	post := &domain.Post{
		ID:          1,
		Title:       "Test Post",
		Slug:        "test-post",
		Content:     "Content",
		Status:      domain.PostStatusPublished,
		PublishedAt: &publishedAt,
	}

	repo.On("GetByID", ctx, int64(1)).Return(post, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(p *domain.Post) bool {
		return p.ID == 1 && p.Status == domain.PostStatusDraft && p.PublishedAt == nil
	})).Return(nil)

	err := service.UnpublishPost(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, domain.PostStatusDraft, post.Status)
	assert.Nil(t, post.PublishedAt)
	repo.AssertExpectations(t)
}

func TestPostService_PublishPost_NotFound(t *testing.T) {
	repo := new(MockPostRepository)
	service := NewPostService(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, int64(999)).Return(nil, sql.ErrNoRows)

	err := service.PublishPost(ctx, 999)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
