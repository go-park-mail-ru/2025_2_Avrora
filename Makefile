.PHONY: test test-with-db lint

DB_NAME := 2025_2_Avrora_test
DB_USER := postgres
DB_PASS := postgres
DB_PORT := 5432
MIGRATIONS_DIR := ./infrastructure/db/migrations

TEST_DB_URL := postgres://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

lint:
	@echo "🔍 Running golangci-lint..."
	golangci-lint run

test: lint test-with-db

test-with-db:
	@echo "🚀 Starting PostgreSQL container..."
	docker run -d \
		--name test-postgres \
		-e POSTGRES_DB=postgres \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-p $(DB_PORT):5432 \
		--health-cmd="pg_isready -U $(DB_USER)" \
		--health-interval=1s \
		--health-timeout=5s \
		--health-retries=10 \
		postgres:15

	@echo "⏳ Waiting for DB to be ready..."
	@until docker exec test-postgres pg_isready -U $(DB_USER); do sleep 1; done

	@echo "🔧 Creating test database..."
	docker exec test-postgres psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"

	@echo "⬆️ Applying migrations using migrate CLI..."
	migrate -database "$(TEST_DB_URL)" -path "$(MIGRATIONS_DIR)" up

	@echo "🧪 Running Go tests..."
	TEST_DB_URL="$(TEST_DB_URL)" go test -v ./...

	@echo "🧹 Cleaning up..."
	docker stop test-postgres > /dev/null 2>&1
	docker rm test-postgres > /dev/null 2>&1