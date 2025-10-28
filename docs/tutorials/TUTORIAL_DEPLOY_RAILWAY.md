# 🚀 Tutorial Deploy SimNikah ke Railway (Langkah demi Langkah)

## 📌 Persiapan (5 menit)

### ✅ Yang Kamu Butuhkan:
- [ ] Akun GitHub (gratis)
- [ ] Akun Railway (gratis - sign up dengan GitHub)
- [ ] Project SimNikah sudah siap di local

---

## 🎯 STEP 1: Push Project ke GitHub (5 menit)

### 1.1 Buat Repository di GitHub

1. Buka [github.com/new](https://github.com/new)
2. Isi nama repository: **`simnikah-backend`**
3. **Jangan centang** "Add README" (karena kita sudah punya)
4. Click **"Create repository"**

### 1.2 Push dari Local ke GitHub

```bash
cd /home/atho/simpadu

# Init git (jika belum)
git init

# Add remote
git remote add origin https://github.com/YOUR_USERNAME/simnikah-backend.git

# Add semua file
git add .

# Commit
git commit -m "Initial commit: SimNikah API with MySQL"

# Push
git push -u origin main
```

> 💡 **Tips:** Ganti `YOUR_USERNAME` dengan username GitHub kamu!

**✅ Berhasil jika:** Repository kamu di GitHub sudah ada file-file project

---

## 🎯 STEP 2: Setup Railway Account (2 menit)

### 2.1 Daftar Railway

1. Buka [railway.app](https://railway.app)
2. Click **"Start a New Project"** atau **"Login"**
3. Pilih **"Login with GitHub"**
4. Authorize Railway untuk akses GitHub kamu
5. Selesai! Kamu punya $5 credit gratis setiap bulan 🎉

**✅ Berhasil jika:** Kamu sudah masuk ke Railway Dashboard

---

## 🎯 STEP 3: Deploy dari GitHub (10 menit)

### 3.1 Create New Project

1. Di Railway Dashboard, click **"New Project"**
2. Pilih **"Deploy from GitHub repo"**
3. Pilih repository **`simnikah-backend`**
4. Railway akan otomatis mulai build

**⏳ Tunggu sekitar 2-3 menit...**

### 3.2 Build Pertama (Akan GAGAL - Normal!)

❌ Build pertama akan **GAGAL** karena belum ada database

```
Error: failed to connect to MySQL database
```

**Ini NORMAL!** Kita akan fix di step berikutnya 👇

---

## 🎯 STEP 4: Tambah MySQL Database (3 menit)

### 4.1 Add MySQL Service

1. Di project Railway kamu, click tombol **"+ New"** (pojok kanan atas)
2. Pilih **"Database"**
3. Pilih **"Add MySQL"**
4. Tunggu 10-20 detik...
5. MySQL siap! ✅

### 4.2 Link Database ke App

Railway otomatis akan link MySQL ke aplikasi kamu.

**Kamu tidak perlu set MYSQL_HOST, MYSQL_PASSWORD, dll secara manual!**

Railway otomatis menyediakan:
- ✅ `MYSQL_HOST`
- ✅ `MYSQL_PORT`
- ✅ `MYSQL_USER`
- ✅ `MYSQL_PASSWORD`
- ✅ `MYSQL_DATABASE`

---

## 🎯 STEP 5: Set Environment Variables (5 menit)

### 5.1 Buka Variables Settings

1. Click **service aplikasi kamu** (bukan MySQL)
2. Pilih tab **"Variables"**
3. Click **"RAW Editor"**

### 5.2 Generate JWT Secret

Buka terminal dan jalankan:

```bash
openssl rand -base64 32
```

Copy hasilnya (contoh: `Xy9pQ2vN8kL5mT7rW3sH6jF1dA4gC9bE8zX2vN5mT0=`)

### 5.3 Paste Variables

Paste kode ini di RAW Editor, lalu **ganti JWT_KEY** dengan hasil generate di atas:

```bash
# Database Configuration (Railway auto-inject dari MySQL service)
DB_HOST=${{MySQL.MYSQL_HOST}}
DB_PORT=${{MySQL.MYSQL_PORT}}
DB_USER=${{MySQL.MYSQL_USER}}
DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
DB_NAME=${{MySQL.MYSQL_DATABASE}}

# JWT Secret (GANTI INI dengan hasil generate!)
JWT_KEY=Xy9pQ2vN8kL5mT7rW3sH6jF1dA4gC9bE8zX2vN5mT0=

# Server Configuration
PORT=8080
GIN_MODE=release

# CORS Configuration (untuk frontend local dulu)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

> ⚠️ **PENTING:**
> - Perhatikan syntax `${{MySQL.MYSQL_HOST}}` dengan huruf besar di "MySQL"!
> - Sesuaikan dengan nama service MySQL kamu (lihat di Railway dashboard)
> - Jika nama service MySQL kamu berbeda, ganti "MySQL" dengan nama yang benar

### 5.4 Save

Click **"Update Variables"** di pojok kanan atas

Railway akan **otomatis redeploy** setelah variables diupdate!

---

## 🎯 STEP 6: Monitor Deployment (3 menit)

### 6.1 Lihat Logs

1. Pilih tab **"Deployments"**
2. Click deployment yang sedang running
3. Scroll ke bawah untuk lihat logs

### 6.2 Tunggu Success

Tunggu sampai kamu lihat log seperti ini:

```
✅ Connecting to MySQL: host=...
✅ Connected to MySQL database successfully
✅ Database migration completed successfully
✅ Info: CORS allowed origins: [...]
✅ Server starting on port 8080
```

**⏳ Biasanya 2-5 menit**

### 6.3 Cek Status

Status deployment harus: **✅ SUCCESS** (warna hijau)

Jika masih **🔄 BUILDING** → Tunggu sebentar
Jika **❌ FAILED** → Lihat section Troubleshooting di bawah

---

## 🎯 STEP 7: Generate Domain & Test (5 menit)

### 7.1 Generate Public Domain

1. Pilih tab **"Settings"**
2. Scroll ke **"Domains"**
3. Click **"Generate Domain"**
4. Railway akan beri URL seperti: `https://simnikah-production-xxxx.up.railway.app`
5. **COPY URL INI!** 📋

### 7.2 Test Health Check

Buka browser atau terminal, test:

```bash
curl https://simnikah-production-xxxx.up.railway.app/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "message": "SimNikah API is running",
  "timestamp": "2025-01-27T..."
}
```

**✅ BERHASIL!** API kamu sudah LIVE! 🎉

### 7.3 Test Register User

```bash
curl -X POST https://simnikah-production-xxxx.up.railway.app/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123",
    "nama": "Test User",
    "no_hp": "08123456789",
    "role": "user_biasa"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "User berhasil didaftarkan"
}
```

---

## 🎯 STEP 8: Update Frontend (2 menit)

### 8.1 Update API URL di Frontend

Jika kamu punya frontend (React/Next.js), update file `.env`:

**React/Vite (.env.local):**
```bash
VITE_API_URL=https://simnikah-production-xxxx.up.railway.app
```

**Next.js (.env.local):**
```bash
NEXT_PUBLIC_API_URL=https://simnikah-production-xxxx.up.railway.app
```

### 8.2 Update CORS di Railway

Setelah frontend di-deploy (misal di Vercel):

1. Kembali ke Railway → Variables
2. Update `ALLOWED_ORIGINS`:

```bash
ALLOWED_ORIGINS=https://your-frontend.vercel.app,https://simnikah-production-xxxx.up.railway.app
```

3. Save → Otomatis redeploy

---

## 🎉 SELESAI! Backend Sudah LIVE!

### ✅ Checklist Final:

- [x] ✅ Project di GitHub
- [x] ✅ Railway project created
- [x] ✅ MySQL database added
- [x] ✅ Environment variables set
- [x] ✅ Deployment SUCCESS
- [x] ✅ Public domain generated
- [x] ✅ Health check working
- [x] ✅ API endpoint tested

### 🌐 URLs Kamu:

- **Backend API**: `https://simnikah-production-xxxx.up.railway.app`
- **Railway Dashboard**: `https://railway.app/project/YOUR_PROJECT_ID`

---

## 🐛 Troubleshooting

### ❌ Build Failed: "cannot find module"

**Solusi:**
```bash
cd /home/atho/simpadu
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
git push
```

Railway akan auto redeploy setelah push.

---

### ❌ Database Connection Failed

**Cek:**
1. MySQL service sudah running? (lihat di Railway dashboard)
2. Variable reference benar?

**Pastikan:**
```bash
DB_HOST=${{MySQL.MYSQL_HOST}}  # ✅ Benar
DB_HOST=${MYSQL_HOST}          # ❌ SALAH!
```

**Jika nama MySQL service kamu bukan "MySQL":**

Lihat nama service di Railway (misal: "mysql-prod"), lalu:
```bash
DB_HOST=${{mysql-prod.MYSQL_HOST}}
```

**Test Manual:**
1. Railway → MySQL service → Connect
2. Copy connection command
3. Test koneksi dari local

---

### ❌ CORS Error di Frontend

**Update ALLOWED_ORIGINS:**

Railway → Variables → Update:
```bash
# Development
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Production
ALLOWED_ORIGINS=https://your-frontend.vercel.app
```

Save → Auto redeploy

---

### ❌ JWT Token Invalid

**Generate JWT secret baru:**
```bash
openssl rand -base64 32
```

Update di Railway Variables → `JWT_KEY`

---

### 🔍 Lihat Logs Detail

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Link project
railway link

# Lihat logs
railway logs
```

---

## 📊 Railway Free Tier Limits

Railway memberikan **$5 credit per bulan** (gratis)

**Estimasi pemakaian SimNikah:**
- Backend API: ~$3-4/month
- MySQL Database: ~$1-2/month
- **Total: $4-6/month** → Masih dalam free tier! ✅

**Jika mau hemat:**
- Sleep service saat tidak dipakai
- Optimize query database
- Enable caching

---

## 🚀 Next Steps

### 1. Custom Domain (Optional)

Railway → Settings → Domains → Add Custom Domain

```
api.your-domain.com → Point CNAME ke Railway
```

### 2. Monitoring

Railway menyediakan:
- ✅ Real-time logs
- ✅ CPU/Memory metrics
- ✅ Network usage
- ✅ Deployment history

### 3. Backup Database

Railway otomatis backup MySQL setiap hari

Manual backup:
```bash
railway connect mysql

# Dump database
mysqldump -u root -p simnikah > backup.sql
```

### 4. CI/CD Auto Deploy

Setiap kali push ke GitHub main branch → Otomatis deploy! ✅

Test:
```bash
# Edit file
echo "// Updated" >> main.go

# Push
git add .
git commit -m "Update"
git push

# Railway otomatis deploy!
```

---

## 📚 Dokumentasi Lengkap

- **[Railway Deployment Guide](RAILWAY_DEPLOYMENT.md)** - Dokumentasi lengkap
- **[API Documentation](API_DOCUMENTATION.md)** - 50+ endpoints
- **[Map Integration](MAP_INTEGRATION.md)** - OpenStreetMap guide
- **[CORS Setup](CORS_SETUP.md)** - CORS troubleshooting

---

## 💬 Butuh Bantuan?

- **Railway Discord**: [discord.gg/railway](https://discord.gg/railway)
- **Railway Docs**: [docs.railway.app](https://docs.railway.app)
- **GitHub Issues**: Buat issue di repository

---

## 🎯 Quick Reference

### Railway CLI Commands

```bash
# Install
npm i -g @railway/cli

# Login
railway login

# Link project
railway link

# Logs
railway logs

# Shell access
railway shell

# Open dashboard
railway open

# Deploy manually
railway up
```

### Test Endpoints

```bash
# Health
curl https://YOUR_APP.up.railway.app/health

# Register
curl -X POST https://YOUR_APP.up.railway.app/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","email":"test@mail.com","nama":"Test","role":"user_biasa"}'

# Login
curl -X POST https://YOUR_APP.up.railway.app/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

---

**🎉 Selamat! Backend SimNikah kamu sudah LIVE di Railway!**

**Total waktu:** ~30-40 menit
**Biaya:** $0 (Free tier $5/month)
**Hasil:** API Production-ready dengan MySQL & HTTPS ✅

---

*Tutorial ini dibuat untuk SimNikah API | Last updated: January 2025*


