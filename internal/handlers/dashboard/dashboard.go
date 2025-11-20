package dashboard

import (
	"fmt"
	"net/http"
	"time"

	structs "simnikah/internal/models"

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
	now := time.Now()
	
	// Get query parameters for date range
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
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			endDate = startDate.Add(24 * time.Hour)
		case "week":
			startDate = now.AddDate(0, 0, -7)
			endDate = now
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(0, 1, 0)
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(1, 0, 0)
		default:
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
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

	// 1. Pending Verifications
	pendingVerifications := h.getPendingVerifications()

	// 2. Documents yang perlu diverifikasi
	pendingDocuments := h.getPendingDocuments()

	// 3. Timeline per pendaftaran (recent activities)
	timeline := h.getRegistrationTimeline(userID.(string))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard staff berhasil diambil",
		"data": gin.H{
			"pending_verifications": pendingVerifications,
			"pending_documents":     pendingDocuments,
			"timeline":              timeline,
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

	// Hari ini
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := todayStart.Add(24 * time.Hour)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", todayStart, todayEnd).
		Count(&hariIni)

	// Bulan ini
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", monthStart, monthEnd).
		Count(&bulanIni)

	// Tahun ini
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := yearStart.AddDate(1, 0, 0)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ?", yearStart, yearEnd).
		Count(&tahunIni)

	// Pernikahan selesai dalam periode
	var selesai int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranSelesai).
		Count(&selesai)

	// Status breakdown sesuai flow baru: Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
	var draft, disetujui, menungguPenugasan, penghuluDitugaskan, ditolak int64
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("created_at >= ? AND created_at < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranDraft).
		Count(&draft)
	h.DB.Model(&structs.PendaftaranNikah{}).
		Where("created_at >= ? AND created_at < ? AND status_pendaftaran = ?", startDate, endDate, structs.StatusPendaftaranDisetujui).
		Count(&disetujui)
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
			"draft":                draft,
			"disetujui":            disetujui,
			"menunggu_penugasan":   menungguPenugasan,
			"penghulu_ditugaskan":  penghuluDitugaskan,
			"selesai":              selesai,
			"ditolak":              ditolak,
		},
	}
}

// getMarriageTrends returns marriage trends for chart
func (h *InDB) getMarriageTrends(startDate, endDate time.Time, period string) []gin.H {
	var trends []gin.H

	// Group by period
	if period == "day" {
		// Daily trends
		current := startDate
		for current.Before(endDate) {
			next := current.Add(24 * time.Hour)
			var count int64
			h.DB.Model(&structs.PendaftaranNikah{}).
				Where("tanggal_nikah >= ? AND tanggal_nikah < ?", current, next).
				Count(&count)

			trends = append(trends, gin.H{
				"date":  current.Format("2006-01-02"),
				"label": current.Format("02 Jan"),
				"count": count,
			})
			current = next
		}
	} else if period == "month" {
		// Monthly trends
		current := startDate
		for current.Before(endDate) {
			next := current.AddDate(0, 1, 0)
			var count int64
			h.DB.Model(&structs.PendaftaranNikah{}).
				Where("tanggal_nikah >= ? AND tanggal_nikah < ?", current, next).
				Count(&count)

			trends = append(trends, gin.H{
				"date":  current.Format("2006-01"),
				"label": current.Format("Jan 2006"),
				"count": count,
			})
			current = next
		}
	} else if period == "year" {
		// Yearly trends
		current := startDate
		for current.Before(endDate) {
			next := current.AddDate(1, 0, 0)
			var count int64
			h.DB.Model(&structs.PendaftaranNikah{}).
				Where("tanggal_nikah >= ? AND tanggal_nikah < ?", current, next).
				Count(&count)

			trends = append(trends, gin.H{
				"date":  current.Format("2006"),
				"label": current.Format("2006"),
				"count": count,
			})
			current = next
		}
	}

	return trends
}

// getStatusDistribution returns status distribution for pie chart
// Menggunakan flow baru: Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
func (h *InDB) getStatusDistribution(startDate, endDate time.Time) []gin.H {
	var distribution []gin.H

	// Status sesuai flow baru
	statuses := []string{
		structs.StatusPendaftaranDraft,
		structs.StatusPendaftaranDisetujui,
		structs.StatusPendaftaranMenungguPenugasan,
		structs.StatusPendaftaranPenghuluDitugaskan,
		structs.StatusPendaftaranSelesai,
		structs.StatusPendaftaranDitolak,
	}

	// Hitung per status (menggunakan created_at untuk tracking pendaftaran baru)
	for _, status := range statuses {
		var count int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("created_at >= ? AND created_at < ? AND status_pendaftaran = ?", startDate, endDate, status).
			Count(&count)

		distribution = append(distribution, gin.H{
			"status": status,
			"count":  count,
			"label":  h.getStatusLabel(status),
		})
	}

	return distribution
}

