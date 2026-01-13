package repositories_test

import (
	"testing"

	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/repositories"
)

func TestUserRepository_Create(t *testing.T) {
	repo := repositories.NewMemoryUserRepository()

	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleUser,
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.ID == 0 {
		t.Error("Create() did not set user ID")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo := repositories.NewMemoryUserRepository()

	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleUser,
	}
	repo.Create(user)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "non-existent user",
			email:   "notfound@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && found == nil {
				t.Error("FindByEmail() returned nil user")
			}
		})
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	repo := repositories.NewMemoryUserRepository()

	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleUser,
	}
	repo.Create(user)

	found, err := repo.FindByUsername("testuser")
	if err != nil {
		t.Fatalf("FindByUsername() error = %v", err)
	}

	if found.Username != "testuser" {
		t.Errorf("FindByUsername() username = %s, want testuser", found.Username)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := repositories.NewMemoryUserRepository()

	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleUser,
	}
	repo.Create(user)

	user.FirstName = "Test"
	user.LastName = "User"

	err := repo.Update(user)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	found, _ := repo.FindByID(user.ID)
	if found.FirstName != "Test" {
		t.Errorf("Update() first_name = %s, want Test", found.FirstName)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := repositories.NewMemoryUserRepository()

	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleUser,
	}
	repo.Create(user)

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("FindByID() id = %d, want %d", found.ID, user.ID)
	}
}
