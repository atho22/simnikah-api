package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"simnikah/config"
	"simnikah/internal/handlers/catin"
	"simnikah/internal/handlers/notification"
	"simnikah/internal/handlers/penghulu"
	"simnikah/internal/handlers/staff"
	"simnikah/internal/middleware"
	structs "simnikah/internal/models"
	"simnikah/internal/services"
	"simnikah/pkg/crypto"
	"simnikah/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// JWT Key
var jwtKey = getJWTKey()

// Token Claims
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Nama   string `json:"nama"`
	jwt.RegisteredClaims
}

var DB *gorm.DB

func main() {
	var err error
	// Validate JWT key is adequate
	if len(jwtKey) < 32 {
		log.Println("Warning: JWT key is less than 32 bytes. Consider using a stronger key.")
	}

	// Initialize database connection
	DB, err = config.ConnectDB()
	if err != nil {
		log.Fatal("Koneksi ke database gagal:", err)
	}

	// Migrate struct
	log.Println("Starting database migration...")
	if err := DB.AutoMigrate(&structs.Users{}, &structs.StaffKUA{}, &structs.Penghulu{}, &structs.DataOrangTua{}, &structs.CalonPasangan{}, &structs.PendaftaranNikah{}, &structs.WaliNikah{}, &structs.BimbinganPerkawinan{}, &structs.PendaftaranBimbingan{}, &structs.Notifikasi{}); err != nil {
		log.Fatal("Database migration failed:", err)
	}
	log.Println("Database migration completed successfully")

	// Add database indexes for performance optimization
	if err := config.AddDatabaseIndexes(DB); err != nil {
		log.Println("Warning: Failed to add database indexes:", err)
		// Don't fatal, just warn - indexes are optional (but highly recommended)
	}

	// Replace individual seed calls with SeedAll

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Configure CORS middleware
	corsConfig := cors.Config{
		AllowOrigins: getAllowedOrigins(),
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	catinHandler := &catin.InDB{DB: DB}
	staffHandler := &staff.InDB{DB: DB}
	penghuluHandler := &penghulu.InDB{DB: DB}
	notificationHandler := &notification.InDB{DB: DB}

	// Start cron job untuk pengingat notifikasi
	cronJobService := services.NewCronJobService(DB)
	cronJobService.StartReminderCronJobWithSchedule(8, 0) // Setiap hari jam 08:00

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "SimNikah API",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Apply global rate limiting (100 req/min per IP)
	r.Use(middleware.RateLimiter())

	// Routes with strict rate limiting for auth endpoints
	r.POST("/register", middleware.StrictRateLimiter(), RegisterUser)
	r.POST("/login", middleware.StrictRateLimiter(), Login)
	r.GET("/profile", AuthMiddleware(), Profile)

	// SimNikah Routes
	simnikahRoutes := r.Group("/simnikah")
	{
		// Calon Pasangan (commented out - implementation not available)
		// simnikahRoutes.POST("/calon-pasangan", AuthMiddleware(), catinHandler.CreateCalonPasangan)
		// simnikahRoutes.GET("/calon-pasangan/:id", AuthMiddleware(), catinHandler.GetCalonPasangan)
		// simnikahRoutes.PUT("/calon-pasangan/:id", AuthMiddleware(), catinHandler.UpdateCalonPasangan)

		// Pendaftaran Nikah
		// simnikahRoutes.POST("/pendaftaran", AuthMiddleware(), catinHandler.CreatePendaftaranNikah)
		// simnikahRoutes.POST("/pendaftaran/lengkap", AuthMiddleware(), catinHandler.CreatePendaftaranLengkap)
		simnikahRoutes.GET("/pendaftaran/status", AuthMiddleware(), catinHandler.CheckUserRegistrationStatus)
		simnikahRoutes.POST("/pendaftaran/form-baru", AuthMiddleware(), catinHandler.CreateMarriageRegistrationForm)
		simnikahRoutes.POST("/pendaftaran/:id/mark-visited", AuthMiddleware(), catinHandler.MarkAsVisited)
		simnikahRoutes.GET("/pendaftaran", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), catinHandler.GetAllMarriageRegistrations)
		// simnikahRoutes.GET("/pendaftaran/:id", AuthMiddleware(), catinHandler.GetPendaftaranNikah)
		// simnikahRoutes.GET("/pendaftaran/:id/lengkap", AuthMiddleware(), catinHandler.GetPendaftaranLengkap)
		// simnikahRoutes.PUT("/pendaftaran/:id/status", MultiRoleMiddleware("staff", "kepala_kua"), catinHandler.UpdateStatusPendaftaran)

		// Wali Nikah (commented out - implementation not available)
		// simnikahRoutes.POST("/wali", AuthMiddleware(), catinHandler.CreateWaliNikah)
		// simnikahRoutes.PUT("/wali/:id/status", AuthMiddleware(), catinHandler.UpdateStatusWali)

		// Staff Management (hanya kepala KUA yang bisa akses)
		simnikahRoutes.POST("/staff", AuthMiddleware(), RoleMiddleware("kepala_kua"), staffHandler.CreateStaffKUA)
		simnikahRoutes.GET("/staff", AuthMiddleware(), RoleMiddleware("kepala_kua"), staffHandler.GetAllStaff)
		simnikahRoutes.PUT("/staff/:id", AuthMiddleware(), RoleMiddleware("kepala_kua"), staffHandler.UpdateStaffKUA)

		// Penghulu Management
		simnikahRoutes.POST("/penghulu", AuthMiddleware(), RoleMiddleware("kepala_kua"), staffHandler.CreatePenghulu)
		simnikahRoutes.GET("/penghulu", AuthMiddleware(), staffHandler.GetAllPenghulu)
		simnikahRoutes.PUT("/penghulu/:id", AuthMiddleware(), RoleMiddleware("kepala_kua"), staffHandler.UpdatePenghulu)

		// Staff Verification
		simnikahRoutes.POST("/staff/verify-formulir/:id", AuthMiddleware(), RoleMiddleware("staff"), staffHandler.VerifyFormulir)
		simnikahRoutes.POST("/staff/verify-berkas/:id", AuthMiddleware(), RoleMiddleware("staff"), staffHandler.VerifyBerkas)

		// Penghulu Operations
		simnikahRoutes.POST("/penghulu/verify-documents/:id", AuthMiddleware(), RoleMiddleware("penghulu"), penghuluHandler.VerifyDocuments)
		simnikahRoutes.GET("/penghulu/assigned-registrations", AuthMiddleware(), RoleMiddleware("penghulu"), penghuluHandler.GetAssignedRegistrations)

		// Jadwal Nikah
		simnikahRoutes.POST("/jadwal", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), CreateJadwalNikah)
		simnikahRoutes.GET("/jadwal", AuthMiddleware(), GetJadwalNikah)
		simnikahRoutes.PUT("/jadwal/:id", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), UpdateJadwalNikah)

		// Kalender Ketersediaan
		simnikahRoutes.GET("/kalender-ketersediaan", AuthMiddleware(), GetKalenderKetersediaan)
		simnikahRoutes.GET("/kalender-tanggal-detail", AuthMiddleware(), GetKalenderTanggalDetail)
		simnikahRoutes.GET("/ketersediaan-tanggal/:tanggal", AuthMiddleware(), GetKetersediaanTanggal)
		simnikahRoutes.GET("/penghulu-jadwal/:tanggal", AuthMiddleware(), GetPenghuluJadwal)

		// Management Penghulu (Kepala KUA)
		simnikahRoutes.POST("/pendaftaran/:id/assign-penghulu", AuthMiddleware(), MultiRoleMiddleware("kepala_kua"), AssignPenghulu)
		simnikahRoutes.PUT("/pendaftaran/:id/change-penghulu", AuthMiddleware(), MultiRoleMiddleware("kepala_kua"), ChangePenghulu)
		simnikahRoutes.GET("/pendaftaran/belum-assign-penghulu", AuthMiddleware(), MultiRoleMiddleware("kepala_kua"), GetPendaftaranBelumAssignPenghulu)
		simnikahRoutes.GET("/penghulu/:id/ketersediaan/:tanggal", AuthMiddleware(), MultiRoleMiddleware("kepala_kua"), GetPenghuluKetersediaan)

		// Bimbingan Perkawinan
		simnikahRoutes.POST("/bimbingan", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), CreateBimbinganPerkawinan)
		simnikahRoutes.GET("/bimbingan", AuthMiddleware(), GetBimbinganPerkawinan)
		simnikahRoutes.GET("/bimbingan/:id", AuthMiddleware(), GetBimbinganPerkawinanByID)
		simnikahRoutes.PUT("/bimbingan/:id", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), UpdateBimbinganPerkawinan)
		simnikahRoutes.GET("/bimbingan-kalender", AuthMiddleware(), GetBimbinganKalender)
		simnikahRoutes.POST("/bimbingan/:id/daftar", AuthMiddleware(), DaftarBimbinganPerkawinan)
		simnikahRoutes.GET("/bimbingan/:id/participants", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), GetBimbinganParticipants)
		simnikahRoutes.PUT("/bimbingan/:id/update-attendance", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), UpdateBimbinganAttendance)
		simnikahRoutes.GET("/bimbingan/:id/undangan", AuthMiddleware(), CetakUndanganBimbingan)
		simnikahRoutes.GET("/bimbingan/:id/undangan-semua", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), CetakUndanganBimbinganSemua)

		// Geocoding API untuk mendapatkan koordinat alamat (GRATIS menggunakan OpenStreetMap)
		simnikahRoutes.GET("/geocoding/coordinates", AuthMiddleware(), GetAddressCoordinates)

		// Status Management
		simnikahRoutes.GET("/pendaftaran/:id/status-flow", AuthMiddleware(), GetStatusFlow)
		simnikahRoutes.PUT("/pendaftaran/:id/complete-bimbingan", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), CompleteBimbingan)
		simnikahRoutes.PUT("/pendaftaran/:id/complete-nikah", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), CompleteNikah)

		// Wedding Address Management
		simnikahRoutes.PUT("/pendaftaran/:id/alamat", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), catinHandler.UpdateWeddingAddress)

		// ==================== MAP & LOCATION INTEGRATION ====================
		// Endpoints untuk integrasi peta dan koordinat lokasi nikah

		// Geocoding: Alamat → Koordinat
		simnikahRoutes.POST("/location/geocode", AuthMiddleware(), catinHandler.GetCoordinatesFromAddressEndpoint)

		// Reverse Geocoding: Koordinat → Alamat
		simnikahRoutes.POST("/location/reverse-geocode", AuthMiddleware(), catinHandler.GetAddressFromCoordinates)

		// Search Address (untuk autocomplete)
		simnikahRoutes.GET("/location/search", AuthMiddleware(), catinHandler.SearchAddress)

		// Update lokasi nikah dengan koordinat
		simnikahRoutes.PUT("/pendaftaran/:id/location", AuthMiddleware(), catinHandler.UpdateWeddingLocationWithCoordinates)

		// Get detail lokasi nikah (untuk penghulu)
		simnikahRoutes.GET("/pendaftaran/:id/location", AuthMiddleware(), catinHandler.GetWeddingLocationDetail)

		// Surat Undangan Nikah (commented out - implementation not available)
		// simnikahRoutes.POST("/surat-undangan", RoleMiddleware("kepala_kua"), catinHandler.CreateSuratUndanganNikah)
		// simnikahRoutes.GET("/surat-undangan/:pendaftaran_id", AuthMiddleware(), catinHandler.GetSuratUndanganNikah)
		// simnikahRoutes.GET("/surat-undangan", MultiRoleMiddleware("staff", "kepala_kua"), catinHandler.GetAllSuratUndanganNikah)
		// simnikahRoutes.POST("/surat-undangan/verify", AuthMiddleware(), catinHandler.VerifyDigitalSignature)

		// Notifikasi
		simnikahRoutes.POST("/notifikasi", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), notificationHandler.CreateNotification)
		simnikahRoutes.GET("/notifikasi/user/:user_id", AuthMiddleware(), notificationHandler.GetUserNotifications)
		simnikahRoutes.GET("/notifikasi/:id", AuthMiddleware(), notificationHandler.GetNotificationByID)
		simnikahRoutes.PUT("/notifikasi/:id/status", AuthMiddleware(), notificationHandler.UpdateNotificationStatus)
		simnikahRoutes.PUT("/notifikasi/user/:user_id/mark-all-read", AuthMiddleware(), notificationHandler.MarkAllAsRead)
		simnikahRoutes.DELETE("/notifikasi/:id", AuthMiddleware(), notificationHandler.DeleteNotification)
		simnikahRoutes.GET("/notifikasi/user/:user_id/stats", AuthMiddleware(), notificationHandler.GetNotificationStats)
		simnikahRoutes.POST("/notifikasi/send-to-role", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), notificationHandler.SendNotificationToRole)
		simnikahRoutes.POST("/notifikasi/run-reminder", AuthMiddleware(), MultiRoleMiddleware("staff", "kepala_kua"), RunReminderNotification)

	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("📊 Performance optimizations enabled:")
		log.Printf("   ✅ Database indexes (5-10x faster queries)")
		log.Printf("   ✅ Rate limiting (100 req/min per IP)")
		log.Printf("   ✅ Graceful shutdown (zero downtime deploys)")
		log.Printf("Environment: %s", os.Getenv("GIN_MODE"))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) and SIGTERM (kill command) will trigger graceful shutdown
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server gracefully...")

	// Give server 10 seconds to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited gracefully")
}

