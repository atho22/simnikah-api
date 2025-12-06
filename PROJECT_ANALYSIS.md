# 📊 Analisis Lengkap Proyek SimNikah API

## 1. 🎯 Ringkasan Proyek

**SimNikah** adalah sebuah sistem digital untuk manajemen pendaftaran nikah (pernikahan) di KUA (Kantor Urusan Agama). Sistem ini memudahkan:
- **Calon Pengantin (Catin)**: Mendaftar pernikahan secara online
- **Staff KUA**: Memverifikasi dan mengelola pendaftaran
- **Penghulu**: Melaksanakan nikah dan mengelola jadwal
- **Kepala KUA**: Monitoring dan assignment penghulu

---

## 2. 🏗️ Struktur Arsitektur

### 2.1 **Tech Stack**
```
Backend:          Go (Golang) 1.23.6
Framework Web:    Gin Web Framework
Database:         MySQL 8.0+
ORM:              GORM v1.26.1
Authentication:   JWT (golang-jwt/jwt v5)
Image Storage:    ImgBB (Cloud Storage)
Encryption:       bcrypt
```

### 2.2 **Struktur Folder**
```
simnikah-api/
├── cmd/api/main.go           # Entry point aplikasi
├── config/                   # Konfigurasi database
├── internal/
│   ├── handlers/             # HTTP request handlers per fitur
│   │   ├── auth/             # Register, Login, GetProfile
│   │   ├── catin/            # Pendaftaran nikah, Status
│   │   ├── staff/            # Verifikasi, Approval
│   │   ├── penghulu/         # Jadwal nikah, Completion
│   │   ├── kepala_kua/       # Assignment, Statistics
│   │   ├── dashboard/        # Dashboard per role
│   │   └── notification/     # Notifikasi
│   ├── middleware/           # Auth, Rate Limit, Role-based
│   ├── models/               # Data models & constants
│   ├── seeders/              # Initial data (Kepala KUA, Staff, Penghulu)
│   └── services/             # Business logic (Cron, Notifications)
├── pkg/                      # Utility packages
│   ├── crypto/               # BCrypt password hashing
│   ├── storage/              # ImgBB image upload
│   ├── utils/                # JWT, String, Date utilities
│   └── validator/            # Marriage validation logic
├── migrations/               # Database schema
├── deployments/              # Docker & Railway configs
└── docs/                     # Documentation & guides
```

---

## 3. 📋 Models/Database Schema

### 3.1 **Users** - Account Management
```go
type Users struct {
    ID              uint      // Primary key
    User_id         string    // Unique user ID (USR + timestamp)
    Username        string    // Unique username
    Email           string    // Unique email
    Password        string    // Hashed password
    Role            string    // user_biasa, penghulu, staff, kepala_kua
    Status          string    // Aktif, Nonaktif, Blokir
    Nama            string    // Full name
    Profile_photo   string    // Profile photo URL (ImgBB)
    Created_at      time.Time
    Updated_at      time.Time
}
```

### 3.2 **StaffKUA** - Staff Member Information
```go
type StaffKUA struct {
    ID           uint
    User_id      string    // FK to Users
    NIP          string    // Employee ID
    Nama_lengkap string
    Jabatan      string    // Staff, Penghulu, Kepala KUA
    Bagian       string    // Department
    No_hp        string
    Email        string
    Alamat       string
    Status       string    // Aktif, Nonaktif
    Created_at   time.Time
    Updated_at   time.Time
}
```

### 3.3 **Penghulu** - Marriage Officer
```go
type Penghulu struct {
    ID           uint
    User_id      string    // FK to Users
    NIP          string    // Employee ID
    Nama_lengkap string
    No_hp        string
    Email        string
    Alamat       string
    Status       string    // Aktif, Nonaktif
    Jumlah_nikah int       // Wedding count
    Rating       float64   // Performance rating
    Created_at   time.Time
    Updated_at   time.Time
}
```

### 3.4 **CalonPasangan** - Bride/Groom Data
```go
type CalonPasangan struct {
    ID                  uint
    User_id             string    // FK to Users
    NIK                 string    // National ID
    Nama_lengkap        string
    Tanggal_lahir       time.Time
    Jenis_kelamin       string    // L/P
    Pendidikan_terakhir string
    Created_at          time.Time
    Updated_at          time.Time
}
```

