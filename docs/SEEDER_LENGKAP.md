# 🌱 Database Seeders - Dokumentasi Lengkap

Dokumentasi untuk semua seeder yang tersedia di SimNikah API.

---

## 📋 Daftar Seeders

1. **Kepala KUA** - User dengan akses penuh
2. **Staff KUA** - User untuk verifikasi formulir dan berkas
3. **Penghulu** - User untuk memimpin nikah

---

## 🚀 Auto Seeding (Default)

Semua seeder akan **otomatis dijalankan** saat aplikasi start jika user belum ada.

**Lokasi:** `cmd/api/main.go` (line 73-87)

```go
// Seed initial data (Kepala KUA, Staff, Penghulu)
if err := seeders.SeedKepalaKUA(DB); err != nil {
    log.Printf("Warning: Failed to seed kepala KUA: %v", err)
}

if err := seeders.SeedStaff(DB); err != nil {
    log.Printf("Warning: Failed to seed staff: %v", err)
}

if err := seeders.SeedPenghulu(DB); err != nil {
    log.Printf("Warning: Failed to seed penghulu: %v", err)
}
```

---

## 1. Kepala KUA Seeder

### Default Credentials:
```
Username: kepalakua
Email: kepalakua@kua.go.id
Password: kepalakua123 ⚠️ CHANGE IN PRODUCTION!
Nama: Kepala KUA Banjarmasin Utara
NIP: 197001011990011000
Role: kepala_kua
Jabatan: Kepala KUA
```

### Functions:

**`SeedKepalaKUA(db *gorm.DB) error`**
- Membuat kepala KUA dengan default credentials
- Auto-skip jika user sudah ada

**`SeedKepalaKUAWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip string) error`**
- Membuat kepala KUA dengan custom credentials
- Semua parameter optional (akan menggunakan default jika kosong)

---

## 2. Staff KUA Seeder

### Default Credentials:
```
Username: staff001
Email: staff@kua.go.id
Password: staff123 ⚠️ CHANGE IN PRODUCTION!
Nama: Staff KUA Banjarmasin Utara
NIP: 197001011990011001
Role: staff
Jabatan: Staff
Bagian: Verifikasi
```

### Functions:

**`SeedStaff(db *gorm.DB) error`**
- Membuat staff dengan default credentials
- Auto-skip jika user sudah ada

**`SeedStaffWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip, jabatan, bagian string) error`**
- Membuat staff dengan custom credentials
- Validasi jabatan: "Staff", "Penghulu", atau "Kepala KUA"
- Semua parameter optional

---

## 3. Penghulu Seeder

### Default Credentials:
```
Username: penghulu001
Email: penghulu@kua.go.id
Password: penghulu123 ⚠️ CHANGE IN PRODUCTION!
Nama: Ustadz Haji Ahmad Wijaya
NIP: 197001011990011002
Role: penghulu
Status: Aktif
Jumlah Nikah: 0
Rating: 0.0
```

### Functions:

**`SeedPenghulu(db *gorm.DB) error`**
- Membuat penghulu dengan default credentials
- Auto-skip jika user sudah ada

**`SeedPenghuluWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip string) error`**
- Membuat penghulu dengan custom credentials
- Semua parameter optional

---

## 🛠️ Manual Seeding

### Menggunakan Command Line Tool

**Jalankan semua seeder:**
```bash
go run cmd/seeder/main.go
```

**Jalankan seeder tertentu:**
```bash
# Hanya Kepala KUA
go run cmd/seeder/main.go -type=kepala_kua

# Hanya Staff
go run cmd/seeder/main.go -type=staff

# Hanya Penghulu
go run cmd/seeder/main.go -type=penghulu

# Semua (default)
go run cmd/seeder/main.go -type=all
```

### Menggunakan Environment Variables

**Kepala KUA:**
```bash
export SEEDER_KEPALA_KUA_USERNAME="kepalakua_custom"
export SEEDER_KEPALA_KUA_EMAIL="kepalakua@custom.kua.go.id"
export SEEDER_KEPALA_KUA_PASSWORD="SecurePassword123!"
export SEEDER_KEPALA_KUA_NAMA="Dr. H. Ahmad Wijaya, S.H., M.H."
export SEEDER_KEPALA_KUA_NIP="197001011990011000"

go run cmd/seeder/main.go -type=kepala_kua
```

