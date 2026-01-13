package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	cosan "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/handlers"
	"github.com/toutaio/toutago-starter-kit-basic/internal/repositories"
	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

func TestAuthHandler_Register(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)
	authHandler := handlers.NewAuthHandler(authService)

	router := cosan.New()
	router.POST("/register", authHandler.Register)

	tests := []struct {
		name           string
		formData       url.Values
		wantStatusCode int
	}{
		{
			name: "valid registration",
			formData: url.Values{
				"email":    {"test@example.com"},
				"username": {"testuser"},
				"password": {"Test123!@#"},
			},
			wantStatusCode: http.StatusFound, // Redirect after success
		},
		{
			name: "weak password",
			formData: url.Values{
				"email":    {"test2@example.com"},
				"username": {"testuser2"},
				"password": {"weak"},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing fields",
			formData: url.Values{
				"email": {"test3@example.com"},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/register", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Register() status = %d, want %d, body = %s", w.Code, tt.wantStatusCode, w.Body.String())
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Register a user first
	authService.Register("test@example.com", "testuser", "Test123!@#", "", "")

	router := cosan.New()
	router.POST("/login", authHandler.Login)

	tests := []struct {
		name           string
		formData       url.Values
		wantStatusCode int
	}{
		{
			name: "valid login",
			formData: url.Values{
				"email":    {"test@example.com"},
				"password": {"Test123!@#"},
			},
			wantStatusCode: http.StatusFound,
		},
		{
			name: "wrong password",
			formData: url.Values{
				"email":    {"test@example.com"},
				"password": {"WrongPassword123!"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "non-existent user",
			formData: url.Values{
				"email":    {"notfound@example.com"},
				"password": {"Test123!@#"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Login() status = %d, want %d", w.Code, tt.wantStatusCode)
			}

			// Check that cookie is set on successful login
			if tt.wantStatusCode == http.StatusFound {
				cookies := w.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "session_id" {
						found = true
						break
					}
				}
				if !found {
					t.Error("Login() did not set session cookie")
				}
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Register and login
	authService.Register("test@example.com", "testuser", "Test123!@#", "", "")
	session, _ := authService.Login("test@example.com", "Test123!@#", "127.0.0.1", "Test Agent")

	router := cosan.New()
	router.POST("/logout", authHandler.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Logout() status = %d, want %d", w.Code, http.StatusFound)
	}

	// Check that session cookie is cleared
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_id" && cookie.MaxAge < 0 {
			return // Success - cookie marked for deletion
		}
	}
	t.Error("Logout() did not clear session cookie")
}
