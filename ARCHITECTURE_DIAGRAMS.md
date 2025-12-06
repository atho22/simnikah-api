# 🗺️ SimNikah Architecture & Flow Diagrams

## 1. System Architecture Overview

```
┌────────────────────────────────────────────────────────────────────┐
│                           SIMNIKAH SYSTEM                          │
│                          (Architecture)                             │
└────────────────────────────────────────────────────────────────────┘

┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│  Browser/   │    │  REST API    │    │  Database   │
│  Mobile     │◄──►│  (Go/Gin)    │◄──►│  (MySQL)    │
│  App        │    │  Port 8080   │    │             │
└─────────────┘    └──────────────┘    └─────────────┘
                         │
                    ┌────┴────┐
                    │          │
              ┌─────▼─┐   ┌────▼──────┐
              │ImgBB  │   │Cron Jobs  │
              │(Cloud)│   │(Scheduler)│
              └───────┘   └───────────┘
```

---

## 2. Request Processing Pipeline

```
HTTP Request (Frontend)
        │
        ▼
┌───────────────────────────────────────┐
│   Gin Router                          │
│   (Route matching)                    │
└─────────────┬───────────────────────┘
              │
              ▼
┌───────────────────────────────────────┐
│   CORS Middleware                     │
│   (Check allowed origins)             │
└─────────────┬───────────────────────┘
              │
              ▼
┌───────────────────────────────────────┐
│   Rate Limiter Middleware             │
│   (100 req/min per IP)                │
└─────────────┬───────────────────────┘
              │
              ▼ (if protected route)
┌───────────────────────────────────────┐
│   Auth Middleware                     │
│   - Extract JWT from header           │
│   - Verify signature                  │
│   - Check expiry                      │
│   - Set user_id in context            │
└─────────────┬───────────────────────┘
              │
              ▼ (if role-required)
┌───────────────────────────────────────┐
│   Role Middleware                     │
│   - Get role from claims              │
│   - Verify required role              │
│   - Return 403 if no match            │
└─────────────┬───────────────────────┘
              │
              ▼
┌───────────────────────────────────────┐
│   Handler Function                    │
│   - Get params/body                   │
│   - Validate input                    │
│   - Query database                    │
│   - Return response                   │
└─────────────┬───────────────────────┘
              │
              ▼
HTTP Response (Frontend)
```

---

## 3. Authentication Flow (JWT)

```
┌─────────────────────────────────────────────────────────┐
│                  LOGIN SEQUENCE                         │
└─────────────────────────────────────────────────────────┘

User enters username & password
        │
        ▼
POST /login {username, password}
        │
        ▼
┌─────────────────────────────────────────────────────────┐
│ Handler: auth.Login()                                   │
│ 1. Find user by username                               │
│ 2. Verify password with bcrypt.VerifyPassword()        │
│ 3. Check user status = "Aktif"                         │
│ 4. Create JWT claims:                                  │
│    {                                                    │
│      UserID: "USR1701929380",                          │
│      Email: "user@example.com",                        │
│      Role: "user_biasa",                               │
│      Nama: "Ahmad",                                    │
│      ExpiresAt: now + 24h                              │
│    }                                                    │
│ 5. Sign with secret key: utils.GenerateToken()        │
│ 6. Return token                                        │
└─────────────────────────────────────────────────────────┘
        │
        ▼
Response: {
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "USR1701929380",
    "email": "user@example.com",
    "role": "user_biasa"
  }
}
        │
        ▼
Frontend saves token to localStorage/Cookie
        │
        ▼
┌─────────────────────────────────────────────────────────┐
│           SUBSEQUENT REQUESTS WITH TOKEN                │
└─────────────────────────────────────────────────────────┘

GET /profile
Header: Authorization: Bearer eyJhbGc...
        │
        ▼
middleware.AuthMiddleware()
        │
        ├─ 1. Extract token from header
        │ 2. Parse JWT signature
        │ 3. Verify not expired
        │ 4. Extract claims (user_id, role)
        │ 5. c.Set("user_id", userID)
        │ 6. c.Next() → proceed to handler
        │
        ▼
Handler reads c.Get("user_id")
        │
        ▼
Return user data
```

