#!/bin/bash

# =====================================================
# SimNikah API - FULL CLEAN RESTRUCTURE
# =====================================================
# Complete reorganization following Go project layout
# https://github.com/golang-standards/project-layout
# =====================================================

set -e

echo "🏗️  Starting FULL CLEAN RESTRUCTURE..."
echo "⚠️  This will reorganize ALL code files"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

# =====================================================
# PHASE 1: Create Complete Folder Structure
# =====================================================
echo "📁 Phase 1: Creating complete folder structure..."

# Application entry points
mkdir -p cmd/api

# Private application code
mkdir -p internal/handlers/{auth,catin,staff,penghulu,kepala_kua,notification}
mkdir -p internal/models
mkdir -p internal/middleware
mkdir -p internal/services
mkdir -p internal/repository

# Public reusable packages
mkdir -p pkg/{validator,utils,crypto,cache,database}

# Configuration stays at root
mkdir -p config

echo "✅ Complete folder structure created!"
echo ""

# =====================================================
# PHASE 2: Move Application Entry Point
# =====================================================
echo "💻 Phase 2: Moving main.go to cmd/api/..."

if [ -f "main.go" ]; then
    mv main.go cmd/api/
    echo "✅ main.go moved to cmd/api/"
else
    echo "  main.go already moved"
fi

echo ""

# =====================================================
# PHASE 3: Move Handlers
# =====================================================
echo "🔀 Phase 3: Moving handlers to internal/handlers/..."

# Move handler directories
for dir in catin staff penghulu kepala_kua notification; do
    if [ -d "$dir" ]; then
        # Move contents, not the directory itself
        mkdir -p "internal/handlers/$dir"
        mv "$dir"/* "internal/handlers/$dir/" 2>/dev/null || true
        rmdir "$dir" 2>/dev/null || true
        echo "  ✅ $dir → internal/handlers/$dir"
    else
        echo "  Already moved: $dir"
    fi
done

echo "✅ Handlers moved!"
echo ""

# =====================================================
# PHASE 4: Move Models
# =====================================================
echo "📊 Phase 4: Moving models to internal/models/..."

if [ -d "structs" ]; then
    mkdir -p internal/models
    mv structs/* internal/models/ 2>/dev/null || true
    rmdir structs 2>/dev/null || true
    echo "✅ Models moved from structs/ to internal/models/"
else
    echo "  Already moved: structs"
fi

echo ""

# =====================================================
# PHASE 5: Move Middleware
# =====================================================
echo "🔧 Phase 5: Moving middleware to internal/middleware/..."

if [ -d "middleware" ]; then
    mkdir -p internal/middleware
    mv middleware/* internal/middleware/ 2>/dev/null || true
    rmdir middleware 2>/dev/null || true
    echo "✅ Middleware moved to internal/middleware/"
else
    echo "  Already moved: middleware"
fi

echo ""

# =====================================================
# PHASE 6: Move Services
# =====================================================
echo "⚙️  Phase 6: Moving services to internal/services/..."

if [ -d "services" ]; then
    mkdir -p internal/services
    mv services/* internal/services/ 2>/dev/null || true
    rmdir services 2>/dev/null || true
    echo "✅ Services moved to internal/services/"
else
    echo "  Already moved: services"
fi

echo ""

# =====================================================
# PHASE 7: Split Helper into PKG
# =====================================================
echo "📦 Phase 7: Splitting helper/ into pkg/..."

if [ -d "helper" ]; then
    mkdir -p pkg/{validator,utils,crypto,cache}
    
    # Move specific files to appropriate packages
    [ -f "helper/validation.go" ] && mv helper/validation.go pkg/validator/ && echo "  ✅ validation.go → pkg/validator/"
    [ -f "helper/marriage_validation.go" ] && mv helper/marriage_validation.go pkg/validator/ && echo "  ✅ marriage_validation.go → pkg/validator/"
    
    [ -f "helper/string_utils.go" ] && mv helper/string_utils.go pkg/utils/ && echo "  ✅ string_utils.go → pkg/utils/"
    [ -f "helper/date_utils.go" ] && mv helper/date_utils.go pkg/utils/ && echo "  ✅ date_utils.go → pkg/utils/"
    [ -f "helper/utils.go" ] && mv helper/utils.go pkg/utils/ && echo "  ✅ utils.go → pkg/utils/"
    
    [ -f "helper/bcrypt.go" ] && mv helper/bcrypt.go pkg/crypto/ && echo "  ✅ bcrypt.go → pkg/crypto/"
    
    [ -f "helper/geocoding_cache.go" ] && mv helper/geocoding_cache.go pkg/cache/ && echo "  ✅ geocoding_cache.go → pkg/cache/"
    
    # JWT stays in helper for now (depends on many things)
    [ -f "helper/jwt.go" ] && mv helper/jwt.go pkg/utils/ && echo "  ✅ jwt.go → pkg/utils/"
    
    # Remove empty helper directory
    rmdir helper 2>/dev/null && echo "✅ helper/ removed" || echo "  helper/ still has files, check manually"
else
    echo "  Already reorganized: helper"
fi

echo ""

# =====================================================
# PHASE 8: Update Makefile
# =====================================================
echo "📝 Phase 8: Updating Makefile for new structure..."

cat > Makefile << 'MAKEFILE_EOF'
.PHONY: build run test clean dev help

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
MAKEFILE_EOF

echo "✅ Makefile updated!"
echo ""

# =====================================================
# PHASE 9: Create .air.toml for hot reload
# =====================================================
echo "🔥 Phase 9: Creating .air.toml for hot reload..."

cat > .air.toml << 'AIR_EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main cmd/api/main.go"
  delay = 1000
  exclude_dir = ["bin", "tmp", "vendor", "testdata", "docs", "scripts", "tests", "deployments", "migrations"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
AIR_EOF

echo "✅ .air.toml created (for hot reload with 'air')"
echo ""

# =====================================================
# SUMMARY
# =====================================================
echo "=================================================="
echo "✅ FULL RESTRUCTURE COMPLETE!"
echo "=================================================="
echo ""
echo "📊 New structure:"
echo ""
echo "simpadu/"
echo "├── cmd/api/              # Application entry point"
echo "│   └── main.go"
echo "├── internal/             # Private application code"
echo "│   ├── handlers/         # HTTP handlers"
echo "│   │   ├── catin/"
echo "│   │   ├── staff/"
echo "│   │   ├── penghulu/"
echo "│   │   ├── kepala_kua/"
echo "│   │   └── notification/"
echo "│   ├── models/           # Database models"
echo "│   ├── middleware/       # Middleware"
echo "│   └── services/         # Business services"
echo "├── pkg/                  # Public reusable packages"
echo "│   ├── validator/"
echo "│   ├── utils/"
echo "│   ├── crypto/"
echo "│   └── cache/"
echo "├── config/               # Configuration"
echo "├── docs/                 # Documentation"
echo "├── scripts/              # Helper scripts"
echo "├── tests/                # Tests"
echo "├── deployments/          # Deployment configs"
echo "└── migrations/           # Database migrations"
echo ""
echo "⚠️  CRITICAL NEXT STEPS:"
echo "  1. Fix import paths: make fix-imports"
echo "  2. Review changes manually"
echo "  3. Test build: make build"
echo "  4. Test run: make dev"
echo ""
echo "📝 If build fails, check import paths in:"
echo "  - cmd/api/main.go"
echo "  - internal/handlers/**/*.go"
echo ""
STRUCTURE_EOF
chmod +x FULL_RESTRUCTURE.sh

