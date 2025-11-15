# 📚 SimNikah API - Dokumentasi Lengkap untuk Frontend

**Base URL:**
```
Development: http://localhost:8080
Production: https://your-app.railway.app
```

**Content-Type:** `application/json`

**Authentication:** 
```
Authorization: Bearer <jwt_token>
```

---

## 📋 Daftar Isi

1. [Authentication](#1-authentication)
2. [Marriage Registration](#2-marriage-registration)
3. [Staff Management](#3-staff-management)
4. [Penghulu Management](#4-penghulu-management)
5. [Calendar & Schedule](#5-calendar--schedule)
6. [Bimbingan Perkawinan](#6-bimbingan-perkawinan)
7. [Location & Map](#7-location--map)
8. [Notifications](#8-notifications)
9. [Status Codes & Error Handling](#9-status-codes--error-handling)

---

## 1. Authentication

### 1.1 Register User

**Endpoint:** `POST /register`  
**Auth Required:** ❌ No  
**Rate Limit:** 5 requests/minute

#### Request Body:
```json
{
  "username": "ahmad123",
  "email": "ahmad@example.com",
  "password": "password123",
  "nama": "Ahmad Wijaya",
  "role": "user_biasa"
}
```

#### Field Validation:
- `username`: Required, string, unique
- `email`: Required, valid email format, unique
- `password`: Required, minimum 6 characters
- `nama`: Required, string
- `role`: Required, one of: `user_biasa`, `staff`, `penghulu`, `kepala_kua`

#### ✅ Success Response (201):
```json
{
  "message": "User berhasil dibuat",
  "user": {
    "user_id": "USR1704067200",
    "username": "ahmad123",
    "email": "ahmad@example.com",
    "nama": "Ahmad Wijaya",
    "role": "user_biasa"
  }
}
```

#### ❌ Error Responses:

**400 - Validation Error:**
```json
{
  "error": "Format registrasi tidak valid"
}
```

**400 - Username/Email Already Exists:**
```json
{
  "error": "Username atau email sudah digunakan"
}
```

**400 - Invalid Role:**
```json
{
  "error": "Role tidak valid. Role yang tersedia: user_biasa, penghulu, staff, kepala_kua"
}
```

**429 - Rate Limit Exceeded:**
```json
{
  "success": false,
  "message": "Terlalu banyak percobaan login",
  "error": "Demi keamanan, Anda harus menunggu sebelum mencoba lagi.",
  "retry_after": "60 detik",
  "type": "security"
}
```

---

### 1.2 Login

**Endpoint:** `POST /login`  
**Auth Required:** ❌ No  
**Rate Limit:** 5 requests/minute

#### Request Body:
```json
{
  "username": "ahmad123",
  "password": "password123"
}
```

#### ✅ Success Response (200):
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiVVNSMTcwNDA2NzIwMCIsImVtYWlsIjoiYWhtYWRAZXhhbXBsZS5jb20iLCJyb2xlIjoi dXNlcl9iaWFzYSIsIm5hbWEiOiJBaG1hZCBXaWpheWEiLCJleHAiOjE3MDQxNTM2MDAsImlhdCI6MTcwNDA2NzIwMCwibmJmIjoxNzA0MDY3MjAwfQ.xxx",
  "user": {
    "user_id": "USR1704067200",
    "email": "ahmad@example.com",
    "role": "user_biasa",
    "nama": "Ahmad Wijaya"
  }
}
```

#### ❌ Error Responses:

**400 - Validation Error:**
```json
{
  "error": "Format login tidak valid"
}
```

**401 - Username Not Found:**
```json
{
  "error": "Username tidak ditemukan"
}
```

**401 - Wrong Password:**
```json
{
  "error": "Password salah"
}
```

**401 - User Inactive:**
```json
{
  "error": "User tidak aktif"
}
```

---

### 1.3 Get Profile

**Endpoint:** `GET /profile`  
**Auth Required:** ✅ Yes

#### Headers:
```
Authorization: Bearer <jwt_token>
```

#### ✅ Success Response (200):
```json
{
  "message": "Profile berhasil diambil",
  "user": {
    "user_id": "USR1704067200",
    "username": "ahmad123",
    "email": "ahmad@example.com",
    "role": "user_biasa",
    "nama": "Ahmad Wijaya"
  }
}
```

#### ❌ Error Responses:

**401 - Unauthorized (No Token):**
```json
{
  "error": "Token tidak ditemukan"
}
```

**401 - Invalid Token:**
```json
{
  "error": "Token tidak valid"
}
```

---

## 2. Marriage Registration

### 2.1 Create Marriage Registration (Form Baru)

**Endpoint:** `POST /simnikah/pendaftaran/form-baru`  
**Auth Required:** ✅ Yes (user_biasa)  
**Complexity:** ⭐⭐⭐⭐⭐ (Sangat Kompleks)

#### Request Body (LENGKAP):
```json
{
  "calon_suami": {
    "nik": "6371012501950001",
    "nama_lengkap": "Ahmad Wijaya",
    "tempat_lahir": "Banjarmasin",
    "tanggal_lahir": "1995-01-25",
    "alamat": "Jl. Pangeran Samudra No. 88, RT 05, RW 02",
    "rt": "05",
    "rw": "02",
    "kelurahan": "Pangeran",
    "kecamatan": "Banjarmasin Utara",
    "kabupaten": "Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "agama": "Islam",
    "status": "Belum Kawin",
    "pekerjaan": "Karyawan Swasta",
    "deskripsi_pekerjaan": "",
    "pendidikan": "S1",
    "penghasilan": 5000000,
    "nomor_telepon": "081234567890",
    "email": "ahmad@example.com",
    "kewarganegaraan": "WNI",
    "nomor_paspor": ""
  },
  "calon_istri": {
    "nik": "6371013001950002",
    "nama_lengkap": "Siti Nurhaliza",
    "tempat_lahir": "Banjarmasin",
    "tanggal_lahir": "1995-01-30",
    "alamat": "Jl. Ahmad Yani No. 123, RT 03, RW 01",
    "rt": "03",
    "rw": "01",
    "kelurahan": "Sungai Miai",
    "kecamatan": "Banjarmasin Utara",
    "kabupaten": "Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "agama": "Islam",
    "status": "Belum Kawin",
    "pekerjaan": "Guru",
    "deskripsi_pekerjaan": "",
    "pendidikan": "S1",
    "penghasilan": 4000000,
    "nomor_telepon": "081234567891",
    "email": "siti@example.com",
    "kewarganegaraan": "WNI",
    "nomor_paspor": ""
  },
  "orang_tua_calon_suami": {
    "ayah": {
      "status_keberadaan": "Hidup",
      "nama": "Bapak Suami",
      "nik": "6371010101700001",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "tempat_lahir": "Banjarmasin",
      "negara_asal": "",
      "pekerjaan": "PNS",
      "deskripsi_pekerjaan": "",
      "alamat": "Jl. Pangeran Samudra No. 88"
    },
    "ibu": {
      "status_keberadaan": "Hidup",
      "nama": "Ibu Suami",
      "nik": "6371010201700002",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "tempat_lahir": "Banjarmasin",
      "negara_asal": "",
      "pekerjaan": "Ibu Rumah Tangga",
      "deskripsi_pekerjaan": "",
      "alamat": "Jl. Pangeran Samudra No. 88"
    }
  },
  "orang_tua_calon_istri": {
    "ayah": {
      "status_keberadaan": "Hidup",
      "nama": "Bapak Istri",
      "nik": "6371010301700003",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "tempat_lahir": "Banjarmasin",
      "negara_asal": "",
      "pekerjaan": "Wiraswasta",
      "deskripsi_pekerjaan": "",
      "alamat": "Jl. Ahmad Yani No. 123"
    },
    "ibu": {
      "status_keberadaan": "Hidup",
      "nama": "Ibu Istri",
      "nik": "6371010401700004",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "tempat_lahir": "Banjarmasin",
      "negara_asal": "",
      "pekerjaan": "Ibu Rumah Tangga",
      "deskripsi_pekerjaan": "",
      "alamat": "Jl. Ahmad Yani No. 123"
    }
  },
  "wali_nikah": {
    "nik_wali": "6371010301700003",
    "nama_lengkap_wali": "Bapak Istri",
    "hubungan_wali": "Ayah Kandung",
    "alamat_wali": "Jl. Ahmad Yani No. 123",
    "nomor_telepon_wali": "081234567892",
    "agama_wali": "Islam",
    "status_wali": "Hidup"
  },
  "jadwal_dan_lokasi": {
    "tanggal_nikah": "2025-02-15",
    "waktu_nikah": "10:00",
    "lokasi_nikah": "Di KUA",
    "alamat_nikah": "",
    "nomor_dispensasi": ""
  }
}
```

#### Field Validations:

**Calon Suami/Istri:**
- `nik`: Required, 16 characters, unique
- `nama_lengkap`: Required, min 3 characters
- `tempat_lahir`: Required, min 2 characters
- `tanggal_lahir`: Required, format YYYY-MM-DD
- `alamat`: Required, min 10 characters
- `agama`: Required
- `status`: Required, one of: "Belum Kawin", "Kawin", "Cerai Hidup", "Cerai Mati"
- `pekerjaan`: Required
- `deskripsi_pekerjaan`: Required if `pekerjaan` = "Lainnya"
- `pendidikan`: Required
- `nomor_telepon`: Required, min 10 digits, starts with "08"
- `email`: Required, valid email format
- `kewarganegaraan`: Required, "WNI" or "WNA"
- `nomor_paspor`: Required if `kewarganegaraan` = "WNA"

**Orang Tua (Conditional):**
- Jika `status_keberadaan` = "Hidup", maka semua field wajib diisi:
  - `nama`, `nik`, `kewarganegaraan`, `agama`, `pekerjaan`, `alamat`
  - Jika WNA: `negara_asal`, `nomor_paspor` wajib
  - Jika pekerjaan = "Lainnya": `deskripsi_pekerjaan` wajib
- Jika `status_keberadaan` = "Meninggal" atau "Tidak Diketahui", field lain tidak wajib

**Wali Nikah:**
- `nik_wali`: Required, 16 characters
- `nama_lengkap_wali`: Required
- `hubungan_wali`: Required, valid sesuai syariat
- `status_wali`: Required, "Hidup" (wajib hidup untuk menjadi wali)
- Validasi syariat:
  - Jika ayah hidup → wali HARUS "Ayah Kandung"
  - NIK wali harus sama dengan NIK ayah (jika ayah hidup)
  - Wali tidak boleh sama dengan calon pengantin

**Jadwal & Lokasi:**
- `tanggal_nikah`: Required, format YYYY-MM-DD, tidak boleh masa lalu
- `waktu_nikah`: Required, format HH:MM (24-hour)
- `lokasi_nikah`: Required, "Di KUA" or "Di Luar KUA"
- `alamat_nikah`: Required if `lokasi_nikah` = "Di Luar KUA"
- `nomor_dispensasi`: Required if:
  - Nikah < 10 hari kerja dari pendaftaran, ATAU
  - Usia suami < 19 tahun, ATAU
  - Usia istri < 19 tahun

#### ✅ Success Response (201):
```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat",
  "data": {
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "pendaftaran_id": 1,
    "calon_suami": {
      "id": 1,
      "user_id": "abc123def456",
      "nik": "6371012501950001",
      "nama_lengkap": "Ahmad Wijaya",
      ...
    },
    "calon_istri": {
      "id": 2,
      "user_id": "def456ghi789",
      "nik": "6371013001950002",
      "nama_lengkap": "Siti Nurhaliza",
      ...
    },
    "pendaftaran": {
      "id": 1,
      "nomor_pendaftaran": "NIKa1b2c3d4e5",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2025-02-15",
      "waktu_nikah": "10:00",
      "tempat_nikah": "Di KUA",
      ...
    },
    "wali_nikah": {
      "id": 1,
      "nik": "6371010301700003",
      "nama_lengkap": "Bapak Istri",
      "hubungan_wali": "Ayah Kandung",
      ...
    }
  }
}
```

#### ❌ Error Responses:

**400 - Format Error:**
```json
{
  "success": false,
  "message": "Format data tidak valid",
  "error": "Format data tidak valid: Key: 'CalonSuami.Nik' Error:Field validation for 'Nik' failed on the 'required' tag",
  "type": "validation"
}
```

**400 - Date Format Error:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format tanggal nikah tidak valid (YYYY-MM-DD)",
  "field": "tanggal_nikah",
  "type": "format"
}
```

**400 - Date in Past:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Tanggal nikah tidak boleh di masa lalu",
  "field": "tanggal_nikah",
  "type": "validation"
}
```

**400 - Time Format Error:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format waktu nikah tidak valid (HH:MM dalam format 24-jam, contoh: 09:00)",
  "field": "waktu_nikah",
  "type": "format"
}
```

**400 - Dispensasi Required:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Nomor dispensasi wajib diisi karena: Pelaksanaan nikah kurang dari 10 hari kerja dan Calon suami berumur kurang dari 19 tahun",
  "field": "nomor_dispensasi",
  "type": "required",
  "details": {
    "reason": "Pelaksanaan nikah kurang dari 10 hari kerja dan Calon suami berumur kurang dari 19 tahun",
    "working_days": 5,
    "groom_age": 18,
    "bride_age": 20
  }
}
```

**400 - Wali Nikah Validation (Wali Meninggal):**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Wali nikah yang telah meninggal dunia tidak dapat menjadi wali. Silakan pilih wali lain yang masih hidup.",
  "field": "status_wali",
  "type": "syariat_validation"
}
```

**400 - Wali Nikah Validation (Ayah Hidup tapi Wali Bukan Ayah):**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Jika wali adalah Ayah Kandung, maka status wali harus sama dengan status ayah catin perempuan (Hidup)",
  "field": "status_wali",
  "type": "syariat_validation"
}
```

**400 - Wali Nikah Validation (NIK Tidak Cocok):**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "NIK wali harus sama dengan NIK ayah kandung catin perempuan",
  "field": "nik_wali",
  "type": "syariat_validation",
  "details": {
    "nik_wali_yang_diinput": "6371010301700003",
    "nik_ayah_catin_perempuan": "6371010301700004"
  }
}
```

**400 - Wali Nikah Validation (Wali = Calon Pengantin):**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "NIK wali tidak boleh sama dengan NIK calon suami",
  "field": "nik_wali",
  "type": "syariat_validation"
}
```

**400 - Duplicate Registration:**
```json
{
  "success": false,
  "message": "Pendaftaran sudah ada",
  "error": "Anda sudah memiliki pendaftaran nikah yang masih aktif",
  "field": "pendaftaran",
  "type": "duplicate",
  "data": {
    "existing_registration_id": 1,
    "status": "Menunggu Verifikasi",
    "nomor_pendaftaran": "NIKa1b2c3d4e5"
  }
}
```

**401 - Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "User ID tidak ditemukan",
  "type": "authentication"
}
```

**500 - Database Error:**
```json
{
  "success": false,
  "message": "Database error",
  "error": "Gagal membuat profile calon suami: ...",
  "type": "database"
}
```

---

### 2.2 Check Registration Status

**Endpoint:** `GET /simnikah/pendaftaran/status`  
**Auth Required:** ✅ Yes (user_biasa)

#### ✅ Success Response - No Registration (200):
```json
{
  "success": true,
  "message": "User belum memiliki pendaftaran nikah",
  "data": {
    "has_registration": false,
    "can_register": true
  }
}
```

#### ✅ Success Response - Has Registration (200):
```json
{
  "success": true,
  "message": "User sudah memiliki pendaftaran nikah",
  "data": {
    "has_registration": true,
    "can_register": false,
    "registration": {
      "id": 1,
      "nomor_pendaftaran": "NIKa1b2c3d4e5",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2025-02-15",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "KUA Kecamatan Banjarmasin Utara",
      "created_at": "2025-01-27T10:00:00Z"
    }
  }
}
```

---

### 2.3 Mark As Visited

**Endpoint:** `POST /simnikah/pendaftaran/:id/mark-visited`  
**Auth Required:** ✅ Yes (user_biasa)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{}
```
*(Tidak ada body, hanya endpoint call)*

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Status berhasil diupdate",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "status_pendaftaran": "Menunggu Penugasan",
    "updated_at": "2025-01-27T10:30:00Z"
  }
}
```

#### ❌ Error Responses:

**404 - Not Found:**
```json
{
  "success": false,
  "message": "Pendaftaran tidak ditemukan",
  "error": "Pendaftaran dengan ID tersebut tidak ditemukan atau bukan milik Anda",
  "type": "not_found"
}
```

**400 - Wrong Status:**
```json
{
  "success": false,
  "message": "Status tidak sesuai",
  "error": "Pendaftaran harus dalam status 'Berkas Diterima' untuk menandai kunjungan",
  "type": "validation"
}
```

---

### 2.4 Update Wedding Address

**Endpoint:** `PUT /simnikah/pendaftaran/:id/alamat`  
**Auth Required:** ✅ Yes (staff, kepala_kua)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{
  "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin"
}
```

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Alamat nikah berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "tempat_nikah": "Di Luar KUA",
    "updated_at": "2025-01-27T10:30:00Z"
  }
}
```

#### ❌ Error Responses:

**400 - Not Allowed (Nikah di KUA):**
```json
{
  "success": false,
  "message": "Alamat tidak dapat diubah",
  "error": "Alamat hanya dapat diubah untuk nikah di luar KUA",
  "type": "validation"
}
```

---

### 2.5 Get All Registrations (Staff/Kepala KUA)

**Endpoint:** `GET /simnikah/pendaftaran`  
**Auth Required:** ✅ Yes (staff, kepala_kua)

#### Query Parameters:
- `page` (optional, default: 1)
- `limit` (optional, default: 10, max: 100)
- `status` (optional): Filter by status
- `date_from` (optional): Format YYYY-MM-DD
- `date_to` (optional): Format YYYY-MM-DD
- `location` (optional): "Di KUA" or "Di Luar KUA"
- `search` (optional): Search in nomor pendaftaran, nama, NIK
- `sort_by` (optional): created_at, tanggal_nikah, status_pendaftaran, nomor_pendaftaran
- `sort_order` (optional): asc or desc

#### Example Request:
```
GET /simnikah/pendaftaran?page=1&limit=10&status=Menunggu Verifikasi&sort_by=created_at&sort_order=desc
```

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": {
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIKa1b2c3d4e5",
        "pendaftar_id": "USR1704067200",
        "status_pendaftaran": "Menunggu Verifikasi",
        "status_bimbingan": "Belum",
        "tanggal_pendaftaran": "2025-01-27T10:00:00Z",
        "tanggal_nikah": "2025-02-15",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "KUA Kecamatan Banjarmasin Utara",
        "nomor_dispensasi": "",
        "penghulu_id": null,
        "catatan": "",
        "created_at": "2025-01-27T10:00:00Z",
        "updated_at": "2025-01-27T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_records": 50,
      "per_page": 10,
      "has_next": true,
      "has_previous": false
    },
    "filters": {
      "status": "Menunggu Verifikasi",
      "date_from": "",
      "date_to": "",
      "location": "",
      "search": "",
      "sort_by": "created_at",
      "sort_order": "desc"
    }
  }
}
```

---

## 3. Staff Management

### 3.1 Create Staff (Kepala KUA Only)

**Endpoint:** `POST /simnikah/staff`  
**Auth Required:** ✅ Yes (kepala_kua only)

#### Request Body:
```json
{
  "username": "staff001",
  "email": "staff@kua.go.id",
  "password": "password123",
  "nama": "Budi Santoso",
  "nip": "197001011990011001",
  "jabatan": "Staff",
  "bagian": "Verifikasi",
  "no_hp": "081234567890",
  "alamat": "Jl. KUA No. 1, Banjarmasin"
}
```

#### Field Validation:
- `username`: Required, unique
- `email`: Required, valid email, unique
- `password`: Required, min 6 characters
- `nama`: Required
- `nip`: Required, unique
- `jabatan`: Required, one of: "Staff", "Penghulu", "Kepala KUA"
- `bagian`: Required
- `no_hp`: Optional
- `alamat`: Optional

#### ✅ Success Response (201):
```json
{
  "message": "Staff KUA berhasil dibuat",
  "data": {
    "user": {
      "user_id": "STF1704067200",
      "username": "staff001",
      "email": "staff@kua.go.id",
      "role": "staff",
      "nama": "Budi Santoso",
      "status": "Aktif"
    },
    "staff": {
      "id": 1,
      "user_id": "STF1704067200",
      "nip": "197001011990011001",
      "nama_lengkap": "Budi Santoso",
      "jabatan": "Staff",
      "bagian": "Verifikasi",
      "status": "Aktif"
    }
  }
}
```

#### ❌ Error Responses:

**403 - Forbidden (Not Kepala KUA):**
```json
{
  "error": "Akses ditolak. Hanya kepala KUA yang dapat mengakses endpoint ini"
}
```

**400 - Username/Email Exists:**
```json
{
  "error": "Username atau email sudah digunakan"
}
```

**400 - NIP Exists:**
```json
{
  "error": "NIP sudah terdaftar"
}
```

---

### 3.2 Verify Formulir (Staff Only)

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`  
**Auth Required:** ✅ Yes (staff only)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir sudah lengkap dan sesuai"
}
```

**OR**

```json
{
  "status": "Formulir Ditolak",
  "catatan": "Data tidak lengkap, silakan lengkapi data orang tua"
}
```

#### Field Validation:
- `status`: Required, "Formulir Disetujui" or "Formulir Ditolak"
- `catatan`: Optional

#### ✅ Success Response - Approved (200):
```json
{
  "success": true,
  "message": "Formulir berhasil disetujui dan status diubah ke Pengumpulan Berkas",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "status_pendaftaran": "Menunggu Pengumpulan Berkas",
    "disetujui_oleh": "STF1704067200",
    "disetujui_pada": "2025-01-27T11:00:00Z",
    "catatan": "Formulir sudah lengkap dan sesuai",
    "updated_at": "2025-01-27T11:00:00Z"
  }
}
```

#### ✅ Success Response - Rejected (200):
```json
{
  "success": true,
  "message": "Formulir berhasil diverifikasi",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "status_pendaftaran": "Formulir Ditolak",
    "disetujui_oleh": "STF1704067200",
    "disetujui_pada": "2025-01-27T11:00:00Z",
    "catatan": "Data tidak lengkap",
    "updated_at": "2025-01-27T11:00:00Z"
  }
}
```

#### ❌ Error Responses:

**400 - Wrong Status:**
```json
{
  "success": false,
  "message": "Status tidak sesuai",
  "error": "Pendaftaran harus dalam status 'Menunggu Verifikasi' untuk diverifikasi"
}
```

**400 - Invalid Status Value:**
```json
{
  "success": false,
  "message": "Status tidak valid",
  "error": "Status harus 'Formulir Disetujui' atau 'Formulir Ditolak'"
}
```

---

### 3.3 Verify Berkas (Staff Only)

**Endpoint:** `POST /simnikah/staff/verify-berkas/:id`  
**Auth Required:** ✅ Yes (staff only)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{
  "status": "Berkas Diterima",
  "catatan": "Berkas lengkap dan valid"
}
```

