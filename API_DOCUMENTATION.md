# 📚 SimNikah API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Semua endpoint (kecuali `/register`, `/login`, `/health`) memerlukan JWT token di header:
```
Authorization: Bearer <jwt_token>
```

---

## 🔐 Authentication Endpoints

### 1. Register User
**POST** `/register`

Mendaftarkan user baru ke sistem.

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

**Response:**
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

**Valid Roles:**
- `user_biasa` - User biasa untuk daftar nikah
- `penghulu` - Penghulu untuk memimpin nikah
- `staff` - Staff KUA untuk verifikasi
- `kepala_kua` - Kepala KUA untuk approval

---

### 2. Login User
**POST** `/login`

Login user dan mendapatkan JWT token.

**Request Body:**
```json
{
  "username": "ahmad123",
  "password": "password123"
}
```

**Response:**
```json
{
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

---

### 3. Get Profile
**GET** `/profile`

Mendapatkan informasi profile user yang sedang login.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
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

---

## 💒 Marriage Registration Endpoints

### 4. Check Registration Status
**GET** `/simnikah/pendaftaran/status`

Cek apakah user sudah memiliki pendaftaran nikah.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (Belum ada pendaftaran):**
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

**Response (Sudah ada pendaftaran):**
```json
{
  "success": true,
  "message": "User sudah memiliki pendaftaran nikah",
  "data": {
    "has_registration": true,
    "can_register": false,
    "registration": {
      "id": 1,
      "nomor_pendaftaran": "NIK1704067200",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2024-02-14T00:00:00Z",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "KUA Banjarmasin Utara",
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

---

### 5. Create Marriage Registration
**POST** `/simnikah/pendaftaran/form-baru`

Membuat pendaftaran nikah baru dengan form lengkap.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "scheduleAndLocation": {
    "weddingLocation": "Di KUA",
    "weddingAddress": "",
    "weddingDate": "2024-02-14",
    "weddingTime": "09:00",
    "dispensationNumber": ""
  },
  "groom": {
    "groomFullName": "Ahmad Wijaya",
    "groomNik": "3201010101010001",
    "groomCitizenship": "WNI",
    "groomPassportNumber": "",
    "groomPlaceOfBirth": "Bandung",
    "groomDateOfBirth": "1995-01-02",
    "groomStatus": "Belum Kawin",
    "groomReligion": "Islam",
    "groomEducation": "SMA",
    "groomOccupation": "Karyawan",
    "groomOccupationDescription": "",
    "groomPhoneNumber": "081234567890",
    "groomEmail": "ahmad@example.com",
    "groomAddress": "Jl. Melati No. 1"
  },
  "bride": {
    "brideFullName": "Aisyah Sari",
    "brideNik": "3201010101010002",
    "brideCitizenship": "WNI",
    "bridePassportNumber": "",
    "bridePlaceOfBirth": "Bandung",
    "brideDateOfBirth": "1996-02-03",
    "brideStatus": "Belum Kawin",
    "brideReligion": "Islam",
    "brideEducation": "SMA",
    "brideOccupation": "Karyawan",
    "brideOccupationDescription": "",
    "bridePhoneNumber": "081298765432",
    "brideEmail": "aisyah@example.com",
    "brideAddress": "Jl. Anggrek No. 2"
  },
  "groomParents": {
    "groomFather": {
      "groomFatherPresenceStatus": "Hidup",
      "groomFatherName": "Bapak Ahmad",
      "groomFatherNik": "3201010101010003",
      "groomFatherCitizenship": "WNI",
      "groomFatherCountryOfOrigin": "",
      "groomFatherPassportNumber": "",
      "groomFatherPlaceOfBirth": "Bandung",
      "groomFatherDateOfBirth": "1970-05-01",
      "groomFatherReligion": "Islam",
      "groomFatherOccupation": "Wiraswasta",
      "groomFatherOccupationDescription": "",
      "groomFatherAddress": "Jl. Melati No. 1"
    },
    "groomMother": {
      "groomMotherPresenceStatus": "Hidup",
      "groomMotherName": "Ibu Ahmad",
      "groomMotherNik": "3201010101010004",
      "groomMotherCitizenship": "WNI",
      "groomMotherCountryOfOrigin": "",
      "groomMotherPassportNumber": "",
      "groomMotherPlaceOfBirth": "Bandung",
      "groomMotherDateOfBirth": "1972-07-01",
      "groomMotherReligion": "Islam",
      "groomMotherOccupation": "Ibu Rumah Tangga",
      "groomMotherOccupationDescription": "",
      "groomMotherAddress": "Jl. Melati No. 1"
    }
  },
  "brideParents": {
    "brideFather": {
      "brideFatherPresenceStatus": "Hidup",
      "brideFatherName": "Bapak Aisyah",
      "brideFatherNik": "3201010101010005",
      "brideFatherCitizenship": "WNI",
      "brideFatherCountryOfOrigin": "",
      "brideFatherPassportNumber": "",
      "brideFatherPlaceOfBirth": "Bandung",
      "brideFatherDateOfBirth": "1971-03-02",
      "brideFatherReligion": "Islam",
      "brideFatherOccupation": "PNS",
      "brideFatherOccupationDescription": "",
      "brideFatherAddress": "Jl. Anggrek No. 2"
    },
    "brideMother": {
      "brideMotherPresenceStatus": "Hidup",
      "brideMotherName": "Ibu Aisyah",
      "brideMotherNik": "3201010101010006",
      "brideMotherCitizenship": "WNI",
      "brideMotherCountryOfOrigin": "",
      "brideMotherPassportNumber": "",
      "brideMotherPlaceOfBirth": "Bandung",
      "brideMotherDateOfBirth": "1973-08-03",
      "brideMotherReligion": "Islam",
      "brideMotherOccupation": "Ibu Rumah Tangga",
      "brideMotherOccupationDescription": "",
      "brideMotherAddress": "Jl. Anggrek No. 2"
    }
  },
  "guardian": {
    "guardianFullName": "Wali Nikah",
    "guardianNik": "3201010101010007",
    "guardianRelationship": "Ayah",
    "guardianStatus": "Hidup",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. Anggrek No. 2",
    "guardianPhoneNumber": "0812000000"
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat",
  "data": {
    "nomor_pendaftaran": "NIK1704067200",
    "pendaftaran_id": 1,
    "calon_suami": {
      "id": 1,
      "user_id": "abc123def456ghi789",
      "nik": "3201010101010001",
      "nama_lengkap": "Ahmad Wijaya"
    },
    "calon_istri": {
      "id": 2,
      "user_id": "xyz789uvw456rst123",
      "nik": "3201010101010002",
      "nama_lengkap": "Aisyah Sari"
    },
    "pendaftaran": {
      "id": 1,
      "nomor_pendaftaran": "NIK1704067200",
      "status_pendaftaran": "Menunggu Verifikasi",
      "tanggal_nikah": "2024-02-14T00:00:00Z",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di KUA"
    }
  }
}
```

---

### 6. Mark as Visited
**POST** `/simnikah/pendaftaran/:id/mark-visited`

Menandai bahwa calon pasangan telah datang ke kantor dengan berkas.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Status berhasil diupdate",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_pendaftaran": "Menunggu Verifikasi",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 7. Update Wedding Address
**PUT** `/simnikah/pendaftaran/:id/alamat`

Mengupdate alamat nikah untuk nikah di luar KUA.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "alamat_akad": "Jl. Merdeka No. 123, Jakarta Selatan"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Alamat nikah berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "alamat_akad": "Jl. Merdeka No. 123, Jakarta Selatan",
    "tempat_nikah": "Di Luar KUA",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 👥 Staff Management Endpoints

### 8. Create Staff KUA
**POST** `/simnikah/staff`

Membuat staff KUA baru (hanya kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "username": "staff001",
  "email": "staff@kua.go.id",
  "password": "password123",
  "nama": "Staff KUA 1",
  "nip": "198012012000031001",
  "jabatan": "Staff",
  "bagian": "Pendaftaran",
  "no_hp": "081234567890",
  "alamat": "Jl. KUA No. 1"
}
```

**Response:**
```json
{
  "message": "Staff KUA berhasil dibuat",
  "data": {
    "user_id": "USR1704067201",
    "username": "staff001",
    "email": "staff@kua.go.id",
    "nama": "Staff KUA 1",
    "nip": "198012012000031001",
    "jabatan": "Staff",
    "bagian": "Pendaftaran"
  }
}
```

---

### 9. Get All Staff
**GET** `/simnikah/staff`

Mendapatkan daftar semua staff KUA.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Data staff berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067201",
      "nip": "198012012000031001",
      "nama_lengkap": "Staff KUA 1",
      "jabatan": "Staff",
      "bagian": "Pendaftaran",
      "no_hp": "081234567890",
      "email": "staff@kua.go.id",
      "alamat": "Jl. KUA No. 1",
      "status": "Aktif"
    }
  ]
}
```

---

### 10. Update Staff KUA
**PUT** `/simnikah/staff/:id`

Mengupdate data staff KUA.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "nama_lengkap": "Staff KUA Updated",
  "jabatan": "Senior Staff",
  "bagian": "Verifikasi",
  "no_hp": "081234567891",
  "email": "staff.updated@kua.go.id",
  "alamat": "Jl. KUA No. 2"
}
```