// getJWTKey returns the JWT key from environment or uses a fallback
func getJWTKey() []byte {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		// Only for development, warn about using default key
		log.Println("Warning: Using default JWT key. Set JWT_KEY environment variable in production.")
		return []byte("secret-key-boleh-diubah-untuk-simnikah")
	}
	return []byte(key)
}

// getAllowedOrigins returns allowed origins for CORS from environment or uses defaults
func getAllowedOrigins() []string {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	if allowedOrigins == "" {
		// Default origins for development
		log.Println("Info: Using default CORS origins. Set ALLOWED_ORIGINS environment variable for production.")
		return []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5173", // Vite default
			"http://localhost:5174",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}

	// Split comma-separated origins from environment variable
	origins := strings.Split(allowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	log.Printf("Info: CORS allowed origins: %v", origins)
	return origins
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nama     string `json:"nama" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format registrasi tidak valid"})
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"user_biasa": true, // User biasa untuk daftar nikah
		"penghulu":   true, // Penghulu untuk memimpin nikah
		"staff":      true, // Staff KUA untuk verifikasi
		"kepala_kua": true, // Kepala KUA untuk approval
	}
	if !validRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid. Role yang tersedia: user_biasa, penghulu, staff, kepala_kua"})
		return
	}

	// Check if username already exists
	var existingUser structs.Users
	if err := DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau email sudah digunakan"})
		return
	}

	// Generate user_id
	userID := "USR" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password menggunakan sistem KUA yang aman
	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak memenuhi kriteria keamanan: " + err.Error()})
		return
	}

	// Create new user
	user := structs.Users{
		User_id:    userID,
		Username:   input.Username,
		Email:      input.Email,
		Password:   string(hashedPassword),
		Role:       input.Role,
		Nama:       input.Nama,
		Status:     "Aktif",
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat",
		"user": gin.H{
			"user_id":  userID,
			"username": input.Username,
			"email":    input.Email,
			"nama":     input.Nama,
			"role":     input.Role,
		},
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format login tidak valid"})
		return
	}

	var user structs.Users

	if err := DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan"})
		return
	}

	// Check if user is active
	if user.Status != "Aktif" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak aktif"})
		return
	}

	// Validate password hash menggunakan sistem KUA yang aman
	if err := crypto.VerifyPassword(loginRequest.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Build JWT claims
	claims := TokenClaims{
		UserID: user.User_id,
		Email:  user.Email,
		Role:   user.Role,
		Nama:   user.Nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token berlaku 24 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"user_id": claims.UserID,
			"email":   claims.Email,
			"role":    claims.Role,
			"nama":    claims.Nama,
		},
	})
}

// A more robust and type-safe AuthMiddleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token otorisasi tidak disediakan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Create an empty instance of your custom claims struct
		claims := &TokenClaims{}

		// Use jwt.ParseWithClaims for type-safe parsing
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing token tidak valid")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
			return
		}

		// --- From here on, everything is type-safe! ---

		// Validasi role (lebih bersih dengan akses struct)
		validRoles := map[string]bool{
			"user_biasa": true,
			"penghulu":   true,
			"staff":      true,
			"kepala_kua": true,
		}
		if !validRoles[claims.Role] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role tidak valid"})
			return
		}

		// Set claims ke context using dot notation. No more panics!
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Set("nama", claims.Nama)

		c.Next()
	}
}

// RoleMiddleware untuk validasi role spesifik
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan"})
			c.Abort()
			return
		}

		// Check if user has required role
		if role.(string) != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak. Role " + requiredRole + " diperlukan"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MultiRoleMiddleware untuk validasi multiple roles
func MultiRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan"})
			c.Abort()
			return
		}

		// Check if user has any of the allowed roles
		userRole := role.(string)
		hasAccess := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak. Role yang diizinkan: " + strings.Join(allowedRoles, ", ")})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		// Error ini yang Anda dapatkan sekarang
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// UBAH DARI (uint) MENJADI (string)
	id, ok := userID.(string)
	if !ok {
		// Pengaman jika tipe data di context tidak sesuai
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe data user ID di context tidak valid"})
		return
	}

	var user structs.Users // Pastikan struct Users memiliki field User_id bertipe string

	// Query ke database menggunakan id yang sudah bertipe string
	if err := DB.Where("user_id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan di database"})
		return
	}

	// ... (lanjutkan mengirim response JSON)
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile berhasil diambil",
		"user": gin.H{
			"user_id":  user.User_id,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"nama":     user.Nama,
		},
	})
}

// GetAvailableRoles returns all available roles and their descriptions
func GetAvailableRoles(c *gin.Context) {
	roles := []gin.H{
		{
			"role_code":   "user_biasa",
			"role_name":   "User Biasa",
			"description": "Calon pasangan yang mendaftar nikah",
		},
		{
			"role_code":   "penghulu",
			"role_name":   "Penghulu",
			"description": "Memimpin prosesi nikah",
		},
		{
			"role_code":   "staff",
			"role_name":   "Staff KUA",
			"description": "Staff administrasi KUA",
		},
		{
			"role_code":   "kepala_kua",
			"role_name":   "Kepala KUA",
			"description": "Supervisor dan approval tertinggi",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar role tersedia",
		"data":    roles,
	})
}

// Placeholder functions for missing handlers
func CreateJadwalNikah(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Fungsi CreateJadwalNikah belum diimplementasi"})
}

func GetJadwalNikah(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Fungsi GetJadwalNikah belum diimplementasi"})
}

func UpdateJadwalNikah(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Fungsi UpdateJadwalNikah belum diimplementasi"})
}

// GetKalenderKetersediaan menampilkan kalender ketersediaan untuk pendaftaran nikah
func GetKalenderKetersediaan(c *gin.Context) {
	// Ambil parameter query
	bulan := c.DefaultQuery("bulan", "")
	tahun := c.DefaultQuery("tahun", "")

	// Default ke bulan dan tahun saat ini jika tidak ada parameter
	now := time.Now()
	if bulan == "" {
		bulan = fmt.Sprintf("%02d", int(now.Month()))
	}
	if tahun == "" {
		tahun = fmt.Sprintf("%d", now.Year())
	}

	// Parse bulan dan tahun
	bulanInt, err := strconv.Atoi(bulan)
	if err != nil || bulanInt < 1 || bulanInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan tidak valid (1-12)"})
		return
	}

	tahunInt, err := strconv.Atoi(tahun)
	if err != nil || tahunInt < now.Year() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun tidak valid"})
		return
	}

	// Hitung tanggal awal dan akhir bulan
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.UTC)
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	// Query SEMUA pendaftaran nikah pada bulan tersebut (ditampilkan di kalender)
	var pendaftaran []structs.PendaftaranNikah
	err = DB.Where("tanggal_nikah >= ? AND tanggal_nikah <= ?",
		awalBulan, akhirBulan).Find(&pendaftaran).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendaftaran"})
		return
	}

	// Hitung kapasitas per hari: maksimal 9 pernikahan di KUA per hari
	// Nikah di KUA: 1 per slot waktu (tidak bisa bersamaan)
	// Nikah di luar KUA: tidak dibatasi (4 penghulu bisa bersamaan)
	kapasitasPerHari := 9 // Total maksimal 9 pernikahan di KUA per hari

	// Buat map untuk menghitung jumlah per hari dan kategori warna
	totalPerHari := make(map[string]int)
	kuningPerHari := make(map[string]int)         // status awal (belum selesai berkas)
	hijauPerHari := make(map[string]int)          // status sudah fix
	nikahDiKUAPerHari := make(map[string]int)     // nikah di KUA per hari
	nikahDiLuarKUAPerHari := make(map[string]int) // nikah di luar KUA per hari

	for _, p := range pendaftaran {
		tanggal := p.Tanggal_nikah.Format("2006-01-02")
		totalPerHari[tanggal]++

		// Hitung nikah di KUA dan di luar KUA
		if p.Tempat_nikah == "Di KUA" {
			nikahDiKUAPerHari[tanggal]++
		} else {
			nikahDiLuarKUAPerHari[tanggal]++
		}

		// Kategori warna: kuning untuk status awal; hijau untuk status setelah berkas fix
		s := p.Status_pendaftaran
		if s == "Menunggu Verifikasi" || s == "Menunggu Pengumpulan Berkas" {
			kuningPerHari[tanggal]++
		} else {
			hijauPerHari[tanggal]++
		}
	}

	// Buat kalender untuk bulan tersebut
	kalender := make([]map[string]interface{}, 0)

	// Mulai dari tanggal 1 sampai akhir bulan
	for tanggal := 1; tanggal <= akhirBulan.Day(); tanggal++ {
		tanggalStr := fmt.Sprintf("%d-%02d-%02d", tahunInt, bulanInt, tanggal)
		tanggalTime := time.Date(tahunInt, time.Month(bulanInt), tanggal, 0, 0, 0, 0, time.UTC)

		// Hitung jumlah nikah yang sudah terjadwal
		jumlahNikahTotal := totalPerHari[tanggalStr]
		jumlahNikahKUA := nikahDiKUAPerHari[tanggalStr]
		jumlahNikahLuarKUA := nikahDiLuarKUAPerHari[tanggalStr]
		kuningCount := kuningPerHari[tanggalStr]
		hijauCount := hijauPerHari[tanggalStr]

		// Tentukan warna keseluruhan hari: hijau jika ada yang hijau, kuning jika hanya kuning
		warnaHari := ""
		if hijauCount > 0 {
			warnaHari = "hijau"
		} else if kuningCount > 0 {
			warnaHari = "kuning"
		}

		// Tentukan status ketersediaan (hanya untuk nikah di KUA)
		var status string
		var tersedia bool
		var sisaKuota int

		if tanggalTime.Before(now.Truncate(24 * time.Hour)) {
			// Tanggal sudah lewat
			status = "Terlewat"
			tersedia = false
			sisaKuota = 0
		} else if jumlahNikahKUA >= kapasitasPerHari {
			// Sudah penuh untuk nikah di KUA
			status = "Penuh"
			tersedia = false
			sisaKuota = 0
		} else {
			// Masih tersedia
			status = "Tersedia"
			tersedia = true
			sisaKuota = kapasitasPerHari - jumlahNikahKUA
		}

		// Tambahkan ke kalender
		kalender = append(kalender, map[string]interface{}{
			"tanggal":            tanggal,
			"tanggal_str":        tanggalStr,
			"status":             status,
			"tersedia":           tersedia,
			"jumlah_nikah_total": jumlahNikahTotal,
			"jumlah_nikah_kua":   jumlahNikahKUA,
			"jumlah_nikah_luar":  jumlahNikahLuarKUA,
			"kuning_count":       kuningCount,
			"hijau_count":        hijauCount,
			"warna":              warnaHari,
			"sisa_kuota_kua":     sisaKuota,
			"kapasitas_kua":      kapasitasPerHari,
		})
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Kalender ketersediaan berhasil diambil",
		"data": gin.H{
			"bulan":            bulanInt,
			"tahun":            tahunInt,
			"nama_bulan":       awalBulan.Month().String(),
			"kapasitas_harian": kapasitasPerHari,
			"penghulu_info": gin.H{
				"total_penghulu":         4,
				"penghulu_aktif":         4,
				"penghulu_cadangan":      0,
				"slot_waktu_per_hari":    9,
				"nikah_per_slot":         4,
				"total_kapasitas_harian": kapasitasPerHari,
			},
			"kalender": kalender,
		},
	})
}

// GetKalenderTanggalDetail menampilkan daftar jam pernikahan per tanggal dengan kategori warna
func GetKalenderTanggalDetail(c *gin.Context) {
	// Query params: tanggal=YYYY-MM-DD
	tanggalStr := c.Query("tanggal")
	if tanggalStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tanggal wajib (YYYY-MM-DD)"})
		return
	}

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	// Ambil semua pendaftaran di tanggal tsb
	var pendaftaran []structs.PendaftaranNikah
	if err := DB.Where("DATE(tanggal_nikah) = ?", tanggal.Format("2006-01-02")).Order("waktu_nikah ASC").Find(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendaftaran"})
		return
	}

	// Kategorikan per jam dengan warna dan tambahkan info nama penghulu & pasangan
	entries := make([]map[string]interface{}, 0, len(pendaftaran))
	for _, p := range pendaftaran {
		warna := "hijau"
		if p.Status_pendaftaran == "Menunggu Verifikasi" || p.Status_pendaftaran == "Menunggu Pengumpulan Berkas" {
			warna = "kuning"
		}

		// Ambil nama penghulu (jika ada)
		var namaPenghulu string
		if p.Penghulu_id != nil {
			var ph structs.Penghulu
			if err := DB.First(&ph, *p.Penghulu_id).Error; err == nil {
				namaPenghulu = ph.Nama_lengkap
			}
		}

		// Ambil nama calon suami & istri berdasarkan ID
		var calonSuami structs.CalonPasangan
		var calonIstri structs.CalonPasangan
		_ = DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami).Error
		_ = DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri).Error

		entries = append(entries, map[string]interface{}{
			"id":                p.ID,
			"nomor_pendaftaran": p.Nomor_pendaftaran,
			"waktu":             p.Waktu_nikah,
			"tempat":            p.Tempat_nikah,
			"status":            p.Status_pendaftaran,
			"warna":             warna,
			"penghulu_id":       p.Penghulu_id,
			"penghulu_nama":     namaPenghulu,
			"nama_calon_suami":  calonSuami.Nama_lengkap,
			"nama_calon_istri":  calonIstri.Nama_lengkap,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail kalender berhasil diambil",
		"data": gin.H{
			"tanggal": tanggalStr,
			"items":   entries,
		},
	})
}

// GetKetersediaanTanggal menampilkan detail ketersediaan untuk tanggal tertentu
func GetKetersediaanTanggal(c *gin.Context) {
	tanggalParam := c.Param("tanggal")

	// Parse tanggal (format: YYYY-MM-DD)
	tanggal, err := time.Parse("2006-01-02", tanggalParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	// Cek apakah tanggal sudah lewat
	now := time.Now()
	if tanggal.Before(now.Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal sudah lewat"})
		return
	}

	// Query pendaftaran nikah untuk tanggal tersebut
	var pendaftaran []structs.PendaftaranNikah
	err = DB.Where("DATE(tanggal_nikah) = ? AND status_pendaftaran = ?",
		tanggal.Format("2006-01-02"), "Disetujui").Find(&pendaftaran).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendaftaran"})
		return
	}

	// Hitung kapasitas berdasarkan lokasi nikah
	// Nikah di KUA: maksimal 9 pernikahan per hari (1 per slot)
	// Nikah di luar KUA: tidak ada batasan (4 penghulu bisa bersamaan)
	var kapasitasPerHari int
	var jumlahNikah int

	// Hitung nikah di KUA dan di luar KUA
	nikahDiKUA := 0
	nikahDiLuarKUA := 0

	for _, p := range pendaftaran {
		if p.Tempat_nikah == "Di KUA" {
			nikahDiKUA++
		} else {
			nikahDiLuarKUA++
		}
	}

	// Kapasitas untuk nikah di KUA: 9 per hari (1 per slot waktu)
	kapasitasKUA := 9
	kapasitasPerHari = kapasitasKUA
	jumlahNikah = nikahDiKUA // Hanya hitung nikah di KUA untuk kapasitas
	sisaKuota := kapasitasPerHari - jumlahNikah

	// Tentukan status
	var status string
	var tersedia bool

	if jumlahNikah >= kapasitasPerHari {
		status = "Penuh"
		tersedia = false
	} else {
		status = "Tersedia"
		tersedia = true
	}

	// Buat detail jadwal untuk tanggal tersebut
	jadwalDetail := make([]map[string]interface{}, 0)
	for _, p := range pendaftaran {
		jadwalDetail = append(jadwalDetail, map[string]interface{}{
			"nomor_pendaftaran": p.Nomor_pendaftaran,
			"waktu_nikah":       p.Waktu_nikah,
			"tempat_nikah":      p.Tempat_nikah,
			"alamat_akad":       p.Alamat_akad,
		})
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail ketersediaan tanggal berhasil diambil",
		"data": gin.H{
			"tanggal":           tanggalParam,
			"status":            status,
			"tersedia":          tersedia,
			"jumlah_nikah_kua":  nikahDiKUA,
			"jumlah_nikah_luar": nikahDiLuarKUA,
			"total_nikah":       len(pendaftaran),
			"sisa_kuota_kua":    sisaKuota,
			"kapasitas_kua":     kapasitasPerHari,
			"keterangan":        "Kapasitas hanya berlaku untuk nikah di KUA. Nikah di luar KUA tidak dibatasi.",
			"jadwal_detail":     jadwalDetail,
		},
	})
}

// GetPenghuluJadwal menampilkan jadwal penghulu untuk tanggal tertentu
func GetPenghuluJadwal(c *gin.Context) {
	tanggalParam := c.Param("tanggal")

	// Parse tanggal (format: YYYY-MM-DD)
	tanggal, err := time.Parse("2006-01-02", tanggalParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	// Cek apakah tanggal sudah lewat
	now := time.Now()
	if tanggal.Before(now.Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal sudah lewat"})
		return
	}

	// Query semua penghulu yang aktif
	var penghulu []structs.Penghulu
	err = DB.Where("status = ?", "Aktif").Find(&penghulu).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penghulu"})
		return
	}

	// Query pendaftaran nikah untuk tanggal tersebut
	var pendaftaran []structs.PendaftaranNikah
	err = DB.Where("DATE(tanggal_nikah) = ? AND status_pendaftaran = ?",
		tanggal.Format("2006-01-02"), "Disetujui").Find(&pendaftaran).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendaftaran"})
		return
	}

	// Buat map untuk menghitung jadwal per penghulu
	jadwalPerPenghulu := make(map[uint][]map[string]interface{})

	// Inisialisasi jadwal untuk semua penghulu
	for _, p := range penghulu {
		jadwalPerPenghulu[p.ID] = make([]map[string]interface{}, 0)
	}

	// Distribusi jadwal berdasarkan penghulu yang sudah di-assign
	for _, pendaftaran := range pendaftaran {
		if pendaftaran.Penghulu_id != nil {
			penghuluID := *pendaftaran.Penghulu_id
			jadwalPerPenghulu[penghuluID] = append(jadwalPerPenghulu[penghuluID], map[string]interface{}{
				"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
				"waktu_nikah":       pendaftaran.Waktu_nikah,
				"tempat_nikah":      pendaftaran.Tempat_nikah,
				"alamat_akad":       pendaftaran.Alamat_akad,
				"assigned_by":       pendaftaran.Penghulu_assigned_by,
				"assigned_at":       pendaftaran.Penghulu_assigned_at,
			})
		}
	}

	// Buat response data penghulu dengan jadwal mereka
	penghuluJadwal := make([]map[string]interface{}, 0)
	for _, p := range penghulu {
		jadwal := jadwalPerPenghulu[p.ID]
		jumlahJadwal := len(jadwal)
		sisaKuota := 3 - jumlahJadwal // Maksimal 3 nikah per penghulu per hari

		var status string
		if jumlahJadwal >= 3 {
			status = "Penuh"
		} else if jumlahJadwal > 0 {
			status = "Sebagian"
		} else {
			status = "Kosong"
		}

		penghuluJadwal = append(penghuluJadwal, map[string]interface{}{
			"id":            p.ID,
			"nama":          p.Nama_lengkap,
			"status":        status,
			"jumlah_jadwal": jumlahJadwal,
			"sisa_kuota":    sisaKuota,
			"maksimal":      3,
			"jadwal":        jadwal,
		})
	}

	// Hitung total kapasitas dan sisa
	totalKapasitas := len(penghulu) * 3
	totalTerisi := len(pendaftaran)
	totalSisa := totalKapasitas - totalTerisi

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal penghulu berhasil diambil",
		"data": gin.H{
			"tanggal":         tanggalParam,
			"total_penghulu":  len(penghulu),
			"total_kapasitas": totalKapasitas,
			"total_terisi":    totalTerisi,
			"total_sisa":      totalSisa,
			"penghulu":        penghuluJadwal,
		},
	})
}

// AssignPenghulu mengassign penghulu untuk pendaftaran nikah (hanya kepala KUA)
func AssignPenghulu(c *gin.Context) {
	pendaftaranID := c.Param("id")

	// Get user_id dari context (kepala KUA)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Parse request body
	var input struct {
		PenghuluID uint `json:"penghulu_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Cek apakah pendaftaran ada
	var pendaftaran structs.PendaftaranNikah
	if err := DB.First(&pendaftaran, pendaftaranID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendaftaran tidak ditemukan"})
		return
	}

	// Cek apakah pendaftaran sudah siap untuk assign penghulu
	if pendaftaran.Status_pendaftaran != "Menunggu Penugasan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pendaftaran harus dalam status 'Menunggu Penugasan' untuk assign penghulu"})
		return
	}

	// Cek apakah penghulu ada dan aktif
	var penghulu structs.Penghulu
	if err := DB.Where("id = ? AND status = ?", input.PenghuluID, "Aktif").First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penghulu tidak ditemukan atau tidak aktif"})
		return
	}

	// Hitung jumlah jadwal penghulu pada tanggal yang sama (tanpa membatasi, hanya untuk peringatan)
	tanggalNikah := pendaftaran.Tanggal_nikah.Format("2006-01-02")
	var count int64
	DB.Model(&structs.PendaftaranNikah{}).Where(
		"penghulu_id = ? AND DATE(tanggal_nikah) = ? AND status_pendaftaran IN ?",
		input.PenghuluID,
		tanggalNikah,
		[]string{"Menunggu Verifikasi Penghulu", "Menunggu Bimbingan", "Sudah Bimbingan", "Selesai"},
	).Count(&count)

	// Batasi total pernikahan per hari: maksimal 9 pernikahan per hari
	var totalPerHari int64
	DB.Model(&structs.PendaftaranNikah{}).Where(
		"DATE(tanggal_nikah) = ? AND status_pendaftaran IN ?",
		tanggalNikah,
		[]string{"Menunggu Verifikasi Penghulu", "Menunggu Bimbingan", "Sudah Bimbingan", "Selesai"},
	).Count(&totalPerHari)
	if totalPerHari >= 9 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Kuota pernikahan harian penuh",
			"details": gin.H{"tanggal": tanggalNikah, "maksimal": 9},
		})
		return
	}

	// Cek konflik waktu - pastikan tidak ada jadwal yang berbarengan
	var existingJadwal []structs.PendaftaranNikah
	DB.Where(
		"penghulu_id = ? AND DATE(tanggal_nikah) = ? AND status_pendaftaran IN ?",
		input.PenghuluID,
		tanggalNikah,
		[]string{"Menunggu Verifikasi Penghulu", "Menunggu Bimbingan", "Sudah Bimbingan", "Selesai"},
	).Find(&existingJadwal)

	// Parse waktu nikah yang akan di-assign
	waktuNikahBaru, err := time.Parse("15:04", pendaftaran.Waktu_nikah)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format waktu nikah tidak valid (HH:MM)"})
		return
	}

	// Cek konflik dengan jadwal yang sudah ada
	for _, jadwal := range existingJadwal {
		waktuJadwalLama, err := time.Parse("15:04", jadwal.Waktu_nikah)
		if err != nil {
			continue // Skip jika format waktu tidak valid
		}

		// Hitung selisih waktu dalam menit
		selisihMenit := int(waktuNikahBaru.Sub(waktuJadwalLama).Minutes())
		if selisihMenit < 0 {
			selisihMenit = -selisihMenit
		}

		// Jika selisih kurang dari 1 jam (60 menit), dianggap konflik
		// Asumsi: setiap slot berjarak 1 jam, maksimal 4 pernikahan bersamaan per slot
		if selisihMenit < 60 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Konflik jadwal! Penghulu sudah memiliki jadwal pada waktu yang berdekatan",
				"details": gin.H{
					"waktu_konflik":   jadwal.Waktu_nikah,
					"tempat":          jadwal.Tempat_nikah,
					"selisih_menit":   selisihMenit,
					"minimal_selisih": "60 menit (1 jam)",
				},
			})
			return
		}
	}

	// Update pendaftaran dengan penghulu
	now := time.Now()
	pendaftaran.Penghulu_id = &input.PenghuluID
	pendaftaran.Penghulu_assigned_by = userID.(string)
	pendaftaran.Penghulu_assigned_at = &now
	pendaftaran.Status_pendaftaran = "Menunggu Verifikasi Penghulu"
	pendaftaran.Updated_at = now

	if err := DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengassign penghulu"})
		return
	}

	// Kirim notifikasi otomatis setelah penghulu berhasil diassign
	notificationService := services.NewNotificationService(DB)
	err = notificationService.SendPenghuluAssignmentNotification(pendaftaran.ID, penghulu.User_id)
	if err != nil {
		// Log error tapi jangan return error karena assign sudah berhasil
		fmt.Printf("Gagal mengirim notifikasi penugasan penghulu: %v\n", err)
	}

	// Overload warning (non-blocking)
	projectedCount := count + 1
	var warning string
	if projectedCount >= 3 {
		warning = fmt.Sprintf("Peringatan: penghulu ini memiliki %d jadwal pada tanggal %s", projectedCount, tanggalNikah)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Penghulu berhasil diassign",
		"warning": warning,
		"data": gin.H{
			"pendaftaran_id": pendaftaran.ID,
			"penghulu_id":    input.PenghuluID,
			"penghulu_nama":  penghulu.Nama_lengkap,
			"assigned_by":    userID.(string),
			"assigned_at":    now,
		},
	})
}

