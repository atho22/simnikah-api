# 📚 SimNikah API - Dokumentasi Lengkap untuk Frontend

**Version:** 1.0.0  
**Last Updated:** 2024-01-15  
**Base URL:** `https://your-api.railway.app` atau `http://localhost:8080`

---

## 📑 Daftar Isi

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [User Roles](#user-roles)
4. [API Endpoints](#api-endpoints)
   - [Authentication](#1-authentication)
   - [Pendaftaran Nikah](#2-pendaftaran-nikah)
   - [Staff Management](#3-staff-management)
   - [Penghulu Management](#4-penghulu-management)
   - [Status Management](#5-status-management)
   - [Jadwal & Kalender](#6-jadwal--kalender)
   - [Bimbingan Perkawinan](#7-bimbingan-perkawinan)
   - [Notifikasi](#8-notifikasi)
   - [Location & Maps](#9-location--maps)
5. [Error Handling](#error-handling)
6. [Status Flow](#status-flow)
7. [Business Rules](#business-rules)

---

## 🚀 Getting Started

### Base URL
```
Development: http://localhost:8080
Production:  https://your-api.railway.app
```

### Headers
```http
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Authentication
Semua endpoint kecuali `/register`, `/login`, dan `/health` memerlukan JWT token di header.

---

## 🔐 Authentication

### 1. Register User
**POST** `/register`

**Auth:** ❌ Tidak diperlukan

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "nama": "John Doe",
  "role": "user_biasa"
}
```

**Response (200 OK):**
```json
{
  "message": "User berhasil didaftarkan",
  "user_id": "USR1704067200",
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user_biasa"
}
```

---

### 2. Login
**POST** `/login`

**Auth:** ❌ Tidak diperlukan

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR1704067200",
    "username": "johndoe",
    "email": "john@example.com",
    "nama": "John Doe",
    "role": "user_biasa"
  }
}
```

---

### 3. Get Profile
**GET** `/profile`

**Auth:** ✅ Required

**Response (200 OK):**
```json
{
  "user_id": "USR1704067200",
  "username": "johndoe",
  "email": "john@example.com",
  "nama": "John Doe",
  "role": "user_biasa",
  "status": "Aktif"
}
```

---

## 👥 User Roles

| Role | Code | Description |
|------|------|-------------|
| User Biasa | `user_biasa` | Calon pengantin |
| Staff | `staff` | Staff KUA |
| Penghulu | `penghulu` | Penghulu |
| Kepala KUA | `kepala_kua` | Kepala KUA |

---

## 📋 API Endpoints

### 1. Authentication

#### 1.1 Register
- **Method:** `POST`
- **Endpoint:** `/register`
- **Auth:** ❌
- **Rate Limit:** Strict (5 req/min)

#### 1.2 Login
- **Method:** `POST`
- **Endpoint:** `/login`
- **Auth:** ❌
- **Rate Limit:** Strict (5 req/min)

#### 1.3 Get Profile
- **Method:** `GET`
- **Endpoint:** `/profile`
- **Auth:** ✅

---

### 2. Pendaftaran Nikah

#### 2.1 Create Marriage Registration
**POST** `/simnikah/pendaftaran/form-baru`

**Auth:** ✅  
**Role:** `user_biasa`

**Request Body:**
```json
{
  "calon_suami": {
    "nik": "3201010101010001",
    "nama_lengkap": "Ahmad Fauzi",
    "tempat_lahir": "Jakarta",
    "tanggal_lahir": "1995-05-15",
    "kewarganegaraan": "WNI",
    "agama": "Islam",
    "status": "Belum Kawin",
    "pekerjaan": "Karyawan Swasta",
    "alamat": "Jl. Merdeka No. 123",
    "no_hp": "081234567890"
  },
  "calon_istri": {
    "nik": "3201010101010002",
    "nama_lengkap": "Siti Aisyah",
    "tempat_lahir": "Bandung",
    "tanggal_lahir": "1997-08-20",
    "kewarganegaraan": "WNI",
    "agama": "Islam",
    "status": "Belum Kawin",
    "pekerjaan": "Guru",
    "alamat": "Jl. Sudirman No. 456",
    "no_hp": "081234567891"
  },
  "data_orang_tua_suami": {
    "ayah": {
      "nik": "3201010101010003",
      "nama_lengkap": "Budi Santoso",
      "tempat_lahir": "Jakarta",
      "tanggal_lahir": "1970-01-01",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "pekerjaan": "PNS",
      "alamat": "Jl. Merdeka No. 123",
      "status_keberadaan": "Hidup"
    },
    "ibu": {
      "nik": "3201010101010004",
      "nama_lengkap": "Siti Fatimah",
      "tempat_lahir": "Jakarta",
      "tanggal_lahir": "1972-03-15",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "pekerjaan": "Ibu Rumah Tangga",
      "alamat": "Jl. Merdeka No. 123",
      "status_keberadaan": "Hidup"
    }
  },
  "data_orang_tua_istri": {
    "ayah": {
      "nik": "3201010101010005",
      "nama_lengkap": "Ahmad Hidayat",
      "tempat_lahir": "Bandung",
      "tanggal_lahir": "1968-06-10",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "pekerjaan": "Wiraswasta",
      "alamat": "Jl. Sudirman No. 456",
      "status_keberadaan": "Hidup"
    },
    "ibu": {
      "nik": "3201010101010006",
      "nama_lengkap": "Rohani",
      "tempat_lahir": "Bandung",
      "tanggal_lahir": "1970-09-25",
      "kewarganegaraan": "WNI",
      "agama": "Islam",
      "pekerjaan": "Ibu Rumah Tangga",
      "alamat": "Jl. Sudirman No. 456",
      "status_keberadaan": "Hidup"
    }
  },
  "wali_nikah": {
    "nik": "3201010101010005",
    "nama_lengkap": "Ahmad Hidayat",
    "tempat_lahir": "Bandung",
    "tanggal_lahir": "1968-06-10",
    "kewarganegaraan": "WNI",
    "agama": "Islam",
    "pekerjaan": "Wiraswasta",
    "alamat": "Jl. Sudirman No. 456",
    "hubungan": "Ayah Kandung",
    "status_keberadaan": "Hidup",
    "status_kehadiran": "Belum Diketahui"
  },
  "jadwal_dan_lokasi": {
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "alamat_akad": "",
    "nomor_dispensasi": ""
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_pendaftaran": "Draft",
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA"
  }
}
```

---

#### 2.2 Check Registration Status
**GET** `/simnikah/pendaftaran/status`

**Auth:** ✅  
**Role:** `user_biasa`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_pendaftaran": "Menunggu Verifikasi",
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA"
  }
}
```

---

#### 2.3 Get All Registrations
**GET** `/simnikah/pendaftaran`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Query Parameters:**
- `status` (optional): Filter by status
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": {
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIK1704067200",
        "pendaftar_id": "USR1704067200",
        "status_pendaftaran": "Menunggu Verifikasi",
        "status_bimbingan": "Belum",
        "tanggal_pendaftaran": "2024-01-15T00:00:00Z",
        "tanggal_nikah": "2024-02-14T00:00:00Z",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "",
        "nomor_dispensasi": "",
        "penghulu_id": null,
        "catatan": "",
        "calon_suami": {
          "id": "CS1704067201",
          "nama_lengkap": "Ahmad Fauzi"
        },
        "calon_istri": {
          "id": "CI1704067202",
          "nama_lengkap": "Siti Aisyah"
        },
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 10,
      "total_records": 100,
      "per_page": 10,
      "has_next": true,
      "has_previous": false
    },
    "filters": {
      "status": "",
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

#### 2.4 Mark as Visited
**POST** `/simnikah/pendaftaran/:id/mark-visited`

**Auth:** ✅  
**Role:** `user_biasa`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Status berhasil diupdate ke Menunggu Penugasan"
}
```

---

#### 2.5 Get Status Flow
**GET** `/simnikah/pendaftaran/:id/status-flow`

**Auth:** ✅  
**Role:** All

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Menunggu Verifikasi",
    "status_flow": [
      {
        "status": "Draft",
        "completed": true,
        "current": false
      },
      {
        "status": "Menunggu Verifikasi",
        "completed": false,
        "current": true
      }
    ]
  }
}
```

---

#### 2.6 Update Wedding Address
**PUT** `/simnikah/pendaftaran/:id/alamat`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin"
}
```

---

#### 2.7 Update Wedding Location
**PUT** `/simnikah/pendaftaran/:id/location`

**Auth:** ✅  
**Role:** All

**Request Body:**
```json
{
  "latitude": -3.3194374,
  "longitude": 114.5900675,
  "alamat_akad": "Jl. Pangeran Samudra No. 88, Banjarmasin"
}
```

---

#### 2.8 Get Wedding Location Detail
**GET** `/simnikah/pendaftaran/:id/location`

**Auth:** ✅  
**Role:** All

**Response (200 OK):**
```json
{
  "success": true,
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

### 3. Staff Management

#### 3.1 Create Staff
**POST** `/simnikah/staff`

**Auth:** ✅  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "username": "staff001",
  "email": "staff@kua.go.id",
  "password": "staff123",
  "nama": "Staff Verifikasi",
  "nip": "198002021995022001",
  "jabatan": "Staff",
  "bagian": "Verifikasi",
  "no_hp": "081234567890",
  "alamat": "Jl. KUA No. 1"
}
```

---

#### 3.2 Get All Staff
**GET** `/simnikah/staff`

**Auth:** ✅  
**Role:** `kepala_kua`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "user_id": "STF1704067201",
      "username": "staff001",
      "nama": "Staff Verifikasi",
      "nip": "198002021995022001",
      "jabatan": "Staff",
      "bagian": "Verifikasi",
      "status": "Aktif"
    }
  ]
}
```

---

#### 3.3 Update Staff
**PUT** `/simnikah/staff/:id`

**Auth:** ✅  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "nama": "Staff Verifikasi Updated",
  "bagian": "Verifikasi & Dokumentasi",
  "no_hp": "081234567891"
}
```

---

#### 3.4 Verify Formulir
**POST** `/simnikah/staff/verify-formulir/:id`

**Auth:** ✅  
**Role:** `staff`

**Request Body:**
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir sudah lengkap dan valid"
}
```

**Status Options:**
- `"Formulir Disetujui"` → Next: "Menunggu Pengumpulan Berkas"
- `"Ditolak"` → Final status

---

#### 3.5 Verify Berkas
**POST** `/simnikah/staff/verify-berkas/:id`

**Auth:** ✅  
**Role:** `staff`

**Request Body:**
```json
{
  "status": "Berkas Diterima",
  "catatan": "Berkas sudah lengkap dan sesuai"
}
```

**Status Options:**
- `"Berkas Diterima"` → Next: "Berkas Diterima"
- `"Ditolak"` → Final status

---

### 4. Penghulu Management

#### 4.1 Create Penghulu
**POST** `/simnikah/penghulu`

**Auth:** ✅  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "username": "penghulu001",
  "email": "penghulu@kua.go.id",
  "password": "penghulu123",
  "nama": "H. Abdul Rahman, S.Ag",
  "nip": "196503031985031002",
  "no_hp": "081234567890",
  "alamat": "Jl. KUA No. 2"
}
```

---

#### 4.2 Get All Penghulu
**GET** `/simnikah/penghulu`

**Auth:** ✅  
**Role:** All

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": "PNG1704067202",
      "nama_lengkap": "H. Abdul Rahman, S.Ag",
      "nip": "196503031985031002",
      "status": "Aktif",
      "jumlah_nikah": 150,
      "rating": 4.8
    }
  ]
}
```

---

#### 4.3 Update Penghulu
**PUT** `/simnikah/penghulu/:id`

**Auth:** ✅  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "nama_lengkap": "H. Abdul Rahman, S.Ag, M.Ag",
  "no_hp": "081234567891",
  "alamat": "Jl. KUA No. 2 Updated"
}
```

---

#### 4.4 Verify Documents (Penghulu)
**POST** `/simnikah/penghulu/verify-documents/:id`

**Auth:** ✅  
**Role:** `penghulu`

**Request Body:**
```json
{
  "status": "Menunggu Pelaksanaan",
  "catatan": "Dokumen sudah lengkap dan valid"
}
```

**Status Options:**
- `"Menunggu Pelaksanaan"` → Next: "Menunggu Bimbingan"
- `"Ditolak"` → Final status

---

#### 4.5 Get Assigned Registrations
**GET** `/simnikah/penghulu/assigned-registrations`

**Auth:** ✅  
**Role:** `penghulu`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nomor_pendaftaran": "NIK1704067200",
      "status_pendaftaran": "Menunggu Verifikasi Penghulu",
      "tanggal_nikah": "2024-02-14",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di KUA"
    }
  ]
}
```

---

#### 4.6 Assign Penghulu
**POST** `/simnikah/pendaftaran/:id/assign-penghulu`

**Auth:** ✅  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "penghulu_id": 1
}
```

**Response (200 OK):**
```json
{
  "message": "Penghulu berhasil diassign",
  "data": {
    "pendaftaran_id": 1,
    "penghulu_id": 1,
    "penghulu_nama": "H. Abdul Rahman, S.Ag",
    "assigned_by": "KKUA1704067203",
    "assigned_at": "2024-01-15T10:30:00Z"
  }
}
```

---

#### 4.7 Change Penghulu
**PUT** `/simnikah/pendaftaran/:id/change-penghulu`

**Auth:** ✅  
**Role:** `kepala_kua`

**Request Body:**
```json
{
  "penghulu_id": 2
}
```

---

#### 4.8 Get Pendaftaran Belum Assign Penghulu
**GET** `/simnikah/pendaftaran/belum-assign-penghulu`

**Auth:** ✅  
**Role:** `kepala_kua`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nomor_pendaftaran": "NIK1704067200",
      "status_pendaftaran": "Menunggu Penugasan",
      "tanggal_nikah": "2024-02-14",
      "waktu_nikah": "09:00"
    }
  ]
}
```

---

#### 4.9 Get Penghulu Ketersediaan
**GET** `/simnikah/penghulu/:id/ketersediaan/:tanggal`

**Auth:** ✅  
**Role:** `kepala_kua`

**URL Parameters:**
- `id`: Penghulu ID
- `tanggal`: Format `YYYY-MM-DD`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "penghulu_id": 1,
    "tanggal": "2024-02-14",
    "jumlah_jadwal": 2,
    "max_per_hari": 3,
    "available": true,
    "jadwal": [
      {
        "waktu": "09:00",
        "pendaftaran_id": 1
      },
      {
        "waktu": "11:00",
        "pendaftaran_id": 2
      }
    ]
  }
}
```

---

### 5. Status Management

#### 5.1 Update Status Flexible
**PUT** `/simnikah/pendaftaran/:id/update-status`

**Auth:** ✅  
**Role:** `staff`, `penghulu`, `kepala_kua`

**Request Body:**
```json
{
  "status": "Menunggu Bimbingan",
  "catatan": "Status diupdate secara manual"
}
```

**Status yang BISA diupdate:**
- `"Draft"`
- `"Menunggu Verifikasi"`
- `"Menunggu Pengumpulan Berkas"`
- `"Berkas Diterima"`
- `"Menunggu Bimbingan"`
- `"Sudah Bimbingan"`
- `"Selesai"`
- `"Ditolak"`

**Status yang TIDAK BISA diupdate (hanya via assign-penghulu):**
- `"Menunggu Penugasan"`
- `"Penghulu Ditugaskan"`
- `"Menunggu Verifikasi Penghulu"`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Status berhasil diupdate",
  "data": {
    "id": 1,
    "status_sebelumnya": "Menunggu Verifikasi",
    "status_sekarang": "Menunggu Bimbingan",
    "updated_by": "STF1704067201",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

#### 5.2 Complete Bimbingan
**PUT** `/simnikah/pendaftaran/:id/complete-bimbingan`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Response (200 OK):**
```json
{
  "message": "Bimbingan perkawinan berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Sudah Bimbingan"
  }
}
```

---

#### 5.3 Complete Nikah
**PUT** `/simnikah/pendaftaran/:id/complete-nikah`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Response (200 OK):**
```json
{
  "message": "Nikah berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Selesai",
    "completed_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 6. Jadwal & Kalender

#### 6.1 Create Jadwal Nikah
**POST** `/simnikah/jadwal`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "tanggal": "2024-02-14",
  "waktu": "09:00",
  "tempat": "Di KUA"
}
```

---

#### 6.2 Get Jadwal Nikah
**GET** `/simnikah/jadwal`

**Auth:** ✅  
**Role:** All

**Query Parameters:**
- `tanggal` (optional): Filter by date `YYYY-MM-DD`
- `penghulu_id` (optional): Filter by penghulu

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "tanggal": "2024-02-14",
      "waktu": "09:00",
      "tempat": "Di KUA",
      "penghulu_id": 1,
      "penghulu_nama": "H. Abdul Rahman, S.Ag"
    }
  ]
}
```

---

#### 6.3 Update Jadwal Nikah
**PUT** `/simnikah/jadwal/:id`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "waktu": "10:00"
}
```

---

#### 6.4 Get Kalender Ketersediaan
**GET** `/simnikah/kalender-ketersediaan`

**Auth:** ❌ Tidak diperlukan

**Query Parameters:**
- `bulan` (optional): Month `YYYY-MM` (default: current month)
- `tempat` (optional): `"Di KUA"` or `"Di Luar KUA"`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "bulan": "2024-02",
    "kalender": [
      {
        "tanggal": "2024-02-14",
        "total_nikah": 5,
        "nikah_di_kua": 3,
        "nikah_di_luar_kua": 2,
        "available": true
      }
    ]
  }
}
```

---

#### 6.5 Get Kalender Tanggal Detail
**GET** `/simnikah/kalender-tanggal-detail`

**Auth:** ❌ Tidak diperlukan

**Query Parameters:**
- `tanggal` (required): `YYYY-MM-DD`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-02-14",
    "total_nikah": 5,
    "nikah_di_kua": 3,
    "nikah_di_luar_kua": 2,
    "slot_waktu": [
      {
        "waktu": "08:00",
        "available": false,
        "pendaftaran_id": 1
      },
      {
        "waktu": "09:00",
        "available": true
      }
    ]
  }
}
```

---

#### 6.6 Get Ketersediaan Tanggal
**GET** `/simnikah/ketersediaan-tanggal/:tanggal`

**Auth:** ❌ Tidak diperlukan

**URL Parameters:**
- `tanggal`: Format `YYYY-MM-DD`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-02-14",
    "total_nikah": 5,
    "max_per_hari": 9,
    "available": true,
    "slot_tersedia": 4
  }
}
```

