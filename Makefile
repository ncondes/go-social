include .env
MIGRATIONS_PATH = ./cmd/db/migrations

.PHONY: migrate-create
migration:
	@echo "Creating migration file..."
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
	@echo "Migration file created successfully"

.PHONY: migrate-up
migrate-up:
	@echo "Applying migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) up
	@echo "Migrations applied successfully"

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) down
	@echo "Migrations rolled back successfully"

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

%:
	@: 
# This is for avoiding errors
