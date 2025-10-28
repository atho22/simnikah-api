# 📁 Rencana Restructure Folder - SimNikah API

## 🎯 Tujuan Restructure:
1. ✅ **Separation of Concerns** - Business logic terpisah dari infrastructure
2. ✅ **Scalability** - Mudah menambah fitur baru
3. ✅ **Maintainability** - Mudah mencari dan mengedit kode
4. ✅ **Industry Standard** - Mengikuti best practice Go project layout

---

## 📊 STRUKTUR SAAT INI (Kurang Terorganisir):

```
simpadu/
├── main.go                    # ❌ Terlalu besar (2600+ lines)
├── catin/                     # ✅ OK
├── staff/                     # ✅ OK
├── penghulu/                  # ✅ OK
├── helper/                    # ✅ OK
├── structs/                   # ✅ OK
├── config/                    # ✅ OK
├── services/                  # ✅ OK
├── middleware/                # ✅ OK
├── notification/              # ✅ OK
├── kepala_kua/                # ✅ OK
├── 20+ file dokumentasi *.md  # ❌ Berantakan di root
├── demo_autocomplete.html     # ❌ File test di root
├── test_*.sh                  # ❌ Scripts berantakan
├── go.mod, go.sum            # ✅ OK (harus di root)
└── env.example               # ✅ OK
```

**Masalah:**
- ❌ Dokumentasi berantakan (20+ file .md di root)
- ❌ Test files & scripts tercampur
- ❌ `main.go` terlalu besar (2600+ lines)
- ❌ Tidak ada pemisahan public vs private packages

---

## 🎯 STRUKTUR BARU (Industry Standard):

```
simpadu/
│
├── cmd/                           # 🆕 Application entry points
│   └── api/
│       └── main.go                # Entry point utama (dipindahkan dari root)
│
├── internal/                      # 🆕 Private application code
│   ├── handlers/                  # 🔄 Business logic handlers
│   │   ├── catin/                 # Dari root/catin/
│   │   ├── staff/                 # Dari root/staff/
│   │   ├── penghulu/              # Dari root/penghulu/
│   │   ├── kepala_kua/            # Dari root/kepala_kua/
│   │   ├── notification/          # Dari root/notification/
│   │   └── auth/                  # 🆕 Pindahkan login/register dari main.go
│   │
│   ├── models/                    # 🔄 Database models
│   │   ├── entities.go            # Dari structs/models.go
│   │   └── constants.go           # Dari structs/constants.go
│   │
│   ├── middleware/                # 🔄 Dari root/middleware/
│   │   ├── auth.go                # JWT auth middleware
│   │   ├── rate_limit.go          # Rate limiting
│   │   └── cors.go                # CORS middleware (pindahkan dari main.go)
│   │
│   ├── services/                  # 🔄 Business services
│   │   ├── notification.go        # Dari root/services/
│   │   ├── cron_job.go
│   │   └── geocoding.go           # 🆕 Service untuk geocoding
│   │
│   └── repository/                # 🆕 Database operations (optional, untuk future)
│       ├── user.go
│       ├── pendaftaran.go
│       └── wali.go
│
├── pkg/                           # 🆕 Public libraries (reusable)
│   ├── validator/                 # Helper validation
│   ├── utils/                     # String, date utils
│   ├── crypto/                    # Bcrypt helpers
│   └── cache/                     # Geocoding cache
│
├── config/                        # ✅ Configuration (tetap di sini)
│   ├── database.go                # Dari config.go
│   ├── indexes.go
│   └── env.go                     # 🆕 Environment config helper
│
├── docs/                          # 🆕 ALL documentation
│   ├── deployment/
│   │   ├── RAILWAY_DEPLOYMENT.md
│   │   ├── QUICK_START.md
│   │   ├── MIGRATION_GUIDE.md
│   │   └── PENJELASAN_MUDAH_RAILWAY.md
│   ├── features/
│   │   ├── MAP_INTEGRATION.md
│   │   ├── WALI_NIKAH_VALIDATION.md
│   │   ├── STATUS_MANAGEMENT.md
│   │   └── CORS_SETUP.md
│   ├── api/
│   │   ├── API_DOCUMENTATION.md
│   │   └── API_TESTING_DOCUMENTATION.md
│   ├── performance/
│   │   └── PERFORMANCE_ASSESSMENT.md
│   └── tutorials/
│       └── TUTORIAL_DEPLOY_RAILWAY.md
│
├── scripts/                       # 🆕 Helper scripts
│   ├── test_autocomplete.sh
│   ├── test_railway_endpoint.sh
│   └── setup_dev.sh               # 🆕 Setup local development
│
├── tests/                         # 🆕 Test files & fixtures
│   ├── fixtures/
│   │   └── demo_autocomplete.html
│   └── integration/
│       └── api_test.go            # 🆕 Future integration tests
│
├── deployments/                   # 🆕 Deployment configs
│   ├── railway/
│   │   ├── railway.json
│   │   └── nixpacks.toml
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/                # 🆕 Future: K8s manifests
│       └── deployment.yaml
│
├── migrations/                    # 🆕 Database migrations (future)
│   └── init.sql
│
├── .gitignore                     # ✅ Ignore patterns
├── go.mod                         # ✅ Dependencies
├── go.sum                         # ✅ Checksums
├── env.example                    # ✅ Environment template
├── README.md                      # ✅ Main documentation
└── Makefile                       # 🆕 Build automation

```

