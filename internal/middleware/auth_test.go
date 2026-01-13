package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	cosan "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/middleware"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/repositories"
	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)

	// Register and login a user
	authService.Register("test@example.com", "testuser", "Test123!@#", "", "")
	session, _ := authService.Login("test@example.com", "Test123!@#", "127.0.0.1", "Test Agent")

	authMiddleware := middleware.NewAuthMiddleware(authService)

	tests := []struct {
		name           string
		sessionCookie  string
		wantStatusCode int
	}{
		{
			name:           "authenticated user",
			sessionCookie:  session.ID,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "no session cookie",
			sessionCookie:  "",
			wantStatusCode: http.StatusFound, // Redirect to login
		},
		{
			name:           "invalid session",
			sessionCookie:  "invalid-session-id",
			wantStatusCode: http.StatusFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := cosan.New()

			// Apply middleware and add test route
			router.Use(cosan.MiddlewareFunc(authMiddleware.RequireAuth))
			router.GET("/protected", func(c cosan.Context) error {
				c.Response().WriteHeader(http.StatusOK)
				c.Response().Write([]byte("Protected resource"))
				return nil
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.sessionCookie != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: tt.sessionCookie,
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("RequireAuth() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAuthMiddleware_RequireRole(t *testing.T) {
	userRepo := repositories.NewMemoryUserRepository()
	sessionStore := services.NewSessionStore()
	authService := services.NewAuthService(userRepo, sessionStore)

	// Create users with different roles
	user, _ := authService.Register("admin@example.com", "admin", "Test123!@#", "", "")
	user.Role = models.RoleAdmin
	userRepo.Update(user)
	adminSession, _ := authService.Login("admin@example.com", "Test123!@#", "127.0.0.1", "Test Agent")

	authService.Register("user@example.com", "regularuser", "Test123!@#", "", "")
	userSession, _ := authService.Login("user@example.com", "Test123!@#", "127.0.0.1", "Test Agent")

	authMiddleware := middleware.NewAuthMiddleware(authService)

	tests := []struct {
		name           string
		sessionCookie  string
		requiredRole   string
		wantStatusCode int
	}{
		{
			name:           "admin accessing admin route",
			sessionCookie:  adminSession.ID,
			requiredRole:   models.RoleAdmin,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "user accessing admin route",
			sessionCookie:  userSession.ID,
			requiredRole:   models.RoleAdmin,
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:           "no session",
			sessionCookie:  "",
			requiredRole:   models.RoleAdmin,
			wantStatusCode: http.StatusFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := cosan.New()

			router.Use(cosan.MiddlewareFunc(authMiddleware.RequireRole(tt.requiredRole)))
			router.GET("/admin", func(c cosan.Context) error {
				c.Response().WriteHeader(http.StatusOK)
				c.Response().Write([]byte("Admin resource"))
				return nil
			})

			req := httptest.NewRequest("GET", "/admin", nil)
			if tt.sessionCookie != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: tt.sessionCookie,
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("RequireRole() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
