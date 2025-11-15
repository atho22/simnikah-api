package catin

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "simnikah/internal/models"
	"simnikah/internal/services"
	"simnikah/pkg/cache"
	"simnikah/pkg/utils"
	"simnikah/pkg/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== NEW MARRIAGE REGISTRATION FORM HANDLER ====================

// CreateMarriageRegistrationForm creates a complete marriage registration using the new form structure
func (h *InDB) CreateMarriageRegistrationForm(c *gin.Context) {
	var dataFormPendaftaran structs.DataFormPendaftaranNikah

	if err := c.ShouldBindJSON(&dataFormPendaftaran); err != nil {
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

	// Parse dates
	tanggalNikah, err := time.Parse("2006-01-02", dataFormPendaftaran.JadwalDanLokasi.TanggalNikah)
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
	_, err = time.Parse("15:04", dataFormPendaftaran.JadwalDanLokasi.WaktuNikah)
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

	tanggalLahirSuami, err := time.Parse("2006-01-02", dataFormPendaftaran.CalonSuami.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Format tanggal lahir suami tidak valid (YYYY-MM-DD)",
			"field":   "tanggal_lahir_suami",
			"type":    "format",
		})
		return
	}

	tanggalLahirIstri, err := time.Parse("2006-01-02", dataFormPendaftaran.CalonIstri.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Format tanggal lahir istri tidak valid (YYYY-MM-DD)",
			"field":   "tanggal_lahir_istri",
			"type":    "format",
		})
		return
	}

	// --- DISPENSASI VALIDATION LOGIC START ---
	// Calculate working days between registration date (now) and wedding date
	now := time.Now()
	workingDays := utils.CalculateWorkingDays(now, tanggalNikah)

	// Calculate ages of bride and groom at registration date
	groomAge := utils.CalculateAge(tanggalLahirSuami, now)
	brideAge := utils.CalculateAge(tanggalLahirIstri, now)

	// Check if dispensation is required
	requiresDispensation := false
	var dispensationReasons []string

	// Check if wedding is less than 10 working days from registration
	if workingDays < 10 {
		requiresDispensation = true
		dispensationReasons = append(dispensationReasons, "Pelaksanaan nikah kurang dari 10 hari kerja")
	}

	// Check if groom is under 19 years old
	if groomAge < 19 {
		requiresDispensation = true
		dispensationReasons = append(dispensationReasons, "Calon suami berumur kurang dari 19 tahun")
	}

	// Check if bride is under 19 years old
	if brideAge < 19 {
		requiresDispensation = true
		dispensationReasons = append(dispensationReasons, "Calon istri berumur kurang dari 19 tahun")
	}

	// Validate dispensation number if required
	if requiresDispensation {
		if strings.TrimSpace(dataFormPendaftaran.JadwalDanLokasi.NomorDispensasi) == "" {
			dispensationReason := strings.Join(dispensationReasons, " dan ")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   "Nomor dispensasi wajib diisi karena: " + dispensationReason,
				"field":   "nomor_dispensasi",
				"type":    "required",
				"details": gin.H{
					"reason":       dispensationReason,
					"working_days": workingDays,
					"groom_age":    groomAge,
					"bride_age":    brideAge,
				},
			})
			return
		}
	}
	// --- DISPENSASI VALIDATION LOGIC END ---

	// Validate citizenship
	if dataFormPendaftaran.CalonSuami.Kewarganegaraan != "WNI" && dataFormPendaftaran.CalonSuami.Kewarganegaraan != "WNA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Kewarganegaraan suami harus WNI atau WNA",
			"field":   "kewarganegaraan_suami",
			"type":    "enum",
		})
		return
	}
	if dataFormPendaftaran.CalonIstri.Kewarganegaraan != "WNI" && dataFormPendaftaran.CalonIstri.Kewarganegaraan != "WNA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Kewarganegaraan istri harus WNI atau WNA",
			"field":   "kewarganegaraan_istri",
			"type":    "enum",
		})
		return
	}

	// Validate marital status
	daftarStatusPerkawinan := []string{
		structs.StatusPerkawinanBelumKawin,
		structs.StatusPerkawinanKawin,
		structs.StatusPerkawinanCeraiHidup,
		structs.StatusPerkawinanCeraiMati,
	}
	if !validator.CheckValidValue(daftarStatusPerkawinan, dataFormPendaftaran.CalonSuami.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status perkawinan suami tidak valid",
			"field":   "status_suami",
			"type":    "enum",
		})
		return
	}
	if !validator.CheckValidValue(daftarStatusPerkawinan, dataFormPendaftaran.CalonIstri.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status perkawinan istri tidak valid",
			"field":   "status_istri",
			"type":    "enum",
		})
		return
	}

	// Validate guardian status
	daftarStatusWali := []string{structs.WaliStatusKeberadaanHidup, structs.WaliStatusKeberadaanMeninggal}
	if !validator.CheckValidValue(daftarStatusWali, dataFormPendaftaran.WaliNikah.StatusWali) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status wali tidak valid",
			"field":   "status_wali",
			"type":    "enum",
		})
		return
	}

	// ==================== VALIDASI WALI NIKAH SESUAI SYARIAT ISLAM ====================
	// Validasi hubungan wali berdasarkan status ayah kandung calon istri
	statusAyahCalonIstri := dataFormPendaftaran.OrangTuaCalonIstri.Ayah.StatusKeberadaan
	hubunganWali := dataFormPendaftaran.WaliNikah.HubunganWali
	statusWali := dataFormPendaftaran.WaliNikah.StatusWali

	// Validasi 1: Wali nikah HARUS hidup untuk dapat menjadi wali
	if statusWali == structs.WaliStatusKeberadaanMeninggal {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi Wali Nikah Gagal",
			"error":   "Wali nikah yang telah meninggal dunia tidak dapat menjadi wali. Silakan pilih wali lain yang masih hidup.",
			"field":   "status_wali",
			"type":    "syariat_validation",
		})
		return
	}

	// Validasi 2: Jika hubungan wali adalah Ayah Kandung, maka status wali harus sama dengan status ayah catin perempuan
	if hubunganWali == structs.WaliHubunganAyahKandung {
		if statusWali != statusAyahCalonIstri {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi Wali Nikah Gagal",
				"error":   fmt.Sprintf("Jika wali adalah Ayah Kandung, maka status wali harus sama dengan status ayah catin perempuan (%s)", statusAyahCalonIstri),
				"field":   "status_wali",
				"type":    "syariat_validation",
			})
			return
		}

		// Validasi 3: Jika ayah kandung masih hidup dan menjadi wali, NIK wali harus sama dengan NIK ayah catin perempuan
		if statusAyahCalonIstri == structs.StatusKeberadaanHidup {
			nikAyahCalonIstri := dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Nik
			nikWali := dataFormPendaftaran.WaliNikah.NikWali

			if nikWali != nikAyahCalonIstri {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi Wali Nikah Gagal",
					"error":   "NIK wali harus sama dengan NIK ayah kandung catin perempuan",
					"field":   "nik_wali",
					"type":    "syariat_validation",
					"details": gin.H{
						"nik_wali_yang_diinput":    nikWali,
						"nik_ayah_catin_perempuan": nikAyahCalonIstri,
					},
				})
				return
			}
		}
	}

	// Validasi 4: Wali nikah tidak boleh sama dengan calon pengantin
	nikWali := dataFormPendaftaran.WaliNikah.NikWali
	nikCalonSuami := dataFormPendaftaran.CalonSuami.Nik
	nikCalonIstri := dataFormPendaftaran.CalonIstri.Nik

	if nikWali == nikCalonSuami {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi Wali Nikah Gagal",
			"error":   "NIK wali tidak boleh sama dengan NIK calon suami",
			"field":   "nik_wali",
			"type":    "syariat_validation",
		})
		return
	}

	if nikWali == nikCalonIstri {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi Wali Nikah Gagal",
			"error":   "NIK wali tidak boleh sama dengan NIK calon istri",
			"field":   "nik_wali",
			"type":    "syariat_validation",
		})
		return
	}

	// Validasi 5: Apakah hubungan wali valid berdasarkan status ayah
	if !structs.IsValidWaliNikah(hubunganWali, statusAyahCalonIstri) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi Wali Nikah Gagal",
			"error":   structs.GetPesanValidasiWali(statusAyahCalonIstri),
			"field":   "hubungan_wali",
			"type":    "syariat_validation",
			"details": gin.H{
				"status_ayah_catin_perempuan": statusAyahCalonIstri,
				"hubungan_wali_yang_dipilih":  hubunganWali,
				"urutan_wali_nasab":           structs.GetUrutanWaliNasab(),
			},
		})
		return
	}
	// ==================== END VALIDASI WALI NIKAH ====================

	// Validate wedding location
	daftarLokasiNikah := []string{"Di KUA", "Di Luar KUA"}
	if !validator.CheckValidValue(daftarLokasiNikah, dataFormPendaftaran.JadwalDanLokasi.LokasiNikah) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Lokasi nikah tidak valid",
			"field":   "lokasi_nikah",
			"type":    "enum",
		})
		return
	}

	// Validate wedding address if location is outside KUA
	if dataFormPendaftaran.JadwalDanLokasi.LokasiNikah == "Di Luar KUA" {
		if dataFormPendaftaran.JadwalDanLokasi.AlamatNikah == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   "Alamat nikah wajib diisi jika lokasi nikah di luar KUA",
				"field":   "alamat_nikah",
				"type":    "required",
			})
			return
		}
	}

	// Validate parent presence status
	daftarStatusKeberadaan := []string{structs.StatusKeberadaanHidup, structs.StatusKeberadaanMeninggal, "Tidak Diketahui"}
	if !validator.CheckValidValue(daftarStatusKeberadaan, dataFormPendaftaran.OrangTuaCalonSuami.Ayah.StatusKeberadaan) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status keberadaan ayah suami tidak valid",
			"field":   "status_keberadaan_ayah_suami",
			"type":    "enum",
		})
		return
	}
	if !validator.CheckValidValue(daftarStatusKeberadaan, dataFormPendaftaran.OrangTuaCalonSuami.Ibu.StatusKeberadaan) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status keberadaan ibu suami tidak valid",
			"field":   "status_keberadaan_ibu_suami",
			"type":    "enum",
		})
		return
	}
	if !validator.CheckValidValue(daftarStatusKeberadaan, dataFormPendaftaran.OrangTuaCalonIstri.Ayah.StatusKeberadaan) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status keberadaan ayah istri tidak valid",
			"field":   "status_keberadaan_ayah_istri",
			"type":    "enum",
		})
		return
	}
	if !validator.CheckValidValue(daftarStatusKeberadaan, dataFormPendaftaran.OrangTuaCalonIstri.Ibu.StatusKeberadaan) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validasi gagal",
			"error":   "Status keberadaan ibu istri tidak valid",
			"field":   "status_keberadaan_ibu_istri",
			"type":    "enum",
		})
		return
	}

	// Validate conditional fields for parents (only if status is "Hidup")
	validator.ValidateParentFields(dataFormPendaftaran.OrangTuaCalonSuami.Ayah, "ayah suami", c)
	validator.ValidateParentFields(dataFormPendaftaran.OrangTuaCalonSuami.Ibu, "ibu suami", c)
	validator.ValidateParentFields(dataFormPendaftaran.OrangTuaCalonIstri.Ayah, "ayah istri", c)
	validator.ValidateParentFields(dataFormPendaftaran.OrangTuaCalonIstri.Ibu, "ibu istri", c)

	// Validate conditional fields for groom and bride
	validator.ValidatePersonFields(dataFormPendaftaran.CalonSuami, "suami", c)
	validator.ValidatePersonFields(dataFormPendaftaran.CalonIstri, "istri", c)

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

	// Generate unique user IDs for groom and bride profiles (max 20 chars)
	userIDStr := userID.(string)
	timestamp := time.Now().Unix()

	// Create unique hash-based IDs that fit within 20 characters
	groomData := fmt.Sprintf("%s_SUAMI_%d", userIDStr, timestamp)
	brideData := fmt.Sprintf("%s_ISTRI_%d", userIDStr, timestamp)

	groomHash := md5.Sum([]byte(groomData))
	brideHash := md5.Sum([]byte(brideData))

	// Use first 20 characters of hash as user ID
	groomUserID := fmt.Sprintf("%x", groomHash)[:20]
	brideUserID := fmt.Sprintf("%x", brideHash)[:20]

	// OPTIMASI: Single timestamp untuk semua records
	createdAt := time.Now()

	// Start database transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create groom profile with unique user_id
	calonSuami := structs.CalonPasangan{
		User_id:             groomUserID,
		NIK:                 dataFormPendaftaran.CalonSuami.Nik,
		Nama_lengkap:        dataFormPendaftaran.CalonSuami.NamaLengkap,
		Tempat_lahir:        dataFormPendaftaran.CalonSuami.TempatLahir,
		Tanggal_lahir:       tanggalLahirSuami,
		Jenis_kelamin:       "L",
		Alamat:              dataFormPendaftaran.CalonSuami.Alamat,
		Agama:               dataFormPendaftaran.CalonSuami.Agama,
		Status_perkawinan:   dataFormPendaftaran.CalonSuami.Status,
		Pekerjaan:           dataFormPendaftaran.CalonSuami.Pekerjaan,
		Deskripsi_pekerjaan: dataFormPendaftaran.CalonSuami.DeskripsiPekerjaan,
		Pendidikan_terakhir: dataFormPendaftaran.CalonSuami.Pendidikan,
		No_hp:               dataFormPendaftaran.CalonSuami.NomorTelepon,
		Email:               dataFormPendaftaran.CalonSuami.Email,
		Warga_negara:        dataFormPendaftaran.CalonSuami.Kewarganegaraan,
		No_paspor:           dataFormPendaftaran.CalonSuami.NomorPaspor,
		Created_at:          createdAt,
		Updated_at:          createdAt,
	}

	if err := tx.Create(&calonSuami).Error; err != nil {
		tx.Rollback()
		// Check if it's a duplicate entry error
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Data sudah ada",
				"error":   "Profile calon suami sudah terdaftar untuk user ini",
				"field":   "profile_suami",
				"type":    "duplicate",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database error",
				"error":   "Gagal membuat profile calon suami: " + err.Error(),
				"type":    "database",
			})
		}
		return
	}

	// Create bride profile with unique user_id
	calonIstri := structs.CalonPasangan{
		User_id:             brideUserID,
		NIK:                 dataFormPendaftaran.CalonIstri.Nik,
		Nama_lengkap:        dataFormPendaftaran.CalonIstri.NamaLengkap,
		Tempat_lahir:        dataFormPendaftaran.CalonIstri.TempatLahir,
		Tanggal_lahir:       tanggalLahirIstri,
		Jenis_kelamin:       "P",
		Alamat:              dataFormPendaftaran.CalonIstri.Alamat,
		Agama:               dataFormPendaftaran.CalonIstri.Agama,
		Status_perkawinan:   dataFormPendaftaran.CalonIstri.Status,
		Pekerjaan:           dataFormPendaftaran.CalonIstri.Pekerjaan,
		Deskripsi_pekerjaan: dataFormPendaftaran.CalonIstri.DeskripsiPekerjaan,
		Pendidikan_terakhir: dataFormPendaftaran.CalonIstri.Pendidikan,
		No_hp:               dataFormPendaftaran.CalonIstri.NomorTelepon,
		Email:               dataFormPendaftaran.CalonIstri.Email,
		Warga_negara:        dataFormPendaftaran.CalonIstri.Kewarganegaraan,
		No_paspor:           dataFormPendaftaran.CalonIstri.NomorPaspor,
		Created_at:          createdAt,
		Updated_at:          createdAt,
	}

	if err := tx.Create(&calonIstri).Error; err != nil {
		tx.Rollback()
		// Check if it's a duplicate entry error
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Data sudah ada",
				"error":   "Profile calon istri sudah terdaftar untuk user ini",
				"field":   "profile_istri",
				"type":    "duplicate",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database error",
				"error":   "Gagal membuat profile calon istri: " + err.Error(),
				"type":    "database",
			})
		}
		return
	}

	// OPTIMASI: Batch insert untuk semua orang tua sekaligus
	var dataOrangTuaList []structs.DataOrangTua

	// Collect all parents data (only if presence status is "Hidup")
	if dataFormPendaftaran.OrangTuaCalonSuami.Ayah.StatusKeberadaan == structs.StatusKeberadaanHidup {
		dataOrangTuaList = append(dataOrangTuaList, structs.DataOrangTua{
			User_id:             userID.(string),
			Jenis_kelamin_calon: "L",
			Hubungan:            structs.HubunganAyah,
			NIK:                 dataFormPendaftaran.OrangTuaCalonSuami.Ayah.Nik,
			Nama_lengkap:        dataFormPendaftaran.OrangTuaCalonSuami.Ayah.Nama,
			Warga_negara:        dataFormPendaftaran.OrangTuaCalonSuami.Ayah.Kewarganegaraan,
			Agama:               dataFormPendaftaran.OrangTuaCalonSuami.Ayah.Agama,
			Tempat_lahir:        dataFormPendaftaran.OrangTuaCalonSuami.Ayah.TempatLahir,
			Negara_asal:         dataFormPendaftaran.OrangTuaCalonSuami.Ayah.NegaraAsal,
			Pekerjaan:           dataFormPendaftaran.OrangTuaCalonSuami.Ayah.Pekerjaan,
			Pekerjaan_lain:      dataFormPendaftaran.OrangTuaCalonSuami.Ayah.DeskripsiPekerjaan,
			Alamat:              dataFormPendaftaran.OrangTuaCalonSuami.Ayah.Alamat,
			Status_keberadaan:   structs.StatusKeberadaanHidup,
			Jenis_kelamin:       "L",
			Created_at:          createdAt,
			Updated_at:          createdAt,
		})
	}

	if dataFormPendaftaran.OrangTuaCalonSuami.Ibu.StatusKeberadaan == structs.StatusKeberadaanHidup {
		dataOrangTuaList = append(dataOrangTuaList, structs.DataOrangTua{
			User_id:             userID.(string),
			Jenis_kelamin_calon: "L",
			Hubungan:            structs.HubunganIbu,
			NIK:                 dataFormPendaftaran.OrangTuaCalonSuami.Ibu.Nik,
			Nama_lengkap:        dataFormPendaftaran.OrangTuaCalonSuami.Ibu.Nama,
			Warga_negara:        dataFormPendaftaran.OrangTuaCalonSuami.Ibu.Kewarganegaraan,
			Agama:               dataFormPendaftaran.OrangTuaCalonSuami.Ibu.Agama,
			Tempat_lahir:        dataFormPendaftaran.OrangTuaCalonSuami.Ibu.TempatLahir,
			Negara_asal:         dataFormPendaftaran.OrangTuaCalonSuami.Ibu.NegaraAsal,
			Pekerjaan:           dataFormPendaftaran.OrangTuaCalonSuami.Ibu.Pekerjaan,
			Pekerjaan_lain:      dataFormPendaftaran.OrangTuaCalonSuami.Ibu.DeskripsiPekerjaan,
			Alamat:              dataFormPendaftaran.OrangTuaCalonSuami.Ibu.Alamat,
			Status_keberadaan:   structs.StatusKeberadaanHidup,
			Jenis_kelamin:       "P",
			Created_at:          createdAt,
			Updated_at:          createdAt,
		})
	}

	if dataFormPendaftaran.OrangTuaCalonIstri.Ayah.StatusKeberadaan == structs.StatusKeberadaanHidup {
		dataOrangTuaList = append(dataOrangTuaList, structs.DataOrangTua{
			User_id:             userID.(string),
			Jenis_kelamin_calon: "P",
			Hubungan:            structs.HubunganAyah,
			NIK:                 dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Nik,
			Nama_lengkap:        dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Nama,
			Warga_negara:        dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Kewarganegaraan,
			Agama:               dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Agama,
			Tempat_lahir:        dataFormPendaftaran.OrangTuaCalonIstri.Ayah.TempatLahir,
			Negara_asal:         dataFormPendaftaran.OrangTuaCalonIstri.Ayah.NegaraAsal,
			Pekerjaan:           dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Pekerjaan,
			Pekerjaan_lain:      dataFormPendaftaran.OrangTuaCalonIstri.Ayah.DeskripsiPekerjaan,
			Alamat:              dataFormPendaftaran.OrangTuaCalonIstri.Ayah.Alamat,
			Status_keberadaan:   structs.StatusKeberadaanHidup,
			Jenis_kelamin:       "L",
			Created_at:          createdAt,
			Updated_at:          createdAt,
		})
	}

	if dataFormPendaftaran.OrangTuaCalonIstri.Ibu.StatusKeberadaan == structs.StatusKeberadaanHidup {
		dataOrangTuaList = append(dataOrangTuaList, structs.DataOrangTua{
			User_id:             userID.(string),
			Jenis_kelamin_calon: "P",
			Hubungan:            structs.HubunganIbu,
			NIK:                 dataFormPendaftaran.OrangTuaCalonIstri.Ibu.Nik,
			Nama_lengkap:        dataFormPendaftaran.OrangTuaCalonIstri.Ibu.Nama,
			Warga_negara:        dataFormPendaftaran.OrangTuaCalonIstri.Ibu.Kewarganegaraan,
			Agama:               dataFormPendaftaran.OrangTuaCalonIstri.Ibu.Agama,
			Tempat_lahir:        dataFormPendaftaran.OrangTuaCalonIstri.Ibu.TempatLahir,
			Negara_asal:         dataFormPendaftaran.OrangTuaCalonIstri.Ibu.NegaraAsal,
			Pekerjaan:           dataFormPendaftaran.OrangTuaCalonIstri.Ibu.Pekerjaan,
			Pekerjaan_lain:      dataFormPendaftaran.OrangTuaCalonIstri.Ibu.DeskripsiPekerjaan,
			Alamat:              dataFormPendaftaran.OrangTuaCalonIstri.Ibu.Alamat,
			Status_keberadaan:   structs.StatusKeberadaanHidup,
			Jenis_kelamin:       "P",
			Created_at:          createdAt,
			Updated_at:          createdAt,
		})
	}

	// OPTIMASI: Batch insert semua orang tua sekaligus jika ada
	if len(dataOrangTuaList) > 0 {
		if err := tx.CreateInBatches(&dataOrangTuaList, len(dataOrangTuaList)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database error",
				"error":   "Gagal membuat data orang tua",
				"type":    "database",
			})
			return
		}
	}

	// Generate registration number
	nomorPendaftaran := utils.GenerateUserID("NIK")

	// Determine wedding address based on location
	var alamatAkad string
	var latitude, longitude *float64

	if dataFormPendaftaran.JadwalDanLokasi.LokasiNikah == "Di KUA" {
		alamatAkad = "KUA Kecamatan Banjarmasin Utara, Kelurahan Pangeran, Kecamatan Banjarmasin Utara, Kota Banjarmasin, Kalimantan Selatan"
		// Koordinat KUA Banjarmasin Utara (contoh)
		lat := -3.3148
		lon := 114.5925
		latitude = &lat
		longitude = &lon
	} else {
		// Use provided address for outside KUA
		alamatAkad = dataFormPendaftaran.JadwalDanLokasi.AlamatNikah
		// OPTIMASI: Koordinat akan diisi nanti secara asynchronous
		// Ini untuk menghindari delay dari external API call
	}

	// Create marriage registration
	pendaftaranNikah := structs.PendaftaranNikah{
		Nomor_pendaftaran:   nomorPendaftaran,
		Pendaftar_id:        userID.(string),
		Calon_suami_id:      fmt.Sprintf("%d", calonSuami.ID),
		Calon_istri_id:      fmt.Sprintf("%d", calonIstri.ID),
		Tanggal_pendaftaran: createdAt,
		Tanggal_nikah:       tanggalNikah,
		Waktu_nikah:         dataFormPendaftaran.JadwalDanLokasi.WaktuNikah,
		Tempat_nikah:        dataFormPendaftaran.JadwalDanLokasi.LokasiNikah,
		Nomor_dispensasi:    dataFormPendaftaran.JadwalDanLokasi.NomorDispensasi,
		Alamat_akad:         alamatAkad,
		Latitude:            latitude,
		Longitude:           longitude,
		Status_pendaftaran:  structs.StatusPendaftaranMenungguVerifikasi,
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

	// Create marriage guardian
	waliNikah := structs.WaliNikah{
		Pendaftaran_id:    pendaftaranNikah.ID,
		NIK:               dataFormPendaftaran.WaliNikah.NikWali,
		Nama_lengkap:      dataFormPendaftaran.WaliNikah.NamaLengkapWali,
		Hubungan_wali:     dataFormPendaftaran.WaliNikah.HubunganWali,
		Alamat:            dataFormPendaftaran.WaliNikah.AlamatWali,
		No_hp:             dataFormPendaftaran.WaliNikah.NomorTeleponWali,
		Agama:             dataFormPendaftaran.WaliNikah.AgamaWali,
		Status_keberadaan: dataFormPendaftaran.WaliNikah.StatusWali,
		Status_kehadiran:  structs.WaliStatusKehadiranBelumDiketahui,
		Created_at:        createdAt,
		Updated_at:        createdAt,
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

	// OPTIMASI: Kirim response ke client terlebih dahulu untuk kecepatan
	response := gin.H{
		"success": true,
		"message": "Pendaftaran nikah berhasil dibuat",
		"data": gin.H{
			"nomor_pendaftaran": nomorPendaftaran,
			"pendaftaran_id":    pendaftaranNikah.ID,
			"calon_suami":       calonSuami,
			"calon_istri":       calonIstri,
			"pendaftaran":       pendaftaranNikah,
			"wali_nikah":        waliNikah,
		},
	}

	c.JSON(http.StatusCreated, response)

	// OPTIMASI: Proses asynchronous setelah response dikirim
	go func() {
		// 1. Update koordinat jika alamat di luar KUA (background geocoding)
		if dataFormPendaftaran.JadwalDanLokasi.LokasiNikah != "Di KUA" {
			lat, lon, err := cache.GetCoordinatesFromAddressCached(alamatAkad)
			if err == nil {
				// Update koordinat ke database
				h.DB.Model(&structs.PendaftaranNikah{}).
					Where("id = ?", pendaftaranNikah.ID).
					Updates(map[string]interface{}{
						"latitude":  lat,
						"longitude": lon,
					})
				fmt.Printf("[ASYNC] Koordinat berhasil diupdate: %.6f, %.6f untuk pendaftaran #%d\n", lat, lon, pendaftaranNikah.ID)
			} else {
				fmt.Printf("[ASYNC] Warning: Gagal mendapatkan koordinat untuk alamat '%s': %v\n", alamatAkad, err)
			}
		}

		// 2. Kirim notifikasi (background notification)
		notificationService := services.NewNotificationService(h.DB)
		if err := notificationService.SendPendaftaranNotification(pendaftaranNikah.ID, userID.(string)); err != nil {
			fmt.Printf("[ASYNC] Gagal mengirim notifikasi pendaftaran: %v\n", err)
		} else {
			fmt.Printf("[ASYNC] Notifikasi berhasil dikirim untuk pendaftaran #%d\n", pendaftaranNikah.ID)
		}
	}()
}

// MarkAsVisited marks that the couple has visited the office with documents
func (h *InDB) MarkAsVisited(c *gin.Context) {
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

	// Check if registration exists and belongs to this user
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

	// Check if registration is in correct status
	if pendaftaran.Status_pendaftaran != structs.StatusPendaftaranBerkasDiterima {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status tidak sesuai",
			"error":   "Pendaftaran harus dalam status 'Berkas Diterima' untuk menandai kunjungan",
			"type":    "validation",
		})
		return
	}

	// Update status to indicate they have visited with documents
	pendaftaran.Status_pendaftaran = structs.StatusPendaftaranMenungguPenugasan
	pendaftaran.Updated_at = time.Now()

	if err := h.DB.Save(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   "Gagal mengupdate status pendaftaran",
			"type":    "database",
		})
		return
	}

	// Send notification to all staff
	notificationService := services.NewNotificationService(h.DB)
	err := notificationService.SendNotificationToRole(
		structs.UserRoleStaff,
		"Catin Datang ke Kantor",
		"Catin dengan nomor pendaftaran "+pendaftaran.Nomor_pendaftaran+" telah datang ke kantor dengan membawa berkas. Silakan verifikasi berkas fisik.",
		structs.NotifikasiTipeInfo,
		"/staff/pendaftaran/"+registrationID,
	)
	if err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Warning: Gagal mengirim notifikasi ke staff: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status berhasil diupdate",
		"data": gin.H{
			"id":                 pendaftaran.ID,
			"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
			"status_pendaftaran": pendaftaran.Status_pendaftaran,
			"updated_at":         pendaftaran.Updated_at,
		},
	})
}

