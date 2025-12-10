# 🔄 Frontend Migration Guide - Update Format Pengumuman Nikah

**Tanggal Update:** Desember 2024  
**Versi:** 2.0.0  
**Status:** ✅ No Breaking Changes

---

## 📋 Overview

Update format pengumuman nikah dari 9 kolom ke 15 kolom **TIDAK memerlukan perubahan setup atau kode frontend**. API endpoint, request format, dan response type tetap sama. Hanya **format HTML output** yang berubah.

---

## ✅ Tidak Ada Perubahan

### 1. API Endpoint
- ✅ URL endpoint **tetap sama**
- ✅ Method (GET/POST) **tetap sama**
- ✅ Query parameters **tetap sama**
- ✅ Request body format **tetap sama**

### 2. Request Format
- ✅ Query parameters: `tanggal_awal`, `tanggal_akhir` (format: YYYY-MM-DD)
- ✅ Request body: JSON dengan field kop surat (opsional)
- ✅ Headers: `Authorization: Bearer <token>`, `Content-Type: application/json`

### 3. Response Format
- ✅ Response type: **`text/html; charset=utf-8`** (tetap sama)
- ✅ Response body: **HTML string** (tetap sama)
- ✅ Error responses: **Format JSON tetap sama**

### 4. Frontend Code
- ✅ `responseType: 'text'` **tetap sama**
- ✅ Cara handle HTML response **tetap sama**
- ✅ Print/PDF generation **tetap sama**

---

## 🔄 Yang Berubah (Hanya Output HTML)

### Format HTML Output

**Sebelum (v1.0.0):**
- Judul: "PENGUMUMAN PERNIKAHAN"
- Tabel: 9 kolom
- Layout: A4 Portrait
- Kolom: No, Nomor Pendaftaran, Tanggal Nikah, Waktu, Tempat, Alamat Akad, Calon Suami, Calon Istri, Wali Nikah

**Sesudah (v2.0.0):**
- Judul: "JADUAL NIKAH [BULAN] [TAHUN]"
- Tabel: **15 kolom** dengan grouping
- Layout: **A4 Landscape**
- Kolom: NO URUT, PRIA/BIN, USIA, PENDK, WANITA/BINTI, USIA, PENDK, HARI, TGL, JAM, TEMPAT, WALINIKAH, PENGHULU, KELURAHAN, KET

### Dampak untuk Frontend

**✅ Tidak Ada Dampak:**
- Kode API call tidak perlu diubah
- Error handling tidak perlu diubah
- Print/PDF generation tetap berfungsi
- HTML masih bisa di-parse/ditampilkan dengan cara yang sama

**⚠️ Perhatian (Opsional):**
- Jika frontend melakukan parsing HTML untuk extract data spesifik, struktur HTML berubah
- Jika frontend memiliki CSS custom untuk styling tabel, mungkin perlu disesuaikan
- Layout landscape mungkin perlu setting printer/browser yang berbeda

---

## 📝 Contoh Kode (Tidak Berubah)

### Request (Tetap Sama)

```javascript
// GET Request
const response = await axios.get(
  `${API_BASE_URL}/staff/pengumuman-nikah/generate`,
  {
    params: {
      tanggal_awal: '2024-12-16',
      tanggal_akhir: '2024-12-22'
    },
    headers: {
      'Authorization': `Bearer ${token}`
    },
    responseType: 'text' // Tetap sama
  }
);

// POST Request dengan custom kop surat
const response = await axios.post(
  `${API_BASE_URL}/staff/pengumuman-nikah/generate`,
  {
    tanggal_awal: '2024-12-16',
    tanggal_akhir: '2024-12-22',
    nama_kua: 'KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA',
    alamat_kua: 'Jl. Contoh No. 123',
    kota: 'Kota Banjarmasin',
    provinsi: 'Kalimantan Selatan',
    kode_pos: '70123',
    telepon: '0511-1234567',
    email: 'kua@example.com',
    logo_url: 'https://example.com/logo.png'
  },
  {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    responseType: 'text' // Tetap sama
  }
);
```

### Handle Response (Tetap Sama)

```javascript
// Print HTML
const handlePrint = async () => {
  try {
    const html = await generatePengumumanHTML(tanggalAwal, tanggalAkhir);
    
    const printWindow = window.open('', '_blank');
    printWindow.document.write(html);
    printWindow.document.close();
    
    setTimeout(() => {
      printWindow.print();
    }, 500);
  } catch (error) {
    console.error('Error:', error);
  }
};

// Export PDF
const handleExportPDF = async () => {
  try {
    const html = await generatePengumumanHTML(tanggalAwal, tanggalAkhir);
    
    const element = document.createElement('div');
    element.innerHTML = html;
    
    const opt = {
      margin: 1,
      filename: 'pengumuman-nikah.pdf',
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { scale: 2, useCORS: true },
      jsPDF: { unit: 'in', format: 'a4', orientation: 'landscape' } // ⚠️ Perhatikan: landscape
    };
    
    await html2pdf().set(opt).from(element).save();
  } catch (error) {
    console.error('Error:', error);
  }
};

// Preview HTML
const handlePreview = async () => {
  try {
    const html = await generatePengumumanHTML(tanggalAwal, tanggalAkhir);
    
    const iframe = document.createElement('iframe');
    iframe.style.width = '100%';
    iframe.style.height = '800px';
    iframe.style.border = 'none';
    
    document.body.appendChild(iframe);
    iframe.contentDocument.write(html);
    iframe.contentDocument.close();
  } catch (error) {
    console.error('Error:', error);
  }
};
```

