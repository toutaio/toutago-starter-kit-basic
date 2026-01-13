package services

import (
	"errors"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/helpers"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/repositories"
)

var (
	// ErrInvalidCredentials is returned when login credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserNotVerified is returned when user email is not verified.
	ErrUserNotVerified = errors.New("email not verified")
)

// AuthService handles authentication operations.
type AuthService struct {
	userRepo     repositories.UserRepository
	sessionStore *SessionStore
}

// NewAuthService creates a new auth service.
func NewAuthService(userRepo repositories.UserRepository, sessionStore *SessionStore) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		sessionStore: sessionStore,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(email, username, password, firstName, lastName string) (*models.User, error) {
	// Validate password
	if err := helpers.ValidatePassword(password); err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:         email,
		Username:      username,
		PasswordHash:  passwordHash,
		FirstName:     firstName,
		LastName:      lastName,
		Role:          models.RoleUser,
		EmailVerified: false, // Will be verified via email
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and creates a session.
func (s *AuthService) Login(email, password, ipAddress, userAgent string) (*models.Session, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !helpers.ComparePasswords(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	// Create session (24 hour duration)
	session, err := s.sessionStore.Create(user.ID, ipAddress, userAgent, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Logout destroys a user session.
func (s *AuthService) Logout(sessionID string) error {
	return s.sessionStore.Delete(sessionID)
}

// GetUserBySession retrieves a user by session ID.
func (s *AuthService) GetUserBySession(sessionID string) (*models.User, error) {
	session, err := s.sessionStore.Get(sessionID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// VerifyEmail marks a user's email as verified.
func (s *AuthService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return err
	}

	// Check if token expired
	if user.VerificationTokenExpiresAt != nil && time.Now().After(*user.VerificationTokenExpiresAt) {
		return errors.New("verification token expired")
	}

	user.EmailVerified = true
	user.VerificationToken = nil
	user.VerificationTokenExpiresAt = nil

	return s.userRepo.Update(user)
}

// GeneratePasswordResetToken creates a password reset token.
func (s *AuthService) GeneratePasswordResetToken(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	// Generate token
	token, err := generateSessionID()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	user.ResetToken = &token
	user.ResetTokenExpiresAt = &expiresAt

	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword resets a user's password using a reset token.
func (s *AuthService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	// Check if token expired
	if user.ResetTokenExpiresAt != nil && time.Now().After(*user.ResetTokenExpiresAt) {
		return errors.New("reset token expired")
	}

	// Validate new password
	if err := helpers.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash password
	passwordHash, err := helpers.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.ResetToken = nil
	user.ResetTokenExpiresAt = nil

	return s.userRepo.Update(user)
}