**Response:**
```json
{
  "message": "Staff KUA berhasil diupdate",
  "data": {
    "id": 1,
    "nama_lengkap": "Staff KUA Updated",
    "jabatan": "Senior Staff",
    "bagian": "Verifikasi",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 🕌 Penghulu Management Endpoints

### 11. Create Penghulu
**POST** `/simnikah/penghulu`

Membuat penghulu baru (hanya kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "username": "penghulu001",
  "email": "penghulu@kua.go.id",
  "password": "password123",
  "nama": "Penghulu 1",
  "nip": "198012012000031002",
  "no_hp": "081234567890",
  "alamat": "Jl. Penghulu No. 1"
}
```

**Response:**
```json
{
  "message": "Penghulu berhasil dibuat",
  "data": {
    "user_id": "USR1704067202",
    "username": "penghulu001",
    "email": "penghulu@kua.go.id",
    "nama": "Penghulu 1",
    "nip": "198012012000031002"
  }
}
```

---

### 12. Get All Penghulu
**GET** `/simnikah/penghulu`

Mendapatkan daftar semua penghulu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Data penghulu berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067202",
      "nip": "198012012000031002",
      "nama_lengkap": "Penghulu 1",
      "no_hp": "081234567890",
      "email": "penghulu@kua.go.id",
      "alamat": "Jl. Penghulu No. 1",
      "status": "Aktif",
      "jumlah_nikah": 0,
      "rating": 0
    }
  ]
}
```

---

### 13. Update Penghulu
**PUT** `/simnikah/penghulu/:id`

Mengupdate data penghulu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "nama_lengkap": "Penghulu Updated",
  "no_hp": "081234567891",
  "email": "penghulu.updated@kua.go.id",
  "alamat": "Jl. Penghulu No. 2"
}
```

