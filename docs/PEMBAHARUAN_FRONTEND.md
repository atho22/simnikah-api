# 📱 Dokumentasi Pembaharuan API untuk Frontend

## 🎯 Ringkasan Pembaharuan

Telah ditambahkan **endpoint baru** untuk update status pendaftaran nikah secara **fleksibel** tanpa validasi ketat pada status sebelumnya. Namun, ada **pembatasan khusus** untuk status terkait assign penghulu.

---

## ✨ Fitur Baru

### 1. Endpoint Update Status Fleksibel
**Endpoint Baru:** `PUT /simnikah/pendaftaran/:id/update-status`

Endpoint ini memungkinkan **Staff, Penghulu, dan Kepala KUA** untuk mengupdate status pendaftaran tanpa harus mengikuti flow yang ketat.

**Keuntungan:**
- ✅ Bisa update dari status apapun ke status apapun (fleksibel)
- ✅ Tidak perlu validasi status sebelumnya
- ✅ Cocok untuk kasus khusus atau kebutuhan administratif

**Pembatasan:**
- ❌ Status terkait penghulu **TIDAK BISA** diupdate via endpoint ini
- ❌ Hanya Kepala KUA yang bisa assign penghulu (melalui endpoint khusus)

---

## 🔌 Endpoint Baru

### PUT `/simnikah/pendaftaran/:id/update-status`

**Method:** `PUT`  
**Auth:** Required (JWT Token)  
**Role:** `staff`, `penghulu`, atau `kepala_kua`

#### Request Headers
```javascript
{
  "Authorization": "Bearer <jwt_token>",
  "Content-Type": "application/json"
}
```

#### Request Body
```json
{
  "status": "Menunggu Bimbingan",
  "catatan": "Status diupdate secara manual oleh staff"
}
```

#### Response Success (200 OK)
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

#### Response Error (403 Forbidden) - Status Terkait Penghulu
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Status 'Menunggu Penugasan' hanya bisa diubah oleh Kepala KUA melalui endpoint assign-penghulu. Gunakan endpoint POST /simnikah/pendaftaran/:id/assign-penghulu untuk menugaskan penghulu."
}
```

---

## 📋 Status yang Bisa Diupdate

### ✅ Status yang BISA diupdate via endpoint fleksibel:
- `"Draft"`
- `"Menunggu Verifikasi"`
- `"Menunggu Pengumpulan Berkas"`
- `"Berkas Diterima"`
- `"Menunggu Bimbingan"`
- `"Sudah Bimbingan"`
- `"Selesai"`
- `"Ditolak"`

### ❌ Status yang TIDAK BISA diupdate via endpoint fleksibel:
Status berikut **hanya bisa diubah oleh Kepala KUA** melalui endpoint khusus:
- `"Menunggu Penugasan"` → Gunakan `POST /simnikah/pendaftaran/:id/assign-penghulu`
- `"Penghulu Ditugaskan"` → Otomatis saat assign penghulu
- `"Menunggu Verifikasi Penghulu"` → Otomatis saat assign penghulu

---

## 💻 Contoh Implementasi Frontend

### 1. React/Next.js dengan Axios

```typescript
// types/status.ts
export type StatusPendaftaran = 
  | "Draft"
  | "Menunggu Verifikasi"
  | "Menunggu Pengumpulan Berkas"
  | "Berkas Diterima"
  | "Menunggu Bimbingan"
  | "Sudah Bimbingan"
  | "Selesai"
  | "Ditolak";

// Status yang tidak bisa diupdate via endpoint fleksibel
export const PENGHULU_RELATED_STATUSES = [
  "Menunggu Penugasan",
  "Penghulu Ditugaskan",
  "Menunggu Verifikasi Penghulu"
] as const;

// services/pendaftaran.ts
import axios from 'axios';
import { StatusPendaftaran } from '@/types/status';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export interface UpdateStatusRequest {
  status: StatusPendaftaran;
  catatan?: string;
}

export interface UpdateStatusResponse {
  success: boolean;
  message: string;
  data: {
    id: number;
    nomor_pendaftaran: string;
    status_sebelumnya: string;
    status_sekarang: string;
    catatan?: string;
    updated_by: string;
    updated_at: string;
  };
}

