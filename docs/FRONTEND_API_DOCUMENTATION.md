# 📱 SimNikah API - Frontend Documentation

**Base URL:** `https://your-api-domain.com` atau `http://localhost:8080`  
**Version:** 1.3.0  
**Format:** JSON  
**Authentication:** JWT Bearer Token

---

## 🚀 Quick Start

### 1. Setup API Client

```javascript
// api.js
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const token = localStorage.getItem('token');

    const config = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Request failed');
      }

      return data;
    } catch (error) {
      console.error('API Error:', error);
      throw error;
    }
  }

  get(endpoint) {
    return this.request(endpoint, { method: 'GET' });
  }

  post(endpoint, body) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(body),
    });
  }

  put(endpoint, body) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(body),
    });
  }

  delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }
}

export const api = new ApiClient();
```

### 2. Authentication Helper

```javascript
// auth.js
import { api } from './api';

export const auth = {
  async login(username, password) {
    const response = await api.post('/login', { username, password });
    if (response.token) {
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
    }
    return response;
  },

  async register(userData) {
    const response = await api.post('/register', userData);
    return response;
  },

  async getProfile() {
    return await api.get('/profile');
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  isAuthenticated() {
    return !!localStorage.getItem('token');
  },

  getCurrentUser() {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  },
};
```

---

## 🔐 Authentication Endpoints

### POST `/login`

Login user dan dapatkan JWT token.

**Request:**
```javascript
const response = await api.post('/login', {
  username: 'ahmad123',
  password: 'password123'
});

// Response:
{
  success: true,
  message: "Login berhasil",
  token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  user: {
    user_id: "USR1704067200",
    email: "ahmad@example.com",
    role: "user_biasa",
    nama: "Ahmad Wijaya"
  }
}
```

### POST `/register`

Daftar user baru.

**Request:**
```javascript
const response = await api.post('/register', {
  username: 'ahmad123',
  email: 'ahmad@example.com',
  password: 'password123',
  nama: 'Ahmad Wijaya',
  role: 'user_biasa' // user_biasa | staff | penghulu | kepala_kua
});
```

### GET `/profile`

Ambil profile user yang sedang login.

**Request:**
```javascript
const response = await api.get('/profile');
```

---

## 💒 Marriage Registration (Catin)

### POST `/simnikah/pendaftaran`

Buat pendaftaran nikah baru.

**Request:**
```javascript
const response = await api.post('/simnikah/pendaftaran', {
  calon_laki_laki: {
    nama_dan_bin: "Ahmad bin Abdullah",
    pendidikan_akhir: "S1",
    umur: 25
  },
  calon_perempuan: {
    nama_dan_binti: "Siti binti Muhammad",
    pendidikan_akhir: "S1",
    umur: 23
  },
  lokasi_nikah: {
    tempat_nikah: "Di KUA", // atau "Di Luar KUA"
    tanggal_nikah: "2024-12-25",
    waktu_nikah: "10:00",
    // Hanya jika tempat_nikah = "Di Luar KUA"
    alamat_nikah: "Jl. Contoh No. 123",
    detail_alamat: "RT 01 RW 02",
    kelurahan: "Sungai Miai"
  },
  wali_nikah: {
    nama_dan_bin: "Abdullah bin Muhammad",
    hubungan_wali: "Ayah Kandung" // Ayah Kandung | Kakek | dll
  }
});

// Response:
{
  success: true,
  message: "Pendaftaran berhasil dibuat",
  data: {
    id: 1,
    nomor_pendaftaran: "REG20240101001",
    status_pendaftaran: "Draft",
    tanggal_nikah: "2024-12-25T10:00:00Z",
    // ... data lengkap
  }
}
```

### GET `/simnikah/pendaftaran/status`

Cek status pendaftaran user.

**Request:**
```javascript
const response = await api.get('/simnikah/pendaftaran/status');

// Response jika ada pendaftaran:
{
  success: true,
  has_registration: true,
  data: {
    id: 1,
    nomor_pendaftaran: "REG20240101001",
    status_pendaftaran: "Disetujui",
    // ... data lengkap
  }
}

// Response jika belum ada:
{
  success: true,
  has_registration: false,
  message: "Belum ada pendaftaran"
}
```

---

## 📅 Calendar & Availability

### GET `/simnikah/kalender-ketersediaan`

Ambil kalender ketersediaan (tidak perlu auth).

**Request:**
```javascript
const response = await api.get('/simnikah/kalender-ketersediaan?bulan=12&tahun=2024');

// Response:
{
  success: true,
  data: {
    "2024-12-01": { available: true, booked: 2, max: 9 },
    "2024-12-02": { available: true, booked: 5, max: 9 },
    // ...
  }
}
```

