package domain

import "time"

type PageStatus string

const (
	PageStatusDraft     PageStatus = "draft"
	PageStatusPublished PageStatus = "published"
	PageStatusArchived  PageStatus = "archived"
)

type Page struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	AuthorID    int64      `json:"author_id"`
	Status      PageStatus `json:"status"`
	MetaTitle   string     `json:"meta_title"`
	MetaDesc    string     `json:"meta_desc"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (ps PageStatus) IsValid() bool {
	switch ps {
	case PageStatusDraft, PageStatusPublished, PageStatusArchived:
		return true
	}
	return false
}
