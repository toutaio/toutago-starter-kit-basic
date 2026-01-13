package utils

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple title",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "title with special characters",
			input:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "title with multiple spaces",
			input:    "Hello   World   Test",
			expected: "hello-world-test",
		},
		{
			name:     "title with numbers",
			input:    "Top 10 Tips",
			expected: "top-10-tips",
		},
		{
			name:     "title with underscores",
			input:    "hello_world_test",
			expected: "hello-world-test",
		},
		{
			name:     "already a slug",
			input:    "hello-world",
			expected: "hello-world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
