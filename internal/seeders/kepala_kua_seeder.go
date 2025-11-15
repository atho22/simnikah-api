package seeders

import (
	"fmt"
	"log"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/crypto"

	"gorm.io/gorm"
)

// SeedKepalaKUA creates initial kepala KUA user if not exists
func SeedKepalaKUA(db *gorm.DB) error {
	log.Println("🌱 Seeding Kepala KUA user...")

	// Default kepala KUA credentials
	defaultUsername := "kepalakua"
	defaultEmail := "kepalakua@kua.go.id"
	defaultPassword := "kepalakua123" // ⚠️ CHANGE THIS IN PRODUCTION!
	defaultNama := "Kepala KUA Banjarmasin Utara"
	defaultNIP := "197001011990011000"

	// Check if kepala KUA already exists
	var existingUser structs.Users
	err := db.Where("username = ? OR email = ?", defaultUsername, defaultEmail).First(&existingUser).Error

	if err == nil {
		// User already exists
		log.Printf("✅ Kepala KUA user already exists (ID: %s, Username: %s)", existingUser.User_id, existingUser.Username)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		// Database error
		return fmt.Errorf("error checking existing kepala KUA: %v", err)
	}

	// Generate user_id
	userID := "KKUA" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password
	hashedPassword, err := crypto.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	// Create user account
	user := structs.Users{
		User_id:    userID,
		Username:   defaultUsername,
		Email:      defaultEmail,
		Password:   hashedPassword,
		Role:       structs.UserRoleKepalaKUA,
		Status:     structs.UserStatusAktif,
		Nama:       defaultNama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating kepala KUA user: %v", err)
	}

	// Create staff profile for kepala KUA
	staff := structs.StaffKUA{
		User_id:      userID,
		NIP:          defaultNIP,
		Nama_lengkap: defaultNama,
		Jabatan:      structs.StaffJabatanKepalaKUA,
		Bagian:       "Kepala KUA",
		Status:       structs.StaffStatusAktif,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := db.Create(&staff).Error; err != nil {
		// If staff creation fails, delete the user
		db.Delete(&user)
		return fmt.Errorf("error creating kepala KUA staff profile: %v", err)
	}

	log.Printf("✅ Kepala KUA user created successfully!")
	log.Printf("   User ID: %s", userID)
	log.Printf("   Username: %s", defaultUsername)
	log.Printf("   Email: %s", defaultEmail)
	log.Printf("   Password: %s (⚠️  CHANGE THIS IN PRODUCTION!)", defaultPassword)
	log.Printf("   Role: %s", structs.UserRoleKepalaKUA)
	log.Println("")
	log.Println("⚠️  IMPORTANT: Change the default password after first login!")
	log.Println("")

	return nil
}

// SeedKepalaKUAWithCustomCredentials creates kepala KUA with custom credentials
func SeedKepalaKUAWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip string) error {
	log.Println("🌱 Seeding Kepala KUA user with custom credentials...")

	// Validate inputs
	if username == "" {
		username = "kepalakua"
	}
	if email == "" {
		email = "kepalakua@kua.go.id"
	}
	if password == "" {
		password = "kepalakua123"
	}
	if nama == "" {
		nama = "Kepala KUA Banjarmasin Utara"
	}
	if nip == "" {
		nip = "197001011990011000"
	}

	// Check if kepala KUA already exists
	var existingUser structs.Users
	err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error

	if err == nil {
		log.Printf("✅ Kepala KUA user already exists (ID: %s, Username: %s)", existingUser.User_id, existingUser.Username)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking existing kepala KUA: %v", err)
	}

	// Generate user_id
	userID := "KKUA" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	// Create user account
	user := structs.Users{
		User_id:    userID,
		Username:   username,
		Email:      email,
		Password:   hashedPassword,
		Role:       structs.UserRoleKepalaKUA,
		Status:     structs.UserStatusAktif,
		Nama:       nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating kepala KUA user: %v", err)
	}

	// Create staff profile for kepala KUA
	staff := structs.StaffKUA{
		User_id:      userID,
		NIP:          nip,
		Nama_lengkap: nama,
		Jabatan:      structs.StaffJabatanKepalaKUA,
		Bagian:       "Kepala KUA",
		Status:       structs.StaffStatusAktif,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := db.Create(&staff).Error; err != nil {
		// If staff creation fails, delete the user
		db.Delete(&user)
		return fmt.Errorf("error creating kepala KUA staff profile: %v", err)
	}

	log.Printf("✅ Kepala KUA user created successfully!")
	log.Printf("   User ID: %s", userID)
	log.Printf("   Username: %s", username)
	log.Printf("   Email: %s", email)
	log.Printf("   Nama: %s", nama)
	log.Printf("   NIP: %s", nip)
	log.Printf("   Role: %s", structs.UserRoleKepalaKUA)
	log.Println("")

	return nil
}

