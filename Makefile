.PHONY: build run test lint migrate

# Build the application
build:
	go build -o loan-service ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Run migrations
migrate:
	psql -U mungkiice -d loan_service -f migrations/001_create_schema.up.sql

migrate-down:
	psql -U mungkiice -d loan_service -f migrations/001_create_schema.down.sql

# Clean build artifacts
clean:
	rm -f loan-service coverage.out coverage.html
