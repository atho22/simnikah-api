package catin

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "simnikah/internal/models"
	services "simnikah/internal/services"

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

var HariLiburNasional = map[string]string{
	// Hari libur tetap (setiap tahun)
	"01-01": "Tahun Baru Masehi",
	"05-01": "Hari Buruh Internasional",
	"06-01": "Hari Lahir Pancasila",
	"08-17": "Hari Kemerdekaan RI",
	"12-25": "Hari Natal",

	// Hari libur 2024 (tanggal berubah)
	"2024-02-08": "Isra Mi'raj",
	"2024-03-11": "Hari Raya Nyepi",
	"2024-03-29": "Wafat Isa Al-Masih",
	"2024-04-10": "Idul Fitri 1445 H",
	"2024-04-11": "Idul Fitri 1445 H",
	"2024-05-09": "Kenaikan Isa Al-Masih",
	"2024-05-23": "Hari Raya Waisak",
	"2024-06-17": "Idul Adha 1445 H",
	"2024-07-07": "Tahun Baru Islam 1446 H",
	"2024-09-16": "Maulid Nabi Muhammad SAW",

	// Hari libur 2025 (tanggal berubah)
	"2025-01-27": "Isra Mi'raj",
	"2025-01-29": "Tahun Baru Imlek",
	"2025-03-29": "Hari Raya Nyepi",
	"2025-03-31": "Idul Fitri 1446 H",
	"2025-04-01": "Idul Fitri 1446 H",
	"2025-04-18": "Wafat Isa Al-Masih",
	"2025-05-12": "Hari Raya Waisak",
	"2025-05-29": "Kenaikan Isa Al-Masih",
	"2025-06-06": "Idul Adha 1446 H",
	"2025-06-27": "Tahun Baru Islam 1447 H",
	"2025-09-05": "Maulid Nabi Muhammad SAW",
}

