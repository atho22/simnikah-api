package testhelper

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"simnikah/internal/handlers/auth"
	"simnikah/internal/handlers/catin"
	"simnikah/internal/handlers/dashboard"
	"simnikah/internal/handlers/kepala_kua"
	"simnikah/internal/handlers/notification"
	"simnikah/internal/handlers/penghulu"
	"simnikah/internal/handlers/staff"
	"simnikah/internal/middleware"
	structs "simnikah/internal/models"
	"simnikah/pkg/crypto"
	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds the test database instance
var TestDB *gorm.DB

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	os.Setenv("JWT_KEY", "test-secret-key-for-unit-testing-only")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	if err := db.AutoMigrate(
		&structs.Users{},
		&structs.StaffKUA{},
		&structs.Penghulu{},
		&structs.PendaftaranNikah{},
		&structs.Notifikasi{},
	); err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	TestDB = db
	return db
}

// SetupRouter creates a Gin router wired to the test DB (no rate limiting)
func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	authHandler := &auth.InDB{DB: db}
	catinHandler := &catin.InDB{DB: db}
	staffHandler := &staff.InDB{DB: db}
	penghuluHandler := &penghulu.InDB{DB: db}
	kepalaKuaHandler := &kepala_kua.InDB{DB: db}
	notificationHandler := &notification.InDB{DB: db}
	dashboardHandler := &dashboard.InDB{DB: db}

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Auth
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)

	// Simnikah
	s := r.Group("/simnikah")
	{
		// Stage 1: Catin
		s.POST("/check-schedule", middleware.AuthMiddleware(), catinHandler.CheckScheduleAvailability)
		s.POST("/pendaftaran", middleware.AuthMiddleware(), catinHandler.CreateRegistration)

		// Stage 2: FC Recommendation
		s.GET("/kepala-kua/forward-chaining/recommendation/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.RecommendPenghuluWithForwardChaining)
		s.GET("/kepala-kua/forward-chaining/evaluation/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.GetDetailedEvaluationReport)
		s.GET("/kepala-kua/forward-chaining/config", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.GetForwardChainingConfig)

		// Stage 3: Assignment
		s.POST("/kepala-kua/forward-chaining/assign/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.AssignPenghuluWithApproval)

		// Stage 4: Penghulu
		s.GET("/penghulu/jadwal-penugasan", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("penghulu", "kepala_kua"), penghuluHandler.GetJadwalPenugasan)

		// Kepala KUA management
		s.GET("/kepala-kua/available-penghulu", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.ListAvailableOfficers)
		s.GET("/kepala-kua/penghulu-tersedia", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.GetPenghuluScheduleForAssignment)

		// Staff
		s.GET("/staff", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), staffHandler.ListStaff)
		s.PUT("/staff/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), staffHandler.UpdateStaff)
		s.POST("/staff/pendaftaran", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), staffHandler.CreateRegistrationForUser)
		s.PUT("/pendaftaran/:id/update-status", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "penghulu", "kepala_kua"), staffHandler.UpdateRegistrationStatus)

		// Location
		s.POST("/location/geocode", middleware.AuthMiddleware(), catinHandler.GetCoordinatesFromAddressEndpoint)
		s.POST("/location/reverse-geocode", middleware.AuthMiddleware(), catinHandler.GetAddressFromCoordinates)
		s.GET("/location/search", middleware.AuthMiddleware(), catinHandler.SearchAddress)
		s.PUT("/pendaftaran/:id/location", middleware.AuthMiddleware(), catinHandler.UpdateWeddingLocationWithCoordinates)
		s.GET("/pendaftaran/:id/location", middleware.AuthMiddleware(), catinHandler.GetWeddingLocationDetail)

		// Notifications
		s.POST("/notifikasi", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), notificationHandler.CreateNotification)
		s.GET("/notifikasi/user/me", middleware.AuthMiddleware(), notificationHandler.GetUserNotifications)
		s.GET("/notifikasi/:id", middleware.AuthMiddleware(), notificationHandler.GetNotificationByID)
		s.PUT("/notifikasi/:id/status", middleware.AuthMiddleware(), notificationHandler.UpdateNotificationStatus)
		s.PUT("/notifikasi/mark-all-read", middleware.AuthMiddleware(), notificationHandler.MarkAllAsRead)
		s.DELETE("/notifikasi/:id", middleware.AuthMiddleware(), notificationHandler.DeleteNotification)
		s.GET("/notifikasi/stats", middleware.AuthMiddleware(), notificationHandler.GetNotificationStats)
		s.POST("/notifikasi/send-to-role", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), notificationHandler.SendNotificationToRole)

		// Dashboard
		s.GET("/dashboard/kepala-kua", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), dashboardHandler.GetKepalaKUADashboard)
		s.GET("/dashboard/staff", middleware.AuthMiddleware(), middleware.RoleMiddleware("staff"), dashboardHandler.GetStaffDashboard)
		s.GET("/dashboard/statistik-pernikahan", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), dashboardHandler.GetMarriageStatistics)
		s.GET("/dashboard/penghulu-performance", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), dashboardHandler.GetPenghuluPerformance)
		s.GET("/dashboard/peak-hours", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), dashboardHandler.GetPeakHoursAnalysis)
	}

	return r
}

