# SimNikah API - Dokumentasi Frontend

Dokumentasi lengkap untuk integrasi frontend dengan SimNikah API.

## Daftar Isi

1. [Informasi Umum](#informasi-umum)
2. [Autentikasi](#autentikasi)
3. [Role & Akses](#role--akses)
4. [API Endpoints](#api-endpoints)
5. [Response Format](#response-format)
6. [Error Handling](#error-handling)
7. [Contoh Implementasi](#contoh-implementasi)

---

## Informasi Umum

### Base URL

```
Production: https://your-api-domain.com
Development: http://localhost:8080
```

### Headers

Semua request yang memerlukan autentikasi harus menyertakan header:

```http
Authorization: Bearer <token>
Content-Type: application/json
```

### Rate Limiting

| Tipe | Limit | Keterangan |
|------|-------|------------|
| Global | 100 req/menit | Semua endpoint |
| Auth (Login/Register) | 5 req/menit | Endpoint sensitif |

---

## Autentikasi

### 1. Register

**Endpoint:** `POST /register`

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "nama": "John Doe",
  "no_hp": "081234567890",
  "alamat": "Jl. Contoh No. 123",
  "role": "user_biasa"
}
```

**Response Success (201):**
```json
{
  "success": true,
  "message": "Registrasi berhasil",
  "data": {
    "user_id": "USR-20241203-XXXXX",
    "username": "johndoe",
    "email": "john@example.com",
    "nama": "John Doe",
    "role": "user_biasa"
  }
}
```

### 2. Login

**Endpoint:** `POST /login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "user_id": "USR-20241203-XXXXX",
      "username": "johndoe",
      "email": "john@example.com",
      "nama": "John Doe",
      "role": "user_biasa"
    }
  }
}
```

### 3. Get Profile

**Endpoint:** `GET /profile`

**Headers:** `Authorization: Bearer <token>`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "user_id": "USR-20241203-XXXXX",
    "username": "johndoe",
    "email": "john@example.com",
    "nama": "John Doe",
    "no_hp": "081234567890",
    "alamat": "Jl. Contoh No. 123",
    "role": "user_biasa",
    "status": "aktif"
  }
}
```

---

## Role & Akses

### Daftar Role

| Role | Kode | Deskripsi |
|------|------|-----------|
| User Biasa | `user_biasa` | Calon pengantin yang mendaftar nikah |
| Staff KUA | `staff` | Staff yang memverifikasi pendaftaran |
| Penghulu | `penghulu` | Petugas yang menikahkan |
| Kepala KUA | `kepala_kua` | Administrator KUA |

### Matriks Akses Endpoint

| Endpoint | user_biasa | staff | penghulu | kepala_kua |
|----------|------------|-------|----------|------------|
| Pendaftaran Nikah | ✅ | ✅ | ❌ | ✅ |
| Verifikasi Formulir | ❌ | ✅ | ❌ | ❌ |
| Verifikasi Dokumen | ❌ | ✅ | ✅ | ❌ |
| Assign Penghulu | ❌ | ❌ | ❌ | ✅ |
| Selesaikan Nikah | ❌ | ❌ | ✅ | ✅ |
| Dashboard | ❌ | ✅ | ❌ | ✅ |
| Kelola Staff | ❌ | ❌ | ❌ | ✅ |
| Kelola Penghulu | ❌ | ❌ | ❌ | ✅ |

---

## API Endpoints

### A. Pendaftaran Nikah (Catin)

#### 1. Buat Pendaftaran Nikah

**Endpoint:** `POST /simnikah/pendaftaran`

**Role:** `user_biasa`, `staff`, `kepala_kua`

**Request Body:**
```json
{
  "nama_pria": "Ahmad",
  "nik_pria": "1234567890123456",
  "tempat_lahir_pria": "Banjarmasin",
  "tanggal_lahir_pria": "1995-05-15",
  "alamat_pria": "Jl. A Yani No. 1",
  "no_hp_pria": "081234567890",
  "pekerjaan_pria": "Karyawan Swasta",
  "status_perkawinan_pria": "Belum Kawin",
  
  "nama_wanita": "Siti",
  "nik_wanita": "6543210987654321",
  "tempat_lahir_wanita": "Banjarmasin",
  "tanggal_lahir_wanita": "1998-08-20",
  "alamat_wanita": "Jl. Sudirman No. 2",
  "no_hp_wanita": "089876543210",
  "pekerjaan_wanita": "PNS",
  "status_perkawinan_wanita": "Belum Kawin",
  
  "tanggal_nikah": "2024-12-25",
  "jam_nikah": "09:00",
  "tempat_nikah": "di_kua",
  "alamat_nikah": "",
  "kelurahan": "Sungai Miai",
  
  "nama_wali": "Bapak Ahmad",
  "hubungan_wali": "Ayah Kandung",
  "no_hp_wali": "081111222333"
}
```

**Response Success (201):**
```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat",
  "data": {
    "pendaftaran_id": "REG-20241203-XXXXX",
    "status": "Draft",
    "tanggal_nikah": "2024-12-25",
    "jam_nikah": "09:00"
  }
}
```

#### 2. Cek Status Pendaftaran

**Endpoint:** `GET /simnikah/pendaftaran/status`

**Role:** `user_biasa`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "pendaftaran_id": "REG-20241203-XXXXX",
    "status": "Menunggu Verifikasi",
    "tanggal_nikah": "2024-12-25",
    "jam_nikah": "09:00",
    "tempat_nikah": "di_kua",
    "penghulu": null,
    "catatan_verifikasi": null
  }
}
```

#### 3. List Pendaftaran (Staff/Kepala KUA)

**Endpoint:** `GET /simnikah/pendaftaran`

**Role:** `staff`, `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `status` | string | Filter berdasarkan status |
| `tanggal_mulai` | date | Filter tanggal mulai (YYYY-MM-DD) |
| `tanggal_akhir` | date | Filter tanggal akhir (YYYY-MM-DD) |
| `page` | int | Halaman (default: 1) |
| `limit` | int | Jumlah per halaman (default: 10) |

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "registrations": [...],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 48,
      "items_per_page": 10
    }
  }
}
```

---

### B. Kalender & Ketersediaan

#### 1. Get Kalender Ketersediaan

**Endpoint:** `GET /simnikah/kalender-ketersediaan`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `bulan` | int | Bulan (1-12) |
| `tahun` | int | Tahun (YYYY) |

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "bulan": 12,
    "tahun": 2024,
    "tanggal": [
      {
        "tanggal": "2024-12-01",
        "hari": "Minggu",
        "tersedia": false,
        "alasan": "Hari Minggu"
      },
      {
        "tanggal": "2024-12-02",
        "hari": "Senin",
        "tersedia": true,
        "slot_tersisa": 7
      }
    ]
  }
}
```

