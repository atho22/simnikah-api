package kepala_kua

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "simnikah/internal/models"
	"simnikah/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type registrationURI struct {
	RegistrationID string `uri:"id" binding:"required"`
}

type assignPenghuluApprovalRequest struct {
	PenghuluID    uint   `json:"penghulu_id" binding:"required"`
	ApprovalNotes string `json:"approval_notes" binding:"required"`
}

type forwardChainingConfigResponse struct {
	Source                  string                 `json:"source"`
	DynamicConfigReady      bool                   `json:"dynamic_config_ready"`
	MinimumRating           float64                `json:"minimum_rating"`
	CapacityPerDay          int                    `json:"capacity_per_day"`
	CapacityPerHour         int                    `json:"capacity_per_hour"`
	KuaLatitude             float64                `json:"kua_latitude"`
	KuaLongitude            float64                `json:"kua_longitude"`
	ScoringWeights          services.ScoringWeights `json:"scoring_weights"`
	RuleConstraintNotes     []string               `json:"rule_constraint_notes"`
	SystemConfigTableName   string                 `json:"system_config_table_name"`
	SystemConfigKeysExample []string               `json:"system_config_keys_example"`
}

func (h *InDB) parseRegistrationID(c *gin.Context) (uint, error) {
	var uri registrationURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return 0, fmt.Errorf("parameter registrasi tidak valid: %w", err)
	}

	parsedID, err := strconv.ParseUint(strings.TrimSpace(uri.RegistrationID), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("ID registrasi harus berupa angka")
	}

	return uint(parsedID), nil
}

func (h *InDB) requireKepalaKua(c *gin.Context) (string, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "user_id tidak ditemukan di context",
		})
		return "", false
	}

	roleValue, roleExists := c.Get("role")
	if !roleExists || roleValue.(string) != structs.UserRoleKepalaKUA {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   "Hanya user dengan role Kepala KUA yang dapat mengakses aksi ini",
		})
		return "", false
	}

	userID := fmt.Sprint(userIDValue)
	var user structs.Users
	if err := h.DB.Select("user_id, role, status").Where("user_id = ?", userID).
		First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   "Akun pengguna tidak valid atau tidak aktif",
		})
		return "", false
	}

	if user.Role != structs.UserRoleKepalaKUA || user.Status != structs.UserStatusAktif {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   "Role atau status akun tidak memenuhi syarat Kepala KUA aktif",
		})
		return "", false
	}

	return userID, true
}

func lockRegistrationForUpdate(tx *gorm.DB, registrationID uint) (*structs.PendaftaranNikah, error) {
	var registration structs.PendaftaranNikah
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", registrationID).First(&registration).Error; err != nil {
		return nil, err
	}
	return &registration, nil
}

func (h *InDB) RecommendPenghuluWithForwardChaining(c *gin.Context) {
	registrationID, err := h.parseRegistrationID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	fcEngine := services.NewForwardChainingEngine(h.DB)
	recommendation, err := fcEngine.GetPenghuluRecommendations(registrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mendapatkan rekomendasi",
			"error":   err.Error(),
		})
		return
	}

	recommendedName := ""
	if recommendation.RecommendedPenghulu != nil {
		recommendedName = recommendation.RecommendedPenghulu.Nama_lengkap
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rekomendasi penghulu berhasil didapatkan",
		"data": gin.H{
			"recommended_penghulu_id":   recommendation.RecommendedPenghuluID,
			"recommended_penghulu_name": recommendedName,
			"selected_score":            recommendation.SelectedScore,
			"confidence":                recommendation.Confidence,
			"decision_reasoning":        recommendation.DecisionReasoning,
			"alternatives":              recommendation.Alternatives,
			"evaluated_at":              recommendation.EvaluatedAt,
			"evaluation_count":          recommendation.EvaluationCount,
		},
	})
}

// GetDetailedEvaluationReport mendapatkan laporan evaluasi lengkap untuk semua penghulu.
// Endpoint: GET /simnikah/kepala-kua/evaluation-report/:id
func (h *InDB) GetDetailedEvaluationReport(c *gin.Context) {
	registrationID, err := h.parseRegistrationID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	fcEngine := services.NewForwardChainingEngine(h.DB)
	evaluation, err := fcEngine.GetDetailedEvaluation(registrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mendapatkan laporan evaluasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Laporan evaluasi berhasil didapatkan",
		"data":    evaluation,
	})
}

