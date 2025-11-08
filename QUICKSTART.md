# Quick Start Guide - Aiki Backend

## Prerequisites
- Docker & Docker Compose installed
- Go 1.24+ installed
- Make (optional)

## Setup (5 minutes)

### 1. Start Database Services
```bash
docker-compose up -d
```

This starts PostgreSQL and Redis in the background.

### 2. Install Tools
```bash
make install-tools
```

Or manually:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Run Migrations
```bash
make migrate-up
```

Or manually:
```bash
migrate -path migrations -database "postgresql://aiki:aiki_password@localhost:5432/aiki_db?sslmode=disable" up
```

### 5. Run the Server
```bash
make run
```

Or:
```bash
go run cmd/api/main.go
```

The API will be running at `http://localhost:8080`

## Test the API

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "Password123!"
  }'
```

Response:
```json
{
  "success": true,
  "message": "user registered successfully",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "Password123!"
  }'
```

### Get Current User (Protected)
```bash
# Save your access token from registration/login
TOKEN="your-access-token-here"

curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Update Current User
```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "phone_number": "+1234567890"
  }'
```

## Run Tests

### All tests
```bash
make test
```

### Unit tests only (skip integration tests)
```bash
go test -short ./...
```

### With coverage
```bash
make test-verbose
open coverage.html
```

## Common Commands

```bash
# Start services
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Build application
make build

# Clean artifacts
make clean
```

## Troubleshooting

### Port already in use
```bash
# Check what's using port 8080
lsof -i :8080

# Or change port in .env
PORT=3000
```

### Database connection failed
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Restart containers
docker-compose down && docker-compose up -d

# Check logs
docker-compose logs postgres
```

### Migration errors
```bash
# Check current migration version
migrate -path migrations -database "postgresql://aiki:aiki_password@localhost:5432/aiki_db?sslmode=disable" version

# Force to specific version
make migrate-force version=1
```

## Next Steps

1. ✅ User authentication is complete
2. Add Lock-In session tracking
3. Add job application tracker
4. Implement AI assistant integration
5. Add notification system

See `README.md` for full documentation.