### 3.5 **PendaftaranNikah** - Marriage Registration
```go
type PendaftaranNikah struct {
    ID                   uint
    Nomor_pendaftaran    string    // Unique registration number
    Pendaftar_id         string    // Who registered (groom or bride)
    Calon_suami_id       string    // FK to Users (groom)
    Calon_istri_id       string    // FK to Users (bride)
    Wali_nikah_id        *uint     // FK to WaliNikah
    Tanggal_pendaftaran  time.Time
    Tanggal_nikah        time.Time
    Waktu_nikah          string    // HH:MM format
    Tempat_nikah         string    // Di KUA or Di Luar KUA
    Alamat_akad          string    // Wedding address if outside KUA
    Latitude             *float64
    Longitude            *float64
    Status_pendaftaran   string    // Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
    Penghulu_id          *uint     // Assigned marriage officer
    Penghulu_assigned_by string    // Who assigned
    Penghulu_assigned_at *time.Time
    Catatan              string    // Notes
    Disetujui_oleh       string    // Approved by
    Disetujui_pada       *time.Time
    Created_at           time.Time
    Updated_at           time.Time
}
```

### 3.6 **WaliNikah** - Guardian for Bride
```go
type WaliNikah struct {
    ID            uint
    Pendaftaran_id uint     // FK to PendaftaranNikah
    Nama_dan_bin  string
    Hubungan_wali string   // Ayah Kandung, Kakek, dll
    Created_at    time.Time
    Updated_at    time.Time
}
```

### 3.7 **FeedbackPernikahan** - User Feedback
```go
type FeedbackPernikahan struct {
    ID               uint
    Pendaftaran_id   uint
    User_id          string
    Jenis_feedback   string    // Rating, Saran, Kritik, Laporan
    Rating           *int      // 1-5
    Judul            string
    Pesan            string
    Status_baca      string    // Belum Dibaca, Sudah Dibaca
    Dibaca_oleh      string    // Read by
    Dibaca_pada      *time.Time
    Created_at       time.Time
    Updated_at       time.Time
}
```

### 3.8 **Notifikasi** - User Notifications
```go
type Notifikasi struct {
    ID          uint
    User_id     string
    Judul       string
    Pesan       string
    Tipe        string    // Info, Success, Warning, Error
    Status_baca string    // Belum Dibaca, Sudah Dibaca
    Link        string    // Optional link
    Created_at  time.Time
    Updated_at  time.Time
}
```

---

## 4. 🔀 Flow Status Pendaftaran Nikah

```
START
  ↓
[Draft] ← User mengirim form pendaftaran
  ↓
[Disetujui] ← Staff verifikasi & approve
  ↓
[Menunggu Penugasan] ← Menunggu Kepala KUA assign penghulu
  ↓
[Penghulu Ditugaskan] ← Kepala KUA sudah assign penghulu
  ↓
[Selesai] ← Penghulu melaksanakan nikah
  ↓
END

OR (Rejected path)
  ↓
[Ditolak] ← Staff atau Kepala KUA tolak
  ↓
END
```

---

## 5. 🔐 Authentication & Authorization

### 5.1 **JWT Token Claims**
```go
type TokenClaims struct {
    UserID string
    Email  string
    Role   string
    Nama   string
    RegisteredClaims jwt.RegisteredClaims
}
```

### 5.2 **Role-Based Access Control**
- **user_biasa**: Catin (Calon Pengantin) - hanya bisa lihat status sendiri
- **staff**: Staff KUA - verifikasi & approve pendaftaran
- **penghulu**: Marriage Officer - manage jadwal & completion
- **kepala_kua**: Head of KUA - assign penghulu, dashboard, feedback management

### 5.3 **Middleware**
- `AuthMiddleware()`: Verify JWT token
- `RoleMiddleware(roles...)`: Check specific role
- `MultiRoleMiddleware(roles...)`: Check multiple roles
- `RateLimiter()`: Global rate limit (100 req/min per IP)
- `StrictRateLimiter()`: Auth endpoints (stricter limits)

---

## 6. 🛣️ Main API Endpoints

### 6.1 **Authentication**
```
POST   /register                    # Create new user account
POST   /login                       # Login & get JWT token
GET    /profile                     # Get current user profile
```

