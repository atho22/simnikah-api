# Validasi Wali Nikah Sesuai Syariat Islam

## 📋 Pendahuluan

Sistem SimNikah mengimplementasikan validasi wali nikah sesuai dengan syariat Islam. Wali nikah adalah salah satu rukun nikah yang sangat penting dalam pernikahan Islam.

## 🕌 Aturan Wali Nikah dalam Islam

### Urutan Wali Nasab

Wali nikah untuk calon pengantin perempuan mengikuti urutan wali nasab sebagai berikut:

1. **Ayah Kandung** - Wali yang paling berhak
2. **Kakek** - Ayah dari ayah (jika ayah meninggal/tidak ada)
3. **Saudara Laki-Laki Kandung** - Saudara sekandung
4. **Saudara Laki-Laki Seayah** - Saudara seayah
5. **Keponakan Laki-Laki** - Anak laki-laki dari saudara laki-laki kandung
6. **Paman Kandung** - Saudara kandung ayah
7. **Paman Seayah** - Saudara seayah dari ayah
8. **Sepupu Laki-Laki** - Anak laki-laki dari paman kandung
9. **Wali Hakim** - Jika tidak ada wali nasab yang memenuhi syarat

### Syarat Wali Nikah

Wali nikah harus memenuhi syarat:
- Laki-laki
- Beragama Islam
- Baligh (dewasa)
- Berakal sehat
- Adil
- Tidak sedang ihram (haji/umrah)

## 🔍 Implementasi Validasi di SimNikah

### 1. Constants untuk Hubungan Wali

File: `structs/constants.go`

```go
const (
    WaliHubunganAyahKandung            = "Ayah Kandung"
    WaliHubunganKakek                  = "Kakek"
    WaliHubunganSaudaraLakiLakiKandung = "Saudara Laki-Laki Kandung"
    WaliHubunganSaudaraLakiLakiSeayah  = "Saudara Laki-Laki Seayah"
    WaliHubunganKeponakanLakiLaki      = "Keponakan Laki-Laki"
    WaliHubunganPamanKandung           = "Paman Kandung"
    WaliHubunganPamanSeayah            = "Paman Seayah"
    WaliHubunganSepupuLakiLaki         = "Sepupu Laki-Laki"
    WaliHubunganWaliHakim              = "Wali Hakim"
    WaliHubunganLainnya                = "Lainnya"
)
```

### 2. Fungsi Validasi

#### `IsValidWaliNikah(hubunganWali, statusAyahCalonIstri string) bool`

Fungsi ini memvalidasi apakah hubungan wali yang dipilih sesuai dengan status ayah calon pengantin perempuan.

**Aturan Validasi:**

1. **Jika Ayah Kandung Masih Hidup:**
   - Wali nikah HARUS "Ayah Kandung"
   - Tidak boleh menggunakan wali lain

2. **Jika Ayah Kandung Meninggal/Tidak Ada:**
   - Wali nikah berpindah ke nasab berikutnya
   - Bisa memilih: Kakek, Saudara Laki-Laki, Paman, dll
   - TIDAK BOLEH memilih "Ayah Kandung"

**Contoh Penggunaan:**

```go
// Contoh 1: Ayah masih hidup
statusAyah := structs.StatusKeberadaanHidup
hubunganWali := structs.WaliHubunganAyahKandung

isValid := structs.IsValidWaliNikah(hubunganWali, statusAyah)
// Result: true ✅

// Contoh 2: Ayah masih hidup, tapi pilih paman (TIDAK VALID)
statusAyah := structs.StatusKeberadaanHidup
hubunganWali := structs.WaliHubunganPamanKandung

isValid := structs.IsValidWaliNikah(hubunganWali, statusAyah)
// Result: false ❌

// Contoh 3: Ayah meninggal, pilih saudara kandung
statusAyah := structs.StatusKeberadaanMeninggal
hubunganWali := structs.WaliHubunganSaudaraLakiLakiKandung

isValid := structs.IsValidWaliNikah(hubunganWali, statusAyah)
// Result: true ✅

// Contoh 4: Ayah meninggal, tapi pilih ayah (TIDAK VALID)
statusAyah := structs.StatusKeberadaanMeninggal
hubunganWali := structs.WaliHubunganAyahKandung

isValid := structs.IsValidWaliNikah(hubunganWali, statusAyah)
// Result: false ❌
```

