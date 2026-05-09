# Phase 2 - REST API (E-Wallet)

REST API for a simple E-Wallet system built with Go (Gin), PostgreSQL, and GORM.

## Folder Structure

```
tahap-2/
├── cmd/
│   └── main.go                         # Application entry point
├── migrations/                         # SQL migration files (golang-migrate)
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_transactions_table.up.sql
│   └── 000002_create_transactions_table.down.sql
├── docs/                               # Swagger documentation (auto-generated)
├── internal/
│   ├── app/                            # Dependency injection & app wiring
│   ├── delivery/http/
│   │   ├── handlers/                   # HTTP request handlers
│   │   ├── middleware/                 # Auth middleware (JWT)
│   │   └── response/                  # Consistent JSON response helper
│   ├── domain/
│   │   ├── models/                     # GORM models
│   │   ├── repository/                 # Repository interfaces
│   │   │   └── mocks/                 # Testify mocks for unit tests
│   │   └── usecase/                    # Usecase interfaces
│   ├── infrastructure/
│   │   ├── auth/                       # JWT service
│   │   ├── config/                     # Environment config loader
│   │   ├── database/                   # Postgres connection & migration runner
│   │   └── worker/                     # Background transfer worker (channel queue)
│   ├── repository/                     # GORM repository implementations
│   └── usecase/                        # Business logic implementations
├── Makefile                            # Common dev commands
├── .env.example                        # Environment variables template
└── README.md
```

## Features

| #   | Feature            | Endpoint                     | Auth       |
| --- | ------------------ | ---------------------------- | ---------- |
| 1   | Register           | `POST /register`             | No         |
| 2   | Login              | `POST /login`                | No         |
| 3   | Top Up             | `POST /topup`                | Bearer JWT |
| 4   | Payment            | `POST /pay`                  | Bearer JWT |
| 5   | Transfer           | `POST /transfer`             | Bearer JWT |
| 6   | Transaction Report | `GET /transactions`          | Bearer JWT |
| 7   | Update Profile     | `PUT /profile`               | Bearer JWT |
| 8   | Queue Stats        | `GET /dashboard/queue-stats` | No         |
| 9   | Recent Jobs        | `GET /dashboard/jobs`        | No         |

## Tech Stack

- **Go** — Gin framework
- **PostgreSQL** — GORM ORM
- **golang-migrate** — versioned SQL migrations
- **JWT** — access token (24h) + refresh token (7d)
- **Bcrypt** — PIN hashing
- **Background Worker** — channel-based worker pool for async transfers
- **Swagger** — API documentation

## How to Run

### 1. Configure environment

**Linux / macOS**

```bash
cp .env.example .env
```

**Windows**

```cmd
copy .env.example .env
```

Edit `.env`:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mnc_test
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=your_jwt_secret_key
SERVER_PORT=8080
```

### 2. Run the application

Migrations run **automatically** on startup.

**Linux / macOS**

```bash
make run
```

**Windows**

```cmd
go run cmd/main.go
```

### 3. Run migrations manually

Install the golang-migrate CLI (one-time):

**Linux / macOS**

```bash
make migrate-install
```

**Windows**

```cmd
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Apply / roll back:

**Linux / macOS**

```bash
make migrate-up          # apply all pending migrations
make migrate-down        # roll back last migration
make migrate-version     # show current version
```

**Windows**

```cmd
migrate -path migrations -database "postgres://postgres:your_password@localhost:5432/mnc_test?sslmode=disable" up
migrate -path migrations -database "postgres://postgres:your_password@localhost:5432/mnc_test?sslmode=disable" down 1
migrate -path migrations -database "postgres://postgres:your_password@localhost:5432/mnc_test?sslmode=disable" version
```

### 4. Load example data

**Linux / macOS**

```bash
make seed
```

**Windows**

```cmd
psql -U postgres -d mnc_test -f ..\example_data.sql
```

Example credentials (PIN: `123456`):

- Phone: `0811255501` — Guntur Saputro (balance: 370,000)
- Phone: `0811255502` — Tom Araya (balance: 30,000)

### 5. Swagger documentation

```
http://localhost:8080/swagger/index.html
```

---

## Running Unit Tests

**Linux / macOS**

```bash
make test              # run tests with verbose output (recommended)
make test-short        # run tests without verbose output
make test-cover        # coverage per package
make test-cover-html   # HTML coverage report (opens browser)
```

**Windows**

```cmd
go test -v mnc-fullstack-technical-test/tahap-2/internal/usecase
```

Tests are located in `internal/usecase/` and cover:

- `user_usecase_test.go` — Register (success, phone already exists), Login (success, wrong PIN)
- `transaction_usecase_test.go` — TopUp, Payment (success + insufficient balance), Transfer (success + insufficient balance + target not found)

---

## Background Transfer Worker

Transfer requests are processed **asynchronously** via an in-memory channel queue:

1. `/transfer` endpoint deducts the sender's balance immediately and creates a `PENDING` transaction.
2. The job is enqueued to a worker pool (3 goroutines, 100-slot buffer).
3. Workers credit the receiver's balance and update the transaction status to `SUCCESS` (or `FAILED` on error).

### Monitoring dashboard

| Endpoint                     | Description                                                      |
| ---------------------------- | ---------------------------------------------------------------- |
| `GET /dashboard/queue-stats` | Live worker stats: pending, processed, failed jobs, worker count |
| `GET /dashboard/jobs`        | Last 100 transfer job records with status                        |

Example response for `/dashboard/queue-stats`:

```json
{
  "jobs_failed": 0,
  "jobs_pending": 0,
  "jobs_processed": 3,
  "worker_count": 3
}
```

---

## Architecture

Clean architecture with four layers:

```
Delivery (HTTP) → Usecase (Business Logic) → Repository (Data Access) → Infrastructure
```

- **Domain** — models and interface contracts (no dependencies on frameworks)
- **Usecase** — business rules, depends only on domain interfaces
- **Repository** — GORM implementations of domain repository interfaces
- **Delivery** — Gin handlers, middleware, JSON response helpers
- **Infrastructure** — JWT, config, database connection, migration runner, worker pool
