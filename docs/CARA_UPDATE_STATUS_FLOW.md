# 🔄 Cara Update Status Flow Pendaftaran Nikah

## ⚠️ Catatan Penting

**Error Validasi Field 'Status':**
Jika Anda mendapat error `Field validation for 'Status' failed on the 'required' tag`, pastikan:
1. Field JSON menggunakan huruf kecil: `"status"` (bukan `"Status"`)
2. Field `status` tidak kosong
3. Format JSON benar

**Contoh Request yang Benar:**
```json
{
  "status": "Selesai",
  "catatan": "Nikah telah dilaksanakan"
}
```

---

## 📋 Daftar Isi
1. [Metode Update Status](#metode-update-status)
2. [Endpoint Fleksibel (Recommended)](#endpoint-fleksibel-recommended)
3. [Endpoint Khusus per Tahap](#endpoint-khusus-per-tahap)
4. [Contoh Implementasi](#contoh-implementasi)
5. [Best Practices](#best-practices)
6. [Flow Sederhana](#flow-sederhana)

---

## 🎯 Metode Update Status

Ada **2 cara** untuk mengupdate status pendaftaran nikah:

### 1. **Endpoint Fleksibel** (Baru - Recommended)
- ✅ Bisa update dari status apapun ke status apapun
- ✅ Tidak perlu validasi status sebelumnya
- ✅ Cocok untuk kasus khusus atau kebutuhan administratif
- ❌ Tidak bisa update status terkait penghulu (hanya Kepala KUA)

### 2. **Endpoint Khusus per Tahap**
- ✅ Validasi ketat sesuai flow
- ✅ Business logic sudah terintegrasi
- ✅ Notifikasi otomatis
- ❌ Harus mengikuti flow yang ketat

---

## 🔌 Endpoint Fleksibel (Recommended)

### PUT `/simnikah/pendaftaran/:id/update-status`

**Auth:** ✅ Required  
**Role:** `staff`, `penghulu`, `kepala_kua`

#### Request
```http
PUT /simnikah/pendaftaran/:id/update-status
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "status": "Menunggu Bimbingan",
  "catatan": "Status diupdate secara manual oleh staff"
}
```

#### Response Success
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

#### Status yang BISA diupdate:
- ✅ `"Draft"`
- ✅ `"Menunggu Verifikasi"`
- ✅ `"Menunggu Pengumpulan Berkas"`
- ✅ `"Berkas Diterima"`
- ✅ `"Menunggu Bimbingan"`
- ✅ `"Sudah Bimbingan"`
- ✅ `"Selesai"`
- ✅ `"Ditolak"`

#### Status yang TIDAK BISA diupdate (hanya via endpoint khusus):
- ❌ `"Menunggu Penugasan"` → Gunakan `POST /simnikah/pendaftaran/:id/assign-penghulu`
- ❌ `"Penghulu Ditugaskan"` → Otomatis saat assign penghulu
- ❌ `"Menunggu Verifikasi Penghulu"` → Otomatis saat assign penghulu

---

## 📝 Endpoint Khusus per Tahap

### 1. Verifikasi Formulir Online (Staff)

**POST** `/simnikah/staff/verify-formulir/:id`

**Role:** `staff`  
**Status Sebelumnya:** Harus `"Menunggu Verifikasi"`

**Request:**
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir sudah lengkap dan valid"
}
```

**Status Options:**
- `"Formulir Disetujui"` → Next: `"Menunggu Pengumpulan Berkas"`
- `"Ditolak"` → Final status

---

### 2. Verifikasi Berkas Fisik (Staff)

**POST** `/simnikah/staff/verify-berkas/:id`

**Role:** `staff`  
**Status Sebelumnya:** Harus `"Menunggu Pengumpulan Berkas"`

**Request:**
```json
{
  "status": "Berkas Diterima",
  "catatan": "Berkas sudah lengkap dan sesuai"
}
```

**Status Options:**
- `"Berkas Diterima"` → Next: `"Berkas Diterima"`
- `"Ditolak"` → Final status

---

### 3. Konfirmasi Kunjungan (User)

**POST** `/simnikah/pendaftaran/:id/mark-visited`

**Role:** `user_biasa`  
**Status Sebelumnya:** Harus `"Berkas Diterima"`

**Response:**
```json
{
  "success": true,
  "message": "Status berhasil diupdate ke Menunggu Penugasan"
}
```

**Next Status:** `"Menunggu Penugasan"`

---

### 4. Assign Penghulu (Kepala KUA)

**POST** `/simnikah/pendaftaran/:id/assign-penghulu`

**Role:** `kepala_kua`  
**Status Sebelumnya:** Harus `"Menunggu Penugasan"`

**Request:**
```json
{
  "penghulu_id": 1
}
```

**Response:**
```json
{
  "message": "Penghulu berhasil diassign",
  "data": {
    "pendaftaran_id": 1,
    "penghulu_id": 1,
    "penghulu_nama": "H. Abdul Rahman, S.Ag",
    "assigned_by": "KKUA1704067203",
    "assigned_at": "2024-01-15T10:30:00Z"
  }
}
```

**Auto Transition:**
- `"Penghulu Ditugaskan"` → `"Menunggu Verifikasi Penghulu"` (otomatis)

---

### 5. Verifikasi Dokumen (Penghulu)

**POST** `/simnikah/penghulu/verify-documents/:id`

**Role:** `penghulu`  
**Status Sebelumnya:** Harus `"Menunggu Verifikasi Penghulu"`

**Request:**
```json
{
  "status": "Menunggu Pelaksanaan",
  "catatan": "Dokumen sudah lengkap dan valid"
}
```

**Status Options:**
- `"Menunggu Pelaksanaan"` → Next: `"Menunggu Bimbingan"`
- `"Ditolak"` → Final status

---

### 6. Complete Bimbingan (Staff/Kepala KUA)

**PUT** `/simnikah/pendaftaran/:id/complete-bimbingan`

**Role:** `staff`, `kepala_kua`  
**Status Sebelumnya:** Harus `"Menunggu Bimbingan"`

**Response:**
```json
{
  "message": "Bimbingan perkawinan berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Sudah Bimbingan"
  }
}
```

**Next Status:** `"Sudah Bimbingan"`

---

### 7. Complete Nikah (Staff/Kepala KUA)

**PUT** `/simnikah/pendaftaran/:id/complete-nikah`

**Role:** `staff`, `kepala_kua`  
**Status Sebelumnya:** Harus `"Sudah Bimbingan"`

**Response:**
```json
{
  "message": "Nikah berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Selesai",
    "completed_at": "2024-01-15T10:30:00Z"
  }
}
```

**Next Status:** `"Selesai"` ✅ (Final)

---

## 💻 Contoh Implementasi

### Contoh 1: Update Status Fleksibel (React/TypeScript)

```typescript
// services/pendaftaran.ts
import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export interface UpdateStatusRequest {
  status: string;
  catatan?: string;
}

export const updateStatusPendaftaran = async (
  pendaftaranId: number,
  data: UpdateStatusRequest,
  token: string
) => {
  const response = await axios.put(
    `${API_BASE_URL}/simnikah/pendaftaran/${pendaftaranId}/update-status`,
    data,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
  return response.data;
};

// components/UpdateStatusButton.tsx
import { updateStatusPendaftaran } from '@/services/pendaftaran';

const UpdateStatusButton = ({ pendaftaranId, currentStatus }) => {
  const handleUpdate = async () => {
    try {
      const token = localStorage.getItem('token');
      const result = await updateStatusPendaftaran(
        pendaftaranId,
        {
          status: 'Menunggu Bimbingan',
          catatan: 'Skip tahap verifikasi karena kasus khusus'
        },
        token
      );
      
      console.log('Status updated:', result);
      // Refresh data atau update UI
    } catch (error) {
      console.error('Failed to update status:', error);
    }
  };

  return (
    <button onClick={handleUpdate}>
      Update ke Menunggu Bimbingan
    </button>
  );
};
```

---

### Contoh 2: Update Status dengan Endpoint Khusus (Flow Normal)

```typescript
// services/staff.ts
export const verifyFormulir = async (
  pendaftaranId: number,
  data: { status: string; catatan?: string },
  token: string
) => {
  const response = await axios.post(
    `${API_BASE_URL}/simnikah/staff/verify-formulir/${pendaftaranId}`,
    data,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
  return response.data;
};

// Usage
await verifyFormulir(
  1,
  {
    status: 'Formulir Disetujui',
    catatan: 'Formulir sudah lengkap'
  },
  token
);
```

---

### Contoh 3: Assign Penghulu (Kepala KUA)

```typescript
// services/kepalaKua.ts
export const assignPenghulu = async (
  pendaftaranId: number,
  penghuluId: number,
  token: string
) => {
  const response = await axios.post(
    `${API_BASE_URL}/simnikah/pendaftaran/${pendaftaranId}/assign-penghulu`,
    { penghulu_id: penghuluId },
    {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
  return response.data;
};

// Usage
await assignPenghulu(1, 2, token);
// Status akan otomatis: "Menunggu Penugasan" → "Penghulu Ditugaskan" → "Menunggu Verifikasi Penghulu"
```

---

## 📊 Tabel Perbandingan

| Aspek | Endpoint Fleksibel | Endpoint Khusus |
|-------|-------------------|-----------------|
| **Fleksibilitas** | ✅ Tinggi | ❌ Ketat |
| **Validasi Status** | ❌ Tidak ada | ✅ Ada |
| **Business Logic** | ❌ Manual | ✅ Terintegrasi |
| **Notifikasi** | ✅ Otomatis | ✅ Otomatis |
| **Use Case** | Kasus khusus, skip tahap | Flow normal |
| **Role** | Staff, Penghulu, Kepala KUA | Sesuai tahap |

---

## 🎯 Kapan Menggunakan Masing-Masing?

### Gunakan Endpoint Fleksibel jika:
- ✅ Perlu skip beberapa tahap (kasus khusus)
- ✅ Perlu kembali ke status sebelumnya (koreksi)
- ✅ Update manual untuk kebutuhan administratif
- ✅ Testing atau debugging

### Gunakan Endpoint Khusus jika:
- ✅ Flow normal berjalan sesuai prosedur
- ✅ Perlu validasi ketat sesuai business rules
- ✅ Ingin memastikan data integrity
- ✅ Production environment

---

## ⚠️ Catatan Penting

### 1. Status Terkait Penghulu
Status berikut **HANYA** bisa diubah oleh **Kepala KUA** melalui endpoint khusus:
- `"Menunggu Penugasan"`
- `"Penghulu Ditugaskan"`
- `"Menunggu Verifikasi Penghulu"`

**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`

### 2. Auto Transition
Beberapa status berubah otomatis:
- `"Penghulu Ditugaskan"` → `"Menunggu Verifikasi Penghulu"` (otomatis setelah assign)

### 3. Status Final
Status berikut adalah final dan tidak bisa diubah lagi:
- `"Selesai"` ✅
- `"Ditolak"` ❌

### 4. Notifikasi Otomatis
Setiap perubahan status akan mengirim notifikasi otomatis ke:
- User (calon pasangan)
- Penghulu (jika di-assign)
- Staff (jika diperlukan)

---

## 🔄 Flow Diagram Update Status

```
┌─────────────────────────────────────────────────────────┐
│              CARA UPDATE STATUS FLOW                     │
└─────────────────────────────────────────────────────────┘

1. Draft
   │
   │ [User: Submit] → Menunggu Verifikasi
   │
   ▼
2. Menunggu Verifikasi
   │
   │ [Staff: verify-formulir] → Menunggu Pengumpulan Berkas
   │ [Staff: update-status] → Status apapun (fleksibel)
   │
   ▼
3. Menunggu Pengumpulan Berkas
   │
   │ [Staff: verify-berkas] → Berkas Diterima
   │ [Staff: update-status] → Status apapun (fleksibel)
   │
   ▼
4. Berkas Diterima
   │
   │ [User: mark-visited] → Menunggu Penugasan
   │ [Staff: update-status] → Status apapun (fleksibel)
   │
   ▼
5. Menunggu Penugasan
   │
   │ [Kepala KUA: assign-penghulu] → Penghulu Ditugaskan
   │ [Kepala KUA: update-status] → Status apapun (kecuali status penghulu)
   │
   ▼
6. Penghulu Ditugaskan (Auto)
   │
   │ [System: Auto] → Menunggu Verifikasi Penghulu
   │
   ▼
7. Menunggu Verifikasi Penghulu
   │
   │ [Penghulu: verify-documents] → Menunggu Bimbingan
   │ [Penghulu: update-status] → Status apapun (kecuali status penghulu)
   │
   ▼
8. Menunggu Bimbingan
   │
   │ [Staff: complete-bimbingan] → Sudah Bimbingan
   │ [Staff: update-status] → Status apapun (fleksibel)
   │
   ▼
9. Sudah Bimbingan
   │
   │ [Staff: complete-nikah] → Selesai ✅
   │ [Staff: update-status] → Status apapun (fleksibel)
   │
   ▼
10. Selesai ✅ (Final)
```

---

## 📝 Contoh Lengkap: Update Status dari Awal sampai Akhir

### Scenario: Flow Normal

```typescript
// 1. User submit formulir
// Status: Draft → Menunggu Verifikasi (otomatis)

// 2. Staff verifikasi formulir
await verifyFormulir(1, {
  status: 'Formulir Disetujui',
  catatan: 'Formulir lengkap'
}, staffToken);
// Status: Menunggu Verifikasi → Menunggu Pengumpulan Berkas

// 3. Staff verifikasi berkas
await verifyBerkas(1, {
  status: 'Berkas Diterima',
  catatan: 'Berkas lengkap'
}, staffToken);
// Status: Menunggu Pengumpulan Berkas → Berkas Diterima

// 4. User konfirmasi kunjungan
await markAsVisited(1, userToken);
// Status: Berkas Diterima → Menunggu Penugasan

// 5. Kepala KUA assign penghulu
await assignPenghulu(1, 2, kepalaKuaToken);
// Status: Menunggu Penugasan → Penghulu Ditugaskan → Menunggu Verifikasi Penghulu (auto)

// 6. Penghulu verifikasi dokumen
await verifyDocuments(1, {
  status: 'Menunggu Pelaksanaan',
  catatan: 'Dokumen lengkap'
}, penghuluToken);
// Status: Menunggu Verifikasi Penghulu → Menunggu Bimbingan

// 7. Staff complete bimbingan
await completeBimbingan(1, staffToken);
// Status: Menunggu Bimbingan → Sudah Bimbingan

// 8. Staff complete nikah
await completeNikah(1, staffToken);
// Status: Sudah Bimbingan → Selesai ✅
```

### Scenario: Skip Tahap (Kasus Khusus)

```typescript
// Dari "Menunggu Verifikasi" langsung ke "Menunggu Bimbingan"
await updateStatusPendaftaran(1, {
  status: 'Menunggu Bimbingan',
  catatan: 'Skip tahap karena kasus khusus'
}, staffToken);
// Status: Menunggu Verifikasi → Menunggu Bimbingan (skip beberapa tahap)
```

---

## 🛡️ Error Handling

### Error: Status Tidak Valid
```json
{
  "success": false,
  "message": "Status tidak valid",
  "error": "Status yang diizinkan: Draft, Menunggu Verifikasi, ..."
}
```

**Solution:** Gunakan status yang valid sesuai daftar.

---

### Error: Status Terkait Penghulu
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Status 'Menunggu Penugasan' hanya bisa diubah oleh Kepala KUA melalui endpoint assign-penghulu"
}
```

**Solution:** Gunakan endpoint `POST /simnikah/pendaftaran/:id/assign-penghulu` untuk status terkait penghulu.

---

### Error: Status Sebelumnya Tidak Sesuai
```json
{
  "success": false,
  "message": "Status tidak sesuai",
  "error": "Pendaftaran harus dalam status 'Menunggu Verifikasi' untuk diverifikasi"
}
```

**Solution:** 
- Gunakan endpoint fleksibel untuk skip validasi, atau
- Pastikan status sebelumnya sesuai dengan yang diharapkan

---

## 🔄 Flow Sederhana

Berdasarkan fokus aplikasi (transparansi jadwal, assign penghulu, alamat nikah), flow bisa disederhanakan menjadi:

### Flow Sederhana (4 Tahap)

1. **Draft** → User submit formulir
2. **Menunggu Penugasan** → Kepala KUA assign penghulu
3. **Penghulu Ditugaskan** → Penghulu lihat alamat & jadwal
4. **Selesai** → Nikah dilaksanakan

**Detail:** Lihat `docs/FLOW_STATUS_SEDERHANA.md`

---

## ✅ Best Practices

1. **Gunakan Endpoint Khusus untuk Flow Normal**
   - Memastikan business logic berjalan dengan benar
   - Validasi otomatis sesuai prosedur

2. **Gunakan Endpoint Fleksibel untuk Kasus Khusus**
   - Hanya jika benar-benar diperlukan
   - Dokumentasikan alasan skip tahap di catatan

3. **Selalu Sertakan Catatan**
   - Memudahkan tracking perubahan
   - Dokumentasi untuk audit trail

4. **Handle Error dengan Baik**
   - Tampilkan pesan error yang jelas ke user
   - Log error untuk debugging

5. **Validasi di Frontend**
   - Cek role user sebelum menampilkan opsi
   - Validasi status sebelum API call

---

## 🔗 Referensi

- **Flow Status Lengkap:** `docs/FLOW_STATUS_LENGKAP.md`
- **Endpoint Fleksibel:** `docs/ENDPOINT_UPDATE_STATUS_FLEKSIBEL.md`
- **API Dokumentasi:** `docs/API_DOKUMENTASI_LENGKAP_FRONTEND.md`

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0