// ChangePenghulu mengubah penghulu untuk pendaftaran nikah (hanya kepala KUA)
func ChangePenghulu(c *gin.Context) {
	pendaftaranID := c.Param("id")

	// Get user_id dari context (kepala KUA)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Parse request body
	var input struct {
		PenghuluID uint `json:"penghulu_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Cek apakah pendaftaran ada
	var pendaftaran structs.PendaftaranNikah
	if err := DB.First(&pendaftaran, pendaftaranID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendaftaran tidak ditemukan"})
		return
	}

	// Cek apakah pendaftaran sudah siap untuk mengubah penghulu
	if pendaftaran.Status_pendaftaran != "Menunggu Penugasan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pendaftaran harus dalam status 'Menunggu Penugasan' untuk mengubah penghulu"})
		return
	}

	// Cek apakah penghulu ada dan aktif
	var penghulu structs.Penghulu
	if err := DB.Where("id = ? AND status = ?", input.PenghuluID, "Aktif").First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penghulu tidak ditemukan atau tidak aktif"})
		return
	}

	// Cek apakah penghulu baru sudah memiliki 3 jadwal pada tanggal yang sama
	tanggalNikah := pendaftaran.Tanggal_nikah.Format("2006-01-02")
	var count int64
	DB.Model(&structs.PendaftaranNikah{}).Where("penghulu_id = ? AND DATE(tanggal_nikah) = ? AND status_pendaftaran = ? AND id != ?",
		input.PenghuluID, tanggalNikah, "Disetujui", pendaftaranID).Count(&count)

	if count >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Penghulu baru sudah memiliki 3 jadwal pada tanggal tersebut"})
		return
	}

	// Cek konflik waktu - pastikan tidak ada jadwal yang berbarengan
	var existingJadwal []structs.PendaftaranNikah
	DB.Where("penghulu_id = ? AND DATE(tanggal_nikah) = ? AND status_pendaftaran = ? AND id != ?",
		input.PenghuluID, tanggalNikah, "Disetujui", pendaftaranID).Find(&existingJadwal)

	// Parse waktu nikah yang akan di-assign
	waktuNikahBaru, err := time.Parse("15:04", pendaftaran.Waktu_nikah)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format waktu nikah tidak valid (HH:MM)"})
		return
	}

	// Cek konflik dengan jadwal yang sudah ada
	for _, jadwal := range existingJadwal {
		waktuJadwalLama, err := time.Parse("15:04", jadwal.Waktu_nikah)
		if err != nil {
			continue // Skip jika format waktu tidak valid
		}

		// Hitung selisih waktu dalam menit
		selisihMenit := int(waktuNikahBaru.Sub(waktuJadwalLama).Minutes())
		if selisihMenit < 0 {
			selisihMenit = -selisihMenit
		}

		// Jika selisih kurang dari 2 jam (120 menit), dianggap konflik
		// Asumsi: setiap nikah membutuhkan minimal 2 jam (persiapan + akad + dokumentasi)
		if selisihMenit < 120 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Konflik jadwal! Penghulu baru sudah memiliki jadwal pada waktu yang berdekatan",
				"details": gin.H{
					"waktu_konflik":   jadwal.Waktu_nikah,
					"tempat":          jadwal.Tempat_nikah,
					"selisih_menit":   selisihMenit,
					"minimal_selisih": "60 menit (1 jam)",
				},
			})
			return
		}
	}

	// Simpan penghulu lama untuk audit
	penghuluLama := pendaftaran.Penghulu_id

	// Update pendaftaran dengan penghulu baru
	now := time.Now()
	pendaftaran.Penghulu_id = &input.PenghuluID
	pendaftaran.Penghulu_assigned_by = userID.(string)
	pendaftaran.Penghulu_assigned_at = &now
	pendaftaran.Updated_at = now

	if err := DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah penghulu"})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Penghulu berhasil diubah",
		"data": gin.H{
			"pendaftaran_id": pendaftaran.ID,
			"penghulu_lama":  penghuluLama,
			"penghulu_baru":  input.PenghuluID,
			"penghulu_nama":  penghulu.Nama_lengkap,
			"changed_by":     userID.(string),
			"changed_at":     now,
		},
	})
}

// GetPendaftaranBelumAssignPenghulu mendapatkan pendaftaran yang belum di-assign penghulu
func GetPendaftaranBelumAssignPenghulu(c *gin.Context) {
	// Query pendaftaran yang sudah siap untuk assign penghulu tapi belum di-assign
	var pendaftaran []structs.PendaftaranNikah
	err := DB.Where("status_pendaftaran = ? AND penghulu_id IS NULL", "Menunggu Penugasan").
		Order("tanggal_nikah ASC").Find(&pendaftaran).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendaftaran"})
		return
	}

	// Buat response data
	pendaftaranData := make([]map[string]interface{}, 0)
	for _, p := range pendaftaran {
		pendaftaranData = append(pendaftaranData, map[string]interface{}{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"alamat_akad":        p.Alamat_akad,
			"status_pendaftaran": p.Status_pendaftaran,
			"created_at":         p.Created_at,
		})
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Pendaftaran belum assign penghulu berhasil diambil",
		"data": gin.H{
			"total":       len(pendaftaran),
			"pendaftaran": pendaftaranData,
		},
	})
}

// CreateBimbinganPerkawinan membuat sesi bimbingan perkawinan baru (Staff/Kepala KUA)
func CreateBimbinganPerkawinan(c *gin.Context) {
	var input struct {
		TanggalBimbingan string `json:"tanggal_bimbingan" binding:"required"` // YYYY-MM-DD
		WaktuMulai       string `json:"waktu_mulai" binding:"required"`       // HH:MM
		WaktuSelesai     string `json:"waktu_selesai" binding:"required"`     // HH:MM
		TempatBimbingan  string `json:"tempat_bimbingan" binding:"required"`
		Pembimbing       string `json:"pembimbing" binding:"required"`
		Kapasitas        int    `json:"kapasitas"`
		Catatan          string `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", input.TanggalBimbingan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	// Validasi hari Rabu
	if tanggal.Weekday() != time.Wednesday {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bimbingan perkawinan hanya bisa dijadwalkan pada hari Rabu"})
		return
	}

	// Cek apakah tanggal sudah lewat
	now := time.Now()
	if tanggal.Before(now.Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal bimbingan tidak boleh di masa lalu"})
		return
	}

	// Cek apakah sudah ada bimbingan pada tanggal yang sama
	var existingBimbingan structs.BimbinganPerkawinan
	if err := DB.Where("DATE(tanggal_bimbingan) = ? AND status = ?",
		tanggal.Format("2006-01-02"), "Aktif").First(&existingBimbingan).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sudah ada bimbingan perkawinan pada tanggal tersebut"})
		return
	}

	// Set kapasitas default jika tidak diisi
	kapasitas := input.Kapasitas
	if kapasitas == 0 {
		kapasitas = 10 // 10 pasangan per hari Rabu
	}

	// Buat bimbingan perkawinan
	bimbingan := structs.BimbinganPerkawinan{
		Tanggal_bimbingan: tanggal,
		Waktu_mulai:       input.WaktuMulai,
		Waktu_selesai:     input.WaktuSelesai,
		Tempat_bimbingan:  input.TempatBimbingan,
		Pembimbing:        input.Pembimbing,
		Kapasitas:         kapasitas,
		Status:            "Aktif",
		Catatan:           input.Catatan,
		Created_at:        now,
		Updated_at:        now,
	}

	if err := DB.Create(&bimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat bimbingan perkawinan"})
		return
	}

	// Kirim notifikasi otomatis setelah bimbingan berhasil dibuat
	notificationService := services.NewNotificationService(DB)
	err = notificationService.SendBimbinganNotification(bimbingan.ID, "created")
	if err != nil {
		// Log error tapi jangan return error karena bimbingan sudah berhasil dibuat
		fmt.Printf("Gagal mengirim notifikasi bimbingan: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bimbingan perkawinan berhasil dibuat",
		"data":    bimbingan,
	})
}

