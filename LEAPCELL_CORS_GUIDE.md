# 🚀 LeapCell Deployment Guide dengan CORS

## 📋 Checklist Deployment

- [ ] Push kode ke GitHub
- [ ] Setup database di LeapCell
- [ ] Set environment variables (termasuk CORS)
- [ ] Deploy aplikasi
- [ ] Test CORS dari frontend
- [ ] Verifikasi health check

---

## Step 1: Persiapan

### 1.1 Push ke GitHub
```bash
cd /home/atho/simpadu
git add .
git commit -m "Add CORS configuration for LeapCell deployment"
git push origin main
```

---

## Step 2: Setup di LeapCell Dashboard

### 2.1 Login & Create Project
1. Buka [LeapCell Dashboard](https://dashboard.leapcell.io)
2. Login dengan akun Anda
3. Klik **"New Project"**
4. Pilih **"Deploy from GitHub"**
5. Connect repository **simpadu**

### 2.2 Create Database
1. Di sidebar, klik **"Database"**
2. Klik **"Create Database"**
3. Pilih **PostgreSQL** (sesuai konfigurasi)
4. Pilih region terdekat
5. Klik **"Create"**
6. **Simpan credentials** (akan digunakan di Step 3)

---

## Step 3: Configure Environment Variables

### 3.1 Di LeapCell Dashboard
1. Pilih project Anda
2. Go to **"Settings"** → **"Environment Variables"**
3. Tambahkan environment variables berikut:

### 🔧 Required Environment Variables

#### Database Configuration
```bash
DB_HOST=<your-leapcell-db-host>
DB_PORT=6438
DB_USER=<your-leapcell-db-user>
DB_PASSWORD=<your-leapcell-db-password>
DB_NAME=<your-leapcell-db-name>
```

#### JWT Configuration
```bash
# Generate dengan: openssl rand -base64 32
JWT_KEY=<your-generated-jwt-key-32-chars-min>
```

#### Server Configuration
```bash
PORT=8080
GIN_MODE=release
```

#### ⭐ CORS Configuration (PENTING!)
```bash
# Ganti dengan domain frontend Anda
# Jika frontend di Vercel:
ALLOWED_ORIGINS=https://your-app.vercel.app

# Jika ada multiple domains (comma-separated):
ALLOWED_ORIGINS=https://your-app.vercel.app,https://your-app.netlify.app,https://custom-domain.com
```

### 3.2 Contoh Lengkap Environment Variables

```bash
# Database (dari LeapCell Database credentials)
DB_HOST=simnikah-db-abc123.leapcell.io
DB_PORT=6438
DB_USER=simnikah_user
DB_PASSWORD=super_secret_password_from_leapcell
DB_NAME=simnikah_db

# JWT (generate sendiri)
JWT_KEY=YourSuperSecretJWTKeyMinimum32CharactersLongForProduction

# Server
PORT=8080
GIN_MODE=release

# CORS - SESUAIKAN DENGAN FRONTEND ANDA!
ALLOWED_ORIGINS=https://simpadu-frontend.vercel.app,https://simpadu.com
```

---

## Step 4: Deploy Application

### 4.1 Deploy via Dashboard
1. Go to **"Deployments"**
2. Klik **"Deploy"**
3. Select branch: **main**
4. Build akan start otomatis

### 4.2 Monitor Deployment
- ✅ Build logs akan muncul real-time
- ✅ Tunggu sampai status: **"Running"**
- ✅ Catat URL deployment: `https://your-app.leapcell.io`

---

## Step 5: Verify Deployment

### 5.1 Test Health Check
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

### 5.2 Test CORS dari Browser

**Open Browser Console (F12) dan jalankan:**
```javascript
fetch('https://your-app.leapcell.io/health', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  },
})
  .then(response => response.json())
  .then(data => console.log('✅ CORS working:', data))
  .catch(error => console.error('❌ CORS error:', error));
```

### 5.3 Check Response Headers
Di Network tab browser, lihat response headers:
```
access-control-allow-origin: https://your-frontend.vercel.app
access-control-allow-methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
access-control-allow-credentials: true
```

---

## Step 6: Update Frontend Configuration

### 6.1 Update API Base URL
**React/Next.js:**
```javascript
// config/api.js
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://your-app.leapcell.io';

export default API_BASE_URL;
```

**Vue/Nuxt:**
```javascript
// nuxt.config.js
export default {
  publicRuntimeConfig: {
    apiURL: process.env.API_URL || 'https://your-app.leapcell.io'
  }
}
```

### 6.2 Set Frontend Environment Variable
**Vercel:**
1. Go to project settings
2. Environment Variables
3. Add: `NEXT_PUBLIC_API_URL=https://your-app.leapcell.io`

**Netlify:**
1. Site settings
2. Build & deploy
3. Environment
4. Add: `VITE_API_URL=https://your-app.leapcell.io`

---

## 🔍 Troubleshooting

### ❌ Problem: CORS Error setelah deploy

**Solusi 1: Check Environment Variable**
```bash
# Di LeapCell Dashboard Logs, cek:
Info: CORS allowed origins: [https://your-frontend.vercel.app]

# Jika muncul default localhost, artinya ALLOWED_ORIGINS tidak terset!
```

**Solusi 2: Update ALLOWED_ORIGINS**
1. Go to Settings → Environment Variables
2. Update `ALLOWED_ORIGINS` dengan domain frontend yang benar
3. **Redeploy aplikasi** (PENTING!)

**Solusi 3: Pastikan HTTPS**
- LeapCell menggunakan HTTPS
- Frontend juga harus HTTPS (Vercel/Netlify otomatis HTTPS)
- Jangan mix HTTP & HTTPS

### ❌ Problem: Database Connection Failed

**Check:**
1. DB credentials benar?
2. DB_PORT = 6438 (LeapCell PostgreSQL default)
3. Database sudah created di LeapCell?

**View Logs:**
```
LeapCell Dashboard → Logs → Application Logs
```

### ❌ Problem: 502 Bad Gateway

**Solusi:**
1. Check health endpoint: `/health`
2. View application logs
3. Pastikan PORT=8080
4. Database migration success?

---

## 📊 Post-Deployment Checklist

- [ ] Health check endpoint working
- [ ] CORS headers present in response
- [ ] Frontend dapat connect ke API
- [ ] Login/Register working
- [ ] Database tables created (auto-migration)
- [ ] JWT authentication working
- [ ] All API endpoints tested

---

## 🔒 Security Best Practices

### 1. Environment Variables
✅ **DO:**
- Set `ALLOWED_ORIGINS` dengan domain spesifik
- Use strong JWT key (32+ characters)
- Keep DB credentials secret

❌ **DON'T:**
- Use wildcard `*` untuk ALLOWED_ORIGINS
- Commit secrets ke Git
- Use default JWT key

### 2. HTTPS
✅ LeapCell provides automatic HTTPS
✅ Use HTTPS URLs di frontend
✅ Set secure cookies if needed

### 3. Database
✅ Enable backup di LeapCell
✅ Monitor database performance
✅ Set proper connection limits

---

## 📈 Monitoring

### Application Logs
```
LeapCell Dashboard → Your Project → Logs
```

Cari log berikut:
```
✅ "Connected to database with SSL"
✅ "Database migration completed successfully"
✅ "Info: CORS allowed origins: [https://...]"
✅ "Server starting on port 8080"
```

### Health Check
Setup monitoring untuk `/health` endpoint:
- LeapCell automatic health checks
- External monitoring (UptimeRobot, etc.)

---

## 🎯 Quick Reference

### Environment Variables Format
```bash
# Single origin
ALLOWED_ORIGINS=https://your-frontend.com

# Multiple origins (NO SPACES!)
ALLOWED_ORIGINS=https://app1.com,https://app2.com,https://app3.com
```

### Common Frontend URLs
- **Vercel**: `https://your-app.vercel.app`
- **Netlify**: `https://your-app.netlify.app`
- **Cloudflare Pages**: `https://your-app.pages.dev`
- **Custom Domain**: `https://your-domain.com`

### Test Commands
```bash
# Health check
curl https://your-app.leapcell.io/health

# Test CORS headers
curl -H "Origin: https://your-frontend.com" \
  -H "Access-Control-Request-Method: POST" \
  -X OPTIONS \
  https://your-app.leapcell.io/login -v

# Test register
curl -X POST https://your-app.leapcell.io/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"test123","nama":"Test User","role":"user_biasa"}'
```

---

## 📞 Support

- **LeapCell Docs**: [docs.leapcell.io](https://docs.leapcell.io)
- **CORS Guide**: `CORS_SETUP.md`
- **API Docs**: `API_DOCUMENTATION.md`

---

## ✅ Success Indicators

Jika semua berhasil, Anda akan melihat:

1. ✅ LeapCell deployment status: **Running**
2. ✅ Health check returns: `{"status":"healthy"}`
3. ✅ Logs show: `"CORS allowed origins: [https://your-frontend.com]"`
4. ✅ Frontend dapat fetch API tanpa CORS error
5. ✅ Login/Register working from frontend
6. ✅ Database connection successful

---

## 🎉 Selesai!

API Anda sekarang sudah live di LeapCell dengan CORS configured!

**Backend URL**: `https://your-app.leapcell.io`

Update URL ini di frontend configuration dan test semua endpoint.

---

*Last updated: 2024 | SimNikah API*

