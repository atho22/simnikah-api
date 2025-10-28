# 🚂 Railway Deployment Guide - SimNikah API

## 📋 Prerequisites

1. **Railway Account** - Daftar di [Railway.app](https://railway.app)
2. **GitHub Repository** - Push project ke GitHub
3. **MySQL Database** - Railway menyediakan MySQL add-on

---

## 🚀 Step-by-Step Deployment

### Step 1: Persiapan Repository

#### 1.1 Push ke GitHub
```bash
cd /home/atho/simpadu
git add .
git commit -m "Switch to Railway with MySQL"
git push origin main
```

#### 1.2 Pastikan File Konfigurasi Ada
- ✅ `Dockerfile`
- ✅ `railway.json`
- ✅ `nixpacks.toml`
- ✅ `env.example`

---

### Step 2: Setup di Railway

#### 2.1 Login & Create New Project
1. Buka [Railway Dashboard](https://railway.app/dashboard)
2. Click **"New Project"**
3. Pilih **"Deploy from GitHub repo"**
4. Authorize Railway untuk akses GitHub
5. Pilih repository **simpadu**

#### 2.2 Tambah MySQL Database
1. Di project Anda, click **"+ New"**
2. Pilih **"Database"**
3. Pilih **"Add MySQL"**
4. Railway akan otomatis provision MySQL instance
5. MySQL akan tersedia dengan environment variables otomatis

---

### Step 3: Configure Environment Variables

#### 3.1 Railway Automatic Variables (MySQL)
Railway otomatis menyediakan:
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DATABASE`
- `MYSQL_URL` (connection string lengkap)

#### 3.2 Manual Environment Variables yang Perlu Ditambahkan

Di Railway Dashboard → Your Service → Variables:

**Click "RAW Editor" dan paste:**

```bash
# Database Configuration (Railway menyediakan otomatis)
DB_HOST=${{MYSQL.MYSQL_HOST}}
DB_PORT=${{MYSQL.MYSQL_PORT}}
DB_USER=${{MYSQL.MYSQL_USER}}
DB_PASSWORD=${{MYSQL.MYSQL_PASSWORD}}
DB_NAME=${{MYSQL.MYSQL_DATABASE}}

# JWT Configuration (GENERATE SENDIRI!)
JWT_KEY=YourSuperSecretJWTKeyMinimum32CharactersLong

# Server Configuration
PORT=8080
GIN_MODE=release

# CORS Configuration
# Untuk development (frontend local)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Untuk production (ganti dengan domain frontend Anda)
# ALLOWED_ORIGINS=https://your-frontend.vercel.app,https://your-app.netlify.app
```

#### 3.3 Generate JWT Secret
```bash
# Generate random JWT secret (32+ characters)
openssl rand -base64 32
# Copy hasil dan paste ke JWT_KEY
```

---

### Step 4: Deploy Application

#### 4.1 Deploy Otomatis
- Railway akan otomatis build & deploy setelah push ke GitHub
- Monitor di **"Deployments"** tab
- Tunggu sampai status: ✅ **"Success"**

#### 4.2 Generate Domain
1. Di service settings, click **"Generate Domain"**
2. Railway akan berikan domain: `https://your-app.up.railway.app`
3. Atau tambahkan custom domain

---

### Step 5: Verify Deployment

#### 5.1 Check Deployment Logs
```
Railway Dashboard → Your Service → Deployments → View Logs
```

Cari log berikut:
```
✅ Connected to MySQL database successfully
✅ Database migration completed successfully
✅ Info: CORS allowed origins: [...]
✅ Server starting on port 8080
```

#### 5.2 Test Health Endpoint
```bash
curl https://your-app.up.railway.app/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 5.3 Test dari Browser
```javascript
fetch('https://your-app.up.railway.app/health')
  .then(res => res.json())
  .then(data => console.log('✅ API Connected:', data))
  .catch(err => console.error('❌ Error:', err));
```

---

### Step 6: Update Frontend Configuration

#### 6.1 Update API Base URL

**React/Vite (.env.local):**
```bash
VITE_API_URL=https://your-app.up.railway.app
```

**Next.js (.env.local):**
```bash
NEXT_PUBLIC_API_URL=https://your-app.up.railway.app
```

#### 6.2 Update ALLOWED_ORIGINS
Setelah deploy frontend:
```bash
# Railway Variables
ALLOWED_ORIGINS=https://your-frontend.vercel.app,https://your-app.up.railway.app
```

---

## 🔧 Railway Configuration Files

### railway.json
```json
{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "NIXPACKS",
    "buildCommand": "go build -o main ."
  },
  "deploy": {
    "startCommand": "./main",
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 10
  }
}
```

### nixpacks.toml
```toml
[phases.setup]
nixPkgs = ["go_1_23"]

[phases.build]
cmds = ["go build -o main ."]

[phases.start]
cmd = "./main"
```

---

## 🗄️ Database Management

### Connect to MySQL
```bash
# Railway CLI
railway connect mysql

# Atau manual dengan credentials dari Railway
mysql -h <MYSQL_HOST> -P <MYSQL_PORT> -u <MYSQL_USER> -p<MYSQL_PASSWORD> <MYSQL_DATABASE>
```

### Auto Migration
Database tables akan dibuat otomatis saat aplikasi pertama kali dijalankan melalui GORM AutoMigrate.

### Manual Migration (jika perlu)
```sql
-- Tables akan dibuat otomatis, tapi jika perlu manual:
SHOW TABLES;
DESCRIBE users;
DESCRIBE pendaftaran_nikahs;
```

---

## 🔍 Troubleshooting

### ❌ Build Failed

**Check:**
1. Go version correct? (1.23.6)
2. Dependencies up to date?
```bash
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
git push
```

### ❌ Database Connection Failed

**Check Railway Variables:**
```bash
# Pastikan variable reference benar:
DB_HOST=${{MYSQL.MYSQL_HOST}}

# Bukan:
DB_HOST=${MYSQL_HOST}  # ❌ SALAH!
```

**View Logs:**
```
Railway Dashboard → Deployments → View Logs
```

### ❌ CORS Error

**Update ALLOWED_ORIGINS:**
```bash
# Frontend local
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Frontend production
ALLOWED_ORIGINS=https://your-frontend.vercel.app
```

**Redeploy setelah update:**
- Railway akan auto redeploy saat environment variables berubah
- Atau manual trigger: Deployments → Redeploy

### ❌ Port Already in Use

Railway otomatis set `PORT` environment variable.
Aplikasi sudah menggunakan: `os.Getenv("PORT")` dengan fallback 8080.

---

## 📊 Railway Features

### ✅ Automatic HTTPS
Railway menyediakan SSL certificate otomatis untuk semua deployments.

### ✅ Auto Deployments
Push ke GitHub → Otomatis deploy di Railway.

### ✅ Environment Variables
Railway menyediakan variable reference untuk MySQL:
```bash
${{MYSQL.MYSQL_HOST}}
${{MYSQL.MYSQL_PORT}}
${{MYSQL.MYSQL_USER}}
${{MYSQL.MYSQL_PASSWORD}}
${{MYSQL.MYSQL_DATABASE}}
```

### ✅ Monitoring
- View logs real-time
- Metrics (CPU, Memory, Network)
- Deployment history

### ✅ Scaling
- Vertical scaling otomatis
- Horizontal scaling available (paid plans)

---

## 💰 Pricing

### Free Tier (Hobby Plan)
- ✅ $5 credit per month
- ✅ Enough untuk small projects
- ✅ Auto sleep setelah inactivity
- ✅ 1 concurrent deployment

### Paid Plans
- 💳 Developer: $20/month
- 💳 Team: $100/month
- 💳 Custom pricing untuk enterprise

**Estimasi biaya SimNikah API:**
- API Service: ~$3-5/month
- MySQL Database: ~$2-3/month
- **Total: ~$5-8/month** (masih dalam free tier!)

---

## 🔐 Security Best Practices

### 1. Environment Variables
✅ **DO:**
- Use Railway variable references untuk database
- Generate strong JWT key (32+ chars)
- Set specific `ALLOWED_ORIGINS`

❌ **DON'T:**
- Hardcode credentials
- Use wildcard `*` untuk CORS
- Commit secrets to Git

### 2. Database
✅ **DO:**
- Use Railway's automatic backups
- Monitor database metrics
- Set proper indexes

### 3. API Security
✅ **DO:**
- Use HTTPS (Railway default)
- Implement rate limiting (future)
- Validate all inputs

---

## 📈 Post-Deployment Checklist

- [ ] Health check working
- [ ] Database connected & migrated
- [ ] CORS headers correct
- [ ] Frontend can connect
- [ ] Login/Register working
- [ ] All endpoints tested
- [ ] Environment variables set
- [ ] Custom domain configured (optional)
- [ ] Monitoring enabled
- [ ] Backups configured

---

## 🎯 Quick Commands

### Railway CLI Installation
```bash
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

# Run command in Railway environment
railway run npm start
```

### Test Endpoints
```bash
# Health check
curl https://your-app.up.railway.app/health

# Register
curl -X POST https://your-app.up.railway.app/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123",
    "nama": "Test User",
    "role": "user_biasa"
  }'

