package penghulu

import (
	"fmt"
	"net/http"
	"time"

	structs "simnikah/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== DOCUMENT VERIFICATION ====================

// VerifyRegistrationDocuments verifies documents for a marriage registration assigned to this penghulu
func (h *InDB) VerifyRegistrationDocuments(c *gin.Context) {
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

// ListMyAssignments gets marriage registrations assigned to this penghulu
// Supports both penghulu and kepala_kua roles (kepala KUA can also act as penghulu)
func (h *InDB) ListMyAssignments(c *gin.Context) {
	// Get user_id and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	role, _ := c.Get("role")
	roleStr := role.(string)

	// Get penghulu info - check if user has penghulu record
	var penghulu structs.Penghulu
	err := h.DB.Where("user_id = ?", userID.(string)).First(&penghulu).Error
	
	// If kepala_kua doesn't have penghulu record, try to get or create one
	if err != nil && roleStr == "kepala_kua" {
		// Get user info to create penghulu record
		var user structs.Users
		if err := h.DB.Where("user_id = ?", userID.(string)).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User tidak ditemukan",
				"error":   "User dengan ID tersebut tidak ditemukan",
			})
			return
		}

		// Check if penghulu record exists, if not create one
		penghulu = structs.Penghulu{
			User_id:      user.User_id,
			NIP:          fmt.Sprintf("KUA-%s", user.User_id), // Generate NIP for kepala KUA
			Nama_lengkap: user.Nama,
			Status:       structs.PenghuluStatusAktif,
			Jumlah_nikah: 0,
			Rating:       0,
		}
		
		// Try to create, but if it fails due to duplicate, try to get existing
		if err := h.DB.Create(&penghulu).Error; err != nil {
			// If creation fails, try to get existing record
			if err := h.DB.Where("user_id = ?", userID.(string)).First(&penghulu).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Database error",
					"error":   "Gagal membuat atau mengambil data penghulu",
				})
				return
			}
		}
	} else if err != nil {
		// For regular penghulu, return error if not found
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data penghulu tidak ditemukan",
			"error":   "Penghulu dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Get query parameters for filtering
	todayOnly := c.Query("today") == "true" || c.Query("today") == "1"
	dateFilter := c.Query("date") // Format: YYYY-MM-DD

	// Build query for assigned registrations
	query := h.DB.Where("penghulu_id = ?", penghulu.ID)
	
	// Filter by date if specified
	if todayOnly || dateFilter != "" {
		var targetDate time.Time
		if dateFilter != "" {
			// Parse custom date
			parsedDate, err := time.Parse("2006-01-02", dateFilter)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Format tanggal tidak valid",
					"error":   "Format tanggal harus YYYY-MM-DD",
				})
				return
			}
			targetDate = parsedDate
		} else {
			// Use today's date
			now := time.Now()
			targetDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		}
		
		// Filter by tanggal_nikah (start of day to end of day)
		startOfDay := targetDate
		endOfDay := targetDate.Add(24 * time.Hour).Add(-time.Second)
		query = query.Where("tanggal_nikah >= ? AND tanggal_nikah <= ?", startOfDay, endOfDay)
	}

	// Get assigned registrations with calon suami and istri names
	var pendaftarans []structs.PendaftaranNikah
	if err := query.Find(&pendaftarans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data pendaftaran",
		})
		return
	}

	var registrations []gin.H
	for _, p := range pendaftarans {
		// Get calon suami and istri names
		var calonSuami, calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		regData := gin.H{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah,
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"alamat_akad":        p.Alamat_akad,
			"catatan":            p.Catatan,
			"calon_suami": gin.H{
				"id":           p.Calon_suami_id,
				"nama_lengkap": calonSuami.Nama_lengkap,
			},
			"calon_istri": gin.H{
				"id":           p.Calon_istri_id,
				"nama_lengkap": calonIstri.Nama_lengkap,
			},
			"created_at":         p.Created_at,
			"updated_at":         p.Updated_at,
		}

		// Tambahkan informasi lokasi jika ada koordinat (untuk nikah di luar KUA)
		if p.Latitude != nil && p.Longitude != nil {
			regData["latitude"] = *p.Latitude
			regData["longitude"] = *p.Longitude
			regData["has_coordinates"] = true
			regData["google_maps_url"] = fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", *p.Latitude, *p.Longitude)
			regData["google_maps_directions_url"] = fmt.Sprintf("https://www.google.com/maps/dir/?api=1&destination=%f,%f", *p.Latitude, *p.Longitude)
			regData["waze_url"] = fmt.Sprintf("https://www.waze.com/ul?ll=%f,%f&navigate=yes", *p.Latitude, *p.Longitude)
			regData["osm_url"] = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", *p.Latitude, *p.Longitude)
		} else {
			regData["has_coordinates"] = false
		}

		// Tambahkan flag untuk nikah di luar KUA
		if p.Tempat_nikah == "Di Luar KUA" {
			regData["is_outside_kua"] = true
			regData["note"] = "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
		} else {
			regData["is_outside_kua"] = false
		}

		registrations = append(registrations, regData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pendaftaran berhasil diambil",
		"data": gin.H{
			"penghulu":      penghulu.Nama_lengkap,
			"registrations": registrations,
			"total":         len(registrations),
			"filter": gin.H{
				"today_only": todayOnly,
				"date":       dateFilter,
			},
		},
	})
}

