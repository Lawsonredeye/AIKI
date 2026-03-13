# Aiki Application Architecture

## Overview

Aiki is a Go-based web application designed to help users maintain focus through timed sessions, track job applications, and build streaks with gamification elements like badges. The application follows Clean Architecture principles to ensure separation of concerns, testability, and maintainability.

## Technology Stack

### Core Technologies
- **Language**: Go 1.25
- **Web Framework**: Echo v4
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **ORM/Query Builder**: SQLC (type-safe SQL generation)

### External Services
- **SerpAPI**: For job search and recommendations
- **LinkedIn OAuth**: For social login (optional)

### Development Tools
- **Containerization**: Docker & Docker Compose
- **API Documentation**: Swagger/OpenAPI
- **Testing**: Go testing with mocks (mockgen)
- **Migrations**: golang-migrate
- **Build Tool**: Makefile

### Key Libraries
- **Database**: `github.com/jackc/pgx/v5` (PostgreSQL driver)
- **Cache**: `github.com/redis/go-redis/v9`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **Validation**: `github.com/go-playground/validator/v10`
- **Password Hashing**: `golang.org/x/crypto/bcrypt`
- **Environment**: `github.com/joho/godotenv`
- **Swagger**: `github.com/swaggo/swag`

## Architecture Overview

The application follows Clean Architecture with clear separation of layers:

```
┌─────────────────┐
│   Handlers      │  HTTP Layer
├─────────────────┤
│   Services      │  Business Logic
├─────────────────┤
│  Repositories   │  Data Access
├─────────────────┤
│    Domain       │  Entities & Business Rules
└─────────────────┘
```

### Layer Responsibilities

#### Domain Layer (`internal/domain/`)
Contains business entities, value objects, and domain logic:
- **User**: User management, authentication
- **Job**: Job application tracking
- **FocusSession**: Pomodoro-style focus sessions
- **Streak**: User streaks and consistency tracking
- **Badge**: Gamification badges
- **Notification**: In-app notifications
- **SerpJob**: Job search results from external APIs

#### Repository Layer (`internal/repository/`)
Handles data persistence and retrieval:
- **UserRepository**: User CRUD operations
- **JobRepository**: Job tracking operations
- **HomeRepository**: Home screen data (sessions, streaks, badges)
- **NotificationRepository**: Notification management
- **SerpJobRepository**: Cached job search results

#### Service Layer (`internal/service/`)
Contains business logic and orchestrates operations:
- **AuthService**: Authentication, JWT token management
- **UserService**: User profile management
- **JobService**: Job application tracking
- **HomeService**: Focus sessions, streaks, badges
- **NotificationService**: Notification creation and delivery
- **SerpJobService**: Job search integration

#### Handler Layer (`internal/handler/`)
HTTP request/response handling:
- **AuthHandler**: Login, registration, token refresh
- **UserHandler**: Profile management, CV upload
- **JobHandler**: Job CRUD operations
- **HomeHandler**: Home screen, sessions, streaks, badges
- **NotificationHandler**: Notification retrieval and management
- **SerpJobHandler**: Job recommendations

## Core Services

### Authentication Service
- JWT-based authentication with access and refresh tokens
- Password hashing with bcrypt
- LinkedIn OAuth integration (optional)
- Password reset with OTP via email
- Session management with Redis

### Focus Session Management
- Start, pause, resume, and end focus sessions
- Real-time session tracking
- Session history and statistics
- Automatic streak calculation

### Job Tracking
- Manual job application tracking
- Status updates (applied, interviewed, etc.)
- Notes and metadata storage
- Job recommendations via SerpAPI

### Gamification System
- Streak tracking (daily focus consistency)
- Badge system with predefined achievements
- Progress statistics and analytics
- Notification triggers for milestones

### Notification System
- In-app notifications for achievements
- Scheduled daily reminders
- Streak warning notifications
- Email-based password reset (OTP)

## Database Schema

The application uses PostgreSQL with the following main tables:

