# 📖 SimNikah Project - Executive Summary

## 🎯 Project Overview

**SimNikah** adalah sebuah aplikasi web dan API berbasis Go untuk mengelola pendaftaran pernikahan di Kantor Urusan Agama (KUA). Sistem ini mendigitalkan proses pendaftaran yang sebelumnya manual.

```
┌─────────────────────────────────────────────────────────────────┐
│                         SIMNIKAH SYSTEM                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Frontend (Next.js/React)     Backend (Go + Gin)   Database      │
│  ┌──────────────────┐        ┌──────────────────┐  ┌─────────┐  │
│  │ User Dashboard   │        │ REST API         │  │ MySQL   │  │
│  │ Registration     │◄──────►│ - Auth           │  │ 8 Tables│  │
│  │ Form             │        │ - Registrations  │  │ Indexes │  │
│  │ Status Tracking  │        │ - Notifications  │  │ Queries │  │
│  │                  │        │ - Dashboard      │  │         │  │
│  └──────────────────┘        │ - Analytics      │  └─────────┘  │
│                              └──────────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 👥 User Roles & Responsibilities

### 1. **Calon Pengantin (User Biasa)**
- Register akun
- Isi form pendaftaran nikah (4-step wizard)
- Cek status pendaftaran
- Lihat feedback setelah nikah selesai
- Lihat kalender ketersediaan

### 2. **Staff KUA**
- Verifikasi kelengkapan data pendaftaran
- Approve atau tolak pendaftaran
- Lihat daftar pendaftaran dengan filter
- Generate laporan pengumuman nikah
- View dashboard dengan pending tasks

### 3. **Penghulu (Marriage Officer)**
- Lihat jadwal nikah yang ditugaskan
- Verifikasi dokumen calon pengantin
- Melaksanakan pernikahan
- Mark pernikahan sebagai selesai
- Lihat detail lokasi nikah (dengan peta jika di luar KUA)

### 4. **Kepala KUA (Head of Office)**
- Assign penghulu ke pendaftaran
- Monitor statistik pernikahan
- Lihat performa penghulu (rating, jumlah nikah)
- Manage staff dan penghulu
- Manage feedback dari calon pengantin
- Dashboard dengan analytics lengkap

---

## 🔄 Marriage Registration Flow

```
CATIN SUBMISSION
┌─────────────────────────────────────────────────┐
│ 1. Register Akun                                │
│    → POST /register                             │
│    → Dapat User ID & Password                   │
└─────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│ 2. Login                                        │
│    → POST /login                                │
│    → Dapat JWT Token                            │
└─────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│ 3. Isi Form Pendaftaran (Status: DRAFT)         │
│    → POST /simnikah/pendaftaran                 │
│    → Pilih jadwal, lokasi, data calon           │
│    → Generate nomor_pendaftaran                 │
└─────────────────────────────────────────────────┘
                    │
        STAFF VERIFICATION PROCESS
                    ▼
┌─────────────────────────────────────────────────┐
│ 4. Verifikasi Formulir (Status: VERIFIKASI)    │
│    → POST /simnikah/staff/verify-formulir/:id   │
│    → Staff check kelengkapan data               │
└─────────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│ 5. Approve Pendaftaran (Status: DISETUJUI)     │
│    → POST /simnikah/staff/approve/:id           │
│    → Lanjut ke penugasan penghulu               │
└─────────────────────────────────────────────────┘
                    │
        KEPALA KUA ASSIGNMENT
                    ▼
┌─────────────────────────────────────────────────┐
│ 6. Assign Penghulu (Status: DITUGASKAN)        │
│    → POST /simnikah/pendaftaran/:id/assign-...  │
│    → Kepala KUA pilih penghulu tersedia         │
│    → Notifikasi ke penghulu                     │
└─────────────────────────────────────────────────┘
                    │
        PENGHULU EXECUTION
                    ▼
