# 🔧 Troubleshooting: Error Update Status

## ❌ Error yang Sering Terjadi

### Error 1: Field validation for 'Status' failed on the 'required' tag

**Penyebab:**
- Field JSON menggunakan huruf besar `"Status"` (salah)
- Field `status` kosong atau tidak ada
- Format JSON tidak valid

**Solusi:**

#### ✅ Request yang Benar:
```json
{
  "status": "Selesai",
  "catatan": "Nikah telah dilaksanakan"
}
```

#### ❌ Request yang Salah:
```json
{
  "Status": "Selesai",  // ❌ Huruf besar
  "catatan": "Nikah telah dilaksanakan"
}
```

```json
{
  "status": "",  // ❌ Kosong
  "catatan": "Nikah telah dilaksanakan"
}
```

---

### Error 2: Status tidak valid

**Penyebab:**
- Status yang diinput tidak ada dalam daftar status yang valid
- Typo pada nama status

**Solusi:**

Gunakan status yang valid:
- `"Draft"`
- `"Menunggu Verifikasi"`
- `"Menunggu Pengumpulan Berkas"`
- `"Berkas Diterima"`
- `"Menunggu Bimbingan"`
- `"Sudah Bimbingan"`
- `"Selesai"`
- `"Ditolak"`

**Catatan:** Status terkait penghulu tidak bisa diupdate via endpoint fleksibel.

---

### Error 3: Akses ditolak - Status terkait penghulu

**Penyebab:**
- Mencoba update status terkait penghulu via endpoint fleksibel
- Bukan Kepala KUA yang mencoba update

**Solusi:**

Untuk status terkait penghulu, gunakan endpoint khusus:
```json
POST /simnikah/pendaftaran/:id/assign-penghulu
{
  "penghulu_id": 1
}
```

---

## ✅ Contoh Request yang Benar

### TypeScript/JavaScript
```typescript
const response = await fetch('/simnikah/pendaftaran/1/update-status', {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    status: 'Selesai',  // ✅ Huruf kecil, tidak kosong
    catatan: 'Nikah telah dilaksanakan'
  })
});
```

### cURL
```bash
curl -X PUT "https://api.example.com/simnikah/pendaftaran/1/update-status" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Selesai",
    "catatan": "Nikah telah dilaksanakan"
  }'
```

### Postman
- Method: `PUT`
- URL: `/simnikah/pendaftaran/1/update-status`
- Headers:
  - `Authorization: Bearer <token>`
  - `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "status": "Selesai",
  "catatan": "Nikah telah dilaksanakan"
}
```

---

## 🔍 Debugging Tips

1. **Cek Format JSON:**
   - Pastikan menggunakan `"status"` (huruf kecil)
   - Pastikan field tidak kosong
   - Validasi JSON di [jsonlint.com](https://jsonlint.com)

2. **Cek Status Value:**
   - Pastikan status sesuai dengan daftar yang valid
   - Case-sensitive: `"Selesai"` bukan `"selesai"`

3. **Cek Token:**
   - Pastikan token masih valid (tidak expired)
   - Pastikan token di header dengan format: `Bearer <token>`

4. **Cek Role:**
   - Pastikan user memiliki role yang sesuai (staff, penghulu, atau kepala_kua)

---

## 📝 Checklist Sebelum Update Status

- [ ] Field JSON menggunakan `"status"` (huruf kecil)
- [ ] Field `status` tidak kosong
- [ ] Status sesuai dengan daftar yang valid
- [ ] Token masih valid
- [ ] User memiliki role yang sesuai
- [ ] Format JSON valid
- [ ] Header Authorization sudah diset

---

**Last Updated:** 2024-01-15

