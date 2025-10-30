# 📚 Dokumentasi API SimNikah - Complete Guide

## 📋 Daftar Isi
- [Informasi Umum](#informasi-umum)
- [Authentication & Authorization](#1-authentication--authorization)
- [Marriage Registration](#2-marriage-registration)
- [Staff Management](#3-staff-management)
- [Penghulu Management](#4-penghulu-management)
- [Calendar & Schedule](#5-calendar--schedule)
- [Bimbingan Perkawinan](#6-bimbingan-perkawinan)
- [Notifications](#7-notifications)
- [Location & Map Integration](#8-location--map-integration)
- [Status Codes](#status-codes)
- [Error Handling](#error-handling)

---

## Informasi Umum

### Base URL
```
Development: http://localhost:8080
Production: https://your-domain.railway.app
```

### Authentication
Sebagian besar endpoint memerlukan JWT token. Tambahkan header berikut pada setiap request:
```
Authorization: Bearer <your_jwt_token>
```

### Content Type
```
Content-Type: application/json
```

### User Roles
- `user_biasa` - Calon pengantin yang mendaftar nikah
- `staff` - Staff KUA untuk verifikasi
- `penghulu` - Penghulu yang memimpin nikah
- `kepala_kua` - Kepala KUA (akses penuh)

---

## 1. Authentication & Authorization

### 1.1 Register User

**Endpoint:** `POST /register`  
**Authorization:** None  
**Description:** Mendaftarkan user baru

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

**Field Validation:**
- `username`: Required, unique
- `email`: Required, valid email format, unique
- `password`: Required, minimum 6 characters
- `nama`: Required
- `role`: Required, one of: `user_biasa`, `staff`, `penghulu`, `kepala_kua`

**Success Response (201):**
```json
{
  "message": "User berhasil dibuat",
  "user": {
    "user_id": "USR1234567890",
    "username": "johndoe",
    "email": "john@example.com",
    "nama": "John Doe",
    "role": "user_biasa"
  }
}
```

**Error Responses:**
```json
// Username/Email sudah ada
{
  "error": "Username atau email sudah digunakan"
}

// Role tidak valid
{
  "error": "Role tidak valid. Role yang tersedia: user_biasa, penghulu, staff, kepala_kua"
}
```

---

### 1.2 Login

**Endpoint:** `POST /login`  
**Authorization:** None  
**Description:** Login dan mendapatkan JWT token

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR1234567890",
    "email": "john@example.com",
    "role": "user_biasa",
    "nama": "John Doe"
  }
}
```

**Token Validity:** 24 jam

**Error Responses:**
```json
// Username tidak ditemukan
{
  "error": "Username tidak ditemukan"
}

// Password salah
{
  "error": "Password salah"
}

// User tidak aktif
{
  "error": "User tidak aktif"
}
```

---

### 1.3 Get Profile

**Endpoint:** `GET /profile`  
**Authorization:** Required (Bearer Token)  
**Description:** Mendapatkan profil user yang sedang login

**Success Response (200):**
```json
{
  "message": "Profile berhasil diambil",
  "user": {
    "user_id": "USR1234567890",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user_biasa",
    "nama": "John Doe"
  }
}
```

---

## 2. Marriage Registration

### 2.1 Create Marriage Registration Form

**Endpoint:** `POST /simnikah/pendaftaran/form-baru`  
**Authorization:** Required (user_biasa)  
**Description:** Membuat pendaftaran nikah baru (form lengkap)

**Request Body:**
```json
{
  "scheduleAndLocation": {
    "weddingLocation": "Di Luar KUA",
    "weddingAddress": "Jl. Merdeka No. 123, RT 05, RW 02, Banjarmasin",
    "weddingDate": "2024-12-25",
    "weddingTime": "09:00",
    "dispensationNumber": ""
  },
  "groom": {
    "groomFullName": "Ahmad Fauzi",
    "groomNik": "6371012501950001",
    "groomCitizenship": "WNI",
    "groomPassportNumber": "",
    "groomPlaceOfBirth": "Banjarmasin",
    "groomDateOfBirth": "1995-01-25",
    "groomStatus": "Belum Kawin",
    "groomReligion": "Islam",
    "groomEducation": "S1",
    "groomOccupation": "Pegawai Swasta",
    "groomOccupationDescription": "Marketing Manager",
    "groomPhoneNumber": "081234567890",
    "groomEmail": "ahmad.fauzi@email.com",
    "groomAddress": "Jl. Sudirman No. 45, Banjarmasin"
  },
  "bride": {
    "brideFullName": "Siti Nurhaliza",
    "brideNik": "6371014203970001",
    "brideCitizenship": "WNI",
    "bridePassportNumber": "",
    "bridePlaceOfBirth": "Banjarmasin",
    "brideDateOfBirth": "1997-03-12",
    "brideStatus": "Belum Kawin",
    "brideReligion": "Islam",
    "brideEducation": "S1",
    "brideOccupation": "Guru",
    "brideOccupationDescription": "Guru SD",
    "bridePhoneNumber": "082345678901",
    "brideEmail": "siti.nurhaliza@email.com",
    "brideAddress": "Jl. Ahmad Yani No. 67, Banjarmasin"
  },
  "groomParents": {
    "groomFather": {
      "groomFatherPresenceStatus": "Hidup",
      "groomFatherName": "Bapak Ahmad",
      "groomFatherNik": "6371011501700001",
      "groomFatherCitizenship": "WNI",
      "groomFatherCountryOfOrigin": "Indonesia",
      "groomFatherPassportNumber": "",
      "groomFatherPlaceOfBirth": "Banjarmasin",
      "groomFatherDateOfBirth": "1970-01-15",
      "groomFatherReligion": "Islam",
      "groomFatherOccupation": "Wiraswasta",
      "groomFatherOccupationDescription": "Pedagang",
      "groomFatherAddress": "Jl. Sudirman No. 45, Banjarmasin"
    },
    "groomMother": {
      "groomMotherPresenceStatus": "Hidup",
      "groomMotherName": "Ibu Fatimah",
      "groomMotherNik": "6371015201720001",
      "groomMotherCitizenship": "WNI",
      "groomMotherCountryOfOrigin": "Indonesia",
      "groomMotherPassportNumber": "",
      "groomMotherPlaceOfBirth": "Banjarmasin",
      "groomMotherDateOfBirth": "1972-12-20",
      "groomMotherReligion": "Islam",
      "groomMotherOccupation": "Ibu Rumah Tangga",
      "groomMotherOccupationDescription": "",
      "groomMotherAddress": "Jl. Sudirman No. 45, Banjarmasin"
    }
  },
  "brideParents": {
    "brideFather": {
      "brideFatherPresenceStatus": "Hidup",
      "brideFatherName": "Bapak Hasan",
      "brideFatherNik": "6371012501680001",
      "brideFatherCitizenship": "WNI",
      "brideFatherCountryOfOrigin": "Indonesia",
      "brideFatherPassportNumber": "",
      "brideFatherPlaceOfBirth": "Banjarmasin",
      "brideFatherDateOfBirth": "1968-01-25",
      "brideFatherReligion": "Islam",
      "brideFatherOccupation": "PNS",
      "brideFatherOccupationDescription": "Guru",
      "brideFatherAddress": "Jl. Ahmad Yani No. 67, Banjarmasin"
    },
    "brideMother": {
      "brideMotherPresenceStatus": "Hidup",
      "brideMotherName": "Ibu Aminah",
      "brideMotherNik": "6371014203700001",
      "brideMotherCitizenship": "WNI",
      "brideMotherCountryOfOrigin": "Indonesia",
      "brideMotherPassportNumber": "",
      "brideMotherPlaceOfBirth": "Banjarmasin",
      "brideMotherDateOfBirth": "1970-03-12",
      "brideMotherReligion": "Islam",
      "brideMotherOccupation": "Ibu Rumah Tangga",
      "brideMotherOccupationDescription": "",
      "brideMotherAddress": "Jl. Ahmad Yani No. 67, Banjarmasin"
    }
  },
  "guardian": {
    "guardianFullName": "Bapak Hasan",
    "guardianNik": "6371012501680001",
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Hidup",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. Ahmad Yani No. 67, Banjarmasin",
    "guardianPhoneNumber": "082345678901"
  }
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Formulir pendaftaran nikah berhasil dibuat",
  "data": {
    "pendaftaran": {
      "id": 1,
      "nomor_pendaftaran": "REG20241120001",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2024-12-25T00:00:00Z",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di Luar KUA",
      "alamat_akad": "Jl. Merdeka No. 123, RT 05, RW 02, Banjarmasin"
    },
    "calon_suami": {
      "id": 1,
      "nama_lengkap": "Ahmad Fauzi",
      "nik": "6371012501950001"
    },
    "calon_istri": {
      "id": 2,
      "nama_lengkap": "Siti Nurhaliza",
      "nik": "6371014203970001"
    },
    "wali_nikah": {
      "id": 1,
      "nama_lengkap": "Bapak Hasan",
      "hubungan_wali": "Ayah Kandung"
    },
    "next_steps": [
      "Tunggu verifikasi formulir oleh staff KUA",
      "Anda akan menerima notifikasi jika formulir disetujui",
      "Setelah disetujui, kumpulkan berkas ke KUA dalam 5 hari kerja"
    ]
  }
}
```

**Field Options:**
- `weddingLocation`: "Di KUA" atau "Di Luar KUA"
- `groomStatus/brideStatus`: "Belum Kawin", "Cerai Hidup", "Cerai Mati"
- `groomReligion/brideReligion`: "Islam", "Kristen", "Katolik", "Hindu", "Buddha", "Konghucu"
- `groomCitizenship/brideCitizenship`: "WNI", "WNA"
- `guardianRelationship`: "Ayah Kandung", "Kakek", "Saudara Kandung", "Paman", "Wali Hakim"
- `presenceStatus`: "Hidup", "Meninggal"

---

### 2.2 Check Registration Status

**Endpoint:** `GET /simnikah/pendaftaran/status`  
**Authorization:** Required (user_biasa)  
**Description:** Mengecek status pendaftaran nikah user yang sedang login

**Success Response (200):**
```json
{
  "success": true,
  "message": "Status pendaftaran berhasil diambil",
  "data": {
    "has_registration": true,
    "registration": {
      "id": 1,
      "nomor_pendaftaran": "REG20241120001",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2024-12-25",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di Luar KUA",
      "alamat_akad": "Jl. Merdeka No. 123, Banjarmasin",
      "created_at": "2024-11-20T10:30:00Z"
    },
    "calon_suami": {
      "nama_lengkap": "Ahmad Fauzi",
      "nik": "6371012501950001"
    },
    "calon_istri": {
      "nama_lengkap": "Siti Nurhaliza",
      "nik": "6371014203970001"
    },
    "status_flow": {
      "current": "Menunggu Verifikasi",
      "next": "Menunggu Pengumpulan Berkas",
      "can_edit": false
    }
  }
}
```

**Response when no registration:**
```json
{
  "success": true,
  "message": "Belum ada pendaftaran",
  "data": {
    "has_registration": false
  }
}
```

---

### 2.3 Get All Registrations

**Endpoint:** `GET /simnikah/pendaftaran`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Mendapatkan semua pendaftaran nikah

**Query Parameters:**
- `status` (optional): Filter berdasarkan status
- `page` (optional): Nomor halaman (default: 1)
- `limit` (optional): Jumlah data per halaman (default: 10)

**Example Request:**
```
GET /simnikah/pendaftaran?status=Menunggu Verifikasi&page=1&limit=20
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": {
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG20241120001",
        "status_pendaftaran": "Menunggu Verifikasi",
        "tanggal_nikah": "2024-12-25",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "calon_suami_nama": "Ahmad Fauzi",
        "calon_istri_nama": "Siti Nurhaliza",
        "created_at": "2024-11-20T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45,
      "total_page": 3
    }
  }
}
```

---

### 2.4 Mark as Visited

**Endpoint:** `POST /simnikah/pendaftaran/:id/mark-visited`  
**Authorization:** Required (staff)  
**Description:** Menandai bahwa calon pengantin sudah datang ke KUA untuk pengumpulan berkas

**Success Response (200):**
```json
{
  "success": true,
  "message": "Pendaftaran berhasil ditandai sudah datang",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "status_pendaftaran": "Menunggu Penugasan",
    "updated_at": "2024-11-21T14:30:00Z"
  }
}
```

---

### 2.5 Update Wedding Address

**Endpoint:** `PUT /simnikah/pendaftaran/:id/alamat`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Update alamat nikah

**Request Body:**
```json
{
  "alamat_akad": "Jl. Merdeka No. 123, RT 05, RW 02, Kelurahan Sungai Baru, Banjarmasin Tengah"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Alamat nikah berhasil diupdate",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "alamat_akad": "Jl. Merdeka No. 123, RT 05, RW 02, Kelurahan Sungai Baru, Banjarmasin Tengah",
    "updated_at": "2024-11-21T15:00:00Z"
  }
}
```

---

## 3. Staff Management

### 3.1 Create Staff KUA

**Endpoint:** `POST /simnikah/staff`  
**Authorization:** Required (kepala_kua)  
**Description:** Membuat staff KUA baru

**Request Body:**
```json
{
  "username": "staff_ali",
  "email": "ali@kua.go.id",
  "password": "staffpass123",
  "nama": "Ali Rahman",
  "nip": "198501012010011001",
  "jabatan": "Staff",
  "bagian": "Pendaftaran",
  "no_hp": "081234567890",
  "alamat": "Jl. Ahmad Yani No. 100, Banjarmasin"
}
```

**Jabatan Options:**
- `Staff`
- `Penghulu`
- `Kepala KUA`

**Success Response (201):**
```json
{
  "message": "Staff KUA berhasil dibuat",
  "data": {
    "user": {
      "user_id": "STF1700567890",
      "username": "staff_ali",
      "email": "ali@kua.go.id",
      "role": "staff",
      "nama": "Ali Rahman"
    },
    "staff": {
      "id": 1,
      "user_id": "STF1700567890",
      "nip": "198501012010011001",
      "nama_lengkap": "Ali Rahman",
      "jabatan": "Staff",
      "bagian": "Pendaftaran",
      "status": "Aktif"
    }
  }
}
```

---

### 3.2 Get All Staff

**Endpoint:** `GET /simnikah/staff`  
**Authorization:** Required (kepala_kua)  
**Description:** Mendapatkan semua staff KUA

**Success Response (200):**
```json
{
  "message": "Data staff berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "STF1700567890",
      "nip": "198501012010011001",
      "nama_lengkap": "Ali Rahman",
      "jabatan": "Staff",
      "bagian": "Pendaftaran",
      "no_hp": "081234567890",
      "email": "ali@kua.go.id",
      "status": "Aktif",
      "created_at": "2024-01-15T08:00:00Z"
    }
  ]
}
```

---

### 3.3 Update Staff

**Endpoint:** `PUT /simnikah/staff/:id`  
**Authorization:** Required (kepala_kua)  
**Description:** Update informasi staff

**Request Body:**
```json
{
  "nama_lengkap": "Ali Rahman Updated",
  "jabatan": "Kepala Bagian",
  "bagian": "Verifikasi",
  "no_hp": "081234567899",
  "alamat": "Jl. Baru No. 200, Banjarmasin",
  "status": "Aktif"
}
```

**Success Response (200):**
```json
{
  "message": "Data staff berhasil diupdate",
  "data": {
    "id": 1,
    "nama_lengkap": "Ali Rahman Updated",
    "jabatan": "Kepala Bagian",
    "bagian": "Verifikasi",
    "status": "Aktif",
    "updated_at": "2024-11-21T16:00:00Z"
  }
}
```

---

### 3.4 Verify Form (Tahap 1)

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`  
**Authorization:** Required (staff)  
**Description:** Verifikasi formulir online pendaftaran nikah

**Request Body:**
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir sudah lengkap dan valid"
}
```

**Status Options:**
- `Formulir Disetujui` - Otomatis ubah status ke "Menunggu Pengumpulan Berkas"
- `Formulir Ditolak` - Pendaftaran ditolak

**Success Response (200):**
```json
{
  "success": true,
  "message": "Formulir berhasil disetujui dan status diubah ke Pengumpulan Berkas",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "status_pendaftaran": "Menunggu Pengumpulan Berkas",
    "disetujui_oleh": "STF1700567890",
    "disetujui_pada": "2024-11-21T10:00:00Z",
    "catatan": "Formulir sudah lengkap dan valid"
  }
}
```

---

### 3.5 Verify Documents (Tahap 2)

**Endpoint:** `POST /simnikah/staff/verify-berkas/:id`  
**Authorization:** Required (staff)  
**Description:** Verifikasi berkas fisik yang dibawa ke KUA

**Request Body:**
```json
{
  "status": "Berkas Diterima",
  "catatan": "Semua berkas lengkap dan sesuai"
}
```

**Status Options:**
- `Berkas Diterima` - Berkas diterima dan valid
- `Berkas Ditolak` - Berkas tidak lengkap atau tidak valid

**Success Response (200):**
```json
{
  "success": true,
  "message": "Berkas berhasil diverifikasi",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "status_pendaftaran": "Berkas Diterima",
    "disetujui_oleh": "STF1700567890",
    "disetujui_pada": "2024-11-22T11:00:00Z",
    "catatan": "Semua berkas lengkap dan sesuai"
  }
}
```

---

## 4. Penghulu Management

### 4.1 Create Penghulu

**Endpoint:** `POST /simnikah/penghulu`  
**Authorization:** Required (kepala_kua)  
**Description:** Membuat penghulu baru

**Request Body:**
```json
{
  "username": "penghulu_usman",
  "email": "usman@kua.go.id",
  "password": "penghulu123",
  "nama": "H. Usman bin Ali",
  "nip": "197501012000011001",
  "no_hp": "081234567891",
  "alamat": "Jl. Masjid Raya No. 50, Banjarmasin"
}
```

**Success Response (201):**
```json
{
  "message": "Penghulu berhasil dibuat",
  "data": {
    "user": {
      "user_id": "PNG1700567891",
      "username": "penghulu_usman",
      "email": "usman@kua.go.id",
      "role": "penghulu",
      "nama": "H. Usman bin Ali"
    },
    "penghulu": {
      "id": 1,
      "user_id": "PNG1700567891",
      "nip": "197501012000011001",
      "nama_lengkap": "H. Usman bin Ali",
      "status": "Aktif",
      "jumlah_nikah": 0,
      "rating": 0.0
    }
  }
}
```

---

### 4.2 Get All Penghulu

**Endpoint:** `GET /simnikah/penghulu`  
**Authorization:** Required  
**Description:** Mendapatkan semua penghulu

**Success Response (200):**
```json
{
  "message": "Data penghulu berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "PNG1700567891",
      "nip": "197501012000011001",
      "nama_lengkap": "H. Usman bin Ali",
      "no_hp": "081234567891",
      "email": "usman@kua.go.id",
      "status": "Aktif",
      "jumlah_nikah": 15,
      "rating": 4.8
    }
  ]
}
```

---

### 4.3 Update Penghulu

**Endpoint:** `PUT /simnikah/penghulu/:id`  
**Authorization:** Required (kepala_kua)  
**Description:** Update informasi penghulu

**Request Body:**
```json
{
  "nama_lengkap": "H. Usman bin Ali, S.Ag",
  "no_hp": "081234567892",
  "alamat": "Jl. Masjid Raya No. 51, Banjarmasin",
  "status": "Aktif",
  "rating": 4.9
}
```

**Success Response (200):**
```json
{
  "message": "Data penghulu berhasil diupdate",
  "data": {
    "id": 1,
    "nama_lengkap": "H. Usman bin Ali, S.Ag",
    "rating": 4.9,
    "updated_at": "2024-11-21T17:00:00Z"
  }
}
```

---

### 4.4 Assign Penghulu

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`  
**Authorization:** Required (kepala_kua)  
**Description:** Assign penghulu ke pendaftaran nikah

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
    "penghulu_nama": "H. Usman bin Ali",
    "assigned_by": "KKA1700567890",
    "assigned_at": "2024-11-23T09:00:00Z"
  }
}
```

**Business Rules:**
- Maksimal 3 nikah per penghulu per hari
- Maksimal 9 nikah per hari di KUA
- Minimal jarak 1 jam antar jadwal penghulu yang sama
- Status harus "Menunggu Penugasan"

**Error Responses:**
```json
// Konflik jadwal
{
  "error": "Konflik jadwal! Penghulu sudah memiliki jadwal pada waktu yang berdekatan",
  "details": {
    "waktu_konflik": "09:00",
    "tempat": "Di Luar KUA",
    "selisih_menit": 30,
    "minimal_selisih": "60 menit (1 jam)"
  }
}