**OR**

```json
{
  "status": "Berkas Ditolak",
  "catatan": "KTP sudah expired, silakan perpanjang"
}
```

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Berkas berhasil diverifikasi",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "status_pendaftaran": "Berkas Diterima",
    "disetujui_oleh": "STF1704067200",
    "disetujui_pada": "2025-01-27T12:00:00Z",
    "catatan": "Berkas lengkap dan valid",
    "updated_at": "2025-01-27T12:00:00Z"
  }
}
```

#### ❌ Error Responses:

**400 - Wrong Status:**
```json
{
  "success": false,
  "message": "Status tidak sesuai",
  "error": "Pendaftaran harus dalam status 'Menunggu Pengumpulan Berkas' untuk verifikasi berkas"
}
```

---

## 4. Penghulu Management

### 4.1 Create Penghulu (Kepala KUA Only)

**Endpoint:** `POST /simnikah/penghulu`  
**Auth Required:** ✅ Yes (kepala_kua only)

#### Request Body:
```json
{
  "username": "penghulu001",
  "email": "penghulu@kua.go.id",
  "password": "password123",
  "nama": "Ustadz Haji Ahmad",
  "nip": "197001011990011002",
  "no_hp": "081234567891",
  "alamat": "Jl. Penghulu No. 1, Banjarmasin"
}
```

#### ✅ Success Response (201):
```json
{
  "message": "Penghulu berhasil dibuat",
  "data": {
    "user": {
      "user_id": "PNG1704067200",
      "username": "penghulu001",
      "email": "penghulu@kua.go.id",
      "role": "penghulu",
      "nama": "Ustadz Haji Ahmad"
    },
    "penghulu": {
      "id": 1,
      "user_id": "PNG1704067200",
      "nip": "197001011990011002",
      "nama_lengkap": "Ustadz Haji Ahmad",
      "status": "Aktif",
      "jumlah_nikah": 0,
      "rating": 0.0
    }
  }
}
```

---

### 4.2 Verify Documents (Penghulu Only)

**Endpoint:** `POST /simnikah/penghulu/verify-documents/:id`  
**Auth Required:** ✅ Yes (penghulu only)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{
  "status": "Menunggu Pelaksanaan",
  "catatan": "Dokumen sudah lengkap dan valid"
}
```

