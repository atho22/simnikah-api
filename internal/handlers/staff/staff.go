package staff

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"simnikah/internal/models"
	"simnikah/internal/services"
	"simnikah/pkg/crypto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== STAFF MANAGEMENT HANDLERS ====================

// CreateStaff creates a new staff member
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Validate jabatan
	validJabatan := map[string]bool{
		structs.StaffJabatanStaff:     true,
		structs.StaffJabatanPenghulu:  true,
		structs.StaffJabatanKepalaKUA: true,
	}
	if !validJabatan[input.Jabatan] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jabatan tidak valid"})
		return
	}

	// Check if username/email already exists
	var existingUser structs.Users
	if err := h.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau email sudah digunakan"})
		return
	}

	// Check if NIP already exists
	var existingStaff structs.StaffKUA
	if err := h.DB.Where("nip = ?", input.NIP).First(&existingStaff).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP sudah terdaftar"})
		return
	}

	// Generate user_id
	userID := "STF" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password
	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Create user account
	user := structs.Users{
		User_id:    userID,
		Username:   input.Username,
		Email:      input.Email,
		Password:   hashedPassword,
		Role:       structs.UserRoleStaff,
		Status:     structs.UserStatusAktif,
		Nama:       input.Nama,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user account"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat profile staff"})
		return
	}

	// Kirim notifikasi otomatis setelah staff berhasil dibuat
	notificationService := services.NewNotificationService(h.DB)
	err = notificationService.SendStaffCreatedNotification(userID, input.Nama, input.Jabatan)
	if err != nil {
		// Log error tapi jangan return error karena staff sudah berhasil dibuat
		fmt.Printf("Gagal mengirim notifikasi pembuatan staff: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff KUA berhasil dibuat",
		"data": gin.H{
			"user":  user,
			"staff": staff,
		},
	})
}

// ListStaff gets all staff members
func (h *InDB) ListStaff(c *gin.Context) {
	var staff []structs.StaffKUA

	if err := h.DB.Find(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data staff berhasil diambil",
		"data":    staff,
	})
}

// UpdateStaff updates staff information
func (h *InDB) UpdateStaff(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Nama_lengkap string `json:"nama_lengkap"`
		Jabatan      string `json:"jabatan"`
		Bagian       string `json:"bagian"`
		No_hp        string `json:"no_hp"`
		Email        string `json:"email"`
		Alamat       string `json:"alamat"`
		Status       string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	var staff structs.StaffKUA
	if err := h.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff tidak ditemukan"})
		return
	}

	// Update fields if provided
	if input.Nama_lengkap != "" {
		staff.Nama_lengkap = input.Nama_lengkap
	}
	if input.Jabatan != "" {
		staff.Jabatan = input.Jabatan
	}
	if input.Bagian != "" {
		staff.Bagian = input.Bagian
	}
	if input.No_hp != "" {
		staff.No_hp = input.No_hp
	}
	if input.Email != "" {
		staff.Email = input.Email
	}
	if input.Alamat != "" {
		staff.Alamat = input.Alamat
	}
	if input.Status != "" {
		staff.Status = input.Status
	}

	staff.Updated_at = time.Now()

	if err := h.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data staff berhasil diupdate",
		"data":    staff,
	})
}

// ==================== PENGHULU MANAGEMENT HANDLERS ====================