// Kuota penuh
{
  "error": "Kuota pernikahan harian penuh",
  "details": {
    "tanggal": "2024-12-25",
    "maksimal": 9
  }
}
```

---

### 4.5 Get Assigned Registrations

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`  
**Authorization:** Required (penghulu)  
**Description:** Mendapatkan pendaftaran yang di-assign ke penghulu yang login

**Success Response (200):**
```json
{
  "success": true,
  "message": "Data pendaftaran berhasil diambil",
  "data": {
    "penghulu": "H. Usman bin Ali",
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG20241120001",
        "status_pendaftaran": "Menunggu Verifikasi Penghulu",
        "tanggal_nikah": "2024-12-25T00:00:00Z",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_akad": "Jl. Merdeka No. 123, Banjarmasin",
        "created_at": "2024-11-20T10:30:00Z"
      }
    ],
    "total": 1
  }
}
```

---

### 4.6 Verify Documents (Penghulu)

**Endpoint:** `POST /simnikah/penghulu/verify-documents/:id`  
**Authorization:** Required (penghulu)  
**Description:** Penghulu memverifikasi dokumen pendaftaran

**Request Body:**
```json
{
  "status": "Menunggu Pelaksanaan",
  "catatan": "Berkas sudah lengkap dan sesuai syariat"
}
```

