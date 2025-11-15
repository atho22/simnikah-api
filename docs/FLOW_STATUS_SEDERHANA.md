# 🔄 Flow Status Sederhana - SimNikah

## 🎯 Fokus Aplikasi

Aplikasi SimNikah fokus pada:
1. **Transparansi Jadwal** - User dan staff bisa lihat jadwal nikah
2. **Mempermudah Kepala KUA** - Menentukan dan assign penghulu dengan mudah
3. **Mempermudah Penghulu** - Mengetahui alamat catin yang nikah di luar KUA

---

## 📊 Flow Status Sederhana (4 Tahap)

```
┌─────────────────────────────────────────────────────────┐
│              FLOW STATUS SEDERHANA                        │
└─────────────────────────────────────────────────────────┘

1. 📝 Draft
   │
   │ [User: Submit Form]
   │ POST /simnikah/pendaftaran/form-baru
   │
   ▼
2. ⏳ Menunggu Penugasan
   │
   │ [Kepala KUA: Assign Penghulu]
   │ POST /simnikah/pendaftaran/:id/assign-penghulu
   │
   ▼
3. 👨‍⚖️ Penghulu Ditugaskan
   │
   │ [System: Auto Transition]
   │ Status: "Penghulu Ditugaskan"
   │ Data: Alamat nikah (jika di luar KUA) sudah tersedia
   │
   ▼
4. ✅ Selesai
   │
   │ [Staff/Kepala KUA: Complete]
   │ PUT /simnikah/pendaftaran/:id/complete-nikah
   │
   └─ Nikah telah dilaksanakan
```

---

## 📝 Status yang Digunakan

### Status Utama (5 Status)

1. **Draft** 📝
   - Formulir masih dalam pengisian
   - User bisa edit dan submit

2. **Menunggu Penugasan** ⏳
   - Formulir sudah disubmit
   - Siap untuk Kepala KUA assign penghulu
   - **Fokus:** Kepala KUA bisa lihat jadwal dan assign penghulu

3. **Penghulu Ditugaskan** 👨‍⚖️
   - Penghulu sudah ditugaskan oleh Kepala KUA
   - **Fokus:** Penghulu bisa lihat alamat nikah (jika di luar KUA)
   - Auto transition setelah assign
   - Status ini bisa langsung di-complete ke "Selesai"

4. **Selesai** ✅
   - Nikah telah dilaksanakan
   - Final status

### Status Tambahan

- **Ditolak** ❌ - Jika ada penolakan

**Catatan:** Status "Siap Nikah" bisa digunakan jika diperlukan, tapi bisa juga langsung dari "Penghulu Ditugaskan" ke "Selesai"

---

## 🔌 Endpoint Update Status

### 1. Update Status Fleksibel (Recommended)

**PUT** `/simnikah/pendaftaran/:id/update-status`

**Request:**
```json
{
  "status": "Selesai",
  "catatan": "Nikah telah dilaksanakan"
}
```

**Catatan:** Field `status` harus diisi dengan benar (case-sensitive). Pastikan menggunakan huruf kapital di awal kata.

**Status yang BISA diupdate:**
- ✅ `"Draft"`
- ✅ `"Menunggu Penugasan"`
- ✅ `"Selesai"`
- ✅ `"Ditolak"`

**Status yang TIDAK BISA diupdate:**
- ❌ `"Penghulu Ditugaskan"` → Otomatis saat assign penghulu

---

### 2. Assign Penghulu (Kepala KUA)

**POST** `/simnikah/pendaftaran/:id/assign-penghulu`

**Role:** `kepala_kua`  
**Status Sebelumnya:** `"Menunggu Penugasan"`

**Request:**
```json
{
  "penghulu_id": 1
}
```

**Auto Transition:**
- `"Menunggu Penugasan"` → `"Penghulu Ditugaskan"` (otomatis)

---

### 3. Complete Nikah

**PUT** `/simnikah/pendaftaran/:id/complete-nikah`

**Role:** `staff`, `kepala_kua`  
**Status Sebelumnya:** `"Penghulu Ditugaskan"` atau status apapun (fleksibel)

**Response:**
```json
{
  "message": "Nikah berhasil diselesaikan",
  "data": {
    "pendaftaran_id": 1,
    "status_sekarang": "Selesai"
  }
}
```

---

## 📍 Fitur Utama

