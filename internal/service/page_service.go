package service

import (
	"context"
	"errors"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
)

type PageRepository interface {
	Create(ctx context.Context, page *domain.Page) error
	GetByID(ctx context.Context, id int64) (*domain.Page, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Page, error)
	Update(ctx context.Context, page *domain.Page) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*domain.Page, error)
	ListByStatus(ctx context.Context, status domain.PageStatus, limit, offset int) ([]*domain.Page, error)
}

type PageService struct {
	repo PageRepository
}

func NewPageService(repo PageRepository) *PageService {
	return &PageService{repo: repo}
}

func (s *PageService) CreatePage(ctx context.Context, page *domain.Page) error {
	if err := s.validatePage(page); err != nil {
		return err
	}

	// Check slug uniqueness
	existing, err := s.repo.GetBySlug(ctx, page.Slug)
	if err == nil && existing != nil {
		return errors.New("slug already exists")
	}

	if page.Status == "" {
		page.Status = domain.PageStatusDraft
	}

	return s.repo.Create(ctx, page)
}

func (s *PageService) GetPageByID(ctx context.Context, id int64) (*domain.Page, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PageService) GetPageBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *PageService) UpdatePage(ctx context.Context, page *domain.Page) error {
	if err := s.validatePage(page); err != nil {
		return err
	}

	// Check slug uniqueness (excluding current page)
	existing, err := s.repo.GetBySlug(ctx, page.Slug)
	if err == nil && existing != nil && existing.ID != page.ID {
		return errors.New("slug already exists")
	}

	return s.repo.Update(ctx, page)
}

func (s *PageService) DeletePage(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *PageService) ListPages(ctx context.Context, limit, offset int) ([]*domain.Page, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}

func (s *PageService) ListPublishedPages(ctx context.Context, limit, offset int) ([]*domain.Page, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListByStatus(ctx, domain.PageStatusPublished, limit, offset)
}

func (s *PageService) PublishPage(ctx context.Context, id int64) error {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	page.Status = domain.PageStatusPublished
	page.PublishedAt = &now

	return s.repo.Update(ctx, page)
}

func (s *PageService) UnpublishPage(ctx context.Context, id int64) error {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	page.Status = domain.PageStatusDraft
	page.PublishedAt = nil

	return s.repo.Update(ctx, page)
}

func (s *PageService) validatePage(page *domain.Page) error {
	if page.Title == "" {
		return errors.New("title is required")
	}
	if page.Slug == "" {
		return errors.New("slug is required")
	}
	if page.Content == "" {
		return errors.New("content is required")
	}

	return nil
}
