package middleware

import "testing"

func TestCanEdit(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		userRole        string
		resourceOwnerID int64
		expected        bool
	}{
		{"Admin can edit any resource", 1, "admin", 2, true},
		{"Editor can edit any resource", 1, "editor", 2, true},
		{"Owner can edit own resource", 1, "user", 1, true},
		{"User cannot edit others' resource", 1, "user", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanEdit(tt.userID, tt.userRole, tt.resourceOwnerID)
			if result != tt.expected {
				t.Errorf("CanEdit() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCanDelete(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		userRole        string
		resourceOwnerID int64
		expected        bool
	}{
		{"Admin can delete any resource", 1, "admin", 2, true},
		{"Editor cannot delete others' resource", 1, "editor", 2, false},
		{"Owner can delete own resource", 1, "user", 1, true},
		{"User cannot delete others' resource", 1, "user", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanDelete(tt.userID, tt.userRole, tt.resourceOwnerID)
			if result != tt.expected {
				t.Errorf("CanDelete() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCanPublish(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		expected bool
	}{
		{"Admin can publish", "admin", true},
		{"Editor can publish", "editor", true},
		{"User cannot publish", "user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanPublish(tt.userRole)
			if result != tt.expected {
				t.Errorf("CanPublish() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCanManageUsers(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		expected bool
	}{
		{"Admin can manage users", "admin", true},
		{"Editor cannot manage users", "editor", false},
		{"User cannot manage users", "user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanManageUsers(tt.userRole)
			if result != tt.expected {
				t.Errorf("CanManageUsers() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