### 6.2 **Catin (User Biasa)**
```
POST   /simnikah/pendaftaran                              # Register marriage
GET    /simnikah/pendaftaran/status                       # Check registration status
POST   /simnikah/feedback-pernikahan                      # Submit feedback
GET    /simnikah/kalender-ketersediaan                    # View availability calendar
GET    /simnikah/ketersediaan-jam                         # Check available time slots
GET    /simnikah/pernikahan-tanggal                       # Weddings by date
POST   /simnikah/location/geocode                         # Address to coordinates
POST   /simnikah/location/reverse-geocode                 # Coordinates to address
```

### 6.3 **Staff KUA**
```
GET    /simnikah/pendaftaran                              # List registrations
POST   /simnikah/staff/verify-formulir/:id                # Verify form
POST   /simnikah/staff/approve/:id                        # Approve registration
POST   /simnikah/staff/pendaftaran                        # Register user manually
PUT    /simnikah/pendaftaran/:id/update-status            # Update status
GET    /simnikah/staff/pengumuman-nikah/list              # Marriage announcement list
```

### 6.4 **Penghulu**
```
GET    /simnikah/penghulu                                 # List marriage officers
PUT    /simnikah/penghulu/:id                             # Update officer info
POST   /simnikah/penghulu/verify-documents/:id            # Verify documents
GET    /simnikah/penghulu/assigned-registrations          # My assignments
GET    /simnikah/penghulu/today-schedule                  # Today schedule
POST   /simnikah/penghulu/complete-marriage/:id           # Complete marriage
```

### 6.5 **Kepala KUA**
```
POST   /simnikah/kepala-kua/staff                         # Create staff member
POST   /simnikah/kepala-kua/penghulu                      # Create marriage officer
POST   /simnikah/pendaftaran/:id/assign-penghulu          # Assign officer to marriage
GET    /simnikah/kepala-kua/available-penghulu            # Available officers
GET    /simnikah/kepala-kua/statistik-penghulu            # Officer statistics
GET    /simnikah/kepala-kua/feedback                      # List feedback
PUT    /simnikah/kepala-kua/feedback/:id/mark-read        # Mark feedback as read
```

### 6.6 **Dashboard**
```
GET    /simnikah/dashboard/kepala-kua                     # Head dashboard
GET    /simnikah/dashboard/staff                          # Staff dashboard
GET    /simnikah/dashboard/statistik-pernikahan           # Marriage statistics
GET    /simnikah/dashboard/penghulu-performance           # Officer performance
```

### 6.7 **Notifications**
```
GET    /simnikah/notifikasi/user/:user_id                 # Get user notifications
GET    /simnikah/notifikasi/:id                           # Get notification detail
PUT    /simnikah/notifikasi/:id/status                    # Mark as read
PUT    /simnikah/notifikasi/user/:user_id/mark-all-read   # Mark all as read
DELETE /simnikah/notifikasi/:id                           # Delete notification
POST   /simnikah/notifikasi/send-to-role                  # Send to role group
```

---

## 7. 🎯 Key Features

### 7.1 **Pendaftaran Nikah Sederhana (Simple Registration)**
- Form wizard 4-step:
  1. Data calon suami & istri (nama, pendidikan, umur)
  2. Data wali nikah
  3. Jadwal & lokasi (Di KUA atau Di Luar KUA)
  4. Review & submit
- Validasi umur, ketersediaan jadwal
- Support lokasi di luar KUA dengan koordinat GPS

### 7.2 **Kalender Ketersediaan (Availability Calendar)**
- 9 slot/hari di KUA (08:00-16:00, per jam)
- Deteksi hari libur nasional otomatis
- Warna status: Hijau (tersedia), Merah (penuh), Kuning (hampir penuh)

### 7.3 **Sistem Penugasan Penghulu (Officer Assignment)**
- Kepala KUA assign penghulu berdasarkan:
  - Ketersediaan jadwal
  - Jumlah nikah yang sudah dilayani hari itu
- Notifikasi otomatis saat assignment

### 7.4 **Feedback & Rating**
- Catin bisa kasih rating (1-5) setelah nikah selesai
- Tipe feedback: Rating, Saran, Kritik, Laporan
- Kepala KUA bisa lihat & manage feedback

### 7.5 **Notifikasi Real-time**
- Cron job setiap hari jam 08:00
- Reminder untuk pendaftaran yang belum diverifikasi
- Tipe notifikasi: Info, Success, Warning, Error
- Mark as read tracking