**OR**

```json
{
  "status": "Ditolak",
  "catatan": "Surat izin orang tua belum lengkap"
}
```

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Verifikasi berkas berhasil",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "status_pendaftaran": "Menunggu Bimbingan",
    "penghulu_id": 1,
    "catatan": "Dokumen sudah lengkap dan valid",
    "updated_at": "2025-01-27T13:00:00Z"
  }
}
```

#### ❌ Error Responses:

**403 - Not Assigned:**
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak ditugaskan untuk pendaftaran ini"
}
```

**400 - Wrong Status:**
```json
{
  "success": false,
  "message": "Status tidak sesuai",
  "error": "Pendaftaran harus dalam status 'Menunggu Verifikasi Penghulu' untuk diverifikasi"
}
```

---

### 4.3 Get Assigned Registrations (Penghulu Only)

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`  
**Auth Required:** ✅ Yes (penghulu only)

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": {
    "penghulu": "Ustadz Haji Ahmad",
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIKa1b2c3d4e5",
        "status_pendaftaran": "Menunggu Verifikasi Penghulu",
        "tanggal_nikah": "2025-02-15",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "KUA Kecamatan Banjarmasin Utara",
        "catatan": "",
        "created_at": "2025-01-27T10:00:00Z",
        "updated_at": "2025-01-27T13:00:00Z"
      }
    ],
    "total": 1
  }
}
```

