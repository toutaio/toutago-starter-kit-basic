// Package helpers provides utility functions.
package helpers

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Password validation constants
const (
	MinPasswordLength = 8
)

var (
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	numberRegex    = regexp.MustCompile(`[0-9]`)
	specialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// ValidatePassword validates password complexity requirements.
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < MinPasswordLength {
		return errors.New("password must be at least 8 characters")
	}

	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if !specialRegex.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePasswords compares a hashed password with a plaintext password.
func ComparePasswords(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
