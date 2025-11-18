# 👨‍💼 Kepala KUA - Staff & Penghulu Management

## 📋 Overview

Kepala KUA sekarang dapat menambahkan staff (seperti staff pendaftaran) dan penghulu beserta akun aplikasinya melalui API endpoints khusus.

---

## 🔑 Endpoints

### 1. Create Staff
**POST** `/simnikah/kepala-kua/staff`

Menambahkan staff baru (staff pendaftaran, penghulu, atau kepala KUA) beserta akun aplikasinya.

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "username": "staff001",
  "email": "staff001@kua.go.id",
  "password": "password123",
  "nama": "Ahmad Staff",
  "nip": "198012345678901234",
  "jabatan": "Staff",
  "bagian": "Pendaftaran",
  "no_hp": "081234567890",
  "alamat": "Jl. Ahmad Yani No. 123, Banjarmasin"
}
```

**Jabatan yang tersedia:**
- `Staff` - Staff pendaftaran
- `Penghulu` - Penghulu (marriage officer)
- `Kepala KUA` - Kepala KUA

**Response:**
```json
{
  "success": true,
  "message": "Staff berhasil dibuat",
  "data": {
    "user": {
      "user_id": "STF1704067200",
      "username": "staff001",
      "email": "staff001@kua.go.id",
      "role": "staff",
      "nama": "Ahmad Staff"
    },
    "staff": {
      "id": 1,
      "nip": "198012345678901234",
      "nama_lengkap": "Ahmad Staff",
      "jabatan": "Staff",
      "bagian": "Pendaftaran",
      "no_hp": "081234567890",
      "email": "staff001@kua.go.id",
      "alamat": "Jl. Ahmad Yani No. 123, Banjarmasin",
      "status": "Aktif"
    }
  }
}
```

**Note:** 
- User role akan otomatis ditentukan berdasarkan jabatan:
  - `Jabatan = "Staff"` → `Role = "staff"`
  - `Jabatan = "Penghulu"` → `Role = "penghulu"`
  - `Jabatan = "Kepala KUA"` → `Role = "kepala_kua"`

---

### 2. Create Marriage Officer (Penghulu)
**POST** `/simnikah/kepala-kua/penghulu`

Menambahkan penghulu baru beserta akun aplikasinya.

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "username": "penghulu001",
  "email": "penghulu001@kua.go.id",
  "password": "password123",
  "nama": "H. Ahmad Penghulu",
  "nip": "197012345678901234",
  "no_hp": "081234567890",
  "email_penghulu": "penghulu001@email.com",
  "alamat": "Jl. Ahmad Yani No. 456, Banjarmasin"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Penghulu berhasil dibuat",
  "data": {
    "user": {
      "user_id": "PNG1704067200",
      "username": "penghulu001",
      "email": "penghulu001@kua.go.id",
      "role": "penghulu",
      "nama": "H. Ahmad Penghulu"
    },
    "penghulu": {
      "id": 1,
      "nip": "197012345678901234",
      "nama_lengkap": "H. Ahmad Penghulu",
      "no_hp": "081234567890",
      "email": "penghulu001@email.com",
      "alamat": "Jl. Ahmad Yani No. 456, Banjarmasin",
      "status": "Aktif"
    }
  }
}
```

**Note:**
- Jika `email_penghulu` tidak diisi, akan menggunakan `email` dari user account
- Password akan di-hash menggunakan bcrypt sebelum disimpan
- User account dan profile penghulu dibuat secara transaksional (jika salah satu gagal, keduanya tidak akan dibuat)

---

## ✅ Features

### 1. **Validasi Data**
- ✅ Username dan email harus unik
- ✅ NIP harus unik
- ✅ Password minimal 6 karakter
- ✅ Email harus format valid
- ✅ Jabatan harus valid (Staff, Penghulu, atau Kepala KUA)

