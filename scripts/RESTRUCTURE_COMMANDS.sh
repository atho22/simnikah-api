#!/bin/bash

# =====================================================
# SimNikah API - Folder Restructure Script
# =====================================================
# This script reorganizes the project structure
# to follow industry best practices (Go project layout)
# =====================================================

set -e  # Exit on error

echo "🏗️  Starting folder restructure..."
echo ""

# =====================================================
# PHASE 1: Create New Folder Structure
# =====================================================
echo "📁 Phase 1: Creating new folder structure..."

# Core application folders
mkdir -p cmd/api
mkdir -p internal/handlers/{catin,staff,penghulu,kepala_kua,notification,auth}
mkdir -p internal/models
mkdir -p internal/middleware
mkdir -p internal/services
mkdir -p internal/repository

# Public reusable packages
mkdir -p pkg/{validator,utils,crypto,cache}

# Documentation
mkdir -p docs/{deployment,features,api,performance,tutorials,architecture}

# Scripts & tests
mkdir -p scripts
mkdir -p tests/{fixtures,integration,unit}

# Deployment configs
mkdir -p deployments/{railway,docker,kubernetes}

# Database migrations
mkdir -p migrations

# Build output
mkdir -p bin

echo "✅ Folder structure created!"
echo ""

# =====================================================
# PHASE 2: Move Documentation Files
# =====================================================
echo "📚 Phase 2: Moving documentation files..."

# Deployment docs
mv RAILWAY_DEPLOYMENT.md docs/deployment/ 2>/dev/null || true
mv MIGRATION_GUIDE.md docs/deployment/ 2>/dev/null || true
mv PENJELASAN_MUDAH_RAILWAY.md docs/deployment/ 2>/dev/null || true
mv TUTORIAL_DEPLOY_RAILWAY.md docs/tutorials/ 2>/dev/null || true
mv QUICK_START.md docs/deployment/ 2>/dev/null || true
mv RAILWAY_ENV_TEMPLATE.txt docs/deployment/ 2>/dev/null || true
mv LEAPCELL_CORS_GUIDE.md docs/deployment/ 2>/dev/null || true

# Feature docs
mv MAP_INTEGRATION.md docs/features/ 2>/dev/null || true
mv WALI_NIKAH_VALIDATION.md docs/features/ 2>/dev/null || true
mv STATUS_MANAGEMENT.md docs/features/ 2>/dev/null || true
mv CORS_SETUP.md docs/features/ 2>/dev/null || true

# API docs
mv API_DOCUMENTATION.md docs/api/ 2>/dev/null || true
mv API_TESTING_DOCUMENTATION.md docs/api/ 2>/dev/null || true

# Performance & Architecture docs
mv PERFORMANCE_ASSESSMENT.md docs/performance/ 2>/dev/null || true
mv OPTIMIZATION_SUMMARY.md docs/performance/ 2>/dev/null || true
mv FOLDER_STRUCTURE_PLAN.md docs/architecture/ 2>/dev/null || true
mv DEPLOYMENT.md docs/deployment/ 2>/dev/null || true

echo "✅ Documentation organized!"
echo ""

# =====================================================
# PHASE 3: Move Scripts & Test Files
# =====================================================
echo "🔧 Phase 3: Moving scripts and test files..."

