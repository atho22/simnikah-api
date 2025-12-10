# Panduan Upload Foto Profil — Frontend

Dokumentasi ini menjelaskan cara frontend mengunggah foto profil pengguna ke SimNikah API.

## Ringkasan Singkat
- Endpoint: `POST /upload-photo`
- Authentication: Required — header `Authorization: Bearer <JWT_TOKEN>`
- Form field: `photo` (harus persis `photo`)
- Content-Type: `multipart/form-data` (jangan set manual saat menggunakan `FormData`)
- Validasi server: max 5MB; allowed MIME types: `image/jpeg`, `image/png`, `image/jpg`, `image/webp`
- Penyimpanan: file diupload ke ImgBB, backend menyimpan URL di field `profile_photo` pada tabel `Users`
- Untuk mengambil URL foto: `GET /profile` (mengembalikan `profile_photo`)

---

## Contoh Penggunaan

### cURL
```bash
TOKEN="<JWT_TOKEN>"
curl -v -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@/path/to/image.jpg"
```

### PowerShell (curl)
```powershell
$token = "<JWT_TOKEN>"
curl -v -X POST "http://localhost:8080/upload-photo" `
  -H "Authorization: Bearer $token" `
  -F "photo=@C:\path\to\image.jpg"
```

> Catatan: Jika `Invoke-RestMethod` tidak cukup fleksibel untuk multipart, gunakan curl atau HttpClient untuk PowerShell.

### Fetch (Browser)
```javascript
const file = document.querySelector('#fileInput').files[0];
const form = new FormData();
form.append('photo', file);

fetch('/upload-photo', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` }, // jangan set Content-Type
  body: form
})
  .then(async res => {
    const json = await res.json();
    if (!res.ok) throw json;
    return json;
  })
  .then(data => console.log('uploaded', data))
  .catch(err => console.error('upload error', err));
```

### Axios (dengan progress)
```javascript
import axios from 'axios';

const upload = async (file, token, onProgress) => {
  const formData = new FormData();
  formData.append('photo', file);

  const res = await axios.post('/upload-photo', formData, {
    headers: { Authorization: `Bearer ${token}` },
    onUploadProgress: (progressEvent) => {
      if (onProgress) {
        const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        onProgress(percent);
      }
    }
  });
  return res.data;
};
```

### React (komponen sederhana dengan preview)
```jsx
import React, { useState } from 'react';
import axios from 'axios';

function ProfilePhotoUploader({ token, onDone }) {
  const [file, setFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState(null);

  const handleSelect = (e) => {
    const f = e.target.files[0];
    if (!f) return;
    setFile(f);
    setPreview(URL.createObjectURL(f));
    setError(null);
  };

  const handleUpload = async () => {
    if (!file) return;
    const fd = new FormData();
    fd.append('photo', file);

    try {
      setProgress(0);
      const res = await axios.post('/upload-photo', fd, {
        headers: { Authorization: `Bearer ${token}` },
        onUploadProgress: (e) => setProgress(Math.round((e.loaded * 100) / e.total))
      });
      onDone && onDone(res.data);
    } catch (err) {
      setError(err.response?.data || err.message);
    }
  };

  return (
    <div>
      <input type="file" accept="image/*" onChange={handleSelect} />
      {preview && <img src={preview} alt="preview" style={{ width: 120 }} />}
      {progress > 0 && <div>Uploading: {progress}%</div>}
      {error && <div style={{ color: 'red' }}>{JSON.stringify(error)}</div>}
      <button onClick={handleUpload} disabled={!file}>Upload</button>
    </div>
  );
}

export default ProfilePhotoUploader;
```

---

## Proxy Dev Server (Vite) jika menggunakan prefix `/api/proxy`
Jika frontend memanggil endpoint lewat prefix (contoh: `/api/proxy/upload-photo`), pastikan proxy meredispatch path ke backend dan melakukan *rewrite* agar backend menerima `/upload-photo`.

Contoh `vite.config.js`:
```js
export default defineConfig({
  server: {
    proxy: {
      '/api/proxy': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api\/proxy/, '')
      }
    }
  }
});
```
Dengan konfigurasi di atas, frontend mengirim ke `/api/proxy/upload-photo` → diteruskan ke `http://localhost:8080/upload-photo`.

---

## Response Sukses
```json
{
  "success": true,
  "message": "Foto profil berhasil diupload",
  "data": {
    "profile_photo": "https://...",
    "user_id": "USR...",
    "username": "john_doe"
  }
}
```

## Error Umum & Penjelasan
- 400 Bad Request
  - Penyebab: field `photo` hilang, file terlalu besar (>5MB), atau MIME type tidak diizinkan.
  - Perbaikan: pastikan `FormData` berisi `photo`, compress file atau pilih file lain.

- 401 Unauthorized
  - Penyebab: header `Authorization` tidak ada/invalid/expired.
  - Perbaikan: sertakan token valid.

- 404 Not Found
  - Penyebab: memanggil URL yang salah (contoh `/profile/photo` atau proxy tidak melakukan rewrite).
  - Perbaikan: gunakan `/upload-photo` atau perbaiki proxy atau minta backend menambahkan alias route.

- 500 Internal Server Error
  - Penyebab: kegagalan upload ke ImgBB atau error DB.
  - Perbaikan: lihat backend logs / Railway logs, cek `IMGBB_API_KEY`.

---

## Troubleshooting Cepat
1. Tes dengan `curl` untuk memastikan backend menerima upload.
2. Periksa Network tab di browser DevTools:
   - Apakah `Authorization` header ada?
   - Apakah `Content-Type` adalah multipart boundary yang dihasilkan browser?
   - Apakah payload menunjukkan field `photo`?
3. Jika response body kosong atau JSON parsing error (`EOF`), periksa apakah proxy atau client men-strip body.
4. Setelah upload berhasil, panggil `GET /profile` untuk mendapatkan `profile_photo` terbaru.

---

## Opsional: Menambahkan Alias Route di Backend
Jika tim frontend sudah menggunakan URL lama seperti `/profile/photo`, kamu bisa menambahkan alias di `cmd/api/main.go`:

```go
r.POST("/profile/photo", middleware.AuthMiddleware(), authHandler.UploadProfilePhoto)
```

Saya bisa tambahkan dan commit perubahan ini jika kamu mau agar frontend tidak perlu diubah.

---

## Environment
- `IMGBB_API_KEY` (opsional) — tambahkan di Railway environment variables jika ingin menggunakan API key ImgBB.

---

Jika mau, saya bisa:
- Buat Postman collection untuk endpoint ini;
- Tambahkan route alias `/profile/photo` di backend dan push commit; atau
- Patch frontend file upload yang kamu tunjukkan.

Pilih aksi yang kamu inginkan dan saya lanjutkan.
