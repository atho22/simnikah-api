package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== LOCATION IDOR TESTS ====================

func TestIDOR_UpdateLocation_OwnerCanUpdate(t *testing.T) {
	testhelper.CleanupDB(testDB)

	owner := testhelper.CreateTestUser(testDB, "loc_owner", "locown@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(owner.User_id, owner.Email, owner.Role, owner.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Update("tempat_nikah", "Di Luar KUA")

	lat := -3.2913
	lon := 114.5881
	body := map[string]interface{}{
		"alamat_akad": "Jl. Test No. 1, Banjarmasin",
		"latitude":    lat,
		"longitude":   lon,
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for owner, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIDOR_UpdateLocation_OtherUserBlocked(t *testing.T) {
	testhelper.CleanupDB(testDB)

	owner := testhelper.CreateTestUser(testDB, "loc_own2", "locown2@test.com", structs.UserRoleUserBiasa)
	other := testhelper.CreateTestUser(testDB, "loc_other", "locoth@test.com", structs.UserRoleUserBiasa)

	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Update("tempat_nikah", "Di Luar KUA")

	otherToken := testhelper.GenerateTestToken(other.User_id, other.Email, other.Role, other.Nama)

	body := map[string]interface{}{
		"alamat_akad": "Hacked address",
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), body, otherToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("IDOR: expected 403 for non-owner, got %d", w.Code)
	}
}

func TestIDOR_UpdateLocation_StaffCanUpdate(t *testing.T) {
	testhelper.CleanupDB(testDB)

	owner := testhelper.CreateTestUser(testDB, "loc_own3", "locown3@test.com", structs.UserRoleUserBiasa)
	staffUser := testhelper.CreateTestUser(testDB, "loc_staff", "locstaff@test.com", structs.UserRoleStaff)

	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Update("tempat_nikah", "Di Luar KUA")

	staffToken := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"alamat_akad": "Updated by staff",
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), body, staffToken)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for staff, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIDOR_GetLocationDetail_OwnerCanView(t *testing.T) {
	testhelper.CleanupDB(testDB)

	owner := testhelper.CreateTestUser(testDB, "loc_view", "locview@test.com", structs.UserRoleUserBiasa)
	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))

	token := testhelper.GenerateTestToken(owner.User_id, owner.Email, owner.Role, owner.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for owner, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIDOR_GetLocationDetail_OtherUserBlocked(t *testing.T) {
	testhelper.CleanupDB(testDB)

	owner := testhelper.CreateTestUser(testDB, "loc_vown", "locvown@test.com", structs.UserRoleUserBiasa)
	other := testhelper.CreateTestUser(testDB, "loc_voth", "locoth2@test.com", structs.UserRoleUserBiasa)

	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	otherToken := testhelper.GenerateTestToken(other.User_id, other.Email, other.Role, other.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), nil, otherToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("IDOR: expected 403 for non-owner, got %d", w.Code)
	}
}

func TestLocation_UpdateFailsForDiKUA(t *testing.T) {
	testhelper.CleanupDB(testDB)

	owner := testhelper.CreateTestUser(testDB, "loc_kua", "lockua@test.com", structs.UserRoleUserBiasa)
	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	// tempat_nikah is already "Di KUA"

	token := testhelper.GenerateTestToken(owner.User_id, owner.Email, owner.Role, owner.Nama)

	body := map[string]interface{}{
		"alamat_akad": "New address",
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for Di KUA, got %d", w.Code)
	}
}

// ==================== PENGHULU TESTS ====================

