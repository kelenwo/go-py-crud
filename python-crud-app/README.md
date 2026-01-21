# Python FastAPI CRUD Application

A RESTful API built with Python FastAPI, featuring user authentication, CRUD operations, and comprehensive security measures. Containerized with Docker and PostgreSQL.

## Features

- **User Authentication**: JWT-based authentication with secure password hashing (bcrypt)
- **CRUD Operations**: Complete Create, Read, Update, Delete functionality for users
- **Security First**:
  - Password strength validation (min 8 chars, uppercase, lowercase, number)
  - Rate limiting to prevent brute force attacks
  - SQL injection prevention via SQLAlchemy ORM
  - CORS configuration
  - Non-root Docker container
  - Environment variable management for secrets
  - Pydantic validation for all inputs
- **Database**: PostgreSQL with SQLAlchemy ORM
- **Testing**: Comprehensive unit and integration tests with pytest
- **Docker**: Multi-stage builds with Docker Compose orchestration
- **API Documentation**: Auto-generated OpenAPI (Swagger) documentation

## Prerequisites

- Docker (version 20.10+)
- Docker Compose (version 2.0+)

**For local development without Docker:**
- Python 3.11+
- PostgreSQL 15+

## Quick Start

### 1. Clone and Setup

```bash
cd python-crud-app
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
docker-compose up --build
```

The application will be available at `http://localhost:8000`

### 4. Verify Installation

```bash
curl http://localhost:8000/health
```

Expected response:
```json
{
  "status": "healthy",
  "time": "2026-01-21T12:00:00Z"
}
```

### 5. Access API Documentation

Open your browser and navigate to:
- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc

## API Documentation

### Base URL
```
http://localhost:8000/api
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
    "created_at": "2026-01-21T12:00:00.000000",
    "updated_at": "2026-01-21T12:00:00.000000"
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
    "created_at": "2026-01-21T12:00:00.000000",
    "updated_at": "2026-01-21T12:00:00.000000"
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
  "created_at": "2026-01-21T12:00:00.000000",
  "updated_at": "2026-01-21T12:00:00.000000"
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
      "created_at": "2026-01-21T12:00:00.000000",
      "updated_at": "2026-01-21T12:00:00.000000"
    }
  ],
  "count": 1
}
```

#### Get User by ID
```http
GET /api/users/{user_id}
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 2,
  "username": "janedoe",
  "email": "jane@example.com",
  "created_at": "2026-01-21T12:00:00.000000",
  "updated_at": "2026-01-21T12:00:00.000000"
}
```

#### Update User (Own Profile Only)
```http
PUT /api/users/{user_id}
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
  "created_at": "2026-01-21T12:00:00.000000",
  "updated_at": "2026-01-21T12:05:00.000000"
}
```

#### Delete User (Own Profile Only)
```http
DELETE /api/users/{user_id}
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
- **Password Hashing**: Bcrypt via passlib
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
- **Pydantic Models**: Automatic request/response validation
- Email format validation via email-validator
- Username validation (3-50 alphanumeric characters and underscores)
- SQL injection prevention via SQLAlchemy parameterization

### 4. Docker Security
- Multi-stage builds for minimal attack surface
- Non-root user in container
- Slim Python base image
- Health checks enabled

### 5. Environment Variables
- No hardcoded secrets
- All sensitive data in environment variables
- `.env.example` template provided

## 🧪 Testing

### Install Test Dependencies
```bash
pip install pytest pytest-cov httpx
```

### Run All Tests
```bash
# With Docker
docker-compose exec app pytest tests/ -v

# Without Docker (local development)
pytest tests/ -v
```

### Run with Coverage
```bash
pytest tests/ -v --cov=app --cov-report=html
```

### Run Specific Test File
```bash
pytest tests/test_security.py -v
pytest tests/test_auth.py -v
```

## 🛠️ Development

### Local Development (Without Docker)

1. **Create Virtual Environment**
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

2. **Install Dependencies**
```bash
pip install -r requirements.txt
```

3. **Setup PostgreSQL**
```bash
# Create database
createdb pythoncrud

