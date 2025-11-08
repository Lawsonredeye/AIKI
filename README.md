# Aiki - Job Search & Productivity MVP Backend

Backend API for Aiki, a productivity and job search management application built with Go, Echo, PostgreSQL, and Redis.

## Features

- ✅ User authentication (JWT-based)
- ✅ User registration and profile management
- ✅ Refresh token mechanism
- ✅ Password strength validation
- ✅ Clean architecture (Repository, Service, Handler layers)
- ✅ Comprehensive unit and integration tests
- ✅ Docker-ready (PostgreSQL + Redis)
- ✅ Database migrations with golang-migrate
- ✅ Type-safe SQL with SQLC

## Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Echo v4
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **SQL Generator**: SQLC
- **Authentication**: JWT
- **Testing**: testify
- **Containerization**: Docker & Docker Compose

## Project Structure

```
aiki/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── database/                   # DB connection setup
│   ├── domain/                     # Domain models & errors
│   ├── handler/                    # HTTP handlers (controllers)
│   ├── middleware/                 # HTTP middleware
│   ├── pkg/                        # Shared utilities
│   │   ├── jwt/                    # JWT helpers
│   │   ├── password/               # Password hashing
│   │   ├── response/               # API responses
│   │   └── validator/              # Request validation
│   ├── repository/                 # Data access layer
│   ├── router/                     # Route definitions
│   └── service/                    # Business logic
├── migrations/                     # Database migrations
├── docker-compose.yml
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make (optional, but recommended)

### Installation

1. **Clone the repository**
```bash
cd /home/lawson/Desktop/projects/AIKI
```

2. **Install development tools**
```bash
make install-tools
```

This installs:
- `sqlc` - SQL code generator
- `golang-migrate` - Database migration tool
- `mockgen` - Mock generator for tests

3. **Copy environment file**
```bash
cp .env.example .env
```

Edit `.env` with your configuration.

4. **Start dependencies (PostgreSQL + Redis)**
```bash
make docker-up
```

5. **Run database migrations**
```bash
make migrate-up
```

6. **Install Go dependencies**
```bash
go mod tidy
```

7. **Run the application**
```bash
make run
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /api/v1/health
```

### Authentication (Public)

#### Register
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone_number": "+1234567890",
  "password": "Password123!"
}
```

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "Password123!"
}
```

#### Refresh Token
```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

#### Logout
```bash
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

### Users (Protected - Requires JWT)

#### Get Current User
```bash
GET /api/v1/users/me
Authorization: Bearer <your-access-token>
```

#### Update Current User
```bash
PUT /api/v1/users/me
Authorization: Bearer <your-access-token>
Content-Type: application/json

{
  "first_name": "Jane",
  "phone_number": "+9876543210"
}
```

## Testing

### Run all tests
```bash
make test
```

### Run tests with coverage
```bash
make test-verbose
```

This generates `coverage.html` which you can open in a browser.

### Run only unit tests (skip integration tests)
```bash
go test -short ./...
```

## Database Migrations

### Create a new migration
```bash
make migrate-create name=create_users_table
```

### Run migrations
```bash
make migrate-up
```

### Rollback migrations
```bash
make migrate-down
```

## Development Workflow

1. **Make changes to SQL queries** in `internal/database/queries/*.sql`
2. **Regenerate SQLC code**
   ```bash
   make sqlc-generate
   ```
3. **Write tests first** (TDD approach)
4. **Implement features**
5. **Run tests**
   ```bash
   make test
   ```

## Docker Commands

```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs
```

## Production Considerations

- [ ] Change `JWT_SECRET` in production
- [ ] Configure CORS with specific allowed origins
- [ ] Set up proper logging (structured logging with Zap/Zerolog)
- [ ] Add rate limiting
- [ ] Set up monitoring (Prometheus + Grafana)
- [ ] Use a proper secrets manager
- [ ] Configure SSL/TLS
- [ ] Set up CI/CD pipeline
- [ ] Add OpenAPI/Swagger documentation

## Password Requirements

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

## License

MIT
