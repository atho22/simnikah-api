package catin

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
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

// GetRegistrationStatus checks the status of the current logged-in catin's registration.
func (h *InDB) GetRegistrationStatus(c *gin.Context) {
	pendaftarID, _ := c.Get("user_id")
	var pendaftarIDStr string
	if uid, ok := pendaftarID.(string); ok {
		pendaftarIDStr = uid
	}

	var pendaftaran structs.PendaftaranNikah
	err := h.DB.Where("pendaftar_id = ?", pendaftarIDStr).Order("created_at desc").First(&pendaftaran).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Belum ada pendaftaran",
				"data": gin.H{
					"has_registration": false,
					"can_register":     true,
					"registration":     nil,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil status pendaftaran",
			"error":   err.Error(),
		})
		return
	}

	canRegister := pendaftaran.Status_pendaftaran == structs.StatusPendaftaranDitolak

	var penghuluObj interface{} = nil
	if pendaftaran.Penghulu_id != nil {
		var p structs.Penghulu
		if err := h.DB.Where("id = ?", *pendaftaran.Penghulu_id).First(&p).Error; err == nil {
			penghuluObj = gin.H{
				"id":           p.ID,
				"nama":         p.Nama_lengkap,
				"nama_lengkap": p.Nama_lengkap,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status pendaftaran berhasil diambil",
		"data": gin.H{
			"has_registration": true,
			"can_register":     canRegister,
			"registration": gin.H{
				"id":                 pendaftaran.ID,
				"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
				"status_pendaftaran": pendaftaran.Status_pendaftaran,
				"tanggal_nikah":      pendaftaran.Tanggal_nikah.Format("2006-01-02"),
				"waktu_nikah":        pendaftaran.Waktu_nikah,
				"tempat_nikah":       pendaftaran.Tempat_nikah,
				"alamat_akad":        pendaftaran.Alamat_akad,
				"created_at":         pendaftaran.Created_at,
				"calon_suami": gin.H{
					"nama_lengkap": pendaftaran.Nama_suami,
					"nama_dan_bin": pendaftaran.Nama_suami,
					"nama":         pendaftaran.Nama_suami,
				},
				"calon_istri": gin.H{
					"nama_lengkap": pendaftaran.Nama_istri,
					"nama_dan_binti": pendaftaran.Nama_istri,
					"nama":         pendaftaran.Nama_istri,
				},
				"penghulu": penghuluObj,
			},
		},
	})
}

// GetRegistrationDetail returns details of a single registration.
func (h *InDB) GetRegistrationDetail(c *gin.Context) {
	id := c.Param("id")

	pendaftarID, _ := c.Get("user_id")
	var pendaftarIDStr string
	if uid, ok := pendaftarID.(string); ok {
		pendaftarIDStr = uid
	}

	role, _ := c.Get("role")
	roleStr, _ := role.(string)

	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ?", id).First(&pendaftaran).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pendaftaran tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil detail pendaftaran",
			"error":   err.Error(),
		})
		return
	}

	// Security: user_biasa can only view their own registration
	if roleStr == structs.UserRoleUserBiasa && pendaftaran.Pendaftar_id != pendaftarIDStr {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Anda tidak memiliki akses untuk melihat pendaftaran ini",
		})
		return
	}

	var penghuluObj interface{} = nil
	if pendaftaran.Penghulu_id != nil {
		var p structs.Penghulu
		if err := h.DB.Where("id = ?", *pendaftaran.Penghulu_id).First(&p).Error; err == nil {
			penghuluObj = gin.H{
				"id":           p.ID,
				"nip":          p.NIP,
				"nama_lengkap": p.Nama_lengkap,
				"no_hp":        p.No_hp,
				"email":        p.Email,
				"alamat":       p.Alamat,
				"status":       p.Status,
			}
		}
	}

	var mapsUrl, dirUrl, osmUrl string
	if pendaftaran.Latitude != nil && pendaftaran.Longitude != nil {
		mapsUrl = fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", *pendaftaran.Latitude, *pendaftaran.Longitude)
		dirUrl = fmt.Sprintf("https://www.google.com/maps/dir/?api=1&destination=%f,%f", *pendaftaran.Latitude, *pendaftaran.Longitude)
		osmUrl = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", *pendaftaran.Latitude, *pendaftaran.Longitude)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail pendaftaran berhasil diambil",
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
			"pendaftar_id":       pendaftaran.Pendaftar_id,
			"status_pendaftaran": pendaftaran.Status_pendaftaran,
			"tanggal_nikah":      pendaftaran.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        pendaftaran.Waktu_nikah,
			"tempat_nikah":       pendaftaran.Tempat_nikah,
			"alamat_akad":        pendaftaran.Alamat_akad,
			"latitude":           pendaftaran.Latitude,
			"longitude":          pendaftaran.Longitude,
			"created_at":         pendaftaran.Created_at,
			"updated_at":         pendaftaran.Updated_at,
			"calon_suami": gin.H{
				"nama_lengkap": pendaftaran.Nama_suami,
				"nama_dan_bin": pendaftaran.Nama_suami,
				"umur":         pendaftaran.Umur_suami,
			},
			"calon_istri": gin.H{
				"nama_lengkap": pendaftaran.Nama_istri,
				"nama_dan_binti": pendaftaran.Nama_istri,
				"umur":         pendaftaran.Umur_istri,
			},
			"wali_nikah": gin.H{
				"nama_dan_bin":  "Wali Nikah",
				"hubungan_wali": "Ayah Kandung",
			},
			"penghulu": penghuluObj,
			"location": gin.H{
				"latitude":                   pendaftaran.Latitude,
				"longitude":                  pendaftaran.Longitude,
				"has_coordinates":            pendaftaran.Latitude != nil && pendaftaran.Longitude != nil,
				"google_maps_url":            mapsUrl,
				"google_maps_directions_url": dirUrl,
				"osm_url":                    osmUrl,
			},
		},
	})
}

