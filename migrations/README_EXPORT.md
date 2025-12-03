# Database Export SQL

## File yang Tersedia

### 1. `export_schema.sql`
File ini berisi struktur database lengkap (schema) untuk SimNikah API, termasuk:
- CREATE DATABASE statement
- CREATE TABLE untuk semua tabel
- CREATE INDEX untuk optimasi performa
- FOREIGN KEY constraints

## Cara Menggunakan

### Import Schema ke Database Baru

```bash
# Masuk ke MySQL
mysql -u root -p

# Atau langsung import
mysql -u root -p < migrations/export_schema.sql
```

### Atau menggunakan MySQL client

```sql
SOURCE migrations/export_schema.sql;
```

## Export Data (jika diperlukan)

Untuk export data dari database yang sudah ada, gunakan `mysqldump`:

```bash
# Export schema saja
mysqldump -u root -p --no-data simnikah > migrations/export_schema_only.sql

# Export data saja
mysqldump -u root -p --no-create-info simnikah > migrations/export_data_only.sql

# Export schema + data (full backup)
mysqldump -u root -p simnikah > migrations/export_full_backup.sql
```

## Struktur Tabel

File `export_schema.sql` mencakup tabel-tabel berikut:

1. **users** - Data pengguna dan autentikasi
2. **staff_kuas** - Data staff KUA
3. **penghulus** - Data penghulu
4. **calon_pasangans** - Data calon pasangan
5. **pendaftaran_nikahs** - Data pendaftaran nikah
6. **wali_nikahs** - Data wali nikah
7. **notifikasis** - Data notifikasi
8. **feedback_pernikahans** - Data feedback pernikahan

## Catatan

- File ini menggunakan `CREATE TABLE IF NOT EXISTS`, jadi aman untuk dijalankan berulang kali
- Semua tabel menggunakan charset `utf8mb4` dan collation `utf8mb4_unicode_ci`
- Indexes sudah dioptimasi untuk performa query yang lebih baik
- Foreign keys menggunakan `ON DELETE CASCADE` atau `ON DELETE SET NULL` sesuai kebutuhan

