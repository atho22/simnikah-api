# ⚡ Quick Start - Deploy ke Railway (5 Langkah)

## 🎯 Langkah Cepat

### 1️⃣ Push ke GitHub
```bash
git add .
git commit -m "Ready for Railway"
git push origin main
```

### 2️⃣ Buat Project di Railway
- Buka [railway.app](https://railway.app)
- Login dengan GitHub
- New Project → Deploy from GitHub → Pilih repo kamu

### 3️⃣ Tambah MySQL
- Di Railway project → "+ New" → Database → Add MySQL
- **Catat nama service MySQL!** (misal: `MySQL` atau `mysql`)

### 4️⃣ Generate JWT Secret
```powershell
# Jalankan script
.\scripts\generate-jwt-key.ps1

# Atau manual
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

### 5️⃣ Set Environment Variables

1. Railway → Your Service → Variables → RAW Editor
2. Buka file `railway-env-template.txt`
3. Copy semua isinya
4. Paste di RAW Editor
5. **GANTI 2 hal:**
   - `YOUR_JWT_SECRET_HERE` → Paste JWT secret dari step 4
   - `MySQL` → Ganti dengan nama service MySQL dari step 3
6. Click "Update Variables"

### 6️⃣ Tunggu Deploy & Test

- Tunggu deployment SUCCESS (2-5 menit)
- Generate domain di Settings → Domains
- Test: `curl https://your-app.up.railway.app/health`

---

## 📚 Dokumentasi Lengkap

Lihat **[RAILWAY_SETUP_GUIDE.md](RAILWAY_SETUP_GUIDE.md)** untuk panduan detail!

---

## 🆘 Troubleshooting Cepat

### Database Connection Failed?
- Cek nama service MySQL di Railway
- Pastikan variable menggunakan: `${{MySQL.MYSQL_HOST}}` (bukan `${MYSQL_HOST}`)

### Build Failed?
```bash
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
git push
```

### CORS Error?
Update `ALLOWED_ORIGINS` di Railway Variables dengan domain frontend kamu.

---

**Butuh bantuan lebih detail?** Baca **[RAILWAY_SETUP_GUIDE.md](RAILWAY_SETUP_GUIDE.md)** 📖


