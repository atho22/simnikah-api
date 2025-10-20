# API Testing Documentation

## Scope
- Authentication & RBAC (AuthMiddleware, RoleMiddleware, MultiRoleMiddleware)
- User flows: pendaftaran sampai selesai dengan status baru
- Staff verification (formulir, berkas)
- Penghulu assignment (kuota harian 9, jarak 60 menit), verifikasi penghulu
- Kalender: bulan (warna kuning/hijau), detail per hari (jam dan warna)

## Conventions
- Base URL: `http://localhost:8080`
- Prefix: `/simnikah`
- Auth: `Authorization: Bearer <jwt>`
- Roles used: `user_biasa`, `staff`, `kepala_kua`, `penghulu`

## End-to-End Flow (Happy Path)
1) user_biasa membuat pendaftaran (status: Menunggu Verifikasi)
2) staff verify-formulir -> Formulir Disetujui -> backend: Menunggu Pengumpulan Berkas
3) staff verify-berkas -> Menunggu Penugasan
4) kepala_kua assign-penghulu -> Menunggu Verifikasi Penghulu
5) penghulu verify-documents -> Menunggu Bimbingan
6) daftar-bimbingan -> tetap Menunggu Bimbingan
7) complete-bimbingan -> Sudah Bimbingan
8) complete-nikah -> Selesai

## Negative Tests & Edge Cases
- verify-formulir saat status ≠ Menunggu Verifikasi -> 400
- verify-berkas saat status ≠ Menunggu Pengumpulan Berkas -> 400
- assign-penghulu saat status ≠ Menunggu Penugasan -> 400
- assign-penghulu jika kuota harian >= 9 -> 400
- assign-penghulu jika selisih waktu penghulu < 60 menit -> 400 + detail konflik
- verify-documents (penghulu) saat status ≠ Menunggu Verifikasi Penghulu -> 400
- daftar-bimbingan ketika status bukan Menunggu Bimbingan -> 400

## Kalender Bulanan
- Endpoint: GET `/simnikah/kalender-ketersediaan?bulan=MM&tahun=YYYY`
- Asersi:
  - Semua pendaftaran dalam bulan itu tercakup
  - Per hari ada: tanggal, status, tersedia, jumlah_nikah, kuning_count, hijau_count, warna, sisa_kuota, kapasitas
  - warna: hijau jika ada booked fix; kuning jika hanya pending; kosong jika 0

## Kalender Detail Harian
- Endpoint: GET `/simnikah/kalender-tanggal-detail?tanggal=YYYY-MM-DD`
- Asersi item per jam:
  - id, nomor_pendaftaran, waktu, tempat, status, warna
  - penghulu_id, nama_penghulu
  - nama_calon_suami, nama_calon_istri

## RBAC Checks
- AuthMiddleware aktif; role salah -> 403/401

## Data Migration & ENUM
- Pastikan status baru valid; update invalid -> error DB

## Performance & Concurrency
- Kalender bulan dengan banyak data: respons < 1s pada 500 entries
- Assign bersamaan pada penghulu sama dan jam berdekatan (<60m) -> salah satu gagal

## Postman Collection (saran)
- Struktur folder per peran, urutkan langkah E2E
- Environment vars: `BASE`, tokens, ids; gunakan test scripts untuk simpan id

## Script Bantuan (opsional)
- `test_profile.sh` untuk export BASE dan tokens; jalankan urutan curl sesuai flow di atas

---

## Full Request Bodies

### 1) POST /simnikah/pendaftaran (user_biasa)
Gunakan saat membuat pendaftaran baru.