export const updateStatusPendaftaran = async (
  pendaftaranId: number,
  data: UpdateStatusRequest,
  token: string
): Promise<UpdateStatusResponse> => {
  // Validasi: cek apakah status terkait penghulu
  if (PENGHULU_RELATED_STATUSES.includes(data.status as any)) {
    throw new Error(
      `Status "${data.status}" hanya bisa diubah oleh Kepala KUA melalui endpoint assign-penghulu`
    );
  }

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
```

### 2. React Component untuk Update Status

```tsx
// components/UpdateStatusModal.tsx
import React, { useState } from 'react';
import { updateStatusPendaftaran, StatusPendaftaran } from '@/services/pendaftaran';
import { PENGHULU_RELATED_STATUSES } from '@/types/status';

interface UpdateStatusModalProps {
  pendaftaranId: number;
  currentStatus: string;
  userRole: 'staff' | 'penghulu' | 'kepala_kua';
  onSuccess: () => void;
  onClose: () => void;
}

const AVAILABLE_STATUSES: StatusPendaftaran[] = [
  'Draft',
  'Menunggu Verifikasi',
  'Menunggu Pengumpulan Berkas',
  'Berkas Diterima',
  'Menunggu Bimbingan',
  'Sudah Bimbingan',
  'Selesai',
  'Ditolak',
];

export const UpdateStatusModal: React.FC<UpdateStatusModalProps> = ({
  pendaftaranId,
  currentStatus,
  userRole,
  onSuccess,
  onClose,
}) => {
  const [selectedStatus, setSelectedStatus] = useState<StatusPendaftaran>('Draft');
  const [catatan, setCatatan] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('Token tidak ditemukan');
      }

      await updateStatusPendaftaran(
        pendaftaranId,
        {
          status: selectedStatus,
          catatan: catatan || undefined,
        },
        token
      );

      onSuccess();
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal mengupdate status');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <h2>Update Status Pendaftaran</h2>
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Status Saat Ini:</label>
            <input type="text" value={currentStatus} disabled />
          </div>

          <div className="form-group">
            <label>Status Baru:</label>
            <select
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value as StatusPendaftaran)}
              required
            >
              <option value="">Pilih Status</option>
              {AVAILABLE_STATUSES.map((status) => (
                <option key={status} value={status}>
                  {status}
                </option>
              ))}
            </select>
            {PENGHULU_RELATED_STATUSES.includes(selectedStatus as any) && (
              <p className="error-text">
                ⚠️ Status ini hanya bisa diubah oleh Kepala KUA melalui menu Assign Penghulu
              </p>
            )}
          </div>

          <div className="form-group">
            <label>Catatan (Opsional):</label>
            <textarea
              value={catatan}
              onChange={(e) => setCatatan(e.target.value)}
              rows={3}
              placeholder="Masukkan catatan untuk perubahan status..."
            />
          </div>

          {error && <div className="error-message">{error}</div>}

          <div className="modal-actions">
            <button type="button" onClick={onClose} disabled={loading}>
              Batal
            </button>
            <button type="submit" disabled={loading || PENGHULU_RELATED_STATUSES.includes(selectedStatus as any)}>
              {loading ? 'Mengupdate...' : 'Update Status'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
```

### 3. Vue.js dengan Composition API

```vue
<!-- components/UpdateStatusModal.vue -->
<template>
  <div class="modal" v-if="isOpen">
    <div class="modal-content">
      <h2>Update Status Pendaftaran</h2>
      
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>Status Saat Ini:</label>
          <input type="text" :value="currentStatus" disabled />
        </div>

        <div class="form-group">
          <label>Status Baru:</label>
          <select v-model="selectedStatus" required>
            <option value="">Pilih Status</option>
            <option v-for="status in availableStatuses" :key="status" :value="status">
              {{ status }}
            </option>
          </select>
          <p v-if="isPenghuluRelatedStatus" class="error-text">
            ⚠️ Status ini hanya bisa diubah oleh Kepala KUA melalui menu Assign Penghulu
          </p>
        </div>

        <div class="form-group">
          <label>Catatan (Opsional):</label>
          <textarea
            v-model="catatan"
            rows="3"
            placeholder="Masukkan catatan untuk perubahan status..."
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div class="modal-actions">
          <button type="button" @click="close" :disabled="loading">
            Batal
          </button>
          <button type="submit" :disabled="loading || isPenghuluRelatedStatus">
            {{ loading ? 'Mengupdate...' : 'Update Status' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { updateStatusPendaftaran } from '@/services/pendaftaran';

const PENGHULU_RELATED_STATUSES = [
  'Menunggu Penugasan',
  'Penghulu Ditugaskan',
  'Menunggu Verifikasi Penghulu',
];

const props = defineProps<{
  isOpen: boolean;
  pendaftaranId: number;
  currentStatus: string;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'success'): void;
}>();

const selectedStatus = ref('');
const catatan = ref('');
const loading = ref(false);
const error = ref<string | null>(null);

const availableStatuses = [
  'Draft',
  'Menunggu Verifikasi',
  'Menunggu Pengumpulan Berkas',
  'Berkas Diterima',
  'Menunggu Bimbingan',
  'Sudah Bimbingan',
  'Selesai',
  'Ditolak',
];

const isPenghuluRelatedStatus = computed(() => {
  return PENGHULU_RELATED_STATUSES.includes(selectedStatus.value);
});

const handleSubmit = async () => {
  error.value = null;
  loading.value = true;

  try {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('Token tidak ditemukan');
    }

    await updateStatusPendaftaran(
      props.pendaftaranId,
      {
        status: selectedStatus.value,
        catatan: catatan.value || undefined,
      },
      token
    );

    emit('success');
    close();
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Gagal mengupdate status';
  } finally {
    loading.value = false;
  }
};

const close = () => {
  selectedStatus.value = '';
  catatan.value = '';
  error.value = null;
  emit('close');
};
</script>
```

---

## 🔄 Perbandingan Endpoint

### Endpoint Lama (Masih Tersedia)

| Endpoint | Role | Validasi Status | Keterangan |
|----------|------|----------------|------------|
| `POST /staff/verify-formulir/:id` | Staff | Harus "Menunggu Verifikasi" | Verifikasi formulir online |
| `POST /staff/verify-berkas/:id` | Staff | Harus "Menunggu Pengumpulan Berkas" | Verifikasi berkas fisik |
| `POST /penghulu/verify-documents/:id` | Penghulu | Harus "Menunggu Verifikasi Penghulu" | Verifikasi dokumen |
| `POST /pendaftaran/:id/assign-penghulu` | Kepala KUA | Harus "Menunggu Penugasan" | Assign penghulu |

### Endpoint Baru

| Endpoint | Role | Validasi Status | Keterangan |
|----------|------|----------------|------------|
| `PUT /pendaftaran/:id/update-status` | Staff, Penghulu, Kepala KUA | ❌ Tidak ada (fleksibel) | Update status fleksibel (kecuali status penghulu) |

---

## ⚠️ Error Handling

### 1. Status Terkait Penghulu (403 Forbidden)
```typescript
try {
  await updateStatusPendaftaran(id, { status: 'Menunggu Penugasan' }, token);
} catch (error) {
  if (error.response?.status === 403) {
    // Tampilkan pesan: "Status ini hanya bisa diubah oleh Kepala KUA"
    // Redirect atau tampilkan tombol ke menu Assign Penghulu
  }
}
```

### 2. Status Tidak Valid (400 Bad Request)
```typescript
try {
  await updateStatusPendaftaran(id, { status: 'Status Invalid' }, token);
} catch (error) {
  if (error.response?.status === 400) {
    // Tampilkan error message dari response
    console.error(error.response.data.error);
  }
}
```

### 3. Unauthorized (401)
```typescript
try {
  await updateStatusPendaftaran(id, data, token);
} catch (error) {
  if (error.response?.status === 401) {
    // Redirect ke login page
    router.push('/login');
  }
}
```

---

## 🎨 UI/UX Recommendations

### 1. Tampilkan Warning untuk Status Penghulu
```tsx
{selectedStatus && PENGHULU_RELATED_STATUSES.includes(selectedStatus) && (
  <div className="alert alert-warning">
    <strong>⚠️ Perhatian:</strong> Status ini hanya bisa diubah oleh Kepala KUA melalui menu 
    <button onClick={() => router.push(`/kepala-kua/assign-penghulu/${pendaftaranId}`)}>
      Assign Penghulu
    </button>
  </div>
)}
```

### 2. Disable Option untuk Status Penghulu
```tsx
<select value={selectedStatus} onChange={handleChange}>
  {AVAILABLE_STATUSES.map((status) => (
    <option 
      key={status} 
      value={status}
      disabled={PENGHULU_RELATED_STATUSES.includes(status) && userRole !== 'kepala_kua'}
    >
      {status}
      {PENGHULU_RELATED_STATUSES.includes(status) && userRole !== 'kepala_kua' && ' (Hanya Kepala KUA)'}
    </option>
  ))}
</select>
```

### 3. Tampilkan Status Sebelumnya dan Sesudahnya
```tsx
<div className="status-change-info">
  <div className="status-before">
    <span className="label">Status Sebelumnya:</span>
    <span className="value">{response.data.status_sebelumnya}</span>
  </div>
  <div className="arrow">→</div>
  <div className="status-after">
    <span className="label">Status Sekarang:</span>
    <span className="value">{response.data.status_sekarang}</span>
  </div>
</div>
```

---

## 📝 Checklist Implementasi

- [ ] Tambahkan endpoint baru ke API service
- [ ] Buat type definitions untuk request/response
- [ ] Implementasikan validasi di frontend (cek status penghulu)
- [ ] Buat UI component untuk update status
- [ ] Handle error responses dengan baik
- [ ] Tampilkan warning untuk status terkait penghulu
- [ ] Update dokumentasi internal tim
- [ ] Test dengan berbagai role (Staff, Penghulu, Kepala KUA)
- [ ] Test error scenarios (403, 400, 401, 404)

---

## 🔗 Referensi

- **Endpoint Detail:** `docs/ENDPOINT_UPDATE_STATUS_FLEKSIBEL.md`
- **Flow Status:** `docs/FLOW_STATUS_LENGKAP.md`
- **Base URL:** `process.env.NEXT_PUBLIC_API_URL` atau sesuai konfigurasi

---

## 💡 Tips

1. **Gunakan endpoint fleksibel untuk:**
   - Skip tahap tertentu (kasus khusus)
   - Kembali ke status sebelumnya (jika ada kesalahan)
   - Update manual untuk kebutuhan administratif

2. **Jangan gunakan endpoint fleksibel untuk:**
   - Assign penghulu (gunakan endpoint khusus)
   - Flow normal (gunakan endpoint yang sudah ada dengan validasi)

3. **Selalu validasi di frontend:**
   - Cek role user sebelum menampilkan opsi status
   - Blokir status terkait penghulu untuk Staff/Penghulu
   - Tampilkan pesan error yang jelas

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0

