package kepala_kua

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	structs "simnikah/internal/models"
	"simnikah/internal/services"
	"simnikah/pkg/crypto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== PENGHULU ASSIGNMENT ====================

// AssignMarriageOfficer assigns a penghulu (marriage officer) to a marriage registration
func (h *InDB) AssignMarriageOfficer(c *gin.Context) {
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
	// Flow sederhana: Disetujui → Penghulu Ditugaskan
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranDisetujui && pendaftaran.Status_pendaftaran != structs.StatusPendaftaranMenungguPenugasan {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Disetujui' atau 'Menunggu Penugasan' untuk ditugaskan ke penghulu",
		})
		return
	}

	// Check if penghulu exists and is active
	var penghulu structs.Penghulu
	if err := h.DB.Where("id = ? AND status = ?", input.PenghuluID, structs.PenghuluStatusAktif).First(&penghulu).Error; err != nil {
		// If penghulu not found, check if kepala KUA is trying to assign themselves
		// In that case, create a penghulu record for kepala KUA if it doesn't exist
		var existingPenghulu structs.Penghulu
		if err := h.DB.Where("user_id = ?", kepalaKuaID.(string)).First(&existingPenghulu).Error; err == nil {
			// Kepala KUA has a penghulu record, check if they're trying to assign themselves
			if existingPenghulu.ID == input.PenghuluID {
				penghulu = existingPenghulu
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Penghulu tidak ditemukan",
					"error":   "Penghulu dengan ID tersebut tidak ditemukan atau tidak aktif",
				})
				return
			}
		} else {
			// Kepala KUA doesn't have a penghulu record yet
			// Check if they're trying to assign themselves (by checking if the requested ID matches their user_id pattern)
			// For now, we'll allow kepala KUA to assign themselves by creating a record
			var user structs.Users
			if err := h.DB.Where("user_id = ?", kepalaKuaID.(string)).First(&user).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "User tidak ditemukan",
					"error":   "User dengan ID tersebut tidak ditemukan",
				})
				return
			}

			// Create penghulu record for kepala KUA
			newPenghulu := structs.Penghulu{
				User_id:      user.User_id,
				NIP:          fmt.Sprintf("KUA-%s", user.User_id),
				Nama_lengkap: user.Nama,
				Status:       structs.PenghuluStatusAktif,
				Jumlah_nikah: 0,
				Rating:       0,
			}

			if err := h.DB.Create(&newPenghulu).Error; err != nil {
				// If creation fails, try to get existing record
				if err := h.DB.Where("user_id = ?", kepalaKuaID.(string)).First(&penghulu).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "Database error",
						"error":   "Gagal membuat atau mengambil data penghulu",
					})
					return
				}
			} else {
				penghulu = newPenghulu
				// Update input.PenghuluID to match the newly created penghulu record ID
				input.PenghuluID = penghulu.ID
			}
		}
	} else {
		// Penghulu found - verify kepala KUA can assign themselves if user_id matches
		if penghulu.User_id == kepalaKuaID.(string) {
			// Kepala KUA is assigning themselves - this is allowed
		}
	}

	// Validasi: Cek apakah penghulu sudah ada jadwal di tanggal dan jam yang sama
	// Satu penghulu hanya bisa menikahkan 1 pasangan per jam
	tanggalNikah := pendaftaran.Tanggal_nikah
	waktuNikah := pendaftaran.Waktu_nikah
	if len(waktuNikah) > 5 {
		waktuNikah = waktuNikah[:5] // Normalize to HH:MM
	}

	// Calculate start and end of day
	startOfDay := time.Date(tanggalNikah.Year(), tanggalNikah.Month(), tanggalNikah.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var existingAssignment structs.PendaftaranNikah
	err := h.DB.Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND waktu_nikah = ? AND id != ? AND status_pendaftaran NOT IN ?",
		input.PenghuluID, startOfDay, endOfDay, waktuNikah, pendaftaran.ID,
		[]string{structs.StatusPendaftaranDitolak}).
		First(&existingAssignment).Error

	if err == nil {
		// Penghulu sudah ada jadwal di jam yang sama
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Penghulu tidak tersedia",
			"error": fmt.Sprintf("Penghulu %s sudah memiliki jadwal pernikahan pada tanggal %s pukul %s. Satu penghulu hanya bisa menikahkan 1 pasangan per jam.",
				penghulu.Nama_lengkap, tanggalNikah.Format("02 Januari 2006"), waktuNikah),
			"type": "schedule_conflict",
			"data": gin.H{
				"penghulu_id":          input.PenghuluID,
				"penghulu_nama":        penghulu.Nama_lengkap,
				"tanggal_nikah":        tanggalNikah.Format("2006-01-02"),
				"waktu_nikah":          waktuNikah,
				"existing_pendaftaran": existingAssignment.Nomor_pendaftaran,
			},
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

// ListAvailableOfficers gets list of available penghulus (marriage officers) for assignment
func (h *InDB) ListAvailableOfficers(c *gin.Context) {
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

// ==================== STAFF MANAGEMENT ====================

// CreateStaff creates a new staff member (accessible by kepala KUA only)
func (h *InDB) CreateStaff(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nama     string `json:"nama" binding:"required"`
		NIP      string `json:"nip" binding:"required"`
		Jabatan  string `json:"jabatan" binding:"required"`
		Bagian   string `json:"bagian" binding:"required"`
		No_hp    string `json:"no_hp"`
		Alamat   string `json:"alamat"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validate jabatan
	validJabatan := map[string]bool{
		structs.StaffJabatanStaff:     true,
		structs.StaffJabatanPenghulu:  true,
		structs.StaffJabatanKepalaKUA: true,
	}
	if !validJabatan[input.Jabatan] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Jabatan tidak valid",
			"error":   "Jabatan harus salah satu dari: Staff, Penghulu, atau Kepala KUA",
		})
		return
	}

	// Check if username/email already exists
	var existingUser structs.Users
	if err := h.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username atau email sudah digunakan",
			"error":   "Username atau email sudah terdaftar di sistem",
		})
		return
	}

	// Check if NIP already exists
	var existingStaff structs.StaffKUA
	if err := h.DB.Where("nip = ?", input.NIP).First(&existingStaff).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "NIP sudah terdaftar",
			"error":   "NIP sudah terdaftar untuk staff lain",
		})
		return
	}

	// Generate user_id
	userID := "STF" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password
	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengenkripsi password",
		})
		return
	}

	// Determine user role based on jabatan
	var userRole string
	if input.Jabatan == structs.StaffJabatanPenghulu {
		userRole = structs.UserRolePenghulu
	} else if input.Jabatan == structs.StaffJabatanKepalaKUA {
		userRole = structs.UserRoleKepalaKUA
	} else {
		userRole = structs.UserRoleStaff
	}

	// Create user account
	user := structs.Users{
		User_id:    userID,
		Username:   input.Username,
		Email:      input.Email,
		Password:   hashedPassword,
		Role:       userRole,
		Status:     structs.UserStatusAktif,
		Nama:       input.Nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat user account",
		})
		return
	}

	// Create staff profile
	staff := structs.StaffKUA{
		User_id:      userID,
		NIP:          input.NIP,
		Nama_lengkap: input.Nama,
		Jabatan:      input.Jabatan,
		Bagian:       input.Bagian,
		No_hp:        input.No_hp,
		Email:        input.Email,
		Alamat:       input.Alamat,
		Status:       structs.StaffStatusAktif,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := h.DB.Create(&staff).Error; err != nil {
		// Rollback user creation if staff creation fails
		h.DB.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat profile staff",
		})
		return
	}

	// Send notification if service is available
	notificationService := services.NewNotificationService(h.DB)
	if err := notificationService.SendStaffCreatedNotification(userID, input.Nama, input.Jabatan); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Gagal mengirim notifikasi pembuatan staff: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Staff berhasil dibuat",
		"data": gin.H{
			"user": gin.H{
				"user_id":  user.User_id,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
				"nama":     user.Nama,
			},
			"staff": gin.H{
				"id":           staff.ID,
				"nip":          staff.NIP,
				"nama_lengkap": staff.Nama_lengkap,
				"jabatan":      staff.Jabatan,
				"bagian":       staff.Bagian,
				"no_hp":        staff.No_hp,
				"email":        staff.Email,
				"alamat":       staff.Alamat,
				"status":       staff.Status,
			},
		},
	})
}

// ==================== MARRIAGE OFFICER (PENGHULU) MANAGEMENT ====================

// CreateMarriageOfficer creates a new marriage officer (penghulu) with user account
func (h *InDB) CreateMarriageOfficer(c *gin.Context) {
	var input struct {
		Username      string `json:"username" binding:"required"`
		Email         string `json:"email" binding:"required,email"`
		Password      string `json:"password" binding:"required,min=6"`
		Nama          string `json:"nama" binding:"required"`
		NIP           string `json:"nip" binding:"required"`
		No_hp         string `json:"no_hp"`
		EmailPenghulu string `json:"email_penghulu"`
		Alamat        string `json:"alamat"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Check if username/email already exists
	var existingUser structs.Users
	if err := h.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username atau email sudah digunakan",
			"error":   "Username atau email sudah terdaftar di sistem",
		})
		return
	}

	// Check if NIP already exists
	var existingPenghulu structs.Penghulu
	if err := h.DB.Where("nip = ?", input.NIP).First(&existingPenghulu).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "NIP sudah terdaftar",
			"error":   "NIP sudah terdaftar untuk penghulu lain",
		})
		return
	}

	// Generate user_id
	userID := "PNG" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password
	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengenkripsi password",
		})
		return
	}

	// Use email_penghulu if provided, otherwise use email
	penghuluEmail := input.Email
	if input.EmailPenghulu != "" {
		penghuluEmail = input.EmailPenghulu
	}

	// Create user account
	user := structs.Users{
		User_id:    userID,
		Username:   input.Username,
		Email:      input.Email,
		Password:   hashedPassword,
		Role:       structs.UserRolePenghulu,
		Status:     structs.UserStatusAktif,
		Nama:       input.Nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat user account",
		})
		return
	}

	// Create penghulu profile
	penghulu := structs.Penghulu{
		User_id:      userID,
		NIP:          input.NIP,
		Nama_lengkap: input.Nama,
		No_hp:        input.No_hp,
		Email:        penghuluEmail,
		Alamat:       input.Alamat,
		Status:       structs.PenghuluStatusAktif,
		Jumlah_nikah: 0,
		Rating:       0.0,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := h.DB.Create(&penghulu).Error; err != nil {
		// Rollback user creation if penghulu creation fails
		h.DB.Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat profile penghulu",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Penghulu berhasil dibuat",
		"data": gin.H{
			"user": gin.H{
				"user_id":  user.User_id,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
				"nama":     user.Nama,
			},
			"penghulu": gin.H{
				"id":           penghulu.ID,
				"nip":          penghulu.NIP,
				"nama_lengkap": penghulu.Nama_lengkap,
				"no_hp":        penghulu.No_hp,
				"email":        penghulu.Email,
				"alamat":       penghulu.Alamat,
				"status":       penghulu.Status,
			},
		},
	})
}

