package notification

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"simnikah/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// CreateNotificationRequest untuk input notifikasi baru
type CreateNotificationRequest struct {
	User_id string `json:"user_id" binding:"required"`
	Judul   string `json:"judul" binding:"required"`
	Pesan   string `json:"pesan" binding:"required"`
	Tipe    string `json:"tipe" binding:"required"`
	Link    string `json:"link"`
}

// UpdateNotificationRequest untuk update notifikasi
type UpdateNotificationRequest struct {
	Status_baca string `json:"status_baca" binding:"required"`
}

// NotificationResponse untuk response notifikasi
type NotificationResponse struct {
	ID          uint      `json:"id"`
	User_id     string    `json:"user_id"`
	Judul       string    `json:"judul"`
	Pesan       string    `json:"pesan"`
	Tipe        string    `json:"tipe"`
	Status_baca string    `json:"status_baca"`
	Link        string    `json:"link"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

// CreateNotification membuat notifikasi baru
func (h *InDB) CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid: " + err.Error()})
		return
	}

	// Validasi tipe notifikasi
	validTypes := map[string]bool{
		structs.NotifikasiTipeInfo:    true,
		structs.NotifikasiTipeWarning: true,
		structs.NotifikasiTipeError:   true,
		structs.NotifikasiTipeSuccess: true,
	}
	if !validTypes[req.Tipe] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe notifikasi tidak valid. Tipe yang tersedia: Info, Warning, Error, Success"})
		return
	}

	// Cek apakah user_id ada
	var user structs.Users
	if err := h.DB.Where("user_id = ?", req.User_id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Buat notifikasi baru
	notification := structs.Notifikasi{
		User_id:     req.User_id,
		Judul:       req.Judul,
		Pesan:       req.Pesan,
		Tipe:        req.Tipe,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        req.Link,
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat notifikasi"})
		return
	}

	// Response
	response := NotificationResponse{
		ID:          notification.ID,
		User_id:     notification.User_id,
		Judul:       notification.Judul,
		Pesan:       notification.Pesan,
		Tipe:        notification.Tipe,
		Status_baca: notification.Status_baca,
		Link:        notification.Link,
		Created_at:  notification.Created_at,
		Updated_at:  notification.Updated_at,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Notifikasi berhasil dibuat",
		"notification": response,
	})
}

// GetUserNotifications mengambil semua notifikasi user
func (h *InDB) GetUserNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID diperlukan"})
		return
	}

	// Cek apakah user ada
	var user structs.Users
	if err := h.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status") // "Belum Dibaca", "Sudah Dibaca", atau kosong untuk semua
	tipe := c.Query("tipe")     // "Info", "Warning", "Error", "Success", atau kosong untuk semua

	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Query notifikasi
	query := h.DB.Where("user_id = ?", userID)

	// Filter berdasarkan status
	if status != "" {
		validStatuses := map[string]bool{
			structs.NotifikasiStatusBelumDibaca: true,
			structs.NotifikasiStatusSudahDibaca: true,
		}
		if validStatuses[status] {
			query = query.Where("status_baca = ?", status)
		}
	}

	// Filter berdasarkan tipe
	if tipe != "" {
		validTypes := map[string]bool{
			structs.NotifikasiTipeInfo:    true,
			structs.NotifikasiTipeWarning: true,
			structs.NotifikasiTipeError:   true,
			structs.NotifikasiTipeSuccess: true,
		}
		if validTypes[tipe] {
			query = query.Where("tipe = ?", tipe)
		}
	}

	// Hitung total
	var total int64
	query.Model(&structs.Notifikasi{}).Count(&total)

	// Ambil data dengan pagination
	var notifications []structs.Notifikasi
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi"})
		return
	}

	// Convert ke response
	var response []NotificationResponse
	for _, notif := range notifications {
		response = append(response, NotificationResponse{
			ID:          notif.ID,
			User_id:     notif.User_id,
			Judul:       notif.Judul,
			Pesan:       notif.Pesan,
			Tipe:        notif.Tipe,
			Status_baca: notif.Status_baca,
			Link:        notif.Link,
			Created_at:  notif.Created_at,
			Updated_at:  notif.Updated_at,
		})
	}

	// Hitung unread count
	var unreadCount int64
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND status_baca = ?", userID, structs.NotifikasiStatusBelumDibaca).Count(&unreadCount)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Notifikasi berhasil diambil",
		"notifications": response,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_page": (total + int64(limit) - 1) / int64(limit),
		},
		"unread_count": unreadCount,
	})
}

// GetNotificationByID mengambil notifikasi berdasarkan ID
func (h *InDB) GetNotificationByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi diperlukan"})
		return
	}

	var notification structs.Notifikasi
	if err := h.DB.Where("id = ?", id).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan"})
		return
	}

	response := NotificationResponse{
		ID:          notification.ID,
		User_id:     notification.User_id,
		Judul:       notification.Judul,
		Pesan:       notification.Pesan,
		Tipe:        notification.Tipe,
		Status_baca: notification.Status_baca,
		Link:        notification.Link,
		Created_at:  notification.Created_at,
		Updated_at:  notification.Updated_at,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Notifikasi berhasil diambil",
		"notification": response,
	})
}

// UpdateNotificationStatus mengupdate status notifikasi
func (h *InDB) UpdateNotificationStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi diperlukan"})
		return
	}

	var req UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid: " + err.Error()})
		return
	}

	// Validasi status
	validStatuses := map[string]bool{
		structs.NotifikasiStatusBelumDibaca: true,
		structs.NotifikasiStatusSudahDibaca: true,
	}
	if !validStatuses[req.Status_baca] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid. Status yang tersedia: Belum Dibaca, Sudah Dibaca"})
		return
	}

	// Cek apakah notifikasi ada
	var notification structs.Notifikasi
	if err := h.DB.Where("id = ?", id).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan"})
		return
	}

	// Update status
	if err := h.DB.Model(&notification).Update("status_baca", req.Status_baca).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status notifikasi"})
		return
	}

	// Ambil data terbaru
	h.DB.Where("id = ?", id).First(&notification)

	response := NotificationResponse{
		ID:          notification.ID,
		User_id:     notification.User_id,
		Judul:       notification.Judul,
		Pesan:       notification.Pesan,
		Tipe:        notification.Tipe,
		Status_baca: notification.Status_baca,
		Link:        notification.Link,
		Created_at:  notification.Created_at,
		Updated_at:  notification.Updated_at,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Status notifikasi berhasil diupdate",
		"notification": response,
	})
}

// MarkAllAsRead menandai semua notifikasi user sebagai sudah dibaca
func (h *InDB) MarkAllAsRead(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID diperlukan"})
		return
	}

	// Cek apakah user ada
	var user structs.Users
	if err := h.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Update semua notifikasi user menjadi sudah dibaca
	result := h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND status_baca = ?", userID, structs.NotifikasiStatusBelumDibaca).Update("status_baca", structs.NotifikasiStatusSudahDibaca)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menandai notifikasi sebagai sudah dibaca"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Semua notifikasi berhasil ditandai sebagai sudah dibaca",
		"updated_count": result.RowsAffected,
	})
}

// DeleteNotification menghapus notifikasi
func (h *InDB) DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi diperlukan"})
		return
	}

	// Cek apakah notifikasi ada
	var notification structs.Notifikasi
	if err := h.DB.Where("id = ?", id).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan"})
		return
	}

	// Hapus notifikasi
	if err := h.DB.Delete(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notifikasi berhasil dihapus",
	})
}

// GetNotificationStats mengambil statistik notifikasi user
func (h *InDB) GetNotificationStats(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID diperlukan"})
		return
	}

	// Cek apakah user ada
	var user structs.Users
	if err := h.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Hitung statistik
	var totalCount, unreadCount, infoCount, warningCount, errorCount, successCount int64

	// Total notifikasi
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ?", userID).Count(&totalCount)

	// Notifikasi belum dibaca
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND status_baca = ?", userID, structs.NotifikasiStatusBelumDibaca).Count(&unreadCount)

	// Notifikasi berdasarkan tipe
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND tipe = ?", userID, structs.NotifikasiTipeInfo).Count(&infoCount)
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND tipe = ?", userID, structs.NotifikasiTipeWarning).Count(&warningCount)
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND tipe = ?", userID, structs.NotifikasiTipeError).Count(&errorCount)
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND tipe = ?", userID, structs.NotifikasiTipeSuccess).Count(&successCount)

	// Notifikasi hari ini
	var todayCount int64
	today := time.Now().Format("2006-01-02")
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&todayCount)

	// Notifikasi minggu ini
	var weekCount int64
	weekStart := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	h.DB.Model(&structs.Notifikasi{}).Where("user_id = ? AND DATE(created_at) >= ?", userID, weekStart).Count(&weekCount)

	c.JSON(http.StatusOK, gin.H{
		"message": "Statistik notifikasi berhasil diambil",
		"stats": gin.H{
			"total":  totalCount,
			"unread": unreadCount,
			"read":   totalCount - unreadCount,
			"by_type": gin.H{
				"info":    infoCount,
				"warning": warningCount,
				"error":   errorCount,
				"success": successCount,
			},
			"today": todayCount,
			"week":  weekCount,
		},
	})
}

// SendNotificationToRole mengirim notifikasi ke semua user dengan role tertentu
func (h *InDB) SendNotificationToRole(c *gin.Context) {
	var req struct {
		Role  string `json:"role" binding:"required"`
		Judul string `json:"judul" binding:"required"`
		Pesan string `json:"pesan" binding:"required"`
		Tipe  string `json:"tipe" binding:"required"`
		Link  string `json:"link"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid: " + err.Error()})
		return
	}

	// Validasi role
	validRoles := map[string]bool{
		structs.UserRoleUserBiasa: true,
		structs.UserRolePenghulu:  true,
		structs.UserRoleStaff:     true,
		structs.UserRoleKepalaKUA: true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid. Role yang tersedia: user_biasa, penghulu, staff, kepala_kua"})
		return
	}

	// Validasi tipe notifikasi
	validTypes := map[string]bool{
		structs.NotifikasiTipeInfo:    true,
		structs.NotifikasiTipeWarning: true,
		structs.NotifikasiTipeError:   true,
		structs.NotifikasiTipeSuccess: true,
	}
	if !validTypes[req.Tipe] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe notifikasi tidak valid. Tipe yang tersedia: Info, Warning, Error, Success"})
		return
	}

	// Ambil semua user dengan role tersebut
	var users []structs.Users
	if err := h.DB.Where("role = ? AND status = ?", req.Role, structs.UserStatusAktif).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada user aktif dengan role " + req.Role})
		return
	}

	// Buat notifikasi untuk setiap user
	var notifications []structs.Notifikasi
	for _, user := range users {
		notification := structs.Notifikasi{
			User_id:     user.User_id,
			Judul:       req.Judul,
			Pesan:       req.Pesan,
			Tipe:        req.Tipe,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        req.Link,
		}
		notifications = append(notifications, notification)
	}

	// Simpan semua notifikasi
	if err := h.DB.Create(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim notifikasi ke semua user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         fmt.Sprintf("Notifikasi berhasil dikirim ke %d user dengan role %s", len(users), req.Role),
		"recipient_count": len(users),
		"role":            req.Role,
	})
}