**Response:**
```json
{
  "message": "Penghulu berhasil diupdate",
  "data": {
    "id": 1,
    "nama_lengkap": "Penghulu Updated",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## ✅ Verification Endpoints

### 14. Verify Formulir
**POST** `/simnikah/staff/verify-formulir/:id`

Verifikasi formulir pendaftaran (hanya staff).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir lengkap dan valid"
}
```

**Valid Status:**
- `"Formulir Disetujui"`
- `"Formulir Ditolak"`

**Response:**
```json
{
  "message": "Status formulir berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "status_sebelum": "Menunggu Verifikasi",
    "status_sesudah": "Menunggu Pengumpulan Berkas",
    "updated_by": "USR1704067201",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 15. Verify Berkas
**POST** `/simnikah/staff/verify-berkas/:id`

Verifikasi berkas fisik (hanya staff).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "status": "Berkas Diterima",
  "catatan": "Semua berkas sudah lengkap (KK, KTP, dll)"
}
```

**Valid Status:**
- `"Berkas Diterima"`
- `"Berkas Ditolak"`

**Response:**
```json
{
  "message": "Status berkas berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "status_sebelum": "Menunggu Pengumpulan Berkas",
    "status_sesudah": "Menunggu Penugasan",
    "updated_by": "USR1704067201",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 16. Verify Documents (Penghulu)
**POST** `/simnikah/penghulu/verify-documents/:id`

Verifikasi dokumen oleh penghulu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "status": "Menunggu Pelaksanaan",
  "catatan": "Dokumen lengkap dan siap dilanjutkan"
}
```

**Valid Status:**
- `"Menunggu Pelaksanaan"`
- `"Ditolak"`

