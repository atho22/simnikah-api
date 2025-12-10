# 📚 SimNikah API - Complete Documentation

**Version:** 1.3.0  
**Last Updated:** December 2024  
**Base URL:** `https://your-api-domain.com` atau `http://localhost:8080`

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Base Information](#base-information)
4. [Authentication Endpoints](#authentication-endpoints)
5. [Catin (User) Endpoints](#catin-user-endpoints)
6. [Calendar & Availability Endpoints](#calendar--availability-endpoints)
7. [Staff Endpoints](#staff-endpoints)
8. [Penghulu Endpoints](#penghulu-endpoints)
9. [Kepala KUA Endpoints](#kepala-kua-endpoints)
10. [Location Endpoints](#location-endpoints)
11. [Notification Endpoints](#notification-endpoints)
12. [Dashboard & Analytics Endpoints](#dashboard--analytics-endpoints)
13. [Error Handling](#error-handling)
14. [Status Flow](#status-flow)
15. [Rate Limiting](#rate-limiting)
16. [Data Models](#data-models)

---

## Overview

SimNikah API adalah RESTful API untuk sistem informasi manajemen pernikahan di KUA (Kantor Urusan Agama). API ini mendukung berbagai role: user biasa (catin), staff, penghulu, dan kepala KUA.

### Features

- ✅ JWT Authentication
- ✅ Role-based Access Control (RBAC)
- ✅ Marriage Registration Management
- ✅ Calendar & Availability System
- ✅ Notification System
- ✅ Geocoding & Maps Integration
- ✅ Dashboard & Analytics
- ✅ Feedback System
- ✅ Rate Limiting

---

## Authentication

### JWT Token

Semua endpoint (kecuali `/register`, `/login`, `/health`, dan endpoint kalender publik) memerlukan JWT token di header:

```
Authorization: Bearer <jwt_token>
```

### Token Details

- **Type:** JWT (JSON Web Token)
- **Algorithm:** HS256
- **Validity:** 24 jam
- **Key:** Diambil dari environment variable `JWT_KEY`

### Token Structure

```json
{
  "user_id": "USR1704067200",
  "email": "ahmad@example.com",
  "role": "user_biasa",
  "nama": "Ahmad Wijaya",
  "exp": 1704153600,
  "iat": 1704067200,
  "nbf": 1704067200
}
```

---

## Base Information

### Base URL

- **Development:** `http://localhost:8080`
- **Production:** `https://your-api-domain.com`

### Content Type

Semua request dan response menggunakan `application/json`.

### Response Format

Semua response mengikuti format standar:

```json
{
  "success": true|false,
  "message": "Deskripsi pesan",
  "data": { /* data object */ },
  "error": "Detail error (jika ada)"
}
```

---

## Authentication Endpoints

### 1. Register User

**Endpoint:** `POST /register`

**Description:** Mendaftarkan user baru ke sistem.

**Auth Required:** ❌ No

**Rate Limit:** 5 requests per minute per IP

**Request Body:**
```json
{
  "username": "ahmad123",
  "email": "ahmad@example.com",
  "password": "password123",
  "nama": "Ahmad Wijaya",
  "role": "user_biasa"
}
```

**Valid Roles:**
- `user_biasa` - User biasa untuk daftar nikah
- `staff` - Staff KUA untuk verifikasi
- `penghulu` - Penghulu untuk memimpin nikah
- `kepala_kua` - Kepala KUA untuk approval dan manajemen

**Response Success (201):**
```json
{
  "success": true,
  "message": "User berhasil dibuat",
  "data": {
    "user_id": "USR1704067200",
    "username": "ahmad123",
    "email": "ahmad@example.com",
    "nama": "Ahmad Wijaya",
    "role": "user_biasa"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Username sudah digunakan",
  "error": "Username sudah terdaftar"
}
```

---

### 2. Login User

**Endpoint:** `POST /login`

**Description:** Login user dan mendapatkan JWT token.

**Auth Required:** ❌ No

**Rate Limit:** 5 requests per minute per IP

**Request Body:**
```json
{
  "username": "ahmad123",
  "password": "password123"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR1704067200",
    "username": "ahmad123",
    "email": "ahmad@example.com",
    "role": "user_biasa",
    "nama": "Ahmad Wijaya"
  },
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "user_id": "USR1704067200",
      "username": "ahmad123",
      "email": "ahmad@example.com",
      "role": "user_biasa",
      "nama": "Ahmad Wijaya"
    }
  }
}
```

**Response Error (401):**
```json
{
  "success": false,
  "message": "Username atau password salah",
  "error": "Kredensial tidak valid"
}
```

---

### 3. Get Profile

**Endpoint:** `GET /profile`

**Description:** Mendapatkan informasi profile user yang sedang login.

**Auth Required:** ✅ Yes

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Profile berhasil diambil",
  "data": {
    "user_id": "USR1704067200",
    "username": "ahmad123",
    "email": "ahmad@example.com",
    "nama": "Ahmad Wijaya",
    "role": "user_biasa",
    "status": "Aktif",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

## Catin (User) Endpoints

### 4. Create Marriage Registration

**Endpoint:** `POST /simnikah/pendaftaran`

**Description:** Membuat pendaftaran nikah dengan form sederhana.

**Auth Required:** ✅ Yes

**Role Required:** `user_biasa`

**Request Body:**
```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad bin Abdullah",
    "pendidikan_akhir": "S1",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti binti Muhammad",
    "pendidikan_akhir": "S1",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di KUA",
    "tanggal_nikah": "2024-12-25",
    "waktu_nikah": "10:00",
    "alamat_nikah": "Jl. Contoh No. 123",
    "detail_alamat": "RT 01 RW 02",
    "kelurahan": "Sungai Miai"
  },
  "wali_nikah": {
    "nama_dan_bin": "Abdullah bin Muhammad",
    "hubungan_wali": "Ayah Kandung"
  }
}
```

**Note:** Field `alamat_nikah`, `detail_alamat`, dan `kelurahan` hanya diperlukan jika `tempat_nikah` = "Di Luar KUA".

**Valid Hubungan Wali:**
- Ayah Kandung
- Kakek
- Saudara Laki-Laki Kandung
- Saudara Laki-Laki Seayah
- Keponakan Laki-Laki
- Paman Kandung
- Paman Seayah
- Sepupu Laki-Laki
- Wali Hakim
- Lainnya

**Valid Kelurahan (Banjarmasin Utara):**
- Sungai Miai
- Sungai Andai
- Surgi Mufti
- Pangeran
- Kuin Utara
- Antasan Kecil Timur
- Alalak Utara
- Alalak Tengah
- Alalak Selatan

**Response Success (201):**
```json
{
  "success": true,
  "message": "Pendaftaran berhasil dibuat",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20240101001",
    "status_pendaftaran": "Draft",
    "tanggal_nikah": "2024-12-25T10:00:00Z",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di KUA",
    "calon_suami": {
      "nama_lengkap": "Ahmad bin Abdullah",
      "umur": 25,
      "pendidikan_terakhir": "S1"
    },
    "calon_istri": {
      "nama_lengkap": "Siti binti Muhammad",
      "umur": 23,
      "pendidikan_terakhir": "S1"
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    }
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Jadwal sudah penuh",
  "error": "Jadwal pernikahan pada tanggal 2024-12-25 pukul 10:00 sudah penuh. Silakan pilih tanggal atau jam lain.",
  "field": "waktu_nikah",
  "type": "schedule_conflict"
}
```

---

### 5. Get Registration Status

**Endpoint:** `GET /simnikah/pendaftaran/status`

**Description:** Cek apakah user sudah memiliki pendaftaran nikah.

**Auth Required:** ✅ Yes

**Role Required:** `user_biasa`

**Response Success (200) - Ada Pendaftaran:**
```json
{
  "success": true,
  "has_registration": true,
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20240101001",
    "status_pendaftaran": "Disetujui",
    "tanggal_nikah": "2024-12-25T10:00:00Z",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di KUA",
    "calon_suami": {
      "nama_lengkap": "Ahmad bin Abdullah"
    },
    "calon_istri": {
      "nama_lengkap": "Siti binti Muhammad"
    }
  }
}
```

**Response Success (200) - Belum Ada Pendaftaran:**
```json
{
  "success": true,
  "has_registration": false,
  "message": "Belum ada pendaftaran"
}
```

---

### 6. List Registrations (Staff/Kepala KUA)

**Endpoint:** `GET /simnikah/pendaftaran`

**Description:** Mendapatkan daftar semua pendaftaran (hanya untuk staff dan kepala KUA).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Query Parameters:**
- `status` (optional) - Filter by status
- `page` (optional) - Page number (default: 1)
- `limit` (optional) - Items per page (default: 10)

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG20240101001",
        "status_pendaftaran": "Disetujui",
        "tanggal_nikah": "2024-12-25T10:00:00Z",
        "calon_suami": {
          "nama_lengkap": "Ahmad bin Abdullah"
        },
        "calon_istri": {
          "nama_lengkap": "Siti binti Muhammad"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

---

### 7. Create Feedback Pernikahan

**Endpoint:** `POST /simnikah/feedback-pernikahan`

**Description:** Membuat feedback, saran, kritik, atau rating setelah pernikahan selesai.

**Auth Required:** ✅ Yes

**Role Required:** `user_biasa`

**Request Body:**
```json
{
  "pendaftaran_id": 1,
  "jenis_feedback": "Rating",
  "rating": 5,
  "judul": "Pelayanan Sangat Baik",
  "pesan": "Terima kasih atas pelayanan yang sangat baik dari KUA. Proses nikah berjalan lancar."
}
```

**Valid Jenis Feedback:**
- `Rating` - Rating 1-5 (wajib field `rating`)
- `Saran` - Saran untuk perbaikan
- `Kritik` - Kritik konstruktif
- `Laporan` - Laporan masalah

**Response Success (201):**
```json
{
  "success": true,
  "message": "Feedback berhasil dibuat",
  "data": {
    "id": 1,
    "pendaftaran_id": 1,
    "jenis_feedback": "Rating",
    "rating": 5,
    "judul": "Pelayanan Sangat Baik",
    "status_baca": "Belum Dibaca"
  }
}
```

---

## Calendar & Availability Endpoints

### 8. Get Calendar Availability

**Endpoint:** `GET /simnikah/kalender-ketersediaan`

**Description:** Mendapatkan kalender ketersediaan slot pernikahan per bulan.

**Auth Required:** ❌ No (Public endpoint)

**Query Parameters:**
- `bulan` (optional) - Bulan (1-12), default: bulan saat ini
- `tahun` (optional) - Tahun (YYYY), default: tahun saat ini

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "2024-12-01": {
      "available": true,
      "booked": 2,
      "max": 9,
      "slots_available": 7
    },
    "2024-12-02": {
      "available": true,
      "booked": 5,
      "max": 9,
      "slots_available": 4
    },
    "2024-12-03": {
      "available": false,
      "booked": 9,
      "max": 9,
      "slots_available": 0
    }
  },
  "month": 12,
  "year": 2024
}
```

---

### 9. Get Available Time Slots

**Endpoint:** `GET /simnikah/ketersediaan-jam`

**Description:** Mendapatkan slot waktu yang tersedia untuk tanggal tertentu.

**Auth Required:** ❌ No (Public endpoint)

**Query Parameters:**
- `tanggal` (required) - Format: YYYY-MM-DD
- `tempat` (required) - "Di KUA" atau "Di Luar KUA"

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-12-25",
    "tempat": "Di KUA",
    "available_slots": ["08:00", "09:00", "11:00", "14:00"],
    "unavailable_slots": ["10:00", "12:00", "13:00", "15:00", "16:00"],
    "total_available": 4,
    "total_unavailable": 5
  }
}
```

---

### 10. Get Weddings by Date

**Endpoint:** `GET /simnikah/pernikahan-tanggal`

**Description:** Mendapatkan daftar pernikahan pada tanggal tertentu.

**Auth Required:** ❌ No (Public endpoint)

**Query Parameters:**
- `tanggal` (required) - Format: YYYY-MM-DD

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-12-25",
    "weddings": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG20240101001",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di KUA",
        "calon_suami": {
          "nama_lengkap": "Ahmad bin Abdullah"
        },
        "calon_istri": {
          "nama_lengkap": "Siti binti Muhammad"
        }
      }
    ],
    "total": 1
  }
}
```

---

## Staff Endpoints

### 11. List Staff

**Endpoint:** `GET /simnikah/staff`

**Description:** Mendapatkan daftar semua staff KUA.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067201",
      "nip": "198001012001011001",
      "nama_lengkap": "Budi Santoso",
      "jabatan": "Staff",
      "bagian": "Pendaftaran",
      "status": "Aktif"
    }
  ]
}
```

---

### 12. Update Staff

**Endpoint:** `PUT /simnikah/staff/:id`

**Description:** Update data staff.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "nama_lengkap": "Budi Santoso",
  "jabatan": "Staff",
  "bagian": "Verifikasi",
  "status": "Aktif"
}
```

---

### 13. Verify Registration Form

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`

**Description:** Verifikasi formulir pendaftaran (online form).

**Auth Required:** ✅ Yes

**Role Required:** `staff`

**Request Body:**
```json
{
  "catatan": "Formulir sudah lengkap dan valid"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Formulir berhasil diverifikasi",
  "data": {
    "id": 1,
    "status": "Menunggu Pengumpulan Berkas"
  }
}
```

---

### 14. Verify Documents

**Endpoint:** `POST /simnikah/staff/verify-berkas/:id`

**Description:** Verifikasi berkas fisik pendaftaran.

**Auth Required:** ✅ Yes

**Role Required:** `staff`

**Request Body:**
```json
{
  "catatan": "Berkas lengkap dan valid"
}
```

---

### 15. Approve Registration

**Endpoint:** `POST /simnikah/staff/approve/:id`

**Description:** Menyetujui pendaftaran nikah.

**Auth Required:** ✅ Yes

**Role Required:** `staff`

**Request Body:**
```json
{
  "catatan": "Pendaftaran disetujui, semua dokumen lengkap"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Pendaftaran berhasil disetujui",
  "data": {
    "id": 1,
    "status_pendaftaran": "Disetujui",
    "disetujui_oleh": "USR1704067201",
    "disetujui_pada": "2024-12-01T10:00:00Z"
  }
}
```

---

### 16. Create Registration for User

**Endpoint:** `POST /simnikah/staff/pendaftaran`

**Description:** Staff membuat pendaftaran nikah atas nama user (untuk user yang tidak bisa menggunakan aplikasi).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Request Body:** (Sama seperti Create Registration, ditambah `pendaftar_username` atau `pendaftar_email`)

---

### 17. Update Registration Status

**Endpoint:** `PUT /simnikah/pendaftaran/:id/update-status`

**Description:** Update status pendaftaran.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `penghulu`, `kepala_kua`

**Request Body:**
```json
{
  "status": "Disetujui",
  "catatan": "Status diupdate oleh staff"
}
```

**Valid Status:**
- `Draft`
- `Disetujui`
- `Selesai`
- `Ditolak`

**Note:** Status terkait penghulu (`Menunggu Penugasan`, `Penghulu Ditugaskan`) hanya bisa diubah melalui endpoint assign-penghulu.

---

### 18. Get Approved Registrations for Announcement

**Endpoint:** `POST /simnikah/staff/pengumuman-nikah/list`

**Description:** Mendapatkan data pendaftaran yang disetujui untuk surat pengumuman nikah.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "tanggal_awal": "2024-12-01",
  "tanggal_akhir": "2024-12-07",
  "kop_surat": {
    "nama_kua": "KUA Kecamatan Banjarmasin Utara",
    "alamat": "Jl. Contoh No. 123",
    "kota": "Banjarmasin",
    "provinsi": "Kalimantan Selatan"
  }
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG20240101001",
        "calon_suami": {
          "nama_lengkap": "Ahmad bin Abdullah",
          "pendidikan_terakhir": "S1"
        },
        "calon_istri": {
          "nama_lengkap": "Siti binti Muhammad",
          "pendidikan_terakhir": "S1"
        },
        "tanggal_nikah": "2024-12-25T10:00:00Z",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di KUA",
        "wali_nikah": {
          "nama_dan_bin": "Abdullah bin Muhammad",
          "hubungan_wali": "Ayah Kandung"
        }
      }
    ],
    "kop_surat": {
      "nama_kua": "KUA Kecamatan Banjarmasin Utara",
      "alamat": "Jl. Contoh No. 123"
    },
    "periode": {
      "tanggal_awal": "2024-12-01",
      "tanggal_akhir": "2024-12-07"
    }
  }
}
```

---

## Penghulu Endpoints

### 19. List Marriage Officers

**Endpoint:** `GET /simnikah/penghulu`

**Description:** Mendapatkan daftar semua penghulu.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067202",
      "nip": "198001012001011002",
      "nama_lengkap": "Ustadz Ahmad",
      "status": "Aktif",
      "jumlah_nikah": 150,
      "rating": 4.8
    }
  ]
}
```

---

### 20. Update Marriage Officer

**Endpoint:** `PUT /simnikah/penghulu/:id`

**Description:** Update data penghulu.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

---

### 21. Verify Registration Documents (Penghulu)

**Endpoint:** `POST /simnikah/penghulu/verify-documents/:id`

**Description:** Penghulu memverifikasi dokumen pendaftaran.

**Auth Required:** ✅ Yes

**Role Required:** `penghulu`

---

### 22. List My Assignments

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`

**Description:** Mendapatkan daftar pendaftaran yang ditugaskan ke penghulu.

**Auth Required:** ✅ Yes

**Role Required:** `penghulu`, `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nomor_pendaftaran": "REG20240101001",
      "calon_suami": {
        "nama_lengkap": "Ahmad bin Abdullah"
      },
      "calon_istri": {
        "nama_lengkap": "Siti binti Muhammad"
      },
      "tanggal_nikah": "2024-12-25T10:00:00Z",
      "waktu_nikah": "10:00",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "KUA Kecamatan Banjarmasin Utara",
      "status_pendaftaran": "Penghulu Ditugaskan"
    }
  ]
}
```

---

### 23. Get Today Schedule

**Endpoint:** `GET /simnikah/penghulu/today-schedule`

**Description:** Mendapatkan jadwal pernikahan hari ini untuk penghulu.

**Auth Required:** ✅ Yes

**Role Required:** `penghulu`, `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-12-25",
    "schedule": [
      {
        "id": 1,
        "waktu_nikah": "10:00",
        "calon_suami": {
          "nama_lengkap": "Ahmad bin Abdullah"
        },
        "calon_istri": {
          "nama_lengkap": "Siti binti Muhammad"
        },
        "tempat_nikah": "Di KUA"
      }
    ],
    "total": 1
  }
}
```

---

### 24. Complete Marriage

**Endpoint:** `POST /simnikah/penghulu/complete-marriage/:id`

**Description:** Menandai pernikahan sudah selesai dilaksanakan.

**Auth Required:** ✅ Yes

**Role Required:** `penghulu`, `kepala_kua`

**Request Body:**
```json
{
  "catatan": "Pernikahan telah dilaksanakan dengan lancar"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Pernikahan berhasil ditandai selesai",
  "data": {
    "id": 1,
    "status_pendaftaran": "Selesai",
    "updated_at": "2024-12-25T10:30:00Z"
  }
}
```

---

## Kepala KUA Endpoints

### 25. Create Staff

**Endpoint:** `POST /simnikah/kepala-kua/staff`

**Description:** Membuat staff baru.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "username": "staff001",
  "email": "staff@kua.go.id",
  "password": "password123",
  "nama": "Budi Santoso",
  "nip": "198001012001011001",
  "jabatan": "Staff",
  "bagian": "Pendaftaran",
  "no_hp": "081234567890",
  "alamat": "Jl. Contoh No. 123"
}
```