---

## 4. Marriage Registration Workflow

```
┌────────────────────────────────────────────────────────────────────┐
│                   MARRIAGE REGISTRATION WORKFLOW                   │
└────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────┐
│ STEP 1: USER REGISTRATION           │
├─────────────────────────────────────┤
│ POST /register                      │
│ {username, email, password, nama}   │
│         │                           │
│         ├─ Hash password with bcrypt
│         ├─ Generate user_id (USR...)
│         ├─ Create Users record      │
│         └─ Status: "Aktif"          │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 2: USER LOGIN                  │
├─────────────────────────────────────┤
│ POST /login                         │
│ {username, password}                │
│         │                           │
│         ├─ Verify password          │
│         ├─ Generate JWT token       │
│         └─ Return token (24h valid) │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 3: FILL REGISTRATION FORM      │
├─────────────────────────────────────┤
│ POST /simnikah/pendaftaran          │
│ Body: 4-step wizard data            │
│   {                                 │
│     calon_laki_laki: {...},         │
│     calon_perempuan: {...},         │
│     wali_nikah: {...},              │
│     lokasi_nikah: {                 │
│       tempat_nikah,                 │
│       tanggal_nikah,                │
│       waktu_nikah                   │
│     }                               │
│   }                                 │
│         │                           │
│         ├─ VALIDATE:                │
│         │  • Age checks             │
│         │  • Holiday detection      │
│         │  • Capacity check (max 9) │
│         │  • Time slot availability │
│         │  • Wali validation        │
│         │                           │
│         ├─ CREATE:                  │
│         │  • PendaftaranNikah       │
│         │    (status: Draft)        │
│         │  • CalonPasangan (2x)     │
│         │  • WaliNikah              │
│         │                           │
│         └─ SEND NOTIFICATION        │
│            to staff                 │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 4: STAFF VERIFICATION          │
├─────────────────────────────────────┤
│ [Staff Dashboard]                   │
│                                     │
│ POST /simnikah/staff/verify-.../:id │
│ {status, catatan}                   │
│         │                           │
│         └─ Update status → Verified │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 5: STAFF APPROVAL              │
├─────────────────────────────────────┤
│ POST /simnikah/staff/approve/:id    │
│ {status, catatan}                   │
│         │                           │
│         ├─ Validate documents       │
│         ├─ Check data completeness  │
│         │                           │
│         └─ Status: "Disetujui"      │
│            (ready for assignment)   │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 6: PENGHULU ASSIGNMENT         │
├─────────────────────────────────────┤
│ [Kepala KUA Dashboard]              │
│                                     │
│ POST /simnikah/pendaftaran/:id/...  │
│   /assign-penghulu                  │
│ {penghulu_id, catatan}              │
│         │                           │
│         ├─ CHECK AVAILABILITY:      │
│         │  • Penghulu schedule      │
│         │  • Workload today         │
│         │  • Conflict check         │
│         │                           │
│         ├─ ASSIGN:                  │
│         │  • Set penghulu_id        │
│         │  • Set penghulu_assigned_by
│         │  • Set penghulu_assigned_at
│         │                           │
│         ├─ CREATE NOTIFICATION      │
│         │  to penghulu              │
│         │                           │
│         └─ Status: "Penghulu..."    │
│            (awaiting execution)     │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 7: MARRIAGE EXECUTION          │
├─────────────────────────────────────┤
│ [Penghulu Dashboard]                │
│                                     │
│ POST /simnikah/penghulu/...         │
│   /complete-marriage/:id            │
│ {catatan}                           │
│         │                           │
│         ├─ VERIFY:                  │
│         │  • Documents check        │
│         │  • Parties present        │
│         │  • Guardian present       │
│         │                           │
│         ├─ COMPLETE:                │
│         │  • Record in book         │
│         │  • Issue certificate      │
│         │                           │
│         ├─ UPDATE DB:               │
│         │  • Status: "Selesai"      │
│         │  • Record timestamp       │
│         │                           │
│         ├─ CREATE NOTIFICATION      │
│         │  to catin (marriage done) │
│         │                           │
│         └─ Update penghulu count    │
│            & rating ready           │
└──────────┬────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│ STEP 8: FEEDBACK & RATING           │
├─────────────────────────────────────┤
│ [Catin Dashboard]                   │
│                                     │
│ POST /simnikah/feedback-pernikahan  │
│ {rating, judul, pesan}              │
│         │                           │
│         ├─ Create FeedbackPernikahan
│         │  record                   │
│         │                           │
│         └─ Available for Kepala KUA │
│            to review                │
└─────────────────────────────────────┘

Final Status: "Selesai" ✓
Marriage completed successfully!
```