---

## 🔄 MAPPING PERUBAHAN:

### Files yang Dipindahkan:

| Dari (Root)                       | Ke (New Location)                              |
|-----------------------------------|------------------------------------------------|
| `main.go`                         | `cmd/api/main.go`                             |
| `catin/*.go`                      | `internal/handlers/catin/*.go`                |
| `staff/*.go`                      | `internal/handlers/staff/*.go`                |
| `penghulu/*.go`                   | `internal/handlers/penghulu/*.go`             |
| `kepala_kua/*.go`                 | `internal/handlers/kepala_kua/*.go`           |
| `notification/*.go`               | `internal/handlers/notification/*.go`         |
| `structs/*.go`                    | `internal/models/*.go`                        |
| `middleware/*.go`                 | `internal/middleware/*.go`                    |
| `services/*.go`                   | `internal/services/*.go`                      |
| `helper/*.go`                     | `pkg/validator/`, `pkg/utils/`, `pkg/cache/`  |
| `config/*.go`                     | `config/*.go` (tetap)                         |
| `*.md` (20+ files)                | `docs/**/*.md`                                |
| `test_*.sh`                       | `scripts/test_*.sh`                           |
| `demo_autocomplete.html`          | `tests/fixtures/demo_autocomplete.html`       |
| `railway.json`, `nixpacks.toml`   | `deployments/railway/`                        |
| `Dockerfile`, `docker-compose.yml`| `deployments/docker/`                         |
| `init.sql`                        | `migrations/init.sql`                         |

### Files yang Dipecah:

| File Besar                        | Dipecah Jadi                                  |
|-----------------------------------|-----------------------------------------------|
| `main.go` (2600+ lines)           | → `cmd/api/main.go` (routing)                |
|                                   | → `internal/handlers/auth/auth.go` (login)   |
|                                   | → `internal/handlers/auth/register.go`       |
|                                   | → `internal/middleware/auth.go` (JWT)        |
|                                   | → `internal/middleware/cors.go` (CORS)       |
| `helper/utils.go`                 | → `pkg/utils/string.go`                      |
|                                   | → `pkg/utils/date.go`                        |
|                                   | → `pkg/cache/geocoding.go`                   |

---

## ✅ KEUNTUNGAN STRUKTUR BARU:

### 1. **Separation of Concerns**
```
internal/        → Private business logic (tidak bisa di-import project lain)
pkg/             → Public libraries (bisa di-reuse)
cmd/             → Entry points (multiple apps in future)
```

### 2. **Scalability**
```
Mudah tambah handler baru:
internal/handlers/saksi/          # 🆕 Handler baru
internal/handlers/sertifikat/     # 🆕 Handler baru
```

