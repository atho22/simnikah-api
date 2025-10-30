# ⚡ SimNikah API - Quick Reference

**Base URL:** `http://localhost:8080` | `https://your-api.railway.app`

---

## 🔐 Authentication

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/register` | POST | ❌ | Register new user |
| `/login` | POST | ❌ | Login & get token |
| `/profile` | GET | ✅ | Get user profile |

**Token:** Send in header: `Authorization: Bearer <token>`  
**Valid for:** 24 hours

---

## 👤 User Roles

| Role | Code | Access |
|------|------|--------|
| User | `user_biasa` | Register marriage, view own status |
| Staff | `staff` | Verify forms & documents |
| Penghulu | `penghulu` | Conduct marriages |
| Head KUA | `kepala_kua` | Full admin access |

---

## 💍 Marriage Registration

| Endpoint | Method | Auth | Role | Description |
|----------|--------|------|------|-------------|
| `/simnikah/pendaftaran/form-baru` | POST | ✅ | user_biasa | Create registration ⚡ |
| `/simnikah/pendaftaran/status` | GET | ✅ | user_biasa | Check own status |
| `/simnikah/pendaftaran/:id/status-flow` | GET | ✅ | All | Get status flow |
| `/simnikah/pendaftaran/:id/mark-visited` | POST | ✅ | user_biasa | Mark as visited |
| `/simnikah/pendaftaran` | GET | ✅ | staff, kepala_kua | Get all registrations |

⚡ **Performance:** 800-1200ms (optimized)

---

## 📊 Status Flow (10 Steps)

```
Draft → Menunggu Verifikasi → Menunggu Pengumpulan Berkas → 
Berkas Diterima → Menunggu Penugasan → Penghulu Ditugaskan → 
Menunggu Verifikasi Penghulu → Menunggu Bimbingan → 
Sudah Bimbingan → Selesai ✅
```

---

## 👨‍💼 Staff Operations

| Endpoint | Method | Auth | Role | Description |
|----------|--------|------|------|-------------|
| `/simnikah/staff` | POST | ✅ | kepala_kua | Create staff |
| `/simnikah/staff` | GET | ✅ | kepala_kua | Get all staff |
| `/simnikah/staff/verify-formulir/:id` | POST | ✅ | staff | Verify online form |
| `/simnikah/staff/verify-berkas/:id` | POST | ✅ | staff | Verify physical docs |

---

## 🕌 Penghulu Management

| Endpoint | Method | Auth | Role | Description |
|----------|--------|------|------|-------------|
| `/simnikah/penghulu` | POST | ✅ | kepala_kua | Create penghulu |
| `/simnikah/penghulu` | GET | ✅ | All | Get all penghulu |
| `/simnikah/pendaftaran/:id/assign-penghulu` | POST | ✅ | kepala_kua | Assign penghulu |
| `/simnikah/penghulu/verify-documents/:id` | POST | ✅ | penghulu | Verify documents |
| `/simnikah/penghulu/assigned-registrations` | GET | ✅ | penghulu | Get assignments |

---

## 📅 Calendar & Schedule

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/simnikah/kalender-ketersediaan` | GET | ✅ | Availability calendar |
| `/simnikah/kalender-tanggal-detail` | GET | ✅ | Date details |
| `/simnikah/ketersediaan-tanggal/:tanggal` | GET | ✅ | Check date availability |
| `/simnikah/penghulu-jadwal/:tanggal` | GET | ✅ | Penghulu schedule |

**Query Params:** `?bulan=12&tahun=2025`

---

## 📚 Counseling (Bimbingan)

| Endpoint | Method | Auth | Role | Description |
|----------|--------|------|------|-------------|
| `/simnikah/bimbingan` | POST | ✅ | staff, kepala_kua | Create session |
| `/simnikah/bimbingan` | GET | ✅ | All | Get sessions |
| `/simnikah/bimbingan/:id/daftar` | POST | ✅ | user_biasa | Register for session |
| `/simnikah/bimbingan/:id/participants` | GET | ✅ | staff, kepala_kua | Get participants |
| `/simnikah/bimbingan/:id/update-attendance` | PUT | ✅ | staff, kepala_kua | Update attendance |
| `/simnikah/bimbingan-kalender` | GET | ✅ | All | Counseling calendar |

**Rules:**
- Only on **Wednesdays**
- Max **10 couples** per session
- **Mandatory** before marriage

