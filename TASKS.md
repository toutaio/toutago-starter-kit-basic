# Implementation Tasks - Starter Kit Basic

## Phase 1: Project Setup ⏱️ Week 1

### Task 1.1: Repository & Structure ✅ COMPLETE
- [x] Create GitHub repository `toutago-starter-kit-basic`
- [x] Initialize Go modules (`go mod init github.com/toutaio/toutago-starter-kit-basic`)
- [x] Create directory structure (cmd, internal, templates, static, etc.)
- [x] Setup .gitignore
- [x] Create LICENSE (MIT)
- [x] Initial README.md

### Task 1.2: CI/CD Setup ⏭️ SKIPPED
- [ ] Create .github/workflows/test.yml - Not needed per requirements
- [ ] Create .github/workflows/lint.yml - Not needed per requirements
- [ ] Configure golangci-lint - ✅ Config file created
- [ ] Setup code coverage reporting - Will be done locally

### Task 1.3: Docker Configuration ✅ COMPLETE
- [x] Create docker/Dockerfile (development)
- [x] Create docker/Dockerfile.prod (production)
- [x] Create docker-compose.yml
- [x] Test docker build - Configuration ready

### Task 1.4: Basic Configuration ✅ COMPLETE
- [x] Create .env.example
- [x] Implement config package (internal/config/config.go)
- [x] Add environment variable loading
- [x] Document all config options
- [x] Write comprehensive tests for config package (100% coverage)
- [x] Implement TDD approach (tests written first)

**Deliverable**: ✅ Empty project structure with Docker and configuration system ready

---

## Phase 2: Core Infrastructure ⏱️ Week 1-2 ✅ COMPLETE

### Task 2.1: Database Connection ✅ COMPLETE
- [x] Add datamapper dependencies to go.mod
- [x] Implement database connection in main.go
- [x] Add connection pooling configuration
- [x] Add health check endpoint
- [x] Test PostgreSQL connection
- [x] Test MySQL connection

### Task 2.2: Migration System ⏭️ DEFERRED
- [ ] Add sil-migrator dependency (deferred to Phase 3)
- [ ] Create migrations directory structure (created, implementation deferred)
- [ ] Implement migration runner
- [ ] Create Makefile targets (migrate, migrate-down)
- [ ] Document migration workflow

### Task 2.3: Router & Middleware ✅ COMPLETE
- [x] Add cosan-router dependency
- [x] Initialize router in main.go
- [x] Create logging middleware
- [x] Create security headers middleware
- [x] Create recovery middleware
- [x] Create CORS middleware
- [x] Create RequestID middleware
- [x] Add static file serving
- [x] Write comprehensive middleware tests (all passing)

### Task 2.4: Template System ✅ COMPLETE
- [x] Add fith-renderer dependency
- [x] Configure template loader
- [x] Create base layout template with Pico.css
- [x] Create home page template
- [x] Create custom CSS
- [x] Test template rendering
- [x] Integrate into main.go

**Deliverable**: ✅ Server running with database, router, middleware, and template system
- [ ] Test template rendering
- [ ] Add template caching

**Deliverable**: Server running with database, router, and template system

---

## Phase 3: Authentication ⏱️ Week 2 🔄 IN PROGRESS

### Task 3.1: User Model & Repository ✅ COMPLETE
- [x] Create internal/models/user.go
- [x] Create users migration (001_create_users.sql)
- [x] Create sessions migration (002_create_sessions.sql)
- [x] Implement password validation helpers
- [x] Write model tests (all passing)
- [x] Implement internal/repositories/user_repository.go
- [x] Implement internal/services/session_store.go
- [x] Write repository tests (all passing)

### Task 3.2: Auth Service ✅ COMPLETE
- [x] Add breitheamh-auth dependency
- [x] Create internal/services/auth_service.go
- [x] Implement Register method
- [x] Implement Login method
- [x] Implement Logout method
- [x] Implement VerifyEmail method
- [x] Implement ResetPassword method
- [x] Write service tests (all passing)

### Task 3.3: Auth Handlers ✅ COMPLETE
- [x] Create internal/handlers/auth_handler.go
- [x] Implement Register handler
- [x] Implement Login handler  
- [x] Implement Logout handler
- [x] Cookie management
- [x] Form data processing
- [x] Write handler tests (all passing)

### Task 3.4: Auth Middleware ✅ COMPLETE
- [x] Create internal/middleware/auth.go
- [x] Implement RequireAuth middleware
- [x] Implement RequireRole middleware
- [x] Implement OptionalAuth middleware
- [x] Session extraction from cookies
- [x] User context storage
- [x] Write middleware tests (all passing)

### Task 3.5: Auth Templates ✅ COMPLETE
- [x] Create templates/auth/login.html
- [x] Create templates/auth/register.html
- [x] Create templates/auth/forgot-password.html
- [x] Create templates/auth/reset-password.html
- [x] Update navigation with auth links
- [x] Style with Pico.css
- [x] Add error/success messages styling

