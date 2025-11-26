# 🔧 Fix CORS Error untuk Frontend Vercel

## 🚨 Masalah
Frontend di `https://kua-ku.vercel.app` mengalami CORS error saat mengakses backend API.

## ✅ Solusi Cepat

### Opsi 1: Set Environment Variable di Railway (RECOMMENDED)

1. **Buka Railway Dashboard**
   - Login ke [railway.app](https://railway.app)
   - Pilih project backend kamu

2. **Tambahkan Environment Variable**
   - Klik tab **"Variables"**
   - Klik **"New Variable"** (atau edit `ALLOWED_ORIGINS` jika sudah ada)
   - **Key:** `ALLOWED_ORIGINS`
   - **Value:** `https://kua-ku.vercel.app`
   - Klik **"Save"**

3. **Tunggu Redeploy**
   - Railway akan otomatis redeploy setelah save
   - Tunggu 2-3 menit sampai deployment selesai
   - Cek status di tab **"Deployments"**

4. **Test**
   - Buka frontend di `https://kua-ku.vercel.app`
   - Coba login atau akses API
   - CORS error seharusnya sudah hilang ✅

### Opsi 2: Multiple Domains (Jika ada staging/production)

Jika kamu punya multiple frontend domains, pisahkan dengan koma:

```
https://kua-ku.vercel.app,https://kua-ku-staging.vercel.app,https://your-backend.up.railway.app
```

**Catatan Penting:**
- ✅ Gunakan `https://` (bukan `http://`)
- ✅ Jangan ada spasi setelah koma
- ✅ Include backend domain juga (untuk testing dari browser)

---

## 🧪 Verifikasi CORS

### Test dengan Browser DevTools

1. Buka `https://kua-ku.vercel.app`
2. Buka **DevTools** (F12) → Tab **Network**
3. Coba login atau akses API
4. Cek response headers, harus ada:
   ```
   Access-Control-Allow-Origin: https://kua-ku.vercel.app
   Access-Control-Allow-Credentials: true
   ```

### Test dengan cURL

```bash
curl -H "Origin: https://kua-ku.vercel.app" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization" \
  -X OPTIONS \
  https://your-backend.up.railway.app/login -v
```

**Expected Response Headers:**
```
Access-Control-Allow-Origin: https://kua-ku.vercel.app
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization, X-Requested-With
Access-Control-Allow-Credentials: true
```

---

## 🔍 Troubleshooting

### ❌ Masih Error CORS Setelah Set Environment Variable

**Solusi 1: Cek Logs Railway**
1. Railway Dashboard → Tab **"Deployments"**
2. Klik deployment terbaru
3. Scroll ke bawah, cari log:
   ```
   Info: CORS allowed origins: [https://kua-ku.vercel.app]
   ```
4. Jika tidak muncul, berarti environment variable belum ter-load

**Solusi 2: Restart Deployment**
1. Railway Dashboard → Tab **"Deployments"**
2. Klik **"Redeploy"** atau **"Deploy"**
3. Tunggu sampai selesai

**Solusi 3: Cek Format Environment Variable**
- ✅ Benar: `https://kua-ku.vercel.app`
- ❌ Salah: `http://kua-ku.vercel.app` (pakai http)
- ❌ Salah: `kua-ku.vercel.app` (tanpa https://)
- ❌ Salah: `https://kua-ku.vercel.app/` (ada trailing slash)

### ❌ CORS Error di Browser Console

**Error:** `Access to fetch at '...' from origin 'https://kua-ku.vercel.app' has been blocked by CORS policy`

**Solusi:**
1. Pastikan `ALLOWED_ORIGINS` sudah di-set di Railway
2. Pastikan format benar: `https://kua-ku.vercel.app` (tanpa trailing slash)
3. Clear browser cache dan hard refresh (Ctrl+Shift+R)
4. Cek Network tab di DevTools untuk melihat response headers

### ❌ Preflight Request Fails

**Error:** `Preflight request doesn't pass access control check`

**Solusi:**
1. Pastikan frontend mengirim header yang benar:
   ```javascript
   fetch('https://your-backend.up.railway.app/api', {
     method: 'POST',
     headers: {
       'Content-Type': 'application/json', // PENTING!
       'Authorization': `Bearer ${token}`
     },
     credentials: 'include' // Jika butuh cookies
   })
   ```

2. Pastikan backend mengizinkan method dan headers yang digunakan

---

## 📋 Checklist

- [ ] Environment variable `ALLOWED_ORIGINS` sudah di-set di Railway
- [ ] Value format benar: `https://kua-ku.vercel.app` (tanpa trailing slash)
- [ ] Railway sudah selesai redeploy (cek di tab Deployments)
- [ ] Logs menunjukkan: `Info: CORS allowed origins: [https://kua-ku.vercel.app]`
- [ ] Test dari browser → CORS error sudah hilang
- [ ] Response headers menunjukkan `Access-Control-Allow-Origin: https://kua-ku.vercel.app`

---

## 🎯 Quick Reference

**Railway Environment Variable:**
```
Key: ALLOWED_ORIGINS
Value: https://kua-ku.vercel.app
```

**Multiple Domains:**
```
Key: ALLOWED_ORIGINS
Value: https://kua-ku.vercel.app,https://staging.vercel.app,https://your-backend.up.railway.app
```

**Backend URL (contoh):**
```
https://simnikah-production-xxxx.up.railway.app
```

---

## 📚 Dokumentasi Lengkap

- **[CORS Setup Guide](features/CORS_SETUP.md)** - Dokumentasi lengkap CORS
- **[Railway Deployment](tutorials/TUTORIAL_DEPLOY_RAILWAY.md)** - Tutorial deploy Railway

---

**Last Updated:** November 2024