┌─────────────────────────────────────────────────┐
│ 7. Melaksanakan Nikah (Status: SELESAI)        │
│    → Penghulu perform nikah ceremony            │
│    → Catat di buku nikah                        │
│    → POST /simnikah/penghulu/complete-marriage  │
└─────────────────────────────────────────────────┘
                    │
        POST-MARRIAGE
                    ▼
┌─────────────────────────────────────────────────┐
│ 8. Feedback & Rating                            │
│    → POST /simnikah/feedback-pernikahan         │
│    → Catin kasih rating 1-5                     │
│    → Kepala KUA bisa view feedback              │
└─────────────────────────────────────────────────┘

Alternative: DITOLAK
    ↓
Catin bisa daftar ulang
```

---

## 📊 Database Schema Overview

```
┌──────────────┐
│    Users     │  Main user account table
├──────────────┤
│ id (PK)      │
│ user_id (UQ) │  USR1701929380 format
│ username (UQ)│
│ email (UQ)   │
│ password     │  bcrypt hashed
│ role         │  user_biasa, staff, penghulu, kepala_kua
│ status       │  Aktif, Nonaktif, Blokir
│ nama         │  Full name
│ profile_photo│  ImgBB URL
│ created_at   │
│ updated_at   │
└──────────────┘
       │
       ├─────────────────┬──────────────┬──────────────┐
       ▼                 ▼              ▼              ▼
┌─────────────┐  ┌───────────┐  ┌──────────┐  ┌─────────────┐
│ StaffKUA    │  │ Penghulu  │  │CalonPsngan│ │Notifikasi   │
├─────────────┤  ├───────────┤  ├──────────┤  ├─────────────┤
│ user_id (FK)│  │user_id(FK)│  │user_id(FK)│ │user_id  (FK) │
│ nip         │  │ nip       │  │ nik      │ │judul        │
│ jabatan     │  │nama_lngkap│  │nama_lngkap│ │pesan        │
│ status      │  │rating     │  │ttl_lahir │ │tipe         │
│ ...         │  │...        │  │...       │ │status_baca  │
└─────────────┘  └───────────┘  └──────────┘  └─────────────┘

┌─────────────────────────┐
│ PendaftaranNikah        │  Marriage registration
├─────────────────────────┤
│ id (PK)                 │
│ nomor_pendaftaran (UQ)  │  NIKAH-20240105-4532
│ pendaftar_id (FK)       │  Who submitted (calon)
│ calon_suami_id (FK)     │  Groom ID
│ calon_istri_id (FK)     │  Bride ID
│ wali_nikah_id (FK)      │  Guardian for bride
│ penghulu_id (FK)        │  Assigned marriage officer
│ tanggal_nikah           │
│ waktu_nikah             │
│ tempat_nikah            │  "Di KUA" atau "Di Luar KUA"
│ alamat_akad             │  Wedding address if outside
│ latitude, longitude     │  GPS coordinates
│ status_pendaftaran      │  Draft → Disetujui → ... → Selesai
│ created_at, updated_at  │
└─────────────────────────┘
         │
         ├──────────────┬─────────────────┐
         ▼              ▼                 ▼
    ┌─────────┐  ┌──────────────┐  ┌──────────────┐
    │WaliNikah│  │FeedbackPrnikh│  │... more data │
    └─────────┘  └──────────────┘  └──────────────┘
```

---

## 🔐 Authentication & Security

```
┌─────────────────────────────────────────────────┐
│           JWT Token Authentication              │
├─────────────────────────────────────────────────┤
│                                                 │
│ 1. User Login                                   │
│    POST /login {username, password}             │
│         │                                       │
│         ├─ Verify password with bcrypt         │
│         ├─ Generate JWT with claims:           │
│         │  - user_id                           │
│         │  - email                             │
│         │  - role                              │
│         │  - nama                              │
│         │  - expiresAt: now + 24h              │
│         └─ Return token                        │
│                                                 │
│ 2. Token Usage                                  │
│    GET /protected-endpoint                      │
│    Header: Authorization: Bearer <token>        │
│         │                                       │
│         ├─ Extract token from header           │
│         ├─ Verify JWT signature                │
│         ├─ Check expiry time                   │
│         ├─ Set context values:                 │
│         │  - c.Set("user_id", claims.UserID)   │
│         │  - c.Set("role", claims.Role)        │
│         └─ Allow request to proceed            │
│                                                 │
│ 3. Authorization Checks                        │
│    Per-route role validation:                  │
│    - RoleMiddleware("staff")                   │
│    - MultiRoleMiddleware("staff", "kepala_kua")│
│                                                 │
└─────────────────────────────────────────────────┘

