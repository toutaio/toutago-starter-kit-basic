package config_test

import (
	"os"
	"testing"

	"github.com/toutaio/toutago-starter-kit-basic/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantErr  bool
		validate func(*testing.T, *config.Config)
	}{
		{
			name: "loads configuration with defaults",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "test_db",
				"DB_USER":     "test_user",
				"DB_PASSWORD": "test_pass",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *config.Config) {
				if cfg.Database.Host != "localhost" {
					t.Errorf("expected DB host localhost, got %s", cfg.Database.Host)
				}
				if cfg.Database.Port != "5432" {
					t.Errorf("expected DB port 5432, got %s", cfg.Database.Port)
				}
				if cfg.Server.Port != "8080" { // default
					t.Errorf("expected server port 8080, got %s", cfg.Server.Port)
				}
			},
		},
		{
			name: "loads all custom values",
			envVars: map[string]string{
				"APP_ENV":     "production",
				"PORT":        "3000",
				"DB_DRIVER":   "mysql",
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "3306",
				"DB_NAME":     "prod_db",
				"DB_USER":     "prod_user",
				"DB_PASSWORD": "prod_pass",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *config.Config) {
				if cfg.Server.Environment != "production" {
					t.Errorf("expected production env, got %s", cfg.Server.Environment)
				}
				if cfg.Server.Port != "3000" {
					t.Errorf("expected port 3000, got %s", cfg.Server.Port)
				}
				if cfg.Database.Driver != "mysql" {
					t.Errorf("expected mysql driver, got %s", cfg.Database.Driver)
				}
			},
		},
		{
			name: "requires database credentials",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_NAME": "test_db",
				// Missing DB_USER and DB_PASSWORD
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test env vars
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set env var %s: %v", k, err)
				}
			}

			cfg, err := config.Load()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestDatabaseConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.DatabaseConfig
		expected string
	}{
		{
			name: "PostgreSQL connection string",
			cfg: config.DatabaseConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     "5432",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_pass",
			},
			expected: "postgres://test_user:test_pass@localhost:5432/test_db?sslmode=disable",
		},
		{
			name: "MySQL connection string",
			cfg: config.DatabaseConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_pass",
			},
			expected: "test_user:test_pass@tcp(localhost:3306)/test_db?parseTime=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cfg.ConnectionString()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsDevelopment(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Environment: "development",
		},
	}

	if !cfg.IsDevelopment() {
		t.Error("expected IsDevelopment() to return true")
	}

	cfg.Server.Environment = "production"
	if cfg.IsDevelopment() {
		t.Error("expected IsDevelopment() to return false")
	}
}