---

### 26. Create Marriage Officer

**Endpoint:** `POST /simnikah/kepala-kua/penghulu`

**Description:** Membuat penghulu baru.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "username": "penghulu001",
  "email": "penghulu@kua.go.id",
  "password": "password123",
  "nama": "Ustadz Ahmad",
  "nip": "198001012001011002",
  "no_hp": "081234567891",
  "alamat": "Jl. Contoh No. 124"
}
```

---

### 27. Assign Penghulu

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`

**Description:** Menugaskan penghulu ke pendaftaran nikah.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "penghulu_id": 1
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Penghulu berhasil ditugaskan",
  "data": {
    "id": 1,
    "penghulu_id": 1,
    "penghulu": {
      "nama_lengkap": "Ustadz Ahmad"
    },
    "status_pendaftaran": "Penghulu Ditugaskan",
    "penghulu_assigned_at": "2024-12-01T10:00:00Z"
  }
}
```

---

### 28. List Available Officers

**Endpoint:** `GET /simnikah/kepala-kua/available-penghulu`

**Description:** Mendapatkan daftar penghulu yang tersedia.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Query Parameters:**
- `tanggal` (optional) - Filter penghulu yang tersedia pada tanggal tertentu
- `waktu` (optional) - Filter penghulu yang tersedia pada waktu tertentu

---

### 29. Get Penghulu Statistics

**Endpoint:** `GET /simnikah/kepala-kua/statistik-penghulu`

**Description:** Mendapatkan statistik performa penghulu.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "penghulu": [
      {
        "id": 1,
        "nama_lengkap": "Ustadz Ahmad",
        "jumlah_nikah": 150,
        "rating": 4.8,
        "nikah_bulan_ini": 12,
        "nikah_tahun_ini": 150
      }
    ]
  }
}
```

