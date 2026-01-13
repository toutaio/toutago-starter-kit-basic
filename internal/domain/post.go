package domain

import "time"

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	AuthorID    int64      `json:"author_id"`
	Status      PostStatus `json:"status"`
	MetaTitle   string     `json:"meta_title"`
	MetaDesc    string     `json:"meta_desc"`
	IsFeatured  bool       `json:"is_featured"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (ps PostStatus) IsValid() bool {
	switch ps {
	case PostStatusDraft, PostStatusPublished, PostStatusArchived:
		return true
	}
	return false
}
