# Cars CRUD

A full-stack project for practicing backend, frontend, and infrastructure skills. The backend is a RESTful API built with **Go**, following **Clean Architecture** principles, with PostgreSQL, Redis, Kafka, and MongoDB.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Clean Architecture](#clean-architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Using Tilt (recommended)](#using-tilt-recommended)
  - [Manual Setup](#manual-setup)
- [Authentication (JWT)](#authentication-jwt)
- [API Endpoints](#api-endpoints)
- [Swagger Documentation](#swagger-documentation)
- [Middlewares](#middlewares)
- [Redis Cache Layer](#redis-cache-layer)
- [Kafka & MongoDB Logging](#kafka--mongodb-logging)
- [Input Validation](#input-validation)
- [Roadmap](#roadmap)

---

## Architecture Overview

```plaintext
Client → chi Router → Middlewares (logging, auth) → Handler → Usecase → Repository → PostgreSQL
                                                      └──→ Redis (cache layer) │ └──→ Kafka (log producer) → Consumer → MongoDB (request_logs)
```


## Clean Architecture

The project follows Clean Architecture to keep concerns separated and the codebase testable and maintainable. Each layer only depends on the layer directly below it:

| Layer | Path | Responsibility |
|---|---|---|
| **Domain** | `internal/domain/` | Entities, DTOs, and business rules. No external dependencies. |
| **Repository** | `internal/repository/` | Interface definitions for data access. Implementations live in subfolders (`postgres/`, `mongo/`). |
| **Usecase** | `internal/usecase/` | Business logic orchestration. Depends on repository interfaces and cache. |
| **Handler** | `internal/handler/` | HTTP transport layer. Parses requests, calls usecases, writes responses. |
| **Middleware** | `internal/middleware/` | Cross-cutting concerns (auth, logging) applied at the router level. |
| **Config** | `pkg/config/` | Application configuration loaded from environment variables. |

Dependencies always point **inward**: `Handler → Usecase → Repository (interface) ← Repository (implementation)`. This makes it easy to swap implementations (e.g., replace PostgreSQL with MySQL) without touching business logic.

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.23 |
| Router | chi v5 |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Message Queue | Apache Kafka (Confluent 7.6) |
| Log Storage | MongoDB 7 |
| Auth | JWT (HS256) |
| Docs | Swagger (swaggo) |
| Dev Tooling | Tilt, Docker Compose |

## Project Structure

```
.
├── backend/
│   └── go/
│       ├── cmd/api/main.go              # Application entrypoint
│       ├── internal/
│       │   ├── domain/                   # Entities & DTOs (Car, RequestLog, Auth)
│       │   ├── repository/               # Repository interfaces
│       │   │   ├── postgres/             # PostgreSQL implementation (cars)
│       │   │   └── mongo/                # MongoDB implementation (logs)
│       │   ├── usecase/                  # Business logic + cache integration
│       │   ├── handler/                  # HTTP handlers (cars, logs, auth)
│       │   ├── middleware/               # JWT auth & request logging
│       │   ├── cache/                    # Redis cache wrapper
│       │   └── queue/                    # Kafka producer & consumer
│       ├── pkg/config/                   # Environment config loader
│       ├── docs/                         # Generated swagger files
│       ├── Dockerfile
│       └── .env
├── frontend/                             # (future) React frontend
├── infra/
│   ├── dev/                              # (future) Dev environment infra
│   └── prod/                             # (future) Prod environment infra
├── docker-compose.yml                    # Infrastructure services
├── Tiltfile                              # Tilt dev orchestration
└── README.md
```


## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- [Tilt](https://tilt.dev/) (optional, recommended)
- swag CLI: `go install github.com/swaggo/swag/cmd/swag@latest`

### Using Tilt (recommended)

Tilt starts all infrastructure services (PostgreSQL, Redis, Kafka, Zookeeper, MongoDB) via Docker Compose and runs the Go API locally with live reload.

```bash
# Start everything (defaults to Go backend)
tilt up

# Or explicitly select the backend
tilt up -- --backend=go
```

Open the Tilt dashboard at http://localhost:10350 to see all resources, logs, and statuses. Press space in the terminal to open it automatically.

## Manual Setup

### 1. Start infrastructure services:

```bash
docker compose up -d
```
This starts: PostgreSQL (:5432), Redis (:6379), Kafka (:9092), Zookeeper (:2181), MongoDB (:27017).

### 2. Configure environment:

```bash
cd backend/go
cp .env.example .env
# Edit .env if needed
```

### 3. Generate Swagger docs:

```bash
swag init -g cmd/api/main.go -o docs
```

### 4. Run the API:

```bash
go mod tidy
go run ./cmd/api
```

The API is available at http://localhost:8080.

## Authentication (JWT)

The API uses API Key + JWT authentication:

1. Send your API key to POST /auth/validate to obtain a JWT token.
2. Include the token in subsequent requests as Authorization: Bearer <token>.
3. Tokens expire after 24 hours.

How it works:

The API_KEY is defined in .env. A client sends it to /auth/validate.
If valid, the server signs a JWT with the JWT_SECRET (HS256) and returns it.
Protected endpoints (cars, logs) require the JWT via the JWTAuth middleware.
Public endpoints (/auth/validate, /health, /swagger/*) do not require authentication.

Example — get a token:

```bash
curl -X POST http://localhost:8080/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"api_key": "my-api-key-12345"}'
```

Response:

```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

Example — use the token:

```bash
curl http://localhost:8080/api/v1/cars \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## API Endpoints

| **Method** | **Path** | **Auth** | **Description** |
|---|---|---|---|
| POST | `/auth/validate` | No | Validate API key, get JWT
| GET | `/api/v1/cars` | Yes | List all cars (paginated)
| GET | `/api/v1/cars/{id}` | Yes | Get a car by ID
| POST | `/api/v1/cars` | Yes | Create a new car
| PUT | `/api/v1/cars/{id}` | Yes | Update a car
| DELETE | `/api/v1/cars/{id}` | Yes | Soft-delete a car
| GET | `/api/v1/logs` | Yes | List request logs (paginated)
| GET | `/health` | No | Health check
| GET | `/swagger/*` | No | Swagger UI

## Swagger Documentation

API documentation is auto-generated using swaggo/swag from Go comments in handlers.

**Swagger UI**: http://localhost:8080/swagger/index.html
To regenerate after editing handler annotations:

```bash
cd backend/go
swag init -g cmd/api/main.go -o docs
```

In Swagger UI, click **"Authorize"** and enter Bearer <your-token> to authenticate protected endpoints.

## Middlewares

The API uses two custom middlewares applied at the router level:

**1. Request Logger (middleware.RequestLogger)**

Applied globally to all routes. Wraps every request to capture method, path, status code, duration, IP, and user agent. After the response is sent, it publishes a log entry to Kafka asynchronously. This means logging never blocks the request.

**2. JWT Auth (middleware.JWTAuth)**

Applied only to protected route groups (`/api/v1/cars, /api/v1/logs`). Extracts the Authorization: Bearer <token> header, parses and validates the JWT using HS256, and injects claims into the request context. Returns 401 Unauthorized if the token is missing, malformed, or expired.

Additionally, the following chi built-in middlewares are used:

- `RequestID` — assigns a unique ID to each request
- `RealIP` — extracts the real client IP from proxy headers
- `Recoverer` — recovers from panics and returns 500
- `Timeout` — sets a 30-second request timeout
- `CORS` — allows cross-origin requests

## Redis Cache Layer

Redis is used as a **read-through** cache for GET endpoints to reduce database load:

| **Cache Key Pattern** | **Endpoint** | **TTL** |
|---|---|---|
| `cars:{uuid}` | `GET /api/v1/cars/{id}` | 5 minutes
| `cars:list:{offset}:{limit}` | `GET /api/v1/cars` | 5 minutes

**Cache behavior:**

- **GET** requests first check Redis. On cache hit, the response is served directly from cache (no DB query).
- **On cache miss**, the data is fetched from PostgreSQL, then stored in Redis for subsequent requests.
- **Create, Update, Delete** operations **invalidate** related cache entries:
- - Single car cache (cars:{uuid}) is deleted on update/delete.
- - All list caches (cars:list:*) are deleted on create/update/delete using pattern-based scan.

## Kafka & MongoDB Logging

Every HTTP request is logged asynchronously through a Kafka → MongoDB pipeline:

```plaintext
Request → Logging Middleware → Kafka Producer → [car-api-logs topic] → Kafka Consumer → MongoDB
```

- The **logging middleware** captures request metadata and publishes it to the car-api-logs Kafka topic.
- The **Kafka consumer** runs as a background goroutine, reads messages from the topic, and inserts them into the request_logs collection in the cars_logs MongoDB database.
- This is fully **asynchronous** — the API response is never delayed by logging.

**MongoDB** cars_logs **database** —> request_logs collection example:

```json
[
  {
    "method": "POST",
    "path": "/auth/validate",
    "status_code": 200,
    "duration_ms": 2,
    "ip": "127.0.0.1:54932",
    "user_agent": "Mozilla/5.0 (X11; Linux x86_64)",
    "timestamp": "2026-02-19T15:30:12.441Z"
  },
  {
    "method": "GET",
    "path": "/api/v1/cars",
    "status_code": 200,
    "duration_ms": 15,
    "ip": "127.0.0.1:54932",
    "user_agent": "Mozilla/5.0 (X11; Linux x86_64)",
    "timestamp": "2026-02-19T15:30:45.112Z"
  },
  {
    "method": "POST",
    "path": "/api/v1/cars",
    "status_code": 201,
    "duration_ms": 8,
    "ip": "127.0.0.1:54932",
    "user_agent": "curl/8.5.0",
    "timestamp": "2026-02-19T15:31:02.887Z"
  },
  {
    "method": "DELETE",
    "path": "/api/v1/cars/550e8400-e29b-41d4-a716-446655440000",
    "status_code": 204,
    "duration_ms": 5,
    "ip": "127.0.0.1:54932",
    "user_agent": "curl/8.5.0",
    "timestamp": "2026-02-19T15:31:30.221Z"
  }
]
```

## Input Validation

Input validation is performed at the handler layer before data reaches the usecase:

**Create Car** (POST /api/v1/cars):

- `brand` — required, non-empty string
- `model` — required, non-empty string
- `year` — required, non-zero integer
- Invalid JSON body returns 400 Bad Request

**Update Car** (PUT /api/v1/cars/{id}):

- All fields are optional (partial update via pointer fields)
- `id` path param must be a valid UUID
- Invalid JSON body returns 400 Bad Request

**Get/Delete Car:**

- `id` path param must be a valid UUID, otherwise 400 Bad Request
- Non-existent car returns 404 Not Found

**List endpoints:**

- `offset` defaults to 0, limit defaults to 10 (cars) or 20 (logs)
- `limit` is capped at 100

**Auth:**

- `api_key` is required and must match the configured API_KEY

---

## Roadmap
- [ ] **Frontend** — React application for managing cars via the API
- [ ] **Infrastructure** — Separate environments inside `infra/` folder:
- - [ ] `infra/dev/` — Development environment (local/docker)
- - [ ] `infra/hml/` — Homologation/staging environment
- - [ ] `infra/prd/` — Production environment (cloud deploy)
- [ ] **Backend: PHP** — Reimplement the same API in PHP (backend/php/)
- [ ] **Backend: Node.js** — Reimplement the same API in Node.js (`backend/javascript/`)
- [ ] **Unit & integration tests** for Go backend
- [ ] **CI/CD pipeline** with GitHub Actions
- [ ] **Rate limiting** middleware
- [ ] **Structured logging** (replace log with slog or zap)