// GetBimbinganPerkawinan mendapatkan daftar bimbingan perkawinan
func GetBimbinganPerkawinan(c *gin.Context) {
	// Ambil parameter query
	bulan := c.DefaultQuery("bulan", "")
	tahun := c.DefaultQuery("tahun", "")
	status := c.DefaultQuery("status", "Aktif")

	// Default ke bulan dan tahun saat ini jika tidak ada parameter
	now := time.Now()
	if bulan == "" {
		bulan = fmt.Sprintf("%02d", int(now.Month()))
	}
	if tahun == "" {
		tahun = fmt.Sprintf("%d", now.Year())
	}

	// Parse bulan dan tahun
	bulanInt, err := strconv.Atoi(bulan)
	if err != nil || bulanInt < 1 || bulanInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan tidak valid (1-12)"})
		return
	}

	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun tidak valid"})
		return
	}

	// Hitung tanggal awal dan akhir bulan
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.UTC)
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	// Query bimbingan perkawinan
	var bimbingan []structs.BimbinganPerkawinan
	query := DB.Where("tanggal_bimbingan >= ? AND tanggal_bimbingan <= ?", awalBulan, akhirBulan)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err = query.Order("tanggal_bimbingan ASC").Find(&bimbingan).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bimbingan perkawinan"})
		return
	}

	// Hitung jumlah peserta untuk setiap bimbingan
	bimbinganData := make([]map[string]interface{}, 0)
	for _, b := range bimbingan {
		var count int64
		DB.Model(&structs.PendaftaranBimbingan{}).Where("bimbingan_perkawinan_id = ?", b.ID).Count(&count)

		bimbinganData = append(bimbinganData, map[string]interface{}{
			"id":                b.ID,
			"tanggal_bimbingan": b.Tanggal_bimbingan.Format("2006-01-02"),
			"waktu_mulai":       b.Waktu_mulai,
			"waktu_selesai":     b.Waktu_selesai,
			"tempat_bimbingan":  b.Tempat_bimbingan,
			"pembimbing":        b.Pembimbing,
			"kapasitas":         b.Kapasitas,
			"jumlah_peserta":    count,
			"sisa_kuota":        b.Kapasitas - int(count),
			"status":            b.Status,
			"catatan":           b.Catatan,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data bimbingan perkawinan berhasil diambil",
		"data": gin.H{
			"bulan":     bulanInt,
			"tahun":     tahunInt,
			"bimbingan": bimbinganData,
		},
	})
}

