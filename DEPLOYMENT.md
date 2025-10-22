# 🚀 SimNikah Deployment Guide - LeapCell

## Prerequisites

1. **LeapCell Account** - Daftar di [LeapCell](https://leapcell.io)
2. **GitHub Repository** - Push project ke GitHub
3. **Database** - MySQL 8.0+ (LeapCell menyediakan)

## Step 1: Prepare Repository

### 1.1 Push ke GitHub
```bash
git add .
git commit -m "Add LeapCell deployment configuration"
git push origin main
```

### 1.2 Pastikan file-file berikut ada:
- ✅ `Dockerfile`
- ✅ `leapcell.yaml`
- ✅ `env.example`
- ✅ `docker-compose.yml` (untuk development)

## Step 2: LeapCell Setup

### 2.1 Login ke LeapCell Dashboard
1. Buka [LeapCell Dashboard](https://dashboard.leapcell.io)
2. Login dengan akun Anda

### 2.2 Create New Project
1. Klik "New Project"
2. Pilih "Deploy from GitHub"
3. Connect GitHub repository
4. Pilih repository SimNikah

### 2.3 Configure Database
1. Di dashboard, pilih "Database"
2. Create new MySQL database
3. Catat credentials:
   - Host: `your-db-host.leapcell.io`
   - Port: `3306`
   - Username: `your-username`
   - Password: `your-password`
   - Database: `your-db-name`

## Step 3: Environment Variables

### 3.1 Set Environment Variables di LeapCell
Di dashboard LeapCell, set environment variables:

```bash
# Database Configuration
DB_HOST=your-db-host.leapcell.io
DB_PORT=3306
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=your-db-name

# JWT Configuration
JWT_KEY=your-super-secret-jwt-key-minimum-32-characters-long

# Server Configuration
PORT=8080
GIN_MODE=release
```

### 3.2 Generate JWT Secret
```bash
# Generate random JWT secret (32+ characters)
openssl rand -base64 32
```

## Step 4: Deploy Application

### 4.1 Deploy via LeapCell Dashboard
1. Go to "Deployments"
2. Click "Deploy"
3. Select branch: `main`
4. Build will start automatically

### 4.2 Monitor Deployment
- Check build logs
- Monitor health checks
- Verify database connection

## Step 5: Post-Deployment Setup

### 5.1 Verify Health Check
```bash
curl https://your-app.leapcell.io/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 5.2 Test API Endpoints
```bash
# Test registration
curl -X POST https://your-app.leapcell.io/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123",
    "nama": "Administrator",
    "role": "kepala_kua"
  }'
```

## Step 6: Database Migration

### 6.1 Auto Migration
Database tables akan dibuat otomatis saat aplikasi pertama kali dijalankan melalui GORM AutoMigrate.

### 6.2 Verify Tables
Login ke database dan pastikan tabel-tabel berikut ada:
- `users`
- `calon_pasangans`
- `pendaftaran_nikahs`
- `data_orang_tuas`
- `wali_nikahs`
- `staff_kuas`
- `penghulus`
- `bimbingan_perkawinans`
- `pendaftaran_bimbingans`
- `notifikasis`

## Step 7: Production Configuration

### 7.1 SSL/HTTPS
LeapCell menyediakan SSL certificate otomatis untuk domain custom.

### 7.2 Custom Domain
1. Di dashboard, go to "Domains"
2. Add custom domain
3. Update DNS records
4. SSL akan diaktifkan otomatis

### 7.3 Monitoring
- Monitor logs di dashboard
- Set up alerts untuk error
- Monitor database performance

## Step 8: Backup & Security

### 8.1 Database Backup
LeapCell menyediakan automatic backup untuk database.

### 8.2 Security Checklist
- ✅ JWT secret key yang kuat
- ✅ Database credentials yang aman
- ✅ HTTPS enabled
- ✅ Environment variables tidak exposed
- ✅ Health checks configured

## Troubleshooting

### Common Issues

#### 1. Database Connection Failed
```bash
# Check database credentials
# Verify network connectivity
# Check database status
```

#### 2. Build Failed
```bash
# Check Dockerfile syntax
# Verify Go version compatibility
# Check dependencies
```

#### 3. Health Check Failed
```bash
# Check application logs
# Verify port configuration
# Check health endpoint
```

### Logs & Debugging
```bash
# View application logs
leapcell logs your-app-name

# View database logs
leapcell logs your-db-name
```

## API Documentation

### Base URL
```
https://your-app.leapcell.io
```

### Key Endpoints
- `GET /health` - Health check
- `POST /register` - User registration
- `POST /login` - User login
- `GET /profile` - User profile
- `POST /simnikah/pendaftaran/form-baru` - Marriage registration

## Support

- **LeapCell Documentation**: [docs.leapcell.io](https://docs.leapcell.io)
- **GitHub Issues**: [github.com/your-repo/issues](https://github.com/your-repo/issues)
- **Email Support**: support@leapcell.io

---

## Quick Start Commands

```bash
# Local development
docker-compose up -d

# Check logs
docker-compose logs -f simnikah-api

# Stop services
docker-compose down

# Rebuild and restart
docker-compose up --build -d
```

## Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | localhost | Yes |
| `DB_PORT` | Database port | 3306 | Yes |
| `DB_USER` | Database username | root | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | simnikah | Yes |
| `JWT_KEY` | JWT secret key | - | Yes |
| `PORT` | Application port | 8080 | No |
| `GIN_MODE` | Gin mode | debug | No |
