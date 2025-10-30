# 📚 SimNikah API Documentation

**Dokumentasi lengkap untuk integrasi frontend dengan SimNikah REST API**

---

## 📁 Files Overview

### 1. **API_DOCUMENTATION_FRONTEND.md** ⭐ (24,000+ baris)
**Complete API documentation untuk frontend developers**

📋 **Contents:**
- Getting Started & Authentication
- 50+ API Endpoints dengan contoh lengkap
- Request/Response examples
- Error handling & codes
- Status flow (10-step marriage registration)
- Business rules & validations
- Testing guide (cURL, Fetch, Axios)
- JavaScript integration examples
- Frontend developer notes
- Quick integration checklist

🎯 **Use this for:**
- Understanding all API endpoints
- Request/response format
- Error handling
- Complete integration guide

---

### 2. **QUICK_REFERENCE.md** ⚡
**Quick lookup guide untuk development**

📋 **Contents:**
- Endpoint summary table
- Valid values & constants
- Status codes
- Business rules cheat sheet
- Quick start code snippets

🎯 **Use this for:**
- Quick endpoint lookup
- Valid values reference
- Code snippets
- Daily development

---

### 3. **POSTMAN_GUIDE.md** 📮
**Postman collection & testing guide**

📋 **Contents:**
- Ready-to-use Postman collection
- Environment variables setup
- Testing workflow
- Auto token management
- Pre-request & test scripts

🎯 **Use this for:**
- API testing dengan Postman
- Manual testing
- QA testing
- API exploration

---

### 4. **example_requests.json** 📝
**Example request bodies untuk semua endpoints**

📋 **Contents:**
- Complete request body examples
- Valid values reference
- Authentication examples
- Registration form example
- All operations examples

🎯 **Use this for:**
- Copy-paste request bodies
- Testing different scenarios
- Understanding required fields
- Valid values lookup

---

### 5. **API_PENDAFTARAN_NIKAH.md** (1,074 baris)
**Detailed documentation untuk marriage registration endpoint**

📋 **Contents:**
- Complete marriage registration API
- 10-step workflow explanation
- Guardian validation rules (Syariat Islam)
- Performance optimizations
- Error handling
- Business logic

🎯 **Use this for:**
- Deep dive into registration flow
- Understanding guardian validation
- Performance optimization details
- Specific registration logic

---

## 🚀 Quick Start

### For Frontend Developers

1. **Read First:** `QUICK_REFERENCE.md` (5 min)
2. **Deep Dive:** `API_DOCUMENTATION_FRONTEND.md` (30 min)
3. **Test API:** Import `POSTMAN_GUIDE.md` collection to Postman
4. **Start Coding:** Use examples from `example_requests.json`

### For QA Engineers

1. **Setup:** Import Postman collection from `POSTMAN_GUIDE.md`
2. **Test Flow:** Follow testing workflow in guide
3. **Reference:** Use `QUICK_REFERENCE.md` for valid values
4. **Examples:** Use `example_requests.json` for test data

### For Backend Developers

1. **Overview:** `API_DOCUMENTATION_FRONTEND.md`
2. **Specific Endpoint:** `API_PENDAFTARAN_NIKAH.md`
3. **Testing:** `POSTMAN_GUIDE.md`

---

## 📡 API Overview

### Base URLs
- **Development:** `http://localhost:8080`
- **Production:** `https://your-api.railway.app`

### Authentication
- **Type:** JWT (JSON Web Token)
- **Header:** `Authorization: Bearer <token>`
- **Validity:** 24 hours

### Total Endpoints
- **50+ REST API endpoints**
- **4 user roles** (user_biasa, staff, penghulu, kepala_kua)
- **10-step** marriage registration workflow

