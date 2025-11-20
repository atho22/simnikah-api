package catin

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "simnikah/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// Kapasitas per hari untuk nikah di KUA
const (
	MaxWeddingsPerDay = 9  // Maksimal 9 pernikahan di KUA per hari
	StartTime         = 8  // Jam mulai (08:00)
	EndTime           = 16 // Jam selesai (16:00)
)

// TimeSlots represents available time slots in a day
var TimeSlots = []string{
	"08:00", "09:00", "10:00", "11:00", "12:00",
	"13:00", "14:00", "15:00", "16:00",
}

// ==================== UPDATE WEDDING ADDRESS ====================

// UpdateMarriageLocation updates the wedding address for a marriage registration
func (h *InDB) UpdateMarriageLocation(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
			"type":    "authentication",
		})
		return
	}

	var input struct {
		AlamatAkad string `json:"alamat_akad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   "Format data tidak valid: " + err.Error(),
			"type":    "validation",
		})
		return
	}

	// Check if registration exists
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ? AND pendaftar_id = ?", registrationID, userID.(string)).First(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pendaftaran tidak ditemukan",
			"error":   "Pendaftaran dengan ID tersebut tidak ditemukan atau bukan milik Anda",
			"type":    "not_found",
		})
		return
	}

	// Check if wedding location is outside KUA
	if pendaftaran.Tempat_nikah != "Di Luar KUA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Alamat tidak dapat diubah",
			"error":   "Alamat hanya dapat diubah untuk nikah di luar KUA",
			"type":    "validation",
		})
		return
	}

	// Update the wedding address
	if err := h.DB.Model(&pendaftaran).Update("alamat_akad", input.AlamatAkad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengupdate alamat nikah",
			"type":    "database",
		})
		return
	}

	// Get updated registration
	h.DB.Where("id = ?", registrationID).First(&pendaftaran)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alamat nikah berhasil diupdate",
		"data": gin.H{
			"pendaftaran_id":    pendaftaran.ID,
			"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
			"alamat_akad":       pendaftaran.Alamat_akad,
			"tempat_nikah":      pendaftaran.Tempat_nikah,
			"updated_at":        pendaftaran.Updated_at,
		},
	})
}

// ==================== FORM PENDAFTARAN NIKAH SEDERHANA ====================

// CreateRegistration creates a simplified marriage registration form
// Hanya memerlukan data minimal: nama+bin/binti, pendidikan, umur, dan lokasi nikah
func (h *InDB) CreateRegistration(c *gin.Context) {
	var formSederhana structs.DataFormPendaftaranSederhana

	if err := c.ShouldBindJSON(&formSederhana); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   "Format data tidak valid: " + err.Error(),
			"type":    "validation",
		})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
			"type":    "authentication",
		})
		return
	}

	// Validate tempat nikah
	if formSederhana.LokasiNikah.TempatNikah != "Di KUA" && formSederhana.LokasiNikah.TempatNikah != "Di Luar KUA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Tempat nikah harus 'Di KUA' atau 'Di Luar KUA'",
			"field":   "tempat_nikah",
			"type":    "enum",
		})
		return
	}

	// Validasi alamat jika di luar KUA
	if formSederhana.LokasiNikah.TempatNikah == "Di Luar KUA" {
		if strings.TrimSpace(formSederhana.LokasiNikah.AlamatNikah) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   "Alamat nikah wajib diisi untuk nikah di luar KUA",
				"field":   "alamat_nikah",
				"type":    "required",
			})
			return
		}

		// Validasi kelurahan harus dalam lingkup Banjarmasin Utara
		if strings.TrimSpace(formSederhana.LokasiNikah.Kelurahan) != "" {
			validKelurahan := false
			for _, kel := range structs.KelurahanBanjarmasinUtara {
				if strings.EqualFold(kel, formSederhana.LokasiNikah.Kelurahan) {
					validKelurahan = true
					break
				}
			}
			if !validKelurahan {
				c.JSON(http.StatusBadRequest, gin.H{
					"success":            false,
					"message":            "Validasi gagal",
					"error":              "Kelurahan harus dalam lingkup Kecamatan Banjarmasin Utara",
					"field":              "kelurahan",
					"type":               "validation",
					"kelurahan_tersedia": structs.KelurahanBanjarmasinUtara,
				})
				return
			}
		}
	}

	// Parse tanggal nikah
	tanggalNikah, err := time.Parse("2006-01-02", formSederhana.LokasiNikah.TanggalNikah)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Format tanggal nikah tidak valid (YYYY-MM-DD)",
			"field":   "tanggal_nikah",
			"type":    "format",
		})
		return
	}

	// Validate that wedding date is not in the past
	if tanggalNikah.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Tanggal nikah tidak boleh di masa lalu",
			"field":   "tanggal_nikah",
			"type":    "validation",
		})
		return
	}

	// Validate wedding time format
	_, err = time.Parse("15:04", formSederhana.LokasiNikah.WaktuNikah)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Format waktu nikah tidak valid (HH:MM dalam format 24-jam, contoh: 09:00)",
			"field":   "waktu_nikah",
			"type":    "format",
		})
		return
	}

	// Validate age (minimum 19 years old)
	if formSederhana.CalonLakiLaki.Umur < 19 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Umur calon laki-laki minimal 19 tahun",
			"field":   "umur_laki_laki",
			"type":    "validation",
		})
		return
	}

	if formSederhana.CalonPerempuan.Umur < 19 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Umur calon perempuan minimal 19 tahun",
			"field":   "umur_perempuan",
			"type":    "validation",
		})
		return
	}

	// Validate wedding schedule availability
	// Normalize waktu nikah format (ensure HH:MM format)
	waktuNikahNormalized := formSederhana.LokasiNikah.WaktuNikah
	if len(waktuNikahNormalized) > 5 {
		waktuNikahNormalized = waktuNikahNormalized[:5] // Take only HH:MM part
	}

	// Check existing registrations with same date and time (include Draft, exclude only rejected)
	// Draft juga dihitung dalam kuota meskipun belum pasti
	// Cek total pernikahan di jam yang sama (baik di KUA maupun di luar KUA)
	var countTotalRegistrations int64
	err = h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah = ? AND waktu_nikah = ? AND status_pendaftaran NOT IN ?",
			tanggalNikah, waktuNikahNormalized,
			[]string{structs.StatusPendaftaranDitolak}).
		Count(&countTotalRegistrations).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengecek ketersediaan jadwal pernikahan",
			"type":    "database",
		})
		return
	}

	// Cek jumlah pernikahan di KUA di jam yang sama (include Draft)
	var countKUA int64
	err = h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah = ? AND waktu_nikah = ? AND tempat_nikah = ? AND status_pendaftaran NOT IN ?",
			tanggalNikah, waktuNikahNormalized, "Di KUA",
			[]string{structs.StatusPendaftaranDitolak}).
		Count(&countKUA).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengecek ketersediaan jadwal pernikahan di KUA",
			"type":    "database",
		})
		return
	}

	// Validate schedule based on location
	const maxTotalWeddings = 3 // Total maksimal pernikahan per jam = 3 penghulu

	if formSederhana.LokasiNikah.TempatNikah == "Di KUA" {
		// Jika nikah di KUA: tidak boleh ada lebih dari 1 pernikahan di KUA di tanggal dan jam yang sama
		if countKUA >= 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal tidak tersedia",
				"error":   fmt.Sprintf("Jadwal pernikahan di KUA pada tanggal %s pukul %s sudah terisi. Silakan pilih tanggal atau jam lain.", 
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
				"data": gin.H{
					"tanggal_nikah":    tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":       waktuNikahNormalized,
					"tempat_nikah":      formSederhana.LokasiNikah.TempatNikah,
					"jumlah_terisi_kua": countKUA,
					"batas_maksimal_kua": 1,
				},
			})
			return
		}
	} else {
		// Jika nikah di luar KUA: 
		// Total maksimal pernikahan per jam = 3 penghulu
		// - Jika sudah ada 1 di KUA, maka slot luar KUA = 2 (total 3)
		// - Jika belum ada di KUA, maka slot luar KUA = 3 (total 3)
		
		// Hitung jumlah pernikahan di luar KUA di jam yang sama
		countLuarKUA := countTotalRegistrations - countKUA
		
		// Hitung slot tersedia untuk luar KUA
		slotTersediaLuarKUA := maxTotalWeddings - countTotalRegistrations
		
		// Validasi 1: Total sudah mencapai maksimal 3
		if countTotalRegistrations >= maxTotalWeddings {
			var errorMsg string
			if countKUA >= 1 {
				errorMsg = fmt.Sprintf("Jadwal pernikahan pada tanggal %s pukul %s sudah penuh (total 3 pernikahan: 1 di KUA dan 2 di luar KUA). Silakan pilih tanggal atau jam lain.", 
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized)
			} else {
				errorMsg = fmt.Sprintf("Jadwal pernikahan pada tanggal %s pukul %s sudah penuh (total 3 pernikahan di luar KUA). Silakan pilih tanggal atau jam lain.", 
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized)
			}
			
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal tidak tersedia",
				"error":   errorMsg,
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
				"data": gin.H{
					"tanggal_nikah":          tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":             waktuNikahNormalized,
					"tempat_nikah":            formSederhana.LokasiNikah.TempatNikah,
					"total_pernikahan":        countTotalRegistrations,
					"pernikahan_di_kua":       countKUA,
					"pernikahan_di_luar_kua":  countLuarKUA,
					"batas_maksimal_total":    maxTotalWeddings,
					"slot_tersedia_luar_kua":  0,
				},
			})
			return
		}
		
		// Validasi 2: Jika sudah ada 1 di KUA, maka slot luar KUA maksimal 2
		if countKUA >= 1 && countLuarKUA >= 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal tidak tersedia",
				"error":   fmt.Sprintf("Jadwal pernikahan di luar KUA pada tanggal %s pukul %s sudah penuh (sudah ada 1 pernikahan di KUA dan 2 di luar KUA). Silakan pilih tanggal atau jam lain.", 
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
				"data": gin.H{
					"tanggal_nikah":          tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":             waktuNikahNormalized,
					"tempat_nikah":            formSederhana.LokasiNikah.TempatNikah,
					"total_pernikahan":        countTotalRegistrations,
					"pernikahan_di_kua":       countKUA,
					"pernikahan_di_luar_kua":  countLuarKUA,
					"batas_maksimal_total":    maxTotalWeddings,
					"slot_tersedia_luar_kua":  slotTersediaLuarKUA,
				},
			})
			return
		}
	}

	// Check if user already has an active marriage registration
	var existingRegistration structs.PendaftaranNikah
	if err := h.DB.Where("pendaftar_id = ? AND status_pendaftaran NOT IN (?)", userID.(string), []string{structs.StatusPendaftaranSelesai, structs.StatusPendaftaranDitolak}).First(&existingRegistration).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Pendaftaran sudah ada",
			"error":   "Anda sudah memiliki pendaftaran nikah yang masih aktif",
			"field":   "pendaftaran",
			"type":    "duplicate",
			"data": gin.H{
				"existing_registration_id": existingRegistration.ID,
				"status":                   existingRegistration.Status_pendaftaran,
				"nomor_pendaftaran":        existingRegistration.Nomor_pendaftaran,
			},
		})
		return
	}

	// Start transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	createdAt := time.Now()

	// Generate unique user IDs untuk calon pasangan
	userIDStr := userID.(string)
	timestamp := time.Now().Unix()

	// Generate unique IDs
	hashGroom := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_groom_%d", userIDStr, timestamp))))
	hashBride := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_bride_%d", timestamp, timestamp+1))))
	groomUserID := hashGroom[:20]
	brideUserID := hashBride[:20]

	// Generate unique NIK yang lebih pendek (max 16 karakter) dan unik
	// Format: T + hash (15 karakter) = 16 karakter total
	nikGroom := fmt.Sprintf("T%s", hashGroom[:15])
	nikBride := fmt.Sprintf("T%s", hashBride[:15])

	// Calculate tanggal lahir dari umur (approximate)
	now := time.Now()
	tanggalLahirGroom := now.AddDate(-formSederhana.CalonLakiLaki.Umur, 0, 0)
	tanggalLahirBride := now.AddDate(-formSederhana.CalonPerempuan.Umur, 0, 0)

	// Create calon suami dengan data minimal (hanya dari form sederhana)
	calonSuami := structs.CalonPasangan{
		User_id:             groomUserID,
		NIK:                 nikGroom, // Temporary NIK untuk search (max 16 karakter)
		Nama_lengkap:        formSederhana.CalonLakiLaki.NamaDanBin,
		Tanggal_lahir:       tanggalLahirGroom,
		Jenis_kelamin:       "L",
		Pendidikan_terakhir: formSederhana.CalonLakiLaki.PendidikanAkhir,
		Created_at:          createdAt,
		Updated_at:          createdAt,
	}

	if err := tx.Create(&calonSuami).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   fmt.Sprintf("Gagal membuat data calon suami: %v", err),
			"type":    "database",
		})
		return
	}

	// Create calon istri dengan data minimal (hanya dari form sederhana)
	calonIstri := structs.CalonPasangan{
		User_id:             brideUserID,
		NIK:                 nikBride, // Temporary NIK untuk search (max 16 karakter)
		Nama_lengkap:        formSederhana.CalonPerempuan.NamaDanBinti,
		Tanggal_lahir:       tanggalLahirBride,
		Jenis_kelamin:       "P",
		Pendidikan_terakhir: formSederhana.CalonPerempuan.PendidikanAkhir,
		Created_at:          createdAt,
		Updated_at:          createdAt,
	}

	if err := tx.Create(&calonIstri).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   fmt.Sprintf("Gagal membuat data calon istri: %v", err),
			"type":    "database",
		})
		return
	}

	// Generate nomor pendaftaran
	nomorPendaftaran := fmt.Sprintf("NIKAH-%s-%d",
		tanggalNikah.Format("20060102"),
		timestamp%10000)

	// Setup alamat akad
	var alamatAkad string
	var latitude, longitude *float64

	if formSederhana.LokasiNikah.TempatNikah == "Di KUA" {
		alamatAkad = "PH5Q+F8C, Jl. Wira Karya, Pangeran, Kec. Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan 70123"
		lat := -3.291304649442475
		lon := 114.58814746634684
		latitude = &lat
		longitude = &lon
	} else {
		// Di luar KUA - gabungkan alamat, detail alamat, dan kelurahan
		alamatParts := []string{}
		if formSederhana.LokasiNikah.AlamatNikah != "" {
			alamatParts = append(alamatParts, formSederhana.LokasiNikah.AlamatNikah)
		}
		if formSederhana.LokasiNikah.DetailAlamat != "" {
			alamatParts = append(alamatParts, formSederhana.LokasiNikah.DetailAlamat)
		}
		if formSederhana.LokasiNikah.Kelurahan != "" {
			alamatParts = append(alamatParts, "Kelurahan "+formSederhana.LokasiNikah.Kelurahan)
		}
		alamatParts = append(alamatParts, "Kecamatan Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan")
		alamatAkad = strings.Join(alamatParts, ", ")
		// Koordinat akan diisi nanti secara asynchronous
	}

	// Create marriage registration
	pendaftaranNikah := structs.PendaftaranNikah{
		Nomor_pendaftaran:   nomorPendaftaran,
		Pendaftar_id:        userID.(string),
		Calon_suami_id:      fmt.Sprintf("%d", calonSuami.ID),
		Calon_istri_id:      fmt.Sprintf("%d", calonIstri.ID),
		Tanggal_pendaftaran: createdAt,
		Tanggal_nikah:       tanggalNikah,
		Waktu_nikah:         formSederhana.LokasiNikah.WaktuNikah,
		Tempat_nikah:        formSederhana.LokasiNikah.TempatNikah,
		Alamat_akad:         alamatAkad,
		Latitude:            latitude,
		Longitude:           longitude,
		Status_pendaftaran:  structs.StatusPendaftaranDraft, // Mulai dari Draft
		Created_at:          createdAt,
		Updated_at:          createdAt,
	}

	if err := tx.Create(&pendaftaranNikah).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat pendaftaran nikah",
			"type":    "database",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal menyimpan data pendaftaran",
			"type":    "database",
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pendaftaran nikah berhasil dibuat (form sederhana)",
		"data": gin.H{
			"id":                 pendaftaranNikah.ID,
			"nomor_pendaftaran":  nomorPendaftaran,
			"status_pendaftaran": pendaftaranNikah.Status_pendaftaran,
			"tanggal_nikah":      pendaftaranNikah.Tanggal_nikah,
			"waktu_nikah":        pendaftaranNikah.Waktu_nikah,
			"tempat_nikah":       pendaftaranNikah.Tempat_nikah,
			"alamat_akad":        pendaftaranNikah.Alamat_akad,
			"calon_suami": gin.H{
				"nama_dan_bin": formSederhana.CalonLakiLaki.NamaDanBin,
				"pendidikan":   formSederhana.CalonLakiLaki.PendidikanAkhir,
				"umur":         formSederhana.CalonLakiLaki.Umur,
			},
			"calon_istri": gin.H{
				"nama_dan_binti": formSederhana.CalonPerempuan.NamaDanBinti,
				"pendidikan":     formSederhana.CalonPerempuan.PendidikanAkhir,
				"umur":           formSederhana.CalonPerempuan.Umur,
			},
			"catatan": "Data minimal telah disimpan. Untuk melengkapi data, silakan datang ke kantor KUA atau lengkapi melalui website SIMNIKAH.",
		},
	})
}

// GetUserRegistrationStatus gets the current registration status for the authenticated user
func (h *InDB) GetUserRegistrationStatus(c *gin.Context) {
	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
			"type":    "authentication",
		})
		return
	}

	// Check if user has any marriage registration
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("pendaftar_id = ?", userID.(string)).First(&pendaftaran).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "User belum memiliki pendaftaran nikah",
				"data": gin.H{
					"has_registration": false,
					"can_register":     true,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database error",
				"error":   "Gagal mengecek status pendaftaran",
				"type":    "database",
			})
		}
		return
	}

	// Build registration response
	registrationData := gin.H{
		"id":                 pendaftaran.ID,
		"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
		"status_pendaftaran": pendaftaran.Status_pendaftaran,
		"tanggal_nikah":      pendaftaran.Tanggal_nikah,
		"waktu_nikah":        pendaftaran.Waktu_nikah,
		"tempat_nikah":       pendaftaran.Tempat_nikah,
		"alamat_akad":        pendaftaran.Alamat_akad,
		"created_at":         pendaftaran.Created_at,
	}

	// Include penghulu info if assigned (for transparency)
	if pendaftaran.Penghulu_id != nil {
		var penghulu structs.Penghulu
		if err := h.DB.Where("id = ?", *pendaftaran.Penghulu_id).First(&penghulu).Error; err == nil {
			registrationData["penghulu"] = gin.H{
				"id":              penghulu.ID,
				"nama_lengkap":    penghulu.Nama_lengkap,
				"nip":             penghulu.NIP,
				"no_hp":           penghulu.No_hp,
				"email":           penghulu.Email,
				"alamat":          penghulu.Alamat,
				"ditugaskan_pada": pendaftaran.Penghulu_assigned_at,
			}
		}
	}

	// User already has registration
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User sudah memiliki pendaftaran nikah",
		"data": gin.H{
			"has_registration": true,
			"can_register":     false,
			"registration":     registrationData,
		},
	})
}

// ListRegistrations retrieves all marriage registrations with filters and pagination for staff
func (h *InDB) ListRegistrations(c *gin.Context) {
	// Get query parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	status := c.Query("status")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	location := c.Query("location")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Parse pagination parameters
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		limitInt = 10
	}
	offset := (pageInt - 1) * limitInt

	// Build query
	query := h.DB.Model(&structs.PendaftaranNikah{})

	// Apply filters
	if status != "" {
		query = query.Where("status_pendaftaran = ?", status)
	}

	if location != "" {
		query = query.Where("tempat_nikah = ?", location)
	}

	if dateFrom != "" {
		if dateFromParsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("tanggal_nikah >= ?", dateFromParsed)
		}
	}

	if dateTo != "" {
		if dateToParsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			// Add one day to include the entire day
			dateToParsed = dateToParsed.Add(24 * time.Hour)
			query = query.Where("tanggal_nikah < ?", dateToParsed)
		}
	}

	// Always join with calon_pasangans to get names
	query = query.Joins("LEFT JOIN calon_pasangans cs ON pendaftaran_nikahs.calon_suami_id = cs.id").
		Joins("LEFT JOIN calon_pasangans ci ON pendaftaran_nikahs.calon_istri_id = ci.id")

	if search != "" {
		// Search in registration number, groom name, bride name, or NIK
		query = query.Where("pendaftaran_nikahs.nomor_pendaftaran LIKE ? OR cs.nama_lengkap LIKE ? OR ci.nama_lengkap LIKE ? OR cs.nik LIKE ? OR ci.nik LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Apply sorting
	validSortFields := map[string]bool{
		"created_at":         true,
		"tanggal_nikah":      true,
		"status_pendaftaran": true,
		"nomor_pendaftaran":  true,
	}
	if validSortFields[sortBy] {
		if sortOrder == "asc" {
			query = query.Order(fmt.Sprintf("%s ASC", sortBy))
		} else {
			query = query.Order(fmt.Sprintf("%s DESC", sortBy))
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Get total count for pagination
	var total int64
	countQuery := h.DB.Model(&structs.PendaftaranNikah{})

	// Apply same filters to count query
	if status != "" {
		countQuery = countQuery.Where("status_pendaftaran = ?", status)
	}
	if location != "" {
		countQuery = countQuery.Where("tempat_nikah = ?", location)
	}
	if dateFrom != "" {
		if dateFromParsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			countQuery = countQuery.Where("tanggal_nikah >= ?", dateFromParsed)
		}
	}
	if dateTo != "" {
		if dateToParsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			dateToParsed = dateToParsed.Add(24 * time.Hour)
			countQuery = countQuery.Where("tanggal_nikah < ?", dateToParsed)
		}
	}
	if search != "" {
		countQuery = countQuery.Joins("LEFT JOIN calon_pasangans cs ON pendaftaran_nikahs.calon_suami_id = cs.id").
			Joins("LEFT JOIN calon_pasangans ci ON pendaftaran_nikahs.calon_istri_id = ci.id").
			Where("pendaftaran_nikahs.nomor_pendaftaran LIKE ? OR cs.nama_lengkap LIKE ? OR ci.nama_lengkap LIKE ? OR cs.nik LIKE ? OR ci.nik LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	countQuery.Count(&total)

	// Get paginated results with join to get calon suami and istri names
	type RegistrationWithNames struct {
		structs.PendaftaranNikah
		CalonSuamiNama string `gorm:"column:calon_suami_nama"`
		CalonIstriNama string `gorm:"column:calon_istri_nama"`
	}

	var results []RegistrationWithNames
	if err := query.
		Select("pendaftaran_nikahs.*, cs.nama_lengkap as calon_suami_nama, ci.nama_lengkap as calon_istri_nama").
		Offset(offset).Limit(limitInt).
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data pendaftaran",
			"type":    "database",
		})
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limitInt) - 1) / int64(limitInt))
	hasNext := pageInt < totalPages
	hasPrev := pageInt > 1

	// Prepare response data
	var registrations []gin.H
	for _, r := range results {
		registrations = append(registrations, gin.H{
			"id":                  r.ID,
			"nomor_pendaftaran":   r.Nomor_pendaftaran,
			"pendaftar_id":        r.Pendaftar_id,
			"status_pendaftaran":  r.Status_pendaftaran,
			"tanggal_pendaftaran": r.Tanggal_pendaftaran,
			"tanggal_nikah":       r.Tanggal_nikah,
			"waktu_nikah":         r.Waktu_nikah,
			"tempat_nikah":        r.Tempat_nikah,
			"alamat_akad":         r.Alamat_akad,
			"penghulu_id":         r.Penghulu_id,
			"catatan":             r.Catatan,
			"calon_suami": gin.H{
				"id":           r.Calon_suami_id,
				"nama_lengkap": r.CalonSuamiNama,
			},
			"calon_istri": gin.H{
				"id":           r.Calon_istri_id,
				"nama_lengkap": r.CalonIstriNama,
			},
			"created_at": r.Created_at,
			"updated_at": r.Updated_at,
		})
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pendaftaran berhasil diambil",
		"data": gin.H{
			"registrations": registrations,
			"pagination": gin.H{
				"current_page":  pageInt,
				"total_pages":   totalPages,
				"total_records": total,
				"per_page":      limitInt,
				"has_next":      hasNext,
				"has_previous":  hasPrev,
			},
			"filters": gin.H{
				"status":     status,
				"date_from":  dateFrom,
			"date_to":    dateTo,
			"location":   location,
			"search":     search,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	},
	})
}

// ==================== CALENDAR AVAILABILITY ====================

// GetCalendarAvailability returns available and unavailable dates for a specific month
// GET /simnikah/kalender-ketersediaan?bulan=01&tahun=2024
func (h *InDB) GetCalendarAvailability(c *gin.Context) {
	// Get query parameters
	bulanStr := c.DefaultQuery("bulan", "")
	tahunStr := c.DefaultQuery("tahun", "")

	// Use current month and year if not provided
	now := time.Now()
	if bulanStr == "" {
		bulanStr = fmt.Sprintf("%02d", int(now.Month()))
	}
	if tahunStr == "" {
		tahunStr = fmt.Sprintf("%d", now.Year())
	}

	// Parse bulan and tahun
	bulanInt, err := strconv.Atoi(bulanStr)
	if err != nil || bulanInt < 1 || bulanInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bulan tidak valid",
			"error":   "Bulan harus antara 1-12",
		})
		return
	}

	tahunInt, err := strconv.Atoi(tahunStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tahun tidak valid",
			"error":   "Format tahun tidak valid",
		})
		return
	}
	
	// Allow viewing past months in current year, but not past years
	if tahunInt < now.Year() {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tahun tidak valid",
			"error":   "Tahun tidak boleh kurang dari tahun sekarang",
		})
		return
	}
	
	// If viewing current year, allow past months for transparency
	// If viewing future year, that's also allowed

	// Calculate start and end of month
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.UTC)
	akhirBulan := awalBulan.AddDate(0, 1, -1)
	// End of month should include the entire last day (23:59:59.999)
	akhirBulanEnd := akhirBulan.Add(24*time.Hour - time.Second)

	// Query all registrations in this month (both KUA and luar KUA, excluding rejected)
	// Include Draft (kuning - belum pasti) dan Disetujui (hijau - sudah pasti)
	var pendaftaran []structs.PendaftaranNikah
	err = h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah <= ? AND status_pendaftaran NOT IN ?",
		awalBulan, akhirBulanEnd, []string{structs.StatusPendaftaranDitolak}).
		Find(&pendaftaran).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data ketersediaan",
			"error":   err.Error(),
		})
		return
	}

	// Count registrations per date, separate Draft and Disetujui, and by location
	registrationsPerDate := make(map[string]int)        // Total (Draft + Disetujui)
	draftPerDate := make(map[string]int)                 // Draft only (kuning - belum pasti)
	disetujuiPerDate := make(map[string]int)             // Disetujui only (hijau - sudah pasti)
	
	// Count per date and time slot, separate by location
	// Format: tanggalStr -> waktu -> {kua: count, luar_kua: count}
	type TimeSlotCount struct {
		KUA     int // Draft + Disetujui di KUA
		LuarKUA int // Draft + Disetujui di luar KUA
		DraftKUA int // Draft di KUA
		DisetujuiKUA int // Disetujui di KUA
		DraftLuarKUA int // Draft di luar KUA
		DisetujuiLuarKUA int // Disetujui di luar KUA
	}
	timeSlotsPerDate := make(map[string]map[string]*TimeSlotCount) // tanggal -> waktu -> count
	
	for _, p := range pendaftaran {
		tanggalStr := p.Tanggal_nikah.Format("2006-01-02")
		
		// Normalize waktu format
		waktu := p.Waktu_nikah
		if len(waktu) > 5 {
			waktu = waktu[:5] // Take only HH:MM part
		}
		
		// Initialize map if needed
		if timeSlotsPerDate[tanggalStr] == nil {
			timeSlotsPerDate[tanggalStr] = make(map[string]*TimeSlotCount)
		}
		if timeSlotsPerDate[tanggalStr][waktu] == nil {
			timeSlotsPerDate[tanggalStr][waktu] = &TimeSlotCount{}
		}
		
		// Count by location
		if p.Tempat_nikah == "Di KUA" {
			registrationsPerDate[tanggalStr]++
			timeSlotsPerDate[tanggalStr][waktu].KUA++
			
			if p.Status_pendaftaran == structs.StatusPendaftaranDraft {
				draftPerDate[tanggalStr]++
				timeSlotsPerDate[tanggalStr][waktu].DraftKUA++
			} else if p.Status_pendaftaran == structs.StatusPendaftaranDisetujui {
				disetujuiPerDate[tanggalStr]++
				timeSlotsPerDate[tanggalStr][waktu].DisetujuiKUA++
			}
		} else {
			// Di luar KUA
			timeSlotsPerDate[tanggalStr][waktu].LuarKUA++
			
			if p.Status_pendaftaran == structs.StatusPendaftaranDraft {
				timeSlotsPerDate[tanggalStr][waktu].DraftLuarKUA++
			} else if p.Status_pendaftaran == structs.StatusPendaftaranDisetujui {
				timeSlotsPerDate[tanggalStr][waktu].DisetujuiLuarKUA++
			}
		}
	}

	// Build calendar data
	calendar := []map[string]interface{}{}
	for tanggal := 1; tanggal <= akhirBulan.Day(); tanggal++ {
		tanggalTime := time.Date(tahunInt, time.Month(bulanInt), tanggal, 0, 0, 0, 0, time.UTC)
		tanggalStr := tanggalTime.Format("2006-01-02")

		jumlahNikah := registrationsPerDate[tanggalStr]
		jumlahDraft := draftPerDate[tanggalStr]           // Kuning - belum pasti
		jumlahDisetujui := disetujuiPerDate[tanggalStr]   // Hijau - sudah pasti
		
		// Sisa kuota dihitung berdasarkan total (Draft + Disetujui)
		// Draft juga dihitung dalam kuota meskipun belum pasti
		totalNikah := jumlahDraft + jumlahDisetujui
		sisaKuota := MaxWeddingsPerDay - totalNikah

		// Determine status
		var status string
		var tersedia bool
		if tanggalTime.Before(now.Truncate(24 * time.Hour)) {
			status = "Terlewat"
			tersedia = false
		} else if totalNikah >= MaxWeddingsPerDay {
			status = "Penuh"
			tersedia = false
		} else {
			status = "Tersedia"
			tersedia = true
		}

		// Build time slots availability for this date
		timeSlotsData := []map[string]interface{}{}
		dateTimeSlots := timeSlotsPerDate[tanggalStr]
		
		// Check each time slot
		for _, slot := range TimeSlots {
			var slotData *TimeSlotCount
			if dateTimeSlots != nil {
				slotData = dateTimeSlots[slot]
			}
			
			if slotData == nil {
				slotData = &TimeSlotCount{}
			}
			
			// Determine availability
			// KUA: maksimal 1 per jam (Draft + Disetujui dihitung)
			// Luar KUA: maksimal 3 total per jam (Draft + Disetujui dihitung)
			// - Jika sudah ada 1 di KUA, maka slot luar KUA = 2 (total 3)
			// - Jika belum ada di KUA, maka slot luar KUA = 3 (total 3)
			const maxTotalPerHour = 3
			const maxKUAPerHour = 1
			
			// Total di KUA (Draft + Disetujui)
			totalKUA := slotData.KUA
			tersediaKUA := totalKUA < maxKUAPerHour && !tanggalTime.Before(now.Truncate(24 * time.Hour))
			
			// Logika untuk luar KUA sesuai dengan CreateRegistration
			// Total semua (Draft + Disetujui) di KUA dan luar KUA
			totalAll := slotData.KUA + slotData.LuarKUA
			var maxLuarKUAPerHour int
			if totalKUA >= 1 {
				// Sudah ada 1 di KUA, maka slot luar KUA maksimal 2
				maxLuarKUAPerHour = 2
			} else {
				// Belum ada di KUA, maka slot luar KUA maksimal 3
				maxLuarKUAPerHour = 3
			}
			
			// Total di luar KUA (Draft + Disetujui)
			totalLuarKUA := slotData.LuarKUA
			tersediaLuarKUA := totalAll < maxTotalPerHour && 
				totalLuarKUA < maxLuarKUAPerHour && 
				!tanggalTime.Before(now.Truncate(24 * time.Hour))
			
			timeSlotsData = append(timeSlotsData, map[string]interface{}{
				"waktu": slot,
				"kua": gin.H{
					"tersedia":         tersediaKUA,
					"terbooking":       totalKUA >= maxKUAPerHour,
					"jumlah_total":     slotData.KUA,
					"jumlah_draft":     slotData.DraftKUA,
					"jumlah_disetujui": slotData.DisetujuiKUA,
				},
				"luar_kua": gin.H{
					"tersedia":         tersediaLuarKUA,
					"terbooking":       totalAll >= maxTotalPerHour || totalLuarKUA >= maxLuarKUAPerHour,
					"jumlah_total":     slotData.LuarKUA,
					"jumlah_draft":     slotData.DraftLuarKUA,
					"jumlah_disetujui": slotData.DisetujuiLuarKUA,
				},
			})
		}

		calendar = append(calendar, map[string]interface{}{
			"tanggal":          tanggal,
			"tanggal_str":      tanggalStr,
			"hari":             tanggalTime.Weekday().String(),
			"status":            status,
			"tersedia":         tersedia,
			"jumlah_nikah":     jumlahNikah,              // Total (Draft + Disetujui) di KUA
			"jumlah_draft":      jumlahDraft,              // Kuning - belum pasti (di KUA)
			"jumlah_disetujui":  jumlahDisetujui,          // Hijau - sudah pasti (di KUA)
			"sisa_kuota":       sisaKuota,                 // Berdasarkan yang sudah pasti (di KUA)
			"kapasitas":        MaxWeddingsPerDay,
			"is_today":         tanggalStr == now.Format("2006-01-02"),
			"is_past":          tanggalTime.Before(now.Truncate(24 * time.Hour)),
			"time_slots":       timeSlotsData,             // Detail jam-jam tersedia
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kalender ketersediaan berhasil diambil",
		"data": gin.H{
			"bulan":            bulanInt,
			"tahun":            tahunInt,
			"nama_bulan":       awalBulan.Month().String(),
			"kapasitas_harian": MaxWeddingsPerDay,
			"calendar":         calendar,
		},
	})
}

// GetAvailableTimeSlots returns available time slots for a specific date
// GET /simnikah/ketersediaan-jam?tanggal=2024-01-15
func (h *InDB) GetAvailableTimeSlots(c *gin.Context) {
	tanggalStr := c.Query("tanggal")
	if tanggalStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tanggal diperlukan",
			"error":   "Format: YYYY-MM-DD",
			"example": "/simnikah/ketersediaan-jam?tanggal=2024-01-15",
		})
		return
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal tidak valid",
			"error":   "Format harus YYYY-MM-DD",
		})
		return
	}

	// Check if date is in the past
	now := time.Now()
	if tanggal.Before(now.Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tanggal sudah lewat",
			"error":   "Tidak bisa melihat ketersediaan untuk tanggal yang sudah lewat",
		})
		return
	}

	// Query registrations on this date (only for KUA weddings, excluding rejected)
	// Include Draft (kuning - belum pasti) dan Disetujui (hijau - sudah pasti)
	// Calculate start and end of day
	startOfDay := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var pendaftaran []structs.PendaftaranNikah
	err = h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND tempat_nikah = ? AND status_pendaftaran NOT IN ?",
		startOfDay, endOfDay, "Di KUA", []string{structs.StatusPendaftaranDitolak}).
		Order("waktu_nikah ASC").
		Find(&pendaftaran).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data ketersediaan",
			"error":   err.Error(),
		})
		return
	}

	// Count registrations per time slot, separate Draft and Disetujui
	registrationsPerTimeSlot := make(map[string]int)    // Total (Draft + Disetujui)
	draftPerTimeSlot := make(map[string]int)            // Draft only (kuning - belum pasti)
	disetujuiPerTimeSlot := make(map[string]int)        // Disetujui only (hijau - sudah pasti)
	
	for _, p := range pendaftaran {
		// Normalize time format (ensure HH:MM format)
		waktu := p.Waktu_nikah
		if len(waktu) > 5 {
			waktu = waktu[:5] // Take only HH:MM part
		}
		registrationsPerTimeSlot[waktu]++
		
		if p.Status_pendaftaran == structs.StatusPendaftaranDraft {
			draftPerTimeSlot[waktu]++
		} else if p.Status_pendaftaran == structs.StatusPendaftaranDisetujui {
			disetujuiPerTimeSlot[waktu]++
		}
	}

	// Build available time slots
	availableSlots := []map[string]interface{}{}
	totalBooked := 0      // Hanya yang sudah pasti (Disetujui)
	totalAvailable := 0

	// Get today's date for comparison (truncated to start of day)
	today := now.Truncate(24 * time.Hour)
	tanggalStartOfDay := tanggal.Truncate(24 * time.Hour)
	
	for _, slot := range TimeSlots {
		jumlahTotal := registrationsPerTimeSlot[slot]
		jumlahDraft := draftPerTimeSlot[slot]           // Kuning - belum pasti
		jumlahDisetujui := disetujuiPerTimeSlot[slot]    // Hijau - sudah pasti
		
		// Slot dianggap terbooking jika ada Draft atau Disetujui (Draft juga dihitung dalam kuota)
		isBooked := jumlahTotal > 0
		// Use time comparison instead of string comparison for accuracy
		isAvailable := !isBooked && !tanggalStartOfDay.Before(today)

		if isBooked {
			totalBooked++
		} else if isAvailable {
			totalAvailable++
		}

		availableSlots = append(availableSlots, map[string]interface{}{
			"waktu":            slot,
			"tersedia":         isAvailable,
			"terbooking":       isBooked,
			"jumlah_nikah":      jumlahTotal,              // Total (Draft + Disetujui)
			"jumlah_draft":      jumlahDraft,              // Kuning - belum pasti
			"jumlah_disetujui":  jumlahDisetujui,          // Hijau - sudah pasti
		})
	}

	// Get summary - sisa kuota berdasarkan total (Draft + Disetujui)
	totalCapacity := len(TimeSlots)
	sisaKuota := totalCapacity - totalBooked

	var status string
	if totalBooked >= totalCapacity {
		status = "Penuh"
	} else if totalBooked > 0 {
		status = "Sebagian Tersedia"
	} else {
		status = "Semua Tersedia"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ketersediaan jam berhasil diambil",
		"data": gin.H{
			"tanggal":      tanggalStr,
			"hari":         tanggal.Weekday().String(),
			"status":       status,
			"summary": gin.H{
				"total_slot":     totalCapacity,
				"terbooking":     totalBooked,
				"tersedia":       totalAvailable,
				"sisa_kuota":     sisaKuota,
			},
			"time_slots": availableSlots,
			"registrations_today": gin.H{
				"total":  len(pendaftaran),
				"detail": pendaftaran,
			},
		},
	})
}

// GetWeddingsByDate returns detailed wedding information for a specific date
// GET /simnikah/pernikahan-tanggal?tanggal=2024-01-15
func (h *InDB) GetWeddingsByDate(c *gin.Context) {
	tanggalStr := c.Query("tanggal")
	if tanggalStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tanggal diperlukan",
			"error":   "Format: YYYY-MM-DD",
			"example": "/simnikah/pernikahan-tanggal?tanggal=2024-01-15",
		})
		return
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal tidak valid",
			"error":   "Format harus YYYY-MM-DD",
		})
		return
	}

	// Calculate start and end of day
	startOfDay := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query registrations on this date with join to get calon suami, istri, and penghulu names
	type WeddingWithNames struct {
		structs.PendaftaranNikah
		CalonSuamiNama  string `gorm:"column:calon_suami_nama"`
		CalonIstriNama  string `gorm:"column:calon_istri_nama"`
		PenghuluNama    string `gorm:"column:penghulu_nama"`
	}

	var results []WeddingWithNames
	err = h.DB.Model(&structs.PendaftaranNikah{}).
		Select("pendaftaran_nikahs.*, cs.nama_lengkap as calon_suami_nama, ci.nama_lengkap as calon_istri_nama, p.nama_lengkap as penghulu_nama").
		Joins("LEFT JOIN calon_pasangans cs ON pendaftaran_nikahs.calon_suami_id = cs.id").
		Joins("LEFT JOIN calon_pasangans ci ON pendaftaran_nikahs.calon_istri_id = ci.id").
		Joins("LEFT JOIN penghulus p ON pendaftaran_nikahs.penghulu_id = p.id").
		Where("pendaftaran_nikahs.tanggal_nikah >= ? AND pendaftaran_nikahs.tanggal_nikah < ? AND pendaftaran_nikahs.status_pendaftaran NOT IN ?",
			startOfDay, endOfDay, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Order("pendaftaran_nikahs.waktu_nikah ASC").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pernikahan",
			"error":   err.Error(),
		})
		return
	}

	// Prepare response data
	var weddings []gin.H
	nikahDiKUA := 0
	nikahDiLuar := 0

	for _, r := range results {
		// Count by location
		if r.Tempat_nikah == "Di KUA" {
			nikahDiKUA++
		} else {
			nikahDiLuar++
		}

		// Build penghulu info
		penghuluInfo := gin.H{
			"id":   r.Penghulu_id,
			"nama": nil,
		}
		if r.Penghulu_id != nil && r.PenghuluNama != "" {
			penghuluInfo["nama"] = r.PenghuluNama
		}

		weddings = append(weddings, gin.H{
			"id":                 r.ID,
			"nomor_pendaftaran":  r.Nomor_pendaftaran,
			"waktu_nikah":        r.Waktu_nikah,
			"tempat_nikah":       r.Tempat_nikah,
			"alamat_akad":        r.Alamat_akad,
			"status_pendaftaran": r.Status_pendaftaran,
			"penghulu":           penghuluInfo,
			"calon_suami": gin.H{
				"id":           r.Calon_suami_id,
				"nama_lengkap": r.CalonSuamiNama,
			},
			"calon_istri": gin.H{
				"id":           r.Calon_istri_id,
				"nama_lengkap": r.CalonIstriNama,
			},
		})
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pernikahan pada tanggal berhasil diambil",
		"data": gin.H{
			"tanggal":        tanggalStr,
			"hari":           tanggal.Weekday().String(),
			"tanggal_format": tanggal.Format("02 Januari 2006"),
			"summary": gin.H{
				"total_nikah":      len(weddings),
				"nikah_di_kua":     nikahDiKUA,
				"nikah_di_luar":    nikahDiLuar,
			},
			"pernikahan": weddings,
		},
	})
}

// ==================== FEEDBACK PERNIKAHAN ====================

// CreateFeedbackPernikahan creates feedback from catin after marriage is completed
// POST /simnikah/feedback-pernikahan
func (h *InDB) CreateFeedbackPernikahan(c *gin.Context) {
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

	var input struct {
		PendaftaranID uint   `json:"pendaftaran_id" binding:"required"`
		JenisFeedback string `json:"jenis_feedback" binding:"required"` // "Rating", "Saran", "Kritik", "Laporan"
		Rating        *int   `json:"rating"`                            // 1-5, hanya untuk jenis "Rating"
		Judul         string `json:"judul" binding:"required"`
		Pesan         string `json:"pesan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validate jenis feedback
	validJenis := map[string]bool{
		structs.FeedbackJenisRating:  true,
		structs.FeedbackJenisSaran:   true,
		structs.FeedbackJenisKritik:  true,
		structs.FeedbackJenisLaporan: true,
	}
	if !validJenis[input.JenisFeedback] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Jenis feedback tidak valid",
			"error":   "Jenis feedback harus salah satu dari: Rating, Saran, Kritik, Laporan",
		})
		return
	}

	// Validate rating if jenis is "Rating"
	if input.JenisFeedback == structs.FeedbackJenisRating {
		if input.Rating == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Rating wajib diisi untuk jenis Rating",
				"error":   "Rating harus diisi antara 1-5",
			})
			return
		}
		if *input.Rating < 1 || *input.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Rating tidak valid",
				"error":   "Rating harus antara 1-5",
			})
			return
		}
	}

	// Check if registration exists and belongs to user
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ? AND pendaftar_id = ?", input.PendaftaranID, userID.(string)).First(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pendaftaran tidak ditemukan",
			"error":   "Pendaftaran dengan ID tersebut tidak ditemukan atau tidak milik Anda",
		})
		return
	}

	// Check if marriage is completed
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranSelesai {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Pernikahan belum selesai",
			"error":   "Feedback hanya bisa diberikan setelah pernikahan selesai",
		})
		return
	}

	// Check if feedback already exists for this registration
	var existingFeedback structs.FeedbackPernikahan
	if err := h.DB.Where("pendaftaran_id = ? AND user_id = ?", input.PendaftaranID, userID.(string)).First(&existingFeedback).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Feedback sudah pernah diberikan",
			"error":   "Anda sudah memberikan feedback untuk pernikahan ini",
		})
		return
	}

	// Create feedback
	feedback := structs.FeedbackPernikahan{
		Pendaftaran_id: input.PendaftaranID,
		User_id:        userID.(string),
		Jenis_feedback: input.JenisFeedback,
		Rating:         input.Rating,
		Judul:          input.Judul,
		Pesan:          input.Pesan,
		Status_baca:    structs.FeedbackStatusBelumDibaca,
		Created_at:     time.Now(),
		Updated_at:     time.Now(),
	}

	if err := h.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat feedback",
			"error":   "Database error",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Feedback berhasil dikirim",
		"data": gin.H{
			"id":              feedback.ID,
			"pendaftaran_id":  feedback.Pendaftaran_id,
			"jenis_feedback":  feedback.Jenis_feedback,
			"rating":          feedback.Rating,
			"judul":           feedback.Judul,
			"status_baca":     feedback.Status_baca,
			"created_at":      feedback.Created_at,
		},
	})
}
