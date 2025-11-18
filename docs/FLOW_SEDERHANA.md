# 📋 Flow Sederhana Pendaftaran Nikah

## 🎯 Overview

Flow pendaftaran nikah yang disederhanakan untuk memudahkan proses. Catin hanya perlu mendaftar, staff verifikasi di belakang layar, kepala KUA tentukan penghulu, dan penghulu laksanakan nikah.

## 📊 Status Flow

```
Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
   ↓
Ditolak
```

### Status Detail:

1. **Draft** 
   - Catin daftar melalui form sederhana (nama+bin/binti, pendidikan, umur, lokasi nikah)
   - Status awal setelah pendaftaran

2. **Disetujui**
   - Staff menyetujui pendaftaran setelah verifikasi dokumen di belakang layar
   - Proses verifikasi dilakukan di kantor secara manual

3. **Menunggu Penugasan** (optional)
   - Status intermediate setelah disetujui
   - Menunggu kepala KUA menentukan penghulu

4. **Penghulu Ditugaskan**
   - Kepala KUA sudah assign penghulu untuk pendaftaran ini
   - Penghulu bisa melihat detail pendaftaran dan lokasi nikah

5. **Selesai**
   - Penghulu sudah melaksanakan nikah
   - Status final setelah pernikahan

6. **Ditolak**
   - Pendaftaran ditolak oleh staff
   - Bisa terjadi dari status Draft

---

## 🔄 Flow Proses

### 1. Catin Daftar (Draft)
**Endpoint:** `POST /simnikah/pendaftaran/form-sederhana`

**Request:**
```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad bin Abdullah",
    "pendidikan_akhir": "SMA",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti binti Abdullah",
    "pendidikan_akhir": "SMA",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di KUA",
    "tanggal_nikah": "2024-02-14",
    "waktu_nikah": "09:00"
  }
}
```

**Status:** `Draft`

---

### 2. Staff Menyetujui Pendaftaran (Disetujui)
**Endpoint:** `POST /simnikah/pendaftaran/:id/approve`

**Request:**
```json
{
  "status": "Disetujui",
  "catatan": "Dokumen lengkap, disetujui"
}
```

**Flow:** `Draft` → `Disetujui`

**Note:** Staff verifikasi dokumen di belakang layar (manual di kantor), kemudian menyetujui melalui aplikasi.

---

### 3. Kepala KUA Menentukan Penghulu (Penghulu Ditugaskan)
**Endpoint:** `POST /simnikah/pendaftaran/:id/assign-penghulu`

**Request:**
```json
{
  "penghulu_id": 1,
  "catatan": "Ditugaskan ke Penghulu Ahmad"
}
```

**Flow:** `Disetujui` / `Menunggu Penugasan` → `Penghulu Ditugaskan`

---

### 4. Penghulu Melaksanakan Nikah (Selesai)
**Endpoint:** `POST /simnikah/pendaftaran/:id/complete-nikah`

**Request:**
```json
{
  "catatan": "Pernikahan telah dilaksanakan dengan lancar"
}
```

**Flow:** `Penghulu Ditugaskan` → `Selesai`

**Note:** Penghulu mengupdate status menjadi "Selesai" setelah melaksanakan nikah.

---

## 🔑 Endpoints Baru

### Staff - Approve Pendaftaran
- **POST** `/simnikah/pendaftaran/:id/approve`
- **Role:** `staff`
- **Flow:** `Draft` → `Disetujui` atau `Draft` → `Ditolak`

### Kepala KUA - Assign Penghulu
- **POST** `/simnikah/pendaftaran/:id/assign-penghulu`
- **Role:** `kepala_kua`
- **Flow:** `Disetujui` / `Menunggu Penugasan` → `Penghulu Ditugaskan`

### Penghulu - Complete Nikah
- **POST** `/simnikah/pendaftaran/:id/complete-nikah`
- **Role:** `penghulu`
- **Flow:** `Penghulu Ditugaskan` → `Selesai`

---

## 📝 Catatan Penting

1. **Verifikasi di Belakang Layar**
   - Catin tidak perlu upload dokumen melalui aplikasi
   - Staff verifikasi dokumen secara manual di kantor
   - Setelah verifikasi selesai, staff menyetujui melalui aplikasi

2. **Flow Sederhana**
   - Hanya 4-5 status utama (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai)
   - Tidak ada verifikasi berkas terpisah, semua dilakukan di belakang layar

3. **Transparansi untuk Catin**
   - Catin bisa melihat status pendaftaran kapan saja
   - Catin bisa melihat penghulu yang ditugaskan dan kontak penghulu

4. **Kemudahan untuk Penghulu**
   - Penghulu bisa melihat lokasi nikah (khusus di luar KUA)
   - Penghulu bisa melihat detail pendaftaran dan catatan

---

## 🔄 Contoh Flow Lengkap

```
1. Catin daftar via aplikasi
   └─ Status: Draft

2. Catin datang ke kantor membawa dokumen
   └─ Staff verifikasi dokumen (manual)

3. Staff menyetujui via aplikasi
   └─ Status: Disetujui
   └─ Notifikasi ke catin: "Pendaftaran disetujui"

4. Kepala KUA melihat daftar pendaftaran yang disetujui
   └─ Kepala KUA assign penghulu
   └─ Status: Penghulu Ditugaskan
   └─ Notifikasi ke penghulu: "Anda ditugaskan untuk nikah"
   └─ Notifikasi ke catin: "Penghulu telah ditugaskan"

5. Penghulu melihat detail pendaftaran dan lokasi nikah
   └─ Penghulu datang ke lokasi (jika di luar KUA)

6. Penghulu melaksanakan nikah
   └─ Penghulu update status via aplikasi
   └─ Status: Selesai
   └─ Notifikasi ke catin: "Pernikahan selesai, selamat!"

7. Selesai
```

---

## ✅ Benefits

1. **Simple**: Flow yang lebih sederhana dan mudah dipahami
2. **Efficient**: Tidak perlu upload dokumen, verifikasi di kantor langsung
3. **Transparent**: Catin bisa melihat status kapan saja
4. **Easy for Staff**: Verifikasi manual lebih fleksibel
5. **Easy for Kepala KUA**: Langsung assign penghulu tanpa banyak tahap
6. **Easy for Penghulu**: Langsung laksanakan tanpa banyak verifikasi

