# Go CRUD Application

A production-ready RESTful API built with Go, featuring user authentication, CRUD operations, and comprehensive security measures. Fully containerized with Docker and PostgreSQL.

## Features

- **User Authentication**: JWT-based authentication with secure password hashing (bcrypt)
- **CRUD Operations**: Complete Create, Read, Update, Delete functionality for users
- **Security First**:
  - Password strength validation (min 8 chars, uppercase, lowercase, number)
  - Rate limiting to prevent brute force attacks
  - SQL injection prevention via GORM ORM
  - CORS configuration
  - Non-root Docker container
  - Environment variable management for secrets
- **Database**: PostgreSQL with GORM ORM
- **Testing**: Comprehensive unit tests for core functionalities
- **Docker**: Multi-stage builds with Docker Compose orchestration

## Prerequisites

- Docker (version 20.10+)
- Docker Compose (version 2.0+)

**For local development without Docker:**
- Go 1.21+
- PostgreSQL 15+

## Quick Start

### 1. Clone and Setup

```bash
cd go-crud-app
cp .env.example .env
```

Edit `.env` file and update the following (especially `JWT_SECRET` for production):

```env
JWT_SECRET=your-very-long-random-secret-key-here
DB_PASSWORD=your-secure-database-password
```

> [!IMPORTANT]
> You **must** create and configure a `.env` file before running Docker Compose. 

### 3. Run with Docker Compose

```bash
docker compose up --build
```

The application will be available at `http://localhost:8080`

### 4. Verify Installation

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "time": "2026-01-21T12:00:00Z"
}
```

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

#### Register a New User
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-01-21T12:00:00Z",
    "updated_at": "2026-01-21T12:00:00Z"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-01-21T12:00:00Z",
    "updated_at": "2026-01-21T12:00:00Z"
  }
}
```

### Protected Endpoints (Require JWT Token)

All endpoints below require the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

#### Get Current User Profile
```http
GET /api/users/me
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2026-01-21T12:00:00Z",
  "updated_at": "2026-01-21T12:00:00Z"
}
```

#### List All Users (Excluding Current User)
```http
GET /api/users
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": 2,
      "username": "janedoe",
      "email": "jane@example.com",
      "created_at": "2026-01-21T12:00:00Z",
      "updated_at": "2026-01-21T12:00:00Z"
    }
  ],
  "count": 1
}
```

#### Get User by ID
```http
GET /api/users/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 2,
  "username": "janedoe",
  "email": "jane@example.com",
  "created_at": "2026-01-21T12:00:00Z",
  "updated_at": "2026-01-21T12:00:00Z"
}
```

#### Update User (Own Profile Only)
```http
PUT /api/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "john_updated",
  "email": "john.new@example.com"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "username": "john_updated",
  "email": "john.new@example.com",
  "created_at": "2026-01-21T12:00:00Z",
  "updated_at": "2026-01-21T12:05:00Z"
}
```

#### Delete User (Own Profile Only)
```http
DELETE /api/users/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

## Security Features

### 1. Authentication & Authorization
- **JWT Tokens**: Secure token-based authentication with 24-hour expiration
- **Password Hashing**: Bcrypt with cost factor 12
- **Password Requirements**:
  - Minimum 8 characters
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one number

### 2. Rate Limiting
- **Registration**: 3 requests per minute per IP
- **Login**: 5 requests per minute per IP
- **General Endpoints**: 100 requests per minute per IP

### 3. Input Validation
- Email format validation
- Username validation (3-50 alphanumeric characters and underscores)
- SQL injection prevention via GORM parameterization

### 4. Docker Security
- Multi-stage builds for minimal attack surface
- Non-root user in container
- Alpine Linux base image
- Health checks enabled

### 5. Environment Variables
- No hardcoded secrets
- All sensitive data in environment variables
- `.env.example` template provided

## Testing

### Run All Tests
```bash
# With Docker
docker compose exec app go test ./tests/... -v

# Without Docker (local development)
go test ./tests/... -v -cover
```

### Run Specific Test
```bash
go test ./tests/password_test.go -v
go test ./tests/jwt_test.go -v
```

### Test Coverage
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Development

### Local Development (Without Docker)

1. **Install Dependencies**
```bash
go mod download
```

2. **Setup PostgreSQL**
```bash
# Create database
createdb gocrud

# Or using psql
psql -U postgres -c "CREATE DATABASE gocrud;"
```

3. **Configure Environment**
```bash
cp .env.example .env
# Edit .env with your local database credentials
```

4. **Run Application**
```bash
go run cmd/server/main.go
```

### Project Structure
```
go-crud-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/
│   │   └── user.go              # User model
│   ├── handlers/
│   │   ├── auth.go              # Authentication handlers
│   │   └── user.go              # CRUD handlers
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   └── ratelimit.go         # Rate limiting
│   ├── database/
│   │   └── database.go          # Database connection
│   └── utils/
│       ├── jwt.go               # JWT utilities
│       └── password.go          # Password utilities
├── tests/                        # Unit tests
├── Dockerfile                    # Docker configuration
├── docker compose.yml           # Docker Compose setup
└── README.md                    # This file
```

## Troubleshooting

### Database Connection Issues

**Problem**: Application can't connect to database

**Solution**:
```bash
# Check if PostgreSQL container is running
docker compose ps

# View PostgreSQL logs
docker compose logs postgres

# Restart services
docker compose down
docker compose up --build
```

### Port Already in Use

**Problem**: Port 8080 or 5432 already in use

**Solution**: Edit `.env` file and change ports:
```env
PORT=8081
DB_PORT=5433
```

Then update `docker compose.yml` port mappings accordingly.

### JWT Token Expired

**Problem**: Getting "Invalid or expired token" error

**Solution**: Login again to get a new token. Tokens expire after 24 hours.

### Rate Limit Exceeded

**Problem**: Getting "Rate limit exceeded" error

**Solution**: Wait for the rate limit window to reset (1 minute) or adjust rate limits in `cmd/server/main.go`.

## Environment Variables

| Variable | Description | Note |
|----------|-------------|---------|
| `DB_HOST` | Database host | Required |
| `DB_PORT` | Database port | Required |
| `DB_USER` | Database user | Required |
| `DB_PASSWORD` | Database password | Required |
| `DB_NAME` | Database name | Required |
| `DB_SSLMODE` | SSL mode for database | Required (default `disable` in app) |
| `JWT_SECRET` | Secret key for JWT signing | ⚠️ **Must change in production** |
| `PORT` | Application port | Required |
| `CORS_ORIGIN` | Allowed CORS origins | Required |

## Production Deployment

### Important Security Considerations

1. **Change JWT Secret**: Use a long, random string
```bash
# Generate a secure secret
openssl rand -base64 64
```

2. **Use Strong Database Password**
```bash
# Generate a secure password
openssl rand -base64 32
```

3. **Enable SSL for Database**
```env
DB_SSLMODE=require
```

4. **Configure CORS Properly**
```env
CORS_ORIGIN=https://yourdomain.com
```

5. **Use Environment-Specific Configurations**
- Never commit `.env` file to version control
- Use secrets management (e.g., Docker secrets, Kubernetes secrets)

## License

This project is licensed under the MIT License.

## 👥 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions, please open an issue on the repository.
