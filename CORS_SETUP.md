# 🔒 CORS Configuration Guide

## Apa itu CORS?

CORS (Cross-Origin Resource Sharing) adalah mekanisme keamanan browser yang membatasi request HTTP dari origin yang berbeda. Misalnya, jika frontend Anda berjalan di `http://localhost:3000` dan backend di `http://localhost:8080`, browser akan memblokir request karena berbeda origin.

## ✅ Solusi yang Sudah Diterapkan

CORS middleware telah ditambahkan ke aplikasi SimNikah dengan konfigurasi berikut:

### Default Allowed Origins (Development)
Jika `ALLOWED_ORIGINS` tidak diset di environment variable, aplikasi akan mengizinkan origin berikut:
- `http://localhost:3000` (React default)
- `http://localhost:3001`
- `http://localhost:5173` (Vite default)
- `http://localhost:5174`
- `http://localhost:8080`
- `http://127.0.0.1:3000`
- `http://127.0.0.1:5173`

### Allowed Methods
- GET
- POST
- PUT
- DELETE
- PATCH
- OPTIONS

### Allowed Headers
- Origin
- Content-Type
- Accept
- Authorization
- X-Requested-With

### Other Settings
- `AllowCredentials`: true (mendukung cookies & auth headers)
- `MaxAge`: 12 hours (cache preflight requests)

---

## 🚀 Cara Menggunakan

### 1. Development (Default)
Tidak perlu konfigurasi tambahan. Jalankan aplikasi:

```bash
go run main.go
```

### 2. Production dengan Custom Origin
Set environment variable `ALLOWED_ORIGINS`:

**Linux/Mac:**
```bash
export ALLOWED_ORIGINS="https://your-frontend.com,https://app.your-domain.com"
go run main.go
```

**Windows (CMD):**
```cmd
set ALLOWED_ORIGINS=https://your-frontend.com,https://app.your-domain.com
go run main.go
```

**Windows (PowerShell):**
```powershell
$env:ALLOWED_ORIGINS="https://your-frontend.com,https://app.your-domain.com"
go run main.go
```

### 3. Docker
Tambahkan ke `docker-compose.yml`:

```yaml
services:
  simnikah-api:
    environment:
      - ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 4. LeapCell Production
Di LeapCell Dashboard, tambahkan environment variable:
- Key: `ALLOWED_ORIGINS`
- Value: `https://your-frontend-domain.com`

---

## 🧪 Testing CORS

### Test dengan cURL
```bash
curl -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization" \
  -X OPTIONS \
  http://localhost:8080/login -v
```

Response headers yang baik:
```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization, X-Requested-With
Access-Control-Allow-Credentials: true
```

### Test dengan Frontend
**JavaScript (Fetch API):**
```javascript
fetch('http://localhost:8080/health', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  },
  credentials: 'include', // Jika butuh cookies
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('CORS Error:', error));
```

**Axios:**
```javascript
import axios from 'axios';

axios.defaults.withCredentials = true;

axios.get('http://localhost:8080/health')
  .then(response => console.log(response.data))
  .catch(error => console.error('CORS Error:', error));
```

---

## 🔍 Troubleshooting

### ❌ Masalah: "Access to fetch has been blocked by CORS policy"

**Solusi 1:** Pastikan origin frontend ada di allowed origins
```bash
# Check logs saat server start
# Akan muncul: "Info: CORS allowed origins: [http://localhost:3000 ...]"
```

**Solusi 2:** Tambahkan origin frontend ke environment variable
```bash
export ALLOWED_ORIGINS="http://localhost:3000,http://localhost:5173"
```

**Solusi 3:** Restart server setelah update environment variable

### ❌ Masalah: "Preflight request doesn't pass access control check"

**Solusi:** Pastikan frontend mengirim header yang benar
```javascript
fetch('http://localhost:8080/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json', // PENTING!
    'Accept': 'application/json',
  },
  body: JSON.stringify({ username: 'test', password: 'test123' })
})
```

### ❌ Masalah: CORS work di development tapi tidak di production

**Solusi:** Set `ALLOWED_ORIGINS` di production environment
```bash
# LeapCell / Production
ALLOWED_ORIGINS=https://your-production-frontend.com,https://app.your-domain.com
```

---

## 📋 Checklist

- [x] CORS middleware installed
- [x] Default origins configured
- [x] Environment variable support
- [x] AllowCredentials enabled
- [x] Proper headers configured
- [ ] Set production origins
- [ ] Test with real frontend
- [ ] Update frontend API base URL

---

## 📚 Resources

- [MDN CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [gin-contrib/cors](https://github.com/gin-contrib/cors)
- [Understanding CORS](https://web.dev/cross-origin-resource-sharing/)

---

## 💡 Tips

1. **Development**: Gunakan default origins (tidak perlu set `ALLOWED_ORIGINS`)
2. **Production**: Selalu set `ALLOWED_ORIGINS` dengan domain frontend yang spesifik
3. **Security**: Jangan gunakan wildcard (`*`) di production
4. **Testing**: Gunakan browser DevTools Network tab untuk debug CORS errors
5. **Credentials**: Set `withCredentials: true` di frontend jika butuh cookies/auth

---

## 🔧 Konfigurasi Lanjutan

Jika butuh konfigurasi CORS lebih spesifik, edit di `main.go`:

```go
corsConfig := cors.Config{
    AllowOrigins:     getAllowedOrigins(),
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

---

*Last updated: 2024*

