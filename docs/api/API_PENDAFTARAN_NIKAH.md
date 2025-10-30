# 📋 API Pendaftaran Nikah - SimNikah

**Version:** 2.0 (Updated with Performance Optimizations & Enhanced Validations)  
**Base URL:** `https://your-api-domain.com/api`  
**Last Updated:** October 30, 2025

---

## 🚀 Performa & Optimasi

API Pendaftaran Nikah telah dioptimasi dengan:
- ✅ **Batch Database Insert** - Mengurangi query dari 8 menjadi 5
- ✅ **Async Geocoding** - Koordinat lokasi diproses di background
- ✅ **Async Notifications** - Notifikasi dikirim setelah response
- ✅ **Response Time** - Rata-rata 800-1200ms (sebelumnya 2500-4000ms)

---

## 📌 Endpoint: Create Marriage Registration

### **POST** `/catin/pendaftaran`

Membuat pendaftaran nikah lengkap dengan data calon pengantin, orang tua, dan wali nikah.

### Authentication
**Required:** Yes  
**Type:** Bearer Token (JWT)  
**Header:** `Authorization: Bearer <your_jwt_token>`

---

## 📥 Request Body

### Content-Type
```
Content-Type: application/json
```

### Complete Request Structure

```json
{
  "scheduleAndLocation": {
    "weddingLocation": "Di KUA",
    "weddingAddress": "Jl. Example No. 123, Banjarmasin",
    "weddingDate": "2025-12-25",
    "weddingTime": "09:00",
    "dispensationNumber": "DISP/123/2025"
  },
  "groom": {
    "groomFullName": "Ahmad Fadli",
    "groomNik": "6371012501950001",
    "groomCitizenship": "WNI",
    "groomPassportNumber": "",
    "groomPlaceOfBirth": "Banjarmasin",
    "groomDateOfBirth": "1995-01-25",
    "groomStatus": "Belum Kawin",
    "groomReligion": "Islam",
    "groomEducation": "S1",
    "groomOccupation": "Pegawai Swasta",
    "groomOccupationDescription": "Staff IT",
    "groomPhoneNumber": "081234567890",
    "groomEmail": "ahmad.fadli@example.com",
    "groomAddress": "Jl. Sutoyo S No. 123, Banjarmasin"
  },
  "bride": {
    "brideFullName": "Siti Nurhaliza",
    "brideNik": "6371015201980002",
    "brideCitizenship": "WNI",
    "bridePassportNumber": "",
    "bridePlaceOfBirth": "Banjarmasin",
    "brideDateOfBirth": "1998-12-15",
    "brideStatus": "Belum Kawin",
    "brideReligion": "Islam",
    "brideEducation": "D3",
    "brideOccupation": "Pegawai Swasta",
    "brideOccupationDescription": "Admin",
    "bridePhoneNumber": "081298765432",
    "brideEmail": "siti.nurhaliza@example.com",
    "brideAddress": "Jl. A. Yani No. 456, Banjarmasin"
  },
  "groomParents": {
    "father": {
      "fatherName": "Budi Santoso",
      "fatherNik": "6371011234567890",
      "fatherCitizenship": "WNI",
      "fatherPassportNumber": "",
      "fatherReligion": "Islam",
      "fatherPlaceOfBirth": "Banjarmasin",
      "fatherDateOfBirth": "1970-05-10",
      "fatherCountry": "Indonesia",
      "fatherOccupation": "PNS",
      "fatherOccupationDescription": "Guru",
      "fatherAddress": "Jl. Veteran No. 789, Banjarmasin",
      "groomFatherPresenceStatus": "Hidup"
    },
    "mother": {
      "motherName": "Dewi Lestari",
      "motherNik": "6371019876543210",
      "motherCitizenship": "WNI",
      "motherPassportNumber": "",
      "motherReligion": "Islam",
      "motherPlaceOfBirth": "Banjarmasin",
      "motherDateOfBirth": "1972-08-20",
      "motherCountry": "Indonesia",
      "motherOccupation": "Ibu Rumah Tangga",
      "motherOccupationDescription": "",
      "motherAddress": "Jl. Veteran No. 789, Banjarmasin",
      "groomMotherPresenceStatus": "Hidup"
    }
  },
  "brideParents": {
    "father": {
      "fatherName": "Hasan Basri",
      "fatherNik": "6371011122334455",
      "fatherCitizenship": "WNI",
      "fatherPassportNumber": "",
      "fatherReligion": "Islam",
      "fatherPlaceOfBirth": "Banjarmasin",
      "fatherDateOfBirth": "1965-03-15",
      "fatherCountry": "Indonesia",
      "fatherOccupation": "Wiraswasta",
      "fatherOccupationDescription": "Pedagang",
      "fatherAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
      "brideFatherPresenceStatus": "Hidup"
    },
    "mother": {
      "motherName": "Fatimah Zahra",
      "motherNik": "6371015566778899",
      "motherCitizenship": "WNI",
      "motherPassportNumber": "",
      "motherReligion": "Islam",
      "motherPlaceOfBirth": "Banjarmasin",
      "motherDateOfBirth": "1968-07-25",
      "motherCountry": "Indonesia",
      "motherOccupation": "Ibu Rumah Tangga",
      "motherOccupationDescription": "",
      "motherAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
      "brideMotherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianFullName": "Hasan Basri",
    "guardianNik": "6371011122334455",
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Hidup",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
    "guardianPhoneNumber": "081298761234"
  }
}
```

