# 🎓 Penjelasan Mudah Deploy ke Railway

## 🤔 Apa itu Railway?

**Analogi sederhana:**

Bayangkan kamu punya **toko online** (aplikasi backend kamu):

```
┌─────────────────────────────────────────┐
│  TANPA RAILWAY (di laptop kamu)         │
├─────────────────────────────────────────┤
│  Toko kamu di RUMAH                     │
│  ❌ Cuma kamu yang bisa akses            │
│  ❌ Kalau laptop mati, toko tutup        │
│  ❌ Cuma bisa buka kalau laptop nyala    │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  DENGAN RAILWAY (di internet)           │
├─────────────────────────────────────────┤
│  Toko kamu di MALL                      │
│  ✅ Semua orang bisa akses               │
│  ✅ Buka 24/7 (selalu online)            │
│  ✅ Punya alamat (URL) yang bisa dibuka  │
└─────────────────────────────────────────┘
```

**Railway = Tempat parkir aplikasi kamu di internet, gratis!**

---

## 📦 Apa yang Kamu Punya Sekarang?

Kamu punya **kode aplikasi** di folder `/home/atho/simpadu`:

```
simpadu/
├── main.go              ← Program utama
├── catin/               ← Fitur pendaftaran nikah
├── helper/              ← Helper functions
├── structs/             ← Database models
├── go.mod               ← Daftar library yang dipakai
└── env.example          ← Contoh konfigurasi
```

**Ini seperti:** Kamu punya toko, tapi masih di rumah. Belum dibuka ke publik.

---

## 🎯 Tujuan Deploy ke Railway

**Dari:**
```
http://localhost:8080  ← Cuma di komputer kamu
```

**Jadi:**
```
https://simnikah-production-xxxx.up.railway.app  ← Bisa diakses dari mana saja!
```

---

## 🚀 STEP-BY-STEP (Bahasa Manusia)

### STEP 1: Taruh Kode di GitHub (5 menit)

**Analogi:** 
- GitHub = Google Drive untuk kode
- Kamu upload kode kamu ke GitHub biar Railway bisa ambil

**Caranya:**

```bash
# 1. Masuk ke folder project
cd /home/atho/simpadu

# 2. Siapkan folder untuk upload
git init

# 3. Daftar di GitHub.com (kalau belum punya akun)
# Buka: https://github.com/signup

# 4. Buat repository baru di GitHub
# Buka: https://github.com/new
# Nama: simnikah-backend
# Klik: Create repository

# 5. Upload kode ke GitHub
git remote add origin https://github.com/USERNAME_KAMU/simnikah-backend.git
git add .
git commit -m "Upload project"
git push -u origin main
```

**✅ Berhasil kalau:**
- Kamu buka `https://github.com/USERNAME_KAMU/simnikah-backend`
- Kamu lihat semua file kamu ada di sana

---

### STEP 2: Daftar di Railway (2 menit)

**Analogi:**
- Railway = Tempat hosting gratis (seperti sewa tempat di mall, tapi gratis!)

**Caranya:**