---

#### 6.7 Get Penghulu Jadwal
**GET** `/simnikah/penghulu-jadwal/:tanggal`

**Auth:** ❌ Tidak diperlukan

**URL Parameters:**
- `tanggal`: Format `YYYY-MM-DD`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "penghulu_id": 1,
      "penghulu_nama": "H. Abdul Rahman, S.Ag",
      "jumlah_jadwal": 2,
      "max_per_hari": 3,
      "available": true
    }
  ]
}
```

---

### 7. Bimbingan Perkawinan

#### 7.1 Create Bimbingan
**POST** `/simnikah/bimbingan`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "tanggal": "2024-02-07",
  "waktu": "09:00",
  "tempat": "Aula KUA",
  "materi": "Bimbingan Perkawinan",
  "kapasitas": 10
}
```

**Note:** Tanggal harus hari Rabu

---

#### 7.2 Get All Bimbingan
**GET** `/simnikah/bimbingan`

**Auth:** ✅  
**Role:** All

**Query Parameters:**
- `status` (optional): Filter by status
- `tanggal` (optional): Filter by date

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "tanggal": "2024-02-07",
      "waktu": "09:00",
      "tempat": "Aula KUA",
      "status": "Aktif",
      "kapasitas": 10,
      "terdaftar": 5
    }
  ]
}
```

---

#### 7.3 Get Bimbingan By ID
**GET** `/simnikah/bimbingan/:id`

**Auth:** ✅  
**Role:** All

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "tanggal": "2024-02-07",
    "waktu": "09:00",
    "tempat": "Aula KUA",
    "materi": "Bimbingan Perkawinan",
    "status": "Aktif",
    "kapasitas": 10,
    "terdaftar": 5,
    "participants": []
  }
}
```