// UpdateWeddingAddress updates the wedding address for a marriage registration
func (h *InDB) UpdateWeddingAddress(c *gin.Context) {
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

	// Get registration ID from URL parameter
	registrationID := c.Param("id")
	if registrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID pendaftaran diperlukan",
			"error":   "ID pendaftaran tidak ditemukan",
			"type":    "validation",
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

// CheckUserRegistrationStatus checks if user already has marriage registration
func (h *InDB) CheckUserRegistrationStatus(c *gin.Context) {
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

	// User already has registration
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User sudah memiliki pendaftaran nikah",
		"data": gin.H{
			"has_registration": true,
			"can_register":     false,
			"registration": gin.H{
				"id":                 pendaftaran.ID,
				"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
				"status_pendaftaran": pendaftaran.Status_pendaftaran,
				"tanggal_nikah":      pendaftaran.Tanggal_nikah,
				"tempat_nikah":       pendaftaran.Tempat_nikah,
				"alamat_akad":        pendaftaran.Alamat_akad,
				"created_at":         pendaftaran.Created_at,
			},
		},
	})
}

// GetAllMarriageRegistrations retrieves all marriage registrations with filters and pagination for staff
func (h *InDB) GetAllMarriageRegistrations(c *gin.Context) {
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
			"status_bimbingan":    r.Status_bimbingan,
			"tanggal_pendaftaran": r.Tanggal_pendaftaran,
			"tanggal_nikah":       r.Tanggal_nikah,
			"waktu_nikah":         r.Waktu_nikah,
			"tempat_nikah":        r.Tempat_nikah,
			"alamat_akad":         r.Alamat_akad,
			"nomor_dispensasi":    r.Nomor_dispensasi,
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