---

### 30. List Feedback Pernikahan

**Endpoint:** `GET /simnikah/kepala-kua/feedback`

**Description:** Mendapatkan daftar feedback pernikahan.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Query Parameters:**
- `status_baca` (optional) - Filter by status baca
- `jenis_feedback` (optional) - Filter by jenis feedback

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pendaftaran_id": 1,
      "jenis_feedback": "Rating",
      "rating": 5,
      "judul": "Pelayanan Sangat Baik",
      "pesan": "Terima kasih...",
      "status_baca": "Belum Dibaca",
      "created_at": "2024-12-01T10:00:00Z"
    }
  ]
}
```

---

### 31. Mark Feedback as Read

**Endpoint:** `PUT /simnikah/kepala-kua/feedback/:id/mark-read`

**Description:** Menandai feedback sudah dibaca.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Feedback berhasil ditandai sudah dibaca",
  "data": {
    "id": 1,
    "status_baca": "Sudah Dibaca",
    "dibaca_oleh": "USR1704067203",
    "dibaca_pada": "2024-12-01T10:00:00Z"
  }
}
```

---

### 32. Get Feedback Statistics

**Endpoint:** `GET /simnikah/kepala-kua/feedback/stats`

**Description:** Mendapatkan statistik feedback.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "total_feedback": 100,
    "belum_dibaca": 20,
    "sudah_dibaca": 80,
    "rating_average": 4.5,
    "by_jenis": {
      "Rating": 50,
      "Saran": 30,
      "Kritik": 15,
      "Laporan": 5
    }
  }
}
```

---

### 33. Get Approved Registrations for Announcement (Kepala KUA)

**Endpoint:** `POST /simnikah/kepala-kua/pengumuman-nikah/list`

**Description:** Sama seperti endpoint staff, untuk mendapatkan data pengumuman nikah.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:** (Sama seperti endpoint staff)

---

## Location Endpoints

### 34. Geocode Address

**Endpoint:** `POST /simnikah/location/geocode`

**Description:** Convert alamat ke koordinat (latitude, longitude).

**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "alamat": "Jl. Contoh No. 123, Banjarmasin"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Koordinat berhasil ditemukan",
  "data": {
    "alamat": "Jl. Contoh No. 123, Banjarmasin",
    "latitude": -3.3145,
    "longitude": 114.5921,
    "map_url": "https://www.google.com/maps?q=-3.3145,114.5921",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3145&mlon=114.5921&zoom=16"
  }
}
```