#### 2. Get Slot Waktu Tersedia

**Endpoint:** `GET /simnikah/ketersediaan-jam`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `tanggal` | date | Tanggal (YYYY-MM-DD) |
| `tempat` | string | `di_kua` atau `luar_kua` |

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-12-25",
    "tempat": "di_kua",
    "slots": [
      { "jam": "08:00", "tersedia": true },
      { "jam": "09:00", "tersedia": false, "alasan": "Sudah terisi" },
      { "jam": "10:00", "tersedia": true },
      { "jam": "11:00", "tersedia": true },
      { "jam": "13:00", "tersedia": true },
      { "jam": "14:00", "tersedia": true },
      { "jam": "15:00", "tersedia": true }
    ]
  }
}
```

---

### C. Verifikasi (Staff)

#### 1. Verifikasi Formulir

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`

**Role:** `staff`

**Request Body:**
```json
{
  "status": "Disetujui",
  "catatan": "Formulir lengkap dan valid"
}
```

**Status yang valid:** `Disetujui`, `Ditolak`

#### 2. Approve Pendaftaran

**Endpoint:** `POST /simnikah/staff/approve/:id`

**Role:** `staff`

**Request Body:**
```json
{
  "status": "Menunggu Penugasan",
  "catatan": "Pendaftaran disetujui, menunggu penugasan penghulu"
}
```

---

### D. Penghulu

#### 1. List Penugasan Saya

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`

**Role:** `penghulu`, `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `filter` | string | `today` atau `tanggal` |
| `tanggal` | date | Jika filter=tanggal (YYYY-MM-DD) |

#### 2. Jadwal Hari Ini

**Endpoint:** `GET /simnikah/penghulu/today-schedule`

