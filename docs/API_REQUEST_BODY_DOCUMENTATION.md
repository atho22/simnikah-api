# Dokumentasi Body Request & Logic Percabangan SimNikah API

## Daftar Isi
1. [Pendaftaran Nikah](#1-pendaftaran-nikah)
2. [Ketersediaan Jadwal](#2-ketersediaan-jadwal)
3. [Verifikasi & Approval](#3-verifikasi--approval)
4. [Assign Penghulu](#4-assign-penghulu)
5. [Aturan Kapasitas](#5-aturan-kapasitas)

---

## 1. Pendaftaran Nikah

### Endpoint: `POST /simnikah/pendaftaran`

```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad bin Abdullah",
    "pendidikan_akhir": "S1",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti binti Muhammad",
    "pendidikan_akhir": "S1",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di KUA",
    "tanggal_nikah": "2024-12-25",
    "waktu_nikah": "09:00",
    "alamat_nikah": "",
    "detail_alamat": "",
    "kelurahan": ""
  },
  "wali_nikah": {
    "nama_dan_bin": "Abdullah bin Muhammad",
    "hubungan_wali": "Ayah Kandung"
  }
}
```

### Percabangan Logic Pendaftaran

```
┌─────────────────────────────────────────────────────────────────┐
│                    VALIDASI PENDAFTARAN                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. CEK UMUR                                                     │
│    - Laki-laki >= 19 tahun                                      │
│    - Perempuan >= 16 tahun                                      │
│    ❌ Error: "Umur tidak memenuhi syarat"                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. CEK TANGGAL                                                  │
│    - Tidak boleh tanggal yang sudah lewat                       │
│    - Minimal H+10 dari hari ini                                 │
│    ❌ Error: "Tanggal tidak valid"                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. CEK HARI LIBUR (jika tempat_nikah = "Di KUA")                │
│    - Minggu = Libur                                             │
│    - Hari libur nasional = Libur                                │
│    ❌ Error: "KUA tutup pada hari libur"                        │
│    ✅ Saran: "Pilih nikah di luar KUA"                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. CEK KAPASITAS                                                │
│    Lihat bagian "Aturan Kapasitas" di bawah                     │
└─────────────────────────────────────────────────────────────────┘
```

### Field `tempat_nikah`

| Nilai | Alamat Wajib? | Hari Libur |
|-------|---------------|------------|
| `"Di KUA"` | ❌ Tidak | ❌ Tidak bisa |
| `"Di Luar KUA"` | ✅ Wajib | ✅ Bisa |

### Kelurahan Valid (Banjarmasin Utara)
- Sungai Miai, Sungai Andai, Surgi Mufti, Pangeran
- Kuin Utara, Antasan Kecil Timur
- Alalak Utara, Alalak Tengah, Alalak Selatan

---

## 2. Ketersediaan Jadwal

### Endpoint: `GET /simnikah/kalender-ketersediaan`

Query: `?bulan=12&tahun=2024`

### Response per Hari

```json
{
  "tanggal": 25,
  "tanggal_str": "2024-12-25",
  "hari": "Wednesday",
  "status": "Tersedia",
  "tersedia": true,
  "tersedia_kua": true,
  "tersedia_luar_kua": true,
  "is_hari_libur": false,
  "nama_hari_libur": "",
  "time_slots": [
    {
      "waktu": "09:00",
      "kua": {
        "tersedia": true,
        "terbooking": false,
        "slot_tersisa": 1
      },
      "luar_kua": {
        "tersedia": true,
        "terbooking": false,
        "slot_tersisa": 3
      },
      "total_pernikahan": 0,
      "slot_tersisa": 3
    }
  ]
}
```

### Percabangan Status Hari

```
┌─────────────────────────────────────────────────────────────────┐
│                    STATUS KETERSEDIAAN HARI                     │
└─────────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            ▼                 ▼                 ▼
    ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
    │ Tanggal Lewat │ │  Hari Libur   │ │  Hari Biasa   │
    └───────────────┘ └───────────────┘ └───────────────┘
            │                 │                 │
            ▼                 ▼                 ▼
    ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
    │status:Terlewat│ │status:Libur   │ │ Cek Kapasitas │
    │tersedia:false │ │tersedia:true  │ │               │
    │kua:false      │ │kua:false      │ │               │
    │luar_kua:false │ │luar_kua:true  │ │               │
    └───────────────┘ └───────────────┘ └───────────────┘
                                                │
                              ┌─────────────────┼─────────────────┐
                              ▼                 ▼                 ▼
                      ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
                      │ Total >= 9    │ │ 0 < Total < 9 │ │  Total = 0    │
                      └───────────────┘ └───────────────┘ └───────────────┘
                              │                 │                 │
                              ▼                 ▼                 ▼
                      ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
                      │status: Penuh  │ │status:Sebagian│ │status:Tersedia│
                      │tersedia:false │ │tersedia:true  │ │tersedia:true  │
                      └───────────────┘ └───────────────┘ └───────────────┘
```

---

## 3. Verifikasi & Approval

### A. Verifikasi Formulir (Staff)

**Endpoint:** `POST /simnikah/staff/verify-formulir/:id`

```json
{
  "status": "Disetujui",
  "catatan": "Formulir lengkap dan valid"
}
```

| Status | Hasil |
|--------|-------|
| `"Disetujui"` | Lanjut ke approval |
| `"Ditolak"` | Pendaftaran ditolak |

### B. Approve Pendaftaran (Staff)

**Endpoint:** `POST /simnikah/staff/approve/:id`

```json
{
  "status": "Disetujui",
  "catatan": "Pendaftaran disetujui"
}
```

### Flow Status Pendaftaran

```
Draft → Disetujui → Penghulu Ditugaskan → Selesai
          ↓
       Ditolak
```

---

## 4. Assign Penghulu

### Endpoint: `POST /simnikah/pendaftaran/:id/assign-penghulu`

```json
{
  "penghulu_id": 1,
  "catatan": "Ditugaskan ke H. Ahmad"
}
```

### Percabangan Validasi Assign

```
┌─────────────────────────────────────────────────────────────────┐
│                    VALIDASI ASSIGN PENGHULU                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. CEK STATUS PENDAFTARAN                                       │
│    - Harus "Disetujui" atau "Menunggu Penugasan"                │
│    ❌ Error: "Status tidak sesuai"                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. CEK PENGHULU AKTIF                                           │
│    - Penghulu harus status "Aktif"                              │
│    ❌ Error: "Penghulu tidak ditemukan"                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. CEK JADWAL PENGHULU                                          │
│    - 1 penghulu = 1 pernikahan per jam                          │
│    - Cek apakah sudah ada jadwal di tanggal+jam yang sama       │
│    ❌ Error: "Penghulu tidak tersedia"                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ ✅ BERHASIL: Status → "Penghulu Ditugaskan"                     │
└─────────────────────────────────────────────────────────────────┘
```

### Lihat Penghulu Tersedia

**Endpoint:** `GET /simnikah/kepala-kua/penghulu-tersedia?tanggal=2024-12-25&waktu=09:00`

---

## 5. Aturan Kapasitas

### Batasan Utama

| Batasan | Nilai |
|---------|-------|
| Maksimal per hari | 9 pernikahan |
| Maksimal per jam | 3 pernikahan |
| Maksimal KUA per jam | 1 pernikahan |
| Maksimal per penghulu per jam | 1 pernikahan |
| Jumlah penghulu | 3 orang |
| Time slots | 08:00 - 16:00 (9 slot) |

### Percabangan Kapasitas per Jam

```
┌─────────────────────────────────────────────────────────────────┐
│              VALIDASI KAPASITAS PER JAM (Contoh 09:00)          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ Total pernikahan di jam ini = ?                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│  Total = 0    │     │ Total = 1-2   │     │  Total = 3    │
│               │     │               │     │               │
│ KUA: ✅ 1     │     │ Cek detail    │     │ KUA: ❌       │
│ Luar: ✅ 3    │     │               │     │ Luar: ❌      │
└───────────────┘     └───────────────┘     └───────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            ▼                 ▼                 ▼
    ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
    │ 1 di KUA      │ │ 1 di Luar     │ │ 2 di Luar     │
    │               │ │               │ │               │
    │ KUA: ❌       │ │ KUA: ✅ 1     │ │ KUA: ✅ 1     │
    │ Luar: ✅ 2    │ │ Luar: ✅ 2    │ │ Luar: ✅ 1    │
    └───────────────┘ └───────────────┘ └───────────────┘
```

### Tabel Skenario Kapasitas

| Kondisi Jam | Di KUA | Di Luar KUA | Total Tersisa |
|-------------|--------|-------------|---------------|
| Kosong | ✅ 1 slot | ✅ 3 slot | 3 |
| 1 di KUA | ❌ Penuh | ✅ 2 slot | 2 |
| 1 di Luar | ✅ 1 slot | ✅ 2 slot | 2 |
| 2 di Luar | ✅ 1 slot | ✅ 1 slot | 1 |
| 3 di Luar | ❌ Penuh | ❌ Penuh | 0 |
| 1 KUA + 2 Luar | ❌ Penuh | ❌ Penuh | 0 |

### Hari Libur

| Hari/Tanggal | Di KUA | Di Luar KUA |
|--------------|--------|-------------|
| Minggu | ❌ | ✅ |
| Hari Libur Nasional | ❌ | ✅ |
| Hari Biasa | ✅ | ✅ |

---

## 6. Body Request Lainnya

### Update Alamat Nikah
**Endpoint:** `PUT /simnikah/pendaftaran/:id/alamat`
```json
{
  "alamat_akad": "Jl. Pangeran No. 10, RT 05 RW 02"
}
```

### Complete Marriage (Penghulu)
**Endpoint:** `POST /simnikah/penghulu/complete-marriage/:id`
```json
{
  "catatan": "Pernikahan telah dilaksanakan"
}
```

### Buat Notifikasi
**Endpoint:** `POST /simnikah/notifikasi`
```json
{
  "user_id": "USR123",
  "judul": "Pengingat",
  "pesan": "Pernikahan Anda besok",
  "tipe": "Info",
  "link": "/pendaftaran/1"
}
```

### Buat Staff (Kepala KUA)
**Endpoint:** `POST /simnikah/kepala-kua/staff`
```json
{
  "username": "staff_baru",
  "email": "staff@kua.go.id",
  "password": "password123",
  "nama": "Staff Baru",
  "nip": "198501012010011001",
  "jabatan": "Staff",
  "bagian": "Pelayanan",
  "no_hp": "081234567890",
  "alamat": "Jl. Contoh No. 1"
}
```

### Buat Penghulu (Kepala KUA)
**Endpoint:** `POST /simnikah/kepala-kua/penghulu`
```json
{
  "username": "penghulu_baru",
  "email": "penghulu@kua.go.id",
  "password": "password123",
  "nama": "H. Penghulu Baru",
  "nip": "198501012010011002",
  "no_hp": "081234567891"
}
```

---

## 7. Error Response Format

```json
{
  "success": false,
  "message": "Pesan singkat error",
  "error": "Penjelasan detail error",
  "type": "validation|schedule_conflict|holiday_restriction|not_found",
  "field": "nama_field_yang_error",
  "saran": "Saran untuk user",
  "data": {
    "info_tambahan": "nilai"
  }
}
```

### Tipe Error

| Type | Deskripsi |
|------|-----------|
| `validation` | Data tidak valid |
| `schedule_conflict` | Jadwal bentrok |
| `holiday_restriction` | Hari libur |
| `not_found` | Data tidak ditemukan |
| `authentication` | Tidak terautentikasi |
| `authorization` | Tidak punya akses |

---

## 8. Surat Pengumuman Nikah

### Endpoint

| Role | Endpoint |
|------|----------|
| Staff | `GET /simnikah/staff/pengumuman-nikah/list` |
| Kepala KUA | `GET /simnikah/kepala-kua/pengumuman-nikah/list` |

### Query Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|-----------|
| `tanggal_awal` | date | ❌ | Tanggal awal periode (YYYY-MM-DD) |
| `tanggal_akhir` | date | ❌ | Tanggal akhir periode (YYYY-MM-DD) |

> Jika tidak diisi, default = minggu ini (Senin - Minggu)

### Request Body (Optional) - Kop Surat Custom

```json
{
  "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
  "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
  "kota": "Kota Banjarmasin",
  "provinsi": "Kalimantan Selatan",
  "kode_pos": "70123",
  "telepon": "(0511) 123456",
  "email": "kua.bjmutara@kemenag.go.id",
  "website": "https://kua-bjmutara.kemenag.go.id",
  "logo_url": "https://example.com/logo-kua.png"
}
```

### Response Success (200)

```json
{
  "success": true,
  "message": "Data pendaftaran disetujui berhasil diambil",
  "data": {
    "tanggal_awal": "2024-12-02",
    "tanggal_akhir": "2024-12-08",
    "periode": "02 Desember 2024 s/d 08 Desember 2024",
    "total": 5,
    "kop_surat": {
      "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
      "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
      "kota": "Kota Banjarmasin",
      "provinsi": "Kalimantan Selatan",
      "kode_pos": "70123",
      "telepon": "-",
      "email": "-",
      "website": "",
      "logo_url": ""
    },
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "REG-20241202-001",
        "tanggal_nikah": "2024-12-05T00:00:00Z",
        "waktu_nikah": "09:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "",
        "calon_suami": {
          "nama_lengkap": "Ahmad bin Abdullah"
        },
        "calon_istri": {
          "nama_lengkap": "Siti binti Muhammad"
        },
        "wali_nikah": {
          "nama_dan_bin": "Abdullah bin Muhammad",
          "hubungan_wali": "Ayah Kandung"
        }
      },
      {
        "id": 2,
        "nomor_pendaftaran": "REG-20241202-002",
        "tanggal_nikah": "2024-12-06T00:00:00Z",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di Luar KUA",
        "alamat_akad": "Jl. Pangeran No. 10, RT 05 RW 02, Kel. Pangeran",
        "calon_suami": {
          "nama_lengkap": "Budi bin Hasan"
        },
        "calon_istri": {
          "nama_lengkap": "Ani binti Ali"
        },
        "wali_nikah": {
          "nama_dan_bin": "Ali bin Umar",
          "hubungan_wali": "Ayah Kandung"
        }
      }
    ]
  }
}
```

### Contoh Penggunaan

#### 1. Ambil data minggu ini (default)
```
GET /simnikah/staff/pengumuman-nikah/list
```

#### 2. Ambil data periode tertentu
```
GET /simnikah/staff/pengumuman-nikah/list?tanggal_awal=2024-12-01&tanggal_akhir=2024-12-31
```

#### 3. Dengan kop surat custom
```bash
curl -X GET "http://localhost:8080/simnikah/staff/pengumuman-nikah/list" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_kua": "KUA Kecamatan Banjarmasin Utara",
    "telepon": "(0511) 123456",
    "email": "kua@example.com"
  }'
```

### Flow Pengumuman Nikah

```
┌─────────────────────────────────────────────────────────────────┐
│                    SURAT PENGUMUMAN NIKAH                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Staff/Kepala KUA request data                                │
│    - Pilih periode (minggu ini / custom)                        │
│    - Optional: custom kop surat                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. API mengambil pendaftaran dengan:                            │
│    - Status = "Disetujui"                                       │
│    - Tanggal nikah dalam periode yang dipilih                   │
│    - Diurutkan berdasarkan tanggal & waktu                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Response berisi:                                             │
│    - Kop surat (untuk header dokumen)                           │
│    - Daftar pernikahan yang disetujui                           │
│    - Data calon pengantin & wali nikah                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Frontend generate dokumen PDF/print                          │
│    - Tampilkan kop surat                                        │
│    - Tampilkan tabel pengumuman nikah                           │
│    - Siap ditempel di papan pengumuman KUA                      │
└─────────────────────────────────────────────────────────────────┘
```

### Catatan
- Hanya pendaftaran dengan status **"Disetujui"** yang ditampilkan
- Data diurutkan berdasarkan tanggal nikah (ASC) lalu waktu nikah (ASC)
- Kop surat memiliki nilai default jika tidak diisi
