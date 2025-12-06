# 🎓 SimNikah Project Study Guide

## Quick Reference - Komponen Kunci

### 1️⃣ Entry Point
```
cmd/api/main.go
├── ConnectDB()             → config/config.go
├── AutoMigrate()           → internal/models/models.go
├── Setup Handlers          → internal/handlers/*/
├── Setup Routes            → cmd/api/main.go (line 109+)
├── CORS Config             → github.com/gin-contrib/cors
├── Rate Limiting           → internal/middleware/rate_limit.go
└── Start Server on :8080
```

### 2️⃣ Authentication Flow
```
Frontend (Browser)
    │
    ├─→ POST /register {username, email, password, nama}
    │     ↓
    │   auth.Register() → Hash password → Create Users record
    │     ↓
    │   Return: {success, user_id, message}
    │
    ├─→ POST /login {username, password}
    │     ↓
    │   auth.Login() → Verify password → Generate JWT
    │     ↓
    │   Return: {success, token, user}
    │
    └─→ GET /profile (with Authorization header)
          ↓
        auth.GetProfile() → Extract user_id from token → Return user data
```

### 3️⃣ Dependency Injection Pattern
```go
// Handler declaration
type InDB struct {
    DB *gorm.DB
}

// Usage in handler
func (h *InDB) GetProfile(c *gin.Context) {
    h.DB.Where("user_id = ?", userID).First(&user)
}

// Initialization in main.go
authHandler := &auth.InDB{DB: DB}
r.GET("/profile", authHandler.GetProfile)
```

### 4️⃣ Middleware Chain
```
Request
  ↓
CORS Middleware
  ↓
Rate Limiter
  ↓
Auth Middleware (jika protected route)
  ├─ Verify JWT token
  ├─ Extract user_id
  └─ Set in context: c.Set("user_id", userID)
  ↓
Role Middleware (jika role-protected)
  ├─ Check user.Role
  └─ Return 403 jika tidak sesuai
  ↓
Handler Function
  ↓
Response
```

### 5️⃣ Marriage Registration Flow
```
User (Catin)
  ↓
POST /simnikah/pendaftaran
  ├─ Validate form data (4-step wizard data)
  ├─ Check jadwal availability
  ├─ Generate nomor_pendaftaran
  ├─ Create PendaftaranNikah record (status: Draft)
  ├─ Create CalonPasangan for bride & groom
  ├─ Create WaliNikah record
  └─ Return: nomor_pendaftaran
  ↓
Status Check
GET /simnikah/pendaftaran/status
  ├─ Find PendaftaranNikah by user_id
  └─ Return: registration + calon pasangan + penghulu (if assigned)
  ↓
Staff Verification
POST /simnikah/staff/verify-formulir/:id
  └─ Update status → Mark verified
  ↓
Staff Approval
POST /simnikah/staff/approve/:id
  └─ Update status → "Menunggu Penugasan"
  ↓
Kepala KUA Assignment
POST /simnikah/pendaftaran/:id/assign-penghulu
  ├─ Check penghulu availability (jadwal + workload)
  ├─ Create notification for penghulu
  └─ Update status → "Penghulu Ditugaskan"
  ↓
Penghulu Complete
POST /simnikah/penghulu/complete-marriage/:id
  └─ Update status → "Selesai"
```

### 6️⃣ Database Relationship Diagram
```
Users
├── user_id (PK)
├── username
├── email
├── password
├── role
├── nama
└── profile_photo

StaffKUA
├── id (PK)
├── user_id (FK → Users)
├── nip
├── jabatan
└── status

Penghulu
├── id (PK)
├── user_id (FK → Users)
├── nip
├── jumlah_nikah
└── rating

CalonPasangan
├── id (PK)
├── user_id (FK → Users)
├── nik
├── nama_lengkap
├── tanggal_lahir
└── jenis_kelamin

PendaftaranNikah
├── id (PK)
├── nomor_pendaftaran (UNIQUE)
├── pendaftar_id (FK → Users)
├── calon_suami_id (FK → Users)
├── calon_istri_id (FK → Users)
├── penghulu_id (FK → Penghulu)
├── wali_nikah_id (FK → WaliNikah)
└── status_pendaftaran

WaliNikah
├── id (PK)
├── pendaftaran_id (FK → PendaftaranNikah)
├── nama_dan_bin
└── hubungan_wali

Notifikasi
├── id (PK)
├── user_id (FK → Users)
├── judul
├── pesan
└── status_baca

FeedbackPernikahan
├── id (PK)
├── pendaftaran_id (FK → PendaftaranNikah)
├── user_id (FK → Users)
├── rating
└── jenis_feedback
```

### 7️⃣ File Structure Navigation

**Core Files to Understand (Priority Order)**
```
1. cmd/api/main.go                    ← Start here! Entry point
2. internal/handlers/auth/auth.go      ← Authentication logic
3. internal/models/models.go           ← Database schemas
4. internal/models/constants.go        ← All constants
5. internal/handlers/catin/daftar.go   ← Marriage registration
6. internal/middleware/auth.go         ← Token verification
7. config/config.go                    ← Database connection
8. pkg/storage/imgbb.go                ← Image upload
9. pkg/utils/jwt.go                    ← JWT token generation
10. internal/handlers/staff/staff.go   ← Staff operations
```