### 7.6 **Dashboard & Analytics**
- Kepala KUA: Total nikah, pending, rating penghulu, grafik trend
- Staff: Pending verifikasi, nikah hari ini
- Penghulu: Jadwal hari ini, assignment list

### 7.7 **Image Upload (Profile Photo)**
- Menggunakan ImgBB (cloud storage gratis)
- API key dari environment variable
- Support anonymous upload

---

## 8. 🔄 Deployment & Configuration

### 8.1 **Environment Variables**
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=simnikah

# Server
PORT=8080
GIN_MODE=debug|release

# CORS
ALLOWED_ORIGINS=http://localhost:3000,https://kua-ku.vercel.app

# ImgBB
IMGBB_API_KEY=your_key_here
```

### 8.2 **Deployment Options**
- **Railway**: Docker + MySQL service
- **Docker Compose**: Local development
- **Direct Go**: `go run cmd/api/main.go`

### 8.3 **Makefile Commands**
```bash
make build      # Compile to binary
make run        # Run production mode
make dev        # Run debug mode (hot reload)
make watch      # Live reload with 'air'
make test       # Run tests
make coverage   # Generate coverage report
make clean      # Clean build artifacts
make deps       # Update dependencies
```

---

## 9. 💡 Koneksi ke Frontend

Frontend sudah disiapkan di repo terpisah dengan flow:
1. User register/login via `/register`, `/login`
2. Get token dari JWT response
3. Simpan token di Cookie (HttpOnly recommended) atau localStorage
4. Kirim token di setiap request: `Authorization: Bearer <token>`
5. Frontend redirect berdasarkan role:
   - user_biasa → `/dashboard/catin`
   - staff → `/dashboard/staff`
   - penghulu → `/dashboard/penghulu`
   - kepala_kua → `/dashboard/kepala-kua`

---

## 10. 📌 Best Practices Implementasi

### 10.1 **Error Handling**
```json
{
  "success": false,
  "message": "User-friendly message",
  "error": "Technical error details",
  "type": "error_type" // schedule_conflict, validation, authentication, etc
}
```

### 10.2 **Rate Limiting**
- Global: 100 req/min per IP
- Auth endpoints (register/login): Stricter limits
- Di-implement via `ulule/limiter` package

### 10.3 **Database Performance**
- GORM AutoMigrate untuk schema management
- Indexes di critical fields
- Connection pooling (10 idle, 100 max open)
- UTC timezone untuk consistency

### 10.4 **Security**
- Password di-hash dengan bcrypt
- JWT token dengan 24-hour expiry
- CORS configuration
- SQL injection protection via GORM parameterized queries
- Rate limiting

---

## 11. 🐛 Testing & Debugging

### 11.1 **Test Database**
- Gunakan seeder untuk initial data
- `SeedKepalaKUA()`: Create 1 Kepala KUA user
- `SeedStaff()`: Create 2 staff users
- `SeedPenghulu()`: Create 3 penghulu users

### 11.2 **API Testing**
- Gunakan Postman/Insomnia
- Collection tersedia di docs
- Test flow: Register → Login → Create Registration → etc

### 11.3 **Common Issues**
- DB connection error: Check DB_PASSWORD, DB_HOST
- CORS error: Add origin ke ALLOWED_ORIGINS
- File upload error: Check IMGBB_API_KEY
- Rate limit hit: Reduce request frequency

---

## 12. 🚀 Roadmap Fitur

- [ ] WhatsApp notification integration
- [ ] E-signature untuk dokumen nikah
- [ ] Export to PDF (Surat Pengumuman Nikah)
- [ ] Mobile app (Flutter/React Native)
- [ ] Advanced analytics & reporting
- [ ] Multi-language support (EN)
- [ ] Two-factor authentication

---

## 📞 Kontak & Support

**Dokumentasi Lengkap:**
- `/docs/API_DOCUMENTATION_COMPLETE.md` - Endpoint detail
- `/docs/FRONTEND_IMPLEMENTATION_GUIDE.md` - Frontend flow
- `/docs/API_REQUEST_BODY_DOCUMENTATION.md` - Request/response examples

**Repository:**
- Backend: `github.com/atho22/simnikah-api` (main branch)
- Frontend: (Link ke frontend repo)

---

**Last Updated:** December 6, 2025
**Version:** 1.0.0
