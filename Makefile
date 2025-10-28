.PHONY: build run test clean dev help

# Build the application
build:
	@echo "🔨 Building SimNikah API..."
	@go build -o bin/simnikah-api .
	@echo "✅ Build complete! Binary: bin/simnikah-api"

# Run the application
run:
	@echo "🚀 Starting SimNikah API..."
	@go run main.go

# Run in development mode
dev:
	@echo "🔧 Starting in development mode..."
	@GIN_MODE=debug go run main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

# Run linter (if available)
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run 2>/dev/null || go vet ./...
	@echo "✅ Lint complete!"

# Show help
help:
	@echo ""
	@echo "SimNikah API - Available Commands:"
	@echo ""
	@echo "  make build     - Build the application binary"
	@echo "  make run       - Run the application (production mode)"
	@echo "  make dev       - Run in development mode with debug output"
	@echo "  make test      - Run all tests"
	@echo "  make coverage  - Run tests with coverage report"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deps      - Install/update dependencies"
	@echo "  make fmt       - Format all Go code"
	@echo "  make lint      - Run linter/vet checks"
	@echo "  make help      - Show this help message"
	@echo ""