---

## 📋 Field Descriptions

### 1. Schedule and Location (`scheduleAndLocation`)

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `weddingLocation` | string | ✅ Yes | Lokasi akad nikah | `"Di KUA"` atau `"Di Luar KUA"` |
| `weddingAddress` | string | ⚠️ Conditional | Alamat lengkap (wajib jika Di Luar KUA) | `"Jl. Example No. 123"` |
| `weddingDate` | string | ✅ Yes | Tanggal akad (YYYY-MM-DD) | `"2025-12-25"` |
| `weddingTime` | string | ✅ Yes | Waktu akad (HH:MM format 24-jam) | `"09:00"` |
| `dispensationNumber` | string | ⚠️ Conditional | Nomor dispensasi (jika diperlukan) | `"DISP/123/2025"` |

**Validations:**
- ✅ `weddingDate` tidak boleh di masa lalu
- ✅ `weddingTime` harus format HH:MM (24-jam)
- ✅ Minimal 10 hari kerja sebelum tanggal nikah (kecuali ada dispensasi)
- ✅ Dispensasi wajib jika:
  - Calon suami < 19 tahun ATAU calon istri < 19 tahun
  - Kurang dari 10 hari kerja sebelum tanggal nikah

---

### 2. Groom Data (`groom`)

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `groomFullName` | string | ✅ Yes | Nama lengkap calon suami | `"Ahmad Fadli"` |
| `groomNik` | string | ✅ Yes | NIK calon suami (16 digit) | `"6371012501950001"` |
| `groomCitizenship` | string | ✅ Yes | Kewarganegaraan | `"WNI"` atau `"WNA"` |
| `groomPassportNumber` | string | ⚠️ Conditional | Nomor paspor (wajib jika WNA) | `"A12345678"` |
| `groomPlaceOfBirth` | string | ✅ Yes | Tempat lahir | `"Banjarmasin"` |
| `groomDateOfBirth` | string | ✅ Yes | Tanggal lahir (YYYY-MM-DD) | `"1995-01-25"` |
| `groomStatus` | string | ✅ Yes | Status perkawinan | `"Belum Kawin"`, `"Kawin"`, `"Cerai Hidup"`, `"Cerai Mati"` |
| `groomReligion` | string | ✅ Yes | Agama | `"Islam"`, `"Kristen"`, `"Katolik"`, `"Hindu"`, `"Buddha"`, `"Konghucu"` |
| `groomEducation` | string | ✅ Yes | Pendidikan terakhir | `"SD"`, `"SMP"`, `"SMA"`, `"D3"`, `"S1"`, `"S2"`, `"S3"` |
| `groomOccupation` | string | ✅ Yes | Pekerjaan | `"PNS"`, `"TNI/Polri"`, `"Pegawai Swasta"`, dll |
| `groomOccupationDescription` | string | ❌ No | Deskripsi pekerjaan | `"Staff IT"` |
| `groomPhoneNumber` | string | ✅ Yes | Nomor telepon | `"081234567890"` |
| `groomEmail` | string | ✅ Yes | Email | `"ahmad@example.com"` |
| `groomAddress` | string | ✅ Yes | Alamat lengkap | `"Jl. Sutoyo S No. 123"` |

---