**Status Options:**
- `Menunggu Pelaksanaan` - Otomatis ubah ke "Menunggu Bimbingan"
- `Ditolak` - Berkas ditolak

**Success Response (200):**
```json
{
  "success": true,
  "message": "Verifikasi berkas berhasil",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "status_pendaftaran": "Menunggu Bimbingan",
    "penghulu_id": 1,
    "catatan": "Berkas sudah lengkap dan sesuai syariat",
    "updated_at": "2024-11-24T10:00:00Z"
  }
}
```

---

## 5. Calendar & Schedule

### 5.1 Get Calendar Availability

**Endpoint:** `GET /simnikah/kalender-ketersediaan`  
**Authorization:** Required  
**Description:** Mendapatkan kalender ketersediaan untuk bulan tertentu

**Query Parameters:**
- `bulan` (optional): Bulan (1-12), default: bulan sekarang
- `tahun` (optional): Tahun (YYYY), default: tahun sekarang

**Example Request:**
```
GET /simnikah/kalender-ketersediaan?bulan=12&tahun=2024
```

**Success Response (200):**
```json
{
  "message": "Kalender ketersediaan berhasil diambil",
  "data": {
    "bulan": 12,
    "tahun": 2024,
    "nama_bulan": "December",
    "kapasitas_harian": 9,
    "penghulu_info": {
      "total_penghulu": 4,
      "penghulu_aktif": 4,
      "slot_waktu_per_hari": 9,
      "nikah_per_slot": 4,
      "total_kapasitas_harian": 9
    },
    "kalender": [
      {
        "tanggal": 25,
        "tanggal_str": "2024-12-25",
        "status": "Tersedia",
        "tersedia": true,
        "jumlah_nikah_total": 5,
        "jumlah_nikah_kua": 4,
        "jumlah_nikah_luar": 1,
        "kuning_count": 2,
        "hijau_count": 3,
        "warna": "hijau",
        "sisa_kuota_kua": 5,
        "kapasitas_kua": 9
      }
    ]
  }
}
```

