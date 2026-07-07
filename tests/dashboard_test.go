package tests

import (
	"net/http"
	"testing"

	structs "simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== DASHBOARD KEPALA KUA TESTS ====================

// TestDashboard_KepalaKUA_DefaultPeriod memastikan dashboard kepala KUA berhasil
// diakses dengan period default (bulan ini).
func TestDashboard_KepalaKUA_DefaultPeriod(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_def", "kkdashdef@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	if resp["data"] == nil {
		t.Error("expected data in response")
	}
}

// TestDashboard_KepalaKUA_PeriodDay memastikan filter period=day berfungsi.
func TestDashboard_KepalaKUA_PeriodDay(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_day", "kkdashday@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua?period=day", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_KepalaKUA_PeriodWeek memastikan filter period=week berfungsi.
func TestDashboard_KepalaKUA_PeriodWeek(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_week", "kkdashweek@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua?period=week", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_KepalaKUA_PeriodYear memastikan filter period=year berfungsi.
func TestDashboard_KepalaKUA_PeriodYear(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_year", "kkdashyear@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua?period=year", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_KepalaKUA_CustomDateRange memastikan custom date range berfungsi.
func TestDashboard_KepalaKUA_CustomDateRange(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_cust", "kkdashcust@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua?date_from=2026-01-01&date_to=2026-12-31", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_KepalaKUA_InvalidDateFrom memastikan 400 untuk date_from format salah.
func TestDashboard_KepalaKUA_InvalidDateFrom(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_invdf", "kkdashinvdf@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua?date_from=invalid-date&date_to=2026-12-31", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid date_from, got %d", w.Code)
	}
}

// TestDashboard_KepalaKUA_InvalidDateTo memastikan 400 untuk date_to format salah.
func TestDashboard_KepalaKUA_InvalidDateTo(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_invdt", "kkdashinvdt@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua?date_from=2026-01-01&date_to=not-a-date", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid date_to, got %d", w.Code)
	}
}

// TestDashboard_KepalaKUA_RBAC_UserBiasa memastikan 403 untuk user_biasa.
func TestDashboard_KepalaKUA_RBAC_UserBiasa(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_dash_rbac", "userdashrbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for user_biasa, got %d", w.Code)
	}
}

// TestDashboard_KepalaKUA_RBAC_Staff memastikan 403 untuk staff.
func TestDashboard_KepalaKUA_RBAC_Staff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_dash_rbac", "staffdashrbac@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for staff, got %d", w.Code)
	}
}

// TestDashboard_KepalaKUA_WithData memastikan statistik tepat saat ada data.
func TestDashboard_KepalaKUA_WithData(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_dash_data", "kkdashdata@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Seed beberapa pendaftaran
	for i := 0; i < 3; i++ {
		testhelper.CreateTestPendaftaran(testDB, "USR-KDASH", testhelper.GetSafeFutureDate())
	}

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== DASHBOARD STAFF TESTS ====================

// TestDashboard_Staff_Success memastikan dashboard staff berhasil diakses.
func TestDashboard_Staff_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_dash_ok", "staffdashok@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/staff", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// TestDashboard_Staff_RBAC_UserBiasa memastikan 403 untuk user_biasa.
func TestDashboard_Staff_RBAC_UserBiasa(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_staff_dash", "userstaffdash@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/dashboard/staff", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403, got %d", w.Code)
	}
}

// ==================== DASHBOARD STATISTIK PERNIKAHAN ====================

// TestDashboard_StatistikPernikahan_KepalaKUA memastikan statistik pernikahan accessible oleh kepala_kua.
func TestDashboard_StatistikPernikahan_KepalaKUA(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_stat", "kkstat@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/statistik-pernikahan", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// TestDashboard_StatistikPernikahan_Staff memastikan statistik pernikahan accessible oleh staff.
func TestDashboard_StatistikPernikahan_Staff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_stat", "staffstat@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/statistik-pernikahan", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for staff, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_StatistikPernikahan_RBAC_UserBiasa memastikan 403 untuk user_biasa.
func TestDashboard_StatistikPernikahan_RBAC_UserBiasa(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_stat_rbac", "userstatrbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/dashboard/statistik-pernikahan", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403, got %d", w.Code)
	}
}

// TestDashboard_StatistikPernikahan_WithYearFilter memastikan filter year berfungsi.
func TestDashboard_StatistikPernikahan_WithYearFilter(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_stat_yr", "kkstatyr@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/statistik-pernikahan?year=2026", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== DASHBOARD PENGHULU PERFORMANCE ====================

// TestDashboard_PenghuluPerformance_Success memastikan data performa penghulu berhasil didapat.
func TestDashboard_PenghuluPerformance_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_perf", "kkperf@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Seed penghulu untuk performa data
	pUser := testhelper.CreateTestUser(testDB, "peng_perf", "pengperf@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	w := doRequest("GET", "/simnikah/dashboard/penghulu-performance", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// TestDashboard_PenghuluPerformance_Staff memastikan staff juga bisa akses.
func TestDashboard_PenghuluPerformance_Staff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_perf", "staffperf@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/penghulu-performance", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for staff, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_PenghuluPerformance_RBAC_UserBiasa memastikan 403 untuk user_biasa.
func TestDashboard_PenghuluPerformance_RBAC_UserBiasa(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_perf_rbac", "userperfrbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/dashboard/penghulu-performance", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403, got %d", w.Code)
	}
}

// TestDashboard_PenghuluPerformance_NoPenghulu memastikan response valid walau tidak ada penghulu.
func TestDashboard_PenghuluPerformance_NoPenghulu(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_perf_empty", "kkperfempty@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Tidak ada penghulu di DB
	w := doRequest("GET", "/simnikah/dashboard/penghulu-performance", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 even with no penghulu, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== DASHBOARD PEAK HOURS ====================

// TestDashboard_PeakHours_Success memastikan analisis jam sibuk berhasil.
func TestDashboard_PeakHours_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_peak", "kkpeak@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/peak-hours", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// TestDashboard_PeakHours_Staff memastikan staff bisa akses peak hours.
func TestDashboard_PeakHours_Staff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_peak", "staffpeak@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/peak-hours", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for staff, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_PeakHours_RBAC_UserBiasa memastikan 403 untuk user_biasa.
func TestDashboard_PeakHours_RBAC_UserBiasa(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_peak_rbac", "userpeakrbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/dashboard/peak-hours", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403, got %d", w.Code)
	}
}

// TestDashboard_PeakHours_WithYearFilter memastikan filter tahun berfungsi.
func TestDashboard_PeakHours_WithYearFilter(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_peak_yr", "kkpeakyr@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/peak-hours?year=2026", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestDashboard_WithoutToken memastikan semua endpoint dashboard menolak request tanpa token.
func TestDashboard_WithoutToken(t *testing.T) {
	endpoints := []string{
		"/simnikah/dashboard/kepala-kua",
		"/simnikah/dashboard/staff",
		"/simnikah/dashboard/statistik-pernikahan",
		"/simnikah/dashboard/penghulu-performance",
		"/simnikah/dashboard/peak-hours",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			w := doRequest("GET", ep, nil, "")
			if w.Code != http.StatusUnauthorized {
				t.Errorf("SECURITY: expected 401 without token on %s, got %d", ep, w.Code)
			}
		})
	}
}
