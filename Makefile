.PHONY: help dev build test lint clean migrate migrate-down seed docker-dev docker-prod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development server with hot reload (requires Air)
	@command -v air > /dev/null || (echo "Air not installed. Install with: go install github.com/air-verse/air@latest" && exit 1)
	air

build: ## Build production binary
	@echo "Building..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server
	@echo "Binary created at bin/server"

test: ## Run all tests
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run linter
	golangci-lint run ./...

clean: ## Clean build artifacts
	rm -rf bin/ tmp/ coverage.out coverage.html

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/server/main.go migrate up

migrate-down: ## Rollback last migration
	@echo "Rolling back migration..."
	@go run cmd/server/main.go migrate down

seed: ## Seed database with demo data
	@echo "Seeding database..."
	@go run scripts/seed.go

docker-dev: ## Start with Docker Compose
	docker compose up --build

docker-prod: ## Build production Docker image
	docker build -f docker/Dockerfile.prod -t starter-kit:latest .

tidy: ## Tidy go modules
	go mod tidy
