# Starter Kit Basic

A production-ready Go web application starter kit demonstrating the Toutā framework with Server-Side Rendering (SSR) and HTMX.

## Features

- 🔐 Complete authentication system (register, login, password reset, email verification)
- 👥 Role-based authorization (Admin, Editor, Author, Reader)
- 📝 Posts and Pages content management
- 🎨 Server-Side Rendering with Fíth templates
- ⚡ Progressive enhancement with HTMX
- 🎯 Clean UI with Pico.css
- 🗄️ PostgreSQL and MySQL support (switchable)
- 📧 Email notifications via message bus
- 🐳 Docker support for development and production
- ✅ Comprehensive test coverage

## Technology Stack

### Toutā Framework Packages
- **toutago-cosan-router** - HTTP routing
- **toutago-fith-renderer** - Server-side template rendering
- **toutago-nasc-dependency-injector** - Dependency injection
- **toutago-breitheamh-auth** - Authentication & authorization
- **toutago-datamapper** - Database abstraction
- **toutago-scela-bus** - Event/message bus
- **toutago-sil-migrator** - Database migrations

### External Dependencies
- HTMX 1.9+ for progressive enhancement
- Pico.css 2.x for styling
- PostgreSQL or MySQL/MariaDB

## Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 16+ or MySQL 8+/MariaDB 11+
- Docker and Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/toutaio/toutago-starter-kit-basic.git
cd toutago-starter-kit-basic
```

2. Configure environment:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Install dependencies:
```bash
go mod download
```

4. Run migrations:
```bash
make migrate
```

5. Start the development server:
```bash
make dev
```

Visit http://localhost:8080

### Using Docker

```bash
# Development with hot reload
make docker-dev

# Production build
make docker-prod
```

## Project Structure

```
toutago-starter-kit-basic/
├── cmd/server/           # Application entry point
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   ├── repositories/     # Data access layer
│   ├── models/           # Domain models
│   ├── middleware/       # HTTP middleware
│   └── config/           # Configuration
├── migrations/           # Database migrations
├── templates/            # Fíth templates
├── static/              # Static assets
├── tests/               # Test files
└── docker/              # Docker configuration
```

## Development

### Make Commands

```bash
make dev          # Start with hot reload
make build        # Build production binary
make test         # Run tests
make lint         # Run linter
make migrate      # Run database migrations
make migrate-down # Rollback migrations
make seed         # Seed demo data
make clean        # Clean build artifacts
```

### Running Tests

```bash
# All tests
make test

# With coverage
go test -cover ./...

# Specific package
go test ./internal/services/...
```

## Configuration

Environment variables (see `.env.example`):

- `DB_DRIVER` - Database driver: postgres or mysql
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `PORT` - Server port (default: 8080)
- `APP_ENV` - Environment: development or production

## Architecture

This starter kit follows clean architecture principles:

- **Handlers** - HTTP request/response handling
- **Services** - Business logic and orchestration
- **Repositories** - Data persistence abstraction
- **Models** - Domain entities

All components use dependency injection for testability and maintainability.

## Documentation

- [Quick Start Guide](docs/QUICK_START.md)
- [Architecture Overview](docs/ARCHITECTURE.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [Extending the Starter Kit](docs/EXTENDING.md)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- GitHub Issues: https://github.com/toutaio/toutago-starter-kit-basic/issues
- Documentation: https://github.com/toutaio/toutago-starter-kit-basic/tree/main/docs

## Acknowledgments

Built with the [Toutā Framework](https://github.com/toutaio) - A modern Go web framework.
