.PHONY: build install test clean run

# Binary name
BINARY_NAME=tasklog

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Build and install to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete. Run 'tasklog' to use."

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	@echo "\nOverall coverage:"
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}'
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "\nCore package coverage:"
	@go tool cover -func=coverage.out | grep -E "(config|storage|timeparse)" | grep -v "\.go:" || true

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Run the application (development)
run: build
	./$(BINARY_NAME)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe
	@echo "Build complete for all platforms"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Setup development environment
setup:
	@echo "Setting up development environment..."
	go mod download
	@echo "Creating example config directory..."
	mkdir -p ~/.tasklog
	@if [ ! -f ~/.tasklog/config.yaml ]; then \
		cp config.example.yaml ~/.tasklog/config.yaml; \
		echo "Config file created at ~/.tasklog/config.yaml"; \
		echo "Please edit it with your credentials."; \
	else \
		echo "Config file already exists at ~/.tasklog/config.yaml"; \
	fi
	@echo "Setup complete"

# Show help
help:
	@echo "Tasklog Makefile commands:"
	@echo "  make build         - Build the application"
	@echo "  make install       - Build and install to /usr/local/bin"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make run           - Build and run the application"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make setup         - Setup development environment"
	@echo "  make help          - Show this help message"
