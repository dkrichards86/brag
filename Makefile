.PHONY: all build test test-verbose test-coverage lint fmt clean install help

# Build variables
BINARY_NAME=brag
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Colors for output
BLUE=\033[0;34m
NC=\033[0m # No Color

all: test lint build ## Run tests, lint, and build

build: ## Build the binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@go build ${LDFLAGS} -o ${BINARY_NAME} .
	@echo "$(BLUE)Build complete: $(BINARY_NAME)$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	@go test -v -race ./...

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(BLUE)Coverage report generated: coverage.html$(NC)"

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	elif [ -f "$(shell go env GOPATH)/bin/golangci-lint" ]; then \
		$(shell go env GOPATH)/bin/golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@go mod tidy

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -f ${BINARY_NAME}
	@rm -f coverage.out coverage.html
	@rm -f brag-*-*
	@go clean

install: ## Install the binary to GOPATH/bin
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	@go install ${LDFLAGS}
	@echo "$(BLUE)Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod verify

build-all: ## Build for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o brag-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o brag-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o brag-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o brag-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o brag-windows-amd64.exe .
	@echo "$(BLUE)Cross-platform builds complete$(NC)"

ci: deps test lint build ## Run CI pipeline locally

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2}'
