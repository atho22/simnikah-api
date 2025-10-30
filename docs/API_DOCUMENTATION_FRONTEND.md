# 📚 SimNikah API Documentation - Frontend Integration Guide

**Version:** 2.0  
**Last Updated:** 30 Oktober 2025  
**Base URL:** `https://your-api.railway.app` (atau `http://localhost:8080` untuk development)

---

## 📑 Table of Contents

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [User Roles](#user-roles)
4. [API Endpoints](#api-endpoints)
   - [Authentication](#1-authentication-endpoints)
   - [Marriage Registration](#2-marriage-registration-endpoints)
   - [Staff Management](#3-staff-management-endpoints)
   - [Penghulu Management](#4-penghulu-management-endpoints)
   - [Calendar & Schedule](#5-calendar--schedule-endpoints)
   - [Counseling (Bimbingan)](#6-counseling-bimbingan-endpoints)
   - [Notifications](#7-notification-endpoints)
   - [Map & Location](#8-map--location-endpoints)
5. [Status Flow](#status-flow)
6. [Error Codes](#error-codes)
7. [Business Rules](#business-rules)
8. [Testing Guide](#testing-guide)

---

## 🚀 Getting Started

### Base URL

```
Development: http://localhost:8080
Production:  https://your-api.railway.app
```

### Headers

**Untuk endpoint yang memerlukan authentication:**
```http
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

### CORS

API sudah dikonfigurasi untuk menerima request dari:
- `http://localhost:3000` (React)
- `http://localhost:5173` (Vite)
- Frontend production URL (configure via `ALLOWED_ORIGINS`)

---

## 🔐 Authentication

### How Authentication Works

1. **Register** atau **Login** untuk mendapatkan JWT token
2. **Simpan token** di localStorage/sessionStorage
3. **Kirim token** di header `Authorization: Bearer <token>` untuk setiap request
4. Token valid selama **24 jam**

### Token Structure

```json
{
  "user_id": "USR1730268000",
  "email": "user@example.com",
  "role": "user_biasa",
  "nama": "John Doe",
  "exp": 1730354400
}
```

---

## 👥 User Roles

| Role | Access Level | Permissions |
|------|--------------|-------------|
| `user_biasa` | User | Register marriage, view own status, join counseling |
| `staff` | Staff KUA | Verify forms & documents, view all registrations |
| `penghulu` | Penghulu | Verify documents, view assigned schedules, conduct marriages |
| `kepala_kua` | Head KUA | Full access, assign penghulu, manage staff, manage counseling |

---

## 📡 API Endpoints

---

## 1. Authentication Endpoints

### 1.1 Register User

**Endpoint:** `POST /register`  
**Auth Required:** No  
**Rate Limit:** 10 requests/minute

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john.doe@example.com",
  "password": "SecurePass123!",
  "nama": "John Doe",
  "role": "user_biasa"
}
```

**Valid Roles:**
- `user_biasa` - Calon pengantin
- `staff` - Staff KUA
- `penghulu` - Penghulu
- `kepala_kua` - Kepala KUA

**Success Response (201):**
```json
{
  "message": "User berhasil dibuat",
  "user": {
    "user_id": "USR1730268000",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "nama": "John Doe",
    "role": "user_biasa"
  }
}
```

**Error Responses:**

*400 - Validation Error:*
```json
{
  "error": "Format registrasi tidak valid"
}
```

*400 - Username/Email Exists:*
```json
{
  "error": "Username atau email sudah digunakan"
}
```

---

### 1.2 Login

**Endpoint:** `POST /login`  
**Auth Required:** No  
**Rate Limit:** 10 requests/minute

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "SecurePass123!"
}
```

**Success Response (200):**
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR1730268000",
    "email": "john.doe@example.com",
    "role": "user_biasa",
    "nama": "John Doe"
  }
}
```

**Error Responses:**

*401 - Invalid Credentials:*
```json
{
  "error": "Username tidak ditemukan"
}
```

*401 - Wrong Password:*
```json
{
  "error": "Password salah"
}
```

*401 - User Inactive:*
```json
{
  "error": "User tidak aktif"
}
```

---

### 1.3 Get Profile

**Endpoint:** `GET /profile`  
**Auth Required:** Yes

**Success Response (200):**
```json
{
  "message": "Profile berhasil diambil",
  "user": {
    "user_id": "USR1730268000",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "role": "user_biasa",
    "nama": "John Doe"
  }
}
```

**Error Responses:**

*401 - Unauthorized:*
```json
{
  "error": "Token otorisasi tidak disediakan"
}
```

*404 - User Not Found:*
```json
{
  "error": "User tidak ditemukan di database"
}
```

---

## 2. Marriage Registration Endpoints

### 2.1 Create Marriage Registration (COMPLETE FORM)

**Endpoint:** `POST /simnikah/pendaftaran/form-baru`  
**Auth Required:** Yes  
**Role:** `user_biasa`  
**Performance:** ⚡ Optimized (800-1200ms)

**Description:**  
Submit lengkap data pendaftaran nikah dalam 1 API call. Includes: jadwal, calon suami, calon istri, orang tua kedua belah pihak, dan wali nikah.

**Request Body:**
```json
{
  "scheduleAndLocation": {
    "weddingLocation": "Di KUA",
    "weddingAddress": "KUA Kecamatan Banjarmasin Tengah",
    "weddingDate": "2025-12-25",
    "weddingTime": "09:00",
    "dispensationNumber": ""
  },
  "groom": {
    "groomFullName": "Ahmad Fauzan",
    "groomNik": "6371012501950001",
    "groomCitizenship": "WNI",
    "groomPassportNumber": "",
    "groomPlaceOfBirth": "Banjarmasin",
    "groomDateOfBirth": "1995-01-25",
    "groomStatus": "Belum Kawin",
    "groomReligion": "Islam",
    "groomEducation": "S1",
    "groomOccupation": "Pegawai Swasta",
    "groomOccupationDescription": "Staff IT di PT. Digital Indonesia",
    "groomPhoneNumber": "081234567890",
    "groomEmail": "ahmad.fauzan@example.com",
    "groomAddress": "Jl. Lambung Mangkurat No. 45, Banjarmasin"
  },
  "bride": {
    "brideFullName": "Siti Aminah",
    "brideNik": "6371016702980001",
    "brideCitizenship": "WNI",
    "bridePassportNumber": "",
    "bridePlaceOfBirth": "Banjarmasin",
    "brideDateOfBirth": "1998-02-27",
    "brideStatus": "Belum Kawin",
    "brideReligion": "Islam",
    "brideEducation": "S1",
    "brideOccupation": "Guru",
    "brideOccupationDescription": "Guru SD Negeri 1 Banjarmasin",
    "bridePhoneNumber": "081298765432",
    "brideEmail": "siti.aminah@example.com",
    "brideAddress": "Jl. A. Yani Km 5 No. 12, Banjarmasin"
  },
  "groomParents": {
    "groomFather": {
      "groomFatherPresenceStatus": "Hidup",
      "groomFatherName": "Muhammad Ali",
      "groomFatherNik": "6371010101700001",
      "groomFatherCitizenship": "WNI",
      "groomFatherCountryOfOrigin": "Indonesia",
      "groomFatherPassportNumber": "",
      "groomFatherPlaceOfBirth": "Banjarmasin",
      "groomFatherDateOfBirth": "1970-01-01",
      "groomFatherReligion": "Islam",
      "groomFatherOccupation": "Wiraswasta",
      "groomFatherOccupationDescription": "Pedagang",
      "groomFatherAddress": "Jl. Veteran No. 20, Banjarmasin"
    },
    "groomMother": {
      "groomMotherPresenceStatus": "Hidup",
      "groomMotherName": "Fatimah",
      "groomMotherNik": "6371014501720001",
      "groomMotherCitizenship": "WNI",
      "groomMotherCountryOfOrigin": "Indonesia",
      "groomMotherPassportNumber": "",
      "groomMotherPlaceOfBirth": "Banjarmasin",
      "groomMotherDateOfBirth": "1972-05-01",
      "groomMotherReligion": "Islam",
      "groomMotherOccupation": "Ibu Rumah Tangga",
      "groomMotherOccupationDescription": "",
      "groomMotherAddress": "Jl. Veteran No. 20, Banjarmasin"
    }
  },
  "brideParents": {
    "brideFather": {
      "brideFatherPresenceStatus": "Hidup",
      "brideFatherName": "Abdullah Rahman",
      "brideFatherNik": "6371015501680001",
      "brideFatherCitizenship": "WNI",
      "brideFatherCountryOfOrigin": "Indonesia",
      "brideFatherPassportNumber": "",
      "brideFatherPlaceOfBirth": "Banjarmasin",
      "brideFatherDateOfBirth": "1968-05-15",
      "brideFatherReligion": "Islam",
      "brideFatherOccupation": "PNS",
      "brideFatherOccupationDescription": "Guru SMA",
      "brideFatherAddress": "Jl. A. Yani Km 5 No. 12, Banjarmasin"
    },
    "brideMother": {
      "brideMotherPresenceStatus": "Hidup",
      "brideMotherName": "Khadijah",
      "brideMotherNik": "6371012201700001",
      "brideMotherCitizenship": "WNI",
      "brideMotherCountryOfOrigin": "Indonesia",
      "brideMotherPassportNumber": "",
      "brideMotherPlaceOfBirth": "Banjarmasin",
      "brideMotherDateOfBirth": "1970-01-22",
      "brideMotherReligion": "Islam",
      "brideMotherOccupation": "Ibu Rumah Tangga",
      "brideMotherOccupationDescription": "",
      "brideMotherAddress": "Jl. A. Yani Km 5 No. 12, Banjarmasin"
    }
  },
  "guardian": {
    "guardianFullName": "Abdullah Rahman",
    "guardianNik": "6371015501680001",
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Hidup",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. A. Yani Km 5 No. 12, Banjarmasin",
    "guardianPhoneNumber": "081234567899"
  }
}
```

**Field Validations:**

**weddingLocation:**
- `"Di KUA"` atau `"Di Luar KUA"`

**groomStatus / brideStatus:**
- `"Belum Kawin"`
- `"Cerai Mati"`
- `"Cerai Hidup"`

**groomFatherPresenceStatus / brideFatherPresenceStatus:**
- `"Hidup"`
- `"Meninggal"`

**guardianRelationship (Urutan Wali Nasab):**
- `"Ayah Kandung"` (prioritas tertinggi jika masih hidup)
- `"Kakek"`
- `"Saudara Laki-Laki Kandung"`
- `"Saudara Laki-Laki Seayah"`
- `"Keponakan Laki-Laki"`
- `"Paman Kandung"`
- `"Paman Seayah"`
- `"Sepupu Laki-Laki"`
- `"Wali Hakim"`

**guardianStatus:**
- `"Hidup"`
- `"Meninggal"` (hanya untuk wali non-ayah kandung)

**Success Response (201):**
```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat",
  "data": {
    "pendaftaran": {
      "id": 1,
      "nomor_pendaftaran": "REG/2025/001",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2025-12-25T00:00:00Z",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "KUA Kecamatan Banjarmasin Tengah"
    },
    "calon_suami": {
      "id": "CS001",
      "nama_lengkap": "Ahmad Fauzan",
      "nik": "6371012501950001"
    },
    "calon_istri": {
      "id": "CI001",
      "nama_lengkap": "Siti Aminah",
      "nik": "6371016702980001"
    },
    "wali_nikah": {
      "id": 1,
      "nama_lengkap": "Abdullah Rahman",
      "hubungan_wali": "Ayah Kandung",
      "status_keberadaan": "Hidup"
    }
  },
  "next_steps": [
    "✅ Formulir berhasil disimpan",
    "⏳ Menunggu verifikasi dari staff KUA",
    "📋 Siapkan berkas fisik untuk diserahkan ke KUA",
    "📱 Pantau status melalui menu 'Status Pendaftaran'"
  ]
}
```

**Error Responses:**

*400 - Validation Error:*
```json
{
  "success": false,
  "error": "Validasi gagal",
  "details": "Tanggal nikah tidak boleh di masa lalu"
}
```

*400 - Guardian Validation Error:*
```json
{
  "success": false,
  "error": "Validasi wali nikah gagal",
  "details": "Wali nikah tidak boleh dalam status Meninggal",
  "validation_rule": "Validasi 1: Wali nikah HARUS hidup"
}
```

*400 - NIK Already Exists:*
```json
{
  "success": false,
  "error": "NIK sudah terdaftar",
  "details": "NIK 6371012501950001 sudah digunakan"
}
```

**Guardian Validation Rules (Automatic):**

1. ✅ **Wali nikah HARUS hidup**
2. ✅ **Konsistensi status wali dengan ayah kandung**
3. ✅ **NIK wali = NIK ayah** (jika wali adalah ayah kandung)
4. ✅ **Wali tidak boleh sama dengan calon pengantin**

---

### 2.2 Check Registration Status

**Endpoint:** `GET /simnikah/pendaftaran/status`  
**Auth Required:** Yes  
**Role:** `user_biasa`

**Description:**  
Check apakah user sudah pernah mendaftar dan status pendaftarannya.

**Success Response (200) - Sudah Terdaftar:**
```json
{
  "has_registration": true,
  "registration": {
    "id": 1,
    "nomor_pendaftaran": "REG/2025/001",
    "status_pendaftaran": "Menunggu Verifikasi",
    "tanggal_nikah": "2025-12-25",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "alamat_akad": "KUA Kecamatan Banjarmasin Tengah",
    "created_at": "2025-10-30T10:00:00Z"
  }
}
```

**Success Response (200) - Belum Terdaftar:**
```json
{
  "has_registration": false,
  "message": "Belum ada pendaftaran nikah"
}
```

---

### 2.3 Mark As Visited (Serahkan Berkas Fisik)

**Endpoint:** `POST /simnikah/pendaftaran/:id/mark-visited`  
**Auth Required:** Yes  
**Role:** `user_biasa`

**Description:**  
User mengkonfirmasi sudah datang ke KUA untuk menyerahkan berkas fisik.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Status berhasil diperbarui menjadi 'Menunggu Penugasan'",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG/2025/001",
    "status_pendaftaran": "Menunggu Penugasan",
    "status_sebelumnya": "Berkas Diterima"
  }
}
```

**Error Responses:**

*400 - Wrong Status:*
```json
{
  "success": false,
  "error": "Hanya pendaftaran dengan status 'Berkas Diterima' yang dapat ditandai sebagai visited"
}
```

---

### 2.4 Get All Registrations (Staff/Kepala KUA)

**Endpoint:** `GET /simnikah/pendaftaran`  
**Auth Required:** Yes  
**Role:** `staff`, `kepala_kua`

**Query Parameters:**
- `status` (optional) - Filter by status
- `search` (optional) - Search by name or NIK
- `page` (optional, default: 1)
- `limit` (optional, default: 20)

**Example:**
```
GET /simnikah/pendaftaran?status=Menunggu Verifikasi&search=Ahmad&page=1&limit=10
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "total": 50,
    "page": 1,
    "limit": 10,
    "total_pages": 5,
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG/2025/001",
        "status_pendaftaran": "Menunggu Verifikasi",
        "tanggal_nikah": "2025-12-25",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "calon_suami": {
          "nama_lengkap": "Ahmad Fauzan",
          "nik": "6371012501950001"
        },
        "calon_istri": {
          "nama_lengkap": "Siti Aminah",
          "nik": "6371016702980001"
        },
        "created_at": "2025-10-30T10:00:00Z"
      }
    ]
  }
}
```

---

### 2.5 Get Status Flow

**Endpoint:** `GET /simnikah/pendaftaran/:id/status-flow`  
**Auth Required:** Yes

**Description:**  
Mendapatkan alur status pendaftaran dengan indikator mana yang sudah completed dan mana yang current.

**Success Response (200):**
```json
{
  "message": "Alur status pendaftaran berhasil diambil",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "REG/2025/001",
    "status_sekarang": "Menunggu Verifikasi",
    "tanggal_nikah": "2025-12-25",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "penghulu_assigned": false,
    "bimbingan_info": null,
    "status_flow": [
      {
        "status": "Draft",
        "description": "Data belum lengkap",
        "completed": true,
        "current": false,
        "can_edit": true
      },
      {
        "status": "Menunggu Verifikasi",
        "description": "Staff verifikasi formulir online",
        "completed": false,
        "current": true,
        "can_edit": false
      },
      {
        "status": "Menunggu Pengumpulan Berkas",
        "description": "Formulir disetujui, siap kumpulkan berkas",
        "completed": false,
        "current": false,
        "can_edit": false
      }
      // ... 7 status lainnya
    ]
  }
}
```

---

### 2.6 Update Wedding Address

**Endpoint:** `PUT /simnikah/pendaftaran/:id/alamat`  
**Auth Required:** Yes  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
  "latitude": -3.3194374,
  "longitude": 114.5900675
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Alamat nikah berhasil diperbarui",
  "data": {
    "id": 1,
    "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "latitude": -3.3194374,
    "longitude": 114.5900675
  }
}
```

---

## 3. Staff Management Endpoints

### 3.1 Create Staff

**Endpoint:** `POST /simnikah/staff`  
**Auth Required:** Yes  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "user_id": "USR1730268000",
  "nip": "199001012020121001",
  "nama_lengkap": "Budi Santoso",
  "jabatan": "Staff",
  "bagian": "Verifikasi",
  "no_hp": "081234567890",
  "email": "budi.santoso@kua.go.id",
  "alamat": "Jl. Jenderal Sudirman No. 10, Banjarmasin"
}
```

**Success Response (201):**
```json
{
  "message": "Staff KUA berhasil dibuat",
  "data": {
    "id": 1,
    "user_id": "USR1730268000",
    "nip": "199001012020121001",
    "nama_lengkap": "Budi Santoso",
    "jabatan": "Staff",
    "bagian": "Verifikasi",
    "status": "Aktif"
  }
}
```

---

### 3.2 Get All Staff

**Endpoint:** `GET /simnikah/staff`  
**Auth Required:** Yes  
**Role:** `kepala_kua`

**Success Response (200):**
```json
{
  "message": "Data staff KUA berhasil diambil",
  "data": {
    "total": 10,
    "staff": [
      {
        "id": 1,
        "user_id": "USR1730268000",
        "nip": "199001012020121001",
        "nama_lengkap": "Budi Santoso",
        "jabatan": "Staff",
        "bagian": "Verifikasi",
        "status": "Aktif",
        "email": "budi.santoso@kua.go.id",
        "no_hp": "081234567890"
      }
    ]
  }
}
```

---

### 3.3 Verify Formulir (Online)

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`  
**Auth Required:** Yes  
**Role:** `staff`

**Request Body:**
```json
{
  "approved": true,
  "catatan": "Formulir sudah lengkap dan sesuai"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Formulir berhasil diverifikasi",
  "data": {
    "pendaftaran_id": 1,
    "status_baru": "Menunggu Pengumpulan Berkas",
    "verified_by": "USR1730268000",
    "verified_at": "2025-10-30T14:30:00Z"
  }
}
```

**If Rejected (approved: false):**
```json
{
  "success": true,
  "message": "Formulir ditolak",
  "data": {
    "pendaftaran_id": 1,
    "status_baru": "Ditolak",
    "catatan": "Data NIK tidak valid, silakan perbaiki"
  }
}
```

---

### 3.4 Verify Berkas (Fisik)

**Endpoint:** `POST /simnikah/staff/verify-berkas/:id`  
**Auth Required:** Yes  
**Role:** `staff`

**Request Body:**
```json
{
  "approved": true,
  "catatan": "Berkas lengkap dan sesuai dengan data online"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Berkas berhasil diverifikasi",
  "data": {
    "pendaftaran_id": 1,
    "status_baru": "Berkas Diterima",
    "verified_by": "USR1730268000",
    "verified_at": "2025-10-31T10:00:00Z"
  }
}
```

---

## 4. Penghulu Management Endpoints

### 4.1 Create Penghulu

**Endpoint:** `POST /simnikah/penghulu`  
**Auth Required:** Yes  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "user_id": "USR1730268100",
  "nip": "198505052010121001",
  "nama_lengkap": "Ustadz Ahmad Ridho",
  "no_hp": "081234567891",
  "email": "ahmad.ridho@kua.go.id",
  "alamat": "Jl. Ahmad Yani No. 25, Banjarmasin"
}
```

**Success Response (201):**
```json
{
  "message": "Penghulu berhasil dibuat",
  "data": {
    "id": 1,
    "user_id": "USR1730268100",
    "nip": "198505052010121001",
    "nama_lengkap": "Ustadz Ahmad Ridho",
    "status": "Aktif",
    "jumlah_nikah": 0,
    "rating": 0
  }
}
```

---

### 4.2 Get All Penghulu

**Endpoint:** `GET /simnikah/penghulu`  
**Auth Required:** Yes

**Success Response (200):**
```json
{
  "message": "Data penghulu berhasil diambil",
  "data": {
    "total": 4,
    "penghulu": [
      {
        "id": 1,
        "nama_lengkap": "Ustadz Ahmad Ridho",
        "status": "Aktif",
        "jumlah_nikah": 15,
        "rating": 4.8,
        "email": "ahmad.ridho@kua.go.id",
        "no_hp": "081234567891"
      }
    ]
  }
}
```

---

### 4.3 Assign Penghulu

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`  
**Auth Required:** Yes  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "penghulu_id": 1
}
```

**Success Response (200):**
```json
{
  "message": "Penghulu berhasil diassign",
  "warning": "",
  "data": {
    "pendaftaran_id": 1,
    "penghulu_id": 1,
    "penghulu_nama": "Ustadz Ahmad Ridho",
    "assigned_by": "USR1730268200",
    "assigned_at": "2025-11-01T08:00:00Z"
  }
}
```

**Error Responses:**

*400 - Schedule Conflict:*
```json
{
  "error": "Konflik jadwal! Penghulu sudah memiliki jadwal pada waktu yang berdekatan",
  "details": {
    "waktu_konflik": "09:00",
    "tempat": "Di Luar KUA",
    "selisih_menit": 30,
    "minimal_selisih": "60 menit (1 jam)"
  }
}
```

*400 - Daily Quota Full:*
```json
{
  "error": "Kuota pernikahan harian penuh",
  "details": {
    "tanggal": "2025-12-25",
    "maksimal": 9
  }
}
```

---

### 4.4 Get Penghulu Schedule

**Endpoint:** `GET /simnikah/penghulu-jadwal/:tanggal`  
**Auth Required:** Yes

**Example:**
```
GET /simnikah/penghulu-jadwal/2025-12-25
```

**Success Response (200):**
```json
{
  "message": "Jadwal penghulu berhasil diambil",
  "data": {
    "tanggal": "2025-12-25",
    "total_penghulu": 4,
    "total_kapasitas": 12,
    "total_terisi": 5,
    "total_sisa": 7,
    "penghulu": [
      {
        "id": 1,
        "nama": "Ustadz Ahmad Ridho",
        "status": "Sebagian",
        "jumlah_jadwal": 2,
        "sisa_kuota": 1,
        "maksimal": 3,
        "jadwal": [
          {
            "nomor_pendaftaran": "REG/2025/001",
            "waktu_nikah": "09:00",
            "tempat_nikah": "Di KUA",
            "alamat_akad": "KUA Kecamatan Banjarmasin Tengah"
          }
        ]
      }
    ]
  }
}
```

---

### 4.5 Verify Documents (Penghulu)

**Endpoint:** `POST /simnikah/penghulu/verify-documents/:id`  
**Auth Required:** Yes  
**Role:** `penghulu`

**Request Body:**
```json
{
  "approved": true,
  "catatan": "Dokumen sudah sesuai dan lengkap"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Dokumen berhasil diverifikasi oleh penghulu",
  "data": {
    "pendaftaran_id": 1,
    "status_baru": "Menunggu Bimbingan",
    "verified_by": "USR1730268100",
    "verified_at": "2025-11-05T09:00:00Z"
  }
}
```

---

### 4.6 Get Assigned Registrations (Penghulu)

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`  
**Auth Required:** Yes  
**Role:** `penghulu`

**Success Response (200):**
```json
{
  "message": "Pendaftaran yang ditugaskan berhasil diambil",
  "data": {
    "total": 5,
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG/2025/001",
        "status_pendaftaran": "Menunggu Verifikasi Penghulu",
        "tanggal_nikah": "2025-12-25",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "calon_suami": {
          "nama_lengkap": "Ahmad Fauzan",
          "nik": "6371012501950001"
        },
        "calon_istri": {
          "nama_lengkap": "Siti Aminah",
          "nik": "6371016702980001"
        }
      }
    ]
  }
}
```

---

## 5. Calendar & Schedule Endpoints

### 5.1 Get Availability Calendar

**Endpoint:** `GET /simnikah/kalender-ketersediaan`  
**Auth Required:** Yes

**Query Parameters:**
- `bulan` (optional, default: current month) - Format: 01-12
- `tahun` (optional, default: current year) - Format: 2025

**Example:**
```
GET /simnikah/kalender-ketersediaan?bulan=12&tahun=2025
```

**Success Response (200):**
```json
{
  "message": "Kalender ketersediaan berhasil diambil",
  "data": {
    "bulan": 12,
    "tahun": 2025,
    "nama_bulan": "December",
    "kapasitas_harian": 9,
    "penghulu_info": {
      "total_penghulu": 4,
      "penghulu_aktif": 4,
      "penghulu_cadangan": 0,
      "slot_waktu_per_hari": 9,
      "nikah_per_slot": 4,
      "total_kapasitas_harian": 9
    },
    "kalender": [
      {
        "tanggal": 1,
        "tanggal_str": "2025-12-01",
        "status": "Tersedia",
        "tersedia": true,
        "jumlah_nikah_total": 2,
        "jumlah_nikah_kua": 2,
        "jumlah_nikah_luar": 0,
        "kuning_count": 0,
        "hijau_count": 2,
        "warna": "hijau",
        "sisa_kuota_kua": 7,
        "kapasitas_kua": 9
      },
      {
        "tanggal": 25,
        "tanggal_str": "2025-12-25",
        "status": "Penuh",
        "tersedia": false,
        "jumlah_nikah_total": 9,
        "jumlah_nikah_kua": 9,
        "jumlah_nikah_luar": 0,
        "kuning_count": 0,
        "hijau_count": 9,
        "warna": "hijau",
        "sisa_kuota_kua": 0,
        "kapasitas_kua": 9
      }
    ]
  }
}
```

**Legend:**
- **warna**: `"hijau"` = ada yang sudah fix, `"kuning"` = masih proses awal
- **status**: `"Tersedia"`, `"Penuh"`, `"Terlewat"`

---

### 5.2 Get Date Detail

**Endpoint:** `GET /simnikah/kalender-tanggal-detail`  
**Auth Required:** Yes

**Query Parameters:**
- `tanggal` (required) - Format: YYYY-MM-DD

**Example:**
```
GET /simnikah/kalender-tanggal-detail?tanggal=2025-12-25
```

**Success Response (200):**
```json
{
  "message": "Detail kalender berhasil diambil",
  "data": {
    "tanggal": "2025-12-25",
    "items": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG/2025/001",
        "waktu": "09:00",
        "tempat": "Di KUA",
        "status": "Sudah Bimbingan",
        "warna": "hijau",
        "penghulu_id": 1,
        "penghulu_nama": "Ustadz Ahmad Ridho",
        "nama_calon_suami": "Ahmad Fauzan",
        "nama_calon_istri": "Siti Aminah"
      }
    ]
  }
}
```

---

### 5.3 Get Date Availability

**Endpoint:** `GET /simnikah/ketersediaan-tanggal/:tanggal`  
**Auth Required:** Yes

**Example:**
```
GET /simnikah/ketersediaan-tanggal/2025-12-25
```

**Success Response (200):**
```json
{
  "message": "Detail ketersediaan tanggal berhasil diambil",
  "data": {
    "tanggal": "2025-12-25",
    "status": "Tersedia",
    "tersedia": true,
    "jumlah_nikah_kua": 5,
    "jumlah_nikah_luar": 3,
    "total_nikah": 8,
    "sisa_kuota_kua": 4,
    "kapasitas_kua": 9,
    "keterangan": "Kapasitas hanya berlaku untuk nikah di KUA. Nikah di luar KUA tidak dibatasi.",
    "jadwal_detail": [
      {
        "nomor_pendaftaran": "REG/2025/001",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "KUA Kecamatan Banjarmasin Tengah"
      }
    ]
  }
}
```

---

## 6. Counseling (Bimbingan) Endpoints

### 6.1 Create Counseling Session

**Endpoint:** `POST /simnikah/bimbingan`  
**Auth Required:** Yes  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "tanggal_bimbingan": "2025-11-13",
  "waktu_mulai": "08:00",
  "waktu_selesai": "12:00",
  "tempat_bimbingan": "Aula KUA Banjarmasin Tengah",
  "pembimbing": "Ustadz Abdullah Hasan",
  "kapasitas": 10,
  "catatan": "Membawa alat tulis dan dress code sopan"
}
```

**Validation:**
- `tanggal_bimbingan` harus hari **Rabu**
- Tidak boleh masa lalu
- Hanya 1 sesi per hari Rabu

**Success Response (201):**
```json
{
  "message": "Bimbingan perkawinan berhasil dibuat",
  "data": {
    "id": 1,
    "tanggal_bimbingan": "2025-11-13T00:00:00Z",
    "waktu_mulai": "08:00",
    "waktu_selesai": "12:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin Tengah",
    "pembimbing": "Ustadz Abdullah Hasan",
    "kapasitas": 10,
    "status": "Aktif"
  }
}
```

---

### 6.2 Get Counseling Sessions

**Endpoint:** `GET /simnikah/bimbingan`  
**Auth Required:** Yes

**Query Parameters:**
- `bulan` (optional, default: current month)
- `tahun` (optional, default: current year)
- `status` (optional, default: "Aktif")

**Example:**
```
GET /simnikah/bimbingan?bulan=11&tahun=2025&status=Aktif
```

**Success Response (200):**
```json
{
  "message": "Data bimbingan perkawinan berhasil diambil",
  "data": {
    "bulan": 11,
    "tahun": 2025,
    "bimbingan": [
      {
        "id": 1,
        "tanggal_bimbingan": "2025-11-13",
        "waktu_mulai": "08:00",
        "waktu_selesai": "12:00",
        "tempat_bimbingan": "Aula KUA Banjarmasin Tengah",
        "pembimbing": "Ustadz Abdullah Hasan",
        "kapasitas": 10,
        "jumlah_peserta": 7,
        "sisa_kuota": 3,
        "status": "Aktif",
        "catatan": "Membawa alat tulis dan dress code sopan"
      }
    ]
  }
}
```

---

### 6.3 Register for Counseling

**Endpoint:** `POST /simnikah/bimbingan/:id/daftar`  
**Auth Required:** Yes  
**Role:** `user_biasa`

**Description:**  
User mendaftar untuk mengikuti bimbingan perkawinan. Hanya bisa jika status pendaftaran = "Menunggu Bimbingan".

**Success Response (201):**
```json
{
  "message": "Berhasil mendaftar bimbingan perkawinan",
  "data": {
    "bimbingan_id": 1,
    "tanggal": "2025-11-13",
    "waktu": "08:00 - 12:00",
    "tempat": "Aula KUA Banjarmasin Tengah",
    "pembimbing": "Ustadz Abdullah Hasan"
  }
}
```

**Error Responses:**

*400 - Already Registered:*
```json
{
  "error": "Anda sudah terdaftar di bimbingan perkawinan ini"
}
```

*400 - Full Capacity:*
```json
{
  "error": "Bimbingan perkawinan sudah penuh"
}
```

*400 - Wrong Status:*
```json
{
  "error": "Anda belum memiliki pendaftaran nikah yang siap untuk bimbingan"
}
```

---

### 6.4 Get Counseling Participants

**Endpoint:** `GET /simnikah/bimbingan/:id/participants`  
**Auth Required:** Yes  
**Role:** `staff`, `kepala_kua`

**Success Response (200):**
```json
{
  "message": "Data peserta bimbingan perkawinan berhasil diambil",
  "data": {
    "bimbingan_id": 1,
    "tanggal": "2025-11-13",
    "waktu": "08:00 - 12:00",
    "tempat": "Aula KUA Banjarmasin Tengah",
    "pembimbing": "Ustadz Abdullah Hasan",
    "kapasitas": 10,
    "jumlah_peserta": 7,
    "peserta": [
      {
        "id": 1,
        "pendaftaran_nikah_id": 1,
        "calon_suami": {
          "nama": "Ahmad Fauzan",
          "nik": "6371012501950001"
        },
        "calon_istri": {
          "nama": "Siti Aminah",
          "nik": "6371016702980001"
        },
        "status_kehadiran": "Belum",
        "status_sertifikat": "Belum",
        "no_sertifikat": null,
        "created_at": "2025-10-30T10:00:00Z"
      }
    ]
  }
}
```

---

### 6.5 Update Attendance

**Endpoint:** `PUT /simnikah/bimbingan/:id/update-attendance`  
**Auth Required:** Yes  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "pendaftaran_nikah_id": 1,
  "status_kehadiran": "Hadir",
  "status_sertifikat": "Sudah",
  "no_sertifikat": "CERT/2025/001"
}
```

**Valid Values:**
- `status_kehadiran`: `"Hadir"`, `"Tidak Hadir"`
- `status_sertifikat`: `"Belum"`, `"Sudah"`

**Success Response (200):**
```json
{
  "message": "Kehadiran bimbingan berhasil diupdate",
  "data": {
    "bimbingan_id": 1,
    "pendaftaran_id": 1,
    "status_kehadiran": "Hadir",
    "status_sertifikat": "Sudah",
    "no_sertifikat": "CERT/2025/001",
    "updated_by": "USR1730268000",
    "updated_at": "2025-11-13T12:00:00Z"
  }
}
```

---

### 6.6 Get Counseling Calendar

**Endpoint:** `GET /simnikah/bimbingan-kalender`  
**Auth Required:** Yes

**Query Parameters:**
- `bulan` (optional)
- `tahun` (optional)

**Example:**
```
GET /simnikah/bimbingan-kalender?bulan=11&tahun=2025
```

**Success Response (200):**
```json
{
  "message": "Kalender bimbingan perkawinan berhasil diambil",
  "data": {
    "bulan": 11,
    "tahun": 2025,
    "nama_bulan": "November",
    "kalender": [
      {
        "tanggal": 13,
        "tanggal_str": "2025-11-13",
        "status": "Tersedia",
        "tersedia": true,
        "sisa_kuota": 3,
        "bimbingan": {
          "id": 1,
          "waktu_mulai": "08:00",
          "waktu_selesai": "12:00",
          "tempat_bimbingan": "Aula KUA Banjarmasin Tengah",
          "pembimbing": "Ustadz Abdullah Hasan",
          "kapasitas": 10,
          "jumlah_peserta": 7
        }
      },
      {
        "tanggal": 20,
        "tanggal_str": "2025-11-20",
        "status": "Belum Dijadwalkan",
        "tersedia": false,
        "sisa_kuota": 0,
        "bimbingan": null
      }
    ]
  }
}
```

---

## 7. Notification Endpoints

### 7.1 Get User Notifications

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`  
**Auth Required:** Yes

**Query Parameters:**
- `status` (optional) - Filter: `"Belum Dibaca"` atau `"Sudah Dibaca"`
- `limit` (optional, default: 20)

**Example:**
```
GET /simnikah/notifikasi/user/USR1730268000?status=Belum Dibaca&limit=10
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "total": 5,
    "unread_count": 3,
    "notifications": [
      {
        "id": 1,
        "judul": "Formulir Diverifikasi",
        "pesan": "Formulir pendaftaran Anda telah diverifikasi oleh staff KUA. Silakan datang ke KUA untuk menyerahkan berkas fisik.",
        "tipe": "Success",
        "status_baca": "Belum Dibaca",
        "link": "/simnikah/pendaftaran/1",
        "created_at": "2025-10-30T14:30:00Z"
      }
    ]
  }
}
```

---

### 7.2 Mark Notification as Read

**Endpoint:** `PUT /simnikah/notifikasi/:id/status`  
**Auth Required:** Yes

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Status notifikasi berhasil diupdate",
  "data": {
    "id": 1,
    "status_baca": "Sudah Dibaca"
  }
}
```

---

### 7.3 Mark All as Read

**Endpoint:** `PUT /simnikah/notifikasi/user/:user_id/mark-all-read`  
**Auth Required:** Yes

**Success Response (200):**
```json
{
  "success": true,
  "message": "Semua notifikasi berhasil ditandai sebagai dibaca",
  "data": {
    "updated_count": 5
  }
}
```

---

### 7.4 Get Notification Stats

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id/stats`  
**Auth Required:** Yes

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "total": 10,
    "unread": 3,
    "read": 7,
    "by_type": {
      "Info": 4,
      "Success": 3,
      "Warning": 2,
      "Error": 1
    }
  }
}
```

---

## 8. Map & Location Endpoints

### 8.1 Geocoding (Address → Coordinates)

**Endpoint:** `POST /simnikah/location/geocode`  
**Auth Required:** Yes

**Description:**  
Convert alamat menjadi koordinat (latitude, longitude).

**Request Body:**
```json
{
  "address": "Jl. Lambung Mangkurat No. 45, Banjarmasin, Kalimantan Selatan"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Koordinat berhasil didapatkan",
  "data": {
    "address": "Jl. Lambung Mangkurat No. 45, Banjarmasin, Kalimantan Selatan",
    "latitude": -3.3194374,
    "longitude": 114.5900675,
    "source": "OpenStreetMap Nominatim API (FREE)",
    "map_url": "https://www.openstreetmap.org/?mlat=-3.319437&mlon=114.590068&zoom=15"
  }
}
```

---

### 8.2 Reverse Geocoding (Coordinates → Address)

**Endpoint:** `POST /simnikah/location/reverse-geocode`  
**Auth Required:** Yes

**Request Body:**
```json
{
  "latitude": -3.3194374,
  "longitude": 114.5900675
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Alamat berhasil didapatkan dari koordinat",
  "data": {
    "latitude": -3.3194374,
    "longitude": 114.5900675,
    "address": "Jalan Lambung Mangkurat, Banjarmasin Tengah, Kalimantan Selatan, Indonesia",
    "source": "OpenStreetMap Nominatim API (FREE)"
  }
}
```

---

### 8.3 Search Address (Autocomplete)

**Endpoint:** `GET /simnikah/location/search`  
**Auth Required:** Yes

**Query Parameters:**
- `q` (required) - Search query
- `limit` (optional, default: 5)

**Example:**
```
GET /simnikah/location/search?q=Lambung Mangkurat Banjarmasin&limit=5
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Hasil pencarian alamat",
  "data": {
    "query": "Lambung Mangkurat Banjarmasin",
    "results": [
      {
        "display_name": "Jalan Lambung Mangkurat, Banjarmasin Tengah, Kalimantan Selatan, Indonesia",
        "latitude": -3.3194374,
        "longitude": 114.5900675,
        "type": "road",
        "importance": 0.8
      }
    ],
    "count": 5
  }
}
```

---

### 8.4 Get Wedding Location Detail

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`  
**Auth Required:** Yes

