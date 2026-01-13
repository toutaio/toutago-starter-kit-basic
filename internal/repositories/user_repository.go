// Package repositories provides data access layer.
package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
)

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists is returned when a user already exists.
	ErrUserExists = errors.New("user already exists")
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByVerificationToken(token string) (*models.User, error)
	FindByResetToken(token string) (*models.User, error)
}

// MemoryUserRepository implements UserRepository in memory.
type MemoryUserRepository struct {
	users  map[int]*models.User
	nextID int
	mu     sync.RWMutex
}

// NewMemoryUserRepository creates a new memory-based user repository.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

// Create creates a new user.
func (r *MemoryUserRepository) Create(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrUserExists
		}
		if u.Username == user.Username {
			return ErrUserExists
		}
	}

	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

// Update updates an existing user.
func (r *MemoryUserRepository) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// FindByID finds a user by ID.
func (r *MemoryUserRepository) FindByID(id int) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindByEmail finds a user by email.
func (r *MemoryUserRepository) FindByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// FindByUsername finds a user by username.
func (r *MemoryUserRepository) FindByUsername(username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// FindByVerificationToken finds a user by verification token.
func (r *MemoryUserRepository) FindByVerificationToken(token string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.VerificationToken != nil && *user.VerificationToken == token {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// FindByResetToken finds a user by reset token.
func (r *MemoryUserRepository) FindByResetToken(token string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ResetToken != nil && *user.ResetToken == token {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}
