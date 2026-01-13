package models

import (
	"testing"
)

func TestPage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		page    Page
		wantErr bool
	}{
		{
			name: "valid page",
			page: Page{
				Title:    "About Us",
				Slug:     "about-us",
				Content:  "This is the about page",
				AuthorID: "user123",
				Status:   StatusDraft,
				Order:    1,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			page: Page{
				Slug:     "about",
				Content:  "Content",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "missing slug",
			page: Page{
				Title:    "About",
				Content:  "Content",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "missing content",
			page: Page{
				Title:    "About",
				Slug:     "about",
				AuthorID: "user123",
				Status:   StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "missing author",
			page: Page{
				Title:   "About",
				Slug:    "about",
				Content: "Content",
				Status:  StatusDraft,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			page: Page{
				Title:    "About",
				Slug:     "about",
				Content:  "Content",
				AuthorID: "user123",
				Status:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.page.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Page.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPage_IsPublished(t *testing.T) {
	tests := []struct {
		name string
		page Page
		want bool
	}{
		{
			name: "published status",
			page: Page{
				Status: StatusPublished,
			},
			want: true,
		},
		{
			name: "draft status",
			page: Page{
				Status: StatusDraft,
			},
			want: false,
		},
		{
			name: "archived status",
			page: Page{
				Status: StatusArchived,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.page.IsPublished(); got != tt.want {
				t.Errorf("Page.IsPublished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		version PageVersion
		wantErr bool
	}{
		{
			name: "valid version",
			version: PageVersion{
				PageID:   "page123",
				Version:  1,
				Title:    "Test",
				Content:  "Content",
				AuthorID: "user123",
			},
			wantErr: false,
		},
		{
			name: "missing page id",
			version: PageVersion{
				Version:  1,
				Title:    "Test",
				Content:  "Content",
				AuthorID: "user123",
			},
			wantErr: true,
		},
		{
			name: "invalid version number",
			version: PageVersion{
				PageID:   "page123",
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
				t.Errorf("PageVersion.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
