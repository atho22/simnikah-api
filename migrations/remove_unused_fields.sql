-- Migration Script: Remove Unused Fields from Database Tables
-- Menghapus field-field yang tidak diperlukan sesuai dengan form sederhana
-- Execute this script manually or via your database management tool

USE simnikah;

-- ==================== CALON PASANGAN TABLE ====================
-- Hapus kolom yang tidak diperlukan dari tabel calon_pasangans
-- Field yang dihapus sudah tidak ada di model CalonPasangan

-- IMPORTANT: MySQL tidak support DROP COLUMN IF EXISTS
-- Jika kolom tidak ada, query akan error tapi tidak masalah
-- Atau jalankan program cleanup: go run cmd/cleanup/main.go

-- Hapus kolom tempat_lahir (tidak diperlukan di form sederhana)
-- Error akan muncul jika kolom sudah tidak ada, tapi itu OK
ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;

-- Hapus kolom alamat (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN alamat;

-- Hapus kolom rt (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN rt;

-- Hapus kolom rw (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN rw;

-- Hapus kolom kelurahan (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN kelurahan;

-- Hapus kolom kecamatan (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN kecamatan;

-- Hapus kolom kabupaten (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN kabupaten;

-- Hapus kolom provinsi (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN provinsi;

-- Hapus kolom agama (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN agama;

-- Hapus kolom status_perkawinan (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN status_perkawinan;

-- Hapus kolom no_hp (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN no_hp;

-- Hapus kolom email (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN email;

-- Hapus kolom warga_negara (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN warga_negara;

-- Hapus kolom deskripsi_pekerjaan (tidak diperlukan di form sederhana)
-- ALTER TABLE calon_pasangans DROP COLUMN deskripsi_pekerjaan;

-- ==================== PENDAFTARAN NIKAH TABLE ====================
-- Hapus kolom yang tidak diperlukan dari tabel pendaftaran_nikahs

-- Hapus kolom nomor_dispensasi (tidak diperlukan di form sederhana)
-- ALTER TABLE pendaftaran_nikahs DROP COLUMN nomor_dispensasi;

-- Hapus kolom status_bimbingan (tidak diperlukan di form sederhana)
-- ALTER TABLE pendaftaran_nikahs DROP COLUMN status_bimbingan;

-- ==================== CARA MENJALANKAN ====================
-- 
-- OPSI 1: Jalankan SQL langsung (CARA CEPAT)
-- mysql -u simnikah_user -p simnikah < migrations/remove_unused_fields.sql
-- 
-- OPSI 2: Gunakan program cleanup Go (LEBIH AMAN - check kolom dulu)
-- go run cmd/cleanup/main.go
-- 
-- OPSI 3: Jalankan manual di MySQL client
-- mysql -u simnikah_user -p
-- USE simnikah;
-- ALTER TABLE calon_pasangans DROP COLUMN tempat_lahir;
--
-- Untuk memeriksa kolom yang ada sebelum menghapus:
-- SHOW COLUMNS FROM calon_pasangans;
-- SHOW COLUMNS FROM pendaftaran_nikahs;

