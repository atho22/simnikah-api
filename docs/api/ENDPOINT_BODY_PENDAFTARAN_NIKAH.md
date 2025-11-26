# 📋 Endpoint Body Request - Pendaftaran Nikah

**Versi:** 1.3.0  
**Update Terakhir:** November 2024

---

## 📌 Endpoint: Create Marriage Registration

### **POST** `/simnikah/pendaftaran`

Membuat pendaftaran nikah dengan form sederhana. Endpoint ini digunakan oleh user biasa (`user_biasa`) untuk mendaftarkan pernikahan mereka.

---

## 🔐 Authentication

**Required:** ✅ Yes  
**Type:** Bearer Token (JWT)  
**Header:** 
```
Authorization: Bearer <your_jwt_token>
```

---

## 📥 Request Body Structure

### Content-Type
```
Content-Type: application/json
```

### Complete Request Body (Untuk Nikah di KUA)

```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
    "pendidikan_akhir": "S1",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
    "pendidikan_akhir": "S1",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di KUA",
    "tanggal_nikah": "2024-12-25",
    "waktu_nikah": "09:00"
  },
  "wali_nikah": {
    "nama_dan_bin": "Abdullah bin Muhammad",
    "hubungan_wali": "Ayah Kandung"
  }
}
```

### Complete Request Body (Untuk Nikah di Luar KUA)

```json
{
  "calon_laki_laki": {
    "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
    "pendidikan_akhir": "S1",
    "umur": 25
  },
  "calon_perempuan": {
    "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
    "pendidikan_akhir": "S1",
    "umur": 23
  },
  "lokasi_nikah": {
    "tempat_nikah": "Di Luar KUA",
    "tanggal_nikah": "2024-12-25",
    "waktu_nikah": "09:00",
    "alamat_nikah": "Jl. Ahmad Yani No. 123",
    "detail_alamat": "Rumah Pengantin Perempuan, RT 05 RW 02",
    "kelurahan": "Pangeran"
  },
  "wali_nikah": {
    "nama_dan_bin": "Abdullah bin Muhammad",
    "hubungan_wali": "Ayah Kandung"
  }
}
```

---

## 📝 Field Descriptions

### 1. `calon_laki_laki` (Required)

Data calon pengantin laki-laki.

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `nama_dan_bin` | string | ✅ Yes | Nama lengkap dengan bin (nama ayah) | `"Ahmad Wijaya bin Abdullah"` |
| `pendidikan_akhir` | string | ✅ Yes | Pendidikan terakhir | `"S1"`, `"SMA"`, `"D3"`, dll |
| `umur` | integer | ✅ Yes | Umur calon laki-laki (minimal 19 tahun) | `25` |

### 2. `calon_perempuan` (Required)

Data calon pengantin perempuan.

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `nama_dan_binti` | string | ✅ Yes | Nama lengkap dengan binti (nama ayah) | `"Siti Nurhaliza binti Muhammad"` |
| `pendidikan_akhir` | string | ✅ Yes | Pendidikan terakhir | `"S1"`, `"SMA"`, `"D3"`, dll |
| `umur` | integer | ✅ Yes | Umur calon perempuan (minimal 19 tahun) | `23` |

### 3. `lokasi_nikah` (Required)

Data lokasi dan waktu pernikahan.

| Field | Type | Required | Condition | Description | Example |
|-------|------|----------|-----------|-------------|---------|
| `tempat_nikah` | string | ✅ Yes | - | Tempat nikah: `"Di KUA"` atau `"Di Luar KUA"` | `"Di KUA"` |
| `tanggal_nikah` | string | ✅ Yes | - | Tanggal nikah (format: `YYYY-MM-DD`). Tidak boleh di masa lalu | `"2024-12-25"` |
| `waktu_nikah` | string | ✅ Yes | - | Waktu nikah (format: `HH:MM`, 24-jam) | `"09:00"`, `"14:30"` |
| `alamat_nikah` | string | ⚠️ Conditional | Required jika `tempat_nikah = "Di Luar KUA"` | Alamat lengkap lokasi nikah | `"Jl. Ahmad Yani No. 123"` |
| `detail_alamat` | string | ❌ Optional | Hanya jika `tempat_nikah = "Di Luar KUA"` | Detail alamat (RT/RW/dll) | `"Rumah Pengantin Perempuan, RT 05 RW 02"` |
| `kelurahan` | string | ⚠️ Conditional | Required jika `tempat_nikah = "Di Luar KUA"` | Nama kelurahan (harus dalam lingkup Kecamatan Banjarmasin Utara) | `"Pangeran"` |