---

### 35. Reverse Geocode

**Endpoint:** `POST /simnikah/location/reverse-geocode`

**Description:** Convert koordinat ke alamat.

**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "latitude": -3.3145,
  "longitude": 114.5921
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "alamat": "Jl. Contoh No. 123, Banjarmasin",
    "latitude": -3.3145,
    "longitude": 114.5921
  }
}
```

---

### 36. Search Address

**Endpoint:** `GET /simnikah/location/search`

**Description:** Mencari alamat menggunakan autocomplete.

**Auth Required:** ✅ Yes

**Query Parameters:**
- `q` (required) - Query search

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "results": [
      {
        "alamat": "Jl. Contoh No. 123, Banjarmasin",
        "latitude": -3.3145,
        "longitude": 114.5921
      }
    ]
  }
}
```

---

### 37. Update Wedding Location with Coordinates

**Endpoint:** `PUT /simnikah/pendaftaran/:id/location`

**Description:** Update lokasi nikah dengan koordinat.

**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "alamat_akad": "Jl. Contoh No. 123",
  "latitude": -3.3145,
  "longitude": 114.5921
}
```

---

### 38. Get Wedding Location Detail

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`

**Description:** Mendapatkan detail lokasi nikah.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Contoh No. 123",
    "latitude": -3.3145,
    "longitude": 114.5921,
    "map_url": "https://www.google.com/maps?q=-3.3145,114.5921"
  }
}
```

---

### 39. Update Marriage Location

**Endpoint:** `PUT /simnikah/pendaftaran/:id/alamat`

**Description:** Update alamat nikah (hanya untuk nikah di luar KUA).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "alamat_akad": "Jl. Contoh Baru No. 456"
}
```

