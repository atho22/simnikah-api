package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== CREATE REGISTRATION TESTS ====================

func TestCatin_CreateRegistration_DiKUA(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Seed active penghulu so FC engine has available slots
	pUser := testhelper.CreateTestUser(testDB, "peng_for_reg", "pengforreg@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	user := testhelper.CreateTestUser(testDB, "catin_reg1", "catin1@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"umur_suami":    28,
		"nama_istri":    "Siti",
		"umur_istri":    25,
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["status_pendaftaran"] != "Menunggu Penugasan" {
		t.Errorf("expected Menunggu Penugasan, got %v", data["status_pendaftaran"])
	}
}

func TestCatin_CreateRegistration_DiLuarKUA_RequiresAddress(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_reg2", "catin2@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di Luar KUA",
		// missing alamat_akad
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing address, got %d", w.Code)
	}
}

func TestCatin_CreateRegistration_InvalidTempat(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_reg3", "catin3@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": "2026-12-01",
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di Masjid", // invalid
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid tempat_nikah, got %d", w.Code)
	}
}

func TestCatin_CreateRegistration_PastDate(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_past", "past@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": "2020-01-01", // past date
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for past date, got %d", w.Code)
	}
}

func TestCatin_CreateRegistration_InvalidTimeFormat(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_time", "time@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "25:99", // invalid time
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid time, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCatin_CreateRegistration_MissingRequiredFields(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_miss", "miss@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	// Missing nama_suami, nama_istri, tanggal_nikah
	body := map[string]interface{}{
		"waktu_nikah":  "09:00",
		"tempat_nikah": "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestCatin_CreateRegistration_PendaftarIdFromJWT(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Seed active penghulu
	pUser := testhelper.CreateTestUser(testDB, "peng_jwt_reg", "pengjwtreg@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	user := testhelper.CreateTestUser(testDB, "catin_jwt", "jwt@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// Verify pendaftar_id was set from JWT
	var pendaftaran structs.PendaftaranNikah
	testDB.Last(&pendaftaran)
	if pendaftaran.Pendaftar_id != user.User_id {
		t.Errorf("SECURITY: pendaftar_id should be %s, got %s", user.User_id, pendaftaran.Pendaftar_id)
	}
}

// ==================== CHECK SCHEDULE TESTS ====================

func TestCatin_CheckSchedule_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Seed active penghulu so slots exist
	pUser := testhelper.CreateTestUser(testDB, "peng_chk_sch", "pengchksch@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	user := testhelper.CreateTestUser(testDB, "catin_chk", "chk@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	body := map[string]interface{}{
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "10:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/check-schedule", body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["available"] != true {
		t.Errorf("expected available=true on empty schedule, got %v", data["available"])
	}
}

func TestCatin_CheckSchedule_MissingDate(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "catin_chk2", "chk2@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"waktu_nikah":  "09:00",
		"tempat_nikah": "Di KUA",
	}

	w := doRequest("POST", "/simnikah/check-schedule", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing date, got %d", w.Code)
	}
}

// ==================== STAFF STATUS TRANSITION TESTS ====================

func TestStaff_UpdateStatus_MenungguToDitugaskan(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_st1", "st1@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-TEST", time.Now().AddDate(0, 1, 0))

	body := map[string]interface{}{
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/update-status", p.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStaff_UpdateStatus_InvalidTransition(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_st2", "st2@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-TEST", time.Now().AddDate(0, 1, 0))

	// Try to go directly from Menunggu Penugasan to Selesai (invalid)
	body := map[string]interface{}{
		"status_pendaftaran": structs.StatusPendaftaranSelesai,
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/update-status", p.ID), body, token)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 for invalid transition, got %d", w.Code)
	}
}

func TestStaff_UpdateStatus_InvalidStatusValue(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_st3", "st3@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-TEST", time.Now().AddDate(0, 1, 0))

	body := map[string]interface{}{
		"status_pendaftaran": "Dibatalkan", // not a valid status
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/update-status", p.ID), body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid status, got %d", w.Code)
	}
}

func TestStaff_UpdateStatus_NotFound(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_st4", "st4@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"status_pendaftaran": structs.StatusPendaftaranSelesai,
	}

	w := doRequest("PUT", "/simnikah/pendaftaran/99999/update-status", body, token)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// ==================== STAFF CRUD TESTS ====================

func TestStaff_ListStaff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_list", "kklist@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/staff", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestStaff_UpdateStaff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_upd", "kkupd@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Create a staff record
	staffRecord := structs.StaffKUA{
		User_id:      "USR-STAFF01",
		NIP:          "NIP-STAFF01",
		Nama_lengkap: "Original Name",
		Jabatan:      "Staff",
		Bagian:       "Admin",
		Status:       "Aktif",
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}
	testDB.Create(&staffRecord)

	body := map[string]interface{}{
		"nama": "Updated Name",
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/staff/%d", staffRecord.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