**Status Types:**
- `Tersedia` - Masih ada kuota
- `Penuh` - Kuota sudah habis
- `Terlewat` - Tanggal sudah lewat

**Warna:**
- `kuning` - Status awal (Menunggu Verifikasi, Menunggu Pengumpulan Berkas)
- `hijau` - Status sudah fix (Berkas Diterima dan seterusnya)

---

### 5.2 Get Calendar Date Detail

**Endpoint:** `GET /simnikah/kalender-tanggal-detail`  
**Authorization:** Required  
**Description:** Mendapatkan detail jadwal per tanggal

**Query Parameters:**
- `tanggal` (required): Tanggal dalam format YYYY-MM-DD

**Example Request:**
```
GET /simnikah/kalender-tanggal-detail?tanggal=2024-12-25
```

**Success Response (200):**
```json
{
  "message": "Detail kalender berhasil diambil",
  "data": {
    "tanggal": "2024-12-25",
    "items": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG20241120001",
        "waktu": "09:00",
        "tempat": "Di Luar KUA",
        "status": "Menunggu Bimbingan",
        "warna": "hijau",
        "penghulu_id": 1,
        "penghulu_nama": "H. Usman bin Ali",
        "nama_calon_suami": "Ahmad Fauzi",
        "nama_calon_istri": "Siti Nurhaliza"
      }
    ]
  }
}
```

---

### 5.3 Get Date Availability

**Endpoint:** `GET /simnikah/ketersediaan-tanggal/:tanggal`  
**Authorization:** Required  
**Description:** Mendapatkan ketersediaan untuk tanggal spesifik

**Example Request:**
```
GET /simnikah/ketersediaan-tanggal/2024-12-25
```