---

#### 7.4 Update Bimbingan
**PUT** `/simnikah/bimbingan/:id`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "waktu": "10:00",
  "tempat": "Aula KUA Updated"
}
```

---

#### 7.5 Get Bimbingan Kalender
**GET** `/simnikah/bimbingan-kalender`

**Auth:** ✅  
**Role:** All

**Query Parameters:**
- `bulan` (optional): `YYYY-MM`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "bulan": "2024-02",
    "kalender": [
      {
        "tanggal": "2024-02-07",
        "bimbingan_id": 1,
        "waktu": "09:00",
        "terdaftar": 5,
        "kapasitas": 10
      }
    ]
  }
}
```

---

#### 7.6 Daftar Bimbingan
**POST** `/simnikah/bimbingan/:id/daftar`

**Auth:** ✅  
**Role:** `user_biasa`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil mendaftar bimbingan perkawinan",
  "data": {
    "bimbingan_id": 1,
    "pendaftaran_id": 1,
    "status_kehadiran": "Belum"
  }
}
```

---

#### 7.7 Get Bimbingan Participants
**GET** `/simnikah/bimbingan/:id/participants`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "pendaftaran_id": 1,
      "nomor_pendaftaran": "NIK1704067200",
      "calon_suami": "Ahmad Fauzi",
      "calon_istri": "Siti Aisyah",
      "status_kehadiran": "Belum",
      "status_sertifikat": "Belum"
    }
  ]
}
```

