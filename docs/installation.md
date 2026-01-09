# Installation Guide

Complete guide for installation and deployment of DPPI ERP API.

## Prerequisites

- Go 1.23 or newer
- MySQL 5.7 or newer (existing database `dppimes`)
- Docker & Docker Compose (optional, for containerized deployment)
- Git

## Option 1: Local Installation

### Step 1: Clone Repository

```bash
git clone <repository-url>
cd dppierp-api
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Configure Environment

```bash
# Copy environment file
cp .env.example .env

# Edit with your database credentials
nano .env
```

**`.env` Configuration:**

```env
# Environment
APP_ENV=development
APP_PORT=8080

# Database (adjust to your MySQL credentials)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=dppimes

# JWT (change to a minimum 32-character secure secret)
JWT_SECRET=your-super-secret-key-minimum-32-characters

# CORS
CORS_ALLOWED_ORIGINS=*
```

### Step 4: Run the Application

```bash
# Development mode
go run ./cmd/api

# Or build and run binary
go build -o dppierp-api ./cmd/api
./dppierp-api
```

### Step 5: Verify Installation

```bash
# Health check
curl http://localhost:8080/health

# Expected output: {"status":"healthy"}
```

## Option 2: Docker Deployment

### Step 1: Configure Environment

```bash
cp .env.example .env
# Edit .env with appropriate credentials
```

### Step 2: Build and Run

```bash
# Build and run containers
docker compose up -d --build

# Check status
docker compose ps

# View logs
docker compose logs -f api
```

### Step 3: Verify

```bash
curl http://localhost:8080/health
```

## Database Setup

This API uses an existing MySQL database (`dppimes`). Ensure the database already has the following tables:

### Check Point Module Tables

- `fabrics`
- `inventories`
- `racks`
- `blocks`
- `buyers`
- `fabric_incomings`
- `orders`

If you need to import the database from a dump:

```bash
mysql -u root -p dppimes < db_dump.sql
```

## Production Deployment

### Recommended Configuration

```env
APP_ENV=production
APP_PORT=8080

# Use appropriate connection pooling
DB_HOST=your-production-db-host
DB_PORT=3306
DB_USER=app_user
DB_PASSWORD=secure_password
DB_NAME=dppimes

# Use a very secure secret (minimum 32 characters)
JWT_SECRET=your-very-secure-production-secret-key-here

# Restrict CORS to valid domains
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://app.your-domain.com
```

### Docker Production Build

```bash
# Build production image
docker build -t dppierp-api:latest .

# Run with environment variables
docker run -d \
  --name dppierp-api \
  -p 8080:8080 \
  -e APP_ENV=production \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e JWT_SECRET=your-secret \
  dppierp-api:latest
```

### Using with Reverse Proxy (Nginx)

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Troubleshooting

### Database Connection Failed

```
Error: Failed to connect to database
```

**Solution:**
1. Ensure MySQL server is running
2. Verify credentials in `.env`
3. Ensure user has access to database

### Port Already in Use

```
Error: listen tcp :8080: bind: address already in use
```

**Solution:**
1. Change `APP_PORT` in `.env`
2. Or stop the service using that port

### Docker: Host.docker.internal Not Found

**Solution (Linux):**
Ensure `extra_hosts` is configured in `docker-compose.yml`:

```yaml
extra_hosts:
  - "host.docker.internal:host-gateway"
```

## Health Check Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Application health check |

## Support

If you encounter issues, please:
1. Check logs: `docker compose logs -f` or check console output
2. Verify database connection
3. Ensure all environment variables are configured correctly