---

## Notification Endpoints

### 40. Create Notification

**Endpoint:** `POST /simnikah/notifikasi`

**Description:** Membuat notifikasi baru.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "user_id": "USR1704067200",
  "judul": "Pendaftaran Disetujui",
  "pesan": "Pendaftaran Anda telah disetujui",
  "tipe": "Success",
  "tautan": "/simnikah/pendaftaran/1"
}
```

**Valid Tipe:**
- `Info`
- `Success`
- `Warning`
- `Error`

---

### 41. Get User Notifications

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`

**Description:** Mendapatkan notifikasi user.

**Auth Required:** ✅ Yes

**Query Parameters:**
- `status_baca` (optional) - Filter by status baca
- `tipe` (optional) - Filter by tipe
- `limit` (optional) - Limit results

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "judul": "Pendaftaran Disetujui",
      "pesan": "Pendaftaran Anda telah disetujui",
      "tipe": "Success",
      "status_baca": "Belum Dibaca",
      "tautan": "/simnikah/pendaftaran/1",
      "created_at": "2024-12-01T10:00:00Z"
    }
  ]
}
```

---

### 42. Get Notification by ID

**Endpoint:** `GET /simnikah/notifikasi/:id`

**Description:** Mendapatkan detail notifikasi.

**Auth Required:** ✅ Yes

---

### 43. Update Notification Status

**Endpoint:** `PUT /simnikah/notifikasi/:id/status`

**Description:** Update status notifikasi (baca/tidak baca).

**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

---

### 44. Mark All as Read

**Endpoint:** `PUT /simnikah/notifikasi/user/:user_id/mark-all-read`

**Description:** Menandai semua notifikasi user sudah dibaca.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "message": "Semua notifikasi berhasil ditandai sudah dibaca",
  "data": {
    "updated_count": 5
  }
}
```