// GetTodaySchedule gets marriage registrations assigned to this penghulu for today
// This is specifically for dashboard display
func (h *InDB) GetTodaySchedule(c *gin.Context) {
	// Get user_id and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	role, _ := c.Get("role")
	roleStr := role.(string)

	// Get penghulu info - check if user has penghulu record
	var penghulu structs.Penghulu
	err := h.DB.Where("user_id = ?", userID.(string)).First(&penghulu).Error
	
	// If kepala_kua doesn't have penghulu record, try to get or create one
	if err != nil && roleStr == "kepala_kua" {
		// Get user info to create penghulu record
		var user structs.Users
		if err := h.DB.Where("user_id = ?", userID.(string)).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User tidak ditemukan",
				"error":   "User dengan ID tersebut tidak ditemukan",
			})
			return
		}

		// Check if penghulu record exists, if not create one
		penghulu = structs.Penghulu{
			User_id:      user.User_id,
			NIP:          fmt.Sprintf("KUA-%s", user.User_id), // Generate NIP for kepala KUA
			Nama_lengkap: user.Nama,
			Status:       structs.PenghuluStatusAktif,
			Jumlah_nikah: 0,
			Rating:       0,
		}
		
		// Try to create, but if it fails due to duplicate, try to get existing
		if err := h.DB.Create(&penghulu).Error; err != nil {
			// If creation fails, try to get existing record
			if err := h.DB.Where("user_id = ?", userID.(string)).First(&penghulu).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Database error",
					"error":   "Gagal membuat atau mengambil data penghulu",
				})
				return
			}
		}
	} else if err != nil {
		// For regular penghulu, return error if not found
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data penghulu tidak ditemukan",
			"error":   "Penghulu dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Get today's date (start and end of day)
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour).Add(-time.Second)

	// Get assigned registrations for today with calon suami and istri names
	var pendaftarans []structs.PendaftaranNikah
	if err := h.DB.Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah <= ?", 
		penghulu.ID, todayStart, todayEnd).Find(&pendaftarans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data pendaftaran",
		})
		return
	}

	var registrations []gin.H
	for _, p := range pendaftarans {
		// Get calon suami and istri names
		var calonSuami, calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		regData := gin.H{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah,
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"alamat_akad":        p.Alamat_akad,
			"catatan":            p.Catatan,
			"calon_suami": gin.H{
				"id":           p.Calon_suami_id,
				"nama_lengkap": calonSuami.Nama_lengkap,
			},
			"calon_istri": gin.H{
				"id":           p.Calon_istri_id,
				"nama_lengkap": calonIstri.Nama_lengkap,
			},
			"created_at":         p.Created_at,
			"updated_at":         p.Updated_at,
		}

		// Tambahkan informasi lokasi jika ada koordinat (untuk nikah di luar KUA)
		if p.Latitude != nil && p.Longitude != nil {
			regData["latitude"] = *p.Latitude
			regData["longitude"] = *p.Longitude
			regData["has_coordinates"] = true
			regData["google_maps_url"] = fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", *p.Latitude, *p.Longitude)
			regData["google_maps_directions_url"] = fmt.Sprintf("https://www.google.com/maps/dir/?api=1&destination=%f,%f", *p.Latitude, *p.Longitude)
			regData["waze_url"] = fmt.Sprintf("https://www.waze.com/ul?ll=%f,%f&navigate=yes", *p.Latitude, *p.Longitude)
			regData["osm_url"] = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", *p.Latitude, *p.Longitude)
		} else {
			regData["has_coordinates"] = false
		}

		// Tambahkan flag untuk nikah di luar KUA
		if p.Tempat_nikah == "Di Luar KUA" {
			regData["is_outside_kua"] = true
			regData["note"] = "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
		} else {
			regData["is_outside_kua"] = false
		}

		registrations = append(registrations, regData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jadwal hari ini berhasil diambil",
		"data": gin.H{
			"penghulu":      penghulu.Nama_lengkap,
			"tanggal":       todayStart.Format("2006-01-02"),
			"registrations": registrations,
			"total":         len(registrations),
		},
	})
}

