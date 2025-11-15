package main

import (
	"flag"
	"log"
	"os"

	"simnikah/config"
	"simnikah/internal/seeders"

	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	seedType := flag.String("type", "all", "Type of seeder to run: all, kepala_kua, staff, penghulu")
	flag.Parse()

	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Koneksi ke database gagal:", err)
	}

	log.Println("🌱 Starting database seeding...")
	log.Println("")

	switch *seedType {
	case "kepala_kua":
		seedKepalaKUA(db)
	case "staff":
		seedStaff(db)
	case "penghulu":
		seedPenghulu(db)
	case "all":
		seedKepalaKUA(db)
		seedStaff(db)
		seedPenghulu(db)
	default:
		log.Fatalf("Invalid seed type: %s. Valid types: all, kepala_kua, staff, penghulu", *seedType)
	}

	log.Println("")
	log.Println("✅ All seeding completed successfully!")
}

func seedKepalaKUA(db *gorm.DB) {
	// Get credentials from environment variables (optional)
	username := os.Getenv("SEEDER_KEPALA_KUA_USERNAME")
	email := os.Getenv("SEEDER_KEPALA_KUA_EMAIL")
	password := os.Getenv("SEEDER_KEPALA_KUA_PASSWORD")
	nama := os.Getenv("SEEDER_KEPALA_KUA_NAMA")
	nip := os.Getenv("SEEDER_KEPALA_KUA_NIP")

	// If custom credentials provided, use them
	if username != "" || email != "" || password != "" {
		log.Println("🌱 Seeding Kepala KUA with custom credentials...")
		if err := seeders.SeedKepalaKUAWithCustomCredentials(db, username, email, password, nama, nip); err != nil {
			log.Printf("⚠️  Warning: Failed to seed kepala KUA: %v", err)
		}
	} else {
		// Use default credentials
		log.Println("🌱 Seeding Kepala KUA with default credentials...")
		if err := seeders.SeedKepalaKUA(db); err != nil {
			log.Printf("⚠️  Warning: Failed to seed kepala KUA: %v", err)
		}
	}
}

func seedStaff(db *gorm.DB) {
	// Get credentials from environment variables (optional)
	username := os.Getenv("SEEDER_STAFF_USERNAME")
	email := os.Getenv("SEEDER_STAFF_EMAIL")
	password := os.Getenv("SEEDER_STAFF_PASSWORD")
	nama := os.Getenv("SEEDER_STAFF_NAMA")
	nip := os.Getenv("SEEDER_STAFF_NIP")
	jabatan := os.Getenv("SEEDER_STAFF_JABATAN")
	bagian := os.Getenv("SEEDER_STAFF_BAGIAN")

	// If custom credentials provided, use them
	if username != "" || email != "" || password != "" {
		log.Println("🌱 Seeding Staff with custom credentials...")
		if err := seeders.SeedStaffWithCustomCredentials(db, username, email, password, nama, nip, jabatan, bagian); err != nil {
			log.Printf("⚠️  Warning: Failed to seed staff: %v", err)
		}
	} else {
		// Use default credentials
		log.Println("🌱 Seeding Staff with default credentials...")
		if err := seeders.SeedStaff(db); err != nil {
			log.Printf("⚠️  Warning: Failed to seed staff: %v", err)
		}
	}
}

func seedPenghulu(db *gorm.DB) {
	// Get credentials from environment variables (optional)
	username := os.Getenv("SEEDER_PENGHULU_USERNAME")
	email := os.Getenv("SEEDER_PENGHULU_EMAIL")
	password := os.Getenv("SEEDER_PENGHULU_PASSWORD")
	nama := os.Getenv("SEEDER_PENGHULU_NAMA")
	nip := os.Getenv("SEEDER_PENGHULU_NIP")

	// If custom credentials provided, use them
	if username != "" || email != "" || password != "" {
		log.Println("🌱 Seeding Penghulu with custom credentials...")
		if err := seeders.SeedPenghuluWithCustomCredentials(db, username, email, password, nama, nip); err != nil {
			log.Printf("⚠️  Warning: Failed to seed penghulu: %v", err)
		}
	} else {
		// Use default credentials
		log.Println("🌱 Seeding Penghulu with default credentials...")
		if err := seeders.SeedPenghulu(db); err != nil {
			log.Printf("⚠️  Warning: Failed to seed penghulu: %v", err)
		}
	}
}

