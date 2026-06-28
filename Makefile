.PHONY: build run test test-report clean dev help

# Application name
APP_NAME=simnikah-api
BUILD_DIR=bin
CMD_DIR=cmd/api

# Build the application
build:
	@echo "🔨 Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "✅ Build complete! Binary: $(BUILD_DIR)/$(APP_NAME)"

# Run the application (production mode)
run:
	@echo "🚀 Starting $(APP_NAME) (production mode)..."
	@go run $(CMD_DIR)/main.go

# Run in development mode with live reload
dev:
	@echo "🔧 Starting $(APP_NAME) (development mode)..."
	@GIN_MODE=debug go run $(CMD_DIR)/main.go

# Run with hot reload (requires air)
watch:
	@echo "👀 Starting with hot reload..."
	@air || echo "Install 'air' first: go install github.com/cosmtrek/air@latest"

# Run all tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Generate HTML test report
test-report:
	@echo "📋 Generating test report..."
	@go run ./cmd/testreport/
	@echo "✅ Test report: test_report.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

# Install/update dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies up to date!"

# Format all Go code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run 2>/dev/null || go vet ./...
	@echo "✅ Lint complete!"

# Update import paths after restructure
fix-imports:
	@echo "🔧 Fixing import paths..."
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/catin|simnikah/internal/handlers/catin|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/staff|simnikah/internal/handlers/staff|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/penghulu|simnikah/internal/handlers/penghulu|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/kepala_kua|simnikah/internal/handlers/kepala_kua|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/notification|simnikah/internal/handlers/notification|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/structs|simnikah/internal/models|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/middleware|simnikah/internal/middleware|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/services|simnikah/internal/services|g' {} +
	@find . -name "*.go" -type f -exec sed -i 's|simnikah/helper|simnikah/pkg/utils|g' {} +
	@echo "✅ Import paths updated!"
	@echo "⚠️  Please review and test!"

# Show project structure
tree:
	@echo "📁 Project structure:"
	@tree -L 3 -I 'bin|go|tests|vendor|node_modules' || ls -R

# Display help
help:
	@echo ""
	@echo "SimNikah API - Available Commands"
	@echo "=================================="
	@echo ""
	@echo "Building & Running:"
	@echo "  make build        - Build production binary"
	@echo "  make run          - Run in production mode"
	@echo "  make dev          - Run in development mode"
	@echo "  make watch        - Run with hot reload (requires air)"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make coverage     - Run tests with coverage report"
	@echo "  make test-report  - Generate HTML test report"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt          - Format all Go code"
	@echo "  make lint         - Run linter checks"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Install/update dependencies"
	@echo "  make fix-imports  - Fix import paths after restructure"
	@echo "  make tree         - Show project structure"
	@echo ""
	@echo "Help:"
	@echo "  make help         - Show this help message"
	@echo ""