**Role:** `penghulu`, `kepala_kua`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "tanggal": "2024-12-03",
    "total_jadwal": 3,
    "jadwal": [
      {
        "pendaftaran_id": "REG-20241203-XXXXX",
        "jam": "09:00",
        "nama_pria": "Ahmad",
        "nama_wanita": "Siti",
        "tempat": "di_kua",
        "alamat": null,
        "status": "Penghulu Ditugaskan"
      }
    ]
  }
}
```

#### 3. Verifikasi Dokumen (Penghulu)

**Endpoint:** `POST /simnikah/penghulu/verify-documents/:id`

**Role:** `penghulu`

**Request Body:**
```json
{
  "status": "Disetujui",
  "catatan": "Dokumen sudah diverifikasi di lokasi"
}
```

#### 4. Selesaikan Pernikahan

**Endpoint:** `POST /simnikah/penghulu/complete-marriage/:id`

**Role:** `penghulu`, `kepala_kua`

**Request Body:**
```json
{
  "catatan": "Pernikahan telah dilaksanakan dengan lancar"
}
```

---

### E. Kepala KUA

#### 1. Assign Penghulu

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`

**Role:** `kepala_kua`

**Request Body:**
```json
{
  "penghulu_id": 1
}
```

#### 2. List Penghulu Tersedia

**Endpoint:** `GET /simnikah/kepala-kua/available-penghulu`

**Role:** `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `tanggal` | date | Tanggal untuk cek ketersediaan |
| `jam` | string | Jam untuk cek ketersediaan (HH:MM) |

#### 3. Buat Staff Baru

**Endpoint:** `POST /simnikah/kepala-kua/staff`

**Role:** `kepala_kua`

**Request Body:**
```json
{
  "username": "staff_baru",
  "email": "staff@kua.go.id",
  "password": "password123",
  "nama": "Staff Baru",
  "nip": "199001012020011001",
  "jabatan": "Staff Administrasi",
  "no_hp": "081234567890"
}
```

#### 4. Buat Penghulu Baru

**Endpoint:** `POST /simnikah/kepala-kua/penghulu`

**Role:** `kepala_kua`

**Request Body:**
```json
{
  "username": "penghulu_baru",
  "email": "penghulu@kua.go.id",
  "password": "password123",
  "nama": "Penghulu Baru",
  "nip": "198501012010011001",
  "no_hp": "081234567890"
}
```

#### 5. Lihat Jadwal Ketersediaan Semua Penghulu

**Endpoint:** `GET /simnikah/kepala-kua/penghulu-schedule`

**Role:** `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `tanggal` | date | Tanggal (YYYY-MM-DD), default: hari ini |

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data ketersediaan jadwal penghulu berhasil diambil",
  "data": {
    "tanggal": "2024-12-25",
    "hari": "Wednesday",
    "tanggal_format": "25 Desember 2024",
    "total_penghulu": 3,
    "total_jadwal_hari_ini": 5,
    "time_slots": ["08:00", "09:00", "10:00", ...],
    "penghulu_availability": [
      {
        "penghulu": {
          "id": 1,
          "nama_lengkap": "H. Ahmad",
          "nip": "198501012010011001",
          "jumlah_nikah": 150,
          "rating": 4.8
        },
        "tanggal": "2024-12-25",
        "jadwal_hari_ini": 2,
        "slot_tersedia": 7,
        "slot_terisi": 2,
        "time_slots": [
          {
            "waktu": "08:00",
            "tersedia": true,
            "status": "Tersedia",
            "jadwal": null
          },
          {
            "waktu": "09:00",
            "tersedia": false,
            "status": "Bertugas",
            "jadwal": {
              "pendaftaran_id": 123,
              "nomor_pendaftaran": "REG-20241225-001",
              "waktu_nikah": "09:00",
              "tempat_nikah": "Di KUA",
              "calon_suami": "Ahmad",
              "calon_istri": "Siti",
              "status": "Penghulu Ditugaskan"
            }
          }
        ],
        "jadwal_detail": [...]
      }
    ]
  }
}
```

#### 6. Lihat Penghulu Tersedia untuk Jam Tertentu

**Endpoint:** `GET /simnikah/kepala-kua/penghulu-tersedia`

**Role:** `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `tanggal` | date | Tanggal (YYYY-MM-DD) - **required** |
| `waktu` | string | Jam (HH:MM) - **required** |