// AssignPenghuluWithApproval melakukan assignment final berbasis input frontend yang eksplisit.
// Frontend wajib mengirim PenghuluID yang valid untuk menghindari race condition akibat perhitungan ulang rekomendasi.
// Transaction dan row lock dipakai agar proses update pendaftaran konsisten bila ada request paralel.
func (h *InDB) AssignPenghuluWithApproval(c *gin.Context) {
	registrationID, err := h.parseRegistrationID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	currentUserID, ok := h.requireKepalaKua(c)
	if !ok {
		return
	}

	var input assignPenghuluApprovalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if input.PenghuluID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Penghulu ID harus disediakan",
			"error":   "frontend wajib mengirim penghulu_id yang valid",
		})
		return
	}

	approvalNotes := strings.TrimSpace(input.ApprovalNotes)
	if approvalNotes == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Approval notes wajib diisi",
			"error":   "approval_notes tidak boleh kosong",
		})
		return
	}

	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memulai transaksi",
			"error":   tx.Error.Error(),
		})
		return
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback().Error
		}
	}()

	registration, err := lockRegistrationForUpdate(tx, registrationID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": "Pendaftaran tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	if registration.Status_pendaftaran != structs.StatusPendaftaranMenungguPenugasan {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Status pendaftaran tidak dapat di-assign",
			"error":   fmt.Sprintf("status saat ini: %s, harus 'Menunggu Penugasan'", registration.Status_pendaftaran),
		})
		return
	}

	var penghulu structs.Penghulu
	if err := tx.Where("id = ? AND status = ?", input.PenghuluID, structs.PenghuluStatusAktif).First(&penghulu).Error; err != nil {
		status := http.StatusBadRequest
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": "Penghulu tidak valid atau tidak aktif",
			"error":   err.Error(),
		})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"penghulu_id":           input.PenghuluID,
		"status_pendaftaran":    structs.StatusPendaftaranPenghuluDitugaskan,
		"updated_at":            now,
	}

	if err := tx.Model(&structs.PendaftaranNikah{}).
		Where("id = ?", registrationID).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui pendaftaran",
			"error":   err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal commit transaksi",
			"error":   err.Error(),
		})
		return
	}
	committed = true

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Penghulu berhasil ditugaskan",
		"data": gin.H{
			"registration_id":    registrationID,
			"nama_suami":         registration.Nama_suami,
			"nama_istri":         registration.Nama_istri,
			"tanggal_nikah":      registration.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        registration.Waktu_nikah,
			"tempat_nikah":       registration.Tempat_nikah,
			"penghulu_id":        penghulu.ID,
			"penghulu_nama":      penghulu.Nama_lengkap,
			"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
			"assigned_by":        currentUserID,
			"assigned_at":        now,
		},
	})
}

// ==================== CONFIGURATION ENDPOINTS ====================

// GetForwardChainingConfig menyiapkan respons konfigurasi yang mudah dipindah ke sumber dinamis.
// Kerangka ini sengaja menyisakan metadata agar parameter rule dapat diambil dari tabel system_configs
// (misalnya: key = forward_chaining.minimum_rating, forward_chaining.capacity_per_day, forward_chaining.weights.rating, dst)
// tanpa perlu re-deploy aplikasi.
func (h *InDB) GetForwardChainingConfig(c *gin.Context) {
	fcEngine := services.NewForwardChainingEngine(h.DB)

	response := forwardChainingConfigResponse{
		Source:             "engine_defaults",
		DynamicConfigReady: true,
		MinimumRating:      fcEngine.Config.MinimumRating,
		CapacityPerDay:     fcEngine.Config.CapacityPerDay,
		CapacityPerHour:    fcEngine.Config.CapacityPerHour,
		KuaLatitude:        fcEngine.Config.KuaLatitude,
		KuaLongitude:       fcEngine.Config.KuaLongitude,
		ScoringWeights:     fcEngine.Config.ScoringWeights,
		RuleConstraintNotes: []string{
			"Flow ini hanya memakai constraint jadwal: status penghulu, bentrok slot, kapasitas harian, kapasitas per jam, dan lokasi nikah",
			"Sumber dinamis yang direkomendasikan: tabel system_configs",
			"Contoh key: forward_chaining.capacity_per_day, forward_chaining.capacity_per_hour, forward_chaining.minimum_rating",
			"Nilai engine saat ini tetap menjadi fallback jika config DB belum tersedia",
		},
		SystemConfigTableName: "system_configs",
		SystemConfigKeysExample: []string{
			"forward_chaining.minimum_rating",
			"forward_chaining.capacity_per_day",
			"forward_chaining.capacity_per_hour",
			"forward_chaining.weights.rating",
			"forward_chaining.weights.availability",
			"forward_chaining.weights.fairness",
			"forward_chaining.weights.location_match",
			"forward_chaining.weights.distance",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Konfigurasi forward chaining engine",
		"data":    response,
	})
}

// ListPenghuluComparison mendapatkan daftar penghulu dengan perbandingan.
// Endpoint: GET /simnikah/kepala-kua/penghulu-comparison
func (h *InDB) ListPenghuluComparison(c *gin.Context) {
	var penghulus []struct {
		ID          uint
		NamaLengkap string
		Rating      float64
		JumlahNikah int
		Status      string
		No_hp       string
		Email       string
	}

	if err := h.DB.Table("penghulus").
		Select("id, nama_lengkap, rating, jumlah_nikah, status, no_hp, email").
		Where("status = ?", structs.PenghuluStatusAktif).
		Order("rating DESC").
		Find(&penghulus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar perbandingan penghulu",
		"data":    penghulus,
	})
}
