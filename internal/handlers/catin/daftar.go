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
	if err != nil || tahunInt < now.Year() {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tahun tidak valid",
			"error":   "Tahun tidak boleh kurang dari tahun sekarang",
		})
		return
	}

	// Calculate start and end of month
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.UTC)
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	// Query all registrations in this month (only for KUA weddings, excluding rejected/draft)
	var pendaftaran []structs.PendaftaranNikah
	err = h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah <= ? AND tempat_nikah = ? AND status_pendaftaran NOT IN ?",
		awalBulan, akhirBulan, "Di KUA", []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Find(&pendaftaran).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data ketersediaan",
			"error":   err.Error(),
		})
		return
	}

	// Count registrations per date
	registrationsPerDate := make(map[string]int)
	for _, p := range pendaftaran {
		tanggalStr := p.Tanggal_nikah.Format("2006-01-02")
		registrationsPerDate[tanggalStr]++
	}

	// Build calendar data
	calendar := []map[string]interface{}{}
	for tanggal := 1; tanggal <= akhirBulan.Day(); tanggal++ {
		tanggalTime := time.Date(tahunInt, time.Month(bulanInt), tanggal, 0, 0, 0, 0, time.UTC)
		tanggalStr := tanggalTime.Format("2006-01-02")

		jumlahNikah := registrationsPerDate[tanggalStr]
		sisaKuota := MaxWeddingsPerDay - jumlahNikah

		// Determine status
		var status string
		var tersedia bool
		if tanggalTime.Before(now.Truncate(24 * time.Hour)) {
			status = "Terlewat"
			tersedia = false
		} else if jumlahNikah >= MaxWeddingsPerDay {
			status = "Penuh"
			tersedia = false
		} else {
			status = "Tersedia"
			tersedia = true
		}

		calendar = append(calendar, map[string]interface{}{
			"tanggal":       tanggal,
			"tanggal_str":   tanggalStr,
			"hari":          tanggalTime.Weekday().String(),
			"status":        status,
			"tersedia":      tersedia,
			"jumlah_nikah":  jumlahNikah,
			"sisa_kuota":    sisaKuota,
			"kapasitas":     MaxWeddingsPerDay,
			"is_today":      tanggalStr == now.Format("2006-01-02"),
			"is_past":       tanggalTime.Before(now.Truncate(24 * time.Hour)),
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

	// Query registrations on this date (only for KUA weddings, excluding rejected/draft)
	// Calculate start and end of day
	startOfDay := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var pendaftaran []structs.PendaftaranNikah
	err = h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND tempat_nikah = ? AND status_pendaftaran NOT IN ?",
		startOfDay, endOfDay, "Di KUA", []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
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

	// Count registrations per time slot
	registrationsPerTimeSlot := make(map[string]int)
	for _, p := range pendaftaran {
		// Normalize time format (ensure HH:MM format)
		waktu := p.Waktu_nikah
		if len(waktu) > 5 {
			waktu = waktu[:5] // Take only HH:MM part
		}
		registrationsPerTimeSlot[waktu]++
	}

	// Build available time slots
	availableSlots := []map[string]interface{}{}
	totalBooked := 0
	totalAvailable := 0

	for _, slot := range TimeSlots {
		isBooked := registrationsPerTimeSlot[slot] > 0
		isAvailable := !isBooked && tanggalStr >= now.Format("2006-01-02")

		if isBooked {
			totalBooked++
		} else if isAvailable {
			totalAvailable++
		}

		availableSlots = append(availableSlots, map[string]interface{}{
			"waktu":         slot,
			"tersedia":      isAvailable,
			"terbooking":    isBooked,
			"jumlah_nikah":  registrationsPerTimeSlot[slot],
		})
	}

	// Get summary
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