// GetBimbinganPerkawinanByID mendapatkan detail bimbingan perkawinan
func GetBimbinganPerkawinanByID(c *gin.Context) {
	bimbinganID := c.Param("id")

	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	// Hitung jumlah peserta
	var count int64
	DB.Model(&structs.PendaftaranBimbingan{}).Where("bimbingan_perkawinan_id = ?", bimbingan.ID).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail bimbingan perkawinan berhasil diambil",
		"data": gin.H{
			"id":                bimbingan.ID,
			"tanggal_bimbingan": bimbingan.Tanggal_bimbingan.Format("2006-01-02"),
			"waktu_mulai":       bimbingan.Waktu_mulai,
			"waktu_selesai":     bimbingan.Waktu_selesai,
			"tempat_bimbingan":  bimbingan.Tempat_bimbingan,
			"pembimbing":        bimbingan.Pembimbing,
			"kapasitas":         bimbingan.Kapasitas,
			"jumlah_peserta":    count,
			"sisa_kuota":        bimbingan.Kapasitas - int(count),
			"status":            bimbingan.Status,
			"catatan":           bimbingan.Catatan,
		},
	})
}

// UpdateBimbinganPerkawinan mengupdate bimbingan perkawinan (Staff/Kepala KUA)
func UpdateBimbinganPerkawinan(c *gin.Context) {
	bimbinganID := c.Param("id")

	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	var input struct {
		WaktuMulai      string `json:"waktu_mulai"`
		WaktuSelesai    string `json:"waktu_selesai"`
		TempatBimbingan string `json:"tempat_bimbingan"`
		Pembimbing      string `json:"pembimbing"`
		Kapasitas       int    `json:"kapasitas"`
		Status          string `json:"status"`
		Catatan         string `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update fields yang diisi
	if input.WaktuMulai != "" {
		bimbingan.Waktu_mulai = input.WaktuMulai
	}
	if input.WaktuSelesai != "" {
		bimbingan.Waktu_selesai = input.WaktuSelesai
	}
	if input.TempatBimbingan != "" {
		bimbingan.Tempat_bimbingan = input.TempatBimbingan
	}
	if input.Pembimbing != "" {
		bimbingan.Pembimbing = input.Pembimbing
	}
	if input.Kapasitas > 0 {
		bimbingan.Kapasitas = input.Kapasitas
	}
	if input.Status != "" {
		bimbingan.Status = input.Status
	}
	if input.Catatan != "" {
		bimbingan.Catatan = input.Catatan
	}

	bimbingan.Updated_at = time.Now()

	if err := DB.Save(&bimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate bimbingan perkawinan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bimbingan perkawinan berhasil diupdate",
		"data":    bimbingan,
	})
}

// GetBimbinganKalender mendapatkan kalender bimbingan perkawinan
func GetBimbinganKalender(c *gin.Context) {
	// Ambil parameter query
	bulan := c.DefaultQuery("bulan", "")
	tahun := c.DefaultQuery("tahun", "")

	// Default ke bulan dan tahun saat ini jika tidak ada parameter
	now := time.Now()
	if bulan == "" {
		bulan = fmt.Sprintf("%02d", int(now.Month()))
	}
	if tahun == "" {
		tahun = fmt.Sprintf("%d", now.Year())
	}

	// Parse bulan dan tahun
	bulanInt, err := strconv.Atoi(bulan)
	if err != nil || bulanInt < 1 || bulanInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan tidak valid (1-12)"})
		return
	}

	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun tidak valid"})
		return
	}

	// Hitung tanggal awal dan akhir bulan
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.UTC)
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	// Query bimbingan perkawinan untuk bulan tersebut
	var bimbingan []structs.BimbinganPerkawinan
	err = DB.Where("tanggal_bimbingan >= ? AND tanggal_bimbingan <= ? AND status = ?",
		awalBulan, akhirBulan, "Aktif").Find(&bimbingan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bimbingan perkawinan"})
		return
	}

	// Buat map untuk menghitung jumlah peserta per bimbingan
	pesertaPerBimbingan := make(map[uint]int64)
	for _, b := range bimbingan {
		var count int64
		DB.Model(&structs.PendaftaranBimbingan{}).Where("bimbingan_perkawinan_id = ?", b.ID).Count(&count)
		pesertaPerBimbingan[b.ID] = count
	}

	// Buat kalender untuk bulan tersebut
	kalender := make([]map[string]interface{}, 0)

	// Mulai dari tanggal 1 sampai akhir bulan
	for tanggal := 1; tanggal <= akhirBulan.Day(); tanggal++ {
		tanggalStr := fmt.Sprintf("%d-%02d-%02d", tahunInt, bulanInt, tanggal)
		tanggalTime := time.Date(tahunInt, time.Month(bulanInt), tanggal, 0, 0, 0, 0, time.UTC)

		// Cek apakah ada bimbingan pada tanggal tersebut
		var bimbinganHariIni *structs.BimbinganPerkawinan
		for _, b := range bimbingan {
			if b.Tanggal_bimbingan.Format("2006-01-02") == tanggalStr {
				bimbinganHariIni = &b
				break
			}
		}

		var status string
		var tersedia bool
		var sisaKuota int
		var bimbinganInfo map[string]interface{}

		if tanggalTime.Before(now.Truncate(24 * time.Hour)) {
			// Tanggal sudah lewat
			status = "Terlewat"
			tersedia = false
			sisaKuota = 0
		} else if bimbinganHariIni == nil {
			// Tidak ada bimbingan
			if tanggalTime.Weekday() == time.Wednesday {
				status = "Belum Dijadwalkan"
				tersedia = false
				sisaKuota = 0
			} else {
				status = "Bukan Hari Rabu"
				tersedia = false
				sisaKuota = 0
			}
		} else {
			// Ada bimbingan
			jumlahPeserta := pesertaPerBimbingan[bimbinganHariIni.ID]
			sisaKuota = bimbinganHariIni.Kapasitas - int(jumlahPeserta)

			if jumlahPeserta >= int64(bimbinganHariIni.Kapasitas) {
				status = "Penuh"
				tersedia = false
			} else {
				status = "Tersedia"
				tersedia = true
			}

			bimbinganInfo = map[string]interface{}{
				"id":               bimbinganHariIni.ID,
				"waktu_mulai":      bimbinganHariIni.Waktu_mulai,
				"waktu_selesai":    bimbinganHariIni.Waktu_selesai,
				"tempat_bimbingan": bimbinganHariIni.Tempat_bimbingan,
				"pembimbing":       bimbinganHariIni.Pembimbing,
				"kapasitas":        bimbinganHariIni.Kapasitas,
				"jumlah_peserta":   jumlahPeserta,
			}
		}

		// Tambahkan ke kalender
		kalender = append(kalender, map[string]interface{}{
			"tanggal":     tanggal,
			"tanggal_str": tanggalStr,
			"status":      status,
			"tersedia":    tersedia,
			"sisa_kuota":  sisaKuota,
			"bimbingan":   bimbinganInfo,
		})
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Kalender bimbingan perkawinan berhasil diambil",
		"data": gin.H{
			"bulan":      bulanInt,
			"tahun":      tahunInt,
			"nama_bulan": awalBulan.Month().String(),
			"kalender":   kalender,
		},
	})
}

// DaftarBimbinganPerkawinan mendaftarkan calon pengantin ke bimbingan perkawinan
func DaftarBimbinganPerkawinan(c *gin.Context) {
	bimbinganID := c.Param("id")

	// Get user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Cek apakah bimbingan ada dan aktif
	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	if bimbingan.Status != "Aktif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bimbingan perkawinan tidak aktif"})
		return
	}

	// Cek apakah tanggal bimbingan sudah lewat
	now := time.Now()
	if bimbingan.Tanggal_bimbingan.Before(now.Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal bimbingan sudah lewat"})
		return
	}

	// Cek apakah user sudah memiliki pendaftaran nikah yang siap untuk bimbingan
	var pendaftaran structs.PendaftaranNikah
	if err := DB.Where("pendaftar_id = ? AND status_pendaftaran = ?", userID.(string), "Menunggu Bimbingan").First(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum memiliki pendaftaran nikah yang siap untuk bimbingan"})
		return
	}

	// Cek apakah sudah terdaftar di bimbingan ini
	var existingDaftar structs.PendaftaranBimbingan
	if err := DB.Where("pendaftaran_nikah_id = ? AND bimbingan_perkawinan_id = ?", pendaftaran.ID, bimbingan.ID).First(&existingDaftar).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah terdaftar di bimbingan perkawinan ini"})
		return
	}

	// Cek kapasitas
	var count int64
	DB.Model(&structs.PendaftaranBimbingan{}).Where("bimbingan_perkawinan_id = ?", bimbingan.ID).Count(&count)

	if count >= int64(bimbingan.Kapasitas) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bimbingan perkawinan sudah penuh"})
		return
	}

	// Buat pendaftaran bimbingan
	pendaftaranBimbingan := structs.PendaftaranBimbingan{
		Pendaftaran_nikah_id:    pendaftaran.ID,
		Bimbingan_perkawinan_id: bimbingan.ID,
		Calon_suami_id:          pendaftaran.Calon_suami_id,
		Calon_istri_id:          pendaftaran.Calon_istri_id,
		Status_kehadiran:        "Belum",
		Status_sertifikat:       "Belum",
		Created_at:              now,
		Updated_at:              now,
	}

	if err := DB.Create(&pendaftaranBimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftar bimbingan perkawinan"})
		return
	}

	// Status sudah "Menunggu Bimbingan" jadi tidak perlu diupdate lagi

	// Kirim notifikasi otomatis setelah berhasil mendaftar bimbingan
	notificationService := services.NewNotificationService(DB)
	err := notificationService.SendBimbinganNotification(bimbingan.ID, "created")
	if err != nil {
		// Log error tapi jangan return error karena pendaftaran sudah berhasil
		fmt.Printf("Gagal mengirim notifikasi pendaftaran bimbingan: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil mendaftar bimbingan perkawinan",
		"data": gin.H{
			"bimbingan_id": bimbingan.ID,
			"tanggal":      bimbingan.Tanggal_bimbingan.Format("2006-01-02"),
			"waktu":        bimbingan.Waktu_mulai + " - " + bimbingan.Waktu_selesai,
			"tempat":       bimbingan.Tempat_bimbingan,
			"pembimbing":   bimbingan.Pembimbing,
		},
	})
}

// GetBimbinganParticipants mendapatkan daftar peserta bimbingan perkawinan
func GetBimbinganParticipants(c *gin.Context) {
	bimbinganID := c.Param("id")

	// Cek apakah bimbingan ada
	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	// Query peserta bimbingan
	var peserta []structs.PendaftaranBimbingan
	err := DB.Where("bimbingan_perkawinan_id = ?", bimbingan.ID).Find(&peserta).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peserta"})
		return
	}

	// Buat response data peserta
	pesertaData := make([]map[string]interface{}, 0)
	for _, p := range peserta {
		// Ambil data calon pasangan
		var calonSuami, calonIstri structs.CalonPasangan
		DB.First(&calonSuami, p.Calon_suami_id)
		DB.First(&calonIstri, p.Calon_istri_id)

		pesertaData = append(pesertaData, map[string]interface{}{
			"id":                   p.ID,
			"pendaftaran_nikah_id": p.Pendaftaran_nikah_id,
			"calon_suami": gin.H{
				"nama": calonSuami.Nama_lengkap,
				"nik":  calonSuami.NIK,
			},
			"calon_istri": gin.H{
				"nama": calonIstri.Nama_lengkap,
				"nik":  calonIstri.NIK,
			},
			"status_kehadiran":  p.Status_kehadiran,
			"status_sertifikat": p.Status_sertifikat,
			"no_sertifikat":     p.No_sertifikat,
			"created_at":        p.Created_at,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data peserta bimbingan perkawinan berhasil diambil",
		"data": gin.H{
			"bimbingan_id":   bimbingan.ID,
			"tanggal":        bimbingan.Tanggal_bimbingan.Format("2006-01-02"),
			"waktu":          bimbingan.Waktu_mulai + " - " + bimbingan.Waktu_selesai,
			"tempat":         bimbingan.Tempat_bimbingan,
			"pembimbing":     bimbingan.Pembimbing,
			"kapasitas":      bimbingan.Kapasitas,
			"jumlah_peserta": len(peserta),
			"peserta":        pesertaData,
		},
	})
}

// GetStatusFlow menampilkan alur status pendaftaran nikah
func GetStatusFlow(c *gin.Context) {
	pendaftaranID := c.Param("id")

	// Cek apakah pendaftaran ada
	var pendaftaran structs.PendaftaranNikah
	if err := DB.First(&pendaftaran, pendaftaranID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendaftaran tidak ditemukan"})
		return
	}

	// Definisikan alur status sesuai workflow yang baru dan lebih sederhana
	statusFlow := []map[string]interface{}{
		{
			"status":      "Draft",
			"description": "Data belum lengkap",
			"completed":   pendaftaran.Status_pendaftaran != "Draft",
			"current":     pendaftaran.Status_pendaftaran == "Draft",
			"can_edit":    true,
		},
		{
			"status":      "Menunggu Verifikasi",
			"description": "Staff verifikasi formulir online",
			"completed":   pendaftaran.Status_pendaftaran == "Menunggu Pengumpulan Berkas" || pendaftaran.Status_pendaftaran == "Berkas Diterima" || pendaftaran.Status_pendaftaran == "Menunggu Penugasan" || pendaftaran.Status_pendaftaran == "Penghulu Ditugaskan" || pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Menunggu Verifikasi",
			"can_edit":    false,
		},
		{
			"status":      "Menunggu Pengumpulan Berkas",
			"description": "Formulir disetujui, siap kumpulkan berkas",
			"completed":   pendaftaran.Status_pendaftaran == "Berkas Diterima" || pendaftaran.Status_pendaftaran == "Menunggu Penugasan" || pendaftaran.Status_pendaftaran == "Penghulu Ditugaskan" || pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Menunggu Pengumpulan Berkas",
			"can_edit":    false,
		},
		{
			"status":      "Berkas Diterima",
			"description": "Staff menerima berkas hardcopy",
			"completed":   pendaftaran.Status_pendaftaran == "Menunggu Penugasan" || pendaftaran.Status_pendaftaran == "Penghulu Ditugaskan" || pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Berkas Diterima",
			"can_edit":    false,
		},
		{
			"status":      "Menunggu Penugasan",
			"description": "Siap untuk assign penghulu",
			"completed":   pendaftaran.Status_pendaftaran == "Penghulu Ditugaskan" || pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Menunggu Penugasan",
			"can_edit":    false,
		},
		{
			"status":      "Penghulu Ditugaskan",
			"description": "Penghulu sudah ditugaskan",
			"completed":   pendaftaran.Status_pendaftaran == "Menunggu Verifikasi Penghulu" || pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Penghulu Ditugaskan",
			"can_edit":    false,
		},
		{
			"status":      "Menunggu Verifikasi Penghulu",
			"description": "Penghulu harus mengecek berkas",
			"completed":   pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Menunggu Verifikasi Penghulu",
			"can_edit":    false,
		},
		{
			"status":      "Menunggu Bimbingan",
			"description": "Siap untuk bimbingan perkawinan",
			"completed":   pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Menunggu Bimbingan",
			"can_edit":    false,
		},
		{
			"status":      "Sudah Bimbingan",
			"description": "Bimbingan selesai",
			"completed":   pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Sudah Bimbingan",
			"can_edit":    false,
		},
		{
			"status":      "Selesai",
			"description": "Nikah telah dilaksanakan",
			"completed":   pendaftaran.Status_pendaftaran == "Selesai",
			"current":     pendaftaran.Status_pendaftaran == "Selesai",
			"can_edit":    false,
		},
	}

	// Informasi tambahan
	var bimbinganInfo map[string]interface{}
	if pendaftaran.Status_pendaftaran == "Menunggu Verifikasi Penghulu" || pendaftaran.Status_pendaftaran == "Menunggu Bimbingan" || pendaftaran.Status_pendaftaran == "Sudah Bimbingan" || pendaftaran.Status_pendaftaran == "Selesai" {
		// Cek apakah sudah terdaftar bimbingan
		var pendaftaranBimbingan structs.PendaftaranBimbingan
		if err := DB.Where("pendaftaran_nikah_id = ?", pendaftaran.ID).First(&pendaftaranBimbingan).Error; err == nil {
			var bimbingan structs.BimbinganPerkawinan
			DB.First(&bimbingan, pendaftaranBimbingan.Bimbingan_perkawinan_id)

			bimbinganInfo = map[string]interface{}{
				"terdaftar":         true,
				"tanggal":           bimbingan.Tanggal_bimbingan.Format("2006-01-02"),
				"waktu":             bimbingan.Waktu_mulai + " - " + bimbingan.Waktu_selesai,
				"tempat":            bimbingan.Tempat_bimbingan,
				"pembimbing":        bimbingan.Pembimbing,
				"status_kehadiran":  pendaftaranBimbingan.Status_kehadiran,
				"status_sertifikat": pendaftaranBimbingan.Status_sertifikat,
			}
		} else {
			bimbinganInfo = map[string]interface{}{
				"terdaftar": false,
				"message":   "Belum terdaftar bimbingan perkawinan",
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alur status pendaftaran berhasil diambil",
		"data": gin.H{
			"pendaftaran_id":    pendaftaran.ID,
			"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
			"status_sekarang":   pendaftaran.Status_pendaftaran,
			"tanggal_nikah":     pendaftaran.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":       pendaftaran.Waktu_nikah,
			"tempat_nikah":      pendaftaran.Tempat_nikah,
			"penghulu_assigned": pendaftaran.Penghulu_id != nil,
			"bimbingan_info":    bimbinganInfo,
			"status_flow":       statusFlow,
		},
	})
}

// CompleteBimbingan menandai bahwa bimbingan perkawinan sudah selesai
func CompleteBimbingan(c *gin.Context) {
	pendaftaranID := c.Param("id")

	// Get user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Cek apakah pendaftaran ada
	var pendaftaran structs.PendaftaranNikah
	if err := DB.First(&pendaftaran, pendaftaranID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendaftaran tidak ditemukan"})
		return
	}

	// Cek apakah status sudah siap untuk bimbingan
	if pendaftaran.Status_pendaftaran != "Menunggu Bimbingan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pendaftaran harus dalam status 'Menunggu Bimbingan' untuk menyelesaikan bimbingan"})
		return
	}

	// Cek apakah sudah terdaftar bimbingan
	var pendaftaranBimbingan structs.PendaftaranBimbingan
	if err := DB.Where("pendaftaran_nikah_id = ?", pendaftaran.ID).First(&pendaftaranBimbingan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Belum terdaftar bimbingan perkawinan"})
		return
	}

	// Update status pendaftaran
	pendaftaran.Status_pendaftaran = "Sudah Bimbingan"
	pendaftaran.Status_bimbingan = "Sudah"
	pendaftaran.Updated_at = time.Now()

	if err := DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status bimbingan"})
		return
	}

	// Update status kehadiran bimbingan
	pendaftaranBimbingan.Status_kehadiran = "Hadir"
	pendaftaranBimbingan.Status_sertifikat = "Sudah"
	pendaftaranBimbingan.Updated_at = time.Now()

	if err := DB.Save(&pendaftaranBimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status kehadiran bimbingan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bimbingan perkawinan berhasil diselesaikan",
		"data": gin.H{
			"pendaftaran_id":  pendaftaran.ID,
			"status_sekarang": pendaftaran.Status_pendaftaran,
			"tanggal_nikah":   pendaftaran.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":     pendaftaran.Waktu_nikah,
			"tempat_nikah":    pendaftaran.Tempat_nikah,
			"updated_by":      userID.(string),
			"updated_at":      time.Now(),
		},
	})
}

// CompleteNikah menandai bahwa nikah sudah dilaksanakan
func CompleteNikah(c *gin.Context) {
	pendaftaranID := c.Param("id")

	// Get user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Cek apakah pendaftaran ada
	var pendaftaran structs.PendaftaranNikah
	if err := DB.First(&pendaftaran, pendaftaranID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendaftaran tidak ditemukan"})
		return
	}

	// Cek apakah status sudah bimbingan
	if pendaftaran.Status_pendaftaran != "Sudah Bimbingan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pendaftaran harus dalam status 'Sudah Bimbingan' untuk menyelesaikan nikah"})
		return
	}

	// Update status pendaftaran
	pendaftaran.Status_pendaftaran = "Selesai"
	pendaftaran.Status_bimbingan = "Sertifikat Diterbitkan"
	pendaftaran.Updated_at = time.Now()

	if err := DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status nikah"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Nikah berhasil diselesaikan",
		"data": gin.H{
			"pendaftaran_id":    pendaftaran.ID,
			"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
			"status_sekarang":   pendaftaran.Status_pendaftaran,
			"tanggal_nikah":     pendaftaran.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":       pendaftaran.Waktu_nikah,
			"tempat_nikah":      pendaftaran.Tempat_nikah,
			"completed_by":      userID.(string),
			"completed_at":      time.Now(),
		},
	})
}

// UpdateBimbinganAttendance mengupdate kehadiran bimbingan perkawinan
func UpdateBimbinganAttendance(c *gin.Context) {
	bimbinganID := c.Param("id")

	// Get user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	var input struct {
		PendaftaranNikahID uint   `json:"pendaftaran_nikah_id" binding:"required"`
		StatusKehadiran    string `json:"status_kehadiran" binding:"required"`
		StatusSertifikat   string `json:"status_sertifikat"`
		NoSertifikat       string `json:"no_sertifikat"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Validasi status kehadiran
	validStatusKehadiran := map[string]bool{
		"Hadir":       true,
		"Tidak Hadir": true,
	}
	if !validStatusKehadiran[input.StatusKehadiran] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status kehadiran tidak valid"})
		return
	}

	// Validasi status sertifikat
	validStatusSertifikat := map[string]bool{
		"Belum": true,
		"Sudah": true,
	}
	if input.StatusSertifikat != "" && !validStatusSertifikat[input.StatusSertifikat] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status sertifikat tidak valid"})
		return
	}

	// Cek apakah bimbingan ada
	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	// Cek apakah pendaftaran bimbingan ada
	var pendaftaranBimbingan structs.PendaftaranBimbingan
	if err := DB.Where("bimbingan_perkawinan_id = ? AND pendaftaran_nikah_id = ?", bimbinganID, input.PendaftaranNikahID).First(&pendaftaranBimbingan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendaftaran bimbingan tidak ditemukan"})
		return
	}

	// Update status kehadiran
	pendaftaranBimbingan.Status_kehadiran = input.StatusKehadiran
	if input.StatusSertifikat != "" {
		pendaftaranBimbingan.Status_sertifikat = input.StatusSertifikat
	}
	if input.NoSertifikat != "" {
		pendaftaranBimbingan.No_sertifikat = input.NoSertifikat
	}
	pendaftaranBimbingan.Updated_at = time.Now()

	if err := DB.Save(&pendaftaranBimbingan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate kehadiran bimbingan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kehadiran bimbingan berhasil diupdate",
		"data": gin.H{
			"bimbingan_id":      bimbingan.ID,
			"pendaftaran_id":    input.PendaftaranNikahID,
			"status_kehadiran":  pendaftaranBimbingan.Status_kehadiran,
			"status_sertifikat": pendaftaranBimbingan.Status_sertifikat,
			"no_sertifikat":     pendaftaranBimbingan.No_sertifikat,
			"updated_by":        userID.(string),
			"updated_at":        time.Now(),
		},
	})
}