### GET `/simnikah/ketersediaan-jam`

Ambil slot waktu yang tersedia untuk tanggal tertentu.

**Request:**
```javascript
const response = await api.get('/simnikah/ketersediaan-jam?tanggal=2024-12-25&tempat=Di KUA');

// Response:
{
  success: true,
  data: {
    available_slots: ["08:00", "09:00", "11:00", "14:00"],
    unavailable_slots: ["10:00", "12:00", "13:00", "15:00", "16:00"],
    tanggal: "2024-12-25",
    tempat: "Di KUA"
  }
}
```

---

## 👨‍💼 Staff Endpoints

### GET `/simnikah/staff/pengumuman-nikah/list`

Ambil data pendaftaran untuk surat pengumuman nikah.

**Request:**
```javascript
// POST request dengan body (optional untuk custom kop surat)
const response = await api.post('/simnikah/staff/pengumuman-nikah/list', {
  tanggal_awal: "2024-12-01",
  tanggal_akhir: "2024-12-07",
  // Optional: custom kop surat
  kop_surat: {
    nama_kua: "KUA Kecamatan Banjarmasin Utara",
    alamat: "Jl. Contoh No. 123",
    // ... data kop surat lainnya
  }
});

// Response:
{
  success: true,
  data: {
    registrations: [
      {
        id: 1,
        nomor_pendaftaran: "REG20240101001",
        calon_suami: { nama_lengkap: "Ahmad bin Abdullah" },
        calon_istri: { nama_lengkap: "Siti binti Muhammad" },
        tanggal_nikah: "2024-12-25T10:00:00Z",
        waktu_nikah: "10:00",
        // ... data lengkap
      }
    ],
    kop_surat: { /* data kop surat */ },
    periode: {
      tanggal_awal: "2024-12-01",
      tanggal_akhir: "2024-12-07"
    }
  }
}
```

### POST `/simnikah/staff/approve/:id`

Approve pendaftaran nikah.

**Request:**
```javascript
const response = await api.post(`/simnikah/staff/approve/${registrationId}`);
```

---

## 🕌 Penghulu Endpoints

### GET `/simnikah/penghulu/assigned-registrations`

Ambil daftar pendaftaran yang ditugaskan ke penghulu.

**Request:**
```javascript
const response = await api.get('/simnikah/penghulu/assigned-registrations');

// Response:
{
  success: true,
  data: [
    {
      id: 1,
      nomor_pendaftaran: "REG20240101001",
      calon_suami: { nama_lengkap: "Ahmad bin Abdullah" },
      calon_istri: { nama_lengkap: "Siti binti Muhammad" },
      tanggal_nikah: "2024-12-25T10:00:00Z",
      waktu_nikah: "10:00",
      tempat_nikah: "Di KUA",
      status_pendaftaran: "Penghulu Ditugaskan"
    }
  ]
}
```

### POST `/simnikah/penghulu/complete-marriage/:id`

Tandai pernikahan selesai.

**Request:**
```javascript
const response = await api.post(`/simnikah/penghulu/complete-marriage/${registrationId}`);
```

---

## 👔 Kepala KUA Endpoints

### POST `/simnikah/pendaftaran/:id/assign-penghulu`

Assign penghulu ke pendaftaran.

**Request:**
```javascript
const response = await api.post(`/simnikah/pendaftaran/${registrationId}/assign-penghulu`, {
  penghulu_id: 1
});
```

### GET `/simnikah/kepala-kua/pengumuman-nikah/list`

Ambil data untuk surat pengumuman (sama seperti staff).

**Request:**
```javascript
const response = await api.post('/simnikah/kepala-kua/pengumuman-nikah/list', {
  tanggal_awal: "2024-12-01",
  tanggal_akhir: "2024-12-07",
  kop_surat: { /* optional */ }
});
```

---

## 📍 Location Endpoints

### POST `/simnikah/location/geocode`

Convert alamat ke koordinat (latitude, longitude).

**Request:**
```javascript
const response = await api.post('/simnikah/location/geocode', {
  alamat: "Jl. Contoh No. 123, Banjarmasin"
});

// Response:
{
  success: true,
  data: {
    latitude: -3.3145,
    longitude: 114.5921,
    formatted_address: "Jl. Contoh No. 123, Banjarmasin"
  }
}
```

### POST `/simnikah/location/reverse-geocode`

Convert koordinat ke alamat.

**Request:**
```javascript
const response = await api.post('/simnikah/location/reverse-geocode', {
  latitude: -3.3145,
  longitude: 114.5921
});
```

---

