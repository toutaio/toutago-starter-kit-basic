package helpers_test

import (
	"testing"

	"github.com/toutaio/toutago-starter-kit-basic/internal/helpers"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Test123!@#",
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "Test1!",
			wantErr:  true,
		},
		{
			name:     "no uppercase",
			password: "test123!@#",
			wantErr:  true,
		},
		{
			name:     "no lowercase",
			password: "TEST123!@#",
			wantErr:  true,
		},
		{
			name:     "no number",
			password: "TestTest!@#",
			wantErr:  true,
		},
		{
			name:     "no special char",
			password: "Test123456",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helpers.ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "Test123!@#"

	hash, err := helpers.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword() returned plaintext password")
	}
}

func TestComparePasswords(t *testing.T) {
	password := "Test123!@#"
	hash, err := helpers.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "correct password",
			hash:     hash,
			password: password,
			want:     true,
		},
		{
			name:     "incorrect password",
			hash:     hash,
			password: "WrongPassword123!",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.ComparePasswords(tt.hash, tt.password); got != tt.want {
				t.Errorf("ComparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}
