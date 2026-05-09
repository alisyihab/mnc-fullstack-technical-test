# MNC Fullstack Technical Test

This repository contains the technical test solutions for the Fullstack Developer position. The project is divided into two main phases:

## Project Structure

```text
.
├── tahap-1/            # Logic & Algorithm Test
│   ├── shared/         # Common utilities (input, stack, etc)
│   ├── soal1/          # Matching Strings
│   ├── soal2/          # Cashier Change
│   ├── soal3/          # Bracket Validation
│   └── soal4/          # Employee Leave
├── tahap-2/            # REST API Implementation
│   ├── cmd/            # Application entry point
│   ├── docs/           # Swagger documentation
│   └── internal/       # Architecture layers (Domain, Usecase, Repo, Delivery)
├── go.mod              # Go module definition
└── README.md           # Root documentation
```

### [Phase 1 - Logic & Algorithm](./tahap-1/README.md)

Contains 4 algorithm problem solutions using Golang:

1.  **Matching Strings**: Case-insensitive string matching.
2.  **Cashier Change**: Calculating cashier change denominations.
3.  **Bracket Validation**: Validating bracket structure without regex.
4.  **Employee Leave**: Calculating pro-rata annual leave quota.

### [Phase 2 - REST API Implementation](./tahap-2/README.md)

REST API implementation using Gin, PostgreSQL, and Swagger:

- **Register & Login**: JWT authentication and PIN encryption.
- **Top Up & Payment**: User balance management.
- **Transfer**: Balance transfer between users with background processing.
- **Transaction History**: Transaction logs with dynamic response formats.
- **Profile Management**: Updating user profile information.
- **Swagger UI**: Interactive API documentation.

---

## Prerequisites

- [Go](https://golang.org/) (version 1.18+ recommended)
- [PostgreSQL](https://www.postgresql.org/) (version 13+ recommended)

## How to Run

### 1. Initial Setup

After cloning the repository, run the following command to download all required dependencies:

```bash
git clone https://github.com/alisyihab/mnc-fullstack-technical-test.git
cd mnc-fullstack-technical-test
go mod tidy
```

### 2. Phase 1 (Logic Test)

Each problem in Phase 1 can be run independently. Please read [Phase 1 README](./tahap-1/README.md) for detailed instructions.

### 3. Phase 2 (REST API)

The REST API in Phase 2 uses Gin, PostgreSQL, and Swagger. Please read [Phase 2 README](./tahap-2/README.md) for technical details, database configuration, and execution steps.
