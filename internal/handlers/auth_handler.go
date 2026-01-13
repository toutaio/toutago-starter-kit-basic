package handlers

import (
	"net/http"
	"time"

	cosan "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(c cosan.Context) error {
	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		c.Response().Write([]byte("Invalid form data"))
		return nil
	}

	email := c.Request().FormValue("email")
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")
	firstName := c.Request().FormValue("first_name")
	lastName := c.Request().FormValue("last_name")

	// Validate required fields
	if email == "" || username == "" || password == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		c.Response().Write([]byte("Email, username, and password are required"))
		return nil
	}

	// Register user
	_, err := h.authService.Register(email, username, password, firstName, lastName)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		c.Response().Write([]byte(err.Error()))
		return nil
	}

	// Redirect to login page
	http.Redirect(c.Response(), c.Request(), "/login", http.StatusFound)
	return nil
}

// Login handles user login.
func (h *AuthHandler) Login(c cosan.Context) error {
	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		c.Response().Write([]byte("Invalid form data"))
		return nil
	}

	email := c.Request().FormValue("email")
	password := c.Request().FormValue("password")

	// Validate required fields
	if email == "" || password == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		c.Response().Write([]byte("Email and password are required"))
		return nil
	}

	// Get client IP and user agent
	ipAddress := c.Request().RemoteAddr
	userAgent := c.Request().UserAgent()

	// Login user
	session, err := h.authService.Login(email, password, ipAddress, userAgent)
	if err != nil {
		c.Response().WriteHeader(http.StatusUnauthorized)
		c.Response().Write([]byte("Invalid email or password"))
		return nil
	}

	// Set session cookie
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})

	// Redirect to dashboard or home
	http.Redirect(c.Response(), c.Request(), "/", http.StatusFound)
	return nil
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c cosan.Context) error {
	// Get session cookie
	cookie, err := c.Request().Cookie("session_id")
	if err == nil {
		// Delete session
		h.authService.Logout(cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	// Redirect to home
	http.Redirect(c.Response(), c.Request(), "/", http.StatusFound)
	return nil
}
