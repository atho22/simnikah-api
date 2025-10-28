#!/bin/bash

# =====================================================
# SimNikah API - LIGHT Folder Restructure
# =====================================================
# Hanya organisasi dokumentasi & scripts
# TIDAK memindahkan kode Go (aman untuk testing)
# =====================================================

set -e

echo "🏗️  Starting LIGHT folder restructure (docs & scripts only)..."
echo ""

# Create folders
echo "📁 Creating folder structure..."
mkdir -p docs/{deployment,features,api,performance,tutorials,architecture}
mkdir -p scripts
mkdir -p tests/fixtures
mkdir -p deployments/{railway,docker}
mkdir -p migrations
mkdir -p bin
echo "✅ Folders created!"
echo ""

# Move documentation
echo "📚 Organizing documentation..."
mv RAILWAY_DEPLOYMENT.md docs/deployment/ 2>/dev/null || echo "  Already moved: RAILWAY_DEPLOYMENT.md"
mv MIGRATION_GUIDE.md docs/deployment/ 2>/dev/null || echo "  Already moved: MIGRATION_GUIDE.md"
mv PENJELASAN_MUDAH_RAILWAY.md docs/deployment/ 2>/dev/null || echo "  Already moved: PENJELASAN_MUDAH_RAILWAY.md"
mv TUTORIAL_DEPLOY_RAILWAY.md docs/tutorials/ 2>/dev/null || echo "  Already moved: TUTORIAL_DEPLOY_RAILWAY.md"
mv QUICK_START.md docs/deployment/ 2>/dev/null || echo "  Already moved: QUICK_START.md"
mv RAILWAY_ENV_TEMPLATE.txt docs/deployment/ 2>/dev/null || echo "  Already moved: RAILWAY_ENV_TEMPLATE.txt"
mv LEAPCELL_CORS_GUIDE.md docs/deployment/ 2>/dev/null || echo "  Already moved: LEAPCELL_CORS_GUIDE.md"
mv DEPLOYMENT.md docs/deployment/ 2>/dev/null || echo "  Already moved: DEPLOYMENT.md"

mv MAP_INTEGRATION.md docs/features/ 2>/dev/null || echo "  Already moved: MAP_INTEGRATION.md"
mv WALI_NIKAH_VALIDATION.md docs/features/ 2>/dev/null || echo "  Already moved: WALI_NIKAH_VALIDATION.md"
mv STATUS_MANAGEMENT.md docs/features/ 2>/dev/null || echo "  Already moved: STATUS_MANAGEMENT.md"
mv CORS_SETUP.md docs/features/ 2>/dev/null || echo "  Already moved: CORS_SETUP.md"

mv API_DOCUMENTATION.md docs/api/ 2>/dev/null || echo "  Already moved: API_DOCUMENTATION.md"
mv API_TESTING_DOCUMENTATION.md docs/api/ 2>/dev/null || echo "  Already moved: API_TESTING_DOCUMENTATION.md"

mv PERFORMANCE_ASSESSMENT.md docs/performance/ 2>/dev/null || echo "  Already moved: PERFORMANCE_ASSESSMENT.md"
mv OPTIMIZATION_SUMMARY.md docs/performance/ 2>/dev/null || echo "  Already moved: OPTIMIZATION_SUMMARY.md"

mv FOLDER_STRUCTURE_PLAN.md docs/architecture/ 2>/dev/null || echo "  Already moved: FOLDER_STRUCTURE_PLAN.md"

echo "✅ Documentation organized!"
echo ""

