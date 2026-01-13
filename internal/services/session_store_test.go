package services_test

import (
	"testing"
	"time"

	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

func TestSessionStore_Create(t *testing.T) {
	store := services.NewSessionStore()

	session, err := store.Create(1, "127.0.0.1", "Test Agent", 24*time.Hour)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if session.ID == "" {
		t.Error("Create() returned session with empty ID")
	}

	if session.UserID != 1 {
		t.Errorf("Create() session.UserID = %d, want 1", session.UserID)
	}

	if session.IsExpired() {
		t.Error("Create() returned expired session")
	}
}

func TestSessionStore_Get(t *testing.T) {
	store := services.NewSessionStore()

	created, _ := store.Create(1, "127.0.0.1", "Test Agent", 24*time.Hour)

	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{
			name:      "existing session",
			sessionID: created.ID,
			wantErr:   false,
		},
		{
			name:      "non-existent session",
			sessionID: "non-existent",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := store.Get(tt.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && session == nil {
				t.Error("Get() returned nil session")
			}
		})
	}
}

func TestSessionStore_Delete(t *testing.T) {
	store := services.NewSessionStore()

	session, _ := store.Create(1, "127.0.0.1", "Test Agent", 24*time.Hour)

	err := store.Delete(session.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	_, err = store.Get(session.ID)
	if err == nil {
		t.Error("Get() after Delete() should return error")
	}
}

func TestSessionStore_DeleteByUserID(t *testing.T) {
	store := services.NewSessionStore()

	// Create multiple sessions for same user
	store.Create(1, "127.0.0.1", "Agent 1", 24*time.Hour)
	store.Create(1, "127.0.0.2", "Agent 2", 24*time.Hour)

	err := store.DeleteByUserID(1)
	if err != nil {
		t.Errorf("DeleteByUserID() error = %v", err)
	}
}

func TestSessionStore_CleanupExpired(t *testing.T) {
	store := services.NewSessionStore()

	// Create expired session
	store.Create(1, "127.0.0.1", "Test Agent", -1*time.Hour)

	// Create valid session
	validSession, _ := store.Create(2, "127.0.0.2", "Test Agent", 24*time.Hour)

	store.CleanupExpired()

	// Valid session should still exist
	_, err := store.Get(validSession.ID)
	if err != nil {
		t.Error("CleanupExpired() deleted valid session")
	}
}