---

## 5. Calendar & Schedule

### 5.1 Get Calendar Availability

**Endpoint:** `GET /simnikah/kalender-ketersediaan`  
**Auth Required:** ✅ Yes

#### Query Parameters:
- `bulan` (optional, default: current month): 1-12
- `tahun` (optional, default: current year)

#### Example Request:
```
GET /simnikah/kalender-ketersediaan?bulan=2&tahun=2025
```

#### ✅ Success Response (200):
```json
{
  "message": "Kalender ketersediaan berhasil diambil",
  "data": {
    "bulan": 2,
    "tahun": 2025,
    "kalender": [
      {
        "tanggal": "2025-02-01",
        "total_nikah": 3,
        "nikah_di_kua": 2,
        "nikah_di_luar_kua": 1,
        "kapasitas_kua": 9,
        "sisa_kuota": 7,
        "status": "Tersedia",
        "warna": "hijau"
      },
      {
        "tanggal": "2025-02-02",
        "total_nikah": 9,
        "nikah_di_kua": 9,
        "nikah_di_luar_kua": 0,
        "kapasitas_kua": 9,
        "sisa_kuota": 0,
        "status": "Penuh",
        "warna": "merah"
      }
    ]
  }
}
```

---

### 5.2 Get Date Availability Detail

