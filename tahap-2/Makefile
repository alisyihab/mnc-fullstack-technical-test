include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATE=migrate -path migrations -database "$(DB_URL)"

.PHONY: run build test test-verbose test-cover \
        migrate-up migrate-down migrate-create migrate-version migrate-force \
        seed swag help

## ── App ──────────────────────────────────────────────────────────────────────

run: ## Run the application (migrations execute automatically on startup)
	go run cmd/main.go

build: ## Build binary to ./bin/app
	go build -o bin/app cmd/main.go

## ── Database ─────────────────────────────────────────────────────────────────

migrate-install: ## Install golang-migrate CLI
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Apply all pending migrations
	$(MIGRATE) up

migrate-down: ## Roll back the last migration
	$(MIGRATE) down 1

migrate-version: ## Show current migration version
	$(MIGRATE) version

migrate-force: ## Force-set migration version (usage: make migrate-force V=1)
	$(MIGRATE) force $(V)

migrate-create: ## Create a new migration (usage: make migrate-create NAME=add_column_x)
	$(MIGRATE) create -ext sql -dir migrations -seq $(NAME)

seed: ## Load example data into the database
	psql "$(DB_URL)" -f ../example_data.sql

## ── Tests ────────────────────────────────────────────────────────────────────

TEST_PKGS = $(shell go list ./... | xargs go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}')

test: ## Run all unit tests (verbose)
	go test -v $(TEST_PKGS)

test-short: ## Run all unit tests without verbose output
	go test $(TEST_PKGS)

test-cover: ## Run tests and show coverage per package
	go test -cover $(TEST_PKGS)

test-cover-html: ## Generate HTML coverage report (opens in browser)
	go test -coverprofile=coverage.out $(TEST_PKGS) && go tool cover -html=coverage.out

## ── Docs ─────────────────────────────────────────────────────────────────────

swag: ## Regenerate Swagger docs (requires swag CLI)
	swag init -g cmd/main.go -o docs

## ── Help ─────────────────────────────────────────────────────────────────────

help: ## List all available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