// CreateMarriageOfficer creates a new marriage officer (penghulu)
func (h *InDB) CreateMarriageOfficer(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nama     string `json:"nama" binding:"required"`
		NIP      string `json:"nip" binding:"required"`
		No_hp    string `json:"no_hp"`
		Alamat   string `json:"alamat"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Check if username/email already exists
	var existingUser structs.Users
	if err := h.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau email sudah digunakan"})
		return
	}

	// Check if NIP already exists
	var existingPenghulu structs.Penghulu
	if err := h.DB.Where("nip = ?", input.NIP).First(&existingPenghulu).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP sudah terdaftar"})
		return
	}

	// Generate user_id
	userID := "PNG" + fmt.Sprintf("%d", time.Now().Unix())

	// Hash password
	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user account"})
		return
	}

	// Create penghulu profile
	penghulu := structs.Penghulu{
		User_id:      userID,
		NIP:          input.NIP,
		Nama_lengkap: input.Nama,
		No_hp:        input.No_hp,
		Email:        input.Email,
		Alamat:       input.Alamat,
		Status:       structs.PenghuluStatusAktif,
		Jumlah_nikah: 0,
		Rating:       0.0,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	if err := h.DB.Create(&penghulu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat profile penghulu"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Penghulu berhasil dibuat",
		"data": gin.H{
			"user":     user,
			"penghulu": penghulu,
		},
	})
}

// ListMarriageOfficers gets all marriage officers (penghulu)
func (h *InDB) ListMarriageOfficers(c *gin.Context) {
	var penghulu []structs.Penghulu

	if err := h.DB.Find(&penghulu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penghulu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penghulu berhasil diambil",
		"data":    penghulu,
	})
}

// UpdateMarriageOfficer updates marriage officer (penghulu) information
func (h *InDB) UpdateMarriageOfficer(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Nama_lengkap string  `json:"nama_lengkap"`
		No_hp        string  `json:"no_hp"`
		Email        string  `json:"email"`
		Alamat       string  `json:"alamat"`
		Status       string  `json:"status"`
		Rating       float64 `json:"rating"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	var penghulu structs.Penghulu
	if err := h.DB.First(&penghulu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penghulu tidak ditemukan"})
		return
	}

	// Update fields if provided
	if input.Nama_lengkap != "" {
		penghulu.Nama_lengkap = input.Nama_lengkap
	}
	if input.No_hp != "" {
		penghulu.No_hp = input.No_hp
	}
	if input.Email != "" {
		penghulu.Email = input.Email
	}
	if input.Alamat != "" {
		penghulu.Alamat = input.Alamat
	}
	if input.Status != "" {
		penghulu.Status = input.Status
	}
	if input.Rating > 0 {
		penghulu.Rating = input.Rating
	}

	penghulu.Updated_at = time.Now()

	if err := h.DB.Save(&penghulu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data penghulu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penghulu berhasil diupdate",
		"data":    penghulu,
	})
}

// ==================== MARRIAGE REGISTRATION VERIFICATION ====================

// VerifyRegistrationForm verifies the registration form data by staff
func (h *InDB) VerifyRegistrationForm(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id from context (staff who is verifying)
	staffID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	var input struct {
		Status  string `json:"status" binding:"required"` // "Formulir Disetujui" or "Formulir Ditolak"
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
	if input.Status != "Formulir Disetujui" && input.Status != "Formulir Ditolak" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak valid",
			"error":   "Status harus 'Formulir Disetujui' atau 'Formulir Ditolak'",
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

	// Check if registration is in correct status for form verification
	// Flow sederhana: Draft bisa diverifikasi
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Draft' untuk diverifikasi",
		})
		return
	}

	// Update registration status
	// Jika formulir disetujui, langsung ubah ke "Disetujui"
	if input.Status == "Formulir Disetujui" {
		pendaftaran.Status_pendaftaran = structs.StatusPendaftaranDisetujui
	} else {
		pendaftaran.Status_pendaftaran = input.Status
	}

	pendaftaran.Catatan = input.Catatan
	pendaftaran.Disetujui_oleh = staffID.(string)
	now := time.Now()
	pendaftaran.Disetujui_pada = &now
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
	if input.Status == "Formulir Disetujui" {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Formulir Disetujui - Silakan Kumpulkan Berkas",
			Pesan:       "Formulir pendaftaran nikah Anda telah disetujui. Silakan datang ke kantor KUA dalam 5 hari kerja dengan membawa berkas yang diperlukan untuk melengkapi pendaftaran.",
			Tipe:        structs.NotifikasiTipeSuccess,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        "/pendaftaran/" + registrationID,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
	} else {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Formulir Pendaftaran Ditolak",
			Pesan:       "Formulir pendaftaran nikah Anda ditolak. " + input.Catatan,
			Tipe:        structs.NotifikasiTipeError,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        "/pendaftaran/" + registrationID,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		// Log error but don't fail the main operation
	}

	// Set message berdasarkan status
	var message string
	if input.Status == "Formulir Disetujui" {
		message = "Formulir berhasil disetujui dan status diubah ke Pengumpulan Berkas"
	} else {
		message = "Formulir berhasil diverifikasi"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
			"status_pendaftaran": pendaftaran.Status_pendaftaran,
			"disetujui_oleh":     pendaftaran.Disetujui_oleh,
			"disetujui_pada":     pendaftaran.Disetujui_pada,
			"catatan":            pendaftaran.Catatan,
			"updated_at":         pendaftaran.Updated_at,
		},
	})
}