**Endpoint:** `GET /simnikah/ketersediaan-tanggal/:tanggal`  
**Auth Required:** ✅ Yes  
**URL Parameter:** `tanggal` (format: YYYY-MM-DD)

#### Example Request:
```
GET /simnikah/ketersediaan-tanggal/2025-02-15
```

#### ✅ Success Response (200):
```json
{
  "message": "Ketersediaan tanggal berhasil diambil",
  "data": {
    "tanggal": "2025-02-15",
    "total_nikah": 3,
    "nikah_di_kua": 2,
    "nikah_di_luar_kua": 1,
    "kapasitas_kua": 9,
    "sisa_kuota": 7,
    "status": "Tersedia",
    "tersedia": true,
    "jadwal_detail": [
      {
        "nomor_pendaftaran": "NIKa1b2c3d4e5",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di KUA",
        "status_pendaftaran": "Menunggu Verifikasi Penghulu"
      }
    ]
  }
}
```

#### ❌ Error Responses:

**400 - Date in Past:**
```json
{
  "error": "Tanggal sudah lewat"
}
```

**400 - Invalid Date Format:**
```json
{
  "error": "Format tanggal tidak valid (YYYY-MM-DD)"
}
```

---

### 5.3 Assign Penghulu (Kepala KUA Only)

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`  
**Auth Required:** ✅ Yes (kepala_kua only)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{
  "penghulu_id": 1,
  "catatan": "Penghulu ditugaskan untuk memimpin nikah"
}
```

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Penghulu berhasil ditugaskan",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "penghulu_id": 1,
    "penghulu_nama": "Ustadz Haji Ahmad",
    "penghulu_assigned_by": "KKUA1704067200",
    "penghulu_assigned_at": "2025-01-27T14:00:00Z",
    "catatan": "Penghulu ditugaskan untuk memimpin nikah",
    "updated_at": "2025-01-27T14:00:00Z"
  }
}
```

#### ❌ Error Responses:

**400 - Wrong Status:**
```json
{
  "success": false,
  "message": "Status tidak sesuai",
  "error": "Pendaftaran harus dalam status 'Berkas Disetujui' untuk ditugaskan ke penghulu"
}
```

**400 - Capacity Exceeded:**
```json
{
  "error": "Kuota pernikahan harian penuh",
  "details": {
    "tanggal": "2025-02-15",
    "maksimal": 9
  }
}
```

**400 - Schedule Conflict:**
```json
{
  "error": "Konflik jadwal! Penghulu sudah memiliki jadwal pada waktu yang berdekatan",
  "details": {
    "waktu_konflik": "10:00",
    "tempat": "Di KUA",
    "selisih_menit": 30,
    "minimal_selisih": "60 menit (1 jam)"
  }
}
```

---

## 6. Bimbingan Perkawinan

### 6.1 Create Bimbingan (Staff/Kepala KUA Only)

**Endpoint:** `POST /simnikah/bimbingan`  
**Auth Required:** ✅ Yes (staff, kepala_kua only)

#### Request Body:
```json
{
  "tanggal_bimbingan": "2025-02-05",
  "waktu_mulai": "09:00",
  "waktu_selesai": "11:00",
  "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
  "pembimbing": "Ustadz Haji Ahmad",
  "kapasitas": 10,
  "catatan": "Bimbingan perkawinan untuk pasangan yang akan menikah"
}
```

#### Field Validation:
- `tanggal_bimbingan`: Required, format YYYY-MM-DD, **HARUS HARI RABU**
- `waktu_mulai`: Required, format HH:MM
- `waktu_selesai`: Required, format HH:MM
- `tempat_bimbingan`: Required
- `pembimbing`: Required
- `kapasitas`: Optional, default: 10
- `catatan`: Optional

#### ✅ Success Response (201):
```json
{
  "message": "Bimbingan perkawinan berhasil dibuat",
  "data": {
    "id": 1,
    "tanggal_bimbingan": "2025-02-05",
    "waktu_mulai": "09:00",
    "waktu_selesai": "11:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
    "pembimbing": "Ustadz Haji Ahmad",
    "kapasitas": 10,
    "status": "Aktif",
    "catatan": "Bimbingan perkawinan untuk pasangan yang akan menikah",
    "created_at": "2025-01-27T15:00:00Z",
    "updated_at": "2025-01-27T15:00:00Z"
  }
}
```

#### ❌ Error Responses:

**400 - Not Wednesday:**
```json
{
  "error": "Bimbingan perkawinan hanya bisa dijadwalkan pada hari Rabu"
}
```

**400 - Date in Past:**
```json
{
  "error": "Tanggal bimbingan tidak boleh di masa lalu"
}
```

**400 - Duplicate Date:**
```json
{
  "error": "Sudah ada bimbingan perkawinan pada tanggal tersebut"
}
```

---

### 6.2 Register for Bimbingan

**Endpoint:** `POST /simnikah/bimbingan/:id/daftar`  
**Auth Required:** ✅ Yes (user_biasa)  
**URL Parameter:** `id` (bimbingan ID)

#### Request Body:
```json
{}
```
*(Tidak ada body, menggunakan pendaftaran nikah user yang sedang login)*

#### ✅ Success Response (201):
```json
{
  "message": "Berhasil mendaftar bimbingan perkawinan",
  "data": {
    "id": 1,
    "bimbingan_perkawinan_id": 1,
    "pendaftaran_nikah_id": 1,
    "calon_suami_id": 1,
    "calon_istri_id": 2,
    "status_kehadiran": "Belum",
    "status_sertifikat": "Belum",
    "created_at": "2025-01-27T16:00:00Z"
  }
}
```

#### ❌ Error Responses:

**400 - Capacity Full:**
```json
{
  "error": "Kuota bimbingan perkawinan sudah penuh"
}
```

**400 - Already Registered:**
```json
{
  "error": "Anda sudah terdaftar untuk bimbingan ini"
}
```

**400 - No Active Registration:**
```json
{
  "error": "Anda belum memiliki pendaftaran nikah yang aktif"
}
```

---

### 6.3 Get Bimbingan Calendar

**Endpoint:** `GET /simnikah/bimbingan-kalender`  
**Auth Required:** ✅ Yes

#### Query Parameters:
- `bulan` (optional, default: current month)
- `tahun` (optional, default: current year)

#### ✅ Success Response (200):
```json
{
  "message": "Kalender bimbingan berhasil diambil",
  "data": {
    "bulan": 2,
    "tahun": 2025,
    "kalender": [
      {
        "tanggal": "2025-02-05",
        "hari": "Rabu",
        "bimbingan": {
          "id": 1,
          "waktu_mulai": "09:00",
          "waktu_selesai": "11:00",
          "tempat_bimbingan": "Aula KUA",
          "pembimbing": "Ustadz Haji Ahmad",
          "kapasitas": 10,
          "jumlah_peserta": 5,
          "sisa_kuota": 5,
          "status": "Aktif"
        }
      }
    ]
  }
}
```

---

## 7. Location & Map

### 7.1 Geocode Address (Get Coordinates)

**Endpoint:** `POST /simnikah/location/geocode`  
**Auth Required:** ✅ Yes

#### Request Body:
```json
{
  "alamat": "Jl. Pangeran Samudra No. 88, Banjarmasin, Kalimantan Selatan"
}
```

#### Field Validation:
- `alamat`: Required, min 10 characters

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Koordinat berhasil ditemukan",
  "data": {
    "alamat": "Jl. Pangeran Samudra No. 88, Banjarmasin, Kalimantan Selatan",
    "latitude": -3.3148,
    "longitude": 114.5925,
    "map_url": "https://www.google.com/maps?q=-3.3148,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3148&mlon=114.5925&zoom=16"
  }
}
```

