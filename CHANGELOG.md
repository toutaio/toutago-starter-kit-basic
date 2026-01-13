# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure and directory layout
- Go module initialization (github.com/toutaio/toutago-starter-kit-basic)
- Configuration management system with environment variable support
- PostgreSQL and MySQL connection string generation
- Development and production Docker configurations
- Docker Compose setup for local development
- Separate docker-compose.mysql.yml for MariaDB/MySQL
- Air configuration for hot reload during development
- Comprehensive test suite for configuration package
- Makefile with common development tasks
- golangci-lint configuration for code quality
- README with project overview and quick start guide
- MIT License
- .gitignore for Go projects
- .env.example with all configuration options
- Placeholder documentation structure
- HTTP middleware (Logger, Recovery, SecurityHeaders, CORS, RequestID)
- Integration with toutago-cosan-router

### Fixed
- Docker Compose healthcheck for PostgreSQL
- Git VCS error in Docker build by adding -buildvcs=false flag
- Database name in healthcheck command

### Testing
- Configuration package tests with 100% coverage
- Environment variable validation tests
- Database connection string generation tests
- Middleware tests with comprehensive coverage

[Unreleased]: https://github.com/toutaio/toutago-starter-kit-basic/commits/main
