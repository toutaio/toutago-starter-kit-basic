package middleware

import (
	"net/http"

	cosan "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

const (
	// SessionCookieName is the name of the session cookie.
	SessionCookieName = "session_id"
	// UserContextKey is the context key for storing the authenticated user.
	UserContextKey = "auth_user"
)

// AuthMiddleware provides authentication middleware.
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware ensures the user is authenticated.
func (m *AuthMiddleware) RequireAuth(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(c cosan.Context) error {
		// Get session cookie
		cookie, err := c.Request().Cookie(SessionCookieName)
		if err != nil {
			http.Redirect(c.Response(), c.Request(), "/login", http.StatusFound)
			return nil
		}

		// Get user by session
		user, err := m.authService.GetUserBySession(cookie.Value)
		if err != nil {
			http.Redirect(c.Response(), c.Request(), "/login", http.StatusFound)
			return nil
		}

		// Store user in context
		c.Set(UserContextKey, user)

		return next(c)
	}
}

// RequireRole middleware ensures the user has the required role.
func (m *AuthMiddleware) RequireRole(role string) func(cosan.HandlerFunc) cosan.HandlerFunc {
	return func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(c cosan.Context) error {
			// Get session cookie
			cookie, err := c.Request().Cookie(SessionCookieName)
			if err != nil {
				http.Redirect(c.Response(), c.Request(), "/login", http.StatusFound)
				return nil
			}

			// Get user by session
			user, err := m.authService.GetUserBySession(cookie.Value)
			if err != nil {
				http.Redirect(c.Response(), c.Request(), "/login", http.StatusFound)
				return nil
			}

			// Check role
			if !hasRole(user, role) {
				c.Response().WriteHeader(http.StatusForbidden)
				c.Response().Write([]byte("Forbidden"))
				return nil
			}

			// Store user in context
			c.Set(UserContextKey, user)

			return next(c)
		}
	}
}

// OptionalAuth middleware loads the user if authenticated, but doesn't require it.
func (m *AuthMiddleware) OptionalAuth(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(c cosan.Context) error {
		// Get session cookie
		cookie, err := c.Request().Cookie(SessionCookieName)
		if err == nil {
			// Try to get user by session
			user, err := m.authService.GetUserBySession(cookie.Value)
			if err == nil {
				// Store user in context if found
				c.Set(UserContextKey, user)
			}
		}

		return next(c)
	}
}

// hasRole checks if user has the required role.
func hasRole(user *models.User, requiredRole string) bool {
	switch requiredRole {
	case models.RoleAdmin:
		return user.IsAdmin()
	case models.RoleEditor:
		return user.IsEditor()
	case models.RoleUser:
		return true // All authenticated users have user role
	default:
		return false
	}
}

// GetAuthUser retrieves the authenticated user from context.
func GetAuthUser(c cosan.Context) *models.User {
	user, ok := c.Get(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
