  # 📚 SimNikah API Documentation

  **Versi:** 1.3.0  
  **Update Terakhir:** November 2024  
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
  10. [Dashboard & Analytics Endpoints](#-dashboard--analytics-endpoints)
  11. [Error Handling](#error-handling)
  12. [Status Flow](#status-flow)

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
  | **Dashboard** | `/simnikah/dashboard/*` | GET | Yes | `staff`, `kepala_kua` |

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
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
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
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    }
  }
  ```

  **Validasi:**
  - Umur minimal: 19 tahun (untuk calon laki-laki dan perempuan)
  - Format tanggal: `YYYY-MM-DD`
  - Format waktu: `HH:MM` (24-jam, contoh: `09:00`, `14:30`)
  - Tanggal nikah tidak boleh di masa lalu
  - Kelurahan harus dalam lingkup **Kecamatan Banjarmasin Utara**
  - **Wali Nikah wajib diisi** (untuk calon pengantin perempuan)
    - `nama_dan_bin`: Nama lengkap wali dengan bin (contoh: "Abdullah bin Muhammad")
    - `hubungan_wali`: Hubungan nasab wali (lihat daftar di bawah)

  **Kelurahan Valid:**
  - Alalak Utara
  - Alalak Tengah
  - Alalak Selatan
  - Antasan Kecil Timur
  - Kuin Utara
  - Pangeran
  - Sungai Miai
  - Sungai Andai
  - Surgi Mufti

  **Hubungan Wali Valid (Urutan Wali Nasab):**
  1. `"Ayah Kandung"` - Wali yang paling berhak (jika ayah masih hidup)
  2. `"Kakek"` - Ayah dari ayah (jika ayah meninggal)
  3. `"Saudara Laki-Laki Kandung"` - Saudara sekandung
  4. `"Saudara Laki-Laki Seayah"` - Saudara seayah
  5. `"Keponakan Laki-Laki"` - Anak laki-laki dari saudara
  6. `"Paman Kandung"` - Saudara kandung ayah
  7. `"Paman Seayah"` - Saudara seayah dari ayah
  8. `"Sepupu Laki-Laki"` - Anak laki-laki dari paman
  9. `"Wali Hakim"` - Jika tidak ada wali nasab yang memenuhi syarat
  10. `"Lainnya"` - Hubungan lainnya

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
      "wali_nikah": {
        "nama_dan_bin": "Abdullah bin Muhammad",
        "hubungan_wali": "Ayah Kandung"
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

  ### 5.5. Get Detail Pendaftaran by ID

  **Endpoint:** `GET /simnikah/pendaftaran/:id`

  **Description:** Mendapatkan detail lengkap pendaftaran nikah berdasarkan ID. Endpoint ini berguna untuk halaman detail pendaftaran di frontend.

  **Auth Required:** ✅ Yes

  **Role Access:**
  - `user_biasa`: Hanya bisa melihat pendaftaran miliknya sendiri
  - `staff`, `penghulu`, `kepala_kua`: Bisa melihat semua pendaftaran

  **URL Parameters:**
  - `id` (required): ID pendaftaran (integer)

  **Example:**
  ```
  GET /simnikah/pendaftaran/1
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

  **Response Error (404):**
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

  **Catatan:**
  - `jumlah_draft`: Jumlah pendaftaran dengan status Draft (kuning - belum pasti)
  - `jumlah_disetujui`: Jumlah pendaftaran dengan status Disetujui (hijau - sudah pasti)
  - `sisa_kuota`: Dihitung berdasarkan yang sudah pasti (Disetujui) saja
  - `time_slots`: Detail ketersediaan jam per tanggal untuk KUA dan luar KUA

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
          "terbooking": false,
          "jumlah_nikah": 0,
          "jumlah_draft": 0,
          "jumlah_disetujui": 0
        },
        {
          "waktu": "09:00",
          "tersedia": false,
          "terbooking": true,
          "jumlah_nikah": 1,
          "jumlah_draft": 0,
          "jumlah_disetujui": 1
        }
      ],
      "registrations_today": {
        "total": 5,
        "detail": []
      }
    }
  }
  ```

  **Catatan:**
  - `jumlah_draft`: Jumlah pendaftaran dengan status Draft (kuning - belum pasti)
  - `jumlah_disetujui`: Jumlah pendaftaran dengan status Disetujui (hijau - sudah pasti)
  - `terbooking`: Dianggap terbooking jika ada yang sudah pasti (Disetujui)
  - `tersedia`: Tersedia jika tidak terbooking dan tanggal tidak di masa lalu

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

  ### 14. Create Registration for User (Staff)

  **Endpoint:** `POST /simnikah/staff/pendaftaran`

  **Description:** Staff dapat membuat pendaftaran nikah atas nama calon pengantin. Fitur ini berguna untuk membantu calon pengantin yang tidak paham teknologi. Sistem akan secara otomatis membuat akun user untuk calon pengantin dan memberikan username serta password default. **Pendaftaran yang dibuat oleh staff otomatis berstatus "Disetujui"** karena staff sudah melakukan verifikasi saat input data.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`, `kepala_kua`

  **Request Body:**
  Sama seperti endpoint pendaftaran biasa (`POST /simnikah/pendaftaran`):
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
      "tanggal_nikah": "2024-12-25",
      "waktu_nikah": "09:00"
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
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
      "tanggal_nikah": "2024-12-25",
      "waktu_nikah": "09:00",
      "alamat_nikah": "Jl. Ahmad Yani No. 123",
      "detail_alamat": "Rumah Pengantin Perempuan",
      "kelurahan": "Pangeran"
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    }
  }
  ```

  **Validasi:**
  - Semua validasi sama seperti endpoint pendaftaran biasa
  - Umur minimal: 19 tahun
  - Format tanggal: `YYYY-MM-DD`
  - Format waktu: `HH:MM` (24-jam)
  - Tanggal nikah tidak boleh di masa lalu
  - Validasi ketersediaan jadwal (kapasitas per jam)
  - **Wali Nikah wajib diisi** (untuk calon pengantin perempuan)

  **Response Success (201):**
  ```json
  {
    "success": true,
    "message": "Pendaftaran nikah berhasil dibuat dan disetujui oleh staff",
    "data": {
      "id": 123,
      "nomor_pendaftaran": "NIKAH-20241225-1234",
      "status_pendaftaran": "Disetujui",
      "tanggal_nikah": "2024-12-25T00:00:00Z",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "PH5Q+F8C, Jl. Wira Karya, Pangeran, Kec. Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan 70123",
      "dibuat_oleh_staff": {
        "nama": "Budi Santoso",
        "nip": "198001012003121001"
      },
      "akun_user": {
        "user_id": "USR1704067200",
        "username": "ahmadwijay1234",
        "email": "ahmadwijay1234@simnikah.local",
        "password_default": "Nikah12345",
        "catatan": "Akun ini dibuat otomatis. User dapat login dan mengubah password."
      },
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
      "catatan": "Pendaftaran dibuat oleh staff. User dapat login menggunakan username dan password default yang diberikan."
    }
  }
  ```

  **Response Error (400):**
  ```json
  {
    "success": false,
    "message": "Format tanggal tidak benar",
    "error": "Format tanggal harus: Tahun-Bulan-Tanggal (contoh: 2024-12-25)",
    "field": "tanggal_nikah",
    "type": "format"
  }
  ```

  **Response Error (403):**
  ```json
  {
    "success": false,
    "message": "Akses ditolak",
    "error": "Hanya staff yang dapat membuat pendaftaran untuk user",
    "type": "authorization"
  }
  ```

  **Catatan Penting:**
  1. **Akun User Otomatis**: Sistem akan membuat akun user otomatis untuk calon pengantin dengan:
     - Username: Generated dari nama calon suami + timestamp
     - Email: `username@simnikah.local`
     - Password: Generated default (format: `Nikah` + angka)
  2. **Informasi Login**: Staff harus memberikan username dan password default kepada calon pengantin
  3. **Ubah Password**: Calon pengantin dapat login dan mengubah password setelah mendapatkan akses
  4. **Catatan Otomatis**: Pendaftaran akan memiliki catatan otomatis yang menyebutkan staff yang membuat pendaftaran
  5. **Status Otomatis Disetujui**: Pendaftaran yang dibuat oleh staff otomatis berstatus "Disetujui" karena staff sudah melakukan verifikasi saat input data. Tidak perlu approval lagi.

  **Use Case:**
  - Calon pengantin datang ke kantor KUA untuk mendaftar
  - Staff membantu mengisi form pendaftaran melalui sistem
  - Staff memberikan username dan password kepada calon pengantin
  - Calon pengantin dapat login nanti untuk melihat status pendaftaran

  ---

  ### 15. Get Approved Registrations Per Week (Staff)

  **Endpoint:** `GET /simnikah/staff/pengumuman-nikah/list`

  **Description:** Mendapatkan daftar pendaftaran nikah untuk periode yang dipilih frontend beserta data kop surat. **Menampilkan semua status pendaftaran kecuali "Ditolak"** (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai). Data ini digunakan oleh frontend untuk generate surat pengumuman nikah dalam format HTML.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`, `kepala_kua`

  **Query Parameters:**
  - `tanggal_awal` (optional): Tanggal awal periode yang dipilih frontend (format: YYYY-MM-DD). Frontend dapat memilih tanggal bebas sesuai kebutuhan (tidak harus mingguan). Default: awal minggu ini (Senin) jika tidak diisi
  - `tanggal_akhir` (optional): Tanggal akhir periode yang dipilih frontend (format: YYYY-MM-DD). Frontend dapat memilih tanggal bebas sesuai kebutuhan (tidak harus mingguan). Default: akhir minggu ini (Minggu) jika tidak diisi
  
  **Note:** Frontend dapat memilih periode tanggal secara bebas sesuai kebutuhan. Contoh: bisa memilih 1 bulan, 2 minggu, atau periode custom lainnya. Tidak terbatas pada periode mingguan.

  **Request Body (Optional - untuk custom kop surat):**
  ```json
  {
    "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
    "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
    "kota": "Kota Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "kode_pos": "70123",
    "telepon": "0511-1234567",
    "email": "kua.banjarmasinutara@kemenag.go.id",
    "website": "https://kua.banjarmasinutara.go.id",
    "logo_url": "https://example.com/logo-kua.png"
  }
  ```

  **Note:** 
  - Jika request body tidak dikirim, akan menggunakan nilai default untuk kop surat
  - Frontend bertanggung jawab untuk generate HTML surat pengumuman nikah menggunakan data yang dikembalikan

  **Examples:**
  
  **Mingguan (default):**
  ```
  GET /simnikah/staff/pengumuman-nikah/list?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
  ```
  
  **Bulanan:**
  ```
  GET /simnikah/staff/pengumuman-nikah/list?tanggal_awal=2024-12-01&tanggal_akhir=2024-12-31
  ```
  
  **Custom periode:**
  ```
  GET /simnikah/staff/pengumuman-nikah/list?tanggal_awal=2024-12-10&tanggal_akhir=2024-12-25
  ```
  
  **Tanpa parameter (menggunakan default - minggu ini):**
  ```
  GET /simnikah/staff/pengumuman-nikah/list
  ```

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Data pendaftaran berhasil diambil",
    "data": {
      "tanggal_awal": "2024-12-16",
      "tanggal_akhir": "2024-12-22",
      "periode": "16 Desember 2024 s/d 22 Desember 2024",
      "total": 5,
      "kop_surat": {
        "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
        "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
        "kota": "Kota Banjarmasin",
        "provinsi": "Kalimantan Selatan",
        "kode_pos": "70123",
        "telepon": "-",
        "email": "-",
        "website": "",
        "logo_url": ""
      },
      "registrations": [
        {
          "id": 1,
          "nomor_pendaftaran": "NIKAH-20241215-1234",
          "status_pendaftaran": "Disetujui",
          "tanggal_nikah": "2024-12-18T10:00:00Z",
          "waktu_nikah": "10:00",
          "tempat_nikah": "Di KUA",
          "alamat_akad": "Kantor KUA Kecamatan Banjarmasin Utara",
          "calon_suami": {
            "nama_lengkap": "Ahmad bin Abdullah"
          },
          "calon_istri": {
            "nama_lengkap": "Siti binti Muhammad"
          },
          "wali_nikah": {
            "nama_dan_bin": "Muhammad bin Ali",
            "hubungan_wali": "Ayah Kandung"
          }
        }
      ]
    }
  }
  ```

  **Catatan Penting:**
  - **Status yang ditampilkan:** Semua status kecuali "Ditolak" (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai)
  - **Status "Ditolak" tidak ditampilkan** karena pendaftaran yang ditolak tidak perlu diumumkan
  - Response includes `status_pendaftaran` untuk setiap pendaftaran agar frontend dapat menampilkan badge/indikator status

  **Use Case:**
  - Frontend memanggil endpoint ini untuk mendapatkan data pendaftaran dan kop surat
  - Frontend menggunakan data tersebut untuk generate HTML surat pengumuman nikah
  - Frontend dapat menampilkan badge status untuk setiap pendaftaran (Draft = kuning, Disetujui = hijau, dll)
  - Surat dapat dicetak atau dikonversi ke PDF di frontend
  - Surat dicetak dan dipasang di papan pengumuman KUA

  ---

  ### 17. Approve Registration

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

  ### 18. Update Registration Status

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

  ### 15.5. Generate Pengumuman Nikah HTML (Staff)

  **Endpoint:** `GET /simnikah/staff/pengumuman-nikah/generate`  
  **Endpoint:** `POST /simnikah/staff/pengumuman-nikah/generate`

  **Description:** Generate surat pengumuman nikah dalam format HTML yang siap dicetak atau dikonversi ke PDF. Endpoint ini mengembalikan HTML document lengkap dengan kop surat, tabel data pendaftaran, dan format surat resmi.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`, `kepala_kua`

  **Query Parameters:**
  - `tanggal_awal` (optional): Tanggal awal periode (format: YYYY-MM-DD). Default: awal minggu ini (Senin)
  - `tanggal_akhir` (optional): Tanggal akhir periode (format: YYYY-MM-DD). Default: akhir minggu ini (Minggu)

  **Request Body (Optional - untuk custom kop surat):**
  ```json
  {
    "tanggal_awal": "2024-12-16",
    "tanggal_akhir": "2024-12-22",
    "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
    "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
    "kota": "Kota Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "kode_pos": "70123",
    "telepon": "0511-1234567",
    "email": "kua.banjarmasinutara@kemenag.go.id",
    "website": "https://kua.banjarmasinutara.go.id",
    "logo_url": "https://example.com/logo-kua.png"
  }
  ```

  **Request Body Fields:**

  | Field | Type | Required | Default | Deskripsi |
  |-------|------|----------|---------|-----------|
  | `tanggal_awal` | string (YYYY-MM-DD) | ❌ No | Awal minggu ini | Tanggal awal periode (fallback jika query param tidak ada) |
  | `tanggal_akhir` | string (YYYY-MM-DD) | ❌ No | Akhir minggu ini | Tanggal akhir periode (fallback jika query param tidak ada) |
  | `nama_kua` | string | ❌ No | "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA" | Nama KUA |
  | `alamat_kua` | string | ❌ No | "PH5Q+F8C, Jl. Wira Karya, Pangeran" | Alamat KUA |
  | `kota` | string | ❌ No | "Kota Banjarmasin" | Kota |
  | `provinsi` | string | ❌ No | "Kalimantan Selatan" | Provinsi |
  | `kode_pos` | string | ❌ No | "70123" | Kode pos |
  | `telepon` | string | ❌ No | "-" | Nomor telepon |
  | `email` | string | ❌ No | "-" | Email |
  | `website` | string | ❌ No | "" | Website (opsional) |
  | `logo_url` | string | ❌ No | "" | URL logo KUA (opsional, harus accessible) |

  **Note:** 
  - Query parameters memiliki prioritas lebih tinggi daripada request body
  - Jika query parameters tidak ada, akan membaca dari request body
  - Jika keduanya tidak ada, akan menggunakan default (minggu ini)
  - Logo URL harus accessible (public URL) dan disarankan menggunakan HTTPS
  - Semua field kop surat bersifat opsional, akan menggunakan default jika tidak diisi

  **Examples:**

  **1. Generate dengan default kop surat (GET):**
  ```
  GET /simnikah/staff/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
  Authorization: Bearer YOUR_TOKEN
  ```

  **2. Generate dengan custom kop surat (POST):**
  ```
  POST /simnikah/staff/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
  Authorization: Bearer YOUR_TOKEN
  Content-Type: application/json

  {
    "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
    "alamat_kua": "Jl. Contoh No. 123",
    "kota": "Kota Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "kode_pos": "70123",
    "telepon": "0511-1234567",
    "email": "kua@example.com",
    "website": "https://kua.example.com",
    "logo_url": "https://example.com/logo.png"
  }
  ```

  **Response Success (200):**

  **Content-Type:** `text/html; charset=utf-8`

  Response berupa **HTML document lengkap** yang siap dicetak atau dikonversi ke PDF.

  **Response Headers:**
  ```
  Content-Type: text/html; charset=utf-8
  Content-Length: <ukuran file>
  ```

  **Response Body:** HTML document lengkap dengan struktur sesuai format Excel standar KUA:

  1. **Kop Surat KUA:**
     - Logo di kiri (jika `logo_url` disediakan)
     - KEMENTERIAN AGAMA REPUBLIK INDONESIA
     - KANTOR KEMENTERIAN AGAMA KOTA [KOTA]
     - Nama KUA
     - Alamat lengkap
     - Kontak (telepon, email)

  2. **Judul:** "JADUAL NIKAH [BULAN] [TAHUN]" (contoh: "JADUAL NIKAH JANUARI 2026")

  3. **Tabel Data Pendaftaran** dengan **15 kolom** (kolom tetap, baris dinamis):
     
     **Struktur Kolom (Tetap 15 kolom):**
     - **NO URUT** - Nomor urut (1, 2, 3, ...)
     - **DATA CALON PENGANTIN:**
       - **PRIA / BIN** - Nama calon suami
       - **USIA** - Usia calon suami (dihitung otomatis dari tanggal lahir)
       - **PENDK** - Pendidikan terakhir calon suami
       - **WANITA / BINTI** - Nama calon istri
       - **USIA** - Usia calon istri (dihitung otomatis dari tanggal lahir)
       - **PENDK** - Pendidikan terakhir calon istri
     - **PELAKSANAAN NIKAH:**
       - **HARI** - Nama hari (SENIN, SELASA, RABU, KAMIS, JUM'AT, SABTU, AHAD)
       - **TGL** - Tanggal (hanya angka: 1, 2, 3, dll)
       - **JAM** - Waktu (format: 08.00, 09.00, dll - dari HH:MM menjadi HH.MM)
       - **TEMPAT** - Tempat nikah
       - **WALINIKAH** - Nama wali nikah
       - **PENGHULU** - Nama penghulu (jika sudah ditugaskan, "-" jika belum)
       - **KELURAHAN** - Kelurahan
       - **KET** - Keterangan
     
     **Catatan Penting:**
     - ✅ **Kolom tetap 15 kolom** (tidak berubah)
     - ✅ **Baris dinamis** - jumlah baris tergantung jumlah data pendaftaran dalam periode yang dipilih
     - ✅ **Setiap baris = 1 pendaftaran nikah**
     - ✅ Data diurutkan berdasarkan tanggal nikah dan waktu (ASC)
     - ✅ Semua pendaftaran dalam periode akan ditampilkan (kecuali status "Ditolak")

  **Catatan Penting:**
  - HTML sudah include CSS untuk print optimization (`@media print`)
  - Format **A4 Landscape** untuk menampung tabel lebar
  - Font: Times New Roman (serif)
  - Font size: 8-9pt untuk tabel, 11-12pt untuk kop surat
  - Tabel dengan border untuk kejelasan
  - Header tabel dengan rowspan/colspan untuk grouping kolom
  - Siap untuk dicetak langsung atau dikonversi ke PDF
  - **Status yang ditampilkan:** Semua status kecuali "Ditolak" (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai)
  - **Usia dihitung otomatis** dari tanggal lahir calon pengantin
  - **Format waktu:** HH:MM diubah menjadi HH.MM (contoh: 08:00 → 08.00)
  - **Nama hari:** Otomatis dalam bahasa Indonesia (SENIN, SELASA, dll)

  **Response Error:**

  **400 Bad Request:**
  ```json
  {
    "success": false,
    "message": "Format tanggal tidak valid",
    "error": "Format tanggal_awal harus YYYY-MM-DD"
  }
  ```

  **401 Unauthorized:**
  ```json
  {
    "success": false,
    "message": "Unauthorized",
    "error": "Token tidak valid atau tidak ada"
  }
  ```

  **403 Forbidden:**
  ```json
  {
    "success": false,
    "message": "Forbidden",
    "error": "Role tidak memiliki akses"
  }
  ```

  **500 Internal Server Error:**
  ```json
  {
    "success": false,
    "message": "Database error",
    "error": "Gagal mengambil data pendaftaran"
  }
  ```

  **Contoh Implementasi Frontend (JavaScript/React):**

  ```javascript
  // Menggunakan Axios
  import axios from 'axios';

  const generatePengumumanHTML = async (tanggalAwal, tanggalAkhir, kopSurat = null) => {
    const token = localStorage.getItem('token');
    
    const config = {
      method: kopSurat ? 'post' : 'get',
      url: `${API_BASE_URL}/staff/pengumuman-nikah/generate`,
      params: {
        tanggal_awal: tanggalAwal,
        tanggal_akhir: tanggalAkhir
      },
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      responseType: 'text' // PENTING: untuk mendapatkan HTML sebagai string
    };

    if (kopSurat) {
      config.data = {
        tanggal_awal: tanggalAwal,
        tanggal_akhir: tanggalAkhir,
        ...kopSurat
      };
    }

    try {
      const response = await axios(config);
      return response.data; // HTML string
    } catch (error) {
      console.error('Error generating pengumuman:', error);
      throw error;
    }
  };

  // Menggunakan response HTML
  const handlePrint = async () => {
    try {
      const html = await generatePengumumanHTML('2024-12-16', '2024-12-22');
      
      // Buka di window baru untuk print
      const printWindow = window.open('', '_blank');
      printWindow.document.write(html);
      printWindow.document.close();
      printWindow.print();
      
      // Atau tampilkan di iframe
      // const iframe = document.createElement('iframe');
      // document.body.appendChild(iframe);
      // iframe.contentDocument.write(html);
      // iframe.contentDocument.close();
      // iframe.contentWindow.print();
    } catch (error) {
      alert('Gagal generate pengumuman');
    }
  };

  // Atau konversi ke PDF menggunakan library seperti jsPDF atau html2pdf
  import html2pdf from 'html2pdf.js';

  const handleExportPDF = async () => {
    try {
      const html = await generatePengumumanHTML('2024-12-16', '2024-12-22');
      
      const element = document.createElement('div');
      element.innerHTML = html;
      
      const opt = {
        margin: 1,
        filename: 'pengumuman-nikah.pdf',
        image: { type: 'jpeg', quality: 0.98 },
        html2canvas: { scale: 2 },
        jsPDF: { unit: 'in', format: 'a4', orientation: 'portrait' }
      };
      
      await html2pdf().set(opt).from(element).save();
    } catch (error) {
      alert('Gagal export PDF');
    }
  };
  ```

  ---

  ## 👨‍⚖️ Penghulu Endpoints

  ### 17. List Marriage Officers

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

  ### 18. Update Marriage Officer

  **Endpoint:** `PUT /simnikah/penghulu/:id`

  **Description:** Mengupdate data penghulu.

  **Auth Required:** ✅ Yes

  **Role Required:** `kepala_kua`

  ---

  ### 19. Get My Assignments

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

  ### 20. Complete Marriage

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

  ### 21. Create Staff

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

  ### 22. Create Marriage Officer

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

  ### 23. Assign Marriage Officer

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

  ### 24. List Available Officers

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

  ### 25. Get Penghulu Statistics

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

  ### 26. Get Approved Registrations Per Week (Kepala KUA)

  **Endpoint:** `GET /simnikah/kepala-kua/pengumuman-nikah/list`

  **Description:** Mendapatkan daftar pendaftaran nikah untuk periode yang dipilih frontend beserta data kop surat. **Menampilkan semua status pendaftaran kecuali "Ditolak"** (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai). Data ini digunakan oleh frontend untuk generate surat pengumuman nikah dalam format HTML.

  **Auth Required:** ✅ Yes

  **Role Required:** `kepala_kua`

  **Query Parameters:**
  - `tanggal_awal` (optional): Tanggal awal periode yang dipilih frontend (format: YYYY-MM-DD). Frontend dapat memilih tanggal bebas sesuai kebutuhan (tidak harus mingguan). Default: awal minggu ini (Senin) jika tidak diisi
  - `tanggal_akhir` (optional): Tanggal akhir periode yang dipilih frontend (format: YYYY-MM-DD). Frontend dapat memilih tanggal bebas sesuai kebutuhan (tidak harus mingguan). Default: akhir minggu ini (Minggu) jika tidak diisi
  
  **Note:** Frontend dapat memilih periode tanggal secara bebas sesuai kebutuhan. Contoh: bisa memilih 1 bulan, 2 minggu, atau periode custom lainnya. Tidak terbatas pada periode mingguan.

  **Request Body (Optional - untuk custom kop surat):**
  ```json
  {
    "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
    "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
    "kota": "Kota Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "kode_pos": "70123",
    "telepon": "0511-1234567",
    "email": "kua.banjarmasinutara@kemenag.go.id",
    "website": "https://kua.banjarmasinutara.go.id",
    "logo_url": "https://example.com/logo-kua.png"
  }
  ```

  **Note:** 
  - Jika request body tidak dikirim, akan menggunakan nilai default untuk kop surat
  - Frontend bertanggung jawab untuk generate HTML surat pengumuman nikah menggunakan data yang dikembalikan

  **Examples:**
  
  **Mingguan (default):**
  ```
  GET /simnikah/kepala-kua/pengumuman-nikah/list?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
  ```
  
  **Bulanan:**
  ```
  GET /simnikah/kepala-kua/pengumuman-nikah/list?tanggal_awal=2024-12-01&tanggal_akhir=2024-12-31
  ```
  
  **Custom periode:**
  ```
  GET /simnikah/kepala-kua/pengumuman-nikah/list?tanggal_awal=2024-12-10&tanggal_akhir=2024-12-25
  ```
  
  **Tanpa parameter (menggunakan default - minggu ini):**
  ```
  GET /simnikah/kepala-kua/pengumuman-nikah/list
  ```

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Data pendaftaran disetujui berhasil diambil",
    "data": {
      "tanggal_awal": "2024-12-16",
      "tanggal_akhir": "2024-12-22",
      "periode": "16 Desember 2024 s/d 22 Desember 2024",
      "total": 5,
      "kop_surat": {
        "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
        "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
        "kota": "Kota Banjarmasin",
        "provinsi": "Kalimantan Selatan",
        "kode_pos": "70123",
        "telepon": "-",
        "email": "-",
        "website": "",
        "logo_url": ""
      },
      "registrations": [
        {
          "id": 1,
          "nomor_pendaftaran": "NIKAH-20241215-1234",
          "status_pendaftaran": "Disetujui",
          "tanggal_nikah": "2024-12-18T10:00:00Z",
          "waktu_nikah": "10:00",
          "tempat_nikah": "Di KUA",
          "alamat_akad": "Kantor KUA Kecamatan Banjarmasin Utara",
          "calon_suami": {
            "nama_lengkap": "Ahmad bin Abdullah"
          },
          "calon_istri": {
            "nama_lengkap": "Siti binti Muhammad"
          },
          "wali_nikah": {
            "nama_dan_bin": "Muhammad bin Ali",
            "hubungan_wali": "Ayah Kandung"
          }
        }
      ]
    }
  }
  ```

  **Catatan Penting:**
  - **Status yang ditampilkan:** Semua status kecuali "Ditolak" (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai)
  - **Status "Ditolak" tidak ditampilkan** karena pendaftaran yang ditolak tidak perlu diumumkan
  - Response includes `status_pendaftaran` untuk setiap pendaftaran agar frontend dapat menampilkan badge/indikator status

  **Use Case:**
  - Frontend memanggil endpoint ini untuk mendapatkan data pendaftaran dan kop surat
  - Frontend menggunakan data tersebut untuk generate HTML surat pengumuman nikah
  - Frontend dapat menampilkan badge status untuk setiap pendaftaran (Draft = kuning, Disetujui = hijau, dll)
  - Surat dapat dicetak atau dikonversi ke PDF di frontend
  - Surat dicetak dan dipasang di papan pengumuman KUA

  ---

  ### 26.5. Generate Pengumuman Nikah HTML (Kepala KUA)

  **Endpoint:** `GET /simnikah/kepala-kua/pengumuman-nikah/generate`  
  **Endpoint:** `POST /simnikah/kepala-kua/pengumuman-nikah/generate`

  **Description:** Generate surat pengumuman nikah dalam format HTML yang siap dicetak atau dikonversi ke PDF. Endpoint ini mengembalikan HTML document lengkap dengan kop surat, tabel data pendaftaran, dan format surat resmi.

  **Auth Required:** ✅ Yes

  **Role Required:** `kepala_kua`

  **Query Parameters:**
  - `tanggal_awal` (optional): Tanggal awal periode (format: YYYY-MM-DD). Default: awal minggu ini (Senin)
  - `tanggal_akhir` (optional): Tanggal akhir periode (format: YYYY-MM-DD). Default: akhir minggu ini (Minggu)

  **Request Body (Optional - untuk custom kop surat):**
  ```json
  {
    "tanggal_awal": "2024-12-16",
    "tanggal_akhir": "2024-12-22",
    "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
    "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
    "kota": "Kota Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "kode_pos": "70123",
    "telepon": "0511-1234567",
    "email": "kua.banjarmasinutara@kemenag.go.id",
    "website": "https://kua.banjarmasinutara.go.id",
    "logo_url": "https://example.com/logo-kua.png"
  }
  ```

  **Request Body Fields:**

  | Field | Type | Required | Default | Deskripsi |
  |-------|------|----------|---------|-----------|
  | `tanggal_awal` | string (YYYY-MM-DD) | ❌ No | Awal minggu ini | Tanggal awal periode (fallback jika query param tidak ada) |
  | `tanggal_akhir` | string (YYYY-MM-DD) | ❌ No | Akhir minggu ini | Tanggal akhir periode (fallback jika query param tidak ada) |
  | `nama_kua` | string | ❌ No | "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA" | Nama KUA |
  | `alamat_kua` | string | ❌ No | "PH5Q+F8C, Jl. Wira Karya, Pangeran" | Alamat KUA |
  | `kota` | string | ❌ No | "Kota Banjarmasin" | Kota |
  | `provinsi` | string | ❌ No | "Kalimantan Selatan" | Provinsi |
  | `kode_pos` | string | ❌ No | "70123" | Kode pos |
  | `telepon` | string | ❌ No | "-" | Nomor telepon |
  | `email` | string | ❌ No | "-" | Email |
  | `website` | string | ❌ No | "" | Website (opsional) |
  | `logo_url` | string | ❌ No | "" | URL logo KUA (opsional, harus accessible) |

  **Note:** 
  - Query parameters memiliki prioritas lebih tinggi daripada request body
  - Jika query parameters tidak ada, akan membaca dari request body
  - Jika keduanya tidak ada, akan menggunakan default (minggu ini)
  - Logo URL harus accessible (public URL) dan disarankan menggunakan HTTPS
  - Semua field kop surat bersifat opsional, akan menggunakan default jika tidak diisi

  **Examples:**

  **1. Generate dengan default kop surat (GET):**
  ```
  GET /simnikah/kepala-kua/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
  Authorization: Bearer YOUR_TOKEN
  ```

  **2. Generate dengan custom kop surat (POST):**
  ```
  POST /simnikah/kepala-kua/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
  Authorization: Bearer YOUR_TOKEN
  Content-Type: application/json

  {
    "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
    "alamat_kua": "Jl. Contoh No. 123",
    "kota": "Kota Banjarmasin",
    "provinsi": "Kalimantan Selatan",
    "kode_pos": "70123",
    "telepon": "0511-1234567",
    "email": "kua@example.com",
    "website": "https://kua.example.com",
    "logo_url": "https://example.com/logo.png"
  }
  ```

  **Response Success (200):**

  **Content-Type:** `text/html; charset=utf-8`

  Response berupa **HTML document lengkap** yang siap dicetak atau dikonversi ke PDF.

  **Response Headers:**
  ```
  Content-Type: text/html; charset=utf-8
  Content-Length: <ukuran file>
  ```

  **Response Body:** HTML document lengkap dengan struktur yang sama seperti endpoint Staff (lihat dokumentasi endpoint 15.5) - format Excel standar KUA dengan 15 kolom dan layout landscape A4.

  **Catatan Penting:**
  - HTML sudah include CSS untuk print optimization (`@media print`)
  - Format **A4 Landscape** untuk menampung tabel lebar
  - Font: Times New Roman (serif)
  - Font size: 8-9pt untuk tabel, 11-12pt untuk kop surat
  - Tabel dengan border untuk kejelasan
  - Header tabel dengan rowspan/colspan untuk grouping kolom
  - Siap untuk dicetak langsung atau dikonversi ke PDF
  - **Status yang ditampilkan:** Semua status kecuali "Ditolak" (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai)
  - **Usia dihitung otomatis** dari tanggal lahir calon pengantin
  - **Format waktu:** HH:MM diubah menjadi HH.MM (contoh: 08:00 → 08.00)
  - **Nama hari:** Otomatis dalam bahasa Indonesia (SENIN, SELASA, dll)

  **Response Error:** Sama seperti endpoint Staff (lihat dokumentasi endpoint 15.5).

  **Contoh Implementasi Frontend:** Sama seperti endpoint Staff (lihat dokumentasi endpoint 15.5).

  ---

  ## 💬 Feedback Endpoints

  ### 28. Create Feedback

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

  ### 29. List Feedback (Kepala KUA)

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

  ### 30. Mark Feedback As Read

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

  ### 31. Get Feedback Statistics

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

  ### 32. Geocode Address

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

  ### 33. Reverse Geocode

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

  ### 34. Get Wedding Location Detail

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

  ### 35. Get User Notifications

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

  ### 36. Create Notification

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

  ## 📊 Dashboard & Analytics Endpoints

  ### 1. Get Kepala KUA Dashboard

  **Endpoint:** `GET /simnikah/dashboard/kepala-kua`

  **Description:** Mendapatkan data dashboard lengkap untuk Kepala KUA termasuk statistik, trends, status distribution, penghulu performance, dan peak hours analysis.

  **Auth Required:** ✅ Yes

  **Role Required:** `kepala_kua`

  **Query Parameters:**
  - `period` (optional): Periode data (`day`, `week`, `month`, `year`). Default: `month`
  - `date_from` (optional): Tanggal mulai (YYYY-MM-DD). Jika tidak diisi, akan menggunakan periode default
  - `date_to` (optional): Tanggal akhir (YYYY-MM-DD). Jika tidak diisi, akan menggunakan periode default

  **Example:**
  ```
  GET /simnikah/dashboard/kepala-kua?period=month
  GET /simnikah/dashboard/kepala-kua?date_from=2024-01-01&date_to=2024-01-31
  ```

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Dashboard data berhasil diambil",
    "data": {
      "period": {
        "type": "month",
        "date_from": "2024-01-01",
        "date_to": "2024-01-31"
      },
      "statistics": {
        "total_periode": 45,
        "hari_ini": 2,
        "bulan_ini": 45,
        "tahun_ini": 450,
        "selesai": 30,
        "pending": 15,
        "status_breakdown": {
          "draft": 5,
          "disetujui": 8,
          "menunggu_penugasan": 2,
          "penghulu_ditugaskan": 5,
          "selesai": 30,
          "ditolak": 2
        }
      },
      "trends": [
        {
          "date": "2024-01-01",
          "label": "01 Jan",
          "count": 2
        },
        {
          "date": "2024-01-02",
          "label": "02 Jan",
          "count": 3
        }
      ],
      "status_distribution": [
        {
          "status": "Draft",
          "count": 5,
          "label": "Draft"
        },
        {
          "status": "Disetujui",
          "count": 8,
          "label": "Disetujui"
        },
        {
          "status": "Menunggu Penugasan",
          "count": 2,
          "label": "Menunggu Penugasan"
        },
        {
          "status": "Penghulu Ditugaskan",
          "count": 5,
          "label": "Penghulu Ditugaskan"
        },
        {
          "status": "Selesai",
          "count": 30,
          "label": "Selesai"
        },
        {
          "status": "Ditolak",
          "count": 2,
          "label": "Ditolak"
        }
      ],
      "penghulu_performance": [
        {
          "penghulu_id": 1,
          "nama_lengkap": "H. Ahmad Fauzi, S.Ag",
          "jumlah_nikah": 15,
          "rating": 4.5,
          "jumlah_rating": 12
        }
      ],
      "peak_hours": [
        {
          "waktu": "08:00",
          "count": 5
        },
        {
          "waktu": "09:00",
          "count": 8
        }
      ]
    }
  }
  ```

  ---

  ### 2. Get Staff Dashboard

  **Endpoint:** `GET /simnikah/dashboard/staff`

  **Description:** Mendapatkan data dashboard untuk Staff termasuk pending verifications, documents yang perlu diverifikasi, dan timeline aktivitas.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Dashboard staff berhasil diambil",
    "data": {
      "pending_verifications": [
        {
          "id": 1,
          "nomor_pendaftaran": "NIKAH-20240101-1234",
          "status_pendaftaran": "Draft",
          "tanggal_nikah": "2024-02-14",
          "calon_suami": "Ahmad Wijaya bin Abdullah",
          "calon_istri": "Siti Nurhaliza binti Muhammad",
          "created_at": "2024-01-01T10:00:00Z"
        }
      ],
      "pending_documents": [
        {
          "id": 2,
          "nomor_pendaftaran": "NIKAH-20240102-5678",
          "status_pendaftaran": "Disetujui",
          "calon_suami": "Budi Santoso bin Ahmad",
          "calon_istri": "Dewi Lestari binti Hasan",
          "created_at": "2024-01-02T11:00:00Z",
          "needs_verification": true
        }
      ],
      "timeline": [
        {
          "id": 1,
          "nomor_pendaftaran": "NIKAH-20240101-1234",
          "status_pendaftaran": "Disetujui",
          "calon_suami": "Ahmad Wijaya bin Abdullah",
          "calon_istri": "Siti Nurhaliza binti Muhammad",
          "updated_at": "2024-01-01T14:00:00Z",
          "action": "Status diubah menjadi Disetujui"
        }
      ]
    }
  }
  ```

  ---

  ### 3. Get Marriage Statistics

  **Endpoint:** `GET /simnikah/dashboard/statistik-pernikahan`

  **Description:** Mendapatkan statistik pernikahan detail dengan trends data untuk chart.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`, `kepala_kua`

  **Query Parameters:**
  - `period` (optional): Periode data (`day`, `month`, `year`). Default: `month`

  **Example:**
  ```
  GET /simnikah/dashboard/statistik-pernikahan?period=year
  ```

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Statistik pernikahan berhasil diambil",
    "data": {
      "statistics": {
        "total_periode": 450,
        "hari_ini": 2,
        "bulan_ini": 45,
        "tahun_ini": 450,
        "selesai": 300,
        "pending": 150,
        "status_breakdown": {
          "draft": 20,
          "disetujui": 50,
          "menunggu_penugasan": 30,
          "penghulu_ditugaskan": 50,
          "selesai": 300,
          "ditolak": 10
        }
      },
      "trends": [
        {
          "date": "2024-01",
          "label": "Jan 2024",
          "count": 35
        }
      ]
    }
  }
  ```

  ---

  ### 4. Get Penghulu Performance

  **Endpoint:** `GET /simnikah/dashboard/penghulu-performance`

  **Description:** Mendapatkan data performance semua penghulu termasuk jumlah nikah dan rating.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`, `kepala_kua`

  **Query Parameters:**
  - `date_from` (optional): Tanggal mulai (YYYY-MM-DD). Default: awal bulan ini
  - `date_to` (optional): Tanggal akhir (YYYY-MM-DD). Default: akhir bulan ini

  **Example:**
  ```
  GET /simnikah/dashboard/penghulu-performance?date_from=2024-01-01&date_to=2024-01-31
  ```

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Statistik penghulu berhasil diambil",
    "data": [
      {
        "penghulu_id": 1,
        "nama_lengkap": "H. Ahmad Fauzi, S.Ag",
        "jumlah_nikah": 15,
        "rating": 4.5,
        "jumlah_rating": 12
      }
    ]
  }
  ```

  ---

  ### 5. Get Peak Hours Analysis

  **Endpoint:** `GET /simnikah/dashboard/peak-hours`

  **Description:** Mendapatkan analisis jam sibuk (peak hours) untuk pernikahan.

  **Auth Required:** ✅ Yes

  **Role Required:** `staff`, `kepala_kua`

  **Query Parameters:**
  - `date_from` (optional): Tanggal mulai (YYYY-MM-DD). Default: awal bulan ini
  - `date_to` (optional): Tanggal akhir (YYYY-MM-DD). Default: akhir bulan ini

  **Example:**
  ```
  GET /simnikah/dashboard/peak-hours?date_from=2024-01-01&date_to=2024-01-31
  ```

  **Response Success (200):**
  ```json
  {
    "success": true,
    "message": "Analisis jam sibuk berhasil diambil",
    "data": [
      {
        "waktu": "08:00",
        "count": 5
      },
      {
        "waktu": "09:00",
        "count": 8
      },
      {
        "waktu": "10:00",
        "count": 12
      }
    ]
  }
  ```

  ---

  ## 📊 Status Flow

  ### Status Pendaftaran Nikah (Flow Baru)

  ```
  1. Draft
    └─> Catin mengisi form pendaftaran sederhana
    └─> Status: Kuning (belum pasti)

  2. Disetujui
    └─> Staff menyetujui setelah verifikasi di belakang layar
    └─> Status: Hijau (sudah pasti)

  3. Menunggu Penugasan
    └─> Menunggu kepala KUA menentukan penghulu

  4. Penghulu Ditugaskan
    └─> Kepala KUA sudah assign penghulu
    └─> Penghulu dapat melihat tugasnya

  5. Selesai
    └─> Penghulu sudah melaksanakan nikah
    └─> Status final

  (6. Ditolak - jika ditolak oleh staff)
    └─> Pendaftaran ditolak dan tidak ditampilkan di kalender
  ```

  ### Status Constants

  - `Draft` - Pendaftaran baru dibuat (form sederhana), status kuning (belum pasti)
  - `Disetujui` - Staff sudah approve, status hijau (sudah pasti)
  - `Menunggu Penugasan` - Menunggu kepala KUA assign penghulu
  - `Penghulu Ditugaskan` - Penghulu sudah ditugaskan oleh kepala KUA
  - `Selesai` - Pernikahan sudah dilaksanakan oleh penghulu
  - `Ditolak` - Pendaftaran ditolak oleh staff (tidak ditampilkan di kalender)

  ### Catatan Penting

  1. **Status Draft vs Disetujui di Kalender:**
    - Draft (kuning): Belum pasti, **tetap dihitung dalam kuota** untuk mencegah double booking
    - Disetujui (hijau): Sudah pasti, dihitung dalam kuota

  2. **Kuota Ketersediaan:**
    - KUA: Maksimal 1 pernikahan per jam
    - Luar KUA: Maksimal 3 total per jam
    - **Draft dan Disetujui dihitung dalam kuota** (hanya "Ditolak" dan "Selesai" yang tidak dihitung)
    - Total maksimal per jam: 3 pernikahan (1 di KUA + 2 di luar KUA, atau 3 di luar KUA)

  3. **Aturan Jam/Waktu:**
    - Jam operasional: 08:00 - 16:00 WITA
    - Time slots: 08:00, 09:00, 10:00, 11:00, 12:00, 13:00, 14:00, 15:00, 16:00
    - Format waktu: `HH:MM` (24 jam)
    - Tanggal nikah tidak boleh di masa lalu

  ---

  ## 🎯 Role-Based Access Control

  | Role | Description | Endpoints Access |
  |------|-------------|------------------|
  | `user_biasa` | Calon pengantin | Pendaftaran, status, kalender, feedback |
  | `staff` | Staff KUA | Verifikasi, approve, update status, dashboard staff |
  | `penghulu` | Penghulu | Lihat tugas, complete nikah |
  | `kepala_kua` | Kepala KUA | Assign penghulu, statistik, manajemen staff/penghulu, feedback management, dashboard lengkap |

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

  - **Alur Pendaftaran & Aturan Jam:** Lihat `docs/ALUR_PENDAFTARAN_DAN_JAM.md`
  - **Status Flow:** Lihat `docs/FLOW_SEDERHANA.md`
  - **Architecture:** Lihat `docs/architecture/`
  - **Testing:** Lihat `docs/api/API_TESTING_DOCUMENTATION.md`

  ---

  **Generated:** 2024  
  **API Version:** 1.0.0

