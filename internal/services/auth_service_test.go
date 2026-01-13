package services_test

import (
	"testing"

	"github.com/toutaio/toutago-starter-kit-basic/internal/repositories"
	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

func TestAuthService_Register(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)

	tests := []struct {
		name     string
		email    string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid registration",
			email:    "test@example.com",
			username: "testuser",
			password: "Test123!@#",
			wantErr:  false,
		},
		{
			name:     "weak password",
			email:    "test2@example.com",
			username: "testuser2",
			password: "weak",
			wantErr:  true,
		},
		{
			name:     "invalid email",
			email:    "invalid-email",
			username: "testuser3",
			password: "Test123!@#",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.Register(tt.email, tt.username, tt.password, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("Register() returned nil user")
				}
				if user.EmailVerified {
					t.Error("Register() should not verify email automatically")
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)

	// Register a user first
	authService.Register("test@example.com", "testuser", "Test123!@#", "", "")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "correct credentials",
			email:    "test@example.com",
			password: "Test123!@#",
			wantErr:  false,
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "WrongPassword123!",
			wantErr:  true,
		},
		{
			name:     "non-existent user",
			email:    "notfound@example.com",
			password: "Test123!@#",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := authService.Login(tt.email, tt.password, "127.0.0.1", "Test Agent")
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && session == nil {
				t.Error("Login() returned nil session")
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)

	authService.Register("test@example.com", "testuser", "Test123!@#", "", "")
	session, _ := authService.Login("test@example.com", "Test123!@#", "127.0.0.1", "Test Agent")

	err := authService.Logout(session.ID)
	if err != nil {
		t.Errorf("Logout() error = %v", err)
	}

	// Session should no longer exist
	_, err = sessionStore.Get(session.ID)
	if err == nil {
		t.Error("Logout() did not delete session")
	}
}

func TestAuthService_GetUserBySession(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)

	authService.Register("test@example.com", "testuser", "Test123!@#", "", "")
	session, _ := authService.Login("test@example.com", "Test123!@#", "127.0.0.1", "Test Agent")

	user, err := authService.GetUserBySession(session.ID)
	if err != nil {
		t.Errorf("GetUserBySession() error = %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("GetUserBySession() email = %s, want test@example.com", user.Email)
	}
}