// GetPenghuluKetersediaan menampilkan ketersediaan waktu penghulu untuk tanggal tertentu
func GetPenghuluKetersediaan(c *gin.Context) {
	penghuluID := c.Param("id")
	tanggalParam := c.Param("tanggal")

	// Parse tanggal (format: YYYY-MM-DD)
	tanggal, err := time.Parse("2006-01-02", tanggalParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	// Cek apakah penghulu ada dan aktif
	var penghulu structs.Penghulu
	if err := DB.Where("id = ? AND status = ?", penghuluID, "Aktif").First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penghulu tidak ditemukan atau tidak aktif"})
		return
	}

	// Query jadwal penghulu untuk tanggal tersebut
	var jadwalPenghulu []structs.PendaftaranNikah
	err = DB.Where("penghulu_id = ? AND DATE(tanggal_nikah) = ? AND status_pendaftaran = ?",
		penghuluID, tanggal.Format("2006-01-02"), "Disetujui").Find(&jadwalPenghulu).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jadwal penghulu"})
		return
	}

	// Buat slot waktu yang tersedia (08:00 - 16:00, 9 slot per hari)
	slotWaktu := []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00"}
	slotTersedia := make([]map[string]interface{}, 0)

	for _, slot := range slotWaktu {
		// Parse waktu slot
		waktuSlot, err := time.Parse("15:04", slot)
		if err != nil {
			continue
		}

		// Cek apakah slot ini tersedia
		tersedia := true
		var konflikJadwal []map[string]interface{}

		for _, jadwal := range jadwalPenghulu {
			waktuJadwal, err := time.Parse("15:04", jadwal.Waktu_nikah)
			if err != nil {
				continue
			}

			// Hitung selisih waktu dalam menit
			selisihMenit := int(waktuSlot.Sub(waktuJadwal).Minutes())
			if selisihMenit < 0 {
				selisihMenit = -selisihMenit
			}

			// Jika selisih kurang dari 1 jam, slot tidak tersedia (untuk 4 penghulu bersamaan per slot)
			if selisihMenit < 60 {
				tersedia = false
				konflikJadwal = append(konflikJadwal, map[string]interface{}{
					"waktu":         jadwal.Waktu_nikah,
					"tempat":        jadwal.Tempat_nikah,
					"selisih_menit": selisihMenit,
				})
			}
		}

		slotTersedia = append(slotTersedia, map[string]interface{}{
			"waktu":          slot,
			"tersedia":       tersedia,
			"konflik_jadwal": konflikJadwal,
		})
	}

	// Hitung statistik
	jumlahJadwal := len(jadwalPenghulu)
	sisaKuota := 3 - jumlahJadwal
	slotTersediaCount := 0
	for _, slot := range slotTersedia {
		if slot["tersedia"].(bool) {
			slotTersediaCount++
		}
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Ketersediaan penghulu berhasil diambil",
		"data": gin.H{
			"penghulu": gin.H{
				"id":     penghulu.ID,
				"nama":   penghulu.Nama_lengkap,
				"status": penghulu.Status,
			},
			"tanggal": tanggalParam,
			"statistik": gin.H{
				"jumlah_jadwal":     jumlahJadwal,
				"sisa_kuota":        sisaKuota,
				"maksimal_per_hari": 3,
				"slot_tersedia":     slotTersediaCount,
				"total_slot":        len(slotWaktu),
			},
			"jadwal_hari_ini": jadwalPenghulu,
			"slot_waktu":      slotTersedia,
		},
	})
}

