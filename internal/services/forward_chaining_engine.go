package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	neturl "net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	structs "simnikah/internal/models"
	"simnikah/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ==================== FORWARD CHAINING ENGINE ====================

const (
	defaultKUALatitude  = -3.291304649442475
	defaultKUALongitude = 114.58814746634684

	derivedFactAdminApproved     = "Lulus Syarat Administrasi"
	derivedFactPenghuluAktif     = "Penghulu Aktif"
	derivedFactJadwalBebas       = "Jadwal Bebas"
	derivedFactKapasitasHarian   = "Kapasitas Harian Memadai"
	derivedFactKapasitasPerJam   = "Kapasitas Per Jam Memadai"
	derivedFactLokasiSesuai      = "Lokasi Sesuai"
	derivedFactRatingMemadai     = "Rating Memadai"
	derivedFactJarakMemadai      = "Jarak Layak"
	derivedFactDirekomendasikan  = "Penghulu Direkomendasikan"
	derivedFactDitolak           = "Penghulu Ditolak"
	allowedTimeFormat            = "15:04"
)

type ForwardChainingEngine struct {
	DB         *gorm.DB
	Config     FCConfig
	HTTPClient *http.Client
}

type FCConfig struct {
	MinimumRating             float64
	CapacityPerDay            int
	CapacityPerHour           int
	KuaLatitude               float64
	KuaLongitude              float64
	OSRMBaseURL               string
	OSRMFailurePenaltyMeters  float64
	DistancePenaltyPerKm      float64
	HTTPTimeout               time.Duration
	GoogleHolidayCalendarID   string
	GoogleAPIKey              string
	ScoringWeights            ScoringWeights
}

type ScoringWeights struct {
	RatingWeight        float64
	AvailabilityWeight   float64
	FairnessWeight      float64
	LocationMatchWeight float64
	DistanceWeight      float64
}

// ==================== FACTS ====================

type MarriageRegistrationFact struct {
	RegistrationID    uint
	TanggalNikah      time.Time
	WaktuNikah        string
	TempatNikah       string
	AlamatNikah       string
	Latitude          float64
	Longitude         float64
	StatusPendaftaran string
	PenghuluID        *uint
}

type PenghuluFact struct {
	PenghuluID      uint
	NamaPenghulu    string
	StatusAktif     string
	JumlahNikah     int
	CapacityPerDay  int
	CapacityPerHour int
	Rating          float64
}

// ==================== INFERENCE AUDIT ====================

type RuleEvaluation struct {
	RuleID      string
	RuleName    string
	IsSatisfied bool
	Reason      string
	Impact      float64
}

type DerivedFact struct {
	Name      string
	Value     any
	RuleID    string
	Reason    string
	CreatedAt time.Time
}

type RuleResult struct {
	PenghuluID     uint
	EvaluatedRules  []RuleEvaluation
	DerivedFacts    []DerivedFact
	AllRulesPassed  bool
	Score           float64
	DistanceMeters  float64
	Conclusion      string
}

// ==================== RECOMMENDATIONS ====================

type AssignmentRecommendation struct {
	RecommendedPenghuluID uint
	RecommendedPenghulu   *structs.Penghulu
	SelectedScore         float64
	Alternatives          []AlternativePenghulu
	EvaluationProcess     []RuleResult
	DecisionReasoning     string
	Confidence            float64
	EvaluatedAt           time.Time
	EvaluationCount       int
}

type AlternativePenghulu struct {
	PenghuluID   uint
	NamaPenghulu string
	Score        float64
	Reason       string
}

// ==================== SCHEDULING-ONLY API ====================

