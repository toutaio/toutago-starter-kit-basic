// Package models defines domain models for the application.
package models

import (
	"errors"
	"regexp"
	"time"
)

// User roles
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleUser   = "user"
)

// User represents a user in the system.
type User struct {
	ID                         int        `db:"id" json:"id"`
	Email                      string     `db:"email" json:"email"`
	Username                   string     `db:"username" json:"username"`
	PasswordHash               string     `db:"password_hash" json:"-"`
	FirstName                  string     `db:"first_name" json:"first_name,omitempty"`
	LastName                   string     `db:"last_name" json:"last_name,omitempty"`
	Role                       string     `db:"role" json:"role"`
	EmailVerified              bool       `db:"email_verified" json:"email_verified"`
	VerificationToken          *string    `db:"verification_token" json:"-"`
	VerificationTokenExpiresAt *time.Time `db:"verification_token_expires_at" json:"-"`
	ResetToken                 *string    `db:"reset_token" json:"-"`
	ResetTokenExpiresAt        *time.Time `db:"reset_token_expires_at" json:"-"`
	LastLoginAt                *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt                  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt                  time.Time  `db:"updated_at" json:"updated_at"`
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Validate validates the user model.
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	if u.Username == "" {
		return errors.New("username is required")
	}

	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if !isValidRole(u.Role) {
		return errors.New("invalid role")
	}

	return nil
}

// IsAdmin returns true if the user is an admin.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEditor returns true if the user is an editor or admin.
func (u *User) IsEditor() bool {
	return u.Role == RoleEditor || u.Role == RoleAdmin
}

// FullName returns the user's full name.
func (u *User) FullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}

// Session represents a user session.
type Session struct {
	ID        string    `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	IPAddress string    `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent string    `db:"user_agent" json:"user_agent,omitempty"`
	Data      string    `db:"data" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// IsExpired returns true if the session has expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func isValidRole(role string) bool {
	return role == RoleAdmin || role == RoleEditor || role == RoleUser
}
