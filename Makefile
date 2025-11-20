include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

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

%:
	@: 
# This is for avoiding errors
