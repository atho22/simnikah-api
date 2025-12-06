# Product Requirements Document (PRD) - Frontend SimNikah

**Versi:** 1.0
**Status:** Draft
**Tanggal:** 4 Desember 2025

---

## 1. Ringkasan Produk
SimNikah adalah sistem informasi pernikahan digital yang menghubungkan Calon Pengantin (Catin), Staff KUA, Penghulu, dan Kepala KUA. Frontend aplikasi ini akan berbasis web (Next.js) dan berfungsi sebagai antarmuka untuk melakukan pendaftaran nikah, verifikasi berkas, penjadwalan penghulu, hingga pelaporan statistik.

## 2. Target Pengguna (User Roles)
1.  **Catin (Calon Pengantin):** Mendaftar nikah, cek status, lihat jadwal.
2.  **Staff KUA:** Verifikasi administrasi, manajemen pendaftaran.
3.  **Penghulu:** Melihat jadwal tugas, input hasil nikah.
4.  **Kepala KUA:** Monitoring, approval, assign penghulu, laporan.
5.  **Publik (Tanpa Login):** Cek ketersediaan jadwal, informasi.

---

## 3. Struktur Halaman & Fitur (Sitemap)

### A. Modul Publik (Tanpa Login)
| Halaman | Fitur Utama | Endpoint API Terkait |
| :--- | :--- | :--- |
| **Landing Page** (`/`) | Info layanan, tombol login/daftar, cek jadwal. | - |
| **Kalender Ketersediaan** (`/jadwal`) | Kalender interaktif menampilkan tanggal penuh/kosong. | `GET /simnikah/kalender-ketersediaan`<br>`GET /simnikah/ketersediaan-jam`<br>`GET /simnikah/pernikahan-tanggal` |
| **Login** (`/login`) | Form login email & password. | `POST /login` |
| **Register** (`/register`) | Form pendaftaran akun baru (Nama, Email, Password). | `POST /register` |

### B. Modul Catin (User Dashboard)
*Base URL: `/dashboard/catin`*

| Halaman | Fitur Utama | Endpoint API Terkait |
| :--- | :--- | :--- |
| **Dashboard Utama** (`/`) | Status pendaftaran terakhir (Draft/Disetujui), Notifikasi. | `GET /simnikah/pendaftaran/status`<br>`GET /simnikah/notifikasi/user/{id}` |
| **Form Pendaftaran** (`/daftar`) | Wizard form (Data Suami, Istri, Wali, Lokasi). Validasi tanggal & jam. | `POST /simnikah/pendaftaran`<br>`POST /simnikah/location/geocode` |
| **Feedback** (`/feedback`) | Form rating & ulasan setelah nikah selesai. | `POST /simnikah/feedback-pernikahan` |

### C. Modul Staff KUA
*Base URL: `/dashboard/staff`*

| Halaman | Fitur Utama | Endpoint API Terkait |
| :--- | :--- | :--- |
| **Dashboard Utama** (`/`) | Statistik singkat (Verifikasi Pending, Nikah Hari Ini). | `GET /simnikah/dashboard/staff` |
| **List Pendaftaran** (`/pendaftaran`) | Tabel data pendaftaran dengan filter (Status, Tanggal). | `GET /simnikah/staff/pendaftaran` |
| **Detail Verifikasi** (`/pendaftaran/[id]`) | View detail data, tombol **Verifikasi Formulir**, tombol **Approve**. | `POST /simnikah/staff/verify-formulir/{id}`<br>`POST /simnikah/staff/approve/{id}` |
| **Pengumuman Nikah** (`/pengumuman`) | Generate daftar nikah mingguan untuk dicetak. | `GET /simnikah/staff/pengumuman-nikah/list` |

### D. Modul Kepala KUA
*Base URL: `/dashboard/kepala-kua`*

