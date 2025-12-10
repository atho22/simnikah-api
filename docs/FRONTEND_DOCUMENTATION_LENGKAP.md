# 📘 Dokumentasi Lengkap Frontend - SimNikah API

**Versi:** 2.0.0  
**Update Terakhir:** Desember 2024  
**Target Audience:** Frontend Developer  

---

## 📋 Daftar Isi

1. [Pengenalan](#pengenalan)
2. [Konfigurasi & Setup](#konfigurasi--setup)
3. [Authentication & Authorization](#authentication--authorization)
4. [Endpoints Berdasarkan Fitur](#endpoints-berdasarkan-fitur)
5. [Request & Response Format](#request--response-format)
6. [Error Handling](#error-handling)
7. [Contoh Implementasi](#contoh-implementasi)
8. [Best Practices](#best-practices)
9. [✅ Validasi Tampilan Frontend (Wajib Diimplementasikan)](#-validasi-tampilan-frontend-wajib-diimplementasikan)
10. [Testing](#testing)
11. [FAQ](#faq)

---

## 🎯 Pengenalan

SimNikah adalah sistem informasi pendaftaran nikah berbasis API untuk KUA (Kantor Urusan Agama) Kecamatan Banjarmasin Utara. Dokumentasi ini menjelaskan cara mengintegrasikan frontend dengan API backend SimNikah.

### Base URL
- **Development:** `http://localhost:8080`
- **Production:** Sesuai dengan domain yang diberikan (Railway/Vercel)

### Teknologi Backend
- Framework: Gin (Golang)
- Database: MySQL/MariaDB
- Authentication: JWT (JSON Web Token)
- File Storage: ImgBB (untuk foto profil)

---

## ⚙️ Konfigurasi & Setup

### 1. Environment Variables

Buat file `.env` atau konfigurasi environment variables:

```env
# API Base URL
REACT_APP_API_URL=http://localhost:8080
# atau untuk production:
# REACT_APP_API_URL=https://your-api-domain.com

# Optional: Timeout untuk request
REACT_APP_API_TIMEOUT=30000
```

### 2. Install Dependencies

#### React/Next.js
```bash
npm install axios
# atau
yarn add axios
```

#### Vue.js
```bash
npm install axios
# atau
yarn add axios
```

### 3. API Client Setup

#### React/Next.js Example
```javascript
// utils/api.js
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor untuk menambahkan token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor untuk handle error
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired atau invalid
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error.response?.data || error);
  }
);

export default apiClient;
```

#### Vue.js Example
```javascript
// utils/api.js
import axios from 'axios';

const API_BASE_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      router.push('/login');
    }
    return Promise.reject(error.response?.data || error);
  }
);

export default apiClient;
```

---

## 🔐 Authentication & Authorization

### Flow Authentication

```
1. User Register/Login
   ↓
2. Backend mengembalikan JWT Token
   ↓
3. Frontend menyimpan token (localStorage/cookie)
   ↓
4. Setiap request mengirim token di header Authorization
   ↓
5. Backend memverifikasi token
   ↓
6. Request berhasil atau ditolak
```

### Menyimpan Token

**LocalStorage (Recommended untuk SPA):**
```javascript
// Setelah login
localStorage.setItem('token', response.token);
localStorage.setItem('user', JSON.stringify(response.user));

// Untuk request
const token = localStorage.getItem('token');
```

**Cookie (Recommended untuk SSR/Next.js):**
```javascript
// Setelah login (gunakan library js-cookie)
import Cookies from 'js-cookie';

Cookies.set('token', response.token, { expires: 7 }); // 7 hari
```

### Endpoints Authentication

#### 1. Register User

**Endpoint:** `POST /register`

**Request:**
```javascript
const registerData = {
  username: "ahmad123",
  email: "ahmad@example.com",
  password: "password123",
  nama: "Ahmad Wijaya",
  role: "user_biasa" // user_biasa, staff, penghulu, kepala_kua
};

const response = await apiClient.post('/register', registerData);
```

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

#### 2. Login

**Endpoint:** `POST /login`

**Request:**
```javascript
const loginData = {
  username: "ahmad123",
  password: "password123"
};

const response = await apiClient.post('/login', loginData);

// Simpan token
localStorage.setItem('token', response.token);
localStorage.setItem('user', JSON.stringify(response.user));
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

#### 3. Get Profile

**Endpoint:** `GET /profile`

**Request:**
```javascript
const response = await apiClient.get('/profile');
// Token otomatis ditambahkan oleh interceptor
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
    "profile_photo": "https://...",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 4. Upload Photo Profile

**Endpoint:** `POST /upload-photo`

**Request:**
```javascript
const formData = new FormData();
formData.append('photo', fileInput.files[0]); // File dari input

const response = await apiClient.post('/upload-photo', formData, {
  headers: {
    'Content-Type': 'multipart/form-data',
  },
});
```

**Validasi:**
- Max file size: 5MB
- Allowed formats: `image/jpeg`, `image/png`, `image/jpg`, `image/webp`

**Response Success (200):**
```json
{
  "success": true,
  "message": "Foto profil berhasil diupload",
  "data": {
    "profile_photo": "https://i.ibb.co/...",
    "user_id": "USR1704067200",
    "username": "ahmad123"
  }
}
```

---

## 📍 Endpoints Berdasarkan Fitur

### 🔵 Publik (Tanpa Authentication)

#### 1. Kalender Ketersediaan

**Endpoint:** `GET /simnikah/kalender-ketersediaan`

**Query Parameters:**
- `bulan` (required): 1-12
- `tahun` (required): YYYY

**Request:**
```javascript
const response = await apiClient.get('/simnikah/kalender-ketersediaan', {
  params: {
    bulan: 12,
    tahun: 2024
  }
});
```

**Response:**
```json
{
  "success": true,
  "message": "Kalender ketersediaan berhasil diambil",
  "data": {
    "bulan": 12,
    "tahun": 2024,
    "nama_bulan": "December",
    "kapasitas_harian": 9,
    "calendar": [
      {
        "tanggal": 1,
        "tanggal_str": "2024-12-01",
        "hari": "Sunday",
        "status": "Tersedia",
        "tersedia": true,
        "jumlah_nikah": 5,
        "jumlah_draft": 2,
        "jumlah_disetujui": 3,
        "sisa_kuota": 6,
        "kapasitas": 9,
        "is_today": false,
        "is_past": false,
        "time_slots": [
          {
            "waktu": "08:00",
            "kua": {
              "tersedia": true,
              "terbooking": false,
              "jumlah_total": 0,
              "jumlah_draft": 0,
              "jumlah_disetujui": 0
            },
            "luar_kua": {
              "tersedia": true,
              "terbooking": false,
              "jumlah_total": 0,
              "jumlah_draft": 0,
              "jumlah_disetujui": 0
            }
          }
        ]
      }
    ]
  }
}
```

#### 2. Ketersediaan Jam

**Endpoint:** `GET /simnikah/ketersediaan-jam`

**Query Parameters:**
- `tanggal` (required): YYYY-MM-DD

**Request:**
```javascript
const response = await apiClient.get('/simnikah/ketersediaan-jam', {
  params: {
    tanggal: '2024-12-15'
  }
});
```

#### 3. Pernikahan per Tanggal

**Endpoint:** `GET /simnikah/pernikahan-tanggal`

**Query Parameters:**
- `tanggal` (required): YYYY-MM-DD

**Request:**
```javascript
const response = await apiClient.get('/simnikah/pernikahan-tanggal', {
  params: {
    tanggal: '2024-12-15'
  }
});
```

---

### 👰 Catin (User Biasa) Endpoints

#### 1. Buat Pendaftaran Nikah

**Endpoint:** `POST /simnikah/pendaftaran`

**Auth Required:** ✅ Yes

**Request Body:**
```javascript
// Untuk nikah di KUA
const pendaftaranData = {
  calon_laki_laki: {
    nama_dan_bin: "Ahmad Wijaya bin Abdullah",
    pendidikan_akhir: "S1",
    umur: 25
  },
  calon_perempuan: {
    nama_dan_binti: "Siti Nurhaliza binti Muhammad",
    pendidikan_akhir: "S1",
    umur: 23
  },
  lokasi_nikah: {
    tempat_nikah: "Di KUA",
    tanggal_nikah: "2024-12-15",
    waktu_nikah: "10:00"
  },
  wali_nikah: {
    nama_dan_bin: "Abdullah bin Muhammad",
    hubungan_wali: "Ayah Kandung"
  }
};

// Untuk nikah di luar KUA
const pendaftaranDataLuar = {
  calon_laki_laki: {
    nama_dan_bin: "Ahmad Wijaya bin Abdullah",
    pendidikan_akhir: "S1",
    umur: 25
  },
  calon_perempuan: {
    nama_dan_binti: "Siti Nurhaliza binti Muhammad",
    pendidikan_akhir: "S1",
    umur: 23
  },
  lokasi_nikah: {
    tempat_nikah: "Di Luar KUA",
    tanggal_nikah: "2024-12-15",
    waktu_nikah: "10:00",
    alamat_nikah: "Jl. Ahmad Yani No. 123",
    alamat_detail: "Rumah Pengantin Perempuan",
    kelurahan: "Pangeran"
  },
  wali_nikah: {
    nama_dan_bin: "Abdullah bin Muhammad",
    hubungan_wali: "Ayah Kandung"
  }
};

const response = await apiClient.post('/simnikah/pendaftaran', pendaftaranData);
```

**Validasi:**
- Umur minimal: 19 tahun
- Format tanggal: `YYYY-MM-DD`
- Format waktu: `HH:MM` (24 jam)
- Tanggal tidak boleh di masa lalu
- Kelurahan harus dalam Kecamatan Banjarmasin Utara
- Wali nikah wajib diisi

**Hubungan Wali Valid:**
- "Ayah Kandung"
- "Kakek"
- "Saudara Laki-Laki Kandung"
- "Saudara Laki-Laki Seayah"
- "Keponakan Laki-Laki"
- "Paman Kandung"
- "Paman Seayah"
- "Sepupu Laki-Laki"
- "Wali Hakim"
- "Lainnya"

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
    "alamat_akad": "PH5Q+F8C, Jl. Wira Karya...",
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
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    },
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 2. Cek Status Pendaftaran

**Endpoint:** `GET /simnikah/pendaftaran/status`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const response = await apiClient.get('/simnikah/pendaftaran/status');
```

**Response:**
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
      "alamat": "Jl. Penghulu No. 123",
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

**Jika belum ada pendaftaran (404):**
```json
{
  "success": true,
  "message": "User belum memiliki pendaftaran nikah",
  "data": null
}
```

#### 3. Get Detail Pendaftaran by ID

**Endpoint:** `GET /simnikah/pendaftaran/:id`

**Auth Required:** ✅ Yes

**Role Access:**
- `user_biasa`: Hanya bisa melihat pendaftaran miliknya sendiri
- `staff`, `penghulu`, `kepala_kua`: Bisa melihat semua pendaftaran

**Request:**
```javascript
const registrationId = 1;
const response = await apiClient.get(`/simnikah/pendaftaran/${registrationId}`);
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Detail pendaftaran berhasil diambil",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "pendaftar_id": "USR1704067200",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "tanggal_pendaftaran": "2024-01-01T10:00:00Z",
    "tanggal_nikah": "2024-12-15T00:00:00Z",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Ahmad Yani No. 123, Banjarmasin",
    "latitude": -3.291304,
    "longitude": 114.588147,
    "catatan": "Catatan tambahan",
    "disetujui_oleh": "USR1704067201",
    "disetujui_pada": "2024-01-02T10:00:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z",
    "calon_suami": {
      "id": 1,
      "user_id": "CP1704067200",
      "nik": "6301010101010001",
      "nama_lengkap": "Ahmad Wijaya bin Abdullah",
      "tanggal_lahir": "1999-01-01T00:00:00Z",
      "jenis_kelamin": "L",
      "pendidikan_terakhir": "S1",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "calon_istri": {
      "id": 2,
      "user_id": "CP1704067201",
      "nik": "6301010101010002",
      "nama_lengkap": "Siti Nurhaliza binti Muhammad",
      "tanggal_lahir": "2001-01-01T00:00:00Z",
      "jenis_kelamin": "P",
      "pendidikan_terakhir": "S1",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "wali_nikah": {
      "id": 1,
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "penghulu": {
      "id": 1,
      "user_id": "PEN1704067202",
      "nip": "198001012003121001",
      "nama_lengkap": "H. Muhammad Amin",
      "no_hp": "081234567890",
      "email": "amin@kua.go.id",
      "alamat": "Jl. Penghulu No. 123, Banjarmasin",
      "status": "Aktif",
      "ditugaskan_oleh": "USR1704067203",
      "ditugaskan_pada": "2024-12-10T10:00:00Z",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "location": {
      "latitude": -3.291304,
      "longitude": 114.588147,
      "has_coordinates": true,
      "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291304,114.588147",
      "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291304,114.588147",
      "waze_url": "https://www.waze.com/ul?ll=-3.291304,114.588147&navigate=yes",
      "osm_url": "https://www.openstreetmap.org/?mlat=-3.291304&mlon=114.588147&zoom=16"
    }
  }
}
```

**Response Error (404) - Pendaftaran tidak ditemukan:**
```json
{
  "success": false,
  "message": "Pendaftaran tidak ditemukan",
  "error": "Pendaftaran dengan ID tersebut tidak ditemukan",
  "type": "not_found"
}
```

**Response Error (403) - Akses ditolak (user_biasa mencoba akses pendaftaran orang lain):**
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak memiliki akses untuk melihat pendaftaran ini",
  "type": "authorization"
}
```

**Use Cases:**
- Halaman detail pendaftaran untuk staff/penghulu/kepala KUA
- Verifikasi detail sebelum approve
- Menampilkan informasi lengkap untuk penghulu sebelum melaksanakan nikah
- User biasa melihat detail pendaftaran miliknya sendiri

#### 4. Buat Feedback Pernikahan

**Endpoint:** `POST /simnikah/feedback-pernikahan`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const feedbackData = {
  pendaftaran_id: 1,
  jenis_feedback: "Rating", // Rating, Saran, Kritik, Laporan
  rating: 5, // 1-5, hanya untuk jenis "Rating"
  judul: "Pelayanan Sangat Baik",
  pesan: "Terima kasih untuk pelayanan yang sangat baik. Proses nikah berjalan lancar."
};

const response = await apiClient.post('/simnikah/feedback-pernikahan', feedbackData);
```

---

### 👔 Staff KUA Endpoints

#### 1. List Semua Pendaftaran

**Endpoint:** `GET /simnikah/pendaftaran`

**Auth Required:** ✅ Yes  
**Role Required:** `staff`, `kepala_kua`

**Query Parameters:**
- `status` (optional): Filter by status
- `page` (optional): Default 1
- `limit` (optional): Default 10

**Request:**
```javascript
const response = await apiClient.get('/simnikah/pendaftaran', {
  params: {
    status: 'Disetujui',
    page: 1,
    limit: 10
  }
});
```

#### 2. Verifikasi Formulir

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`

**Auth Required:** ✅ Yes  
**Role Required:** `staff`

**Request:**
```javascript
const response = await apiClient.post(`/simnikah/staff/verify-formulir/${id}`, {
  status: "Formulir Disetujui",
  catatan: "Formulir sudah lengkap dan valid"
});
```

#### 3. Approve Pendaftaran

**Endpoint:** `POST /simnikah/staff/approve/:id`

**Auth Required:** ✅ Yes  
**Role Required:** `staff`

**Request:**
```javascript
const response = await apiClient.post(`/simnikah/staff/approve/${id}`, {
  status: "Disetujui",
  catatan: "Pendaftaran disetujui"
});
```

#### 4. Buat Pendaftaran untuk User (Staff)

**Endpoint:** `POST /simnikah/staff/pendaftaran`

**Auth Required:** ✅ Yes  
**Role Required:** `staff`, `kepala_kua`

**Request:** Sama seperti endpoint pendaftaran biasa

**Catatan:** Pendaftaran yang dibuat oleh staff otomatis berstatus "Disetujui"

#### 5. List Pengumuman Nikah

**Endpoint:** `GET /simnikah/staff/pengumuman-nikah/list`

**Auth Required:** ✅ Yes  
**Role Required:** `staff`, `kepala_kua`

**Query Parameters:**
- `tanggal_awal` (optional): YYYY-MM-DD
- `tanggal_akhir` (optional): YYYY-MM-DD

**Request:**
```javascript
const response = await apiClient.get('/simnikah/staff/pengumuman-nikah/list', {
  params: {
    tanggal_awal: '2024-12-16',
    tanggal_akhir: '2024-12-22'
  }
});
```

---

### 👨‍⚖️ Penghulu Endpoints

#### 1. List Tugas Saya

**Endpoint:** `GET /simnikah/penghulu/assigned-registrations`

**Auth Required:** ✅ Yes  
**Role Required:** `penghulu`

**Request:**
```javascript
const response = await apiClient.get('/simnikah/penghulu/assigned-registrations');
```

**Response:**
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
        "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291304,114.588147",
        "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291304,114.588147",
        "waze_url": "https://www.waze.com/ul?ll=-3.291304,114.588147&navigate=yes",
        "osm_url": "https://www.openstreetmap.org/?mlat=-3.291304&mlon=114.588147&zoom=16"
      }
    ],
    "total": 1
  }
}
```

#### 2. Selesaikan Pernikahan

**Endpoint:** `POST /simnikah/penghulu/complete-marriage/:id`

**Auth Required:** ✅ Yes  
**Role Required:** `penghulu`

**Request:**
```javascript
const response = await apiClient.post(`/simnikah/penghulu/complete-marriage/${id}`, {
  catatan: "Pernikahan sudah dilaksanakan dengan lancar"
});
```

---

### 👔 Kepala KUA Endpoints

#### 1. Assign Penghulu

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`

**Auth Required:** ✅ Yes  
**Role Required:** `kepala_kua`

**Request:**
```javascript
const response = await apiClient.post(`/simnikah/pendaftaran/${id}/assign-penghulu`, {
  penghulu_id: 1,
  catatan: "Penghulu ditugaskan untuk menikahkan pasangan ini"
});
```

#### 2. List Penghulu Tersedia

**Endpoint:** `GET /simnikah/kepala-kua/available-penghulu`

**Auth Required:** ✅ Yes  
**Role Required:** `kepala_kua`

**Request:**
```javascript
const response = await apiClient.get('/simnikah/kepala-kua/available-penghulu');
```

#### 3. Buat Staff

**Endpoint:** `POST /simnikah/kepala-kua/staff`

**Auth Required:** ✅ Yes  
**Role Required:** `kepala_kua`

**Request:**
```javascript
const staffData = {
  username: "staff001",
  email: "staff001@kua.go.id",
  password: "password123",
  nama: "Budi Santoso",
  nip: "198001012003121001",
  jabatan: "Staff Pendaftaran",
  no_hp: "081234567890",
  alamat: "Jl. Staff No. 123"
};

const response = await apiClient.post('/simnikah/kepala-kua/staff', staffData);
```

#### 4. Buat Penghulu

**Endpoint:** `POST /simnikah/kepala-kua/penghulu`

**Auth Required:** ✅ Yes  
**Role Required:** `kepala_kua`

**Request:**
```javascript
const penghuluData = {
  username: "penghulu001",
  email: "penghulu001@kua.go.id",
  password: "password123",
  nama: "H. Muhammad Amin",
  nip: "198001012003121002",
  no_hp: "081234567891",
  alamat: "Jl. Penghulu No. 123"
};

const response = await apiClient.post('/simnikah/kepala-kua/penghulu', penghuluData);
```

#### 5. Dashboard Kepala KUA

**Endpoint:** `GET /simnikah/dashboard/kepala-kua`

**Auth Required:** ✅ Yes  
**Role Required:** `kepala_kua`

**Query Parameters:**
- `period` (optional): `day`, `week`, `month`, `year` (default: `month`)
- `date_from` (optional): YYYY-MM-DD
- `date_to` (optional): YYYY-MM-DD

**Request:**
```javascript
const response = await apiClient.get('/simnikah/dashboard/kepala-kua', {
  params: {
    period: 'month'
  }
});
```

#### 6. List Feedback

**Endpoint:** `GET /simnikah/kepala-kua/feedback`

**Auth Required:** ✅ Yes  
**Role Required:** `kepala_kua`

**Query Parameters:**
- `jenis_feedback` (optional): Rating, Saran, Kritik, Laporan
- `status_baca` (optional): Belum Dibaca, Sudah Dibaca
- `page` (optional)
- `limit` (optional)

**Request:**
```javascript
const response = await apiClient.get('/simnikah/kepala-kua/feedback', {
  params: {
    status_baca: 'Belum Dibaca',
    page: 1,
    limit: 10
  }
});
```

---

### 📍 Location Endpoints

#### 1. Geocode (Alamat ke Koordinat)

**Endpoint:** `POST /simnikah/location/geocode`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const response = await apiClient.post('/simnikah/location/geocode', {
  alamat: "Jl. Ahmad Yani No. 123, Banjarmasin"
});
```

**Response:**
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

#### 2. Reverse Geocode (Koordinat ke Alamat)

**Endpoint:** `POST /simnikah/location/reverse-geocode`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const response = await apiClient.post('/simnikah/location/reverse-geocode', {
  latitude: -3.291304,
  longitude: 114.588147
});
```

#### 3. Detail Lokasi Pernikahan

**Endpoint:** `GET /simnikah/pendaftaran/:id/location`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const response = await apiClient.get(`/simnikah/pendaftaran/${id}/location`);
```

---

### 🔔 Notification Endpoints

#### 1. List Notifikasi User

**Endpoint:** `GET /simnikah/notifikasi/user/:user_id`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const userId = localStorage.getItem('user_id');
const response = await apiClient.get(`/simnikah/notifikasi/user/${userId}`);
```

#### 2. Update Status Notifikasi

**Endpoint:** `PUT /simnikah/notifikasi/:id/status`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const response = await apiClient.put(`/simnikah/notifikasi/${id}/status`, {
  status_baca: "Sudah Dibaca"
});
```

#### 3. Mark All as Read

**Endpoint:** `PUT /simnikah/notifikasi/user/:user_id/mark-all-read`

**Auth Required:** ✅ Yes

**Request:**
```javascript
const userId = localStorage.getItem('user_id');
const response = await apiClient.put(`/simnikah/notifikasi/user/${userId}/mark-all-read`);
```

---

## 📦 Request & Response Format

### Standard Request Headers

```javascript
{
  'Content-Type': 'application/json',
  'Authorization': 'Bearer <token>', // Untuk protected endpoints
  'Accept': 'application/json'
}
```

### Standard Success Response

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": {
    // Data response
  }
}
```

### Standard Error Response

```json
{
  "success": false,
  "message": "Pesan error umum",
  "error": "Detail error",
  "type": "validation", // atau authentication, authorization, dll
  "field": "nama_field" // optional, jika error validasi field spesifik
}
```

---

## ❌ Error Handling

### HTTP Status Codes

| Code | Description | Action |
|------|-------------|--------|
| `200` | Success | Continue normal flow |
| `201` | Created | Resource berhasil dibuat |
| `400` | Bad Request | Validasi error, tampilkan pesan error |
| `401` | Unauthorized | Token invalid/expired, redirect ke login |
| `403` | Forbidden | Tidak punya permission |
| `404` | Not Found | Resource tidak ditemukan |
| `429` | Too Many Requests | Rate limit exceeded, tunggu sebentar |
| `500` | Internal Server Error | Error server, tampilkan pesan umum |

### Error Types

- `authentication` - Error authentication (token invalid/expired)
- `authorization` - Error authorization/permission
- `validation` - Error validasi input
- `not_found` - Resource tidak ditemukan
- `database` - Error database
- `format` - Error format data
- `schedule_conflict` - Konflik jadwal (untuk pendaftaran)

### Error Handler Utility

```javascript
// utils/errorHandler.js
export const handleApiError = (error) => {
  if (!error.response) {
    return {
      message: 'Tidak dapat terhubung ke server',
      type: 'network'
    };
  }

  const { status, data } = error.response;

  switch (status) {
    case 400:
      return {
        message: data.message || 'Data tidak valid',
        error: data.error,
        type: data.type || 'validation',
        field: data.field
      };
    
    case 401:
      // Redirect ke login
      localStorage.removeItem('token');
      window.location.href = '/login';
      return {
        message: 'Sesi Anda telah berakhir, silakan login kembali',
        type: 'authentication'
      };
    
    case 403:
      return {
        message: 'Anda tidak memiliki akses untuk melakukan aksi ini',
        type: 'authorization'
      };
    
    case 404:
      return {
        message: 'Data tidak ditemukan',
        type: 'not_found'
      };
    
    case 429:
      return {
        message: 'Terlalu banyak request, silakan tunggu sebentar',
        type: 'rate_limit'
      };
    
    case 500:
      return {
        message: 'Terjadi kesalahan pada server',
        type: 'server'
      };
    
    default:
      return {
        message: data.message || 'Terjadi kesalahan',
        type: 'unknown'
      };
  }
};
```

### Contoh Penggunaan

```javascript
try {
  const response = await apiClient.post('/simnikah/pendaftaran', data);
  // Success handling
} catch (error) {
  const errorInfo = handleApiError(error);
  
  // Tampilkan error ke user
  if (errorInfo.type === 'validation') {
    // Tampilkan error validasi spesifik
    setFieldError(errorInfo.field, errorInfo.error);
  } else {
    // Tampilkan error umum
    toast.error(errorInfo.message);
  }
}
```

---

## 💻 Contoh Implementasi

### 1. Login Component (React)

```jsx
import React, { useState } from 'react';
import apiClient from '../utils/api';
import { useNavigate } from 'react-router-dom';

const LoginPage = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await apiClient.post('/login', formData);
      
      // Simpan token dan user info
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      
      // Redirect berdasarkan role
      const role = response.user.role;
      switch (role) {
        case 'user_biasa':
          navigate('/dashboard/catin');
          break;
        case 'staff':
          navigate('/dashboard/staff');
          break;
        case 'penghulu':
          navigate('/dashboard/penghulu');
          break;
        case 'kepala_kua':
          navigate('/dashboard/kepala-kua');
          break;
        default:
          navigate('/');
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Login gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Username"
        value={formData.username}
        onChange={(e) => setFormData({ ...formData, username: e.target.value })}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
        required
      />
      {error && <div className="error">{error}</div>}
      <button type="submit" disabled={loading}>
        {loading ? 'Loading...' : 'Login'}
      </button>
    </form>
  );
};

export default LoginPage;
```

### 2. Calendar Component (React)

```jsx
import React, { useState, useEffect } from 'react';
import apiClient from '../utils/api';

const CalendarAvailability = () => {
  const [calendar, setCalendar] = useState(null);
  const [month, setMonth] = useState(new Date().getMonth() + 1);
  const [year, setYear] = useState(new Date().getFullYear());
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchCalendar();
  }, [month, year]);

  const fetchCalendar = async () => {
    setLoading(true);
    try {
      const response = await apiClient.get('/simnikah/kalender-ketersediaan', {
        params: { bulan: month, tahun: year }
      });
      setCalendar(response.data);
    } catch (error) {
      console.error('Error fetching calendar:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (day) => {
    if (day.is_past) return 'gray';
    if (!day.tersedia) return 'red';
    if (day.sisa_kuota < 3) return 'yellow';
    return 'green';
  };

  if (loading) return <div>Loading...</div>;
  if (!calendar) return null;

  return (
    <div>
      <div>
        <button onClick={() => setMonth(month - 1)}>Prev</button>
        <span>{calendar.nama_bulan} {calendar.tahun}</span>
        <button onClick={() => setMonth(month + 1)}>Next</button>
      </div>
      
      <div className="calendar-grid">
        {calendar.calendar.map((day) => (
          <div
            key={day.tanggal}
            className={`day ${getStatusColor(day)}`}
          >
            <div>{day.tanggal}</div>
            <div>{day.status}</div>
            <div>Sisa: {day.sisa_kuota}/{day.kapasitas}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CalendarAvailability;
```

### 3. Form Pendaftaran (React dengan Multi-Step)

```jsx
import React, { useState } from 'react';
import apiClient from '../utils/api';

const RegistrationForm = () => {
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState({
    calon_laki_laki: {
      nama_dan_bin: '',
      pendidikan_akhir: '',
      umur: ''
    },
    calon_perempuan: {
      nama_dan_binti: '',
      pendidikan_akhir: '',
      umur: ''
    },
    lokasi_nikah: {
      tempat_nikah: 'Di KUA',
      tanggal_nikah: '',
      waktu_nikah: '',
      alamat_nikah: '',
      alamat_detail: '',
      kelurahan: ''
    },
    wali_nikah: {
      nama_dan_bin: '',
      hubungan_wali: ''
    }
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    setLoading(true);
    setError('');

    try {
      const response = await apiClient.post('/simnikah/pendaftaran', formData);
      alert('Pendaftaran berhasil!');
      // Redirect atau reset form
    } catch (err) {
      setError(err.response?.data?.error || 'Pendaftaran gagal');
      if (err.response?.data?.field) {
        // Scroll ke field yang error
        const fieldElement = document.querySelector(`[name="${err.response.data.field}"]`);
        if (fieldElement) fieldElement.scrollIntoView();
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <form>
      {step === 1 && (
        <div>
          <h2>Data Calon Suami</h2>
          <input
            name="nama_dan_bin"
            placeholder="Nama dan Bin"
            value={formData.calon_laki_laki.nama_dan_bin}
            onChange={(e) => setFormData({
              ...formData,
              calon_laki_laki: {
                ...formData.calon_laki_laki,
                nama_dan_bin: e.target.value
              }
            })}
          />
          {/* Field lainnya */}
          <button onClick={() => setStep(2)}>Next</button>
        </div>
      )}
      
      {step === 2 && (
        <div>
          <h2>Data Calon Istri</h2>
          {/* Form fields */}
          <button onClick={() => setStep(1)}>Back</button>
          <button onClick={() => setStep(3)}>Next</button>
        </div>
      )}
      
      {step === 3 && (
        <div>
          <h2>Data Wali Nikah</h2>
          {/* Form fields */}
          <button onClick={() => setStep(2)}>Back</button>
          <button onClick={() => setStep(4)}>Next</button>
        </div>
      )}
      
      {step === 4 && (
        <div>
          <h2>Jadwal & Lokasi</h2>
          {/* Form fields */}
          {error && <div className="error">{error}</div>}
          <button onClick={() => setStep(3)}>Back</button>
          <button onClick={handleSubmit} disabled={loading}>
            {loading ? 'Mengirim...' : 'Kirim Pendaftaran'}
          </button>
        </div>
      )}
    </form>
  );
};

export default RegistrationForm;
```

### 4. Upload Photo Component (React)

```jsx
import React, { useState, useRef } from 'react';
import apiClient from '../utils/api';

const UploadPhoto = ({ onUploadSuccess }) => {
  const [file, setFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const fileInputRef = useRef(null);

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    
    // Validasi file size (5MB)
    if (selectedFile.size > 5 * 1024 * 1024) {
      setError('Ukuran file maksimal 5MB');
      return;
    }
    
    // Validasi file type
    const validTypes = ['image/jpeg', 'image/png', 'image/jpg', 'image/webp'];
    if (!validTypes.includes(selectedFile.type)) {
      setError('Format file tidak didukung. Gunakan JPG, PNG, atau WEBP');
      return;
    }
    
    setFile(selectedFile);
    setError('');
    
    // Preview image
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreview(reader.result);
    };
    reader.readAsDataURL(selectedFile);
  };

  const handleUpload = async () => {
    if (!file) return;
    
    setLoading(true);
    setError('');
    
    try {
      const formData = new FormData();
      formData.append('photo', file);
      
      const response = await apiClient.post('/upload-photo', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      
      alert('Foto berhasil diupload!');
      if (onUploadSuccess) {
        onUploadSuccess(response.data.profile_photo);
      }
      
      // Reset
      setFile(null);
      setPreview(null);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Upload gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/jpg,image/webp"
        onChange={handleFileChange}
        style={{ display: 'none' }}
      />
      
      <button onClick={() => fileInputRef.current?.click()}>
        Pilih Foto
      </button>
      
      {preview && (
        <div>
          <img src={preview} alt="Preview" style={{ maxWidth: '200px' }} />
          <button onClick={handleUpload} disabled={loading}>
            {loading ? 'Uploading...' : 'Upload'}
          </button>
        </div>
      )}
      
      {error && <div className="error">{error}</div>}
    </div>
  );
};

export default UploadPhoto;
```

---

## ✅ Best Practices

### 1. Token Management

- **Jangan hardcode token** di code
- **Gunakan interceptor** untuk otomatis menambahkan token
- **Refresh token** sebelum expired (jika ada mekanisme refresh)
- **Hapus token** saat logout

### 2. Error Handling

- **Selalu handle error** di setiap API call
- **Tampilkan pesan error yang user-friendly**
- **Log error** untuk debugging (tapi jangan tampilkan ke user)
- **Handle network errors** (tidak ada koneksi)

### 3. Loading States

- **Tampilkan loading indicator** saat fetch data
- **Disable button** saat submit form
- **Gunakan skeleton loader** untuk UX yang lebih baik

### 4. Form Validation

- **Validasi client-side** sebelum submit (untuk UX)
- **Jangan hanya andalkan client-side validation** (server selalu validasi)
- **Tampilkan error di field yang sesuai**
- **Gunakan library validasi** seperti Yup, Zod, atau Formik

### 5. State Management

- **Gunakan state management** untuk data global (Redux, Zustand, Context)
- **Cache API responses** jika memungkinkan
- **Invalidate cache** saat data berubah

### 6. Security

- **Jangan expose sensitive data** di frontend
- **Sanitize user input** sebelum display
- **Gunakan HTTPS** di production
- **Hapus token** saat logout

### 7. Performance

- **Lazy load** komponen yang tidak critical
- **Debounce** untuk search/autocomplete
- **Pagination** untuk list data panjang
- **Optimize images** sebelum upload

### 8. Code Organization

```
src/
├── api/
│   └── client.js          # API client setup
├── components/
│   ├── Auth/
│   ├── Calendar/
│   ├── Forms/
│   └── ...
├── hooks/
│   ├── useAuth.js
│   ├── useApi.js
│   └── ...
├── pages/
│   ├── Login/
│   ├── Dashboard/
│   └── ...
├── utils/
│   ├── errorHandler.js
│   └── validators.js
└── constants/
    └── status.js
```

---

## ✅ Validasi Tampilan Frontend (Wajib Diimplementasikan)

**PENTING:** Semua validasi di bawah ini **WAJIB** diimplementasikan di frontend untuk:
1. **Keamanan** - Mencegah input berbahaya
2. **UX yang baik** - Memberikan feedback langsung ke user
3. **Konsistensi dengan backend** - Memastikan data yang dikirim sesuai dengan validasi backend
4. **Mengurangi beban server** - Validasi client-side mengurangi request yang tidak valid

> ⚠️ **CATATAN PENTING:** Validasi frontend **TIDAK menggantikan** validasi backend. Backend akan selalu melakukan validasi ulang. Validasi frontend hanya untuk UX dan mengurangi request yang tidak perlu.

---

### 1. Validasi Form Pendaftaran Nikah

#### 1.1 Validasi Calon Laki-Laki

```javascript
// utils/validators/registration.js

export const validateCalonLakiLaki = (data) => {
  const errors = {};

  // Nama dan Bin - Required, min 3 karakter, max 100 karakter
  if (!data.nama_dan_bin || data.nama_dan_bin.trim() === '') {
    errors.nama_dan_bin = 'Nama dan bin wajib diisi';
  } else if (data.nama_dan_bin.trim().length < 3) {
    errors.nama_dan_bin = 'Nama minimal 3 karakter';
  } else if (data.nama_dan_bin.length > 100) {
    errors.nama_dan_bin = 'Nama maksimal 100 karakter';
  } else if (!data.nama_dan_bin.includes('bin')) {
    errors.nama_dan_bin = 'Format nama harus menggunakan "bin" (contoh: Ahmad bin Abdullah)';
  }

  // Umur - Required, minimal 19 tahun
  if (!data.umur || data.umur === '') {
    errors.umur = 'Umur wajib diisi';
  } else {
    const umur = parseInt(data.umur);
    if (isNaN(umur) || umur < 19) {
      errors.umur = 'Umur calon laki-laki minimal 19 tahun';
    } else if (umur > 100) {
      errors.umur = 'Umur tidak valid';
    }
  }

  // Pendidikan - Optional tapi jika diisi harus valid
  if (data.pendidikan_akhir && data.pendidikan_akhir.length > 50) {
    errors.pendidikan_akhir = 'Pendidikan maksimal 50 karakter';
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

#### 1.2 Validasi Calon Perempuan

```javascript
export const validateCalonPerempuan = (data) => {
  const errors = {};

  // Nama dan Binti - Required, min 3 karakter, max 100 karakter
  if (!data.nama_dan_binti || data.nama_dan_binti.trim() === '') {
    errors.nama_dan_binti = 'Nama dan binti wajib diisi';
  } else if (data.nama_dan_binti.trim().length < 3) {
    errors.nama_dan_binti = 'Nama minimal 3 karakter';
  } else if (data.nama_dan_binti.length > 100) {
    errors.nama_dan_binti = 'Nama maksimal 100 karakter';
  } else if (!data.nama_dan_binti.includes('binti')) {
    errors.nama_dan_binti = 'Format nama harus menggunakan "binti" (contoh: Siti binti Muhammad)';
  }

  // Umur - Required, minimal 19 tahun
  if (!data.umur || data.umur === '') {
    errors.umur = 'Umur wajib diisi';
  } else {
    const umur = parseInt(data.umur);
    if (isNaN(umur) || umur < 19) {
      errors.umur = 'Umur calon perempuan minimal 19 tahun';
    } else if (umur > 100) {
      errors.umur = 'Umur tidak valid';
    }
  }

  // Pendidikan - Optional tapi jika diisi harus valid
  if (data.pendidikan_akhir && data.pendidikan_akhir.length > 50) {
    errors.pendidikan_akhir = 'Pendidikan maksimal 50 karakter';
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

#### 1.3 Validasi Wali Nikah

```javascript
export const validateWaliNikah = (data) => {
  const errors = {};

  // Nama dan Bin - Required
  if (!data.nama_dan_bin || data.nama_dan_bin.trim() === '') {
    errors.nama_dan_bin = 'Nama wali nikah wajib diisi';
  } else if (data.nama_dan_bin.trim().length < 3) {
    errors.nama_dan_bin = 'Nama minimal 3 karakter';
  } else if (data.nama_dan_bin.length > 100) {
    errors.nama_dan_bin = 'Nama maksimal 100 karakter';
  } else if (!data.nama_dan_bin.includes('bin')) {
    errors.nama_dan_bin = 'Format nama harus menggunakan "bin"';
  }

  // Hubungan Wali - Required, harus dari daftar valid
  const validHubunganWali = [
    'Ayah Kandung',
    'Kakek',
    'Saudara Laki-Laki Kandung',
    'Saudara Laki-Laki Seayah',
    'Keponakan Laki-Laki',
    'Paman Kandung',
    'Paman Seayah',
    'Sepupu Laki-Laki',
    'Wali Hakim',
    'Lainnya'
  ];

  if (!data.hubungan_wali || data.hubungan_wali.trim() === '') {
    errors.hubungan_wali = 'Hubungan wali wajib diisi';
  } else if (!validHubunganWali.includes(data.hubungan_wali)) {
    errors.hubungan_wali = 'Hubungan wali tidak valid. Pilih dari daftar yang tersedia';
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

#### 1.4 Validasi Lokasi & Jadwal Nikah

```javascript
export const validateLokasiNikah = (data, availableSlots = null) => {
  const errors = {};

  // Tempat Nikah - Required, harus "Di KUA" atau "Di Luar KUA"
  const validTempat = ['Di KUA', 'Di Luar KUA'];
  if (!data.tempat_nikah || !validTempat.includes(data.tempat_nikah)) {
    errors.tempat_nikah = 'Pilih tempat nikah: Di KUA atau Di Luar KUA';
  }

  // Tanggal Nikah - Required, format YYYY-MM-DD, tidak boleh di masa lalu
  if (!data.tanggal_nikah || data.tanggal_nikah.trim() === '') {
    errors.tanggal_nikah = 'Tanggal nikah wajib diisi';
  } else {
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
    if (!dateRegex.test(data.tanggal_nikah)) {
      errors.tanggal_nikah = 'Format tanggal harus YYYY-MM-DD (contoh: 2024-12-15)';
    } else {
      const selectedDate = new Date(data.tanggal_nikah);
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      
      if (isNaN(selectedDate.getTime())) {
        errors.tanggal_nikah = 'Tanggal tidak valid';
      } else if (selectedDate < today) {
        errors.tanggal_nikah = 'Tanggal nikah tidak boleh di masa lalu';
      }
      
      // Cek hari libur (Minggu)
      if (selectedDate.getDay() === 0) {
        errors.tanggal_nikah = 'Tidak bisa memilih hari Minggu';
      }
    }
  }

  // Waktu Nikah - Required, format HH:MM, dalam range 08:00-16:00
  if (!data.waktu_nikah || data.waktu_nikah.trim() === '') {
    errors.waktu_nikah = 'Waktu nikah wajib diisi';
  } else {
    const timeRegex = /^([0-1][0-9]|2[0-3]):[0-5][0-9]$/;
    if (!timeRegex.test(data.waktu_nikah)) {
      errors.waktu_nikah = 'Format waktu harus HH:MM (contoh: 09:00, 14:30)';
    } else {
      const [hours, minutes] = data.waktu_nikah.split(':').map(Number);
      const timeInMinutes = hours * 60 + minutes;
      const startTime = 8 * 60; // 08:00
      const endTime = 16 * 60; // 16:00
      
      if (timeInMinutes < startTime || timeInMinutes > endTime) {
        errors.waktu_nikah = 'Waktu nikah harus antara 08:00 - 16:00';
      }
      
      // Cek apakah waktu dalam time slots yang valid
      const validTimeSlots = [
        '08:00', '09:00', '10:00', '11:00', '12:00',
        '13:00', '14:00', '15:00', '16:00'
      ];
      if (!validTimeSlots.includes(data.waktu_nikah)) {
        errors.waktu_nikah = 'Waktu harus sesuai dengan slot yang tersedia (08:00, 09:00, ..., 16:00)';
      }
    }
  }

  // Jika nikah di luar KUA, validasi alamat
  if (data.tempat_nikah === 'Di Luar KUA') {
    if (!data.alamat_nikah || data.alamat_nikah.trim() === '') {
      errors.alamat_nikah = 'Alamat nikah wajib diisi untuk nikah di luar KUA';
    } else if (data.alamat_nikah.trim().length < 10) {
      errors.alamat_nikah = 'Alamat minimal 10 karakter';
    } else if (data.alamat_nikah.length > 200) {
      errors.alamat_nikah = 'Alamat maksimal 200 karakter';
    }

    if (data.alamat_detail && data.alamat_detail.length > 200) {
      errors.alamat_detail = 'Detail alamat maksimal 200 karakter';
    }

    // Kelurahan - Required, harus dalam daftar valid
    const validKelurahan = [
      'Alalak Utara',
      'Alalak Tengah',
      'Alalak Selatan',
      'Antasan Kecil Timur',
      'Kuin Utara',
      'Pangeran',
      'Sungai Miai',
      'Sungai Andai',
      'Surgi Mufti'
    ];

    if (!data.kelurahan || data.kelurahan.trim() === '') {
      errors.kelurahan = 'Kelurahan wajib diisi';
    } else if (!validKelurahan.includes(data.kelurahan)) {
      errors.kelurahan = 'Kelurahan harus dalam Kecamatan Banjarmasin Utara';
    }
  }

  // Validasi ketersediaan slot (jika data availableSlots tersedia)
  if (availableSlots && data.tanggal_nikah && data.waktu_nikah) {
    const slot = availableSlots.find(
      s => s.waktu === data.waktu_nikah
    );
    
    if (slot) {
      if (data.tempat_nikah === 'Di KUA') {
        if (!slot.kua.tersedia || slot.kua.terbooking) {
          errors.waktu_nikah = 'Slot waktu di KUA sudah terbooking';
        }
      } else {
        if (!slot.luar_kua.tersedia || slot.luar_kua.terbooking) {
          errors.waktu_nikah = 'Slot waktu di luar KUA sudah penuh (maksimal 3 per jam)';
        }
      }
    }
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

---

### 2. Validasi Authentication & Authorization

#### 2.1 Validasi Login

```javascript
// utils/validators/auth.js

export const validateLogin = (data) => {
  const errors = {};

  // Username - Required, min 3 karakter
  if (!data.username || data.username.trim() === '') {
    errors.username = 'Username wajib diisi';
  } else if (data.username.trim().length < 3) {
    errors.username = 'Username minimal 3 karakter';
  } else if (data.username.length > 50) {
    errors.username = 'Username maksimal 50 karakter';
  } else if (!/^[a-zA-Z0-9_]+$/.test(data.username)) {
    errors.username = 'Username hanya boleh mengandung huruf, angka, dan underscore';
  }

  // Password - Required, min 6 karakter
  if (!data.password || data.password === '') {
    errors.password = 'Password wajib diisi';
  } else if (data.password.length < 6) {
    errors.password = 'Password minimal 6 karakter';
  } else if (data.password.length > 255) {
    errors.password = 'Password terlalu panjang';
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

#### 2.2 Validasi Register

```javascript
export const validateRegister = (data) => {
  const errors = {};

  // Username
  if (!data.username || data.username.trim() === '') {
    errors.username = 'Username wajib diisi';
  } else if (data.username.trim().length < 3) {
    errors.username = 'Username minimal 3 karakter';
  } else if (data.username.length > 50) {
    errors.username = 'Username maksimal 50 karakter';
  } else if (!/^[a-zA-Z0-9_]+$/.test(data.username)) {
    errors.username = 'Username hanya boleh mengandung huruf, angka, dan underscore';
  }

  // Email - Required, format email valid
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!data.email || data.email.trim() === '') {
    errors.email = 'Email wajib diisi';
  } else if (!emailRegex.test(data.email)) {
    errors.email = 'Format email tidak valid';
  } else if (data.email.length > 100) {
    errors.email = 'Email maksimal 100 karakter';
  }

  // Password
  if (!data.password || data.password === '') {
    errors.password = 'Password wajib diisi';
  } else if (data.password.length < 6) {
    errors.password = 'Password minimal 6 karakter';
  }

  // Nama - Required
  if (!data.nama || data.nama.trim() === '') {
    errors.nama = 'Nama wajib diisi';
  } else if (data.nama.trim().length < 3) {
    errors.nama = 'Nama minimal 3 karakter';
  } else if (data.nama.length > 100) {
    errors.nama = 'Nama maksimal 100 karakter';
  }

  // Role - Required, harus dari daftar valid
  const validRoles = ['user_biasa', 'staff', 'penghulu', 'kepala_kua'];
  if (!data.role || !validRoles.includes(data.role)) {
    errors.role = 'Role tidak valid';
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

#### 2.3 Validasi Role-Based Access

```javascript
// utils/validators/authorization.js

export const checkRoleAccess = (userRole, requiredRoles) => {
  if (!userRole) {
    return {
      hasAccess: false,
      message: 'User tidak terautentikasi'
    };
  }

  if (Array.isArray(requiredRoles)) {
    if (!requiredRoles.includes(userRole)) {
      return {
        hasAccess: false,
        message: `Akses ditolak. Role yang diizinkan: ${requiredRoles.join(', ')}`
      };
    }
  } else {
    if (userRole !== requiredRoles) {
      return {
        hasAccess: false,
        message: `Akses ditolak. Role yang diizinkan: ${requiredRoles}`
      };
    }
  }

  return {
    hasAccess: true,
    message: 'Akses diizinkan'
  };
};

// Contoh penggunaan di component
export const useRoleGuard = (requiredRoles) => {
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  const access = checkRoleAccess(user.role, requiredRoles);
  
  if (!access.hasAccess) {
    // Redirect atau tampilkan error
    return false;
  }
  
  return true;
};
```

---

### 3. Validasi File Upload

```javascript
// utils/validators/file.js

export const validateProfilePhoto = (file) => {
  const errors = [];

  if (!file) {
    errors.push('File wajib dipilih');
    return { isValid: false, errors };
  }

  // Validasi file size (5MB)
  const maxSize = 5 * 1024 * 1024; // 5MB dalam bytes
  if (file.size > maxSize) {
    errors.push('Ukuran file maksimal 5MB');
  }

  // Validasi file type
  const validTypes = [
    'image/jpeg',
    'image/jpg',
    'image/png',
    'image/webp'
  ];
  
  if (!validTypes.includes(file.type)) {
    errors.push('Format file tidak didukung. Gunakan JPG, PNG, atau WEBP');
  }

  // Validasi file extension (double check)
  const validExtensions = ['.jpg', '.jpeg', '.png', '.webp'];
  const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
  if (!validExtensions.includes(fileExtension)) {
    errors.push('Ekstensi file tidak valid');
  }

  return {
    isValid: errors.length === 0,
    errors
  };
};
```

---

### 4. Validasi Input Sanitization (Security)

```javascript
// utils/validators/sanitize.js

export const sanitizeString = (input) => {
  if (typeof input !== 'string') return '';
  
  return input
    .trim()
    .replace(/[<>]/g, '') // Remove potential HTML tags
    .replace(/javascript:/gi, '') // Remove javascript: protocol
    .replace(/on\w+=/gi, '') // Remove event handlers
    .substring(0, 1000); // Limit length
};

export const sanitizeEmail = (email) => {
  if (typeof email !== 'string') return '';
  
  return email.trim().toLowerCase().substring(0, 100);
};

export const sanitizeNumber = (input) => {
  if (typeof input === 'number') return input;
  if (typeof input !== 'string') return null;
  
  const cleaned = input.replace(/[^0-9]/g, '');
  return cleaned === '' ? null : parseInt(cleaned, 10);
};

// XSS Protection untuk display
export const escapeHtml = (text) => {
  if (typeof text !== 'string') return '';
  
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  
  return text.replace(/[&<>"']/g, (m) => map[m]);
};
```

---

### 5. Validasi Business Logic

#### 5.1 Validasi Status Pendaftaran

```javascript
// utils/validators/status.js

export const validateStatusTransition = (currentStatus, newStatus, userRole) => {
  // Status yang valid
  const validStatuses = [
    'Draft',
    'Disetujui',
    'Menunggu Penugasan',
    'Penghulu Ditugaskan',
    'Selesai',
    'Ditolak'
  ];

  if (!validStatuses.includes(newStatus)) {
    return {
      isValid: false,
      message: 'Status tidak valid'
    };
  }

  // Flow status yang diizinkan
  const allowedTransitions = {
    'Draft': ['Disetujui', 'Ditolak'],
    'Disetujui': ['Menunggu Penugasan', 'Ditolak'],
    'Menunggu Penugasan': ['Penghulu Ditugaskan'],
    'Penghulu Ditugaskan': ['Selesai'],
    'Selesai': [], // Final status
    'Ditolak': [] // Final status
  };

  const allowed = allowedTransitions[currentStatus] || [];
  if (!allowed.includes(newStatus)) {
    return {
      isValid: false,
      message: `Tidak bisa mengubah status dari "${currentStatus}" ke "${newStatus}"`
    };
  }

  // Role-based validation
  if (newStatus === 'Disetujui' && userRole !== 'staff') {
    return {
      isValid: false,
      message: 'Hanya staff yang dapat menyetujui pendaftaran'
    };
  }

  if (newStatus === 'Penghulu Ditugaskan' && userRole !== 'kepala_kua') {
    return {
      isValid: false,
      message: 'Hanya kepala KUA yang dapat menugaskan penghulu'
    };
  }

  return {
    isValid: true,
    message: 'Status transition valid'
  };
};
```

#### 5.2 Validasi Ketersediaan Jadwal

```javascript
// utils/validators/schedule.js

export const validateScheduleAvailability = async (tanggal, waktu, tempat, apiClient) => {
  try {
    // Fetch ketersediaan jam untuk tanggal tersebut
    const response = await apiClient.get('/simnikah/ketersediaan-jam', {
      params: { tanggal }
    });

    const timeSlot = response.data.time_slots.find(slot => slot.waktu === waktu);
    
    if (!timeSlot) {
      return {
        isValid: false,
        message: 'Waktu tidak tersedia'
      };
    }

    if (tempat === 'Di KUA') {
      if (!timeSlot.kua.tersedia || timeSlot.kua.terbooking) {
        return {
          isValid: false,
          message: 'Slot waktu di KUA sudah terbooking'
        };
      }
    } else {
      if (!timeSlot.luar_kua.tersedia || timeSlot.luar_kua.terbooking) {
        return {
          isValid: false,
          message: 'Slot waktu di luar KUA sudah penuh (maksimal 3 per jam)'
        };
      }
    }

    return {
      isValid: true,
      message: 'Jadwal tersedia'
    };
  } catch (error) {
    return {
      isValid: false,
      message: 'Gagal memvalidasi ketersediaan jadwal'
    };
  }
};
```

---

### 6. Validasi Feedback

```javascript
// utils/validators/feedback.js

export const validateFeedback = (data) => {
  const errors = {};

  // Pendaftaran ID - Required, harus number
  if (!data.pendaftaran_id) {
    errors.pendaftaran_id = 'ID pendaftaran wajib diisi';
  } else if (isNaN(parseInt(data.pendaftaran_id))) {
    errors.pendaftaran_id = 'ID pendaftaran harus berupa angka';
  }

  // Jenis Feedback - Required, harus dari daftar valid
  const validJenis = ['Rating', 'Saran', 'Kritik', 'Laporan'];
  if (!data.jenis_feedback || !validJenis.includes(data.jenis_feedback)) {
    errors.jenis_feedback = 'Jenis feedback tidak valid';
  }

  // Rating - Required jika jenis "Rating", harus 1-5
  if (data.jenis_feedback === 'Rating') {
    if (!data.rating) {
      errors.rating = 'Rating wajib diisi untuk jenis Rating';
    } else {
      const rating = parseInt(data.rating);
      if (isNaN(rating) || rating < 1 || rating > 5) {
        errors.rating = 'Rating harus antara 1-5';
      }
    }
  }

  // Judul - Required, min 5 karakter, max 200 karakter
  if (!data.judul || data.judul.trim() === '') {
    errors.judul = 'Judul wajib diisi';
  } else if (data.judul.trim().length < 5) {
    errors.judul = 'Judul minimal 5 karakter';
  } else if (data.judul.length > 200) {
    errors.judul = 'Judul maksimal 200 karakter';
  }

  // Pesan - Required, min 10 karakter
  if (!data.pesan || data.pesan.trim() === '') {
    errors.pesan = 'Pesan wajib diisi';
  } else if (data.pesan.trim().length < 10) {
    errors.pesan = 'Pesan minimal 10 karakter';
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
```

---

### 7. Contoh Implementasi Validasi di Component

```jsx
// components/RegistrationForm.jsx
import React, { useState } from 'react';
import { 
  validateCalonLakiLaki, 
  validateCalonPerempuan, 
  validateWaliNikah, 
  validateLokasiNikah 
} from '../utils/validators/registration';
import apiClient from '../utils/api';

const RegistrationForm = () => {
  const [formData, setFormData] = useState({ /* ... */ });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const validateForm = () => {
    const newErrors = {};

    // Validasi calon laki-laki
    const calonLakiLakiValidation = validateCalonLakiLaki(formData.calon_laki_laki);
    if (!calonLakiLakiValidation.isValid) {
      newErrors.calon_laki_laki = calonLakiLakiValidation.errors;
    }

    // Validasi calon perempuan
    const calonPerempuanValidation = validateCalonPerempuan(formData.calon_perempuan);
    if (!calonPerempuanValidation.isValid) {
      newErrors.calon_perempuan = calonPerempuanValidation.errors;
    }

    // Validasi wali nikah
    const waliValidation = validateWaliNikah(formData.wali_nikah);
    if (!waliValidation.isValid) {
      newErrors.wali_nikah = waliValidation.errors;
    }

    // Validasi lokasi (tanpa availableSlots dulu, akan divalidasi lagi sebelum submit)
    const lokasiValidation = validateLokasiNikah(formData.lokasi_nikah);
    if (!lokasiValidation.isValid) {
      newErrors.lokasi_nikah = lokasiValidation.errors;
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validasi client-side
    if (!validateForm()) {
      // Scroll ke error pertama
      const firstErrorField = document.querySelector('.error-field');
      if (firstErrorField) {
        firstErrorField.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
      return;
    }

    // Validasi ketersediaan jadwal sebelum submit
    try {
      const scheduleValidation = await validateScheduleAvailability(
        formData.lokasi_nikah.tanggal_nikah,
        formData.lokasi_nikah.waktu_nikah,
        formData.lokasi_nikah.tempat_nikah,
        apiClient
      );

      if (!scheduleValidation.isValid) {
        setErrors({
          ...errors,
          lokasi_nikah: {
            ...errors.lokasi_nikah,
            waktu_nikah: scheduleValidation.message
          }
        });
        return;
      }
    } catch (error) {
      alert('Gagal memvalidasi ketersediaan jadwal. Silakan coba lagi.');
      return;
    }

    // Submit ke API
    setLoading(true);
    try {
      const response = await apiClient.post('/simnikah/pendaftaran', formData);
      alert('Pendaftaran berhasil!');
      // Reset form atau redirect
    } catch (error) {
      // Handle error dari backend
      if (error.field) {
        setErrors({
          ...errors,
          [error.field]: error.error
        });
      } else {
        alert(error.message || 'Pendaftaran gagal');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields dengan error display */}
      {errors.calon_laki_laki?.nama_dan_bin && (
        <div className="error">{errors.calon_laki_laki.nama_dan_bin}</div>
      )}
      {/* ... */}
    </form>
  );
};
```

---

### 8. Checklist Validasi Wajib

#### ✅ Form Pendaftaran
- [ ] Validasi nama calon laki-laki (required, min 3, max 100, format "bin")
- [ ] Validasi nama calon perempuan (required, min 3, max 100, format "binti")
- [ ] Validasi umur (required, min 19 tahun, max 100)
- [ ] Validasi wali nikah (required, format "bin", hubungan wali valid)
- [ ] Validasi tanggal (required, format YYYY-MM-DD, tidak boleh masa lalu, bukan Minggu)
- [ ] Validasi waktu (required, format HH:MM, 08:00-16:00, dalam time slots)
- [ ] Validasi alamat (jika luar KUA: required, min 10, max 200)
- [ ] Validasi kelurahan (jika luar KUA: required, dalam daftar valid)
- [ ] Validasi ketersediaan slot sebelum submit

#### ✅ Authentication
- [ ] Validasi username (required, min 3, max 50, alphanumeric + underscore)
- [ ] Validasi email (required, format email valid, max 100)
- [ ] Validasi password (required, min 6 karakter)
- [ ] Validasi nama (required, min 3, max 100)
- [ ] Validasi role (required, dalam daftar valid)

#### ✅ File Upload
- [ ] Validasi file size (max 5MB)
- [ ] Validasi file type (JPG, PNG, WEBP)
- [ ] Validasi file extension

#### ✅ Security
- [ ] Sanitize semua input string
- [ ] Escape HTML untuk display
- [ ] Validasi role-based access
- [ ] Validasi status transition

#### ✅ Business Logic
- [ ] Validasi ketersediaan jadwal real-time
- [ ] Validasi status flow
- [ ] Validasi feedback (rating 1-5, judul, pesan)

---

### 9. Rekomendasi Library Validasi

#### React
```bash
# Yup (Schema validation)
npm install yup

# React Hook Form (Form management + validation)
npm install react-hook-form @hookform/resolvers yup

# Zod (TypeScript-first validation)
npm install zod @hookform/resolvers
```

#### Vue.js
```bash
# VeeValidate (Form validation)
npm install vee-validate yup

# atau
npm install vee-validate zod
```

---

### 10. Contoh dengan Yup Schema

```javascript
// utils/schemas/registration.js
import * as yup from 'yup';

export const registrationSchema = yup.object().shape({
  calon_laki_laki: yup.object().shape({
    nama_dan_bin: yup
      .string()
      .required('Nama dan bin wajib diisi')
      .min(3, 'Nama minimal 3 karakter')
      .max(100, 'Nama maksimal 100 karakter')
      .matches(/bin/, 'Format nama harus menggunakan "bin"'),
    pendidikan_akhir: yup
      .string()
      .max(50, 'Pendidikan maksimal 50 karakter'),
    umur: yup
      .number()
      .required('Umur wajib diisi')
      .min(19, 'Umur minimal 19 tahun')
      .max(100, 'Umur tidak valid')
  }),
  calon_perempuan: yup.object().shape({
    nama_dan_binti: yup
      .string()
      .required('Nama dan binti wajib diisi')
      .min(3, 'Nama minimal 3 karakter')
      .max(100, 'Nama maksimal 100 karakter')
      .matches(/binti/, 'Format nama harus menggunakan "binti"'),
    pendidikan_akhir: yup
      .string()
      .max(50, 'Pendidikan maksimal 50 karakter'),
    umur: yup
      .number()
      .required('Umur wajib diisi')
      .min(19, 'Umur minimal 19 tahun')
      .max(100, 'Umur tidak valid')
  }),
  lokasi_nikah: yup.object().shape({
    tempat_nikah: yup
      .string()
      .required('Tempat nikah wajib diisi')
      .oneOf(['Di KUA', 'Di Luar KUA'], 'Pilih tempat nikah yang valid'),
    tanggal_nikah: yup
      .string()
      .required('Tanggal nikah wajib diisi')
      .matches(/^\d{4}-\d{2}-\d{2}$/, 'Format tanggal harus YYYY-MM-DD')
      .test('not-past', 'Tanggal tidak boleh di masa lalu', function(value) {
        if (!value) return true;
        const date = new Date(value);
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        return date >= today;
      })
      .test('not-sunday', 'Tidak bisa memilih hari Minggu', function(value) {
        if (!value) return true;
        const date = new Date(value);
        return date.getDay() !== 0;
      }),
    waktu_nikah: yup
      .string()
      .required('Waktu nikah wajib diisi')
      .matches(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'Format waktu harus HH:MM')
      .oneOf([
        '08:00', '09:00', '10:00', '11:00', '12:00',
        '13:00', '14:00', '15:00', '16:00'
      ], 'Waktu harus sesuai dengan slot yang tersedia'),
    alamat_nikah: yup
      .string()
      .when('tempat_nikah', {
        is: 'Di Luar KUA',
        then: (schema) => schema
          .required('Alamat nikah wajib diisi untuk nikah di luar KUA')
          .min(10, 'Alamat minimal 10 karakter')
          .max(200, 'Alamat maksimal 200 karakter'),
        otherwise: (schema) => schema
      }),
    kelurahan: yup
      .string()
      .when('tempat_nikah', {
        is: 'Di Luar KUA',
        then: (schema) => schema
          .required('Kelurahan wajib diisi')
          .oneOf([
            'Alalak Utara', 'Alalak Tengah', 'Alalak Selatan',
            'Antasan Kecil Timur', 'Kuin Utara', 'Pangeran',
            'Sungai Miai', 'Sungai Andai', 'Surgi Mufti'
          ], 'Kelurahan harus dalam Kecamatan Banjarmasin Utara'),
        otherwise: (schema) => schema
      })
  }),
  wali_nikah: yup.object().shape({
    nama_dan_bin: yup
      .string()
      .required('Nama wali nikah wajib diisi')
      .min(3, 'Nama minimal 3 karakter')
      .max(100, 'Nama maksimal 100 karakter')
      .matches(/bin/, 'Format nama harus menggunakan "bin"'),
    hubungan_wali: yup
      .string()
      .required('Hubungan wali wajib diisi')
      .oneOf([
        'Ayah Kandung', 'Kakek', 'Saudara Laki-Laki Kandung',
        'Saudara Laki-Laki Seayah', 'Keponakan Laki-Laki',
        'Paman Kandung', 'Paman Seayah', 'Sepupu Laki-Laki',
        'Wali Hakim', 'Lainnya'
      ], 'Hubungan wali tidak valid')
  })
});
```

---

## 🧪 Testing

### Unit Testing

```javascript
// __tests__/api.test.js
import apiClient from '../utils/api';

describe('API Client', () => {
  it('should add token to request headers', () => {
    localStorage.setItem('token', 'test-token');
    
    apiClient.interceptors.request.use((config) => {
      expect(config.headers.Authorization).toBe('Bearer test-token');
    });
  });
});
```

### Integration Testing

```javascript
// __tests__/registration.test.js
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import RegistrationForm from '../components/RegistrationForm';

describe('Registration Form', () => {
  it('should submit form successfully', async () => {
    render(<RegistrationForm />);
    
    // Fill form
    fireEvent.change(screen.getByLabelText('Nama'), {
      target: { value: 'Ahmad' }
    });
    
    // Submit
    fireEvent.click(screen.getByText('Kirim Pendaftaran'));
    
    // Wait for success
    await waitFor(() => {
      expect(screen.getByText('Pendaftaran berhasil!')).toBeInTheDocument();
    });
  });
});
```

---

## ❓ FAQ

### Q: Bagaimana cara handle token expired?
**A:** Gunakan response interceptor untuk detect 401 error, lalu redirect ke login dan hapus token.

### Q: Apakah bisa menggunakan fetch instead of axios?
**A:** Bisa, tapi axios lebih mudah untuk interceptor dan error handling.

### Q: Bagaimana cara handle rate limiting?
**A:** Tampilkan pesan error 429 ke user, dan implementasi retry dengan exponential backoff.

### Q: Apakah perlu refresh token?
**A:** Untuk saat ini, token berlaku 24 jam. User perlu login ulang jika expired.

### Q: Bagaimana cara test API di development?
**A:** Gunakan Postman, Insomnia, atau curl untuk test endpoint sebelum implementasi di frontend.

### Q: Apakah semua endpoint perlu authentication?
**A:** Tidak. Endpoint kalender ketersediaan (`/simnikah/kalender-ketersediaan`) dan beberapa endpoint lainnya adalah public.

### Q: Bagaimana format waktu dan tanggal?
**A:** 
- Tanggal: `YYYY-MM-DD` (contoh: `2024-12-15`)
- Waktu: `HH:MM` (24 jam, contoh: `09:00`, `14:30`)
- Timestamp: ISO 8601 (contoh: `2024-12-15T10:00:00Z`)

---

## 📚 Resources Tambahan

- **API Documentation Lengkap:** `docs/api/API_DOCUMENTATION_LENGKAP.md`
- **Alur Pendaftaran:** `docs/ALUR_PENDAFTARAN_DAN_JAM.md`
- **Implementation Guide:** `docs/FRONTEND_IMPLEMENTATION_GUIDE.md`

---

## 📞 Support

Jika ada pertanyaan atau masalah terkait integrasi frontend:

1. Cek dokumentasi API lengkap
2. Cek error message dari API response
3. Pastikan token valid dan tidak expired
4. Pastikan role user sesuai dengan endpoint yang diakses

---

**Last Updated:** Desember 2024  
**API Version:** 1.3.0  
**Frontend Documentation Version:** 2.0.0

