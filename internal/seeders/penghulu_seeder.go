package seeders

import (
	"fmt"
	"log"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/crypto"

	"gorm.io/gorm"
)

// SeedPenghulu creates initial penghulu user if not exists
func SeedPenghulu(db *gorm.DB) error {
	log.Println("🌱 Seeding Penghulu user...")

	// Default penghulu credentials
	defaultUsername := "penghulu001"
	defaultEmail := "penghulu@kua.go.id"
	defaultPassword := "penghulu123" // ⚠️ CHANGE THIS IN PRODUCTION!
	defaultNama := "Ustadz Haji Ahmad Wijaya"
	defaultNIP := "197001011990011002"

	// Check if penghulu already exists
	var existingUser structs.Users
	err := db.Where("username = ? OR email = ?", defaultUsername, defaultEmail).First(&existingUser).Error

	if err == nil {
		// User already exists
		log.Printf("✅ Penghulu user already exists (ID: %s, Username: %s)", existingUser.User_id, existingUser.Username)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		// Database error
		return fmt.Errorf("error checking existing penghulu: %v", err)
	}

	// Generate user_id
	userID := "PNG" + fmt.Sprintf("%d", time.Now().Unix())

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
		Role:       structs.UserRolePenghulu,
		Status:     structs.UserStatusAktif,
		Nama:       defaultNama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating penghulu user: %v", err)
	}

	// Create penghulu profile
	penghulu := structs.Penghulu{
		User_id:      userID,
		NIP:          defaultNIP,
		Nama_lengkap: defaultNama,
		Status:       structs.PenghuluStatusAktif,
		Jumlah_nikah: 0,
		Rating:       0.0,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := db.Create(&penghulu).Error; err != nil {
		// If penghulu creation fails, delete the user
		db.Delete(&user)
		return fmt.Errorf("error creating penghulu profile: %v", err)
	}

	log.Printf("✅ Penghulu user created successfully!")
	log.Printf("   User ID: %s", userID)
	log.Printf("   Username: %s", defaultUsername)
	log.Printf("   Email: %s", defaultEmail)
	log.Printf("   Password: %s (⚠️  CHANGE THIS IN PRODUCTION!)", defaultPassword)
	log.Printf("   Role: %s", structs.UserRolePenghulu)
	log.Println("")
	log.Println("⚠️  IMPORTANT: Change the default password after first login!")
	log.Println("")

	return nil
}

// SeedPenghuluWithCustomCredentials creates penghulu with custom credentials
func SeedPenghuluWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip string) error {
	log.Println("🌱 Seeding Penghulu user with custom credentials...")

	// Validate inputs
	if username == "" {
		username = "penghulu001"
	}
	if email == "" {
		email = "penghulu@kua.go.id"
	}
	if password == "" {
		password = "penghulu123"
	}
	if nama == "" {
		nama = "Ustadz Haji Ahmad Wijaya"
	}
	if nip == "" {
		nip = "197001011990011002"
	}

	// Check if penghulu already exists
	var existingUser structs.Users
	err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error

	if err == nil {
		log.Printf("✅ Penghulu user already exists (ID: %s, Username: %s)", existingUser.User_id, existingUser.Username)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking existing penghulu: %v", err)
	}

	// Check if NIP already exists
	var existingPenghulu structs.Penghulu
	if err := db.Where("nip = ?", nip).First(&existingPenghulu).Error; err == nil {
		return fmt.Errorf("NIP %s already exists", nip)
	}

	// Generate user_id
	userID := "PNG" + fmt.Sprintf("%d", time.Now().Unix())

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
		Role:       structs.UserRolePenghulu,
		Status:     structs.UserStatusAktif,
		Nama:       nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating penghulu user: %v", err)
	}

	// Create penghulu profile
	penghulu := structs.Penghulu{
		User_id:      userID,
		NIP:          nip,
		Nama_lengkap: nama,
		Status:       structs.PenghuluStatusAktif,
		Jumlah_nikah: 0,
		Rating:       0.0,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := db.Create(&penghulu).Error; err != nil {
		// If penghulu creation fails, delete the user
		db.Delete(&user)
		return fmt.Errorf("error creating penghulu profile: %v", err)
	}

	log.Printf("✅ Penghulu user created successfully!")
	log.Printf("   User ID: %s", userID)
	log.Printf("   Username: %s", username)
	log.Printf("   Email: %s", email)
	log.Printf("   Nama: %s", nama)
	log.Printf("   NIP: %s", nip)
	log.Printf("   Role: %s", structs.UserRolePenghulu)
	log.Println("")

	return nil
}

