package helpers

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple text",
			input: "Hello World",
			want:  "hello-world",
		},
		{
			name:  "with special characters",
			input: "Hello, World!",
			want:  "hello-world",
		},
		{
			name:  "with numbers",
			input: "Top 10 Tips",
			want:  "top-10-tips",
		},
		{
			name:  "with unicode",
			input: "Café Münchën",
			want:  "cafe-munchen",
		},
		{
			name:  "already lowercase",
			input: "already-a-slug",
			want:  "already-a-slug",
		},
		{
			name:  "multiple spaces",
			input: "Multiple   Spaces   Here",
			want:  "multiple-spaces-here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSlug(tt.input)
			if got != tt.want {
				t.Errorf("GenerateSlug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "plain text",
			markdown: "Hello",
			want:     "<p>Hello</p>\n",
		},
		{
			name:     "bold text",
			markdown: "**Bold**",
			want:     "<p><strong>Bold</strong></p>\n",
		},
		{
			name:     "heading",
			markdown: "# Heading 1",
			want:     "<h1 id=\"heading-1\">Heading 1</h1>\n",
		},
		{
			name:     "link",
			markdown: "[Link](https://example.com)",
			want:     "<p><a href=\"https://example.com\" target=\"_blank\">Link</a></p>\n",
		},
		{
			name:     "code",
			markdown: "`code`",
			want:     "<p><code>code</code></p>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderMarkdown(tt.markdown)
			if got != tt.want {
				t.Errorf("RenderMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "safe HTML",
			input: "<p>Hello</p>",
			want:  "<p>Hello</p>",
		},
		{
			name:  "with script tag",
			input: "<p>Hello</p><script>alert('xss')</script>",
			want:  "<p>Hello</p>",
		},
		{
			name:  "with onclick",
			input: "<p onclick='alert(1)'>Hello</p>",
			want:  "<p>Hello</p>",
		},
		{
			name:  "safe links",
			input: "<a href='https://example.com'>Link</a>",
			want:  "<a href=\"https://example.com\" rel=\"nofollow\">Link</a>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeHTML(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