**Contoh:** `/simnikah/kepala-kua/penghulu-tersedia?tanggal=2024-12-25&waktu=09:00`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Data ketersediaan penghulu berhasil diambil",
  "data": {
    "tanggal": "2024-12-25",
    "waktu": "09:00",
    "hari": "Wednesday",
    "tanggal_format": "25 Desember 2024",
    "total_penghulu": 3,
    "jumlah_tersedia": 2,
    "jumlah_tidak_tersedia": 1,
    "penghulu_tersedia": [
      {
        "id": 2,
        "nama_lengkap": "H. Mahmud",
        "nip": "198601012010011002",
        "jumlah_nikah": 120,
        "rating": 4.5,
        "jadwal_hari_ini": 1,
        "status": "Tersedia"
      }
    ],
    "penghulu_tidak_tersedia": [
      {
        "id": 1,
        "nama_lengkap": "H. Ahmad",
        "nip": "198501012010011001",
        "jumlah_nikah": 150,
        "rating": 4.8,
        "jadwal_hari_ini": 2,
        "status": "Tidak Tersedia",
        "alasan": "Sudah ada jadwal pada pukul 09:00"
      }
    ]
  }
}
```

---

### F. Dashboard

#### 1. Dashboard Kepala KUA

**Endpoint:** `GET /simnikah/dashboard/kepala-kua`

**Role:** `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `periode` | string | `minggu`, `bulan`, `tahun` |
| `tanggal_mulai` | date | Custom range start |
| `tanggal_akhir` | date | Custom range end |

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "statistik": {
      "total_pendaftaran": 150,
      "menunggu_verifikasi": 12,
      "disetujui": 45,
      "selesai": 80,
      "ditolak": 13
    },
    "tren_pernikahan": [...],
    "performa_penghulu": [...],
    "jam_sibuk": [...]
  }
}
```

#### 2. Dashboard Staff

**Endpoint:** `GET /simnikah/dashboard/staff`

**Role:** `staff`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "pending_verifikasi": 5,
    "pending_dokumen": 3,
    "timeline_pendaftaran": [...]
  }
}
```

#### 3. Statistik Pernikahan

**Endpoint:** `GET /simnikah/dashboard/statistik-pernikahan`

**Role:** `staff`, `kepala_kua`

#### 4. Performa Penghulu

**Endpoint:** `GET /simnikah/dashboard/penghulu-performance`

**Role:** `staff`, `kepala_kua`

#### 5. Analisis Jam Sibuk

**Endpoint:** `GET /simnikah/dashboard/peak-hours`

**Role:** `staff`, `kepala_kua`

---

### G. Lokasi

#### 1. Geocode (Alamat ke Koordinat)

**Endpoint:** `POST /simnikah/location/geocode`

**Request Body:**
```json
{
  "alamat": "Jl. A Yani No. 1, Banjarmasin"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "latitude": -3.3194374,
    "longitude": 114.5900474,
    "alamat_lengkap": "Jl. A Yani No. 1, Banjarmasin, Kalimantan Selatan"
  }
}
```

#### 2. Reverse Geocode (Koordinat ke Alamat)

**Endpoint:** `POST /simnikah/location/reverse-geocode`

**Request Body:**
```json
{
  "latitude": -3.3194374,
  "longitude": 114.5900474
}
```

#### 3. Search Alamat (Autocomplete)

**Endpoint:** `GET /simnikah/location/search`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `q` | string | Query pencarian alamat |
| `limit` | int | Jumlah hasil (default: 5) |

#### 4. Update Lokasi Pernikahan

**Endpoint:** `PUT /simnikah/pendaftaran/:id/location`

**Request Body:**
```json
{
  "alamat": "Jl. Baru No. 123",
  "latitude": -3.3194374,
  "longitude": 114.5900474
}
```

