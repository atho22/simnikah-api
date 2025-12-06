# Frontend Task Breakdown - SimNikah

Dokumen ini berisi daftar tugas teknis untuk pengembangan frontend SimNikah menggunakan Next.js.

## Sprint 1: Setup & Auth (Pondasi)
Target: User bisa login, register, dan melihat layout dasar.

- [ ] **[F-001] Setup Project**
  - Initialize Next.js 14 (App Router, TypeScript).
  - Install Tailwind CSS, Shadcn/ui, Lucide React, Axios/Fetch wrapper.
  - Konfigurasi Environment (`.env.local`).
  - Setup Folder Structure (`app/(auth)`, `app/dashboard`, `components`, `lib`).

- [ ] **[F-002] Authentication Pages**
  - Halaman Login (`/login`): Form email/pass, handle API `/login`, simpan Token di Cookie/LocalStorage.
  - Halaman Register (`/register`): Form register catin, handle API `/register`.
  - Middleware: Proteksi route `/dashboard/*`, redirect jika tidak ada token.

- [ ] **[F-003] Layouts**
  - Layout Auth (untuk login/register).
  - Layout Dashboard (Sidebar, Header dengan User Profile, Logout button).
  - Sidebar dinamis berdasarkan Role (Catin, Staff, Penghulu, Kepala KUA).

---

## Sprint 2: Publik & Pendaftaran (Core Catin)
Target: User bisa melihat jadwal dan mendaftar nikah.

- [ ] **[F-004] Kalender Ketersediaan (Publik)**
  - Component Calendar (bisa pakai `react-day-picker` atau custom).
  - Fetch API `/simnikah/kalender-ketersediaan`.
  - Warna tanggal: Hijau (Tersedia), Merah (Penuh/Libur).
  - Klik tanggal -> Tampilkan detail slot jam (Modal/Pop-up).

- [ ] **[F-005] Wizard Form Pendaftaran (Step 1: Data Diri)**
  - Form Data Suami & Istri (Nama, NIK, dll).
  - Form Wali Nikah.
  - State management untuk menyimpan data antar-step.

- [ ] **[F-006] Wizard Form Pendaftaran (Step 2: Jadwal & Lokasi)**
  - Pilih Tanggal (DatePicker).
  - Pilih Lokasi (Radio: "Di KUA" / "Di Luar KUA").
  - Validasi Jam: Fetch `/ketersediaan-jam`, disable jam yang penuh.
  - Integrasi Peta (Leaflet) untuk "Di Luar KUA": Pin lokasi, get lat/long.

- [ ] **[F-007] Submit Pendaftaran**
  - Review data sebelum submit.
  - POST ke `/simnikah/pendaftaran`.
  - Handle error API (misal: jadwal tiba-tiba penuh).
  - Success Page -> Redirect ke Dashboard.

---

## Sprint 3: Staff Dashboard (Verifikasi)
Target: Staff bisa memproses pendaftaran masuk.

- [ ] **[F-008] List Pendaftaran (Staff)**
  - Tabel data pendaftaran (`/dashboard/staff/pendaftaran`).
  - Filter: Status, Tanggal.
  - Pagination.
  - Badge status warna-warni.

- [ ] **[F-009] Detail & Verifikasi**
  - Halaman Detail Pendaftaran (`[id]`).
  - Tampilkan data lengkap Catin & Jadwal.
  - Tombol Action: "Verifikasi Formulir" (Modal input catatan).
  - Tombol Action: "Approve" (lanjut ke Kepala KUA).
  - Hit API `/staff/verify-formulir` dan `/staff/approve`.

- [ ] **[F-010] Cetak Pengumuman**
  - Halaman List Pengumuman Nikah.
  - Filter tanggal minggu ini.
  - Button "Print / Export PDF".

---

## Sprint 4: Kepala KUA & Penghulu
Target: Penugasan penghulu dan laporan.

- [ ] **[F-011] Assign Penghulu (Kepala KUA)**
  - List pendaftaran status "Menunggu Penugasan".
  - Modal Assign: Dropdown Penghulu.
  - Validasi: Cek ketersediaan penghulu sebelum submit.
  - Hit API `/kepala-kua/penghulu-tersedia` dan `/assign-penghulu`.

- [ ] **[F-012] Dashboard Statistik (Kepala KUA)**
  - Chart JS / Recharts.
  - Fetch API statistik (`/dashboard/statistik-pernikahan`).
  - Widget: Total Nikah, Performa Penghulu.

- [ ] **[F-013] Dashboard Penghulu**
  - Card "Jadwal Hari Ini".
  - List "Tugas Saya".
  - Halaman detail tugas -> Tombol "Selesaikan Pernikahan" (Input catatan).

---

## Sprint 5: Finishing & UX (Polish)
Target: Aplikasi siap rilis dengan experience yang baik.

- [ ] **[F-014] Notifikasi System**
  - Component Notification Bell di Header.
  - Polling API `/notifikasi/user`.
  - Halaman "Semua Notifikasi".

- [ ] **[F-015] Feedback Catin**
  - Form Feedback muncul di dashboard Catin jika status = "Selesai".
  - Star Rating component.

- [ ] **[F-016] Profile & Settings**
  - Halaman Edit Profile.
  - Ganti Password.

- [ ] **[F-017] Landing Page**
  - Hero section, Features, CTA to Login/Register.
  - Footer.

## Estimasi Waktu (Perkiraan 1 Developer)
- **Sprint 1:** 3-4 Hari
- **Sprint 2:** 5-6 Hari
- **Sprint 3:** 3-4 Hari
- **Sprint 4:** 3-4 Hari
- **Sprint 5:** 2-3 Hari
**Total:** ~3-4 Minggu