**Description:**  
Mendapatkan detail lokasi nikah beserta koordinat dan link navigasi.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Detail lokasi nikah berhasil diambil",
  "data": {
    "pendaftaran_id": 1,
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "latitude": -3.3194374,
    "longitude": 114.5900675,
    "maps": {
      "google_maps": "https://www.google.com/maps?q=-3.319437,114.590068",
      "waze": "https://waze.com/ul?ll=-3.319437,114.590068",
      "openstreetmap": "https://www.openstreetmap.org/?mlat=-3.319437&mlon=114.590068&zoom=15"
    }
  }
}
```

---

## 🔄 Status Flow

### Marriage Registration Status Flow (10 Steps)

```
1. Draft
   ↓ (User submit form)
2. Menunggu Verifikasi
   ↓ (Staff verify formulir online - approved)
3. Menunggu Pengumpulan Berkas
   ↓ (User datang ke KUA bawa berkas)
4. Berkas Diterima
   ↓ (Staff verify berkas fisik - approved)
5. Menunggu Penugasan
   ↓ (User confirm: mark as visited)
6. Menunggu Penugasan (siap assign)
   ↓ (Kepala KUA assign penghulu)
7. Penghulu Ditugaskan
   ↓ (Automatically transitions)
8. Menunggu Verifikasi Penghulu
   ↓ (Penghulu verify - approved)
