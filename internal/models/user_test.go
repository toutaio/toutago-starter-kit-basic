package models_test

import (
	"testing"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: models.User{
				Email:    "test@example.com",
				Username: "testuser",
				Role:     models.RoleUser,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: models.User{
				Email:    "invalid-email",
				Username: "testuser",
				Role:     models.RoleUser,
			},
			wantErr: true,
		},
		{
			name: "empty username",
			user: models.User{
				Email:    "test@example.com",
				Username: "",
				Role:     models.RoleUser,
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			user: models.User{
				Email:    "test@example.com",
				Username: "testuser",
				Role:     "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"admin role", models.RoleAdmin, true},
		{"editor role", models.RoleEditor, false},
		{"user role", models.RoleUser, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{Role: tt.role}
			if got := u.IsAdmin(); got != tt.want {
				t.Errorf("User.IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsEditor(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"admin role", models.RoleAdmin, true},
		{"editor role", models.RoleEditor, true},
		{"user role", models.RoleUser, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{Role: tt.role}
			if got := u.IsEditor(); got != tt.want {
				t.Errorf("User.IsEditor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "expired session",
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      true,
		},
		{
			name:      "valid session",
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.Session{ExpiresAt: tt.expiresAt}
			if got := s.IsExpired(); got != tt.want {
				t.Errorf("Session.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