---

### 45. Delete Notification

**Endpoint:** `DELETE /simnikah/notifikasi/:id`

**Description:** Menghapus notifikasi.

**Auth Required:** ✅ Yes

---

### 46. Get Notification Statistics

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id/stats`

**Description:** Mendapatkan statistik notifikasi user.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "total": 20,
    "belum_dibaca": 5,
    "sudah_dibaca": 15,
    "by_tipe": {
      "Info": 10,
      "Success": 5,
      "Warning": 3,
      "Error": 2
    }
  }
}
```

---

### 47. Send Notification to Role

**Endpoint:** `POST /simnikah/notifikasi/send-to-role`

**Description:** Mengirim notifikasi ke semua user dengan role tertentu.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "role": "user_biasa",
  "judul": "Pengumuman Penting",
  "pesan": "Ada pengumuman penting untuk semua user",
  "tipe": "Info",
  "tautan": "/pengumuman"
}
```

---

### 48. Run Reminder Notification

**Endpoint:** `POST /simnikah/notifikasi/run-reminder`

**Description:** Menjalankan reminder notifikasi secara manual (untuk testing).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

---

## Dashboard & Analytics Endpoints

### 49. Get Kepala KUA Dashboard

**Endpoint:** `GET /simnikah/dashboard/kepala-kua`

**Description:** Mendapatkan data dashboard untuk kepala KUA.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Query Parameters:**
- `period` (optional) - `day`, `week`, `month`, `year` (default: `month`)
- `date_from` (optional) - Format: YYYY-MM-DD
- `date_to` (optional) - Format: YYYY-MM-DD

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "period": {
      "type": "month",
      "date_from": "2024-12-01",
      "date_to": "2024-12-31"
    },
    "statistics": {
      "total_pendaftaran": 150,
      "pending_approval": 10,
      "pernikahan_hari_ini": 5,
      "pernikahan_bulan_ini": 50
    },
    "trends": {
      "daily": [
        {
          "date": "2024-12-01",
          "count": 3
        }
      ]
    },
    "status_distribution": {
      "Draft": 5,
      "Disetujui": 20,
      "Penghulu Ditugaskan": 15,
      "Selesai": 100
    },
    "penghulu_performance": [
      {
        "nama_lengkap": "Ustadz Ahmad",
        "jumlah_nikah": 50,
        "rating": 4.8
      }
    ],
    "peak_hours": {
      "10:00": 20,
      "14:00": 15
    }
  }
}
```