### 3. Bride Data (`bride`)

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `brideFullName` | string | ✅ Yes | Nama lengkap calon istri | `"Siti Nurhaliza"` |
| `brideNik` | string | ✅ Yes | NIK calon istri (16 digit) | `"6371015201980002"` |
| `brideCitizenship` | string | ✅ Yes | Kewarganegaraan | `"WNI"` atau `"WNA"` |
| `bridePassportNumber` | string | ⚠️ Conditional | Nomor paspor (wajib jika WNA) | `"B98765432"` |
| `bridePlaceOfBirth` | string | ✅ Yes | Tempat lahir | `"Banjarmasin"` |
| `brideDateOfBirth` | string | ✅ Yes | Tanggal lahir (YYYY-MM-DD) | `"1998-12-15"` |
| `brideStatus` | string | ✅ Yes | Status perkawinan | `"Belum Kawin"`, `"Kawin"`, `"Cerai Hidup"`, `"Cerai Mati"` |
| `brideReligion` | string | ✅ Yes | Agama | `"Islam"`, `"Kristen"`, `"Katolik"`, dll |
| `brideEducation` | string | ✅ Yes | Pendidikan terakhir | `"SD"`, `"SMP"`, `"SMA"`, `"D3"`, `"S1"`, dll |
| `brideOccupation` | string | ✅ Yes | Pekerjaan | `"Pegawai Swasta"`, `"Ibu Rumah Tangga"`, dll |
| `brideOccupationDescription` | string | ❌ No | Deskripsi pekerjaan | `"Admin"` |
| `bridePhoneNumber` | string | ✅ Yes | Nomor telepon | `"081298765432"` |
| `brideEmail` | string | ✅ Yes | Email | `"siti@example.com"` |
| `brideAddress` | string | ✅ Yes | Alamat lengkap | `"Jl. A. Yani No. 456"` |

---

### 4. Parents Data (`groomParents` & `brideParents`)

**Structure sama untuk father dan mother:**

| Field | Type | Required | Description | Valid Values |
|-------|------|----------|-------------|--------------|
| `fatherName` / `motherName` | string | ✅ Yes | Nama lengkap orang tua | - |
| `fatherNik` / `motherNik` | string | ✅ Yes | NIK orang tua | 16 digit |
| `*Citizenship` | string | ✅ Yes | Kewarganegaraan | `"WNI"`, `"WNA"` |
| `*PassportNumber` | string | ⚠️ Conditional | Nomor paspor (wajib jika WNA) | - |
| `*Religion` | string | ✅ Yes | Agama | `"Islam"`, `"Kristen"`, dll |
| `*PlaceOfBirth` | string | ✅ Yes | Tempat lahir | - |
| `*DateOfBirth` | string | ✅ Yes | Tanggal lahir (YYYY-MM-DD) | - |
| `*Country` | string | ✅ Yes | Negara asal | `"Indonesia"` |
| `*Occupation` | string | ✅ Yes | Pekerjaan | - |
| `*OccupationDescription` | string | ❌ No | Deskripsi pekerjaan | - |
| `*Address` | string | ✅ Yes | Alamat lengkap | - |
| `*PresenceStatus` | string | ✅ Yes | Status keberadaan | `"Hidup"`, `"Meninggal"` |

**Note:** Data orang tua hanya disimpan ke database jika status keberadaan = `"Hidup"`

---

### 5. Guardian Data (`guardian`) 🔒 **Enhanced Validations**

| Field | Type | Required | Description | Valid Values |
|-------|------|----------|-------------|--------------|
| `guardianFullName` | string | ✅ Yes | Nama lengkap wali nikah | - |
| `guardianNik` | string | ✅ Yes | NIK wali nikah (16 digit) | - |
| `guardianRelationship` | string | ✅ Yes | Hubungan wali dengan calon istri | Lihat tabel di bawah |
| `guardianStatus` | string | ✅ Yes | Status keberadaan wali | `"Hidup"`, `"Meninggal"` |
| `guardianReligion` | string | ✅ Yes | Agama wali | `"Islam"`, dll |
| `guardianAddress` | string | ✅ Yes | Alamat lengkap wali | - |
| `guardianPhoneNumber` | string | ✅ Yes | Nomor telepon wali | - |

#### Valid Guardian Relationships (Urutan Wali Nasab):

| Value | Description | Prioritas |
|-------|-------------|-----------|
| `"Ayah Kandung"` | Ayah kandung (wajib jika ayah hidup) | 1 (Tertinggi) |
| `"Kakek"` | Kakek dari pihak ayah | 2 |
| `"Saudara Laki-Laki Kandung"` | Saudara laki-laki sekandung | 3 |
| `"Saudara Laki-Laki Seayah"` | Saudara laki-laki seayah | 4 |
| `"Keponakan Laki-Laki"` | Anak laki-laki dari saudara | 5 |
| `"Paman Kandung"` | Saudara laki-laki ayah kandung | 6 |
| `"Paman Seayah"` | Saudara laki-laki ayah seayah | 7 |
| `"Sepupu Laki-Laki"` | Anak laki-laki dari paman | 8 |
| `"Wali Hakim"` | Wali hakim dari KUA | 9 |
| `"Lainnya"` | Hubungan lainnya | 10 |

