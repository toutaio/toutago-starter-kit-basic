# Architecture Overview

**Status**: Coming in Phase 2+

This document will explain:
- System architecture and design patterns
- Package organization
- Data flow
- Why SSR + HTMX?
- Dependency injection patterns
- Testing strategy

## Current Architecture (Phase 1)

### Package Structure
```
toutago-starter-kit-basic/
├── cmd/server/           # Application entry point
├── internal/
│   └── config/          # Configuration management ✅
└── tests/               # Test files
```

### Configuration System

The configuration package provides centralized management of application settings:

- Environment-based configuration
- Validation of required settings
- Database connection string generation
- Development vs Production modes

Check back for complete architecture documentation after Phase 5.