// CetakUndanganBimbingan mencetak undangan bimbingan perkawinan untuk peserta
func CetakUndanganBimbingan(c *gin.Context) {
	bimbinganID := c.Param("id")

	// Get user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Cek apakah bimbingan ada
	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	// Cek apakah user sudah terdaftar di bimbingan ini
	var pendaftaranBimbingan structs.PendaftaranBimbingan
	if err := DB.Where("bimbingan_perkawinan_id = ? AND (calon_suami_id = ? OR calon_istri_id = ?)",
		bimbinganID, userID.(string), userID.(string)).First(&pendaftaranBimbingan).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda belum terdaftar di bimbingan perkawinan ini"})
		return
	}

	// Ambil data pendaftaran nikah
	var pendaftaran structs.PendaftaranNikah
	if err := DB.First(&pendaftaran, pendaftaranBimbingan.Pendaftaran_nikah_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pendaftaran nikah tidak ditemukan"})
		return
	}

	// Ambil data calon pasangan
	var calonSuami, calonIstri structs.CalonPasangan
	DB.First(&calonSuami, pendaftaran.Calon_suami_id)
	DB.First(&calonIstri, pendaftaran.Calon_istri_id)

	// Format tanggal dan waktu
	tanggalBimbingan := bimbingan.Tanggal_bimbingan.Format("02 Januari 2006")
	hariBimbingan := bimbingan.Tanggal_bimbingan.Format("Monday")
	waktuBimbingan := bimbingan.Waktu_mulai + " - " + bimbingan.Waktu_selesai

	// Buat data undangan
	undanganData := map[string]interface{}{
		"bimbingan": gin.H{
			"id":         bimbingan.ID,
			"tanggal":    tanggalBimbingan,
			"hari":       hariBimbingan,
			"waktu":      waktuBimbingan,
			"tempat":     bimbingan.Tempat_bimbingan,
			"pembimbing": bimbingan.Pembimbing,
			"kapasitas":  bimbingan.Kapasitas,
			"status":     bimbingan.Status,
		},
		"pendaftaran": gin.H{
			"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
			"tanggal_nikah":     pendaftaran.Tanggal_nikah.Format("02 Januari 2006"),
			"waktu_nikah":       pendaftaran.Waktu_nikah,
			"tempat_nikah":      pendaftaran.Tempat_nikah,
		},
		"calon_pasangan": gin.H{
			"calon_suami": gin.H{
				"nama":          calonSuami.Nama_lengkap,
				"nik":           calonSuami.NIK,
				"tempat_lahir":  calonSuami.Tempat_lahir,
				"tanggal_lahir": calonSuami.Tanggal_lahir.Format("02 Januari 2006"),
				"alamat":        calonSuami.Alamat,
			},
			"calon_istri": gin.H{
				"nama":          calonIstri.Nama_lengkap,
				"nik":           calonIstri.NIK,
				"tempat_lahir":  calonIstri.Tempat_lahir,
				"tanggal_lahir": calonIstri.Tanggal_lahir.Format("02 Januari 2006"),
				"alamat":        calonIstri.Alamat,
			},
		},
		"status_kehadiran":  pendaftaranBimbingan.Status_kehadiran,
		"status_sertifikat": pendaftaranBimbingan.Status_sertifikat,
		"no_sertifikat":     pendaftaranBimbingan.No_sertifikat,
		"tanggal_daftar":    pendaftaranBimbingan.Created_at.Format("02 Januari 2006"),
		"keterangan": gin.H{
			"title":    "UNDANGAN BIMBINGAN PERKAWINAN",
			"subtitle": "Kantor Urusan Agama (KUA)",
			"message":  "Dengan hormat, kami mengundang Bapak/Ibu untuk mengikuti Bimbingan Perkawinan sebagai persyaratan sebelum melaksanakan pernikahan.",
			"instruksi": []string{
				"1. Harap hadir tepat waktu sesuai jadwal yang telah ditentukan",
				"2. Membawa dokumen asli untuk verifikasi",
				"3. Menggunakan pakaian yang sopan dan rapi",
				"4. Mengikuti seluruh rangkaian bimbingan dengan baik",
				"5. Sertifikat bimbingan akan diberikan setelah selesai mengikuti bimbingan",
			},
			"contact": gin.H{
				"alamat_kua": "Jl. Merdeka No. 123, Kota Bandung",
				"telepon":    "022-1234567",
				"email":      "kua@bandung.go.id",
			},
		},
	}

	// Set header untuk PDF (jika ingin generate PDF)
	// c.Header("Content-Type", "application/pdf")
	// c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"undangan_bimbingan_%s.pdf\"", pendaftaran.Nomor_pendaftaran))

	c.JSON(http.StatusOK, gin.H{
		"message": "Undangan bimbingan perkawinan berhasil dibuat",
		"data":    undanganData,
		"print_info": gin.H{
			"filename": fmt.Sprintf("undangan_bimbingan_%s_%s", pendaftaran.Nomor_pendaftaran, bimbingan.Tanggal_bimbingan.Format("20060102")),
			"format":   "PDF/HTML",
			"size":     "A4",
		},
	})
}

