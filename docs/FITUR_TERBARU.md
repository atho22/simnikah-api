# 🆕 Dokumentasi Fitur Terbaru

**Versi:** 1.0.0  
**Tanggal Update:** Desember 2024  
**Status:** ✅ Production Ready

---

## 📋 Daftar Isi

1. [Overview](#overview)
2. [Fitur 1: Get Detail Pendaftaran by ID](#fitur-1-get-detail-pendaftaran-by-id)
3. [Fitur 2: Generate Pengumuman Nikah HTML](#fitur-2-generate-pengumuman-nikah-html)
4. [Quick Start Guide](#quick-start-guide)
5. [Contoh Implementasi Lengkap](#contoh-implementasi-lengkap)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)

---

## 🎯 Overview

Dokumentasi ini menjelaskan **2 fitur terbaru** yang telah ditambahkan ke API Simnikah:

1. **Get Detail Pendaftaran by ID** - Endpoint untuk mendapatkan detail lengkap pendaftaran nikah
2. **Generate Pengumuman Nikah HTML** - Endpoint untuk generate surat pengumuman nikah dalam format HTML siap cetak

**Base URL:** `https://your-api-domain.com/simnikah`

---

## 📝 Fitur 1: Get Detail Pendaftaran by ID

### Deskripsi

Endpoint ini memungkinkan frontend untuk mendapatkan **detail lengkap** dari sebuah pendaftaran nikah berdasarkan ID. Sangat berguna untuk halaman detail pendaftaran.

### Endpoint

```
GET /simnikah/pendaftaran/:id
```

### Authentication

✅ **Required** - Bearer Token

### Role Access

| Role | Akses |
|------|-------|
| `user_biasa` | Hanya bisa melihat pendaftaran miliknya sendiri |
| `staff` | Bisa melihat semua pendaftaran |
| `penghulu` | Bisa melihat semua pendaftaran |
| `kepala_kua` | Bisa melihat semua pendaftaran |

### URL Parameters

| Parameter | Type | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id` | integer | ✅ Yes | ID pendaftaran |

### Request Example

```javascript
// Menggunakan Axios
import axios from 'axios';

const getRegistrationDetail = async (registrationId) => {
  const token = localStorage.getItem('token');
  
  try {
    const response = await axios.get(
      `${API_BASE_URL}/pendaftaran/${registrationId}`,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    );
    
    return response.data;
  } catch (error) {
    console.error('Error fetching registration detail:', error);
    throw error;
  }
};

// Usage
const detail = await getRegistrationDetail(1);
console.log(detail.data);
```

### Response Success (200)

```json
{
  "success": true,
  "message": "Detail pendaftaran berhasil diambil",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241215-1234",
    "pendaftar_id": "USR1704067200",
    "status_pendaftaran": "Penghulu Ditugaskan",
    "tanggal_pendaftaran": "2024-01-01T10:00:00Z",
    "tanggal_nikah": "2024-12-15T00:00:00Z",
    "waktu_nikah": "10:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Ahmad Yani No. 123, Banjarmasin",
    "latitude": -3.291304,
    "longitude": 114.588147,
    "catatan": "Catatan tambahan",
    "disetujui_oleh": "USR1704067201",
    "disetujui_pada": "2024-01-02T10:00:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z",
    "calon_suami": {
      "id": 1,
      "user_id": "CP1704067200",
      "nik": "6301010101010001",
      "nama_lengkap": "Ahmad Wijaya bin Abdullah",
      "tanggal_lahir": "1999-01-01T00:00:00Z",
      "jenis_kelamin": "L",
      "pendidikan_terakhir": "S1",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "calon_istri": {
      "id": 2,
      "user_id": "CP1704067201",
      "nik": "6301010101010002",
      "nama_lengkap": "Siti Nurhaliza binti Muhammad",
      "tanggal_lahir": "2001-01-01T00:00:00Z",
      "jenis_kelamin": "P",
      "pendidikan_terakhir": "S1",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "wali_nikah": {
      "id": 1,
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "penghulu": {
      "id": 1,
      "user_id": "PEN1704067202",
      "nip": "198001012003121001",
      "nama_lengkap": "H. Muhammad Amin",
      "no_hp": "081234567890",
      "email": "amin@kua.go.id",
      "alamat": "Jl. Penghulu No. 123, Banjarmasin",
      "status": "Aktif",
      "ditugaskan_oleh": "USR1704067203",
      "ditugaskan_pada": "2024-12-10T10:00:00Z",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    "location": {
      "latitude": -3.291304,
      "longitude": 114.588147,
      "has_coordinates": true,
      "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.291304,114.588147",
      "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.291304,114.588147",
      "waze_url": "https://www.waze.com/ul?ll=-3.291304,114.588147&navigate=yes",
      "osm_url": "https://www.openstreetmap.org/?mlat=-3.291304&mlon=114.588147&zoom=16"
    }
  }
}
```

### Response Error

#### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token tidak valid atau tidak ada"
}
```

#### 403 Forbidden
```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak memiliki izin untuk melihat pendaftaran ini",
  "type": "authorization"
}
```

#### 404 Not Found
```json
{
  "success": false,
  "message": "Pendaftaran tidak ditemukan",
  "error": "Data pendaftaran dengan ID tersebut tidak ada"
}
```

### Use Cases

1. **Halaman Detail Pendaftaran** - Menampilkan semua informasi lengkap pendaftaran
2. **Preview Data** - Sebelum melakukan update atau action tertentu
3. **Print Detail** - Mencetak detail pendaftaran untuk keperluan administrasi
4. **Mobile App** - Menampilkan detail pendaftaran di aplikasi mobile

---

## 📄 Fitur 2: Generate Pengumuman Nikah HTML

### Deskripsi

Endpoint ini mengembalikan **HTML document lengkap** untuk surat pengumuman nikah yang siap dicetak atau dikonversi ke PDF. Format HTML mengikuti format Excel standar KUA dengan **15 kolom** dan layout **landscape A4**. HTML sudah include CSS untuk print optimization.

### Endpoints

**Staff & Kepala KUA:**
```
GET /simnikah/staff/pengumuman-nikah/generate
POST /simnikah/staff/pengumuman-nikah/generate
```

**Kepala KUA Only:**
```
GET /simnikah/kepala-kua/pengumuman-nikah/generate
POST /simnikah/kepala-kua/pengumuman-nikah/generate
```

> ⚠️ **PENTING:** Jika Anda mendapatkan error 404, pastikan server backend sudah di-restart setelah route ditambahkan. Route baru tidak akan aktif sampai server di-restart.

### Authentication

✅ **Required** - Bearer Token

### Role Access

| Endpoint | Role Required |
|----------|---------------|
| `/staff/pengumuman-nikah/generate` | `staff`, `kepala_kua` |
| `/kepala-kua/pengumuman-nikah/generate` | `kepala_kua` |

### Query Parameters

| Parameter | Type | Required | Default | Deskripsi |
|-----------|------|----------|---------|-----------|
| `tanggal_awal` | string (YYYY-MM-DD) | ❌ No | Awal minggu ini (Senin) | Tanggal awal periode |
| `tanggal_akhir` | string (YYYY-MM-DD) | ❌ No | Akhir minggu ini (Minggu) | Tanggal akhir periode |

### Request Body (Optional)

Untuk custom kop surat, kirim JSON body dengan method `POST`:

```json
{
  "tanggal_awal": "2024-12-16",
  "tanggal_akhir": "2024-12-22",
  "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
  "alamat_kua": "PH5Q+F8C, Jl. Wira Karya, Pangeran",
  "kota": "Kota Banjarmasin",
  "provinsi": "Kalimantan Selatan",
  "kode_pos": "70123",
  "telepon": "0511-1234567",
  "email": "kua.banjarmasinutara@kemenag.go.id",
  "website": "https://kua.banjarmasinutara.go.id",
  "logo_url": "https://example.com/logo-kua.png"
}
```

**Catatan:**
- Query parameters memiliki **prioritas lebih tinggi** daripada request body
- Jika query parameters tidak ada, akan membaca dari request body
- Jika keduanya tidak ada, akan menggunakan default (minggu ini)
- Logo URL harus accessible (public URL) dan disarankan menggunakan HTTPS

### Request Examples

#### 1. Generate dengan default kop surat (GET)

```javascript
const generatePengumuman = async (tanggalAwal, tanggalAkhir) => {
  const token = localStorage.getItem('token');
  
  try {
    const response = await axios.get(
      `${API_BASE_URL}/staff/pengumuman-nikah/generate`,
      {
        params: {
          tanggal_awal: tanggalAwal,
          tanggal_akhir: tanggalAkhir
        },
        headers: {
          'Authorization': `Bearer ${token}`
        },
        responseType: 'text' // PENTING: untuk mendapatkan HTML sebagai string
      }
    );
    
    return response.data; // HTML string
  } catch (error) {
    console.error('Error generating pengumuman:', error);
    throw error;
  }
};
```

#### 2. Generate dengan custom kop surat (POST)

```javascript
const generatePengumumanCustom = async (tanggalAwal, tanggalAkhir, kopSurat) => {
  const token = localStorage.getItem('token');
  
  try {
    const response = await axios.post(
      `${API_BASE_URL}/staff/pengumuman-nikah/generate`,
      {
        tanggal_awal: tanggalAwal,
        tanggal_akhir: tanggalAkhir,
        ...kopSurat
      },
      {
        params: {
          tanggal_awal: tanggalAwal,
          tanggal_akhir: tanggalAkhir
        },
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        responseType: 'text' // PENTING: untuk mendapatkan HTML sebagai string
      }
    );
    
    return response.data; // HTML string
  } catch (error) {
    console.error('Error generating pengumuman:', error);
    throw error;
  }
};
```

### Response Success (200)

**Content-Type:** `text/html; charset=utf-8`

Response berupa **HTML document lengkap** dengan struktur sesuai format Excel standar KUA:

1. **Kop Surat KUA:**
   - Logo di kiri (jika `logo_url` disediakan)
   - KEMENTERIAN AGAMA REPUBLIK INDONESIA
   - KANTOR KEMENTERIAN AGAMA KOTA [KOTA]
   - Nama KUA
   - Alamat lengkap
   - Kontak (telepon, email)

2. **Judul:** "JADUAL NIKAH [BULAN] [TAHUN]" (contoh: "JADUAL NIKAH JANUARI 2026")

3. **Tabel Data Pendaftaran** dengan **15 kolom** (kolom tetap, baris dinamis):
   
   **Struktur Kolom (Tetap 15 kolom):**
   - **NO URUT** - Nomor urut (1, 2, 3, ...)
   - **DATA CALON PENGANTIN:**
     - **PRIA / BIN** - Nama calon suami
     - **USIA** - Usia calon suami (dihitung otomatis)
     - **PENDK** - Pendidikan terakhir calon suami
     - **WANITA / BINTI** - Nama calon istri
     - **USIA** - Usia calon istri (dihitung otomatis)
     - **PENDK** - Pendidikan terakhir calon istri
   - **PELAKSANAAN NIKAH:**
     - **HARI** - Nama hari (SENIN, SELASA, RABU, KAMIS, JUM'AT, SABTU, AHAD)
     - **TGL** - Tanggal (hanya angka: 1, 2, 3, dll)
     - **JAM** - Waktu (format: 08.00, 09.00, dll)
     - **TEMPAT** - Tempat nikah
     - **WALINIKAH** - Nama wali nikah
     - **PENGHULU** - Nama penghulu (jika sudah ditugaskan)
     - **KELURAHAN** - Kelurahan
     - **KET** - Keterangan
   
   **Catatan Penting:**
   - ✅ **Kolom tetap 15 kolom** (tidak berubah)
   - ✅ **Baris dinamis** - jumlah baris tergantung jumlah data pendaftaran dalam periode yang dipilih
   - ✅ **Setiap baris = 1 pendaftaran nikah**
   - ✅ Jika ada 10 pendaftaran, akan ada 10 baris data
   - ✅ Jika ada 50 pendaftaran, akan ada 50 baris data
   - ✅ Data diurutkan berdasarkan tanggal nikah dan waktu (ASC)

**Fitur HTML:**
- ✅ CSS untuk print optimization (`@media print`)
- ✅ Format **A4 Landscape** untuk menampung tabel lebar
- ✅ Font: Times New Roman (serif)
- ✅ Font size: 8-9pt untuk tabel, 11-12pt untuk kop surat
- ✅ Tabel dengan border untuk kejelasan
- ✅ Header tabel dengan rowspan/colspan untuk grouping kolom
- ✅ Siap untuk dicetak langsung atau dikonversi ke PDF

**Catatan Penting:**
- **Status yang ditampilkan:** Semua status kecuali "Ditolak" (Draft, Disetujui, Menunggu Penugasan, Penghulu Ditugaskan, Selesai)
- Status "Ditolak" tidak ditampilkan karena tidak perlu diumumkan
- **Usia dihitung otomatis** dari tanggal lahir calon pengantin
- **Format waktu:** HH:MM diubah menjadi HH.MM (contoh: 08:00 → 08.00)
- **Nama hari:** Otomatis dalam bahasa Indonesia (SENIN, SELASA, dll)

### Response Error

#### 400 Bad Request
```json
{
  "success": false,
  "message": "Format tanggal tidak valid",
  "error": "Format tanggal_awal harus YYYY-MM-DD"
}
```

#### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token tidak valid atau tidak ada"
}
```

#### 403 Forbidden
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Role tidak memiliki akses"
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Database error",
  "error": "Gagal mengambil data pendaftaran"
}
```

### Use Cases

1. **Print Pengumuman** - Mencetak surat pengumuman untuk dipasang di papan pengumuman KUA
2. **Export PDF** - Mengkonversi HTML ke PDF untuk arsip digital
3. **Preview** - Preview surat sebelum dicetak
4. **Email** - Mengirim surat pengumuman via email (setelah dikonversi ke PDF)

---

## 🚀 Quick Start Guide

### Setup API Client

```javascript
// api/pengumuman.js
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'https://your-api-domain.com/simnikah';

// Get token dari localStorage atau context
const getToken = () => {
  return localStorage.getItem('token');
};

// Get Detail Pendaftaran
export const getRegistrationDetail = async (registrationId) => {
  const token = getToken();
  
  const response = await axios.get(
    `${API_BASE_URL}/pendaftaran/${registrationId}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    }
  );
  
  return response.data;
};

// Generate Pengumuman HTML
export const generatePengumumanHTML = async (tanggalAwal, tanggalAkhir, kopSurat = null) => {
  const token = getToken();
  
  const config = {
    method: kopSurat ? 'post' : 'get',
    url: `${API_BASE_URL}/staff/pengumuman-nikah/generate`,
    params: {
      tanggal_awal: tanggalAwal,
      tanggal_akhir: tanggalAkhir
    },
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    responseType: 'text' // PENTING: untuk mendapatkan HTML sebagai string
  };

  if (kopSurat) {
    config.data = {
      tanggal_awal: tanggalAwal,
      tanggal_akhir: tanggalAkhir,
      ...kopSurat
    };
  }

  const response = await axios(config);
  return response.data; // HTML string
};
```

---

## 💻 Contoh Implementasi Lengkap

### React Component: Detail Pendaftaran

```jsx
import React, { useState, useEffect } from 'react';
import { getRegistrationDetail } from '../api/pengumuman';

const RegistrationDetail = ({ registrationId }) => {
  const [detail, setDetail] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDetail = async () => {
      try {
        setLoading(true);
        const response = await getRegistrationDetail(registrationId);
        setDetail(response.data);
      } catch (err) {
        setError(err.response?.data?.message || 'Gagal mengambil detail pendaftaran');
      } finally {
        setLoading(false);
      }
    };

    if (registrationId) {
      fetchDetail();
    }
  }, [registrationId]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!detail) return <div>Data tidak ditemukan</div>;

  return (
    <div className="registration-detail">
      <h2>Detail Pendaftaran: {detail.nomor_pendaftaran}</h2>
      
      <div className="info-section">
        <h3>Informasi Pendaftaran</h3>
        <p><strong>Status:</strong> {detail.status_pendaftaran}</p>
        <p><strong>Tanggal Nikah:</strong> {new Date(detail.tanggal_nikah).toLocaleDateString('id-ID')}</p>
        <p><strong>Waktu:</strong> {detail.waktu_nikah}</p>
        <p><strong>Tempat:</strong> {detail.tempat_nikah}</p>
        <p><strong>Alamat Akad:</strong> {detail.alamat_akad}</p>
      </div>

      <div className="info-section">
        <h3>Calon Suami</h3>
        <p><strong>Nama:</strong> {detail.calon_suami.nama_lengkap}</p>
        <p><strong>NIK:</strong> {detail.calon_suami.nik}</p>
        <p><strong>Tanggal Lahir:</strong> {new Date(detail.calon_suami.tanggal_lahir).toLocaleDateString('id-ID')}</p>
      </div>

      <div className="info-section">
        <h3>Calon Istri</h3>
        <p><strong>Nama:</strong> {detail.calon_istri.nama_lengkap}</p>
        <p><strong>NIK:</strong> {detail.calon_istri.nik}</p>
        <p><strong>Tanggal Lahir:</strong> {new Date(detail.calon_istri.tanggal_lahir).toLocaleDateString('id-ID')}</p>
      </div>

      {detail.wali_nikah && (
        <div className="info-section">
          <h3>Wali Nikah</h3>
          <p><strong>Nama:</strong> {detail.wali_nikah.nama_dan_bin}</p>
          <p><strong>Hubungan:</strong> {detail.wali_nikah.hubungan_wali}</p>
        </div>
      )}

      {detail.penghulu && (
        <div className="info-section">
          <h3>Penghulu</h3>
          <p><strong>Nama:</strong> {detail.penghulu.nama_lengkap}</p>
          <p><strong>NIP:</strong> {detail.penghulu.nip}</p>
          <p><strong>No. HP:</strong> {detail.penghulu.no_hp}</p>
        </div>
      )}

      {detail.location?.has_coordinates && (
        <div className="info-section">
          <h3>Lokasi</h3>
          <a 
            href={detail.location.google_maps_url} 
            target="_blank" 
            rel="noopener noreferrer"
          >
            Buka di Google Maps
          </a>
        </div>
      )}
    </div>
  );
};

export default RegistrationDetail;
```

### React Component: Generate & Print Pengumuman

```jsx
import React, { useState } from 'react';
import { generatePengumumanHTML } from '../api/pengumuman';
import html2pdf from 'html2pdf.js';

const PengumumanGenerator = () => {
  const [tanggalAwal, setTanggalAwal] = useState('');
  const [tanggalAkhir, setTanggalAkhir] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Custom kop surat (opsional)
  const [kopSurat, setKopSurat] = useState({
    nama_kua: '',
    alamat_kua: '',
    kota: '',
    provinsi: '',
    kode_pos: '',
    telepon: '',
    email: '',
    website: '',
    logo_url: ''
  });

  const handlePrint = async () => {
    try {
      setLoading(true);
      setError(null);

      // Generate HTML
      const html = await generatePengumumanHTML(
        tanggalAwal,
        tanggalAkhir,
        Object.values(kopSurat).some(v => v) ? kopSurat : null
      );

      // Buka di window baru untuk print
      const printWindow = window.open('', '_blank');
      printWindow.document.write(html);
      printWindow.document.close();
      
      // Tunggu sebentar untuk memastikan HTML sudah ter-render
      setTimeout(() => {
        printWindow.print();
        setLoading(false);
      }, 500);
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal generate pengumuman');
      setLoading(false);
    }
  };

  const handleExportPDF = async () => {
    try {
      setLoading(true);
      setError(null);

      // Generate HTML
      const html = await generatePengumumanHTML(
        tanggalAwal,
        tanggalAkhir,
        Object.values(kopSurat).some(v => v) ? kopSurat : null
      );

      // Buat element untuk konversi
      const element = document.createElement('div');
      element.innerHTML = html;

      // Konfigurasi PDF
      const opt = {
        margin: 1,
        filename: `pengumuman-nikah-${tanggalAwal}-${tanggalAkhir}.pdf`,
        image: { type: 'jpeg', quality: 0.98 },
        html2canvas: { scale: 2, useCORS: true },
        jsPDF: { 
          unit: 'in', 
          format: 'a4', 
          orientation: 'portrait' 
        }
      };

      // Generate dan download PDF
      await html2pdf().set(opt).from(element).save();
      setLoading(false);
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal export PDF');
      setLoading(false);
    }
  };

  const handlePreview = async () => {
    try {
      setLoading(true);
      setError(null);

      // Generate HTML
      const html = await generatePengumumanHTML(
        tanggalAwal,
        tanggalAkhir,
        Object.values(kopSurat).some(v => v) ? kopSurat : null
      );

      // Tampilkan di iframe
      const iframe = document.createElement('iframe');
      iframe.style.width = '100%';
      iframe.style.height = '800px';
      iframe.style.border = 'none';
      
      document.body.appendChild(iframe);
      iframe.contentDocument.write(html);
      iframe.contentDocument.close();
      
      setLoading(false);
    } catch (err) {
      setError(err.response?.data?.message || 'Gagal preview pengumuman');
      setLoading(false);
    }
  };

  return (
    <div className="pengumuman-generator">
      <h2>Generate Pengumuman Nikah</h2>

      <div className="form-group">
        <label>Tanggal Awal:</label>
        <input
          type="date"
          value={tanggalAwal}
          onChange={(e) => setTanggalAwal(e.target.value)}
        />
      </div>

      <div className="form-group">
        <label>Tanggal Akhir:</label>
        <input
          type="date"
          value={tanggalAkhir}
          onChange={(e) => setTanggalAkhir(e.target.value)}
        />
      </div>

      {/* Form Custom Kop Surat (opsional) */}
      <div className="kop-surat-section">
        <h3>Custom Kop Surat (Opsional)</h3>
        <input
          type="text"
          placeholder="Nama KUA"
          value={kopSurat.nama_kua}
          onChange={(e) => setKopSurat({ ...kopSurat, nama_kua: e.target.value })}
        />
        <input
          type="text"
          placeholder="Alamat KUA"
          value={kopSurat.alamat_kua}
          onChange={(e) => setKopSurat({ ...kopSurat, alamat_kua: e.target.value })}
        />
        <input
          type="url"
          placeholder="Logo URL (https://...)"
          value={kopSurat.logo_url}
          onChange={(e) => setKopSurat({ ...kopSurat, logo_url: e.target.value })}
        />
        {/* Tambahkan field lainnya sesuai kebutuhan */}
      </div>

      {error && <div className="error">{error}</div>}

      <div className="actions">
        <button 
          onClick={handlePreview} 
          disabled={loading || !tanggalAwal || !tanggalAkhir}
        >
          {loading ? 'Loading...' : 'Preview'}
        </button>
        <button 
          onClick={handlePrint} 
          disabled={loading || !tanggalAwal || !tanggalAkhir}
        >
          {loading ? 'Loading...' : 'Print'}
        </button>
        <button 
          onClick={handleExportPDF} 
          disabled={loading || !tanggalAwal || !tanggalAkhir}
        >
          {loading ? 'Loading...' : 'Export PDF'}
        </button>
      </div>
    </div>
  );
};

export default PengumumanGenerator;
```

---

## ✅ Best Practices

### 1. Error Handling

```javascript
try {
  const detail = await getRegistrationDetail(id);
  // Handle success
} catch (error) {
  if (error.response?.status === 401) {
    // Token expired, redirect to login
    window.location.href = '/login';
  } else if (error.response?.status === 403) {
    // Access denied
    alert('Anda tidak memiliki izin untuk melihat pendaftaran ini');
  } else if (error.response?.status === 404) {
    // Not found
    alert('Pendaftaran tidak ditemukan');
  } else {
    // Other errors
    console.error('Error:', error);
    alert('Terjadi kesalahan. Silakan coba lagi.');
  }
}
```

### 2. Loading States

```javascript
const [loading, setLoading] = useState(false);

// Show loading indicator
{loading && <div>Loading...</div>}

// Disable buttons during loading
<button disabled={loading}>Generate</button>
```

### 3. Date Formatting

```javascript
// Format tanggal untuk display
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

// Format tanggal untuk API (YYYY-MM-DD)
const formatDateForAPI = (date) => {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};
```

### 4. Response Type untuk HTML

```javascript
// PENTING: Set responseType ke 'text' untuk mendapatkan HTML sebagai string
const response = await axios.get(url, {
  responseType: 'text' // Bukan 'json'!
});
```

### 5. Print Optimization

```javascript
// Tunggu HTML ter-render sebelum print
const printWindow = window.open('', '_blank');
printWindow.document.write(html);
printWindow.document.close();

// Tunggu sebentar untuk memastikan CSS dan images sudah ter-load
setTimeout(() => {
  printWindow.print();
}, 500);
```

### 6. PDF Export dengan html2pdf.js

```javascript
// Install: npm install html2pdf.js
import html2pdf from 'html2pdf.js';

// Konfigurasi untuk hasil terbaik
const opt = {
  margin: 1,
  filename: 'pengumuman-nikah.pdf',
  image: { type: 'jpeg', quality: 0.98 },
  html2canvas: { 
    scale: 2, 
    useCORS: true, // Penting untuk load images dari external URL
    logging: false 
  },
  jsPDF: { 
    unit: 'in', 
    format: 'a4', 
    orientation: 'portrait' 
  }
};

await html2pdf().set(opt).from(element).save();
```

---

## 🔧 Troubleshooting

### Problem 1: Error 404 - Endpoint tidak ditemukan

**Gejala:**
```
Failed to load resource: the server responded with a status of 404 (Not Found)
{"error":"Path '/simnikah/kepala-kua/pengumuman-nikah/generate' tidak ditemukan"}
```

**Solusi:**
1. **Restart Server Backend** - Route baru memerlukan restart server:
   ```bash
   # Stop server yang sedang berjalan (Ctrl+C)
   # Kemudian start ulang
   go run cmd/api/main.go
   # atau jika menggunakan build
   ./simnikah-api
   ```

2. **Verifikasi Route Terdaftar** - Pastikan route sudah terdaftar di `cmd/api/main.go`:
   ```go
   simnikahRoutes.GET("/kepala-kua/pengumuman-nikah/generate", ...)
   simnikahRoutes.POST("/kepala-kua/pengumuman-nikah/generate", ...)
   ```

3. **Cek Proxy Configuration** - Jika menggunakan proxy di frontend, pastikan konfigurasi proxy sudah benar:
   ```javascript
   // Contoh: Jika menggunakan proxy /api/proxy
   // Pastikan proxy mengarah ke backend yang benar
   const API_BASE_URL = '/api/proxy/simnikah';
   ```

4. **Test dengan curl atau Postman** - Test langsung ke backend tanpa proxy:
   ```bash
   curl -X GET "http://localhost:8080/simnikah/kepala-kua/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22" \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

5. **Cek Log Server** - Lihat log server untuk melihat apakah request sampai ke server:
   ```
   [GIN] 2024/12/16 - 10:00:00 | 404 | 1ms | 127.0.0.1 | GET /simnikah/kepala-kua/pengumuman-nikah/generate
   ```

### Problem 2: HTML tidak ter-render dengan benar

**Solusi:**
- Pastikan `responseType: 'text'` sudah diset
- Pastikan HTML di-write ke document dengan benar
- Tunggu beberapa saat sebelum print untuk memastikan CSS ter-load

### Problem 3: Logo tidak muncul di PDF

**Solusi:**
- Pastikan logo URL accessible (public URL)
- Gunakan HTTPS untuk logo URL
- Set `useCORS: true` di html2canvas config
- Pastikan server logo mengizinkan CORS

### Problem 4: Error 403 Forbidden

**Solusi:**
- Pastikan user memiliki role yang sesuai (`staff` atau `kepala_kua`)
- Pastikan token masih valid
- Cek apakah user memiliki akses ke endpoint tersebut

### Problem 5: Format tanggal tidak valid

**Solusi:**
- Pastikan format tanggal adalah `YYYY-MM-DD` (contoh: `2024-12-16`)
- Gunakan date picker yang otomatis format ke format yang benar
- Validasi tanggal di frontend sebelum kirim request

### Problem 6: PDF tidak sesuai dengan print preview

**Solusi:**
- Pastikan `@media print` CSS sudah include di HTML
- Set `scale: 2` di html2canvas untuk kualitas lebih baik
- Pastikan format A4 sudah diset di jsPDF config

---

## 📚 Referensi

- [API Documentation Lengkap](./api/API_DOCUMENTATION_LENGKAP.md)
- [Parsing API Pengumuman Nikah](./api/PARSING_API_PENGUMUMAN_NIKAH.md)
- [Axios Documentation](https://axios-http.com/)
- [html2pdf.js Documentation](https://ekoopmans.github.io/html2pdf.js/)

---

## 📞 Support

Jika ada pertanyaan atau masalah, silakan hubungi tim development atau buat issue di repository project.

---

**Last Updated:** Desember 2024  
**Version:** 1.0.0


