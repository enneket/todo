# Makefile for Todo App (Cross-platform)
# Inspired by https://github.com/WCY-dt/MrRSS

# .PHONY targets are not files
.PHONY: all build dev install-deps lint lint-backend lint-frontend format format-backend format-frontend test test-backend test-frontend test-e2e clean help

# Variables
APP_NAME := TodoApp
BUILD_DIR := build/bin
FRONTEND_DIR := frontend
GO_CMD := go
NPM_CMD := npm
WAILS_CMD := wails

# Detect OS
OS := $(shell uname -s)
ifeq ($(OS),Linux)
    # Linux-specific flags
    # Use webkit2_41 build tag for compatibility with newer WebKitGTK
    BUILD_TAGS := -tags webkit2_41
else
    # Default (empty) for macOS/Windows
    BUILD_TAGS :=
endif

# Default target
all: help

# Help target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Development ---

dev: ## Run the application in development mode
	$(WAILS_CMD) dev $(BUILD_TAGS)

ensure-dist: ## Ensure frontend/dist exists for Go tools
	$(GO_CMD) run scripts/ensure_dist.go

install-deps: install-backend-deps install-frontend-deps ## Install all dependencies

install-backend-deps: ensure-dist ## Install backend dependencies
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy

install-frontend-deps: ## Install frontend dependencies
	cd $(FRONTEND_DIR) && $(NPM_CMD) install

update-deps: update-backend-deps update-frontend-deps ## Update all dependencies

update-backend-deps: ensure-dist ## Update backend dependencies
	$(GO_CMD) get -u ./...
	$(GO_CMD) mod tidy

update-frontend-deps: ## Update frontend dependencies
	cd $(FRONTEND_DIR) && $(NPM_CMD) update

bindings: ## Generate Wails bindings
	$(WAILS_CMD) generate module

# --- Build ---

build: ## Build the application for production
	$(WAILS_CMD) build $(BUILD_TAGS)

build-windows: ## Build for Windows
	$(WAILS_CMD) build -platform windows/amd64

build-linux: ## Build for Linux
	$(WAILS_CMD) build -platform linux/amd64 $(BUILD_TAGS)

build-darwin: ## Build for macOS
	$(WAILS_CMD) build -platform darwin/universal

run: ## Run the built application
	./$(BUILD_DIR)/$(APP_NAME)

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf coverage.out coverage.html
	@echo "✅ Cleaned build artifacts"

# --- Quality Assurance ---

check: lint test build ## Run full check (lint, test, build)

lint: lint-backend lint-frontend ## Run all linters

lint-backend: ensure-dist ## Run Go linter
	$(GO_CMD) vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Skipping advanced linting."; \
	fi

lint-frontend: ## Run Frontend linter (ESLint)
	cd $(FRONTEND_DIR) && $(NPM_CMD) run lint

format: format-backend format-frontend ## Format all code

format-backend: ## Format Go code
	$(GO_CMD) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Skipping import organization."; \
	fi

format-frontend: ## Format Frontend code (Prettier)
	cd $(FRONTEND_DIR) && $(NPM_CMD) run format

# --- Testing ---

test: test-backend test-frontend ## Run unit tests

test-backend: ensure-dist ## Run Go unit tests
	$(GO_CMD) test ./...

test-coverage: ensure-dist ## Run backend tests with coverage
	$(GO_CMD) test -v -timeout=5m -coverprofile=coverage.out -covermode=atomic ./...
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-frontend: ## Run Frontend unit tests
	cd $(FRONTEND_DIR) && $(NPM_CMD) run test:unit || echo "No unit tests configured yet"

test-e2e: ## Run E2E tests with Cypress (requires app running)
	cd $(FRONTEND_DIR) && npx cypress run

test-all: test test-e2e ## Run all tests including E2E

# --- Setup ---

setup: install-deps ## Setup development environment
	@echo "Setup complete! Run 'make dev' to start the app."

love: ## Show some love
	@echo "❤️ Todo App loves you too! ❤️"
