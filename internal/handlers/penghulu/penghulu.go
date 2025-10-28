package penghulu

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

// ==================== DOCUMENT VERIFICATION ====================

// VerifyDocuments verifies documents for a marriage registration
func (h *InDB) VerifyDocuments(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id from context (penghulu who is verifying)
	penghuluID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	var input struct {
		Status  string `json:"status" binding:"required"` // "Menunggu Pelaksanaan" or "Ditolak"
		Catatan string `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validate status
	if input.Status != "Menunggu Pelaksanaan" && input.Status != "Ditolak" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak valid",
			"error":   "Status harus 'Menunggu Pelaksanaan' atau 'Ditolak'",
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

	// Check if registration is in correct status for document verification
	if pendaftaran.Status_pendaftaran != "Menunggu Verifikasi Penghulu" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Menunggu Verifikasi Penghulu' untuk diverifikasi",
		})
		return
	}

	// Check if this penghulu is assigned to this registration
	if pendaftaran.Penghulu_id == nil || *pendaftaran.Penghulu_id == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Anda tidak ditugaskan untuk pendaftaran ini",
		})
		return
	}

	// Get penghulu info to verify assignment
	var penghulu structs.Penghulu
	if err := h.DB.Where("user_id = ?", penghuluID.(string)).First(&penghulu).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Data penghulu tidak ditemukan",
		})
		return
	}

	if penghulu.ID != *pendaftaran.Penghulu_id {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Anda tidak ditugaskan untuk pendaftaran ini",
		})
		return
	}

	// Update registration status
	// Jika disetujui, ubah ke "Menunggu Bimbingan"
	if input.Status == "Menunggu Pelaksanaan" {
		pendaftaran.Status_pendaftaran = "Menunggu Bimbingan"
	} else {
		pendaftaran.Status_pendaftaran = input.Status
	}
	pendaftaran.Catatan = input.Catatan
	pendaftaran.Updated_at = time.Now()

	if err := h.DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengupdate status pendaftaran",
		})
		return
	}

	// Create notification for the couple
	var notification structs.Notifikasi
	if input.Status == "Menunggu Pelaksanaan" {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Berkas Disetujui - Siap Bimbingan",
			Pesan:       "Berkas Anda telah disetujui oleh penghulu. Sekarang Anda dapat mendaftar bimbingan perkawinan yang dilaksanakan setiap hari Rabu.",
			Tipe:        "Success",
			Status_baca: "Belum Dibaca",
			Link:        "/bimbingan",
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
	} else {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Berkas Ditolak",
			Pesan:       "Berkas Anda ditolak oleh penghulu. " + input.Catatan,
			Tipe:        "Error",
			Status_baca: "Belum Dibaca",
			Link:        "/pendaftaran/" + registrationID,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		// Log error but don't fail the main operation
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verifikasi berkas berhasil",
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
			"status_pendaftaran": pendaftaran.Status_pendaftaran,
			"penghulu_id":        pendaftaran.Penghulu_id,
			"catatan":            pendaftaran.Catatan,
			"updated_at":         pendaftaran.Updated_at,
		},
	})
}

// GetAssignedRegistrations gets marriage registrations assigned to this penghulu
func (h *InDB) GetAssignedRegistrations(c *gin.Context) {
	// Get user_id from context
	penghuluID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	// Get penghulu info
	var penghulu structs.Penghulu
	if err := h.DB.Where("user_id = ?", penghuluID.(string)).First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data penghulu tidak ditemukan",
			"error":   "Penghulu dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Get assigned registrations
	var pendaftarans []structs.PendaftaranNikah
	if err := h.DB.Where("penghulu_id = ?", penghulu.ID).Find(&pendaftarans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data pendaftaran",
		})
		return
	}

	var registrations []gin.H
	for _, p := range pendaftarans {
		registrations = append(registrations, gin.H{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah,
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"alamat_akad":        p.Alamat_akad,
			"catatan":            p.Catatan,
			"created_at":         p.Created_at,
			"updated_at":         p.Updated_at,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pendaftaran berhasil diambil",
		"data": gin.H{
			"penghulu":      penghulu.Nama_lengkap,
			"registrations": registrations,
			"total":         len(registrations),
		},
	})
}