9. Menunggu Bimbingan
   ↓ (User daftar & ikut bimbingan)
10. Sudah Bimbingan
    ↓ (Staff/Kepala KUA mark as complete)
11. Selesai ✅
```

### Status Transitions Table

| Current Status | Action | Next Status | Who |
|---------------|--------|-------------|-----|
| Draft | Submit form | Menunggu Verifikasi | User |
| Menunggu Verifikasi | Approve formulir | Menunggu Pengumpulan Berkas | Staff |
| Menunggu Verifikasi | Reject formulir | Ditolak | Staff |
| Menunggu Pengumpulan Berkas | User datang ke KUA | Berkas Diterima | Staff |
| Berkas Diterima | Verify berkas | Menunggu Penugasan | Staff |
| Menunggu Penugasan | Assign penghulu | Penghulu Ditugaskan | Kepala KUA |
| Penghulu Ditugaskan | Auto transition | Menunggu Verifikasi Penghulu | System |
| Menunggu Verifikasi Penghulu | Approve documents | Menunggu Bimbingan | Penghulu |
| Menunggu Bimbingan | Complete counseling | Sudah Bimbingan | Staff/Kepala KUA |
| Sudah Bimbingan | Complete marriage | Selesai | Staff/Kepala KUA |

---

## ⚠️ Error Codes

### HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid input or validation error |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error |

### Common Error Response Format

```json
{
  "error": "Error message",
  "details": "Detailed explanation (optional)"
}
```

### Validation Errors

**Date Validation:**
```json
{
  "error": "Tanggal nikah tidak boleh di masa lalu",
  "details": "Tanggal yang Anda masukkan: 2024-01-01"
}
```

**Time Format:**
```json
{
  "error": "Format waktu tidak valid",
  "details": "Gunakan format HH:MM (24-jam), contoh: 09:00"
}
```

**NIK Duplicate:**
```json
{
  "error": "NIK sudah terdaftar",
  "details": "NIK 6371012501950001 sudah digunakan untuk pendaftaran lain"
}
```

**Guardian Validation:**
```json
{
  "error": "Validasi wali nikah gagal",
  "details": "Jika ayah kandung masih hidup, wali nikah harus Ayah Kandung",
  "validation_rule": "Validasi 2: Konsistensi Status Wali dengan Ayah Kandung"
}
```

---

## 📊 Business Rules

### Marriage Capacity Rules

| Rule | Value | Description |
|------|-------|-------------|
| Max marriages per day at KUA | 9 | Maximum capacity for KUA venue |
| Max marriages outside KUA | Unlimited | No limit for external venues |
| Max marriages per penghulu/day | 3 | One penghulu can handle 3 marriages max |
| Min gap between marriages | 60 minutes | Minimum time between two consecutive marriages for same penghulu |
| Operating hours | 08:00 - 16:00 | 9 time slots available |

### Counseling Rules

| Rule | Value | Description |
|------|-------|-------------|
| Counseling day | Wednesday only | Bimbingan only scheduled on Wednesdays |
| Max participants per session | 10 couples | Maximum 10 couples per session |
| Sessions per Wednesday | 1 | Only one session per Wednesday |
| Mandatory attendance | Yes | Must attend to proceed to marriage |

### Date & Age Rules

| Rule | Value | Description |
|------|-------|-------------|
| Min working days before marriage | 10 days | Must register at least 10 working days before wedding date |
| Min age (groom) | 19 years | Minimum age for groom |
| Min age (bride) | 19 years | Minimum age for bride |
| Dispensation required if | Age < 19 or Days < 10 | Court dispensation needed |

### Guardian Rules (Wali Nikah)

1. ✅ **Wali must be alive** (`status_keberadaan` = `"Hidup"`)
2. ✅ **If father alive** → Wali MUST be `"Ayah Kandung"`
3. ✅ **If father deceased** → Wali moves to next nasab order
4. ✅ **NIK validation** - If wali is father, NIK must match bride's father NIK
5. ✅ **Wali cannot be bride or groom** - Different NIK required

**Nasab Order (Priority):**
1. Ayah Kandung (biological father)
2. Kakek (grandfather)
3. Saudara Laki-Laki Kandung (full brother)
4. Saudara Laki-Laki Seayah (half-brother)
5. Keponakan Laki-Laki (nephew)
6. Paman Kandung (full uncle)
7. Paman Seayah (half uncle)
8. Sepupu Laki-Laki (male cousin)
9. Wali Hakim (if no wali nasab available)

---

## 🧪 Testing Guide

### Using cURL

**1. Health Check:**
```bash
curl http://localhost:8080/health
```

**2. Register:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "nama": "Test User",
    "role": "user_biasa"
  }'
```