// ==================== PENGHULU STATISTICS ====================

// GetPenghuluStatistics returns statistics for all penghulu or a specific penghulu
// GET /simnikah/kepala-kua/statistik-penghulu?penghulu_id=1&bulan=01&tahun=2024
func (h *InDB) GetPenghuluStatistics(c *gin.Context) {
	// Get query parameters
	penghuluIDStr := c.Query("penghulu_id")
	bulanStr := c.Query("bulan")
	tahunStr := c.Query("tahun")

	// Parse tahun and bulan (optional filters)
	now := time.Now()
	tahun := now.Year()
	bulan := int(now.Month())

	if tahunStr != "" {
		if tahunInt, err := strconv.Atoi(tahunStr); err == nil && tahunInt > 0 {
			tahun = tahunInt
		}
	}

	if bulanStr != "" {
		if bulanInt, err := strconv.Atoi(bulanStr); err == nil && bulanInt >= 1 && bulanInt <= 12 {
			bulan = bulanInt
		}
	}

	// Calculate start and end of month
	startOfMonth := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// If penghulu_id is provided, get statistics for specific penghulu
	if penghuluIDStr != "" {
		penghuluID, err := strconv.ParseUint(penghuluIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID penghulu tidak valid",
				"error":   "Format ID harus berupa angka",
			})
			return
		}

		h.getSinglePenghuluStatistics(c, uint(penghuluID), startOfMonth, endOfMonth, bulan, tahun)
		return
	}

	// Get statistics for all penghulu
	h.getAllPenghuluStatistics(c, startOfMonth, endOfMonth, bulan, tahun)
}

