# 📋 Alur Pendaftaran dan Aturan Jam/Waktu Pernikahan

**Versi:** 1.1.0  
**Update Terakhir:** November 2024

---

## 🔄 Alur Pendaftaran Pernikahan

### Status Flow (Urutan Status)

```
1. Draft
   ↓ (Staff verifikasi & approve)
2. Disetujui
   ↓ (Kepala KUA assign penghulu)
3. Menunggu Penugasan
   ↓ (Kepala KUA menentukan penghulu)
4. Penghulu Ditugaskan
   ↓ (Penghulu melaksanakan nikah)
5. Selesai ✅
```

**Status Khusus:**
- **Ditolak** - Dapat terjadi kapan saja jika pendaftaran ditolak oleh staff

---

## 📅 Aturan Jam/Waktu Pernikahan

### 1. **Jam Operasional**
- **Jam Mulai:** 08:00 WITA
- **Jam Selesai:** 16:00 WITA
- **Time Slots Tersedia:** 
  - 08:00, 09:00, 10:00, 11:00, 12:00
  - 13:00, 14:00, 15:00, 16:00
  - **Total:** 9 slot waktu per hari

### 2. **Kapasitas Pernikahan per Jam**

#### **A. Nikah di KUA**
- **Maksimal:** 1 pernikahan per jam
- **Alasan:** Hanya ada 1 ruang nikah di KUA
- **Contoh:** 
  - Jam 09:00 → Maksimal 1 pernikahan di KUA
  - Jika sudah ada 1 pernikahan di KUA jam 09:00, maka tidak bisa daftar lagi di jam yang sama

#### **B. Nikah di Luar KUA**
- **Maksimal:** 3 pernikahan per jam (total semua penghulu)
- **Alasan:** Ada 3 penghulu yang bisa ditugaskan ke lokasi berbeda
- **Contoh:**
  - Jam 09:00 → Maksimal 3 pernikahan di luar KUA
  - Jika sudah ada 3 pernikahan di luar KUA jam 09:00, maka tidak bisa daftar lagi di jam yang sama

#### **C. Kombinasi (KUA + Luar KUA)**
- **Total maksimal per jam:** 3 pernikahan
- **Aturan:**
  - Jika ada 1 pernikahan di KUA → Slot luar KUA = 2 (total 3)
  - Jika belum ada di KUA → Slot luar KUA = 3 (total 3)
- **Contoh:**
  - Jam 10:00:
    - 1 pernikahan di KUA ✅
    - 2 pernikahan di luar KUA ✅
    - Total = 3 pernikahan (PENUH)
  - Jam 11:00:
    - 0 pernikahan di KUA ✅
    - 3 pernikahan di luar KUA ✅
    - Total = 3 pernikahan (PENUH)

---

## ✅ Validasi Waktu saat Pendaftaran

### 1. **Format Waktu**
- **Format:** `HH:MM` (24 jam)
- **Contoh Valid:** `09:00`, `14:30`, `16:00`
- **Contoh Invalid:** `9:00`, `09:0`, `25:00`

### 2. **Validasi Tanggal**
- ✅ Tanggal nikah **tidak boleh** di masa lalu
- ✅ Format tanggal: `YYYY-MM-DD` (contoh: `2024-11-20`)
- ✅ Tanggal harus valid (tidak boleh 31 Februari, dll)

### 3. **Validasi Ketersediaan Slot**
Sistem akan mengecek:
- ✅ Apakah slot waktu sudah terisi?
- ✅ Apakah kapasitas maksimal sudah tercapai?
- ✅ Apakah ada konflik dengan pernikahan lain?

### 4. **Status yang Diperhitungkan dalam Kuota**
- ✅ **Draft** - Diperhitungkan (meskipun belum pasti)
- ✅ **Disetujui** - Diperhitungkan
- ✅ **Menunggu Penugasan** - Diperhitungkan
- ✅ **Penghulu Ditugaskan** - Diperhitungkan
- ❌ **Ditolak** - **TIDAK** diperhitungkan
- ❌ **Selesai** - **TIDAK** diperhitungkan (sudah lewat)

---

## 📊 Contoh Skenario

### Skenario 1: Pendaftaran di KUA
```
User A mendaftar:
- Tanggal: 2024-11-20
- Waktu: 09:00
- Tempat: Di KUA
→ ✅ BERHASIL (slot KUA jam 09:00 masih kosong)

User B mendaftar:
- Tanggal: 2024-11-20
- Waktu: 09:00
- Tempat: Di KUA
→ ❌ GAGAL (slot KUA jam 09:00 sudah terisi oleh User A)
```