// VerifyDocuments verifies the physical documents by staff
func (h *InDB) VerifyDocuments(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id from context (staff who is verifying)
	staffID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	var input struct {
		Status  string `json:"status" binding:"required"` // "Berkas Diterima" or "Berkas Ditolak"
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
	if input.Status != "Berkas Diterima" && input.Status != "Berkas Ditolak" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak valid",
			"error":   "Status harus 'Berkas Diterima' atau 'Berkas Ditolak'",
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
	// Dalam flow sederhana, verifikasi berkas bisa dilakukan pada status Draft
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Draft' untuk verifikasi berkas",
		})
		return
	}

	// Update registration status
	// Jika berkas diterima, ubah ke "Disetujui" (flow sederhana)
	if input.Status == "Berkas Diterima" {
		pendaftaran.Status_pendaftaran = structs.StatusPendaftaranDisetujui
	} else {
		pendaftaran.Status_pendaftaran = input.Status
	}
	pendaftaran.Catatan = input.Catatan
	pendaftaran.Disetujui_oleh = staffID.(string)
	now := time.Now()
	pendaftaran.Disetujui_pada = &now
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
	if input.Status == "Berkas Diterima" {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Berkas Diterima",
			Pesan:       "Berkas fisik Anda telah diterima. Pendaftaran akan diproses untuk verifikasi oleh penghulu.",
			Tipe:        structs.NotifikasiTipeSuccess,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        "/pendaftaran/" + registrationID,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
	} else {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Berkas Ditolak",
			Pesan:       "Berkas fisik Anda ditolak. " + input.Catatan,
			Tipe:        structs.NotifikasiTipeError,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
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
		"message": "Berkas berhasil diverifikasi",
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
			"status_pendaftaran": pendaftaran.Status_pendaftaran,
			"disetujui_oleh":     pendaftaran.Disetujui_oleh,
			"disetujui_pada":     pendaftaran.Disetujui_pada,
			"catatan":            pendaftaran.Catatan,
			"updated_at":         pendaftaran.Updated_at,
		},
	})
}

// ==================== APPROVE PENDAFTARAN (FLOW SEDERHANA) ====================

// ApproveRegistration approves the registration after staff verifies documents offline
// Flow: Draft → Disetujui
func (h *InDB) ApproveRegistration(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id from context (staff who is approving)
	staffID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return
	}

	var input struct {
		Status  string `json:"status" binding:"required"` // "Disetujui" or "Ditolak"
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
	if input.Status != structs.StatusPendaftaranDisetujui && input.Status != structs.StatusPendaftaranDitolak {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak valid",
			"error":   "Status harus 'Disetujui' atau 'Ditolak'",
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

	// Check if registration is in Draft status
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Draft' untuk disetujui",
		})
		return
	}

	// Update registration status
	pendaftaran.Status_pendaftaran = input.Status
	pendaftaran.Catatan = input.Catatan
	pendaftaran.Disetujui_oleh = staffID.(string)
	now := time.Now()
	pendaftaran.Disetujui_pada = &now
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
	if input.Status == structs.StatusPendaftaranDisetujui {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Pendaftaran Disetujui",
			Pesan:       "Pendaftaran nikah Anda telah disetujui. Menunggu penugasan penghulu oleh kepala KUA.",
			Tipe:        structs.NotifikasiTipeSuccess,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        "/pendaftaran/" + registrationID,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
	} else {
		notification = structs.Notifikasi{
			User_id:     pendaftaran.Pendaftar_id,
			Judul:       "Pendaftaran Ditolak",
			Pesan:       "Pendaftaran nikah Anda ditolak. " + input.Catatan,
			Tipe:        structs.NotifikasiTipeError,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
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
		"message": "Pendaftaran berhasil disetujui",
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
			"status_pendaftaran": pendaftaran.Status_pendaftaran,
			"disetujui_oleh":     pendaftaran.Disetujui_oleh,
			"disetujui_pada":     pendaftaran.Disetujui_pada,
			"catatan":            pendaftaran.Catatan,
			"updated_at":         pendaftaran.Updated_at,
		},
	})
}