#### ❌ Error Responses:

**400 - Address Too Short:**
```json
{
  "success": false,
  "message": "Alamat terlalu pendek",
  "error": "Masukkan alamat lengkap minimal 10 karakter"
}
```

**500 - Address Not Found:**
```json
{
  "success": false,
  "message": "Gagal mendapatkan koordinat",
  "error": "Tidak dapat menemukan lokasi untuk alamat: ...",
  "details": "..."
}
```

---

### 7.2 Reverse Geocode (Get Address from Coordinates)

**Endpoint:** `POST /simnikah/location/reverse-geocode`  
**Auth Required:** ✅ Yes

#### Request Body:
```json
{
  "latitude": -3.3148,
  "longitude": 114.5925
}
```

#### Field Validation:
- `latitude`: Required, -90 to 90
- `longitude`: Required, -180 to 180

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Alamat berhasil ditemukan",
  "data": {
    "latitude": -3.3148,
    "longitude": 114.5925,
    "alamat": "Jl. Pangeran Samudra, Banjarmasin Utara, Banjarmasin, Kalimantan Selatan, Indonesia",
    "detail": {
      "road": "Jl. Pangeran Samudra",
      "suburb": "Banjarmasin Utara",
      "city": "Banjarmasin",
      "state": "Kalimantan Selatan",
      "country": "Indonesia"
    },
    "map_url": "https://www.google.com/maps?q=-3.3148,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3148&mlon=114.5925&zoom=16"
  }
}
```

---

### 7.3 Search Address (Autocomplete)

**Endpoint:** `GET /simnikah/location/search?q=banjarmasin`  
**Auth Required:** ✅ Yes

#### Query Parameters:
- `q`: Required, min 3 characters

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Hasil pencarian alamat",
  "data": {
    "query": "banjarmasin",
    "results": [
      {
        "display_name": "Banjarmasin, Kalimantan Selatan, Indonesia",
        "latitude": "-3.3194",
        "longitude": "114.5901",
        "address": {
          "city": "Banjarmasin",
          "state": "Kalimantan Selatan",
          "country": "Indonesia"
        }
      }
    ],
    "count": 1
  }
}
```

