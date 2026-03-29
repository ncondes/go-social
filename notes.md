# Go Social App

## Set up Database Locally

```bash
docker compose up -d
```

## Run with hot reload locally

```bash
air
```

## Create migration

```bash
# -seq: create migration with sequential number
# -ext: specify the file extension for the migration files (.sql)
# -dir: specify the directory for the migration files (migrations)

migrate create -seq -ext sql -dir ./cmd/migrate/migrations <name_of_migration>
```

Example:

```bash
migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users_table
```

## Run migration

```bash
migrate -path ./cmd/migrate/migrations -database "postgres://<username>:<password>@localhost:5432/<database_name>?sslmode=disable" up
```

Example:

```bash
migrate -path ./cmd/migrate/migrations -database "postgres://postgres:password@localhost:5432/social?sslmode=disable" up
```

## With Makefile

### Up

```bash
make migrate-up
```

### Down

```bash
make migrate-down
```

## Run tests

Start test database

```bash
docker compose -f docker-compose.test.yml up -d
```

Run tests

```bash
make test
```
