.PHONY: help install-tools migrate-up migrate-down migrate-create sqlc-generate test test-verbose run docker-up docker-down clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## Install development tools
	@echo "Installing sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing mockgen..."
	@go install go.uber.org/mock/mockgen@latest
	@echo "Tools installed successfully!"

sqlc-generate: ## Generate SQLC code
	@echo "Generating SQLC code..."
	@sqlc generate

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@migrate create -ext sql -dir migrations -seq $(name)

migrate-test-up: ## Test database migrations up
	@migrate -path migrations -database "postgresql://aiki_test:aiki_test_password@localhost:5433/aiki_test_db?sslmode=disable" -verbose up

migrate-up: ## Run database migrations up
	@migrate -path migrations -database "postgresql://aiki:aiki_password@localhost:5432/aiki_db?sslmode=disable" -verbose up

migrate-cleanup: ## Clean up failed migrations
	@migrate -path migrations -database "postgresql://aiki:aiki_password@localhost:5432/aiki_db?sslmode=disable" force 1

migrate-down: ## Run database migrations down
	@migrate -path migrations -database "postgresql://aiki:aiki_password@localhost:5432/aiki_db?sslmode=disable" -verbose down

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@migrate -path migrations -database "postgresql://aiki:aiki_password@localhost:5432/aiki_db?sslmode=disable" force $(version)

test: ## Run tests
	@go test -v -race -cover ./...

test-verbose: ## Run tests with verbose output
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

run: ## Run the application
	@go run cmd/api/main.go

docker-up: ## Start docker containers
	@docker compose up -d

docker-down: ## Stop docker containers
	@docker compose down

docker-logs: ## Show docker logs
	@docker compose logs -f

clean: ## Clean build artifacts
	@rm -f coverage.out coverage.html
	@go clean

build: ## Build the application
	@go build -o bin/aiki cmd/api/main.go

tidy: ## Tidy go modules
	@go mod tidy