// getSinglePenghuluStatistics gets statistics for a specific penghulu
func (h *InDB) getSinglePenghuluStatistics(c *gin.Context, penghuluID uint, startOfMonth, endOfMonth time.Time, bulan, tahun int) {
	// Check if penghulu exists
	var penghulu structs.Penghulu
	if err := h.DB.Where("id = ?", penghuluID).First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Penghulu tidak ditemukan",
			"error":   "Penghulu dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Get total registrations (all time)
	var totalAllTime int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id = ? AND status_pendaftaran NOT IN ?", penghuluID, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Count(&totalAllTime)

	// Get completed registrations (all time)
	var totalCompleted int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id = ? AND status_pendaftaran = ?", penghuluID, structs.StatusPendaftaranSelesai).
		Count(&totalCompleted)

	// Get in progress registrations
	var totalInProgress int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id = ? AND status_pendaftaran IN ?", penghuluID, []string{structs.StatusPendaftaranPenghuluDitugaskan, structs.StatusPendaftaranMenungguPenugasan}).
		Count(&totalInProgress)

	// Get statistics for this month
	var totalThisMonth int64
	var completedThisMonth int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?",
			penghuluID, startOfMonth, endOfMonth, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Count(&totalThisMonth)

	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
			penghuluID, startOfMonth, endOfMonth, structs.StatusPendaftaranSelesai).
		Count(&completedThisMonth)

	// Get statistics per month for last 6 months
	var monthlyStats []gin.H
	for i := 5; i >= 0; i-- {
		monthStart := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0)

		var monthTotal int64
		var monthCompleted int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?",
				penghuluID, monthStart, monthEnd, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
			Count(&monthTotal)

		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
				penghuluID, monthStart, monthEnd, structs.StatusPendaftaranSelesai).
			Count(&monthCompleted)

		monthlyStats = append(monthlyStats, gin.H{
			"bulan":      int(monthStart.Month()),
			"tahun":      monthStart.Year(),
			"nama_bulan": monthStart.Month().String(),
			"total":      monthTotal,
			"selesai":    monthCompleted,
			"progress":   monthTotal - monthCompleted,
		})
	}

	// Get recent registrations (last 10)
	var recentRegistrations []structs.PendaftaranNikah
	h.DB.Where("penghulu_id = ? AND status_pendaftaran NOT IN ?", penghuluID, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Order("tanggal_nikah DESC").
		Limit(10).
		Find(&recentRegistrations)

	var recentData []gin.H
	for _, reg := range recentRegistrations {
		recentData = append(recentData, gin.H{
			"id":                reg.ID,
			"nomor_pendaftaran": reg.Nomor_pendaftaran,
			"tanggal_nikah":     reg.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":       reg.Waktu_nikah,
			"tempat_nikah":      reg.Tempat_nikah,
			"status":            reg.Status_pendaftaran,
		})
	}

	// Calculate average per month (last 6 months)
	var avgPerMonth float64
	if len(monthlyStats) > 0 {
		total := int64(0)
		for _, stat := range monthlyStats {
			if totalVal, ok := stat["total"].(int64); ok {
				total += totalVal
			}
		}
		avgPerMonth = float64(total) / float64(len(monthlyStats))
	}

	// Calculate completion percentage
	var completionPctAllTime float64
	var completionPctThisMonth float64
	if totalAllTime > 0 {
		completionPctAllTime = float64(totalCompleted) / float64(totalAllTime) * 100
	}
	if totalThisMonth > 0 {
		completionPctThisMonth = float64(completedThisMonth) / float64(totalThisMonth) * 100
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik penghulu berhasil diambil",
		"data": gin.H{
			"penghulu": gin.H{
				"id":           penghulu.ID,
				"nip":          penghulu.NIP,
				"nama_lengkap": penghulu.Nama_lengkap,
				"no_hp":        penghulu.No_hp,
				"email":        penghulu.Email,
				"status":       penghulu.Status,
			},
			"periode": gin.H{
				"bulan":      bulan,
				"tahun":      tahun,
				"nama_bulan": startOfMonth.Month().String(),
			},
			"statistik_keseluruhan": gin.H{
				"total_semua":        totalAllTime,
				"selesai_semua":      totalCompleted,
				"progress_semua":     totalInProgress,
				"persentase_selesai": completionPctAllTime,
			},
			"statistik_bulan_ini": gin.H{
				"total":              totalThisMonth,
				"selesai":            completedThisMonth,
				"progress":           totalThisMonth - completedThisMonth,
				"persentase_selesai": completionPctThisMonth,
			},
			"statistik_per_bulan": monthlyStats,
			"rata_rata_per_bulan": avgPerMonth,
			"registrasi_terbaru":  recentData,
		},
	})
}

// getAllPenghuluStatistics gets statistics for all penghulu
func (h *InDB) getAllPenghuluStatistics(c *gin.Context, startOfMonth, endOfMonth time.Time, bulan, tahun int) {
	// Get all active penghulu
	var penghulus []structs.Penghulu
	if err := h.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data penghulu",
			"error":   err.Error(),
		})
		return
	}

	var penghuluStats []gin.H

	for _, penghulu := range penghulus {
		// Get total registrations (all time)
		var totalAllTime int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND status_pendaftaran NOT IN ?", penghulu.ID, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
			Count(&totalAllTime)

		// Get completed registrations (all time)
		var totalCompleted int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND status_pendaftaran = ?", penghulu.ID, structs.StatusPendaftaranSelesai).
			Count(&totalCompleted)

		// Get statistics for this month
		var totalThisMonth int64
		var completedThisMonth int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?",
				penghulu.ID, startOfMonth, endOfMonth, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
			Count(&totalThisMonth)

		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
				penghulu.ID, startOfMonth, endOfMonth, structs.StatusPendaftaranSelesai).
			Count(&completedThisMonth)

		// Calculate percentages
		var completionPctAllTime float64
		var completionPctThisMonth float64
		if totalAllTime > 0 {
			completionPctAllTime = float64(totalCompleted) / float64(totalAllTime) * 100
		}
		if totalThisMonth > 0 {
			completionPctThisMonth = float64(completedThisMonth) / float64(totalThisMonth) * 100
		}

		penghuluStats = append(penghuluStats, gin.H{
			"penghulu": gin.H{
				"id":           penghulu.ID,
				"nip":          penghulu.NIP,
				"nama_lengkap": penghulu.Nama_lengkap,
				"status":       penghulu.Status,
			},
			"statistik_semua": gin.H{
				"total":              totalAllTime,
				"selesai":            totalCompleted,
				"progress":           totalAllTime - totalCompleted,
				"persentase_selesai": completionPctAllTime,
			},
			"statistik_bulan_ini": gin.H{
				"total":              totalThisMonth,
				"selesai":            completedThisMonth,
				"progress":           totalThisMonth - completedThisMonth,
				"persentase_selesai": completionPctThisMonth,
			},
		})
	}

	// Get overall statistics
	var totalAllPenghulu int64
	var totalCompletedAllPenghulu int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id IS NOT NULL AND status_pendaftaran NOT IN ?", []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Count(&totalAllPenghulu)

	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id IS NOT NULL AND status_pendaftaran = ?", structs.StatusPendaftaranSelesai).
		Count(&totalCompletedAllPenghulu)

	var totalThisMonthAllPenghulu int64
	var completedThisMonthAllPenghulu int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id IS NOT NULL AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?",
			startOfMonth, endOfMonth, []string{structs.StatusPendaftaranDraft, structs.StatusPendaftaranDitolak}).
		Count(&totalThisMonthAllPenghulu)

	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id IS NOT NULL AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
			startOfMonth, endOfMonth, structs.StatusPendaftaranSelesai).
		Count(&completedThisMonthAllPenghulu)

	// Calculate overall completion percentage
	var completionPctAllTime float64
	var completionPctThisMonth float64
	if totalAllPenghulu > 0 {
		completionPctAllTime = float64(totalCompletedAllPenghulu) / float64(totalAllPenghulu) * 100
	}
	if totalThisMonthAllPenghulu > 0 {
		completionPctThisMonth = float64(completedThisMonthAllPenghulu) / float64(totalThisMonthAllPenghulu) * 100
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik semua penghulu berhasil diambil",
		"data": gin.H{
			"periode": gin.H{
				"bulan":      bulan,
				"tahun":      tahun,
				"nama_bulan": startOfMonth.Month().String(),
			},
			"statistik_keseluruhan": gin.H{
				"total_penghulu":           len(penghulus),
				"total_semua":              totalAllPenghulu,
				"selesai_semua":            totalCompletedAllPenghulu,
				"progress_semua":           totalAllPenghulu - totalCompletedAllPenghulu,
				"persentase_selesai":       completionPctAllTime,
				"total_bulan_ini":          totalThisMonthAllPenghulu,
				"selesai_bulan_ini":        completedThisMonthAllPenghulu,
				"progress_bulan_ini":       totalThisMonthAllPenghulu - completedThisMonthAllPenghulu,
				"persentase_selesai_bulan": completionPctThisMonth,
			},
			"statistik_per_penghulu": penghuluStats,
		},
	})
}