#### `GetUrutanWaliNasab() []string`

Mengembalikan daftar urutan wali nasab dari yang paling berhak.

**Contoh Penggunaan:**

```go
urutan := structs.GetUrutanWaliNasab()
// Result:
// [
//   "Ayah Kandung",
//   "Kakek",
//   "Saudara Laki-Laki Kandung",
//   "Saudara Laki-Laki Seayah",
//   "Keponakan Laki-Laki",
//   "Paman Kandung",
//   "Paman Seayah",
//   "Sepupu Laki-Laki",
//   "Wali Hakim"
// ]
```

#### `GetPesanValidasiWali(statusAyahCalonIstri string) string`

Memberikan pesan validasi yang informatif berdasarkan status ayah.

**Contoh Penggunaan:**

```go
// Jika ayah masih hidup
pesan := structs.GetPesanValidasiWali(structs.StatusKeberadaanHidup)
// Result: "Jika ayah kandung masih hidup, maka wali nikah harus Ayah Kandung sesuai syariat Islam"

// Jika ayah meninggal
pesan := structs.GetPesanValidasiWali(structs.StatusKeberadaanMeninggal)
// Result: "Jika ayah kandung meninggal/tidak ada, wali nikah berpindah ke nasab berikutnya: Kakek, Saudara Laki-Laki Kandung, Paman, atau Wali Hakim"
```

### 3. Validasi di Endpoint Pendaftaran Nikah

File: `catin/daftar.go`

Validasi wali nikah otomatis dijalankan saat user melakukan pendaftaran nikah:

```go
// Validasi hubungan wali berdasarkan status ayah kandung calon istri
statusAyahCalonIstri := dataFormPendaftaran.OrangTuaCalonIstri.Ayah.StatusKeberadaan
hubunganWali := dataFormPendaftaran.WaliNikah.HubunganWali

// Validasi apakah hubungan wali valid berdasarkan status ayah
if !structs.IsValidWaliNikah(hubunganWali, statusAyahCalonIstri) {
    c.JSON(http.StatusBadRequest, gin.H{
        "success": false,
        "message": "Validasi Wali Nikah Gagal",
        "error":   structs.GetPesanValidasiWali(statusAyahCalonIstri),
        "field":   "hubungan_wali",
        "type":    "syariat_validation",
        "details": gin.H{
            "status_ayah_catin_perempuan": statusAyahCalonIstri,
            "hubungan_wali_yang_dipilih":  hubunganWali,
            "urutan_wali_nasab":           structs.GetUrutanWaliNasab(),
        },
    })
    return
}
```

## 📝 Contoh Request/Response

### ✅ Request Valid - Ayah Masih Hidup

```json
{
  "orang_tua_catin_istri": {
    "ayah": {
      "nama": "Bapak Ahmad",
      "status_keberadaan": "Hidup"
    }
  },
  "wali_nikah": {
    "nama_lengkap_wali": "Bapak Ahmad",
    "hubungan_wali": "Ayah Kandung",
    "status_wali": "Hidup"
  }
}
```

**Response:** ✅ Success

### ❌ Request Tidak Valid - Ayah Hidup tapi Pilih Paman

```json
{
  "orang_tua_catin_istri": {
    "ayah": {
      "nama": "Bapak Ahmad",
      "status_keberadaan": "Hidup"
    }
  },
  "wali_nikah": {
    "nama_lengkap_wali": "Bapak Usman",
    "hubungan_wali": "Paman Kandung",
    "status_wali": "Hidup"
  }
}
```

