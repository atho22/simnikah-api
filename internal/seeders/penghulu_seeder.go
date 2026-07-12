package seeders

import (
	"fmt"
	"log"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/crypto"

	"gorm.io/gorm"
)

// SeedPenghulu creates initial penghulu users if they do not exist (minimum 4)
func SeedPenghulu(db *gorm.DB) error {
	log.Println("🌱 Seeding default Penghulu users (minimum 4)...")

	defaultPenghulus := []struct {
		Username string
		Email    string
		Password string
		Nama     string
		NIP      string
	}{
		{
			Username: "penghulu001",
			Email:    "penghulu@kua.go.id",
			Password: "penghulu123",
			Nama:     "Ustadz Haji Ahmad Wijaya",
			NIP:      "197001011990011002",
		},
		{
			Username: "penghulu002",
			Email:    "penghulu002@kua.go.id",
			Password: "penghulu123",
			Nama:     "Ustadz Muhammad Yusuf",
			NIP:      "197503122002121003",
		},
		{
			Username: "penghulu003",
			Email:    "penghulu003@kua.go.id",
			Password: "penghulu123",
			Nama:     "Ustadz H. Abdul Rahman",
			NIP:      "198008202008011004",
		},
		{
			Username: "penghulu004",
			Email:    "penghulu004@kua.go.id",
			Password: "penghulu123",
			Nama:     "Ustadz Syarif Hidayatullah",
			NIP:      "198512052014031005",
		},
	}

	baseUnix := time.Now().Unix()

	for i, p := range defaultPenghulus {
		// Check if penghulu already exists
		var existingUser structs.Users
		err := db.Where("username = ? OR email = ?", p.Username, p.Email).First(&existingUser).Error

		if err == nil {
			// User already exists
			log.Printf("✅ Penghulu user already exists: %s (ID: %s)", p.Username, existingUser.User_id)
			continue
		}

		if err != gorm.ErrRecordNotFound {
			// Database error
			return fmt.Errorf("error checking existing penghulu: %v", err)
		}

		// Generate unique user_id
		userID := fmt.Sprintf("PNG%d%d", baseUnix, i)

		// Hash password
		hashedPassword, err := crypto.HashPassword(p.Password)
		if err != nil {
			return fmt.Errorf("error hashing password for %s: %v", p.Username, err)
		}

		// Create user account
		user := structs.Users{
			User_id:    userID,
			Username:   p.Username,
			Email:      p.Email,
			Password:   hashedPassword,
			Role:       structs.UserRolePenghulu,
			Status:     structs.UserStatusAktif,
			Nama:       p.Nama,
			Created_at: time.Now(),
			Updated_at: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("error creating user %s: %v", p.Username, err)
		}

		// Create penghulu profile
		penghuluProfile := structs.Penghulu{
			User_id:      userID,
			NIP:          p.NIP,
			Nama_lengkap: p.Nama,
			Status:       structs.PenghuluStatusAktif,
			Jumlah_nikah: 0,
			Rating:       0.0,
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		}

		if err := db.Create(&penghuluProfile).Error; err != nil {
			// Clean up user
			db.Delete(&user)
			return fmt.Errorf("error creating profile for %s: %v", p.Username, err)
		}

		log.Printf("   [+] Created Penghulu user: %s (ID: %s, Username: %s, NIP: %s)", p.Nama, userID, p.Username, p.NIP)
	}

	log.Println("✅ Default Penghulu users seeded successfully!")
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

