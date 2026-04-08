include .env
MIGRATIONS_PATH = ./cmd/db/migrations

.DEFAULT_GOAL := help

# ==================== Docker & Project Setup ====================

.PHONY: setup
setup:
	@echo "Setting up the project..."
	@echo "Starting Docker services..."
	@docker compose up -d
	@sleep 3
	@echo "Resetting database migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) drop -f || true
	@echo "Running migrations..."
	@$(MAKE) migrate-up
	@echo "Seeding database..."
	@$(MAKE) db-seed
	@echo "Setup complete! Run \033[1;32mmake dev\033[0m to start the API"

.PHONY: docker-up
docker-up:
	@echo "Starting all Docker services..."
	@docker compose up -d
	@echo "Services started:"
	@echo "   - PostgreSQL (dev):  localhost:5432"
	@echo "   - PostgreSQL (test): localhost:5433"
	@echo "   - Redis:             localhost:6379"
	@echo "   - Web UI:            http://localhost:5173"
	@echo "   - Redis Commander:   http://localhost:8081"

.PHONY: docker-down
docker-down:
	@echo "Stopping all Docker services..."
	@docker compose down
	@echo "All services stopped"

.PHONY: docker-restart
docker-restart:
	@echo "Restarting all Docker services..."
	@docker compose restart
	@echo "Services restarted"

.PHONY: docker-logs
docker-logs:
	@docker compose logs -f

.PHONY: docker-ps
docker-ps:
	@docker compose ps

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

.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	@$(MAKE) db-flush
	@$(MAKE) migrate-up
	@$(MAKE) db-seed
	@echo "Database reset complete"

.PHONY: db-reset-hard
db-reset-hard:
	@echo "Hard resetting database (drops all tables and migrations)..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) drop -f || true
	@$(MAKE) migrate-up
	@$(MAKE) db-seed
	@echo "Database hard reset complete"

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
	@echo "Starting API with hot reload (air)..."
	@air

.PHONY: run
run:
	@echo "Starting API (without hot reload)..."
	@go run cmd/api/main.go

# ==================== Web UI (Activation) ====================

.PHONY: web-up
web-up:
	@echo "Starting web UI..."
	@docker compose up -d social-web
	@echo "Web UI started at http://localhost:5173"

.PHONY: web-down
web-down:
	@echo "Stopping web UI..."
	@docker compose stop social-web
	@echo "Web UI stopped"

.PHONY: web-logs
web-logs:
	@docker compose logs -f social-web

.PHONY: web-rebuild
web-rebuild:
	@echo "Rebuilding web UI..."
	@docker compose build social-web
	@docker compose up -d social-web
	@echo "Web UI rebuilt and restarted"

# ==================== Swagger ====================
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -d ./cmd/api/,./internal/handlers,./internal/dtos,./internal/domain -o ./docs && swag fmt
	@echo "Swagger documentation generated successfully"

# ==================== Utilities ====================

.PHONY: help
help:
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "  Social API - Makefile Commands"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make setup          - First-time setup (Docker + DB + seed)"
	@echo "  make docker-up      - Start all Docker services"
	@echo "  make dev            - Run API with hot reload"
	@echo "  make docker-down    - Stop all Docker services"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-up      - Start all services (DB, Redis, Web)"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-restart - Restart all services"
	@echo "  make docker-logs    - View all service logs"
	@echo "  make docker-ps      - List running containers"
	@echo ""
	@echo "💾 Database:"
	@echo "  make db-start       - Start development database"
	@echo "  make db-stop        - Stop development database"
	@echo "  make db-seed        - Seed database with test data"
	@echo "  make db-flush       - Remove all data from database"
	@echo "  make db-reset       - Flush + migrate + seed"
	@echo "  make db-reset-hard  - Drop all tables + migrate + seed"
	@echo ""
	@echo "🔄 Migrations:"
	@echo "  make migration      - Create new migration file"
	@echo "  make migrate-up     - Apply all pending migrations"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test-setup     - Setup test database (run once)"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make test-coverage  - Generate coverage report"
	@echo ""
	@echo "💻 Development:"
	@echo "  make dev            - Run API with hot reload (air)"
	@echo "  make run            - Run API without hot reload"
	@echo ""
	@echo "🌐 Web UI (Activation):"
	@echo "  make web-up         - Start web UI (http://localhost:5173)"
	@echo "  make web-down       - Stop web UI"
	@echo "  make web-logs       - View web UI logs"
	@echo "  make web-rebuild    - Rebuild web UI container"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  make swagger        - Generate Swagger docs"
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
%:
	@: 
# This is for avoiding errors