---

## 🗺️ Map & Location (FREE)

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/simnikah/location/geocode` | POST | ✅ | Address → Coordinates |
| `/simnikah/location/reverse-geocode` | POST | ✅ | Coordinates → Address |
| `/simnikah/location/search` | GET | ✅ | Address autocomplete |
| `/simnikah/pendaftaran/:id/location` | GET | ✅ | Get wedding location |
| `/simnikah/pendaftaran/:id/alamat` | PUT | ✅ | Update wedding address |

**Powered by:** OpenStreetMap Nominatim API

---

## 🔔 Notifications

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/simnikah/notifikasi/user/:user_id` | GET | ✅ | Get user notifications |
| `/simnikah/notifikasi/:id/status` | PUT | ✅ | Mark as read |
| `/simnikah/notifikasi/user/:user_id/mark-all-read` | PUT | ✅ | Mark all as read |
| `/simnikah/notifikasi/user/:user_id/stats` | GET | ✅ | Get notification stats |

**Query Params:** `?status=Belum Dibaca&limit=10`

---

## 📋 Valid Values

### Wedding Location
- `"Di KUA"` - At KUA office (max 9/day)
- `"Di Luar KUA"` - External venue (unlimited)

### Marital Status
- `"Belum Kawin"` - Never married
- `"Cerai Mati"` - Widowed
- `"Cerai Hidup"` - Divorced

### Guardian Relationship (Wali Nikah)
Priority order (nasab):
1. `"Ayah Kandung"` - Biological father (if alive)
2. `"Kakek"` - Grandfather
3. `"Saudara Laki-Laki Kandung"` - Full brother
4. `"Paman Kandung"` - Uncle
5. `"Wali Hakim"` - Judge (if no nasab)

### Presence Status
- `"Hidup"` - Alive
- `"Meninggal"` - Deceased

### Attendance Status
- `"Belum"` - Not yet
- `"Hadir"` - Present
- `"Tidak Hadir"` - Absent

---

## ⚠️ Common Error Codes

| Code | Meaning | Action |
|------|---------|--------|
| 400 | Bad Request | Check request body |
| 401 | Unauthorized | Login again |
| 403 | Forbidden | Check role permissions |
| 404 | Not Found | Check ID/URL |
| 500 | Server Error | Contact admin |

---

## 📊 Business Rules

| Rule | Value |
|------|-------|
| Max marriages per day (KUA) | 9 |
| Max marriages per penghulu | 3/day |
| Min gap between marriages | 60 min |
| Min registration days before wedding | 10 days |
| Min age (groom/bride) | 19 years |
| Counseling day | Wednesday only |
| Max counseling participants | 10 couples |

---

## 🚀 Quick Start (JavaScript)

```javascript
// Setup
const API_BASE = 'http://localhost:8080';
const token = localStorage.getItem('token');

// Login
const login = async (username, password) => {
  const res = await fetch(`${API_BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await res.json();
  localStorage.setItem('token', data.token);
  return data;
};

// Get with Auth
const getProfile = async () => {
  const res = await fetch(`${API_BASE}/profile`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return res.json();
};

// Post with Auth
const createRegistration = async (formData) => {
  const res = await fetch(`${API_BASE}/simnikah/pendaftaran/form-baru`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(formData)
  });
  return res.json();
};
```

---

## 📚 Complete Documentation

- 📖 **Full API Docs:** `/docs/API_DOCUMENTATION_FRONTEND.md`
- 📮 **Postman Guide:** `/docs/api/POSTMAN_GUIDE.md`
- 📋 **Example Requests:** `/docs/api/example_requests.json`
- 🗺️ **Map Integration:** `/docs/MAP_INTEGRATION.md`
- ☪️ **Guardian Validation:** `/docs/WALI_NIKAH_VALIDATION.md`

---

## 💡 Tips

1. **Always include token** in authenticated requests
2. **Check token expiry** - Token valid 24h
3. **Use environment variables** for base URL
4. **Handle 401 errors** - Redirect to login
5. **Validate client-side** before sending to API
6. **Use optimized endpoint** - Registration is 60-70% faster now ⚡
7. **Poll notifications** every 30s for updates

---

## 📱 Contact

- **Backend Team:** backend@simnikah.go.id
- **API Issues:** GitHub Issues
- **Documentation:** This file

---

**Version:** 2.0 | **Last Updated:** 30 Oktober 2025

