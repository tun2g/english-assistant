# Variables
APP_NAME := app-backend
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_COMPOSE_DEV_FILE := docker-compose.dev.yml
MIGRATION_DIR := migrations
DATABASE_URL := postgres://postgres:postgres@localhost:5434/app_backend_dev?sslmode=disable

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help build run dev test clean docker-up docker-down docker-logs migrate-up migrate-down migrate-force migrate-version swagger lint format deps install-tools

# Default target
help: ## Show this help message
	@echo "$(BLUE)Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build commands
build: ## Build the application
	@echo "$(BLUE)Building $(APP_NAME)...$(NC)"
	go build -o bin/$(APP_NAME) ./cmd/server

build-linux: ## Build for Linux
	@echo "$(BLUE)Building $(APP_NAME) for Linux...$(NC)"
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux ./cmd/server

# Run commands
run: ## Run the application
	@echo "$(BLUE)Running $(APP_NAME)...$(NC)"
	go run ./cmd/server

dev: ## Run in development mode with live reload (requires air)
	@echo "$(BLUE)Starting development server...$(NC)"
	air

# Test commands
test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Docker commands
docker-up: ## Start development environment with Docker
	@echo "$(BLUE)Starting development environment...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_DEV_FILE) up -d
	@echo "$(GREEN)Development environment started$(NC)"

docker-down: ## Stop development environment
	@echo "$(BLUE)Stopping development environment...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_DEV_FILE) down
	@echo "$(GREEN)Development environment stopped$(NC)"

docker-logs: ## Show Docker logs
	docker-compose -f $(DOCKER_COMPOSE_DEV_FILE) logs -f

docker-prod-up: ## Start production environment
	@echo "$(BLUE)Starting production environment...$(NC)"
	cd ../../devops/docker && docker-compose -f docker-compose.prod.yml up -d
	@echo "$(GREEN)Production environment started$(NC)"

docker-prod-down: ## Stop production environment
	@echo "$(BLUE)Stopping production environment...$(NC)"
	cd ../../devops/docker && docker-compose -f docker-compose.prod.yml down
	@echo "$(GREEN)Production environment stopped$(NC)"

# Deployment commands
deploy-prod: ## Deploy to production using DevOps scripts
	@echo "$(BLUE)Deploying to production...$(NC)"
	cd ../../devops/scripts && ./deploy.sh production
	@echo "$(GREEN)Production deployment completed$(NC)"

# Database migration commands
migrate-up: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Database migrations completed$(NC)"

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" down 1

migrate-force: ## Force migration version (use with VERSION=N)
	@echo "$(YELLOW)Forcing migration version $(VERSION)...$(NC)"
	migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" force $(VERSION)

migrate-version: ## Show current migration version
	@echo "$(BLUE)Current migration version:$(NC)"
	migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" version

migrate-create: ## Create new migration (use with NAME=migration_name)
	@echo "$(BLUE)Creating migration: $(NAME)$(NC)"
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

# Code quality commands
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	golangci-lint run

format: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	goimports -w .

# Swagger commands
swagger: ## Generate Swagger documentation
	@echo "$(BLUE)Generating Swagger documentation...$(NC)"
	swag init -g cmd/server/main.go -o docs/swagger
	@echo "$(GREEN)Swagger documentation generated$(NC)"

# Dependency commands
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy

deps-update: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	go get -u ./...
	go mod tidy

# Installation commands
install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)Development tools installed$(NC)"

# Clean commands
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache
	@echo "$(GREEN)Cleaned$(NC)"

# Environment setup
setup: install-tools deps ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	cp configs/app.yaml.example configs/app.yaml 2>/dev/null || true
	@echo "$(GREEN)Development environment setup complete$(NC)"
	@echo "$(YELLOW)Don't forget to:$(NC)"
	@echo "  1. Update configs/app.yaml with your settings"
	@echo "  2. Start the database: make docker-up"
	@echo "  3. Run migrations: make migrate-up"

# Database commands
db-reset: ## Reset database (drop and recreate)
	@echo "$(YELLOW)Resetting database...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_DEV_FILE) down -v
	docker-compose -f $(DOCKER_COMPOSE_DEV_FILE) up -d postgres
	sleep 5
	$(MAKE) migrate-up
	@echo "$(GREEN)Database reset completed$(NC)"

# Security scan
security: ## Run security scan
	@echo "$(BLUE)Running security scan...$(NC)"
	gosec ./...

# Performance profiling
profile: ## Run with profiling
	@echo "$(BLUE)Running with profiling...$(NC)"
	go run -race ./cmd/server