Encryption:
├─ Passwords: bcrypt (salted hash)
├─ Images: Stored on ImgBB cloud
├─ Database: Plain text (use SSL for connections)
└─ API: CORS restricted to allowed origins
```

---

## 🛣️ API Endpoint Categories

```
🔓 PUBLIC ENDPOINTS (No Auth Required)
├─ GET  /health                           Health check
├─ POST /register                         Create account
├─ POST /login                            Login
├─ GET  /simnikah/kalender-ketersediaan   View calendar
└─ GET  /simnikah/ketersediaan-jam        View time slots

🔒 PROTECTED ENDPOINTS (Auth Required)
├─ GET  /profile                          Get profile
├─ POST /simnikah/pendaftaran             Register marriage
├─ GET  /simnikah/pendaftaran/status      Check status
└─ ... 40+ more endpoints

👮 ROLE-RESTRICTED ENDPOINTS
├─ STAFF ONLY
│  ├─ POST /simnikah/staff/verify-formulir/:id
│  └─ POST /simnikah/staff/approve/:id
│
├─ PENGHULU ONLY
│  ├─ POST /simnikah/penghulu/verify-documents/:id
│  └─ POST /simnikah/penghulu/complete-marriage/:id
│
└─ KEPALA KUA ONLY
   ├─ POST /simnikah/pendaftaran/:id/assign-penghulu
   └─ GET  /simnikah/kepala-kua/feedback
```

---

## 📊 Key Features Comparison

| Feature | Catin | Staff | Penghulu | Kepala KUA |
|---------|-------|-------|----------|-----------|
| Register nikah | ✅ | ✅ | ❌ | ❌ |
| Check status | ✅ | ✅ | ✅ | ✅ |
| Verify form | ❌ | ✅ | ❌ | ❌ |
| Approve | ❌ | ✅ | ❌ | ❌ |
| Assign penghulu | ❌ | ❌ | ❌ | ✅ |
| Execute marriage | ❌ | ❌ | ✅ | ❌ |
| View dashboard | ❌ | ✅ | ✅ | ✅ |
| View feedback | ❌ | ❌ | ❌ | ✅ |
| View analytics | ❌ | ❌ | ❌ | ✅ |
| Get notifications | ✅ | ✅ | ✅ | ✅ |

---

## 🚀 Quick Start Commands

```bash
# 1. Setup environment
export DB_HOST=127.0.0.1
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=simnikah
export GIN_MODE=debug

# 2. Run application
make dev
# or
go run cmd/api/main.go

# 3. Test with curl
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456","nama":"Test User"}'

# 4. Build for production
make build
./bin/simnikah-api

