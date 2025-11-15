# 🌱 Seeder Kepala KUA

Dokumentasi untuk seeder user Kepala KUA yang akan dibuat otomatis saat aplikasi pertama kali dijalankan.

---

## 📋 Overview

Seeder ini akan membuat user **Kepala KUA** dengan akses penuh ke sistem. User ini dapat:
- Membuat staff baru
- Membuat penghulu baru
- Assign penghulu ke pendaftaran nikah
- Mengelola semua aspek sistem

---

## 🚀 Cara Kerja

### 1. Automatic Seeding (Default)

Seeder akan **otomatis dijalankan** saat aplikasi start jika user kepala KUA belum ada.

**Lokasi:** `cmd/api/main.go` (line 73-77)

```go
// Seed initial data (Kepala KUA)
if err := seeders.SeedKepalaKUA(DB); err != nil {
    log.Printf("Warning: Failed to seed kepala KUA: %v", err)
    // Don't fatal, just warn - seeder is optional
}
```

**Default Credentials:**
- **Username:** `kepalakua`
- **Email:** `kepalakua@kua.go.id`
- **Password:** `kepalakua123` ⚠️ **CHANGE THIS IN PRODUCTION!**
- **Nama:** `Kepala KUA Banjarmasin Utara`
- **NIP:** `197001011990011000`
- **Role:** `kepala_kua`

**Log Output:**
```
🌱 Seeding Kepala KUA user...
✅ Kepala KUA user created successfully!
   User ID: KKUA1704067200
   Username: kepalakua
   Email: kepalakua@kua.go.id
   Password: kepalakua123 (⚠️  CHANGE THIS IN PRODUCTION!)
   Role: kepala_kua

⚠️  IMPORTANT: Change the default password after first login!
```

---

### 2. Manual Seeding dengan Custom Credentials

Anda juga bisa menjalankan seeder secara manual dengan credentials custom.

#### A. Menggunakan Environment Variables

**Buat file `.env` atau set environment variables:**
```bash
SEEDER_USERNAME=kepalakua_custom
SEEDER_EMAIL=kepalakua@custom.kua.go.id
SEEDER_PASSWORD=SecurePassword123!
SEEDER_NAMA=Dr. H. Ahmad Wijaya, S.H., M.H.
SEEDER_NIP=197001011990011001
```

**Jalankan seeder:**
```bash
go run cmd/seeder/main.go
```

#### B. Menggunakan Command Line Arguments (Future)

*(Bisa ditambahkan jika diperlukan)*

---

## 📝 File Structure

```
simnikah-api/
├── internal/
│   └── seeders/
│       └── kepala_kua_seeder.go    # Seeder implementation
├── cmd/
│   └── seeder/
│       └── main.go                  # Manual seeder runner
└── cmd/
    └── api/
        └── main.go                  # Auto seeder call
```

---

## 🔧 Functions

### `SeedKepalaKUA(db *gorm.DB) error`

Membuat kepala KUA dengan **default credentials**.

**Features:**
- ✅ Check jika user sudah ada (skip jika sudah ada)
- ✅ Create user account dengan role `kepala_kua`
- ✅ Create staff profile dengan jabatan `Kepala KUA`
- ✅ Hash password dengan bcrypt (cost 12)
- ✅ Generate unique user_id: `KKUA{timestamp}`

**Returns:**
- `error` - Error jika gagal, `nil` jika sukses atau sudah ada

---

### `SeedKepalaKUAWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip string) error`

Membuat kepala KUA dengan **custom credentials**.

**Parameters:**
- `username` - Username untuk login (default: "kepalakua" jika kosong)
- `email` - Email user (default: "kepalakua@kua.go.id" jika kosong)
- `password` - Password (default: "kepalakua123" jika kosong)
- `nama` - Nama lengkap (default: "Kepala KUA Banjarmasin Utara" jika kosong)
- `nip` - NIP (default: "197001011990011000" jika kosong)

**Features:**
- ✅ Semua parameter optional (akan menggunakan default jika kosong)
- ✅ Check jika user sudah ada
- ✅ Create user + staff profile
- ✅ Hash password dengan bcrypt