# Move scripts
echo "🔧 Organizing scripts..."
mv test_autocomplete.sh scripts/ 2>/dev/null || echo "  Already moved: test_autocomplete.sh"
mv test_railway_endpoint.sh scripts/ 2>/dev/null || echo "  Already moved: test_railway_endpoint.sh"
mv RESTRUCTURE_COMMANDS.sh scripts/ 2>/dev/null || echo "  Already moved: RESTRUCTURE_COMMANDS.sh"
chmod +x scripts/*.sh 2>/dev/null || true
echo "✅ Scripts organized!"
echo ""

# Move test fixtures
echo "🧪 Organizing test files..."
mv demo_autocomplete.html tests/fixtures/ 2>/dev/null || echo "  Already moved: demo_autocomplete.html"
echo "✅ Test files organized!"
echo ""

# Move deployment configs
echo "🚀 Organizing deployment configs..."
mv railway.json deployments/railway/ 2>/dev/null || echo "  Already moved: railway.json"
mv nixpacks.toml deployments/railway/ 2>/dev/null || echo "  Already moved: nixpacks.toml"
mv Dockerfile deployments/docker/ 2>/dev/null || echo "  Already moved: Dockerfile"
mv docker-compose.yml deployments/docker/ 2>/dev/null || echo "  Already moved: docker-compose.yml"
mv init.sql migrations/ 2>/dev/null || echo "  Already moved: init.sql"
echo "✅ Deployment configs organized!"
echo ""

# Create Makefile
echo "📝 Creating Makefile..."
cat > Makefile << 'MAKEFILE_EOF'
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
MAKEFILE_EOF
echo "✅ Makefile created!"
echo ""

# Update .gitignore
echo "🔒 Updating .gitignore..."
if ! grep -q "^bin/" .gitignore 2>/dev/null; then
cat >> .gitignore << 'GITIGNORE_EOF'

# === Added by restructure ===
# Build artifacts
bin/
*.exe
*.exe~

# Test artifacts
*.test
*.out
coverage.html

# IDE
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
GITIGNORE_EOF
echo "✅ .gitignore updated!"
else
echo "  .gitignore already up-to-date"
fi
echo ""

# Create structure documentation
cat > docs/architecture/FOLDER_STRUCTURE.md << 'STRUCTURE_EOF'
# 📁 SimNikah API - Current Folder Structure

**Last Updated:** $(date +"%Y-%m-%d")

## Root Structure

```
simpadu/
├── main.go                    # Application entry point
├── Makefile                   # Build automation
├── go.mod, go.sum             # Go dependencies
├── env.example                # Environment template
├── README.md                  # Main documentation
│
├── catin/                     # Marriage registration handlers
├── staff/                     # Staff management handlers
├── penghulu/                  # Penghulu (officiant) handlers
├── kepala_kua/                # KUA head handlers
├── notification/              # Notification handlers
│
├── structs/                   # Data models & constants
├── helper/                    # Helper utilities
├── middleware/                # HTTP middleware
├── services/                  # Business services
├── config/                    # Configuration & database
│
├── docs/                      # 📚 All documentation
│   ├── deployment/            # Deployment guides
│   ├── features/              # Feature documentation
│   ├── api/                   # API documentation
│   ├── performance/           # Performance reports
│   ├── tutorials/             # Tutorials & guides
│   └── architecture/          # Architecture docs
│
├── scripts/                   # 🔧 Helper scripts
├── tests/                     # 🧪 Tests & fixtures
├── deployments/               # 🚀 Deployment configs
│   ├── railway/               # Railway specific
│   └── docker/                # Docker configs
│
├── migrations/                # 💾 Database migrations
└── bin/                       # 📦 Build output (gitignored)
```

## Documentation Organization

All documentation files are now organized in `docs/`:

### Deployment (`docs/deployment/`)
- `RAILWAY_DEPLOYMENT.md` - Railway deployment guide
- `QUICK_START.md` - Quick start guide
- `MIGRATION_GUIDE.md` - Migration from LeapCell
- `PENJELASAN_MUDAH_RAILWAY.md` - Easy Railway tutorial (Indonesian)
- `RAILWAY_ENV_TEMPLATE.txt` - Environment variable template

### Features (`docs/features/`)
- `MAP_INTEGRATION.md` - Map & geocoding features
- `WALI_NIKAH_VALIDATION.md` - Guardian validation logic
- `STATUS_MANAGEMENT.md` - Status management patterns
- `CORS_SETUP.md` - CORS configuration

### API (`docs/api/`)
- `API_DOCUMENTATION.md` - Complete API reference
- `API_TESTING_DOCUMENTATION.md` - API testing guide

### Performance (`docs/performance/`)
- `PERFORMANCE_ASSESSMENT.md` - Performance analysis
- `OPTIMIZATION_SUMMARY.md` - Optimization summary

### Architecture (`docs/architecture/`)
- `FOLDER_STRUCTURE_PLAN.md` - Full restructure plan
- `FOLDER_STRUCTURE.md` - This file

## Scripts

All helper scripts are in `scripts/`:
- `test_autocomplete.sh` - Test address autocomplete
- `test_railway_endpoint.sh` - Test Railway deployment
- `RESTRUCTURE_COMMANDS.sh` - Full restructure script
- `RESTRUCTURE_LIGHT.sh` - Light restructure (this script)

## Deployment Configs

All deployment configurations in `deployments/`:

**Railway** (`deployments/railway/`):
- `railway.json` - Railway build config
- `nixpacks.toml` - Nixpacks config

**Docker** (`deployments/docker/`):
- `Dockerfile` - Docker image
- `docker-compose.yml` - Docker Compose setup

## Database Migrations

SQL migrations in `migrations/`:
- `init.sql` - Initial database schema

## Build Artifacts

Build output goes to `bin/` (gitignored):
- `bin/simnikah-api` - Built binary

## Quick Commands

```bash
# Build
make build

# Run
make run

# Development mode
make dev

# Test
make test

# Clean
make clean
```

See `Makefile` for all available commands.
STRUCTURE_EOF
echo "✅ Folder structure documentation created!"
echo ""

# Final summary
echo "=================================================="
echo "✅ LIGHT RESTRUCTURE COMPLETE!"
echo "=================================================="
echo ""
echo "📊 What was reorganized:"
echo "  ✅ 20+ documentation files → docs/"
echo "  ✅ Test files → tests/fixtures/"
echo "  ✅ Scripts → scripts/"
echo "  ✅ Deployment configs → deployments/"
echo "  ✅ Database migrations → migrations/"
echo "  ✅ Makefile created"
echo "  ✅ .gitignore updated"
echo ""
echo "📝 New Makefile commands:"
echo "  make build   - Build aplikasi"
echo "  make run     - Jalankan aplikasi"
echo "  make dev     - Development mode"
echo "  make test    - Run tests"
echo "  make help    - Lihat semua commands"
echo ""
echo "📁 Root directory now:"
ls -1 | head -20
echo ""
echo "✅ Struktur folder lebih rapih & professional!"
echo ""
STRUCTURE_EOF
chmod +x RESTRUCTURE_LIGHT.sh

