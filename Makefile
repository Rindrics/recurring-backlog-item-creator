.PHONY: build test clean install help

# Variable definitions
BINARY_NAME=gh-issue-config-filter
MAIN_PACKAGE=.
BUILD_DIR=bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
.DEFAULT_GOAL := help

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## test-coverage: Generate test coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## install: Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(MAIN_PACKAGE)
	@echo "Install complete"

## lint: Run linters
lint:
	@echo "Running linters..."
	@go vet ./...
	@golangci-lint run ./... || echo "golangci-lint not installed, skipping..."

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## help: Show this help message
help: Makefile
	@echo "Available targets:"
	@sed -n 's/^##//p' ${<} | column -t -s ':' | sed -e 's/^/ /'