**Response:**
```json
{
  "message": "Status dokumen berhasil diupdate",
  "data": {
    "pendaftaran_id": 1,
    "status_sebelum": "Menunggu Verifikasi Penghulu",
    "status_sesudah": "Menunggu Bimbingan",
    "updated_by": "USR1704067202",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 📅 Calendar & Scheduling Endpoints

### 17. Get Calendar Availability
**GET** `/simnikah/kalender-ketersediaan`

Mendapatkan kalender ketersediaan untuk pendaftaran nikah.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `bulan` (optional): Bulan (01-12), default: bulan saat ini
- `tahun` (optional): Tahun, default: tahun saat ini

**Example:**
```
GET /simnikah/kalender-ketersediaan?bulan=02&tahun=2024
```

**Response:**
```json
{
  "message": "Kalender ketersediaan berhasil diambil",
  "data": {
    "bulan": 2,
    "tahun": 2024,
    "nama_bulan": "February",
    "kapasitas_harian": 9,
    "penghulu_info": {
      "total_penghulu": 4,
      "penghulu_aktif": 3,
      "penghulu_cadangan": 1,
      "nikah_per_penghulu": 3
    },
    "kalender": [
      {
        "tanggal": 1,
        "tanggal_str": "2024-02-01",
        "status": "Tersedia",
        "tersedia": true,
        "jumlah_nikah": 2,
        "kuning_count": 1,
        "hijau_count": 1,
        "warna": "hijau",
        "sisa_kuota": 7,
        "kapasitas": 9
      }
    ]
  }
}
```

---

### 18. Get Calendar Date Detail
**GET** `/simnikah/kalender-tanggal-detail`

Mendapatkan detail jadwal per tanggal.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `tanggal` (required): Tanggal dalam format YYYY-MM-DD

**Example:**
```
GET /simnikah/kalender-tanggal-detail?tanggal=2024-02-14
```

**Response:**
```json
{
  "message": "Detail kalender berhasil diambil",
  "data": {
    "tanggal": "2024-02-14",
    "items": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIK1704067200",
        "waktu": "09:00",
        "tempat": "Di KUA",
        "status": "Menunggu Verifikasi",
        "warna": "kuning",
        "penghulu_id": 1,
        "penghulu_nama": "Penghulu 1",
        "nama_calon_suami": "Ahmad Wijaya",
        "nama_calon_istri": "Aisyah Sari"
      }
    ]
  }
}
```

---

### 19. Get Date Availability
**GET** `/simnikah/ketersediaan-tanggal/:tanggal`

Mendapatkan detail ketersediaan untuk tanggal tertentu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example:**
```
GET /simnikah/ketersediaan-tanggal/2024-02-14
```

**Response:**
```json
{
  "message": "Detail ketersediaan tanggal berhasil diambil",
  "data": {
    "tanggal": "2024-02-14",
    "status": "Tersedia",
    "tersedia": true,
    "jumlah_nikah": 2,
    "sisa_kuota": 7,
    "kapasitas": 9,
    "jadwal_detail": [
      {
        "nomor_pendaftaran": "NIK1704067200",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "KUA Banjarmasin Utara"
      }
    ]
  }
}
```

---

### 20. Get Penghulu Schedule
**GET** `/simnikah/penghulu-jadwal/:tanggal`

Mendapatkan jadwal penghulu untuk tanggal tertentu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example:**
```
GET /simnikah/penghulu-jadwal/2024-02-14
```

**Response:**
```json
{
  "message": "Jadwal penghulu berhasil diambil",
  "data": {
    "tanggal": "2024-02-14",
    "total_penghulu": 3,
    "total_kapasitas": 9,
    "total_terisi": 2,
    "total_sisa": 7,
    "penghulu": [
      {
        "id": 1,
        "nama": "Penghulu 1",
        "status": "Sebagian",
        "jumlah_jadwal": 1,
        "sisa_kuota": 2,
        "maksimal": 3,
        "jadwal": [
          {
            "nomor_pendaftaran": "NIK1704067200",
            "waktu_nikah": "09:00",
            "tempat_nikah": "Di KUA",
            "alamat_akad": "KUA Banjarmasin Utara"
          }
        ]
      }
    ]
  }
}
```

---

## 👨‍💼 Penghulu Assignment Endpoints

### 21. Assign Penghulu
**POST** `/simnikah/pendaftaran/:id/assign-penghulu`

Menugaskan penghulu untuk pendaftaran nikah (hanya kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "penghulu_id": 1
}
```

**Response:**
```json
{
  "message": "Penghulu berhasil diassign",
  "warning": "Peringatan: penghulu ini memiliki 2 jadwal pada tanggal 2024-02-14",
  "data": {
    "pendaftaran_id": 1,
    "penghulu_id": 1,
    "penghulu_nama": "Penghulu 1",
    "assigned_by": "USR1704067203",
    "assigned_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 22. Change Penghulu
**PUT** `/simnikah/pendaftaran/:id/change-penghulu`

Mengubah penghulu untuk pendaftaran nikah (hanya kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "penghulu_id": 2
}
```

