# 📋 Update Format Pengumuman Nikah

**Tanggal Update:** Desember 2024  
**Versi:** 2.0.0  
**Status:** ✅ Production Ready

---

## 🎯 Overview

Format pengumuman nikah telah diupdate untuk mengikuti format Excel standar KUA dengan **15 kolom** dan layout **landscape A4**. Update ini membuat format HTML output lebih sesuai dengan format dokumen resmi yang digunakan di KUA.

---

## 🔄 Perubahan Format

### Format Lama (v1.0.0)
- Judul: "PENGUMUMAN PERNIKAHAN"
- Tabel: 9 kolom
- Layout: A4 Portrait
- Kolom: No, Nomor Pendaftaran, Tanggal Nikah, Waktu, Tempat, Alamat Akad, Calon Suami, Calon Istri, Wali Nikah

### Format Baru (v2.0.0) ✨
- Judul: "JADUAL NIKAH [BULAN] [TAHUN]" (contoh: "JADUAL NIKAH JANUARI 2026")
- Tabel: **15 kolom** dengan grouping
- Layout: **A4 Landscape**
- Kop surat: Logo di kiri, informasi KUA di kanan

---

## 📊 Struktur Tabel Baru (15 Kolom)

### Header Tabel

| NO URUT | DATA CALON PENGANTIN | | | | | | PELAKSANAAN NIKAH | | | | | | | | |
|---------|----------------------|-|-|-|-|-|-------------------|-|-|-|-|-|-|-|-|
| | PRIA / BIN | USIA | PENDK | WANITA / BINTI | USIA | PENDK | HARI | TGL | JAM | TEMPAT | WALINIKAH | PENGHULU | KELURAHAN | KET |

**Catatan:**
- ✅ **Kolom tetap 15 kolom** (tidak berubah)
- ✅ **Baris dinamis** - jumlah baris tergantung jumlah data pendaftaran dalam periode yang dipilih
- ✅ **Setiap baris = 1 pendaftaran nikah**
- ✅ Semua pendaftaran dalam periode akan ditampilkan (kecuali status "Ditolak")

### Detail Kolom

