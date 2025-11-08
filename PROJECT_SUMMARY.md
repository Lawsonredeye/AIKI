# Aiki Backend - Project Summary

## What We Built

A production-ready, test-driven backend for the Aiki MVP focusing on **user account creation and authentication**. The architecture is clean, scalable, and follows Go best practices.

## ✅ Completed Features

### 1. **Authentication System**
- User registration with email/password
- Login with JWT access tokens
- Refresh token mechanism (7-day expiry)
- Logout functionality
- Password strength validation (uppercase, lowercase, number, special char)
- Bcrypt password hashing (cost 12)

### 2. **User Profile Management**
- Get current user profile
- Update user profile (name, phone number)
- Soft delete capability (is_active flag)

### 3. **Security**
- JWT-based authentication
- Bearer token validation middleware
- Password hashing with bcrypt
- CORS middleware
- Input validation
- SQL injection prevention (parameterized queries)

### 4. **Testing (TDD Approach)**
- Unit tests for all layers
- Integration tests for repository layer
- Mock-based testing for services and handlers
- Middleware tests
- **Test Coverage**: All critical paths tested
- **Test Results**: All tests passing ✅

## Architecture

### Clean Architecture (3 Layers)

```
┌─────────────────────────────────────────┐
│           Handler Layer                  │  ← HTTP Handlers
│  (auth_handler, user_handler)           │     Request/Response
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│          Service Layer                   │  ← Business Logic
│  (auth_service, user_service)           │     Validation, Processing
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│        Repository Layer                  │  ← Data Access
│  (user_repository)                       │     SQL Queries
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│         PostgreSQL + Redis               │  ← Data Storage
└─────────────────────────────────────────┘
```

### Directory Structure

```
aiki/
├── cmd/api/main.go                    # Application entry point
├── internal/
│   ├── config/                        # Configuration management
│   ├── database/                      # DB connections (Postgres, Redis)
│   ├── domain/                        # Domain models & errors
│   │   ├── user.go
│   │   ├── jwt.go
│   │   └── errors.go
│   ├── handler/                       # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── auth_handler_test.go
│   │   ├── user_handler.go
│   │   └── user_handler_test.go
│   ├── middleware/                    # HTTP middleware
│   │   ├── auth.go
│   │   ├── auth_test.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recover.go
│   ├── pkg/                           # Shared utilities
│   │   ├── jwt/                       # JWT generation/validation
│   │   ├── password/                  # Password hashing/validation
│   │   ├── response/                  # Standardized API responses
│   │   └── validator/                 # Request validation
│   ├── repository/                    # Data access layer
│   │   ├── user_repository.go
│   │   └── user_repository_test.go
│   ├── service/                       # Business logic
│   │   ├── auth_service.go
│   │   ├── auth_service_test.go
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   └── router/                        # Route definitions
│       └── router.go
├── migrations/                        # Database migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_refresh_tokens_table.up.sql
│   └── 000002_create_refresh_tokens_table.down.sql
├── docker-compose.yml                 # Local development setup
├── Makefile                           # Development commands
└── README.md                          # Full documentation
```

## Tech Stack Choices & Rationale

| Component | Choice | Why? |
|-----------|--------|------|
| **Language** | Go 1.24 | Concurrency, performance, type safety |
| **Web Framework** | Echo v4 | Lightweight, fast, great middleware support |
| **Database** | PostgreSQL 16 | ACID compliance, relational data, strong indexes |
| **Cache** | Redis 7 | Fast token storage, future rate limiting |
| **SQL Tool** | SQLC | Type-safe SQL, compile-time checks |
| **Migrations** | golang-migrate | Standard tool, SQL-based migrations |
| **Auth** | JWT | Stateless, scalable, industry standard |
| **Testing** | testify | Rich assertions, mocking support |

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## API Endpoints

### Public Endpoints
```
POST   /api/v1/auth/register      # Register new user
POST   /api/v1/auth/login         # Login
POST   /api/v1/auth/refresh       # Refresh access token
POST   /api/v1/auth/logout        # Logout
GET    /api/v1/health             # Health check
```

