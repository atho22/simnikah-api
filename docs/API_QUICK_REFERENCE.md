# 🚀 SimNikah API - Quick Reference Guide

## 📌 Base URL
```
Local: http://localhost:8080
Production: https://your-app.railway.app
```

## 🔐 Authentication Header
```
Authorization: Bearer <your_jwt_token>
```

---

## Quick Links by User Role

### 👤 User Biasa (Calon Pengantin)

| Action | Method | Endpoint |
|--------|--------|----------|
| Register | POST | `/register` |
| Login | POST | `/login` |
| Create Registration | POST | `/simnikah/pendaftaran/form-baru` |
| Check Status | GET | `/simnikah/pendaftaran/status` |
| Update Location | PUT | `/simnikah/pendaftaran/:id/location` |
| Register Bimbingan | POST | `/simnikah/bimbingan/:id/daftar` |
| Get Notifications | GET | `/simnikah/notifikasi/user/:user_id` |

### 👨‍💼 Staff KUA

| Action | Method | Endpoint |
|--------|--------|----------|
| Get All Registrations | GET | `/simnikah/pendaftaran` |
| Verify Form | POST | `/simnikah/staff/verify-formulir/:id` |
| Verify Documents | POST | `/simnikah/staff/verify-berkas/:id` |
| Mark as Visited | POST | `/simnikah/pendaftaran/:id/mark-visited` |
| Create Bimbingan | POST | `/simnikah/bimbingan` |
| Update Attendance | PUT | `/simnikah/bimbingan/:id/update-attendance` |

### 👨‍🏫 Penghulu

| Action | Method | Endpoint |
|--------|--------|----------|
| Get Assigned Registrations | GET | `/simnikah/penghulu/assigned-registrations` |
| Verify Documents | POST | `/simnikah/penghulu/verify-documents/:id` |
| Get Location Detail | GET | `/simnikah/pendaftaran/:id/location` |
| Get Schedule | GET | `/simnikah/penghulu-jadwal/:tanggal` |

### 👨‍💼 Kepala KUA

| Action | Method | Endpoint |
|--------|--------|----------|
| Create Staff | POST | `/simnikah/staff` |
| Create Penghulu | POST | `/simnikah/penghulu` |
| Assign Penghulu | POST | `/simnikah/pendaftaran/:id/assign-penghulu` |
| Get Penghulu Availability | GET | `/simnikah/penghulu/:id/ketersediaan/:tanggal` |
| Send Bulk Notification | POST | `/simnikah/notifikasi/send-to-role` |

---

## 📋 Common Request Examples

### 1. Register & Login Flow

**Step 1: Register**
```bash
POST /register
{
  "username": "ahmad_fauzi",
  "email": "ahmad@example.com",
  "password": "secure123",
  "nama": "Ahmad Fauzi",
  "role": "user_biasa"
}
```

**Step 2: Login**
```bash
POST /login
{
  "username": "ahmad_fauzi",
  "password": "secure123"
}

# Response includes JWT token
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Step 3: Use Token**
```bash
GET /profile
Headers:
  Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

### 2. Marriage Registration Flow

**Step 1: Create Registration (User)**
```bash
POST /simnikah/pendaftaran/form-baru
Headers:
  Authorization: Bearer <user_token>
Body: {
  "scheduleAndLocation": { ... },
  "groom": { ... },
  "bride": { ... },
  "groomParents": { ... },
  "brideParents": { ... },
  "guardian": { ... }
}
```

**Step 2: Check Status (User)**
```bash
GET /simnikah/pendaftaran/status
Headers:
  Authorization: Bearer <user_token>
```

**Step 3: Staff Verify Form**
```bash
POST /simnikah/staff/verify-formulir/1
Headers:
  Authorization: Bearer <staff_token>
Body: {
  "status": "Formulir Disetujui",
  "catatan": "Formulir sudah lengkap"
}
```

**Step 4: Staff Verify Documents**
```bash
POST /simnikah/staff/verify-berkas/1
Headers:
  Authorization: Bearer <staff_token>
Body: {
  "status": "Berkas Diterima",
  "catatan": "Semua berkas lengkap"
}
```

**Step 5: Mark as Visited**
```bash
POST /simnikah/pendaftaran/1/mark-visited
Headers:
  Authorization: Bearer <staff_token>
```

