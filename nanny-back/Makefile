.PHONY: help build run test test-coverage test-unit test-integration clean docker-build docker-up docker-down docker-logs migrate-up migrate-down migrate-create lint fmt vet db-reset db-seed install-tools

# Variables
BINARY_NAME=nanny-backend
MAIN_PATH=./cmd/api
MIGRATION_DIR=./migrations
DATABASE_URL=postgres://postgres:Ana4aBada$$@localhost:5432/nanny_db?sslmode=disable

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "$(BLUE)Nanny Platform - Available Commands:$(NC)"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  make run              - Run the application locally"
	@echo "  make build            - Build the application binary"
	@echo "  make clean            - Remove binary and temporary files"
	@echo "  make install-tools    - Install development tools"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test             - Run all tests"
	@echo "  make test-unit        - Run unit tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-verbose     - Run tests with verbose output"
	@echo ""
	@echo "$(GREEN)Code Quality:$(NC)"
	@echo "  make lint             - Run golangci-lint"
	@echo "  make fmt              - Format code with gofmt"
	@echo "  make vet              - Run go vet"
	@echo "  make check            - Run all checks (fmt, vet, lint)"
	@echo ""
	@echo "$(GREEN)Database:$(NC)"
	@echo "  make migrate-up       - Apply all pending migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-create   - Create new migration (use: make migrate-create name=add_users)"
	@echo "  make db-reset         - Drop and recreate database"
	@echo "  make db-seed          - Seed database with test data"
	@echo ""
	@echo "$(GREEN)Docker:$(NC)"
	@echo "  make docker-build     - Build Docker images"
	@echo "  make docker-up        - Start all services with docker-compose"
	@echo "  make docker-down      - Stop all services"
	@echo "  make docker-logs      - View logs from all services"
	@echo "  make docker-clean     - Remove all containers and volumes"
	@echo ""

## build: Build the application binary
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: bin/$(BINARY_NAME)$(NC)"

## run: Run the application locally
run:
	@echo "$(BLUE)Starting $(BINARY_NAME)...$(NC)"
	@go run $(MAIN_PATH)

## clean: Remove binary and temporary files
clean:
	@echo "$(BLUE)Cleaning up...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

## test: Run all tests
test:
	@echo "$(BLUE)Running all tests...$(NC)"
	@go test -v -race -timeout 30s ./...
	@echo "$(GREEN)✓ All tests passed$(NC)"

## test-unit: Run unit tests only
test-unit:
	@echo "$(BLUE)Running unit tests...$(NC)"
	@go test -v -short -race ./...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

## test-integration: Run integration tests only
test-integration:
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -v -run Integration ./...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

## test-verbose: Run tests with verbose output
test-verbose:
	@go test -v -race -count=1 ./... | grep -v "no test files"

## lint: Run golangci-lint
lint:
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@golangci-lint run --timeout 5m
	@echo "$(GREEN)✓ Linting passed$(NC)"

## fmt: Format code with gofmt
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@gofmt -s -w .
	@goimports -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

## check: Run all checks (fmt, vet, lint)
check: fmt vet lint
	@echo "$(GREEN)✓ All checks passed$(NC)"

## migrate-up: Apply all pending migrations
migrate-up:
	@echo "$(BLUE)Applying migrations...$(NC)"
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) up
	@echo "$(GREEN)✓ Migrations applied$(NC)"

## migrate-down: Rollback last migration
migrate-down:
	@echo "$(BLUE)Rolling back migration...$(NC)"
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) down 1
	@echo "$(GREEN)✓ Migration rolled back$(NC)"

## migrate-create: Create new migration (use: make migrate-create name=add_users)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: name parameter is required$(NC)"; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@echo "$(BLUE)Creating migration: $(name)$(NC)"
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)
	@echo "$(GREEN)✓ Migration created$(NC)"

## db-reset: Drop and recreate database
db-reset:
	@echo "$(BLUE)Resetting database...$(NC)"
	@psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS nanny_db;"
	@psql -h localhost -U postgres -c "CREATE DATABASE nanny_db;"
	@make migrate-up
	@echo "$(GREEN)✓ Database reset complete$(NC)"

## db-seed: Seed database with test data
db-seed:
	@echo "$(BLUE)Seeding database...$(NC)"
	@psql -h localhost -U postgres -d nanny_db -f scripts/seeds.sql
	@echo "$(GREEN)✓ Database seeded$(NC)"

## docker-build: Build Docker images
docker-build:
	@echo "$(BLUE)Building Docker images...$(NC)"
	@docker-compose build
	@echo "$(GREEN)✓ Docker images built$(NC)"

## docker-up: Start all services with docker-compose
docker-up:
	@echo "$(BLUE)Starting services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"
	@echo ""
	@echo "Access the application at: http://localhost:8080"

## docker-down: Stop all services
docker-down:
	@echo "$(BLUE)Stopping services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

## docker-logs: View logs from all services
docker-logs:
	@docker-compose logs -f

## docker-clean: Remove all containers and volumes
docker-clean:
	@echo "$(BLUE)Cleaning up Docker resources...$(NC)"
	@docker-compose down -v --remove-orphans
	@echo "$(GREEN)✓ Docker cleanup complete$(NC)"

## swagger: Generate Swagger documentation
swagger:
	@echo "$(BLUE)Generating Swagger documentation...$(NC)"
	@swag init -g cmd/api/main.go -o docs
	@echo "$(GREEN)✓ Swagger docs generated$(NC)"

## deps: Download and tidy dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## vendor: Create vendor directory
vendor:
	@echo "$(BLUE)Creating vendor directory...$(NC)"
	@go mod vendor
	@echo "$(GREEN)✓ Vendor directory created$(NC)"

## ci: Run CI pipeline (checks + tests)
ci: check test
	@echo "$(GREEN)✓ CI pipeline passed$(NC)"

## pre-commit: Run pre-commit checks
pre-commit: fmt vet lint test-unit
	@echo "$(GREEN)✓ Pre-commit checks passed$(NC)"