---

## 🔒 Guardian Validation Rules (Syariat Islam)

### Validasi 1: Wali Harus Hidup ✅
```
❌ Error jika: guardianStatus = "Meninggal"
✅ Valid jika: guardianStatus = "Hidup"
```

**Error Response:**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Wali nikah yang telah meninggal dunia tidak dapat menjadi wali. Silakan pilih wali lain yang masih hidup.",
  "field": "status_wali",
  "type": "syariat_validation"
}
```

---

### Validasi 2: Konsistensi Status Wali Ayah ✅
```
Jika guardianRelationship = "Ayah Kandung"
MAKA guardianStatus HARUS sama dengan brideFatherPresenceStatus
```

**Example:**
```json
// ❌ INVALID - Inconsistent status
{
  "brideParents": {
    "father": {
      "brideFatherPresenceStatus": "Hidup"  // Ayah hidup
    }
  },
  "guardian": {
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Meninggal"  // ❌ Bertentangan!
  }
}

// ✅ VALID
{
  "brideParents": {
    "father": {
      "brideFatherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Hidup"  // ✅ Konsisten
  }
}
```

---

### Validasi 3: NIK Wali = NIK Ayah (Jika Wali Ayah Kandung) ✅
```
Jika guardianRelationship = "Ayah Kandung" DAN brideFatherPresenceStatus = "Hidup"
MAKA guardianNik HARUS sama dengan fatherNik (dari brideParents)
```

**Example:**
```json
// ❌ INVALID - NIK tidak sama
{
  "brideParents": {
    "father": {
      "fatherNik": "6371011122334455",
      "brideFatherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianRelationship": "Ayah Kandung",
    "guardianNik": "1234567890123456",  // ❌ NIK berbeda!
    "guardianStatus": "Hidup"
  }
}

// ✅ VALID
{
  "brideParents": {
    "father": {
      "fatherNik": "6371011122334455",
      "brideFatherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianRelationship": "Ayah Kandung",
    "guardianNik": "6371011122334455",  // ✅ NIK sama
    "guardianStatus": "Hidup"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "NIK wali harus sama dengan NIK ayah kandung catin perempuan",
  "field": "nik_wali",
  "type": "syariat_validation",
  "details": {
    "nik_wali_yang_diinput": "1234567890123456",
    "nik_ayah_catin_perempuan": "6371011122334455"
  }
}
```

---

### Validasi 4: Wali Tidak Boleh Sama dengan Calon Pengantin ✅
```
guardianNik ≠ groomNik
guardianNik ≠ brideNik
```

**Error Response:**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "NIK wali tidak boleh sama dengan NIK calon suami",
  "field": "nik_wali",
  "type": "syariat_validation"
}
```

---

### Validasi 5: Urutan Wali Nasab Berdasarkan Status Ayah ✅

#### Scenario A: Ayah Kandung Masih Hidup
```
Jika brideFatherPresenceStatus = "Hidup"
MAKA guardianRelationship HARUS = "Ayah Kandung"
```

**Example:**
```json
// ❌ INVALID
{
  "brideParents": {
    "father": {
      "brideFatherPresenceStatus": "Hidup"  // Ayah masih hidup
    }
  },
  "guardian": {
    "guardianRelationship": "Kakek"  // ❌ Tidak boleh! Harus ayah
  }
}

// ✅ VALID
{
  "brideParents": {
    "father": {
      "brideFatherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianRelationship": "Ayah Kandung"  // ✅ Correct
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Jika ayah kandung masih hidup, maka wali nikah harus Ayah Kandung sesuai syariat Islam",
  "field": "hubungan_wali",
  "type": "syariat_validation",
  "details": {
    "status_ayah_catin_perempuan": "Hidup",
    "hubungan_wali_yang_dipilih": "Kakek",
    "urutan_wali_nasab": [
      "Ayah Kandung (prioritas 1)",
      "Kakek (prioritas 2)",
      "Saudara Laki-Laki Kandung (prioritas 3)",
      "..."
    ]
  }
}
```

---

#### Scenario B: Ayah Kandung Meninggal
```
Jika brideFatherPresenceStatus = "Meninggal"
MAKA guardianRelationship TIDAK BOLEH = "Ayah Kandung"
DAN harus pilih dari urutan wali nasab berikutnya
```

**Example:**
```json
// ❌ INVALID
{
  "brideParents": {
    "father": {
      "brideFatherPresenceStatus": "Meninggal"  // Ayah meninggal
    }
  },
  "guardian": {
    "guardianRelationship": "Ayah Kandung"  // ❌ Tidak mungkin!
  }
}

// ✅ VALID
{
  "brideParents": {
    "father": {
      "brideFatherPresenceStatus": "Meninggal"
    }
  },
  "guardian": {
    "guardianRelationship": "Kakek"  // ✅ Urutan wali nasab berikutnya
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Jika ayah kandung meninggal/tidak ada, wali nikah berpindah ke nasab berikutnya: Kakek, Saudara Laki-Laki Kandung, Paman, atau Wali Hakim",
  "field": "hubungan_wali",
  "type": "syariat_validation"
}
```

---

## 📤 Success Response

### Status Code: `201 Created`

```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat",
  "data": {
    "nomor_pendaftaran": "NIK2025103012345",
    "pendaftaran_id": 42,
    "calon_suami": {
      "id": 123,
      "user_id": "a1b2c3d4e5f6g7h8i9j0",
      "nik": "6371012501950001",
      "nama_lengkap": "Ahmad Fadli",
      "tempat_lahir": "Banjarmasin",
      "tanggal_lahir": "1995-01-25T00:00:00Z",
      "jenis_kelamin": "L",
      "alamat": "Jl. Sutoyo S No. 123, Banjarmasin",
      "agama": "Islam",
      "status_perkawinan": "Belum Kawin",
      "pekerjaan": "Pegawai Swasta",
      "deskripsi_pekerjaan": "Staff IT",
      "pendidikan_terakhir": "S1",
      "nomor_telepon": "081234567890",
      "email": "ahmad.fadli@example.com",
      "warga_negara": "WNI",
      "nomor_paspor": "",
      "dibuat_pada": "2025-10-30T10:15:30Z",
      "diperbarui_pada": "2025-10-30T10:15:30Z"
    },
    "calon_istri": {
      "id": 124,
      "user_id": "k1l2m3n4o5p6q7r8s9t0",
      "nik": "6371015201980002",
      "nama_lengkap": "Siti Nurhaliza",
      "tempat_lahir": "Banjarmasin",
      "tanggal_lahir": "1998-12-15T00:00:00Z",
      "jenis_kelamin": "P",
      "alamat": "Jl. A. Yani No. 456, Banjarmasin",
      "agama": "Islam",
      "status_perkawinan": "Belum Kawin",
      "pekerjaan": "Pegawai Swasta",
      "deskripsi_pekerjaan": "Admin",
      "pendidikan_terakhir": "D3",
      "nomor_telepon": "081298765432",
      "email": "siti.nurhaliza@example.com",
      "warga_negara": "WNI",
      "nomor_paspor": "",
      "dibuat_pada": "2025-10-30T10:15:30Z",
      "diperbarui_pada": "2025-10-30T10:15:30Z"
    },
    "pendaftaran": {
      "id": 42,
      "nomor_pendaftaran": "NIK2025103012345",
      "pendaftar_id": "user_12345",
      "calon_suami_id": "123",
      "calon_istri_id": "124",
      "tanggal_pendaftaran": "2025-10-30T10:15:30Z",
      "tanggal_nikah": "2025-12-25T00:00:00Z",
      "waktu_nikah": "09:00",
      "tempat_nikah": "Di KUA",
      "alamat_akad": "KUA Kecamatan Banjarmasin Utara, Kelurahan Pangeran, Kecamatan Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan",
      "latitude": -3.3148,
      "longitude": 114.5925,
      "nomor_dispensasi": "DISP/123/2025",
      "status_pendaftaran": "Menunggu Verifikasi",
      "dibuat_pada": "2025-10-30T10:15:30Z",
      "diperbarui_pada": "2025-10-30T10:15:30Z"
    },
    "wali_nikah": {
      "id": 67,
      "id_pendaftaran": 42,
      "nik": "6371011122334455",
      "nama_lengkap": "Hasan Basri",
      "hubungan_wali": "Ayah Kandung",
      "alamat": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
      "nomor_telepon": "081298761234",
      "email": "",
      "warga_negara": "WNI",
      "agama": "Islam",
      "tempat_lahir": "",
      "tanggal_lahir": null,
      "negara_asal": "Indonesia",
      "pekerjaan": "",
      "nomor_paspor": "",
      "pekerjaan_lain": "",
      "status_keberadaan": "Hidup",
      "status_kehadiran": "Belum Diketahui",
      "dibuat_pada": "2025-10-30T10:15:30Z",
      "diperbarui_pada": "2025-10-30T10:15:30Z"
    }
  }
}
```

### 📝 Response Notes:
- ✅ Response dikirim **immediately** (800-1200ms)
- 🔄 Geocoding (jika Di Luar KUA) diproses **asynchronous** di background
- 🔄 Notifikasi dikirim **asynchronous** setelah response
- 📊 `latitude` dan `longitude` akan `null` untuk "Di Luar KUA" saat pertama kali, kemudian diupdate otomatis dalam beberapa detik

---

## ❌ Error Responses

### 1. Validation Error - Format Tanggal Salah
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format tanggal nikah tidak valid (YYYY-MM-DD)",
  "field": "tanggal_nikah",
  "type": "format"
}
```

---

### 2. Validation Error - Tanggal Nikah Di Masa Lalu
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Tanggal nikah tidak boleh di masa lalu",
  "field": "tanggal_nikah",
  "type": "validation"
}
```

---

### 3. Validation Error - Format Waktu Salah
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format waktu nikah tidak valid (HH:MM dalam format 24-jam, contoh: 09:00)",
  "field": "waktu_nikah",
  "type": "format"
}
```

---

### 4. Validation Error - Dispensasi Diperlukan
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Dispensasi diperlukan",
  "error": "Kurang dari 10 hari kerja menuju tanggal nikah. Dispensasi diperlukan.",
  "field": "dispensation_number",
  "type": "dispensation_required",
  "details": {
    "working_days_remaining": 7,
    "minimum_required": 10,
    "groom_age": 20,
    "bride_age": 17,
    "reason": "Kurang dari 10 hari kerja"
  }
}
```

---

### 5. Validation Error - Umur Minimum
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Dispensasi diperlukan",
  "error": "Calon pengantin belum memenuhi syarat umur minimum (19 tahun). Diperlukan surat dispensasi dari pengadilan.",
  "field": "dispensation_number",
  "type": "dispensation_required",
  "details": {
    "working_days_remaining": 15,
    "minimum_required": 10,
    "groom_age": 20,
    "bride_age": 17,
    "reason": "Umur calon istri belum 19 tahun"
  }
}
```

---

### 6. Validation Error - Wali Meninggal
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Wali nikah yang telah meninggal dunia tidak dapat menjadi wali. Silakan pilih wali lain yang masih hidup.",
  "field": "status_wali",
  "type": "syariat_validation"
}
```

---

### 7. Validation Error - Status Wali Tidak Konsisten
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Jika wali adalah Ayah Kandung, maka status wali harus sama dengan status ayah catin perempuan (Hidup)",
  "field": "status_wali",
  "type": "syariat_validation"
}
```

---

### 8. Validation Error - NIK Wali Berbeda
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "NIK wali harus sama dengan NIK ayah kandung catin perempuan",
  "field": "nik_wali",
  "type": "syariat_validation",
  "details": {
    "nik_wali_yang_diinput": "1234567890123456",
    "nik_ayah_catin_perempuan": "6371011122334455"
  }
}
```