### Key Features
- ✅ Authentication & Authorization (JWT + bcrypt)
- ✅ Marriage Registration (Complete form - 1 API call)
- ✅ Staff Management (CRUD operations)
- ✅ Penghulu Management (Assignment & scheduling)
- ✅ Calendar System (Availability & scheduling)
- ✅ Counseling System (Bimbingan perkawinan)
- ✅ Notification System (Auto notifications)
- ✅ Map Integration (Geocoding - 100% FREE)
- ✅ Guardian Validation (Syariat-compliant)
- ⚡ **Performance Optimized** (60-70% faster)

---

## 📊 Performance Metrics

### Marriage Registration API ⚡
- **Before Optimization:** 2500-4000ms
- **After Optimization:** 800-1200ms
- **Improvement:** 60-70% faster

### Optimizations Applied:
1. ✅ Batch database insert (4 queries → 1 query)
2. ✅ Async geocoding (non-blocking)
3. ✅ Async notifications (background)
4. ✅ Single timestamp (consistency + performance)

---

## 🔐 Security Features

- ✅ JWT Authentication (24h validity)
- ✅ Bcrypt Password Hashing (10 rounds)
- ✅ Role-based Authorization (RBAC)
- ✅ SQL Injection Protection (GORM)
- ✅ CORS Configuration (customizable)
- ✅ Rate Limiting (100 req/min)
- ✅ Token Expiration Handling
- ✅ Input Validation

---

## 📋 Marriage Registration Workflow

```
1. Draft (User creates)
   ↓
2. Menunggu Verifikasi (Staff verifies online form)
   ↓
3. Menunggu Pengumpulan Berkas (User brings documents)
   ↓
4. Berkas Diterima (Staff verifies physical docs)
   ↓
5. Menunggu Penugasan (User confirms visit)
   ↓
6. Menunggu Penugasan (Ready for penghulu assignment)
   ↓
7. Penghulu Ditugaskan (Head KUA assigns penghulu)
   ↓
8. Menunggu Verifikasi Penghulu (Penghulu checks)
   ↓
9. Menunggu Bimbingan (Join counseling)
   ↓
10. Sudah Bimbingan (Counseling completed)
    ↓
11. Selesai ✅ (Marriage completed)
```

---

## ☪️ Guardian Validation (Wali Nikah)

### Automatic Validations:

1. ✅ **Wali must be alive** (status = "Hidup")
2. ✅ **Father alive → Wali MUST be father**
3. ✅ **NIK consistency check** (wali = father → same NIK)
4. ✅ **Wali ≠ Bride/Groom** (different NIK)
5. ✅ **Nasab order validation** (according to Islamic law)

### Priority Order (Nasab):
1. Ayah Kandung (biological father)
2. Kakek (grandfather)
3. Saudara Laki-Laki Kandung (full brother)
4. Paman Kandung (uncle)
5. Wali Hakim (if no nasab available)

**Full details:** See `API_PENDAFTARAN_NIKAH.md`

---

## 📊 Business Rules

| Rule | Value | Description |
|------|-------|-------------|
| Max marriages/day (KUA) | 9 | At KUA office |
| Max marriages/day (outside) | Unlimited | External venues |
| Max per penghulu/day | 3 | One penghulu |
| Min gap between marriages | 60 min | Same penghulu |
| Min days before wedding | 10 days | Working days |
| Min age | 19 years | Both parties |
| Counseling day | Wednesday | Only |
| Max counseling participants | 10 couples | Per session |

---

## 🗺️ Map Integration (100% FREE)

**Powered by:** OpenStreetMap Nominatim API

### Features:
- ✅ **Geocoding** - Address → Coordinates
- ✅ **Reverse Geocoding** - Coordinates → Address
- ✅ **Address Search** - Autocomplete
- ✅ **Navigation Links** - Google Maps, Waze, OSM
- ✅ **Caching** - For performance
- ✅ **Async Processing** - Non-blocking

**No API key required!** 🎉

---

## 🔔 Notification System

