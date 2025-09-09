# VIES Query - Makefile

.PHONY: build test clean lint fmt install help run

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY_NAME = viesquery
BUILD_DIR = bin
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -s -w"
PACKAGE = l22.io/viesquery

# Go environment
GO ?= go
GOFMT ?= gofmt
GOLINT ?= golangci-lint

# Default target
all: build

## help: Display this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary
build: clean
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/viesquery

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	$(GOFMT) -s -w .

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)/
	@rm -f coverage.out coverage.html

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

## run: Build and run with example
run: build
	@echo "Running example validation..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --help

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying modules..."
	$(GO) mod tidy

## mod-verify: Verify go modules
mod-verify:
	@echo "Verifying modules..."
	$(GO) mod verify

## release: Build binaries for multiple platforms
release: clean test lint
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/viesquery
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/viesquery
	# macOS
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/viesquery
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/viesquery
	# Windows
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/viesquery
	@echo "Release binaries built in $(BUILD_DIR)/"

## check: Run all checks (format, lint, test)
check: fmt lint test
	@echo "All checks passed!"

## dev-setup: Set up development environment
dev-setup:
	@echo "Setting up development environment..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment ready!"

# Show build info
info:
	@echo "Build Information:"
	@echo "  Version: $(VERSION)"
	@echo "  Binary:  $(BINARY_NAME)"
	@echo "  Package: $(PACKAGE)"
	@echo "  Go:      $(shell $(GO) version)"