// ListRegistrations lists all registrations (for staff and kepala KUA).
func (h *InDB) ListRegistrations(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")
	
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	query := h.DB.Model(&structs.PendaftaranNikah{})

	if status != "" {
		query = query.Where("status_pendaftaran = ?", status)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("nama_suami LIKE ? OR nama_istri LIKE ? OR nomor_pendaftaran LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung jumlah pendaftaran",
			"error":   err.Error(),
		})
		return
	}

	var registrations []structs.PendaftaranNikah
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&registrations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil daftar pendaftaran",
			"error":   err.Error(),
		})
		return
	}

	// Fetch penghulu names map to avoid N+1 queries
	var penghuluIDs []uint
	for _, r := range registrations {
		if r.Penghulu_id != nil {
			penghuluIDs = append(penghuluIDs, *r.Penghulu_id)
		}
	}

	penghuluMap := make(map[uint]structs.Penghulu)
	if len(penghuluIDs) > 0 {
		var penghulus []structs.Penghulu
		if err := h.DB.Where("id IN ?", penghuluIDs).Find(&penghulus).Error; err == nil {
			for _, p := range penghulus {
				penghuluMap[p.ID] = p
			}
		}
	}

	// Format output to match frontend expectations
	var formattedRegs []gin.H
	for _, r := range registrations {
		var penghuluObj interface{} = nil
		if r.Penghulu_id != nil {
			if p, ok := penghuluMap[*r.Penghulu_id]; ok {
				penghuluObj = gin.H{
					"id":           p.ID,
					"nama":         p.Nama_lengkap,
					"nama_lengkap": p.Nama_lengkap,
				}
			}
		}

		formattedRegs = append(formattedRegs, gin.H{
			"id":                 r.ID,
			"nomor_pendaftaran": r.Nomor_pendaftaran,
			"status_pendaftaran": r.Status_pendaftaran,
			"tanggal_nikah":      r.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        r.Waktu_nikah,
			"tempat_nikah":       r.Tempat_nikah,
			"alamat_akad":        r.Alamat_akad,
			"created_at":         r.Created_at,
			"calon_suami": gin.H{
				"nama_lengkap": r.Nama_suami,
				"nama_dan_bin": r.Nama_suami,
				"nama":         r.Nama_suami,
			},
			"calon_istri": gin.H{
				"nama_lengkap": r.Nama_istri,
				"nama_dan_binti": r.Nama_istri,
				"nama":         r.Nama_istri,
			},
			"penghulu": penghuluObj,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar pendaftaran berhasil diambil",
		"data": gin.H{
			"registrations": formattedRegs,
			"pagination": gin.H{
				"current_page":  page,
				"per_page":      limit,
				"total_records": totalRecords,
				"total_pages":   totalPages,
			},
		},
	})
}