// IsHariLibur mengecek apakah tanggal adalah hari libur (Minggu atau hari libur nasional)
func IsHariLibur(tanggal time.Time) (bool, string) {
	// Cek hari Minggu
	if tanggal.Weekday() == time.Sunday {
		return true, "Hari Minggu"
	}

	// Cek hari libur nasional dengan format YYYY-MM-DD
	tanggalFull := tanggal.Format("2006-01-02")
	if nama, exists := HariLiburNasional[tanggalFull]; exists {
		return true, nama
	}

	// Cek hari libur nasional dengan format MM-DD (hari libur tetap)
	tanggalShort := tanggal.Format("01-02")
	if nama, exists := HariLiburNasional[tanggalShort]; exists {
		return true, nama
	}

	return false, ""
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
			"success":      false,
			"message":      "Format tanggal tidak benar",
			"error":        "Format tanggal harus: Tahun-Bulan-Tanggal (contoh: 2024-12-25 untuk tanggal 25 Desember 2024)",
			"field":        "tanggal_nikah",
			"type":         "format",
			"contoh_benar": []string{"2024-12-25", "2024-01-05", "2024-11-20"},
			"contoh_salah": []string{"25-12-2024", "25/12/2024", "25 Desember 2024"},
		})
		return
	}

	// Validate that wedding date is not in the past
	if tanggalNikah.Before(time.Now().Truncate(24 * time.Hour)) {
		today := time.Now().Format("02 Januari 2006")
		c.JSON(http.StatusBadRequest, gin.H{
			"success":          false,
			"message":          "Tanggal tidak boleh sudah lewat",
			"error":            fmt.Sprintf("Tanggal nikah tidak boleh di masa lalu. Hari ini adalah %s. Silakan pilih tanggal hari ini atau tanggal yang akan datang.", today),
			"field":            "tanggal_nikah",
			"type":             "validation",
			"tanggal_hari_ini": today,
		})
		return
	}

	// Validate wedding time format
	_, err = time.Parse("15:04", formSederhana.LokasiNikah.WaktuNikah)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"message":      "Format jam tidak benar",
			"error":        "Format jam harus: Jam:Menit dengan 2 angka (contoh: 09:00 untuk jam 9 pagi, 14:30 untuk jam 2 siang)",
			"field":        "waktu_nikah",
			"type":         "format",
			"contoh_benar": []string{"08:00", "09:00", "14:30", "16:00"},
			"contoh_salah": []string{"9:00", "9:0", "2:30 PM", "14.30"},
			"jam_tersedia": TimeSlots,
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

	// Validate wali nikah
	if strings.TrimSpace(formSederhana.WaliNikah.NamaDanBin) == "" {
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success":        false,
			"message":        "Validasi gagal",
			"error":          "Hubungan wali tidak valid",
			"field":          "wali_nikah.hubungan_wali",
			"type":           "enum",
			"hubungan_valid": structs.ValidHubunganWali,
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

	// Cek apakah hari libur
	isLibur, namaLibur := IsHariLibur(tanggalNikah)

	// Validate schedule based on location
	const maxTotalWeddings = 3 // Total maksimal pernikahan per jam = 3 penghulu

	if formSederhana.LokasiNikah.TempatNikah == "Di KUA" {
		// Jika hari libur, tidak bisa nikah di KUA
		if isLibur {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "KUA tutup pada hari libur",
				"error": fmt.Sprintf("Maaf, KUA tutup pada tanggal %s (%s). Pernikahan di KUA tidak tersedia pada hari libur. Silakan pilih tanggal lain atau pilih nikah di luar KUA.",
					tanggalNikah.Format("02 Januari 2006"), namaLibur),
				"field": "tanggal_nikah",
				"type":  "holiday_restriction",
				"saran": "Anda masih bisa memilih nikah di luar KUA pada hari libur, atau pilih tanggal lain untuk nikah di KUA.",
				"data": gin.H{
					"tanggal_nikah":     tanggalNikah.Format("2006-01-02"),
					"hari":              tanggalNikah.Weekday().String(),
					"is_hari_libur":     true,
					"nama_hari_libur":   namaLibur,
					"tersedia_luar_kua": true,
				},
			})
			return
		}
		// Validasi 1: Total sudah mencapai maksimal 3 (semua penghulu terpakai)
		// Jika sudah ada 3 nikah di luar KUA, maka KUA juga tidak bisa
		if countTotalRegistrations >= maxTotalWeddings {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal sudah penuh",
				"error": fmt.Sprintf("Maaf, jadwal pernikahan pada tanggal %s pukul %s sudah penuh. Sudah ada %d pernikahan (semua penghulu sudah bertugas). Maksimal hanya 3 pernikahan per jam. Silakan pilih tanggal atau jam lain yang masih tersedia.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized, countTotalRegistrations),
				"field": "waktu_nikah",
				"type":  "schedule_conflict",
				"saran": "Coba pilih jam lain pada tanggal yang sama, atau pilih tanggal lain. Gunakan fitur kalender untuk melihat jadwal yang tersedia.",
				"data": gin.H{
					"tanggal_nikah":        tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":          waktuNikahNormalized,
					"tempat_nikah":         formSederhana.LokasiNikah.TempatNikah,
					"total_pernikahan":     countTotalRegistrations,
					"batas_maksimal_total": maxTotalWeddings,
				},
			})
			return
		}

		// Validasi 2: Jika nikah di KUA, tidak boleh ada lebih dari 1 pernikahan di KUA di jam yang sama
		if countKUA >= 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal di KUA sudah terisi",
				"error": fmt.Sprintf("Maaf, jadwal pernikahan di KUA pada tanggal %s pukul %s sudah terisi oleh pendaftar lain. Di KUA hanya bisa 1 pernikahan per jam. Silakan pilih tanggal atau jam lain, atau pilih nikah di luar KUA.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field": "waktu_nikah",
				"type":  "schedule_conflict",
				"saran": "Anda masih bisa memilih nikah di luar KUA pada jam yang sama, atau pilih jam/tanggal lain untuk nikah di KUA.",
				"data": gin.H{
					"tanggal_nikah":          tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":            waktuNikahNormalized,
					"tempat_nikah":           formSederhana.LokasiNikah.TempatNikah,
					"jumlah_terisi_kua":      countKUA,
					"batas_maksimal_kua":     1,
					"slot_luar_kua_tersedia": maxTotalWeddings - countTotalRegistrations,
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
				errorMsg = fmt.Sprintf("Maaf, jadwal pernikahan pada tanggal %s pukul %s sudah penuh. Sudah ada 3 pernikahan (1 di KUA dan 2 di luar KUA). Maksimal hanya 3 pernikahan per jam. Silakan pilih tanggal atau jam lain yang masih tersedia.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized)
			} else {
				errorMsg = fmt.Sprintf("Maaf, jadwal pernikahan pada tanggal %s pukul %s sudah penuh. Sudah ada 3 pernikahan di luar KUA. Maksimal hanya 3 pernikahan per jam. Silakan pilih tanggal atau jam lain yang masih tersedia.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized)
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal sudah penuh",
				"error":   errorMsg,
				"field":   "waktu_nikah",
				"type":    "schedule_conflict",
				"saran":   "Coba pilih jam lain pada tanggal yang sama, atau pilih tanggal lain. Gunakan fitur kalender untuk melihat jadwal yang tersedia.",
				"data": gin.H{
					"tanggal_nikah":          tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":            waktuNikahNormalized,
					"tempat_nikah":           formSederhana.LokasiNikah.TempatNikah,
					"total_pernikahan":       countTotalRegistrations,
					"pernikahan_di_kua":      countKUA,
					"pernikahan_di_luar_kua": countLuarKUA,
					"batas_maksimal_total":   maxTotalWeddings,
					"slot_tersedia_luar_kua": 0,
				},
			})
			return
		}

		// Validasi 2: Jika sudah ada 1 di KUA, maka slot luar KUA maksimal 2
		if countKUA >= 1 && countLuarKUA >= 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Jadwal di luar KUA sudah penuh",
				"error": fmt.Sprintf("Maaf, jadwal pernikahan di luar KUA pada tanggal %s pukul %s sudah penuh. Sudah ada 1 pernikahan di KUA dan 2 di luar KUA (total maksimal 3 pernikahan per jam). Silakan pilih tanggal atau jam lain yang masih tersedia.",
					tanggalNikah.Format("02 Januari 2006"), waktuNikahNormalized),
				"field": "waktu_nikah",
				"type":  "schedule_conflict",
				"saran": "Coba pilih jam lain pada tanggal yang sama, atau pilih tanggal lain. Gunakan fitur kalender untuk melihat jadwal yang tersedia.",
				"data": gin.H{
					"tanggal_nikah":          tanggalNikah.Format("2006-01-02"),
					"waktu_nikah":            waktuNikahNormalized,
					"tempat_nikah":           formSederhana.LokasiNikah.TempatNikah,
					"total_pernikahan":       countTotalRegistrations,
					"pernikahan_di_kua":      countKUA,
					"pernikahan_di_luar_kua": countLuarKUA,
					"batas_maksimal_total":   maxTotalWeddings,
					"slot_tersedia_luar_kua": slotTersediaLuarKUA,
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
	hashBride := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_bride_%d", userIDStr, timestamp+1))))
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

	// Run forward-chaining recommendation asynchronously
	go func(regID uint, pendaftarID string) {
		fc := services.NewForwardChainingEngine(h.DB)
		rec, err := fc.GetPenghuluRecommendations(regID)
		if err != nil {
			fmt.Println("forward chaining error:", err)
			return
		}
		if rec != nil && rec.RecommendedPenghuluID != 0 {
			// Notify the registrant that a penghulu recommendation is available
			notif := structs.Notifikasi{
				User_id:    pendaftarID,
				Judul:      "Rekomendasi Penghulu Tersedia",
				Pesan:      fmt.Sprintf("Sistem merekomendasikan penghulu %s (score: %.2f).", rec.RecommendedPenghulu.Nama_lengkap, rec.SelectedScore),
				Tipe:       structs.NotifikasiTipeInfo,
				Status_baca: structs.NotifikasiStatusBelumDibaca,
				Link:       fmt.Sprintf("/pendaftaran/%d/rekomendasi", regID),
				Created_at: time.Now(),
				Updated_at: time.Now(),
			}
			if err := h.DB.Create(&notif).Error; err != nil {
				fmt.Println("failed to create recommendation notification:", err)
			}
		}
	}(pendaftaranNikah.ID, pendaftaranNikah.Pendaftar_id)

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
			"wali_nikah": gin.H{
				"nama_dan_bin":  waliNikah.Nama_dan_bin,
				"hubungan_wali": waliNikah.Hubungan_wali,
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

	// Fetch Calon Suami
	var calonSuami structs.CalonPasangan
	if err := h.DB.Where("id = ?", pendaftaran.Calon_suami_id).First(&calonSuami).Error; err == nil {
		registrationData["calon_suami"] = gin.H{
			"nama_lengkap":        calonSuami.Nama_lengkap,
			"pendidikan_terakhir": calonSuami.Pendidikan_terakhir,
			"tanggal_lahir":       calonSuami.Tanggal_lahir,
		}
	}

	// Fetch Calon Istri
	var calonIstri structs.CalonPasangan
	if err := h.DB.Where("id = ?", pendaftaran.Calon_istri_id).First(&calonIstri).Error; err == nil {
		registrationData["calon_istri"] = gin.H{
			"nama_lengkap":        calonIstri.Nama_lengkap,
			"pendidikan_terakhir": calonIstri.Pendidikan_terakhir,
			"tanggal_lahir":       calonIstri.Tanggal_lahir,
		}
	}

	// Fetch Wali Nikah
	if pendaftaran.Wali_nikah_id != nil {
		var waliNikah structs.WaliNikah
		if err := h.DB.Where("id = ?", *pendaftaran.Wali_nikah_id).First(&waliNikah).Error; err == nil {
			registrationData["wali_nikah"] = gin.H{
				"nama_dan_bin":  waliNikah.Nama_dan_bin,
				"hubungan_wali": waliNikah.Hubungan_wali,
			}
		}
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

// GetRegistrationDetail gets detailed information of a specific registration by ID
// GET /simnikah/pendaftaran/:id
func (h *InDB) GetRegistrationDetail(c *gin.Context) {
	registrationID := c.Param("id")

	// Get user_id and role from context
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

	userRole, roleExists := c.Get("role")
	if !roleExists {
		// Try to get from user data
		var user structs.Users
		if err := h.DB.Where("user_id = ?", userID.(string)).First(&user).Error; err == nil {
			userRole = user.Role
		} else {
			userRole = ""
		}
	}

	// Get registration
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ?", registrationID).First(&pendaftaran).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pendaftaran tidak ditemukan",
				"error":   "Pendaftaran dengan ID tersebut tidak ditemukan",
				"type":    "not_found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database error",
				"error":   "Gagal mengambil data pendaftaran",
				"type":    "database",
			})
		}
		return
	}

	// Authorization check: user_biasa hanya bisa lihat pendaftaran miliknya sendiri
	if userRole == "user_biasa" {
		if pendaftaran.Pendaftar_id != userID.(string) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak",
				"error":   "Anda tidak memiliki akses untuk melihat pendaftaran ini",
				"type":    "authorization",
			})
			return
		}
	}
	// Staff, penghulu, dan kepala_kua bisa melihat semua pendaftaran

	// Build detailed registration response
	registrationData := gin.H{
		"id":                 pendaftaran.ID,
		"nomor_pendaftaran":  pendaftaran.Nomor_pendaftaran,
		"pendaftar_id":       pendaftaran.Pendaftar_id,
		"status_pendaftaran": pendaftaran.Status_pendaftaran,
		"tanggal_pendaftaran": pendaftaran.Tanggal_pendaftaran,
		"tanggal_nikah":      pendaftaran.Tanggal_nikah,
		"waktu_nikah":        pendaftaran.Waktu_nikah,
		"tempat_nikah":       pendaftaran.Tempat_nikah,
		"alamat_akad":        pendaftaran.Alamat_akad,
		"latitude":           pendaftaran.Latitude,
		"longitude":          pendaftaran.Longitude,
		"catatan":            pendaftaran.Catatan,
		"disetujui_oleh":     pendaftaran.Disetujui_oleh,
		"disetujui_pada":     pendaftaran.Disetujui_pada,
		"created_at":         pendaftaran.Created_at,
		"updated_at":         pendaftaran.Updated_at,
	}

	// Fetch Calon Suami (full details)
	var calonSuami structs.CalonPasangan
	if err := h.DB.Where("id = ?", pendaftaran.Calon_suami_id).First(&calonSuami).Error; err == nil {
		registrationData["calon_suami"] = gin.H{
			"id":                 calonSuami.ID,
			"user_id":            calonSuami.User_id,
			"nik":                calonSuami.NIK,
			"nama_lengkap":       calonSuami.Nama_lengkap,
			"tanggal_lahir":      calonSuami.Tanggal_lahir,
			"jenis_kelamin":      calonSuami.Jenis_kelamin,
			"pendidikan_terakhir": calonSuami.Pendidikan_terakhir,
			"created_at":         calonSuami.Created_at,
			"updated_at":         calonSuami.Updated_at,
		}
	}

	// Fetch Calon Istri (full details)
	var calonIstri structs.CalonPasangan
	if err := h.DB.Where("id = ?", pendaftaran.Calon_istri_id).First(&calonIstri).Error; err == nil {
		registrationData["calon_istri"] = gin.H{
			"id":                 calonIstri.ID,
			"user_id":            calonIstri.User_id,
			"nik":                calonIstri.NIK,
			"nama_lengkap":       calonIstri.Nama_lengkap,
			"tanggal_lahir":      calonIstri.Tanggal_lahir,
			"jenis_kelamin":      calonIstri.Jenis_kelamin,
			"pendidikan_terakhir": calonIstri.Pendidikan_terakhir,
			"created_at":         calonIstri.Created_at,
			"updated_at":         calonIstri.Updated_at,
		}
	}

	// Fetch Wali Nikah
	if pendaftaran.Wali_nikah_id != nil {
		var waliNikah structs.WaliNikah
		if err := h.DB.Where("id = ?", *pendaftaran.Wali_nikah_id).First(&waliNikah).Error; err == nil {
			registrationData["wali_nikah"] = gin.H{
				"id":             waliNikah.ID,
				"nama_dan_bin":   waliNikah.Nama_dan_bin,
				"hubungan_wali":  waliNikah.Hubungan_wali,
				"created_at":     waliNikah.Created_at,
				"updated_at":     waliNikah.Updated_at,
			}
		}
	}

	// Fetch Penghulu info if assigned
	if pendaftaran.Penghulu_id != nil {
		var penghulu structs.Penghulu
		if err := h.DB.Where("id = ?", *pendaftaran.Penghulu_id).First(&penghulu).Error; err == nil {
			registrationData["penghulu"] = gin.H{
				"id":                penghulu.ID,
				"user_id":           penghulu.User_id,
				"nip":               penghulu.NIP,
				"nama_lengkap":      penghulu.Nama_lengkap,
				"no_hp":             penghulu.No_hp,
				"email":             penghulu.Email,
				"alamat":            penghulu.Alamat,
				"status":            penghulu.Status,
				"ditugaskan_oleh":    pendaftaran.Penghulu_assigned_by,
				"ditugaskan_pada":    pendaftaran.Penghulu_assigned_at,
				"created_at":         penghulu.Created_at,
				"updated_at":         penghulu.Updated_at,
			}
		}
	}

	// Add location URLs if coordinates exist
	if pendaftaran.Latitude != nil && pendaftaran.Longitude != nil {
		lat := *pendaftaran.Latitude
		lon := *pendaftaran.Longitude
		registrationData["location"] = gin.H{
			"latitude":                  lat,
			"longitude":                 lon,
			"has_coordinates":           true,
			"google_maps_url":           fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", lat, lon),
			"google_maps_directions_url": fmt.Sprintf("https://www.google.com/maps/dir/?api=1&destination=%f,%f", lat, lon),
			"waze_url":                  fmt.Sprintf("https://www.waze.com/ul?ll=%f,%f&navigate=yes", lat, lon),
			"osm_url":                   fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", lat, lon),
		}
	} else {
		registrationData["location"] = gin.H{
			"has_coordinates": false,
		}
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail pendaftaran berhasil diambil",
		"data":    registrationData,
	})
}

