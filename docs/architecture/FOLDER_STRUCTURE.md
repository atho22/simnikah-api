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