**3. Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }'
```

**4. Get Profile (with token):**
```bash
curl http://localhost:8080/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**5. Create Marriage Registration:**
```bash
curl -X POST http://localhost:8080/simnikah/pendaftaran/form-baru \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d @request_body.json
```

### Using JavaScript (Fetch)

**Login Example:**
```javascript
const login = async () => {
  const response = await fetch('http://localhost:8080/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      username: 'testuser',
      password: 'SecurePass123!'
    })
  });
  
  const data = await response.json();
  
  if (response.ok) {
    // Save token
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    console.log('Login successful:', data);
  } else {
    console.error('Login failed:', data.error);
  }
};
```

**Get Profile with Token:**
```javascript
const getProfile = async () => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('http://localhost:8080/profile', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const data = await response.json();
  
  if (response.ok) {
    console.log('Profile:', data.user);
  } else {
    console.error('Error:', data.error);
  }
};
```

**Create Marriage Registration:**
```javascript
const createRegistration = async (formData) => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('http://localhost:8080/simnikah/pendaftaran/form-baru', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(formData)
  });
  
  const data = await response.json();
  
  if (response.ok) {
    console.log('Registration created:', data);
  } else {
    console.error('Registration failed:', data.error);
  }
};
```

### Using Axios

**Setup Axios Interceptor:**
```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Request interceptor - add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

**Usage:**
```javascript
import api from './api';

