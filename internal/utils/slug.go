package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug converts a string into a URL-friendly slug
func GenerateSlug(s string) string {
	// Convert to lowercase
	slug := strings.ToLower(s)
	
	// Replace underscores with hyphens
	slug = strings.ReplaceAll(slug, "_", "-")
	
	// Remove all non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "-")
	
	// Replace multiple hyphens with single hyphen
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}