**Response:**
```json
{
  "message": "Penghulu berhasil diubah",
  "data": {
    "pendaftaran_id": 1,
    "penghulu_lama": 1,
    "penghulu_baru": 2,
    "penghulu_nama": "Penghulu 2",
    "changed_by": "USR1704067203",
    "changed_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 23. Get Unassigned Registrations
**GET** `/simnikah/pendaftaran/belum-assign-penghulu`

Mendapatkan pendaftaran yang belum di-assign penghulu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Pendaftaran belum assign penghulu berhasil diambil",
  "data": {
    "total": 5,
    "pendaftaran": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIK1704067200",
        "tanggal_nikah": "2024-02-14",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "KUA Banjarmasin Utara",
        "status_pendaftaran": "Menunggu Penugasan",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 24. Get Penghulu Availability
**GET** `/simnikah/penghulu/:id/ketersediaan/:tanggal`

Mendapatkan ketersediaan waktu penghulu untuk tanggal tertentu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example:**
```
GET /simnikah/penghulu/1/ketersediaan/2024-02-14
```

**Response:**
```json
{
  "message": "Ketersediaan penghulu berhasil diambil",
  "data": {
    "penghulu": {
      "id": 1,
      "nama": "Penghulu 1",
      "status": "Aktif"
    },
    "tanggal": "2024-02-14",
    "statistik": {
      "jumlah_jadwal": 1,
      "sisa_kuota": 2,
      "maksimal_per_hari": 3,
      "slot_tersedia": 4,
      "total_slot": 5
    },
    "jadwal_hari_ini": [
      {
        "nomor_pendaftaran": "NIK1704067200",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA"
      }
    ],
    "slot_waktu": [
      {
        "waktu": "08:00",
        "tersedia": true,
        "konflik_jadwal": []
      },
      {
        "waktu": "10:00",
        "tersedia": false,
        "konflik_jadwal": [
          {
            "waktu": "09:00",
            "tempat": "Di KUA",
            "selisih_menit": 60
          }
        ]
      }
    ]
  }
}
```

---

## 📚 Marriage Counseling Endpoints

### 25. Create Marriage Counseling
**POST** `/simnikah/bimbingan`

Membuat sesi bimbingan perkawinan baru (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "tanggal_bimbingan": "2024-02-07",
  "waktu_mulai": "08:00",
  "waktu_selesai": "12:00",
  "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
  "pembimbing": "Ustadz Ahmad",
  "kapasitas": 10,
  "catatan": "Bimbingan untuk calon pasangan yang akan menikah"
}
```

**Response:**
```json
{
  "message": "Bimbingan perkawinan berhasil dibuat",
  "data": {
    "id": 1,
    "tanggal_bimbingan": "2024-02-07T00:00:00Z",
    "waktu_mulai": "08:00",
    "waktu_selesai": "12:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
    "pembimbing": "Ustadz Ahmad",
    "kapasitas": 10,
    "status": "Aktif",
    "catatan": "Bimbingan untuk calon pasangan yang akan menikah"
  }
}
```

---

### 26. Get Marriage Counseling
**GET** `/simnikah/bimbingan`

Mendapatkan daftar bimbingan perkawinan.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `bulan` (optional): Bulan (01-12)
- `tahun` (optional): Tahun
- `status` (optional): Status bimbingan

**Example:**
```
GET /simnikah/bimbingan?bulan=02&tahun=2024&status=Aktif
```

**Response:**
```json
{
  "message": "Data bimbingan perkawinan berhasil diambil",
  "data": {
    "bulan": 2,
    "tahun": 2024,
    "bimbingan": [
      {
        "id": 1,
        "tanggal_bimbingan": "2024-02-07",
        "waktu_mulai": "08:00",
        "waktu_selesai": "12:00",
        "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
        "pembimbing": "Ustadz Ahmad",
        "kapasitas": 10,
        "jumlah_peserta": 5,
        "sisa_kuota": 5,
        "status": "Aktif",
        "catatan": "Bimbingan untuk calon pasangan yang akan menikah"
      }
    ]
  }
}
```

---

### 27. Get Marriage Counseling by ID
**GET** `/simnikah/bimbingan/:id`

Mendapatkan detail bimbingan perkawinan.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Detail bimbingan perkawinan berhasil diambil",
  "data": {
    "id": 1,
    "tanggal_bimbingan": "2024-02-07",
    "waktu_mulai": "08:00",
    "waktu_selesai": "12:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
    "pembimbing": "Ustadz Ahmad",
    "kapasitas": 10,
    "jumlah_peserta": 5,
    "sisa_kuota": 5,
    "status": "Aktif",
    "catatan": "Bimbingan untuk calon pasangan yang akan menikah"
  }
}
```

---

### 28. Update Marriage Counseling
**PUT** `/simnikah/bimbingan/:id`

Mengupdate bimbingan perkawinan (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "waktu_mulai": "09:00",
  "waktu_selesai": "13:00",
  "tempat_bimbingan": "Aula KUA Banjarmasin Utara (Updated)",
  "pembimbing": "Ustadz Ahmad (Updated)",
  "kapasitas": 15,
  "status": "Aktif",
  "catatan": "Bimbingan yang telah diupdate"
}
```

