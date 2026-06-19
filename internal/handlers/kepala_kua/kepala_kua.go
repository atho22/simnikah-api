package kepala_kua

import (
	"net/http"

	structs "simnikah/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== LIST AVAILABLE PENGHULU ====================

// ListAvailableOfficers menampilkan daftar penghulu aktif yang tersedia.
// Dipakai Kepala KUA untuk melihat penghulu sebelum melakukan assignment.
func (h *InDB) ListAvailableOfficers(c *gin.Context) {
	var penghulus []structs.Penghulu
	if err := h.DB.Where("status = ?", structs.PenghuluStatusAktif).
		Order("rating DESC").
		Find(&penghulus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data penghulu",
			"error":   err.Error(),
		})
		return
	}

	type penghuluSummary struct {
		ID           uint    `json:"id"`
		NamaLengkap  string  `json:"nama_lengkap"`
		Rating       float64 `json:"rating"`
		JumlahNikah  int     `json:"jumlah_nikah"`
		Status       string  `json:"status"`
	}

	summaries := make([]penghuluSummary, 0, len(penghulus))
	for _, p := range penghulus {
		summaries = append(summaries, penghuluSummary{
			ID:          p.ID,
			NamaLengkap: p.Nama_lengkap,
			Rating:      p.Rating,
			JumlahNikah: p.Jumlah_nikah,
			Status:      p.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar penghulu aktif",
		"data": gin.H{
			"total":     len(summaries),
			"penghulus": summaries,
		},
	})
}

// ==================== PENGHULU SCHEDULE AVAILABILITY ====================

// GetPenghuluScheduleForAssignment menampilkan jadwal penghulu untuk tanggal tertentu.
// Dipakai Kepala KUA untuk menentukan penghulu mana yang tersedia pada slot waktu tertentu.
func (h *InDB) GetPenghuluScheduleForAssignment(c *gin.Context) {
	dateStr := c.Query("tanggal")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tanggal diperlukan",
			"error":   "Query parameter tanggal (format YYYY-MM-DD) wajib diisi",
		})
		return
	}

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

	// Single GROUP BY query untuk menghindari N+1
	type countRow struct {
		PenghuluID uint
		BookedSlots int64
	}
	var counts []countRow
	h.DB.Model(&structs.PendaftaranNikah{}).
		Select("penghulu_id as penghulu_id, COUNT(*) as booked_slots").
		Where("DATE(tanggal_nikah) = ? AND status_pendaftaran NOT IN ?",
			dateStr, []string{structs.StatusPendaftaranDitolak}).
		Group("penghulu_id").Find(&counts)

	countMap := map[uint]int64{}
	for _, c := range counts {
		countMap[c.PenghuluID] = c.BookedSlots
	}

	type scheduleInfo struct {
		PenghuluID  uint     `json:"penghulu_id"`
		NamaLengkap string  `json:"nama_lengkap"`
		Rating      float64 `json:"rating"`
		JumlahNikah int     `json:"jumlah_nikah"`
		BookedSlots int64   `json:"booked_slots"`
		Available   bool    `json:"available"`
	}

	capacityPerDay := 3 // default CapacityPerDay
	schedules := make([]scheduleInfo, 0, len(penghulus))
	for _, p := range penghulus {
		booked := countMap[p.ID]
		schedules = append(schedules, scheduleInfo{
			PenghuluID:  p.ID,
			NamaLengkap: p.Nama_lengkap,
			Rating:      p.Rating,
			JumlahNikah: p.Jumlah_nikah,
			BookedSlots: booked,
			Available:   booked < int64(capacityPerDay),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jadwal penghulu berhasil diambil",
		"data": gin.H{
			"tanggal":   dateStr,
			"penghulus": schedules,
		},
	})
}