func TestPenghulu_GetJadwalPenugasan_Empty(t *testing.T) {
	testhelper.CleanupDB(testDB)

	pUser := testhelper.CreateTestUser(testDB, "peng_jdw", "pengjdw@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)
	token := testhelper.GenerateTestToken(pUser.User_id, pUser.Email, pUser.Role, pUser.Nama)

	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["total"].(float64) != 0 {
		t.Errorf("expected 0 jadwal, got %v", data["total"])
	}
}

func TestPenghulu_GetJadwalPenugasan_WithAssignment(t *testing.T) {
	testhelper.CleanupDB(testDB)

	pUser := testhelper.CreateTestUser(testDB, "peng_asgn", "pengasgn@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)
	token := testhelper.GenerateTestToken(pUser.User_id, pUser.Email, pUser.Role, pUser.Nama)

	// Create pendaftaran assigned to this penghulu
	p := testhelper.CreateTestPendaftaran(testDB, "USR-CATIN", time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Updates(map[string]interface{}{
		"penghulu_id":        penghulu.ID,
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
	})

	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Errorf("expected 1 jadwal, got %v", data["total"])
	}
}

func TestPenghulu_GetJadwal_OnlyShowsOwnAssignments(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Create two penghulu
	pUser1 := testhelper.CreateTestUser(testDB, "peng_own1", "pown1@test.com", structs.UserRolePenghulu)
	penghulu1 := testhelper.CreateTestPenghulu(testDB, pUser1.User_id)
	token1 := testhelper.GenerateTestToken(pUser1.User_id, pUser1.Email, pUser1.Role, pUser1.Nama)

	pUser2 := testhelper.CreateTestUser(testDB, "peng_own2", "pown2@test.com", structs.UserRolePenghulu)
	penghulu2 := testhelper.CreateTestPenghulu(testDB, pUser2.User_id)

	// Assign registration to penghulu1
	p1 := testhelper.CreateTestPendaftaran(testDB, "USR-CAT1", time.Now().AddDate(0, 1, 0))
	testDB.Model(&p1).Updates(map[string]interface{}{
		"penghulu_id":        penghulu1.ID,
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
	})

	// Assign registration to penghulu2
	p2 := testhelper.CreateTestPendaftaran(testDB, "USR-CAT2", time.Now().AddDate(0, 2, 0))
	testDB.Model(&p2).Updates(map[string]interface{}{
		"penghulu_id":        penghulu2.ID,
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
	})

	// Penghulu1 should only see 1 assignment
	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, token1)
	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["total"].(float64) != 1 {
		t.Errorf("SECURITY: penghulu1 should see only 1 assignment, got %v", data["total"])
	}
}

func TestPenghulu_RequiresActiveStatus(t *testing.T) {
	testhelper.CleanupDB(testDB)

	pUser := testhelper.CreateTestUser(testDB, "peng_inact", "penginact@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)
	// Deactivate user
	testDB.Model(&pUser).Update("status", "Nonaktif")

	token := testhelper.GenerateTestToken(pUser.User_id, pUser.Email, pUser.Role, pUser.Nama)

	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for inactive penghulu, got %d", w.Code)
	}
}

// ==================== KEPALA KUA TESTS ====================

func TestKepalaKUA_ListAvailablePenghulu(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_avail", "kkavail@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_avl", "pengavl@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	w := doRequest("GET", "/simnikah/kepala-kua/available-penghulu", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["total"].(float64) < 1 {
		t.Errorf("expected at least 1 penghulu, got %v", data["total"])
	}
}

func TestKepalaKUA_PenghuluTersedia(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_tersd", "kktersd@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/penghulu-tersedia?tanggal=2026-12-01", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestKepalaKUA_PenghuluTersedia_MissingDate(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_nodate", "kknodate@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/penghulu-tersedia", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing date, got %d", w.Code)
	}
}

func TestKepalaKUA_GetFCConfig(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_conf", "kkconf@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/config", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["minimum_rating"] == nil {
		t.Error("expected minimum_rating in config")
	}
	if data["scoring_weights"] == nil {
		t.Error("expected scoring_weights in config")
	}
}

func TestPenghulu_UpdateCoordinates(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Create user & corresponding penghulu
	user := testhelper.CreateTestUser(testDB, "penghulu_user", "penghulu@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, user.User_id)

	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	lat := -3.2913
	lon := 114.5881
	body := map[string]interface{}{
		"latitude":  lat,
		"longitude": lon,
	}

	w := doRequest("PUT", "/simnikah/penghulu/coordinates", body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for updating coordinates, got %d: %s", w.Code, w.Body.String())
	}

	// Verify coordinates are saved in database
	var dbPenghulu structs.Penghulu
	if err := testDB.Where("user_id = ?", user.User_id).First(&dbPenghulu).Error; err != nil {
		t.Fatalf("failed to find penghulu in database: %v", err)
	}

	if dbPenghulu.Latitude == nil || *dbPenghulu.Latitude != lat {
		t.Errorf("expected latitude %f, got %v", lat, dbPenghulu.Latitude)
	}

	if dbPenghulu.Longitude == nil || *dbPenghulu.Longitude != lon {
		t.Errorf("expected longitude %f, got %v", lon, dbPenghulu.Longitude)
	}
}

func TestPenghulu_UpdateCoordinates_ForbiddenForOtherRoles(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_user", "catin@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude":  -3.2913,
		"longitude": 114.5881,
	}

	w := doRequest("PUT", "/simnikah/penghulu/coordinates", body, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for user role catin, got %d", w.Code)
	}
}