---

#### 7.8 Update Attendance
**PUT** `/simnikah/bimbingan/:id/update-attendance`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "pendaftaran_id": 1,
  "status_kehadiran": "Hadir"
}
```

**Status Options:**
- `"Belum"`
- `"Hadir"`
- `"Tidak Hadir"`

---

#### 7.9 Cetak Undangan Bimbingan
**GET** `/simnikah/bimbingan/:id/undangan`

**Auth:** ✅  
**Role:** All

**Response:** PDF file atau HTML

---

#### 7.10 Cetak Undangan Semua
**GET** `/simnikah/bimbingan/:id/undangan-semua`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Response:** PDF file atau HTML

---

### 8. Notifikasi

#### 8.1 Create Notification
**POST** `/simnikah/notifikasi`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "user_id": "USR1704067200",
  "judul": "Pemberitahuan Penting",
  "pesan": "Formulir Anda telah disetujui",
  "tipe": "Success",
  "link": "/pendaftaran/1"
}
```

**Tipe Options:**
- `"Info"`
- `"Success"`
- `"Warning"`
- `"Error"`

---

#### 8.2 Get User Notifications
**GET** `/simnikah/notifikasi/user/:user_id`

**Auth:** ✅  
**Role:** All (hanya bisa akses notifikasi sendiri)

