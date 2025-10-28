package kepala_kua

import (
	"net/http"
	"time"

	"simnikah/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== PENGHULU ASSIGNMENT ====================

// AssignPenghulu assigns a penghulu to a marriage registration
func (h *InDB) AssignPenghulu(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id from context (kepala KUA who is assigning)
	kepalaKuaID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	var input struct {
		PenghuluID uint   `json:"penghulu_id" binding:"required"`
		Catatan    string `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Check if registration exists
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ?", registrationID).First(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pendaftaran tidak ditemukan",
			"error":   "Pendaftaran dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Check if registration is in correct status for penghulu assignment
	// Note: "Berkas Disetujui" status might be a custom status or should be added to constants
	if pendaftaran.Status_pendaftaran != "Berkas Disetujui" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Berkas Disetujui' untuk ditugaskan ke penghulu",
		})
		return
	}

	// Check if penghulu exists and is active
	var penghulu structs.Penghulu
	if err := h.DB.Where("id = ? AND status = ?", input.PenghuluID, structs.PenghuluStatusAktif).First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Penghulu tidak ditemukan",
			"error":   "Penghulu dengan ID tersebut tidak ditemukan atau tidak aktif",
		})
		return
	}

	// Update registration with penghulu assignment
	pendaftaran.Status_pendaftaran = structs.StatusPendaftaranPenghuluDitugaskan
	pendaftaran.Penghulu_id = &input.PenghuluID
	pendaftaran.Penghulu_assigned_by = kepalaKuaID.(string)
	now := time.Now()
	pendaftaran.Penghulu_assigned_at = &now
	pendaftaran.Catatan = input.Catatan
	pendaftaran.Updated_at = time.Now()

	if err := h.DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal menugaskan penghulu",
		})
		return
	}

	// Create notification for the assigned penghulu
	notification := structs.Notifikasi{
		User_id:     penghulu.User_id,
		Judul:       "Penugasan Nikah Baru",
		Pesan:       "Anda ditugaskan untuk memimpin nikah dengan nomor pendaftaran: " + pendaftaran.Nomor_pendaftaran + ". Silakan periksa berkas dan verifikasi.",
		Tipe:        structs.NotifikasiTipeInfo,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        "/penghulu/pendaftaran/" + registrationID,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		// Log error but don't fail the main operation
	}

	// Create notification for the couple
	coupleNotification := structs.Notifikasi{
		User_id:     pendaftaran.Pendaftar_id,
		Judul:       "Penghulu Ditugaskan",
		Pesan:       "Penghulu telah ditugaskan untuk menikahkan Anda. Penghulu akan memeriksa berkas Anda.",
		Tipe:        structs.NotifikasiTipeInfo,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        "/pendaftaran/" + registrationID,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	if err := h.DB.Create(&coupleNotification).Error; err != nil {
		// Log error but don't fail the main operation
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Penghulu berhasil ditugaskan",
		"data": gin.H{
			"id":                   pendaftaran.ID,
			"nomor_pendaftaran":    pendaftaran.Nomor_pendaftaran,
			"status_pendaftaran":   pendaftaran.Status_pendaftaran,
			"penghulu_id":          pendaftaran.Penghulu_id,
			"penghulu_nama":        penghulu.Nama_lengkap,
			"penghulu_assigned_by": pendaftaran.Penghulu_assigned_by,
			"penghulu_assigned_at": pendaftaran.Penghulu_assigned_at,
			"catatan":              pendaftaran.Catatan,
			"updated_at":           pendaftaran.Updated_at,
		},
	})
}

// GetAvailablePenghulus gets list of available penghulus for assignment
func (h *InDB) GetAvailablePenghulus(c *gin.Context) {
	var penghulus []structs.Penghulu
	if err := h.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data penghulu",
		})
		return
	}

	var penghuluList []gin.H
	for _, p := range penghulus {
		penghuluList = append(penghuluList, gin.H{
			"id":           p.ID,
			"nama_lengkap": p.Nama_lengkap,
			"nip":          p.NIP,
			"no_hp":        p.No_hp,
			"email":        p.Email,
			"alamat":       p.Alamat,
			"jumlah_nikah": p.Jumlah_nikah,
			"rating":       p.Rating,
			"status":       p.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data penghulu berhasil diambil",
		"data":    penghuluList,
	})
}