// ==================== SEED HELPERS ====================

func uid() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// CreateTestUser creates a user with hashed password and returns it
func CreateTestUser(db *gorm.DB, username, email, role string) structs.Users {
	hashed, _ := crypto.HashPassword("password123")
	user := structs.Users{
		User_id:    fmt.Sprintf("USR-%s-%s", username, uid()),
		Username:   username,
		Email:      email,
		Password:   hashed,
		Role:       role,
		Status:     structs.UserStatusAktif,
		Nama:       "Test " + username,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	db.Create(&user)
	return user
}

// CreateTestPenghulu creates a penghulu record linked to a user
func CreateTestPenghulu(db *gorm.DB, userID string) structs.Penghulu {
	p := structs.Penghulu{
		User_id:      userID,
		NIP:          fmt.Sprintf("NIP-%s", userID[:8]),
		Nama_lengkap: "Penghulu " + userID,
		No_hp:        "081234567890",
		Email:        fmt.Sprintf("%s@test.com", userID),
		Status:       structs.PenghuluStatusAktif,
		Jumlah_nikah: 5,
		Rating:       4.5,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}
	db.Create(&p)
	return p
}

// CreateTestNotification creates a notification for a user
func CreateTestNotification(db *gorm.DB, userID, tipe string) structs.Notifikasi {
	n := structs.Notifikasi{
		User_id:     userID,
		Judul:       "Test Notification",
		Pesan:       "This is a test notification",
		Tipe:        tipe,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}
	db.Create(&n)
	return n
}

// CreateTestPendaftaran creates a test pendaftaran record
func CreateTestPendaftaran(db *gorm.DB, pendaftarID string, tanggal time.Time) structs.PendaftaranNikah {
	p := structs.PendaftaranNikah{
		Nama_suami:         "Test Suami",
		Umur_suami:         28,
		Nama_istri:         "Test Istri",
		Umur_istri:         25,
		Tanggal_nikah:      tanggal,
		Waktu_nikah:        "09:00",
		Tempat_nikah:       structs.TempatNikahDiKUA,
		Status_pendaftaran: structs.StatusPendaftaranMenungguPenugasan,
		Pendaftar_id:       pendaftarID,
		Created_at:         time.Now(),
		Updated_at:         time.Now(),
	}
	db.Create(&p)
	return p
}

// ==================== JWT HELPERS ====================

// GenerateTestToken creates a valid JWT token for testing
func GenerateTestToken(userID, email, role, nama string) string {
	claims := utils.TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Nama:   nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, _ := utils.GenerateToken(claims)
	return token
}

// AuthHeader wraps a token in the Authorization header format
func AuthHeader(token string) string {
	return "Bearer " + token
}

// CleanupDB removes all data from all tables (use between tests)
func CleanupDB(db *gorm.DB) {
	db.Exec("DELETE FROM notifikasis")
	db.Exec("DELETE FROM pendaftaran_nikahs")
	db.Exec("DELETE FROM penghulus")
	db.Exec("DELETE FROM staff_k_u_a_s")
	db.Exec("DELETE FROM users")
}