**Kelurahan Valid (Kecamatan Banjarmasin Utara):**
- Alalak Utara
- Alalak Tengah
- Alalak Selatan
- Antasan Kecil Timur
- Kuin Utara
- Pangeran
- Sungai Miai
- Sungai Andai
- Surgi Mufti

### 4. `wali_nikah` (Required)

Data wali nikah untuk calon pengantin perempuan.

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `nama_dan_bin` | string | ✅ Yes | Nama lengkap wali dengan bin (nama ayah) | `"Abdullah bin Muhammad"` |
| `hubungan_wali` | string | ✅ Yes | Hubungan nasab wali (lihat daftar di bawah) | `"Ayah Kandung"` |

**Hubungan Wali Valid (Urutan Wali Nasab):**

1. `"Ayah Kandung"` - Wali yang paling berhak (jika ayah masih hidup)
2. `"Kakek"` - Ayah dari ayah (jika ayah meninggal)
3. `"Saudara Laki-Laki Kandung"` - Saudara sekandung
4. `"Saudara Laki-Laki Seayah"` - Saudara seayah
5. `"Keponakan Laki-Laki"` - Anak laki-laki dari saudara
6. `"Paman Kandung"` - Saudara kandung ayah
7. `"Paman Seayah"` - Saudara seayah dari ayah
8. `"Sepupu Laki-Laki"` - Anak laki-laki dari paman
9. `"Wali Hakim"` - Jika tidak ada wali nasab yang memenuhi syarat
10. `"Lainnya"` - Hubungan lainnya

---

## ✅ Validations

### Umur
- **Calon Laki-Laki:** Minimal 19 tahun
- **Calon Perempuan:** Minimal 19 tahun

### Format Tanggal
- Format: `YYYY-MM-DD`
- Contoh: `"2024-12-25"`
- **Tidak boleh di masa lalu**

### Format Waktu
- Format: `HH:MM` (24-jam)
- Contoh: `"09:00"`, `"14:30"`, `"16:00"`

### Tempat Nikah
- Harus salah satu dari: `"Di KUA"` atau `"Di Luar KUA"`
- Jika `"Di Luar KUA"`, field `alamat_nikah` dan `kelurahan` wajib diisi

### Kelurahan
- Harus dalam lingkup **Kecamatan Banjarmasin Utara**
- Lihat daftar kelurahan valid di atas

### Wali Nikah
- **Wajib diisi** untuk semua pendaftaran
- `nama_dan_bin` tidak boleh kosong
- `hubungan_wali` harus dari daftar hubungan wali valid

### Ketersediaan Jadwal
- Sistem akan mengecek ketersediaan jadwal berdasarkan:
  - Tanggal dan waktu yang dipilih
  - Kapasitas per jam (maksimal 3 pernikahan per jam)
  - Jika `"Di KUA"`: maksimal 1 pernikahan per jam
  - Jika `"Di Luar KUA"`: maksimal 2 pernikahan per jam (jika sudah ada 1 di KUA)

---

## 📤 Response Examples

### Success Response (201 Created)

```json
{
  "success": true,
  "message": "Pendaftaran nikah berhasil dibuat (form sederhana)",
  "data": {
    "id": 1,
    "nomor_pendaftaran": "NIKAH-20241225-1234",
    "status_pendaftaran": "Draft",
    "tanggal_nikah": "2024-12-25T00:00:00Z",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di KUA",
    "alamat_akad": "PH5Q+F8C, Jl. Wira Karya, Pangeran, Kec. Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan 70123",
    "calon_suami": {
      "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
      "pendidikan": "S1",
      "umur": 25
    },
    "calon_istri": {
      "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
      "pendidikan": "S1",
      "umur": 23
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    },
    "created_at": "2024-11-25T10:00:00Z"
  }
}
```

### Error Response (400 Bad Request)

**Contoh: Umur tidak valid**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Umur calon laki-laki minimal 19 tahun",
  "field": "umur_laki_laki",
  "type": "validation"
}
```

**Contoh: Format tanggal tidak valid**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD (contoh: 2024-12-25)",
  "field": "tanggal_nikah",
  "type": "validation"
}
```

**Contoh: Jadwal sudah penuh**
```json
{
  "success": false,
  "message": "Jadwal sudah penuh",
  "error": "Jadwal pernikahan pada tanggal 2024-12-25 pukul 09:00 sudah penuh. Maksimal 3 pernikahan per jam.",
  "field": "waktu_nikah",
  "type": "schedule_conflict"
}
```

