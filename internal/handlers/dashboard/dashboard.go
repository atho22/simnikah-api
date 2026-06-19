package dashboard

import (
	"fmt"
	"net/http"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// ==================== DASHBOARD KEPALA KUA ====================

// GetKepalaKUADashboard returns comprehensive dashboard data for Kepala KUA
func (h *InDB) GetKepalaKUADashboard(c *gin.Context) {
	now := time.Now().In(utils.WITA)
	period := c.DefaultQuery("period", "month") // day, week, month, year
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Calculate date range based on period
	var startDate, endDate time.Time
	if dateFrom != "" && dateTo != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_from tidak valid",
				"error":   "Format harus YYYY-MM-DD",
			})
			return
		}
		endDate, err = time.Parse("2006-01-02", dateTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_to tidak valid",
				"error":   "Format harus YYYY-MM-DD",
			})
			return
		}
		endDate = endDate.Add(24 * time.Hour) // Include entire day
	} else {
		switch period {
		case "day":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utils.WITA)
			endDate = startDate.Add(24 * time.Hour)
		case "week":
			startDate = time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, utils.WITA)
			endDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utils.WITA).Add(24 * time.Hour)
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
			endDate = startDate.AddDate(0, 1, 0)
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, utils.WITA)
			endDate = startDate.AddDate(1, 0, 0)
		default:
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
			endDate = startDate.AddDate(0, 1, 0)
		}
	}

	// 1. Statistik Pernikahan (Harian, Bulanan, Tahunan)
	stats := h.getMarriageStatistics(startDate, endDate, now)

	// 2. Grafik Tren Pernikahan
	trends := h.getMarriageTrends(startDate, endDate, period)

	// 3. Status Pendaftaran (Pie Chart)
	statusDistribution := h.getStatusDistribution(startDate, endDate)

	// 4. Penghulu Performance
	penghuluPerformance := h.getPenghuluPerformance(startDate, endDate)

	// 5. Peak Hours Analysis
	peakHours := h.getPeakHoursAnalysis(startDate, endDate)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard data berhasil diambil",
		"data": gin.H{
			"period": gin.H{
				"type":       period,
				"date_from":  startDate.Format("2006-01-02"),
				"date_to":    endDate.Add(-24 * time.Hour).Format("2006-01-02"),
			},
			"statistics":           stats,
			"trends":               trends,
			"status_distribution":  statusDistribution,
			"penghulu_performance": penghuluPerformance,
			"peak_hours":           peakHours,
		},
	})
}

// ==================== DASHBOARD STAFF ====================

// GetStaffDashboard returns dashboard data for Staff
func (h *InDB) GetStaffDashboard(c *gin.Context) {
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

	// 1. Pending assignments (Menunggu Penugasan)
	pendingAssignments := h.getPendingAssignments()

	// 2. Recent registrations
	recentRegistrations := h.getRecentRegistrations()

	// 3. Timeline per pendaftaran (recent activities)
	timeline := h.getRegistrationTimeline(userID.(string))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard staff berhasil diambil",
		"data": gin.H{
			"pending_assignments":    pendingAssignments,
			"recent_registrations":   recentRegistrations,
			"timeline":               timeline,
		},
	})
}

// ==================== HELPER FUNCTIONS ====================