#### 5. Get Detail Lokasi Pernikahan

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`

---

### H. Notifikasi

#### 1. Get Notifikasi User

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `status` | string | `belum_dibaca`, `dibaca` |
| `page` | int | Halaman |
| `limit` | int | Jumlah per halaman |

#### 2. Get Detail Notifikasi

**Endpoint:** `GET /simnikah/notifikasi/:id`

#### 3. Update Status Notifikasi

**Endpoint:** `PUT /simnikah/notifikasi/:id/status`

**Request Body:**
```json
{
  "status": "dibaca"
}
```

#### 4. Tandai Semua Dibaca

**Endpoint:** `PUT /simnikah/notifikasi/user/:user_id/mark-all-read`

#### 5. Hapus Notifikasi

**Endpoint:** `DELETE /simnikah/notifikasi/:id`

#### 6. Statistik Notifikasi

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id/stats`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "total": 25,
    "belum_dibaca": 5,
    "dibaca": 20
  }
}
```

#### 7. Kirim Notifikasi ke Role (Staff/Kepala KUA)

**Endpoint:** `POST /simnikah/notifikasi/send-to-role`

**Role:** `staff`, `kepala_kua`

**Request Body:**
```json
{
  "role": "user_biasa",
  "judul": "Pengumuman Penting",
  "pesan": "Isi pengumuman...",
  "tipe": "pengumuman"
}
```

---

### I. Feedback

#### 1. Buat Feedback Pernikahan

**Endpoint:** `POST /simnikah/feedback-pernikahan`

**Role:** `user_biasa`

**Request Body:**
```json
{
  "pendaftaran_id": "REG-20241203-XXXXX",
  "rating": 5,
  "komentar": "Pelayanan sangat baik dan ramah",
  "jenis": "pujian"
}
```

#### 2. List Feedback (Kepala KUA)

**Endpoint:** `GET /simnikah/kepala-kua/feedback`

**Role:** `kepala_kua`

**Query Parameters:**
| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `jenis` | string | `pujian`, `saran`, `keluhan` |
| `status` | string | `belum_dibaca`, `dibaca` |

#### 3. Tandai Feedback Dibaca

**Endpoint:** `PUT /simnikah/kepala-kua/feedback/:id/mark-read`

**Role:** `kepala_kua`

#### 4. Statistik Feedback

**Endpoint:** `GET /simnikah/kepala-kua/feedback/stats`

**Role:** `kepala_kua`

---

## Response Format

### Success Response

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Pesan error",
  "error": "Detail error"
}
```

---

## Error Handling

### HTTP Status Codes

| Code | Deskripsi |
|------|-----------|
| 200 | OK - Request berhasil |
| 201 | Created - Resource berhasil dibuat |
| 400 | Bad Request - Request tidak valid |
| 401 | Unauthorized - Token tidak valid/expired |
| 403 | Forbidden - Tidak memiliki akses |
| 404 | Not Found - Resource tidak ditemukan |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Error server |

### Contoh Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token tidak valid atau sudah expired"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak memiliki akses untuk endpoint ini"
}
```

**429 Rate Limit:**
```json
{
  "success": false,
  "message": "Terlalu banyak request",
  "error": "Silakan coba lagi dalam 1 menit"
}
```

---

## Contoh Implementasi

### React/Next.js dengan Axios

```javascript
// lib/api.js
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor untuk menambahkan token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor untuk handle error
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

### Contoh Login

```javascript
// hooks/useAuth.js
import { useState } from 'react';
import api from '../lib/api';

export function useAuth() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const login = async (email, password) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await api.post('/login', { email, password });
      const { token, user } = response.data.data;
      
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      return { success: true, user };
    } catch (err) {
      setError(err.response?.data?.message || 'Login gagal');
      return { success: false, error: err.response?.data };
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
  };

  return { login, logout, loading, error };
}
```

### Contoh Fetch Pendaftaran

```javascript
// hooks/usePendaftaran.js
import { useState, useEffect } from 'react';
import api from '../lib/api';

export function usePendaftaran() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchStatus = async () => {
    try {
      const response = await api.get('/simnikah/pendaftaran/status');
      setData(response.data.data);
    } catch (err) {
      setError(err.response?.data?.message);
    } finally {
      setLoading(false);
    }
  };

  const createPendaftaran = async (formData) => {
    try {
      const response = await api.post('/simnikah/pendaftaran', formData);
      return { success: true, data: response.data.data };
    } catch (err) {
      return { success: false, error: err.response?.data?.message };
    }
  };

  useEffect(() => {
    fetchStatus();
  }, []);

  return { data, loading, error, createPendaftaran, refetch: fetchStatus };
}
```

### Contoh Component Dashboard