**Query Parameters:**
- `status_baca` (optional): `"Belum Dibaca"` or `"Sudah Dibaca"`
- `limit` (optional): Number of notifications

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "judul": "Pemberitahuan Penting",
      "pesan": "Formulir Anda telah disetujui",
      "tipe": "Success",
      "status_baca": "Belum Dibaca",
      "link": "/pendaftaran/1",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 10,
  "unread": 5
}
```

---

#### 8.3 Get Notification By ID
**GET** `/simnikah/notifikasi/:id`

**Auth:** ✅  
**Role:** All

---

#### 8.4 Update Notification Status
**PUT** `/simnikah/notifikasi/:id/status`

**Auth:** ✅  
**Role:** All

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

---

#### 8.5 Mark All as Read
**PUT** `/simnikah/notifikasi/user/:user_id/mark-all-read`

**Auth:** ✅  
**Role:** All

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Semua notifikasi telah ditandai sebagai sudah dibaca",
  "updated_count": 5
}
```

---

#### 8.6 Delete Notification
**DELETE** `/simnikah/notifikasi/:id`

**Auth:** ✅  
**Role:** All

---

#### 8.7 Get Notification Stats
**GET** `/simnikah/notifikasi/user/:user_id/stats`

**Auth:** ✅  
**Role:** All

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total": 10,
    "unread": 5,
    "read": 5,
    "by_type": {
      "Info": 3,
      "Success": 4,
      "Warning": 2,
      "Error": 1
    }
  }
}
```

---

#### 8.8 Send Notification to Role
**POST** `/simnikah/notifikasi/send-to-role`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "role": "user_biasa",
  "judul": "Pemberitahuan Umum",
  "pesan": "Pendaftaran nikah online sudah tersedia",
  "tipe": "Info"
}
```