## 🔔 Notification Endpoints

### GET `/simnikah/notifikasi/user/:user_id`

Ambil notifikasi user.

**Request:**
```javascript
const user = auth.getCurrentUser();
const response = await api.get(`/simnikah/notifikasi/user/${user.user_id}`);

// Response:
{
  success: true,
  data: [
    {
      id: 1,
      judul: "Pendaftaran Disetujui",
      pesan: "Pendaftaran Anda telah disetujui",
      tipe: "Success",
      status_baca: "Belum Dibaca",
      created_at: "2024-12-01T10:00:00Z"
    }
  ]
}
```

### PUT `/simnikah/notifikasi/:id/status`

Update status notifikasi (baca/tidak baca).

**Request:**
```javascript
const response = await api.put(`/simnikah/notifikasi/${notificationId}/status`, {
  status_baca: "Sudah Dibaca"
});
```

---

## 📊 Dashboard Endpoints

### GET `/simnikah/dashboard/kepala-kua`

Dashboard untuk kepala KUA.

**Request:**
```javascript
const response = await api.get('/simnikah/dashboard/kepala-kua');

// Response:
{
  success: true,
  data: {
    total_pendaftaran: 150,
    pending_approval: 10,
    pernikahan_hari_ini: 5,
    // ... statistik lainnya
  }
}
```

### GET `/simnikah/dashboard/staff`

Dashboard untuk staff.

**Request:**
```javascript
const response = await api.get('/simnikah/dashboard/staff');
```

---

## ❌ Error Handling

Semua error mengikuti format standar:

```javascript
{
  success: false,
  message: "Error message",
  error: "Detailed error description",
  type: "error_type" // validation | authentication | authorization | not_found | database
}
```

**Contoh Error Handling:**

```javascript
try {
  const response = await api.post('/simnikah/pendaftaran', data);
  // Success
} catch (error) {
  if (error.message.includes('401') || error.message.includes('Unauthorized')) {
    // Token expired atau tidak valid
    auth.logout();
    // Redirect ke login
  } else if (error.message.includes('403')) {
    // Tidak punya akses
    alert('Anda tidak memiliki akses untuk fitur ini');
  } else if (error.message.includes('400')) {
    // Validation error
    alert('Data yang diinput tidak valid');
  } else {
    // Server error
    alert('Terjadi kesalahan pada server');
  }
}
```

---

## 🔑 Status Pendaftaran

Status flow pendaftaran nikah:

1. **Draft** - Pendaftaran baru dibuat
2. **Disetujui** - Staff sudah approve
3. **Menunggu Penugasan** - Menunggu kepala KUA assign penghulu
4. **Penghulu Ditugaskan** - Penghulu sudah di-assign
5. **Selesai** - Pernikahan sudah dilaksanakan
6. **Ditolak** - Pendaftaran ditolak

---

## 📝 Tips & Best Practices

### 1. Token Management
- Simpan token di `localStorage` atau `sessionStorage`
- Refresh token sebelum expired (jika ada fitur refresh)
- Hapus token saat logout

### 2. Error Handling
- Selalu handle error dengan try-catch
- Tampilkan pesan error yang user-friendly
- Log error untuk debugging

### 3. Loading States
- Tampilkan loading indicator saat request
- Gunakan skeleton loader untuk UX yang lebih baik

### 4. Caching
- Cache data yang jarang berubah (misal: kalender ketersediaan)
- Invalidate cache saat data di-update

### 5. Rate Limiting
- API memiliki rate limiting (100 req/min per IP)
- Handle 429 Too Many Requests dengan retry logic

---

## 🔗 Quick Reference

| Endpoint | Method | Auth | Role |
|----------|--------|------|------|
| `/login` | POST | ❌ | - |
| `/register` | POST | ❌ | - |
| `/profile` | GET | ✅ | All |
| `/simnikah/pendaftaran` | POST | ✅ | user_biasa |
| `/simnikah/pendaftaran/status` | GET | ✅ | user_biasa |
| `/simnikah/kalender-ketersediaan` | GET | ❌ | - |
| `/simnikah/ketersediaan-jam` | GET | ❌ | - |
| `/simnikah/staff/pengumuman-nikah/list` | POST | ✅ | staff, kepala_kua |
| `/simnikah/penghulu/assigned-registrations` | GET | ✅ | penghulu |
| `/simnikah/kepala-kua/pengumuman-nikah/list` | POST | ✅ | kepala_kua |
| `/simnikah/notifikasi/user/:user_id` | GET | ✅ | All |

---

## 📞 Support

Jika ada pertanyaan atau issue, hubungi tim backend development.

**Last Updated:** December 2024

