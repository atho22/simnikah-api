package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"simnikah/config"
	"simnikah/internal/handlers/auth"
	"simnikah/internal/handlers/catin"
	"simnikah/internal/handlers/dashboard"
	"simnikah/internal/handlers/kepala_kua"
	"simnikah/internal/handlers/notification"
	"simnikah/internal/handlers/penghulu"
	"simnikah/internal/handlers/staff"
	"simnikah/internal/middleware"
	structs "simnikah/internal/models"
	"simnikah/internal/seeders"
	"simnikah/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	var err error

	// Initialize database connection
	DB, err = config.ConnectDB()
	if err != nil {
		log.Fatal("Koneksi ke database gagal:", err)
	}

	// Migrate struct (scheduling-only models)
	log.Println("Starting database migration...")
	if err := DB.AutoMigrate(
		&structs.Users{},
		&structs.StaffKUA{},
		&structs.Penghulu{},
		&structs.PendaftaranNikah{},
		&structs.Notifikasi{},
	); err != nil {
		log.Fatal("Database migration failed:", err)
	}
	log.Println("Database migration completed successfully")

	// Add database indexes for performance optimization
	if err := config.AddDatabaseIndexes(DB); err != nil {
		log.Println("Warning: Failed to add database indexes:", err)
	}

	// Seed initial data (Kepala KUA, Staff, Penghulu)
	if err := seeders.SeedKepalaKUA(DB); err != nil {
		log.Printf("Warning: Failed to seed kepala KUA: %v", err)
	}
	if err := seeders.SeedStaff(DB); err != nil {
		log.Printf("Warning: Failed to seed staff: %v", err)
	}
	if err := seeders.SeedPenghulu(DB); err != nil {
		log.Printf("Warning: Failed to seed penghulu: %v", err)
	}

	// Set Gin to release mode in production
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	r := createGinEngine(ginMode)

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

	// Initialize handlers
	authHandler := &auth.InDB{DB: DB}
	catinHandler := &catin.InDB{DB: DB}
	staffHandler := &staff.InDB{DB: DB}
	penghuluHandler := &penghulu.InDB{DB: DB}
	kepalaKuaHandler := &kepala_kua.InDB{DB: DB}
	notificationHandler := &notification.InDB{DB: DB}
	dashboardHandler := &dashboard.InDB{DB: DB}

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

	// ==================== AUTH ROUTES ====================
	r.POST("/register", middleware.StrictRateLimiter(), authHandler.Register)
	r.POST("/login", middleware.StrictRateLimiter(), authHandler.Login)
	r.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
	r.POST("/upload-photo", middleware.AuthMiddleware(), authHandler.UploadProfilePhoto)

	// ==================== SIMNIKAH ROUTES ====================
	simnikahRoutes := r.Group("/simnikah")
	{
		// ==================== STAGE 1: CATIN CHECK SCHEDULE & REGISTER ====================
		simnikahRoutes.POST("/check-schedule", middleware.AuthMiddleware(), catinHandler.CheckScheduleAvailability)
		simnikahRoutes.POST("/pendaftaran", middleware.AuthMiddleware(), catinHandler.CreateRegistration)

		// ==================== STAGE 2: FORWARD CHAINING RECOMMENDATION (KEPALA KUA) ====================
		simnikahRoutes.GET("/kepala-kua/forward-chaining/recommendation/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.RecommendPenghuluWithForwardChaining)
		simnikahRoutes.GET("/kepala-kua/forward-chaining/evaluation/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.GetDetailedEvaluationReport)
		simnikahRoutes.GET("/kepala-kua/forward-chaining/config", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.GetForwardChainingConfig)

		// ==================== STAGE 3: ASSIGNMENT APPROVAL (KEPALA KUA) ====================
		simnikahRoutes.POST("/kepala-kua/forward-chaining/assign/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.AssignPenghuluWithApproval)

		// ==================== STAGE 4: PENGHULU VIEW ASSIGNMENTS ====================
		simnikahRoutes.GET("/penghulu/jadwal-penugasan", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("penghulu", "kepala_kua"), penghuluHandler.GetJadwalPenugasan)

		// ==================== KEPALA KUA MANAGEMENT ====================
		simnikahRoutes.GET("/kepala-kua/available-penghulu", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.ListAvailableOfficers)
		simnikahRoutes.GET("/kepala-kua/penghulu-tersedia", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), kepalaKuaHandler.GetPenghuluScheduleForAssignment)

		// ==================== STAFF MANAGEMENT ====================
		simnikahRoutes.GET("/staff", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), staffHandler.ListStaff)
		simnikahRoutes.PUT("/staff/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), staffHandler.UpdateStaff)
		simnikahRoutes.POST("/staff/pendaftaran", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), staffHandler.CreateRegistrationForUser)
		simnikahRoutes.PUT("/pendaftaran/:id/update-status", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "penghulu", "kepala_kua"), staffHandler.UpdateRegistrationStatus)

		// ==================== LOCATION & GEOCODING ====================
		simnikahRoutes.POST("/location/geocode", middleware.AuthMiddleware(), catinHandler.GetCoordinatesFromAddressEndpoint)
		simnikahRoutes.POST("/location/reverse-geocode", middleware.AuthMiddleware(), catinHandler.GetAddressFromCoordinates)
		simnikahRoutes.GET("/location/search", middleware.AuthMiddleware(), catinHandler.SearchAddress)
		simnikahRoutes.PUT("/pendaftaran/:id/location", middleware.AuthMiddleware(), catinHandler.UpdateWeddingLocationWithCoordinates)
		simnikahRoutes.GET("/pendaftaran/:id/location", middleware.AuthMiddleware(), catinHandler.GetWeddingLocationDetail)

		// ==================== NOTIFICATION ROUTES ====================
		simnikahRoutes.POST("/notifikasi", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), notificationHandler.CreateNotification)
		simnikahRoutes.GET("/notifikasi/user/me", middleware.AuthMiddleware(), notificationHandler.GetUserNotifications)
		simnikahRoutes.GET("/notifikasi/:id", middleware.AuthMiddleware(), notificationHandler.GetNotificationByID)
		simnikahRoutes.PUT("/notifikasi/:id/status", middleware.AuthMiddleware(), notificationHandler.UpdateNotificationStatus)
		simnikahRoutes.PUT("/notifikasi/mark-all-read", middleware.AuthMiddleware(), notificationHandler.MarkAllAsRead)
		simnikahRoutes.DELETE("/notifikasi/:id", middleware.AuthMiddleware(), notificationHandler.DeleteNotification)
		simnikahRoutes.GET("/notifikasi/stats", middleware.AuthMiddleware(), notificationHandler.GetNotificationStats)
		simnikahRoutes.POST("/notifikasi/send-to-role", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), notificationHandler.SendNotificationToRole)
		simnikahRoutes.POST("/notifikasi/run-reminder", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), RunReminderNotification)

		// ==================== DASHBOARD ROUTES ====================
		simnikahRoutes.GET("/dashboard/kepala-kua", middleware.AuthMiddleware(), middleware.RoleMiddleware("kepala_kua"), dashboardHandler.GetKepalaKUADashboard)
		simnikahRoutes.GET("/dashboard/staff", middleware.AuthMiddleware(), middleware.RoleMiddleware("staff"), dashboardHandler.GetStaffDashboard)
		simnikahRoutes.GET("/dashboard/statistik-pernikahan", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), dashboardHandler.GetMarriageStatistics)
		simnikahRoutes.GET("/dashboard/penghulu-performance", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), dashboardHandler.GetPenghuluPerformance)
		simnikahRoutes.GET("/dashboard/peak-hours", middleware.AuthMiddleware(), middleware.MultiRoleMiddleware("staff", "kepala_kua"), dashboardHandler.GetPeakHoursAnalysis)
	}

	// Custom 404 and 405 handlers for JSON responses
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Endpoint tidak ditemukan",
			"error":   fmt.Sprintf("Path '%s' tidak ditemukan", c.Request.URL.Path),
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"success": false,
			"message": "Method tidak diizinkan",
			"error":   fmt.Sprintf("Method '%s' tidak diizinkan untuk path ini", c.Request.Method),
		})
	})

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
		log.Printf("Server starting on port %s", port)
		log.Printf("SIPENA - Forward Chaining Scheduling Engine")
		log.Printf("4-Stage Flow: Catin Check -> FC Recommend -> Kepala KUA Assign -> Penghulu View")
		log.Printf("Environment: %s", os.Getenv("GIN_MODE"))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	cronJobService.StopReminderCronJob()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}

