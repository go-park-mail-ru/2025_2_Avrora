.PHONY: test test-with-db lint

DB_NAME := 2025_2_Avrora
DB_USER := postgres
DB_PASS := postgres
DB_PORT := 5432

DB_URL := postgres://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Path to migrations directory
MIGRATIONS_DIR := internal/db/migrations

MIGRATE := migrate -database $(DB_URL) -path $(MIGRATIONS_DIR)

.PHONY: migrate-up
migrate-up:
	@echo "🚀 Applying database migrations..."
	@$(MIGRATE) up
	@echo "✅ Migrations applied successfully."

.PHONY: migrate-down
migrate-down:
	@echo "⏳ Rolling back last migration..."
	@$(MIGRATE) down
	@echo "✅ Last migration rolled back."

.PHONY: migrate-force-clean
migrate-force-clean:
	@echo "⚠️ Forcing clean state (version 0) — use only in dev!"
	@$(MIGRATE) force 0
	@echo "✅ Database reset to version 0."

.PHONY: migrate-status
migrate-status:
	@echo "📊 Migration status:"
	@$(MIGRATE) version

run:
	@echo "🚀 Starting server..."
	go run ./cmd/app/main.go &
	go run ./cmd/auth/main.go &
	go run ./cmd/fileserver/main.go &

PORTS := 8080 50051 50052

.PHONY: killports clean

# Kill processes on specified ports
killports:
	@echo "🔍 Killing processes on ports: $(PORTS)"
	@for port in $(PORTS); do \
		echo "➡️  Checking port $$port..."; \
		pids=$$(lsof -ti:$$port 2>/dev/null); \
		if [ -n "$$pids" ]; then \
			echo "   🚫 Killing PID(s): $$pids"; \
			kill -9 $$pids; \
		else \
			echo "   ✅ No process found on port $$port"; \
		fi; \
	done
	@echo "✅ Done."

# Alias for convenience
clean: killports

build_proto:
	@echo "🔧 Generating proto files..."
	find ./proto -type f -name "*.proto" -exec protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative {} +

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