---

### 9. Validation Error - Wali = Calon Pengantin
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "NIK wali tidak boleh sama dengan NIK calon suami",
  "field": "nik_wali",
  "type": "syariat_validation"
}
```

---

### 10. Validation Error - Urutan Wali Nasab Salah
**Status Code:** `400 Bad Request`
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Jika ayah kandung masih hidup, maka wali nikah harus Ayah Kandung sesuai syariat Islam",
  "field": "hubungan_wali",
  "type": "syariat_validation",
  "details": {
    "status_ayah_catin_perempuan": "Hidup",
    "hubungan_wali_yang_dipilih": "Kakek",
    "urutan_wali_nasab": [
      "Ayah Kandung (prioritas 1)",
      "Kakek (prioritas 2)",
      "Saudara Laki-Laki Kandung (prioritas 3)",
      "Saudara Laki-Laki Seayah (prioritas 4)",
      "Keponakan Laki-Laki (prioritas 5)",
      "Paman Kandung (prioritas 6)",
      "Paman Seayah (prioritas 7)",
      "Sepupu Laki-Laki (prioritas 8)",
      "Wali Hakim (prioritas 9)",
      "Lainnya (prioritas 10)"
    ]
  }
}
```

---

### 11. Authentication Error
**Status Code:** `401 Unauthorized`
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "User ID tidak ditemukan",
  "type": "authentication"
}
```

---

### 12. Database Error
**Status Code:** `500 Internal Server Error`
```json
{
  "success": false,
  "message": "Database error",
  "error": "Gagal membuat pendaftaran nikah",
  "type": "database"
}
```

---

## 🧪 Test Cases

### Test Case 1: Valid Registration (Normal Case)
```bash
curl -X POST https://your-api-domain.com/api/catin/pendaftaran \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scheduleAndLocation": {
      "weddingLocation": "Di KUA",
      "weddingAddress": "",
      "weddingDate": "2025-12-25",
      "weddingTime": "09:00",
      "dispensationNumber": ""
    },
    "groom": {
      "groomFullName": "Ahmad Fadli",
      "groomNik": "6371012501950001",
      "groomCitizenship": "WNI",
      "groomPassportNumber": "",
      "groomPlaceOfBirth": "Banjarmasin",
      "groomDateOfBirth": "1995-01-25",
      "groomStatus": "Belum Kawin",
      "groomReligion": "Islam",
      "groomEducation": "S1",
      "groomOccupation": "Pegawai Swasta",
      "groomOccupationDescription": "Staff IT",
      "groomPhoneNumber": "081234567890",
      "groomEmail": "ahmad.fadli@example.com",
      "groomAddress": "Jl. Sutoyo S No. 123, Banjarmasin"
    },
    "bride": {
      "brideFullName": "Siti Nurhaliza",
      "brideNik": "6371015201980002",
      "brideCitizenship": "WNI",
      "bridePassportNumber": "",
      "bridePlaceOfBirth": "Banjarmasin",
      "brideDateOfBirth": "1998-12-15",
      "brideStatus": "Belum Kawin",
      "brideReligion": "Islam",
      "brideEducation": "D3",
      "brideOccupation": "Pegawai Swasta",
      "brideOccupationDescription": "Admin",
      "bridePhoneNumber": "081298765432",
      "brideEmail": "siti.nurhaliza@example.com",
      "brideAddress": "Jl. A. Yani No. 456, Banjarmasin"
    },
    "groomParents": {
      "father": {
        "fatherName": "Budi Santoso",
        "fatherNik": "6371011234567890",
        "fatherCitizenship": "WNI",
        "fatherPassportNumber": "",
        "fatherReligion": "Islam",
        "fatherPlaceOfBirth": "Banjarmasin",
        "fatherDateOfBirth": "1970-05-10",
        "fatherCountry": "Indonesia",
        "fatherOccupation": "PNS",
        "fatherOccupationDescription": "Guru",
        "fatherAddress": "Jl. Veteran No. 789, Banjarmasin",
        "groomFatherPresenceStatus": "Hidup"
      },
      "mother": {
        "motherName": "Dewi Lestari",
        "motherNik": "6371019876543210",
        "motherCitizenship": "WNI",
        "motherPassportNumber": "",
        "motherReligion": "Islam",
        "motherPlaceOfBirth": "Banjarmasin",
        "motherDateOfBirth": "1972-08-20",
        "motherCountry": "Indonesia",
        "motherOccupation": "Ibu Rumah Tangga",
        "motherOccupationDescription": "",
        "motherAddress": "Jl. Veteran No. 789, Banjarmasin",
        "groomMotherPresenceStatus": "Hidup"
      }
    },
    "brideParents": {
      "father": {
        "fatherName": "Hasan Basri",
        "fatherNik": "6371011122334455",
        "fatherCitizenship": "WNI",
        "fatherPassportNumber": "",
        "fatherReligion": "Islam",
        "fatherPlaceOfBirth": "Banjarmasin",
        "fatherDateOfBirth": "1965-03-15",
        "fatherCountry": "Indonesia",
        "fatherOccupation": "Wiraswasta",
        "fatherOccupationDescription": "Pedagang",
        "fatherAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
        "brideFatherPresenceStatus": "Hidup"
      },
      "mother": {
        "motherName": "Fatimah Zahra",
        "motherNik": "6371015566778899",
        "motherCitizenship": "WNI",
        "motherPassportNumber": "",
        "motherReligion": "Islam",
        "motherPlaceOfBirth": "Banjarmasin",
        "motherDateOfBirth": "1968-07-25",
        "motherCountry": "Indonesia",
        "motherOccupation": "Ibu Rumah Tangga",
        "motherOccupationDescription": "",
        "motherAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
        "brideMotherPresenceStatus": "Hidup"
      }
    },
    "guardian": {
      "guardianFullName": "Hasan Basri",
      "guardianNik": "6371011122334455",
      "guardianRelationship": "Ayah Kandung",
      "guardianStatus": "Hidup",
      "guardianReligion": "Islam",
      "guardianAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
      "guardianPhoneNumber": "081298761234"
    }
  }'
