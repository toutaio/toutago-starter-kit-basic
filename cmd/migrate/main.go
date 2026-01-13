package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/toutaio/toutago-sil-migrator/pkg/sil"
	"github.com/toutaio/toutago-sil-migrator/pkg/sil/adapters"

	// Import migrations to register them
	_ "github.com/toutaio/toutago-starter-kit-basic/internal/migrations"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get command from args
	command := "migrate"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	// Build DATABASE_URL from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Build from individual variables
		dbDriver := os.Getenv("DB_DRIVER")
		if dbDriver == "" {
			dbDriver = "postgres"
		}
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")

		if dbHost == "" || dbName == "" || dbUser == "" {
			log.Fatal("Database configuration is incomplete. Set DATABASE_URL or DB_HOST, DB_NAME, DB_USER variables")
		}

		switch dbDriver {
		case "postgres":
			if dbPort == "" {
				dbPort = "5432"
			}
			databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
		case "mysql":
			if dbPort == "" {
				dbPort = "3306"
			}
			databaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
		default:
			log.Fatalf("Unsupported database driver: %s", dbDriver)
		}
	}

	// Configure database connection
	config := sil.DefaultConfig()
	config.DatabaseURL = databaseURL
	config.MigrationsDir = "./internal/migrations"
	config.Verbose = true

	// Determine database type
	var adapter sil.DatabaseAdapter
	var err error

	dbType := os.Getenv("DB_DRIVER")
	if dbType == "" {
		dbType = os.Getenv("DB_TYPE")
	}
	if dbType == "" {
		dbType = "postgres" // default
	}

	switch dbType {
	case "postgres":
		adapter, err = adapters.NewPostgresAdapter(config)
	case "mysql":
		adapter, err = adapters.NewMySQLAdapter(config)
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := adapter.Connect(ctx, config); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer adapter.Close()

	// Create migrator
	migrator, err := sil.NewMigrator(config, adapter)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	// Execute command
	switch command {
	case "migrate", "up":
		fmt.Println("Running migrations...")
		if err := migrator.Migrate(ctx); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✅ Migrations complete!")

	case "rollback", "down":
		fmt.Println("Rolling back last batch...")
		if err := migrator.Rollback(ctx); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("✅ Rollback complete!")

	case "status":
		fmt.Println("Migration status:")
		statuses, err := migrator.Status(ctx)
		if err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}

		for _, status := range statuses {
			appliedStatus := "❌ Pending"
			if status.Applied {
				appliedStatus = fmt.Sprintf("✅ Applied (batch %d)", status.Batch)
			}
			fmt.Printf("  %s - %s - %s\n", status.Version, status.Description, appliedStatus)
		}

	case "reset":
		fmt.Println("Resetting all migrations...")
		if err := migrator.Reset(ctx); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		fmt.Println("✅ Reset complete!")

	case "fresh":
		fmt.Println("Resetting and re-running all migrations...")
		if err := migrator.Reset(ctx); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		if err := migrator.Migrate(ctx); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✅ Fresh migration complete!")

	default:
		log.Fatalf("Unknown command: %s. Available: migrate, rollback, status, reset, fresh", command)
	}
}
