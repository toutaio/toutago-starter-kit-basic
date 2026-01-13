package helpers

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	return slug.Make(text)
}

// RenderMarkdown converts markdown to HTML
func RenderMarkdown(md string) string {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// Create HTML renderer with options
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

// SanitizeHTML removes potentially dangerous HTML elements and attributes
func SanitizeHTML(rawHTML string) string {
	// Use UGC (User Generated Content) policy which allows safe HTML
	policy := bluemonday.UGCPolicy()
	return policy.Sanitize(rawHTML)
}

// RenderMarkdownSafe renders markdown and sanitizes the output
func RenderMarkdownSafe(md string) string {
	html := RenderMarkdown(md)
	return SanitizeHTML(html)
}