```json
{
  "scheduleAndLocation": {
    "weddingLocation": "Di KUA",               
    "weddingAddress": "",                      
    "weddingDate": "2025-10-21",              
    "weddingTime": "09:00",                   
    "dispensationNumber": ""                   
  },
  "groom": {
    "groomFullName": "Ahmad",
    "groomNik": "3201010101010001",
    "groomCitizenship": "WNI",
    "groomPassportNumber": "",
    "groomPlaceOfBirth": "Bandung",
    "groomDateOfBirth": "1995-01-02",
    "groomStatus": "Belum Kawin",
    "groomReligion": "Islam",
    "groomEducation": "SMA",
    "groomOccupation": "Karyawan",
    "groomOccupationDescription": "",
    "groomPhoneNumber": "081234567890",
    "groomEmail": "ahmad@example.com",
    "groomAddress": "Jl. Melati No. 1"
  },
  "bride": {
    "brideFullName": "Aisyah",
    "brideNik": "3201010101010002",
    "brideCitizenship": "WNI",
    "bridePassportNumber": "",
    "bridePlaceOfBirth": "Bandung",
    "brideDateOfBirth": "1996-02-03",
    "brideStatus": "Belum Kawin",
    "brideReligion": "Islam",
    "brideEducation": "SMA",
    "brideOccupation": "Karyawan",
    "brideOccupationDescription": "",
    "bridePhoneNumber": "081298765432",
    "brideEmail": "aisyah@example.com",
    "brideAddress": "Jl. Anggrek No. 2"
  },
  "groomParents": {
    "groomFather": {
      "groomFatherPresenceStatus": "Hidup",
      "groomFatherName": "Bapak Ahmad",
      "groomFatherNik": "",
      "groomFatherCitizenship": "WNI",
      "groomFatherCountryOfOrigin": "",
      "groomFatherPassportNumber": "",
      "groomFatherPlaceOfBirth": "Bandung",
      "groomFatherDateOfBirth": "1970-05-01",
      "groomFatherReligion": "Islam",
      "groomFatherOccupation": "Wiraswasta",
      "groomFatherOccupationDescription": "",
      "groomFatherAddress": "Jl. Melati No. 1"
    },
    "groomMother": {
      "groomMotherPresenceStatus": "Hidup",
      "groomMotherName": "Ibu Ahmad",
      "groomMotherNik": "",
      "groomMotherCitizenship": "WNI",
      "groomMotherCountryOfOrigin": "",
      "groomMotherPassportNumber": "",
      "groomMotherPlaceOfBirth": "Bandung",
      "groomMotherDateOfBirth": "1972-07-01",
      "groomMotherReligion": "Islam",
      "groomMotherOccupation": "Ibu Rumah Tangga",
      "groomMotherOccupationDescription": "",
      "groomMotherAddress": "Jl. Melati No. 1"
    }
  },
  "brideParents": {
    "brideFather": {
      "brideFatherPresenceStatus": "Hidup",
      "brideFatherName": "Bapak Aisyah",
      "brideFatherNik": "",
      "brideFatherCitizenship": "WNI",
      "brideFatherCountryOfOrigin": "",
      "brideFatherPassportNumber": "",
      "brideFatherPlaceOfBirth": "Bandung",
      "brideFatherDateOfBirth": "1971-03-02",
      "brideFatherReligion": "Islam",
      "brideFatherOccupation": "PNS",
      "brideFatherOccupationDescription": "",
      "brideFatherAddress": "Jl. Anggrek No. 2"
    },
    "brideMother": {
      "brideMotherPresenceStatus": "Hidup",
      "brideMotherName": "Ibu Aisyah",
      "brideMotherNik": "",
      "brideMotherCitizenship": "WNI",
      "brideMotherCountryOfOrigin": "",
      "brideMotherPassportNumber": "",
      "brideMotherPlaceOfBirth": "Bandung",
      "brideMotherDateOfBirth": "1973-08-03",
      "brideMotherReligion": "Islam",
      "brideMotherOccupation": "Ibu Rumah Tangga",
      "brideMotherOccupationDescription": "",
      "brideMotherAddress": "Jl. Anggrek No. 2"
    }
  },
  "guardian": {
    "guardianFullName": "Wali Nikah",
    "guardianNik": "",
    "guardianRelationship": "Ayah",
    "guardianStatus": "Hidup",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. Anggrek No. 2",
    "guardianPhoneNumber": "0812000000"
  }
}
```

### 2) POST /simnikah/staff/verify-formulir/:id (staff)
```json
{
  "status": "Formulir Disetujui",
  "catatan": "Formulir lengkap dan valid"
}
```
- Nilai status valid: `"Formulir Disetujui"`, `"Formulir Ditolak"`

### 3) POST /simnikah/staff/verify-berkas/:id (staff)
```json
{
  "status": "Berkas Diterima",
  "catatan": "Semua berkas sudah lengkap (KK, KTP, dll)"
}
```
- Nilai status valid: `"Berkas Diterima"`, `"Berkas Ditolak"`

### 4) POST /simnikah/pendaftaran/:id/assign-penghulu (kepala_kua)
```json
{
  "penghulu_id": 1
}
```

### 5) POST /simnikah/penghulu/verify-documents/:id (penghulu)
```json
{
  "status": "Menunggu Pelaksanaan",
  "catatan": "Dokumen lengkap dan siap dilanjutkan"
}
```
- Nilai status valid: `"Menunggu Pelaksanaan"`, `"Ditolak"`

### 6) POST /simnikah/bimbingan/:id/daftar (user_biasa)
- Tanpa body

### 7) POST /simnikah/bimbingan/:id/complete (staff)
- Tanpa body

### 8) POST /simnikah/pendaftaran/:id/complete (staff/kepala)
- Tanpa body

### 9) GET /simnikah/kalender-ketersediaan?bulan=MM&tahun=YYYY
- Tanpa body (query params di URL)

### 10) GET /simnikah/kalender-tanggal-detail?tanggal=YYYY-MM-DD
- Tanpa body (query params di URL)
