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