**Success Response (200):**
```json
{
  "message": "Detail ketersediaan tanggal berhasil diambil",
  "data": {
    "tanggal": "2024-12-25",
    "status": "Tersedia",
    "tersedia": true,
    "jumlah_nikah_kua": 4,
    "jumlah_nikah_luar": 2,
    "total_nikah": 6,
    "sisa_kuota_kua": 5,
    "kapasitas_kua": 9,
    "keterangan": "Kapasitas hanya berlaku untuk nikah di KUA. Nikah di luar KUA tidak dibatasi.",
    "jadwal_detail": [
      {
        "nomor_pendaftaran": "REG20241120001",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_akad": "Jl. Merdeka No. 123, Banjarmasin"
      }
    ]
  }
}
```

---

### 5.4 Get Penghulu Schedule

**Endpoint:** `GET /simnikah/penghulu-jadwal/:tanggal`  
**Authorization:** Required  
**Description:** Mendapatkan jadwal semua penghulu untuk tanggal tertentu

**Example Request:**
```
GET /simnikah/penghulu-jadwal/2024-12-25
```

**Success Response (200):**
```json
{
  "message": "Jadwal penghulu berhasil diambil",
  "data": {
    "tanggal": "2024-12-25",
    "total_penghulu": 4,
    "total_kapasitas": 12,
    "total_terisi": 5,
    "total_sisa": 7,
    "penghulu": [
      {
        "id": 1,
        "nama": "H. Usman bin Ali",
        "status": "Sebagian",
        "jumlah_jadwal": 2,
        "sisa_kuota": 1,
        "maksimal": 3,
        "jadwal": [
          {
            "nomor_pendaftaran": "REG20241120001",
            "waktu_nikah": "09:00",
            "tempat_nikah": "Di Luar KUA",
            "alamat_akad": "Jl. Merdeka No. 123",
            "assigned_by": "KKA1700567890",
            "assigned_at": "2024-11-23T09:00:00Z"
          }
        ]
      }
    ]
  }
}
```

**Penghulu Status:**
- `Kosong` - Belum ada jadwal
- `Sebagian` - Ada jadwal tapi belum penuh
- `Penuh` - Sudah 3 jadwal (maksimal)

---

### 5.5 Get Penghulu Availability

**Endpoint:** `GET /simnikah/penghulu/:id/ketersediaan/:tanggal`  
**Authorization:** Required (kepala_kua)  
**Description:** Mendapatkan ketersediaan slot waktu penghulu spesifik

**Example Request:**
```
GET /simnikah/penghulu/1/ketersediaan/2024-12-25
```

**Success Response (200):**
```json
{
  "message": "Ketersediaan penghulu berhasil diambil",
  "data": {
    "penghulu": {
      "id": 1,
      "nama": "H. Usman bin Ali",
      "status": "Aktif"
    },
    "tanggal": "2024-12-25",
    "statistik": {
      "jumlah_jadwal": 2,
      "sisa_kuota": 1,
      "maksimal_per_hari": 3,
      "slot_tersedia": 7,
      "total_slot": 9
    },
    "jadwal_hari_ini": [
      {
        "id": 1,
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di Luar KUA"
      }
    ],
    "slot_waktu": [
      {
        "waktu": "08:00",
        "tersedia": true,
        "konflik_jadwal": []
      },
      {
        "waktu": "09:00",
        "tersedia": false,
        "konflik_jadwal": [
          {
            "waktu": "09:00",
            "tempat": "Di Luar KUA",
            "selisih_menit": 0
          }
        ]
      }
    ]
  }
}
```

---

## 6. Bimbingan Perkawinan

### 6.1 Create Bimbingan Session

**Endpoint:** `POST /simnikah/bimbingan`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Membuat sesi bimbingan perkawinan baru

**Request Body:**
```json
{
  "tanggal_bimbingan": "2024-12-18",
  "waktu_mulai": "08:00",
  "waktu_selesai": "12:00",
  "tempat_bimbingan": "Aula KUA Banjarmasin",
  "pembimbing": "H. Ahmad Dahlan, M.Ag",
  "kapasitas": 10,
  "catatan": "Bimbingan bulan Desember"
}
```

**Business Rules:**
- Hanya bisa dijadwalkan hari Rabu
- Kapasitas default: 10 pasangan
- Hanya 1 sesi per hari

**Success Response (201):**
```json
{
  "message": "Bimbingan perkawinan berhasil dibuat",
  "data": {
    "id": 1,
    "tanggal_bimbingan": "2024-12-18T00:00:00Z",
    "waktu_mulai": "08:00",
    "waktu_selesai": "12:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin",
    "pembimbing": "H. Ahmad Dahlan, M.Ag",
    "kapasitas": 10,
    "status": "Aktif",
    "catatan": "Bimbingan bulan Desember"
  }
}
```

**Error Responses:**
```json
// Bukan hari Rabu
{
  "error": "Bimbingan perkawinan hanya bisa dijadwalkan pada hari Rabu"
}

// Sudah ada bimbingan pada tanggal tersebut
{
  "error": "Sudah ada bimbingan perkawinan pada tanggal tersebut"
}
```

---

### 6.2 Get Bimbingan List

**Endpoint:** `GET /simnikah/bimbingan`  
**Authorization:** Required  
**Description:** Mendapatkan daftar bimbingan perkawinan

**Query Parameters:**
- `bulan` (optional): Bulan (1-12), default: bulan sekarang
- `tahun` (optional): Tahun (YYYY), default: tahun sekarang
- `status` (optional): Status ("Aktif", "Selesai"), default: "Aktif"

**Example Request:**
```
GET /simnikah/bimbingan?bulan=12&tahun=2024&status=Aktif
```

**Success Response (200):**
```json
{
  "message": "Data bimbingan perkawinan berhasil diambil",
  "data": {
    "bulan": 12,
    "tahun": 2024,
    "bimbingan": [
      {
        "id": 1,
        "tanggal_bimbingan": "2024-12-18",
        "waktu_mulai": "08:00",
        "waktu_selesai": "12:00",
        "tempat_bimbingan": "Aula KUA Banjarmasin",
        "pembimbing": "H. Ahmad Dahlan, M.Ag",
        "kapasitas": 10,
        "jumlah_peserta": 7,
        "sisa_kuota": 3,
        "status": "Aktif",
        "catatan": "Bimbingan bulan Desember"
      }
    ]
  }
}
```

---

### 6.3 Get Bimbingan Detail

**Endpoint:** `GET /simnikah/bimbingan/:id`  
**Authorization:** Required  
**Description:** Mendapatkan detail bimbingan spesifik