// ==================== FEEDBACK MANAGEMENT ====================

// ListFeedbackPernikahan gets all feedback from catin (for kepala KUA)
// GET /simnikah/kepala-kua/feedback?jenis=Rating&status=Belum Dibaca&bulan=01&tahun=2024
func (h *InDB) ListFeedbackPernikahan(c *gin.Context) {
	// Get query parameters
	jenisFeedback := c.Query("jenis") // "Rating", "Saran", "Kritik", "Laporan"
	statusBaca := c.Query("status")   // "Belum Dibaca", "Sudah Dibaca"
	bulanStr := c.Query("bulan")
	tahunStr := c.Query("tahun")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse pagination
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Build query
	query := h.DB.Model(&structs.FeedbackPernikahan{})

	// Apply filters
	if jenisFeedback != "" {
		query = query.Where("jenis_feedback = ?", jenisFeedback)
	}

	if statusBaca != "" {
		query = query.Where("status_baca = ?", statusBaca)
	}

	// Filter by date (bulan and tahun)
	if bulanStr != "" && tahunStr != "" {
		if bulanInt, err := strconv.Atoi(bulanStr); err == nil && bulanInt >= 1 && bulanInt <= 12 {
			if tahunInt, err := strconv.Atoi(tahunStr); err == nil && tahunInt > 0 {
				startOfMonth := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.UTC)
				endOfMonth := startOfMonth.AddDate(0, 1, 0)
				query = query.Where("created_at >= ? AND created_at < ?", startOfMonth, endOfMonth)
			}
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results with join to get registration and user info
	type FeedbackWithDetails struct {
		structs.FeedbackPernikahan
		NomorPendaftaran string    `gorm:"column:nomor_pendaftaran"`
		NamaUser         string    `gorm:"column:nama_user"`
		EmailUser        string    `gorm:"column:email_user"`
		TanggalNikah     time.Time `gorm:"column:tanggal_nikah"`
	}

	var results []FeedbackWithDetails
	err = query.
		Select("feedback_pernikahans.*, p.nomor_pendaftaran, u.nama as nama_user, u.email as email_user, p.tanggal_nikah").
		Joins("LEFT JOIN pendaftaran_nikahs p ON feedback_pernikahans.pendaftaran_id = p.id").
		Joins("LEFT JOIN users u ON feedback_pernikahans.user_id = u.user_id").
		Order("feedback_pernikahans.created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data feedback",
			"error":   err.Error(),
		})
		return
	}

	// Prepare response data
	var feedbackList []gin.H
	for _, r := range results {
		feedbackList = append(feedbackList, gin.H{
			"id":                r.ID,
			"pendaftaran_id":    r.Pendaftaran_id,
			"nomor_pendaftaran": r.NomorPendaftaran,
			"user": gin.H{
				"user_id": r.User_id,
				"nama":    r.NamaUser,
				"email":   r.EmailUser,
			},
			"jenis_feedback": r.Jenis_feedback,
			"rating":         r.Rating,
			"judul":          r.Judul,
			"pesan":          r.Pesan,
			"status_baca":    r.Status_baca,
			"dibaca_oleh":    r.Dibaca_oleh,
			"dibaca_pada":    r.Dibaca_pada,
			"tanggal_nikah":  r.TanggalNikah.Format("2006-01-02"),
			"created_at":     r.Created_at,
		})
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data feedback berhasil diambil",
		"data": gin.H{
			"feedback": feedbackList,
			"pagination": gin.H{
				"current_page":  page,
				"total_pages":   totalPages,
				"total_records": total,
				"per_page":      limit,
				"has_next":      hasNext,
				"has_previous":  hasPrev,
			},
			"filters": gin.H{
				"jenis":  jenisFeedback,
				"status": statusBaca,
				"bulan":  bulanStr,
				"tahun":  tahunStr,
			},
		},
	})
}

// MarkFeedbackAsRead marks feedback as read by kepala KUA
// PUT /simnikah/kepala-kua/feedback/:id/mark-read
func (h *InDB) MarkFeedbackAsRead(c *gin.Context) {
	feedbackID := c.Param("id")

	// Get user_id from context (kepala KUA)
	kepalaKuaID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	// Check if feedback exists
	var feedback structs.FeedbackPernikahan
	if err := h.DB.Where("id = ?", feedbackID).First(&feedback).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Feedback tidak ditemukan",
			"error":   "Feedback dengan ID tersebut tidak ditemukan",
		})
		return
	}

	// Check if already read
	if feedback.Status_baca == structs.FeedbackStatusSudahDibaca {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Feedback sudah dibaca",
			"error":   "Feedback ini sudah ditandai sebagai sudah dibaca",
		})
		return
	}

	// Update feedback
	now := time.Now()
	feedback.Status_baca = structs.FeedbackStatusSudahDibaca
	feedback.Dibaca_oleh = kepalaKuaID.(string)
	feedback.Dibaca_pada = &now
	feedback.Updated_at = now

	if err := h.DB.Save(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate feedback",
			"error":   "Database error",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Feedback berhasil ditandai sebagai sudah dibaca",
		"data": gin.H{
			"id":          feedback.ID,
			"status_baca": feedback.Status_baca,
			"dibaca_oleh": feedback.Dibaca_oleh,
			"dibaca_pada": feedback.Dibaca_pada,
			"updated_at":  feedback.Updated_at,
		},
	})
}