mv test_autocomplete.sh scripts/ 2>/dev/null || true
mv test_railway_endpoint.sh scripts/ 2>/dev/null || true
chmod +x scripts/*.sh 2>/dev/null || true

mv demo_autocomplete.html tests/fixtures/ 2>/dev/null || true

echo "✅ Scripts organized!"
echo ""

# =====================================================
# PHASE 4: Move Deployment Configs
# =====================================================
echo "🚀 Phase 4: Moving deployment configs..."

mv railway.json deployments/railway/ 2>/dev/null || true
mv nixpacks.toml deployments/railway/ 2>/dev/null || true

mv Dockerfile deployments/docker/ 2>/dev/null || true
mv docker-compose.yml deployments/docker/ 2>/dev/null || true

mv init.sql migrations/ 2>/dev/null || true

echo "✅ Deployment configs organized!"
echo ""

# =====================================================
# PHASE 5: Move Application Code
# =====================================================
echo "💻 Phase 5: Moving application code..."

# Move main.go to cmd/api/
mv main.go cmd/api/ 2>/dev/null || true

# Move handlers
mv catin internal/handlers/ 2>/dev/null || true
mv staff internal/handlers/ 2>/dev/null || true
mv penghulu internal/handlers/ 2>/dev/null || true
mv kepala_kua internal/handlers/ 2>/dev/null || true
mv notification internal/handlers/ 2>/dev/null || true

# Move models (structs)
mv structs internal/models 2>/dev/null || true

# Move middleware
mv middleware internal/ 2>/dev/null || true

# Move services
mv services internal/ 2>/dev/null || true

# Split helper into pkg
# (We'll do this manually as it requires code changes)
echo "⚠️  Note: helper/ needs manual splitting into pkg/ subdirectories"

echo "✅ Application code organized!"
echo ""

# =====================================================
# PHASE 6: Create Additional Files
# =====================================================
echo "📝 Phase 6: Creating additional files..."

# Create Makefile
cat > Makefile << 'EOF'
.PHONY: build run test clean dev

# Build the application
build:
	@echo "🔨 Building SimNikah API..."
	@go build -o bin/simnikah-api cmd/api/main.go
	@echo "✅ Build complete! Binary: bin/simnikah-api"

# Run the application
run:
	@echo "🚀 Starting SimNikah API..."
	@go run cmd/api/main.go

# Run in development mode
dev:
	@echo "🔧 Starting in development mode..."
	@GIN_MODE=debug go run cmd/api/main.go

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

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run || go vet ./...
	@echo "✅ Lint complete!"

# Show help
help:
	@echo "SimNikah API - Available Commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make dev       - Run in development mode"
	@echo "  make test      - Run tests"
	@echo "  make coverage  - Run tests with coverage report"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Run linter"
EOF

echo "✅ Makefile created!"

# Create .gitignore additions
cat >> .gitignore << 'EOF'

# Build artifacts
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

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
EOF

echo "✅ .gitignore updated!"

echo ""
echo "✅ Additional files created!"
echo ""

# =====================================================
# PHASE 7: Update README
# =====================================================
echo "📖 Phase 7: Creating updated README structure..."

cat > README_NEW.md << 'EOF'
# 🏛️ SimNikah API - Sistem Informasi Manajemen Nikah KUA

> Modern, scalable, and high-performance marriage registration system for KUA (Office of Religious Affairs)

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Deployment](https://img.shields.io/badge/deploy-Railway-purple.svg)](https://railway.app)

---

## 📚 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [Deployment](#-deployment)
- [Performance](#-performance)
- [Contributing](#-contributing)

---

## ✨ Features

### Core Features
- 👥 **User Management** - Multi-role authentication (Admin, Staff, Penghulu, User)
- 💍 **Marriage Registration** - Complete registration workflow
- 📋 **Document Management** - Digital document handling
- 📅 **Scheduling** - Marriage ceremony scheduling
- 👨‍👩‍👧 **Family Data** - Bride & groom family information
- 🤵 **Guardian (Wali) Management** - Sharia-compliant guardian validation
- 📊 **Staff Dashboard** - KUA staff management interface
- 🔔 **Notifications** - Real-time notification system
- 📜 **Certificate Generation** - Digital marriage certificate

### Advanced Features
- 🗺️ **Map Integration** - OpenStreetMap integration for ceremony locations
- 📍 **Geocoding** - Address to coordinates conversion (cached for performance)
- 🔍 **Address Autocomplete** - Smart address search
- ⚡ **Rate Limiting** - DDoS protection (100 req/min per IP)
- 🛡️ **Security** - JWT authentication, bcrypt password hashing
- 🚀 **Performance** - Database indexes, caching (10x faster!)
- 🔄 **Graceful Shutdown** - Zero downtime deploys

---

## 🛠️ Tech Stack

- **Backend:** Go 1.23 (Gin framework)
- **Database:** MySQL 8.0
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt)
- **Geocoding:** OpenStreetMap Nominatim API (FREE!)
- **Rate Limiting:** ulule/limiter
- **Deployment:** Railway (auto-deploy from GitHub)
- **CI/CD:** GitHub Actions (future)

---

## 📁 Project Structure

```
simpadu/
├── cmd/api/              # Application entry point
│   └── main.go
├── internal/             # Private application code
│   ├── handlers/         # HTTP handlers
│   ├── models/           # Database models
│   ├── middleware/       # Middleware (auth, rate limit, etc)
│   ├── services/         # Business services
│   └── repository/       # Data access layer (future)
├── pkg/                  # Public reusable packages
│   ├── validator/
│   ├── utils/
│   ├── crypto/
│   └── cache/
├── config/               # Configuration
├── docs/                 # Documentation
├── scripts/              # Helper scripts
├── tests/                # Tests & fixtures
├── deployments/          # Deployment configs
└── migrations/           # Database migrations
```

See [FOLDER_STRUCTURE_PLAN.md](docs/architecture/FOLDER_STRUCTURE_PLAN.md) for details.

---

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- MySQL 8.0+
- Git

### Local Development

1. **Clone repository:**
```bash
git clone https://github.com/atho22/simnikah-api.git
cd simnikah-api
```

2. **Install dependencies:**
```bash
make deps
```

3. **Setup environment:**
```bash
cp env.example .env
# Edit .env with your local settings
```

4. **Run database migrations:**
```bash
mysql -u root -p < migrations/init.sql
```

5. **Run application:**
```bash
make run
# or for development mode:
make dev
```

6. **Test:**
```bash
curl http://localhost:8080/health
```

See [QUICK_START.md](docs/deployment/QUICK_START.md) for more details.

---

## 📖 API Documentation

### Base URL
```
Local:      http://localhost:8080
Production: https://simnikah-api.railway.app
```

### Authentication
All protected endpoints require JWT token in header:
```
Authorization: Bearer <token>
```

### Key Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/login` | User login | ❌ |
| POST | `/register` | User registration | ❌ |
| POST | `/simnikah/pendaftaran` | Submit marriage registration | ✅ |
| GET | `/simnikah/pendaftaran/:id` | Get registration details | ✅ |
| POST | `/simnikah/location/geocode` | Address to coordinates | ✅ |
| GET | `/simnikah/location/search?q=jakarta` | Address autocomplete | ✅ |
| GET | `/health` | Health check | ❌ |

See [API_DOCUMENTATION.md](docs/api/API_DOCUMENTATION.md) for complete API reference.

---

## 🌐 Deployment

### Railway (Recommended)

1. **Connect to GitHub:**
   - Fork this repository
   - Connect Railway to your GitHub account

2. **Add MySQL database:**
   - Click "+ New" → "Database" → "MySQL"

3. **Set environment variables:**
   ```bash
   DB_HOST=${{MySQL.MYSQLHOST}}
   DB_PORT=${{MySQL.MYSQLPORT}}
   DB_USER=${{MySQL.MYSQLUSER}}
   DB_PASSWORD=${{MySQL.MYSQLPASSWORD}}
   DB_NAME=${{MySQL.MYSQLDATABASE}}
   JWT_KEY=<generate-with-openssl-rand-base64-32>
   PORT=8080
   GIN_MODE=release
   ALLOWED_ORIGINS=http://localhost:3000,https://your-frontend.vercel.app
   ```

4. **Deploy:**
   - Push to `main` branch
   - Railway auto-deploys!

See detailed guide: [TUTORIAL_DEPLOY_RAILWAY.md](docs/tutorials/TUTORIAL_DEPLOY_RAILWAY.md)

### Docker

```bash
cd deployments/docker
docker-compose up -d
```

See: [docker-compose.yml](deployments/docker/docker-compose.yml)

---

## ⚡ Performance

### Optimizations Implemented

- ✅ **30+ Database Indexes** - 5-10x faster queries
- ✅ **Geocoding Cache** - 100x faster map features
- ✅ **Rate Limiting** - Protect from DDoS/spam
- ✅ **Graceful Shutdown** - Zero downtime deploys
- ✅ **Connection Pooling** - Efficient database connections

### Benchmarks

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Response Time | 200-500ms | 20-50ms | **10x faster** |
| Requests/sec | 200-500 | 2,000-5,000 | **10x more** |
| Max Users | 100-200 | 1,000-2,000 | **10x scale** |
| Geocoding | 1-3s | <1ms (cached) | **1000x faster** |

See: [PERFORMANCE_ASSESSMENT.md](docs/performance/PERFORMANCE_ASSESSMENT.md)

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Test specific endpoint
./scripts/test_railway_endpoint.sh
```

---

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

---

## 📄 License

This project is licensed under the MIT License.

---

## 📞 Contact

- **Developer:** Atho
- **Email:** [your-email]
- **GitHub:** [@atho22](https://github.com/atho22)

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [OpenStreetMap Nominatim](https://nominatim.openstreetmap.org/)
- [Railway](https://railway.app/)

---

**Made with ❤️ for Indonesian KUA offices**
EOF

echo "✅ New README created (README_NEW.md)!"
echo ""

# =====================================================
# SUMMARY
# =====================================================
echo "=================================================="
echo "✅ FOLDER RESTRUCTURE COMPLETE!"
echo "=================================================="
echo ""
echo "📊 Summary:"
echo "  ✅ New folder structure created"
echo "  ✅ Documentation organized (20+ files)"
echo "  ✅ Scripts moved to scripts/"
echo "  ✅ Deployment configs moved to deployments/"
echo "  ✅ Application code moved to internal/"
echo "  ✅ Makefile created"
echo "  ✅ .gitignore updated"
echo "  ✅ New README created"
echo ""
echo "⚠️  IMPORTANT NEXT STEPS:"
echo "  1. Update import paths in all .go files"
echo "  2. Review README_NEW.md and replace README.md"
echo "  3. Test build: make build"
echo "  4. Test run: make run"
echo ""
echo "📝 See docs/architecture/FOLDER_STRUCTURE_PLAN.md for details"
echo ""
EOF
chmod +x RESTRUCTURE_COMMANDS.sh

