package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AllowedImageTypes defines the allowed image MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// MaxUploadSize is the maximum file size allowed (10MB)
const MaxUploadSize = 10 << 20 // 10 MB

// UploadConfig holds configuration for file uploads
type UploadConfig struct {
	UploadDir    string
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

// DefaultUploadConfig returns the default upload configuration
func DefaultUploadConfig() *UploadConfig {
	return &UploadConfig{
		UploadDir:    "./static/uploads",
		MaxFileSize:  MaxUploadSize,
		AllowedTypes: AllowedImageTypes,
	}
}

// SaveUploadedFile saves an uploaded file and returns its path
func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, config *UploadConfig) (string, error) {
	if config == nil {
		config = DefaultUploadConfig()
	}

	// Check file size
	if header.Size > config.MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", config.MaxFileSize)
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	if !config.AllowedTypes[contentType] {
		return "", fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// Create year/month directory structure
	now := time.Now()
	subDir := fmt.Sprintf("%d/%02d", now.Year(), now.Month())
	fullDir := filepath.Join(config.UploadDir, subDir)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Create destination file
	dst := filepath.Join(fullDir, filename)
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	// Copy file content
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path from static directory
	relativePath := filepath.Join("/uploads", subDir, filename)
	return strings.ReplaceAll(relativePath, "\\", "/"), nil
}

// DeleteUploadedFile deletes an uploaded file
func DeleteUploadedFile(path string) error {
	if path == "" {
		return nil
	}

	// Remove leading slash and prepend static directory
	path = strings.TrimPrefix(path, "/")
	fullPath := filepath.Join("./static", path)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ValidateImageFile validates an image file without saving it
func ValidateImageFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxUploadSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxUploadSize)
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return fmt.Errorf("file type %s is not allowed", contentType)
	}

	return nil
}