**Contoh: Wali nikah tidak valid**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Hubungan wali tidak valid",
  "field": "wali_nikah.hubungan_wali",
  "type": "validation"
}
```

---

## 🔄 Status Flow

Setelah pendaftaran dibuat, status akan mengikuti alur berikut:

1. **Draft** - Pendaftaran baru dibuat (default)
2. **Disetujui** - Staff/Kepala KUA menyetujui setelah verifikasi
3. **Menunggu Penugasan** - Menunggu kepala KUA menentukan penghulu
4. **Penghulu Ditugaskan** - Kepala KUA sudah assign penghulu
5. **Selesai** - Penghulu sudah melaksanakan nikah

**Catatan:** Status "Draft" juga dihitung dalam kuota jadwal untuk mencegah double booking.

---

## 📋 Complete Example (cURL)

### Nikah di KUA

```bash
curl -X POST https://your-api-domain.com/simnikah/pendaftaran \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "calon_laki_laki": {
      "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
      "pendidikan_akhir": "S1",
      "umur": 25
    },
    "calon_perempuan": {
      "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
      "pendidikan_akhir": "S1",
      "umur": 23
    },
    "lokasi_nikah": {
      "tempat_nikah": "Di KUA",
      "tanggal_nikah": "2024-12-25",
      "waktu_nikah": "09:00"
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    }
  }'
```

### Nikah di Luar KUA

```bash
curl -X POST https://your-api-domain.com/simnikah/pendaftaran \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "calon_laki_laki": {
      "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
      "pendidikan_akhir": "S1",
      "umur": 25
    },
    "calon_perempuan": {
      "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
      "pendidikan_akhir": "S1",
      "umur": 23
    },
    "lokasi_nikah": {
      "tempat_nikah": "Di Luar KUA",
      "tanggal_nikah": "2024-12-25",
      "waktu_nikah": "09:00",
      "alamat_nikah": "Jl. Ahmad Yani No. 123",
      "detail_alamat": "Rumah Pengantin Perempuan, RT 05 RW 02",
      "kelurahan": "Pangeran"
    },
    "wali_nikah": {
      "nama_dan_bin": "Abdullah bin Muhammad",
      "hubungan_wali": "Ayah Kandung"
    }
  }'
```

---

## 📋 Complete Example (JavaScript/Fetch)

```javascript
const response = await fetch('https://your-api-domain.com/simnikah/pendaftaran', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${jwtToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    calon_laki_laki: {
      nama_dan_bin: "Ahmad Wijaya bin Abdullah",
      pendidikan_akhir: "S1",
      umur: 25
    },
    calon_perempuan: {
      nama_dan_binti: "Siti Nurhaliza binti Muhammad",
      pendidikan_akhir: "S1",
      umur: 23
    },
    lokasi_nikah: {
      tempat_nikah: "Di KUA",
      tanggal_nikah: "2024-12-25",
      waktu_nikah: "09:00"
    },
    wali_nikah: {
      nama_dan_bin: "Abdullah bin Muhammad",
      hubungan_wali: "Ayah Kandung"
    }
  })
});

const data = await response.json();
console.log(data);
```

---

## 📋 Complete Example (Python/Requests)

```python
import requests

url = "https://your-api-domain.com/simnikah/pendaftaran"
headers = {
    "Authorization": f"Bearer {jwt_token}",
    "Content-Type": "application/json"
}
payload = {
    "calon_laki_laki": {
        "nama_dan_bin": "Ahmad Wijaya bin Abdullah",
        "pendidikan_akhir": "S1",
        "umur": 25
    },
    "calon_perempuan": {
        "nama_dan_binti": "Siti Nurhaliza binti Muhammad",
        "pendidikan_akhir": "S1",
        "umur": 23
    },
    "lokasi_nikah": {
        "tempat_nikah": "Di KUA",
        "tanggal_nikah": "2024-12-25",
        "waktu_nikah": "09:00"
    },
    "wali_nikah": {
        "nama_dan_bin": "Abdullah bin Muhammad",
        "hubungan_wali": "Ayah Kandung"
    }
}

response = requests.post(url, json=payload, headers=headers)
data = response.json()
print(data)
```

---

## 📋 Complete Example (PHP/cURL)

```php
<?php
$url = "https://your-api-domain.com/simnikah/pendaftaran";
$token = "YOUR_JWT_TOKEN";

$data = [
    "calon_laki_laki" => [
        "nama_dan_bin" => "Ahmad Wijaya bin Abdullah",
        "pendidikan_akhir" => "S1",
        "umur" => 25
    ],
    "calon_perempuan" => [
        "nama_dan_binti" => "Siti Nurhaliza binti Muhammad",
        "pendidikan_akhir" => "S1",
        "umur" => 23
    ],
    "lokasi_nikah" => [
        "tempat_nikah" => "Di KUA",
        "tanggal_nikah" => "2024-12-25",
        "waktu_nikah" => "09:00"
    ],
    "wali_nikah" => [
        "nama_dan_bin" => "Abdullah bin Muhammad",
        "hubungan_wali" => "Ayah Kandung"
    ]
];