---

#### 8.9 Run Reminder Notification
**POST** `/simnikah/notifikasi/run-reminder`

**Auth:** ✅  
**Role:** `staff`, `kepala_kua`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Reminder notification berhasil dijalankan",
  "sent_count": 5
}
```

---

### 9. Location & Maps

#### 9.1 Geocode (Address → Coordinates)
**POST** `/simnikah/location/geocode`

**Auth:** ✅  
**Role:** All

**Request Body:**
```json
{
  "alamat": "Jl. Pangeran Samudra No. 88, Banjarmasin"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "alamat": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "latitude": -3.3194374,
    "longitude": 114.5900675,
    "formatted_address": "Jl. Pangeran Samudra No.88, Banjarmasin Utara, Banjarmasin, Kalimantan Selatan"
  }
}
```

---

#### 9.2 Reverse Geocode (Coordinates → Address)
**POST** `/simnikah/location/reverse-geocode`

**Auth:** ✅  
**Role:** All

**Request Body:**
```json
{
  "latitude": -3.3194374,
  "longitude": 114.5900675
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "latitude": -3.3194374,
    "longitude": 114.5900675,
    "alamat": "Jl. Pangeran Samudra No. 88, Banjarmasin",
    "formatted_address": "Jl. Pangeran Samudra No.88, Banjarmasin Utara, Banjarmasin, Kalimantan Selatan"
  }
}
```

---

#### 9.3 Search Address
**GET** `/simnikah/location/search`

**Auth:** ✅  
**Role:** All

**Query Parameters:**
- `q` (required): Search query
- `limit` (optional): Max results (default: 5)

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "display_name": "Jl. Pangeran Samudra, Banjarmasin",
      "latitude": -3.3194374,
      "longitude": 114.5900675
    }
  ]
}
```