---

### 7.4 Update Wedding Location with Coordinates

**Endpoint:** `PUT /simnikah/pendaftaran/:id/location`  
**Auth Required:** ✅ Yes (user_biasa)  
**URL Parameter:** `id` (pendaftaran ID)

#### Request Body:
```json
{
  "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
  "latitude": -3.3148,
  "longitude": 114.5925
}
```

**OR (Auto-geocoding jika koordinat tidak disediakan):**
```json
{
  "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin"
}
```

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Lokasi nikah berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "tempat_nikah": "Di Luar KUA",
    "latitude": -3.3148,
    "longitude": 114.5925,
    "map_url": "https://www.google.com/maps?q=-3.3148,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3148&mlon=114.5925&zoom=16",
    "updated_at": "2025-01-27T17:00:00Z"
  }
}
```

---

### 7.5 Get Wedding Location Detail (For Penghulu)

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`  
**Auth Required:** ✅ Yes  
**URL Parameter:** `id` (pendaftaran ID)

#### ✅ Success Response (200):
```json
{
  "success": true,
  "message": "Detail lokasi nikah berhasil diambil",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIKa1b2c3d4e5",
    "tanggal_nikah": "2025-02-15",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "latitude": -3.3148,
    "longitude": 114.5925,
    "has_coordinates": true,
    "map_url": "https://www.google.com/maps?q=-3.3148,114.5925",
    "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.3148,114.5925",
    "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.3148,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3148&mlon=114.5925&zoom=16",
    "waze_url": "https://www.waze.com/ul?ll=-3.3148,114.5925&navigate=yes",
    "is_outside_kua": true,
    "note": "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
  }
}
```

