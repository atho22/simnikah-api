package tests

import (
	"fmt"
	"net/http"
	"testing"

	structs "simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== STAFF CREATE REGISTRATION FOR USER ====================

// TestStaffReg_Success memastikan staff dapat membuat pendaftaran untuk catin.
func TestStaffReg_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_creg_ok", "staffcregok@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	// Perlu penghulu aktif agar FC engine menyetujui jadwal
	pUser := testhelper.CreateTestUser(testDB, "peng_creg", "pengcreg@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad Pendaftar",
		"nama_istri":    "Siti Binti Hasan",
		"umur_suami":    28,
		"umur_istri":    25,
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// TestStaffReg_KepalaKUACanCreate memastikan kepala_kua juga bisa buat pendaftaran untuk catin.
func TestStaffReg_KepalaKUACanCreate(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_creg_ok", "kkcregok@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_kk_creg", "pengkkcreg@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Budi Santoso",
		"nama_istri":    "Dewi Lestari",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "10:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 for kepala_kua, got %d: %s", w.Code, w.Body.String())
	}
}

// TestStaffReg_RBAC_UserBiasaRejected memastikan user_biasa tidak bisa pakai endpoint staff.
func TestStaffReg_RBAC_UserBiasaRejected(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_creg_rbac", "usercregrbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for user_biasa, got %d", w.Code)
	}
}

// TestStaffReg_RBAC_PenghuluRejected memastikan penghulu tidak bisa pakai endpoint staff.
func TestStaffReg_RBAC_PenghuluRejected(t *testing.T) {
	testhelper.CleanupDB(testDB)

	pUser := testhelper.CreateTestUser(testDB, "peng_creg_rbac", "pengcregrbac@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)
	token := testhelper.GenerateTestToken(pUser.User_id, pUser.Email, pUser.Role, pUser.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for penghulu, got %d", w.Code)
	}
}

// TestStaffReg_MissingNamaSuami memastikan 400 jika nama_suami tidak dikirim.
func TestStaffReg_MissingNamaSuami(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_miss_sm", "staffmisssm@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		// nama_suami tidak ada
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing nama_suami, got %d", w.Code)
	}
}

// TestStaffReg_MissingNamaIstri memastikan 400 jika nama_istri tidak dikirim.
func TestStaffReg_MissingNamaIstri(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_miss_ist", "staffmissist@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami": "Ahmad",
		// nama_istri tidak ada
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing nama_istri, got %d", w.Code)
	}
}

// TestStaffReg_InvalidTempatNikah memastikan 400 untuk tempat_nikah tidak valid.
func TestStaffReg_InvalidTempatNikah(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_inv_tmp", "staffinvtmp@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di Mushola", // tidak valid
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid tempat_nikah, got %d", w.Code)
	}
}

// TestStaffReg_InvalidDateFormat memastikan 400 untuk tanggal_nikah format salah.
func TestStaffReg_InvalidDateFormat(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_inv_dt", "staffinvdt@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": "01-12-2026", // format salah, harusnya YYYY-MM-DD
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid date format, got %d", w.Code)
	}
}

// TestStaffReg_DiLuarKUA_WithAddress memastikan pendaftaran di luar KUA berhasil dengan alamat.
func TestStaffReg_DiLuarKUA_WithAddress(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_luar_ok", "staffluarok@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_luar", "pengluar@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di Luar KUA",
		"alamat_akad":   "Jl. Merdeka No. 10, Banjarmasin",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	// Bisa 201 (berhasil) atau 409 (jadwal penuh) — keduanya valid
	if w.Code != http.StatusCreated && w.Code != http.StatusConflict {
		t.Errorf("expected 201 or 409, got %d: %s", w.Code, w.Body.String())
	}
}

// TestStaffReg_WithoutToken memastikan 401 tanpa token.
func TestStaffReg_WithoutToken(t *testing.T) {
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": "2026-12-01",
		"waktu_nikah":   "09:00",
		"tempat_nikah":  "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

// TestStaffReg_MissingWaktuNikah memastikan 400 jika waktu_nikah tidak dikirim.
func TestStaffReg_MissingWaktuNikah(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_miss_wkt", "staffmisswkt@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	futureDate := testhelper.GetSafeFutureDate().Format("2006-01-02")
	body := map[string]interface{}{
		"nama_suami":    "Ahmad",
		"nama_istri":    "Siti",
		"tanggal_nikah": futureDate,
		// waktu_nikah tidak ada
		"tempat_nikah": "Di KUA",
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing waktu_nikah, got %d", w.Code)
	}
}

// TestStaffReg_MissingTanggalNikah memastikan 400 jika tanggal_nikah tidak dikirim.
func TestStaffReg_MissingTanggalNikah(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_miss_tgl", "staffmisstgl@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"nama_suami":   "Ahmad",
		"nama_istri":   "Siti",
		"waktu_nikah":  "09:00",
		"tempat_nikah": "Di KUA",
		// tanggal_nikah tidak ada
	}

	w := doRequest("POST", "/simnikah/staff/pendaftaran", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing tanggal_nikah, got %d", w.Code)
	}
}

// TestStaffReg_UpdateStaff_NotFound memastikan 404 saat update staff yang tidak ada.
func TestStaffReg_UpdateStaff_NotFound(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_upd_nf", "kkupdnf@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	body := map[string]interface{}{
		"nama": "Updated Name",
	}

	w := doRequest("PUT", "/simnikah/staff/99999", body, token)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent staff, got %d", w.Code)
	}
}

// TestStaffReg_UpdateStaff_AllFields memastikan semua field staff bisa diupdate sekaligus.
func TestStaffReg_UpdateStaff_AllFields(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_upd_all", "kkupdall@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Create a staff record
	staffRecord := structs.StaffKUA{
		User_id:      "USR-STAFFALL01",
		NIP:          "NIP-STAFFALL01",
		Nama_lengkap: "Original Name",
		Jabatan:      "Staff",
		Bagian:       "Admin",
		Status:       "Aktif",
	}
	testDB.Create(&staffRecord)

	body := map[string]interface{}{
		"nama":    "Updated Nama",
		"jabatan": "Kepala Bagian",
		"bagian":  "Pelayanan",
		"no_hp":   "081234567890",
		"email":   "updated@kua.go.id",
		"alamat":  "Jl. Test No. 1",
		"status":  "Aktif",
	}

	w := doRequest("PUT", "/simnikah/staff/"+itoa(int(staffRecord.ID)), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// itoa adalah helper untuk convert int ke string URL parameter
func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
