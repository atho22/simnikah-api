# 🔧 Troubleshooting: 403 Forbidden Error di Frontend Vercel

## 🚨 Masalah

Frontend di `https://kua-ku.vercel.app` mengalami error **403 Forbidden** saat login:
```
POST https://kua-ku.vercel.app/api/proxy/login 403 (Forbidden)
Error: Akses ditolak oleh server API. Pastikan environment variable NEXT_PUBLIC_API_URL sudah di-set dengan benar di Vercel.
```

## 🔍 Analisis Masalah

Error ini **BUKAN** dari backend API, tapi dari **Next.js API Route Proxy** di frontend. Frontend menggunakan proxy route (`/api/proxy/login`) yang mem-forward request ke backend.

### Penyebab Umum:

1. ❌ **Environment Variable `NEXT_PUBLIC_API_URL` tidak di-set di Vercel**
2. ❌ **Next.js API route proxy memblokir request karena validasi gagal**
3. ❌ **Backend URL tidak valid atau tidak bisa diakses dari Vercel**

---

## ✅ Solusi

### Step 1: Set Environment Variable di Vercel

1. **Buka Vercel Dashboard**
   - Login ke [vercel.com](https://vercel.com)
   - Pilih project `kua-ku` (atau nama project frontend kamu)

2. **Tambahkan Environment Variable**
   - Klik tab **"Settings"**
   - Klik **"Environment Variables"** di sidebar kiri
   - Klik **"Add New"**
   - **Key:** `NEXT_PUBLIC_API_URL`
   - **Value:** `https://your-backend.up.railway.app` (ganti dengan URL backend Railway kamu)
   - **Environment:** Pilih semua (Production, Preview, Development)
   - Klik **"Save"**

3. **Redeploy Frontend**
   - Setelah save environment variable, Vercel akan otomatis trigger redeploy
   - Atau manual: Klik tab **"Deployments"** → Klik **"Redeploy"** pada deployment terbaru
   - Tunggu 2-3 menit sampai deployment selesai

4. **Test**
   - Buka `https://kua-ku.vercel.app`
   - Coba login lagi
   - Error 403 seharusnya sudah hilang ✅

### Step 2: Verifikasi Backend URL

Pastikan backend URL yang di-set di `NEXT_PUBLIC_API_URL` adalah URL yang benar dan bisa diakses:

```bash
# Test backend health check
curl https://your-backend.up.railway.app/health

# Expected response:
# {"status":"healthy","service":"SimNikah API","timestamp":"..."}
```

Jika backend tidak bisa diakses, cek:
- ✅ Backend sudah deployed di Railway?
- ✅ Backend service status "Running"?
- ✅ Public domain sudah di-generate?

### Step 3: Cek Next.js API Route Proxy

Jika masih error setelah set environment variable, kemungkinan ada masalah di Next.js API route proxy. Cek file:

**File:** `pages/api/proxy/login.ts` atau `app/api/proxy/login/route.ts`

**Pastikan:**
1. ✅ Environment variable di-read dengan benar:
   ```typescript
   const apiUrl = process.env.NEXT_PUBLIC_API_URL;
   if (!apiUrl) {
     return res.status(500).json({ error: 'API URL not configured' });
   }
   ```

2. ✅ Request di-forward dengan benar:
   ```typescript
   const response = await fetch(`${apiUrl}/login`, {
     method: 'POST',
     headers: {
       'Content-Type': 'application/json',
     },
     body: JSON.stringify(req.body),
   });
   ```

3. ✅ Error handling yang benar:
   ```typescript
   if (!response.ok) {
     const error = await response.json();
     return res.status(response.status).json(error);
   }
   ```

---

## 🔍 Debugging

### 1. Cek Environment Variable di Vercel

**Vercel Dashboard → Settings → Environment Variables**

Pastikan:
- ✅ `NEXT_PUBLIC_API_URL` sudah ada
- ✅ Value adalah URL backend yang benar (contoh: `https://simnikah-production-xxxx.up.railway.app`)
- ✅ Environment sudah di-set untuk Production, Preview, dan Development

### 2. Cek Logs Vercel

**Vercel Dashboard → Deployments → Klik deployment terbaru → Tab "Functions"**

Cari error logs dari API route proxy:
- ❌ `API URL not configured` → Environment variable tidak ter-load
- ❌ `Failed to fetch` → Backend tidak bisa diakses
- ❌ `403 Forbidden` → Backend mengembalikan 403 (bisa karena CORS atau auth)

### 3. Test Backend Langsung

Test backend API langsung dari browser atau Postman:

```bash
# Test login endpoint
curl -X POST https://your-backend.up.railway.app/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

**Expected Response (200):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "token": "...",
  "user": {...}
}
```

Jika backend mengembalikan 403, berarti masalahnya di backend (bisa CORS atau rate limiting).

### 4. Cek CORS di Backend

Pastikan backend sudah mengizinkan origin frontend:

**Railway Dashboard → Variables → `ALLOWED_ORIGINS`**

Value harus include:
```
https://kua-ku.vercel.app
```

Lihat dokumentasi: [FIX_CORS_VERCEL.md](FIX_CORS_VERCEL.md)

---

## 📋 Checklist Troubleshooting

### Frontend (Vercel)
- [ ] Environment variable `NEXT_PUBLIC_API_URL` sudah di-set di Vercel
- [ ] Value adalah URL backend yang benar (dengan `https://`)
- [ ] Environment variable di-set untuk Production, Preview, dan Development
- [ ] Frontend sudah di-redeploy setelah set environment variable
- [ ] Next.js API route proxy tidak memblokir request

### Backend (Railway)
- [ ] Backend sudah deployed dan running
- [ ] Public domain sudah di-generate
- [ ] Health check endpoint (`/health`) bisa diakses
- [ ] CORS sudah dikonfigurasi dengan `ALLOWED_ORIGINS` include `https://kua-ku.vercel.app`
- [ ] Login endpoint (`/login`) bisa diakses langsung

### Network
- [ ] Backend URL bisa diakses dari browser
- [ ] Tidak ada firewall atau network restriction
- [ ] SSL certificate valid (https://)

---

## 🎯 Quick Fix

**Jika error: "Akses ditolak oleh server API. Pastikan environment variable NEXT_PUBLIC_API_URL sudah di-set dengan benar di Vercel."**

1. **Vercel Dashboard** → Settings → Environment Variables
2. **Add New:**
   - Key: `NEXT_PUBLIC_API_URL`
   - Value: `https://your-backend.up.railway.app` (ganti dengan URL backend kamu)
3. **Save** → Tunggu redeploy otomatis
4. **Test** → Refresh frontend dan coba login lagi

---

## 🔗 Related Documentation

- **[Fix CORS Error](FIX_CORS_VERCEL.md)** - Troubleshooting CORS
- **[CORS Setup Guide](features/CORS_SETUP.md)** - Dokumentasi lengkap CORS
- **[Railway Deployment](tutorials/TUTORIAL_DEPLOY_RAILWAY.md)** - Tutorial deploy backend

---

## 💡 Tips

1. **Environment Variable Naming:**
   - Next.js hanya expose environment variable yang dimulai dengan `NEXT_PUBLIC_` ke client-side
   - Untuk API routes, bisa pakai `NEXT_PUBLIC_` atau tanpa prefix (server-side only)

2. **Multiple Environments:**
   - Set environment variable untuk Production, Preview, dan Development
   - Bisa pakai value berbeda untuk setiap environment

3. **Debugging:**
   - Gunakan Vercel Function Logs untuk debug API route proxy
   - Gunakan Railway Logs untuk debug backend
   - Gunakan Browser DevTools Network tab untuk lihat request/response

---

**Last Updated:** November 2024

