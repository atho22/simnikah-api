# Frontend Implementation Guide - SimNikah
**Panduan Implementasi Lengkap untuk Frontend Developer**

Dokumen ini menjelaskan secara detail halaman apa yang harus dibuat, endpoint API yang digunakan, dan flow interaksi user.

---

## 📋 Daftar Halaman Berdasarkan Role

### 🌐 Publik (Tanpa Login)

#### 1. Landing Page (`/`)
**Tujuan:** Halaman utama untuk menarik user mendaftar atau login.

**Komponen:**
- Hero Section (Judul, Deskripsi, CTA Button)
- Fitur Unggulan (3-4 card)
- Footer (Kontak KUA)

**API:** Tidak ada (static content).

---

#### 2. Kalender Ketersediaan (`/jadwal`)
**Tujuan:** User bisa cek tanggal mana yang masih tersedia untuk nikah.

**Flow:**
1. User pilih bulan & tahun (default: bulan ini).
2. Tampilkan kalender dengan warna:
   - 🟢 Hijau: Tersedia
   - 🔴 Merah: Penuh / Libur
   - 🟡 Kuning: Hampir penuh
3. User klik tanggal → Modal muncul menampilkan detail jam tersedia.

**API Endpoint:**
```
GET /simnikah/kalender-ketersediaan?bulan=12&tahun=2024
```

**Response Example:**
```json
{
  "data": {
    "calendar": [
      {
        "tanggal": 25,
        "status": "Tersedia",
        "tersedia_kua": true,
        "tersedia_luar_kua": true,
        "time_slots": [
          { "waktu": "09:00", "kua": { "tersedia": true }, "luar_kua": { "tersedia": true } }
        ]
      }
    ]
  }
}
```

**Komponen UI:**
- `<Calendar />` component (bisa pakai `react-day-picker` atau custom).
- `<TimeSlotModal />` untuk detail jam.

---

#### 3. Login (`/login`)
**Tujuan:** User login dengan email & password.

**Flow:**
1. User input email & password.
2. Submit form → Hit API `/login`.
3. Simpan token JWT di Cookie atau LocalStorage.
4. Redirect ke dashboard sesuai role.

**API Endpoint:**
```
POST /login
Body: { "email": "user@example.com", "password": "password123" }
```

