package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Get database connection string from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback to default connection string
		dsn = "simnikah_user:simnikah_password@tcp(localhost:3306)/simnikah?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// Connect to database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🔧 Starting cleanup of unused database fields...")
	fmt.Println()

	// List of columns to remove from calon_pasangans table
	calonPasanganColumns := []string{
		"tempat_lahir",
		"alamat",
		"rt",
		"rw",
		"kelurahan",
		"kecamatan",
		"kabupaten",
		"provinsi",
		"agama",
		"status_perkawinan",
		"no_hp",
		"email",
		"warga_negara",
		"deskripsi_pekerjaan",
	}

	// List of columns to remove from pendaftaran_nikahs table
	pendaftaranNikahColumns := []string{
		"nomor_dispensasi",
		"status_bimbingan",
	}

	// Remove columns from calon_pasangans table
	fmt.Println("📋 Checking calon_pasangans table...")
	for _, col := range calonPasanganColumns {
		if columnExists(db, "calon_pasangans", col) {
			if err := dropColumn(db, "calon_pasangans", col); err != nil {
				log.Printf("⚠️  Warning: Failed to drop column %s from calon_pasangans: %v", col, err)
			} else {
				fmt.Printf("   ✅ Dropped column: calon_pasangans.%s\n", col)
			}
		} else {
			fmt.Printf("   ⏭️  Column calon_pasangans.%s does not exist, skipping\n", col)
		}
	}

	fmt.Println()
	fmt.Println("📋 Checking pendaftaran_nikahs table...")
	for _, col := range pendaftaranNikahColumns {
		if columnExists(db, "pendaftaran_nikahs", col) {
			if err := dropColumn(db, "pendaftaran_nikahs", col); err != nil {
				log.Printf("⚠️  Warning: Failed to drop column %s from pendaftaran_nikahs: %v", col, err)
			} else {
				fmt.Printf("   ✅ Dropped column: pendaftaran_nikahs.%s\n", col)
			}
		} else {
			fmt.Printf("   ⏭️  Column pendaftaran_nikahs.%s does not exist, skipping\n", col)
		}
	}

	fmt.Println()
	fmt.Println("✅ Cleanup completed!")
}

// columnExists checks if a column exists in a table
func columnExists(db *gorm.DB, tableName, columnName string) bool {
	var count int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = ? 
		AND COLUMN_NAME = ?
	`, tableName, columnName).Scan(&count)
	return count > 0
}

// dropColumn drops a column from a table
func dropColumn(db *gorm.DB, tableName, columnName string) error {
	sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, columnName)
	return db.Exec(sql).Error
}


