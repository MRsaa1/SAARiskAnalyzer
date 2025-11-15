.PHONY: help up down build rebuild logs clean test lint seed migrate

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

build: ## Build all services
	docker-compose build

rebuild: ## Rebuild and restart all services
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

logs: ## Show logs (use: make logs SERVICE=api)
	docker-compose logs -f $(SERVICE)

logs-api: ## Show API logs
	docker-compose logs -f api

logs-web: ## Show frontend logs
	docker-compose logs -f web

clean: ## Clean up containers, volumes, and build artifacts
	docker-compose down -v
	rm -rf frontend/dist frontend/node_modules
	rm -rf backend/tmp

test: ## Run all tests
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd frontend && npm test

test-backend: ## Run backend tests only
	cd backend && go test -v -race -coverprofile=coverage.out ./...

test-frontend: ## Run frontend tests only
	cd frontend && npm test

lint: ## Run linters
	cd backend && golangci-lint run
	cd frontend && npm run lint

lint-fix: ## Fix linting issues
	cd backend && golangci-lint run --fix
	cd frontend && npm run lint:fix

seed: ## Seed database with demo data
	./scripts/seed.sh

migrate: ## Run database migrations
	./scripts/migrate.sh

dev-backend: ## Run backend in development mode
	cd backend && go run cmd/api/main.go

dev-frontend: ## Run frontend in development mode
	cd frontend && npm run dev

install-backend: ## Install backend dependencies
	cd backend && go mod download

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

setup: install-backend install-frontend ## Install all dependencies

ps: ## Show running containers
	docker-compose ps

restart: ## Restart all services
	docker-compose restart

restart-api: ## Restart API service
	docker-compose restart api

restart-web: ## Restart web service
	docker-compose restart web
