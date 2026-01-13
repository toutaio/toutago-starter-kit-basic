// Package services provides business logic services.
package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
)

var (
	// ErrSessionNotFound is returned when a session is not found.
	ErrSessionNotFound = errors.New("session not found")
)

// SessionStore manages user sessions in memory.
type SessionStore struct {
	sessions map[string]*models.Session
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*models.Session),
	}
}

// Create creates a new session.
func (s *SessionStore) Create(userID int, ipAddress, userAgent string, duration time.Duration) (*models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.sessions[sessionID] = session
	return session, nil
}

// Get retrieves a session by ID.
func (s *SessionStore) Get(sessionID string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	if session.IsExpired() {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// Delete deletes a session by ID.
func (s *SessionStore) Delete(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
	return nil
}

// DeleteByUserID deletes all sessions for a user.
func (s *SessionStore) DeleteByUserID(userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		if session.UserID == userID {
			delete(s.sessions, id)
		}
	}

	return nil
}

// CleanupExpired removes expired sessions.
func (s *SessionStore) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, session := range s.sessions {
		if session.ExpiresAt.Before(now) {
			delete(s.sessions, id)
		}
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