### Protected Endpoints (Require JWT)
```
GET    /api/v1/users/me           # Get current user
PUT    /api/v1/users/me           # Update current user
```

## Test Coverage

```
✓ JWT Generation & Validation     (4/4 tests passing)
✓ Password Hashing & Validation   (6/6 tests passing)
✓ Auth Service                    (7/7 tests passing)
✓ User Service                    (6/6 tests passing)
✓ Auth Handler                    (7/7 tests passing)
✓ User Handler                    (6/6 tests passing)
✓ Auth Middleware                 (5/5 tests passing)
```

**Total: 41 tests passing**

## Security Features

1. **Password Security**
   - Bcrypt hashing (cost 12)
   - Strength validation (8+ chars, upper, lower, number, special)
   - Never returned in API responses

2. **JWT Security**
   - HS256 signing
   - 15-minute access token expiry
   - 7-day refresh token expiry
   - Secure token validation

3. **SQL Injection Prevention**
   - Parameterized queries everywhere
   - SQLC compile-time query validation

4. **CORS Protection**
   - Configurable allowed origins
   - Proper headers

5. **Error Handling**
   - No sensitive data in error messages
   - Proper HTTP status codes
   - Centralized error mapping

## Performance Considerations

1. **Database**
   - Indexes on email, is_active
   - Connection pooling (5-25 connections)
   - Prepared statements via pgx

2. **Redis**
   - Future token caching
   - Rate limiting support

3. **Middleware Chain**
   - Logger → Recover → CORS → Auth
   - Minimal overhead

## Development Workflow

```bash
# 1. Start dependencies
make docker-up

# 2. Run migrations
make migrate-up

# 3. Run tests
make test

# 4. Start server
make run
```

## What's Next? (Future Features)

Based on the PRD, here's what to build next:

### 1. Lock-In Sessions
- POST /api/v1/sessions/lock-in
- PUT /api/v1/sessions/:id/complete
- GET /api/v1/sessions/streak
- Database table for sessions & streaks

### 2. Job Application Tracker
- CRUD endpoints for job applications
- Filter by status, platform
- Weekly stats endpoint

### 3. AI Assistant Integration
- OpenAI/Claude API integration
- CV analysis endpoint
- Cover letter generation
- Proposal improvement

### 4. User Profiles Extended
- Education, skills, links tables
- CV upload (S3 integration)
- Profile completion tracking

### 5. Notifications
- Email notifications (SendGrid)
- Push notifications (Firebase)
- Daily "Lock In" reminders

## Configuration

All configuration is environment-based:

```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=aiki
DB_PASSWORD=aiki_password
DB_NAME=aiki_db

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

## Production Deployment Checklist

- [ ] Change JWT_SECRET to strong random value
- [ ] Configure CORS with specific origins
- [ ] Set up SSL/TLS certificates
- [ ] Enable database SSL mode
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure logging (structured JSON logs)
- [ ] Set up CI/CD pipeline
- [ ] Add rate limiting middleware
- [ ] Set up backup strategy
- [ ] Configure alerting
- [ ] Add health check monitoring
- [ ] Set up log aggregation (ELK/Loki)

## Key Design Decisions

### 1. PostgreSQL over MongoDB
**Why?** The data is highly relational (users → profiles → skills → jobs). PostgreSQL gives us:
- Foreign key constraints
- ACID transactions
- Better query performance for relations
- Strong typing

### 2. SQLC over ORM (GORM)
**Why?**
- Type safety at compile time
- Better performance (no reflection)
- Full SQL control
- Easier to optimize queries

### 3. JWT over Sessions
**Why?**
- Stateless (no session storage)
- Scalable across multiple servers
- Mobile-friendly
- Industry standard

### 4. Repository Pattern
**Why?**
- Testability (easy to mock)
- Separation of concerns
- Can swap DB implementations
- Clean dependencies

## Contributors

Built with TDD approach using:
- Go 1.24
- Echo v4
- PostgreSQL 16
- Redis 7
- Docker

## License

MIT