**Response:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR123",
    "nama": "Ahmad",
    "role": "user_biasa"
  }
}
```

**Redirect Logic:**
- `role: "user_biasa"` → `/dashboard/catin`
- `role: "staff"` → `/dashboard/staff`
- `role: "penghulu"` → `/dashboard/penghulu`
- `role: "kepala_kua"` → `/dashboard/kepala-kua`

---

#### 4. Register (`/register`)
**Tujuan:** User baru (Catin) mendaftar akun.

**Flow:**
1. User input username, email, password, nama lengkap.
2. Submit → Hit API `/register`.
3. Auto login atau redirect ke `/login`.

**API Endpoint:**
```
POST /register
Body: { "username": "ahmad123", "email": "ahmad@example.com", "password": "pass123", "nama": "Ahmad bin Abdullah" }
```

---

## 👤 Role: Catin (Calon Pengantin)

### Dashboard Catin (`/dashboard/catin`)
**Tujuan:** Melihat status pendaftaran nikah.

**Flow:**
1. Fetch status pendaftaran user yang login.
2. Jika belum daftar → Tampilkan tombol "Daftar Nikah".
3. Jika sudah daftar → Tampilkan card status (Draft/Disetujui/dll).

**API Endpoint:**
```
GET /simnikah/pendaftaran/status
```

**Response:**
```json
{
  "data": {
    "has_registration": true,
    "registration": {
      "nomor_pendaftaran": "NIKAH-20241225-1234",
      "status_pendaftaran": "Draft",
      "tanggal_nikah": "2024-12-25",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di KUA",
      "calon_suami": { "nama_lengkap": "Ahmad bin Abdullah" },
      "calon_istri": { "nama_lengkap": "Siti binti Muhammad" }
    }
  }
}
```

**Komponen UI:**
- Card Status dengan badge warna (Draft = Kuning, Disetujui = Hijau).
- Timeline Progress (Draft → Verifikasi → Penghulu Ditugaskan → Selesai).

---

### Form Pendaftaran Nikah (`/dashboard/catin/daftar`)
**Tujuan:** User mengisi form pendaftaran nikah.

**Flow (Wizard Multi-Step):**

#### **Step 1: Data Calon Suami & Istri**
- Input: Nama & Bin/Binti, Pendidikan, Umur.

#### **Step 2: Data Wali Nikah**
- Input: Nama & Bin, Hubungan Wali (Ayah Kandung, Kakek, dll).

#### **Step 3: Jadwal & Lokasi**
1. User pilih tanggal (DatePicker).
2. Sistem fetch `/ketersediaan-jam?tanggal=YYYY-MM-DD`.
3. Tampilkan dropdown jam yang tersedia.
4. User pilih lokasi:
   - Radio: "Di KUA" atau "Di Luar KUA".
   - Jika "Di Luar KUA" → Tampilkan form alamat + peta (Leaflet).

#### **Step 4: Review & Submit**
- Tampilkan ringkasan data.
- Button "Kirim Pendaftaran".

**API Endpoint:**
```
POST /simnikah/pendaftaran
Body: {
  "calon_laki_laki": { "nama_dan_bin": "Ahmad bin Abdullah", "pendidikan_akhir": "S1", "umur": 25 },
  "calon_perempuan": { "nama_dan_binti": "Siti binti Muhammad", "pendidikan_akhir": "S1", "umur": 23 },
  "lokasi_nikah": {
    "tempat_nikah": "Di KUA",
    "tanggal_nikah": "2024-12-25",
    "waktu_nikah": "09:00"
  },
  "wali_nikah": { "nama_dan_bin": "Abdullah bin Muhammad", "hubungan_wali": "Ayah Kandung" }
}
```

**Validasi Penting:**
- Cek ketersediaan jam sebelum submit (bisa saja orang lain baru booking).
- Handle error `schedule_conflict` → Alert: "Maaf, jadwal baru saja diambil orang lain."

---

### Feedback (`/dashboard/catin/feedback`)
**Tujuan:** Catin memberi rating setelah nikah selesai.

**Flow:**
1. Halaman ini hanya muncul jika status = "Selesai".
2. User input rating (1-5 bintang), judul, pesan.
3. Submit → Hit API.

**API Endpoint:**
```
POST /simnikah/feedback-pernikahan
Body: {
  "pendaftaran_id": 1,
  "jenis_feedback": "Rating",
  "rating": 5,
  "judul": "Pelayanan Memuaskan",
  "pesan": "Terima kasih"
}
```

---

## 👔 Role: Staff KUA

### Dashboard Staff (`/dashboard/staff`)
**Tujuan:** Overview tugas staff hari ini.

**API Endpoint:**
```
GET /simnikah/dashboard/staff
```

**Response:**
```json
{
  "data": {
    "pending_verifikasi": 5,
    "nikah_hari_ini": 3,
    "total_pendaftaran_bulan_ini": 120
  }
}
```

---

### List Pendaftaran (`/dashboard/staff/pendaftaran`)
**Tujuan:** Melihat semua pendaftaran dengan filter.

**Flow:**
1. Tampilkan tabel data pendaftaran.
2. Filter: Status (Draft, Disetujui), Tanggal, Lokasi.
3. Klik row → Redirect ke halaman detail.

**API Endpoint:**
```
GET /simnikah/staff/pendaftaran?status=Draft&page=1&limit=10
```

**Response:**
```json
{
  "data": {
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIKAH-20241225-1234",
        "status_pendaftaran": "Draft",
        "tanggal_nikah": "2024-12-25",
        "calon_suami_nama": "Ahmad",
        "calon_istri_nama": "Siti"
      }
    ],
    "pagination": { "total": 50, "page": 1, "limit": 10 }
  }
}
```

**Komponen UI:**
- Table dengan kolom: Nomor, Nama Suami, Nama Istri, Tanggal, Status, Aksi.
- Badge status warna-warni.

---

### Detail & Verifikasi (`/dashboard/staff/pendaftaran/[id]`)
**Tujuan:** Verifikasi data pendaftaran.

**Flow:**
1. Tampilkan detail lengkap pendaftaran.
2. Tombol "Verifikasi Formulir" → Modal input catatan → Hit API.
3. Tombol "Approve" → Konfirmasi → Hit API.

**API Endpoint:**
```
POST /simnikah/staff/verify-formulir/:id
Body: { "status": "Disetujui", "catatan": "Data lengkap" }

