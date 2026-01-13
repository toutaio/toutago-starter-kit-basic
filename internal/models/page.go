package models

import (
	"errors"
	"time"
)

// Page represents a static page
type Page struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Content   string     `json:"content"` // Markdown content
	AuthorID  string     `json:"author_id"`
	Status    string     `json:"status"`
	Order     int        `json:"order"`      // For menu ordering
	DeletedAt *time.Time `json:"deleted_at"` // Soft delete (trash)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	
	// Relations (not stored in DB)
	Author *User `json:"author,omitempty"`
}

// PageVersion stores historical versions of pages
type PageVersion struct {
	ID        string    `json:"id"`
	PageID    string    `json:"page_id"`
	Version   int       `json:"version"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"` // Who made this version
	CreatedAt time.Time `json:"created_at"`
}

// Validate validates the page
func (p *Page) Validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Slug == "" {
		return errors.New("slug is required")
	}
	if p.Content == "" {
		return errors.New("content is required")
	}
	if p.AuthorID == "" {
		return errors.New("author is required")
	}
	if p.Status != StatusDraft && p.Status != StatusPublished && p.Status != StatusArchived {
		return errors.New("invalid status")
	}
	return nil
}

// IsPublished checks if the page is published
func (p *Page) IsPublished() bool {
	return p.Status == StatusPublished
}

// IsDeleted checks if the page is in trash
func (p *Page) IsDeleted() bool {
	return p.DeletedAt != nil
}

// Validate validates the page version
func (pv *PageVersion) Validate() error {
	if pv.PageID == "" {
		return errors.New("page_id is required")
	}
	if pv.Version < 1 {
		return errors.New("version must be positive")
	}
	if pv.Title == "" {
		return errors.New("title is required")
	}
	if pv.Content == "" {
		return errors.New("content is required")
	}
	if pv.AuthorID == "" {
		return errors.New("author_id is required")
	}
	return nil
}
