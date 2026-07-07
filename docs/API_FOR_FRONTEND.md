**Dokumentasi API untuk Frontend**

Base URL
- Gunakan base URL sesuai environment: `https://{API_HOST}` (contoh: `https://api.example.com`)

Header Umum
- `Authorization: Bearer <JWT>` — diperlukan untuk endpoint yang memerlukan autentikasi.
- `Content-Type: application/json`

Format error
- Respon error umum:
  - HTTP status non-2xx
  - Body:
    {
      "success": false,
      "message": "...",
      "error": "detail teknis (ringkas)"
    }

-------------------------
1) Auth

- POST /register
  - Body: {"email":"...","password":"...","nama":"...","role":"pengguna"}
  - Response: {"success": true, "user": {...}}

- POST /login
  - Body: {"email":"...","password":"..."}
  - Response: {"success": true, "token": "<JWT>", "user": {...}}

- GET /profile
  - Auth required
  - Response: {"success": true, "data": {...}}

- POST /upload-photo
  - Auth required, multipart/form-data, field `photo`
  - Response: {"success": true, "url": "..."}

-------------------------
2) Pendaftaran (Catin)

- POST /simnikah/check-schedule
  - Auth required
  - Body: {"tanggal_nikah":"YYYY-MM-DD","waktu_nikah":"HH:MM","tempat_nikah":"di_kua|di_luar","alamat_nikah":"...","latitude":...,"longitude":...}
  - Response: `ScheduleCheckResult` (lihat PRD) — kembalian berisi `available`, `recommended_penghulu`, `alternatives`.

- POST /simnikah/pendaftaran
  - Auth required
  - Body: pendaftaran lengkap (nama, tanggal, waktu, tempat, alamat, pasangan, dokumen)
  - Response: created pendaftaran detail

- GET /simnikah/pendaftaran/status
  - Auth required
  - Query: ?user_id=... (opsional)
  - Response: status ringkasan

- GET /simnikah/pendaftaran/:id
  - Auth required
  - Response: detail pendaftaran

- GET /simnikah/pendaftaran
  - Auth required, role `staff` atau `kepala_kua`
  - Query params: paging, filter status

- PUT /simnikah/pendaftaran/:id/update-status
  - Auth required, role `staff|penghulu|kepala_kua`
  - Body: {"status_pendaftaran":"...","catatan":"..."}

-------------------------
3) Forward Chaining (Kepala KUA)

- GET /simnikah/kepala-kua/forward-chaining/recommendation/:id
  - Auth required, role `kepala_kua`
  - Response: `AssignmentRecommendation` — recommended_penghulu, alternatives, evaluation_process

- GET /simnikah/kepala-kua/forward-chaining/evaluation/:id
  - Auth required, role `kepala_kua`
  - Response: trace evaluasi per penghulu (rule results, derived facts)

- GET /simnikah/kepala-kua/forward-chaining/config
  - Auth required, role `kepala_kua`
  - Response: konfigurasi FC (minimum_rating, capacity, weights)

- POST /simnikah/kepala-kua/forward-chaining/assign/:id
  - Auth required, role `kepala_kua`
  - Body: {"penghulu_id": <id>, "notes": "..."}
  - Response: {"success": true, "assignment_id": ...}

-------------------------
4) Penghulu

- GET /simnikah/penghulu/jadwal-penugasan
  - Auth required, role `penghulu|kepala_kua`
  - Response: list jadwal penugasan

- PUT /simnikah/penghulu/coordinates
  - Auth required, role `penghulu`
  - Body: {"latitude":..., "longitude":...}

-------------------------
5) Kepala KUA management

- GET /simnikah/kepala-kua/available-penghulu
  - Auth required, role `kepala_kua`

- GET /simnikah/kepala-kua/penghulu-tersedia
  - Auth required, role `kepala_kua`

-------------------------
6) Staff

- GET /simnikah/staff
  - Auth required, role `kepala_kua`