// GetFeedbackStats gets statistics for feedback (for kepala KUA dashboard)
// GET /simnikah/kepala-kua/feedback/stats
func (h *InDB) GetFeedbackStats(c *gin.Context) {
	// Get total feedback
	var totalFeedback int64
	h.DB.Model(&structs.FeedbackPernikahan{}).Count(&totalFeedback)

	// Get unread feedback
	var unreadFeedback int64
	h.DB.Model(&structs.FeedbackPernikahan{}).
		Where("status_baca = ?", structs.FeedbackStatusBelumDibaca).
		Count(&unreadFeedback)

	// Get feedback by jenis
	var statsByJenis []gin.H
	jenisList := []string{structs.FeedbackJenisRating, structs.FeedbackJenisSaran, structs.FeedbackJenisKritik, structs.FeedbackJenisLaporan}

	for _, jenis := range jenisList {
		var count int64
		h.DB.Model(&structs.FeedbackPernikahan{}).
			Where("jenis_feedback = ?", jenis).
			Count(&count)
		statsByJenis = append(statsByJenis, gin.H{
			"jenis": jenis,
			"total": count,
		})
	}

	// Get average rating (only for Rating jenis)
	var avgRating float64
	var ratingCount int64
	var sumRating int64
	h.DB.Model(&structs.FeedbackPernikahan{}).
		Where("jenis_feedback = ? AND rating IS NOT NULL", structs.FeedbackJenisRating).
		Count(&ratingCount)

	if ratingCount > 0 {
		// Get all ratings and calculate sum
		var feedbacks []structs.FeedbackPernikahan
		h.DB.Where("jenis_feedback = ? AND rating IS NOT NULL", structs.FeedbackJenisRating).
			Find(&feedbacks)

		sumRating = 0
		for _, fb := range feedbacks {
			if fb.Rating != nil {
				sumRating += int64(*fb.Rating)
			}
		}

		avgRating = float64(sumRating) / float64(ratingCount)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik feedback berhasil diambil",
		"data": gin.H{
			"total_feedback":  totalFeedback,
			"unread_feedback": unreadFeedback,
			"read_feedback":   totalFeedback - unreadFeedback,
			"stats_by_jenis":  statsByJenis,
			"rating_stats": gin.H{
				"average_rating": avgRating,
				"total_ratings":  ratingCount,
			},
		},
	})
}

// ==================== SURAT PENGUMUMAN NIKAH ====================

// GetApprovedRegistrationsPerWeek gets approved registrations for a specific week
// Used for generating pengumuman nikah (marriage announcement)
// Supports custom kop surat via request body (optional)
func (h *InDB) GetApprovedRegistrationsPerWeek(c *gin.Context) {
	// Get query parameters (prioritas pertama)
	tanggalAwal := c.Query("tanggal_awal")   // Format: YYYY-MM-DD (start of week)
	tanggalAkhir := c.Query("tanggal_akhir") // Format: YYYY-MM-DD (end of week)

	// Get from request body (fallback jika query params tidak ada)
	var requestBody struct {
		TanggalAwal  string `json:"tanggal_awal"`  // Format: YYYY-MM-DD
		TanggalAkhir string `json:"tanggal_akhir"` // Format: YYYY-MM-DD
		NamaKUA      string `json:"nama_kua"`      // Nama KUA
		AlamatKUA    string `json:"alamat_kua"`    // Alamat lengkap KUA
		Kota         string `json:"kota"`          // Kota
		Provinsi     string `json:"provinsi"`      // Provinsi
		KodePos      string `json:"kode_pos"`      // Kode pos
		Telepon      string `json:"telepon"`        // Nomor telepon
		Email        string `json:"email"`         // Email
		Website      string `json:"website"`       // Website (optional)
		LogoURL      string `json:"logo_url"`     // URL logo KUA (optional)
	}

	// Try to bind JSON body, if not provided use default
	if c.Request.ContentLength > 0 {
		c.ShouldBindJSON(&requestBody)
	}

	// Gunakan request body jika query params tidak ada
	if tanggalAwal == "" && requestBody.TanggalAwal != "" {
		tanggalAwal = requestBody.TanggalAwal
	}
	if tanggalAkhir == "" && requestBody.TanggalAkhir != "" {
		tanggalAkhir = requestBody.TanggalAkhir
	}

	// Set default kop surat values
	namaKUA := requestBody.NamaKUA
	alamatKUA := requestBody.AlamatKUA
	kota := requestBody.Kota
	provinsi := requestBody.Provinsi
	kodePos := requestBody.KodePos
	telepon := requestBody.Telepon
	email := requestBody.Email
	website := requestBody.Website
	logoURL := requestBody.LogoURL

	// Set default values if not provided
	if namaKUA == "" {
		namaKUA = "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA"
	}
	if alamatKUA == "" {
		alamatKUA = "PH5Q+F8C, Jl. Wira Karya, Pangeran"
	}
	if kota == "" {
		kota = "Kota Banjarmasin"
	}
	if provinsi == "" {
		provinsi = "Kalimantan Selatan"
	}
	if kodePos == "" {
		kodePos = "70123"
	}
	if telepon == "" {
		telepon = "-"
	}
	if email == "" {
		email = "-"
	}

	// Prepare kop surat for response
	kopSurat := gin.H{
		"nama_kua":   namaKUA,
		"alamat_kua": alamatKUA,
		"kota":       kota,
		"provinsi":   provinsi,
		"kode_pos":   kodePos,
		"telepon":    telepon,
		"email":      email,
		"website":    website,
		"logo_url":   logoURL,
	}

	// If not provided, use current week
	now := time.Now()
	var startOfWeek, endOfWeek time.Time

	if tanggalAwal != "" && tanggalAkhir != "" {
		// Parse provided dates
		start, err := time.Parse("2006-01-02", tanggalAwal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal tidak valid",
				"error":   "Format tanggal_awal harus YYYY-MM-DD",
			})
			return
		}
		end, err := time.Parse("2006-01-02", tanggalAkhir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal tidak valid",
				"error":   "Format tanggal_akhir harus YYYY-MM-DD",
			})
			return
		}
		startOfWeek = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		endOfWeek = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, time.UTC)
	} else {
		// Default: current week (Monday to Sunday)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		daysFromMonday := weekday - 1
		startOfWeek = time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday, 0, 0, 0, 0, time.UTC)
		endOfWeek = startOfWeek.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Get registrations within the week (all status except "Ditolak")
	var pendaftaran []structs.PendaftaranNikah
	err := h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah <= ? AND status_pendaftaran != ?",
		startOfWeek, endOfWeek, structs.StatusPendaftaranDitolak).
		Order("tanggal_nikah ASC, waktu_nikah ASC").
		Find(&pendaftaran).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data pendaftaran",
		})
		return
	}

	// Get calon pasangan and wali nikah data
	var registrations []gin.H
	for _, p := range pendaftaran {
		// Get calon suami
		var calonSuami structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)

		// Get calon istri
		var calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		// Get wali nikah
		var waliNikah structs.WaliNikah
		if p.Wali_nikah_id != nil {
			h.DB.Where("id = ?", *p.Wali_nikah_id).First(&waliNikah)
		}

		regData := gin.H{
			"id":                p.ID,
			"nomor_pendaftaran": p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":     p.Tanggal_nikah,
			"waktu_nikah":       p.Waktu_nikah,
			"tempat_nikah":      p.Tempat_nikah,
			"alamat_akad":       p.Alamat_akad,
			"calon_suami": gin.H{
				"nama_lengkap": calonSuami.Nama_lengkap,
			},
			"calon_istri": gin.H{
				"nama_lengkap": calonIstri.Nama_lengkap,
			},
		}

		// Add wali nikah if exists
		if waliNikah.ID != 0 {
			regData["wali_nikah"] = gin.H{
				"nama_dan_bin":  waliNikah.Nama_dan_bin,
				"hubungan_wali": waliNikah.Hubungan_wali,
			}
		}

		registrations = append(registrations, regData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pendaftaran berhasil diambil",
		"data": gin.H{
			"tanggal_awal":  startOfWeek.Format("2006-01-02"),
			"tanggal_akhir": endOfWeek.Format("2006-01-02"),
			"periode":       fmt.Sprintf("%s s/d %s", startOfWeek.Format("02 Januari 2006"), endOfWeek.Format("02 Januari 2006")),
			"total":         len(registrations),
			"kop_surat":     kopSurat,
			"registrations": registrations,
		},
	})
}