// Login
const login = async (username, password) => {
  try {
    const response = await api.post('/login', { username, password });
    localStorage.setItem('token', response.data.token);
    return response.data;
  } catch (error) {
    console.error('Login error:', error.response?.data);
    throw error;
  }
};

// Get profile
const getProfile = async () => {
  try {
    const response = await api.get('/profile');
    return response.data;
  } catch (error) {
    console.error('Profile error:', error.response?.data);
    throw error;
  }
};

// Create registration
const createRegistration = async (formData) => {
  try {
    const response = await api.post('/simnikah/pendaftaran/form-baru', formData);
    return response.data;
  } catch (error) {
    console.error('Registration error:', error.response?.data);
    throw error;
  }
};
```

---

## 📝 Notes for Frontend Developers

### 1. Token Management

- **Save token** after successful login
- **Include token** in every authenticated request
- **Handle 401** responses (token expired/invalid) → redirect to login
- **Refresh token** before expiration (optional)

### 2. Date & Time Formats

- **Date Input**: Use `YYYY-MM-DD` format (ISO 8601)
- **Time Input**: Use `HH:MM` format (24-hour)
- **Date Display**: Format as needed for user (e.g., "25 Desember 2025")
- **API Returns**: Dates in ISO 8601 format with timezone (e.g., `2025-12-25T00:00:00Z`)

### 3. Status Display

Use colors and icons to indicate status:

**Status Colors:**
- 🟡 Yellow: `"Menunggu Verifikasi"`, `"Menunggu Pengumpulan Berkas"`
- 🟢 Green: `"Berkas Diterima"`, `"Sudah Bimbingan"`, `"Selesai"`
- 🔵 Blue: `"Menunggu Penugasan"`, `"Menunggu Bimbingan"`
- 🔴 Red: `"Ditolak"`

### 4. Form Validation (Client-Side)

Before sending to API, validate:
- ✅ NIK format (16 digits)
- ✅ Email format
- ✅ Phone number (min 10 digits)
- ✅ Date not in the past
- ✅ Time format HH:MM
- ✅ All required fields filled

### 5. Error Handling

```javascript
const handleApiError = (error) => {
  if (error.response) {
    // Server responded with error
    const { status, data } = error.response;
    
    switch (status) {
      case 400:
        return `Validasi gagal: ${data.error}`;
      case 401:
        return 'Sesi Anda telah berakhir, silakan login kembali';
      case 403:
        return 'Anda tidak memiliki akses untuk operasi ini';
      case 404:
        return 'Data tidak ditemukan';
      case 500:
        return 'Terjadi kesalahan server, silakan coba lagi';
      default:
        return data.error || 'Terjadi kesalahan';
    }
  } else if (error.request) {
    // No response received
    return 'Tidak dapat terhubung ke server';
  } else {
    // Other errors
    return error.message;
  }
};
```

### 6. Loading States

Show loading indicators for:
- Form submissions
- Data fetching
- File uploads (if implemented)

### 7. Real-time Updates

Consider implementing:
- **Polling** for notification count (every 30s)
- **WebSocket** for real-time notifications (future enhancement)

### 8. Mobile Responsiveness

API is mobile-ready, ensure your frontend is too:
- Responsive forms
- Touch-friendly buttons
- Mobile-optimized calendar views

---

## 🎯 Quick Integration Checklist

- [ ] Setup axios/fetch with base URL
- [ ] Implement token storage (localStorage)
- [ ] Create authentication flow (login/register)
- [ ] Handle 401 responses (redirect to login)
- [ ] Implement registration form (10 sections)
- [ ] Add client-side validation
- [ ] Display status flow with progress indicator
- [ ] Implement calendar views
- [ ] Add notification system
- [ ] Test all user flows
- [ ] Handle errors gracefully
- [ ] Add loading states

---

## 📞 Support

- **API Base URL**: `https://your-api.railway.app`
- **Documentation**: This file
- **Contact**: Backend team

---

**Last Updated:** 30 Oktober 2025  
**Version:** 2.0  
**Total Endpoints:** 50+

---

Made with ❤️ for SimNikah Frontend Team

