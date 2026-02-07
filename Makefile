.PHONY: help install build test clean dev docker

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Installing backend dependencies..."
	cd backend && go mod download

build: ## Build frontend and backend
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Building backend..."
	cd backend && go build -o server cmd/server/main.go

test: ## Run all tests
	@echo "Running frontend tests..."
	cd frontend && npm run test
	@echo "Running backend tests..."
	cd backend && go test ./...

test-coverage: ## Run tests with coverage
	@echo "Running frontend tests with coverage..."
	cd frontend && npm run test:coverage
	@echo "Running backend tests with coverage..."
	cd backend && go test -cover ./...

lint: ## Run linters
	@echo "Linting frontend..."
	cd frontend && npm run lint
	@echo "Linting backend..."
	cd backend && golangci-lint run

format: ## Format code
	@echo "Formatting frontend..."
	cd frontend && npm run format
	@echo "Formatting backend..."
	cd backend && go fmt ./...

dev-frontend: ## Run frontend in development mode
	cd frontend && npm run dev

dev-backend: ## Run backend in development mode
	cd backend && go run cmd/server/main.go

dev: ## Run both frontend and backend in development mode
	@echo "Starting development servers..."
	@make -j2 dev-frontend dev-backend

docker-build: ## Build Docker image
	docker build -t flow:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 flow:latest

docker-compose-up: ## Start services with docker-compose
	docker-compose up --build

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/server
	rm -rf backend/coverage.out

e2e: ## Run E2E tests
	cd frontend && npm run e2e

security-scan: ## Run security scans
	@echo "Running frontend security scan..."
	cd frontend && npm audit
	@echo "Running backend security scan..."
	cd backend && gosec ./...