---

#### 9.4 Get Address Coordinates
**GET** `/simnikah/geocoding/coordinates`

**Auth:** ✅  
**Role:** All

**Query Parameters:**
- `address` (required): Address string

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "latitude": -3.3194374,
    "longitude": 114.5900675
  }
}
```

---

## ⚠️ Error Handling

### Error Response Format
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description",
  "type": "validation|authentication|authorization|not_found|database"
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 500 | Internal Server Error |

### Common Errors

#### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token tidak valid atau sudah expired"
}
```

#### 403 Forbidden
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak memiliki akses untuk endpoint ini"
}
```

#### 400 Bad Request
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Field 'tanggal_nikah' tidak boleh di masa lalu",
  "field": "tanggal_nikah",
  "type": "validation"
}
```

---

## 🔄 Status Flow

### Status Pendaftaran Nikah

1. **Draft** - Formulir masih dalam pengisian
2. **Menunggu Verifikasi** - Menunggu verifikasi formulir online oleh Staff
3. **Menunggu Pengumpulan Berkas** - User harus datang ke KUA dengan berkas
4. **Berkas Diterima** - Berkas fisik sudah diterima
5. **Menunggu Penugasan** - Menunggu Kepala KUA assign penghulu
6. **Penghulu Ditugaskan** - Penghulu sudah ditugaskan (auto transition)
7. **Menunggu Verifikasi Penghulu** - Penghulu verifikasi dokumen
8. **Menunggu Bimbingan** - User harus ikut bimbingan perkawinan
9. **Sudah Bimbingan** - Bimbingan sudah selesai
10. **Selesai** - Nikah telah dilaksanakan ✅
11. **Ditolak** - Pendaftaran ditolak ❌

