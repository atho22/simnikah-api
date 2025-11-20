# 🗑️ Cara Drop Tabel/Kolom di Railway

Panduan lengkap untuk menghapus tabel atau kolom database di Railway.

---

## 🎯 METODE 1: Via Railway Dashboard (TERMUDAH)

### Langkah-langkah:

1. **Buka Railway Dashboard**
   - Login ke https://railway.app
   - Pilih project Anda
   - Klik **service MySQL** (bukan service aplikasi)

2. **Buka Query Editor**
   - Klik tab **"Query"** di bagian atas
   - Atau klik tombol **"Connect"** → **"Query"**

3. **Jalankan SQL Query**

   **Untuk DROP KOLOM (yang Anda butuhkan sekarang):**
   ```sql
   USE simnikah;
   ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;
   ```

   **Untuk DROP TABEL (jika diperlukan):**
   ```sql
   USE simnikah;
   DROP TABLE IF EXISTS nama_tabel;
   ```

4. **Klik "Run"** atau tekan `Ctrl+Enter`

5. **Selesai!** ✅

---

## 🎯 METODE 2: Via Railway CLI + MySQL Client

### Langkah-langkah:

1. **Install Railway CLI** (jika belum)
   ```bash
   npm install -g @railway/cli
   ```

2. **Login ke Railway**
   ```bash
   railway login
   ```

3. **Link ke Project**
   ```bash
   railway link
   ```

4. **Dapatkan Connection String**
   ```bash
   railway variables
   ```
   
   Copy nilai dari `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE`

5. **Connect ke MySQL**
   ```bash
   railway connect mysql
   ```
   
   Atau gunakan MySQL client langsung:
   ```bash
   mysql -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE
   ```

6. **Jalankan SQL**
   ```sql
   USE simnikah;
   ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;
   ```

---

## 🎯 METODE 3: Via Program Cleanup Go (OTOMATIS)

### Langkah-langkah:

1. **Update Program Cleanup untuk Railway**
   
   Program `cmd/cleanup/main.go` sudah siap, tapi perlu disesuaikan untuk Railway.

2. **Set Environment Variables di Railway**
   
   Di Railway Dashboard → Variables tab, pastikan ada:
   ```
   DATABASE_URL=mysql://user:password@host:port/database
   ```
   
   Atau gunakan format Railway:
   ```
   DATABASE_URL=${{MySQL.MYSQLCONNECTIONSTRING}}
   ```

3. **Jalankan via Railway CLI**
   ```bash
   railway run go run cmd/cleanup/main.go
   ```

4. **Atau Deploy sebagai One-off Service**
   
   - Buat service baru di Railway
   - Set command: `go run cmd/cleanup/main.go`
   - Set environment variables
   - Deploy

---

## 🎯 METODE 4: Via Koneksi Langsung dari Local

### Langkah-langkah:

1. **Dapatkan Database Credentials dari Railway**
   
   - Buka Railway Dashboard
   - Klik service MySQL
   - Klik tab **"Connect"**
   - Copy **"Public Network"** connection string atau credentials

2. **Connect dari Local Machine**
   
   ```bash
   mysql -h <MYSQL_HOST> -P <MYSQL_PORT> -u <MYSQL_USER> -p<MYSQL_PASSWORD> <MYSQL_DATABASE>
   ```

   Atau gunakan MySQL Workbench/TablePlus:
   - Host: `<MYSQL_HOST>`
   - Port: `<MYSQL_PORT>`
   - User: `<MYSQL_USER>`
   - Password: `<MYSQL_PASSWORD>`
   - Database: `<MYSQL_DATABASE>`

3. **Jalankan SQL**
   ```sql
   USE simnikah;
   ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;
   ```

---

## 🎯 METODE 5: Via Aplikasi API Endpoint (UNTUK DEVELOPMENT)

### Buat Endpoint Admin untuk Cleanup:

**Tidak disarankan untuk production!** Hanya untuk development/testing.

---

## 📋 SQL Scripts yang Siap Pakai

### 1. Drop Kolom `tempat_lahir` (SOLUSI CEPAT)
```sql
USE simnikah;
ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;
```

### 2. Drop Semua Kolom yang Tidak Diperlukan
```sql
USE simnikah;

-- Drop kolom dari calon_pasangans
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS tempat_lahir;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS alamat;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS rt;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS rw;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS kelurahan;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS kecamatan;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS kabupaten;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS provinsi;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS agama;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS status_perkawinan;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS no_hp;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS email;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS warga_negara;
ALTER TABLE calon_pasangans DROP COLUMN IF EXISTS deskripsi_pekerjaan;

-- Drop kolom dari pendaftaran_nikahs
ALTER TABLE pendaftaran_nikahs DROP COLUMN IF EXISTS nomor_dispensasi;
ALTER TABLE pendaftaran_nikahs DROP COLUMN IF EXISTS status_bimbingan;
```

> **Note:** MySQL tidak support `DROP COLUMN IF EXISTS` secara langsung.
> Jika kolom tidak ada, query akan error tapi tidak masalah.
> Atau gunakan program cleanup Go untuk check kolom dulu.

### 3. Check Kolom yang Ada (Sebelum Drop)
```sql
USE simnikah;

-- Lihat semua kolom di calon_pasangans
SHOW COLUMNS FROM calon_pasangans;

-- Lihat semua kolom di pendaftaran_nikahs
SHOW COLUMNS FROM pendaftaran_nikahs;

-- Cek apakah kolom ada
SELECT COLUMN_NAME 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'simnikah' 
  AND TABLE_NAME = 'calon_pasangans' 
  AND COLUMN_NAME = 'tempat_lahir';
```

---

## ✅ Rekomendasi: METODE 1 (Via Railway Dashboard)

**Ini cara TERMUDAH dan TERCEPAT:**

1. Buka Railway Dashboard
2. Klik service MySQL
3. Klik tab "Query"
4. Jalankan:
   ```sql
   USE simnikah;
   ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;
   ```
5. Done! ✅

---

## 🔒 Security Note

⚠️ **PENTING:** 
- Hati-hati saat drop kolom/tabel, pastikan tidak ada data penting
- Selalu backup database sebelum melakukan perubahan besar
- Gunakan `DROP COLUMN IF EXISTS` jika tersedia atau check dulu dengan `SHOW COLUMNS`
- Untuk production, pertimbangkan membuat migration script yang bisa di-rollback

---

## 🆘 Troubleshooting

### Error: "Column doesn't exist"
- **Solusi:** Kolom sudah tidak ada, tidak masalah. Lanjutkan.

### Error: "Access denied"
- **Solusi:** Pastikan menggunakan user yang punya permission DROP COLUMN

### Error: "Table is locked"
- **Solusi:** Pastikan tidak ada transaction yang sedang berjalan atau restart aplikasi

---

**Created:** 2024  
**Last Updated:** 2024