// ==================== COMPLETE NIKAH (FLOW SEDERHANA) ====================

// CompleteMarriage updates status to Selesai after penghulu conducts the marriage ceremony
// Flow: Penghulu Ditugaskan → Selesai
// Supports both penghulu and kepala_kua roles (kepala KUA can also act as penghulu)
func (h *InDB) CompleteMarriage(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	role, _ := c.Get("role")
	roleStr := role.(string)

	var input struct {
		Catatan string `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Catatan is optional, so don't fail if it's missing
		input.Catatan = ""
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

	// Check if registration is in correct status
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranPenghuluDitugaskan {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Penghulu Ditugaskan' untuk diselesaikan",
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
	// Support both penghulu and kepala_kua roles
	var penghulu structs.Penghulu
	err := h.DB.Where("user_id = ?", userID.(string)).First(&penghulu).Error
	
	// If kepala_kua doesn't have penghulu record, try to get or create one
	if err != nil && roleStr == "kepala_kua" {
		// Get user info to create penghulu record
		var user structs.Users
		if err := h.DB.Where("user_id = ?", userID.(string)).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User tidak ditemukan",
				"error":   "User dengan ID tersebut tidak ditemukan",
			})
			return
		}

		// Check if penghulu record exists, if not create one
		penghulu = structs.Penghulu{
			User_id:      user.User_id,
			NIP:          fmt.Sprintf("KUA-%s", user.User_id), // Generate NIP for kepala KUA
			Nama_lengkap: user.Nama,
			Status:       structs.PenghuluStatusAktif,
			Jumlah_nikah: 0,
			Rating:       0,
		}
		
		// Try to create, but if it fails due to duplicate, try to get existing
		if err := h.DB.Create(&penghulu).Error; err != nil {
			// If creation fails, try to get existing record
			if err := h.DB.Where("user_id = ?", userID.(string)).First(&penghulu).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Database error",
					"error":   "Gagal membuat atau mengambil data penghulu",
				})
				return
			}
		}
	} else if err != nil {
		// For regular penghulu, return error if not found
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Data penghulu tidak ditemukan",
		})
		return
	}

	// Verify that this penghulu is assigned to this registration
	if penghulu.ID != *pendaftaran.Penghulu_id {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Anda tidak ditugaskan untuk pendaftaran ini",
		})
		return
	}

	// Update registration status to Selesai
	pendaftaran.Status_pendaftaran = structs.StatusPendaftaranSelesai
	if input.Catatan != "" {
		pendaftaran.Catatan = input.Catatan
	}
	pendaftaran.Updated_at = time.Now()

	// Update penghulu's jumlah_nikah (increment)
	if err := h.DB.Model(&penghulu).Update("jumlah_nikah", penghulu.Jumlah_nikah+1).Error; err != nil {
		// Log error but don't fail the main operation
	}

	if err := h.DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengupdate status pendaftaran",
		})
		return
	}

	// Create notification for the couple
	notification := structs.Notifikasi{
		User_id:     pendaftaran.Pendaftar_id,
		Judul:       "Pernikahan Selesai",
		Pesan:       "Pernikahan Anda telah selesai dilaksanakan. Selamat menempuh hidup baru!",
		Tipe:        structs.NotifikasiTipeSuccess,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        "/pendaftaran/" + registrationID,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		// Log error but don't fail the main operation
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pernikahan berhasil diselesaikan",
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
