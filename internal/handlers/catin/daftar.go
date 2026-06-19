package catin

import (
	"net/http"
	"strings"
	"time"

	structs "simnikah/internal/models"
	services "simnikah/internal/services"
	"simnikah/pkg/cache"
	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== CHECK SCHEDULE AVAILABILITY ====================

type checkScheduleAvailabilityRequest struct {
	TanggalNikah string   `json:"tanggal_nikah" binding:"required"`
	WaktuNikah   string   `json:"waktu_nikah"`
	TempatNikah  string   `json:"tempat_nikah"`
	AlamatNikah  string   `json:"alamat_nikah"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

// CheckScheduleAvailability mengecek apakah slot nikah masih tersedia.
// Catin memanggil endpoint ini sebelum mengirim pendaftaran.
// Delegasikan perhitungan ke Forward Chaining Engine.
func (h *InDB) CheckScheduleAvailability(c *gin.Context) {
	var input checkScheduleAvailabilityRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	engine := services.NewForwardChainingEngine(h.DB)

	result, err := engine.CheckScheduleAvailability(services.ScheduleCheckInput{
		TanggalNikah: input.TanggalNikah,
		WaktuNikah:   input.WaktuNikah,
		TempatNikah:  input.TempatNikah,
		AlamatNikah:  input.AlamatNikah,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengecek jadwal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengecekan jadwal selesai",
		"data":    result,
	})
}

// ==================== CREATE REGISTRATION (SCHEDULING-ONLY) ====================

// createRegistrationRequest input untuk pendaftaran nikah.
// Data pasangan (nama/umur) adalah referensi display dari Excel Kepala KUA.
// Data scheduling (tanggal/waktu/tempat/koordinat) adalah variabel keputusan FC.
type createRegistrationRequest struct {
	NamaSuami    string   `json:"nama_suami" binding:"required"`    // Nama calon suami (data referensi)
	UmurSuami    int      `json:"umur_suami"`                      // Umur calon suami (data referensi)
	NamaIstri    string   `json:"nama_istri" binding:"required"`    // Nama calon istri (data referensi)
	UmurIstri    int      `json:"umur_istri"`                      // Umur calon istri (data referensi)
	TanggalNikah string   `json:"tanggal_nikah" binding:"required"` // Format: YYYY-MM-DD (variabel FC)
	WaktuNikah   string   `json:"waktu_nikah" binding:"required"`   // Format: HH:MM (variabel FC)
	TempatNikah  string   `json:"tempat_nikah" binding:"required"`  // "Di KUA" atau "Di Luar KUA" (variabel FC)
	AlamatAkad   string   `json:"alamat_akad"`                      // Wajib jika Di Luar KUA (variabel FC)
	Latitude     *float64 `json:"latitude"`                         // Opsional, auto-geocode jika kosong (variabel FC)
	Longitude    *float64 `json:"longitude"`                        // Opsional, auto-geocode jika kosong (variabel FC)
}

// CreateRegistration membuat pendaftaran nikah baru dengan data scheduling-only.
// Status langsung "Menunggu Penugasan" (tidak ada alur Draft/Disetujui SIMKAH).
func (h *InDB) CreateRegistration(c *gin.Context) {
	var input createRegistrationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validate tempat nikah
	if input.TempatNikah != structs.TempatNikahDiKUA && input.TempatNikah != structs.TempatNikahDiLuarKUA {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "tempat_nikah harus 'Di KUA' atau 'Di Luar KUA'",
		})
		return
	}

	// Alamat wajib jika di luar KUA
	if input.TempatNikah == structs.TempatNikahDiLuarKUA && strings.TrimSpace(input.AlamatAkad) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "alamat_akad wajib diisi untuk nikah di luar KUA",
		})
		return
	}

	// Parse dan validasi tanggal nikah (WITA)
	tanggalNikah, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(input.TanggalNikah), utils.WITA)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal tidak valid",
			"error":   "tanggal_nikah harus berformat YYYY-MM-DD",
		})
		return
	}

	// Tanggal tidak boleh di masa lalu (bandingkan dalam WITA)
	today := time.Now().In(utils.WITA).Truncate(24 * time.Hour)
	if tanggalNikah.Before(today) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tanggal tidak valid",
			"error":   "tanggal_nikah tidak boleh di masa lalu",
		})
		return
	}

	// Validasi format waktu (WITA)
	waktuNikah := strings.TrimSpace(input.WaktuNikah)
	if _, err := time.ParseInLocation("15:04", waktuNikah, utils.WITA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format waktu tidak valid",
			"error":   "waktu_nikah harus berformat HH:MM (contoh: 09:00)",
		})
		return
	}

	// Cek ketersediaan jadwal via Forward Chaining Engine
	engine := services.NewForwardChainingEngine(h.DB)
	scheduleResult, err := engine.CheckScheduleAvailability(services.ScheduleCheckInput{
		TanggalNikah: input.TanggalNikah,
		WaktuNikah:   waktuNikah,
		TempatNikah:  input.TempatNikah,
		AlamatNikah:  input.AlamatAkad,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengecek ketersediaan jadwal",
			"error":   err.Error(),
		})
		return
	}

	if !scheduleResult.Available {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Jadwal tidak tersedia",
			"error":   scheduleResult.Reason,
			"data": gin.H{
				"available":      false,
				"total_booked":   scheduleResult.TotalBooked,
				"slot_remaining": scheduleResult.SlotRemaining,
			},
		})
		return
	}

	// Auto-geocode jika alamat ada tapi koordinat belum disediakan
	var latitude, longitude *float64
	if input.Latitude != nil && input.Longitude != nil {
		latitude = input.Latitude
		longitude = input.Longitude
	} else if input.TempatNikah == structs.TempatNikahDiLuarKUA && strings.TrimSpace(input.AlamatAkad) != "" {
		geoLat, geoLon, geoErr := cache.GetCoordinatesFromAddressCached(input.AlamatAkad)
		if geoErr == nil {
			latitude = &geoLat
			longitude = &geoLon
		}
	}

	// Get pendaftar (catin) user_id from JWT context
	pendaftarID, _ := c.Get("user_id")
	var pendaftarIDStr string
	if uid, ok := pendaftarID.(string); ok {
		pendaftarIDStr = uid
	}

	// Buat pendaftaran
	registration := structs.PendaftaranNikah{
		Nama_suami:         strings.TrimSpace(input.NamaSuami),
		Umur_suami:         input.UmurSuami,
		Nama_istri:         strings.TrimSpace(input.NamaIstri),
		Umur_istri:         input.UmurIstri,
		Tanggal_nikah:      tanggalNikah,
		Waktu_nikah:        waktuNikah,
		Tempat_nikah:       input.TempatNikah,
		Alamat_akad:        input.AlamatAkad,
		Latitude:           latitude,
		Longitude:          longitude,
		Status_pendaftaran: structs.StatusPendaftaranMenungguPenugasan,
		Pendaftar_id:       pendaftarIDStr,
		Created_at:         time.Now().In(utils.WITA),
		Updated_at:         time.Now().In(utils.WITA),
	}

	if err := h.DB.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat pendaftaran",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pendaftaran berhasil dibuat, menunggu penugasan penghulu",
		"data": gin.H{
			"id":                 registration.ID,
			"nama_suami":         registration.Nama_suami,
			"umur_suami":         registration.Umur_suami,
			"nama_istri":         registration.Nama_istri,
			"umur_istri":         registration.Umur_istri,
			"tanggal_nikah":      registration.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        registration.Waktu_nikah,
			"tempat_nikah":       registration.Tempat_nikah,
			"alamat_akad":        registration.Alamat_akad,
			"latitude":           registration.Latitude,
			"longitude":          registration.Longitude,
			"status_pendaftaran": registration.Status_pendaftaran,
		},
	})
}