- PUT /simnikah/staff/:id
  - Auth required, role `kepala_kua`

- POST /simnikah/staff/pendaftaran
  - Auth required, role `staff|kepala_kua`

-------------------------
7) Lokasi & Geocoding (penting untuk FE peta)

- POST /simnikah/location/geocode
  - Auth required
  - Body: {"address":"..."}
  - Response: {"latitude":..., "longitude":..., "display_name":"..."}

- POST /simnikah/location/reverse-geocode
  - Auth required
  - Body: {"latitude":..., "longitude":...}
  - Response: {"road":"...","kelurahan":"...","kecamatan":"...","kota":"...","provinsi":"..."}
  - NOTE: FE harus memvalidasi `kecamatan == "Banjarmasin Utara"` sebelum menyimpan lokasi.

- GET /simnikah/location/search?query=...
  - Auth required
  - Response: array hasil pencarian alamat dengan koordinat

- PUT /simnikah/pendaftaran/:id/location
  - Auth required
  - Body: {"latitude":...,"longitude":...,"alamat":"...","jalan":"...","kelurahan":"...","kecamatan":"...","kota":"...","provinsi":"..."}

- GET /simnikah/pendaftaran/:id/location
  - Auth required
  - Response: lokasi yang tersimpan untuk pendaftaran

Catatan untuk FE (peta):
- Semua reverse-geocoding harus divalidasi berada di Kecamatan Banjarmasin Utara (lihat PRD peta). Jika diluar, tampilkan pesan: "Lokasi berada di luar wilayah pelayanan KUA Kecamatan Banjarmasin Utara." dan batalkan penyimpanan.
- Semua panggilan OSRM (routing) hanya boleh dipakai untuk titik dalam wilayah Kecamatan Banjarmasin Utara.

-------------------------
8) Notifikasi

- POST /simnikah/notifikasi
  - Auth required, role `staff|kepala_kua`

- GET /simnikah/notifikasi/user/me
  - Auth required

- GET /simnikah/notifikasi/:id
  - Auth required

- PUT /simnikah/notifikasi/:id/status
  - Auth required

- PUT /simnikah/notifikasi/mark-all-read
  - Auth required

- DELETE /simnikah/notifikasi/:id
  - Auth required

- POST /simnikah/notifikasi/send-to-role
  - Auth required, role `staff|kepala_kua`

- POST /simnikah/notifikasi/run-reminder
  - Auth required, role `staff|kepala_kua` (manual trigger)

-------------------------
9) Dashboard

- GET /simnikah/dashboard/kepala-kua
  - Auth required, role `kepala_kua`

- GET /simnikah/dashboard/staff
  - Auth required, role `staff`

- GET /simnikah/dashboard/statistik-pernikahan
  - Auth required, role `staff|kepala_kua`

- GET /simnikah/dashboard/penghulu-performance
  - Auth required, role `staff|kepala_kua`

- GET /simnikah/dashboard/peak-hours
  - Auth required, role `staff|kepala_kua`

-------------------------
10) Health

- GET /health
  - No auth
  - Response: {"status":"healthy","service":"SimNikah API","timestamp":"..."}

-------------------------
Petunjuk integrasi cepat untuk FE
- Sertakan header `Authorization: Bearer <token>` untuk semua endpoint yang memerlukan auth.
- Untuk peta:
  - Ambil polygon Kecamatan Banjarmasin Utara (jika tersedia dari backend) atau sediakan GeoJSON statis dari FE assets.
  - Lakukan reverse-geocode via `POST /simnikah/location/reverse-geocode` untuk mendapatkan `kecamatan` dan validasi.
  - Gunakan OSRM endpoint (backend memanggil OSRM) — FE hanya memanggil backend yang mem-proxy request routing.

File ini adalah ringkasan endpoint utama. Untuk detail request/response penuh (field model, contoh payload lebih lengkap), saya bisa tambahkan koleksi Postman / OpenAPI (YAML) jika Anda mau.