---

## 5. Database Query Flow

```
Handler receives request
        │
        ▼
┌──────────────────────────────────────┐
│  Build GORM Query                    │
│  h.DB.Where(condition, args...)      │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  Add Conditions (if multiple)        │
│  .Where(...).Where(...)              │
│  .Order("field DESC")                │
│  .Limit(10).Offset(0)                │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  Execute Query (.First, .Find, etc) │
│  GORM converts to SQL                │
│  SQL with ? placeholders             │
│  (prevents SQL injection)            │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  MySQL Driver                        │
│  Binds parameters to ?               │
│  Executes SQL safely                 │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  MySQL Database                      │
│  Retrieves data from table           │
│  Returns result set                  │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  GORM Maps Rows                      │
│  Scans columns to struct fields      │
│  Type conversion                     │
│  Return populated struct(s)          │
└──────┬───────────────────────────────┘
       │
       ▼
Handler uses data
    └─ Return JSON response
```

---

## 6. Notification System

```
┌────────────────────────────────────────────────────────┐
│              NOTIFICATION SYSTEM                       │
└────────────────────────────────────────────────────────┘

Type 1: IMMEDIATE NOTIFICATIONS (on action)
        │
        ├─ User register → create Users
        │  └─ Send "Akun Berhasil Dibuat"
        │
        ├─ Marriage registered → create PendaftaranNikah
        │  └─ Send "Pendaftaran Diterima" to catin & staff
        │
        ├─ Staff approves → update status
        │  └─ Send "Pendaftaran Disetujui" to catin
        │
        ├─ Penghulu assigned → update penghulu_id
        │  └─ Send "Anda Ditugaskan" to penghulu
        │
        └─ Marriage completed → update status
           └─ Send "Pernikahan Selesai" to catin

Type 2: SCHEDULED NOTIFICATIONS (Cron Job)
        │
        Every day at 08:00
        │
        ▼
        ┌─────────────────────────────────────┐
        │ CronJobService                      │
        │ (services/cron_job.go)              │
        │                                     │
        │ StartReminderCronJobWithSchedule()  │
        │ - Find pending registrations        │
        │ - Find unverified registrations     │
        │ - Send reminder to staff            │
        └─────────────────────────────────────┘
        │
        ▼
        Notifications created for:
        - Staff: "X pending registrations"
        - Staff: "X unverified registrations"
        - Penghulu: "Schedule reminders"

Type 3: MANUAL NOTIFICATIONS (by admin)
        │
        POST /simnikah/notifikasi
        {user_id, judul, pesan, tipe}
        │
        ▼
        Create notification immediately
        (for ad-hoc messages)
```

---

## 7. User Role Access Matrix