// GeneratePengumumanNikahHTML generates HTML document for pengumuman nikah
// GET /simnikah/kepala-kua/pengumuman-nikah/generate
func (h *InDB) GeneratePengumumanNikahHTML(c *gin.Context) {
	// Get query parameters (prioritas pertama)
	tanggalAwal := c.Query("tanggal_awal")
	tanggalAkhir := c.Query("tanggal_akhir")

	// Get from request body (fallback jika query params tidak ada)
	var requestBody struct {
		TanggalAwal  string `json:"tanggal_awal"`
		TanggalAkhir string `json:"tanggal_akhir"`
		NamaKUA      string `json:"nama_kua"`
		AlamatKUA    string `json:"alamat_kua"`
		Kota         string `json:"kota"`
		Provinsi     string `json:"provinsi"`
		KodePos      string `json:"kode_pos"`
		Telepon      string `json:"telepon"`
		Email        string `json:"email"`
		Website      string `json:"website"`
		LogoURL      string `json:"logo_url"`
	}

	// Try to bind JSON body
	if c.Request.ContentLength > 0 {
		c.ShouldBindJSON(&requestBody)
	}

	// Gunakan request body jika query params tidak ada
	if tanggalAwal == "" && requestBody.TanggalAwal != "" {
		tanggalAwal = requestBody.TanggalAwal
	}
	if tanggalAkhir == "" && requestBody.TanggalAkhir != "" {
		tanggalAkhir = requestBody.TanggalAkhir
	}

	// Set kop surat values
	namaKUA := requestBody.NamaKUA
	alamatKUA := requestBody.AlamatKUA
	kota := requestBody.Kota
	provinsi := requestBody.Provinsi
	kodePos := requestBody.KodePos
	telepon := requestBody.Telepon
	email := requestBody.Email
	website := requestBody.Website
	logoURL := requestBody.LogoURL

	// Set default values if not provided
	if namaKUA == "" {
		namaKUA = "KANTOR URUSAN AGAMA KECAMATAN BANJARMASIN UTARA"
	}
	if alamatKUA == "" {
		alamatKUA = "PH5Q+F8C, Jl. Wira Karya, Pangeran"
	}
	if kota == "" {
		kota = "Kota Banjarmasin"
	}
	if provinsi == "" {
		provinsi = "Kalimantan Selatan"
	}
	if kodePos == "" {
		kodePos = "70123"
	}
	if telepon == "" {
		telepon = "-"
	}
	if email == "" {
		email = "-"
	}

	// Parse dates
	now := time.Now()
	var startOfWeek, endOfWeek time.Time

	if tanggalAwal != "" && tanggalAkhir != "" {
		start, err := time.Parse("2006-01-02", tanggalAwal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal tidak valid",
				"error":   "Format tanggal_awal harus YYYY-MM-DD",
			})
			return
		}
		end, err := time.Parse("2006-01-02", tanggalAkhir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal tidak valid",
				"error":   "Format tanggal_akhir harus YYYY-MM-DD",
			})
			return
		}
		startOfWeek = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		endOfWeek = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, time.UTC)
	} else {
		// Default: current week
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		daysFromMonday := weekday - 1
		startOfWeek = time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday, 0, 0, 0, 0, time.UTC)
		endOfWeek = startOfWeek.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Get registrations
	var pendaftaran []structs.PendaftaranNikah
	err := h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah <= ? AND status_pendaftaran != ?",
		startOfWeek, endOfWeek, structs.StatusPendaftaranDitolak).
		Order("tanggal_nikah ASC, waktu_nikah ASC").
		Find(&pendaftaran).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data pendaftaran",
		})
		return
	}

	// Build registrations data
	type RegData struct {
		No              int
		NomorPendaftaran string
		TanggalNikah     string
		WaktuNikah       string
		TempatNikah      string
		AlamatAkad       string
		CalonSuami       string
		CalonIstri       string
		WaliNikah        string
		HubunganWali     string
	}

	var regDataList []RegData
	for i, p := range pendaftaran {
		// Get calon suami
		var calonSuami structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)

		// Get calon istri
		var calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		// Get wali nikah
		var waliNikah structs.WaliNikah
		waliNikahStr := "-"
		hubunganWaliStr := "-"
		if p.Wali_nikah_id != nil {
			h.DB.Where("id = ?", *p.Wali_nikah_id).First(&waliNikah)
			if waliNikah.ID != 0 {
				waliNikahStr = waliNikah.Nama_dan_bin
				hubunganWaliStr = waliNikah.Hubungan_wali
			}
		}

		regDataList = append(regDataList, RegData{
			No:               i + 1,
			NomorPendaftaran: p.Nomor_pendaftaran,
			TanggalNikah:     p.Tanggal_nikah.Format("02 Januari 2006"),
			WaktuNikah:       p.Waktu_nikah,
			TempatNikah:      p.Tempat_nikah,
			AlamatAkad:       p.Alamat_akad,
			CalonSuami:       calonSuami.Nama_lengkap,
			CalonIstri:       calonIstri.Nama_lengkap,
			WaliNikah:        waliNikahStr,
			HubunganWali:     hubunganWaliStr,
		})
	}

	// Generate HTML
	periodeStr := fmt.Sprintf("%s s/d %s", startOfWeek.Format("02 Januari 2006"), endOfWeek.Format("02 Januari 2006"))
	tanggalSurat := now.Format("02 Januari 2006")

	html := h.generatePengumumanHTML(namaKUA, alamatKUA, kota, provinsi, kodePos, telepon, email, website, logoURL, periodeStr, tanggalSurat, regDataList)

	// Set response headers
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// generatePengumumanHTML generates HTML string for pengumuman nikah
func (h *InDB) generatePengumumanHTML(namaKUA, alamatKUA, kota, provinsi, kodePos, telepon, email, website, logoURL, periode, tanggalSurat string, registrations []RegData) string {
	logoHTML := ""
	if logoURL != "" {
		logoHTML = fmt.Sprintf(`<img src="%s" alt="Logo KUA" style="max-width: 100px; max-height: 100px; margin-bottom: 10px;">`, logoURL)
	}

	websiteHTML := ""
	if website != "" {
		websiteHTML = fmt.Sprintf(`<p style="margin: 2px 0;">Website: %s</p>`, website)
	}

	tableRows := ""
	for _, reg := range registrations {
		tableRows += fmt.Sprintf(`
			<tr>
				<td style="border: 1px solid #000; padding: 8px; text-align: center;">%d</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
				<td style="border: 1px solid #000; padding: 8px;">%s</td>
			</tr>`,
			reg.No,
			reg.NomorPendaftaran,
			reg.TanggalNikah,
			reg.WaktuNikah,
			reg.TempatNikah,
			reg.AlamatAkad,
			reg.CalonSuami,
			reg.CalonIstri,
			reg.WaliNikah,
		)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Pengumuman Pernikahan - %s</title>
	<style>
		@media print {
			@page {
				size: A4;
				margin: 2cm;
			}
			.no-print {
				display: none;
			}
		}
		body {
			font-family: 'Times New Roman', Times, serif;
			font-size: 12pt;
			line-height: 1.6;
			margin: 0;
			padding: 20px;
			color: #000;
		}
		.header {
			text-align: center;
			margin-bottom: 30px;
		}
		.logo {
			margin-bottom: 10px;
		}
		.kop-surat {
			border-bottom: 3px solid #000;
			padding-bottom: 10px;
			margin-bottom: 20px;
		}
		.kop-surat h1 {
			margin: 5px 0;
			font-size: 14pt;
			font-weight: bold;
			text-transform: uppercase;
		}
		.kop-surat p {
			margin: 2px 0;
			font-size: 11pt;
		}
		.content {
			margin: 20px 0;
		}
		.content p {
			text-align: justify;
			margin: 10px 0;
			text-indent: 50px;
		}
		table {
			width: 100%%;
			border-collapse: collapse;
			margin: 20px 0;
			font-size: 10pt;
		}
		table th {
			background-color: #f0f0f0;
			border: 1px solid #000;
			padding: 10px;
			text-align: center;
			font-weight: bold;
		}
		table td {
			border: 1px solid #000;
			padding: 8px;
		}
		.signature {
			margin-top: 50px;
			text-align: right;
		}
		.signature p {
			margin: 5px 0;
		}
		.footer {
			margin-top: 30px;
			text-align: center;
			font-size: 10pt;
		}
	</style>
</head>
<body>
	<div class="kop-surat">
		<div class="header">
			<div class="logo">%s</div>
			<h1>%s</h1>
			<p>%s</p>
			<p>%s, %s %s</p>
			<p style="margin: 2px 0;">Telp: %s | Email: %s</p>
			%s
		</div>
	</div>

	<div class="content">
		<p style="text-align: center; font-weight: bold; font-size: 14pt; margin: 20px 0;">PENGUMUMAN PERNIKAHAN</p>
		
		<p>Berdasarkan data pendaftaran pernikahan yang telah diterima oleh Kantor Urusan Agama, dengan ini diumumkan kepada masyarakat bahwa akan dilaksanakan pernikahan pada periode <strong>%s</strong> dengan rincian sebagai berikut:</p>

		<table>
			<thead>
				<tr>
					<th style="width: 3%%;">No</th>
					<th style="width: 10%%;">Nomor Pendaftaran</th>
					<th style="width: 10%%;">Tanggal Nikah</th>
					<th style="width: 8%%;">Waktu</th>
					<th style="width: 10%%;">Tempat</th>
					<th style="width: 15%%;">Alamat Akad</th>
					<th style="width: 12%%;">Calon Suami</th>
					<th style="width: 12%%;">Calon Istri</th>
					<th style="width: 10%%;">Wali Nikah</th>
				</tr>
			</thead>
			<tbody>
				%s
			</tbody>
		</table>

		<p>Pengumuman ini dibuat untuk memberikan kesempatan kepada masyarakat yang berkeberatan terhadap pernikahan tersebut untuk menyampaikan keberatannya sesuai dengan ketentuan yang berlaku.</p>

		<p>Demikian pengumuman ini dibuat dengan sebenarnya untuk dapat dipergunakan sebagaimana mestinya.</p>
	</div>

	<div class="signature">
		<p>%s, %s</p>
		<br><br><br>
		<p><strong>Kepala KUA</strong></p>
		<p style="margin-top: 50px;"><strong>(_______________________)</strong></p>
	</div>

	<div class="footer">
		<p>Surat ini dicetak pada: %s</p>
	</div>
</body>
</html>`,
		namaKUA,
		logoHTML,
		namaKUA,
		alamatKUA,
		kota,
		provinsi,
		kodePos,
		telepon,
		email,
		websiteHTML,
		periode,
		tableRows,
		kota,
		tanggalSurat,
		tanggalSurat,
	)

	return html
}

// ==================== PENGHULU SCHEDULE AVAILABILITY ====================

// TimeSlots untuk jadwal pernikahan
var TimeSlots = []string{
	"08:00", "09:00", "10:00", "11:00", "12:00",
	"13:00", "14:00", "15:00", "16:00",
}

// GetPenghuluScheduleAvailability returns schedule availability for all penghulus on a specific date
// GET /simnikah/kepala-kua/penghulu-schedule?tanggal=2024-12-25
func (h *InDB) GetPenghuluScheduleAvailability(c *gin.Context) {
	tanggalStr := c.Query("tanggal")
	if tanggalStr == "" {
		// Default ke hari ini
		tanggalStr = time.Now().Format("2006-01-02")
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

	// Get all active penghulus
	var penghulus []structs.Penghulu
	if err := h.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data penghulu",
		})
		return
	}

	// Calculate start and end of day
	startOfDay := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get all registrations on this date that have penghulu assigned
	var pendaftarans []structs.PendaftaranNikah
	err = h.DB.Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND penghulu_id IS NOT NULL AND status_pendaftaran NOT IN ?",
		startOfDay, endOfDay, []string{structs.StatusPendaftaranDitolak}).
		Find(&pendaftarans).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data jadwal",
		})
		return
	}

	// Build schedule map: penghulu_id -> waktu -> pendaftaran
	type ScheduleItem struct {
		PendaftaranID    uint   `json:"pendaftaran_id"`
		NomorPendaftaran string `json:"nomor_pendaftaran"`
		WaktuNikah       string `json:"waktu_nikah"`
		TempatNikah      string `json:"tempat_nikah"`
		AlamatAkad       string `json:"alamat_akad"`
		CalonSuami       string `json:"calon_suami"`
		CalonIstri       string `json:"calon_istri"`
		Status           string `json:"status"`
	}

	penghuluSchedules := make(map[uint][]ScheduleItem)
	penghuluTodayCount := make(map[uint]int)

	for _, p := range pendaftarans {
		if p.Penghulu_id == nil {
			continue
		}

		penghuluID := *p.Penghulu_id

		// Get calon names
		var calonSuami, calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		// Normalize waktu
		waktu := p.Waktu_nikah
		if len(waktu) > 5 {
			waktu = waktu[:5]
		}

		schedule := ScheduleItem{
			PendaftaranID:    p.ID,
			NomorPendaftaran: p.Nomor_pendaftaran,
			WaktuNikah:       waktu,
			TempatNikah:      p.Tempat_nikah,
			AlamatAkad:       p.Alamat_akad,
			CalonSuami:       calonSuami.Nama_lengkap,
			CalonIstri:       calonIstri.Nama_lengkap,
			Status:           p.Status_pendaftaran,
		}

		penghuluSchedules[penghuluID] = append(penghuluSchedules[penghuluID], schedule)
		penghuluTodayCount[penghuluID]++
	}

	// Build response for each penghulu
	var penghuluAvailability []gin.H
	for _, penghulu := range penghulus {
		schedules := penghuluSchedules[penghulu.ID]
		todayCount := penghuluTodayCount[penghulu.ID]

		// Build time slots availability
		scheduledTimes := make(map[string]ScheduleItem)
		for _, s := range schedules {
			scheduledTimes[s.WaktuNikah] = s
		}

		var timeSlotDetails []gin.H
		availableCount := 0
		busyCount := 0

		for _, slot := range TimeSlots {
			if schedule, exists := scheduledTimes[slot]; exists {
				// Penghulu sudah ada jadwal di jam ini
				timeSlotDetails = append(timeSlotDetails, gin.H{
					"waktu":    slot,
					"tersedia": false,
					"status":   "Bertugas",
					"jadwal":   schedule,
				})
				busyCount++
			} else {
				// Penghulu tersedia di jam ini
				timeSlotDetails = append(timeSlotDetails, gin.H{
					"waktu":    slot,
					"tersedia": true,
					"status":   "Tersedia",
					"jadwal":   nil,
				})
				availableCount++
			}
		}

		penghuluAvailability = append(penghuluAvailability, gin.H{
			"penghulu": gin.H{
				"id":           penghulu.ID,
				"nama_lengkap": penghulu.Nama_lengkap,
				"nip":          penghulu.NIP,
				"no_hp":        penghulu.No_hp,
				"email":        penghulu.Email,
				"jumlah_nikah": penghulu.Jumlah_nikah,
				"rating":       penghulu.Rating,
			},
			"tanggal":         tanggalStr,
			"jadwal_hari_ini": todayCount,
			"slot_tersedia":   availableCount,
			"slot_terisi":     busyCount,
			"time_slots":      timeSlotDetails,
			"jadwal_detail":   schedules,
		})
	}

	// Summary
	totalPenghulu := len(penghulus)
	totalJadwalHariIni := len(pendaftarans)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data ketersediaan jadwal penghulu berhasil diambil",
		"data": gin.H{
			"tanggal":               tanggalStr,
			"hari":                  tanggal.Weekday().String(),
			"tanggal_format":        tanggal.Format("02 Januari 2006"),
			"total_penghulu":        totalPenghulu,
			"total_jadwal_hari_ini": totalJadwalHariIni,
			"time_slots":            TimeSlots,
			"penghulu_availability": penghuluAvailability,
		},
	})
}

// GetPenghuluScheduleForAssignment returns penghulu availability for a specific date and time
// Used when assigning penghulu to a registration
// GET /simnikah/kepala-kua/penghulu-tersedia?tanggal=2024-12-25&waktu=09:00
func (h *InDB) GetPenghuluScheduleForAssignment(c *gin.Context) {
	tanggalStr := c.Query("tanggal")
	waktuStr := c.Query("waktu")

	if tanggalStr == "" || waktuStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tidak lengkap",
			"error":   "Parameter tanggal dan waktu diperlukan",
			"example": "/simnikah/kepala-kua/penghulu-tersedia?tanggal=2024-12-25&waktu=09:00",
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

	// Normalize waktu
	if len(waktuStr) > 5 {
		waktuStr = waktuStr[:5]
	}

	// Get all active penghulus
	var penghulus []structs.Penghulu
	if err := h.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data penghulu",
		})
		return
	}

	// Calculate start and end of day
	startOfDay := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get registrations on this date and time that have penghulu assigned
	var busyPenghuluIDs []uint
	err = h.DB.Model(&structs.PendaftaranNikah{}).
		Select("penghulu_id").
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND waktu_nikah = ? AND penghulu_id IS NOT NULL AND status_pendaftaran NOT IN ?",
			startOfDay, endOfDay, waktuStr, []string{structs.StatusPendaftaranDitolak}).
		Pluck("penghulu_id", &busyPenghuluIDs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengambil data jadwal",
		})
		return
	}

	// Create map of busy penghulus
	busyMap := make(map[uint]bool)
	for _, id := range busyPenghuluIDs {
		busyMap[id] = true
	}

	// Get today's schedule count for each penghulu
	type PenghuluCount struct {
		PenghuluID uint
		Count      int64
	}
	var penghuluCounts []PenghuluCount
	h.DB.Model(&structs.PendaftaranNikah{}).
		Select("penghulu_id, COUNT(*) as count").
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND penghulu_id IS NOT NULL AND status_pendaftaran NOT IN ?",
			startOfDay, endOfDay, []string{structs.StatusPendaftaranDitolak}).
		Group("penghulu_id").
		Scan(&penghuluCounts)

	countMap := make(map[uint]int64)
	for _, pc := range penghuluCounts {
		countMap[pc.PenghuluID] = pc.Count
	}

	// Build response
	var tersedia []gin.H
	var tidakTersedia []gin.H

	for _, penghulu := range penghulus {
		jadwalHariIni := countMap[penghulu.ID]
		isBusy := busyMap[penghulu.ID]

		penghuluData := gin.H{
			"id":              penghulu.ID,
			"nama_lengkap":    penghulu.Nama_lengkap,
			"nip":             penghulu.NIP,
			"no_hp":           penghulu.No_hp,
			"email":           penghulu.Email,
			"jumlah_nikah":    penghulu.Jumlah_nikah,
			"rating":          penghulu.Rating,
			"jadwal_hari_ini": jadwalHariIni,
		}

		if isBusy {
			penghuluData["status"] = "Tidak Tersedia"
			penghuluData["alasan"] = fmt.Sprintf("Sudah ada jadwal pada pukul %s", waktuStr)
			tidakTersedia = append(tidakTersedia, penghuluData)
		} else {
			penghuluData["status"] = "Tersedia"
			tersedia = append(tersedia, penghuluData)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data ketersediaan penghulu berhasil diambil",
		"data": gin.H{
			"tanggal":                 tanggalStr,
			"waktu":                   waktuStr,
			"hari":                    tanggal.Weekday().String(),
			"tanggal_format":          tanggal.Format("02 Januari 2006"),
			"total_penghulu":          len(penghulus),
			"jumlah_tersedia":         len(tersedia),
			"jumlah_tidak_tersedia":   len(tidakTersedia),
			"penghulu_tersedia":       tersedia,
			"penghulu_tidak_tersedia": tidakTersedia,
		},
	})
}