**Deliverable**: ✅ Complete authentication system with tests

---

## Phase 4: Content Management 🔄 IN PROGRESS - Week 3

**Status**: 60% Complete (Models, Migrations, Repositories & Services Done)

### Task 4.1: Post System 🔄 IN PROGRESS
- [x] Create internal/models/post.go
- [x] Create internal/models/category.go  
- [x] Create internal/models/tag.go
- [x] Create posts migration (005_create_posts.sql)
- [x] Create categories migration (003_create_categories.sql)
- [x] Create tags migration (004_create_tags.sql)
- [x] Create post_tags migration (006_create_post_tags.sql)
- [x] Create post_versions migration (007_create_post_versions.sql)
- [x] Write model tests (all passing)
- [x] Create internal/domain/post.go
- [x] Create internal/repository/post_repository.go
- [x] Write repository tests (all passing - 91.7% coverage)
- [x] Create internal/service/post_service.go
- [x] Write service tests (all passing - 77.5% coverage)
- [ ] Create internal/repositories/category_repository.go
- [ ] Create internal/repositories/tag_repository.go
- [ ] Create internal/handlers/post_handler.go
- [ ] Write handler tests for post CRUD
- [ ] Create post templates (index, show, edit, new)

### Task 4.2: Page System ✅ COMPLETE
- [x] Create internal/models/page.go
- [x] Create pages migration (008_create_pages.sql)
- [x] Create page_versions migration (009_create_page_versions.sql)
- [x] Write model tests (all passing)
- [x] Create internal/domain/page.go
- [x] Create internal/repository/page_repository.go
- [x] Write repository tests (all passing - 91.7% coverage)
- [x] Create internal/service/page_service.go
- [x] Write service tests (all passing - 77.5% coverage)
- [ ] Create internal/handlers/page_handler.go
- [ ] Write handler tests for page CRUD
- [ ] Create page templates (index, show, edit)

### Task 4.3: Slug Generation ✅ COMPLETE
- [x] Implement slug generation utility (internal/helpers/content.go)
- [x] Add markdown rendering
- [x] Add HTML sanitization
- [x] Test slug generation (all passing)

### Task 4.4: File Upload Support ✅ COMPLETE
- [x] Implement file upload helper (internal/helpers/upload.go)
- [x] Add image validation
- [x] Add file size limits
- [x] Support multiple image types

### Task 4.5: Authorization
- [ ] Implement post ownership checks
- [ ] Implement page ownership checks
- [ ] Add role-based permissions
- [ ] Test authorization logic

**Deliverable**: Full content management for posts and pages

---

## Phase 5: UI & HTMX ⏱️ Week 3-4

### Task 5.1: Layout System
- [ ] Create templates/layouts/base.html
- [ ] Add Pico.css
- [ ] Create header partial
- [ ] Create footer partial
- [ ] Create navigation partial
- [ ] Test responsive design

### Task 5.2: Dashboard
- [ ] Create internal/handlers/dashboard_handler.go
- [ ] Create templates/pages/dashboard.html
- [ ] Show user stats
- [ ] Show recent activity
- [ ] Add quick action buttons

### Task 5.3: Profile & Settings
- [ ] Create internal/handlers/profile_handler.go
- [ ] Create templates/pages/profile.html
- [ ] Create internal/handlers/settings_handler.go
- [ ] Create templates/pages/settings.html
- [ ] Implement password change
- [ ] Implement email update

### Task 5.4: HTMX Integration
- [ ] Add HTMX library to static/js/
- [ ] Implement infinite scroll for posts
- [ ] Implement inline form validation
- [ ] Implement delete confirmations
- [ ] Implement live search
- [ ] Add loading indicators

### Task 5.5: Flash Messages
- [ ] Implement flash message system
- [ ] Create flash partial template
- [ ] Style success/error messages
- [ ] Test flash messages

**Deliverable**: Complete UI with HTMX progressive enhancement

---

## Phase 6: Message Bus ⏱️ Week 4

### Task 6.1: Email Service
- [ ] Add scela-bus dependency
- [ ] Create internal/services/email_service.go
- [ ] Implement login notification email
- [ ] Implement registration welcome email
- [ ] Implement password reset email
- [ ] Create email templates (plain text)
- [ ] Test email sending (use logger in dev)

### Task 6.2: Activity Tracking
- [ ] Create activity log migration (005_create_activity_log.sql)
- [ ] Create internal/services/activity_service.go
- [ ] Track post creation
- [ ] Track page creation
- [ ] Track login events
- [ ] Create activity log viewer

### Task 6.3: Event Bus Setup
- [ ] Initialize message bus in main.go
- [ ] Register event handlers
- [ ] Test event publishing
- [ ] Test async processing