### 2. **Auto Role Assignment**
- ✅ Role otomatis ditentukan berdasarkan jabatan
- ✅ User role dan staff profile konsisten

### 3. **Transaction Safety**
- ✅ Jika pembuatan staff/penghulu gagal, user account akan di-rollback
- ✅ Jika pembuatan user account gagal, tidak ada data yang tersimpan

### 4. **Notification**
- ✅ Notifikasi otomatis dikirim ke staff baru setelah akun dibuat
- ✅ Error pada notifikasi tidak akan mempengaruhi proses pembuatan akun

---

## 📝 Contoh Penggunaan

### Menambahkan Staff Pendaftaran

```bash
curl -X POST https://api.example.com/simnikah/kepala-kua/staff \
  -H "Authorization: Bearer <kepala_kua_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "staff001",
    "email": "staff001@kua.go.id",
    "password": "password123",
    "nama": "Ahmad Staff",
    "nip": "198012345678901234",
    "jabatan": "Staff",
    "bagian": "Pendaftaran",
    "no_hp": "081234567890",
    "alamat": "Jl. Ahmad Yani No. 123"
  }'
```

### Menambahkan Penghulu

```bash
curl -X POST https://api.example.com/simnikah/kepala-kua/penghulu \
  -H "Authorization: Bearer <kepala_kua_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "penghulu001",
    "email": "penghulu001@kua.go.id",
    "password": "password123",
    "nama": "H. Ahmad Penghulu",
    "nip": "197012345678901234",
    "no_hp": "081234567890",
    "email_penghulu": "penghulu001@email.com",
    "alamat": "Jl. Ahmad Yani No. 456"
  }'
```

---

## 🔒 Security

1. **Authentication**: Hanya kepala KUA yang dapat mengakses endpoint ini
2. **Password Hashing**: Password di-hash menggunakan bcrypt sebelum disimpan
3. **Validation**: Semua input divalidasi sebelum disimpan ke database
4. **Transaction**: Operasi database menggunakan transaction untuk memastikan konsistensi data

---

## ❌ Error Handling

### Username atau Email Sudah Digunakan
```json
{
  "success": false,
  "message": "Username atau email sudah digunakan",
  "error": "Username atau email sudah terdaftar di sistem"
}
```

### NIP Sudah Terdaftar
```json
{
  "success": false,
  "message": "NIP sudah terdaftar",
  "error": "NIP sudah terdaftar untuk staff/penghulu lain"
}
```

### Jabatan Tidak Valid
```json
{
  "success": false,
  "message": "Jabatan tidak valid",
  "error": "Jabatan harus salah satu dari: Staff, Penghulu, atau Kepala KUA"
}
```

---

## 📌 Catatan Penting

1. **User ID Format:**
   - Staff: `STF{timestamp}`
   - Penghulu: `PNG{timestamp}`

2. **Default Values:**
   - Status: `Aktif` (untuk staff dan penghulu)
   - User Status: `Aktif`
   - Jumlah Nikah: `0` (untuk penghulu)
   - Rating: `0.0` (untuk penghulu)

3. **Field Optional:**
   - `no_hp`: Nomor HP (opsional)
   - `alamat`: Alamat (opsional)
   - `email_penghulu`: Email penghulu khusus (opsional, akan menggunakan `email` jika tidak diisi)

4. **Rollback Mechanism:**
   - Jika pembuatan staff/penghulu gagal setelah user account dibuat, user account akan dihapus otomatis
   - Ini memastikan tidak ada "orphan" user account di database

---

## ✅ Benefits

1. **Efficient**: Kepala KUA dapat menambahkan staff dan penghulu langsung dari aplikasi
2. **Consistent**: User account dan profile dibuat secara konsisten
3. **Secure**: Password di-hash dan semua input divalidasi
4. **Safe**: Transaction mechanism memastikan data konsisten
5. **User-Friendly**: Notifikasi otomatis dikirim ke staff/penghulu baru