type ScheduleCheckInput struct {
	TanggalNikah string   `json:"tanggal_nikah" binding:"required"`
	WaktuNikah   string   `json:"waktu_nikah" binding:"required"`
	TempatNikah  string   `json:"tempat_nikah" binding:"required"`
	AlamatNikah  string   `json:"alamat_nikah"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type ScheduleCandidate struct {
	PenghuluID        uint    `json:"penghulu_id"`
	NamaPenghulu      string  `json:"nama_penghulu"`
	JumlahNikah       int     `json:"jumlah_nikah"`
	Rating            float64 `json:"rating"`
	SlotHarianTerisi  int     `json:"slot_harian_terisi"`
	SlotJamTerisi     int     `json:"slot_jam_terisi"`
	Tersedia          bool    `json:"tersedia"`
	Alasan            string  `json:"alasan"`
}

type ScheduleCheckResult struct {
	TanggalNikah        time.Time          `json:"tanggal_nikah"`
	WaktuNikah          string             `json:"waktu_nikah"`
	TempatNikah         string             `json:"tempat_nikah"`
	IsHoliday           bool               `json:"is_holiday"`
	HolidayName         string             `json:"holiday_name"`
	TotalBooked         int64              `json:"total_booked"`
	BookedInKUA         int64              `json:"booked_in_kua"`
	SlotRemaining       int64              `json:"slot_remaining"`
	Available           bool               `json:"available"`
	Reason              string             `json:"reason"`
	RecommendedPenghulu *ScheduleCandidate  `json:"recommended_penghulu,omitempty"`
	Alternatives        []ScheduleCandidate `json:"alternatives"`
}

func (fc *ForwardChainingEngine) CheckScheduleAvailability(input ScheduleCheckInput) (*ScheduleCheckResult, error) {
	tanggalNikah, err := time.Parse("2006-01-02", strings.TrimSpace(input.TanggalNikah))
	if err != nil {
		return nil, fmt.Errorf("format tanggal_nikah harus YYYY-MM-DD")
	}

	waktuNikah, err := normalizeClock(strings.TrimSpace(input.WaktuNikah))
	if err != nil {
		return nil, fmt.Errorf("waktu_nikah tidak valid: %w", err)
	}

	tempatNikah := strings.TrimSpace(input.TempatNikah)
	if tempatNikah == "" {
		tempatNikah = structs.TempatNikahDiKUA
	}

	isHoliday, holidayName := fc.IsHoliday(tanggalNikah)
	startOfDay := time.Date(tanggalNikah.Year(), tanggalNikah.Month(), tanggalNikah.Day(), 0, 0, 0, 0, tanggalNikah.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalBooked int64
	if err := fc.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah = ? AND waktu_nikah = ? AND status_pendaftaran NOT IN ?", tanggalNikah, waktuNikah, []string{structs.StatusPendaftaranDitolak}).
		Count(&totalBooked).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung total jadwal: %w", err)
	}

	var bookedInKUA int64
	if err := fc.DB.Model(&structs.PendaftaranNikah{}).
		Where("tanggal_nikah = ? AND waktu_nikah = ? AND tempat_nikah = ? AND status_pendaftaran NOT IN ?", tanggalNikah, waktuNikah, structs.TempatNikahDiKUA, []string{structs.StatusPendaftaranDitolak}).
		Count(&bookedInKUA).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung jadwal KUA: %w", err)
	}

	var penghulus []structs.Penghulu
	if err := fc.DB.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data penghulu aktif: %w", err)
	}

	var dayRegistrations []structs.PendaftaranNikah
	if err := fc.DB.Where("tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?", startOfDay, endOfDay, []string{structs.StatusPendaftaranDitolak}).Find(&dayRegistrations).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil jadwal harian: %w", err)
	}

	dayCounts := map[uint]int{}
	hourCounts := map[uint]int{}
	for _, registration := range dayRegistrations {
		if registration.Penghulu_id == nil {
			continue
		}
		penghuluID := *registration.Penghulu_id
		dayCounts[penghuluID]++
		clock, clockErr := normalizeClock(registration.Waktu_nikah)
		if clockErr != nil {
			clock = registration.Waktu_nikah
		}
		if clock == waktuNikah {
			hourCounts[penghuluID]++
		}
	}

	totalCapacity := int64(len(penghulus) * fc.Config.CapacityPerDay)

	slotRemaining := totalCapacity - totalBooked
	if slotRemaining < 0 {
		slotRemaining = 0
	}

	result := &ScheduleCheckResult{
		TanggalNikah:  tanggalNikah,
		WaktuNikah:    waktuNikah,
		TempatNikah:   tempatNikah,
		IsHoliday:     isHoliday,
		HolidayName:   holidayName,
		TotalBooked:   totalBooked,
		BookedInKUA:   bookedInKUA,
		SlotRemaining: slotRemaining,
		Available:     true,
	}

	if tempatNikah == structs.TempatNikahDiKUA && isHoliday {
		result.Available = false
		result.Reason = fmt.Sprintf("KUA tutup pada hari libur: %s", holidayName)
		return result, nil
	}

	if totalBooked >= totalCapacity {
		result.Available = false
		result.Reason = "Semua slot penghulu pada jam tersebut sudah penuh"
		return result, nil
	}

	if tempatNikah == structs.TempatNikahDiKUA && bookedInKUA >= 1 {
		result.Available = false
		result.Reason = "Slot nikah di KUA pada jam tersebut sudah terisi"
		return result, nil
	}

	for _, penghulu := range penghulus {
		candidate := ScheduleCandidate{
			PenghuluID:       penghulu.ID,
			NamaPenghulu:     penghulu.Nama_lengkap,
			JumlahNikah:      penghulu.Jumlah_nikah,
			Rating:           penghulu.Rating,
			SlotHarianTerisi: dayCounts[penghulu.ID],
			SlotJamTerisi:    hourCounts[penghulu.ID],
			Tersedia:         dayCounts[penghulu.ID] < fc.Config.CapacityPerDay && hourCounts[penghulu.ID] < fc.Config.CapacityPerHour,
		}
		if candidate.Tersedia {
			candidate.Alasan = "Memenuhi kapasitas harian dan kapasitas per jam"
			if result.RecommendedPenghulu == nil {
				result.RecommendedPenghulu = &candidate
			} else {
				result.Alternatives = append(result.Alternatives, candidate)
			}
		} else {
			candidate.Alasan = "Kapasitas penghulu pada slot tersebut sudah terpakai"
			result.Alternatives = append(result.Alternatives, candidate)
		}
	}

	if result.RecommendedPenghulu == nil {
		result.Available = false
		result.Reason = "Tidak ada penghulu aktif yang tersedia pada slot tersebut"
	}

	return result, nil
}

func (fc *ForwardChainingEngine) ListApprovedAssignmentsForPenghulu(userID string, includeCompleted bool) ([]structs.PendaftaranJadwal, error) {
	var penghulu structs.Penghulu
	if err := fc.DB.Where("user_id = ? AND status = ?", userID, structs.PenghuluStatusAktif).First(&penghulu).Error; err != nil {
		return nil, err
	}

	query := fc.DB.Where("penghulu_id = ? AND status_pendaftaran = ?", penghulu.ID, structs.StatusPendaftaranPenghuluDitugaskan)
	if includeCompleted {
		query = fc.DB.Where("penghulu_id = ? AND status_pendaftaran IN ?", penghulu.ID, []string{structs.StatusPendaftaranPenghuluDitugaskan, structs.StatusPendaftaranSelesai})
	}

	var assignments []structs.PendaftaranJadwal
	if err := query.
		Select("id, nama_suami, umur_suami, nama_istri, umur_istri, tanggal_nikah, waktu_nikah, tempat_nikah, alamat_akad, latitude, longitude, status_pendaftaran, penghulu_id").
		Order("tanggal_nikah ASC, waktu_nikah ASC").Find(&assignments).Error; err != nil {
		return nil, err
	}

	return assignments, nil
}

// ==================== INTERNAL SNAPSHOT ====================

type scheduledAssignment struct {
	RegistrationID uint
	PenghuluID     uint
	WaktuNikah     string
	ScheduledAt    time.Time
	Latitude       *float64
	Longitude      *float64
}

type evaluationSnapshot struct {
	RegistrationFact    *MarriageRegistrationFact
	Holiday             bool
	HolidayName         string
	ActivePenghuluFacts []*PenghuluFact
	DayCounts           map[uint]int
	HourCounts          map[uint]map[string]int
	AssignmentsByPenghulu map[uint][]scheduledAssignment
	MaxHistoricalWorkload int
	MinHistoricalWorkload int
}

type candidateInferenceState struct {
	Facts          map[string]bool
	FiredRules     map[string]bool
	Evaluations    []RuleEvaluation
	DerivedFacts   []DerivedFact
	Rejected       bool
	RejectReason   string
	DistanceMeters float64
	Score          float64
}

type inferenceRule struct {
	ID       string
	Name     string
	Blocking bool
	CanFire  func() bool
	Fire     func() (RuleEvaluation, []DerivedFact)
}

// osrmDistanceCache meng-cache hasil OSRM route berdasarkan pasangan koordinat.
// Mengurangi HTTP call duplikat untuk rute yang sama dalam satu sesi evaluasi.
type osrmDistanceCache struct {
	mu    sync.RWMutex
	cache map[string]float64
}

func newOSRMDistanceCache() *osrmDistanceCache {
	return &osrmDistanceCache{cache: make(map[string]float64)}
}

func (c *osrmDistanceCache) get(fromLat, fromLon, toLat, toLon float64) (float64, bool) {
	key := fmt.Sprintf("%.6f,%.6f->%.6f,%.6f", fromLat, fromLon, toLat, toLon)
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.cache[key]
	return val, ok
}

func (c *osrmDistanceCache) set(fromLat, fromLon, toLat, toLon, distance float64) {
	key := fmt.Sprintf("%.6f,%.6f->%.6f,%.6f", fromLat, fromLon, toLat, toLon)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = distance
}

// ==================== CONSTRUCTOR ====================

func NewForwardChainingEngine(db *gorm.DB) *ForwardChainingEngine {
	return NewForwardChainingEngineWithConfig(db, FCConfig{
		MinimumRating:            3.0,
		CapacityPerDay:           3,
		CapacityPerHour:          1,
		KuaLatitude:              defaultKUALatitude,
		KuaLongitude:             defaultKUALongitude,
		OSRMBaseURL:              "https://router.project-osrm.org",
		OSRMFailurePenaltyMeters: 1500,
		DistancePenaltyPerKm:     8,
		HTTPTimeout:              6 * time.Second,
		GoogleHolidayCalendarID:  os.Getenv("GOOGLE_HOLIDAYS_CALENDAR_ID"),
		GoogleAPIKey:             os.Getenv("GOOGLE_API_KEY"),
		ScoringWeights: ScoringWeights{
			RatingWeight:        0.35,
			AvailabilityWeight:   0.25,
			FairnessWeight:       0.20,
			LocationMatchWeight:  0.10,
			DistanceWeight:       0.10,
		},
	})
}

func NewForwardChainingEngineWithConfig(db *gorm.DB, cfg FCConfig) *ForwardChainingEngine {
	cfg = normalizeFCConfig(cfg)
	return &ForwardChainingEngine{
		DB:     db,
		Config: cfg,
		HTTPClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}
}

func normalizeFCConfig(cfg FCConfig) FCConfig {
	if cfg.MinimumRating <= 0 {
		cfg.MinimumRating = 3.0
	}
	if cfg.CapacityPerDay <= 0 {
		cfg.CapacityPerDay = 3
	}
	if cfg.CapacityPerHour <= 0 {
		cfg.CapacityPerHour = 1
	}
	if cfg.KuaLatitude == 0 && cfg.KuaLongitude == 0 {
		cfg.KuaLatitude = defaultKUALatitude
		cfg.KuaLongitude = defaultKUALongitude
	}
	if cfg.OSRMBaseURL == "" {
		cfg.OSRMBaseURL = "https://router.project-osrm.org"
	}
	if cfg.OSRMFailurePenaltyMeters <= 0 {
		cfg.OSRMFailurePenaltyMeters = 1500
	}
	if cfg.DistancePenaltyPerKm <= 0 {
		cfg.DistancePenaltyPerKm = 8
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = 6 * time.Second
	}
	return cfg
}

// ==================== HELPERS ====================

func clampFloat(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func normalizeClock(clock string) (string, error) {
	trimmed := strings.TrimSpace(clock)
	if trimmed == "" {
		return "", fmt.Errorf("waktu tidak boleh kosong")
	}
	parsed, err := time.Parse(allowedTimeFormat, trimmed)
	if err != nil {
		parts := strings.Split(trimmed, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("format waktu tidak valid, gunakan HH:MM")
		}
		hour, errHour := strconv.Atoi(parts[0])
		minute, errMinute := strconv.Atoi(parts[1])
		if errHour != nil || errMinute != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
			return "", fmt.Errorf("format waktu tidak valid, gunakan HH:MM")
		}
		return fmt.Sprintf("%02d:%02d", hour, minute), nil
	}
	return parsed.Format(allowedTimeFormat), nil
}

func combineDateAndClock(date time.Time, clock string) (time.Time, error) {
	normalized, err := normalizeClock(clock)
	if err != nil {
		return time.Time{}, err
	}
	parts := strings.Split(normalized, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid clock format: %s", clock)
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location()), nil
}

func haversineMeters(fromLat, fromLon, toLat, toLon float64) float64 {
	const earthRadius = 6371000.0
	lat1 := fromLat * math.Pi / 180.0
	lat2 := toLat * math.Pi / 180.0
	dLat := (toLat - fromLat) * math.Pi / 180.0
	dLon := (toLon - fromLon) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func isAllowedAssignmentStatus(status string) bool {
	switch status {
	case structs.StatusPendaftaranMenungguPenugasan:
		return true
	default:
		return false
	}
}

func (s *candidateInferenceState) hasFact(name string) bool {
	return s.Facts[name]
}

func (s *candidateInferenceState) addFact(name string, value any, ruleID, reason string) bool {
	if s.Facts == nil {
		s.Facts = map[string]bool{}
	}
	if s.Facts[name] {
		return false
	}
	s.Facts[name] = true
	s.DerivedFacts = append(s.DerivedFacts, DerivedFact{
		Name:      name,
		Value:     value,
		RuleID:    ruleID,
		Reason:    reason,
		CreatedAt: time.Now(),
	})
	return true
}

func (s *candidateInferenceState) appendEvaluation(e RuleEvaluation) {
	s.Evaluations = append(s.Evaluations, e)
}

func (fc *ForwardChainingEngine) routeDistanceMeters(fromLat, fromLon, toLat, toLon float64) (float64, error) {
	baseURL := strings.TrimRight(fc.Config.OSRMBaseURL, "/")
	reqURL := fmt.Sprintf(
		"%s/route/v1/driving/%f,%f;%f,%f?overview=false&alternatives=false&steps=false",
		baseURL,
		fromLon, fromLat, toLon, toLat,
	)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := fc.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("osrm returned status %s", resp.Status)
	}

	var parsed struct {
		Routes []struct {
			Distance float64 `json:"distance"`
		} `json:"routes"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, err
	}
	if len(parsed.Routes) == 0 || parsed.Routes[0].Distance <= 0 {
		return 0, fmt.Errorf("osrm returned empty route")
	}
	return parsed.Routes[0].Distance, nil
}