| Halaman | Fitur Utama | Endpoint API Terkait |
| :--- | :--- | :--- |
| **Dashboard Eksekutif** (`/`) | Grafik statistik nikah, performa penghulu, jam sibuk. | `GET /simnikah/dashboard/kepala-kua`<br>`GET /simnikah/dashboard/statistik-pernikahan` |
| **Assign Penghulu** (`/assign`) | List pendaftaran "Menunggu Penugasan", modal pilih penghulu tersedia. | `GET /simnikah/kepala-kua/penghulu-tersedia`<br>`POST /simnikah/pendaftaran/{id}/assign-penghulu` |
| **Manajemen Pegawai** (`/pegawai`) | Tambah/Edit Staff dan Penghulu. | `POST /simnikah/kepala-kua/staff`<br>`POST /simnikah/kepala-kua/penghulu` |
| **Laporan Feedback** (`/laporan`) | List masukan dari masyarakat. | `GET /simnikah/kepala-kua/feedback` |

### E. Modul Penghulu
*Base URL: `/dashboard/penghulu`*

| Halaman | Fitur Utama | Endpoint API Terkait |
| :--- | :--- | :--- |
| **Jadwal Saya** (`/`) | List jadwal nikah hari ini & mendatang. | `GET /simnikah/penghulu/today-schedule`<br>`GET /simnikah/penghulu/assigned-registrations` |
| **Proses Nikah** (`/tugas/[id]`) | Detail lokasi, verifikasi berkas fisik, selesaikan nikah. | `POST /simnikah/penghulu/verify-documents/{id}`<br>`POST /simnikah/penghulu/complete-marriage/{id}` |

---

## 4. Kebutuhan Fungsional Detail

### 1. Kalender & Pendaftaran (UX Critical)
*   **Kalender:** Harus membedakan warna tanggal:
    *   Merah: Penuh / Libur.
    *   Hijau: Tersedia.
    *   Kuning: Hampir penuh.
*   **Validasi Jam:** Saat user memilih tanggal di form pendaftaran:
    *   Sistem harus fetch `GET /ketersediaan-jam`.
    *   Dropdown jam hanya menampilkan slot yang tersedia.
    *   Jika lokasi "Di KUA" dipilih, slot yang sudah ada nikah "Di KUA" harus didisable.

### 2. Lokasi & Peta
*   Integrasi dengan Leaflet/Google Maps.
*   Fitur pencarian alamat (`GET /location/search`) untuk mengisi koordinat otomatis (`POST /geocode`).

### 3. Notifikasi Real-time (Optional/Polling)
*   Frontend melakukan polling ke `GET /notifikasi/user/{id}` setiap X detik atau saat pindah halaman.
*   Tampilkan badge merah di icon lonceng header.

---

## 5. Tech Stack Recommendation
*   **Framework:** Next.js 14 (App Router)
*   **Language:** TypeScript
*   **Styling:** Tailwind CSS + shadcn/ui (untuk komponen cepat seperti Calendar, Table, Dialog).
*   **State Management:** React Query (TanStack Query) -> *Sangat disarankan untuk caching data API*.
*   **Forms:** React Hook Form + Zod (untuk validasi sesuai body request).
*   **Maps:** React Leaflet.
*   **Icons:** Lucide React.

---

## 6. Timeline & Prioritas Pengembangan
1.  **Fase 1 (Core):** Login, Register, Form Pendaftaran Catin (Draft), Kalender Publik.
2.  **Fase 2 (Staff Workflow):** Dashboard Staff, List Pendaftaran, Fitur Verifikasi & Approve.
3.  **Fase 3 (Kepala KUA Workflow):** Penugasan Penghulu, Manajemen Staff.
4.  **Fase 4 (Penghulu Workflow):** Dashboard Penghulu, Penyelesaian Nikah.
5.  **Fase 5 (Finishing):** Dashboard Statistik, Notifikasi, Feedback, Cetak Laporan.

---

## 7. Error Handling Requirement
Frontend harus menangani response error standard dari API:
```json
{
  "success": false,
  "message": "Error message summary",
  "type": "schedule_conflict" // Gunakan ini untuk UI yang spesifik
}
```
*   Jika `type: "schedule_conflict"`, tampilkan alert: "Maaf, jadwal tersebut baru saja diambil orang lain."
*   Jika `type: "holiday_restriction"`, tampilkan alert: "KUA tutup pada hari libur, silakan pilih lokasi 'Di Luar KUA' atau ganti tanggal."
