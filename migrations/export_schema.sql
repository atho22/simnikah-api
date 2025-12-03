-- ============================================================
-- SimNikah Database Schema Export
-- Generated from GORM models
-- ============================================================

-- Create database if not exists
CREATE DATABASE IF NOT EXISTS simnikah CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE simnikah;

-- Set timezone
SET time_zone = '+07:00';

-- ============================================================
-- DROP TABLES (if exists) - Uncomment if needed for fresh start
-- ============================================================
-- DROP TABLE IF EXISTS feedback_pernikahans;
-- DROP TABLE IF EXISTS wali_nikahs;
-- DROP TABLE IF EXISTS pendaftaran_nikahs;
-- DROP TABLE IF EXISTS notifikasis;
-- DROP TABLE IF EXISTS calon_pasangans;
-- DROP TABLE IF EXISTS penghulus;
-- DROP TABLE IF EXISTS staff_kuas;
-- DROP TABLE IF EXISTS users;

-- ============================================================
-- TABLE: users
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Aktif',
    nama VARCHAR(100) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_user_id (user_id),
    INDEX idx_users_status (status),
    INDEX idx_users_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: staff_kuas
-- ============================================================
CREATE TABLE IF NOT EXISTS staff_kuas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL UNIQUE,
    n_ip VARCHAR(30) UNIQUE,
    nama_lengkap VARCHAR(100) NOT NULL,
    jabatan VARCHAR(50) NOT NULL,
    bagian VARCHAR(50) NOT NULL,
    no_hp VARCHAR(15),
    email VARCHAR(100),
    alamat VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'Aktif',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_staff_kua_user_id (user_id),
    INDEX idx_staff_kua_nip (n_ip),
    INDEX idx_staff_kua_status (status),
    INDEX idx_staff_kua_jabatan (jabatan)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: penghulus
-- ============================================================
CREATE TABLE IF NOT EXISTS penghulus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL UNIQUE,
    n_ip VARCHAR(30) UNIQUE,
    nama_lengkap VARCHAR(100) NOT NULL,
    no_hp VARCHAR(15),
    email VARCHAR(100),
    alamat VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'Aktif',
    jumlah_nikah INT DEFAULT 0,
    rating DOUBLE DEFAULT 0,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_penghulu_user_id (user_id),
    INDEX idx_penghulu_nip (n_ip),
    INDEX idx_penghulu_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: calon_pasangans
-- ============================================================
CREATE TABLE IF NOT EXISTS calon_pasangans (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL UNIQUE,
    nik VARCHAR(16) UNIQUE,
    nama_lengkap VARCHAR(100) NOT NULL,
    tanggal_lahir DATETIME(3) NOT NULL,
    jenis_kelamin VARCHAR(1) NOT NULL,
    pendidikan_terakhir VARCHAR(50),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_calon_pasangan_user_id (user_id),
    INDEX idx_calon_pasangan_nik (nik)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: pendaftaran_nikahs
-- ============================================================
CREATE TABLE IF NOT EXISTS pendaftaran_nikahs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nomor_pendaftaran VARCHAR(20) NOT NULL UNIQUE,
    pendaftar_id VARCHAR(20) NOT NULL,
    calon_suami_id VARCHAR(20) NOT NULL,
    calon_istri_id VARCHAR(20) NOT NULL,
    wali_nikah_id BIGINT UNSIGNED,
    tanggal_pendaftaran DATETIME(3) NOT NULL,
    tanggal_nikah DATETIME(3) NOT NULL,
    waktu_nikah VARCHAR(10) NOT NULL,
    tempat_nikah VARCHAR(100) NOT NULL,
    alamat_akad VARCHAR(200),
    latitude DOUBLE,
    longitude DOUBLE,
    status_pendaftaran VARCHAR(40) NOT NULL DEFAULT 'Draft',
    penghulu_id BIGINT UNSIGNED,
    penghulu_assigned_by VARCHAR(20),
    penghulu_assigned_at DATETIME(3),
    catatan VARCHAR(500),
    disetujui_oleh VARCHAR(20),
    disetujui_pada DATETIME(3),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_pendaftaran_pendaftar_id (pendaftar_id),
    INDEX idx_pendaftaran_calon_suami_id (calon_suami_id),
    INDEX idx_pendaftaran_calon_istri_id (calon_istri_id),
    INDEX idx_pendaftaran_penghulu_id (penghulu_id),
    INDEX idx_pendaftaran_status_pendaftaran (status_pendaftaran),
    INDEX idx_pendaftaran_tanggal_nikah (tanggal_nikah),
    INDEX idx_pendaftaran_tanggal_pendaftaran (tanggal_pendaftaran),
    INDEX idx_pendaftaran_status_tanggal (status_pendaftaran, tanggal_nikah)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: wali_nikahs
-- ============================================================
CREATE TABLE IF NOT EXISTS wali_nikahs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pendaftaran_id BIGINT UNSIGNED NOT NULL,
    nama_dan_bin VARCHAR(100) NOT NULL,
    hubungan_wali VARCHAR(50) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_wali_nikah_pendaftaran_id (pendaftaran_id),
    FOREIGN KEY (pendaftaran_id) REFERENCES pendaftaran_nikahs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: notifikasis
-- ============================================================
CREATE TABLE IF NOT EXISTS notifikasis (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    judul VARCHAR(100) NOT NULL,
    pesan VARCHAR(500) NOT NULL,
    tipe VARCHAR(10) NOT NULL DEFAULT 'Info',
    status_baca VARCHAR(20) NOT NULL DEFAULT 'Belum Dibaca',
    tautan VARCHAR(200),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_notifikasi_user_id (user_id),
    INDEX idx_notifikasi_status_baca (status_baca),
    INDEX idx_notifikasi_tipe (tipe),
    INDEX idx_notifikasi_created_at (created_at),
    INDEX idx_notifikasi_user_status (user_id, status_baca)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TABLE: feedback_pernikahans
-- ============================================================
CREATE TABLE IF NOT EXISTS feedback_pernikahans (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pendaftaran_id BIGINT UNSIGNED NOT NULL,
    user_id VARCHAR(20) NOT NULL,
    jenis_feedback VARCHAR(20) NOT NULL,
    rating INT,
    judul VARCHAR(200) NOT NULL,
    pesan TEXT NOT NULL,
    status_baca VARCHAR(20) NOT NULL DEFAULT 'Belum Dibaca',
    dibaca_oleh VARCHAR(20),
    dibaca_pada DATETIME(3),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    INDEX idx_feedback_pendaftaran_id (pendaftaran_id),
    INDEX idx_feedback_user_id (user_id),
    INDEX idx_feedback_status_baca (status_baca),
    INDEX idx_feedback_jenis_feedback (jenis_feedback),
    FOREIGN KEY (pendaftaran_id) REFERENCES pendaftaran_nikahs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- FOREIGN KEY CONSTRAINTS
-- ============================================================
-- Note: Foreign keys are already defined in table creation above
-- Additional foreign keys can be added here if needed

-- Update pendaftaran_nikahs to reference wali_nikahs
ALTER TABLE pendaftaran_nikahs 
ADD CONSTRAINT fk_pendaftaran_wali_nikah 
FOREIGN KEY (wali_nikah_id) REFERENCES wali_nikahs(id) ON DELETE SET NULL;

-- ============================================================
-- END OF SCHEMA EXPORT
-- ============================================================