**Success Response (200):**
```json
{
  "message": "Detail bimbingan perkawinan berhasil diambil",
  "data": {
    "id": 1,
    "tanggal_bimbingan": "2024-12-18",
    "waktu_mulai": "08:00",
    "waktu_selesai": "12:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin",
    "pembimbing": "H. Ahmad Dahlan, M.Ag",
    "kapasitas": 10,
    "jumlah_peserta": 7,
    "sisa_kuota": 3,
    "status": "Aktif",
    "catatan": "Bimbingan bulan Desember"
  }
}
```

---

### 6.4 Register for Bimbingan

**Endpoint:** `POST /simnikah/bimbingan/:id/daftar`  
**Authorization:** Required (user_biasa)  
**Description:** Mendaftar ke sesi bimbingan perkawinan

**Request Body:** No body required

**Success Response (201):**
```json
{
  "message": "Berhasil mendaftar bimbingan perkawinan",
  "data": {
    "bimbingan_id": 1,
    "tanggal": "2024-12-18",
    "waktu": "08:00 - 12:00",
    "tempat": "Aula KUA Banjarmasin",
    "pembimbing": "H. Ahmad Dahlan, M.Ag"
  }
}
```

**Business Rules:**
- User harus punya pendaftaran dengan status "Menunggu Bimbingan"
- Belum terdaftar di bimbingan lain
- Masih ada kuota tersedia

**Error Responses:**
```json
// Belum siap bimbingan
{
  "error": "Anda belum memiliki pendaftaran nikah yang siap untuk bimbingan"
}

// Sudah terdaftar
{
  "error": "Anda sudah terdaftar di bimbingan perkawinan ini"
}

// Kuota penuh
{
  "error": "Bimbingan perkawinan sudah penuh"
}
```

---

### 6.5 Get Bimbingan Participants

**Endpoint:** `GET /simnikah/bimbingan/:id/participants`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Mendapatkan daftar peserta bimbingan

**Success Response (200):**
```json
{
  "message": "Data peserta bimbingan perkawinan berhasil diambil",
  "data": {
    "bimbingan_id": 1,
    "tanggal": "2024-12-18",
    "waktu": "08:00 - 12:00",
    "tempat": "Aula KUA Banjarmasin",
    "pembimbing": "H. Ahmad Dahlan, M.Ag",
    "kapasitas": 10,
    "jumlah_peserta": 7,
    "peserta": [
      {
        "id": 1,
        "pendaftaran_nikah_id": 1,
        "calon_suami": {
          "nama": "Ahmad Fauzi",
          "nik": "6371012501950001"
        },
        "calon_istri": {
          "nama": "Siti Nurhaliza",
          "nik": "6371014203970001"
        },
        "status_kehadiran": "Belum",
        "status_sertifikat": "Belum",
        "no_sertifikat": null,
        "created_at": "2024-11-25T10:00:00Z"
      }
    ]
  }
}
```

---

### 6.6 Update Bimbingan Attendance

**Endpoint:** `PUT /simnikah/bimbingan/:id/update-attendance`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Update kehadiran peserta bimbingan

**Request Body:**
```json
{
  "pendaftaran_nikah_id": 1,
  "status_kehadiran": "Hadir",
  "status_sertifikat": "Sudah",
  "no_sertifikat": "CERT-2024-001"
}
```

**Status Kehadiran Options:**
- `Hadir`
- `Tidak Hadir`

**Status Sertifikat Options:**
- `Belum`
- `Sudah`

**Success Response (200):**
```json
{
  "message": "Kehadiran bimbingan berhasil diupdate",
  "data": {
    "bimbingan_id": 1,
    "pendaftaran_id": 1,
    "status_kehadiran": "Hadir",
    "status_sertifikat": "Sudah",
    "no_sertifikat": "CERT-2024-001",
    "updated_by": "STF1700567890",
    "updated_at": "2024-12-18T12:30:00Z"
  }
}
```

---

### 6.7 Get Bimbingan Calendar

**Endpoint:** `GET /simnikah/bimbingan-kalender`  
**Authorization:** Required  
**Description:** Mendapatkan kalender bimbingan untuk bulan tertentu

**Query Parameters:**
- `bulan` (optional): Bulan (1-12)
- `tahun` (optional): Tahun (YYYY)

**Success Response (200):**
```json
{
  "message": "Kalender bimbingan perkawinan berhasil diambil",
  "data": {
    "bulan": 12,
    "tahun": 2024,
    "nama_bulan": "December",
    "kalender": [
      {
        "tanggal": 18,
        "tanggal_str": "2024-12-18",
        "status": "Tersedia",
        "tersedia": true,
        "sisa_kuota": 3,
        "bimbingan": {
          "id": 1,
          "waktu_mulai": "08:00",
          "waktu_selesai": "12:00",
          "tempat_bimbingan": "Aula KUA Banjarmasin",
          "pembimbing": "H. Ahmad Dahlan, M.Ag",
          "kapasitas": 10,
          "jumlah_peserta": 7
        }
      },
      {
        "tanggal": 25,
        "tanggal_str": "2024-12-25",
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

## 7. Notifications

### 7.1 Create Notification

**Endpoint:** `POST /simnikah/notifikasi`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Membuat notifikasi baru

**Request Body:**
```json
{
  "user_id": "USR1234567890",
  "judul": "Formulir Disetujui",
  "pesan": "Formulir pendaftaran nikah Anda telah disetujui. Silakan kumpulkan berkas ke KUA.",
  "tipe": "Success",
  "link": "/pendaftaran/1"
}
```

**Tipe Options:**
- `Info`
- `Warning`
- `Error`
- `Success`

**Success Response (201):**
```json
{
  "message": "Notifikasi berhasil dibuat",
  "notification": {
    "id": 1,
    "user_id": "USR1234567890",
    "judul": "Formulir Disetujui",
    "pesan": "Formulir pendaftaran nikah Anda telah disetujui.",
    "tipe": "Success",
    "status_baca": "Belum Dibaca",
    "link": "/pendaftaran/1",
    "created_at": "2024-11-21T10:00:00Z"
  }
}
```

---

### 7.2 Get User Notifications

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`  
**Authorization:** Required  
**Description:** Mendapatkan semua notifikasi user

**Query Parameters:**
- `page` (optional): Nomor halaman (default: 1)
- `limit` (optional): Data per halaman (default: 10, max: 100)
- `status` (optional): Filter status ("Belum Dibaca", "Sudah Dibaca")
- `tipe` (optional): Filter tipe ("Info", "Warning", "Error", "Success")

**Example Request:**
```
GET /simnikah/notifikasi/user/USR1234567890?page=1&limit=20&status=Belum Dibaca
```