# Or using psql
psql -U postgres -c "CREATE DATABASE pythoncrud;"
```

4. **Configure Environment**
```bash
cp .env.example .env
# Edit .env with your local database credentials
```

5. **Run Application**
```bash
uvicorn app.main:app --reload --port 8000
```

The application will be available at `http://localhost:8000`

### Project Structure
```
python-crud-app/
├── app/
│   ├── __init__.py
│   ├── main.py                  # FastAPI application
│   ├── models/
│   │   ├── __init__.py
│   │   └── user.py              # SQLAlchemy models
│   ├── schemas/
│   │   ├── __init__.py
│   │   └── user.py              # Pydantic schemas
│   ├── routers/
│   │   ├── __init__.py
│   │   ├── auth.py              # Authentication routes
│   │   └── users.py             # User CRUD routes
│   ├── dependencies/
│   │   ├── __init__.py
│   │   ├── auth.py              # Auth dependencies
│   │   └── database.py          # DB dependencies
│   ├── core/
│   │   ├── __init__.py
│   │   ├── config.py            # Configuration
│   │   ├── security.py          # Security utilities
│   │   └── database.py          # Database connection
│   └── middleware/
│       ├── __init__.py
│       └── ratelimit.py         # Rate limiting
├── tests/                        # Unit tests
├── Dockerfile                    # Docker configuration
├── docker-compose.yml           # Docker Compose setup
├── requirements.txt             # Python dependencies
└── README.md                    # This file
```

## Troubleshooting

### Database Connection Issues

**Problem**: Application can't connect to database

**Solution**:
```bash
# Check if PostgreSQL container is running
docker-compose ps

# View PostgreSQL logs
docker-compose logs postgres

# Restart services
docker-compose down
docker-compose up --build
```

### Port Already in Use

**Problem**: Port 8000 or 5432 already in use

**Solution**: Edit `.env` file and change ports:
```env
PORT=8001
DB_PORT=5433
```

Then update `docker-compose.yml` port mappings accordingly.

### JWT Token Expired

**Problem**: Getting "Invalid or expired token" error

**Solution**: Login again to get a new token. Tokens expire after 24 hours.

### Rate Limit Exceeded

**Problem**: Getting "Rate limit exceeded" error

**Solution**: Wait for the rate limit window to reset (1 minute) or adjust rate limits in `app/main.py`.

### Import Errors

**Problem**: Getting import errors when running tests

**Solution**: Make sure you're in the project root directory and have installed all dependencies:
```bash
pip install -r requirements.txt
```

## Environment Variables

| Variable | Description | Note |
|----------|-------------|---------|
| `DB_HOST` | Database host | Required |
| `DB_PORT` | Database port | Required |
| `DB_USER` | Database user | Required |
| `DB_PASSWORD` | Database password | Required |
| `DB_NAME` | Database name | Required |
| `JWT_SECRET` | Secret key for JWT signing | ⚠️ **Must change in production** |
| `JWT_ALGORITHM` | JWT algorithm | Optional (default `HS256`) |
| `JWT_EXPIRATION_HOURS` | Token expiration in hours | Optional (default `24`) |
| `PORT` | Application port | Required |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | Required |
| `ENV` | Environment (development/production) | Required |

## Production Deployment

### Important Security Considerations

1. **Change JWT Secret**: Use a long, random string
```bash
# Generate a secure secret
python -c "import secrets; print(secrets.token_urlsafe(64))"
```

2. **Use Strong Database Password**
```bash
# Generate a secure password
python -c "import secrets; print(secrets.token_urlsafe(32))"
```

3. **Configure CORS Properly**
```env
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

4. **Set Environment to Production**
```env
ENV=production
```

5. **Use Environment-Specific Configurations**
- Never commit `.env` file to version control
- Use secrets management (e.g., Docker secrets, Kubernetes secrets)
- Enable SSL/TLS for database connections

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions, please open an issue on the repository.

## Additional Resources

- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [SQLAlchemy Documentation](https://docs.sqlalchemy.org/)
- [Pydantic Documentation](https://docs.pydantic.dev/)
- [Docker Documentation](https://docs.docker.com/)
