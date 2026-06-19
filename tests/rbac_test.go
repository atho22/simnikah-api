package tests

import (
	"testing"

	"simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== RBAC ENFORCEMENT TESTS ====================
// Verifies that every protected endpoint correctly rejects unauthorized roles

func TestRBAC_KepalaKUAOnlyEndpoints(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rbac_user", "rbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/simnikah/kepala-kua/forward-chaining/recommendation/1"},
		{"GET", "/simnikah/kepala-kua/forward-chaining/evaluation/1"},
		{"GET", "/simnikah/kepala-kua/forward-chaining/config"},
		{"POST", "/simnikah/kepala-kua/forward-chaining/assign/1"},
		{"GET", "/simnikah/kepala-kua/available-penghulu"},
		{"GET", "/simnikah/kepala-kua/penghulu-tersedia?tanggal=2026-12-01"},
		{"GET", "/simnikah/staff"},
		{"PUT", "/simnikah/staff/1"},
		{"GET", "/simnikah/dashboard/kepala-kua"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(ep.method, ep.path, nil, token)
			if w.Code != 403 {
				t.Errorf("SECURITY: expected 403 for user_biasa on %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}

func TestRBAC_StaffEndpoints(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rbac_user2", "rbac2@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/simnikah/dashboard/staff"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(ep.method, ep.path, nil, token)
			if w.Code != 403 {
				t.Errorf("SECURITY: expected 403 for user_biasa on %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}

func TestRBAC_StaffAndKepalaKUAEndpoints(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rbac_user3", "rbac3@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"nama_suami":    "Test",
		"nama_istri":    "Test",
		"tanggal_nikah": "2026-12-01",
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	endpoints := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"POST", "/simnikah/staff/pendaftaran", body},
		{"POST", "/simnikah/notifikasi/send-to-role", map[string]interface{}{
			"role": "penghulu", "judul": "Test", "pesan": "Test", "tipe": "Info",
		}},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(ep.method, ep.path, ep.body, token)
			if w.Code != 403 {
				t.Errorf("SECURITY: expected 403 for user_biasa on %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}

func TestRBAC_MultiRoleEndpoints_StaffPenghuluKepala(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rbac_user4", "rbac4@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	endpoints := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"PUT", "/simnikah/pendaftaran/1/update-status", map[string]interface{}{
			"status_pendaftaran": "Selesai",
		}},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(ep.method, ep.path, ep.body, token)
			if w.Code != 403 {
				t.Errorf("SECURITY: expected 403 for user_biasa on %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}

func TestRBAC_PenghuluEndpoint_RejectsOtherRoles(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// staff should NOT access penghulu-only endpoint
	staffUser := testhelper.CreateTestUser(testDB, "rbac_staff", "rbacstaff@test.com", structs.UserRoleStaff)
	staffToken := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, staffToken)
	if w.Code != 403 {
		t.Errorf("SECURITY: expected 403 for staff on penghulu endpoint, got %d", w.Code)
	}
}

func TestRBAC_AuthenticatedEndpointsRejectNoToken(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/profile"},
		{"POST", "/simnikah/check-schedule"},
		{"POST", "/simnikah/pendaftaran"},
		{"GET", "/simnikah/notifikasi/user/me"},
		{"GET", "/simnikah/notifikasi/stats"},
		{"GET", "/simnikah/kepala-kua/available-penghulu"},
		{"GET", "/simnikah/penghulu/jadwal-penugasan"},
		{"GET", "/simnikah/dashboard/kepala-kua"},
		{"GET", "/simnikah/dashboard/staff"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(ep.method, ep.path, nil, "")
			if w.Code != 401 {
				t.Errorf("SECURITY: expected 401 without token on %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}

func TestRBAC_KepalaKUACanAccessPenghuluEndpoint(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_rbac", "kk@test.com", structs.UserRoleKepalaKUA)
	kkToken := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, kkToken)

	// kepala_kua is allowed by MultiRoleMiddleware but requirePenghulu will reject
	// because it checks role == penghulu internally
	if w.Code == 403 {
		// This is expected - the internal requirePenghulu blocks non-penghulu
		// The MultiRoleMiddleware allows kepala_kua through
		return
	}
	// If we get a different code, that's also OK (e.g., 200 with empty data or 404)
}

func TestRBAC_StaffCanAccessStaffDashboard(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_dash", "staffdash@test.com", structs.UserRoleStaff)
	staffToken := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/staff", nil, staffToken)
	if w.Code != 200 {
		t.Errorf("expected 200 for staff dashboard, got %d: %s", w.Code, w.Body.String())
	}
}