// routeDistanceCached wraps routeDistanceMeters dengan in-memory cache.
// Jika OSRM gagal, fallback ke Haversine * 1.35 (road factor).
func (fc *ForwardChainingEngine) routeDistanceCached(cache *osrmDistanceCache, fromLat, fromLon, toLat, toLon float64) float64 {
	if dist, found := cache.get(fromLat, fromLon, toLat, toLon); found {
		return dist
	}

	distance, err := fc.routeDistanceMeters(fromLat, fromLon, toLat, toLon)
	if err != nil {
		// Fallback Haversine dengan road factor 1.35 + penalty
		distance = haversineMeters(fromLat, fromLon, toLat, toLon)*1.35 + fc.Config.OSRMFailurePenaltyMeters
	}

	cache.set(fromLat, fromLon, toLat, toLon, distance)
	return distance
}

func (fc *ForwardChainingEngine) IsHoliday(date time.Time) (bool, string) {
	if date.Weekday() == time.Sunday {
		return true, "Minggu"
	}
	if fc.Config.GoogleHolidayCalendarID == "" || fc.Config.GoogleAPIKey == "" {
		return false, ""
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).UTC().Format(time.RFC3339)
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location()).UTC().Format(time.RFC3339)
	reqURL := fmt.Sprintf(
		"https://www.googleapis.com/calendar/v3/calendars/%s/events?timeMin=%s&timeMax=%s&singleEvents=true&key=%s",
		neturl.QueryEscape(fc.Config.GoogleHolidayCalendarID),
		neturl.QueryEscape(start),
		neturl.QueryEscape(end),
		neturl.QueryEscape(fc.Config.GoogleAPIKey),
	)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return false, ""
	}
	resp, err := fc.HTTPClient.Do(req)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, ""
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, ""
	}

	var parsed struct {
		Items []struct {
			Summary string `json:"summary"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return false, ""
	}
	if len(parsed.Items) > 0 {
		return true, parsed.Items[0].Summary
	}
	return false, ""
}

// ==================== FACT EXTRACTION ====================

func (fc *ForwardChainingEngine) ExtractMarriageRegistrationFacts(registrationID uint) (*MarriageRegistrationFact, error) {
	return fc.extractMarriageRegistrationFactsWithDB(fc.DB, registrationID)
}

func (fc *ForwardChainingEngine) extractMarriageRegistrationFactsWithDB(db *gorm.DB, registrationID uint) (*MarriageRegistrationFact, error) {
	var reg structs.PendaftaranNikah
	if err := db.Where("id = ?", registrationID).First(&reg).Error; err != nil {
		return nil, fmt.Errorf("registrasi tidak ditemukan: %w", err)
	}

	waktu, _ := normalizeClock(reg.Waktu_nikah)
	if waktu == "" {
		waktu = reg.Waktu_nikah
	}

	fact := &MarriageRegistrationFact{
		RegistrationID:    reg.ID,
		TanggalNikah:      reg.Tanggal_nikah,
		WaktuNikah:        waktu,
		TempatNikah:       strings.TrimSpace(reg.Tempat_nikah),
		AlamatNikah:       reg.Alamat_akad,
		StatusPendaftaran: reg.Status_pendaftaran,
		PenghuluID:        reg.Penghulu_id,
	}
	if reg.Latitude != nil {
		fact.Latitude = *reg.Latitude
	}
	if reg.Longitude != nil {
		fact.Longitude = *reg.Longitude
	}
	return fact, nil
}

func (fc *ForwardChainingEngine) ExtractPenghuluFacts(penghuluID uint) (*PenghuluFact, error) {
	return fc.extractPenghuluFactsWithDB(fc.DB, penghuluID)
}

func (fc *ForwardChainingEngine) extractPenghuluFactsWithDB(db *gorm.DB, penghuluID uint) (*PenghuluFact, error) {
	var penghulu structs.Penghulu
	if err := db.Where("id = ?", penghuluID).First(&penghulu).Error; err != nil {
		return nil, fmt.Errorf("penghulu tidak ditemukan: %w", err)
	}

	return &PenghuluFact{
		PenghuluID:      penghulu.ID,
		NamaPenghulu:    penghulu.Nama_lengkap,
		StatusAktif:     penghulu.Status,
		JumlahNikah:     penghulu.Jumlah_nikah,
		CapacityPerDay:  fc.Config.CapacityPerDay,
		CapacityPerHour: fc.Config.CapacityPerHour,
		Rating:          penghulu.Rating,
	}, nil
}

// loadEvaluationSnapshot mengambil snapshot data untuk inferensi FC.
// Parameter lockRows=false (read-only) untuk rekomendasi, true untuk assignment.
func (fc *ForwardChainingEngine) loadEvaluationSnapshot(registrationID uint, lockRows bool) (*evaluationSnapshot, error) {
	db := fc.DB
	var tx *gorm.DB
	if lockRows {
		tx = fc.DB.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}
		db = tx
		defer func() {
			tx.Rollback()
		}()
	}

	regFact, err := fc.extractMarriageRegistrationFactsWithDB(db, registrationID)
	if err != nil {
		return nil, err
	}

	holiday, holidayName := fc.IsHoliday(regFact.TanggalNikah)

	var penghulus []structs.Penghulu
	if err := db.Where("status = ?", structs.PenghuluStatusAktif).Find(&penghulus).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil data penghulu aktif: %w", err)
	}

	startOfDay := time.Date(regFact.TanggalNikah.Year(), regFact.TanggalNikah.Month(), regFact.TanggalNikah.Day(), 0, 0, 0, 0, regFact.TanggalNikah.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.Model(&structs.PendaftaranNikah{}).
		Where("id <> ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?", registrationID, startOfDay, endOfDay, []string{structs.StatusPendaftaranDitolak})
	if lockRows {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var assignments []structs.PendaftaranNikah
	if err := query.Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil jadwal harian: %w", err)
	}

	dayCounts := map[uint]int{}
	hourCounts := map[uint]map[string]int{}
	assignmentsByPenghulu := map[uint][]scheduledAssignment{}
	for _, assignment := range assignments {
		if assignment.Penghulu_id == nil {
			continue
		}
		penghuluID := *assignment.Penghulu_id
		dayCounts[penghuluID]++
		if hourCounts[penghuluID] == nil {
			hourCounts[penghuluID] = map[string]int{}
		}
		clock, clockErr := normalizeClock(assignment.Waktu_nikah)
		if clockErr != nil {
			clock = assignment.Waktu_nikah
		}
		hourCounts[penghuluID][clock]++

		scheduledAt, err := combineDateAndClock(regFact.TanggalNikah, clock)
		if err != nil {
			scheduledAt = regFact.TanggalNikah
		}
		assignmentsByPenghulu[penghuluID] = append(assignmentsByPenghulu[penghuluID], scheduledAssignment{
			RegistrationID: assignment.ID,
			PenghuluID:     penghuluID,
			WaktuNikah:     clock,
			ScheduledAt:    scheduledAt,
			Latitude:       assignment.Latitude,
			Longitude:      assignment.Longitude,
		})
	}

	minWorkload := 0
	maxWorkload := 0
	for i, penghulu := range penghulus {
		if i == 0 || penghulu.Jumlah_nikah < minWorkload {
			minWorkload = penghulu.Jumlah_nikah
		}
		if penghulu.Jumlah_nikah > maxWorkload {
			maxWorkload = penghulu.Jumlah_nikah
		}
	}

	activeFacts := make([]*PenghuluFact, 0, len(penghulus))
	for _, penghulu := range penghulus {
		p := penghulu
		activeFacts = append(activeFacts, &PenghuluFact{
			PenghuluID:      p.ID,
			NamaPenghulu:    p.Nama_lengkap,
			StatusAktif:     p.Status,
			JumlahNikah:     p.Jumlah_nikah,
			CapacityPerDay:  fc.Config.CapacityPerDay,
			CapacityPerHour: fc.Config.CapacityPerHour,
			Rating:          p.Rating,
		})
	}

	return &evaluationSnapshot{
		RegistrationFact:      regFact,
		Holiday:               holiday,
		HolidayName:           holidayName,
		ActivePenghuluFacts:   activeFacts,
		DayCounts:             dayCounts,
		HourCounts:            hourCounts,
		AssignmentsByPenghulu: assignmentsByPenghulu,
		MaxHistoricalWorkload: maxWorkload,
		MinHistoricalWorkload: minWorkload,
	}, nil
}

// ==================== RULE ENGINE ====================

func (fc *ForwardChainingEngine) buildRuleBase(snapshot *evaluationSnapshot, candidate *PenghuluFact, state *candidateInferenceState, osrmCache *osrmDistanceCache) []inferenceRule {
	registration := snapshot.RegistrationFact
	dayCount := snapshot.DayCounts[candidate.PenghuluID]
	hourCount := snapshot.HourCounts[candidate.PenghuluID][registration.WaktuNikah]

	return []inferenceRule{
		{
			ID:       "RULE_001",
			Name:     "Validasi Administrasi",
			Blocking: true,
			CanFire: func() bool { return !state.hasFact(derivedFactAdminApproved) && !state.Rejected },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				if isAllowedAssignmentStatus(registration.StatusPendaftaran) {
					eval := RuleEvaluation{RuleID: "RULE_001", RuleName: "Validasi Administrasi", IsSatisfied: true, Reason: "Status pendaftaran memenuhi syarat penugasan", Impact: 10}
					return eval, []DerivedFact{{Name: derivedFactAdminApproved, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
				}
				eval := RuleEvaluation{RuleID: "RULE_001", RuleName: "Validasi Administrasi", IsSatisfied: false, Reason: fmt.Sprintf("Status pendaftaran %s belum siap untuk penugasan", registration.StatusPendaftaran)}
				return eval, nil
			},
		},
		{
			ID:       "RULE_002",
			Name:     "Validasi Status Penghulu",
			Blocking: true,
			CanFire: func() bool { return state.hasFact(derivedFactAdminApproved) && !state.hasFact(derivedFactPenghuluAktif) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				if candidate.StatusAktif == structs.PenghuluStatusAktif {
					eval := RuleEvaluation{RuleID: "RULE_002", RuleName: "Validasi Status Penghulu", IsSatisfied: true, Reason: "Penghulu aktif", Impact: 10}
					return eval, []DerivedFact{{Name: derivedFactPenghuluAktif, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
				}
				eval := RuleEvaluation{RuleID: "RULE_002", RuleName: "Validasi Status Penghulu", IsSatisfied: false, Reason: fmt.Sprintf("Penghulu berstatus %s", candidate.StatusAktif)}
				return eval, nil
			},
		},
		{
			ID:       "RULE_003",
			Name:     "Konteks Hari Libur",
			Blocking: false,
			CanFire: func() bool { return state.hasFact(derivedFactPenghuluAktif) && !state.hasFact("Konteks Hari") },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				factName := "Hari Libur"
				reason := "Hari kerja"
				value := false
				if snapshot.Holiday {
					factName = "Hari Libur"
					reason = fmt.Sprintf("Tanggal jatuh pada hari libur: %s", snapshot.HolidayName)
					value = true
				}
				eval := RuleEvaluation{RuleID: "RULE_003", RuleName: "Konteks Hari Libur", IsSatisfied: true, Reason: reason, Impact: 0}
				return eval, []DerivedFact{{Name: "Konteks Hari", Value: value, RuleID: eval.RuleID, Reason: reason, CreatedAt: time.Now()}, {Name: factName, Value: value, RuleID: eval.RuleID, Reason: reason, CreatedAt: time.Now()}}
			},
		},
		{
			ID:       "RULE_004",
			Name:     "Cek Konflik Jadwal",
			Blocking: true,
			CanFire: func() bool { return state.hasFact(derivedFactPenghuluAktif) && !state.hasFact(derivedFactJadwalBebas) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				if dayCount == 0 || hourCount == 0 {
					eval := RuleEvaluation{RuleID: "RULE_004", RuleName: "Cek Konflik Jadwal", IsSatisfied: true, Reason: "Tidak ada bentrok pada slot waktu yang sama", Impact: 20}
					return eval, []DerivedFact{{Name: derivedFactJadwalBebas, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
				}
				eval := RuleEvaluation{RuleID: "RULE_004", RuleName: "Cek Konflik Jadwal", IsSatisfied: false, Reason: fmt.Sprintf("Penghulu sudah terisi pada %s pukul %s", registration.TanggalNikah.Format("02-01-2006"), registration.WaktuNikah)}
				return eval, nil
			},
		},
		{
			ID:       "RULE_005",
			Name:     "Cek Kapasitas Harian",
			Blocking: true,
			CanFire: func() bool { return state.hasFact(derivedFactJadwalBebas) && !state.hasFact(derivedFactKapasitasHarian) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				if dayCount < candidate.CapacityPerDay {
					eval := RuleEvaluation{RuleID: "RULE_005", RuleName: "Cek Kapasitas Harian", IsSatisfied: true, Reason: fmt.Sprintf("Kapasitas harian tersedia (%d/%d)", dayCount, candidate.CapacityPerDay), Impact: 15}
					return eval, []DerivedFact{{Name: derivedFactKapasitasHarian, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
				}
				eval := RuleEvaluation{RuleID: "RULE_005", RuleName: "Cek Kapasitas Harian", IsSatisfied: false, Reason: fmt.Sprintf("Kapasitas harian penuh (%d/%d)", dayCount, candidate.CapacityPerDay)}
				return eval, nil
			},
		},
		{
			ID:       "RULE_006",
			Name:     "Cek Kapasitas Per Jam",
			Blocking: true,
			CanFire: func() bool { return state.hasFact(derivedFactKapasitasHarian) && !state.hasFact(derivedFactKapasitasPerJam) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				if hourCount < candidate.CapacityPerHour {
					eval := RuleEvaluation{RuleID: "RULE_006", RuleName: "Cek Kapasitas Per Jam", IsSatisfied: true, Reason: fmt.Sprintf("Kapasitas per jam tersedia (%d/%d)", hourCount, candidate.CapacityPerHour), Impact: 25}
					return eval, []DerivedFact{{Name: derivedFactKapasitasPerJam, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
				}
				eval := RuleEvaluation{RuleID: "RULE_006", RuleName: "Cek Kapasitas Per Jam", IsSatisfied: false, Reason: fmt.Sprintf("Kapasitas per jam penuh (%d/%d)", hourCount, candidate.CapacityPerHour)}
				return eval, nil
			},
		},
		{
			ID:       "RULE_007",
			Name:     "Cek Kesesuaian Lokasi",
			Blocking: true,
			CanFire: func() bool { return state.hasFact(derivedFactKapasitasPerJam) && !state.hasFact(derivedFactLokasiSesuai) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				eval := RuleEvaluation{RuleID: "RULE_007", RuleName: "Cek Kesesuaian Lokasi", IsSatisfied: true, Reason: "Semua penghulu dapat melayani di dalam maupun luar KUA", Impact: 5}
				return eval, []DerivedFact{{Name: derivedFactLokasiSesuai, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
			},
		},
		{
			ID:       "RULE_008",
			Name:     "Cek Batas Rating",
			Blocking: true,
			CanFire: func() bool { return state.hasFact(derivedFactLokasiSesuai) && !state.hasFact(derivedFactRatingMemadai) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				if candidate.Rating >= fc.Config.MinimumRating {
					eval := RuleEvaluation{RuleID: "RULE_008", RuleName: "Cek Batas Rating", IsSatisfied: true, Reason: fmt.Sprintf("Rating %.2f memenuhi minimum %.2f", candidate.Rating, fc.Config.MinimumRating), Impact: 30}
					return eval, []DerivedFact{{Name: derivedFactRatingMemadai, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
				}
				eval := RuleEvaluation{RuleID: "RULE_008", RuleName: "Cek Batas Rating", IsSatisfied: false, Reason: fmt.Sprintf("Rating %.2f di bawah minimum %.2f", candidate.Rating, fc.Config.MinimumRating)}
				return eval, nil
			},
		},
		{
			ID:       "RULE_009",
			Name:     "Estimasi Jarak",
			Blocking: false,
			CanFire: func() bool { return state.hasFact(derivedFactRatingMemadai) && !state.hasFact(derivedFactJarakMemadai) },
			Fire: func() (RuleEvaluation, []DerivedFact) {
				distanceMeters, distanceReason := fc.estimateRouteDistance(snapshot, candidate, osrmCache)
				state.DistanceMeters = distanceMeters
				eval := RuleEvaluation{RuleID: "RULE_009", RuleName: "Estimasi Jarak", IsSatisfied: true, Reason: distanceReason, Impact: 0}
				if distanceMeters <= 30000 {
					return eval, []DerivedFact{{Name: derivedFactJarakMemadai, Value: true, RuleID: eval.RuleID, Reason: distanceReason, CreatedAt: time.Now()}}
				}
				return eval, []DerivedFact{{Name: derivedFactJarakMemadai, Value: false, RuleID: eval.RuleID, Reason: distanceReason, CreatedAt: time.Now()}}
			},
		},
		{
			ID:       "RULE_010",
			Name:     "Konklusi Akhir",
			Blocking: false,
			CanFire: func() bool {
				return state.hasFact(derivedFactAdminApproved) &&
					state.hasFact(derivedFactPenghuluAktif) &&
					state.hasFact(derivedFactJadwalBebas) &&
					state.hasFact(derivedFactKapasitasHarian) &&
					state.hasFact(derivedFactKapasitasPerJam) &&
					state.hasFact(derivedFactLokasiSesuai) &&
					state.hasFact(derivedFactRatingMemadai) &&
					!state.hasFact(derivedFactDirekomendasikan)
			},
			Fire: func() (RuleEvaluation, []DerivedFact) {
				eval := RuleEvaluation{RuleID: "RULE_010", RuleName: "Konklusi Akhir", IsSatisfied: true, Reason: "Semua fakta pendukung telah terpenuhi, penghulu dapat direkomendasikan", Impact: 0}
				return eval, []DerivedFact{{Name: derivedFactDirekomendasikan, Value: true, RuleID: eval.RuleID, Reason: eval.Reason, CreatedAt: time.Now()}}
			},
		},
	}
}

func (fc *ForwardChainingEngine) estimateRouteDistance(snapshot *evaluationSnapshot, candidate *PenghuluFact, osrmCache *osrmDistanceCache) (float64, string) {
	registration := snapshot.RegistrationFact
	targetLat := registration.Latitude
	targetLon := registration.Longitude
	targetFallbackUsed := false
	if targetLat == 0 && targetLon == 0 {
		targetLat = fc.Config.KuaLatitude
		targetLon = fc.Config.KuaLongitude
		targetFallbackUsed = registration.TempatNikah == structs.TempatNikahDiLuarKUA
	}

	originLat := fc.Config.KuaLatitude
	originLon := fc.Config.KuaLongitude

	type routePoint struct {
		latitude   float64
		longitude  float64
		scheduledAt time.Time
	}

	points := make([]routePoint, 0, len(snapshot.AssignmentsByPenghulu[candidate.PenghuluID])+1)
	for _, assignment := range snapshot.AssignmentsByPenghulu[candidate.PenghuluID] {
		if assignment.Latitude == nil || assignment.Longitude == nil {
			continue
		}
		points = append(points, routePoint{latitude: *assignment.Latitude, longitude: *assignment.Longitude, scheduledAt: assignment.ScheduledAt})
	}
	points = append(points, routePoint{latitude: targetLat, longitude: targetLon, scheduledAt: registration.TanggalNikah})
	sort.Slice(points, func(i, j int) bool {
		return points[i].scheduledAt.Before(points[j].scheduledAt)
	})

	totalDistance := 0.0
	currentLat := originLat
	currentLon := originLon
	for _, point := range points {
		 distance := fc.routeDistanceCached(osrmCache, currentLat, currentLon, point.latitude, point.longitude)
		totalDistance += distance
		currentLat = point.latitude
		currentLon = point.longitude
	}

	if targetFallbackUsed {
		totalDistance += fc.Config.OSRMFailurePenaltyMeters
	}

	return totalDistance, fmt.Sprintf("Estimasi rute %.2f km", totalDistance/1000.0)
}

func (fc *ForwardChainingEngine) evaluateCandidate(snapshot *evaluationSnapshot, candidate *PenghuluFact, osrmCache *osrmDistanceCache) RuleResult {
	state := &candidateInferenceState{
		Facts:      map[string]bool{},
		FiredRules: map[string]bool{},
	}
	rules := fc.buildRuleBase(snapshot, candidate, state, osrmCache)

	for iteration := 0; iteration < 10; iteration++ {
		changed := false
		for _, rule := range rules {
			if state.Rejected || state.FiredRules[rule.ID] || !rule.CanFire() {
				continue
			}
			evaluation, derivedFacts := rule.Fire()
			state.FiredRules[rule.ID] = true
			state.appendEvaluation(evaluation)

			if !evaluation.IsSatisfied && rule.Blocking {
				state.Rejected = true
				state.RejectReason = evaluation.Reason
				state.addFact(derivedFactDitolak, true, rule.ID, evaluation.Reason)
				break
			}

			for _, derivedFact := range derivedFacts {
				if state.addFact(derivedFact.Name, derivedFact.Value, derivedFact.RuleID, derivedFact.Reason) {
					changed = true
				}
			}
		}
		if state.Rejected || !changed {
			break
		}
	}

	result := RuleResult{
		PenghuluID:    candidate.PenghuluID,
		EvaluatedRules: state.Evaluations,
		DerivedFacts:   state.DerivedFacts,
		AllRulesPassed: !state.Rejected && state.hasFact(derivedFactDirekomendasikan),
		DistanceMeters: state.DistanceMeters,
	}

	if state.Rejected {
		result.Conclusion = state.RejectReason
		result.Score = 0
		return result
	}

	if state.hasFact(derivedFactDirekomendasikan) {
		result.Conclusion = "Penghulu direkomendasikan"
	} else {
		result.Conclusion = "Belum memenuhi seluruh fakta untuk rekomendasi akhir"
	}

	result.Score = fc.calculateScore(snapshot, candidate, state.DistanceMeters)
	return result
}

func (fc *ForwardChainingEngine) calculateScore(snapshot *evaluationSnapshot, candidate *PenghuluFact, distanceMeters float64) float64 {
	registration := snapshot.RegistrationFact
	dayCount := snapshot.DayCounts[candidate.PenghuluID]

	ratingScore := clampFloat(candidate.Rating/5.0, 0, 1) * 100 * fc.Config.ScoringWeights.RatingWeight

	availabilityRatio := 1.0
	if candidate.CapacityPerDay > 0 {
		availabilityRatio = 1.0 - (float64(dayCount) / float64(candidate.CapacityPerDay))
	}
	availabilityScore := clampFloat(availabilityRatio, 0, 1) * 100 * fc.Config.ScoringWeights.AvailabilityWeight

	fairnessRatio := 1.0
	if snapshot.MaxHistoricalWorkload > snapshot.MinHistoricalWorkload {
		fairnessRatio = 1.0 - (float64(candidate.JumlahNikah-snapshot.MinHistoricalWorkload) / float64(snapshot.MaxHistoricalWorkload-snapshot.MinHistoricalWorkload))
	}
	fairnessScore := clampFloat(fairnessRatio, 0, 1) * 100 * fc.Config.ScoringWeights.FairnessWeight

	locationScore := 100.0 * fc.Config.ScoringWeights.LocationMatchWeight
	if registration.TempatNikah == structs.TempatNikahDiLuarKUA {
		locationScore = 85.0 * fc.Config.ScoringWeights.LocationMatchWeight
	}

	distanceScoreBase := math.Max(0, 100-(distanceMeters/1000.0)*fc.Config.DistancePenaltyPerKm)
	distanceScore := distanceScoreBase * fc.Config.ScoringWeights.DistanceWeight

	return ratingScore + availabilityScore + fairnessScore + locationScore + distanceScore
}

// ==================== PUBLIC EVALUATION APIS ====================

// EvaluateRules dipertahankan untuk kompatibilitas; implementasi produksi memakai snapshot inference.
func (fc *ForwardChainingEngine) EvaluateRules(regFact *MarriageRegistrationFact, penghuluFact *PenghuluFact) (*RuleResult, error) {
	if regFact == nil || penghuluFact == nil {
		return nil, fmt.Errorf("facts tidak boleh kosong")
	}
	dayCounts := map[uint]int{}
	hourCounts := map[uint]map[string]int{}
	assignmentsByPenghulu := map[uint][]scheduledAssignment{}
	if err := fc.DB.Model(&structs.PendaftaranNikah{}).
		Where("penghulu_id = ? AND tanggal_nikah >= ? AND tanggal_nikah < ? AND status_pendaftaran NOT IN ?", penghuluFact.PenghuluID,
			time.Date(regFact.TanggalNikah.Year(), regFact.TanggalNikah.Month(), regFact.TanggalNikah.Day(), 0, 0, 0, 0, regFact.TanggalNikah.Location()),
			time.Date(regFact.TanggalNikah.Year(), regFact.TanggalNikah.Month(), regFact.TanggalNikah.Day(), 0, 0, 0, 0, regFact.TanggalNikah.Location()).Add(24*time.Hour),
			[]string{structs.StatusPendaftaranDitolak}).Find(&[]structs.PendaftaranNikah{}).Error; err != nil {
		return nil, err
	}
	dayCounts[penghuluFact.PenghuluID] = 0
	hourCounts[penghuluFact.PenghuluID] = map[string]int{}
	snapshot := &evaluationSnapshot{
		RegistrationFact:      regFact,
		Holiday:               false,
		HolidayName:           "",
		ActivePenghuluFacts:   []*PenghuluFact{penghuluFact},
		DayCounts:             dayCounts,
		HourCounts:            hourCounts,
		AssignmentsByPenghulu: assignmentsByPenghulu,
		MaxHistoricalWorkload: penghuluFact.JumlahNikah,
		MinHistoricalWorkload: penghuluFact.JumlahNikah,
	}
	result := fc.evaluateCandidate(snapshot, penghuluFact, newOSRMDistanceCache())
	return &result, nil
}

// GetPenghuluRecommendations mengeksekusi snapshot-based forward chaining satu kali di awal.
func (fc *ForwardChainingEngine) GetPenghuluRecommendations(registrationID uint) (*AssignmentRecommendation, error) {
	snapshot, err := fc.loadEvaluationSnapshot(registrationID, false)
	if err != nil {
		return nil, err
	}

	osrmCache := newOSRMDistanceCache()

	results := make([]RuleResult, 0, len(snapshot.ActivePenghuluFacts))
	for _, penghuluFact := range snapshot.ActivePenghuluFacts {
		results = append(results, fc.evaluateCandidate(snapshot, penghuluFact, osrmCache))
	}

	passedResults := make([]RuleResult, 0, len(results))
	for _, result := range results {
		if result.AllRulesPassed {
			passedResults = append(passedResults, result)
		}
	}
	sort.Slice(passedResults, func(i, j int) bool {
		return passedResults[i].Score > passedResults[j].Score
	})

	rec := &AssignmentRecommendation{
		EvaluatedAt:        time.Now().In(utils.WITA),
		EvaluationCount:    len(snapshot.ActivePenghuluFacts),
		EvaluationProcess:  results,
	}

	if len(passedResults) == 0 {
		rec.DecisionReasoning = "Tidak ada penghulu yang memenuhi seluruh fakta hasil inference cycle"
		rec.Confidence = 0
		return rec, nil
	}

	var bestPenghulu structs.Penghulu
	if err := fc.DB.Where("id = ?", passedResults[0].PenghuluID).First(&bestPenghulu).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil detail penghulu terbaik: %w", err)
	}

	rec.RecommendedPenghuluID = passedResults[0].PenghuluID
	rec.RecommendedPenghulu = &bestPenghulu
	rec.SelectedScore = passedResults[0].Score
	rec.Confidence = clampFloat(50+(passedResults[0].Score/2), 0, 99)
	rec.DecisionReasoning = fmt.Sprintf(
		"Forward chaining menetapkan %s sebagai kandidat terbaik dengan score %.2f, jarak %.2f km, dan seluruh fakta utama terpenuhi.",
		bestPenghulu.Nama_lengkap,
		passedResults[0].Score,
		passedResults[0].DistanceMeters/1000.0,
	)

	for i := 1; i < len(passedResults) && i < 3; i++ {
		var altPenghulu structs.Penghulu
		if err := fc.DB.Where("id = ?", passedResults[i].PenghuluID).First(&altPenghulu).Error; err != nil {
			continue
		}
		rec.Alternatives = append(rec.Alternatives, AlternativePenghulu{
			PenghuluID:   passedResults[i].PenghuluID,
			NamaPenghulu: altPenghulu.Nama_lengkap,
			Score:        passedResults[i].Score,
			Reason:       fmt.Sprintf("Score %.2f, jarak %.2f km", passedResults[i].Score, passedResults[i].DistanceMeters/1000.0),
		})
	}

	return rec, nil
}

// GetDetailedEvaluation menampilkan trace inference per penghulu secara terurut.
func (fc *ForwardChainingEngine) GetDetailedEvaluation(registrationID uint) (map[string]interface{}, error) {
	snapshot, err := fc.loadEvaluationSnapshot(registrationID, false)
	if err != nil {
		return nil, err
	}

	osrmCache := newOSRMDistanceCache()

	evaluations := make([]map[string]interface{}, 0, len(snapshot.ActivePenghuluFacts))
	for _, penghuluFact := range snapshot.ActivePenghuluFacts {
		ruleResult := fc.evaluateCandidate(snapshot, penghuluFact, osrmCache)
		evaluations = append(evaluations, map[string]interface{}{
			"penghulu_id":      penghuluFact.PenghuluID,
			"penghulu_nama":    penghuluFact.NamaPenghulu,
			"rating":           penghuluFact.Rating,
			"jumlah_nikah":     penghuluFact.JumlahNikah,
			"status":           penghuluFact.StatusAktif,
			"all_rules_passed": ruleResult.AllRulesPassed,
			"score":            ruleResult.Score,
			"distance_meters":  ruleResult.DistanceMeters,
			"conclusion":       ruleResult.Conclusion,
			"evaluated_rules":  ruleResult.EvaluatedRules,
			"derived_facts":    ruleResult.DerivedFacts,
		})
	}

	sort.Slice(evaluations, func(i, j int) bool {
		passedI := evaluations[i]["all_rules_passed"].(bool)
		passedJ := evaluations[j]["all_rules_passed"].(bool)
		if passedI != passedJ {
			return passedI
		}
		return evaluations[i]["score"].(float64) > evaluations[j]["score"].(float64)
	})

	return map[string]interface{}{
		"registration_id":     registrationID,
		"tanggal_nikah":       snapshot.RegistrationFact.TanggalNikah,
		"waktu_nikah":         snapshot.RegistrationFact.WaktuNikah,
		"tempat_nikah":        snapshot.RegistrationFact.TempatNikah,
		"alamat_nikah":        snapshot.RegistrationFact.AlamatNikah,
		"latitude":            snapshot.RegistrationFact.Latitude,
		"longitude":           snapshot.RegistrationFact.Longitude,
		"status_pendaftaran":  snapshot.RegistrationFact.StatusPendaftaran,
		"penghulu_id":         snapshot.RegistrationFact.PenghuluID,
		"holiday":             snapshot.Holiday,
		"holiday_name":        snapshot.HolidayName,
		"evaluated_at":        time.Now().In(utils.WITA),
		"penghulu_evaluations": evaluations,
	}, nil
}