1. **NO URUT** - Nomor urut pendaftaran
2. **PRIA / BIN** - Nama lengkap calon suami
3. **USIA** - Usia calon suami (dihitung otomatis dari tanggal lahir)
4. **PENDK** - Pendidikan terakhir calon suami
5. **WANITA / BINTI** - Nama lengkap calon istri
6. **USIA** - Usia calon istri (dihitung otomatis dari tanggal lahir)
7. **PENDK** - Pendidikan terakhir calon istri
8. **HARI** - Nama hari (SENIN, SELASA, RABU, KAMIS, JUM'AT, SABTU, AHAD)
9. **TGL** - Tanggal (hanya angka: 1, 2, 3, dll)
10. **JAM** - Waktu (format: 08.00, 09.00, dll - dari HH:MM menjadi HH.MM)
11. **TEMPAT** - Tempat nikah:
    - Jika `tempat_nikah` = "Di KUA" → menampilkan "Di KUA"
    - Jika `tempat_nikah` = "Di Luar KUA" → menampilkan alamat lengkap dari `alamat_akad`
12. **WALINIKAH** - Nama wali nikah
13. **PENGHULU** - Nama penghulu (jika sudah ditugaskan, "-" jika belum)
14. **KELURAHAN** - Kelurahan
15. **KET** - Keterangan

---

## 🆕 Fitur Baru

### 1. Perhitungan Usia Otomatis
- Usia calon pengantin dihitung otomatis dari tanggal lahir
- Format: Angka bulat (contoh: 29, 24, 26)

### 2. Format Waktu
- Format input: `HH:MM` (contoh: 08:00, 09:00)
- Format output: `HH.MM` (contoh: 08.00, 09.00)

### 3. Nama Hari Otomatis
- Nama hari otomatis dalam bahasa Indonesia
- Format: SENIN, SELASA, RABU, KAMIS, JUM'AT, SABTU, AHAD

### 4. Format Tanggal
- Hanya menampilkan angka tanggal (1, 2, 3, dll)
- Tidak menampilkan bulan dan tahun di kolom TGL

### 6. Judul Dinamis
- Format: "JADUAL NIKAH [BULAN] [TAHUN]"
- Contoh: "JADUAL NIKAH JANUARI 2026"
- Bulan otomatis berdasarkan tanggal awal periode

### 7. Kop Surat Baru
- Logo di kiri (jika disediakan)
- Informasi KUA di kanan
- Format: KEMENTERIAN AGAMA REPUBLIK INDONESIA
- KANTOR KEMENTERIAN AGAMA KOTA [KOTA]
- Nama KUA

---

## 📐 Layout & Styling

### Layout
- **Format:** A4 Landscape
- **Margin:** 1cm
- **Font:** Times New Roman (serif)

### Font Size
- Kop surat: 11-12pt
- Judul: 14pt (bold)
- Tabel header: 8pt (bold)
- Tabel content: 8-9pt

### Tabel
- Border: 1px solid black
- Padding: 4-5px per cell
- Header dengan background color: #e0e0e0
- Rowspan/colspan untuk grouping kolom

---

## 🔧 Technical Changes

### Backend Changes

1. **RegData Struct** - Diupdate untuk 15 kolom:
   ```go
   type RegData struct {
       NoUrut           int
       PriaBin          string
       UsiaPria         int
       PendidikanPria   string
       WanitaBinti      string
       UsiaWanita       int
       PendidikanWanita string
       Hari             string
       Tanggal          string
       Jam              string
       Tempat           string
       WaliNikah        string
       Penghulu         string
       Kelurahan        string
       Keterangan       string
   }
   ```

2. **Helper Functions** - Ditambahkan:
   - `calculateAge()` - Menghitung usia dari tanggal lahir
   - `getDayName()` - Mengembalikan nama hari dalam bahasa Indonesia
   - `formatWaktu()` - Mengubah format waktu dari HH:MM ke HH.MM

3. **HTML Generation** - Diupdate untuk:
   - Layout landscape A4
   - Tabel 15 kolom dengan header kompleks
   - Kop surat dengan logo di kiri

### Files Modified
- `internal/handlers/staff/staff.go`
- `internal/handlers/kepala_kua/kepala_kua.go`

---

## 📝 Contoh Output

### Kop Surat
```
[LOGO]  KEMENTERIAN AGAMA REPUBLIK INDONESIA
        KANTOR KEMENTERIAN AGAMA KOTA BANJARMASIN
        KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA
        Jalan Brigjen H. Hasan Basri Komplek Kejaksaan RT.16 Banjarmasin 70124
        Telpon (0511) 3301966
        Kuabjmutara@gmail.com
```

### Judul
```
JADUAL NIKAH JANUARI 2026
```

### Tabel (Sample - Multiple Rows)

Tabel akan menampilkan **semua pendaftaran** dalam periode yang dipilih. Contoh jika ada 3 pendaftaran:

| NO | PRIA / BIN | USIA | PENDK | WANITA / BINTI | USIA | PENDK | HARI | TGL | JAM | TEMPAT | WALINIKAH | PENGHULU | KELURAHAN | KET |
|----|------------|------|-------|----------------|------|-------|------|-----|-----|---------|-----------|----------|-----------|-----|
| 1 | AHMAD RIFANI / SALIM | 29 | S1 | NUR AIDA ANGGRAINI / SADARUDDIN | 23 | S1 | KAMIS | 1 | 08.00 | KOMP DASAMAYA 2 BLOK D RT.16 RW.02 ALALAK SELATAN | SADARUDDIN | - | ALALAK SELATAN | - |
| 2 | MUHAMMAD RIZKY YULIANTO / SUPRAPTO (ALM) | 24 | S1 | ALFINA NUR ISLAMY SABILA RIADY / RIADI (ALM) | 24 | S1 | KAMIS | 1 | 09.00 | JL. TEMBUS PERUMNAS GG. NURUL FATA NO.4 | MUHAMMAD NUR ALAMSYAHBANA | - | ALALAK UTARA | SDR KDG |
| 3 | AKMAL MAULANA RAHMAN / MUHAMMAD IHSAN RAHMATILLAH | 26 | S1 | DINDA FINA SHOLEHA / JUANDA (ALM) | 26 | S1 | AHAD | 4 | 07.00 | GEDUNG HAFIYYUN | MUHAMMAD ALFI MUGRAHA | - | SURGI MUFTI | SDR KDG |

**Catatan:**
- Jumlah baris akan sesuai dengan jumlah pendaftaran dalam periode
- Jika ada 10 pendaftaran = 10 baris
- Jika ada 50 pendaftaran = 50 baris
- Data diurutkan berdasarkan tanggal nikah dan waktu (ASC)

---

## 🔄 Migration Guide

### Untuk Frontend Developer

**✅ TIDAK ADA PERUBAHAN SETUP YANG DIPERLUKAN**

Tidak ada perubahan pada API endpoint atau request/response format. Perubahan hanya pada format HTML output.

**Sebelum:**
- HTML dengan 9 kolom
- Format portrait A4
- Judul "PENGUMUMAN PERNIKAHAN"

**Sesudah:**
- HTML dengan 15 kolom
- Format landscape A4
- Judul "JADUAL NIKAH [BULAN] [TAHUN]"

**Action Required:**
- ✅ **Tidak ada action required untuk frontend**
- ✅ API endpoint, request format, response type **tetap sama**
- ✅ HTML output tetap bisa di-parse/ditampilkan dengan cara yang sama
- ✅ Print/PDF generation tetap berfungsi
- ⚠️ Opsional: Set PDF export orientation ke landscape untuk hasil lebih baik
- ⚠️ Pastikan printer/browser support landscape printing

**Lihat [Frontend Migration Guide](./FRONTEND_MIGRATION_GUIDE.md) untuk detail lengkap.**

---

## 📚 Dokumentasi Terkait

- [Fitur Terbaru](./FITUR_TERBARU.md) - Dokumentasi lengkap fitur generate pengumuman nikah
- [API Documentation Lengkap](./api/API_DOCUMENTATION_LENGKAP.md) - Dokumentasi endpoint lengkap
- [Parsing API Pengumuman Nikah](./api/PARSING_API_PENGUMUMAN_NIKAH.md) - Panduan parsing HTML response

---

## ✅ Testing Checklist

- [x] Format HTML sesuai Excel standar KUA
- [x] Tabel 15 kolom dengan header yang benar
- [x] Usia dihitung otomatis dengan benar
- [x] Format waktu HH.MM benar
- [x] Nama hari dalam bahasa Indonesia
- [x] Layout landscape A4
- [x] Print optimization CSS
- [x] Kop surat dengan logo di kiri
- [x] Judul dinamis berdasarkan bulan/tahun

---

**Last Updated:** Desember 2024  
**Version:** 2.0.0