1. **Buka:** [railway.app](https://railway.app)
2. **Klik:** "Login with GitHub"
3. **Klik:** "Authorize Railway" (kasih izin Railway akses GitHub kamu)
4. **Selesai!** Kamu punya akun Railway gratis ($5 credit per bulan)

**✅ Berhasil kalau:**
- Kamu masuk ke Railway Dashboard
- Kamu lihat tombol "New Project"

---

### STEP 3: Buat Project di Railway (3 menit)

**Analogi:**
- Project = Toko kamu yang mau dibuka di mall

**Caranya:**

1. **Klik:** "New Project"
2. **Pilih:** "Deploy from GitHub repo"
3. **Pilih:** Repository `simnikah-backend` (yang tadi kamu buat)
4. **Klik:** Deploy

Railway akan mulai **build** (seperti renovasi toko sebelum buka).

**⏳ Tunggu 2-3 menit...**

**Hasilnya:**
- ❌ Build pertama akan **GAGAL** (ini NORMAL!)
- Kenapa? Karena aplikasi kamu butuh **database MySQL**
- Railway belum tau kamu butuh database

---

### STEP 4: Kasih Database MySQL (3 menit)

**Analogi:**
- Database = Gudang untuk simpan data
- Aplikasi kamu butuh tempat simpan data nikah, user, dll

**Caranya:**

1. **Di Railway Dashboard, klik:** tombol "+ New" (pojok kanan atas)
2. **Pilih:** "Database"
3. **Pilih:** "Add MySQL"
4. **Tunggu 20 detik...**
5. **Selesai!** MySQL ready

**Yang terjadi di balik layar:**
```
Railway otomatis bikin:
- Database server MySQL ✅
- Username & password ✅
- Database kosong siap pakai ✅
```

**✅ Berhasil kalau:**
- Kamu lihat ada 2 kotak di Railway:
  - 1 kotak = Aplikasi kamu
  - 1 kotak = MySQL database

---

### STEP 5: Kasih Konfigurasi (5 menit)

**Analogi:**
- Aplikasi kamu perlu "petunjuk" untuk connect ke database
- Seperti kasih alamat gudang ke pegawai toko

**Caranya:**

1. **Klik:** Kotak aplikasi kamu (bukan kotak MySQL)
2. **Pilih tab:** "Variables"
3. **Klik:** "RAW Editor"
4. **Buka terminal baru**, jalankan:
   ```bash
   openssl rand -base64 32
   ```
   Copy hasil (contoh: `Xy9pQ2vN8kL5mT7rW3sH6jF1dA4gC9bE8zX2vN5mT0=`)

5. **Paste kode ini di RAW Editor:**

```bash
DB_HOST=${{MySQL.MYSQL_HOST}}
DB_PORT=${{MySQL.MYSQL_PORT}}
DB_USER=${{MySQL.MYSQL_USER}}
DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
DB_NAME=${{MySQL.MYSQL_DATABASE}}
JWT_KEY=PASTE_HASIL_OPENSSL_DI_SINI
PORT=8080
GIN_MODE=release
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

6. **Ganti** `PASTE_HASIL_OPENSSL_DI_SINI` dengan hasil copy tadi
7. **Klik:** "Add Variables"

**⚠️ PENTING:** 
- Perhatikan huruf besar di `MySQL` (bukan `mysql`)
- Syntax `${{MySQL.MYSQL_HOST}}` artinya: "Railway, tolong isi otomatis dari MySQL service"

**Yang terjadi:**
```
Railway akan:
1. Otomatis ambil password MySQL ✅
2. Otomatis connect aplikasi ke database ✅
3. Otomatis redeploy (build ulang) ✅
```

**⏳ Tunggu 2-5 menit... Railway sedang build ulang**

---

### STEP 6: Cek Apakah Berhasil (2 menit)

**Caranya:**

1. **Pilih tab:** "Deployments"
2. **Lihat status:**
   - 🔄 Building = Masih proses
   - ✅ Success = BERHASIL!
   - ❌ Failed = Ada error

3. **Klik deployment terakhir**
4. **Scroll ke bawah**, cari log seperti ini:

```
✅ Connected to MySQL database successfully
✅ Database migration completed successfully
✅ Server starting on port 8080
```

**✅ Kalau kamu lihat log di atas = BERHASIL!**

---

### STEP 7: Dapatkan URL (2 menit)

**Analogi:**
- Sekarang toko kamu sudah buka, tapi belum ada alamat
- Kita perlu kasih alamat biar orang bisa datang

**Caranya:**

1. **Klik:** Kotak aplikasi kamu
2. **Pilih tab:** "Settings"
3. **Scroll ke:** "Domains"
4. **Klik:** "Generate Domain"
5. **Railway kasih URL:** `https://simnikah-production-xxxx.up.railway.app`
6. **COPY URL INI!** 📋

---

### STEP 8: TEST! (3 menit)

**Buka browser atau terminal:**

```bash
# Test 1: Health check
curl https://simnikah-production-xxxx.up.railway.app/health

# Harusnya muncul:
{
  "status": "healthy",
  "message": "SimNikah API is running"
}
```

**ATAU buka di browser:**
```
https://simnikah-production-xxxx.up.railway.app/health
```

**✅ BERHASIL kalau:**
- Muncul JSON response seperti di atas
- Tidak ada error

---

## 🎉 SELESAI!

**Sekarang aplikasi kamu:**
- ✅ Online 24/7
- ✅ Bisa diakses dari mana saja
- ✅ Punya URL sendiri
- ✅ Database MySQL ready
- ✅ HTTPS (aman)

**URL aplikasi kamu:**
```
https://simnikah-production-xxxx.up.railway.app
```

---

## 🤔 Pertanyaan Umum

### Q: Kenapa harus ke GitHub dulu?
**A:** Railway ambil kode dari GitHub. Seperti tukang bangunan perlu lihat blueprint sebelum bangun.

### Q: Kenapa build pertama gagal?
**A:** Normal! Aplikasi kamu butuh database. Setelah kasih MySQL, build kedua akan sukses.

### Q: Apa itu `${{MySQL.MYSQL_HOST}}`?
**A:** Ini cara bilang ke Railway: "Tolong isi otomatis dengan IP MySQL kamu."

Railway tau MySQL password, tapi kamu tidak (dan tidak perlu tau).

### Q: Gratis berapa lama?
**A:** Railway kasih $5 credit per bulan GRATIS selamanya. Aplikasi kamu cuma pakai ~$4-6/bulan = masih gratis!

### Q: Gimana update kode?
**A:** Tinggal push ke GitHub:
```bash
git add .
git commit -m "Update"
git push
```
Railway otomatis deploy ulang!

---

## 🆘 Kalau Ada Masalah

### ❌ Build Failed

**Lihat log error di Railway:**
1. Deployments → Click deployment → Scroll ke bawah
2. Cari error berwarna merah
3. Biasanya:
   - Database not connected → Cek environment variables
   - Module not found → Push ulang ke GitHub

### ❌ Database Connection Error

**Cek:**
1. MySQL service sudah running? (lihat kotak MySQL di Railway)
2. Environment variables benar? (cek syntax `${{MySQL.MYSQL_HOST}}`)

### ❌ Masih Bingung?

**Hubungi:**
- Railway Discord: [discord.gg/railway](https://discord.gg/railway)
- Atau screenshot error kirim ke saya

---

## 📊 Rangkuman Visual

```
┌─────────────┐
│ KODE KAMU   │
│ (di laptop) │
└──────┬──────┘
       │ git push
       ▼
┌─────────────┐
│   GITHUB    │
│ (simpan kode)│
└──────┬──────┘
       │ Railway ambil kode
       ▼
┌─────────────┐         ┌─────────────┐
│   RAILWAY   │────────▶│   MySQL     │
│ (run app)   │         │ (database)  │
└──────┬──────┘         └─────────────┘
       │
       ▼
   INTERNET
   (semua orang bisa akses)
   
   https://simnikah-xxxx.up.railway.app
```

---

## ✅ Checklist Sederhana

Ikuti step ini satu-satu:

- [ ] Daftar GitHub (kalau belum)
- [ ] Upload kode ke GitHub
- [ ] Daftar Railway (pakai akun GitHub)
- [ ] Buat project di Railway dari GitHub
- [ ] Tambah MySQL database
- [ ] Generate JWT secret (`openssl rand -base64 32`)
- [ ] Set environment variables
- [ ] Tunggu deploy sukses
- [ ] Generate domain
- [ ] Test URL di browser

**Total waktu: ~30 menit**

---

## 🎯 Mau Coba Sekarang?

**Copy-paste command ini satu per satu:**

```bash
# 1. Masuk ke project
cd /home/atho/simpadu

# 2. Git init (kalau belum)
git init

# 3. Ganti USERNAME_KAMU dengan username GitHub kamu
git remote add origin https://github.com/USERNAME_KAMU/simnikah-backend.git

# 4. Upload
git add .
git commit -m "Deploy to Railway"
git push -u origin main

# 5. Generate JWT secret
openssl rand -base64 32

# COPY hasilnya, nanti pakai di Railway Variables
```

**Setelah itu:**
1. Buka [railway.app](https://railway.app)
2. Login dengan GitHub
3. Ikuti step 3-8 di atas

---

**Masih bingung step mana? Bilang aja step berapa yang belum paham! 😊**