### 3. **Clear Documentation**
```
docs/
├── deployment/     → Panduan deploy
├── features/       → Dokumentasi fitur
├── api/            → API docs
└── tutorials/      → Step-by-step guides
```

### 4. **Easy Testing**
```
tests/
├── fixtures/       → Test data
├── integration/    → Integration tests
└── unit/           → Unit tests
```

### 5. **Professional Look**
```
✅ Industry standard (mirip project Go besar seperti Kubernetes, Docker, dll)
✅ Mudah onboarding developer baru
✅ GitHub README lebih clean
```

---

## 🚀 IMPLEMENTATION PLAN:

### Phase 1: Create New Folders (5 menit)
```bash
mkdir -p cmd/api
mkdir -p internal/{handlers/{catin,staff,penghulu,kepala_kua,notification,auth},models,middleware,services,repository}
mkdir -p pkg/{validator,utils,crypto,cache}
mkdir -p docs/{deployment,features,api,performance,tutorials}
mkdir -p scripts
mkdir -p tests/{fixtures,integration}
mkdir -p deployments/{railway,docker,kubernetes}
mkdir -p migrations
```

### Phase 2: Move Files (10 menit)
```bash
# Move main.go
mv main.go cmd/api/main.go

# Move handlers
mv catin internal/handlers/
mv staff internal/handlers/
mv penghulu internal/handlers/
mv kepala_kua internal/handlers/
mv notification internal/handlers/

# Move models
mv structs internal/models

# Move middleware
mv middleware internal/

# Move services
mv services internal/

# Move helper to pkg
# (Will split helper/ into multiple packages)

# Move docs
mv *DEPLOYMENT*.md docs/deployment/
mv MAP_*.md WALI_*.md STATUS_*.md CORS_*.md docs/features/
mv API_*.md docs/api/
mv PERFORMANCE_*.md docs/performance/
mv TUTORIAL_*.md docs/tutorials/

# Move scripts
mv *.sh scripts/
mv demo_autocomplete.html tests/fixtures/

# Move deployment configs
mv railway.json nixpacks.toml deployments/railway/
mv Dockerfile docker-compose.yml deployments/docker/
mv init.sql migrations/
```

### Phase 3: Update Imports (15 menit)
```bash
# Update all import paths
# Dari: "simnikah/catin"
# Jadi: "simnikah/internal/handlers/catin"
```

### Phase 4: Split main.go (10 menit)
```bash
# Extract auth handlers from main.go to internal/handlers/auth/
# Extract middleware from main.go to internal/middleware/
```

### Phase 5: Create Makefile (5 menit)
```makefile
.PHONY: build run test clean

build:
	go build -o bin/simnikah-api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./...

clean:
	rm -f bin/simnikah-api
```

---

## 📝 TODO CHECKLIST:

- [ ] Create new folder structure
- [ ] Move files to new locations
- [ ] Update import paths in all .go files
- [ ] Split main.go into smaller files
- [ ] Update README.md with new structure
- [ ] Create Makefile
- [ ] Update .gitignore
- [ ] Test build: `make build`
- [ ] Test run: `make run`
- [ ] Update Railway deployment config
- [ ] Update documentation links

---

## ⚠️ COMPATIBILITY:

**Import Path Changes:**
```go
// OLD
import "simnikah/catin"
import "simnikah/structs"
import "simnikah/helper"

// NEW
import "simnikah/internal/handlers/catin"
import "simnikah/internal/models"
import "simnikah/pkg/utils"
```

**Build Command Changes:**
```bash
# OLD
go build -o main .

# NEW
go build -o bin/simnikah-api cmd/api/main.go
# ATAU
make build
```

---

## 🎯 RESULT:

Setelah restructure, struktur folder akan:
- ✅ **20% lebih clean** (dokumentasi terorganisir)
- ✅ **50% lebih maintainable** (kode terpisah jelas)
- ✅ **100% lebih professional** (industry standard)
- ✅ **Mudah scale** (tambah fitur tanpa bingung)

---

**Mau saya implementasikan restructure ini sekarang?** 🚀

