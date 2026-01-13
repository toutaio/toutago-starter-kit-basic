// Package main is the application entry point.
package main

import (
	"log"
	"net/http"
	"os"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-starter-kit-basic/internal/config"
	"github.com/toutaio/toutago-starter-kit-basic/internal/database"
	"github.com/toutaio/toutago-starter-kit-basic/internal/handlers"
	"github.com/toutaio/toutago-starter-kit-basic/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Starter Kit Basic")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Server Port: %s", cfg.Server.Port)

	// Check for migration command
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Migration command not yet implemented")
		os.Exit(0)
	}

	// Initialize database connection
	log.Printf("Connecting to %s database at %s:%s", cfg.Database.Driver, cfg.Database.Host, cfg.Database.Port)
	sqlDB, err := database.Connect(cfg.Database)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Continuing without database connection...")
	} else {
		log.Println("Database connected successfully")
		defer database.Close(sqlDB)
	}

	// Initialize template renderer
	rendererConfig := &fith.Config{
		TemplateDir: "templates",
	}
	renderer, err := fith.New(rendererConfig)
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}
	log.Println("Template renderer initialized")

	// Initialize router
	r := router.New()

	// Apply global middleware
	r.Use(router.MiddlewareFunc(middleware.RequestID))
	r.Use(router.MiddlewareFunc(middleware.Logger))
	r.Use(router.MiddlewareFunc(middleware.Recovery))
	r.Use(router.MiddlewareFunc(middleware.SecurityHeaders))

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(sqlDB)
	homeHandler := handlers.NewHomeHandler(renderer)

	// Register routes
	r.GET("/", homeHandler.Index)
	r.GET("/health", healthHandler.Check)

	// Serve static files (using GET for now since Static might not be available)
	r.GET("/static/*", func(ctx router.Context) error {
		http.FileServer(http.Dir("static")).ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	// Start HTTP server
	addr := ":" + cfg.Server.Port
	log.Printf("Server listening on %s", addr)
	log.Printf("Visit http://localhost%s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