// createGinEngine creates Gin engine with environment-aware logging
func createGinEngine(ginMode string) *gin.Engine {
	if ginMode == "release" {
		gin.DisableConsoleColor()
		r := gin.New()
		r.Use(gin.Recovery())
		r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/health"},
			Formatter: func(param gin.LogFormatterParams) string {
				if param.StatusCode >= 400 || param.Latency > 2*time.Second {
					return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %#v\n%s",
						param.TimeStamp.Format("2006/01/02 - 15:04:05"),
						param.StatusCode,
						param.Latency,
						param.ClientIP,
						param.Method,
						param.Path,
						param.ErrorMessage,
					)
				}
				return ""
			},
		}))
		return r
	}
	return gin.Default()
}

// getAllowedOrigins returns allowed origins for CORS from environment or uses defaults
func getAllowedOrigins() []string {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Println("Info: Using default CORS origins. Set ALLOWED_ORIGINS environment variable for production.")
		return []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"https://kua-ku.vercel.app",
		}
	}

	origins := strings.Split(allowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	log.Printf("Info: CORS allowed origins: %v", origins)
	return origins
}

// RunReminderNotification menjalankan pengingat notifikasi secara manual (untuk testing)
func RunReminderNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	notificationService := services.NewNotificationService(DB)
	err := notificationService.SendReminderNotification()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menjalankan pengingat notifikasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Pengingat notifikasi berhasil dijalankan",
		"executed_by": userID,
		"executed_at": time.Now(),
	})
}