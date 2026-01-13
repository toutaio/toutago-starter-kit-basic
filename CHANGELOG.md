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
- Database connection management with connection pooling
- PostgreSQL and MySQL driver support
- Health check endpoint with database status
- Template rendering with toutago-fith-renderer
- Home page with Pico.css styling
- Custom CSS with flash message support
- Static file serving
- Complete Phase 2 infrastructure
- Sil migrator integration for database migrations
- Go-based migrations with Up/Down methods
- Migration runner command (cmd/migrate) with support for migrate, rollback, status, reset, and fresh
- Database migrations for users, posts, and pages tables
- Makefile targets for migration management (migrate, migrate-down, migrate-status, migrate-fresh, migrate-reset)
- Database-agnostic migrations (PostgreSQL & MySQL)
- Users table migration with authentication fields
- Posts table migration with author relationship
- Pages table migration with author relationship
- User and Session models with validation
- Database migrations for users and sessions (PostgreSQL & MySQL)
- Password validation with complexity requirements
- Password hashing with bcrypt
- User roles (admin, editor, user)
- breitheamh-auth integration
- scela-bus message bus integration
- UserRepository with in-memory implementation
- SessionStore for session management  
- AuthService with register, login, logout, password reset
- Email verification support
- Auth handlers (register, login, logout)
- Auth middleware (RequireAuth, RequireRole, OptionalAuth)
- Auth templates (login, register, forgot/reset password)
- Content management models (Post, Page, Category, Tag)
- Post domain model with status management (draft, published, archived)
- Page domain model with status management (draft, published, archived)
- PostRepository with full CRUD operations
- PageRepository with full CRUD operations
- Repository tests with 91.7% coverage
- PostService with business logic and validation
- PageService with business logic and validation
- Service tests with 77.5% coverage
- Post publishing/unpublishing functionality
- Page publishing/unpublishing functionality
- Slug-based content retrieval
- Author-based post filtering
- Status-based content filtering
- Enhanced navigation with user dropdown
- Session-based authentication with cookies
- Pico.css styling for auth pages
- Comprehensive auth tests (30 new tests)
- Post and Page models with soft delete (trash) support
- Category and Tag models for post organization
- PostVersion and PageVersion models for content history
- Slug generation utility with comprehensive tests
- URL-friendly slug conversion from any string
- Fixed handler tests to accept HTTP 303 redirects
- Content helper functions (slug generation, markdown rendering, HTML sanitization)
- Dashboard handler with user statistics
- Dashboard template showing posts, pages, and activity
- Profile handler for viewing and updating user information
- Profile template with email editing capability
- Settings handler with password change functionality
- Settings template for account management
- UpdatePassword method in AuthService
- ListByAuthor method in PageRepository
- File upload helper with image validation
- Database migrations for posts, pages, categories, tags (PostgreSQL & MySQL)
- Support for version history tracking
- Markdown to HTML conversion with sanitization
- URL-friendly slug generation with unicode support
- Image upload with type and size validation
- Post handler with full CRUD operations (Index, Show, Create, Update, Delete)
- Page handler with full CRUD operations (Index, Show, Create, Update, Delete)
- Post templates (index, show, new, edit) with Pico.css styling
- Page templates (index, show, new, edit) with Pico.css styling
- Publish/unpublish actions for posts and pages
- Authorization checks for post/page editing and deletion
- Integration with helpers for slug generation
- AuthorID field to Page domain model for consistency
- Authorization middleware with role-based permissions (admin, editor, user)
- Permission helpers (CanEdit, CanDelete, CanPublish, CanManageUsers)
- RequireOwnership middleware for resource protection
- Flash message system with middleware
- Flash message template partial with auto-dismiss
- Flash message styling (success, error, warning, info)
- Slug uniqueness validation in PostService and PageService
- Authorization tests with 100% coverage

### Fixed
- Docker Compose healthcheck for PostgreSQL
- Git VCS error in Docker build by adding -buildvcs=false flag
- Database name in healthcheck command
- Migration foreign key constraints by standardizing ID types to INTEGER/SERIAL

### Testing
- Configuration package tests with 100% coverage
- Environment variable validation tests
- Database connection string generation tests
- Middleware tests with comprehensive coverage
- Database connection tests
- Health handler tests
- Home handler tests with template rendering

[Unreleased]: https://github.com/toutaio/toutago-starter-kit-basic/commits/main
