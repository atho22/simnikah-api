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
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nama     string `json:"nama" binding:"required"`
		NIP      string `json:"nip" binding:"required"`
		No_hp    string `json:"no_hp"`
		EmailPenghulu string `json:"email_penghulu"`
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
				"id":          penghulu.ID,
				"nip":         penghulu.NIP,
				"nama_lengkap": penghulu.Nama_lengkap,
				"no_hp":       penghulu.No_hp,
				"email":       penghulu.Email,
				"status":      penghulu.Status,
			},
			"periode": gin.H{
				"bulan":      bulan,
				"tahun":      tahun,
				"nama_bulan": startOfMonth.Month().String(),
			},
			"statistik_keseluruhan": gin.H{
				"total_semua":       totalAllTime,
				"selesai_semua":     totalCompleted,
				"progress_semua":    totalInProgress,
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
			"registrasi_terbaru":   recentData,
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
				"id":          penghulu.ID,
				"nip":         penghulu.NIP,
				"nama_lengkap": penghulu.Nama_lengkap,
				"status":      penghulu.Status,
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
				"total_penghulu":         len(penghulus),
				"total_semua":             totalAllPenghulu,
				"selesai_semua":           totalCompletedAllPenghulu,
				"progress_semua":          totalAllPenghulu - totalCompletedAllPenghulu,
				"persentase_selesai":      completionPctAllTime,
				"total_bulan_ini":         totalThisMonthAllPenghulu,
				"selesai_bulan_ini":       completedThisMonthAllPenghulu,
				"progress_bulan_ini":      totalThisMonthAllPenghulu - completedThisMonthAllPenghulu,
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
	jenisFeedback := c.Query("jenis")      // "Rating", "Saran", "Kritik", "Laporan"
	statusBaca := c.Query("status")        // "Belum Dibaca", "Sudah Dibaca"
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
			"judul":           r.Judul,
			"pesan":           r.Pesan,
			"status_baca":     r.Status_baca,
			"dibaca_oleh":     r.Dibaca_oleh,
			"dibaca_pada":     r.Dibaca_pada,
			"tanggal_nikah":   r.TanggalNikah.Format("2006-01-02"),
			"created_at":      r.Created_at,
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
				"jenis": jenisFeedback,
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
			"id":            feedback.ID,
			"status_baca":   feedback.Status_baca,
			"dibaca_oleh":   feedback.Dibaca_oleh,
			"dibaca_pada":   feedback.Dibaca_pada,
			"updated_at":    feedback.Updated_at,
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
			"total_feedback":   totalFeedback,
			"unread_feedback":  unreadFeedback,
			"read_feedback":    totalFeedback - unreadFeedback,
			"stats_by_jenis":   statsByJenis,
			"rating_stats": gin.H{
				"average_rating": avgRating,
				"total_ratings":  ratingCount,
			},
		},
	})
}