---

## ⚠️ Perhatian (Opsional Update)

### 1. PDF Export - Orientation

Jika menggunakan `html2pdf.js` atau library serupa, pastikan set orientation ke **landscape**:

```javascript
// Sebelum (Portrait)
jsPDF: { unit: 'in', format: 'a4', orientation: 'portrait' }

// Sesudah (Landscape) - OPSIONAL, bisa tetap portrait jika diinginkan
jsPDF: { unit: 'in', format: 'a4', orientation: 'landscape' }
```

**Catatan:** HTML sudah include CSS `@media print` dengan `@page { size: A4 landscape; }`, jadi browser akan otomatis menggunakan landscape saat print. Tapi untuk PDF export, mungkin perlu set manual.

### 2. Print Preview

Browser print preview akan otomatis menggunakan landscape karena CSS sudah include. Tidak perlu perubahan kode.

### 3. HTML Parsing (Jika Ada)

Jika frontend melakukan parsing HTML untuk extract data spesifik (tidak disarankan, lebih baik gunakan API `/list`), struktur HTML berubah:

**Sebelum:**
- Tabel dengan 9 kolom
- Selector: `table tr td:nth-child(1)` = No, `td:nth-child(2)` = Nomor Pendaftaran, dll.

**Sesudah:**
- Tabel dengan 15 kolom
- Selector: `table tr td:nth-child(1)` = NO URUT, `td:nth-child(2)` = PRIA/BIN, dll.

**Rekomendasi:** Gunakan endpoint `/list` untuk mendapatkan data dalam format JSON, bukan parsing HTML.

---

## ✅ Checklist Migration

- [x] **Tidak ada perubahan kode API call** - Endpoint, method, params tetap sama
- [x] **Tidak ada perubahan response handling** - Masih HTML string
- [x] **Tidak ada perubahan error handling** - Format error tetap sama
- [ ] **Opsional: Update PDF export orientation** - Set ke landscape jika diinginkan
- [ ] **Opsional: Update HTML parsing** - Jika ada custom parsing logic
- [ ] **Opsional: Test print/PDF** - Pastikan output sesuai ekspektasi

---

## 🧪 Testing

### Test Checklist

1. **API Call**
   - [ ] GET request dengan query params berhasil
   - [ ] POST request dengan custom kop surat berhasil
   - [ ] Error handling masih berfungsi

2. **HTML Display**
   - [ ] HTML ditampilkan dengan benar di iframe/preview
   - [ ] Tabel 15 kolom terlihat lengkap
   - [ ] Kop surat dengan logo (jika ada) terlihat benar

3. **Print**
   - [ ] Print preview menggunakan landscape
   - [ ] Tabel tidak terpotong
   - [ ] Font dan styling sesuai

4. **PDF Export**
   - [ ] PDF ter-generate dengan benar
   - [ ] Layout landscape (jika set manual)
   - [ ] Tabel lengkap dan tidak terpotong

---

## 📚 Dokumentasi Terkait

- [Fitur Terbaru](./FITUR_TERBARU.md) - Dokumentasi lengkap fitur generate pengumuman nikah
- [Update Format Pengumuman Nikah](./UPDATE_FORMAT_PENGUMUMAN_NIKAH.md) - Detail perubahan format
- [API Documentation Lengkap](./api/API_DOCUMENTATION_LENGKAP.md) - Dokumentasi endpoint lengkap
- [Parsing API Pengumuman Nikah](./api/PARSING_API_PENGUMUMAN_NIKAH.md) - Panduan parsing HTML response

---

## ❓ FAQ

### Q: Apakah perlu update kode frontend?
**A:** Tidak perlu. API endpoint, request format, dan response type tetap sama. Hanya format HTML output yang berubah.

### Q: Apakah print masih berfungsi?
**A:** Ya, masih berfungsi. Browser akan otomatis menggunakan landscape karena CSS sudah include `@page { size: A4 landscape; }`.

### Q: Apakah PDF export perlu diubah?
**A:** Opsional. Jika menggunakan library seperti `html2pdf.js`, bisa set orientation ke landscape untuk hasil yang lebih baik. Tapi tidak wajib.

### Q: Apakah ada breaking changes?
**A:** Tidak ada breaking changes. Semua API call tetap berfungsi dengan kode yang sama.

### Q: Bagaimana jika frontend melakukan parsing HTML?
**A:** Struktur HTML berubah, jadi parsing logic perlu diupdate. Tapi lebih baik gunakan endpoint `/list` untuk mendapatkan data dalam format JSON.

---

**Last Updated:** Desember 2024  
**Version:** 2.0.0