### Auto Notifications for:
- ✅ Form verification (approved/rejected)
- ✅ Document verification
- ✅ Penghulu assignment
- ✅ Counseling schedule
- ✅ Daily reminders (08:00 AM)
- ✅ Status changes

### Features:
- Role-based notifications
- Read/unread status
- Notification stats
- Mark all as read
- Daily cron job (08:00 AM)

---

## 🧪 Testing

### Using Postman
1. Import collection from `POSTMAN_GUIDE.md`
2. Setup environment variables
3. Run **Login** to get token
4. Test endpoints sequentially

### Using cURL
```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Get profile (with token)
curl http://localhost:8080/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Using JavaScript
```javascript
// Example in QUICK_REFERENCE.md
// Complete examples in API_DOCUMENTATION_FRONTEND.md
```

---

## ⚠️ Common Errors & Solutions

### 401 Unauthorized
- **Cause:** Token missing or invalid
- **Solution:** Login again to get new token

### 403 Forbidden
- **Cause:** Insufficient permissions
- **Solution:** Check role requirements for endpoint

### 400 Bad Request
- **Cause:** Invalid input data
- **Solution:** Check request body format & required fields

### 404 Not Found
- **Cause:** Resource doesn't exist
- **Solution:** Verify ID/URL parameter

---

## 📱 Frontend Integration Checklist

- [ ] Setup axios/fetch with base URL
- [ ] Implement token storage (localStorage)
- [ ] Create authentication flow
- [ ] Handle 401 responses (redirect to login)
- [ ] Implement registration form (10 sections)
- [ ] Add client-side validation
- [ ] Display status flow with progress indicator
- [ ] Implement calendar views
- [ ] Add notification system
- [ ] Test all user flows
- [ ] Handle errors gracefully
- [ ] Add loading states
- [ ] Implement role-based UI

---

## 📞 Support & Contact

### Documentation
- 📖 Full API Docs: `API_DOCUMENTATION_FRONTEND.md`
- ⚡ Quick Ref: `QUICK_REFERENCE.md`
- 📮 Postman: `POSTMAN_GUIDE.md`
- 📋 Examples: `example_requests.json`

### Contact
- **Backend Team:** backend@simnikah.go.id
- **API Issues:** GitHub Issues
- **Emergency:** +62 xxx xxxx xxxx

### Links
- **GitHub:** https://github.com/your-org/simpadu
- **Production API:** https://your-api.railway.app
- **Frontend Demo:** https://your-frontend.vercel.app

---

## 📈 Version History

### Version 2.0 (30 Oktober 2025)
- ⚡ Performance optimization (60-70% faster)
- ☪️ Guardian validation enhancement (5 rules)
- 🗺️ Map integration (FREE geocoding)
- 🔔 Improved notification system
- 📚 Complete API documentation for frontend
- 📮 Postman collection
- 📋 Example request bodies

### Version 1.0 (Initial Release)
- ✅ Core API functionality
- ✅ Authentication & Authorization
- ✅ Basic marriage registration
- ✅ Staff & penghulu management
- ✅ Calendar system
- ✅ Counseling system

---

## 🎯 Next Steps

1. **Frontend Developers:**
   - Read `QUICK_REFERENCE.md`
   - Study `API_DOCUMENTATION_FRONTEND.md`
   - Start integration

2. **QA Engineers:**
   - Import Postman collection
   - Test all flows
   - Report bugs

3. **Product Managers:**
   - Review business rules
   - Verify workflow
   - Approve features

---

## 🙏 Credits

**Made with ❤️ for SimNikah Project**

- **Backend Team:** Go, Gin, GORM, MySQL
- **Geocoding:** OpenStreetMap Nominatim API
- **Deployment:** Railway
- **Version Control:** Git & GitHub

---

**Last Updated:** 30 Oktober 2025  
**Version:** 2.0  
**Total Documentation:** 30,000+ lines

---

Happy Coding! 🚀