// ==================== CALENDAR AVAILABILITY ====================

type checkScheduleAvailabilityRequest struct {
	TanggalNikah string   `json:"tanggal_nikah" binding:"required"`
	WaktuNikah   string   `json:"waktu_nikah"`
	TempatNikah  string   `json:"tempat_nikah"`
	AlamatNikah  string   `json:"alamat_nikah"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

// CheckScheduleAvailability mengecek apakah slot nikah masih tersedia.
// Endpoint ini dipakai catin sebelum mengirim pendaftaran final.
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

	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(input.TanggalNikah))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format tanggal tidak valid",
			"error":   "tanggal_nikah harus berformat YYYY-MM-DD",
		})
		return
	}

	engine := services.NewForwardChainingEngine(h.DB)

	var activePenghuluCount int64
	if err := h.DB.Model(&structs.Penghulu{}).
		Where("status = ?", structs.PenghuluStatusAktif).
		Count(&activePenghuluCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung penghulu aktif",
			"error":   err.Error(),
		})
		return
	}

	if activePenghuluCount == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Belum ada penghulu aktif",
			"data": gin.H{
				"is_available": false,
				"total_slots_available": 0,
				"message": "Tidak ada penghulu aktif yang dapat ditugaskan",
			},
		})
		return
	}

	selectedTime := strings.TrimSpace(input.WaktuNikah)
	selectedPlace := strings.TrimSpace(input.TempatNikah)
	if selectedPlace == "" {
		selectedPlace = structs.TempatNikahDiKUA
	}

	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var occupiedDaily int64
	if err := h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran IN ?",
			startOfDay,
			endOfDay,
			[]string{structs.StatusPendaftaranDisetujui, structs.StatusPendaftaranMenungguPenugasan, structs.StatusPendaftaranPenghuluDitugaskan, structs.StatusPendaftaranSelesai}).
		Count(&occupiedDaily).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengecek kapasitas jadwal",
			"error":   err.Error(),
		})
		return
	}

	dailyCapacity := activePenghuluCount * int64(engine.Config.CapacityPerDay)
	totalSlotsAvailable := dailyCapacity - occupiedDaily
	if totalSlotsAvailable < 0 {
		totalSlotsAvailable = 0
	}

	isAvailable := occupiedDaily < dailyCapacity
	message := "Tanggal masih tersedia"
	if !isAvailable {
		message = "Tanggal sudah penuh"
	}

	response := gin.H{
		"tanggal_nikah":           parsedDate.Format("2006-01-02"),
		"waktu_nikah":             selectedTime,
		"tempat_nikah":            selectedPlace,
		"occupied_slots":          occupiedDaily,
		"daily_capacity":          dailyCapacity,
		"total_slots_available":    totalSlotsAvailable,
		"is_available":            isAvailable,
		"message":                 message,
	}

	if selectedTime != "" {
		var occupiedPerSlot int64
		if err := h.DB.Model(&structs.PendaftaranNikah{}).
			Where("tanggal_nikah = ? AND waktu_nikah = ? AND status_pendaftaran IN ?",
				parsedDate,
				selectedTime,
				[]string{structs.StatusPendaftaranDisetujui, structs.StatusPendaftaranMenungguPenugasan, structs.StatusPendaftaranPenghuluDitugaskan, structs.StatusPendaftaranSelesai}).
			Count(&occupiedPerSlot).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal mengecek kapasitas per jam",
				"error":   err.Error(),
			})
			return
		}

		hourCapacity := activePenghuluCount * int64(engine.Config.CapacityPerHour)
		hourSlotsAvailable := hourCapacity - occupiedPerSlot
		if hourSlotsAvailable < 0 {
			hourSlotsAvailable = 0
		}

		response["occupied_hour_slots"] = occupiedPerSlot
		response["hour_capacity"] = hourCapacity
		response["hour_slots_available"] = hourSlotsAvailable

		if occupiedPerSlot >= hourCapacity {
			response["is_available"] = false
			response["message"] = "Slot waktu tersebut sudah penuh"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengecekan jadwal selesai",
		"data":    response,
	})
}

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
	registrationsPerDate := make(map[string]int) // Total (Draft + Disetujui)
	draftPerDate := make(map[string]int)         // Draft only (kuning - belum pasti)
	disetujuiPerDate := make(map[string]int)     // Disetujui only (hijau - sudah pasti)

	// Count per date and time slot, separate by location
	// Format: tanggalStr -> waktu -> {kua: count, luar_kua: count}
	type TimeSlotCount struct {
		KUA              int // Draft + Disetujui di KUA
		LuarKUA          int // Draft + Disetujui di luar KUA
		DraftKUA         int // Draft di KUA
		DisetujuiKUA     int // Disetujui di KUA
		DraftLuarKUA     int // Draft di luar KUA
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
		jumlahDraft := draftPerDate[tanggalStr]         // Kuning - belum pasti
		jumlahDisetujui := disetujuiPerDate[tanggalStr] // Hijau - sudah pasti

		// Sisa kuota dihitung berdasarkan total (Draft + Disetujui)
		// Draft juga dihitung dalam kuota meskipun belum pasti
		totalNikah := jumlahDraft + jumlahDisetujui
		sisaKuota := MaxWeddingsPerDay - totalNikah

		// Cek apakah hari libur
		isLibur, namaLibur := IsHariLibur(tanggalTime)

		// Determine status
		var status string
		var tersedia bool
		var tersediaKUAHariIni bool = true
		var tersediaLuarKUAHariIni bool = true

		if tanggalTime.Before(now.Truncate(24 * time.Hour)) {
			status = "Terlewat"
			tersedia = false
			tersediaKUAHariIni = false
			tersediaLuarKUAHariIni = false
		} else if isLibur {
			// Hari libur: KUA tutup, tapi luar KUA tetap tersedia
			status = "Hari Libur (Luar KUA Tersedia)"
			tersedia = true               // Masih tersedia untuk luar KUA
			tersediaKUAHariIni = false    // KUA tutup
			tersediaLuarKUAHariIni = true // Luar KUA tetap buka
		} else if totalNikah >= MaxWeddingsPerDay {
			status = "Penuh"
			tersedia = false
			tersediaKUAHariIni = false
			tersediaLuarKUAHariIni = false
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

			// Total semua (Draft + Disetujui) di KUA dan luar KUA
			totalAll := slotData.KUA + slotData.LuarKUA
			totalKUA := slotData.KUA

			// KUA tidak tersedia jika:
			// 1. Sudah ada 1 nikah di KUA (slot KUA penuh)
			// 2. Total sudah 3 (semua penghulu terpakai, termasuk jika 3 di luar KUA)
			// 3. Tanggal sudah lewat
			// 4. Hari libur
			tersediaKUA := totalKUA < maxKUAPerHour &&
				totalAll < maxTotalPerHour &&
				!tanggalTime.Before(now.Truncate(24*time.Hour)) &&
				tersediaKUAHariIni

			// Logika untuk luar KUA
			// Luar KUA tidak tersedia jika total sudah 3 (semua penghulu terpakai)
			tersediaLuarKUA := totalAll < maxTotalPerHour &&
				!tanggalTime.Before(now.Truncate(24*time.Hour)) &&
				tersediaLuarKUAHariIni

			// Hitung slot tersisa
			slotKUATersisa := 0
			if tersediaKUA {
				slotKUATersisa = 1
			}
			slotLuarKUATersisa := maxTotalPerHour - totalAll
			if slotLuarKUATersisa < 0 {
				slotLuarKUATersisa = 0
			}

			timeSlotsData = append(timeSlotsData, map[string]interface{}{
				"waktu": slot,
				"kua": gin.H{
					"tersedia":         tersediaKUA,
					"terbooking":       totalKUA >= maxKUAPerHour || totalAll >= maxTotalPerHour,
					"jumlah_total":     slotData.KUA,
					"jumlah_draft":     slotData.DraftKUA,
					"jumlah_disetujui": slotData.DisetujuiKUA,
					"slot_tersisa":     slotKUATersisa,
				},
				"luar_kua": gin.H{
					"tersedia":         tersediaLuarKUA,
					"terbooking":       totalAll >= maxTotalPerHour,
					"jumlah_total":     slotData.LuarKUA,
					"jumlah_draft":     slotData.DraftLuarKUA,
					"jumlah_disetujui": slotData.DisetujuiLuarKUA,
					"slot_tersisa":     slotLuarKUATersisa,
				},
				"total_pernikahan": totalAll,
				"slot_tersisa":     maxTotalPerHour - totalAll,
			})
		}

		calendar = append(calendar, map[string]interface{}{
			"tanggal":           tanggal,
			"tanggal_str":       tanggalStr,
			"hari":              tanggalTime.Weekday().String(),
			"status":            status,
			"tersedia":          tersedia,
			"tersedia_kua":      tersediaKUAHariIni,
			"tersedia_luar_kua": tersediaLuarKUAHariIni,
			"is_hari_libur":     isLibur,
			"nama_hari_libur":   namaLibur,
			"jumlah_nikah":      jumlahNikah,     // Total (Draft + Disetujui) di KUA
			"jumlah_draft":      jumlahDraft,     // Kuning - belum pasti (di KUA)
			"jumlah_disetujui":  jumlahDisetujui, // Hijau - sudah pasti (di KUA)
			"sisa_kuota":        sisaKuota,       // Berdasarkan yang sudah pasti (di KUA)
			"kapasitas":         MaxWeddingsPerDay,
			"is_today":          tanggalStr == now.Format("2006-01-02"),
			"is_past":           tanggalTime.Before(now.Truncate(24 * time.Hour)),
			"time_slots":        timeSlotsData, // Detail jam-jam tersedia
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

	// Cek apakah hari libur
	isLibur, namaLibur := IsHariLibur(tanggal)

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
	registrationsPerTimeSlot := make(map[string]int) // Total (Draft + Disetujui)
	draftPerTimeSlot := make(map[string]int)         // Draft only (kuning - belum pasti)
	disetujuiPerTimeSlot := make(map[string]int)     // Disetujui only (hijau - sudah pasti)

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
	totalBooked := 0 // Hanya yang sudah pasti (Disetujui)
	totalAvailable := 0

	// Get today's date for comparison (truncated to start of day)
	today := now.Truncate(24 * time.Hour)
	tanggalStartOfDay := tanggal.Truncate(24 * time.Hour)

	for _, slot := range TimeSlots {
		jumlahTotal := registrationsPerTimeSlot[slot]
		jumlahDraft := draftPerTimeSlot[slot]         // Kuning - belum pasti
		jumlahDisetujui := disetujuiPerTimeSlot[slot] // Hijau - sudah pasti

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
			"jumlah_nikah":     jumlahTotal,     // Total (Draft + Disetujui)
			"jumlah_draft":     jumlahDraft,     // Kuning - belum pasti
			"jumlah_disetujui": jumlahDisetujui, // Hijau - sudah pasti
		})
	}

	// Get summary - sisa kuota berdasarkan total (Draft + Disetujui)
	totalCapacity := len(TimeSlots)
	sisaKuota := totalCapacity - totalBooked

	var status string
	var tersediaKUA bool = true
	var tersediaLuarKUA bool = true

	if isLibur {
		// Hari libur: KUA tutup, luar KUA tetap tersedia
		status = "Hari Libur (Luar KUA Tersedia)"
		tersediaKUA = false
		tersediaLuarKUA = true
	} else if totalBooked >= totalCapacity {
		status = "Penuh"
		tersediaKUA = false
		tersediaLuarKUA = false
	} else if totalBooked > 0 {
		status = "Sebagian Tersedia"
	} else {
		status = "Semua Tersedia"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ketersediaan jam berhasil diambil",
		"data": gin.H{
			"tanggal":           tanggalStr,
			"hari":              tanggal.Weekday().String(),
			"status":            status,
			"is_hari_libur":     isLibur,
			"nama_hari_libur":   namaLibur,
			"tersedia_kua":      tersediaKUA,
			"tersedia_luar_kua": tersediaLuarKUA,
			"summary": gin.H{
				"total_slot": totalCapacity,
				"terbooking": totalBooked,
				"tersedia":   totalAvailable,
				"sisa_kuota": sisaKuota,
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
		CalonSuamiNama string `gorm:"column:calon_suami_nama"`
		CalonIstriNama string `gorm:"column:calon_istri_nama"`
		PenghuluNama   string `gorm:"column:penghulu_nama"`
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
				"total_nikah":   len(weddings),
				"nikah_di_kua":  nikahDiKUA,
				"nikah_di_luar": nikahDiLuar,
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
			"id":             feedback.ID,
			"pendaftaran_id": feedback.Pendaftaran_id,
			"jenis_feedback": feedback.Jenis_feedback,
			"rating":         feedback.Rating,
			"judul":          feedback.Judul,
			"status_baca":    feedback.Status_baca,
			"created_at":     feedback.Created_at,
		},
	})
}