### 8️⃣ Key Constants & Enums

**User Roles**
```
user_biasa    → Calon pengantin (regular user)
staff         → Staff KUA (verification officer)
penghulu      → Marriage officer (executor)
kepala_kua    → Head of KUA (assignment & approval)
```

**Registration Status Flow**
```
Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
                                                                  ↓
                                                            OR Ditolak
```

**User Status**
```
Aktif       → Active account
Nonaktif    → Inactive account
Blokir      → Blocked account
```

**Feedback Types**
```
Rating   → Star rating 1-5
Saran    → Suggestion
Kritik   → Criticism
Laporan  → Report/Complaint
```

**Notification Types**
```
Info      → Informational
Success   → Success notification
Warning   → Warning message
Error     → Error message
```

### 9️⃣ Common Code Patterns

**Pattern 1: Handler with Dependency Injection**
```go
type InDB struct {
    DB *gorm.DB
}

func (h *InDB) HandlerName(c *gin.Context) {
    // Access database via h.DB
}
```

**Pattern 2: JSON Response Format**
```go
c.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "Operation succeeded",
    "data": gin.H{
        "key": "value",
    },
})
```

**Pattern 3: Error Response**
```go
c.JSON(http.StatusBadRequest, gin.H{
    "success": false,
    "message": "User-friendly message",
    "error": "Technical error details",
    "type": "error_type",
})
```

**Pattern 4: Get User ID from Context**
```go
userID, exists := c.Get("user_id")
if !exists {
    // Return unauthorized
}
// Use userID as string
```

**Pattern 5: Database Query with GORM**
```go
var user structs.Users
h.DB.Where("username = ?", username).First(&user)

if err := h.DB.Create(&user).Error; err != nil {
    // Handle error
}
```

### 🔟 Important Validations

**Marriage Registration Validation**
```
✓ Umur calon ≥ 19 tahun (laki-laki) & ≥ 16 tahun (perempuan)
✓ Tanggal nikah ≥ hari ini
✓ Tanggal nikah ≠ hari libur nasional/minggu
✓ Waktu nikah dalam TimeSlots (08:00-16:00)
✓ Jadwal tidak penuh (max 9 per hari di KUA)
✓ Wali nikah sesuai urutan wali nasab
✓ Lokasi valid (Di KUA atau Di Luar KUA with address)
```

**User Validation**
```
✓ Username unique & alphanumeric
✓ Email valid & unique
✓ Password ≥ 6 characters
✓ Role dari daftar valid roles
```

### 1️⃣1️⃣ Testing Checklist

**Manual Testing Steps**
```
1. [ ] Register user: POST /register
2. [ ] Login: POST /login
3. [ ] Get token & check JWT claims
4. [ ] Call protected endpoint: GET /profile
5. [ ] Create marriage registration: POST /simnikah/pendaftaran
6. [ ] Check status: GET /simnikah/pendaftaran/status
7. [ ] Staff verify: POST /simnikah/staff/verify-formulir/:id
8. [ ] Staff approve: POST /simnikah/staff/approve/:id
9. [ ] Assign penghulu: POST /simnikah/pendaftaran/:id/assign-penghulu
10. [ ] Complete marriage: POST /simnikah/penghulu/complete-marriage/:id
```

---

## 🔧 Troubleshooting Guide

| Problem | Solution |
|---------|----------|
| 404 Endpoint tidak ditemukan | Periksa route di main.go, method POST/GET/PUT/DELETE, URL path |
| 401 Unauthorized | Periksa Authorization header, JWT token valid?, token expired? |
| 403 Forbidden | Role tidak sesuai, periksa middleware.RoleMiddleware() |
| 500 Database error | Periksa DB_PASSWORD, DB_HOST, migration successful?, table exists? |
| CORS error | Add origin ke ALLOWED_ORIGINS, periksa HTTP method (OPTIONS preflight) |
| Rate limit error | Tunggu atau reduce request frequency, check middleware config |
| Image upload failed | Periksa IMGBB_API_KEY, file size, format |

---

## 📚 Documentation References

| File | Purpose |
|------|---------|
| `/docs/FULL_API_USAGE.md` | All API endpoints with examples |
| `/docs/FRONTEND_IMPLEMENTATION_GUIDE.md` | Frontend developer guide |
| `/docs/API_REQUEST_BODY_DOCUMENTATION.md` | Request/response formats |
| `/docs/KEPALA_KUA_STAFF_MANAGEMENT.md` | Staff management specific |
| `/docs/STATUS_MANAGEMENT.md` | Status flow details |
| `/migrations/init.sql` | Database schema |
| `README.md` (if exists) | Project overview |

---

## 🎯 Next Steps

1. **Read**: `cmd/api/main.go` - Understand entry point
2. **Read**: `internal/models/models.go` - Understand data structure
3. **Read**: `internal/handlers/auth/auth.go` - Understand auth flow
4. **Test**: Use Postman to test endpoints
5. **Read**: Handler files in order of interest
6. **Experiment**: Make changes, test locally, understand behavior

---

**Happy Learning! 🚀**
