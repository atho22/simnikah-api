# 🔄 Endpoint Update Status Fleksibel

## 📋 Deskripsi

Endpoint ini memungkinkan **Staff, Penghulu, dan Kepala KUA** untuk mengupdate status pendaftaran nikah secara **fleksibel tanpa validasi ketat** pada status sebelumnya.

**Perbedaan dengan endpoint lain:**
- ❌ Endpoint lain: Harus status tertentu dulu baru bisa update (contoh: harus "Menunggu Verifikasi" dulu baru bisa approve)
- ✅ Endpoint ini: Bisa update ke status apapun dari status manapun (fleksibel)

---

## 🔌 Endpoint

```
PUT /simnikah/pendaftaran/:id/update-status
```

**Method:** `PUT`  
**Auth:** Required (JWT Token)  
**Role:** `staff`, `penghulu`, atau `kepala_kua`

---

## 📥 Request

### Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### URL Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | ID pendaftaran nikah |

### Request Body
```json
{
  "status": "Menunggu Bimbingan",
  "catatan": "Status diupdate secara manual oleh staff"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `status` | string | Yes | Status baru yang diinginkan |
| `catatan` | string | No | Catatan untuk perubahan status |

### Status yang Valid
- `"Draft"`
- `"Menunggu Verifikasi"`
- `"Menunggu Pengumpulan Berkas"`
- `"Berkas Diterima"`
- `"Menunggu Penugasan"`
- `"Penghulu Ditugaskan"`
- `"Menunggu Verifikasi Penghulu"`
- `"Menunggu Bimbingan"`
- `"Sudah Bimbingan"`
- `"Selesai"`
- `"Ditolak"`

---

## 📤 Response

### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Status berhasil diupdate",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_sebelumnya": "Menunggu Verifikasi",
    "status_sekarang": "Menunggu Bimbingan",
    "catatan": "Status diupdate secara manual oleh staff",
    "updated_by": "STF1704067201",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Error Responses

#### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "User ID tidak ditemukan"
}
```

#### 403 Forbidden
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Hanya Staff, Penghulu, atau Kepala KUA yang dapat mengupdate status"
}
```

#### 400 Bad Request
```json
{
  "success": false,
  "message": "Status tidak valid",
  "error": "Status yang diizinkan: Draft, Menunggu Verifikasi, Menunggu Pengumpulan Berkas, Berkas Diterima, Menunggu Penugasan, Penghulu Ditugaskan, Menunggu Verifikasi Penghulu, Menunggu Bimbingan, Sudah Bimbingan, Selesai, Ditolak"
}
```

#### 404 Not Found
```json
{
  "success": false,
  "message": "Pendaftaran tidak ditemukan",
  "error": "Pendaftaran dengan ID tersebut tidak ditemukan"
}
```

---

## 📝 Contoh Penggunaan

### Contoh 1: Update dari "Menunggu Verifikasi" ke "Menunggu Bimbingan"
```bash
curl -X PUT "https://api.example.com/simnikah/pendaftaran/1/update-status" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Menunggu Bimbingan",
    "catatan": "Langsung skip ke bimbingan karena berkas sudah lengkap"
  }'
```

### Contoh 2: Update ke "Selesai" dari status apapun
```bash
curl -X PUT "https://api.example.com/simnikah/pendaftaran/1/update-status" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Selesai",
    "catatan": "Nikah sudah dilaksanakan, update manual"
  }'
```

### Contoh 3: Update ke "Ditolak" dari status apapun
```bash
curl -X PUT "https://api.example.com/simnikah/pendaftaran/1/update-status" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Ditolak",
    "catatan": "Data tidak valid, perlu perbaikan"
  }'
```

---

## ⚠️ Catatan Penting

1. **Fleksibilitas**: Endpoint ini **tidak memvalidasi** status sebelumnya. Bisa update dari status apapun ke status apapun.

2. **Notifikasi Otomatis**: Setiap perubahan status akan mengirim notifikasi otomatis ke user (calon pasangan).

3. **Audit Trail**: Semua perubahan dicatat dengan:
   - `updated_by`: User ID yang melakukan update
   - `updated_at`: Timestamp perubahan
   - `disetujui_oleh`: User ID yang melakukan update
   - `disetujui_pada`: Timestamp persetujuan

4. **Role-Based Access**: Hanya Staff, Penghulu, dan Kepala KUA yang bisa menggunakan endpoint ini.

5. **Validasi Status**: Meskipun fleksibel, status yang diinput harus valid (ada dalam daftar status yang diizinkan).

---

## 🔄 Perbandingan dengan Endpoint Lain

| Endpoint | Validasi Status Sebelumnya | Fleksibilitas |
|----------|---------------------------|----------------|
| `/staff/verify-formulir/:id` | ✅ Harus "Menunggu Verifikasi" | ❌ Ketat |
| `/staff/verify-berkas/:id` | ✅ Harus "Menunggu Pengumpulan Berkas" | ❌ Ketat |
| `/penghulu/verify-documents/:id` | ✅ Harus "Menunggu Verifikasi Penghulu" | ❌ Ketat |
| `/pendaftaran/:id/update-status` | ❌ Tidak ada validasi | ✅ **Fleksibel** |

---

## 🎯 Use Cases

### Use Case 1: Skip Tahap
Jika ada kasus khusus dimana perlu skip beberapa tahap:
```
Draft → Menunggu Bimbingan (skip verifikasi, berkas, dll)
```

### Use Case 2: Kembali ke Status Sebelumnya
Jika ada kesalahan dan perlu kembali:
```
Menunggu Bimbingan → Menunggu Verifikasi (kembali ke tahap sebelumnya)
```

### Use Case 3: Update Manual
Jika ada kasus khusus yang tidak ter-cover oleh flow normal:
```
Status apapun → Selesai (jika nikah sudah dilaksanakan)
```

---

## 📚 Contoh Response Lengkap

### Success
```json
{
  "success": true,
  "message": "Status berhasil diupdate",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIK1704067200",
    "status_sebelumnya": "Menunggu Verifikasi",
    "status_sekarang": "Menunggu Bimbingan",
    "catatan": "Langsung skip ke bimbingan karena berkas sudah lengkap",
    "updated_by": "STF1704067201",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Error - Status Tidak Valid
```json
{
  "success": false,
  "message": "Status tidak valid",
  "error": "Status yang diizinkan: Draft, Menunggu Verifikasi, Menunggu Pengumpulan Berkas, Berkas Diterima, Menunggu Penugasan, Penghulu Ditugaskan, Menunggu Verifikasi Penghulu, Menunggu Bimbingan, Sudah Bimbingan, Selesai, Ditolak"
}
```

---

## 🔐 Security

- ✅ JWT Authentication required
- ✅ Role-based authorization (staff, penghulu, kepala_kua)
- ✅ Input validation (status harus valid)
- ✅ Audit trail (mencatat siapa dan kapan update)

---

## 📝 Notes

- Endpoint ini dibuat untuk memberikan fleksibilitas kepada pegawai dalam mengelola status pendaftaran
- Tetap disarankan untuk menggunakan endpoint yang sudah ada (dengan validasi ketat) untuk flow normal
- Gunakan endpoint ini hanya untuk kasus khusus atau kebutuhan administratif

