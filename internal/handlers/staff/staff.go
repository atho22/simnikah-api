package staff

import (
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

// CreateStaffKUA creates a new staff KUA
func (h *InDB) CreateStaffKUA(c *gin.Context) {
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

// GetAllStaff gets all staff KUA
func (h *InDB) GetAllStaff(c *gin.Context) {
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

// UpdateStaffKUA updates staff KUA information
func (h *InDB) UpdateStaffKUA(c *gin.Context) {
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

// CreatePenghulu creates a new penghulu
func (h *InDB) CreatePenghulu(c *gin.Context) {
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

// GetAllPenghulu gets all penghulu
func (h *InDB) GetAllPenghulu(c *gin.Context) {
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

// UpdatePenghulu updates penghulu information
func (h *InDB) UpdatePenghulu(c *gin.Context) {
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

// VerifyFormulir verifies the form data by staff (Tahap 1)
func (h *InDB) VerifyFormulir(c *gin.Context) {
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
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranMenungguVerifikasi {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Menunggu Verifikasi' untuk diverifikasi",
		})
		return
	}

	// Update registration status
	// Jika formulir disetujui, langsung ubah ke "Menunggu Pengumpulan Berkas"
	if input.Status == "Formulir Disetujui" {
		pendaftaran.Status_pendaftaran = structs.StatusPendaftaranMenungguPengumpulanBerkas
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

// VerifyBerkas verifies the physical documents by staff (Tahap 2)
func (h *InDB) VerifyBerkas(c *gin.Context) {
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
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranMenungguPengumpulanBerkas {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Menunggu Pengumpulan Berkas' untuk verifikasi berkas",
		})
		return
	}

	// Update registration status
	// Jika berkas diterima, ubah ke "Berkas Diterima"
	if input.Status == "Berkas Diterima" {
		pendaftaran.Status_pendaftaran = structs.StatusPendaftaranBerkasDiterima
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

// ==================== FLEKSIBEL STATUS UPDATE ====================

// UpdateStatusFlexible - Update status pendaftaran secara fleksibel tanpa validasi ketat
// Bisa digunakan oleh Staff, Penghulu, dan Kepala KUA untuk update status secara manual
func (h *InDB) UpdateStatusFlexible(c *gin.Context) {
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
		structs.StatusPendaftaranMenungguPenugasan:         true,
		structs.StatusPendaftaranPenghuluDitugaskan:         true,
		structs.StatusPendaftaranMenungguVerifikasiPenghulu: true,
	}

	if penghuluRelatedStatuses[input.Status] {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak",
			"error":   "Status '" + input.Status + "' hanya bisa diubah oleh Kepala KUA melalui endpoint assign-penghulu. Gunakan endpoint POST /simnikah/pendaftaran/:id/assign-penghulu untuk menugaskan penghulu.",
		})
		return
	}

	// Validasi status yang diizinkan (tanpa status terkait penghulu)
	validStatuses := map[string]bool{
		structs.StatusPendaftaranDraft:                     true,
		structs.StatusPendaftaranMenungguVerifikasi:        true,
		structs.StatusPendaftaranMenungguPengumpulanBerkas: true,
		structs.StatusPendaftaranBerkasDiterima:            true,
		structs.StatusPendaftaranMenungguBimbingan:         true,
		structs.StatusPendaftaranSudahBimbingan:             true,
		structs.StatusPendaftaranSelesai:                   true,
		structs.StatusPendaftaranDitolak:                   true,
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