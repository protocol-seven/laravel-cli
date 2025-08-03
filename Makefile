BINARY_NAME=laravel
VERSION?=dev
BUILD_DIR=dist

.PHONY: build clean test install uninstall cross-compile help

# Default target
build: ## Build the binary
	go build -ldflags="-s -w" -o $(BINARY_NAME)

# Clean build artifacts
clean: ## Clean build artifacts
	rm -rf $(BINARY_NAME) $(BUILD_DIR)/ *.deb *.rpm

# Run tests
test: ## Run tests
	go test -v ./...

# Run tests in the tests directory only
test-matrix: ## Run the comprehensive test matrix
	go test -v ./tests/...

# Run specific test categories
test-flags: ## Test command-line flags
	go test -v ./tests/ -run TestCommandFlags

test-validation: ## Test input validation
	go test -v ./tests/ -run TestProjectValidation

test-database: ## Test database functionality
	go test -v ./tests/ -run TestDatabaseConfiguration

test-starterkits: ## Test starter kits
	go test -v ./tests/ -run TestStarterKits

test-git: ## Test Git operations
	go test -v ./tests/ -run TestGitOperations

test-env: ## Test environment file operations
	go test -v ./tests/ -run TestEnvironmentFileOperations

test-interactive: ## Test interactive features
	go test -v ./tests/ -run "Test.*Port.*|Test.*Input.*"

test-utils: ## Test utility functions
	go test -v ./tests/ -run TestUtilityFunctions

test-integration: ## Test integration scenarios
	go test -v ./tests/ -run TestIntegrationScenarios

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	go test -v -cover ./...
	go test -v -cover ./tests/...

# Run only unit tests (excluding integration tests)
test-unit: ## Run unit tests only
	go test -v ./tests/ -run "Test.*" -skip "Integration"

# Run only integration tests
test-integration-only: ## Run integration tests only
	go test -v ./tests/ -run "TestIntegration.*"

# Run original main tests
test-main: ## Run tests in main package
	go test -v .

# Install locally
install: build ## Install the binary to /usr/local/bin
	sudo cp $(BINARY_NAME) /usr/local/bin/

# Uninstall
uninstall: ## Remove the binary from /usr/local/bin
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Cross-compile for all platforms
cross-compile: ## Build for all platforms
	mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64

# Initialize go module
mod-init: ## Initialize go module and download dependencies
	go mod tidy

# Format code
fmt: ## Format Go code
	go fmt ./...

# Lint code (requires golangci-lint)
lint: ## Lint Go code
	golangci-lint run

# Check for security issues (requires gosec)
security: ## Run security checks
	gosec ./...

# Create a release archive
release: cross-compile ## Create release archives
	cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-*; do \
		if [[ $$file == *.exe ]]; then \
			zip $${file%.*}.zip $$file; \
		else \
			tar -czf $$file.tar.gz $$file; \
		fi \
	done

# Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