### 1. Transparansi Jadwal

**Endpoint:**
- `GET /simnikah/kalender-ketersediaan` - Kalender bulanan
- `GET /simnikah/kalender-tanggal-detail` - Detail per tanggal
- `GET /simnikah/ketersediaan-tanggal/:tanggal` - Ketersediaan tanggal
- `GET /simnikah/penghulu-jadwal/:tanggal` - Jadwal penghulu

**Fitur:**
- ✅ User bisa lihat jadwal nikah
- ✅ Staff bisa lihat semua jadwal
- ✅ Kepala KUA bisa lihat jadwal untuk assign penghulu

---

### 2. Assign Penghulu (Kepala KUA)

**Endpoint:**
- `POST /simnikah/pendaftaran/:id/assign-penghulu` - Assign penghulu
- `GET /simnikah/pendaftaran/belum-assign-penghulu` - List belum assign
- `GET /simnikah/penghulu/:id/ketersediaan/:tanggal` - Cek ketersediaan

**Fitur:**
- ✅ Lihat pendaftaran yang belum assign penghulu
- ✅ Cek ketersediaan penghulu per tanggal
- ✅ Validasi kapasitas (max 3/hari per penghulu)
- ✅ Validasi konflik waktu (min 60 menit gap)

---

### 3. Alamat Nikah (Penghulu)

**Endpoint:**
- `GET /simnikah/pendaftaran/:id/location` - Detail lokasi nikah
- `PUT /simnikah/pendaftaran/:id/location` - Update lokasi dengan koordinat
- `POST /simnikah/location/geocode` - Alamat → Koordinat
- `POST /simnikah/location/reverse-geocode` - Koordinat → Alamat

**Fitur:**
- ✅ Penghulu bisa lihat alamat nikah (jika di luar KUA)
- ✅ Link Google Maps, Waze, OpenStreetMap
- ✅ Koordinat GPS untuk navigasi
- ✅ Geocoding otomatis saat input alamat

---

## 🔄 Contoh Flow Lengkap

### Scenario: Nikah di KUA

```
1. User submit formulir
   → Status: Draft → Menunggu Penugasan
   → POST /simnikah/pendaftaran/form-baru

2. Kepala KUA lihat jadwal
   → GET /simnikah/kalender-ketersediaan?bulan=2024-02
   → GET /simnikah/pendaftaran/belum-assign-penghulu

3. Kepala KUA assign penghulu
   → POST /simnikah/pendaftaran/:id/assign-penghulu
   → Body: { "penghulu_id": 1 }
   → Status: Menunggu Penugasan → Penghulu Ditugaskan (auto)

4. Penghulu lihat jadwal
   → GET /simnikah/penghulu/assigned-registrations
   → Lihat jadwal nikah yang ditugaskan

5. Staff complete nikah
   → PUT /simnikah/pendaftaran/:id/complete-nikah
   → Status: Penghulu Ditugaskan → Selesai ✅
   
   Atau bisa langsung:
   → PUT /simnikah/pendaftaran/:id/update-status
   → Body: { "status": "Selesai" }
```

---

### Scenario: Nikah di Luar KUA

```
1. User submit formulir (dengan alamat nikah)
   → Status: Draft → Menunggu Penugasan
   → POST /simnikah/pendaftaran/form-baru
   → Body: { ..., "tempat_nikah": "Di Luar KUA", "alamat_akad": "Jl. ..." }

2. System geocoding otomatis
   → Alamat → Koordinat GPS
   → Simpan latitude & longitude

3. Kepala KUA assign penghulu
   → POST /simnikah/pendaftaran/:id/assign-penghulu
   → Status: Menunggu Penugasan → Penghulu Ditugaskan

4. Penghulu lihat alamat & lokasi
   → GET /simnikah/pendaftaran/:id/location
   → Response: {
       "alamat_akad": "Jl. Pangeran Samudra No. 88",
       "latitude": -3.3194374,
       "longitude": 114.5900675,
       "maps": {
         "google_maps": "https://www.google.com/maps?q=...",
         "waze": "https://waze.com/ul?ll=..."
       }
     }

5. Penghulu bisa navigasi ke lokasi
   → Buka link Google Maps atau Waze
   → Navigasi ke alamat nikah

6. Staff complete nikah
   → PUT /simnikah/pendaftaran/:id/complete-nikah
   → Status: Penghulu Ditugaskan → Selesai ✅
   
   Atau bisa langsung:
   → PUT /simnikah/pendaftaran/:id/update-status
   → Body: { "status": "Selesai" }
```