```

**Expected:** `201 Created` with complete registration data

---

### Test Case 2: Ayah Meninggal, Wali Kakek
```json
{
  "brideParents": {
    "father": {
      "fatherName": "Hasan Basri (Almarhum)",
      "fatherNik": "6371011122334455",
      "brideFatherPresenceStatus": "Meninggal"
    }
  },
  "guardian": {
    "guardianFullName": "Abdul Rahman",
    "guardianNik": "6371019988776655",
    "guardianRelationship": "Kakek",
    "guardianStatus": "Hidup",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. Sudirman No. 100, Banjarmasin",
    "guardianPhoneNumber": "081298765555"
  }
}
```

**Expected:** `201 Created`

---

### Test Case 3: ERROR - Wali Meninggal
```json
{
  "guardian": {
    "guardianFullName": "Hasan Basri",
    "guardianNik": "6371011122334455",
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Meninggal",
    "guardianReligion": "Islam",
    "guardianAddress": "Jl. Lambung Mangkurat No. 321, Banjarmasin",
    "guardianPhoneNumber": "081298761234"
  }
}
```

**Expected:** `400 Bad Request` - "Wali nikah yang telah meninggal dunia tidak dapat menjadi wali"

---

### Test Case 4: ERROR - Ayah Hidup, Tapi Pilih Kakek
```json
{
  "brideParents": {
    "father": {
      "fatherName": "Hasan Basri",
      "fatherNik": "6371011122334455",
      "brideFatherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianFullName": "Abdul Rahman",
    "guardianNik": "6371019988776655",
    "guardianRelationship": "Kakek",
    "guardianStatus": "Hidup"
  }
}
```

**Expected:** `400 Bad Request` - "Jika ayah kandung masih hidup, maka wali nikah harus Ayah Kandung"

---

### Test Case 5: ERROR - NIK Wali Berbeda dari NIK Ayah
```json
{
  "brideParents": {
    "father": {
      "fatherName": "Hasan Basri",
      "fatherNik": "6371011122334455",
      "brideFatherPresenceStatus": "Hidup"
    }
  },
  "guardian": {
    "guardianFullName": "Hasan Basri",
    "guardianNik": "1234567890123456",
    "guardianRelationship": "Ayah Kandung",
    "guardianStatus": "Hidup"
  }
}
```

**Expected:** `400 Bad Request` - "NIK wali harus sama dengan NIK ayah kandung catin perempuan"

---

## 📊 Performance Metrics

| Metric | Value | Notes |
|--------|-------|-------|
| **Average Response Time** | 800-1200ms | Optimized with async processing |
| **Database Queries** | 5 queries | Batch insert for parents |
| **External API Calls** | 0 (async) | Geocoding runs in background |
| **Transaction Size** | ~10-15KB | Complete registration data |
| **Throughput** | ~50 req/sec | Depends on server specs |

---

## 🔐 Security Notes

1. **JWT Authentication Required** - Semua request harus menyertakan valid JWT token
2. **User Isolation** - Pendaftaran hanya bisa dibuat oleh user yang authenticated
3. **Data Validation** - Semua input divalidasi sebelum disimpan
4. **SQL Injection Protection** - Menggunakan parameterized queries (GORM)
5. **Rate Limiting** - Implementasi rate limit untuk mencegah abuse

---

## 🎯 Best Practices

1. **Always validate data on client-side first** - Kurangi round-trip ke server
2. **Use proper date format** - YYYY-MM-DD untuk tanggal, HH:MM untuk waktu
3. **Check guardian relationship carefully** - Ikuti urutan wali nasab
4. **Provide complete parent data** - Untuk konsistensi validasi wali
5. **Handle async geocoding** - Koordinat mungkin `null` saat pertama kali, refresh setelah beberapa detik

---

## 📞 Support

Jika ada pertanyaan atau masalah, silakan hubungi:
- **Email:** support@simnikah.go.id
- **Documentation:** https://docs.simnikah.go.id
- **API Status:** https://status.simnikah.go.id

---

**© 2025 SimNikah - Sistem Informasi Manajemen Nikah**