---

## 📐 Business Rules

### Kapasitas Nikah
- Max 9 nikah di KUA per hari
- Max 3 nikah per penghulu per hari
- Min 60 menit gap waktu antar nikah
- Jam operasional: 08:00 - 16:00

### Dispensasi
Wajib dispensasi jika:
- Nikah < 10 hari kerja dari pendaftaran
- Usia calon suami < 19 tahun
- Usia calon istri < 19 tahun

### Bimbingan Perkawinan
- Hanya hari Rabu
- Max 10 pasangan per sesi
- 1 sesi per Rabu

### Wali Nikah
Urutan wali nasab (syariat Islam):
1. Ayah Kandung
2. Kakek
3. Saudara Laki-Laki Kandung
4. Paman Kandung
5. Wali Hakim (jika tidak ada wali nasab)

---

## 🔗 Quick Reference

### Endpoint by Role

#### User Biasa
- `POST /simnikah/pendaftaran/form-baru`
- `GET /simnikah/pendaftaran/status`
- `POST /simnikah/pendaftaran/:id/mark-visited`
- `POST /simnikah/bimbingan/:id/daftar`
- `GET /simnikah/notifikasi/user/:user_id`

#### Staff
- `GET /simnikah/pendaftaran`
- `POST /simnikah/staff/verify-formulir/:id`
- `POST /simnikah/staff/verify-berkas/:id`
- `PUT /simnikah/pendaftaran/:id/update-status`
- `POST /simnikah/bimbingan`
- `PUT /simnikah/bimbingan/:id/update-attendance`

#### Penghulu
- `GET /simnikah/penghulu/assigned-registrations`
- `POST /simnikah/penghulu/verify-documents/:id`
- `GET /simnikah/pendaftaran/:id/location`
- `PUT /simnikah/pendaftaran/:id/update-status`

#### Kepala KUA
- `POST /simnikah/staff`
- `POST /simnikah/penghulu`
- `POST /simnikah/pendaftaran/:id/assign-penghulu`
- `GET /simnikah/pendaftaran/belum-assign-penghulu`
- `GET /simnikah/penghulu/:id/ketersediaan/:tanggal`

---

## 📝 Notes

1. **Rate Limiting:** 100 requests per minute per IP (global), 5 requests per minute untuk auth endpoints
2. **Token Expiry:** JWT token valid untuk 24 jam
3. **CORS:** Sudah dikonfigurasi untuk development dan production
4. **Date Format:** Gunakan format `YYYY-MM-DD` untuk tanggal
5. **Time Format:** Gunakan format `HH:MM` (24-hour) untuk waktu

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0

