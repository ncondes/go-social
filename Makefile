include .env
MIGRATIONS_PATH = ./cmd/db/migrations

# ==================== Migrations ====================

.PHONY: migration
migration:
	@echo "Creating migration file..."
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
	@echo "Migration file created successfully"

.PHONY: migrate-up
migrate-up:
	@echo "Applying migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) up
	@echo "Migrations applied successfully"

.PHONY: migrate-test-up
migrate-test-up:
	@echo "Applying migrations to test database..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_TEST_ADDR) up
	@echo "Migrations applied successfully"

.PHONY: migrate-test-down
migrate-test-down:
	@echo "Rolling back migrations to test database..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_TEST_ADDR) down
	@echo "Migrations rolled back successfully"

# ==================== Database ====================

.PHONY: db-start
db-start:
	@echo "Starting development database..."
	@docker compose up -d social-db
	@echo "Development database started"

.PHONY: db-stop
db-stop:
	@echo "Stopping development database..."
	@docker compose stop social-db
	@echo "Development database stopped"

.PHONY: db-start-test
db-start-test:
	@echo "Starting test database..."
	@docker compose up -d social-test-db
	@echo "Test database started"

.PHONY: db-stop-test
db-stop-test:
	@echo "Stopping test database..."
	@docker compose stop social-test-db
	@echo "Test database stopped"

.PHONY: db-start-all
db-start-all:
	@echo "Starting all databases..."
	@docker compose up -d
	@echo "All databases started"

.PHONY: db-stop-all
db-stop-all:
	@echo "Stopping all databases..."
	@docker compose stop
	@echo "All databases stopped"

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	@go run cmd/db/seed/main.go
	@echo "Database seeded successfully"

.PHONY: db-flush
db-flush:
	@echo "Flushing database..."
	@go run cmd/db/flush/main.go
	@echo "Database flushed successfully"

# ==================== Testing ====================

.PHONY: test-setup
test-setup:
	@echo "Setting up test database..."
	@docker compose up -d social-test-db
	@sleep 2
	@echo "Dropping existing test database..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_TEST_ADDR) drop -f || true
	@echo "Running migrations..."
	@$(MAKE) migrate-test-up
	@echo "Test database ready!"

.PHONY: test
test:
	@echo "Running all tests..."
	@$(MAKE) test-unit
	@$(MAKE) test-integration
	@echo "All tests completed successfully"

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -count=1 \
		./internal/handlers/... \
		./internal/services/... \
		./packages/...
	@echo "Unit tests completed"

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@go test -v -race -count=1 -tags=integration \
		./internal/repositories/...
	@echo "Integration tests completed"

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# ==================== Development ====================

.PHONY: dev
dev:
	@echo "Starting development server with hot reload..."
	@air

# ==================== Swagger ====================
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -d ./cmd/api/,./internal/handlers,./internal/dtos,./internal/domain -o ./docs && swag fmt
	@echo "Swagger documentation generated successfully"

# ==================== Utilities ====================

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  Migrations:"
	@echo "    migration         - Create a new migration"
	@echo "    migrate-up        - Apply migrations to dev database"
	@echo "    migrate-down      - Rollback one migration from dev database"
	@echo "    migrate-test-up   - Apply migrations to test database"
	@echo "    migrate-test-down - Rollback one migration from test database"
	@echo ""
	@echo "  Database:"
	@echo "    db-start          - Start development database"
	@echo "    db-stop           - Stop development database"
	@echo "    db-start-test     - Start test database"
	@echo "    db-stop-test      - Stop test database"
	@echo "    db-start-all      - Start all databases"
	@echo "    db-seed           - Seed development database"
	@echo "    db-flush          - Flush development database"
	@echo ""
	@echo "  Testing:"
	@echo "    test-setup        - Setup test database (run once)"
	@echo "    test              - Run all tests"
	@echo "    test-unit         - Run unit tests only"
	@echo "    test-integration  - Run integration tests only"
	@echo "    test-coverage     - Run tests with coverage report"
	@echo ""
	@echo "  Development:"
	@echo "    dev               - Run with hot reload (requires air)"
	@echo ""
	@echo "  Swagger:"
	@echo "    swagger           - Generate Swagger documentation"
	@echo ""
%:
	@: 
# This is for avoiding errors