```jsx
// components/Dashboard.jsx
import { useEffect, useState } from 'react';
import api from '../lib/api';

export default function Dashboard() {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const user = JSON.parse(localStorage.getItem('user'));
        const endpoint = user.role === 'kepala_kua' 
          ? '/simnikah/dashboard/kepala-kua'
          : '/simnikah/dashboard/staff';
        
        const response = await api.get(endpoint);
        setStats(response.data.data);
      } catch (error) {
        console.error('Error fetching dashboard:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboard();
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Pendaftaran</h3>
          <p>{stats?.statistik?.total_pendaftaran || 0}</p>
        </div>
        <div className="stat-card">
          <h3>Menunggu Verifikasi</h3>
          <p>{stats?.statistik?.menunggu_verifikasi || 0}</p>
        </div>
        <div className="stat-card">
          <h3>Selesai</h3>
          <p>{stats?.statistik?.selesai || 0}</p>
        </div>
      </div>
    </div>
  );
}
```

---

## Status Pendaftaran Flow

```
Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
                ↓
             Ditolak
```

| Status | Deskripsi |
|--------|-----------|
| `Draft` | Pendaftaran baru dibuat |
| `Disetujui` | Formulir sudah diverifikasi staff |
| `Menunggu Penugasan` | Dokumen disetujui, menunggu penghulu |
| `Penghulu Ditugaskan` | Penghulu sudah ditugaskan |
| `Selesai` | Pernikahan sudah dilaksanakan |
| `Ditolak` | Pendaftaran ditolak |

---

## Aturan Hari Libur

### Hari Libur Nasional

Sistem mengenali hari libur nasional Indonesia:

**Hari Libur Tetap (Setiap Tahun):**
- 1 Januari - Tahun Baru Masehi
- 1 Mei - Hari Buruh Internasional
- 1 Juni - Hari Lahir Pancasila
- 17 Agustus - Hari Kemerdekaan RI
- 25 Desember - Hari Natal

**Hari Libur Berubah (2024-2025):**
- Idul Fitri, Idul Adha, Nyepi, Waisak, dll.

### Aturan Ketersediaan pada Hari Libur

| Lokasi | Hari Libur | Hari Biasa |
|--------|------------|------------|
| **Di KUA** | ❌ Tidak Tersedia | ✅ Tersedia |
| **Di Luar KUA** | ✅ Tersedia | ✅ Tersedia |

### Response API pada Hari Libur

**GET /simnikah/kalender-ketersediaan**
```json
{
  "tanggal": 25,
  "tanggal_str": "2024-12-25",
  "hari": "Wednesday",
  "status": "Hari Libur (Luar KUA Tersedia)",
  "tersedia": true,
  "tersedia_kua": false,
  "tersedia_luar_kua": true,
  "is_hari_libur": true,
  "nama_hari_libur": "Hari Natal"
}
```

**Error saat mendaftar di KUA pada hari libur:**
```json
{
  "success": false,
  "message": "KUA tutup pada hari libur",
  "error": "Maaf, KUA tutup pada tanggal 25 Desember 2024 (Hari Natal)...",
  "type": "holiday_restriction",
  "saran": "Anda masih bisa memilih nikah di luar KUA pada hari libur...",
  "data": {
    "tanggal_nikah": "2024-12-25",
    "is_hari_libur": true,
    "nama_hari_libur": "Hari Natal",
    "tersedia_luar_kua": true
  }
}
```

---

## Environment Variables (Frontend)

```env
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_GOOGLE_MAPS_KEY=your_google_maps_key
```

---

## Tips & Best Practices

1. **Token Storage**: Simpan token di `localStorage` atau `httpOnly cookie`
2. **Error Handling**: Selalu handle error 401 untuk redirect ke login
3. **Loading States**: Tampilkan loading indicator saat fetch data
4. **Optimistic Updates**: Update UI sebelum response untuk UX yang lebih baik
5. **Caching**: Gunakan SWR atau React Query untuk caching data
6. **Form Validation**: Validasi input di frontend sebelum submit

---

## Kontak & Support

Jika ada pertanyaan atau masalah, silakan hubungi:
- Email: support@simnikah.go.id
- Dokumentasi API: `/docs` (jika tersedia Swagger)

---

*Dokumentasi ini dibuat untuk SimNikah API v1.0*