# 5. Run tests
make test
make coverage
```

---

## 📁 Important Files to Know

```
Project Structure:
├── cmd/api/main.go                    ← START HERE! Entry point
│   └── Initialize DB, routes, handlers
│
├── internal/
│   ├── handlers/
│   │   ├── auth/auth.go              ← Authentication logic
│   │   ├── catin/daftar.go           ← Marriage registration
│   │   ├── staff/staff.go            ← Staff verification
│   │   ├── penghulu/penghulu.go      ← Officer operations
│   │   ├── kepala_kua/kepala_kua.go  ← Head of office
│   │   ├── dashboard/dashboard.go    ← Dashboards
│   │   └── notification/notification.go ← Notifications
│   │
│   ├── models/
│   │   ├── models.go                 ← Database schemas
│   │   └── constants.go              ← All constants
│   │
│   ├── middleware/
│   │   ├── auth.go                   ← JWT verification
│   │   └── rate_limit.go             ← Rate limiting
│   │
│   ├── services/
│   │   ├── cron_job.go               ← Scheduled tasks
│   │   └── notification_service.go   ← Notification logic
│   │
│   └── seeders/                       ← Initial data
│
├── pkg/
│   ├── crypto/bcrypt.go              ← Password hashing
│   ├── storage/imgbb.go              ← Image upload
│   ├── utils/
│   │   ├── jwt.go                    ← Token generation
│   │   ├── string_utils.go
│   │   └── date_utils.go
│   └── validator/marriage_validation.go
│
├── config/
│   ├── config.go                     ← DB connection
│   ├── db.go
│   └── indexes.go                    ← Database indexes
│
├── migrations/
│   └── init.sql                      ← Database schema
│
└── docs/                              ← Documentation
    ├── API_DOCUMENTATION_COMPLETE.md
    ├── FRONTEND_IMPLEMENTATION_GUIDE.md
    └── ... many more guides
```

---

## 🔗 Key Connections

```
Frontend Request Flow:
┌──────────────┐
│ React App    │
└──────┬───────┘
       │
       │ Fetch with JWT
       │
       ▼
┌──────────────────────┐
│ Gin Router (main.go) │
│ - CORS Middleware    │
│ - Auth Middleware    │
│ - Role Middleware    │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ Handlers             │
│ - auth.go            │
│ - daftar.go          │
│ - staff.go           │
│ - penghulu.go        │
│ - kepala_kua.go      │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ GORM ORM             │
│ - Build queries      │
│ - Execute with ?     │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ MySQL Database       │
│ 8 Tables            │
│ with Indexes        │
└──────────────────────┘
```

---

## ✨ Notable Implementation Details

1. **Marriage Registration Validation**
   - Age check (19+ for groom, 16+ for bride)
   - Holiday detection (Sundays + national holidays)
   - Capacity check (max 9 marriages/day at KUA)
   - Time slot availability (08:00-16:00, hourly slots)

2. **Image Storage**
   - Uses ImgBB (free cloud storage)
   - Stores URL in `Users.profile_photo`
   - Optional API key (anonymous upload if not provided)

3. **Notification System**
   - Cron job runs daily at 8:00 AM
   - Automatically reminds about pending registrations
   - Can be manually triggered via API

4. **Database Indexes**
   - Optimized queries for common operations
   - 5-10x faster than without indexes

5. **Rate Limiting**
   - Global: 100 req/min per IP
   - Auth endpoints: Stricter limits
   - Uses `ulule/limiter` package

---

## 📞 Documentation Files

| File | Purpose |
|------|---------|
| `PROJECT_ANALYSIS.md` | Complete project overview (THIS FILE) |
| `LEARNING_GUIDE.md` | Study guide with quick references |
| `DEVELOPER_REFERENCE.md` | How-to guide for common tasks |
| `/docs/FULL_API_USAGE.md` | All endpoints with examples |
| `/docs/FRONTEND_IMPLEMENTATION_GUIDE.md` | Frontend developer guide |
| `/docs/API_REQUEST_BODY_DOCUMENTATION.md` | Request/response formats |

---

## 🎯 Summary

**SimNikah** adalah sistem lengkap untuk mengelola pendaftaran pernikahan di KUA dengan:
- ✅ Modern Go backend dengan REST API
- ✅ JWT-based authentication & role-based authorization
- ✅ Complete marriage registration workflow
- ✅ Staff verification & approval process
- ✅ Officer assignment & scheduling
- ✅ Real-time notifications
- ✅ Analytics & reporting
- ✅ Responsive design for all user roles
- ✅ Production-ready with security & performance

**Current Status:** ✅ Fully functional and ready for deployment

**Tech Stack:** Go 1.23 + Gin + MySQL + JWT + ImgBB

---

**Created:** December 6, 2025  
**Updated:** Today  
**Author:** Project Analysis Team