```
┌──────────────────────────────────────────────────────────┐
│             ROLE-BASED ACCESS CONTROL MATRIX             │
└──────────────────────────────────────────────────────────┘

Endpoint                          user_biasa  staff  penghulu  kepala_kua
───────────────────────────────────────────────────────────────────────
/register                              ✓       ✓       ✓         ✓
/login                                 ✓       ✓       ✓         ✓
/profile                               ✓       ✓       ✓         ✓
/simnikah/pendaftaran (POST)           ✓       ✓       ✗         ✗
/simnikah/pendaftaran/status           ✓       ✓       ✓         ✓
/simnikah/feedback                     ✓       ✗       ✗         ✗
/simnikah/kalender-ketersediaan        ✓       ✓       ✓         ✓
/simnikah/staff/verify-formulir        ✗       ✓       ✗         ✗
/simnikah/staff/approve                ✗       ✓       ✗         ✗
/simnikah/penghulu/verify-docs         ✗       ✗       ✓         ✗
/simnikah/penghulu/complete-marriage   ✗       ✗       ✓         ✓
/simnikah/pendaftaran/assign-penghulu  ✗       ✗       ✗         ✓
/simnikah/kepala-kua/staff (POST)      ✗       ✗       ✗         ✓
/simnikah/kepala-kua/feedback          ✗       ✗       ✗         ✓
/simnikah/dashboard/*                  ✗       ✓       ✓         ✓
/simnikah/notifikasi/*                 ✓       ✓       ✓         ✓

Legend: ✓ = Allowed, ✗ = Denied, ? = Conditional
```

---

## 8. Data Model Relationships

```
Users (Central Hub)
├── 1:1 → StaffKUA (if role = staff/penghulu/kepala_kua)
├── 1:1 → Penghulu (if role = penghulu)
├── 1:1 → CalonPasangan (if bride/groom)
├── 1:N → PendaftaranNikah (as pendaftar)
├── 1:N → Notifikasi
└── 1:N → FeedbackPernikahan


PendaftaranNikah (Central to Flow)
├── N:1 ← Users (pendaftar_id)
├── N:1 ← Users (calon_suami_id)
├── N:1 ← Users (calon_istri_id)
├── N:1 ← Penghulu (penghulu_id)
├── 1:1 → WaliNikah
└── 1:N ← FeedbackPernikahan

        Example:
        Pendaftaran 1 has:
        - pendaftar: Ahmad (submitted)
        - calon_suami: Ahmad
        - calon_istri: Siti
        - wali_nikah: Abdullah (father of Siti)
        - penghulu: Haji Budi (assigned officer)
```

---

## 9. Error Handling Flow

```
Request → Handler
    │
    ├─ Validation Error
    │  └─ 400 Bad Request
    │     └─ {success: false, message: "...", error: "..."}
    │
    ├─ Authentication Error
    │  └─ 401 Unauthorized
    │     └─ {success: false, message: "Invalid token"}
    │
    ├─ Authorization Error
    │  └─ 403 Forbidden
    │     └─ {success: false, message: "Role tidak sesuai"}
    │
    ├─ Resource Not Found
    │  └─ 404 Not Found
    │     └─ {success: false, message: "Data tidak ditemukan"}
    │
    ├─ Database Error
    │  └─ 500 Internal Server Error
    │     └─ {success: false, message: "Gagal query database"}
    │
    ├─ Business Logic Error
    │  └─ 400/409 Conflict
    │     └─ {success: false, message: "Jadwal penuh", type: "schedule_conflict"}
    │
    └─ Success
       └─ 200/201 OK/Created
          └─ {success: true, message: "...", data: {...}}
```

---

## 10. Caching Strategy

```
┌─────────────────────────────────────────────────────┐
│           CACHING CONSIDERATIONS                    │
└─────────────────────────────────────────────────────┘

Current (No Explicit Cache Layer):
├─ Database indexes (fast queries)
├─ Connection pooling (reuse connections)
└─ In-memory constants:
   ├─ HariLiburNasional map
   ├─ TimeSlots array
   ├─ ValidHubunganWali array
   └─ User roles/permissions

Potential Future Caching:
├─ Redis for session tokens
├─ Cache user profile data
├─ Cache calendar availability
└─ Cache statistics (dashboard data)
```

---

**These diagrams illustrate the complete flow and architecture of the SimNikah system.**

*For detailed code implementation, refer to the respective handler files.*