$ch = curl_init($url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
curl_setopt($ch, CURLOPT_HTTPHEADER, [
    "Authorization: Bearer " . $token,
    "Content-Type: application/json"
]);

$response = curl_exec($ch);
curl_close($ch);

$result = json_decode($response, true);
print_r($result);
?>
```

---

## 📋 Complete Example (Go)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type RequestBody struct {
    CalonLakiLaki struct {
        NamaDanBin      string `json:"nama_dan_bin"`
        PendidikanAkhir string `json:"pendidikan_akhir"`
        Umur            int    `json:"umur"`
    } `json:"calon_laki_laki"`
    CalonPerempuan struct {
        NamaDanBinti    string `json:"nama_dan_binti"`
        PendidikanAkhir string `json:"pendidikan_akhir"`
        Umur            int    `json:"umur"`
    } `json:"calon_perempuan"`
    LokasiNikah struct {
        TempatNikah  string `json:"tempat_nikah"`
        TanggalNikah string `json:"tanggal_nikah"`
        WaktuNikah   string `json:"waktu_nikah"`
    } `json:"lokasi_nikah"`
    WaliNikah struct {
        NamaDanBin   string `json:"nama_dan_bin"`
        HubunganWali string `json:"hubungan_wali"`
    } `json:"wali_nikah"`
}

func main() {
    url := "https://your-api-domain.com/simnikah/pendaftaran"
    token := "YOUR_JWT_TOKEN"

    reqBody := RequestBody{
        CalonLakiLaki: struct {
            NamaDanBin      string `json:"nama_dan_bin"`
            PendidikanAkhir string `json:"pendidikan_akhir"`
            Umur            int    `json:"umur"`
        }{
            NamaDanBin:      "Ahmad Wijaya bin Abdullah",
            PendidikanAkhir: "S1",
            Umur:            25,
        },
        CalonPerempuan: struct {
            NamaDanBinti    string `json:"nama_dan_binti"`
            PendidikanAkhir string `json:"pendidikan_akhir"`
            Umur            int    `json:"umur"`
        }{
            NamaDanBinti:    "Siti Nurhaliza binti Muhammad",
            PendidikanAkhir: "S1",
            Umur:            23,
        },
        LokasiNikah: struct {
            TempatNikah  string `json:"tempat_nikah"`
            TanggalNikah string `json:"tanggal_nikah"`
            WaktuNikah   string `json:"waktu_nikah"`
        }{
            TempatNikah:  "Di KUA",
            TanggalNikah: "2024-12-25",
            WaktuNikah:   "09:00",
        },
        WaliNikah: struct {
            NamaDanBin   string `json:"nama_dan_bin"`
            HubunganWali string `json:"hubungan_wali"`
        }{
            NamaDanBin:   "Abdullah bin Muhammad",
            HubunganWali: "Ayah Kandung",
        },
    }

    jsonData, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Println(result)
}
```

---

## ⚠️ Important Notes

1. **Wali Nikah Wajib:** Field `wali_nikah` adalah **wajib** untuk semua pendaftaran nikah
2. **Status Draft di Kuota:** Status "Draft" juga dihitung dalam kuota jadwal untuk mencegah double booking
3. **Format Tanggal:** Gunakan format `YYYY-MM-DD` (contoh: `2024-12-25`)
4. **Format Waktu:** Gunakan format `HH:MM` dalam 24-jam (contoh: `09:00`, `14:30`)
5. **Tanggal Tidak Boleh Masa Lalu:** Sistem akan menolak jika tanggal nikah di masa lalu
6. **Kapasitas Per Jam:** 
   - Maksimal 3 pernikahan per jam
   - Jika "Di KUA": maksimal 1 per jam
   - Jika "Di Luar KUA": maksimal 2 per jam (jika sudah ada 1 di KUA)

---

## 🔗 Related Endpoints

- **GET** `/simnikah/kalender-ketersediaan` - Cek ketersediaan tanggal
- **GET** `/simnikah/ketersediaan-jam` - Cek ketersediaan jam untuk tanggal tertentu
- **GET** `/simnikah/pendaftaran/status` - Cek status pendaftaran user
- **POST** `/simnikah/staff/pendaftaran` - Staff membuat pendaftaran untuk user (otomatis Disetujui)

---

**Last Updated:** November 2024  
**API Version:** 1.3.0



