# Dokumentasi Lengkap Penggunaan API SimNikah

Dokumentasi ini mencakup seluruh endpoint API yang tersedia di sistem SimNikah, dikelompokkan berdasarkan role pengguna.

## Base URL
`http://localhost:8080`

## Daftar Isi
1. [Authentication](#1-authentication)
2. [Catin (Calon Pengantin)](#2-catin-calon-pengantin)
3. [Staff KUA](#3-staff-kua)
4. [Penghulu](#4-penghulu)
5. [Kepala KUA](#5-kepala-kua)
6. [Dashboard & Statistik](#6-dashboard--statistik)
7. [Publik / Umum](#7-publik--umum)

---

## 1. Authentication

Semua role menggunakan endpoint ini untuk masuk ke sistem.

### Register (User Biasa/Catin)
`POST /register`
```json
{
  "username": "ahmad123",
  "email": "ahmad@example.com",
  "password": "password123",
  "nama": "Ahmad bin Abdullah"
}
```

### Login
`POST /login`
```json
{
  "email": "ahmad@example.com",
  "password": "password123"
}
```
**Response:** Token JWT (gunakan di header `Authorization: Bearer <token>`)

### Get Profile
`GET /profile`
- Mendapatkan data user yang sedang login.

---

## 2. Catin (Calon Pengantin)

Modul untuk pendaftaran nikah dan pengecekan status.

### Pendaftaran Nikah
`POST /simnikah/pendaftaran`
- Mendaftarkan pernikahan baru (Formulir Sederhana).
- **Body:** Lihat [Dokumentasi Body Request](API_REQUEST_BODY_DOCUMENTATION.md#1-pendaftaran-nikah)

### Cek Status Pendaftaran Saya
`GET /simnikah/pendaftaran/status`
- Melihat status pendaftaran milik user yang login.
- **Response:** Detail pendaftaran, status (Draft/Disetujui/dll), data calon, dan penghulu (jika ada).

### Beri Feedback (Setelah Nikah)
`POST /simnikah/feedback-pernikahan`
```json
{
  "pendaftaran_id": 1,
  "jenis_feedback": "Rating",
  "rating": 5,
  "judul": "Pelayanan Memuaskan",
  "pesan": "Terima kasih bapak penghulu"
}
```

---

## 3. Staff KUA

Modul untuk manajemen pendaftaran dan administrasi.

### List Pendaftaran
`GET /simnikah/staff/pendaftaran`
- Filter: `?status=Draft`, `?date_from=...`, `?location=Di KUA`

### Buat Pendaftaran (Atas Nama User)
`POST /simnikah/staff/pendaftaran`
- Staff mendaftarkan user manual (datang ke kantor).

### Verifikasi Formulir
`POST /simnikah/staff/verify-formulir/:id`
- Memverifikasi kelengkapan data formulir.
```json
{ "status": "Disetujui", "catatan": "Data lengkap" }
```

### Approve Pendaftaran
`POST /simnikah/staff/approve/:id`
- Menyetujui pendaftaran untuk lanjut ke penugasan penghulu.
```json
{ "status": "Menunggu Penugasan", "catatan": "Siap ditugaskan" }
```

### Update Status Pendaftaran
`PUT /simnikah/pendaftaran/:id/update-status`
- Update status manual jika diperlukan.

### List Pengumuman Nikah
`GET /simnikah/staff/pengumuman-nikah/list`
- Mengambil data pernikahan yang disetujui untuk dicetak/diumumkan.

### Manajemen Staff & Penghulu
- `GET /simnikah/staff` (List staff)
- `GET /simnikah/penghulu` (List penghulu)

---

## 4. Penghulu

Modul khusus untuk petugas penghulu.

### List Tugas Saya
`GET /simnikah/penghulu/assigned-registrations`
- Melihat daftar pernikahan yang ditugaskan ke penghulu yang login.

### Jadwal Hari Ini
`GET /simnikah/penghulu/today-schedule`
- Shortcut untuk melihat jadwal nikah hari ini.

### Verifikasi Dokumen Fisik
`POST /simnikah/penghulu/verify-documents/:id`
- Penghulu memverifikasi dokumen fisik saat pemeriksaan.

### Selesaikan Pernikahan
`POST /simnikah/penghulu/complete-marriage/:id`
- Menandai pernikahan telah dilaksanakan.
```json
{ "catatan": "Pernikahan lancar, buku nikah diserahkan" }
```

---

## 5. Kepala KUA

Modul untuk pimpinan KUA (Monitoring & Approval).

### Manajemen User Internal
- `POST /simnikah/kepala-kua/staff` (Tambah Staff)
- `POST /simnikah/kepala-kua/penghulu` (Tambah Penghulu)

### Penugasan Penghulu (Assign)
`POST /simnikah/pendaftaran/:id/assign-penghulu`
- Menugaskan penghulu ke pendaftaran tertentu.
- **Validasi:** 1 Penghulu max 1 nikah/jam.
```json
{ "penghulu_id": 2, "catatan": "Tugas untuk Pak Budi" }
```

### Cek Ketersediaan Penghulu
- `GET /simnikah/kepala-kua/available-penghulu` (List semua penghulu aktif)
- `GET /simnikah/kepala-kua/penghulu-schedule` (Jadwal harian semua penghulu)
- `GET /simnikah/kepala-kua/penghulu-tersedia` (Cek siapa kosong di jam X)

### Statistik & Laporan
- `GET /simnikah/kepala-kua/statistik-penghulu` (Beban kerja penghulu)
- `GET /simnikah/kepala-kua/pengumuman-nikah/list` (Laporan mingguan)

### Feedback Management
- `GET /simnikah/kepala-kua/feedback` (Lihat masukan masyarakat)
- `GET /simnikah/kepala-kua/feedback/stats` (Statistik kepuasan)
- `PUT /simnikah/kepala-kua/feedback/:id/mark-read` (Tandai dibaca)

---

## 6. Dashboard & Statistik

Endpoint untuk data visualisasi dashboard.

### Dashboard Kepala KUA
`GET /simnikah/dashboard/kepala-kua`
- Ringkasan total nikah, performa, feedback, dll.

### Dashboard Staff
`GET /simnikah/dashboard/staff`
- Ringkasan tugas verifikasi, pendaftaran baru hari ini.

### Analisis
- `GET /simnikah/dashboard/statistik-pernikahan` (Grafik bulanan)
- `GET /simnikah/dashboard/penghulu-performance` (Rating & jumlah nikah)
- `GET /simnikah/dashboard/peak-hours` (Analisis jam sibuk)

---

## 7. Publik / Umum

Endpoint informasi yang bisa diakses tanpa login (atau Catin).

### Kalender Ketersediaan
`GET /simnikah/kalender-ketersediaan?bulan=12&tahun=2024`
- Melihat tanggal mana yang penuh/tersedia dalam satu bulan.

### Cek Jam Tersedia
`GET /simnikah/ketersediaan-jam?tanggal=2024-12-25`
- Melihat slot jam (08:00 - 16:00) yang masih kosong pada tanggal tertentu.
- Membedakan kuota "Di KUA" dan "Di Luar KUA".

### List Pernikahan Tanggal Tertentu
`GET /simnikah/pernikahan-tanggal?tanggal=2024-12-25`
- Melihat siapa saja yang menikah di tanggal tersebut (Jadwal Publik).

### Geocoding & Lokasi
`POST /simnikah/location/geocode`
- Mengubah alamat teks menjadi koordinat.

`POST /simnikah/location/reverse-geocode`
- Mengubah koordinat menjadi alamat.

`GET /simnikah/location/search`
- Mencari saran alamat (autocomplete).

`PUT /simnikah/pendaftaran/:id/location`
- Update koordinat lokasi nikah.

### Notifikasi (Semua Role)
`GET /simnikah/notifikasi/user/:user_id`
- Mengambil daftar notifikasi user.

`PUT /simnikah/notifikasi/user/:user_id/mark-all-read`
- Menandai semua notifikasi sudah dibaca.

`PUT /simnikah/notifikasi/:id/status`
- Menandai satu notifikasi sudah dibaca.

`POST /simnikah/notifikasi`
- Membuat notifikasi baru (Staff/Kepala KUA).

`POST /simnikah/notifikasi/send-to-role`
- Mengirim notifikasi broadcast ke role tertentu (contoh: semua penghulu).

---

## Alur Sistem Utama

1. **Catin** mendaftar online (`POST /pendaftaran`). Status: **Draft**.
2. **Staff** memverifikasi formulir (`POST /verify-formulir`). Status: **Disetujui**.
3. **Staff** melakukan approval (`POST /approve`). Status: **Menunggu Penugasan**.
4. **Kepala KUA** melihat jadwal & menugaskan penghulu (`POST /assign-penghulu`). Status: **Penghulu Ditugaskan**.
5. **Penghulu** memverifikasi dokumen fisik & melaksanakan akad (`POST /complete-marriage`). Status: **Selesai**.
6. **Catin** memberikan feedback (`POST /feedback`).

## Catatan Penting
- **Limitasi:** Maksimal 3 pernikahan per jam (total). Maksimal 1 "Di KUA" per jam.
- **Hari Libur:** Sabtu, Minggu, dan Hari Libur Nasional (KUA tutup, tapi "Di Luar KUA" bisa tetap jalan jika ada penghulu).
