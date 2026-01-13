package models

import (
	"testing"
	"time"
)

func TestPost_Validate(t *testing.T) {
	tests := []struct {
		name    string
		post    Post
		wantErr bool
	}{
		{
			name: "valid post",
			post: Post{
				Title:    "Test Post",
				Slug:     "test-post",
				Content:  "This is test content",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			post: Post{
				Slug:     "test-post",
				Content:  "Content",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "missing slug",
			post: Post{
				Title:    "Test Post",
				Content:  "Content",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "missing content",
			post: Post{
				Title:    "Test Post",
				Slug:     "test-post",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "missing author",
			post: Post{
				Title:   "Test Post",
				Slug:    "test-post",
				Content: "Content",
				Status:  StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			post: Post{
				Title:    "Test Post",
				Slug:     "test-post",
				Content:  "Content",
				AuthorID: "user123",
				Status:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPost_IsPublished(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		post Post
		want bool
	}{
		{
			name: "published status with date",
			post: Post{
				Status:      StatusPublished,
				PublishedAt: &now,
			},
			want: true,
		},
		{
			name: "published status without date",
			post: Post{
				Status:      StatusPublished,
				PublishedAt: nil,
			},
			want: false,
		},
		{
			name: "draft status",
			post: Post{
				Status:      StatusDraft,
				PublishedAt: nil,
			},
			want: false,
		},
		{
			name: "archived status",
			post: Post{
				Status:      StatusArchived,
				PublishedAt: &now,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.post.IsPublished(); got != tt.want {
				t.Errorf("Post.IsPublished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost_GenerateExcerpt(t *testing.T) {
	tests := []struct {
		name    string
		content string
		length  int
		want    string
	}{
		{
			name:    "short content",
			content: "Short content",
			length:  200,
			want:    "Short content",
		},
		{
			name:    "long content truncated",
			content: "This is a very long piece of content that should be truncated to the specified length and should end with an ellipsis to indicate that there is more content available.",
			length:  50,
			want:    "This is a very long piece of content that should b...",
		},
		{
			name:    "exactly at length",
			content: "This is exactly fifty characters long for test.",
			length:  50,
			want:    "This is exactly fifty characters long for test.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{Content: tt.content}
			got := post.GenerateExcerpt(tt.length)
			if got != tt.want {
				t.Errorf("Post.GenerateExcerpt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		version PostVersion
		wantErr bool
	}{
		{
			name: "valid version",
			version: PostVersion{
				PostID:   "post123",
				Version:  1,
				Title:    "Test",
				Content:  "Content",
				AuthorID: "user123",
			},
			wantErr: false,
		},
		{
			name: "missing post id",
			version: PostVersion{
				Version:  1,
				Title:    "Test",
				Content:  "Content",
				AuthorID: "user123",
			},
			wantErr: true,
		},
		{
			name: "invalid version number",
			version: PostVersion{
				PostID:   "post123",
				Version:  0,
				Title:    "Test",
				Content:  "Content",
				AuthorID: "user123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.version.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PostVersion.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		wantErr  bool
	}{
		{
			name: "valid category",
			category: Category{
				Name: "Technology",
				Slug: "technology",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			category: Category{
				Slug: "technology",
			},
			wantErr: true,
		},
		{
			name: "missing slug",
			category: Category{
				Name: "Technology",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Category.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTag_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tag     Tag
		wantErr bool
	}{
		{
			name: "valid tag",
			tag: Tag{
				Name: "golang",
				Slug: "golang",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			tag: Tag{
				Slug: "golang",
			},
			wantErr: true,
		},
		{
			name: "missing slug",
			tag: Tag{
				Name: "golang",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tag.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tag.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