---

### 50. Get Staff Dashboard

**Endpoint:** `GET /simnikah/dashboard/staff`

**Description:** Mendapatkan data dashboard untuk staff.

**Auth Required:** ✅ Yes

**Role Required:** `staff`

---

### 51. Get Marriage Statistics

**Endpoint:** `GET /simnikah/dashboard/statistik-pernikahan`

**Description:** Mendapatkan statistik pernikahan.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Query Parameters:**
- `period` (optional) - `day`, `week`, `month`, `year`
- `date_from` (optional)
- `date_to` (optional)

---

### 52. Get Penghulu Performance

**Endpoint:** `GET /simnikah/dashboard/penghulu-performance`

**Description:** Mendapatkan statistik performa penghulu.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

---

### 53. Get Peak Hours Analysis

**Endpoint:** `GET /simnikah/dashboard/peak-hours`

**Description:** Mendapatkan analisis jam-jam sibuk.

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "peak_hours": [
      {
        "waktu": "10:00",
        "jumlah": 20,
        "percentage": 25
      },
      {
        "waktu": "14:00",
        "jumlah": 15,
        "percentage": 18.75
      }
    ],
    "quiet_hours": [
      {
        "waktu": "08:00",
        "jumlah": 2,
        "percentage": 2.5
      }
    ]
  }
}
```

---

## Error Handling

### Error Response Format

Semua error mengikuti format standar:

```json
{
  "success": false,
  "message": "Deskripsi error",
  "error": "Detail error",
  "type": "error_type",
  "field": "field_name" // Optional, untuk validation errors
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | OK - Request berhasil |
| `201` | Created - Resource berhasil dibuat |
| `400` | Bad Request - Format data tidak valid |
| `401` | Unauthorized - Token tidak valid atau tidak ada |
| `403` | Forbidden - Role tidak memiliki akses |
| `404` | Not Found - Resource tidak ditemukan |
| `429` | Too Many Requests - Rate limit exceeded |
| `500` | Internal Server Error - Server error |

### Error Types

- `validation` - Validasi gagal
- `authentication` - Error autentikasi
- `authorization` - Error autorisasi
- `not_found` - Resource tidak ditemukan
- `database` - Database error
- `schedule_conflict` - Konflik jadwal
- `security` - Security error

### Common Error Examples

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format data tidak valid",
  "type": "validation",
  "field": "tanggal_nikah"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token otorisasi tidak disediakan",
  "type": "authentication"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Role kepala_kua diperlukan",
  "type": "authorization"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Data tidak ditemukan",
  "error": "Pendaftaran dengan ID tersebut tidak ditemukan",
  "type": "not_found"
}
```

**429 Too Many Requests:**
```json
{
  "success": false,
  "message": "Rate limit exceeded",
  "error": "Terlalu banyak request. Silakan coba lagi nanti.",
  "retry_after": "60 detik"
}
```

---

## Status Flow

### Marriage Registration Status Flow

```
Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
   ↓
