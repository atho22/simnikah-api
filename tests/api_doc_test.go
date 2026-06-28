package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// This test file systematically documents ALL API endpoints
// with full request/response capture for the PDF report.

// ==================== HEALTH CHECK ====================

func TestDoc_HealthCheck(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Health Check", "System")

	w := doRequest("GET", "/health", nil, "")
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// ==================== AUTHENTICATION ====================

func TestDoc_Register_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("User Registration - Success", "Authentication")

	body := map[string]interface{}{
		"username": "demo_user",
		"email":    "demo@example.com",
		"password": "password123",
		"nama":     "Demo User",
	}

	w := doRequest("POST", "/register", body, "")
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_Register_ValidationError(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("User Registration - Validation Error", "Authentication")

	body := map[string]interface{}{
		"username": "x",
	}

	w := doRequest("POST", "/register", body, "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDoc_Login_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("User Login - Success", "Authentication")

	testhelper.CreateTestUser(testDB, "login_demo", "logindemo@example.com", "user_biasa")

	body := map[string]interface{}{
		"username": "login_demo",
		"password": "password123",
	}

	w := doRequest("POST", "/login", body, "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_Login_WrongPassword(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("User Login - Wrong Password", "Authentication")

	testhelper.CreateTestUser(testDB, "login_fail", "loginfail@example.com", "user_biasa")

	body := map[string]interface{}{
		"username": "login_fail",
		"password": "wrong_password",
	}

	w := doRequest("POST", "/login", body, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestDoc_GetProfile(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get User Profile", "Authentication")

	user := testhelper.CreateTestUser(testDB, "profile_demo", "profiledemo@example.com", "user_biasa")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/profile", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_GetProfile_Unauthorized(t *testing.T) {
	testhelper.StartTest("Get Profile - No Token (401)", "Authentication")

	w := doRequest("GET", "/profile", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// ==================== PENDAFTARAN / REGISTRATION ====================

func TestDoc_CheckSchedule(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Check Schedule Availability", "Registration")

	pUser := testhelper.CreateTestUser(testDB, "doc_peng1", "docpeng1@example.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	user := testhelper.CreateTestUser(testDB, "doc_catin1", "doccatin1@example.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "10:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/check-schedule", body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_CreateRegistration_DiKUA(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Create Registration - Di KUA", "Registration")

	pUser := testhelper.CreateTestUser(testDB, "doc_peng2", "docpeng2@example.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	user := testhelper.CreateTestUser(testDB, "doc_catin2", "doccatin2@example.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad Fauzi",
		"umur_suami":    28,
		"nama_istri":    "Siti Nurhaliza",
		"umur_istri":    25,
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_CreateRegistration_DiLuarKUA(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Create Registration - Di Luar KUA", "Registration")

	user := testhelper.CreateTestUser(testDB, "doc_catin3", "doccatin3@example.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Budi Santoso",
		"umur_suami":    30,
		"nama_istri":    "Dewi Lestari",
		"umur_istri":    27,
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "10:00",
		"tempat_nikah":  "Di Luar KUA",
		"alamat_akad":   "Jl. Pramuka No. 10, Banjarmasin",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	// May get 201 or 400 depending on FC engine - both are valid for documentation
	t.Logf("Di Luar KUA registration: status %d", w.Code)
}

func TestDoc_CreateRegistration_PastDate(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Create Registration - Past Date (Rejected)", "Registration")

	user := testhelper.CreateTestUser(testDB, "doc_catin4", "doccatin4@example.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"nama_suami":    "Test",
		"nama_istri":    "Test",
		"tanggal_nikah": "2020-01-01",
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ==================== STATUS TRANSITIONS ====================

func TestDoc_UpdateStatus_MenungguToDitugaskan(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Update Status - Menunggu to Ditugaskan", "Status Transition")

	staffUser := testhelper.CreateTestUser(testDB, "doc_staff1", "docstaff1@example.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-DOC", time.Now().AddDate(0, 1, 0))

	body := map[string]interface{}{
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/update-status", p.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_UpdateStatus_InvalidTransition(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Update Status - Invalid Transition (409)", "Status Transition")

	staffUser := testhelper.CreateTestUser(testDB, "doc_staff2", "docstaff2@example.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-DOC2", time.Now().AddDate(0, 1, 0))

	body := map[string]interface{}{
		"status_pendaftaran": structs.StatusPendaftaranSelesai,
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/update-status", p.ID), body, token)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

// ==================== NOTIFICATIONS ====================

func TestDoc_Notification_Create(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Create Notification", "Notification")

	staffUser := testhelper.CreateTestUser(testDB, "doc_nstaff", "docnstaff@example.com", structs.UserRoleStaff)
	targetUser := testhelper.CreateTestUser(testDB, "doc_ntarget", "docntarget@example.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"user_id": targetUser.User_id,
		"judul":   "Pendaftaran Diterima",
		"pesan":   "Pendaftaran pernikahan Anda telah diterima dan sedang diproses.",
		"tipe":    "Info",
	}

	w := doRequest("POST", "/simnikah/notifikasi", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_Notification_GetMyNotifications(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get My Notifications", "Notification")

	user := testhelper.CreateTestUser(testDB, "doc_nuser", "docnuser@example.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")
	testhelper.CreateTestNotification(testDB, user.User_id, "Error")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/notifikasi/user/me?page=1&limit=10", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_Notification_GetByID(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get Notification By ID", "Notification")

	user := testhelper.CreateTestUser(testDB, "doc_nget", "docnget@example.com", structs.UserRoleUserBiasa)
	notif := testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/notifikasi/%d", notif.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_Notification_UpdateStatus(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Update Notification Status (Mark Read)", "Notification")

	user := testhelper.CreateTestUser(testDB, "doc_nupd", "docnupd@example.com", structs.UserRoleUserBiasa)
	notif := testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"status_baca": "Sudah Dibaca",
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/notifikasi/%d/status", notif.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_Notification_MarkAllRead(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Mark All Notifications Read", "Notification")

	user := testhelper.CreateTestUser(testDB, "doc_nmark", "docnmark@example.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("PUT", "/simnikah/notifikasi/mark-all-read", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_Notification_Stats(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get Notification Statistics", "Notification")

	user := testhelper.CreateTestUser(testDB, "doc_nstat", "docnstat@example.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")
	testhelper.CreateTestNotification(testDB, user.User_id, "Error")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/notifikasi/stats", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_Notification_Delete(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Delete Notification", "Notification")

	user := testhelper.CreateTestUser(testDB, "doc_ndel", "docndel@example.com", structs.UserRoleUserBiasa)
	notif := testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("DELETE", fmt.Sprintf("/simnikah/notifikasi/%d", notif.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_Notification_SendToRole(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Send Notification to Role", "Notification")

	staffUser := testhelper.CreateTestUser(testDB, "doc_nsender", "docnsender@example.com", structs.UserRoleStaff)
	recv := testhelper.CreateTestUser(testDB, "doc_nrecv", "docnrecv@example.com", structs.UserRolePenghulu)
	testDB.Model(&recv).Update("status", structs.UserStatusAktif)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"role":  "penghulu",
		"judul": "Rapat Koordinasi",
		"pesan": "Rapat koordinasi penghulu pada hari Senin pukul 09:00 WIB.",
		"tipe":  "Info",
	}

	w := doRequest("POST", "/simnikah/notifikasi/send-to-role", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== LOCATION ====================

func TestDoc_Location_Update(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Update Wedding Location", "Location")

	owner := testhelper.CreateTestUser(testDB, "doc_locown", "doclocown@example.com", structs.UserRoleUserBiasa)
	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Update("tempat_nikah", "Di Luar KUA")
	token := testhelper.GenerateTestToken(owner.User_id, owner.Email, owner.Role, owner.Nama)

	body := map[string]interface{}{
		"alamat_akad": "Jl. Pramuka No. 10, Banjarmasin",
		"latitude":    -3.2913,
		"longitude":   114.5881,
	}

	w := doRequest("PUT", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_Location_GetDetail(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get Wedding Location Detail", "Location")

	owner := testhelper.CreateTestUser(testDB, "doc_locget", "doclocget@example.com", structs.UserRoleUserBiasa)
	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	token := testhelper.GenerateTestToken(owner.User_id, owner.Email, owner.Role, owner.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== PENGHULU ====================

func TestDoc_Penghulu_JadwalPenugasan(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get Penghulu Schedule", "Penghulu")

	pUser := testhelper.CreateTestUser(testDB, "doc_pengj", "docpengj@example.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)
	token := testhelper.GenerateTestToken(pUser.User_id, pUser.Email, pUser.Role, pUser.Nama)

	// Create an assignment
	p := testhelper.CreateTestPendaftaran(testDB, "USR-DOCCATIN", time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Updates(map[string]interface{}{
		"penghulu_id":        penghulu.ID,
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
	})

	w := doRequest("GET", "/simnikah/penghulu/jadwal-penugasan", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== KEPALA KUA ====================

func TestDoc_KepalaKUA_AvailablePenghulu(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("List Available Penghulu", "Kepala KUA")

	kkUser := testhelper.CreateTestUser(testDB, "doc_kk1", "dockk1@example.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "doc_pengav", "docpengav@example.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	w := doRequest("GET", "/simnikah/kepala-kua/available-penghulu", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_KepalaKUA_PenghuluTersedia(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get Penghulu Tersedia", "Kepala KUA")

	kkUser := testhelper.CreateTestUser(testDB, "doc_kk2", "dockk2@example.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/penghulu-tersedia?tanggal=2026-12-01", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_KepalaKUA_FCConfig(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Get Forward Chaining Config", "Kepala KUA")

	kkUser := testhelper.CreateTestUser(testDB, "doc_kk3", "dockk3@example.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/config", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_KepalaKUA_Dashboard(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Kepala KUA Dashboard", "Kepala KUA")

	kkUser := testhelper.CreateTestUser(testDB, "doc_kk4", "dockk4@example.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/kepala-kua", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== STAFF ====================

func TestDoc_Staff_List(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("List Staff", "Staff")

	kkUser := testhelper.CreateTestUser(testDB, "doc_kk5", "dockk5@example.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/staff", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDoc_Staff_Dashboard(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Staff Dashboard", "Staff")

	staffUser := testhelper.CreateTestUser(testDB, "doc_staffd", "docstaffd@example.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/staff", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ==================== SECURITY ====================

func TestDoc_RBAC_ForbiddenAccess(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("RBAC - User Biasa Access Kepala KUA (403)", "Security")

	user := testhelper.CreateTestUser(testDB, "doc_rbac", "docrbac@example.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/available-penghulu", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDoc_IDOR_CrossUserAccess(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("IDOR Prevention - Cross User Notification Access", "Security")

	userA := testhelper.CreateTestUser(testDB, "doc_idorA", "docidora@example.com", structs.UserRoleUserBiasa)
	userB := testhelper.CreateTestUser(testDB, "doc_idorB", "docidorb@example.com", structs.UserRoleUserBiasa)
	notifA := testhelper.CreateTestNotification(testDB, userA.User_id, "Info")
	tokenB := testhelper.GenerateTestToken(userB.User_id, userB.Email, userB.Role, userB.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/notifikasi/%d", notifA.ID), nil, tokenB)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 (IDOR blocked), got %d", w.Code)
	}
}

func TestDoc_IDOR_LocationAccess(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("IDOR Prevention - Cross User Location Access", "Security")

	owner := testhelper.CreateTestUser(testDB, "doc_locsec", "doclocsec@example.com", structs.UserRoleUserBiasa)
	other := testhelper.CreateTestUser(testDB, "doc_locoth", "doclocoth@example.com", structs.UserRoleUserBiasa)
	p := testhelper.CreateTestPendaftaran(testDB, owner.User_id, time.Now().AddDate(0, 1, 0))
	otherToken := testhelper.GenerateTestToken(other.User_id, other.Email, other.Role, other.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/pendaftaran/%d/location", p.ID), nil, otherToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 (IDOR blocked), got %d", w.Code)
	}
}

func TestDoc_DashboardStats(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Marriage Statistics Dashboard", "Dashboard")

	staffUser := testhelper.CreateTestUser(testDB, "doc_dstaff", "docdstaff@example.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/statistik-pernikahan", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDoc_PeakHours(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.StartTest("Peak Hours Analysis", "Dashboard")

	staffUser := testhelper.CreateTestUser(testDB, "doc_phstaff", "docphstaff@example.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/dashboard/peak-hours", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
