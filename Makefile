# Go REST API Makefile

# Variables
APP_NAME := go-rest-api
BINARY_NAME := $(APP_NAME)
DB_NAME := api.db

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GORUN := $(GOCMD) run

# Build flags
BUILD_FLAGS := -v
LDFLAGS := -w -s

.PHONY: help build run clean test deps tidy fmt vet lint
.PHONY: dev-setup docker-build docker-run

# Default target
all: build

## Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

build-release: ## Build optimized release binary
	@echo "Building release version of $(BINARY_NAME)..."
	CGO_ENABLED=1 GOOS=linux $(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY_NAME) .

## Development
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GORUN) main.go

## Testing
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

## Code Quality
fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

lint: ## Run golangci-lint (requires golangci-lint)
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOGET) -d ./...

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOGET) -u ./...

tidy: ## Clean up module dependencies
	@echo "Tidying module dependencies..."
	$(GOMOD) tidy

## Database
db-reset: ## Reset database (delete and recreate)
	@echo "Resetting database..."
	@if [ -f $(DB_NAME) ]; then \
		rm $(DB_NAME); \
		echo "Database $(DB_NAME) removed"; \
	fi

db-backup: ## Backup database
	@if [ -f $(DB_NAME) ]; then \
		cp $(DB_NAME) $(DB_NAME).backup.$$(date +%Y%m%d_%H%M%S); \
		echo "Database backed up"; \
	else \
		echo "Database $(DB_NAME) not found"; \
	fi

## Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@if [ -f $(BINARY_NAME) ]; then rm $(BINARY_NAME); fi
	@if [ -f coverage.out ]; then rm coverage.out; fi
	@if [ -f coverage.html ]; then rm coverage.html; fi

clean-all: clean ## Clean everything including database
	@if [ -f $(DB_NAME) ]; then \
		rm $(DB_NAME); \
		echo "Database removed"; \
	fi

## Docker (optional)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run: ## Run application in Docker container
	@echo "Running $(APP_NAME) in Docker..."
	docker run -p 8080:8080 --name $(APP_NAME) $(APP_NAME):latest

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

## Quick commands for common workflows
check: fmt vet lint test ## Run all code quality checks

release: clean test build-release ## Build release version

# Development workflow
dev-workflow: clean deps tidy fmt vet test build ## Complete development workflow

# Show available make targets
list: ## List all available targets
	@$(MAKE) help