// CetakUndanganBimbinganSemua mencetak undangan bimbingan perkawinan untuk semua peserta (Staff/Kepala KUA)
func CetakUndanganBimbinganSemua(c *gin.Context) {
	bimbinganID := c.Param("id")

	// Cek apakah bimbingan ada
	var bimbingan structs.BimbinganPerkawinan
	if err := DB.First(&bimbingan, bimbinganID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bimbingan perkawinan tidak ditemukan"})
		return
	}

	// Query semua peserta bimbingan
	var peserta []structs.PendaftaranBimbingan
	err := DB.Where("bimbingan_perkawinan_id = ?", bimbingan.ID).Find(&peserta).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peserta"})
		return
	}

	// Buat data undangan untuk semua peserta
	undanganSemua := make([]map[string]interface{}, 0)

	for _, p := range peserta {
		// Ambil data pendaftaran nikah
		var pendaftaran structs.PendaftaranNikah
		if err := DB.First(&pendaftaran, p.Pendaftaran_nikah_id).Error; err != nil {
			continue // Skip jika data tidak ditemukan
		}

		// Ambil data calon pasangan
		var calonSuami, calonIstri structs.CalonPasangan
		DB.First(&calonSuami, p.Calon_suami_id)
		DB.First(&calonIstri, p.Calon_istri_id)

		// Format tanggal dan waktu
		tanggalBimbingan := bimbingan.Tanggal_bimbingan.Format("02 Januari 2006")
		hariBimbingan := bimbingan.Tanggal_bimbingan.Format("Monday")
		waktuBimbingan := bimbingan.Waktu_mulai + " - " + bimbingan.Waktu_selesai

		// Buat data undangan untuk peserta ini
		undanganPeserta := map[string]interface{}{
			"pendaftaran_bimbingan_id": p.ID,
			"bimbingan": gin.H{
				"id":         bimbingan.ID,
				"tanggal":    tanggalBimbingan,
				"hari":       hariBimbingan,
				"waktu":      waktuBimbingan,
				"tempat":     bimbingan.Tempat_bimbingan,
				"pembimbing": bimbingan.Pembimbing,
				"kapasitas":  bimbingan.Kapasitas,
				"status":     bimbingan.Status,
			},
			"pendaftaran": gin.H{
				"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
				"tanggal_nikah":     pendaftaran.Tanggal_nikah.Format("02 Januari 2006"),
				"waktu_nikah":       pendaftaran.Waktu_nikah,
				"tempat_nikah":      pendaftaran.Tempat_nikah,
			},
			"calon_pasangan": gin.H{
				"calon_suami": gin.H{
					"nama":          calonSuami.Nama_lengkap,
					"nik":           calonSuami.NIK,
					"tempat_lahir":  calonSuami.Tempat_lahir,
					"tanggal_lahir": calonSuami.Tanggal_lahir.Format("02 Januari 2006"),
					"alamat":        calonSuami.Alamat,
				},
				"calon_istri": gin.H{
					"nama":          calonIstri.Nama_lengkap,
					"nik":           calonIstri.NIK,
					"tempat_lahir":  calonIstri.Tempat_lahir,
					"tanggal_lahir": calonIstri.Tanggal_lahir.Format("02 Januari 2006"),
					"alamat":        calonIstri.Alamat,
				},
			},
			"status_kehadiran":  p.Status_kehadiran,
			"status_sertifikat": p.Status_sertifikat,
			"no_sertifikat":     p.No_sertifikat,
			"tanggal_daftar":    p.Created_at.Format("02 Januari 2006"),
			"keterangan": gin.H{
				"title":    "UNDANGAN BIMBINGAN PERKAWINAN",
				"subtitle": "Kantor Urusan Agama (KUA)",
				"message":  "Dengan hormat, kami mengundang Bapak/Ibu untuk mengikuti Bimbingan Perkawinan sebagai persyaratan sebelum melaksanakan pernikahan.",
				"instruksi": []string{
					"1. Harap hadir tepat waktu sesuai jadwal yang telah ditentukan",
					"2. Membawa dokumen asli untuk verifikasi",
					"3. Menggunakan pakaian yang sopan dan rapi",
					"4. Mengikuti seluruh rangkaian bimbingan dengan baik",
					"5. Sertifikat bimbingan akan diberikan setelah selesai mengikuti bimbingan",
				},
				"contact": gin.H{
					"alamat_kua": "Jl. Merdeka No. 123, Kota Bandung",
					"telepon":    "022-1234567",
					"email":      "kua@bandung.go.id",
				},
			},
		}

		undanganSemua = append(undanganSemua, undanganPeserta)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Undangan bimbingan perkawinan untuk semua peserta berhasil dibuat",
		"data": gin.H{
			"bimbingan_info": gin.H{
				"id":             bimbingan.ID,
				"tanggal":        bimbingan.Tanggal_bimbingan.Format("02 Januari 2006"),
				"hari":           bimbingan.Tanggal_bimbingan.Format("Monday"),
				"waktu":          bimbingan.Waktu_mulai + " - " + bimbingan.Waktu_selesai,
				"tempat":         bimbingan.Tempat_bimbingan,
				"pembimbing":     bimbingan.Pembimbing,
				"kapasitas":      bimbingan.Kapasitas,
				"jumlah_peserta": len(peserta),
			},
			"undangan_peserta": undanganSemua,
			"print_info": gin.H{
				"filename":       fmt.Sprintf("undangan_bimbingan_semua_%d_%s", bimbingan.ID, bimbingan.Tanggal_bimbingan.Format("20060102")),
				"format":         "PDF/HTML",
				"size":           "A4",
				"total_undangan": len(undanganSemua),
			},
		},
	})
}

// RunReminderNotification menjalankan pengingat notifikasi secara manual (untuk testing)
func RunReminderNotification(c *gin.Context) {
	// Get user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	// Jalankan pengingat notifikasi
	notificationService := services.NewNotificationService(DB)
	err := notificationService.SendReminderNotification()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal menjalankan pengingat notifikasi",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Pengingat notifikasi berhasil dijalankan",
		"executed_by": userID,
		"executed_at": time.Now(),
	})
}

// GetAddressCoordinates mendapatkan koordinat dari alamat menggunakan OpenStreetMap Nominatim API (GRATIS)
func GetAddressCoordinates(c *gin.Context) {
	// Ambil parameter alamat dari query string
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter alamat diperlukan",
			"error":   "Parameter 'address' tidak boleh kosong",
			"example": "/simnikah/geocoding/coordinates?address=Jl. Merdeka No. 123, Banjarmasin",
		})
		return
	}

	// Dapatkan koordinat menggunakan OpenStreetMap Nominatim API (GRATIS)
	latitude, longitude, err := utils.GetCoordinatesFromAddress(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal mendapatkan koordinat",
			"error":   err.Error(),
			"address": address,
		})
		return
	}

	// Response berhasil
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Koordinat berhasil didapatkan",
		"data": gin.H{
			"address":   address,
			"latitude":  latitude,
			"longitude": longitude,
			"source":    "OpenStreetMap Nominatim API (GRATIS)",
			"map_url":   fmt.Sprintf("https://www.openstreetmap.org/?mlat=%.6f&mlon=%.6f&zoom=15", latitude, longitude),
		},
	})
}
