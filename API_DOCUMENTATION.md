# 📘 SIPENA API — Dokumentasi Lengkap

> **SIPENA** — Sistem Pendukung Keputusan untuk distribusi beban kerja penghulu menggunakan metode **Forward Chaining**.
>
> Bukan aplikasi pendaftaran nikah yang memvalidasi berkas (seperti SIMKAH), melainkan murni SPK untuk penjadwalan dan optimasi penugasan penghulu.

**Versi:** 3.0.0
**Update Terakhir:** Juli 2026
**Base URL:** `http://localhost:8080`
**Production:** Sesuai domain Railway/Vercel
**Authentication:** Bearer Token (JWT), dikirim via header `Authorization: Bearer <token>`
**Database:** MySQL/MariaDB (GORM ORM)
**Framework:** Gin (Golang)

---

## 📋 Daftar Isi

1. [Architecture Overview](#1-architecture-overview)
2. [Data Models](#2-data-models)
3. [Authentication & User Management](#3-authentication--user-management)
4. [Stage 1 — Catin: Check Schedule & Register](#4-stage-1--catin-check-schedule--register)
5. [Stage 2 — Kepala KUA: Forward Chaining Recommendation](#5-stage-2--kepala-kua-forward-chaining-recommendation)
6. [Stage 3 — Kepala KUA: Assignment Approval](#6-stage-3--kepala-kua-assignment-approval)
7. [Stage 4 — Penghulu: View Assignments](#7-stage-4--penghulu-view-assignments)
8. [Location & Geocoding](#8-location--geocoding)
9. [Staff Management](#9-staff-management)
10. [Notification System](#10-notification-system)
11. [Dashboard & Analytics](#11-dashboard--analytics)
12. [Forward Chaining Rules Reference](#12-forward-chaining-rules-reference)
13. [Error Handling & Status Codes](#13-error-handling--status-codes)
14. [Rate Limiting](#14-rate-limiting)
15. [Environment Variables & Configuration](#15-environment-variables--configuration)
16. [Health Check](#16-health-check)
17. [PRD v2.0 (ringkas)](docs/PRD_v2.md)

---

## 1. Architecture Overview

### 4-Stage Scheduling Flow

```
Stage 1: CATIN                    Stage 2: KEPALA KUA             Stage 3: KEPALA KUA            Stage 4: PENGHULU
┌─────────────────┐               ┌──────────────────┐            ┌──────────────────┐          ┌──────────────────┐
│ POST /check-    │               │ GET /recom-      │            │ POST /assign/:id │          │ GET /jadwal-     │
│   schedule      │──> input ──>  │   mendation/:id  │──> FC ──>  │  (Transaction +  │──> DB ──>│   penugasan      │
│ POST /pendaftaran│              │ GET /evaluation/ │            │   Row Lock)      │          │ (Maps URL, Geo)  │
└─────────────────┘               │   :id            │            └──────────────────┘          └──────────────────┘
                                  └──────────────────┘
```

### Roles & RBAC

| Role           | Description                        | Key Endpoints                                        |
|----------------|------------------------------------|------------------------------------------------------|
| `user_biasa`   | Calon Pengantin (Catin)            | Check schedule, Create registration, View own status |
| `kepala_kua`   | Kepala KUA — decision maker        | FC recommendation, Assign penghulu, Dashboard, Staff mgmt |
| `staff`        | Staff KUA — assists operations     | Create registration, Update status, Notifications, Dashboard |
| `penghulu`     | Penghulu — receives assignments    | View jadwal penugasan, Update coordinates            |

### Status Flow (State Machine)

```
Menunggu Penugasan  ──>  Penghulu Ditugaskan  ──>  Selesai
       │                        │
       └──────> Ditolak  <──────┘ (tidak valid)
```

**Valid State Transitions:**

| Dari                   | Ke                          |
|------------------------|-----------------------------|
| `Menunggu Penugasan`   | `Penghulu Ditugaskan`       |
| `Menunggu Penugasan`   | `Ditolak`                   |
| `Penghulu Ditugaskan`  | `Selesai`                   |

> **Note:** Status `Draft` dan `Disetujui` telah dihapus. Pendaftaran langsung masuk ke `Menunggu Penugasan` saat dibuat.

### Project Structure

```
simnikah-api/
├── cmd/api/main.go                          # Entry point & route registration
├── config/
│   ├── config.go                            # Database connection (MySQL)
│   ├── db.go                                # DB helper
│   └── indexes.go                           # Database index optimization
├── internal/
│   ├── handlers/
│   │   ├── auth/auth.go                     # Register, Login, Profile, Upload Photo
│   │   ├── catin/
│   │   │   ├── daftar.go                    # Check Schedule, Registration CRUD
│   │   │   └── location.go                  # Geocoding, Reverse Geocoding, Search
│   │   ├── kepala_kua/
│   │   │   ├── forward_chaining_handlers.go # FC Recommendation, Evaluation, Assign
│   │   │   └── kepala_kua.go                # Available Penghulu, Schedule View
│   │   ├── penghulu/penghulu.go             # Jadwal Penugasan, Update Coordinates
│   │   ├── staff/staff.go                   # Staff CRUD, Create Registration, Update Status
│   │   ├── notification/notification.go     # Full notification CRUD
│   │   └── dashboard/dashboard.go           # Kepala KUA & Staff Dashboard
│   ├── middleware/
│   │   ├── auth.go                          # JWT Auth, Role, MultiRole middleware
│   │   └── rate_limit.go                    # Global & Strict rate limiter
│   ├── models/
│   │   ├── models.go                        # PendaftaranNikah, Users, Penghulu, StaffKUA, Notifikasi
│   │   └── constants.go                     # Status, Role, Tempat constants
│   ├── seeders/                             # Initial data seeding
│   └── services/
│       ├── forward_chaining_engine.go       # Core FC Engine (~47KB)
│       ├── notification_service.go          # Notification business logic
│       └── cron_job.go                      # Scheduled reminder notifications
├── pkg/
│   ├── cache/                               # Geocoding cache
│   ├── crypto/                              # Password hashing (bcrypt)
│   ├── storage/                             # ImgBB file upload
│   ├── utils/                               # JWT, random string, timezone (WITA)
│   └── validator/                           # Input validation helpers
├── migrations/                              # Database migrations
└── docs/                                    # Documentation files
```

---

## 2. Data Models

### PendaftaranNikah (Marriage Registration)

| Field                | Type        | Required | FC Variable | JSON Key             | Description                                    |
|----------------------|-------------|----------|-------------|----------------------|------------------------------------------------|
| `id`                 | `uint`      | auto     | —           | `id`                 | Primary key                                    |
| `nomor_pendaftaran`  | `string(20)`| auto     | —           | `nomor_pendaftaran`  | Auto-generated: `REG-{year}-{6digits}`         |
| `pendaftar_id`       | `string(20)`| auto     | —           | `pendaftar_id`       | User ID registrant (from JWT)                  |
| `nama_suami`         | `string(100)`| **Yes** | **No**      | `nama_suami`         | Groom name (display reference from Excel)      |
| `umur_suami`         | `int`       | No       | **No**      | `umur_suami`         | Groom age (display reference)                  |
| `nama_istri`         | `string(100)`| **Yes** | **No**      | `nama_istri`         | Bride name (display reference from Excel)      |
| `umur_istri`         | `int`       | No       | **No**      | `umur_istri`         | Bride age (display reference)                  |
| `tanggal_nikah`      | `date`      | **Yes**  | **Yes**     | `tanggal_nikah`      | Wedding date (`YYYY-MM-DD`)                    |
| `waktu_nikah`        | `string(10)`| **Yes**  | **Yes**     | `waktu_nikah`        | Wedding time (`HH:MM`)                         |
| `tempat_nikah`       | `string(100)`| **Yes** | **Yes**     | `tempat_nikah`       | `"Di KUA"` or `"Di Luar KUA"`                 |
| `alamat_akad`        | `string(200)`| Cond.   | **Yes**     | `alamat_akad`        | Full address (required if Di Luar KUA)         |
| `latitude`           | `*float64`  | No       | **Yes**     | `latitude`           | Latitude (auto-geocoded if empty)              |
| `longitude`          | `*float64`  | No       | **Yes**     | `longitude`          | Longitude (auto-geocoded if empty)             |
| `status_pendaftaran` | `string(40)`| auto     | —           | `status_pendaftaran` | Current status                                 |
| `penghulu_id`        | `*uint`     | auto     | —           | `penghulu_id`        | Assigned penghulu ID                           |
| `created_at`         | `timestamp` | auto     | —           | `dibuat_pada`        | Creation timestamp                             |
| `updated_at`         | `timestamp` | auto     | —           | `diperbarui_pada`    | Last update timestamp                          |

> **Critical:** Fields marked **FC Variable = Yes** are the ONLY inputs to Forward Chaining rules. Couple data (nama/umur) is display-only and never enters the inference engine.

### Users

| Field           | Type        | JSON Key       | Description                                            |
|-----------------|-------------|----------------|--------------------------------------------------------|
| `id`            | `uint`      | `id`           | Primary key (auto)                                     |
| `user_id`       | `string(20)`| `id_pengguna`  | Unique user ID (`USR-{12chars}`)                       |
| `username`      | `string(50)`| `nama_pengguna`| Login username (unique)                                |
| `email`         | `string(100)`| `email`       | Email (unique)                                         |
| `password`      | `string(255)`| `kata_sandi`  | Bcrypt hashed password                                 |
| `role`          | `string(20)`| `peran`        | `user_biasa`, `penghulu`, `staff`, `kepala_kua`        |
| `status`        | `string(20)`| `status`       | `Aktif`, `Nonaktif`, `Blokir`                         |
| `nama`          | `string(100)`| `nama`        | Full name                                              |
| `profile_photo` | `string(500)`| `foto_profil` | Photo URL (ImgBB)                                      |
| `created_at`    | `timestamp` | `dibuat_pada`  | Creation timestamp                                     |
| `updated_at`    | `timestamp` | `diperbarui_pada` | Last update timestamp                              |

### StaffKUA

| Field           | Type        | JSON Key        | Description                     |
|-----------------|-------------|-----------------|----------------------------------|
| `id`            | `uint`      | `id`            | Primary key                      |
| `user_id`       | `string(20)`| `id_pengguna`   | Linked user ID                   |
| `nip`           | `string(30)`| `nip`           | NIP (unique)                     |
| `nama_lengkap`  | `string(100)`| `nama_lengkap` | Full name                        |
| `jabatan`       | `string(50)`| `jabatan`       | Position                         |
| `bagian`        | `string(50)`| `bagian`        | Department/Section               |
| `no_hp`         | `string(15)`| `nomor_telepon` | Phone number                     |
| `email`         | `string(100)`| `email`        | Email                            |
| `alamat`        | `string(200)`| `alamat`       | Address                          |
| `status`        | `string(20)`| `status`        | `Aktif`, `Nonaktif`             |
| `created_at`    | `timestamp` | `dibuat_pada`   | Creation timestamp               |
| `updated_at`    | `timestamp` | `diperbarui_pada` | Last update timestamp          |

### Penghulu

| Field           | Type        | JSON Key        | Description                     |
|-----------------|-------------|-----------------|----------------------------------|
| `id`            | `uint`      | `id`            | Primary key                      |
| `user_id`       | `string(20)`| `id_pengguna`   | Linked user ID                   |
| `nip`           | `string(30)`| `nip`           | NIP (unique)                     |
| `nama_lengkap`  | `string(100)`| `nama_lengkap` | Full name                        |
| `no_hp`         | `string(15)`| `nomor_telepon` | Phone number                     |
| `email`         | `string(100)`| `email`        | Email                            |
| `alamat`        | `string(200)`| `alamat`       | Address                          |
| `latitude`      | `*float64`  | `latitude`      | Home latitude (for holiday routing) |
| `longitude`     | `*float64`  | `longitude`     | Home longitude (for holiday routing) |
| `status`        | `string(20)`| `status`        | `Aktif`, `Nonaktif`             |
| `jumlah_nikah`  | `int`       | `jumlah_nikah`  | Historical marriage count        |
| `rating`        | `float64`   | `rating`        | Performance rating (0-5)         |
| `created_at`    | `timestamp` | `dibuat_pada`   | Creation timestamp               |
| `updated_at`    | `timestamp` | `diperbarui_pada` | Last update timestamp          |

### Notifikasi

| Field           | Type        | JSON Key        | Description                     |
|-----------------|-------------|-----------------|----------------------------------|
| `id`            | `uint`      | `id`            | Primary key                      |
| `user_id`       | `string(20)`| `id_pengguna`   | Target user ID                   |
| `judul`         | `string(100)`| `judul`        | Notification title               |
| `pesan`         | `string(500)`| `pesan`        | Notification message             |
| `tipe`          | `string(10)`| `tipe`          | `Info`, `Warning`, `Error`, `Success` |
| `status_baca`   | `string(20)`| `status_dibaca` | `Belum Dibaca`, `Sudah Dibaca`  |
| `link`          | `string(200)`| `tautan`       | Optional navigation link         |
| `created_at`    | `timestamp` | `dibuat_pada`   | Creation timestamp               |
| `updated_at`    | `timestamp` | `diperbarui_pada` | Last update timestamp          |

---

## 3. Authentication & User Management

### POST `/register`

Register a new user account. Role is always `user_biasa` via public registration.

**Auth:** None
**Rate Limit:** Strict (5 req/min per IP)

**Request Body:**
```json
{
  "username": "ahmad_fauzi",
  "email": "ahmad@example.com",
  "password": "secret123",
  "nama": "Ahmad Fauzi"
}
```

**Validation Rules:**
- `username`: required, unique
- `email`: required, valid email format, unique
- `password`: required, min 6 characters
- `nama`: required

**Response `201 Created`:**
```json
{
  "success": true,
  "message": "User berhasil dibuat",
  "data": {
    "user_id": "USR-a3Bf9xK2mNp1",
    "username": "ahmad_fauzi",
    "email": "ahmad@example.com",
    "nama": "Ahmad Fauzi",
    "role": "user_biasa"
  }
}
```

**Error Responses:**

| Status | Condition                | Error Message                    |
|--------|--------------------------|----------------------------------|
| `400`  | Invalid input format     | `"Format data tidak valid"`      |
| `500`  | Username already exists  | `"username sudah digunakan"`     |
| `500`  | Email already registered | `"email sudah terdaftar"`        |

> **Security Note:** Username & email uniqueness is checked atomically within a GORM transaction using `SELECT ... FOR UPDATE` to prevent race conditions.

---

### POST `/login`

Authenticate and receive JWT token.

**Auth:** None
**Rate Limit:** Strict (5 req/min per IP)

**Request Body:**
```json
{
  "username": "ahmad_fauzi",
  "password": "secret123"
}
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "user_id": "USR-a3Bf9xK2mNp1",
    "username": "ahmad_fauzi",
    "email": "ahmad@example.com",
    "role": "user_biasa",
    "nama": "Ahmad Fauzi"
  },
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "user_id": "USR-a3Bf9xK2mNp1",
      "username": "ahmad_fauzi",
      "email": "ahmad@example.com",
      "role": "user_biasa",
      "nama": "Ahmad Fauzi"
    }
  }
}
```

> **Token expires in 24 hours.** Include in all subsequent requests: `Authorization: Bearer <token>`

**JWT Claims:**
```json
{
  "user_id": "USR-a3Bf9xK2mNp1",
  "email": "ahmad@example.com",
  "role": "user_biasa",
  "nama": "Ahmad Fauzi",
  "exp": 1737475200,
  "iat": 1737388800,
  "nbf": 1737388800
}
```

**Error Responses:**

| Status | Condition                 | Error Message                     |
|--------|---------------------------|-----------------------------------|
| `400`  | Invalid input format      | `"Format data tidak valid"`       |
| `401`  | User not found            | `"Username atau password salah"`  |
| `401`  | Wrong password            | `"Username atau password salah"`  |
| `401`  | Account not active        | `"Akun tidak aktif"`             |

---

### GET `/profile`

Get current user's profile.

**Auth:** Required (Bearer Token)

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Profile berhasil diambil",
  "data": {
    "user_id": "USR-a3Bf9xK2mNp1",
    "username": "ahmad_fauzi",
    "email": "ahmad@example.com",
    "role": "user_biasa",
    "nama": "Ahmad Fauzi",
    "status": "Aktif",
    "profile_photo": "https://i.ibb.co/..."
  }
}
```

---

### POST `/upload-photo`

Upload profile photo. Uploaded to ImgBB cloud storage.

**Auth:** Required
**Content-Type:** `multipart/form-data`

**Form Data:**

| Field   | Type | Max Size | Allowed Types     | Description              |
|---------|------|----------|-------------------|--------------------------|
| `photo` | file | 5 MB     | JPG, PNG, WebP    | Profile photo file       |

> **Security:** File type is validated using magic bytes (`http.DetectContentType`), not just the Content-Type header.

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Foto profil berhasil diupload",
  "data": {
    "profile_photo": "https://i.ibb.co/...",
    "user_id": "USR-a3Bf9xK2mNp1",
    "username": "ahmad_fauzi"
  }
}
```

**Error Responses:**

| Status | Condition            | Error Message                    |
|--------|----------------------|----------------------------------|
| `400`  | No file provided     | `"File tidak ditemukan"`         |
| `400`  | File too large       | `"Ukuran file terlalu besar"`    |
| `400`  | Invalid file type    | `"Tipe file tidak didukung"`     |

---

## 4. Stage 1 — Catin: Check Schedule & Register

### POST `/simnikah/check-schedule`

Check if a schedule slot is available before registering. Delegates to the Forward Chaining Engine to evaluate all active penghulu.

**Auth:** Required
**Roles:** All authenticated users

**Request Body:**
```json
{
  "tanggal_nikah": "2025-03-15",
  "waktu_nikah": "09:00",
  "tempat_nikah": "Di KUA",
  "alamat_nikah": "",
  "latitude": null,
  "longitude": null
}
```

**Request Fields:**

| Field          | Type     | Required | Description                         |
|----------------|----------|----------|-------------------------------------|
| `tanggal_nikah`| string   | **Yes**  | Date (`YYYY-MM-DD`)                 |
| `waktu_nikah`  | string   | No       | Time (`HH:MM`)                      |
| `tempat_nikah` | string   | No       | `"Di KUA"` or `"Di Luar KUA"`       |
| `alamat_nikah` | string   | No       | Full address (for geocoding)         |
| `latitude`     | *float64 | No       | Latitude coordinate                  |
| `longitude`    | *float64 | No       | Longitude coordinate                 |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Pengecekan jadwal selesai",
  "data": {
    "tanggal_nikah": "2025-03-15",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "is_holiday": false,
    "holiday_name": "",
    "total_booked": 1,
    "booked_in_kua": 0,
    "slot_remaining": 2,
    "available": true,
    "reason": "",
    "recommended_penghulu": {
      "penghulu_id": 1,
      "nama_penghulu": "H. Abdul Rahman, S.Ag",
      "jumlah_nikah": 12,
      "rating": 4.8,
      "slot_harian_terisi": 1,
      "slot_jam_terisi": 0,
      "tersedia": true,
      "alasan": "Memenuhi kapasitas harian dan kapasitas per jam"
    },
    "alternatives": [
      {
        "penghulu_id": 2,
        "nama_penghulu": "H. Muhammad Said, S.Ag",
        "jumlah_nikah": 8,
        "rating": 4.5,
        "slot_harian_terisi": 0,
        "slot_jam_terisi": 0,
        "tersedia": true,
        "alasan": "Memenuhi kapasitas harian dan kapasitas per jam"
      }
    ]
  }
}
```

**Response `409 Conflict` (slot full):**
```json
{
  "success": false,
  "message": "Jadwal tidak tersedia",
  "error": "Semua slot penghulu pada jam tersebut sudah penuh",
  "data": {
    "available": false,
    "total_booked": 3,
    "slot_remaining": 0
  }
}
```

---

### POST `/simnikah/pendaftaran`

Create a new marriage registration with couple data + scheduling data.

**Auth:** Required
**Roles:** All authenticated users

**Request Body:**
```json
{
  "nama_suami": "Ahmad Fauzi",
  "umur_suami": 28,
  "nama_istri": "Siti Aminah",
  "umur_istri": 25,
  "tanggal_nikah": "2025-03-15",
  "waktu_nikah": "09:00",
  "tempat_nikah": "Di Luar KUA",
  "alamat_akad": "Jl. Gatot Subroto No. 10, Banjarmasin Utara",
  "latitude": -3.2913,
  "longitude": 114.5881
}
```

**Validation Rules:**

| Field          | Rule                                          |
|----------------|-----------------------------------------------|
| `nama_suami`   | Required                                      |
| `nama_istri`   | Required                                      |
| `tanggal_nikah`| Required, `YYYY-MM-DD`, cannot be in the past (WITA timezone) |
| `waktu_nikah`  | Required, `HH:MM` format                     |
| `tempat_nikah` | Required, must be `"Di KUA"` or `"Di Luar KUA"` |
| `alamat_akad`  | Required if `tempat_nikah` is `"Di Luar KUA"` |
| `latitude`     | Optional, auto-geocoded from `alamat_akad` if empty |
| `longitude`    | Optional, auto-geocoded from `alamat_akad` if empty |

**Business Logic:**
1. Validate all input fields
2. Check schedule availability via Forward Chaining Engine
3. If slot unavailable → return `409 Conflict`
4. Auto-geocode address if coordinates not provided
5. Create registration with status `Menunggu Penugasan`
6. Auto-generate `nomor_pendaftaran` as `REG-{year}-{6digits}`

**Response `201 Created`:**
```json
{
  "success": true,
  "message": "Pendaftaran berhasil dibuat, menunggu penugasan penghulu",
  "data": {
    "id": 42,
    "nama_suami": "Ahmad Fauzi",
    "umur_suami": 28,
    "nama_istri": "Siti Aminah",
    "umur_istri": 25,
    "tanggal_nikah": "2025-03-15",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Gatot Subroto No. 10, Banjarmasin Utara",
    "latitude": -3.2913,
    "longitude": 114.5881,
    "status_pendaftaran": "Menunggu Penugasan"
  }
}
```

---

### GET `/simnikah/pendaftaran/status`

Get current logged-in catin's latest registration status.

**Auth:** Required

**Response `200 OK` (has registration):**
```json
{
  "success": true,
  "message": "Status pendaftaran berhasil diambil",
  "data": {
    "has_registration": true,
    "can_register": false,
    "registration": {
      "id": 42,
      "nomor_pendaftaran": "REG-2025-123456",
      "status_pendaftaran": "Penghulu Ditugaskan",
      "tanggal_nikah": "2025-03-15",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di Luar KUA",
      "alamat_akad": "Jl. Gatot Subroto No. 10",
      "created_at": "2025-01-20T10:30:00Z",
      "calon_suami": {
        "nama_lengkap": "Ahmad Fauzi",
        "nama_dan_bin": "Ahmad Fauzi",
        "nama": "Ahmad Fauzi"
      },
      "calon_istri": {
        "nama_lengkap": "Siti Aminah",
        "nama_dan_binti": "Siti Aminah",
        "nama": "Siti Aminah"
      },
      "penghulu": {
        "id": 1,
        "nama": "H. Abdul Rahman, S.Ag",
        "nama_lengkap": "H. Abdul Rahman, S.Ag"
      }
    }
  }
}
```

**Response `200 OK` (no registration):**
```json
{
  "success": true,
  "message": "Belum ada pendaftaran",
  "data": {
    "has_registration": false,
    "can_register": true,
    "registration": null
  }
}
```

> **Note:** `can_register` is `true` only when there is no registration or when the last registration has status `Ditolak`.

---

### GET `/simnikah/pendaftaran/:id`

Get detailed information of a single registration.

**Auth:** Required
**Access Control:** `user_biasa` can only view their own registration. `staff`, `kepala_kua` can view any.

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Detail pendaftaran berhasil diambil",
  "data": {
    "id": 42,
    "nomor_pendaftaran": "REG-2025-123456",
    "pendaftar_id": "USR-a3Bf9xK2mNp1",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "tanggal_nikah": "2025-03-15",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Gatot Subroto No. 10",
    "latitude": -3.2913,
    "longitude": 114.5881,
    "created_at": "2025-01-20T10:30:00Z",
    "updated_at": "2025-01-20T10:35:00Z",
    "calon_suami": {
      "nama_lengkap": "Ahmad Fauzi",
      "nama_dan_bin": "Ahmad Fauzi",
      "umur": 28
    },
    "calon_istri": {
      "nama_lengkap": "Siti Aminah",
      "nama_dan_binti": "Siti Aminah",
      "umur": 25
    },
    "wali_nikah": {
      "nama_dan_bin": "Wali Nikah",
      "hubungan_wali": "Ayah Kandung"
    },
    "penghulu": {
      "id": 1,
      "nip": "197001011990031001",
      "nama_lengkap": "H. Abdul Rahman, S.Ag",
      "no_hp": "081234567890",
      "email": "penghulu@kua.go.id",
      "alamat": "Banjarmasin",
      "status": "Aktif"
    },
    "location": {
      "latitude": -3.2913,
      "longitude": 114.5881,
      "has_coordinates": true,
      "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291300,114.588100",
      "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291300,114.588100",
      "osm_url": "https://www.openstreetmap.org/?mlat=-3.291300&mlon=114.588100&zoom=16"
    }
  }
}
```

---

### GET `/simnikah/pendaftaran`

List all registrations with pagination, search, and status filtering.

**Auth:** Required
**Roles:** `staff`, `kepala_kua` only

**Query Parameters:**

| Parameter | Type   | Default | Description                               |
|-----------|--------|---------|-------------------------------------------|
| `page`    | int    | `1`     | Page number                               |
| `limit`   | int    | `10`    | Items per page (max 100)                  |
| `status`  | string | —       | Filter by status pendaftaran              |
| `search`  | string | —       | Search nama_suami, nama_istri, or nomor_pendaftaran |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Daftar pendaftaran berhasil diambil",
  "data": {
    "registrations": [
      {
        "id": 42,
        "nomor_pendaftaran": "REG-2025-123456",
        "status_pendaftaran": "Menunggu Penugasan",
        "tanggal_nikah": "2025-03-15",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_akad": "Jl. Gatot Subroto No. 10",
        "created_at": "2025-01-20T10:30:00Z",
        "calon_suami": {
          "nama_lengkap": "Ahmad Fauzi",
          "nama_dan_bin": "Ahmad Fauzi",
          "nama": "Ahmad Fauzi"
        },
        "calon_istri": {
          "nama_lengkap": "Siti Aminah",
          "nama_dan_binti": "Siti Aminah",
          "nama": "Siti Aminah"
        },
        "penghulu": null
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_records": 45,
      "total_pages": 5
    }
  }
}
```

---

## 5. Stage 2 — Kepala KUA: Forward Chaining Recommendation

### GET `/simnikah/kepala-kua/forward-chaining/recommendation/:id`

Trigger Forward Chaining Engine to get penghulu recommendation for a specific registration.

**Auth:** Required
**Roles:** `kepala_kua` only

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Rekomendasi penghulu berhasil didapatkan",
  "data": {
    "recommended_penghulu_id": 1,
    "recommended_penghulu_name": "H. Abdul Rahman, S.Ag",
    "selected_score": 87.5,
    "confidence": 93.75,
    "decision_reasoning": "Forward chaining menetapkan H. Abdul Rahman, S.Ag sebagai kandidat terbaik dengan score 87.50, jarak 2.35 km, dan seluruh fakta utama terpenuhi.",
    "alternatives": [
      {
        "penghulu_id": 2,
        "nama_penghulu": "H. Muhammad Said, S.Ag",
        "score": 82.3,
        "reason": "Score 82.30, jarak 4.10 km"
      }
    ],
    "evaluated_at": "2025-01-20T10:30:00Z",
    "evaluation_count": 5
  }
}
```

---

### GET `/simnikah/kepala-kua/forward-chaining/evaluation/:id`

Get detailed evaluation report showing inference trace for ALL penghulu candidates.

**Auth:** Required
**Roles:** `kepala_kua` only

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Laporan evaluasi berhasil didapatkan",
  "data": {
    "registration_id": 42,
    "tanggal_nikah": "2025-03-15",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_nikah": "Jl. Gatot Subroto No. 10",
    "latitude": -3.2913,
    "longitude": 114.5881,
    "status_pendaftaran": "Menunggu Penugasan",
    "holiday": false,
    "holiday_name": "",
    "evaluated_at": "2025-01-20T10:31:00Z",
    "penghulu_evaluations": [
      {
        "penghulu_id": 1,
        "penghulu_nama": "H. Abdul Rahman, S.Ag",
        "rating": 4.8,
        "jumlah_nikah": 12,
        "status": "Aktif",
        "all_rules_passed": true,
        "score": 87.5,
        "distance_meters": 2350,
        "conclusion": "Penghulu direkomendasikan",
        "evaluated_rules": [
          {
            "rule_id": "RULE_001",
            "rule_name": "Validasi Administrasi",
            "is_satisfied": true,
            "reason": "Status pendaftaran memenuhi syarat penugasan",
            "impact": 10
          },
          {
            "rule_id": "RULE_004",
            "rule_name": "Cek Konflik Jadwal",
            "is_satisfied": true,
            "reason": "Tidak ada bentrok pada slot waktu yang sama",
            "impact": 20
          }
        ],
        "derived_facts": [
          {
            "name": "Lulus Syarat Administrasi",
            "value": true,
            "rule_id": "RULE_001",
            "reason": "Status pendaftaran memenuhi syarat penugasan"
          }
        ]
      }
    ]
  }
}
```

---

### GET `/simnikah/kepala-kua/forward-chaining/config`

Get Forward Chaining Engine configuration parameters.

**Auth:** Required
**Roles:** `kepala_kua` only

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Konfigurasi forward chaining engine",
  "data": {
    "source": "engine_defaults",
    "dynamic_config_ready": true,
    "minimum_rating": 3.0,
    "capacity_per_day": 3,
    "capacity_per_hour": 1,
    "kua_latitude": -3.291304649442475,
    "kua_longitude": 114.58814746634684,
    "scoring_weights": {
      "rating_weight": 0.35,
      "availability_weight": 0.25,
      "fairness_weight": 0.20,
      "location_match_weight": 0.10,
      "distance_weight": 0.10
    },
    "rule_constraint_notes": [
      "Flow ini hanya memakai constraint jadwal: status penghulu, bentrok slot, kapasitas harian, kapasitas per jam, dan lokasi nikah",
      "Sumber dinamis yang direkomendasikan: tabel system_configs",
      "Contoh key: forward_chaining.capacity_per_day, forward_chaining.capacity_per_hour, forward_chaining.minimum_rating",
      "Nilai engine saat ini tetap menjadi fallback jika config DB belum tersedia"
    ],
    "system_config_table_name": "system_configs",
    "system_config_keys_example": [
      "forward_chaining.minimum_rating",
      "forward_chaining.capacity_per_day",
      "forward_chaining.capacity_per_hour",
      "forward_chaining.weights.rating",
      "forward_chaining.weights.availability",
      "forward_chaining.weights.fairness",
      "forward_chaining.weights.location_match",
      "forward_chaining.weights.distance"
    ]
  }
}
```

---

## 6. Stage 3 — Kepala KUA: Assignment Approval

### POST `/simnikah/kepala-kua/forward-chaining/assign/:id`

Assign a penghulu to a registration. Uses **GORM Transaction + `SELECT ... FOR UPDATE` row locking** to prevent race conditions during concurrent approval.

**Auth:** Required
**Roles:** `kepala_kua` only

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Request Body:**
```json
{
  "penghulu_id": 1,
  "approval_notes": "Disetujui, penghulu tersedia dan lokasi dekat"
}
```

**Request Fields:**

| Field            | Type   | Required | Description                     |
|------------------|--------|----------|---------------------------------|
| `penghulu_id`    | uint   | **Yes**  | ID of the penghulu to assign    |
| `approval_notes` | string | **Yes**  | Reason/notes for the approval   |

**Business Logic:**
1. Validate role = `kepala_kua` (double-check against DB, not just JWT)
2. Validate `penghulu_id` is not 0
3. Validate `approval_notes` is not empty
4. Begin GORM transaction
5. `SELECT ... FOR UPDATE` on the registration row (row lock)
6. Verify status is `Menunggu Penugasan` (reject if not)
7. Verify penghulu exists and has status `Aktif`
8. Update: `penghulu_id`, `status_pendaftaran` → `Penghulu Ditugaskan`
9. Commit transaction

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Penghulu berhasil ditugaskan",
  "data": {
    "registration_id": 42,
    "nama_suami": "Ahmad Fauzi",
    "nama_istri": "Siti Aminah",
    "tanggal_nikah": "2025-03-15",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "penghulu_id": 1,
    "penghulu_nama": "H. Abdul Rahman, S.Ag",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "assigned_by": "USR-k7Hd2Xm9pQ4w",
    "assigned_at": "2025-01-20T10:35:00Z"
  }
}
```

**Error Responses:**

| Status | Condition                    | Error Message                                              |
|--------|------------------------------|------------------------------------------------------------|
| `400`  | Missing penghulu_id          | `"Penghulu ID harus disediakan"`                           |
| `400`  | Empty approval_notes         | `"Approval notes wajib diisi"`                             |
| `404`  | Registration not found       | `"Pendaftaran tidak ditemukan"`                            |
| `404`  | Penghulu not found/inactive  | `"Penghulu tidak valid atau tidak aktif"`                  |
| `409`  | Status not Menunggu Penugasan| `"status saat ini: X, harus 'Menunggu Penugasan'"`         |
| `403`  | Not Kepala KUA               | `"Hanya user dengan role Kepala KUA yang dapat mengakses"` |

---

### GET `/simnikah/kepala-kua/available-penghulu`

List all active penghulu sorted by rating.

**Auth:** Required
**Roles:** `kepala_kua` only

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Daftar penghulu aktif",
  "data": {
    "total": 5,
    "penghulus": [
      {
        "id": 1,
        "nama_lengkap": "H. Abdul Rahman, S.Ag",
        "rating": 4.8,
        "jumlah_nikah": 12,
        "status": "Aktif"
      }
    ]
  }
}
```

---

### GET `/simnikah/kepala-kua/penghulu-tersedia`

View penghulu schedule availability for a specific date. Shows how many slots each penghulu has booked.

**Auth:** Required
**Roles:** `kepala_kua` only

**Query Parameters:**

| Parameter  | Type   | Required | Description         |
|------------|--------|----------|---------------------|
| `tanggal`  | string | **Yes**  | Date `YYYY-MM-DD`   |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Jadwal penghulu berhasil diambil",
  "data": {
    "tanggal": "2025-03-15",
    "penghulus": [
      {
        "penghulu_id": 1,
        "nama_lengkap": "H. Abdul Rahman, S.Ag",
        "rating": 4.8,
        "jumlah_nikah": 12,
        "booked_slots": 2,
        "available": true
      },
      {
        "penghulu_id": 3,
        "nama_lengkap": "Ustadz Ahmad Yusuf",
        "rating": 4.2,
        "jumlah_nikah": 15,
        "booked_slots": 3,
        "available": false
      }
    ]
  }
}
```

> **Note:** `available` is `false` when `booked_slots >= capacity_per_day` (default: 3).

---

## 7. Stage 4 — Penghulu: View Assignments

### GET `/simnikah/penghulu/jadwal-penugasan`

View assigned tasks with full address, geolocation, and navigation URLs.

**Auth:** Required
**Roles:** `penghulu`, `kepala_kua`

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Jadwal penugasan berhasil diambil",
  "data": {
    "penghulu_id": 1,
    "penghulu_nama": "H. Abdul Rahman, S.Ag",
    "total": 2,
    "jadwal": [
      {
        "id": 42,
        "nama_suami": "Ahmad Fauzi",
        "umur_suami": 28,
        "nama_istri": "Siti Aminah",
        "umur_istri": 25,
        "tanggal_nikah": "2025-03-15",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_lengkap": "Jl. Gatot Subroto No. 10, Banjarmasin Utara",
        "status_pendaftaran": "Penghulu Ditugaskan",
        "has_coordinates": true,
        "latitude": -3.2913,
        "longitude": 114.5881,
        "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291300,114.588100",
        "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291300,114.588100",
        "waze_url": "https://www.waze.com/ul?ll=-3.291300,114.588100&navigate=yes",
        "osm_url": "https://www.openstreetmap.org/?mlat=-3.291300&mlon=114.588100&zoom=16",
        "is_outside_kua": true
      }
    ]
  }
}
```

> **Note:** Only shows assignments with status `Penghulu Ditugaskan`. Sorted by `tanggal_nikah ASC, waktu_nikah ASC`.

---

### PUT `/simnikah/penghulu/coordinates`

Update the home address coordinates for the currently logged-in Penghulu. These coordinates are used as the origin point for routing and distance calculations on Saturdays, Sundays, and holidays.

**Auth:** Required
**Roles:** `penghulu` only

**Request Body:**
```json
{
  "latitude": -3.2913,
  "longitude": 114.5881
}
```

**Request Fields:**

| Field       | Type     | Required | Description              |
|-------------|----------|----------|--------------------------|
| `latitude`  | *float64 | **Yes**  | Home latitude coordinate |
| `longitude` | *float64 | **Yes**  | Home longitude coordinate|

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Koordinat alamat penghulu berhasil diperbarui",
  "data": {
    "id": 1,
    "nama": "H. Abdul Rahman, S.Ag",
    "latitude": -3.2913,
    "longitude": 114.5881
  }
}
```

---

## 8. Location & Geocoding

All location endpoints use **OpenStreetMap Nominatim API** (100% free, no API key required).

### POST `/simnikah/location/geocode`

Convert address text to coordinates (Nominatim, with caching).

**Auth:** Required

**Request Body:**
```json
{
  "alamat": "Jl. Gatot Subroto No. 10, Banjarmasin"
}
```

**Validation:** Address must be at least 10 characters.

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Koordinat berhasil ditemukan",
  "data": {
    "alamat": "Jl. Gatot Subroto No. 10, Banjarmasin",
    "latitude": -3.2913,
    "longitude": 114.5881,
    "map_url": "https://www.google.com/maps?q=-3.291300,114.588100",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.291300&mlon=114.588100&zoom=16"
  }
}
```

---

### POST `/simnikah/location/reverse-geocode`

Convert coordinates to address text.

**Auth:** Required

**Request Body:**
```json
{
  "latitude": -3.2913,
  "longitude": 114.5881
}
```

**Validation:**
- `latitude`: must be between -90 and 90
- `longitude`: must be between -180 and 180

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Alamat berhasil ditemukan",
  "data": {
    "latitude": -3.2913,
    "longitude": 114.5881,
    "alamat": "Jl. Gatot Subroto No. 10, Banjarmasin Utara, Banjarmasin, ...",
    "detail": {
      "road": "Jl. Gatot Subroto",
      "suburb": "Banjarmasin Utara",
      "city_district": "Banjarmasin Utara",
      "city": "Banjarmasin",
      "state": "Kalimantan Selatan",
      "postcode": "70116",
      "country": "Indonesia"
    },
    "map_url": "https://www.google.com/maps?q=-3.291300,114.588100",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.291300&mlon=114.588100&zoom=16"
  }
}
```

---

### GET `/simnikah/location/search`

Search address with autocomplete (Nominatim, Indonesia only, max 5 results).

**Auth:** Required

**Query Parameters:**

| Parameter | Type   | Required | Description                |
|-----------|--------|----------|----------------------------|
| `q`       | string | **Yes**  | Min 3 characters           |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Hasil pencarian alamat",
  "data": {
    "query": "Gatot Subroto Banjarmasin",
    "results": [
      {
        "display_name": "Jl. Gatot Subroto, Banjarmasin Utara, ...",
        "latitude": "-3.2913",
        "longitude": "114.5881",
        "address": {
          "road": "Jl. Gatot Subroto",
          "suburb": "Banjarmasin Utara",
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

### PUT `/simnikah/pendaftaran/:id/location`

Update wedding location address and coordinates for a registration.

**Auth:** Required
**Access Control:** Only the registration owner (pendaftar), `staff`, or `kepala_kua` can update.

**Restriction:** Only works for registrations with `tempat_nikah = "Di Luar KUA"`.

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Request Body:**
```json
{
  "alamat_akad": "Jl. Merdeka No. 5, Banjarmasin",
  "latitude": -3.3200,
  "longitude": 114.5900
}
```

**Request Fields:**

| Field        | Type     | Required | Description                            |
|--------------|----------|----------|----------------------------------------|
| `alamat_akad`| string   | **Yes**  | New full address                        |
| `latitude`   | *float64 | No       | Latitude (auto-geocoded if not provided)|
| `longitude`  | *float64 | No       | Longitude (auto-geocoded if not provided)|

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Lokasi nikah berhasil diupdate",
  "data": {
    "pendaftaran_id": 42,
    "alamat_akad": "Jl. Merdeka No. 5, Banjarmasin",
    "tempat_nikah": "Di Luar KUA",
    "latitude": -3.3200,
    "longitude": 114.5900,
    "map_url": "https://www.google.com/maps?q=-3.320000,114.590000",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.320000&mlon=114.590000&zoom=16",
    "updated_at": "2025-01-20T10:40:00Z"
  }
}
```

---

### GET `/simnikah/pendaftaran/:id/location`

Get detailed wedding location with maps URLs for navigation.

**Auth:** Required
**Access Control:** Only the registration owner, the assigned penghulu, `staff`, or `kepala_kua` can view.

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Detail lokasi nikah berhasil diambil",
  "data": {
    "pendaftaran_id": 42,
    "nama_suami": "Ahmad Fauzi",
    "nama_istri": "Siti Aminah",
    "tanggal_nikah": "2025-03-15T00:00:00Z",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Gatot Subroto No. 10",
    "latitude": -3.2913,
    "longitude": 114.5881,
    "has_coordinates": true,
    "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291300,114.588100",
    "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291300,114.588100",
    "waze_url": "https://www.waze.com/ul?ll=-3.291300,114.588100&navigate=yes",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.291300&mlon=114.588100&zoom=16",
    "is_outside_kua": true,
    "note": "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
  }
}
```

---

## 9. Staff Management

### GET `/simnikah/staff`

List all staff members.

**Auth:** Required
**Roles:** `kepala_kua` only

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Data staff berhasil diambil",
  "data": [
    {
      "id": 1,
      "id_pengguna": "USR-s1Bf9xK2mNp1",
      "nip": "197501011990031002",
      "nama_lengkap": "Staff KUA",
      "jabatan": "Staff",
      "bagian": "Administrasi",
      "nomor_telepon": "081234567891",
      "email": "staff@kua.go.id",
      "alamat": "Banjarmasin",
      "status": "Aktif",
      "dibuat_pada": "2025-01-01T00:00:00Z",
      "diperbarui_pada": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

### PUT `/simnikah/staff/:id`

Update staff information. All fields are optional — only provided fields will be updated.

**Auth:** Required
**Roles:** `kepala_kua` only

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | uint | Staff ID    |

**Request Body:**
```json
{
  "nama": "Staff Updated",
  "jabatan": "Staff",
  "bagian": "Administrasi",
  "no_hp": "081234567890",
  "email": "staff@kua.go.id",
  "alamat": "Banjarmasin",
  "status": "Aktif"
}
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Staff berhasil diupdate",
  "data": { /* updated StaffKUA object */ }
}
```

---

### POST `/simnikah/staff/pendaftaran`

Create registration on behalf of a catin. Uses the same schema as catin register. Schedule availability is validated before creation via Forward Chaining Engine.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Request Body:** Same as [POST `/simnikah/pendaftaran`](#post-simnikahpendaftaran)

**Response `201 Created`:**
```json
{
  "success": true,
  "message": "Pendaftaran berhasil dibuat",
  "data": { /* full PendaftaranNikah object */ }
}
```

**Response `409 Conflict` (schedule unavailable):**
```json
{
  "success": false,
  "message": "Jadwal tidak tersedia",
  "error": "Semua slot penghulu pada jam tersebut sudah penuh"
}
```

---

### PUT `/simnikah/pendaftaran/:id/update-status`

Update registration status manually. Enforces valid state transitions.

**Auth:** Required
**Roles:** `staff`, `penghulu`, `kepala_kua`

**Path Parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | uint | Registration ID   |

**Request Body:**
```json
{
  "status_pendaftaran": "Selesai"
}
```

**Valid status values:** `Menunggu Penugasan`, `Penghulu Ditugaskan`, `Selesai`, `Ditolak`

**Valid transitions (enforced as state machine):**

| From                      | Allowed To                                |
|---------------------------|-------------------------------------------|
| `Menunggu Penugasan`      | `Penghulu Ditugaskan`, `Ditolak`          |
| `Penghulu Ditugaskan`     | `Selesai`                                 |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Status berhasil diupdate",
  "data": { /* updated PendaftaranNikah object */ }
}
```

**Response `409 Conflict` (invalid transition):**
```json
{
  "success": false,
  "message": "Transisi status tidak valid",
  "error": "Tidak dapat mengubah status dari 'Selesai' ke 'Menunggu Penugasan'"
}
```

---

## 10. Notification System

### POST `/simnikah/notifikasi`

Create a notification for a specific user.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "user_id": "USR-a3Bf9xK2mNp1",
  "judul": "Jadwal Diperbarui",
  "pesan": "Jadwal nikah Anda telah diperbarui",
  "tipe": "Info",
  "link": "/simnikah/pendaftaran/42"
}
```

**Request Fields:**

| Field     | Type   | Required | Description                                    |
|-----------|--------|----------|------------------------------------------------|
| `user_id` | string | **Yes**  | Target user's user_id                          |
| `judul`   | string | **Yes**  | Notification title                             |
| `pesan`   | string | **Yes**  | Notification message                           |
| `tipe`    | string | **Yes**  | `Info`, `Warning`, `Error`, `Success`          |
| `link`    | string | No       | Optional navigation link                       |

**Response `201 Created`:**
```json
{
  "message": "Notifikasi berhasil dibuat",
  "notification": {
    "id": 1,
    "user_id": "USR-a3Bf9xK2mNp1",
    "judul": "Jadwal Diperbarui",
    "pesan": "Jadwal nikah Anda telah diperbarui",
    "tipe": "Info",
    "status_baca": "Belum Dibaca",
    "link": "/simnikah/pendaftaran/42",
    "created_at": "2025-01-20T10:30:00Z",
    "updated_at": "2025-01-20T10:30:00Z"
  }
}
```

---

### GET `/simnikah/notifikasi/user/me`

Get paginated notifications for the currently authenticated user (user ID from JWT).

**Auth:** Required

**Query Parameters:**

| Parameter | Type | Default | Description                              |
|-----------|------|---------|------------------------------------------|
| `page`    | int  | `1`     | Page number                              |
| `limit`   | int  | `10`    | Items per page (max 100)                 |
| `status`  | str  | —       | `Belum Dibaca` or `Sudah Dibaca`         |
| `tipe`    | str  | —       | `Info`, `Warning`, `Error`, `Success`    |

**Response `200 OK`:**
```json
{
  "message": "Notifikasi berhasil diambil",
  "notifications": [
    {
      "id": 1,
      "user_id": "USR-a3Bf9xK2mNp1",
      "judul": "Jadwal Diperbarui",
      "pesan": "Jadwal nikah Anda telah diperbarui",
      "tipe": "Info",
      "status_baca": "Belum Dibaca",
      "link": "/simnikah/pendaftaran/42",
      "created_at": "2025-01-20T10:30:00Z",
      "updated_at": "2025-01-20T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_page": 5
  },
  "unread_count": 5
}
```

---

### GET `/simnikah/notifikasi/:id`

Get a single notification by ID. Only accessible by the notification owner.

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description        |
|-----------|------|--------------------|
| `id`      | uint | Notification ID    |

**Response `200 OK`:**
```json
{
  "message": "Notifikasi berhasil diambil",
  "notification": {
    "id": 1,
    "user_id": "USR-a3Bf9xK2mNp1",
    "judul": "Jadwal Diperbarui",
    "pesan": "Jadwal nikah Anda telah diperbarui",
    "tipe": "Info",
    "status_baca": "Belum Dibaca",
    "link": "/simnikah/pendaftaran/42",
    "created_at": "2025-01-20T10:30:00Z",
    "updated_at": "2025-01-20T10:30:00Z"
  }
}
```

---

### PUT `/simnikah/notifikasi/:id/status`

Mark notification as read/unread. Only accessible by the notification owner.

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description        |
|-----------|------|--------------------|
| `id`      | uint | Notification ID    |

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

**Valid values:** `Belum Dibaca`, `Sudah Dibaca`

**Response `200 OK`:**
```json
{
  "message": "Status notifikasi berhasil diupdate",
  "notification": { /* updated notification object */ }
}
```

---

### PUT `/simnikah/notifikasi/mark-all-read`

Mark all notifications as read for the currently authenticated user.

**Auth:** Required

**Response `200 OK`:**
```json
{
  "message": "Semua notifikasi berhasil ditandai sebagai sudah dibaca",
  "updated_count": 5
}
```

---

### DELETE `/simnikah/notifikasi/:id`

Delete a notification. Only accessible by the notification owner.

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description        |
|-----------|------|--------------------|
| `id`      | uint | Notification ID    |

**Response `200 OK`:**
```json
{
  "message": "Notifikasi berhasil dihapus"
}
```

---

### GET `/simnikah/notifikasi/stats`

Get notification statistics for the currently authenticated user.

**Auth:** Required

**Response `200 OK`:**
```json
{
  "message": "Statistik notifikasi berhasil diambil",
  "stats": {
    "total": 50,
    "unread": 5,
    "read": 45,
    "by_type": {
      "info": 20,
      "warning": 10,
      "error": 5,
      "success": 15
    },
    "today": 3,
    "week": 12
  }
}
```

---

### POST `/simnikah/notifikasi/send-to-role`

Send notification to all active users with a specific role.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "role": "penghulu",
  "judul": "Pengumuman",
  "pesan": "Rapat rutin hari ini pukul 14:00",
  "tipe": "Warning",
  "link": "/dashboard"
}
```

**Request Fields:**

| Field  | Type   | Required | Description                                                 |
|--------|--------|----------|-------------------------------------------------------------|
| `role` | string | **Yes**  | `user_biasa`, `penghulu`, `staff`, `kepala_kua`             |
| `judul`| string | **Yes**  | Notification title                                          |
| `pesan`| string | **Yes**  | Notification message                                        |
| `tipe` | string | **Yes**  | `Info`, `Warning`, `Error`, `Success`                       |
| `link` | string | No       | Optional navigation link                                    |

**Response `201 Created`:**
```json
{
  "message": "Notifikasi berhasil dikirim ke 5 user dengan role penghulu",
  "recipient_count": 5,
  "role": "penghulu"
}
```

---

### POST `/simnikah/notifikasi/run-reminder`

Manually trigger reminder notifications (for testing). In production, reminders are sent automatically via cron job at 08:00 WITA daily.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Pengingat notifikasi berhasil dijalankan",
  "executed_by": "USR-k7Hd2Xm9pQ4w",
  "executed_at": "2025-01-20T08:00:00Z"
}
```

---

## 11. Dashboard & Analytics

### GET `/simnikah/dashboard/kepala-kua`

Comprehensive dashboard for Kepala KUA with statistics, trends, status distribution, penghulu performance, and peak hours analysis.

**Auth:** Required
**Roles:** `kepala_kua` only

**Query Parameters:**

| Parameter    | Type   | Default | Description                              |
|--------------|--------|---------|------------------------------------------|
| `period`     | string | `month` | `day`, `week`, `month`, `year`           |
| `date_from`  | string | —       | Override start date (`YYYY-MM-DD`)       |
| `date_to`    | string | —       | Override end date (`YYYY-MM-DD`)         |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Dashboard data berhasil diambil",
  "data": {
    "period": {
      "type": "month",
      "date_from": "2025-01-01",
      "date_to": "2025-01-31"
    },
    "statistics": {
      "total_periode": 45,
      "hari_ini": 3,
      "bulan_ini": 45,
      "tahun_ini": 120,
      "selesai": 30,
      "pending": 15,
      "status_breakdown": {
        "menunggu_penugasan": 10,
        "penghulu_ditugaskan": 5,
        "selesai": 30,
        "ditolak": 0
      }
    },
    "trends": [
      { "date": "2025-01", "count": 45 }
    ],
    "status_distribution": [
      { "status": "Menunggu Penugasan", "count": 10, "label": "Menunggu Penugasan" },
      { "status": "Penghulu Ditugaskan", "count": 5, "label": "Penghulu Ditugaskan" },
      { "status": "Selesai", "count": 30, "label": "Selesai" },
      { "status": "Ditolak", "count": 0, "label": "Ditolak" }
    ],
    "penghulu_performance": [
      {
        "penghulu_id": 1,
        "nama_lengkap": "H. Abdul Rahman, S.Ag",
        "jumlah_nikah": 12,
        "rating": 4.8
      }
    ],
    "peak_hours": [
      { "waktu": "08:00", "count": 2 },
      { "waktu": "09:00", "count": 8 },
      { "waktu": "10:00", "count": 5 },
      { "waktu": "11:00", "count": 3 },
      { "waktu": "12:00", "count": 1 },
      { "waktu": "13:00", "count": 4 },
      { "waktu": "14:00", "count": 6 },
      { "waktu": "15:00", "count": 2 },
      { "waktu": "16:00", "count": 1 }
    ]
  }
}
```

---

### GET `/simnikah/dashboard/staff`

Staff dashboard showing pending assignments, recent registrations, and activity timeline.

**Auth:** Required
**Roles:** `staff` only

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Dashboard staff berhasil diambil",
  "data": {
    "pending_assignments": [
      {
        "id": 42,
        "nama_suami": "Ahmad Fauzi",
        "nama_istri": "Siti Aminah",
        "status_pendaftaran": "Menunggu Penugasan",
        "tanggal_nikah": "2025-03-15",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_akad": "Jl. Gatot Subroto No. 10",
        "created_at": "2025-01-20T10:30:00Z"
      }
    ],
    "recent_registrations": [ /* last 10 registrations */ ],
    "timeline": [
      {
        "id": 42,
        "nama_suami": "Ahmad Fauzi",
        "nama_istri": "Siti Aminah",
        "status_pendaftaran": "Penghulu Ditugaskan",
        "tanggal_nikah": "2025-03-15",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "updated_at": "2025-01-20T10:35:00Z",
        "action": "Status diubah menjadi Penghulu Ditugaskan"
      }
    ]
  }
}
```

---

### GET `/simnikah/dashboard/statistik-pernikahan`

Detailed marriage statistics with trends.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Query Parameters:**

| Parameter | Type   | Default | Description          |
|-----------|--------|---------|----------------------|
| `period`  | string | `month` | `day`, `month`, `year` |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Statistik pernikahan berhasil diambil",
  "data": {
    "statistics": { /* same as kepala-kua dashboard statistics */ },
    "trends": [ /* array of { date, count } */ ]
  }
}
```

---

### GET `/simnikah/dashboard/penghulu-performance`

Detailed penghulu performance metrics for a date range.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Query Parameters:**

| Parameter   | Type   | Default     | Description            |
|-------------|--------|-------------|------------------------|
| `date_from` | string | month start | Start date `YYYY-MM-DD`|
| `date_to`   | string | month end   | End date `YYYY-MM-DD`  |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Statistik penghulu berhasil diambil",
  "data": [
    {
      "penghulu_id": 1,
      "nama_lengkap": "H. Abdul Rahman, S.Ag",
      "jumlah_nikah": 12,
      "rating": 4.8
    }
  ]
}
```

---

### GET `/simnikah/dashboard/peak-hours`

Peak hours analysis showing busiest time slots.

**Auth:** Required
**Roles:** `staff`, `kepala_kua`

**Query Parameters:**

| Parameter   | Type   | Default     | Description            |
|-------------|--------|-------------|------------------------|
| `date_from` | string | month start | Start date `YYYY-MM-DD`|
| `date_to`   | string | month end   | End date `YYYY-MM-DD`  |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Analisis jam sibuk berhasil diambil",
  "data": [
    { "waktu": "08:00", "count": 2 },
    { "waktu": "09:00", "count": 8 },
    { "waktu": "10:00", "count": 5 },
    { "waktu": "11:00", "count": 3 },
    { "waktu": "12:00", "count": 1 },
    { "waktu": "13:00", "count": 4 },
    { "waktu": "14:00", "count": 6 },
    { "waktu": "15:00", "count": 2 },
    { "waktu": "16:00", "count": 1 }
  ]
}
```

---

## 12. Forward Chaining Rules Reference

The FC Engine evaluates **10 rules** in sequence per penghulu candidate:

| Rule ID  | Rule Name              | Blocking | Decision Variable Used          | Description                                     |
|----------|------------------------|----------|---------------------------------|-------------------------------------------------|
| RULE_001 | Validasi Administrasi  | **Yes**  | `status_pendaftaran`            | Status must be `Menunggu Penugasan`             |
| RULE_002 | Validasi Status Penghulu| **Yes** | —                               | Penghulu must be `Aktif`                        |
| RULE_003 | Konteks Hari Libur     | No       | `tanggal_nikah`                 | Context: is it a holiday?                       |
| RULE_004 | Cek Konflik Jadwal     | **Yes**  | `tanggal_nikah`, `waktu_nikah`  | No schedule conflict at same date/time          |
| RULE_005 | Cek Kapasitas Harian   | **Yes**  | `tanggal_nikah`                 | Daily capacity not exceeded (default: 3/day)    |
| RULE_006 | Cek Kapasitas Per Jam  | **Yes**  | `tanggal_nikah`, `waktu_nikah`  | Hourly capacity not exceeded (default: 1/hour)  |
| RULE_007 | Cek Kesesuaian Lokasi  | **Yes**  | `tempat_nikah`                  | All penghulu can serve inside/outside KUA       |
| RULE_008 | Cek Batas Rating       | **Yes**  | — (penghulu.rating)             | Rating >= minimum (default: 3.0)                |
| RULE_009 | Estimasi Jarak         | No       | `latitude`, `longitude`         | OSRM route distance + Haversine fallback        |
| RULE_010 | Konklusi Akhir         | No       | —                               | All facts satisfied → recommend                 |

### RULE_009 Distance Estimation Note

- **Hari kerja (weekday):** Origin point is always KUA office coordinates.
- **Hari libur nasional / akhir pekan:** If the penghulu has home coordinates (`latitude` & `longitude` not null), the origin uses the penghulu's home coordinates. Falls back to KUA coordinates if not set.
- **Method:** OSRM (Open Source Routing Machine) for actual route distance, with Haversine formula as fallback.

### Scoring Formula

```
Score = (RatingWeight × RatingScore)
      + (AvailabilityWeight × AvailabilityScore)
      + (FairnessWeight × FairnessScore)
      + (LocationMatchWeight × LocationScore)
      + (DistanceWeight × DistanceScore)
```

| Component          | Weight | Calculation                              |
|--------------------|--------|------------------------------------------|
| Rating             | 0.35   | `(rating/5) × 100`                       |
| Availability       | 0.25   | `(1 - dayBooked/capacityPerDay) × 100`   |
| Fairness           | 0.20   | Workload balance vs min/max penghulu     |
| Location Match     | 0.10   | 100 for Di KUA, 85 for Di Luar KUA       |
| Distance           | 0.10   | `max(0, 100 - km × penalty)` (penalty=8) |

### Inference Flow (Forward Chaining)

```
Initial Facts (Working Memory)           Rules                              Derived Facts
┌─────────────────────────┐               ┌──────────────────────┐           ┌────────────────────────┐
│ • tanggal_nikah         │──> RULE_001 ──>│ Lulus Administrasi  │           │ ✓ Slot tersedia        │
│ • waktu_nikah           │──> RULE_002 ──>│ Penghulu Aktif      │           │ ✓ Penghulu aktif       │
│ • tempat_nikah          │──> RULE_003 ──>│ Konteks Hari Libur  │           │ ✓ Kapasitas harian OK  │
│ • alamat_akad           │──> RULE_004 ──>│ Tidak Bentrok       │──────────>│ ✓ Kapasitas jam OK     │
│ • latitude / longitude  │──> RULE_005 ──>│ Kapasitas Harian OK │           │ ✓ Lokasi sesuai        │
│ • status_pendaftaran    │──> RULE_006 ──>│ Kapasitas Jam OK    │           │ ✓ Jarak layak          │
│ • penghulu status/data  │──> RULE_007 ──>│ Lokasi Sesuai       │           │ ✓ DIREKOMENDASIKAN     │
│ • beban kerja jadwal    │──> RULE_008 ──>│ Rating Memadai      │           └────────────────────────┘
│                         │──> RULE_009 ──>│ Jarak Optimal       │
│                         │──> RULE_010 ──>│ Konklusi Akhir      │
└─────────────────────────┘               └──────────────────────┘
```

> If any **blocking** rule fails, the engine immediately stops evaluation for that penghulu candidate and provides a clear rejection reason.

---

## 13. Error Handling & Status Codes

All errors follow a consistent JSON structure:

```json
{
  "success": false,
  "message": "Human-readable error description",
  "error": "Technical error detail"
}
```

### Common HTTP Status Codes

| Status | Meaning                    | When                                          |
|--------|----------------------------|-----------------------------------------------|
| `200`  | OK                         | Successful GET/PUT                            |
| `201`  | Created                    | Successful POST (create)                      |
| `400`  | Bad Request                | Invalid input, validation failure             |
| `401`  | Unauthorized               | Missing or invalid/expired JWT token          |
| `403`  | Forbidden                  | Role not permitted for this action            |
| `404`  | Not Found                  | Resource doesn't exist                        |
| `405`  | Method Not Allowed         | Wrong HTTP method                             |
| `409`  | Conflict                   | Slot full, status mismatch, invalid transition|
| `429`  | Too Many Requests          | Rate limit exceeded                           |
| `500`  | Internal Server Error      | Database or server failure                    |

### Custom 404 / 405 JSON Responses

The API returns JSON (not HTML) for unknown routes:

**404 Not Found:**
```json
{
  "success": false,
  "message": "Endpoint tidak ditemukan",
  "error": "Path '/unknown/path' tidak ditemukan"
}
```

**405 Method Not Allowed:**
```json
{
  "success": false,
  "message": "Method tidak diizinkan",
  "error": "Method 'DELETE' tidak diizinkan untuk path ini"
}
```

---

## 14. Rate Limiting

### Global Rate Limit

- **Rate:** 100 requests/min per IP
- **Applied to:** All endpoints (after CORS middleware)

### Strict Rate Limit

- **Rate:** 5 requests/min per IP
- **Applied to:** `/login`, `/register`

### Rate Limit Headers

All responses include:

| Header                  | Description                          |
|-------------------------|--------------------------------------|
| `X-RateLimit-Limit`     | Maximum requests allowed per window  |
| `X-RateLimit-Remaining` | Remaining requests in current window |
| `X-RateLimit-Reset`     | Unix timestamp when window resets    |

### Rate Limit Exceeded Response `429`

**Global:**
```json
{
  "success": false,
  "message": "Rate limit exceeded",
  "error": "Terlalu banyak request. Silakan coba lagi nanti.",
  "retry_after": "45 detik"
}
```

**Strict (Auth endpoints):**
```json
{
  "success": false,
  "message": "Terlalu banyak percobaan login",
  "error": "Demi keamanan, Anda harus menunggu sebelum mencoba lagi.",
  "retry_after": "30 detik",
  "type": "security"
}
```

---

## 15. Environment Variables & Configuration

### Required Environment Variables

| Variable        | Default        | Description                              |
|-----------------|----------------|------------------------------------------|
| `DB_USER`       | `root`         | MySQL username                           |
| `DB_PASSWORD`   | —              | MySQL password (required for non-localhost)|
| `DB_HOST`       | `127.0.0.1`    | MySQL host                               |
| `DB_PORT`       | `3306`         | MySQL port                               |
| `DB_NAME`       | `simnikah`     | MySQL database name                      |
| `PORT`          | `8080`         | Server port                              |
| `GIN_MODE`      | —              | `release` for production                 |
| `ENVIRONMENT`   | —              | `production` for production              |
| `JWT_SECRET`    | —              | JWT signing secret key                   |
| `IMGBB_API_KEY` | —              | ImgBB API key for photo uploads          |
| `ALLOWED_ORIGINS` | (defaults)  | Comma-separated CORS origins             |

### Default CORS Origins (when `ALLOWED_ORIGINS` not set)

```
http://localhost:3000
http://localhost:3001
http://localhost:5173
http://localhost:5174
http://localhost:8080
http://127.0.0.1:3000
http://127.0.0.1:5173
https://kua-ku.vercel.app
```

### Database Connection Pool

| Setting             | Value    |
|---------------------|----------|
| Max Idle Connections| 10       |
| Max Open Connections| 100      |
| Connection Max Life | 1 hour   |

### Timezone

All date/time operations use **WITA (Waktu Indonesia Tengah / UTC+8)**.

### Cron Job

| Schedule      | Description                                  |
|---------------|----------------------------------------------|
| Daily 08:00   | Send reminder notifications for upcoming weddings |

---

## 16. Health Check

### GET `/health`

No auth required. Returns server health status.

**Response `200 OK`:**
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2025-01-20T10:30:00Z"
}
```

---

## Endpoint Quick Reference

### Auth (No prefix)

| Method | Path            | Auth     | Roles     | Description                    |
|--------|-----------------|----------|-----------|--------------------------------|
| POST   | `/register`     | None     | —         | Register new user (user_biasa) |
| POST   | `/login`        | None     | —         | Login & get JWT token          |
| GET    | `/profile`      | Required | All       | Get current user profile       |
| POST   | `/upload-photo` | Required | All       | Upload profile photo           |
| GET    | `/health`       | None     | —         | Health check                   |

### SimNikah — Registration

| Method | Path                                    | Auth     | Roles                          | Description                         |
|--------|-----------------------------------------|----------|--------------------------------|-------------------------------------|
| POST   | `/simnikah/check-schedule`              | Required | All                            | Check schedule availability         |
| POST   | `/simnikah/pendaftaran`                 | Required | All                            | Create registration (catin)         |
| GET    | `/simnikah/pendaftaran/status`          | Required | All                            | Get own registration status         |
| GET    | `/simnikah/pendaftaran/:id`             | Required | All (own) / staff, kepala_kua  | Get registration detail             |
| GET    | `/simnikah/pendaftaran`                 | Required | staff, kepala_kua              | List all registrations (paginated)  |

### SimNikah — Forward Chaining (Kepala KUA)

| Method | Path                                                            | Auth     | Roles       | Description                         |
|--------|-----------------------------------------------------------------|----------|-------------|-------------------------------------|
| GET    | `/simnikah/kepala-kua/forward-chaining/recommendation/:id`     | Required | kepala_kua  | Get FC recommendation               |
| GET    | `/simnikah/kepala-kua/forward-chaining/evaluation/:id`         | Required | kepala_kua  | Get detailed evaluation report      |
| GET    | `/simnikah/kepala-kua/forward-chaining/config`                 | Required | kepala_kua  | Get FC engine configuration         |
| POST   | `/simnikah/kepala-kua/forward-chaining/assign/:id`             | Required | kepala_kua  | Assign penghulu (transactional)     |
| GET    | `/simnikah/kepala-kua/available-penghulu`                      | Required | kepala_kua  | List active penghulu                |
| GET    | `/simnikah/kepala-kua/penghulu-tersedia`                       | Required | kepala_kua  | Penghulu schedule by date           |

### SimNikah — Penghulu

| Method | Path                                      | Auth     | Roles               | Description                    |
|--------|--------------------------------------------|----------|----------------------|--------------------------------|
| GET    | `/simnikah/penghulu/jadwal-penugasan`     | Required | penghulu, kepala_kua | View assigned tasks            |
| PUT    | `/simnikah/penghulu/coordinates`          | Required | penghulu             | Update home coordinates        |

### SimNikah — Staff

| Method | Path                                       | Auth     | Roles               | Description                    |
|--------|---------------------------------------------|----------|----------------------|--------------------------------|
| GET    | `/simnikah/staff`                          | Required | kepala_kua           | List all staff                 |
| PUT    | `/simnikah/staff/:id`                      | Required | kepala_kua           | Update staff info              |
| POST   | `/simnikah/staff/pendaftaran`              | Required | staff, kepala_kua    | Create registration for user   |
| PUT    | `/simnikah/pendaftaran/:id/update-status`  | Required | staff, penghulu, kepala_kua | Update registration status |

### SimNikah — Location

| Method | Path                                       | Auth     | Roles     | Description                    |
|--------|---------------------------------------------|----------|-----------|--------------------------------|
| POST   | `/simnikah/location/geocode`               | Required | All       | Address → Coordinates          |
| POST   | `/simnikah/location/reverse-geocode`       | Required | All       | Coordinates → Address          |
| GET    | `/simnikah/location/search`                | Required | All       | Address search (autocomplete)  |
| PUT    | `/simnikah/pendaftaran/:id/location`       | Required | Owner/Staff/Kepala | Update wedding location  |
| GET    | `/simnikah/pendaftaran/:id/location`       | Required | Owner/Penghulu/Staff/Kepala | Get wedding location detail |

### SimNikah — Notifications

| Method | Path                                       | Auth     | Roles               | Description                    |
|--------|---------------------------------------------|----------|----------------------|--------------------------------|
| POST   | `/simnikah/notifikasi`                     | Required | staff, kepala_kua    | Create notification            |
| GET    | `/simnikah/notifikasi/user/me`             | Required | All                  | Get own notifications          |
| GET    | `/simnikah/notifikasi/:id`                 | Required | Owner                | Get single notification        |
| PUT    | `/simnikah/notifikasi/:id/status`          | Required | Owner                | Update read/unread status      |
| PUT    | `/simnikah/notifikasi/mark-all-read`       | Required | All                  | Mark all as read               |
| DELETE | `/simnikah/notifikasi/:id`                 | Required | Owner                | Delete notification            |
| GET    | `/simnikah/notifikasi/stats`               | Required | All                  | Get notification statistics    |
| POST   | `/simnikah/notifikasi/send-to-role`        | Required | staff, kepala_kua    | Send to all users by role      |
| POST   | `/simnikah/notifikasi/run-reminder`        | Required | staff, kepala_kua    | Manual reminder trigger        |

### SimNikah — Dashboard

| Method | Path                                           | Auth     | Roles               | Description                    |
|--------|--------------------------------------------------|----------|----------------------|--------------------------------|
| GET    | `/simnikah/dashboard/kepala-kua`                | Required | kepala_kua           | Full Kepala KUA dashboard      |
| GET    | `/simnikah/dashboard/staff`                     | Required | staff                | Staff dashboard                |
| GET    | `/simnikah/dashboard/statistik-pernikahan`      | Required | staff, kepala_kua    | Marriage statistics            |
| GET    | `/simnikah/dashboard/penghulu-performance`      | Required | staff, kepala_kua    | Penghulu performance metrics   |
| GET    | `/simnikah/dashboard/peak-hours`                | Required | staff, kepala_kua    | Peak hours analysis            |

---

**Total Endpoints: 35**

| Category                | Count |
|-------------------------|-------|
| Auth & Profile          | 5     |
| Registration (Catin)    | 5     |
| Forward Chaining        | 4     |
| Kepala KUA Management   | 2     |
| Penghulu                | 2     |
| Staff Management        | 4     |
| Location & Geocoding    | 5     |
| Notifications           | 9     |
| Dashboard & Analytics   | 5     |