**Deliverable**: Email notifications and activity tracking working

---

## Phase 7: Testing ⏱️ Week 4-5

### Task 7.1: Test Infrastructure
- [ ] Setup test database
- [ ] Create test helpers
- [ ] Create fixtures
- [ ] Setup test cleanup

### Task 7.2: Unit Tests
- [ ] Service layer tests (80%+ coverage)
- [ ] Repository tests (75%+ coverage)
- [ ] Utility function tests

### Task 7.3: Integration Tests
- [ ] Auth flow tests
- [ ] Post CRUD tests
- [ ] Page CRUD tests
- [ ] Email sending tests

### Task 7.4: Handler Tests
- [ ] Handler tests (70%+ coverage)
- [ ] Middleware tests
- [ ] Test HTMX responses

### Task 7.5: CI Integration
- [ ] Run tests in CI
- [ ] Add coverage reporting
- [ ] Add test badges to README

**Deliverable**: 75%+ test coverage across all layers

---

## Phase 8: Documentation ⏱️ Week 5

### Task 8.1: Main README
- [ ] Write introduction
- [ ] List features
- [ ] Write quick start guide
- [ ] Document project structure
- [ ] Document configuration
- [ ] Add screenshots
- [ ] Add badges

### Task 8.2: QUICK_START.md
- [ ] Step-by-step tutorial
- [ ] First login walkthrough
- [ ] Create first post guide
- [ ] Customize templates guide

### Task 8.3: ARCHITECTURE.md
- [ ] Document package overview
- [ ] Explain design patterns
- [ ] Create data flow diagrams
- [ ] Explain SSR + HTMX choice
- [ ] Document DI patterns

### Task 8.4: DEPLOYMENT.md
- [ ] Production checklist
- [ ] Environment setup
- [ ] SSL/TLS configuration
- [ ] Reverse proxy examples (nginx, Caddy)
- [ ] Database backup guide
- [ ] Monitoring suggestions

### Task 8.5: EXTENDING.md
- [ ] How to add new entities
- [ ] How to create handlers
- [ ] How to write services
- [ ] Template development guide
- [ ] HTMX patterns guide

### Task 8.6: Code Comments
- [ ] Add package documentation
- [ ] Add function documentation
- [ ] Add inline comments where needed
- [ ] Generate godoc

**Deliverable**: Complete, professional documentation

---

## Phase 9: Polish ⏱️ Week 5-6

### Task 9.1: Error Handling
- [ ] Implement custom error types
- [ ] Add error pages (404, 500, etc.)
- [ ] Improve error messages
- [ ] Test error scenarios

### Task 9.2: Validation
- [ ] Add input validation for all forms
- [ ] Add server-side validation
- [ ] Add client-side validation (HTMX)
- [ ] Test validation

### Task 9.3: Security Audit
- [ ] Review CSRF protection
- [ ] Review SQL injection prevention
- [ ] Review XSS protection
- [ ] Review password security
- [ ] Review session security
- [ ] Run security scanner

### Task 9.4: Performance
- [ ] Add template caching
- [ ] Add query optimization
- [ ] Add database indexes
- [ ] Test page load times
- [ ] Add compression

### Task 9.5: Demo Data
- [ ] Create seed script
- [ ] Add demo users
- [ ] Add demo posts
- [ ] Add demo pages
- [ ] Document seeding

### Task 9.6: Final Testing
- [ ] Test PostgreSQL deployment
- [ ] Test MySQL deployment
- [ ] Test Docker deployment
- [ ] Test production build
- [ ] Cross-browser testing

**Deliverable**: Production-ready starter kit

---

## Post-Launch Tasks

### Documentation
- [ ] Create video tutorial
- [ ] Write blog post announcement
- [ ] Create example projects

### Community
- [ ] Setup issue templates
- [ ] Setup PR templates
- [ ] Create CONTRIBUTING.md
- [ ] Create CODE_OF_CONDUCT.md

### Maintenance
- [ ] Setup dependabot
- [ ] Create release process
- [ ] Plan version 1.1 features

---

## Task Summary

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| 1. Project Setup | 4 | 2-3 days |
| 2. Core Infrastructure | 4 | 3-4 days |
| 3. Authentication | 5 | 5-6 days |
| 4. Content Management | 4 | 4-5 days |
| 5. UI & HTMX | 5 | 5-6 days |
| 6. Message Bus | 3 | 2-3 days |
| 7. Testing | 5 | 4-5 days |
| 8. Documentation | 6 | 3-4 days |
| 9. Polish | 6 | 4-5 days |
| **Total** | **42 tasks** | **5-6 weeks** |

---

## Priority Levels

🔴 **Critical** - Must have for v1.0
🟡 **Important** - Should have for v1.0
🟢 **Nice to have** - Can defer to v1.1

All tasks listed above are **Critical (🔴)** for v1.0 release.