// getMarriageStatistics returns marriage statistics
func (h *InDB) getMarriageStatistics(startDate, endDate, now time.Time) gin.H {
	var total, hariIni, bulanIni, tahunIni int64

	// Total pernikahan dalam periode
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", startDate, endDate).
		Count(&total)

	// Hari ini (WITA)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utils.WITA)
	todayEnd := todayStart.Add(24 * time.Hour)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", todayStart, todayEnd).
		Count(&hariIni)

	// Bulan ini (WITA)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
	monthEnd := monthStart.AddDate(0, 1, 0)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", monthStart, monthEnd).
		Count(&bulanIni)

	// Tahun ini (WITA)
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, utils.WITA)
	yearEnd := yearStart.AddDate(1, 0, 0)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", yearStart, yearEnd).
		Count(&tahunIni)

	// Pernikahan selesai dalam periode
	var selesai int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranSelesai).
		Count(&selesai)

	// Status breakdown sesuai flow scheduling-only: Menunggu Penugasan → Penghulu Ditugaskan → Selesai
	var menungguPenugasan, penghuluDitugaskan, ditolak int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("created_at >= ? AND created_at < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranMenungguPenugasan).
		Count(&menungguPenugasan)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("created_at >= ? AND created_at < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranPenghuluDitugaskan).
		Count(&penghuluDitugaskan)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("created_at >= ? AND created_at < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranDitolak).
		Count(&ditolak)

	return gin.H{
		"total_periode": total,
		"hari_ini":      hariIni,
		"bulan_ini":     bulanIni,
		"tahun_ini":     tahunIni,
		"selesai":       selesai,
		"pending":       total - selesai,
		"status_breakdown": gin.H{
			"menunggu_penugasan":  menungguPenugasan,
			"penghulu_ditugaskan": penghuluDitugaskan,
			"selesai":             selesai,
			"ditolak":             ditolak,
		},
	}
}

// getMarriageTrends menggunakan GROUP BY tunggal untuk efisiensi (mengganti N query berurutan)
func (h *InDB) getMarriageTrends(startDate, endDate time.Time, period string) []gin.H {
	var dateFormat string
	switch period {
	case "day":
		dateFormat = "%Y-%m-%d"
	case "month":
		dateFormat = "%Y-%m"
	case "year":
		dateFormat = "%Y"
	default:
		dateFormat = "%Y-%m"
	}

	var results []struct {
		DateStr string
		Count   int64
	}
	h.DB.Model(&structs.PendaftaranNikah{}).
		Select(fmt.Sprintf("DATE_FORMAT(tanggal_nikah, '%s') as date_str, COUNT(*) as count", dateFormat)).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", startDate, endDate).
		Group("date_str").
		Order("date_str ASC").
		Find(&results)

	trends := make([]gin.H, 0, len(results))
	for _, r := range results {
		trends = append(trends, gin.H{
			"date":  r.DateStr,
			"count": r.Count,
		})
	}
	return trends
}

// getStatusDistribution menggunakan GROUP BY tunggal untuk efisiensi
func (h *InDB) getStatusDistribution(startDate, endDate time.Time) []gin.H {
	var results []struct {
		Status string
		Count  int64
	}
	h.DB.Model(&structs.PendaftaranNikah{}).
		Select("status_pendaftaran as status, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Group("status_pendaftaran").
		Find(&results)

	countMap := map[string]int64{}
	for _, r := range results {
		countMap[r.Status] = r.Count
	}

	statuses := []string{
		structs.StatusPendaftaranMenungguPenugasan,
		structs.StatusPendaftaranPenghuluDitugaskan,
		structs.StatusPendaftaranSelesai,
		structs.StatusPendaftaranDitolak,
	}

	distribution := make([]gin.H, 0, len(statuses))
	for _, status := range statuses {
		distribution = append(distribution, gin.H{
			"status": status,
			"count":  countMap[status],
			"label":  h.getStatusLabel(status),
		})
	}
	return distribution
}

// getStatusLabel returns human-readable label for status
func (h *InDB) getStatusLabel(status string) string {
	labels := map[string]string{
		structs.StatusPendaftaranMenungguPenugasan:  "Menunggu Penugasan",
		structs.StatusPendaftaranPenghuluDitugaskan: "Penghulu Ditugaskan",
		structs.StatusPendaftaranSelesai:            "Selesai",
		structs.StatusPendaftaranDitolak:            "Ditolak",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

// getPenghuluPerformance menggunakan GROUP BY tunggal untuk efisiensi (mengganti N+1)
func (h *InDB) getPenghuluPerformance(startDate, endDate time.Time) []gin.H {
	var penghulus []structs.Penghulu
	h.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus)

	// Single GROUP BY query
	type countRow struct {
		PenghuluID uint
		Jumlah     int64
	}
	var counts []countRow
	h.DB.Model(&structs.PendaftaranNikah{}).
		Select("penghulu_id as penghulu_id, COUNT(*) as jumlah").
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
			startDate, endDate, structs.StatusPendaftaranSelesai).
		Group("penghulu_id").Find(&counts)

	countMap := map[uint]int64{}
	for _, c := range counts {
		countMap[c.PenghuluID] = c.Jumlah
	}

	performance := make([]gin.H, 0, len(penghulus))
	for _, penghulu := range penghulus {
		performance = append(performance, gin.H{
			"penghulu_id":   penghulu.ID,
			"nama_lengkap":  penghulu.Nama_lengkap,
			"jumlah_nikah":  countMap[penghulu.ID],
			"rating":        penghulu.Rating,
		})
	}
	return performance
}

// getPeakHoursAnalysis menggunakan GROUP BY tunggal untuk efisiensi
func (h *InDB) getPeakHoursAnalysis(startDate, endDate time.Time) []gin.H {
	var results []struct {
		Hour  string
		Count int64
	}
	h.DB.Model(&structs.PendaftaranNikah{}).
		Select("SUBSTRING(waktu_nikah, 1, 2) as hour, COUNT(*) as count").
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
			startDate, endDate, structs.StatusPendaftaranSelesai).
		Group("hour").
		Find(&results)

	countMap := map[string]int64{}
	for _, r := range results {
		countMap[r.Hour] = r.Count
	}

	timeSlots := []string{"08", "09", "10", "11", "12", "13", "14", "15", "16"}
	peakHours := make([]gin.H, 0, len(timeSlots))
	for _, slotHour := range timeSlots {
		peakHours = append(peakHours, gin.H{
			"waktu": slotHour + ":00",
			"count": countMap[slotHour],
		})
	}
	return peakHours
}

// getPendingAssignments returns registrations awaiting penghulu assignment
func (h *InDB) getPendingAssignments() []gin.H {
	var pendaftarans []structs.PendaftaranNikah
	h.DB.Where("status_pendaftaran = ?", structs.StatusPendaftaranMenungguPenugasan).
		Order("created_at DESC").
		Limit(10).
		Find(&pendaftarans)

	var pending []gin.H
	for _, p := range pendaftarans {
		pending = append(pending, gin.H{
			"id":                 p.ID,
			"nama_suami":         p.Nama_suami,
			"nama_istri":         p.Nama_istri,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"alamat_akad":        p.Alamat_akad,
			"created_at":         p.Created_at,
		})
	}

	return pending
}

// getRecentRegistrations returns recent registrations
func (h *InDB) getRecentRegistrations() []gin.H {
	var pendaftarans []structs.PendaftaranNikah
	h.DB.Order("created_at DESC").
		Limit(10).
		Find(&pendaftarans)

	var recent []gin.H
	for _, p := range pendaftarans {
		recent = append(recent, gin.H{
			"id":                 p.ID,
			"nama_suami":         p.Nama_suami,
			"nama_istri":         p.Nama_istri,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"alamat_akad":        p.Alamat_akad,
			"created_at":         p.Created_at,
		})
	}

	return recent
}

// getRegistrationTimeline returns recent registration activities
func (h *InDB) getRegistrationTimeline(staffID string) []gin.H {
	var pendaftarans []structs.PendaftaranNikah
	h.DB.Order("updated_at DESC").
		Limit(20).
		Find(&pendaftarans)

	var timeline []gin.H
	for _, p := range pendaftarans {
		timeline = append(timeline, gin.H{
			"id":                 p.ID,
			"nama_suami":         p.Nama_suami,
			"nama_istri":         p.Nama_istri,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        p.Waktu_nikah,
			"tempat_nikah":       p.Tempat_nikah,
			"updated_at":         p.Updated_at,
			"action":             fmt.Sprintf("Status diubah menjadi %s", p.Status_pendaftaran),
		})
	}

	return timeline
}

// ==================== DETAILED STATISTICS ENDPOINTS ====================

// GetMarriageStatistics returns detailed marriage statistics
func (h *InDB) GetMarriageStatistics(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	now := time.Now().In(utils.WITA)

	var startDate, endDate time.Time
	switch period {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utils.WITA)
		endDate = startDate.Add(24 * time.Hour)
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
		endDate = startDate.AddDate(0, 1, 0)
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, utils.WITA)
		endDate = startDate.AddDate(1, 0, 0)
	default:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
		endDate = startDate.AddDate(0, 1, 0)
	}

	stats := h.getMarriageStatistics(startDate, endDate, now)
	trends := h.getMarriageTrends(startDate, endDate, period)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik pernikahan berhasil diambil",
		"data": gin.H{
			"statistics": stats,
			"trends":     trends,
		},
	})
}