POST /simnikah/staff/approve/:id
Body: { "status": "Menunggu Penugasan", "catatan": "Siap ditugaskan" }
```

---

### Pengumuman Nikah (`/dashboard/staff/pengumuman`)
**Tujuan:** Generate daftar nikah untuk dicetak/diumumkan.

**Flow:**
1. User pilih tanggal awal & akhir (default: minggu ini).
2. Fetch data → Tampilkan tabel.
3. Button "Print" → Window print browser.

**API Endpoint:**
```
GET /simnikah/staff/pengumuman-nikah/list?tanggal_awal=2024-12-01&tanggal_akhir=2024-12-31
```

**Response:**
```json
{
  "data": {
    "periode": "01 Desember 2024 s/d 31 Desember 2024",
    "total": 10,
    "registrations": [
      {
        "nomor_pendaftaran": "NIKAH-001",
        "tanggal_nikah": "2024-12-05",
        "waktu_nikah": "09:00",
        "calon_suami": { "nama_lengkap": "Ahmad" },
        "calon_istri": { "nama_lengkap": "Siti" }
      }
    ]
  }
}
```

---

## 👨‍⚖️ Role: Kepala KUA

### Dashboard Kepala KUA (`/dashboard/kepala-kua`)
**Tujuan:** Monitoring statistik & performa.

**API Endpoint:**
```
GET /simnikah/dashboard/kepala-kua
GET /simnikah/dashboard/statistik-pernikahan
GET /simnikah/dashboard/penghulu-performance
```

**Komponen UI:**
- Chart (Line/Bar) untuk grafik nikah per bulan.
- Card statistik: Total Nikah, Pending, Rating Rata-rata.

---

### Assign Penghulu (`/dashboard/kepala-kua/assign`)
**Tujuan:** Menugaskan penghulu ke pendaftaran.

**Flow:**
1. Tampilkan list pendaftaran status "Menunggu Penugasan".
2. Klik "Assign" → Modal pilih penghulu.
3. Sistem fetch penghulu yang tersedia di jam tersebut.
4. Submit → Hit API.

**API Endpoint:**
```
GET /simnikah/kepala-kua/penghulu-tersedia?tanggal=2024-12-25&waktu=09:00

Response:
{
  "data": {
    "tersedia": [
      { "id": 1, "nama_lengkap": "H. Ahmad", "jumlah_nikah_hari_ini": 2 }
    ],
    "tidak_tersedia": [
      { "id": 2, "nama_lengkap": "H. Budi", "alasan": "Sudah ada jadwal di jam ini" }
    ]
  }
}

POST /simnikah/pendaftaran/:id/assign-penghulu
Body: { "penghulu_id": 1, "catatan": "Tugas untuk Pak Ahmad" }
```

**Validasi:**
- Jika penghulu sudah ada jadwal di jam yang sama → Tampilkan error.

---

## 🕌 Role: Penghulu

### Dashboard Penghulu (`/dashboard/penghulu`)
**Tujuan:** Melihat jadwal nikah hari ini & mendatang.

**API Endpoint:**
```
GET /simnikah/penghulu/today-schedule
GET /simnikah/penghulu/assigned-registrations
```

**Response:**
```json
{
  "data": {
    "today": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIKAH-001",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "calon_suami": "Ahmad",
        "calon_istri": "Siti"
      }
    ]
  }
}
```

---

### Detail Tugas (`/dashboard/penghulu/tugas/[id]`)
**Tujuan:** Penghulu menyelesaikan pernikahan.

**Flow:**
1. Tampilkan detail pendaftaran + lokasi (peta jika di luar KUA).
2. Tombol "Selesaikan Pernikahan" → Modal input catatan → Hit API.

**API Endpoint:**
```
POST /simnikah/penghulu/complete-marriage/:id
Body: { "catatan": "Pernikahan lancar, buku nikah diserahkan" }
```

---

## 🔔 Notifikasi (Semua Role)

### Bell Icon di Header
**Tujuan:** Tampilkan notifikasi real-time.

**Flow:**
1. Polling setiap 30 detik atau saat pindah halaman.
2. Fetch notifikasi user yang login.
3. Tampilkan badge merah jika ada notifikasi belum dibaca.

**API Endpoint:**
```
GET /simnikah/notifikasi/user/:user_id
```

**Response:**
```json
{
  "data": {
    "notifikasi": [
      {
        "id": 1,
        "judul": "Pendaftaran Disetujui",
        "pesan": "Pendaftaran nikah Anda telah disetujui",
        "status_baca": "Belum Dibaca",
        "created_at": "2024-12-04T10:00:00Z"
      }
    ]
  }
}
```

**Mark as Read:**
```
PUT /simnikah/notifikasi/:id/status
Body: { "status_baca": "Sudah Dibaca" }
```

---

## 📌 Catatan Penting untuk Frontend Developer

1. **Token Management:** Simpan JWT token di Cookie (HttpOnly) atau LocalStorage. Kirim di header `Authorization: Bearer <token>` untuk setiap request.

2. **Error Handling:** Semua response error dari API memiliki format:
   ```json
   {
     "success": false,
     "message": "Error summary",
     "error": "Detail error",
     "type": "schedule_conflict"
   }
   ```
   Gunakan field `type` untuk menampilkan pesan error yang spesifik.

3. **Loading State:** Tampilkan spinner/skeleton saat fetch data.

4. **Validasi Form:** Gunakan library seperti `Zod` atau `Yup` untuk validasi sebelum submit.

5. **Responsive Design:** Pastikan semua halaman mobile-friendly (Tailwind CSS memudahkan ini).