### Core Tables
- `users`: User accounts and authentication
- `refresh_tokens`: JWT refresh token storage
- `user_profile`: Extended user information and CV storage

### Job Tracking
- `jobs`: Manual job application tracking

### Focus & Gamification
- `focus_sessions`: Focus session records
- `streaks`: User streak data
- `badge_definitions`: Available badges
- `user_badges`: Earned badges
- `daily_progress`: Daily statistics

### Notifications
- `notifications`: In-app notification storage

### External Integrations
- `serp_job_cache`: Cached job search results from SerpAPI

## External Integrations

### SerpAPI Integration
- Job search using Google Jobs API
- Caching with 24-hour TTL
- Location and experience level filtering
- Job saving to personal tracker

### LinkedIn OAuth (Optional)
- Social login integration
- Profile data enrichment

## Infrastructure

### Development Environment
- Docker Compose with PostgreSQL and Redis
- Hot reload with `air` or manual rebuild
- Environment-based configuration

### Production Deployment
- Docker containerization
- Environment variables for configuration
- Health checks and graceful shutdown

### Background Jobs
- Notification scheduler running daily reminders and warnings
- Automatic badge awarding on milestones
- Streak maintenance and calculation

## API Design

### RESTful Endpoints
- `/api/v1/auth/*`: Authentication endpoints
- `/api/v1/users/*`: User management
- `/api/v1/jobs/*`: Job tracking
- `/api/v1/home/*`: Home screen data
- `/api/v1/sessions/*`: Focus sessions
- `/api/v1/streaks/*`: Streak data
- `/api/v1/badges/*`: Badge system
- `/api/v1/progress/*`: Progress statistics
- `/api/v1/notifications/*`: Notification management

### Middleware
- **Auth**: JWT token validation
- **CORS**: Cross-origin resource sharing
- **Logger**: Request logging
- **Recover**: Panic recovery
- **Validator**: Request validation

## Security

### Authentication
- JWT access tokens (15-minute expiry)
- Refresh tokens (7-day expiry)
- Password hashing with bcrypt
- OTP for password reset

### Data Protection
- Input validation and sanitization
- SQL injection prevention via prepared statements
- XSS protection via Echo framework
- Secure headers and CORS configuration

## Testing

### Unit Tests
- Service layer testing with mocks
- Repository layer testing with test database
- Handler testing with Echo test utilities

### Integration Tests
- Full API testing with test database
- External API mocking for SerpAPI

## Development Workflow

### Code Generation
- SQLC for type-safe database queries
- Mock generation with mockgen
- Swagger documentation generation

### Database Migrations
- Version-controlled schema changes
- Up/down migration support
- Test database isolation

### Build Process
- Makefile for common tasks
- Docker-based builds
- Multi-stage Dockerfiles for optimization

## Hosting Costs

This section outlines the estimated monthly hosting costs for deploying the Aiki application across different user scales. Costs are based on current pricing (as of March 2026) from popular cloud providers. The application uses:

- **Render**: For Golang application hosting
- **Neon**: For PostgreSQL database (as configured in the project)
- **Redis Cloud**: For Redis caching and queuing
- **Firebase**: Assumed for additional services like authentication or real-time features (not currently implemented)
- **Queue Service**: Using Redis for background job queuing (included in Redis costs)

### Cost Estimation Methodology
- **User Tiers**: Based on concurrent users and expected load
- **Assumptions**: 
  - Average 10-20% of users active simultaneously
  - Database storage scales with users (100KB per user average)
  - Bandwidth and compute scale with user activity
- **Pricing Sources**: Official provider pricing pages
- **Currency**: USD per month

### Hosting Cost Comparison Table

