# 🔄 Migration Guide: LeapCell (PostgreSQL) → Railway (MySQL)

## 📋 Overview

Project ini telah di-migrate dari:
- **Platform**: LeapCell → Railway
- **Database**: PostgreSQL → MySQL

## ✅ Perubahan yang Dilakukan

### 1. Database Driver
- ❌ Removed: `gorm.io/driver/postgres`
- ✅ Added: `gorm.io/driver/mysql`

### 2. Database Connection
File: `config/config.go`
- Update connection string dari PostgreSQL ke MySQL format
- Update default port: 5432 → 3306
- Update default user: postgres → root

### 3. Docker Configuration
File: `docker-compose.yml`
- Replace PostgreSQL service dengan MySQL 8.0
- Update environment variables
- Add MySQL healthcheck

### 4. Railway Configuration
Files baru:
- ✅ `railway.json` - Railway deployment config
- ✅ `nixpacks.toml` - Build configuration
- ✅ `RAILWAY_DEPLOYMENT.md` - Deployment guide

Files dihapus:
- ❌ `leapcell.yaml` - Tidak lagi diperlukan

### 5. Environment Variables
File: `env.example`
- Update comment untuk Railway
- Add Railway MySQL variable references

---

## 🚀 Deployment Steps

### Quick Start

```bash
# 1. Update dependencies
go mod tidy

# 2. Test locally dengan Docker
docker-compose up -d

# 3. Push ke GitHub
git add .
git commit -m "Migrate to Railway with MySQL"
git push origin main

# 4. Deploy di Railway
# Follow: RAILWAY_DEPLOYMENT.md
```

### Detailed Steps

Baca panduan lengkap di **[RAILWAY_DEPLOYMENT.md](RAILWAY_DEPLOYMENT.md)**

---

## 🔧 Environment Variables untuk Railway

### Required Variables (Manual Setup)

```bash
# Database (otomatis dari Railway MySQL)
DB_HOST=${{MYSQL.MYSQL_HOST}}
DB_PORT=${{MYSQL.MYSQL_PORT}}
DB_USER=${{MYSQL.MYSQL_USER}}
DB_PASSWORD=${{MYSQL.MYSQL_PASSWORD}}
DB_NAME=${{MYSQL.MYSQL_DATABASE}}

# JWT (generate sendiri)
JWT_KEY=your-32-char-secret-key

# Server
PORT=8080
GIN_MODE=release

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

---

## 📊 Data Migration (Jika Ada Data di LeapCell)

### Option 1: Fresh Start (Recommended)
Railway akan create database baru dengan auto-migration.

### Option 2: Migrate Existing Data

#### Step 1: Export dari LeapCell PostgreSQL
```bash
# Export schema
pg_dump -h leapcell-host -U user -d simnikah --schema-only > schema.sql

# Export data
pg_dump -h leapcell-host -U user -d simnikah --data-only > data.sql
```

#### Step 2: Convert PostgreSQL to MySQL
```bash
# Manual conversion diperlukan karena perbedaan syntax
# Tools: pgloader, atau manual editing

# Perubahan utama:
# - SERIAL → AUTO_INCREMENT
# - BOOLEAN → TINYINT(1)
# - TEXT → LONGTEXT
# - timestamps → DATETIME
```

#### Step 3: Import ke Railway MySQL
```bash
# Connect to Railway MySQL
railway connect mysql

# Import
mysql -h railway-host -u user -p database < converted_schema.sql
mysql -h railway-host -u user -p database < converted_data.sql
```

**Note**: Karena GORM auto-migration, lebih mudah fresh start dan manual input data penting.

---

## 🧪 Testing

### Local Testing
```bash
# Start services
docker-compose up -d

# Check logs
docker-compose logs -f simnikah-api

# Test health
curl http://localhost:8080/health

# Stop services
docker-compose down
```

### Production Testing
```bash
# Test health
curl https://your-app.up.railway.app/health

# Test register
curl -X POST https://your-app.up.railway.app/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123",
    "nama": "Test User",
    "role": "user_biasa"
  }'
```

---

## ⚠️ Breaking Changes

### None for API Users
API endpoints dan response format tetap sama. Migration hanya internal.

### For Developers
- Database driver changed
- Environment variable references untuk Railway
- LeapCell config files removed

---

## 🔍 Troubleshooting

### Build Error: "cannot find package gorm.io/driver/postgres"
```bash
# Solution: Update dependencies
go mod tidy
```

### Database Connection Error
```bash
# Check environment variables di Railway
# Pastikan menggunakan Railway variable references:
DB_HOST=${{MYSQL.MYSQL_HOST}}
```

### CORS Error
```bash
# Update ALLOWED_ORIGINS di Railway
ALLOWED_ORIGINS=http://localhost:3000,https://your-frontend.com
```

---

## 📝 Checklist Migration

- [x] Update database driver (PostgreSQL → MySQL)
- [x] Update config/config.go
- [x] Update docker-compose.yml
- [x] Create railway.json
- [x] Create nixpacks.toml
- [x] Create RAILWAY_DEPLOYMENT.md
- [x] Update env.example
- [x] Update go.mod
- [x] Remove leapcell.yaml
- [ ] Deploy to Railway
- [ ] Test all endpoints
- [ ] Update frontend API URL
- [ ] Verify CORS working
- [ ] Test authentication
- [ ] Test database operations

---

## 📞 Support

Jika ada masalah:
1. Check **[RAILWAY_DEPLOYMENT.md](RAILWAY_DEPLOYMENT.md)**
2. Check **[CORS_SETUP.md](CORS_SETUP.md)**
3. Railway Discord: [discord.gg/railway](https://discord.gg/railway)

---

*Migration completed: 2024*

