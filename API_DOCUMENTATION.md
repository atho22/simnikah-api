# SIPENA API Documentation

> **SIPENA** — Sistem Pendukung Keputusan untuk distribusi beban kerja penghulu menggunakan metode Forward Chaining.
>
> Bukan aplikasi pendaftaran nikah yang memvalidasi berkas (seperti SIMKAH), melainkan murni SPK untuk penjadwalan dan optimasi penugasan penghulu.

**Base URL:** `http://localhost:8080`

**Authentication:** Bearer Token (JWT), dikirim via header `Authorization: Bearer <token>`

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Data Models](#2-data-models)
3. [Authentication](#3-authentication)
4. [Stage 1 — Catin: Check Schedule & Register](#4-stage-1--catin-check-schedule--register)
5. [Stage 2 — Kepala KUA: Forward Chaining Recommendation](#5-stage-2--kepala-kua-forward-chaining-recommendation)
6. [Stage 3 — Kepala KUA: Assignment Approval](#6-stage-3--kepala-kua-assignment-approval)
7. [Stage 4 — Penghulu: View Assignments](#7-stage-4--penghulu-view-assignments)
8. [Location & Geocoding](#8-location--geocoding)
9. [Staff Management](#9-staff-management)
10. [Notification](#10-notification)
11. [Dashboard & Analytics](#11-dashboard--analytics)
12. [Forward Chaining Rules Reference](#12-forward-chaining-rules-reference)
13. [Error Handling](#13-error-handling)

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
| `user_biasa`   | Calon Pengantin (Catin)            | Check schedule, Create registration                  |
| `kepala_kua`   | Kepala KUA — decision maker        | FC recommendation, Assign penghulu, Dashboard        |
| `staff`        | Staff KUA — assists operations     | Create registration, Update status, Notifications    |
| `penghulu`     | Penghulu — receives assignments    | View jadwal penugasan with geolocation               |

### Status Flow

```
Menunggu Penugasan  ──>  Penghulu Ditugaskan  ──>  Selesai
       │                        │
       └──────> Ditolak  <──────┘
```

> **Note:** Draft and Disetujui statuses have been removed. Registrations go directly to `Menunggu Penugasan` upon creation.

---

## 2. Data Models

### PendaftaranNikah

| Field                | Type       | Required | FC Variable | Description                                    |
|----------------------|------------|----------|-------------|------------------------------------------------|
| `id`                 | uint       | auto     | —           | Primary key                                    |
| `pendaftar_id`       | string(20) | auto     | —           | User ID of the registrant (from JWT)           |
| `nama_suami`         | string(100)| **Yes**  | **No**      | Calon suami name (Excel reference / display)   |
| `umur_suami`         | int        | No       | **No**      | Calon suami age (Excel reference / display)    |
| `nama_istri`         | string(100)| **Yes**  | **No**      | Calon istri name (Excel reference / display)   |
| `umur_istri`         | int        | No       | **No**      | Calon istri age (Excel reference / display)    |
| `tanggal_nikah`      | date       | **Yes**  | **Yes**     | Wedding date (`YYYY-MM-DD`)                    |
| `waktu_nikah`        | string(10) | **Yes**  | **Yes**     | Wedding time (`HH:MM`)                         |
| `tempat_nikah`       | string(100)| **Yes**  | **Yes**     | `"Di KUA"` or `"Di Luar KUA"`                  |
| `alamat_akad`        | string(200)| Cond.    | **Yes**     | Full address (required if Di Luar KUA)         |
| `latitude`           | *float64   | No       | **Yes**     | Latitude (auto-geocoded if empty)              |
| `longitude`          | *float64   | No       | **Yes**     | Longitude (auto-geocoded if empty)             |
| `status_pendaftaran` | string(40) | auto     | —           | Current status                                 |
| `penghulu_id`        | *uint      | auto     | —           | Assigned penghulu ID                           |
| `created_at`         | timestamp  | auto     | —           |                                                |
| `updated_at`         | timestamp  | auto     | —           |                                                |

> **Critical:** Fields marked **FC Variable = Yes** are the ONLY inputs to Forward Chaining rules. Couple data (nama/umur) is display-only and never enters the inference engine.

### Users

| Field           | Type       | Description            |
|-----------------|------------|------------------------|
| `user_id`       | string(20) | Unique user ID         |
| `username`      | string(50) | Login username         |
| `email`         | string(100)| Email                  |
| `role`          | string(20) | `user_biasa`, `penghulu`, `staff`, `kepala_kua` |
| `status`        | string(20) | `Aktif`, `Nonaktif`, `Blokir` |
| `nama`          | string(100)| Full name              |
| `profile_photo` | string(500)| Photo URL (ImgBB)      |

### Penghulu

| Field           | Type       | Description                     |
|-----------------|------------|---------------------------------|
| `id`            | uint       | Primary key                     |
| `user_id`       | string(20) | Linked user ID                  |
| `nip`           | string(30) | NIP                             |
| `nama_lengkap`  | string(100)| Full name                       |
| `latitude`      | float64    | Latitude of home address        |
| `longitude`     | float64    | Longitude of home address       |
| `status`        | string(20) | `Aktif`, `Nonaktif`             |
| `jumlah_nikah`  | int        | Historical marriage count       |
| `rating`        | float64    | Performance rating (0-5)        |

---

## 3. Authentication

### POST `/register`

Register a new user account.

**Auth:** None  
**Rate Limit:** Strict (5 req/min)

**Request Body:**
```json
{
  "username": "ahmad_fauzi",
  "email": "ahmad@example.com",
  "password": "secret123",
  "nama": "Ahmad Fauzi",
  "role": "user_biasa"
}
```

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

---

### POST `/login`

Authenticate and receive JWT token.

**Auth:** None  
**Rate Limit:** Strict (5 req/min)

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
  }
}
```

> Token expires in 24 hours. Include in all subsequent requests: `Authorization: Bearer <token>`

---

### GET `/profile`

Get current user's profile.

**Auth:** Required

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
    "profile_photo": "https://ibb.co/..."
  }
}
```

---

### POST `/upload-photo`

Upload profile photo (ImgBB).

**Auth:** Required  
**Content-Type:** `multipart/form-data`

**Form Data:**
| Field  | Type | Description           |
|--------|------|-----------------------|
| `photo`| file | JPG/PNG/WebP, max 5MB |

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

---

## 4. Stage 1 — Catin: Check Schedule & Register

### POST `/simnikah/check-schedule`

Check if a schedule slot is available before registering. Delegates to Forward Chaining Engine.

**Auth:** Required  
**Roles:** All

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

Create a new marriage registration with couple + scheduling data.

**Auth:** Required  
**Roles:** All

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
- `nama_suami` and `nama_istri`: required
- `tanggal_nikah`: `YYYY-MM-DD`, cannot be in the past
- `waktu_nikah`: `HH:MM` format
- `tempat_nikah`: must be `"Di KUA"` or `"Di Luar KUA"`
- `alamat_akad`: required if `tempat_nikah` is `"Di Luar KUA"`
- `latitude`/`longitude`: optional, auto-geocoded from `alamat_akad` if empty

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

## 5. Stage 2 — Kepala KUA: Forward Chaining Recommendation

### GET `/simnikah/kepala-kua/forward-chaining/recommendation/:id`

Trigger Forward Chaining Engine to get penghulu recommendation for a specific registration.

**Auth:** Required  
**Roles:** `kepala_kua` only

**Path Parameters:**
| Parameter | Type | Description             |
|-----------|------|-------------------------|
| `id`      | uint | Registration ID         |

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
      "Flow ini hanya memakai constraint jadwal...",
      "Sumber dinamis yang direkomendasikan: tabel system_configs"
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
| Parameter | Type | Description             |
|-----------|------|-------------------------|
| `id`      | uint | Registration ID         |

**Request Body:**
```json
{
  "penghulu_id": 1,
  "approval_notes": "Disetujui, penghulu tersedia dan lokasi dekat"
}
```

**Business Logic:**
1. Validate role = `kepala_kua`
2. Begin GORM transaction
3. `SELECT ... FOR UPDATE` on the registration row
4. Verify status is `Menunggu Penugasan` (reject if not)
5. Verify penghulu is active
6. Update: `penghulu_id`, `status_pendaftaran` → `Penghulu Ditugaskan`
7. Commit transaction

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

**Response `409 Conflict` (already assigned):**
```json
{
  "success": false,
  "message": "Status pendaftaran tidak dapat di-assign",
  "error": "status saat ini: Penghulu Ditugaskan, harus 'Menunggu Penugasan'"
}
```

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

View penghulu schedule availability for a specific date.

**Auth:** Required  
**Roles:** `kepala_kua` only

**Query Parameters:**
| Parameter  | Type   | Required | Description        |
|------------|--------|----------|--------------------|
| `tanggal`  | string | **Yes**  | `YYYY-MM-DD`       |

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

### POST `/simnikah/location/geocode`

Convert address text to coordinates (Nominatim, cached).

**Auth:** Required

**Request Body:**
```json
{
  "alamat": "Jl. Gatot Subroto No. 10, Banjarmasin"
}
```

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

---

### GET `/simnikah/location/search`

Search address with autocomplete (Nominatim, Indonesia only).

**Auth:** Required

**Query Parameters:**
| Parameter | Type   | Required | Description              |
|-----------|--------|----------|--------------------------|
| `q`       | string | **Yes**  | Min 3 characters         |

---

### PUT `/simnikah/pendaftaran/:id/location`

Update wedding location address and coordinates.

**Auth:** Required

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

> Only works for registrations with `tempat_nikah = "Di Luar KUA"`. Coordinates are auto-geocoded if not provided.
>
> **Access Control:** Only the registration owner (pendaftar), staff, or kepala_kua can update the location.

---

### GET `/simnikah/pendaftaran/:id/location`

Get detailed wedding location with maps URLs for navigation.

**Auth:** Required

> **Access Control:** Only the registration owner, the assigned penghulu, staff, or kepala_kua can view the location detail.

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
    "google_maps_url": "https://www.google.com/maps/search/?api=1&query=...",
    "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=...",
    "waze_url": "https://www.waze.com/ul?ll=...&navigate=yes",
    "osm_url": "https://www.openstreetmap.org/?mlat=...&mlon=...&zoom=16",
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

---

### PUT `/simnikah/staff/:id`

Update staff information.

**Auth:** Required  
**Roles:** `kepala_kua` only

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

---

### POST `/simnikah/staff/pendaftaran`

Create registration on behalf of a catin (same schema as Catin register). Schedule availability is validated before creation.

**Auth:** Required  
**Roles:** `staff`, `kepala_kua`

**Request Body:** Same as [POST `/simnikah/pendaftaran`](#post-simnikahpendaftaran)

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

Update registration status manually.

**Auth:** Required  
**Roles:** `staff`, `penghulu`, `kepala_kua`

**Request Body:**
```json
{
  "status_pendaftaran": "Selesai"
}
```

**Valid statuses:** `Menunggu Penugasan`, `Penghulu Ditugaskan`, `Selesai`, `Ditolak`

---

## 10. Notification

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

**Valid `tipe` values:** `Info`, `Warning`, `Error`, `Success`

---

### GET `/simnikah/notifikasi/user/me`

Get paginated notifications for the currently authenticated user (user ID from JWT).

**Auth:** Required

**Query Parameters:**
| Parameter | Type | Default | Description                         |
|-----------|------|---------|-------------------------------------|
| `page`    | int  | 1       | Page number                         |
| `limit`   | int  | 10      | Items per page (max 100)            |
| `status`  | str  | —       | `Belum Dibaca` or `Sudah Dibaca`    |
| `tipe`    | str  | —       | `Info`, `Warning`, `Error`, `Success` |

**Response includes `unread_count` and `pagination` metadata.**

---

### GET `/simnikah/notifikasi/:id`

Get a single notification by ID. Only accessible by the notification owner.

**Auth:** Required

---

### PUT `/simnikah/notifikasi/:id/status`

Mark notification as read/unread. Only accessible by the notification owner.

**Auth:** Required

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

---

### PUT `/simnikah/notifikasi/mark-all-read`

Mark all notifications as read for the currently authenticated user.

**Auth:** Required

---

### DELETE `/simnikah/notifikasi/:id`

Delete a notification. Only accessible by the notification owner.

**Auth:** Required

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

Send notification to all users with a specific role.

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

---

### POST `/simnikah/notifikasi/run-reminder`

Manually trigger reminder notifications (for testing).

**Auth:** Required  
**Roles:** `staff`, `kepala_kua`

---

## 11. Dashboard & Analytics

### GET `/simnikah/dashboard/kepala-kua`

Comprehensive dashboard for Kepala KUA.

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
    "period": { "type": "month", "date_from": "2025-01-01", "date_to": "2025-01-31" },
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
      { "date": "2025-01", "label": "Jan 2025", "count": 45 }
    ],
    "status_distribution": [
      { "status": "Menunggu Penugasan", "count": 10, "label": "Menunggu Penugasan" }
    ],
    "penghulu_performance": [
      { "penghulu_id": 1, "nama_lengkap": "H. Abdul Rahman", "jumlah_nikah": 12, "rating": 4.8 }
    ],
    "peak_hours": [
      { "waktu": "08:00", "count": 2 },
      { "waktu": "09:00", "count": 8 }
    ]
  }
}
```

---

### GET `/simnikah/dashboard/staff`

Staff dashboard showing pending assignments and recent activity.

**Auth:** Required  
**Roles:** `staff` only

---

### GET `/simnikah/dashboard/statistik-pernikahan`

Detailed marriage statistics.

**Auth:** Required  
**Roles:** `staff`, `kepala_kua`

**Query:** `?period=day|month|year`

---

### GET `/simnikah/dashboard/penghulu-performance`

Detailed penghulu performance metrics.

**Auth:** Required  
**Roles:** `staff`, `kepala_kua`

**Query:** `?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD`

---

### GET `/simnikah/dashboard/peak-hours`

Peak hours analysis showing busiest time slots.

**Auth:** Required  
**Roles:** `staff`, `kepala_kua`

**Query:** `?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD`

---

## 12. Forward Chaining Rules Reference

The FC Engine evaluates 10 rules in sequence per penghulu candidate:

| Rule ID  | Rule Name              | Blocking | Decision Variable Used         | Description                                     |
|----------|------------------------|----------|--------------------------------|-------------------------------------------------|
| RULE_001 | Validasi Administrasi  | **Yes**  | `status_pendaftaran`           | Status must be `Menunggu Penugasan`             |
| RULE_002 | Validasi Status Penghulu| **Yes** | —                              | Penghulu must be `Aktif`                        |
| RULE_003 | Konteks Hari Libur     | No       | `tanggal_nikah`                | Context: is it a holiday?                       |
| RULE_004 | Cek Konflik Jadwal     | **Yes**  | `tanggal_nikah`, `waktu_nikah` | No schedule conflict at same date/time          |
| RULE_005 | Cek Kapasitas Harian   | **Yes**  | `tanggal_nikah`                | Daily capacity not exceeded (default: 3/day)    |
| RULE_006 | Cek Kapasitas Per Jam  | **Yes**  | `tanggal_nikah`, `waktu_nikah` | Hourly capacity not exceeded (default: 1/hour)  |
| RULE_007 | Cek Kesesuaian Lokasi  | **Yes**  | `tempat_nikah`                 | All penghulu can serve inside/outside KUA       |
| RULE_008 | Cek Batas Rating       | **Yes**  | — (penghulu.rating)            | Rating >= minimum (default: 3.0)                |
| RULE_009 | Estimasi Jarak         | No       | `latitude`, `longitude`        | OSRM route distance + Haversine fallback [1]    |
| RULE_010 | Konklusi Akhir         | No       | —                              | All facts satisfied → recommend                 |

**[1] Catatan untuk RULE_009:** Pada hari kerja, titik awal (`origin`) estimasi jarak selalu menggunakan koordinat kantor KUA. Pada hari libur nasional atau akhir pekan, jika penghulu yang dievaluasi memiliki koordinat rumah (`latitude` & `longitude` tidak null), titik awal estimasi rute akan menggunakan koordinat rumah penghulu tersebut (fallback ke KUA jika tidak diatur).

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

---

## 13. Error Handling

All errors follow a consistent JSON structure:

```json
{
  "success": false,
  "message": "Human-readable error description",
  "error": "Technical error detail"
}
```

### Common HTTP Status Codes

| Status | Meaning                          | When                                          |
|--------|----------------------------------|-----------------------------------------------|
| `200`  | OK                               | Successful GET/PUT                            |
| `201`  | Created                          | Successful POST (create)                      |
| `400`  | Bad Request                      | Invalid input, validation failure             |
| `401`  | Unauthorized                     | Missing or invalid JWT token                  |
| `403`  | Forbidden                        | Role not permitted for this action            |
| `404`  | Not Found                        | Resource doesn't exist                        |
| `405`  | Method Not Allowed               | Wrong HTTP method                             |
| `409`  | Conflict                         | Slot full, status mismatch                    |
| `429`  | Too Many Requests                | Rate limit exceeded                           |
| `500`  | Internal Server Error            | Database or server failure                    |

### Rate Limiting

- **Global:** 100 requests/min per IP
- **Auth endpoints** (`/login`, `/register`): 5 requests/min per IP (strict)

---

## Health Check

### GET `/health`

No auth required. Returns server health status.

```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2025-01-20T10:30:00Z"
}
```
