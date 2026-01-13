package service

import (
	"context"
	"errors"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	GetByID(ctx context.Context, id int64) (*domain.Post, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Post, error)
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*domain.Post, error)
	ListByStatus(ctx context.Context, status domain.PostStatus, limit, offset int) ([]*domain.Post, error)
	ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domain.Post, error)
}

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) error {
	if err := s.validatePost(post); err != nil {
		return err
	}

	// Check slug uniqueness
	existing, err := s.repo.GetBySlug(ctx, post.Slug)
	if err == nil && existing != nil {
		return errors.New("slug already exists")
	}

	if post.Status == "" {
		post.Status = domain.PostStatusDraft
	}

	return s.repo.Create(ctx, post)
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) GetPostBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *PostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	if err := s.validatePost(post); err != nil {
		return err
	}

	// Check slug uniqueness (excluding current post)
	existing, err := s.repo.GetBySlug(ctx, post.Slug)
	if err == nil && existing != nil && existing.ID != post.ID {
		return errors.New("slug already exists")
	}

	return s.repo.Update(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *PostService) ListPosts(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}

func (s *PostService) ListPublishedPosts(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListByStatus(ctx, domain.PostStatusPublished, limit, offset)
}

func (s *PostService) ListPostsByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domain.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListByAuthor(ctx, authorID, limit, offset)
}

func (s *PostService) PublishPost(ctx context.Context, id int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	post.Status = domain.PostStatusPublished
	post.PublishedAt = &now

	return s.repo.Update(ctx, post)
}

func (s *PostService) UnpublishPost(ctx context.Context, id int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	post.Status = domain.PostStatusDraft
	post.PublishedAt = nil

	return s.repo.Update(ctx, post)
}

func (s *PostService) validatePost(post *domain.Post) error {
	if post.Title == "" {
		return errors.New("title is required")
	}
	if post.Slug == "" {
		return errors.New("slug is required")
	}
	if post.Content == "" {
		return errors.New("content is required")
	}

	return nil
}
