package staff

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	structs "simnikah/internal/models"
	services "simnikah/internal/services"
	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== STAFF CRUD ====================

// ListStaff gets all staff members
func (h *InDB) ListStaff(c *gin.Context) {
	var staff []structs.StaffKUA
	if err := h.DB.Find(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data staff", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data staff berhasil diambil", "data": staff})
}

// UpdateStaff updates staff information
func (h *InDB) UpdateStaff(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Nama     string `json:"nama"`
		Jabatan  string `json:"jabatan"`
		Bagian   string `json:"bagian"`
		No_hp    string `json:"no_hp"`
		Email    string `json:"email"`
		Alamat   string `json:"alamat"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format data tidak valid", "error": err.Error()})
		return
	}

	var staff structs.StaffKUA
	if err := h.DB.Where("id = ?", id).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Staff tidak ditemukan"})
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if input.Nama != "" {
		updates["nama_lengkap"] = input.Nama
	}
	if input.Jabatan != "" {
		updates["jabatan"] = input.Jabatan
	}
	if input.Bagian != "" {
		updates["bagian"] = input.Bagian
	}
	if input.No_hp != "" {
		updates["no_hp"] = input.No_hp
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Alamat != "" {
		updates["alamat"] = input.Alamat
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}

	if err := h.DB.Model(&staff).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengupdate staff", "error": err.Error()})
		return
	}

	h.DB.Where("id = ?", id).First(&staff)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Staff berhasil diupdate", "data": staff})
}

// ==================== CREATE REGISTRATION FOR USER (SCHEDULING-ONLY) ====================

// CreateRegistrationForUser allows staff/kepala_kua to create a registration on behalf of a catin.
// Simplified to scheduling-only data (no CalonPasangan/WaliNikah).
func (h *InDB) CreateRegistrationForUser(c *gin.Context) {
	var input struct {
		NamaSuami    string   `json:"nama_suami" binding:"required"`
		UmurSuami    int      `json:"umur_suami"`
		NamaIstri    string   `json:"nama_istri" binding:"required"`
		UmurIstri    int      `json:"umur_istri"`
		TanggalNikah string   `json:"tanggal_nikah" binding:"required"`
		WaktuNikah   string   `json:"waktu_nikah" binding:"required"`
		TempatNikah  string   `json:"tempat_nikah" binding:"required"`
		AlamatAkad   string   `json:"alamat_akad"`
		Latitude     *float64 `json:"latitude"`
		Longitude    *float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format data tidak valid", "error": err.Error()})
		return
	}

	if input.TempatNikah != structs.TempatNikahDiKUA && input.TempatNikah != structs.TempatNikahDiLuarKUA {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tempat_nikah tidak valid", "error": "Harus 'Di KUA' atau 'Di Luar KUA'"})
		return
	}

	tanggalNikah, err := time.ParseInLocation("2006-01-02", input.TanggalNikah, utils.WITA)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format tanggal tidak valid", "error": "Gunakan format YYYY-MM-DD"})
		return
	}

	// Cek ketersediaan jadwal via Forward Chaining Engine (sama seperti catin handler)
	engine := services.NewForwardChainingEngine(h.DB)
	scheduleResult, scheduleErr := engine.CheckScheduleAvailability(services.ScheduleCheckInput{
		TanggalNikah: input.TanggalNikah,
		WaktuNikah:   strings.TrimSpace(input.WaktuNikah),
		TempatNikah:  input.TempatNikah,
		AlamatNikah:  input.AlamatAkad,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
	})
	if scheduleErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengecek ketersediaan jadwal", "error": scheduleErr.Error()})
		return
	}
	if !scheduleResult.Available {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Jadwal tidak tersedia", "error": scheduleResult.Reason})
		return
	}

	registration := structs.PendaftaranNikah{
		Nama_suami:         input.NamaSuami,
		Umur_suami:         input.UmurSuami,
		Nama_istri:         input.NamaIstri,
		Umur_istri:         input.UmurIstri,
		Tanggal_nikah:      tanggalNikah,
		Waktu_nikah:        input.WaktuNikah,
		Tempat_nikah:       input.TempatNikah,
		Alamat_akad:        input.AlamatAkad,
		Latitude:           input.Latitude,
		Longitude:          input.Longitude,
		Status_pendaftaran: structs.StatusPendaftaranMenungguPenugasan,
		Created_at:         time.Now().In(utils.WITA),
		Updated_at:         time.Now().In(utils.WITA),
	}

	if err := h.DB.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat pendaftaran", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pendaftaran berhasil dibuat",
		"data":    registration,
	})
}

// ==================== UPDATE REGISTRATION STATUS ====================

// UpdateRegistrationStatus updates the status of a registration.
// Simplified status flow: Menunggu Penugasan -> Penghulu Ditugaskan -> Selesai / Ditolak
func (h *InDB) UpdateRegistrationStatus(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Status string `json:"status_pendaftaran" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format data tidak valid", "error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		structs.StatusPendaftaranMenungguPenugasan:  true,
		structs.StatusPendaftaranPenghuluDitugaskan: true,
		structs.StatusPendaftaranSelesai:            true,
		structs.StatusPendaftaranDitolak:            true,
	}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak valid",
			"error":   "Status harus salah satu: Menunggu Penugasan, Penghulu Ditugaskan, Selesai, Ditolak",
		})
		return
	}

	var registration structs.PendaftaranNikah
	if err := h.DB.Where("id = ?", id).First(&registration).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Pendaftaran tidak ditemukan"})
		return
	}

	// Validasi transisi status (state machine)
	// Menunggu Penugasan -> Penghulu Ditugaskan / Ditolak
	// Penghulu Ditugaskan -> Selesai
	var validTransitions = map[string]map[string]bool{
		structs.StatusPendaftaranMenungguPenugasan: {
			structs.StatusPendaftaranPenghuluDitugaskan: true,
			structs.StatusPendaftaranDitolak:            true,
		},
		structs.StatusPendaftaranPenghuluDitugaskan: {
			structs.StatusPendaftaranSelesai: true,
		},
	}

	allowedTargets, ok := validTransitions[registration.Status_pendaftaran]
	if !ok || !allowedTargets[input.Status] {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Transisi status tidak valid",
			"error":   fmt.Sprintf("Tidak dapat mengubah status dari '%s' ke '%s'", registration.Status_pendaftaran, input.Status),
		})
		return
	}

	if err := h.DB.Model(&registration).Updates(map[string]interface{}{
		"status_pendaftaran": input.Status,
		"updated_at":         time.Now().In(utils.WITA),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengupdate status", "error": err.Error()})
		return
	}

	h.DB.Where("id = ?", id).First(&registration)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status berhasil diupdate", "data": registration})
}

