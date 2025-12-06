package auth

import (
	"fmt"
	"io"
	"net/http"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/crypto"
	"simnikah/pkg/storage"
	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// RegisterRequest struct for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nama     string `json:"nama" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// LoginRequest struct for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ==================== REGISTER ====================

// Register handles user registration
func (h *InDB) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validate role
	validRoles := map[string]bool{
		structs.UserRoleUserBiasa: true,
		structs.UserRolePenghulu:  true,
		structs.UserRoleStaff:     true,
		structs.UserRoleKepalaKUA: true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Role tidak valid",
			"error":   "Role harus salah satu dari: user_biasa, penghulu, staff, kepala_kua",
		})
		return
	}

	// Check if username already exists
	var existingUser structs.Users
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username sudah digunakan",
			"error":   "Username sudah terdaftar",
		})
		return
	}

	// Check if email already exists
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email sudah digunakan",
			"error":   "Email sudah terdaftar",
		})
		return
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengenkripsi password",
			"error":   err.Error(),
		})
		return
	}

	// Generate user_id
	userID := fmt.Sprintf("USR%d", time.Now().Unix())

	// Create user
	user := structs.Users{
		User_id:    userID,
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       req.Role,
		Status:     structs.UserStatusAktif,
		Nama:       req.Nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User berhasil dibuat",
		"data": gin.H{
			"user_id":  user.User_id,
			"username": user.Username,
			"email":    user.Email,
			"nama":     user.Nama,
			"role":     user.Role,
		},
	})
}

// ==================== LOGIN ====================

// Login handles user login and returns JWT token
func (h *InDB) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Find user by username
	var user structs.Users
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Username atau password salah",
			"error":   "User tidak ditemukan",
		})
		return
	}

	// Check if user is active
	if user.Status != structs.UserStatusAktif {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Akun tidak aktif",
			"error":   "Akun Anda telah dinonaktifkan",
		})
		return
	}

	// Verify password
	if err := crypto.VerifyPassword(req.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Username atau password salah",
			"error":   "Password tidak valid",
		})
		return
	}

	// Generate JWT token
	claims := utils.TokenClaims{
		UserID: user.User_id,
		Email:  user.Email,
		Role:   user.Role,
		Nama:   user.Nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := utils.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"user_id":  user.User_id,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"nama":     user.Nama,
		},
		// Juga include di data untuk backward compatibility
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"user_id":  user.User_id,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
				"nama":     user.Nama,
			},
		},
	})
}

// ==================== GET PROFILE ====================

// GetProfile gets the current user's profile
func (h *InDB) GetProfile(c *gin.Context) {
	// Get user_id from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	// Find user
	var user structs.Users
	if err := h.DB.Where("user_id = ?", userID.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
			"error":   "User dengan ID tersebut tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": gin.H{
			"user_id":       user.User_id,
			"username":      user.Username,
			"email":         user.Email,
			"role":          user.Role,
			"nama":          user.Nama,
			"status":        user.Status,
			"profile_photo": user.Profile_photo,
		},
	})
}

// ==================== UPLOAD PROFILE PHOTO ====================

// UploadProfilePhoto handles profile photo upload
func (h *InDB) UploadProfilePhoto(c *gin.Context) {
	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	// Get file from request
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File tidak ditemukan",
			"error":   "Pastikan field 'photo' ada dalam form",
		})
		return
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Ukuran file terlalu besar",
			"error":   "Maksimal 5MB",
		})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
		"image/webp": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tipe file tidak didukung",
			"error":   "Gunakan JPG, PNG, atau WebP",
		})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal membaca file",
			"error":   err.Error(),
		})
		return
	}
	defer src.Close()

	// Upload to ImgBB
	photoURL, err := storage.UploadFileFromMultipart(src, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal upload foto",
			"error":   err.Error(),
		})
		return
	}

	// Update user profile photo in database
	user := structs.Users{}
	if err := h.DB.Where("user_id = ?", userID.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
			"error":   "User dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Save new photo URL
	if err := h.DB.Model(&user).Update("profile_photo", photoURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan foto ke database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Foto profil berhasil diupload",
		"data": gin.H{
			"profile_photo": photoURL,
			"user_id":       user.User_id,
			"username":      user.Username,
		},
	})
}