// ==================== FLEKSIBEL STATUS UPDATE ====================

// UpdateRegistrationStatus - Updates registration status flexibly without strict validation
// Can be used by Staff, Penghulu, and Kepala KUA for manual status updates
func (h *InDB) UpdateRegistrationStatus(c *gin.Context) {
	registrationID := c.Param("id")

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

	// Get user role from context
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "Role tidak ditemukan",
		})
		return
	}

	// Validasi role - hanya staff, penghulu, atau kepala_kua yang bisa update
	allowedRoles := map[string]bool{
		structs.UserRoleStaff:     true,
		structs.UserRolePenghulu: true,
		structs.UserRoleKepalaKUA: true,
	}
	if !allowedRoles[userRole.(string)] {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Hanya Staff, Penghulu, atau Kepala KUA yang dapat mengupdate status",
		})
		return
	}

	var input struct {
		Status  string `json:"status"`
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

	// Validate status is not empty (manual validation)
	if strings.TrimSpace(input.Status) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   "Field 'status' wajib diisi",
		})
		return
	}

	// Status yang terkait dengan assign penghulu hanya bisa diubah oleh Kepala KUA
	// melalui endpoint khusus assign-penghulu
	penghuluRelatedStatuses := map[string]bool{
		structs.StatusPendaftaranMenungguPenugasan:  true,
		structs.StatusPendaftaranPenghuluDitugaskan: true,
	}

	if penghuluRelatedStatuses[input.Status] {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Status '" + input.Status + "' hanya bisa diubah oleh Kepala KUA melalui endpoint assign-penghulu. Gunakan endpoint POST /simnikah/pendaftaran/:id/assign-penghulu untuk menugaskan penghulu.",
		})
		return
	}

	// Validasi status yang diizinkan (flow sederhana)
	validStatuses := map[string]bool{
		structs.StatusPendaftaranDraft:     true,
		structs.StatusPendaftaranDisetujui: true,
		structs.StatusPendaftaranSelesai:   true,
		structs.StatusPendaftaranDitolak:   true,
	}

	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak valid",
			"error":   "Status yang diizinkan: Draft, Menunggu Verifikasi, Menunggu Pengumpulan Berkas, Berkas Diterima, Menunggu Bimbingan, Sudah Bimbingan, Selesai, Ditolak. Untuk status terkait penghulu, gunakan endpoint assign-penghulu.",
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

	// Simpan status lama untuk logging
	statusLama := pendaftaran.Status_pendaftaran

	// Update status tanpa validasi ketat (fleksibel)
	pendaftaran.Status_pendaftaran = input.Status
	pendaftaran.Catatan = input.Catatan
	pendaftaran.Disetujui_oleh = userID.(string)
	now := time.Now()
	pendaftaran.Disetujui_pada = &now
	pendaftaran.Updated_at = time.Now()

	if err := h.DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengupdate status pendaftaran",
		})
		return
	}

	// Create notification untuk user
	notification := structs.Notifikasi{
		User_id:     pendaftaran.Pendaftar_id,
		Judul:       "Status Pendaftaran Diupdate",
		Pesan:       fmt.Sprintf("Status pendaftaran Anda telah diubah dari '%s' menjadi '%s'. %s", statusLama, input.Status, input.Catatan),
		Tipe:        structs.NotifikasiTipeInfo,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        "/pendaftaran/" + registrationID,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Gagal membuat notifikasi: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status berhasil diupdate",
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
			"status_sebelumnya":   statusLama,
			"status_sekarang":     pendaftaran.Status_pendaftaran,
			"catatan":             pendaftaran.Catatan,
			"updated_by":          userID.(string),
			"updated_at":          pendaftaran.Updated_at,
		},
	})
}

