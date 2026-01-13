package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/toutaio/toutago-starter-kit-basic/internal/config"
	"github.com/toutaio/toutago-starter-kit-basic/internal/database"
)

func main() {
	action := flag.String("action", "up", "Migration action: up, down, status")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "postgres"
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	switch *action {
	case "up":
		fmt.Println("Running migrations...")
		if err := database.RunMigrations(db, dbType); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		fmt.Println("Rolling back last migration...")
		if err := database.RollbackMigrations(db, dbType); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Rollback completed successfully")

	case "status":
		statuses, err := database.GetMigrationStatus(db, dbType)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		fmt.Println("Migration Status:")
		for _, status := range statuses {
			appliedStatus := "❌ Pending"
			if status.Applied {
				appliedStatus = "✅ Applied"
			}
			fmt.Printf("  %s - %s\n", status.Version, appliedStatus)
		}

	default:
		log.Fatalf("Unknown action: %s (use: up, down, status)", *action)
	}
}
