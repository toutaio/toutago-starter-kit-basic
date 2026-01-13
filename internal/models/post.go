package models

import (
	"errors"
	"time"
)

// Content status constants
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

// Post represents a blog post
type Post struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`      // Markdown content
	Excerpt     string     `json:"excerpt"`      // Manual or auto-generated
	AuthorID    string     `json:"author_id"`
	CategoryID  *string    `json:"category_id"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
	DeletedAt   *time.Time `json:"deleted_at"` // Soft delete (trash)
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relations (not stored in DB)
	Author   *User      `json:"author,omitempty"`
	Category *Category  `json:"category,omitempty"`
	Tags     []Tag      `json:"tags,omitempty"`
}

// PostVersion stores historical versions of posts
type PostVersion struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	Version   int       `json:"version"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Excerpt   string    `json:"excerpt"`
	AuthorID  string    `json:"author_id"` // Who made this version
	CreatedAt time.Time `json:"created_at"`
}

// Category represents a post category
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag represents a post tag
type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostTag is the join table between posts and tags
type PostTag struct {
	PostID string `json:"post_id"`
	TagID  string `json:"tag_id"`
}

// Validate validates the post
func (p *Post) Validate() error {
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

// IsPublished checks if the post is published
func (p *Post) IsPublished() bool {
	return p.Status == StatusPublished && p.PublishedAt != nil
}

// IsDeleted checks if the post is in trash
func (p *Post) IsDeleted() bool {
	return p.DeletedAt != nil
}

// GenerateExcerpt creates an excerpt from content
func (p *Post) GenerateExcerpt(maxLength int) string {
	if len(p.Content) <= maxLength {
		return p.Content
	}
	return p.Content[:maxLength] + "..."
}

// Validate validates the post version
func (pv *PostVersion) Validate() error {
	if pv.PostID == "" {
		return errors.New("post_id is required")
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

// Validate validates the category
func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Slug == "" {
		return errors.New("slug is required")
	}
	return nil
}

// Validate validates the tag
func (t *Tag) Validate() error {
	if t.Name == "" {
		return errors.New("name is required")
	}
	if t.Slug == "" {
		return errors.New("slug is required")
	}
	return nil
}