**Response:**
```json
{
  "message": "Bimbingan perkawinan berhasil diupdate",
  "data": {
    "id": 1,
    "waktu_mulai": "09:00",
    "waktu_selesai": "13:00",
    "tempat_bimbingan": "Aula KUA Banjarmasin Utara (Updated)",
    "pembimbing": "Ustadz Ahmad (Updated)",
    "kapasitas": 15,
    "status": "Aktif",
    "catatan": "Bimbingan yang telah diupdate",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 29. Get Marriage Counseling Calendar
**GET** `/simnikah/bimbingan-kalender`

Mendapatkan kalender bimbingan perkawinan.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `bulan` (optional): Bulan (01-12)
- `tahun` (optional): Tahun

**Example:**
```
GET /simnikah/bimbingan-kalender?bulan=02&tahun=2024
```

**Response:**
```json
{
  "message": "Kalender bimbingan perkawinan berhasil diambil",
  "data": {
    "bulan": 2,
    "tahun": 2024,
    "nama_bulan": "February",
    "kalender": [
      {
        "tanggal": 7,
        "tanggal_str": "2024-02-07",
        "status": "Tersedia",
        "tersedia": true,
        "sisa_kuota": 5,
        "bimbingan": {
          "id": 1,
          "waktu_mulai": "08:00",
          "waktu_selesai": "12:00",
          "tempat_bimbingan": "Aula KUA Banjarmasin Utara",
          "pembimbing": "Ustadz Ahmad",
          "kapasitas": 10,
          "jumlah_peserta": 5
        }
      }
    ]
  }
}
```

---

### 30. Register for Marriage Counseling
**POST** `/simnikah/bimbingan/:id/daftar`

Mendaftarkan calon pengantin ke bimbingan perkawinan.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Berhasil mendaftar bimbingan perkawinan",
  "data": {
    "bimbingan_id": 1,
    "tanggal": "2024-02-07",
    "waktu": "08:00 - 12:00",
    "tempat": "Aula KUA Banjarmasin Utara",
    "pembimbing": "Ustadz Ahmad"
  }
}
```

---

### 31. Get Marriage Counseling Participants
**GET** `/simnikah/bimbingan/:id/participants`