// GetPenghuluPerformance returns detailed penghulu performance
func (h *InDB) GetPenghuluPerformance(c *gin.Context) {
	// Get query parameters
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	var startDate, endDate time.Time
	now := time.Now().In(utils.WITA)

	if dateFrom != "" && dateTo != "" {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", dateFrom, utils.WITA)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_from tidak valid",
			})
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", dateTo, utils.WITA)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_to tidak valid",
			})
			return
		}
		endDate = endDate.Add(24 * time.Hour)
	} else {
		// Default: bulan ini
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
		endDate = startDate.AddDate(0, 1, 0)
	}

	performance := h.getPenghuluPerformance(startDate, endDate)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik penghulu berhasil diambil",
		"data":    performance,
	})
}

// GetPeakHoursAnalysis returns peak hours analysis
func (h *InDB) GetPeakHoursAnalysis(c *gin.Context) {
	// Get query parameters
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	var startDate, endDate time.Time
	now := time.Now().In(utils.WITA)

	if dateFrom != "" && dateTo != "" {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", dateFrom, utils.WITA)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_from tidak valid",
			})
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", dateTo, utils.WITA)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_to tidak valid",
			})
			return
		}
		endDate = endDate.Add(24 * time.Hour)
	} else {
		// Default: bulan ini
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, utils.WITA)
		endDate = startDate.AddDate(0, 1, 0)
	}

	peakHours := h.getPeakHoursAnalysis(startDate, endDate)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analisis jam sibuk berhasil diambil",
		"data":    peakHours,
	})
}