Ditolak
```

### Status Descriptions

1. **Draft**
   - Status awal setelah pendaftaran dibuat
   - Catin sudah mengisi form pendaftaran

2. **Disetujui**
   - Staff sudah menyetujui pendaftaran
   - Setelah verifikasi dokumen

3. **Menunggu Penugasan**
   - Menunggu kepala KUA menentukan penghulu
   - Status intermediate

4. **Penghulu Ditugaskan**
   - Kepala KUA sudah assign penghulu
   - Penghulu bisa melihat detail pendaftaran

5. **Selesai**
   - Penghulu sudah melaksanakan nikah
   - Status final

6. **Ditolak**
   - Pendaftaran ditolak oleh staff
   - Bisa terjadi dari status Draft

---

## Rate Limiting

### Global Rate Limit

- **Limit:** 100 requests per minute per IP
- **Headers:**
  - `X-RateLimit-Limit`: Limit maksimal
  - `X-RateLimit-Remaining`: Sisa request
  - `X-RateLimit-Reset`: Waktu reset (Unix timestamp)

### Strict Rate Limit

Untuk endpoint autentikasi (`/register`, `/login`):
- **Limit:** 5 requests per minute per IP
- **Reason:** Security untuk mencegah brute force

### Rate Limit Response

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

## Data Models

### User Model

```json
{
  "id": 1,
  "user_id": "USR1704067200",
  "username": "ahmad123",
  "email": "ahmad@example.com",
  "nama": "Ahmad Wijaya",
  "role": "user_biasa",
  "status": "Aktif",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### PendaftaranNikah Model

```json
{
  "id": 1,
  "nomor_pendaftaran": "REG20240101001",
  "pendaftar_id": "USR1704067200",
  "calon_suami_id": "USR1704067201",
  "calon_istri_id": "USR1704067202",
  "wali_nikah_id": 1,
  "tanggal_pendaftaran": "2024-01-01T10:00:00Z",
  "tanggal_nikah": "2024-12-25T10:00:00Z",
  "waktu_nikah": "10:00",
  "tempat_nikah": "Di KUA",
  "alamat_akad": "KUA Kecamatan Banjarmasin Utara",
  "latitude": null,
  "longitude": null,
  "status_pendaftaran": "Disetujui",
  "penghulu_id": 1,
  "penghulu_assigned_by": "USR1704067203",
  "penghulu_assigned_at": "2024-12-01T10:00:00Z",
  "catatan": "Pendaftaran disetujui",
  "disetujui_oleh": "USR1704067201",
  "disetujui_pada": "2024-12-01T10:00:00Z",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-12-01T10:00:00Z"
}
```

### CalonPasangan Model

```json
{
  "id": 1,
  "user_id": "USR1704067201",
  "nik": "T123456789012345",
  "nama_lengkap": "Ahmad bin Abdullah",
  "tanggal_lahir": "1999-01-01T00:00:00Z",
  "jenis_kelamin": "L",
  "pendidikan_terakhir": "S1",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### WaliNikah Model

```json
{
  "id": 1,
  "pendaftaran_id": 1,
  "nama_dan_bin": "Abdullah bin Muhammad",
  "hubungan_wali": "Ayah Kandung",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### Notifikasi Model

```json
{
  "id": 1,
  "user_id": "USR1704067200",
  "judul": "Pendaftaran Disetujui",
  "pesan": "Pendaftaran Anda telah disetujui",
  "tipe": "Success",
  "status_baca": "Belum Dibaca",
  "tautan": "/simnikah/pendaftaran/1",
  "created_at": "2024-12-01T10:00:00Z",
  "updated_at": "2024-12-01T10:00:00Z"
}
```

---

## Additional Information

### Health Check

**Endpoint:** `GET /health`

**Description:** Health check endpoint untuk monitoring.

**Auth Required:** ❌ No

**Response:**
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2024-12-01T10:00:00Z"
}
```

### Kapasitas Pernikahan

**Aturan Kapasitas:**
- **Di KUA:** 1 pernikahan per jam (1 ruang nikah)
- **Di Luar KUA:** 3 pernikahan per jam (3 penghulu)
- **Total Maksimal:** 3 pernikahan per jam
- **Jam Operasional:** 08:00 - 16:00 WITA
- **Total Slot per Hari:** 9 slot (08:00, 09:00, 10:00, 11:00, 12:00, 13:00, 14:00, 15:00, 16:00)

### Kelurahan Valid

Hanya kelurahan di Kecamatan Banjarmasin Utara:
- Sungai Miai
- Sungai Andai
- Surgi Mufti
- Pangeran
- Kuin Utara
- Antasan Kecil Timur
- Alalak Utara
- Alalak Tengah
- Alalak Selatan

---

## Support & Contact

Untuk pertanyaan atau bantuan, silakan hubungi tim development.

**Last Updated:** December 2024  
**API Version:** 1.3.0







