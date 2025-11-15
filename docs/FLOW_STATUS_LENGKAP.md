# 🔄 Flow Proses Status Pendaftaran Nikah

## 📋 Daftar Isi
1. [Diagram Flow Status](#diagram-flow-status)
2. [Detail Setiap Status](#detail-setiap-status)
3. [Transisi Status](#transisi-status)
4. [Business Rules & Validasi](#business-rules--validasi)
5. [Endpoint API per Status](#endpoint-api-per-status)
6. [Contoh Alur Lengkap](#contoh-alur-lengkap)

---

## 📊 Diagram Flow Status

```
┌─────────────────────────────────────────────────────────────────┐
│                    FLOW STATUS PENDAFTARAN NIKAH                │
└─────────────────────────────────────────────────────────────────┘

1. 📝 Draft
   │
   │ [User: Submit Form]
   │ POST /simnikah/pendaftaran
   │
   ▼
2. ⏳ Menunggu Verifikasi
   │
   │ [Staff: Verifikasi Formulir Online]
   │ PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir
   │ Status: "Formulir Disetujui" → Next
   │ Status: "Ditolak" → ❌ END
   │
   ▼
3. 📦 Menunggu Pengumpulan Berkas
   │
   │ [User: Datang ke KUA dengan Berkas Fisik]
   │ [Staff: Verifikasi Berkas Fisik]
   │ PUT /simnikah/staff/pendaftaran/:id/verifikasi-berkas
   │ Status: "Berkas Diterima" → Next
   │ Status: "Ditolak" → ❌ END
   │
   ▼
4. ✅ Berkas Diterima
   │
   │ [User: Konfirmasi Kunjungan]
   │ PUT /simnikah/pendaftaran/:id/konfirmasi-kunjungan
   │
   ▼
5. 🎯 Menunggu Penugasan
   │
   │ [Kepala KUA: Assign Penghulu]
   │ PUT /simnikah/kepala-kua/pendaftaran/:id/assign-penghulu
   │ Validasi: Max 3 nikah/penghulu/hari
   │ Validasi: Min 60 menit gap waktu
   │
   ▼
6. 👨‍⚖️ Penghulu Ditugaskan
   │
   │ [System: Auto Transition]
   │
   ▼
7. 🔍 Menunggu Verifikasi Penghulu
   │
   │ [Penghulu: Verifikasi Dokumen]
   │ PUT /simnikah/penghulu/pendaftaran/:id/verifikasi-dokumen
   │ Status: "Menunggu Pelaksanaan" → Next
   │ Status: "Ditolak" → ❌ END
   │
   ▼
8. 📚 Menunggu Bimbingan
   │
   │ [User: Daftar Bimbingan Perkawinan]
   │ POST /simnikah/bimbingan/:id/daftar
   │ Validasi: Hanya hari Rabu
   │ Validasi: Max 10 pasangan/sesi
   │
   │ [User: Ikut Bimbingan]
   │ [Staff/Kepala KUA: Update Kehadiran]
   │ PUT /simnikah/bimbingan/:id/update-kehadiran
   │
   ▼
9. ✅ Sudah Bimbingan
   │
   │ [Staff/Kepala KUA: Complete Nikah]
   │ PUT /simnikah/pendaftaran/:id/complete-nikah
   │
   ▼
10. 🎉 Selesai ✅
    └─ Nikah telah dilaksanakan
```

---

## 📝 Detail Setiap Status

### 1. **Draft** 📝
- **Deskripsi**: Formulir pendaftaran masih dalam tahap pengisian
- **Aksi yang Bisa Dilakukan**:
  - ✅ User: Mengisi formulir pendaftaran
  - ✅ User: Submit formulir (ubah ke "Menunggu Verifikasi")
- **Validasi**:
  - Semua field wajib harus terisi
  - Validasi dispensasi (jika diperlukan)
  - Validasi wali nikah sesuai syariat
- **Endpoint**:
  - `POST /simnikah/pendaftaran` - Create pendaftaran (status: Draft)
  - `PUT /simnikah/pendaftaran/:id` - Update pendaftaran (jika masih Draft)
  - `PUT /simnikah/pendaftaran/:id/submit` - Submit formulir

---

### 2. **Menunggu Verifikasi** ⏳
- **Deskripsi**: Formulir online sudah disubmit, menunggu verifikasi oleh Staff
- **Aksi yang Bisa Dilakukan**:
  - ✅ Staff: Verifikasi formulir online
  - ✅ Staff: Setujui → "Menunggu Pengumpulan Berkas"
  - ✅ Staff: Tolak → "Ditolak"
- **Validasi**:
  - Hanya Staff yang bisa verifikasi
  - Formulir harus lengkap
- **Endpoint**:
  - `PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir`
- **Request Body**:
```json
{
  "status": "Formulir Disetujui",  // atau "Ditolak"
  "catatan": "Formulir sudah lengkap dan valid"
}
```

---

### 3. **Menunggu Pengumpulan Berkas** 📦
- **Deskripsi**: Formulir online sudah disetujui, user harus datang ke KUA dengan berkas fisik
- **Aksi yang Bisa Dilakukan**:
  - ✅ User: Datang ke KUA dengan membawa berkas
  - ✅ Staff: Verifikasi berkas fisik
  - ✅ Staff: Setujui → "Berkas Diterima"
  - ✅ Staff: Tolak → "Ditolak"
- **Validasi**:
  - User harus datang ke KUA dalam 5 hari kerja
  - Berkas fisik harus lengkap
- **Endpoint**:
  - `PUT /simnikah/staff/pendaftaran/:id/verifikasi-berkas`
- **Request Body**:
```json
{
  "status": "Berkas Diterima",  // atau "Ditolak"
  "catatan": "Berkas sudah lengkap dan sesuai"
}
```

---

### 4. **Berkas Diterima** ✅
- **Deskripsi**: Berkas fisik sudah diterima dan diverifikasi oleh Staff
- **Aksi yang Bisa Dilakukan**:
  - ✅ User: Konfirmasi kunjungan (ubah ke "Menunggu Penugasan")
- **Validasi**:
  - Hanya user pemilik pendaftaran yang bisa konfirmasi
- **Endpoint**:
  - `PUT /simnikah/pendaftaran/:id/konfirmasi-kunjungan`

---

### 5. **Menunggu Penugasan** 🎯
- **Deskripsi**: Berkas sudah diterima, menunggu Kepala KUA untuk menugaskan Penghulu
- **Aksi yang Bisa Dilakukan**:
  - ✅ Kepala KUA: Assign Penghulu (ubah ke "Penghulu Ditugaskan")
- **Validasi**:
  - Max 3 nikah per penghulu per hari
  - Min 60 menit gap waktu antar nikah
  - Max 9 nikah di KUA per hari (jika tempat = "Di KUA")
  - Penghulu harus aktif
- **Endpoint**:
  - `PUT /simnikah/kepala-kua/pendaftaran/:id/assign-penghulu`
- **Request Body**:
```json
{
  "penghulu_id": 1,
  "catatan": "Penghulu ditugaskan untuk tanggal dan waktu yang diminta"
}
```

---

### 6. **Penghulu Ditugaskan** 👨‍⚖️
- **Deskripsi**: Penghulu sudah ditugaskan oleh Kepala KUA
- **Aksi yang Bisa Dilakukan**:
  - ✅ System: Auto transition ke "Menunggu Verifikasi Penghulu"
- **Catatan**:
  - Status ini bersifat sementara
  - System otomatis mengubah ke status berikutnya
- **Endpoint**:
  - Tidak ada endpoint manual, auto transition

---

### 7. **Menunggu Verifikasi Penghulu** 🔍
- **Deskripsi**: Penghulu yang ditugaskan harus memverifikasi dokumen sebelum pelaksanaan
- **Aksi yang Bisa Dilakukan**:
  - ✅ Penghulu: Verifikasi dokumen
  - ✅ Penghulu: Setujui → "Menunggu Bimbingan"
  - ✅ Penghulu: Tolak → "Ditolak"
- **Validasi**:
  - Hanya Penghulu yang ditugaskan yang bisa verifikasi
  - Penghulu harus sesuai dengan yang di-assign
- **Endpoint**:
  - `PUT /simnikah/penghulu/pendaftaran/:id/verifikasi-dokumen`
- **Request Body**:
```json
{
  "status": "Menunggu Pelaksanaan",  // atau "Ditolak"
  "catatan": "Dokumen sudah lengkap dan valid"
}
```

---

### 8. **Menunggu Bimbingan** 📚
- **Deskripsi**: Dokumen sudah diverifikasi Penghulu, user harus mengikuti bimbingan perkawinan
- **Aksi yang Bisa Dilakukan**:
  - ✅ User: Daftar bimbingan perkawinan
  - ✅ User: Ikut bimbingan
  - ✅ Staff: Update kehadiran bimbingan
  - ✅ System: Auto update ke "Sudah Bimbingan" setelah bimbingan selesai
- **Validasi**:
  - Bimbingan hanya diadakan hari Rabu
  - Max 10 pasangan per sesi bimbingan
  - User harus terdaftar dulu sebelum ikut bimbingan
- **Endpoint**:
  - `POST /simnikah/bimbingan/:id/daftar` - Daftar bimbingan
  - `PUT /simnikah/bimbingan/:id/update-kehadiran` - Update kehadiran
  - `PUT /simnikah/pendaftaran/:id/complete-bimbingan` - Complete bimbingan

---

### 9. **Sudah Bimbingan** ✅
- **Deskripsi**: Bimbingan perkawinan sudah selesai diikuti
- **Aksi yang Bisa Dilakukan**:
  - ✅ Staff/Kepala KUA: Complete nikah (ubah ke "Selesai")
- **Validasi**:
  - Bimbingan harus sudah selesai
  - Sertifikat bimbingan harus sudah diterbitkan
- **Endpoint**:
  - `PUT /simnikah/pendaftaran/:id/complete-nikah`

---

### 10. **Selesai** 🎉
- **Deskripsi**: Proses pendaftaran nikah sudah selesai, nikah telah dilaksanakan
- **Aksi yang Bisa Dilakukan**:
  - ✅ Tidak ada aksi lebih lanjut (Final Status)
- **Catatan**:
  - Status final, tidak bisa diubah lagi
  - Sertifikat nikah sudah diterbitkan

---

### ❌ **Ditolak**
- **Deskripsi**: Pendaftaran ditolak pada tahap verifikasi
- **Kapan Bisa Terjadi**:
  - Staff menolak formulir online (dari "Menunggu Verifikasi")
  - Staff menolak berkas fisik (dari "Menunggu Pengumpulan Berkas")
  - Penghulu menolak dokumen (dari "Menunggu Verifikasi Penghulu")
- **Aksi yang Bisa Dilakukan**:
  - ✅ User: Lihat catatan penolakan
  - ✅ User: Perbaiki dan daftar ulang (jika diperlukan)
- **Catatan**:
  - Status final untuk pendaftaran yang ditolak
  - User bisa membuat pendaftaran baru

---

## 🔄 Transisi Status

### Tabel Transisi Status

| Status Saat Ini | Aksi | Status Berikutnya | Actor | Endpoint |
|----------------|------|-------------------|-------|----------|
| **Draft** | Submit form | Menunggu Verifikasi | User | `PUT /simnikah/pendaftaran/:id/submit` |
| **Menunggu Verifikasi** | Approve formulir | Menunggu Pengumpulan Berkas | Staff | `PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir` |
| **Menunggu Verifikasi** | Reject formulir | Ditolak | Staff | `PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir` |
| **Menunggu Pengumpulan Berkas** | Approve berkas | Berkas Diterima | Staff | `PUT /simnikah/staff/pendaftaran/:id/verifikasi-berkas` |
| **Menunggu Pengumpulan Berkas** | Reject berkas | Ditolak | Staff | `PUT /simnikah/staff/pendaftaran/:id/verifikasi-berkas` |
| **Berkas Diterima** | Konfirmasi kunjungan | Menunggu Penugasan | User | `PUT /simnikah/pendaftaran/:id/konfirmasi-kunjungan` |
| **Menunggu Penugasan** | Assign penghulu | Penghulu Ditugaskan | Kepala KUA | `PUT /simnikah/kepala-kua/pendaftaran/:id/assign-penghulu` |
| **Penghulu Ditugaskan** | Auto transition | Menunggu Verifikasi Penghulu | System | Auto |
| **Menunggu Verifikasi Penghulu** | Approve dokumen | Menunggu Bimbingan | Penghulu | `PUT /simnikah/penghulu/pendaftaran/:id/verifikasi-dokumen` |
| **Menunggu Verifikasi Penghulu** | Reject dokumen | Ditolak | Penghulu | `PUT /simnikah/penghulu/pendaftaran/:id/verifikasi-dokumen` |
| **Menunggu Bimbingan** | Complete bimbingan | Sudah Bimbingan | Staff/System | `PUT /simnikah/pendaftaran/:id/complete-bimbingan` |
| **Sudah Bimbingan** | Complete nikah | Selesai | Staff/Kepala KUA | `PUT /simnikah/pendaftaran/:id/complete-nikah` |

---

## 🛡️ Business Rules & Validasi

### 1. **Kapasitas Nikah**

| Rule | Value | Validasi |
|------|-------|----------|
| Max nikah di KUA/hari | 9 | Validasi saat assign penghulu |
| Max nikah di luar KUA | Unlimited | Tidak ada batasan |
| Max nikah per penghulu/hari | 3 | Validasi saat assign penghulu |
| Min gap waktu antar nikah | 60 menit | Validasi konflik waktu |
| Jam operasional | 08:00 - 16:00 | 9 slot waktu (08:00, 09:00, ..., 16:00) |

### 2. **Dispensasi**

**Wajib dispensasi jika:**
- Nikah < 10 hari kerja dari pendaftaran
- Usia calon suami < 19 tahun
- Usia calon istri < 19 tahun

**Validasi:**
- Nomor dispensasi wajib diisi jika salah satu kondisi terpenuhi
- Validasi di `CreateMarriageRegistrationForm()`

### 3. **Bimbingan Perkawinan**

| Rule | Value | Validasi |
|------|-------|----------|
| Hari bimbingan | Hanya Rabu | Validasi saat create bimbingan |
| Max pasangan/sesi | 10 | Validasi saat daftar bimbingan |
| 1 sesi per Rabu | Ya | Validasi duplicate tanggal |

### 4. **Wali Nikah (Syariat Islam)**

**Urutan Wali Nasab:**
1. Ayah Kandung
2. Kakek (Ayah dari Ayah)
3. Saudara Laki-Laki Kandung
4. Saudara Laki-Laki Seayah
5. Keponakan Laki-Laki
6. Paman Kandung
7. Paman Seayah
8. Sepupu Laki-Laki
9. Wali Hakim (jika tidak ada wali nasab)

**Validasi:**
- Jika ayah masih hidup → Wali HARUS ayah kandung
- Jika ayah meninggal → Wali tidak boleh ayah kandung
- Wali harus berbeda NIK dengan calon suami/istri

---

## 🔌 Endpoint API per Status

### Status: Draft
```http
POST /simnikah/pendaftaran
PUT /simnikah/pendaftaran/:id
PUT /simnikah/pendaftaran/:id/submit
GET /simnikah/pendaftaran/:id
```

### Status: Menunggu Verifikasi
```http
GET /simnikah/pendaftaran/:id
PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir
```

### Status: Menunggu Pengumpulan Berkas
```http
GET /simnikah/pendaftaran/:id
PUT /simnikah/staff/pendaftaran/:id/verifikasi-berkas
```

### Status: Berkas Diterima
```http
GET /simnikah/pendaftaran/:id
PUT /simnikah/pendaftaran/:id/konfirmasi-kunjungan
```

### Status: Menunggu Penugasan
```http
GET /simnikah/pendaftaran/:id
GET /simnikah/kepala-kua/penghulu/available
PUT /simnikah/kepala-kua/pendaftaran/:id/assign-penghulu
```

### Status: Penghulu Ditugaskan
```http
GET /simnikah/pendaftaran/:id
# Auto transition ke "Menunggu Verifikasi Penghulu"
```

### Status: Menunggu Verifikasi Penghulu
```http
GET /simnikah/pendaftaran/:id
PUT /simnikah/penghulu/pendaftaran/:id/verifikasi-dokumen
```

### Status: Menunggu Bimbingan
```http
GET /simnikah/pendaftaran/:id
GET /simnikah/bimbingan
POST /simnikah/bimbingan/:id/daftar
PUT /simnikah/bimbingan/:id/update-kehadiran
PUT /simnikah/pendaftaran/:id/complete-bimbingan
```

### Status: Sudah Bimbingan
```http
GET /simnikah/pendaftaran/:id
PUT /simnikah/pendaftaran/:id/complete-nikah
```

### Status: Selesai
```http
GET /simnikah/pendaftaran/:id
GET /simnikah/pendaftaran/:id/detail
```

---

## 📖 Contoh Alur Lengkap

### Scenario 1: Alur Normal (Tanpa Penolakan)

```
1. User membuat pendaftaran
   → Status: Draft
   → POST /simnikah/pendaftaran

2. User submit formulir
   → Status: Menunggu Verifikasi
   → PUT /simnikah/pendaftaran/:id/submit

3. Staff verifikasi formulir online
   → Status: Menunggu Pengumpulan Berkas
   → PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir
   → Body: { "status": "Formulir Disetujui" }

4. User datang ke KUA dengan berkas
   → Staff verifikasi berkas fisik
   → Status: Berkas Diterima
   → PUT /simnikah/staff/pendaftaran/:id/verifikasi-berkas
   → Body: { "status": "Berkas Diterima" }

5. User konfirmasi kunjungan
   → Status: Menunggu Penugasan
   → PUT /simnikah/pendaftaran/:id/konfirmasi-kunjungan

6. Kepala KUA assign penghulu
   → Status: Penghulu Ditugaskan
   → PUT /simnikah/kepala-kua/pendaftaran/:id/assign-penghulu
   → Body: { "penghulu_id": 1 }

7. System auto transition
   → Status: Menunggu Verifikasi Penghulu
   → (Auto)

8. Penghulu verifikasi dokumen
   → Status: Menunggu Bimbingan
   → PUT /simnikah/penghulu/pendaftaran/:id/verifikasi-dokumen
   → Body: { "status": "Menunggu Pelaksanaan" }

9. User daftar bimbingan
   → POST /simnikah/bimbingan/:id/daftar

10. User ikut bimbingan
    → Staff update kehadiran
    → PUT /simnikah/bimbingan/:id/update-kehadiran
    → Status: Sudah Bimbingan

11. Staff complete nikah
    → Status: Selesai ✅
    → PUT /simnikah/pendaftaran/:id/complete-nikah
```

### Scenario 2: Alur dengan Penolakan

```
1. User membuat pendaftaran
   → Status: Draft

2. User submit formulir
   → Status: Menunggu Verifikasi

3. Staff verifikasi formulir online
   → Status: Ditolak ❌
   → PUT /simnikah/staff/pendaftaran/:id/verifikasi-formulir
   → Body: { "status": "Ditolak", "catatan": "Data tidak lengkap" }

4. User melihat catatan penolakan
   → GET /simnikah/pendaftaran/:id

5. User perbaiki dan daftar ulang
   → POST /simnikah/pendaftaran (pendaftaran baru)
```

### Scenario 3: Alur dengan Dispensasi

```
1. User membuat pendaftaran
   → Tanggal nikah: < 10 hari kerja dari sekarang
   → Validasi: Wajib nomor dispensasi
   → POST /simnikah/pendaftaran
   → Body: { ..., "nomor_dispensasi": "DISP/2024/001" }

2. (Lanjut seperti Scenario 1)
```

---

## 📊 Status Constants (Go)

```go
// Status Pendaftaran Nikah
const (
    StatusPendaftaranDraft                      = "Draft"
    StatusPendaftaranMenungguVerifikasi         = "Menunggu Verifikasi"
    StatusPendaftaranMenungguPengumpulanBerkas   = "Menunggu Pengumpulan Berkas"
    StatusPendaftaranBerkasDiterima             = "Berkas Diterima"
    StatusPendaftaranMenungguPenugasan          = "Menunggu Penugasan"
    StatusPendaftaranPenghuluDitugaskan         = "Penghulu Ditugaskan"
    StatusPendaftaranMenungguVerifikasiPenghulu = "Menunggu Verifikasi Penghulu"
    StatusPendaftaranMenungguBimbingan          = "Menunggu Bimbingan"
    StatusPendaftaranSudahBimbingan             = "Sudah Bimbingan"
    StatusPendaftaranSelesai                    = "Selesai"
    StatusPendaftaranDitolak                    = "Ditolak"
)
```

---

## 🎯 Kesimpulan

**Total Status**: 11 status (10 normal + 1 ditolak)

**Flow Normal**: 10 tahap dari Draft sampai Selesai

**Actor yang Terlibat**:
- 👤 User (Calon Pasangan)
- 👨‍💼 Staff
- 👨‍⚖️ Penghulu
- 👔 Kepala KUA
- 🤖 System (Auto transition)

**Waktu Estimasi**:
- Minimal: ~10 hari kerja (jika semua lancar)
- Normal: 2-3 minggu
- Dengan dispensasi: Bisa lebih cepat

---

**📝 Catatan Penting:**
- Setiap transisi status memiliki validasi ketat
- Notifikasi otomatis dikirim ke user pada setiap perubahan status
- Status tidak bisa di-skip atau diubah mundur (kecuali ditolak)
- Semua aksi dicatat dengan timestamp dan actor yang melakukan