| Service | Tier 1: 0-100 Users | Tier 2: 100-1,000 Users | Tier 3: 1,000-10,000 Users | Tier 4: 10,000+ Users |
|---------|---------------------|-------------------------|----------------------------|----------------------|
| **Render (Golang App)** | $7/month<br/>(Starter: 1GB RAM, 10GB SSD, 100GB bandwidth) | $25/month<br/>(Standard: 2GB RAM, 20GB SSD, 200GB bandwidth) | $85/month<br/>(Pro: 4GB RAM, 40GB SSD, 500GB bandwidth) | $170+/month<br/>(Pro+ or multiple instances) |
| **Neon (PostgreSQL)** | $15/month<br/>(1GB storage, 100 compute hours) | $50/month<br/>(5GB storage, 500 compute hours) | $150/month<br/>(25GB storage, 2,500 compute hours) | $500+/month<br/>(100GB+ storage, custom compute) |
| **Redis Cloud** | $6/month<br/>(1GB memory, basic plan) | $30/month<br/>(5GB memory, standard plan) | $120/month<br/>(25GB memory, pro plan) | $400+/month<br/>(100GB+ memory, enterprise) |
| **Firebase** | $0/month<br/>(Spark plan: 1GB storage, 100 concurrent connections) | $25/month<br/>(Blaze: ~$0.026/GB storage + usage) | $100/month<br/>(Blaze: higher usage + $0.15/GB downloads) | $500+/month<br/>(Enterprise custom pricing) |
| **Queue (Redis-based)** | Included in Redis | Included in Redis | Included in Redis | Included in Redis |
| **Total Estimated Cost** | **$28/month** | **$130/month** | **$455/month** | **$1,570+/month** |

### Cost Breakdown Notes

#### Render (Application Hosting)
- **Free Tier Available**: 750 hours/month for development/testing
- **Scaling**: Horizontal scaling available for higher tiers
- **Bandwidth**: Additional bandwidth costs may apply beyond included limits
- **Custom Domains**: Included in paid plans

#### Neon (PostgreSQL Database)
- **Free Tier**: 512MB storage, 100 compute hours (suitable for Tier 1)
- **Compute Hours**: Based on database activity; higher user loads increase compute needs
- **Storage**: $0.00015 per GB per hour; scales with user data growth
- **Connection Pooling**: Included for better performance

#### Redis Cloud
- **Free Tier**: 30MB memory for development
- **Memory Scaling**: Critical for caching user sessions and job data
- **Queue Functionality**: Redis can handle background jobs without additional services
- **High Availability**: Available in higher tiers

#### Firebase
- **Not Currently Used**: The application uses JWT and Redis for auth/session management
- **Potential Use Cases**: Could replace custom auth, add real-time notifications, or file storage
- **Free Tier**: Generous limits for small applications
- **Pay-as-you-go**: Costs scale with actual usage (storage, bandwidth, function invocations)

### Additional Cost Considerations

#### Bandwidth and Data Transfer
- **Render**: Included bandwidth may be insufficient for high-traffic apps
- **External APIs**: SerpAPI costs ~$50-200/month depending on search volume
- **CDN**: Consider Cloudflare or similar for static assets (~$20-100/month)

#### Monitoring and Observability
- **Application Monitoring**: DataDog, New Relic (~$15-100/month)
- **Error Tracking**: Sentry (~$26-200/month)
- **Uptime Monitoring**: UptimeRobot (free) or Pingdom (~$10-50/month)

#### Development and CI/CD
- **GitHub Actions**: Free for public repos, $0.008/minute for private
- **Docker Registry**: Docker Hub free tier, or GitHub Container Registry

#### Scaling Strategies
- **Tier 1-2**: Single instance with managed services
- **Tier 3+**: Load balancing, database read replicas, CDN
- **Cost Optimization**: Use reserved instances, monitor usage, implement caching

### Total Cost of Ownership (TCO)
For a production deployment serving 1,000-10,000 users:
- **Infrastructure**: ~$455/month (from table)
- **Third-party APIs**: ~$100/month (SerpAPI, monitoring)
- **Domain/SSL**: ~$20/month
- **Development Tools**: ~$50/month
- **Total**: ~$625/month