---

## 🔒 Security Notes

### ⚠️ IMPORTANT: Change Default Password!

**Default password adalah:** `kepalakua123`

**Ini HANYA untuk development/testing!**

**Untuk Production:**
1. Set environment variables dengan password yang kuat
2. Atau jalankan manual seeder dengan password custom
3. **WAJIB** ganti password setelah first login
4. Gunakan password yang kuat (min 12 karakter, kombinasi huruf, angka, simbol)

### Best Practices:

1. **Development:**
   - Gunakan default credentials untuk kemudahan
   - Ganti password setelah testing

2. **Staging:**
   - Gunakan environment variables
   - Password: `StagingPassword123!`

3. **Production:**
   - **WAJIB** gunakan environment variables
   - Password: Strong password (min 16 karakter)
   - Jangan commit password ke git!

---

## 🧪 Testing

### Test Auto Seeder:

```bash
# 1. Pastikan database kosong (atau user belum ada)
# 2. Jalankan aplikasi
go run cmd/api/main.go

# Expected output:
# 🌱 Seeding Kepala KUA user...
# ✅ Kepala KUA user created successfully!
```

### Test Manual Seeder:

```bash
# Dengan default credentials
go run cmd/seeder/main.go

# Dengan custom credentials (via env)
export SEEDER_USERNAME="kepalakua_test"
export SEEDER_EMAIL="test@kua.go.id"
export SEEDER_PASSWORD="TestPassword123!"
export SEEDER_NAMA="Test Kepala KUA"
export SEEDER_NIP="197001011990011002"
go run cmd/seeder/main.go
```

### Test Login:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "kepalakua",
    "password": "kepalakua123"
  }'
```

**Expected Response:**
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "KKUA1704067200",
    "email": "kepalakua@kua.go.id",
    "role": "kepala_kua",
    "nama": "Kepala KUA Banjarmasin Utara"
  }
}
```

---

## 🔍 Verification

### Check User di Database:

```sql
-- Check user
SELECT user_id, username, email, role, status, nama 
FROM users 
WHERE role = 'kepala_kua';

-- Check staff profile
SELECT s.*, u.username, u.email 
FROM staff_kuas s
JOIN users u ON s.user_id = u.user_id
WHERE s.jabatan = 'Kepala KUA';
```

### Check via API:

```bash
# Login dulu untuk dapat token
TOKEN=$(curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"kepalakua","password":"kepalakua123"}' \
  | jq -r '.token')

# Get profile
curl http://localhost:8080/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🐛 Troubleshooting

### Error: "Username atau email sudah digunakan"

**Penyebab:** User kepala KUA sudah ada di database.

**Solusi:**
- Seeder akan skip jika user sudah ada (tidak error)
- Jika ingin reset, hapus user manual dari database

### Error: "Failed to seed kepala KUA"

**Penyebab:** Database error atau constraint violation.

**Solusi:**
1. Check database connection
2. Check database permissions
3. Check logs untuk detail error

### User Tidak Bisa Login

**Penyebab:** Password tidak match atau user tidak aktif.

**Solusi:**
1. Check password yang digunakan
2. Check user status di database (harus "Aktif")
3. Reset password jika perlu

---

## 📚 Related Files

- `internal/seeders/kepala_kua_seeder.go` - Seeder implementation
- `cmd/seeder/main.go` - Manual seeder runner
- `cmd/api/main.go` - Auto seeder call
- `internal/models/models.go` - User & StaffKUA models
- `internal/models/constants.go` - Role & status constants

---

## 🎯 Next Steps

Setelah seeder berhasil:

1. ✅ **Login dengan credentials default**
2. ✅ **Ganti password** (via API atau manual update database)
3. ✅ **Buat staff baru** (via API endpoint `/simnikah/staff`)
4. ✅ **Buat penghulu baru** (via API endpoint `/simnikah/penghulu`)
5. ✅ **Start managing registrations!**

---

**Last Updated:** 27 Januari 2025