# Login
curl -X POST https://your-app.up.railway.app/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123"
  }'
```

---

## 🆚 Railway vs LeapCell

| Feature | Railway | LeapCell |
|---------|---------|----------|
| **Database** | MySQL, PostgreSQL, Redis, MongoDB | MySQL, PostgreSQL |
| **Free Tier** | $5/month credit | Limited free tier |
| **Deployment** | GitHub auto-deploy | GitHub auto-deploy |
| **Custom Domain** | ✅ Free | ✅ Free |
| **SSL** | ✅ Automatic | ✅ Automatic |
| **CLI** | ✅ Excellent | ❌ Limited |
| **Logs** | ✅ Real-time | ✅ Real-time |
| **Scaling** | ✅ Easy | ✅ Easy |
| **Pricing** | Pay-as-you-go | Fixed tiers |

---

## 📞 Support & Resources

- **Railway Docs**: [docs.railway.app](https://docs.railway.app)
- **Railway Discord**: [discord.gg/railway](https://discord.gg/railway)
- **Railway Blog**: [blog.railway.app](https://blog.railway.app)
- **API Docs**: [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
- **CORS Guide**: [CORS_SETUP.md](CORS_SETUP.md)

---

## ✅ Success Indicators

Jika deployment berhasil, Anda akan melihat:

1. ✅ Railway deployment status: **"Success"**
2. ✅ Logs show: `"Connected to MySQL database successfully"`
3. ✅ Logs show: `"Database migration completed successfully"`
4. ✅ Health endpoint returns: `{"status":"healthy"}`
5. ✅ Frontend dapat connect tanpa CORS error
6. ✅ Login/Register working
7. ✅ Custom domain working (jika dikonfigurasi)

---

## 🎉 Selesai!

API Anda sekarang sudah live di Railway dengan MySQL!

**Backend URL**: `https://your-app.up.railway.app`

Update URL ini di frontend dan test semua endpoint.

---

## 📝 Migration from LeapCell

Jika sebelumnya deploy di LeapCell:

### Data Migration
1. Export data dari LeapCell PostgreSQL:
```bash
pg_dump -h leapcell-host -U user -d simnikah > backup.sql
```

2. Convert PostgreSQL dump to MySQL compatible (manual)

3. Import ke Railway MySQL:
```bash
mysql -h railway-mysql-host -u user -p database < backup.sql
```

### DNS Update
- Update DNS records dari LeapCell ke Railway
- Update frontend API URL
- Test thoroughly

---

*Last updated: 2024 | SimNikah API on Railway*

