# 🚀 Panduan Setup Railway - SimNikah API

## 📋 Checklist Persiapan

- [ ] Akun GitHub (gratis)
- [ ] Akun Railway (gratis - daftar di railway.app)
- [ ] Project sudah di-push ke GitHub
- [ ] Terminal/Command Prompt siap

---

## 🎯 STEP 1: Persiapan Project (5 menit)

### 1.1 Pastikan Project Siap

```bash
# Pastikan kamu di folder project
cd C:\Users\user\simnikah-api

# Cek apakah ada file penting
ls railway.json
ls deployments/railway/nixpacks.toml
```

### 1.2 Push ke GitHub (jika belum)

```bash
# Init git (jika belum)
git init

# Add remote (ganti YOUR_USERNAME)
git remote add origin https://github.com/YOUR_USERNAME/simnikah-api.git

# Add semua file
git add .

# Commit
git commit -m "Ready for Railway deployment"

# Push
git push -u origin main
```

---

## 🎯 STEP 2: Setup Railway Project (10 menit)

### 2.1 Login ke Railway

1. Buka [railway.app](https://railway.app)
2. Click **"Start a New Project"** atau **"Login"**
3. Pilih **"Login with GitHub"**
4. Authorize Railway untuk akses GitHub

### 2.2 Create New Project

1. Di Railway Dashboard, click **"New Project"**
2. Pilih **"Deploy from GitHub repo"**
3. Pilih repository **`simnikah-api`** (atau nama repo kamu)
4. Railway akan otomatis mulai build

**⏳ Tunggu 2-3 menit...** (Build pertama mungkin gagal - NORMAL!)

---

## 🎯 STEP 3: Tambah MySQL Database (3 menit)

### 3.1 Add MySQL Service

1. Di project Railway kamu, click tombol **"+ New"** (pojok kanan atas)
2. Pilih **"Database"**
3. Pilih **"Add MySQL"**
4. Tunggu 10-20 detik sampai MySQL siap ✅

### 3.2 Cek Nama Service MySQL

**PENTING!** Catat nama service MySQL kamu:
- Buka tab **"Services"** di Railway
- Lihat nama service MySQL (biasanya: `MySQL` atau `mysql` atau `mysql-prod`)
- **Catat nama ini!** (akan dipakai di environment variables)

---

## 🎯 STEP 4: Generate JWT Secret (1 menit)

### 4.1 Generate Secret Key

**Windows PowerShell:**
```powershell
# Generate random JWT secret
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

**Atau gunakan online generator:**
- Buka: https://www.random.org/strings/
- Length: 32
- Characters: Alphanumeric
- Copy hasilnya

**Atau gunakan script yang sudah saya buat:**
```powershell
# Jalankan script generate-jwt-key.ps1
.\scripts\generate-jwt-key.ps1
```

**✅ Copy hasil JWT secret!** (contoh: `Xy9pQ2vN8kL5mT7rW3sH6jF1dA4gC9bE8zX2vN5mT0=`)

---

## 🎯 STEP 5: Setup Environment Variables (5 menit)

### 5.1 Buka Variables Settings

1. Di Railway Dashboard, click **service aplikasi kamu** (bukan MySQL)
2. Pilih tab **"Variables"**
3. Click **"RAW Editor"** (pojok kanan atas)

### 5.2 Paste Environment Variables

**⚠️ PENTING:** Ganti 2 hal:
1. `YOUR_JWT_SECRET_HERE` → Paste JWT secret yang sudah di-generate
2. `MySQL` → Ganti dengan nama service MySQL kamu (lihat Step 3.2)

**Copy dan paste kode ini:**

```bash
# ============================================
# DATABASE CONFIGURATION
# Railway otomatis inject dari MySQL service
# ============================================
# ⚠️ GANTI "MySQL" dengan nama service MySQL kamu!
DB_HOST=${{MySQL.MYSQL_HOST}}
DB_PORT=${{MySQL.MYSQL_PORT}}
DB_USER=${{MySQL.MYSQL_USER}}
DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
DB_NAME=${{MySQL.MYSQL_DATABASE}}

# ============================================
# JWT CONFIGURATION
# ============================================
# ⚠️ GANTI dengan JWT secret yang sudah di-generate!
JWT_KEY=YOUR_JWT_SECRET_HERE

# ============================================
# SERVER CONFIGURATION
# ============================================
PORT=8080
GIN_MODE=release

# ============================================
# CORS CONFIGURATION
# ============================================
# Untuk development (frontend local)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:8080

# Untuk production (ganti setelah frontend di-deploy)
# ALLOWED_ORIGINS=https://your-frontend.vercel.app,https://your-app.netlify.app
```

### 5.3 Cara Mengetahui Nama Service MySQL

Jika tidak yakin nama service MySQL:

1. Di Railway Dashboard → Tab **"Services"**
2. Lihat service yang bertipe **"MySQL"**
3. Nama service biasanya di bawah icon MySQL
4. Contoh nama: `MySQL`, `mysql`, `mysql-prod`, `database`

**Jika nama service adalah `mysql-prod`, maka:**
```bash
DB_HOST=${{mysql-prod.MYSQL_HOST}}
DB_PORT=${{mysql-prod.MYSQL_PORT}}
# dst...
```

### 5.4 Save Variables

1. Click **"Update Variables"** (pojok kanan atas)
2. Railway akan **otomatis redeploy** setelah variables diupdate! ✅

---

## 🎯 STEP 6: Monitor Deployment (5 menit)

### 6.1 Lihat Logs

1. Pilih tab **"Deployments"**
2. Click deployment yang sedang running (yang terbaru)
3. Scroll ke bawah untuk lihat logs real-time

### 6.2 Tunggu Success

Tunggu sampai kamu lihat log seperti ini:

```
✅ === DATABASE CONFIGURATION ===
✅ DB_HOST: ...
✅ DB_PORT: 3306
✅ DB_USER: root
✅ Connected to MySQL database successfully
✅ Database migration completed successfully
✅ Info: CORS allowed origins: [...]
✅ 🚀 Server starting on port 8080
```

**⏳ Biasanya 2-5 menit**

### 6.3 Cek Status

Status deployment harus: **✅ SUCCESS** (warna hijau)

- ✅ **SUCCESS** → Lanjut ke Step 7
- 🔄 **BUILDING** → Tunggu sebentar
- ❌ **FAILED** → Lihat Troubleshooting di bawah

---

## 🎯 STEP 7: Generate Domain & Test (5 menit)

### 7.1 Generate Public Domain

1. Pilih tab **"Settings"**
2. Scroll ke section **"Domains"**
3. Click **"Generate Domain"**
4. Railway akan beri URL seperti: `https://simnikah-production-xxxx.up.railway.app`
5. **COPY URL INI!** 📋

### 7.2 Test Health Check

**Test dari Browser:**
```
https://simnikah-production-xxxx.up.railway.app/health
```

**Atau dari Terminal/PowerShell:**
```powershell
# Test health endpoint
curl https://simnikah-production-xxxx.up.railway.app/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2025-01-27T..."
}
```

**✅ BERHASIL!** API kamu sudah LIVE! 🎉

### 7.3 Test Register User

```powershell
# Test register endpoint
curl -X POST https://simnikah-production-xxxx.up.railway.app/register `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"test123\",\"nama\":\"Test User\",\"role\":\"user_biasa\"}'
```

**Expected Response:**
```json
{
  "message": "User berhasil dibuat",
  "user": {
    "user_id": "USR...",
    "username": "testuser",
    "email": "test@example.com",
    "nama": "Test User",
    "role": "user_biasa"
  }
}
```

### 7.4 Test Login

```powershell
# Test login endpoint
curl -X POST https://simnikah-production-xxxx.up.railway.app/login `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"testuser\",\"password\":\"test123\"}'
```

**Expected Response:**
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR...",
    "email": "test@example.com",
    "role": "user_biasa",
    "nama": "Test User"
  }
}
```

---

## 🎯 STEP 8: Update Frontend (2 menit)

### 8.1 Update API URL di Frontend

Jika kamu punya frontend, update file `.env`:

**React/Vite (.env.local):**
```bash
VITE_API_URL=https://simnikah-production-xxxx.up.railway.app
```

**Next.js (.env.local):**
```bash
NEXT_PUBLIC_API_URL=https://simnikah-production-xxxx.up.railway.app
```

### 8.2 Update CORS di Railway

Setelah frontend di-deploy (misal di Vercel/Netlify):

1. Kembali ke Railway → Variables
2. Update `ALLOWED_ORIGINS`:

```bash
ALLOWED_ORIGINS=https://your-frontend.vercel.app,https://simnikah-production-xxxx.up.railway.app
```

3. Save → Otomatis redeploy

---

## 🎉 SELESAI! Checklist Final

- [x] ✅ Project di GitHub
- [x] ✅ Railway project created
- [x] ✅ MySQL database added
- [x] ✅ Environment variables set
- [x] ✅ Deployment SUCCESS
- [x] ✅ Public domain generated
- [x] ✅ Health check working
- [x] ✅ Register/Login tested

---

## 🐛 Troubleshooting

### ❌ Build Failed: "cannot find module"

**Solusi:**
```bash
cd C:\Users\user\simnikah-api
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
git push
```

Railway akan auto redeploy setelah push.

---

### ❌ Database Connection Failed

**Cek 1: MySQL Service Running?**
- Railway Dashboard → Services → MySQL
- Pastikan status: **Running** ✅

**Cek 2: Variable Reference Benar?**
```bash
# ✅ BENAR (dengan ${{...}})
DB_HOST=${{MySQL.MYSQL_HOST}}

# ❌ SALAH (tanpa ${{...}})
DB_HOST=${MYSQL_HOST}
```

**Cek 3: Nama Service MySQL Benar?**
- Railway Dashboard → Services
- Lihat nama service MySQL
- Pastikan di variables menggunakan nama yang sama

**Contoh jika nama service adalah `mysql-prod`:**
```bash
DB_HOST=${{mysql-prod.MYSQL_HOST}}
DB_PORT=${{mysql-prod.MYSQL_PORT}}
DB_USER=${{mysql-prod.MYSQL_USER}}
DB_PASSWORD=${{mysql-prod.MYSQL_PASSWORD}}
DB_NAME=${{mysql-prod.MYSQL_DATABASE}}
```

**Cek 4: Lihat Logs Detail**
- Railway Dashboard → Deployments → View Logs
- Cari error message spesifik

---

### ❌ CORS Error di Frontend

**Update ALLOWED_ORIGINS:**

Railway → Variables → Update:
```bash
# Development (frontend local)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Production (setelah frontend di-deploy)
ALLOWED_ORIGINS=https://your-frontend.vercel.app,https://simnikah-production-xxxx.up.railway.app
```

Save → Auto redeploy

---

### ❌ JWT Token Invalid

**Generate JWT secret baru:**
```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

Update di Railway Variables → `JWT_KEY`

---

### ❌ Port Already in Use

Railway otomatis set `PORT` environment variable.
Aplikasi sudah menggunakan: `os.Getenv("PORT")` dengan fallback 8080.

**Tidak perlu set manual!**

---

## 📊 Railway Free Tier

Railway memberikan **$5 credit per bulan** (gratis)

**Estimasi pemakaian SimNikah:**
- Backend API: ~$3-4/month
- MySQL Database: ~$1-2/month
- **Total: $4-6/month** → Masih dalam free tier! ✅

---

## 🚀 Next Steps

### 1. Custom Domain (Optional)
Railway → Settings → Domains → Add Custom Domain

### 2. Monitoring
Railway menyediakan:
- ✅ Real-time logs
- ✅ CPU/Memory metrics
- ✅ Network usage
- ✅ Deployment history

### 3. Backup Database
Railway otomatis backup MySQL setiap hari

### 4. CI/CD Auto Deploy
Setiap push ke GitHub main branch → Otomatis deploy! ✅

---

## 📞 Butuh Bantuan?

- **Railway Discord**: [discord.gg/railway](https://discord.gg/railway)
- **Railway Docs**: [docs.railway.app](https://docs.railway.app)
- **Railway Status**: [status.railway.app](https://status.railway.app)

---

## 🎯 Quick Reference

### Railway CLI (Optional)

```powershell
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Link project
railway link

# View logs
railway logs

# Open dashboard
railway open
```

### Test Endpoints

```powershell
# Health
curl https://YOUR_APP.up.railway.app/health

# Register
curl -X POST https://YOUR_APP.up.railway.app/register `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"test\",\"email\":\"test@mail.com\",\"password\":\"test123\",\"nama\":\"Test\",\"role\":\"user_biasa\"}'

# Login
curl -X POST https://YOUR_APP.up.railway.app/login `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"test\",\"password\":\"test123\"}'
```

---

**🎉 Selamat! Backend SimNikah kamu sudah LIVE di Railway!**

**Total waktu:** ~30-40 menit  
**Biaya:** $0 (Free tier $5/month)  
**Hasil:** API Production-ready dengan MySQL & HTTPS ✅

---

*Panduan ini dibuat untuk SimNikah API | Last updated: January 2025*


