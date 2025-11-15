package seeders

import (
	"fmt"
	"log"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/crypto"

	"gorm.io/gorm"
)

// SeedStaff creates initial staff KUA user if not exists
func SeedStaff(db *gorm.DB) error {
	log.Println("🌱 Seeding Staff KUA user...")

	// Default staff credentials
	defaultUsername := "staff001"
	defaultEmail := "staff@kua.go.id"
	defaultPassword := "staff123" // ⚠️ CHANGE THIS IN PRODUCTION!
	defaultNama := "Staff KUA Banjarmasin Utara"
	defaultNIP := "197001011990011001"
	defaultJabatan := structs.StaffJabatanStaff
	defaultBagian := "Verifikasi"

	// Check if staff already exists
	var existingUser structs.Users
	err := db.Where("username = ? OR email = ?", defaultUsername, defaultEmail).First(&existingUser).Error

	if err == nil {
		// User already exists
		log.Printf("✅ Staff KUA user already exists (ID: %s, Username: %s)", existingUser.User_id, existingUser.Username)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		// Database error
		return fmt.Errorf("error checking existing staff: %v", err)
	}

	// Generate user_id
	userID := "STF" + fmt.Sprintf("%d", time.Now().Unix())

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
		Role:       structs.UserRoleStaff,
		Status:     structs.UserStatusAktif,
		Nama:       defaultNama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating staff user: %v", err)
	}

	// Create staff profile
	staff := structs.StaffKUA{
		User_id:      userID,
		NIP:          defaultNIP,
		Nama_lengkap: defaultNama,
		Jabatan:      defaultJabatan,
		Bagian:       defaultBagian,
		Status:       structs.StaffStatusAktif,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := db.Create(&staff).Error; err != nil {
		// If staff creation fails, delete the user
		db.Delete(&user)
		return fmt.Errorf("error creating staff profile: %v", err)
	}

	log.Printf("✅ Staff KUA user created successfully!")
	log.Printf("   User ID: %s", userID)
	log.Printf("   Username: %s", defaultUsername)
	log.Printf("   Email: %s", defaultEmail)
	log.Printf("   Password: %s (⚠️  CHANGE THIS IN PRODUCTION!)", defaultPassword)
	log.Printf("   Role: %s", structs.UserRoleStaff)
	log.Printf("   Jabatan: %s", defaultJabatan)
	log.Println("")
	log.Println("⚠️  IMPORTANT: Change the default password after first login!")
	log.Println("")

	return nil
}

// SeedStaffWithCustomCredentials creates staff with custom credentials
func SeedStaffWithCustomCredentials(db *gorm.DB, username, email, password, nama, nip, jabatan, bagian string) error {
	log.Println("🌱 Seeding Staff KUA user with custom credentials...")

	// Validate inputs
	if username == "" {
		username = "staff001"
	}
	if email == "" {
		email = "staff@kua.go.id"
	}
	if password == "" {
		password = "staff123"
	}
	if nama == "" {
		nama = "Staff KUA Banjarmasin Utara"
	}
	if nip == "" {
		nip = "197001011990011001"
	}
	if jabatan == "" {
		jabatan = structs.StaffJabatanStaff
	}
	if bagian == "" {
		bagian = "Verifikasi"
	}

	// Validate jabatan
	validJabatan := map[string]bool{
		structs.StaffJabatanStaff:     true,
		structs.StaffJabatanPenghulu:  true,
		structs.StaffJabatanKepalaKUA: true,
	}
	if !validJabatan[jabatan] {
		return fmt.Errorf("invalid jabatan: %s. Valid values: Staff, Penghulu, Kepala KUA", jabatan)
	}

	// Check if staff already exists
	var existingUser structs.Users
	err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error

	if err == nil {
		log.Printf("✅ Staff KUA user already exists (ID: %s, Username: %s)", existingUser.User_id, existingUser.Username)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking existing staff: %v", err)
	}

	// Check if NIP already exists
	var existingStaff structs.StaffKUA
	if err := db.Where("nip = ?", nip).First(&existingStaff).Error; err == nil {
		return fmt.Errorf("NIP %s already exists", nip)
	}

	// Generate user_id
	userID := "STF" + fmt.Sprintf("%d", time.Now().Unix())

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
		Role:       structs.UserRoleStaff,
		Status:     structs.UserStatusAktif,
		Nama:       nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating staff user: %v", err)
	}

	// Create staff profile
	staff := structs.StaffKUA{
		User_id:      userID,
		NIP:          nip,
		Nama_lengkap: nama,
		Jabatan:      jabatan,
		Bagian:       bagian,
		Status:       structs.StaffStatusAktif,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := db.Create(&staff).Error; err != nil {
		// If staff creation fails, delete the user
		db.Delete(&user)
		return fmt.Errorf("error creating staff profile: %v", err)
	}

	log.Printf("✅ Staff KUA user created successfully!")
	log.Printf("   User ID: %s", userID)
	log.Printf("   Username: %s", username)
	log.Printf("   Email: %s", email)
	log.Printf("   Nama: %s", nama)
	log.Printf("   NIP: %s", nip)
	log.Printf("   Jabatan: %s", jabatan)
	log.Printf("   Bagian: %s", bagian)
	log.Printf("   Role: %s", structs.UserRoleStaff)
	log.Println("")

	return nil
}