Mendapatkan daftar peserta bimbingan perkawinan (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Data peserta bimbingan perkawinan berhasil diambil",
  "data": {
    "bimbingan_id": 1,
    "tanggal": "2024-02-07",
    "waktu": "08:00 - 12:00",
    "tempat": "Aula KUA Banjarmasin Utara",
    "pembimbing": "Ustadz Ahmad",
    "kapasitas": 10,
    "jumlah_peserta": 5,
    "peserta": [
      {
        "id": 1,
        "pendaftaran_nikah_id": 1,
        "calon_suami": {
          "nama": "Ahmad Wijaya",
          "nik": "3201010101010001"
        },
        "calon_istri": {
          "nama": "Aisyah Sari",
          "nik": "3201010101010002"
        },
        "status_kehadiran": "Belum",
        "status_sertifikat": "Belum",
        "no_sertifikat": "",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 32. Update Marriage Counseling Attendance
**PUT** `/simnikah/bimbingan/:id/update-attendance`

Mengupdate kehadiran bimbingan perkawinan (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "pendaftaran_nikah_id": 1,
  "status_kehadiran": "Hadir",
  "status_sertifikat": "Sudah",
  "no_sertifikat": "SERT-001-2024"
}
```

**Valid Status:**
- `status_kehadiran`: `"Hadir"`, `"Tidak Hadir"`
- `status_sertifikat`: `"Belum"`, `"Sudah"`

**Response:**
```json
{
  "message": "Kehadiran bimbingan berhasil diupdate",
  "data": {
    "bimbingan_id": 1,
    "pendaftaran_id": 1,
    "status_kehadiran": "Hadir",
    "status_sertifikat": "Sudah",
    "no_sertifikat": "SERT-001-2024",
    "updated_by": "USR1704067201",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 📋 Status Management Endpoints

### 33. Get Status Flow
**GET** `/simnikah/pendaftaran/:id/status-flow`

Mendapatkan alur status pendaftaran nikah.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Alur status pendaftaran berhasil diambil",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_sekarang": "Menunggu Verifikasi",
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "penghulu_assigned": false,
    "bimbingan_info": {
      "terdaftar": false,
      "message": "Belum terdaftar bimbingan perkawinan"
    },
    "status_flow": [
      {
        "status": "Draft",
        "description": "Data belum lengkap",
        "completed": true,
        "current": false,
        "can_edit": false
      },
      {
        "status": "Menunggu Verifikasi",
        "description": "Staff verifikasi formulir online",
        "completed": false,
        "current": true,
        "can_edit": false
      }
    ]
  }
}
```

---

### 34. Complete Marriage Counseling
**PUT** `/simnikah/pendaftaran/:id/complete-bimbingan`

Menandai bahwa bimbingan perkawinan sudah selesai (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Bimbingan perkawinan berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Sudah Bimbingan",
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "updated_by": "USR1704067201",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 35. Complete Marriage
**PUT** `/simnikah/pendaftaran/:id/complete-nikah`

Menandai bahwa nikah sudah dilaksanakan (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Nikah berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_sekarang": "Selesai",
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "completed_by": "USR1704067201",
    "completed_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 🔔 Notification Endpoints

### 36. Create Notification
**POST** `/simnikah/notifikasi`

Membuat notifikasi baru (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "user_id": "USR1704067200",
  "judul": "Pengingat Nikah",
  "pesan": "Nikah Anda akan dilaksanakan besok pukul 09:00",
  "tipe": "Warning",
  "link": "/simnikah/pendaftaran/1"
}
```

**Valid Types:**
- `"Info"`
- `"Success"`
- `"Warning"`
- `"Error"`

**Response:**
```json
{
  "message": "Notifikasi berhasil dibuat",
  "data": {
    "id": 1,
    "user_id": "USR1704067200",
    "judul": "Pengingat Nikah",
    "pesan": "Nikah Anda akan dilaksanakan besok pukul 09:00",
    "tipe": "Warning",
    "status_baca": "Belum Dibaca",
    "link": "/simnikah/pendaftaran/1",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 37. Get User Notifications
**GET** `/simnikah/notifikasi/user/:user_id`

Mendapatkan notifikasi untuk user tertentu.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Notifikasi user berhasil diambil",
  "data": [
    {
      "id": 1,
      "user_id": "USR1704067200",
      "judul": "Pengingat Nikah",
      "pesan": "Nikah Anda akan dilaksanakan besok pukul 09:00",
      "tipe": "Warning",
      "status_baca": "Belum Dibaca",
      "link": "/simnikah/pendaftaran/1",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### 38. Get Notification by ID
**GET** `/simnikah/notifikasi/:id`

Mendapatkan detail notifikasi.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Notifikasi berhasil diambil",
  "data": {
    "id": 1,
    "user_id": "USR1704067200",
    "judul": "Pengingat Nikah",
    "pesan": "Nikah Anda akan dilaksanakan besok pukul 09:00",
    "tipe": "Warning",
    "status_baca": "Belum Dibaca",
    "link": "/simnikah/pendaftaran/1",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 39. Update Notification Status
**PUT** `/simnikah/notifikasi/:id/status`

Mengupdate status notifikasi (dibaca/belum dibaca).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "status_baca": "Sudah Dibaca"
}
```

**Valid Status:**
- `"Belum Dibaca"`
- `"Sudah Dibaca"`

**Response:**
```json
{
  "message": "Status notifikasi berhasil diupdate",
  "data": {
    "id": 1,
    "status_baca": "Sudah Dibaca",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 40. Mark All Notifications as Read
**PUT** `/simnikah/notifikasi/user/:user_id/mark-all-read`

Menandai semua notifikasi user sebagai sudah dibaca.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Semua notifikasi berhasil ditandai sebagai sudah dibaca",
  "data": {
    "user_id": "USR1704067200",
    "updated_count": 5,
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 41. Delete Notification
**DELETE** `/simnikah/notifikasi/:id`

Menghapus notifikasi.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Notifikasi berhasil dihapus",
  "data": {
    "id": 1,
    "deleted_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 42. Get Notification Stats
**GET** `/simnikah/notifikasi/user/:user_id/stats`

Mendapatkan statistik notifikasi user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Statistik notifikasi berhasil diambil",
  "data": {
    "user_id": "USR1704067200",
    "total_notifications": 10,
    "unread_count": 3,
    "read_count": 7,
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

### 43. Send Notification to Role
**POST** `/simnikah/notifikasi/send-to-role`

Mengirim notifikasi ke semua user dengan role tertentu (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "role": "staff",
  "judul": "Pemberitahuan Penting",
  "pesan": "Ada perubahan jadwal bimbingan perkawinan",
  "tipe": "Info",
  "link": "/simnikah/bimbingan"
}
```

**Valid Roles:**
- `"user_biasa"`
- `"staff"`
- `"penghulu"`
- `"kepala_kua"`

**Response:**
```json
{
  "message": "Notifikasi berhasil dikirim ke role staff",
  "data": {
    "role": "staff",
    "total_sent": 5,
    "judul": "Pemberitahuan Penting",
    "tipe": "Info",
    "sent_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 44. Run Reminder Notification
**POST** `/simnikah/notifikasi/run-reminder`

Menjalankan pengingat notifikasi secara manual (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Pengingat notifikasi berhasil dijalankan",
  "executed_by": "USR1704067201",
  "executed_at": "2024-01-01T12:00:00Z"
}
```

---

## 📊 Data Management Endpoints

### 45. Get All Marriage Registrations
**GET** `/simnikah/pendaftaran`

Mendapatkan semua pendaftaran nikah dengan filter dan pagination (Staff/Kepala KUA).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `page` (optional): Halaman, default: 1
- `limit` (optional): Jumlah data per halaman, default: 10, max: 100
- `status` (optional): Filter berdasarkan status
- `date_from` (optional): Filter dari tanggal (YYYY-MM-DD)
- `date_to` (optional): Filter sampai tanggal (YYYY-MM-DD)
- `location` (optional): Filter berdasarkan lokasi
- `search` (optional): Pencarian berdasarkan nomor pendaftaran, nama, atau NIK
- `sort_by` (optional): Field untuk sorting, default: created_at
- `sort_order` (optional): Urutan sorting (asc/desc), default: desc

**Example:**
```
GET /simnikah/pendaftaran?page=1&limit=10&status=Menunggu Verifikasi&search=Ahmad
```

**Response:**
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
        "tanggal_pendaftaran": "2024-01-01T00:00:00Z",
        "tanggal_nikah": "2024-02-14T00:00:00Z",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "KUA Banjarmasin Utara",
        "nomor_dispensasi": "",
        "penghulu_id": null,
        "catatan": "",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
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
      "search": "Ahmad",
      "sort_by": "created_at",
      "sort_order": "desc"
    }
  }
}
```

---

## 🏥 Health Check

### 46. Health Check
**GET** `/health`

Mengecek status kesehatan aplikasi.

**Response:**
```json
{
  "status": "healthy",
  "service": "SimNikah API",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## 📝 Error Responses

### Common Error Formats

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format data tidak valid",
  "type": "validation"
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

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Database error",
  "error": "Gagal mengambil data pendaftaran",
  "type": "database"
}
```

---

## 🔐 Authentication & Authorization

### JWT Token Structure
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

### Role Permissions

| Endpoint | user_biasa | staff | penghulu | kepala_kua |
|----------|------------|-------|----------|-------------|
| `/register` | ✅ | ✅ | ✅ | ✅ |
| `/login` | ✅ | ✅ | ✅ | ✅ |
| `/profile` | ✅ | ✅ | ✅ | ✅ |
| `/simnikah/pendaftaran/form-baru` | ✅ | ❌ | ❌ | ❌ |
| `/simnikah/staff/*` | ❌ | ❌ | ❌ | ✅ |
| `/simnikah/penghulu/*` | ❌ | ❌ | ✅ | ✅ |
| `/simnikah/staff/verify-*` | ❌ | ✅ | ❌ | ✅ |
| `/simnikah/penghulu/verify-*` | ❌ | ❌ | ✅ | ❌ |

---

## 📋 Status Codes Reference

| Status | Description |
|--------|-------------|
| `200` | OK - Request berhasil |
| `201` | Created - Resource berhasil dibuat |
| `400` | Bad Request - Format data tidak valid |
| `401` | Unauthorized - Token tidak valid atau tidak ada |
| `403` | Forbidden - Role tidak memiliki akses |
| `404` | Not Found - Resource tidak ditemukan |
| `409` | Conflict - Data sudah ada |
| `422` | Unprocessable Entity - Validasi bisnis gagal |
| `500` | Internal Server Error - Error server |

---

## 🚀 Getting Started

### 1. Register Admin User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@kua.go.id",
    "password": "admin123",
    "nama": "Administrator",
    "role": "kepala_kua"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 3. Use Token in Requests
```bash
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer <jwt_token>"
```

---

## 📞 Support

- **API Version**: 1.0.0
- **Base URL**: `http://localhost:8080`
- **Documentation**: [API Documentation](./API_DOCUMENTATION.md)
- **Testing**: [API Testing Documentation](./API_TESTING_DOCUMENTATION.md)

---

*Dokumentasi ini dibuat untuk SimNikah API v1.0.0. Untuk pertanyaan atau bantuan, silakan hubungi tim development.*