**Step 6: Kepala KUA Assign Penghulu**
```bash
POST /simnikah/pendaftaran/1/assign-penghulu
Headers:
  Authorization: Bearer <kepala_token>
Body: {
  "penghulu_id": 1
}
```

**Step 7: Penghulu Verify**
```bash
POST /simnikah/penghulu/verify-documents/1
Headers:
  Authorization: Bearer <penghulu_token>
Body: {
  "status": "Menunggu Pelaksanaan",
  "catatan": "Berkas sesuai syariat"
}
```

---

### 3. Bimbingan Flow

**Step 1: Create Bimbingan Session (Staff)**
```bash
POST /simnikah/bimbingan
Headers:
  Authorization: Bearer <staff_token>
Body: {
  "tanggal_bimbingan": "2024-12-18",
  "waktu_mulai": "08:00",
  "waktu_selesai": "12:00",
  "tempat_bimbingan": "Aula KUA",
  "pembimbing": "H. Ahmad Dahlan",
  "kapasitas": 10
}
```

**Step 2: User Register for Bimbingan**
```bash
POST /simnikah/bimbingan/1/daftar
Headers:
  Authorization: Bearer <user_token>
```

**Step 3: Update Attendance (Staff)**
```bash
PUT /simnikah/bimbingan/1/update-attendance
Headers:
  Authorization: Bearer <staff_token>
Body: {
  "pendaftaran_nikah_id": 1,
  "status_kehadiran": "Hadir",
  "status_sertifikat": "Sudah",
  "no_sertifikat": "CERT-2024-001"
}
```

---

### 4. Location/Map Integration

**Geocoding (Address → Coordinates)**
```bash
POST /simnikah/location/geocode
Headers:
  Authorization: Bearer <token>
Body: {
  "alamat": "Jl. Merdeka No. 123, Banjarmasin"
}

# Response
{
  "latitude": -3.3194,
  "longitude": 114.5903,
  "map_url": "https://www.google.com/maps?q=-3.319400,114.590300"
}
```

**Search Address (Autocomplete)**
```bash
GET /simnikah/location/search?q=Jl. Merdeka
Headers:
  Authorization: Bearer <token>
```

**Update Wedding Location**
```bash
PUT /simnikah/pendaftaran/1/location
Headers:
  Authorization: Bearer <token>
Body: {
  "alamat_akad": "Jl. Merdeka No. 123, Banjarmasin",
  "latitude": -3.3194,
  "longitude": 114.5903
}
```

---

## 📅 Calendar & Schedule Quick Access

### Get Monthly Calendar
```bash
GET /simnikah/kalender-ketersediaan?bulan=12&tahun=2024
```

### Get Date Detail
```bash
GET /simnikah/kalender-tanggal-detail?tanggal=2024-12-25
```

### Get Penghulu Schedule
```bash
GET /simnikah/penghulu-jadwal/2024-12-25
```

### Get Penghulu Availability
```bash
GET /simnikah/penghulu/1/ketersediaan/2024-12-25
```

---

## 🔔 Notification Quick Access

### Get User Notifications
```bash
GET /simnikah/notifikasi/user/USR1234567890?page=1&limit=20&status=Belum Dibaca
```

### Mark as Read
```bash
PUT /simnikah/notifikasi/1/status
Body: {
  "status_baca": "Sudah Dibaca"
}
```

### Mark All as Read
```bash
PUT /simnikah/notifikasi/user/USR1234567890/mark-all-read
```

### Get Stats
```bash
GET /simnikah/notifikasi/user/USR1234567890/stats
```

### Send to All Users with Role
```bash
POST /simnikah/notifikasi/send-to-role
Body: {
  "role": "user_biasa",
  "judul": "Pengumuman Penting",
  "pesan": "KUA akan tutup tanggal 1-3 Januari 2025",
  "tipe": "Info"
}
```

---

## 🎯 Status Codes Quick Reference

| Code | Meaning | Action |
|------|---------|--------|
| 200 | OK | Success |
| 201 | Created | Resource created |
| 400 | Bad Request | Check request format |
| 401 | Unauthorized | Check token |
| 403 | Forbidden | Check role permission |
| 404 | Not Found | Resource doesn't exist |
| 429 | Rate Limited | Wait before retry |
| 500 | Server Error | Contact support |

---

## 🔑 Valid Values Reference