**Success Response (200):**
```json
{
  "message": "Notifikasi berhasil diambil",
  "notifications": [
    {
      "id": 1,
      "user_id": "USR1234567890",
      "judul": "Formulir Disetujui",
      "pesan": "Formulir pendaftaran nikah Anda telah disetujui.",
      "tipe": "Success",
      "status_baca": "Belum Dibaca",
      "link": "/pendaftaran/1",
      "created_at": "2024-11-21T10:00:00Z",
      "updated_at": "2024-11-21T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 15,
    "total_page": 1
  },
  "unread_count": 5
}
```

---

### 7.3 Get Notification Detail

**Endpoint:** `GET /simnikah/notifikasi/:id`  
**Authorization:** Required  
**Description:** Mendapatkan detail notifikasi

**Success Response (200):**
```json
{
  "message": "Notifikasi berhasil diambil",
  "notification": {
    "id": 1,
    "user_id": "USR1234567890",
    "judul": "Formulir Disetujui",
    "pesan": "Formulir pendaftaran nikah Anda telah disetujui.",
    "tipe": "Success",
    "status_baca": "Belum Dibaca",
    "link": "/pendaftaran/1",
    "created_at": "2024-11-21T10:00:00Z",
    "updated_at": "2024-11-21T10:00:00Z"
  }
}
```

---

### 7.4 Update Notification Status

**Endpoint:** `PUT /simnikah/notifikasi/:id/status`  
**Authorization:** Required  
**Description:** Update status notifikasi (mark as read/unread)

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

**Success Response (200):**
```json
{
  "message": "Status notifikasi berhasil diupdate",
  "notification": {
    "id": 1,
    "status_baca": "Sudah Dibaca",
    "updated_at": "2024-11-21T15:00:00Z"
  }
}
```

---

### 7.5 Mark All as Read

**Endpoint:** `PUT /simnikah/notifikasi/user/:user_id/mark-all-read`  
**Authorization:** Required  
**Description:** Tandai semua notifikasi sebagai sudah dibaca

**Success Response (200):**
```json
{
  "message": "Semua notifikasi berhasil ditandai sebagai sudah dibaca",
  "updated_count": 5
}
```

---

### 7.6 Delete Notification

**Endpoint:** `DELETE /simnikah/notifikasi/:id`  
**Authorization:** Required  
**Description:** Menghapus notifikasi

**Success Response (200):**
```json
{
  "message": "Notifikasi berhasil dihapus"
}
```

---

### 7.7 Get Notification Stats

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id/stats`  
**Authorization:** Required  
**Description:** Mendapatkan statistik notifikasi user

**Success Response (200):**
```json
{
  "message": "Statistik notifikasi berhasil diambil",
  "stats": {
    "total": 25,
    "unread": 5,
    "read": 20,
    "by_type": {
      "info": 10,
      "warning": 5,
      "error": 3,
      "success": 7
    },
    "today": 3,
    "week": 12
  }
}
```

---

### 7.8 Send Notification to Role

**Endpoint:** `POST /simnikah/notifikasi/send-to-role`  
**Authorization:** Required (staff, kepala_kua)  
**Description:** Mengirim notifikasi ke semua user dengan role tertentu

**Request Body:**
```json
{
  "role": "user_biasa",
  "judul": "Pengumuman Penting",
  "pesan": "KUA akan tutup pada tanggal 1-3 Januari 2025 untuk libur tahun baru.",
  "tipe": "Info",
  "link": "/pengumuman"
}
```

**Success Response (201):**
```json
{
  "message": "Notifikasi berhasil dikirim ke 45 user dengan role user_biasa",
  "recipient_count": 45,
  "role": "user_biasa"
}
```

---

## 8. Location & Map Integration

### 8.1 Geocoding (Address to Coordinates)

**Endpoint:** `POST /simnikah/location/geocode`  
**Authorization:** Required  
**Description:** Mendapatkan koordinat dari alamat (100% GRATIS)

**Request Body:**
```json
{
  "alamat": "Jl. Merdeka No. 123, Banjarmasin, Kalimantan Selatan"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Koordinat berhasil ditemukan",
  "data": {
    "alamat": "Jl. Merdeka No. 123, Banjarmasin, Kalimantan Selatan",
    "latitude": -3.3194,
    "longitude": 114.5903,
    "map_url": "https://www.google.com/maps?q=-3.319400,114.590300",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.319400&mlon=114.590300&zoom=16"
  }
}
```

**Features:**
- Menggunakan OpenStreetMap Nominatim API (GRATIS)
- Caching untuk performa
- Validasi alamat minimal 10 karakter

---

### 8.2 Reverse Geocoding (Coordinates to Address)

**Endpoint:** `POST /simnikah/location/reverse-geocode`  
**Authorization:** Required  
**Description:** Mendapatkan alamat dari koordinat

**Request Body:**
```json
{
  "latitude": -3.3194,
  "longitude": 114.5903
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Alamat berhasil ditemukan",
  "data": {
    "latitude": -3.3194,
    "longitude": 114.5903,
    "alamat": "Jalan Merdeka, Kelurahan Kertak Baru Ilir, Kecamatan Banjarmasin Tengah, Kota Banjarmasin, Kalimantan Selatan, 70111, Indonesia",
    "detail": {
      "road": "Jalan Merdeka",
      "suburb": "Kelurahan Kertak Baru Ilir",
      "city": "Kota Banjarmasin",
      "state": "Kalimantan Selatan",
      "postcode": "70111",
      "country": "Indonesia"
    },
    "map_url": "https://www.google.com/maps?q=-3.319400,114.590300",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.319400&mlon=114.590300&zoom=16"
  }
}
```

---

### 8.3 Search Address (Autocomplete)

**Endpoint:** `GET /simnikah/location/search`  
**Authorization:** Required  
**Description:** Mencari alamat untuk fitur autocomplete

**Query Parameters:**
- `q` (required): Query pencarian (minimal 3 karakter)

**Example Request:**
```
GET /simnikah/location/search?q=Jl. Merdeka Banjarmasin
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Hasil pencarian alamat",
  "data": {
    "query": "Jl. Merdeka Banjarmasin",
    "results": [
      {
        "display_name": "Jalan Merdeka, Kelurahan Kertak Baru Ilir, Kecamatan Banjarmasin Tengah, Kota Banjarmasin, Kalimantan Selatan, 70111, Indonesia",
        "latitude": "-3.3194",
        "longitude": "114.5903",
        "address": {
          "road": "Jalan Merdeka",
          "suburb": "Kelurahan Kertak Baru Ilir",
          "city": "Kota Banjarmasin",
          "state": "Kalimantan Selatan"
        }
      }
    ],
    "count": 1
  }
}
```

**Features:**
- Filter khusus Indonesia (countrycodes=id)
- Limit 5 hasil teratas
- Detail alamat lengkap

---

### 8.4 Update Wedding Location with Coordinates

**Endpoint:** `PUT /simnikah/pendaftaran/:id/location`  
**Authorization:** Required (user_biasa)  
**Description:** Update lokasi nikah beserta koordinatnya

**Request Body:**
```json
{
  "alamat_akad": "Jl. Merdeka No. 123, RT 05, RW 02, Banjarmasin",
  "latitude": -3.3194,
  "longitude": 114.5903
}
```

**Notes:**
- Jika koordinat tidak disediakan, sistem akan auto-geocoding
- Hanya untuk nikah "Di Luar KUA"

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lokasi nikah berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "alamat_akad": "Jl. Merdeka No. 123, RT 05, RW 02, Banjarmasin",
    "tempat_nikah": "Di Luar KUA",
    "latitude": -3.3194,
    "longitude": 114.5903,
    "map_url": "https://www.google.com/maps?q=-3.319400,114.590300",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.319400&mlon=114.590300&zoom=16",
    "updated_at": "2024-11-21T16:00:00Z"
  }
}
```