**Staff:**
```bash
export SEEDER_STAFF_USERNAME="staff_custom"
export SEEDER_STAFF_EMAIL="staff@custom.kua.go.id"
export SEEDER_STAFF_PASSWORD="SecurePassword123!"
export SEEDER_STAFF_NAMA="Budi Santoso"
export SEEDER_STAFF_NIP="197001011990011001"
export SEEDER_STAFF_JABATAN="Staff"
export SEEDER_STAFF_BAGIAN="Verifikasi"

go run cmd/seeder/main.go -type=staff
```

**Penghulu:**
```bash
export SEEDER_PENGHULU_USERNAME="penghulu_custom"
export SEEDER_PENGHULU_EMAIL="penghulu@custom.kua.go.id"
export SEEDER_PENGHULU_PASSWORD="SecurePassword123!"
export SEEDER_PENGHULU_NAMA="Ustadz Haji Ahmad Wijaya"
export SEEDER_PENGHULU_NIP="197001011990011002"

go run cmd/seeder/main.go -type=penghulu
```

---

## 🧪 Testing

### Test Login dengan Default Credentials

**1. Kepala KUA:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "kepalakua",
    "password": "kepalakua123"
  }'
```

**2. Staff:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "staff001",
    "password": "staff123"
  }'
```

**3. Penghulu:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "penghulu001",
    "password": "penghulu123"
  }'
```

---

## 📊 Default Users Summary

| Role | Username | Email | Password | NIP |
|------|----------|-------|----------|-----|
| **Kepala KUA** | `kepalakua` | `kepalakua@kua.go.id` | `kepalakua123` | `197001011990011000` |
| **Staff** | `staff001` | `staff@kua.go.id` | `staff123` | `197001011990011001` |
| **Penghulu** | `penghulu001` | `penghulu@kua.go.id` | `penghulu123` | `197001011990011002` |

---

## 🔒 Security Notes

### ⚠️ IMPORTANT: Change Default Passwords!

**Default passwords:**
- Kepala KUA: `kepalakua123`
- Staff: `staff123`
- Penghulu: `penghulu123`

**Ini HANYA untuk development/testing!**

**Untuk Production:**
1. Set environment variables dengan password yang kuat
2. Atau jalankan manual seeder dengan password custom
3. **WAJIB** ganti password setelah first login
4. Gunakan password yang kuat (min 12 karakter)

---

## 📁 File Structure

```
simnikah-api/
├── internal/
│   └── seeders/
│       ├── kepala_kua_seeder.go    # Kepala KUA seeder
│       ├── staff_seeder.go         # Staff seeder
│       └── penghulu_seeder.go      # Penghulu seeder
├── cmd/
│   ├── api/
│   │   └── main.go                 # Auto seeder call
│   └── seeder/
│       └── main.go                 # Manual seeder runner
└── docs/
    ├── SEEDER_KEPALA_KUA.md        # Kepala KUA docs
    └── SEEDER_LENGKAP.md           # This file
```

---

## 🔍 Verification

### Check Users di Database:

```sql
-- Check all seeded users
SELECT user_id, username, email, role, status, nama 
FROM users 
WHERE role IN ('kepala_kua', 'staff', 'penghulu')
ORDER BY role;

-- Check staff profiles
SELECT s.*, u.username, u.email 
FROM staff_kuas s
JOIN users u ON s.user_id = u.user_id
ORDER BY s.jabatan;

-- Check penghulu profiles
SELECT p.*, u.username, u.email 
FROM penghulus p
JOIN users u ON p.user_id = u.user_id;
```

---

## 🐛 Troubleshooting

### Error: "Username atau email sudah digunakan"

**Penyebab:** User sudah ada di database.

**Solusi:**
- Seeder akan skip jika user sudah ada (tidak error)
- Jika ingin reset, hapus user manual dari database

### Error: "NIP already exists"

**Penyebab:** NIP sudah digunakan oleh user lain.

**Solusi:**
- Gunakan NIP yang berbeda
- Atau hapus user dengan NIP tersebut

### Error: "Invalid jabatan"

**Penyebab:** Jabatan tidak valid.

**Solusi:**
- Gunakan salah satu: "Staff", "Penghulu", atau "Kepala KUA"

---

## 📚 Related Files

- `internal/seeders/kepala_kua_seeder.go` - Kepala KUA seeder
- `internal/seeders/staff_seeder.go` - Staff seeder
- `internal/seeders/penghulu_seeder.go` - Penghulu seeder
- `cmd/seeder/main.go` - Manual seeder runner
- `cmd/api/main.go` - Auto seeder call

---

**Last Updated:** 27 Januari 2025

