// Package config provides application configuration management.
package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Session  SessionConfig
	Email    EmailConfig
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Environment string
	Port        string
	LogLevel    string
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// SessionConfig holds session management configuration.
type SessionConfig struct {
	Secret string
}

// EmailConfig holds email service configuration.
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("PORT", "8080"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
		},
		Session: SessionConfig{
			Secret: getEnv("SESSION_SECRET", ""),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "localhost"),
			SMTPPort:     getEnv("SMTP_PORT", "1025"),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromAddress:  getEnv("SMTP_FROM", "noreply@example.com"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required configuration is present.
func (c *Config) validate() error {
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// ConnectionString returns the database connection string.
func (d *DatabaseConfig) ConnectionString() string {
	switch d.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			d.User, d.Password, d.Host, d.Port, d.Name)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			d.User, d.Password, d.Host, d.Port, d.Name)
	default:
		return ""
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