---

## 8. Notifications

### 8.1 Get User Notifications

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`  
**Auth Required:** ✅ Yes  
**URL Parameter:** `user_id`

#### Query Parameters:
- `page` (optional, default: 1)
- `limit` (optional, default: 10, max: 100)
- `status` (optional): "Belum Dibaca" or "Sudah Dibaca"
- `tipe` (optional): "Info", "Success", "Warning", "Error"

#### ✅ Success Response (200):
```json
{
  "message": "Notifikasi berhasil diambil",
  "notifications": [
    {
      "id": 1,
      "user_id": "USR1704067200",
      "judul": "Formulir Disetujui - Silakan Kumpulkan Berkas",
      "pesan": "Formulir pendaftaran nikah Anda telah disetujui...",
      "tipe": "Success",
      "status_baca": "Belum Dibaca",
      "link": "/pendaftaran/1",
      "created_at": "2025-01-27T11:00:00Z",
      "updated_at": "2025-01-27T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_page": 3
  },
  "unread_count": 5
}
```

---

### 8.2 Mark Notification as Read

**Endpoint:** `PUT /simnikah/notifikasi/:id/status`  
**Auth Required:** ✅ Yes  
**URL Parameter:** `id` (notification ID)

#### Request Body:
```json
{
  "status_baca": "Sudah Dibaca"
}
```

#### ✅ Success Response (200):
```json
{
  "message": "Status notifikasi berhasil diupdate",
  "notification": {
    "id": 1,
    "user_id": "USR1704067200",
    "judul": "Formulir Disetujui",
    "pesan": "...",
    "tipe": "Success",
    "status_baca": "Sudah Dibaca",
    "link": "/pendaftaran/1",
    "created_at": "2025-01-27T11:00:00Z",
    "updated_at": "2025-01-27T18:00:00Z"
  }
}
```

---

### 8.3 Mark All as Read

**Endpoint:** `PUT /simnikah/notifikasi/user/:user_id/mark-all-read`  
**Auth Required:** ✅ Yes  
**URL Parameter:** `user_id`

#### Request Body:
```json
{}
```

#### ✅ Success Response (200):
```json
{
  "message": "Semua notifikasi berhasil ditandai sebagai sudah dibaca",
  "updated_count": 5
}
```

---

## 9. Status Codes & Error Handling

### HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid input or validation error |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |

### Error Response Format

**Standard Error Format:**
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description",
  "type": "validation|authentication|database|not_found",
  "field": "field_name (optional)"
}
```

**Error Types:**
- `validation` - Input validation error
- `authentication` - Auth/token error
- `authorization` - Permission error
- `database` - Database error
- `not_found` - Resource not found
- `duplicate` - Duplicate entry
- `format` - Format error
- `required` - Required field missing
- `enum` - Invalid enum value
- `syariat_validation` - Islamic law validation error

### Common Error Scenarios

**1. Missing Token:**
```json
{
  "error": "Token tidak ditemukan"
}
```

**2. Invalid Token:**
```json
{
  "error": "Token tidak valid"
}
```

**3. Expired Token:**
```json
{
  "error": "Token sudah expired"
}
```

**4. Rate Limit Exceeded:**
```json
{
  "success": false,
  "message": "Rate limit exceeded",
  "error": "Terlalu banyak request. Silakan coba lagi nanti.",
  "retry_after": "60 detik"
}
```

**5. Validation Error:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Field validation error message",
  "field": "field_name",
  "type": "validation"
}
```

---

## 📝 Notes untuk Frontend Developer

### 1. Authentication Flow
1. User register → Get user data
2. User login → Get JWT token
3. Store token in localStorage/sessionStorage
4. Include token in every request: `Authorization: Bearer <token>`
5. Handle 401 errors → Redirect to login

### 2. Error Handling Best Practices
- Always check `success` field in response
- Display user-friendly error messages
- Log detailed errors for debugging
- Handle network errors separately

### 3. Form Validation
- Validate on frontend before submit
- But always handle backend validation errors
- Show field-specific errors when available

### 4. Loading States
- Show loading indicators for async operations
- Disable buttons during submission
- Handle timeout scenarios

### 5. Date Format
- Always use `YYYY-MM-DD` format
- Time format: `HH:MM` (24-hour)

### 6. Pagination
- Always check `has_next` and `has_previous`
- Implement infinite scroll or page navigation

### 7. File Upload
- Currently not implemented
- Will be added in future version

---

**Last Updated:** 27 Januari 2025  
**API Version:** 1.0.0