**Response:** ❌ Error
```json
{
  "success": false,
  "message": "Validasi Wali Nikah Gagal",
  "error": "Jika ayah kandung masih hidup, maka wali nikah harus Ayah Kandung sesuai syariat Islam",
  "field": "hubungan_wali",
  "type": "syariat_validation",
  "details": {
    "status_ayah_catin_perempuan": "Hidup",
    "hubungan_wali_yang_dipilih": "Paman Kandung",
    "urutan_wali_nasab": [
      "Ayah Kandung",
      "Kakek",
      "Saudara Laki-Laki Kandung",
      "..."
    ]
  }
}
```

### ✅ Request Valid - Ayah Meninggal, Pilih Saudara

```json
{
  "orang_tua_catin_istri": {
    "ayah": {
      "nama": "Bapak Ahmad (Alm)",
      "status_keberadaan": "Meninggal"
    }
  },
  "wali_nikah": {
    "nama_lengkap_wali": "Bapak Ali",
    "hubungan_wali": "Saudara Laki-Laki Kandung",
    "status_wali": "Hidup"
  }
}
```

**Response:** ✅ Success

## 🎯 Best Practices

### 1. UI/Frontend

Pada form pendaftaran nikah di frontend:

```javascript
// Saat user memilih status ayah kandung
const statusAyah = form.orangTuaCatinIstri.ayah.statusKeberadaan;

// Tampilkan pilihan wali yang sesuai
if (statusAyah === "Hidup") {
  // HANYA tampilkan pilihan "Ayah Kandung"
  waliOptions = ["Ayah Kandung"];
  showWarning("Ayah kandung masih hidup, wali nikah harus Ayah Kandung");
} else if (statusAyah === "Meninggal") {
  // Tampilkan pilihan wali nasab selain ayah
  waliOptions = [
    "Kakek",
    "Saudara Laki-Laki Kandung",
    "Saudara Laki-Laki Seayah",
    "Keponakan Laki-Laki",
    "Paman Kandung",
    "Paman Seayah",
    "Sepupu Laki-Laki",
    "Wali Hakim"
  ];
  showInfo("Pilih wali nikah sesuai urutan wali nasab");
}
```

### 2. Informasi untuk User

Berikan informasi yang jelas kepada user:

- Tampilkan tooltip/help text tentang urutan wali nasab
- Berikan warning jika pilihan tidak sesuai syariat
- Sediakan link ke dokumentasi/artikel tentang wali nikah

### 3. Edge Cases

**Kasus Wali Hakim:**

Wali hakim diperlukan jika:
- Tidak ada wali nasab yang memenuhi syarat
- Wali nasab menolak/tidak mau menjadi wali (wali adlal)
- Wali nasab tidak berada di tempat dan sulit dihubungi

**Kasus Khusus:**

Sistem mendukung pilihan "Lainnya" untuk kasus-kasus khusus yang perlu verifikasi manual oleh staff KUA.

## 📚 Referensi

1. **Al-Quran Surah An-Nisa ayat 25**
2. **Hadits Riwayat Tirmidzi**: "Tidak sah nikah kecuali dengan wali"
3. **Fiqih Islam** - Kitab An-Nikah
4. **Kompilasi Hukum Islam (KHI)** Pasal 19-23

## ⚠️ Catatan Penting

- Validasi ini membantu memastikan pendaftaran nikah sesuai syariat Islam
- Staff KUA tetap melakukan verifikasi manual saat user datang ke kantor
- Untuk kasus-kasus khusus yang memerlukan pertimbangan ulama, konsultasikan dengan penghulu atau kepala KUA
- Sistem ini adalah alat bantu, keputusan akhir tetap di tangan petugas KUA yang berwenang

---

**Terakhir diupdate:** 2025-01-27
**Versi:** 1.0

