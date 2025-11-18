# 📚 SimNikah API Documentation

**Versi:** 1.0.0  
**Update Terakhir:** 2024  
**Base URL:** `http://localhost:8080` (Development) atau `https://your-domain.com` (Production)

---

## 📋 Daftar Isi

1. [Authentication](#authentication)
2. [Endpoints Overview](#endpoints-overview)
3. [Authentication Endpoints](#authentication-endpoints)
4. [Catin Endpoints](#catin-endpoints)
5. [Staff Endpoints](#staff-endpoints)
6. [Penghulu Endpoints](#penghulu-endpoints)
7. [Kepala KUA Endpoints](#kepala-kua-endpoints)
8. [Location Endpoints](#location-endpoints)
9. [Notification Endpoints](#notification-endpoints)
10. [Error Handling](#error-handling)
11. [Status Flow](#status-flow)

---

## 🔐 Authentication

Semua endpoint (kecuali `/register`, `/login`, dan endpoint kalender publik) memerlukan JWT token di header:

```
Authorization: Bearer <jwt_token>
```

**Token Format:**
- Type: JWT (JSON Web Token)
- Validity: 24 jam
- Algorithm: HS256

**Cara Mendapatkan Token:**
1. Register user baru: `POST /register`
2. Login: `POST /login`
3. Token akan dikembalikan dalam response `login`

---

## 📊 Endpoints Overview

| Category | Endpoints | Method | Auth Required | Role Required |
|----------|-----------|--------|---------------|---------------|
| **Authentication** | `/register`, `/login`, `/profile` | POST, POST, GET | No, No, Yes | - |
| **Catin** | `/simnikah/pendaftaran` | POST, GET | Yes | `user_biasa` |
| **Calendar** | `/simnikah/kalender-ketersediaan` | GET | No | - |
| **Staff** | `/simnikah/staff/*` | GET, POST, PUT | Yes | `staff`, `kepala_kua` |
| **Penghulu** | `/simnikah/penghulu/*` | GET, POST | Yes | `penghulu` |
| **Kepala KUA** | `/simnikah/kepala-kua/*` | GET, POST, PUT | Yes | `kepala_kua` |
| **Feedback** | `/simnikah/feedback-pernikahan` | POST, GET | Yes | `user_biasa`, `kepala_kua` |
| **Location** | `/simnikah/location/*` | GET, POST, PUT | Yes | `user_biasa` |
| **Notification** | `/simnikah/notifikasi/*` | GET, POST, PUT, DELETE | Yes | All |

---

## 🔑 Authentication Endpoints

### 1. Register User

**Endpoint:** `POST /register`

**Description:** Mendaftarkan user baru ke sistem.

**Auth Required:** ❌ No

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
  "user": {
    "user_id": "USR1704067200",
    "username": "ahmad123",
    "email": "ahmad@example.com",
    "nama": "Ahmad Wijaya",
    "role": "user_biasa",
    "created_at": "2024-01-01T10:00:00Z"
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
    "email": "ahmad@example.com",
    "role": "user_biasa",
    "nama": "Ahmad Wijaya"
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
  "user": {
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

## 👰 Catin Endpoints

### 4. Create Registration

**Endpoint:** `POST /simnikah/pendaftaran`

**Description:** Membuat pendaftaran nikah dengan form sederhana.

**Auth Required:** ✅ Yes

**Role Required:** `user_biasa`

**Request Body:**
```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
    "pendidikan_akhir": "S1",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
    "pendidikan_akhir": "S1",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di KUA",
    "tanggal_nikah": "2024-12-15",
    "waktu_nikah": "10:00"
  }
}
```

**Untuk Nikah di Luar KUA:**
```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
    "pendidikan_akhir": "S1",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
    "pendidikan_akhir": "S1",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di Luar KUA",
    "tanggal_nikah": "2024-12-15",
    "waktu_nikah": "10:00",
    "alamat_nikah": "Jl. Ahmad Yani No. 123",
    "alamat_detail": "Rumah Pengantin Perempuan",
    "kelurahan": "Pangeran"
  }
}
```

**Validasi:**
- Umur minimal: 19 tahun (untuk calon laki-laki dan perempuan)
- Format tanggal: `YYYY-MM-DD`
- Format waktu: `HH:MM` (24-jam, contoh: `09:00`, `14:30`)
- Tanggal nikah tidak boleh di masa lalu
- Kelurahan harus dalam lingkup **Kecamatan Banjarmasin Utara**

**Kelurahan Valid:**
- Alalak
- Antasan Kecil Timur
- Antasan Besar
- Basirih
- Basirih Selatan
- Kuin Cerucuk
- Kuin Selatan
- Mantuil
- Pangeran
- Pelambuan
- Sungai Miai
- Surgi Mufti
- Teluk Dalam
- Telawang

**Response Success (201):**
```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat (form sederhana)",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "status_pendaftaran": "Draft",
    "tanggal_nikah": "2024-12-15T00:00:00Z",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di KUA",
    "alamat_akad": "PH5Q+F8C, Jl. Wira Karya, Pangeran, Kec. Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan 70123",
    "calon_suami": {
      "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
      "pendidikan": "S1",
      "umur": 25
    },
    "calon_istri": {
      "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
      "pendidikan": "S1",
      "umur": 23
    },
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Umur calon laki-laki minimal 19 tahun",
  "field": "umur_laki_laki",
  "type": "validation"
}
```

---

### 5. Get User Registration Status

**Endpoint:** `GET /simnikah/pendaftaran/status`

**Description:** Cek status pendaftaran nikah user yang sedang login.

**Auth Required:** ✅ Yes

**Role Required:** `user_biasa`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Status pendaftaran berhasil diambil",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "tanggal_nikah": "2024-12-15T00:00:00Z",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di KUA",
    "penghulu": {
      "nama": "H. Muhammad Amin",
      "nip": "198001012003121001",
      "no_hp": "081234567890",
      "email": "amin@kua.go.id",
      "alamat": "Jl. Penghulu No. 123, Banjarmasin Utara",
      "ditugaskan_pada": "2024-12-10T10:00:00Z"
    },
    "calon_suami": {
      "id": "abc123",
      "nama_lengkap": "Ahmad Wijaya bin Abdullah"
    },
    "calon_istri": {
      "id": "def456",
      "nama_lengkap": "Siti Nurhaliza binti Muhammad"
    }
  }
}
```

**Response (404) - Tidak ada pendaftaran:**
```json
{
  "success": true,
  "message": "User belum memiliki pendaftaran nikah",
  "data": null
}
```

---

### 6. List All Registrations

**Endpoint:** `GET /simnikah/pendaftaran`

**Description:** Mendapatkan semua pendaftaran nikah (hanya untuk Staff dan Kepala KUA).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

**Query Parameters:**
- `status` (optional): Filter by status (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai, Ditolak)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Example:**
```
GET /simnikah/pendaftaran?status=Disetujui&page=1&limit=10
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": [
    {
      "id": 1,
      "nomor_pendaftaran": "NIKAH-20241215-1234",
      "pendaftar_id": "USR1704067200",
      "status_pendaftaran": "Disetujui",
      "tanggal_pendaftaran": "2024-01-01T10:00:00Z",
      "tanggal_nikah": "2024-12-15T00:00:00Z",
      "waktu_nikah": "10:00",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "PH5Q+F8C, Jl. Wira Karya...",
      "penghulu_id": null,
      "catatan": "",
      "calon_suami": {
        "id": "abc123",
        "nama_lengkap": "Ahmad Wijaya bin Abdullah"
      },
      "calon_istri": {
        "id": "def456",
        "nama_lengkap": "Siti Nurhaliza binti Muhammad"
      },
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

---

## 📅 Calendar & Availability Endpoints

### 7. Get Calendar Availability

**Endpoint:** `GET /simnikah/kalender-ketersediaan`

**Description:** Mendapatkan ketersediaan tanggal untuk bulan tertentu (PUBLIC - tidak perlu auth).

**Auth Required:** ❌ No

**Query Parameters:**
- `bulan` (required): Bulan (1-12)
- `tahun` (required): Tahun (YYYY)

**Example:**
```
GET /simnikah/kalender-ketersediaan?bulan=12&tahun=2024
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Ketersediaan kalender berhasil diambil",
  "data": {
    "bulan": 12,
    "tahun": 2024,
    "kapasitas_per_hari": 9,
    "hari": [
      {
        "tanggal": "2024-12-01",
        "hari": "Sunday",
        "total_nikah": 5,
        "sisa_kuota": 4,
        "status": "Sebagian Tersedia"
      },
      {
        "tanggal": "2024-12-02",
        "hari": "Monday",
        "total_nikah": 0,
        "sisa_kuota": 9,
        "status": "Semua Tersedia"
      },
      {
        "tanggal": "2024-12-03",
        "hari": "Tuesday",
        "total_nikah": 9,
        "sisa_kuota": 0,
        "status": "Penuh"
      }
    ]
  }
}
```

---

### 8. Get Available Time Slots

**Endpoint:** `GET /simnikah/ketersediaan-jam`

**Description:** Mendapatkan ketersediaan slot waktu untuk tanggal tertentu (PUBLIC - tidak perlu auth).

**Auth Required:** ❌ No

**Query Parameters:**
- `tanggal` (required): Tanggal dalam format `YYYY-MM-DD`

**Example:**
```
GET /simnikah/ketersediaan-jam?tanggal=2024-12-15
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Ketersediaan jam berhasil diambil",
  "data": {
    "tanggal": "2024-12-15",
    "hari": "Sunday",
    "status": "Sebagian Tersedia",
    "summary": {
      "total_slot": 9,
      "terbooking": 5,
      "tersedia": 4,
      "sisa_kuota": 4
    },
    "time_slots": [
      {
        "waktu": "08:00",
        "tersedia": true,
        "status": "Tersedia"
      },
      {
        "waktu": "09:00",
        "tersedia": false,
        "status": "Terbooking"
      },
      {
        "waktu": "10:00",
        "tersedia": true,
        "status": "Tersedia"
      }
    ],
    "registrations_today": {
      "total": 5,
      "detail": [
        {
          "waktu": "09:00",
          "nama_suami": "Ahmad Wijaya",
          "nama_istri": "Siti Nurhaliza"
        }
      ]
    }
  }
}
```

---

### 9. Get Weddings By Date

**Endpoint:** `GET /simnikah/pernikahan-tanggal`

**Description:** Mendapatkan informasi detail pernikahan pada tanggal tertentu (PUBLIC - tidak perlu auth).

**Auth Required:** ❌ No

**Query Parameters:**
- `tanggal` (required): Tanggal dalam format `YYYY-MM-DD`

**Example:**
```
GET /simnikah/pernikahan-tanggal?tanggal=2024-12-15
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data pernikahan pada tanggal berhasil diambil",
  "data": {
    "tanggal": "2024-12-15",
    "hari": "Sunday",
    "tanggal_format": "15 Desember 2024",
    "summary": {
      "total_nikah": 5,
      "nikah_di_kua": 4,
      "nikah_di_luar": 1
    },
    "pernikahan": [
      {
        "nomor_pendaftaran": "NIKAH-20241215-1234",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "PH5Q+F8C, Jl. Wira Karya...",
        "status_pendaftaran": "Penghulu Ditugaskan",
        "penghulu": {
          "nama": "H. Muhammad Amin",
          "nip": "198001012003121001"
        },
        "calon_suami": {
          "id": "abc123",
          "nama_lengkap": "Ahmad Wijaya bin Abdullah"
        },
        "calon_istri": {
          "id": "def456",
          "nama_lengkap": "Siti Nurhaliza binti Muhammad"
        }
      }
    ]
  }
}
```

---

## 👨‍💼 Staff Endpoints

### 10. List Staff

**Endpoint:** `GET /simnikah/staff`

**Description:** Mendapatkan daftar semua staff KUA.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data staff berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067200",
      "nip": "198001012003121001",
      "nama_lengkap": "Budi Santoso",
      "jabatan": "Staff Pendaftaran",
      "no_hp": "081234567890",
      "email": "budi@kua.go.id",
      "status": "Aktif",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### 11. Update Staff

**Endpoint:** `PUT /simnikah/staff/:id`

**Description:** Mengupdate data staff KUA.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "nama_lengkap": "Budi Santoso",
  "jabatan": "Staff Pendaftaran",
  "no_hp": "081234567890",
  "email": "budi@kua.go.id",
  "status": "Aktif"
}
```

---

### 12. Verify Registration Form

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`

**Description:** Verifikasi formulir pendaftaran (mengubah status dari Draft ke Disetujui).

**Auth Required:** ✅ Yes

**Role Required:** `staff`

**Request Body:**
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir sudah lengkap dan valid"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Formulir berhasil disetujui dan status diubah ke Disetujui",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "status_pendaftaran": "Disetujui",
    "disetujui_oleh": "USR1704067200",
    "disetujui_pada": "2024-01-01T10:00:00Z",
    "catatan": "Formulir sudah lengkap dan valid",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### 13. Verify Documents

**Endpoint:** `POST /simnikah/staff/verify-berkas/:id`

**Description:** Verifikasi berkas fisik pendaftaran.

**Auth Required:** ✅ Yes

**Role Required:** `staff`

**Request Body:**
```json
{
  "status": "Berkas Diterima",
  "catatan": "Semua berkas lengkap"
}
```

---

### 14. Approve Registration

**Endpoint:** `POST /simnikah/staff/approve/:id`

**Description:** Approve atau reject pendaftaran nikah.

**Auth Required:** ✅ Yes

**Role Required:** `staff`

**Request Body:**
```json
{
  "status": "Disetujui",
  "catatan": "Pendaftaran disetujui"
}
```

**Atau untuk reject:**
```json
{
  "status": "Ditolak",
  "catatan": "Data tidak lengkap"
}
```

---

### 15. Update Registration Status

**Endpoint:** `PUT /simnikah/pendaftaran/:id/update-status`

**Description:** Update status pendaftaran secara fleksibel (tidak bisa update status terkait penghulu).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `penghulu`, `kepala_kua`

**Request Body:**
```json
{
  "status": "Disetujui",
  "catatan": "Status diupdate"
}
```

**Note:** Status berikut TIDAK bisa diupdate melalui endpoint ini (hanya melalui endpoint khusus):
- `Menunggu Penugasan` (hanya via auto setelah approve)
- `Penghulu Ditugaskan` (hanya via assign-penghulu oleh kepala KUA)

---

## 👨‍⚖️ Penghulu Endpoints

### 16. List Marriage Officers

**Endpoint:** `GET /simnikah/penghulu`

**Description:** Mendapatkan daftar semua penghulu.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data penghulu berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067201",
      "nip": "198001012003121002",
      "nama_lengkap": "H. Muhammad Amin",
      "no_hp": "081234567891",
      "email": "amin@kua.go.id",
      "alamat": "Jl. Penghulu No. 123",
      "status": "Aktif",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### 17. Update Marriage Officer

**Endpoint:** `PUT /simnikah/penghulu/:id`

**Description:** Mengupdate data penghulu.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

---

### 18. Get My Assignments

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`

**Description:** Mendapatkan daftar pendaftaran yang ditugaskan ke penghulu yang sedang login.

**Auth Required:** ✅ Yes

**Role Required:** `penghulu`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": {
    "penghulu": "H. Muhammad Amin",
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIKAH-20241215-1234",
        "status_pendaftaran": "Penghulu Ditugaskan",
        "tanggal_nikah": "2024-12-15T00:00:00Z",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_akad": "Jl. Ahmad Yani No. 123",
        "latitude": -3.291304,
        "longitude": 114.588147,
        "has_coordinates": true,
        "is_outside_kua": true,
        "note": "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi.",
        "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291304,114.588147",
        "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291304,114.588147",
        "waze_url": "https://www.waze.com/ul?ll=-3.291304,114.588147&navigate=yes",
        "osm_url": "https://www.openstreetmap.org/?mlat=-3.291304&mlon=114.588147&zoom=16",
        "catatan": "",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 1
  }
}
```

---

### 19. Complete Marriage

**Endpoint:** `POST /simnikah/penghulu/complete-marriage/:id`

**Description:** Menandai bahwa pernikahan sudah selesai dilaksanakan.

**Auth Required:** ✅ Yes

**Role Required:** `penghulu`

**Request Body:**
```json
{
  "catatan": "Pernikahan sudah dilaksanakan dengan lancar"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Pernikahan berhasil ditandai selesai",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "status_pendaftaran": "Selesai",
    "updated_at": "2024-12-15T10:00:00Z"
  }
}
```

---

## 👔 Kepala KUA Endpoints

### 20. Create Staff

**Endpoint:** `POST /simnikah/kepala-kua/staff`

**Description:** Membuat akun staff baru.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "username": "staff001",
  "email": "staff001@kua.go.id",
  "password": "password123",
  "nama": "Budi Santoso",
  "nip": "198001012003121001",
  "jabatan": "Staff Pendaftaran",
  "no_hp": "081234567890",
  "alamat": "Jl. Staff No. 123"
}
```

---

### 21. Create Marriage Officer

**Endpoint:** `POST /simnikah/kepala-kua/penghulu`

**Description:** Membuat akun penghulu baru.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "username": "penghulu001",
  "email": "penghulu001@kua.go.id",
  "password": "password123",
  "nama": "H. Muhammad Amin",
  "nip": "198001012003121002",
  "no_hp": "081234567891",
  "alamat": "Jl. Penghulu No. 123"
}
```

---

### 22. Assign Marriage Officer

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`

**Description:** Menugaskan penghulu ke pendaftaran nikah (HANYA KEPALA KUA).

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Request Body:**
```json
{
  "penghulu_id": 1,
  "catatan": "Penghulu ditugaskan untuk menikahkan pasangan ini"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Penghulu berhasil ditugaskan",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "penghulu_id": 1,
    "penghulu_nama": "H. Muhammad Amin",
    "penghulu_assigned_by": "USR1704067202",
    "penghulu_assigned_at": "2024-12-10T10:00:00Z",
    "catatan": "Penghulu ditugaskan untuk menikahkan pasangan ini",
    "updated_at": "2024-12-10T10:00:00Z"
  }
}
```

---

### 23. List Available Officers

**Endpoint:** `GET /simnikah/kepala-kua/available-penghulu`

**Description:** Mendapatkan daftar penghulu yang tersedia (aktif dan tidak terlalu sibuk).

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data penghulu tersedia berhasil diambil",
  "data": [
    {
      "id": 1,
      "nama_lengkap": "H. Muhammad Amin",
      "nip": "198001012003121002",
      "status": "Aktif",
      "total_assigned": 5,
      "available": true
    }
  ]
}
```

---

### 24. Get Penghulu Statistics

**Endpoint:** `GET /simnikah/kepala-kua/statistik-penghulu`

**Description:** Mendapatkan statistik penghulu (total nikah, per bulan, dll).

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Query Parameters:**
- `penghulu_id` (optional): Filter by penghulu ID
- `bulan` (optional): Filter by bulan (1-12)
- `tahun` (optional): Filter by tahun (YYYY)

**Example:**
```
GET /simnikah/kepala-kua/statistik-penghulu?penghulu_id=1&bulan=12&tahun=2024
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Statistik penghulu berhasil diambil",
  "data": {
    "penghulu": {
      "id": 1,
      "nama_lengkap": "H. Muhammad Amin",
      "nip": "198001012003121002"
    },
    "overall": {
      "total_nikah": 150,
      "nikah_di_kua": 120,
      "nikah_di_luar": 30
    },
    "monthly": {
      "bulan": 12,
      "tahun": 2024,
      "total_nikah": 15,
      "nikah_di_kua": 12,
      "nikah_di_luar": 3
    },
    "recent": [
      {
        "tanggal": "2024-12-15",
        "total": 2
      }
    ]
  }
}
```

---

## 💬 Feedback Endpoints

### 25. Create Feedback

**Endpoint:** `POST /simnikah/feedback-pernikahan`

**Description:** Membuat feedback pernikahan dari catin (setelah pernikahan selesai).

**Auth Required:** ✅ Yes

**Role Required:** `user_biasa`

**Request Body:**
```json
{
  "pendaftaran_id": 1,
  "jenis_feedback": "Rating",
  "rating": 5,
  "judul": "Pelayanan Sangat Baik",
  "pesan": "Terima kasih untuk pelayanan yang sangat baik. Proses nikah berjalan lancar."
}
```

**Jenis Feedback:**
- `Rating` - Rating 1-5 (harus include `rating`)
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
    "pesan": "Terima kasih untuk pelayanan yang sangat baik...",
    "status_baca": "Belum Dibaca",
    "created_at": "2024-12-15T10:00:00Z"
  }
}
```

---

### 26. List Feedback (Kepala KUA)

**Endpoint:** `GET /simnikah/kepala-kua/feedback`

**Description:** Mendapatkan daftar semua feedback (hanya untuk Kepala KUA).

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Query Parameters:**
- `jenis_feedback` (optional): Filter by jenis (Rating, Saran, Kritik, Laporan)
- `status_baca` (optional): Filter by status (Belum Dibaca, Sudah Dibaca)
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data feedback berhasil diambil",
  "data": [
    {
      "id": 1,
      "pendaftaran_id": 1,
      "user_id": "USR1704067200",
      "jenis_feedback": "Rating",
      "rating": 5,
      "judul": "Pelayanan Sangat Baik",
      "pesan": "Terima kasih untuk pelayanan yang sangat baik...",
      "status_baca": "Belum Dibaca",
      "dibaca_oleh": null,
      "dibaca_pada": null,
      "created_at": "2024-12-15T10:00:00Z"
    }
  ]
}
```

---

### 27. Mark Feedback As Read

**Endpoint:** `PUT /simnikah/kepala-kua/feedback/:id/mark-read`

**Description:** Menandai feedback sebagai sudah dibaca.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Feedback berhasil ditandai sebagai sudah dibaca",
  "data": {
    "id": 1,
    "status_baca": "Sudah Dibaca",
    "dibaca_oleh": "USR1704067202",
    "dibaca_pada": "2024-12-16T10:00:00Z"
  }
}
```

---

### 28. Get Feedback Statistics

**Endpoint:** `GET /simnikah/kepala-kua/feedback/stats`

**Description:** Mendapatkan statistik feedback.

**Auth Required:** ✅ Yes

**Role Required:** `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Statistik feedback berhasil diambil",
  "data": {
    "total_feedback": 100,
    "total_belum_dibaca": 25,
    "total_sudah_dibaca": 75,
    "rating_rata_rata": 4.5,
    "per_jenis": {
      "Rating": 50,
      "Saran": 30,
      "Kritik": 15,
      "Laporan": 5
    }
  }
}
```

---

## 📍 Location Endpoints

### 29. Geocode Address

**Endpoint:** `POST /simnikah/location/geocode`

**Description:** Convert alamat ke koordinat (latitude, longitude).

**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "alamat": "Jl. Ahmad Yani No. 123, Banjarmasin"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Koordinat berhasil didapatkan",
  "data": {
    "alamat": "Jl. Ahmad Yani No. 123, Banjarmasin",
    "latitude": -3.291304,
    "longitude": 114.588147
  }
}
```

---

### 30. Reverse Geocode

**Endpoint:** `POST /simnikah/location/reverse-geocode`

**Description:** Convert koordinat ke alamat.

**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "latitude": -3.291304,
  "longitude": 114.588147
}
```

---

### 31. Get Wedding Location Detail

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`

**Description:** Mendapatkan detail lokasi pernikahan lengkap dengan link navigasi.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "message": "Detail lokasi berhasil diambil",
  "data": {
    "pendaftaran_id": 1,
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Ahmad Yani No. 123",
    "latitude": -3.291304,
    "longitude": 114.588147,
    "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291304,114.588147",
    "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291304,114.588147",
    "waze_url": "https://www.waze.com/ul?ll=-3.291304,114.588147&navigate=yes",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.291304&mlon=114.588147&zoom=16"
  }
}
```

---

## 🔔 Notification Endpoints

### 32. Get User Notifications

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`

**Description:** Mendapatkan semua notifikasi untuk user tertentu.

**Auth Required:** ✅ Yes

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data notifikasi berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067200",
      "judul": "Pendaftaran Disetujui",
      "pesan": "Pendaftaran nikah Anda telah disetujui.",
      "tipe": "Success",
      "status_baca": "Belum Dibaca",
      "link": "/pendaftaran/1",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

---

### 33. Create Notification

**Endpoint:** `POST /simnikah/notifikasi`

**Description:** Membuat notifikasi baru (hanya untuk Staff dan Kepala KUA).

**Auth Required:** ✅ Yes

**Role Required:** `staff`, `kepala_kua`

---

## ❌ Error Handling

### Error Response Format

Semua error mengikuti format berikut:

```json
{
  "success": false,
  "message": "Pesan error umum",
  "error": "Detail error atau pesan spesifik",
  "type": "jenis_error",
  "field": "field_yang_error" // optional
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `201` | Created |
| `400` | Bad Request (validation error) |
| `401` | Unauthorized (tidak ada token atau token invalid) |
| `403` | Forbidden (tidak punya permission) |
| `404` | Not Found |
| `405` | Method Not Allowed |
| `500` | Internal Server Error |

### Error Types

- `authentication` - Error authentication
- `authorization` - Error authorization/permission
- `validation` - Error validasi input
- `not_found` - Resource tidak ditemukan
- `database` - Error database
- `format` - Error format data

### Contoh Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token tidak valid atau sudah kadaluarsa",
  "type": "authentication"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak memiliki permission untuk mengakses endpoint ini",
  "type": "authorization"
}
```

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Umur calon laki-laki minimal 19 tahun",
  "field": "umur_laki_laki",
  "type": "validation"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Endpoint tidak ditemukan",
  "error": "Path '/api/invalid' tidak ditemukan"
}
```

---

## 📊 Status Flow

### Status Pendaftaran Nikah

```
1. Draft
   └─> Catin mengisi form pendaftaran

2. Disetujui
   └─> Staff approve pendaftaran

3. Menunggu Penugasan (opsional)
   └─> Menunggu kepala KUA assign penghulu

4. Penghulu Ditugaskan
   └─> Kepala KUA assign penghulu

5. Selesai
   └─> Penghulu complete nikah

(6. Ditolak - jika ditolak oleh staff)
```

### Status Constants

- `Draft` - Pendaftaran baru dibuat
- `Disetujui` - Staff sudah approve
- `Menunggu Penugasan` - Menunggu kepala KUA assign penghulu
- `Penghulu Ditugaskan` - Penghulu sudah ditugaskan
- `Selesai` - Pernikahan sudah dilaksanakan
- `Ditolak` - Pendaftaran ditolak

---

## 🎯 Role-Based Access Control

| Role | Description | Endpoints Access |
|------|-------------|------------------|
| `user_biasa` | Calon pengantin | Pendaftaran, status, kalender, feedback |
| `staff` | Staff KUA | Verifikasi, approve, update status |
| `penghulu` | Penghulu | Lihat tugas, complete nikah |
| `kepala_kua` | Kepala KUA | Assign penghulu, statistik, manajemen staff/penghulu, feedback management |

---

## 📝 Notes

1. **Rate Limiting:**
   - Global: 100 requests/minute per IP
   - Auth endpoints: 10 requests/minute per IP

2. **CORS:**
   - Allowed origins dikonfigurasi di environment variable
   - Credentials: enabled

3. **Database:**
   - MySQL/MariaDB
   - Auto migration pada startup

4. **Time Format:**
   - Tanggal: `YYYY-MM-DD` (ISO 8601)
   - Waktu: `HH:MM` (24-jam format)
   - Timestamp: ISO 8601 (contoh: `2024-01-01T10:00:00Z`)

5. **KUA Location:**
   - Alamat: `PH5Q+F8C, Jl. Wira Karya, Pangeran, Kec. Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan 70123`
   - Koordinat: `-3.291304649442475, 114.58814746634684`

---

## 🔗 Additional Resources

- **Status Flow:** Lihat `docs/FLOW_SEDERHANA.md`
- **Architecture:** Lihat `docs/architecture/`
- **Testing:** Lihat `docs/api/API_TESTING_DOCUMENTATION.md`

---

**Generated:** 2024  
**API Version:** 1.0.0