// getStatusLabel returns human-readable label for status
func (h *InDB) getStatusLabel(status string) string {
	labels := map[string]string{
		structs.StatusPendaftaranDraft:              "Draft",
		structs.StatusPendaftaranDisetujui:          "Disetujui",
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

// getPenghuluPerformance returns penghulu performance data
func (h *InDB) getPenghuluPerformance(startDate, endDate time.Time) []gin.H {
	type PenghuluStats struct {
		PenghuluID   uint
		NamaLengkap  string
		JumlahNikah  int64
		Rating       float64
		TotalRating  float64
		JumlahRating int64
	}

	var penghulus []structs.Penghulu
	h.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus)

	var performance []gin.H
	for _, penghulu := range penghulus {
		// Count pernikahan dalam periode
		var jumlahNikah int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran = ?",
				penghulu.ID, startDate, endDate, structs.StatusPendaftaranSelesai).
			Count(&jumlahNikah)

		// Get average rating from feedback
		// First, get all pendaftaran IDs for this penghulu
		var pendaftaranIDs []uint
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("penghulu_id = ? AND status_pendaftaran = ?", penghulu.ID, structs.StatusPendaftaranSelesai).
			Pluck("id", &pendaftaranIDs)

		var avgRating float64
		var ratingCount int64
		if len(pendaftaranIDs) > 0 {
			type RatingResult struct {
				AvgRating float64
				Count     int64
			}
			var result RatingResult
			h.DB.Model(&structs.FeedbackPernikahan{}).
				Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
				Where("pendaftaran_id IN ? AND jenis_feedback = ? AND rating IS NOT NULL",
					pendaftaranIDs, structs.FeedbackJenisRating).
				Scan(&result)
			avgRating = result.AvgRating
			ratingCount = result.Count
		}

		performance = append(performance, gin.H{
			"penghulu_id":   penghulu.ID,
			"nama_lengkap":  penghulu.Nama_lengkap,
			"jumlah_nikah":  jumlahNikah,
			"rating":        avgRating,
			"jumlah_rating": ratingCount,
		})
	}

	return performance
}

// getPeakHoursAnalysis returns peak hours analysis
func (h *InDB) getPeakHoursAnalysis(startDate, endDate time.Time) []gin.H {
	// Time slots
	timeSlots := []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00"}

	var peakHours []gin.H
	for _, slot := range timeSlots {
		var count int64
		h.DB.Model(&structs.PendaftaranNikah{}).
			Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND waktu_nikah LIKE ? AND status_pendaftaran = ?",
				startDate, endDate, slot+"%", structs.StatusPendaftaranSelesai).
			Count(&count)

		peakHours = append(peakHours, gin.H{
			"waktu": slot,
			"count": count,
		})
	}

	return peakHours
}

// getPendingVerifications returns pending verifications for staff
// Sesuai flow baru: Draft dan Disetujui perlu diverifikasi
func (h *InDB) getPendingVerifications() []gin.H {
	var pendaftarans []structs.PendaftaranNikah
	// Status yang perlu diverifikasi: Draft (baru daftar) dan Disetujui (perlu verifikasi berkas)
	h.DB.Where("status_pendaftaran IN ?", []string{
		structs.StatusPendaftaranDraft,
		structs.StatusPendaftaranDisetujui,
	}).
		Order("created_at DESC").
		Limit(10).
		Find(&pendaftarans)

	var pending []gin.H
	for _, p := range pendaftarans {
		// Get calon suami and istri names
		var calonSuami, calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		pending = append(pending, gin.H{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"tanggal_nikah":      p.Tanggal_nikah.Format("2006-01-02"),
			"calon_suami":        calonSuami.Nama_lengkap,
			"calon_istri":        calonIstri.Nama_lengkap,
			"created_at":         p.Created_at,
		})
	}

	return pending
}

// getPendingDocuments returns documents that need verification
// Sesuai flow baru: Draft perlu verifikasi form, Disetujui perlu verifikasi dokumen fisik
func (h *InDB) getPendingDocuments() []gin.H {
	var pendaftarans []structs.PendaftaranNikah
	// Status yang perlu verifikasi dokumen: Draft (verifikasi form) dan Disetujui (verifikasi dokumen fisik)
	h.DB.Where("status_pendaftaran IN ?", []string{
		structs.StatusPendaftaranDraft,
		structs.StatusPendaftaranDisetujui,
	}).
		Order("created_at DESC").
		Limit(10).
		Find(&pendaftarans)

	var pending []gin.H
	for _, p := range pendaftarans {
		var calonSuami, calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		pending = append(pending, gin.H{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"calon_suami":        calonSuami.Nama_lengkap,
			"calon_istri":        calonIstri.Nama_lengkap,
			"created_at":         p.Created_at,
			"needs_verification": true,
		})
	}

	return pending
}

// getRegistrationTimeline returns recent registration activities
func (h *InDB) getRegistrationTimeline(staffID string) []gin.H {
	var pendaftarans []structs.PendaftaranNikah
	h.DB.Order("updated_at DESC").
		Limit(20).
		Find(&pendaftarans)

	var timeline []gin.H
	for _, p := range pendaftarans {
		var calonSuami, calonIstri structs.CalonPasangan
		h.DB.Where("id = ?", p.Calon_suami_id).First(&calonSuami)
		h.DB.Where("id = ?", p.Calon_istri_id).First(&calonIstri)

		timeline = append(timeline, gin.H{
			"id":                 p.ID,
			"nomor_pendaftaran":  p.Nomor_pendaftaran,
			"status_pendaftaran": p.Status_pendaftaran,
			"calon_suami":        calonSuami.Nama_lengkap,
			"calon_istri":        calonIstri.Nama_lengkap,
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
	now := time.Now()

	var startDate, endDate time.Time
	switch period {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endDate = startDate.Add(24 * time.Hour)
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(1, 0, 0)
	default:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
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
	now := time.Now()

	if dateFrom != "" && dateTo != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_from tidak valid",
			})
			return
		}
		endDate, err = time.Parse("2006-01-02", dateTo)
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
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
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
	now := time.Now()

	if dateFrom != "" && dateTo != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format date_from tidak valid",
			})
			return
		}
		endDate, err = time.Parse("2006-01-02", dateTo)
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
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)
	}

	peakHours := h.getPeakHoursAnalysis(startDate, endDate)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analisis jam sibuk berhasil diambil",
		"data":    peakHours,
	})
}