### User Roles
- `user_biasa` - Calon pengantin
- `staff` - Staff KUA
- `penghulu` - Penghulu
- `kepala_kua` - Kepala KUA

### Registration Status
1. `Draft`
2. `Menunggu Verifikasi`
3. `Menunggu Pengumpulan Berkas`
4. `Berkas Diterima`
5. `Menunggu Penugasan`
6. `Penghulu Ditugaskan`
7. `Menunggu Verifikasi Penghulu`
8. `Menunggu Bimbingan`
9. `Sudah Bimbingan`
10. `Selesai`

### Marital Status
- `Belum Kawin`
- `Cerai Hidup`
- `Cerai Mati`

### Religions
- `Islam`
- `Kristen`
- `Katolik`
- `Hindu`
- `Buddha`
- `Konghucu`

### Citizenship
- `WNI` - Warga Negara Indonesia
- `WNA` - Warga Negara Asing

### Guardian Relationships
- `Ayah Kandung`
- `Kakek`
- `Saudara Kandung`
- `Paman`
- `Wali Hakim`

### Presence Status
- `Hidup`
- `Meninggal`

### Notification Types
- `Info`
- `Warning`
- `Error`
- `Success`

### Notification Status
- `Belum Dibaca`
- `Sudah Dibaca`

### Staff Positions
- `Staff`
- `Penghulu`
- `Kepala KUA`

### Active Status
- `Aktif`
- `Tidak Aktif`

### Wedding Location
- `Di KUA`
- `Di Luar KUA`

### Attendance Status
- `Belum`
- `Hadir`
- `Tidak Hadir`

### Certificate Status
- `Belum`
- `Sudah`

---

## 🚦 Business Rules Quick Reference

### Calendar Rules
- ✅ Max 9 marriages per day at KUA
- ✅ Unlimited marriages outside KUA
- ✅ Max 3 marriages per penghulu per day
- ✅ Minimum 1 hour gap between penghulu schedules

### Bimbingan Rules
- ✅ Only on Wednesdays
- ✅ Max 10 couples per session
- ✅ Only 1 session per day

### Document Verification
- ✅ Staff verifies form first
- ✅ Then staff verifies physical documents
- ✅ Finally penghulu verifies everything

### Age Requirements
- ✅ Groom: Minimum 19 years old
- ✅ Bride: Minimum 19 years old
- ℹ️ Dispensation available with court order

---

## 💡 Pro Tips

### 1. Token Management
- Store JWT token securely
- Token expires after 24 hours
- Refresh by logging in again

### 2. Pagination Best Practices
```bash
# Start with small pages
GET /endpoint?page=1&limit=10

# Increase if needed
GET /endpoint?page=1&limit=50
```

### 3. Filtering
```bash
# Combine multiple filters
GET /simnikah/pendaftaran?status=Menunggu Verifikasi&page=1&limit=20
```

### 4. Date Formats
- Always use: `YYYY-MM-DD` for dates
- Always use: `HH:MM` for times (24-hour format)

### 5. Location Data
- Use geocoding endpoint first to get coordinates
- Then save both address and coordinates
- Coordinates make map integration easier

### 6. Error Handling
```javascript
// Always check response status
if (response.status === 401) {
  // Token expired, redirect to login
}
```

---

## 🧪 Testing Checklist

### User Registration Flow
- [ ] Register new user
- [ ] Login with credentials
- [ ] Get profile
- [ ] Create marriage registration
- [ ] Check registration status

### Staff Verification Flow
- [ ] Get all registrations
- [ ] Verify form
- [ ] Verify documents
- [ ] Mark as visited

### Penghulu Assignment Flow
- [ ] Get penghulu list
- [ ] Check penghulu availability
- [ ] Assign penghulu to registration
- [ ] Penghulu verifies documents

### Bimbingan Flow
- [ ] Create bimbingan session
- [ ] User registers for bimbingan
- [ ] Update attendance
- [ ] Complete bimbingan

### Notification Flow
- [ ] Get notifications
- [ ] Mark as read
- [ ] Get notification stats
- [ ] Send bulk notification

---

## 📞 Support

**Need Help?**
- 📖 Full Documentation: [API_DOCUMENTATION_COMPLETE.md](./API_DOCUMENTATION_COMPLETE.md)
- 📧 Email: support@simnikah.id
- 🐛 Report Issues: GitHub Issues

---

**Last Updated:** November 2024  
**API Version:** 1.0.0

