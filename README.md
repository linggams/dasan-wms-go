# DPPI ERP API

API System for DPPI ERP, built with Go using clean and modern architecture.

## Modules

### 1. Check Point API

REST API for Check Point fabric management system that includes tracking fabrics through various production stages.

#### Features

- **JWT Authentication** - Secure bearer token authentication
- **Fabric Tracking** - Track fabrics through multiple production stages
- **Rack Management** - Scan racks and relocate items
- **Docker Ready** - Containerized deployment with Docker Compose
- **Clean Architecture** - Repository, Service, Handler pattern
- **Secure** - CORS, rate limiting, input validation

#### API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/login` | User login | ❌ |
| POST | `/auth/token/refresh` | Refresh access token | ❌ |
| POST | `/auth/logout` | Logout user | ✅ |
| POST | `/auth/forgot-password/request` | Request password reset | ❌ |
| POST | `/auth/forgot-password/reset` | Reset password | ❌ |
| GET | `/auth/me` | Get current user | ✅ |
| GET | `/profile` | Get user profile | ✅ |
| POST | `/profile/change-password` | Change user password | ✅ |
| GET | `/check-point/v1/overview` | Get all stages | ✅ |
| POST | `/check-point/v1/scan` | Scan fabric QR | ✅ |
| POST | `/check-point/v1/move?stage={stage}` | Move items to stage | ✅ |
| POST | `/check-point/v1/scan-rack` | Scan rack QR | ✅ |
| POST | `/check-point/v1/relocation` | Relocate rack items | ✅ |
| GET | `/check-point/v1/master/blocks` | Get all blocks | ✅ |
| GET | `/check-point/v1/master/racks` | Get all racks | ✅ |
| GET | `/check-point/v1/master/relaxation-blocks` | Get all relaxation blocks | ✅ |
| GET | `/check-point/v1/master/relaxation-racks` | Get all relaxation racks | ✅ |

---

## Tech Stack

- **Language**: Go 1.23+
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: MySQL 5.7+
- **Authentication**: JWT
- **Logging**: Zerolog

## Quick Start

### Prerequisites

- Go 1.23+
- MySQL 5.7+
- Docker & Docker Compose (for containerized deployment)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd dppierp-api
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run the API**
   ```bash
   go run ./cmd/api
   ```

### Docker Deployment

```bash
# Build and run
docker compose up -d

# View logs
docker compose logs -f api

# Stop
docker compose down
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/production) | development |
| `APP_PORT` | Server port | 8080 |
| `DB_HOST` | MySQL host | localhost |
| `DB_PORT` | MySQL port | 3306 |
| `DB_USER` | MySQL username | root |
| `DB_PASSWORD` | MySQL password | - |
| `DB_NAME` | MySQL database name | dppimes |
| `JWT_SECRET` | JWT signing secret | - |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | * |

## Project Structure

```
.
├── cmd/api/            # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── domain/         # Domain models
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # Auth, CORS, logging
│   ├── repository/     # Database layer
│   └── service/        # Business logic
├── pkg/
│   └── database/       # Database connection
├── docs/               # Documentation
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Authentication

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dppi.com","password":"password123"}'
```

### Use Token

```bash
curl http://localhost:8080/check-point/v1/overview \
  -H "Authorization: Bearer <your-token>"
```

## Documentation

- [Installation Guide](docs/installation.md) - Setup and deployment instructions
- [Testing Guide](docs/testing.md) - Unit testing and API testing documentation

## License

MIT License