---

### 8.5 Get Wedding Location Detail

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`  
**Authorization:** Required  
**Description:** Mendapatkan detail lokasi nikah (untuk penghulu)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Detail lokasi nikah berhasil diambil",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "REG20241120001",
    "tanggal_nikah": "2024-12-25T00:00:00Z",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Merdeka No. 123, RT 05, RW 02, Banjarmasin",
    "has_coordinates": true,
    "latitude": -3.3194,
    "longitude": 114.5903,
    "map_url": "https://www.google.com/maps?q=-3.319400,114.590300",
    "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.319400,114.590300",
    "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.319400,114.590300",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.319400&mlon=114.590300&zoom=16",
    "waze_url": "https://www.waze.com/ul?ll=-3.319400,114.590300&navigate=yes",
    "is_outside_kua": true,
    "note": "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
  }
}
```

**Navigation Links:**
- `google_maps_url` - Lihat di Google Maps
- `google_maps_directions_url` - Dapatkan arah ke lokasi
- `waze_url` - Navigasi dengan Waze
- `osm_url` - Lihat di OpenStreetMap

---

## Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request berhasil |
| 201 | Created - Resource berhasil dibuat |
| 400 | Bad Request - Format request tidak valid |
| 401 | Unauthorized - Token tidak valid/expired |
| 403 | Forbidden - Tidak punya akses |
| 404 | Not Found - Resource tidak ditemukan |
| 500 | Internal Server Error - Error server |

---

## Error Handling

### Standard Error Response Format

```json
{
  "error": "Deskripsi error"
}
```

### Detailed Error Response

```json
{
  "success": false,
  "message": "Judul error",
  "error": "Detail error",
  "details": {
    "field": "Informasi tambahan"
  }
}
```

### Common Errors

**Authentication Errors:**
```json
{
  "error": "Token otorisasi tidak disediakan"
}

{
  "error": "Token tidak valid atau kedaluwarsa"
}
```

**Authorization Errors:**
```json
{
  "error": "Akses ditolak. Role kepala_kua diperlukan"
}

{
  "error": "Akses ditolak. Role yang diizinkan: staff, kepala_kua"
}
```

**Validation Errors:**
```json
{
  "error": "Format data tidak valid"
}

{
  "error": "NIK sudah terdaftar"
}
```

---

## Workflow Status Pendaftaran Nikah

### Status Flow Lengkap:

```
1. Draft
   ↓
2. Menunggu Verifikasi (Staff verify form online)
   ↓
3. Menunggu Pengumpulan Berkas (User datang ke KUA)
   ↓
4. Berkas Diterima (Staff verify berkas fisik)
   ↓
5. Menunggu Penugasan (Ready untuk assign penghulu)
   ↓
6. Penghulu Ditugaskan (Kepala KUA assign penghulu)
   ↓
7. Menunggu Verifikasi Penghulu (Penghulu cek berkas)
   ↓
8. Menunggu Bimbingan (Siap daftar bimbingan)
   ↓
9. Sudah Bimbingan (Bimbingan selesai)
   ↓
10. Selesai (Nikah telah dilaksanakan)
```

### Transisi Status:

| From | To | Actor | Action |
|------|-----|-------|--------|
| Draft | Menunggu Verifikasi | User | Submit form |
| Menunggu Verifikasi | Menunggu Pengumpulan Berkas | Staff | Approve form |
| Menunggu Pengumpulan Berkas | Berkas Diterima | Staff | Verify documents |
| Berkas Diterima | Menunggu Penugasan | Staff | Mark as visited |
| Menunggu Penugasan | Penghulu Ditugaskan | Kepala KUA | Assign penghulu |
| Penghulu Ditugaskan | Menunggu Bimbingan | Penghulu | Verify documents |
| Menunggu Bimbingan | Sudah Bimbingan | User | Complete bimbingan |
| Sudah Bimbingan | Selesai | Staff/Kepala KUA | Complete nikah |

---

## Rate Limiting

### Global Rate Limit
- **100 requests per minute** per IP address

### Strict Rate Limit (Auth Endpoints)
- **10 requests per minute** untuk `/register` dan `/login`

### Rate Limit Headers
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1700567890
```

### Rate Limit Exceeded Response (429)
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

---

## Testing & Development

### Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2024-11-21T10:00:00Z"
}
```

### Example cURL Commands

**Register:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "nama": "John Doe",
    "role": "user_biasa"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Authenticated Request:**
```bash
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Changelog

### Version 1.0.0 (Current)
- ✅ Authentication & Authorization
- ✅ Marriage Registration System
- ✅ Staff Management
- ✅ Penghulu Management
- ✅ Calendar & Scheduling
- ✅ Bimbingan Perkawinan
- ✅ Notification System
- ✅ Location & Map Integration (100% FREE)
- ✅ Rate Limiting
- ✅ Graceful Shutdown
- ✅ Database Indexes

---

## Support & Contact

**Documentation Issues:**
- Email: support@simnikah.id
- GitHub: https://github.com/your-org/simpadu

**Production URL:**
- Railway: https://your-app.railway.app

**Development URL:**
- Local: http://localhost:8080

---

**© 2024 SimNikah - Sistem Manajemen Pendaftaran Nikah KUA**

Made with ❤️ for Indonesian KUA (Kantor Urusan Agama)