// ==================== STAFF CREATE REGISTRATION FOR USER ====================

// CreateRegistrationForUser allows staff to create marriage registration on behalf of users
// This is useful for users who are not tech-savvy and need help from staff
func (h *InDB) CreateRegistrationForUser(c *gin.Context) {
	// Get staff user_id from context
	staffID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
			"type":    "authentication",
		})
		return
	}

	// Get staff info
	var staff structs.StaffKUA
	if err := h.DB.Where("user_id = ?", staffID.(string)).First(&staff).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Hanya staff yang dapat membuat pendaftaran untuk user",
			"type":    "authorization",
		})
		return
	}

	// Parse form data (same structure as user registration)
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

	// Create user account automatically for calon pengantin
	timestamp := time.Now().Unix()
	pendaftarUserID := "USR" + fmt.Sprintf("%d", timestamp)

	// Generate username from calon suami name
	usernameBase := strings.ToLower(strings.ReplaceAll(formSederhana.CalonLakiLaki.NamaDanBin, " ", ""))
	if len(usernameBase) > 15 {
		usernameBase = usernameBase[:15]
	}
	username := fmt.Sprintf("%s%d", usernameBase, timestamp%10000)
	email := fmt.Sprintf("%s@simnikah.local", username)
	defaultPassword := fmt.Sprintf("Nikah%d", timestamp%100000)
	hashedPassword, err := crypto.HashPassword(defaultPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat password default",
			"type":    "database",
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

	// Validate tempat nikah
	if formSederhana.LokasiNikah.TempatNikah != "Di KUA" && formSederhana.LokasiNikah.TempatNikah != "Di Luar KUA" {
		tx.Rollback()
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
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   "Alamat nikah wajib diisi untuk nikah di luar KUA",
				"field":   "alamat_nikah",
				"type":    "required",
			})
			return
		}
	}

	// Parse tanggal nikah
	tanggalNikah, err := time.Parse("2006-01-02", formSederhana.LokasiNikah.TanggalNikah)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal tidak benar",
			"error":   "Format tanggal harus: Tahun-Bulan-Tanggal (contoh: 2024-12-25)",
			"field":   "tanggal_nikah",
			"type":    "format",
		})
		return
	}

	// Validate that wedding date is not in the past
	if tanggalNikah.Before(time.Now().Truncate(24 * time.Hour)) {
		today := time.Now().Format("02 Januari 2006")
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tanggal tidak boleh sudah lewat",
			"error":   fmt.Sprintf("Tanggal nikah tidak boleh di masa lalu. Hari ini adalah %s.", today),
			"field":   "tanggal_nikah",
			"type":    "validation",
		})
		return
	}

	// Validate wedding time format
	_, err = time.Parse("15:04", formSederhana.LokasiNikah.WaktuNikah)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format jam tidak benar",
			"error":   "Format jam harus: Jam:Menit dengan 2 angka (contoh: 09:00)",
			"field":   "waktu_nikah",
			"type":    "format",
		})
		return
	}

	// Validate age
	if formSederhana.CalonLakiLaki.Umur < 19 {
		tx.Rollback()
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
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Umur calon perempuan minimal 19 tahun",
			"field":   "umur_perempuan",
			"type":    "validation",
		})
		return
	}

	// Validate wali nikah
	if strings.TrimSpace(formSederhana.WaliNikah.NamaDanBin) == "" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Nama wali nikah wajib diisi",
			"field":   "wali_nikah.nama_dan_bin",
			"type":    "required",
		})
		return
	}

	// Validate hubungan wali
	validHubungan := false
	for _, hubungan := range structs.ValidHubunganWali {
		if formSederhana.WaliNikah.HubunganWali == hubungan {
			validHubungan = true
			break
		}
	}

	if !validHubungan {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Hubungan wali tidak valid",
			"field":   "wali_nikah.hubungan_wali",
			"type":    "enum",
			"hubungan_valid": structs.ValidHubunganWali,
		})
		return
	}

	// Validate schedule availability (same logic as user registration)
	waktuNikahNormalized := formSederhana.LokasiNikah.WaktuNikah
	if len(waktuNikahNormalized) > 5 {
		waktuNikahNormalized = waktuNikahNormalized[:5]
	}

	var countTotalRegistrations int64
	err = tx.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah = ? AND waktu_nikah = ? AND status_pendaftaran NOT IN ?",
			tanggalNikah, waktuNikahNormalized,
			[]string{structs.StatusPendaftaranDitolak}).
		Count(&countTotalRegistrations).Error

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengecek ketersediaan jadwal",
			"type":    "database",
		})
		return
	}

	var countKUA int64
	err = tx.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah = ? AND waktu_nikah = ? AND tempat_nikah = ? AND status_pendaftaran NOT IN ?",
			tanggalNikah, waktuNikahNormalized, "Di KUA",
			[]string{structs.StatusPendaftaranDitolak}).
		Count(&countKUA).Error

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengecek ketersediaan jadwal di KUA",
			"type":    "database",
		})
		return
	}

	const maxTotalWeddings = 3

	if formSederhana.LokasiNikah.TempatNikah == "Di KUA" {
		if countKUA >= 1 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal di KUA sudah terisi",
				"error":   fmt.Sprintf("Jadwal pernikahan di KUA pada tanggal %s pukul %s sudah terisi. Silakan pilih tanggal atau jam lain.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
			})
			return
		}
	} else {
		countLuarKUA := countTotalRegistrations - countKUA
		if countTotalRegistrations >= maxTotalWeddings {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal sudah penuh",
				"error":   fmt.Sprintf("Jadwal pernikahan pada tanggal %s pukul %s sudah penuh. Maksimal 3 pernikahan per jam.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
			})
			return
		}
		if countKUA >= 1 && countLuarKUA >= 2 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal di luar KUA sudah penuh",
				"error":   fmt.Sprintf("Jadwal pernikahan di luar KUA pada tanggal %s pukul %s sudah penuh. Silakan pilih tanggal atau jam lain.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
			})
			return
		}
	}

	// Create user account
	pendaftarUser := structs.Users{
		User_id:    pendaftarUserID,
		Username:   username,
		Email:      email,
		Password:   hashedPassword,
		Role:       structs.UserRoleUserBiasa,
		Status:     structs.UserStatusAktif,
		Nama:       formSederhana.CalonLakiLaki.NamaDanBin,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := tx.Create(&pendaftarUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat akun user",
			"type":    "database",
		})
		return
	}

	createdAt := time.Now()

	// Generate unique IDs for calon pasangan
	hashGroom := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_groom_%d", pendaftarUserID, timestamp))))
	hashBride := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_bride_%d", timestamp, timestamp+1))))
	groomUserID := hashGroom[:20]
	brideUserID := hashBride[:20]
	nikGroom := fmt.Sprintf("T%s", hashGroom[:15])
	nikBride := fmt.Sprintf("T%s", hashBride[:15])

	// Calculate tanggal lahir dari umur
	now := time.Now()
	tanggalLahirGroom := now.AddDate(-formSederhana.CalonLakiLaki.Umur, 0, 0)
	tanggalLahirBride := now.AddDate(-formSederhana.CalonPerempuan.Umur, 0, 0)

	// Create calon suami
	calonSuami := structs.CalonPasangan{
		User_id:             groomUserID,
		NIK:                 nikGroom,
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
			"error":   "Gagal membuat data calon suami",
			"type":    "database",
		})
		return
	}

	// Create calon istri
	calonIstri := structs.CalonPasangan{
		User_id:             brideUserID,
		NIK:                 nikBride,
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
			"error":   "Gagal membuat data calon istri",
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
	}

	// Create marriage registration with note that it was created by staff
	// Status otomatis "Disetujui" karena staff sudah melakukan verifikasi saat input
	catatanStaff := fmt.Sprintf("Pendaftaran dibuat dan disetujui oleh staff: %s (NIP: %s) pada %s",
		staff.Nama_lengkap, staff.NIP, time.Now().Format("02 Januari 2006 15:04"))

	pendaftaranNikah := structs.PendaftaranNikah{
		Nomor_pendaftaran:   nomorPendaftaran,
		Pendaftar_id:        pendaftarUserID,
		Calon_suami_id:      fmt.Sprintf("%d", calonSuami.ID),
		Calon_istri_id:      fmt.Sprintf("%d", calonIstri.ID),
		Tanggal_pendaftaran: createdAt,
		Tanggal_nikah:       tanggalNikah,
		Waktu_nikah:         formSederhana.LokasiNikah.WaktuNikah,
		Tempat_nikah:        formSederhana.LokasiNikah.TempatNikah,
		Alamat_akad:         alamatAkad,
		Latitude:            latitude,
		Longitude:           longitude,
		Status_pendaftaran:  structs.StatusPendaftaranDisetujui, // Otomatis disetujui karena dibuat oleh staff
		Disetujui_oleh:      staffID.(string), // Staff yang membuat pendaftaran
		Disetujui_pada:      &createdAt, // Waktu disetujui = waktu dibuat
		Catatan:             catatanStaff,
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

	// Create wali nikah
	waliNikah := structs.WaliNikah{
		Pendaftaran_id: pendaftaranNikah.ID,
		Nama_dan_bin:   formSederhana.WaliNikah.NamaDanBin,
		Hubungan_wali:  formSederhana.WaliNikah.HubunganWali,
		Created_at:     createdAt,
		Updated_at:     createdAt,
	}

	if err := tx.Create(&waliNikah).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal membuat data wali nikah",
			"type":    "database",
		})
		return
	}

	// Update pendaftaran dengan wali_nikah_id
	pendaftaranNikah.Wali_nikah_id = &waliNikah.ID
	if err := tx.Save(&pendaftaranNikah).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengupdate pendaftaran dengan wali nikah",
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
		"message": "Pendaftaran nikah berhasil dibuat dan disetujui oleh staff",
		"data": gin.H{
			"id":                 pendaftaranNikah.ID,
			"nomor_pendaftaran":  nomorPendaftaran,
			"status_pendaftaran": pendaftaranNikah.Status_pendaftaran,
			"disetujui_oleh":     pendaftaranNikah.Disetujui_oleh,
			"disetujui_pada":     pendaftaranNikah.Disetujui_pada,
			"tanggal_nikah":      pendaftaranNikah.Tanggal_nikah,
			"waktu_nikah":        pendaftaranNikah.Waktu_nikah,
			"tempat_nikah":       pendaftaranNikah.Tempat_nikah,
			"alamat_akad":        pendaftaranNikah.Alamat_akad,
			"dibuat_oleh_staff": gin.H{
				"nama": staff.Nama_lengkap,
				"nip":  staff.NIP,
			},
			"akun_user": gin.H{
				"user_id":         pendaftarUserID,
				"username":        username,
				"email":           email,
				"password_default": defaultPassword,
				"catatan":         "Akun ini dibuat otomatis. User dapat login dan mengubah password.",
			},
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
			"wali_nikah": gin.H{
				"nama_dan_bin":   waliNikah.Nama_dan_bin,
				"hubungan_wali":  waliNikah.Hubungan_wali,
			},
			"catatan": "Pendaftaran dibuat oleh staff. User dapat login menggunakan username dan password default yang diberikan.",
		},
	})
}