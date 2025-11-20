package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"simnikah/config"
	"simnikah/internal/seeders"

	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	seedKepalaKUA bool
	seedStaff     bool
	seedPenghulu  bool
	seedAll       bool
)

func init() {
	flag.BoolVar(&seedKepalaKUA, "kepala-kua", false, "Seed kepala KUA user")
	flag.BoolVar(&seedStaff, "staff", false, "Seed staff KUA user")
	flag.BoolVar(&seedPenghulu, "penghulu", false, "Seed penghulu user")
	flag.BoolVar(&seedAll, "all", false, "Seed all users (kepala KUA, staff, penghulu)")
	flag.Parse()
}

func main() {
	// Initialize database connection
	var err error
	db, err = config.ConnectDB()
	if err != nil {
		log.Fatal("❌ Koneksi ke database gagal:", err)
	}

	log.Println("✅ Koneksi database berhasil")
	log.Println("")

	// Check if no flags provided
	if !seedKepalaKUA && !seedStaff && !seedPenghulu && !seedAll {
		printUsage()
		os.Exit(1)
	}

	// Seed based on flags
	if seedAll || seedKepalaKUA {
		log.Println("🌱 Seeding Kepala KUA...")
		if err := seeders.SeedKepalaKUA(db); err != nil {
			log.Printf("❌ Gagal seeding kepala KUA: %v\n", err)
		} else {
			log.Println("✅ Kepala KUA seeding selesai")
		}
		log.Println("")
	}

	if seedAll || seedStaff {
		log.Println("🌱 Seeding Staff KUA...")
		if err := seeders.SeedStaff(db); err != nil {
			log.Printf("❌ Gagal seeding staff: %v\n", err)
		} else {
			log.Println("✅ Staff seeding selesai")
		}
		log.Println("")
	}

	if seedAll || seedPenghulu {
		log.Println("🌱 Seeding Penghulu...")
		if err := seeders.SeedPenghulu(db); err != nil {
			log.Printf("❌ Gagal seeding penghulu: %v\n", err)
		} else {
			log.Println("✅ Penghulu seeding selesai")
		}
		log.Println("")
	}

	log.Println("🎉 Proses seeding selesai!")
	log.Println("")
	log.Println("📝 Kredensial default:")
	log.Println("   Kepala KUA: username=kepalakua, password=kepalakua123")
	log.Println("   Staff:      username=staff001, password=staff123")
	log.Println("   Penghulu:   username=penghulu001, password=penghulu123")
	log.Println("")
	log.Println("⚠️  PENTING: Ganti password default setelah login pertama kali!")
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Penggunaan: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -kepala-kua    Seed kepala KUA user\n")
	fmt.Fprintf(os.Stderr, "  -staff         Seed staff KUA user\n")
	fmt.Fprintf(os.Stderr, "  -penghulu      Seed penghulu user\n")
	fmt.Fprintf(os.Stderr, "  -all           Seed semua user (kepala KUA, staff, penghulu)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Contoh:\n")
	fmt.Fprintf(os.Stderr, "  %s -all\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -kepala-kua -staff\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -penghulu\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
}

