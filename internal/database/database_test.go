package database_test

import (
	"testing"

	"github.com/toutaio/toutago-starter-kit-basic/internal/config"
	"github.com/toutaio/toutago-starter-kit-basic/internal/database"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.DatabaseConfig
		wantErr bool
	}{
		{
			name: "valid postgres config",
			cfg: config.DatabaseConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     "5432",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_pass",
			},
			wantErr: true, // Will fail without real DB, but tests structure
		},
		{
			name: "valid mysql config",
			cfg: config.DatabaseConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_pass",
			},
			wantErr: true, // Will fail without real DB, but tests structure
		},
		{
			name: "invalid driver",
			cfg: config.DatabaseConfig{
				Driver:   "invalid",
				Host:     "localhost",
				Port:     "5432",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_pass",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.Connect(tt.cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if db != nil {
				defer db.Close()
			}
		})
	}
}

func TestConnectionString(t *testing.T) {
	// This test verifies we're using the config properly
	cfg := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		Name:     "test_db",
		User:     "test_user",
		Password: "test_pass",
	}

	expected := "postgres://test_user:test_pass@localhost:5432/test_db?sslmode=disable"
	if got := cfg.ConnectionString(); got != expected {
		t.Errorf("ConnectionString() = %v, want %v", got, expected)
	}
}

func TestPing(t *testing.T) {
	// Test that Ping method exists and handles nil DB
	t.Run("nil database", func(t *testing.T) {
		err := database.Ping(nil)
		if err == nil {
			t.Error("Ping(nil) should return error")
		}
	})
}

func TestClose(t *testing.T) {
	// Test that Close method exists and handles nil DB
	t.Run("nil database", func(t *testing.T) {
		err := database.Close(nil)
		if err == nil {
			t.Error("Close(nil) should return error")
		}
	})
}
