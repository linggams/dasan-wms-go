# Testing Guide

Testing guide for DPPI ERP API.

## Running Tests

### Run All Tests

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

### Run Tests with Coverage Report

```bash
# Generate coverage file
go test ./... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### Run Specific Package Tests

```bash
# Domain tests
go test ./internal/domain/... -v

# Handler tests
go test ./internal/handler/... -v

# Middleware tests
go test ./internal/middleware/... -v

# Service tests
go test ./internal/service/... -v
```

## Test Structure

```
dppierp-api/
├── internal/
│   ├── domain/
│   │   └── models_test.go         # Domain model tests
│   ├── handler/
│   │   └── handler_test.go        # HTTP handler tests
│   ├── middleware/
│   │   └── middleware_test.go     # Middleware tests
│   └── service/
│       └── checkpoint_service_test.go  # Service tests
```

## Test Categories

### Domain Tests (`models_test.go`)

| Test | Description |
|------|-------------|
| `TestGetAllStages` | Verifies all 9 stages are returned |
| `TestIsValidStage` | Validates stage validation logic |
| `TestStageConstants` | Checks stage constant values |

### Handler Tests (`handler_test.go`)

| Test | Description |
|------|-------------|
| `TestSuccessResponse` | Tests success response format |
| `TestErrorResponse` | Tests error response format |
| `TestValidationErrorResponse` | Tests validation error format |
| `TestScanQRHandler_ValidationError` | Tests scan QR validation |
| `TestMoveStageHandler_MissingStageQuery` | Tests missing stage parameter |

### Middleware Tests (`middleware_test.go`)

| Test | Description |
|------|-------------|
| `TestAuthMiddleware_GenerateAndValidateToken` | JWT token generation/validation |
| `TestAuthMiddleware_InvalidToken` | Invalid token handling |
| `TestAuthMiddleware_MissingAuthHeader` | Missing auth header handling |
| `TestAuthMiddleware_InvalidAuthFormat` | Invalid auth format handling |
| `TestAuthMiddleware_ValidAuth` | Valid authentication flow |
| `TestCORSMiddleware` | CORS header configuration |
| `TestLogger` | Request logging middleware |

### Service Tests (`checkpoint_service_test.go`)

| Test | Description |
|------|-------------|
| `TestGetOverview` | Tests overview endpoint logic |
| `TestIsValidStage` | Tests stage validation in service |

## Manual API Testing

### Check Point Module

#### 1. Start the Server

```bash
go run ./cmd/api
```

#### 2. Health Check

```bash
curl http://localhost:8080/health
```

Expected output:
```json
{"status":"healthy"}
```

#### 3. Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dppi.com","password":"password123"}'
```

#### 4. Get Overview (with token)

```bash
TOKEN="<token-from-login>"
curl http://localhost:8080/check-point/v1/overview \
  -H "Authorization: Bearer $TOKEN"
```

#### 5. Scan QR

```bash
curl -X POST http://localhost:8080/check-point/v1/scan \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code":"F24120001"}'
```

#### 6. Move Stage

```bash
curl -X POST "http://localhost:8080/check-point/v1/move?stage=inventory" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "block_id": 1,
    "rack_id": 2,
    "entries": [{"code": "F24120001", "yard": 12.5}]
  }'
```

#### 7. Scan Rack

```bash
curl -X POST http://localhost:8080/check-point/v1/scan-rack \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code":"RACK-001"}'
```

#### 8. Relocation

```bash
curl -X POST http://localhost:8080/check-point/v1/relocation \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"current_rack_id": 1, "new_rack_id": 2}'
```

## Continuous Integration

Add to your CI/CD pipeline:

```yaml
# GitHub Actions example
- name: Run Tests
  run: go test ./... -v -race -coverprofile=coverage.out

- name: Upload Coverage
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage.out
```

## Writing New Tests

### Example Test Structure

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    expected := "expected value"

    // Act
    result := FunctionToTest()

    // Assert
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Table-Driven Tests

```go
func TestMultipleCases(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected bool
    }{
        {"valid case", "valid", true},
        {"invalid case", "invalid", false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := FunctionToTest(tc.input)
            if result != tc.expected {
                t.Errorf("Expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```