### Skenario 2: Pendaftaran di Luar KUA
```
User A mendaftar:
- Tanggal: 2024-11-20
- Waktu: 10:00
- Tempat: Di Luar KUA
→ ✅ BERHASIL (slot luar KUA jam 10:00 = 1/3)

User B mendaftar:
- Tanggal: 2024-11-20
- Waktu: 10:00
- Tempat: Di Luar KUA
→ ✅ BERHASIL (slot luar KUA jam 10:00 = 2/3)

User C mendaftar:
- Tanggal: 2024-11-20
- Waktu: 10:00
- Tempat: Di Luar KUA
→ ✅ BERHASIL (slot luar KUA jam 10:00 = 3/3)

User D mendaftar:
- Tanggal: 2024-11-20
- Waktu: 10:00
- Tempat: Di Luar KUA
→ ❌ GAGAL (slot luar KUA jam 10:00 sudah penuh 3/3)
```

### Skenario 3: Kombinasi KUA + Luar KUA
```
User A mendaftar:
- Tanggal: 2024-11-20
- Waktu: 11:00
- Tempat: Di KUA
→ ✅ BERHASIL (slot KUA jam 11:00 = 1/1)

User B mendaftar:
- Tanggal: 2024-11-20
- Waktu: 11:00
- Tempat: Di Luar KUA
→ ✅ BERHASIL (slot luar KUA jam 11:00 = 1/2, total = 2/3)

User C mendaftar:
- Tanggal: 2024-11-20
- Waktu: 11:00
- Tempat: Di Luar KUA
→ ✅ BERHASIL (slot luar KUA jam 11:00 = 2/2, total = 3/3)

User D mendaftar:
- Tanggal: 2024-11-20
- Waktu: 11:00
- Tempat: Di Luar KUA
→ ❌ GAGAL (total sudah 3/3, tidak ada slot tersedia)
```

---

## 🔍 Endpoint untuk Cek Ketersediaan

### 1. **Cek Kalender Ketersediaan (Bulanan)**
```
GET /simnikah/kalender-ketersediaan?bulan=11&tahun=2024
```
Menampilkan:
- Tanggal yang tersedia/tidak tersedia
- Jumlah draft dan disetujui per hari
- Time slots detail per hari

### 2. **Cek Time Slots Tersedia (Harian)**
```
GET /simnikah/time-slots?tanggal=2024-11-20
```
Menampilkan:
- Daftar jam yang tersedia (08:00 - 16:00)
- Status setiap slot (tersedia/terisi)
- Jumlah draft dan disetujui per slot
- Detail untuk KUA dan luar KUA

---

## 📝 Catatan Penting

1. **Draft dihitung dalam kuota** - Meskipun status masih "Draft" (belum pasti), pendaftaran tersebut tetap menghitung slot waktu untuk mencegah double booking.

2. **Hanya "Ditolak" yang tidak dihitung** - Status "Ditolak" tidak menghitung slot karena pendaftaran tersebut dibatalkan.

3. **Selesai tidak dihitung** - Status "Selesai" tidak menghitung slot karena pernikahan sudah selesai dilaksanakan.

4. **Validasi real-time** - Sistem melakukan validasi ketersediaan saat pendaftaran dibuat, bukan saat approval.

5. **Tidak bisa ubah waktu setelah disetujui** - Setelah status "Disetujui", waktu nikah tidak bisa diubah untuk menjaga konsistensi jadwal.

---

## 🎯 Ringkasan Aturan

| Tempat Nikah | Maksimal per Jam | Total Maksimal per Jam |
|--------------|------------------|------------------------|
| **Di KUA** | 1 pernikahan | - |
| **Di Luar KUA** | 3 pernikahan | - |
| **Kombinasi** | - | 3 pernikahan (1 KUA + 2 Luar KUA) |

| Status | Diperhitungkan dalam Kuota? |
|--------|----------------------------|
| **Draft** | ✅ Ya |
| **Disetujui** | ✅ Ya |
| **Menunggu Penugasan** | ✅ Ya |
| **Penghulu Ditugaskan** | ✅ Ya |
| **Selesai** | ❌ Tidak |
| **Ditolak** | ❌ Tidak |

---

**Dokumen ini menjelaskan alur pendaftaran dan aturan jam/waktu pernikahan di sistem SimNikah.**

