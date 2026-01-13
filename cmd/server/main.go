// Package main is the application entry point.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/toutaio/toutago-starter-kit-basic/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log startup information
	log.Printf("Starting Starter Kit Basic")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Server Port: %s", cfg.Server.Port)
	log.Printf("Database Driver: %s", cfg.Database.Driver)
	log.Printf("Database Host: %s", cfg.Database.Host)

	// Check for migration command
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Migration command not yet implemented")
		os.Exit(0)
	}

	// TODO: Initialize database connection
	// TODO: Initialize router
	// TODO: Initialize services
	// TODO: Start HTTP server

	fmt.Println("Server initialization complete (Phase 1)")
	fmt.Println("Phase 2 will add database connectivity and routing")
}
