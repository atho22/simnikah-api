Sistem Rekomendasi Keputusan Penentuan Penghulu — PRD v2.0
=========================================================

1. Tujuan Produk
-----------------
Backend bertugas menyediakan REST API untuk membantu Kepala KUA Kecamatan Banjarmasin Utara dalam menentukan penghulu berdasarkan aturan (rule) menggunakan metode Forward Chaining.

Backend bertanggung jawab untuk:

- Autentikasi pengguna
- Manajemen data penghulu
- Manajemen jadwal akad
- Manajemen rule Forward Chaining
- Proses inferensi
- Integrasi OSRM
- Riwayat rekomendasi
- Dashboard

2. User Role
------------
Admin

Hak akses

- ✅ Login
- ✅ Kelola User
- ✅ Kelola Penghulu
- ✅ Kelola Rule
- ✅ Kelola Jadwal
- ✅ Kelola Lokasi
- ✅ Melihat Dashboard

Kepala KUA

Hak akses

- ✅ Login
- ✅ Membuat permintaan rekomendasi
- ✅ Melihat seluruh kandidat
- ✅ Melihat alasan rekomendasi
- ✅ Menetapkan penghulu
- ✅ Melihat histori

3. Modul Backend
-----------------
Authentication

POST /api/v1/auth/login

POST /api/v1/auth/logout

POST /api/v1/auth/refresh

GET /api/v1/auth/profile

User
GET /api/v1/users
POST /api/v1/users
PUT /api/v1/users/{id}
DELETE /api/v1/users/{id}

Penghulu
GET /api/v1/penghulu
POST /api/v1/penghulu
PUT /api/v1/penghulu/{id}
DELETE /api/v1/penghulu/{id}

Field

- id
- nip
- nama
- status
- latitude
- longitude
- created_at
- updated_at

Jadwal
GET /api/v1/jadwal
POST /api/v1/jadwal
PUT /api/v1/jadwal/{id}
DELETE /api/v1/jadwal/{id}

Lokasi Akad
GET /api/v1/lokasi
POST /api/v1/lokasi
PUT /api/v1/lokasi/{id}
DELETE /api/v1/lokasi/{id}

4. Rule Management ⭐
--------------------
Karena penelitian menggunakan Forward Chaining. Harus ada.

Rule
GET /api/v1/rules

GET /api/v1/rules/{id}

POST /api/v1/rules

PUT /api/v1/rules/{id}

DELETE /api/v1/rules/{id}

Contoh Rule

R001

IF
Jadwal Bentrok

THEN

Tidak Direkomendasikan

R002

IF

Jadwal Tersedia

AND

Beban Rendah

THEN

Prioritas Tinggi

R003

IF

Prioritas Tinggi

AND

Jarak Dekat

THEN

Sangat Direkomendasikan

5. Recommendation Engine
-------------------------
Endpoint

POST /api/v1/recommendations

Input

- tanggal
- jam
- latitude
- longitude

Output

{
  "recommendation": [
    {
      "penghulu": "Ahmad",
      "status": "Sangat Direkomendasikan",
      "matched_rules": [
        "R002",
        "R003",
        "R006"
      ],
      "reason": [
        "Jadwal tersedia",
        "Beban kerja rendah",
        "Lokasi dekat"
      ]
    }
  ]
}

⚠ Tidak ada `score`, `confidence`, atau `ranking_value` karena Forward Chaining tidak menggunakan pembobotan.

6. Recommendation Simulation ⭐
-----------------------------
Untuk testing.

POST /api/v1/recommendations/simulate

Input

- tanggal
- jam
- latitude
- longitude

Output

Semua kandidat beserta Rule yang terpenuhi, Rule yang gagal, dan Status

7. Assignment
-------------
POST /api/v1/assignments

Body

- recommendation_id
- penghulu_id

Jika penghulu sudah ditetapkan -> Return `409 Conflict`

8. Evaluation
--------------
POST /api/v1/evaluations

Output

- Matched Rule
- Failed Rule
- Final Status
- Reason

9. Dashboard
------------
GET /api/v1/dashboard

Widget

- Jumlah Penghulu
- Jadwal
- Akad
- Penugasan

Workload
GET /api/v1/dashboard/workload

Output

Penghulu A — 10 tugas
Penghulu B — 8 tugas
Penghulu C — 5 tugas

Assignment History
GET /api/v1/dashboard/history

10. Audit Log ⭐
----------------
GET /api/v1/audit-log

Isi

- User
- Action
- Timestamp
- IP
- Device

11. Notification
-----------------
GET /api/v1/notifications
PUT /api/v1/notifications/{id}/read

12. Calendar
-------------
GET /api/v1/penghulu/{id}/calendar

Output contoh:

08:00 — Akad
10:00 — Akad
13:00 — Kosong

13. OSRM
-----------
GET /api/v1/osrm/route

Output

- Distance
- Duration

Data ini menjadi salah satu fakta dalam proses Forward Chaining.

14. Logging
-----------
Disimpan untuk:

- Login
- Logout
- CRUD
- Rule Evaluation
- Assignment
- Recommendation
- Error

15. Non Functional Requirement
-----------------------------
- JWT Authentication
- REST API
- MySQL
- GORM
- Gin
- Soft Delete
- Audit Trail
- Response < 500 ms untuk operasi CRUD
- Response rekomendasi < 3 detik
- JSON Response
- Structured Logging
- Rate Limiter
- Environment Configuration
- API Versioning (`/api/v1`)

---

File ini adalah versi ringkas PRD v2.0 untuk `Sistem Rekomendasi Keputusan Penentuan Penghulu`.
