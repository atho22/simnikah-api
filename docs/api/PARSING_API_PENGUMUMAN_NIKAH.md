# 📄 Dokumentasi Parsing API Surat Pengumuman Nikah

**Versi:** 1.0.0  
**Tanggal:** Desember 2024  
**Status:** ✅ Production Ready

---

## 📋 Daftar Isi

1. [Overview](#overview)
2. [Endpoint List](#endpoint-list)
3. [Authentication](#authentication)
4. [Get Approved Registrations Per Week](#get-approved-registrations-per-week)
5. [Generate Pengumuman Nikah](#generate-pengumuman-nikah)
6. [Parsing HTML Response](#parsing-html-response)
7. [Custom Kop Surat](#custom-kop-surat)
8. [Error Handling](#error-handling)
9. [Contoh Implementasi](#contoh-implementasi)

---

## 🎯 Overview

API Surat Pengumuman Nikah memungkinkan **Staff** dan **Kepala KUA** untuk:
- Mengambil daftar pendaftaran nikah yang telah disetujui per minggu
- Generate surat pengumuman nikah dalam format HTML yang siap dicetak atau dikonversi ke PDF
- Mengkustomisasi kop surat (logo, alamat, kontak, dll)

**Base URL:** `https://your-api-domain.com/simnikah`

---

## 📍 Endpoint List

| Method | Endpoint | Role | Deskripsi |
|--------|----------|------|-----------|
| `GET` | `/staff/pengumuman-nikah/list` | `staff`, `kepala_kua` | Daftar pendaftaran disetujui per minggu |
| `GET` | `/staff/pengumuman-nikah/generate` | `staff`, `kepala_kua` | Generate surat pengumuman (HTML) |
| `GET` | `/kepala-kua/pengumuman-nikah/list` | `kepala_kua` | Daftar pendaftaran disetujui per minggu |
| `GET` | `/kepala-kua/pengumuman-nikah/generate` | `kepala_kua` | Generate surat pengumuman (HTML) |

---

## 🔐 Authentication

Semua endpoint memerlukan authentication token. Tambahkan header berikut:

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

**Contoh:**
```javascript
const headers = {
  'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
  'Content-Type': 'application/json'
};
```

---

## 📊 Get Approved Registrations Per Week

### Endpoint

```
GET /simnikah/staff/pengumuman-nikah/list
GET /simnikah/kepala-kua/pengumuman-nikah/list
```

### Query Parameters

| Parameter | Type | Required | Default | Deskripsi |
|-----------|------|----------|---------|-----------|
| `tanggal_awal` | string (YYYY-MM-DD) | ❌ No | Awal minggu ini (Senin) | Tanggal awal periode |
| `tanggal_akhir` | string (YYYY-MM-DD) | ❌ No | Akhir minggu ini (Minggu) | Tanggal akhir periode |

### Request Example

```bash
# Menggunakan periode default (minggu ini)
GET /simnikah/staff/pengumuman-nikah/list
Authorization: Bearer YOUR_TOKEN

# Menggunakan periode custom
GET /simnikah/staff/pengumuman-nikah/list?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
Authorization: Bearer YOUR_TOKEN
```

### Response Success (200)

```json
{
  "success": true,
  "message": "Data pendaftaran disetujui berhasil diambil",
  "data": {
    "tanggal_awal": "2024-12-16",
    "tanggal_akhir": "2024-12-22",
    "periode": "16 Desember 2024 s/d 22 Desember 2024",
    "total": 5,
    "registrations": [
      {
        "id": 1,
        "nomor_pendaftaran": "NIKAH-20241215-1234",
        "tanggal_nikah": "2024-12-18T10:00:00Z",
        "waktu_nikah": "10:00",
        "tempat_nikah": "Di KUA",
        "alamat_akad": "Kantor KUA Kecamatan Banjarmasin Utara",
        "calon_suami": {
          "nama_lengkap": "Ahmad bin Abdullah"
        },
        "calon_istri": {
          "nama_lengkap": "Siti binti Muhammad"
        },
        "wali_nikah": {
          "nama_dan_bin": "Muhammad bin Ali",
          "hubungan_wali": "Ayah Kandung"
        }
      }
    ]
  }
}
```

### Response Fields

| Field | Type | Deskripsi |
|-------|------|-----------|
| `success` | boolean | Status sukses |
| `message` | string | Pesan response |
| `data.tanggal_awal` | string | Tanggal awal periode (YYYY-MM-DD) |
| `data.tanggal_akhir` | string | Tanggal akhir periode (YYYY-MM-DD) |
| `data.periode` | string | Periode dalam format Indonesia |
| `data.total` | number | Total pendaftaran |
| `data.registrations[]` | array | Array pendaftaran |
| `data.registrations[].id` | number | ID pendaftaran |
| `data.registrations[].nomor_pendaftaran` | string | Nomor pendaftaran |
| `data.registrations[].tanggal_nikah` | string (ISO 8601) | Tanggal nikah |
| `data.registrations[].waktu_nikah` | string (HH:MM) | Waktu nikah |
| `data.registrations[].tempat_nikah` | string | Tempat nikah (Di KUA / Di Luar KUA) |
| `data.registrations[].alamat_akad` | string | Alamat lengkap akad nikah |
| `data.registrations[].calon_suami.nama_lengkap` | string | Nama calon suami |
| `data.registrations[].calon_istri.nama_lengkap` | string | Nama calon istri |
| `data.registrations[].wali_nikah.nama_dan_bin` | string | Nama wali nikah dengan bin |
| `data.registrations[].wali_nikah.hubungan_wali` | string | Hubungan wali |

---

## 📝 Generate Pengumuman Nikah

### Endpoint

```
GET /simnikah/staff/pengumuman-nikah/generate
GET /simnikah/kepala-kua/pengumuman-nikah/generate
```

### Query Parameters

| Parameter | Type | Required | Default | Deskripsi |
|-----------|------|----------|---------|-----------|
| `tanggal_awal` | string (YYYY-MM-DD) | ❌ No | Awal minggu ini (Senin) | Tanggal awal periode |
| `tanggal_akhir` | string (YYYY-MM-DD) | ❌ No | Akhir minggu ini (Minggu) | Tanggal akhir periode |

### Request Body (Optional - untuk custom kop surat)

Jika ingin mengkustomisasi kop surat, kirim JSON body dengan method `GET` (atau gunakan `POST` jika lebih sesuai):

```json
{
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

**Note:** Jika request body tidak dikirim, akan menggunakan nilai default.

### Request Example

#### 1. Generate dengan default kop surat

```bash
GET /simnikah/staff/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
Authorization: Bearer YOUR_TOKEN
```

#### 2. Generate dengan custom kop surat (menggunakan POST)

```bash
POST /simnikah/staff/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "nama_kua": "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
  "alamat_kua": "Jl. Contoh No. 123",
  "kota": "Kota Banjarmasin",
  "provinsi": "Kalimantan Selatan",
  "kode_pos": "70123",
  "telepon": "0511-1234567",
  "email": "kua@example.com",
  "website": "https://kua.example.com",
  "logo_url": "https://example.com/logo.png"
}
```

### Response Success (200)

**Content-Type:** `text/html; charset=utf-8`

Response berupa HTML document yang siap dicetak atau dikonversi ke PDF.

**Response Headers:**
```
Content-Type: text/html; charset=utf-8
Content-Length: 12345
```

**Response Body:** HTML document lengkap dengan:
- Kop surat KUA (dengan logo jika disediakan)
- Nomor surat
- Judul "PENGUMUMAN PERNIKAHAN"
- Periode
- Paragraf pembuka
- Tabel data pendaftaran
- Paragraf penutup
- Tanda tangan Kepala KUA

---

## 🔍 Parsing HTML Response

### Menggunakan JavaScript/TypeScript

```javascript
// Fetch API
async function generatePengumuman(tanggalAwal, tanggalAkhir, kopSurat = null) {
  const url = `/simnikah/staff/pengumuman-nikah/generate?tanggal_awal=${tanggalAwal}&tanggal_akhir=${tanggalAkhir}`;
  
  const options = {
    method: kopSurat ? 'POST' : 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  };
  
  if (kopSurat) {
    options.body = JSON.stringify(kopSurat);
  }
  
  const response = await fetch(url, options);
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  // Get HTML content
  const html = await response.text();
  
  // Option 1: Display in iframe
  const iframe = document.createElement('iframe');
  iframe.srcdoc = html;
  document.body.appendChild(iframe);
  
  // Option 2: Open in new window for printing
  const printWindow = window.open('', '_blank');
  printWindow.document.write(html);
  printWindow.document.close();
  printWindow.print();
  
  // Option 3: Download as HTML file
  const blob = new Blob([html], { type: 'text/html' });
  const url_blob = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url_blob;
  a.download = `pengumuman-nikah-${tanggalAwal}-${tanggalAkhir}.html`;
  a.click();
  
  return html;
}

// Usage
generatePengumuman('2024-12-16', '2024-12-22', {
  nama_kua: "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
  alamat_kua: "Jl. Contoh No. 123",
  kota: "Kota Banjarmasin",
  provinsi: "Kalimantan Selatan",
  kode_pos: "70123",
  telepon: "0511-1234567",
  email: "kua@example.com"
});
```

### Menggunakan Axios

```javascript
import axios from 'axios';

async function generatePengumuman(tanggalAwal, tanggalAkhir, kopSurat = null) {
  const config = {
    method: kopSurat ? 'post' : 'get',
    url: `/simnikah/staff/pengumuman-nikah/generate`,
    params: {
      tanggal_awal: tanggalAwal,
      tanggal_akhir: tanggalAkhir
    },
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    responseType: 'text' // Important: untuk mendapatkan HTML sebagai string
  };
  
  if (kopSurat) {
    config.data = kopSurat;
  }
  
  try {
    const response = await axios(config);
    return response.data; // HTML string
  } catch (error) {
    console.error('Error generating pengumuman:', error);
    throw error;
  }
}
```

### Menggunakan Fetch dengan React

```jsx
import React, { useState } from 'react';

function PengumumanGenerator() {
  const [html, setHtml] = useState('');
  const [loading, setLoading] = useState(false);
  
  const generatePengumuman = async () => {
    setLoading(true);
    try {
      const response = await fetch(
        '/simnikah/staff/pengumuman-nikah/generate?tanggal_awal=2024-12-16&tanggal_akhir=2024-12-22',
        {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        }
      );
      
      if (!response.ok) {
        throw new Error('Failed to generate pengumuman');
      }
      
      const htmlContent = await response.text();
      setHtml(htmlContent);
    } catch (error) {
      console.error(error);
      alert('Gagal generate pengumuman');
    } finally {
      setLoading(false);
    }
  };
  
  const printPengumuman = () => {
    if (!html) return;
    
    const printWindow = window.open('', '_blank');
    printWindow.document.write(html);
    printWindow.document.close();
    printWindow.print();
  };
  
  return (
    <div>
      <button onClick={generatePengumuman} disabled={loading}>
        {loading ? 'Generating...' : 'Generate Pengumuman'}
      </button>
      
      {html && (
        <>
          <button onClick={printPengumuman}>Print</button>
          <iframe
            srcDoc={html}
            style={{ width: '100%', height: '800px', border: '1px solid #ccc' }}
            title="Pengumuman Nikah"
          />
        </>
      )}
    </div>
  );
}
```

### Konversi ke PDF (Client-side)

#### Menggunakan html2pdf.js

```javascript
import html2pdf from 'html2pdf.js';

async function generatePDF(html) {
  const options = {
    margin: [10, 10, 10, 10],
    filename: 'pengumuman-nikah.pdf',
    image: { type: 'jpeg', quality: 0.98 },
    html2canvas: { scale: 2 },
    jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' }
  };
  
  await html2pdf().set(options).from(html).save();
}

// Usage
const html = await generatePengumuman('2024-12-16', '2024-12-22');
await generatePDF(html);
```

#### Menggunakan Puppeteer (Server-side)

```javascript
const puppeteer = require('puppeteer');

async function generatePDF(html) {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  
  await page.setContent(html, { waitUntil: 'networkidle0' });
  
  const pdf = await page.pdf({
    format: 'A4',
    margin: {
      top: '15mm',
      right: '20mm',
      bottom: '15mm',
      left: '20mm'
    }
  });
  
  await browser.close();
  return pdf;
}
```

---

## 🎨 Custom Kop Surat

### Field Kop Surat

| Field | Type | Required | Default | Deskripsi |
|-------|------|----------|---------|-----------|
| `nama_kua` | string | ❌ No | "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA" | Nama lengkap KUA |
| `alamat_kua` | string | ❌ No | "PH5Q+F8C, Jl. Wira Karya, Pangeran" | Alamat lengkap KUA |
| `kota` | string | ❌ No | "Kota Banjarmasin" | Kota |
| `provinsi` | string | ❌ No | "Kalimantan Selatan" | Provinsi |
| `kode_pos` | string | ❌ No | "70123" | Kode pos |
| `telepon` | string | ❌ No | "-" | Nomor telepon |
| `email` | string | ❌ No | "-" | Email |
| `website` | string | ❌ No | "" | Website (optional) |
| `logo_url` | string | ❌ No | "" | URL logo KUA (optional) |

### Contoh Custom Kop Surat

```javascript
const customKopSurat = {
  nama_kua: "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA",
  alamat_kua: "Jl. Ahmad Yani No. 123, Kelurahan Pangeran",
  kota: "Kota Banjarmasin",
  provinsi: "Kalimantan Selatan",
  kode_pos: "70123",
  telepon: "0511-1234567",
  email: "kua.banjarmasinutara@kemenag.go.id",
  website: "https://kua.banjarmasinutara.go.id",
  logo_url: "https://cdn.example.com/logo-kua.png" // URL harus accessible
};

// Generate dengan custom kop surat
const html = await generatePengumuman(
  '2024-12-16',
  '2024-12-22',
  customKopSurat
);
```

**Catatan Penting:**
- `logo_url` harus berupa URL yang dapat diakses (tidak bisa menggunakan base64 atau file lokal)
- Jika logo tidak dapat di-load, surat tetap akan di-generate tanpa logo
- Semua field optional, jika tidak dikirim akan menggunakan default

---

## ⚠️ Error Handling

### Error Response Format

```json
{
  "success": false,
  "message": "Pesan error umum",
  "error": "Detail error atau pesan spesifik",
  "type": "jenis_error"
}
```

### HTTP Status Codes

| Code | Deskripsi |
|------|-----------|
| `200` | Success (untuk generate HTML) |
| `400` | Bad Request (format tanggal tidak valid) |
| `401` | Unauthorized (token tidak valid atau tidak ada) |
| `403` | Forbidden (tidak punya permission/role) |
| `500` | Internal Server Error |

### Contoh Error Responses

#### 400 Bad Request - Format Tanggal Tidak Valid

```json
{
  "success": false,
  "message": "Format tanggal tidak valid",
  "error": "Format tanggal_awal harus YYYY-MM-DD",
  "type": "validation"
}
```

#### 401 Unauthorized

```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Token tidak valid atau sudah kadaluarsa",
  "type": "authentication"
}
```

#### 403 Forbidden

```json
{
  "success": false,
  "message": "Akses ditolak",
  "error": "Anda tidak memiliki permission untuk mengakses endpoint ini",
  "type": "authorization"
}
```

#### 500 Internal Server Error

```json
{
  "success": false,
  "message": "Database error",
  "error": "Gagal mengambil data pendaftaran",
  "type": "database"
}
```

### Error Handling Example

```javascript
async function generatePengumumanWithErrorHandling(tanggalAwal, tanggalAkhir) {
  try {
    const response = await fetch(
      `/simnikah/staff/pengumuman-nikah/generate?tanggal_awal=${tanggalAwal}&tanggal_akhir=${tanggalAkhir}`,
      {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }
    );
    
    if (!response.ok) {
      const errorData = await response.json();
      
      switch (response.status) {
        case 400:
          console.error('Validation error:', errorData.error);
          alert(`Format tanggal tidak valid: ${errorData.error}`);
          break;
        case 401:
          console.error('Authentication error');
          // Redirect to login
          window.location.href = '/login';
          break;
        case 403:
          console.error('Authorization error');
          alert('Anda tidak memiliki akses untuk fitur ini');
          break;
        case 500:
          console.error('Server error:', errorData.error);
          alert('Terjadi kesalahan pada server. Silakan coba lagi nanti.');
          break;
        default:
          console.error('Unknown error:', errorData);
          alert('Terjadi kesalahan. Silakan coba lagi.');
      }
      
      throw new Error(errorData.message);
    }
    
    const html = await response.text();
    return html;
    
  } catch (error) {
    console.error('Error generating pengumuman:', error);
    throw error;
  }
}
```

---

## 💻 Contoh Implementasi Lengkap

### React Component

```jsx
import React, { useState } from 'react';
import axios from 'axios';

function PengumumanNikahGenerator() {
  const [tanggalAwal, setTanggalAwal] = useState('');
  const [tanggalAkhir, setTanggalAkhir] = useState('');
  const [html, setHtml] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showCustomKop, setShowCustomKop] = useState(false);
  
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
  
  const generatePengumuman = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('token');
      const url = '/simnikah/staff/pengumuman-nikah/generate';
      
      const config = {
        method: showCustomKop ? 'post' : 'get',
        url,
        params: {
          tanggal_awal: tanggalAwal || undefined,
          tanggal_akhir: tanggalAkhir || undefined
        },
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        responseType: 'text'
      };
      
      if (showCustomKop) {
        config.data = kopSurat;
      }
      
      const response = await axios(config);
      setHtml(response.data);
      
    } catch (err) {
      if (err.response) {
        setError(err.response.data.error || err.response.data.message);
      } else {
        setError('Terjadi kesalahan. Silakan coba lagi.');
      }
    } finally {
      setLoading(false);
    }
  };
  
  const printPengumuman = () => {
    if (!html) return;
    
    const printWindow = window.open('', '_blank');
    printWindow.document.write(html);
    printWindow.document.close();
    printWindow.print();
  };
  
  const downloadHTML = () => {
    if (!html) return;
    
    const blob = new Blob([html], { type: 'text/html' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `pengumuman-nikah-${tanggalAwal || 'current'}-${tanggalAkhir || 'current'}.html`;
    a.click();
    URL.revokeObjectURL(url);
  };
  
  return (
    <div className="pengumuman-generator">
      <h2>Generate Surat Pengumuman Nikah</h2>
      
      <div className="form-group">
        <label>Tanggal Awal (YYYY-MM-DD):</label>
        <input
          type="date"
          value={tanggalAwal}
          onChange={(e) => setTanggalAwal(e.target.value)}
        />
      </div>
      
      <div className="form-group">
        <label>Tanggal Akhir (YYYY-MM-DD):</label>
        <input
          type="date"
          value={tanggalAkhir}
          onChange={(e) => setTanggalAkhir(e.target.value)}
        />
      </div>
      
      <div className="form-group">
        <label>
          <input
            type="checkbox"
            checked={showCustomKop}
            onChange={(e) => setShowCustomKop(e.target.checked)}
          />
          Gunakan Kop Surat Custom
        </label>
      </div>
      
      {showCustomKop && (
        <div className="custom-kop-form">
          <h3>Kop Surat Custom</h3>
          <input
            type="text"
            placeholder="Nama KUA"
            value={kopSurat.nama_kua}
            onChange={(e) => setKopSurat({...kopSurat, nama_kua: e.target.value})}
          />
          <input
            type="text"
            placeholder="Alamat KUA"
            value={kopSurat.alamat_kua}
            onChange={(e) => setKopSurat({...kopSurat, alamat_kua: e.target.value})}
          />
          <input
            type="text"
            placeholder="Kota"
            value={kopSurat.kota}
            onChange={(e) => setKopSurat({...kopSurat, kota: e.target.value})}
          />
          <input
            type="text"
            placeholder="Provinsi"
            value={kopSurat.provinsi}
            onChange={(e) => setKopSurat({...kopSurat, provinsi: e.target.value})}
          />
          <input
            type="text"
            placeholder="Kode Pos"
            value={kopSurat.kode_pos}
            onChange={(e) => setKopSurat({...kopSurat, kode_pos: e.target.value})}
          />
          <input
            type="text"
            placeholder="Telepon"
            value={kopSurat.telepon}
            onChange={(e) => setKopSurat({...kopSurat, telepon: e.target.value})}
          />
          <input
            type="email"
            placeholder="Email"
            value={kopSurat.email}
            onChange={(e) => setKopSurat({...kopSurat, email: e.target.value})}
          />
          <input
            type="url"
            placeholder="Website (optional)"
            value={kopSurat.website}
            onChange={(e) => setKopSurat({...kopSurat, website: e.target.value})}
          />
          <input
            type="url"
            placeholder="Logo URL (optional)"
            value={kopSurat.logo_url}
            onChange={(e) => setKopSurat({...kopSurat, logo_url: e.target.value})}
          />
        </div>
      )}
      
      <button onClick={generatePengumuman} disabled={loading}>
        {loading ? 'Generating...' : 'Generate Pengumuman'}
      </button>
      
      {error && (
        <div className="error-message">
          <strong>Error:</strong> {error}
        </div>
      )}
      
      {html && (
        <div className="result-actions">
          <button onClick={printPengumuman}>Print</button>
          <button onClick={downloadHTML}>Download HTML</button>
        </div>
      )}
      
      {html && (
        <div className="preview">
          <h3>Preview:</h3>
          <iframe
            srcDoc={html}
            style={{
              width: '100%',
              height: '800px',
              border: '1px solid #ccc',
              marginTop: '20px'
            }}
            title="Pengumuman Nikah Preview"
          />
        </div>
      )}
    </div>
  );
}

export default PengumumanNikahGenerator;
```

### Vue.js Component

```vue
<template>
  <div class="pengumuman-generator">
    <h2>Generate Surat Pengumuman Nikah</h2>
    
    <div class="form-group">
      <label>Tanggal Awal:</label>
      <input type="date" v-model="tanggalAwal" />
    </div>
    
    <div class="form-group">
      <label>Tanggal Akhir:</label>
      <input type="date" v-model="tanggalAkhir" />
    </div>
    
    <button @click="generatePengumuman" :disabled="loading">
      {{ loading ? 'Generating...' : 'Generate' }}
    </button>
    
    <div v-if="error" class="error">{{ error }}</div>
    
    <div v-if="html" class="preview">
      <button @click="printPengumuman">Print</button>
      <iframe :srcdoc="html" style="width: 100%; height: 800px;"></iframe>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      tanggalAwal: '',
      tanggalAkhir: '',
      html: '',
      loading: false,
      error: null
    };
  },
  methods: {
    async generatePengumuman() {
      this.loading = true;
      this.error = null;
      
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get(
          '/simnikah/staff/pengumuman-nikah/generate',
          {
            params: {
              tanggal_awal: this.tanggalAwal || undefined,
              tanggal_akhir: this.tanggalAkhir || undefined
            },
            headers: {
              'Authorization': `Bearer ${token}`
            },
            responseType: 'text'
          }
        );
        
        this.html = response.data;
      } catch (err) {
        this.error = err.response?.data?.error || 'Terjadi kesalahan';
      } finally {
        this.loading = false;
      }
    },
    
    printPengumuman() {
      if (!this.html) return;
      
      const printWindow = window.open('', '_blank');
      printWindow.document.write(this.html);
      printWindow.document.close();
      printWindow.print();
    }
  }
};
</script>
```

### Node.js/Express Backend

```javascript
const express = require('express');
const axios = require('axios');
const router = express.Router();

// Proxy endpoint untuk generate pengumuman
router.get('/generate-pengumuman', async (req, res) => {
  try {
    const { tanggal_awal, tanggal_akhir } = req.query;
    const token = req.headers.authorization;
    
    const response = await axios.get(
      'https://your-api-domain.com/simnikah/staff/pengumuman-nikah/generate',
      {
        params: { tanggal_awal, tanggal_akhir },
        headers: { Authorization: token },
        responseType: 'text'
      }
    );
    
    res.setHeader('Content-Type', 'text/html; charset=utf-8');
    res.send(response.data);
    
  } catch (error) {
    res.status(error.response?.status || 500).json({
      success: false,
      error: error.response?.data?.error || error.message
    });
  }
});

module.exports = router;
```

---

## 📌 Best Practices

1. **Caching:** Cache HTML response jika data tidak berubah untuk mengurangi load server
2. **Error Handling:** Selalu handle error dengan baik dan berikan feedback ke user
3. **Loading State:** Tampilkan loading indicator saat generate
4. **Validation:** Validasi format tanggal sebelum mengirim request
5. **Security:** Jangan expose token di client-side code yang bisa diakses public
6. **Logo URL:** Pastikan logo URL accessible dan menggunakan HTTPS
7. **Print Optimization:** Gunakan CSS print media queries untuk optimasi print

---

## 🔗 Related Documentation

- [API Documentation Lengkap](./API_DOCUMENTATION_LENGKAP.md)
- [Alur Pendaftaran dan Jam](./ALUR_PENDAFTARAN_DAN_JAM.md)
- [User Guide Sederhana](./USER_GUIDE_SEDERHANA.md)

---

## 📞 Support

Jika ada pertanyaan atau masalah, silakan hubungi:
- Email: support@simnikah.com
- Documentation: https://docs.simnikah.com

---

**Last Updated:** Desember 2024  
**Version:** 1.0.0