---

## 📋 Tabel Status Sederhana

| Status | Deskripsi | Siapa yang Bisa Update | Endpoint |
|--------|-----------|------------------------|----------|
| **Draft** | Formulir dalam pengisian | User (submit) | `POST /simnikah/pendaftaran/form-baru` |
| **Menunggu Penugasan** | Siap untuk assign penghulu | Kepala KUA (assign) | `POST /simnikah/pendaftaran/:id/assign-penghulu` |
| **Penghulu Ditugaskan** | Penghulu sudah ditugaskan | System (auto) | Auto setelah assign |
| **Selesai** | Nikah telah dilaksanakan | Staff/Kepala KUA | `PUT /simnikah/pendaftaran/:id/complete-nikah` atau `PUT /simnikah/pendaftaran/:id/update-status` |
| **Ditolak** | Pendaftaran ditolak | Staff/Kepala KUA | `PUT /simnikah/pendaftaran/:id/update-status` |

---

## 🎯 Endpoint Utama

### Untuk User (Calon Pasangan)
- `POST /simnikah/pendaftaran/form-baru` - Daftar nikah
- `GET /simnikah/pendaftaran/status` - Cek status
- `GET /simnikah/kalender-ketersediaan` - Lihat jadwal

### Untuk Kepala KUA
- `GET /simnikah/pendaftaran/belum-assign-penghulu` - List belum assign
- `GET /simnikah/kalender-ketersediaan` - Lihat jadwal
- `GET /simnikah/penghulu/:id/ketersediaan/:tanggal` - Cek ketersediaan
- `POST /simnikah/pendaftaran/:id/assign-penghulu` - Assign penghulu
- `PUT /simnikah/pendaftaran/:id/update-status` - Update status

### Untuk Penghulu
- `GET /simnikah/penghulu/assigned-registrations` - Jadwal yang ditugaskan
- `GET /simnikah/pendaftaran/:id/location` - **Alamat nikah (jika di luar KUA)**
- `GET /simnikah/penghulu-jadwal/:tanggal` - Jadwal per tanggal

### Untuk Staff
- `GET /simnikah/pendaftaran` - List semua pendaftaran
- `PUT /simnikah/pendaftaran/:id/update-status` - Update status
- `PUT /simnikah/pendaftaran/:id/complete-nikah` - Complete nikah

---

## 💡 Tips Penggunaan

### 1. Update Status Fleksibel
Gunakan endpoint fleksibel untuk update status dengan mudah:
```json
PUT /simnikah/pendaftaran/:id/update-status
{
  "status": "Siap Nikah",
  "catatan": "Semua sudah siap untuk pelaksanaan"
}
```

### 2. Assign Penghulu
Kepala KUA bisa assign penghulu dengan mudah:
```json
POST /simnikah/pendaftaran/:id/assign-penghulu
{
  "penghulu_id": 1
}
```

### 3. Lihat Alamat Nikah
Penghulu bisa lihat alamat nikah yang di luar KUA:
```json
GET /simnikah/pendaftaran/:id/location
```

Response akan berisi:
- Alamat lengkap
- Koordinat GPS (latitude, longitude)
- Link Google Maps, Waze, OpenStreetMap

---

## ⚠️ Catatan Penting

1. **Status "Penghulu Ditugaskan"** otomatis setelah assign penghulu
2. **Alamat nikah** hanya relevan jika `tempat_nikah = "Di Luar KUA"`
3. **Geocoding** otomatis saat input alamat (gratis via OpenStreetMap)
4. **Validasi kapasitas** penghulu saat assign (max 3/hari)
5. **Validasi konflik waktu** saat assign (min 60 menit gap)

---

## 🔗 Referensi

- **API Dokumentasi Lengkap:** `docs/API_DOKUMENTASI_LENGKAP_FRONTEND.md`
- **Cara Update Status:** `docs/CARA_UPDATE_STATUS_FLOW.md`
- **Endpoint Fleksibel:** `docs/ENDPOINT_UPDATE_STATUS_FLEKSIBEL.md`

---

**Last Updated:** 2024-01-15  
**Version:** 2.0.0 (Simplified)

