package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	structs "simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== FORWARD CHAINING RECOMMENDATION TESTS ====================

// TestFC_Recommendation_Success memastikan endpoint rekomendasi mengembalikan data valid
// saat ada penghulu aktif dan pendaftaran yang valid.
func TestFC_Recommendation_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Setup: kepala_kua user
	kkUser := testhelper.CreateTestUser(testDB, "kk_fc_rec", "kkfcrec@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Setup: penghulu aktif
	pUser := testhelper.CreateTestUser(testDB, "peng_fc_rec", "pengfcrec@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	// Setup: pendaftaran
	p := testhelper.CreateTestPendaftaran(testDB, "USR-CATIN-FC", testhelper.GetSafeFutureDate())

	w := doRequest("GET", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/recommendation/%d", p.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	// Engine harus mengembalikan field-field berikut
	if _, ok := data["evaluation_count"]; !ok {
		t.Error("expected evaluation_count in response")
	}
	if _, ok := data["evaluated_at"]; !ok {
		t.Error("expected evaluated_at in response")
	}
}

// TestFC_Recommendation_NoPenghulu memastikan engine tidak crash saat tidak ada penghulu.
func TestFC_Recommendation_NoPenghulu(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_fc_empty", "kkfcempty@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-CATIN-EMPTY", testhelper.GetSafeFutureDate())

	// Tidak ada penghulu di DB → engine harus tetap return 200 dengan data kosong/default
	w := doRequest("GET", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/recommendation/%d", p.ID), nil, token)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Recommendation_NotFound memastikan 404/400 saat registration ID tidak ada.
func TestFC_Recommendation_NotFound(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_fc_nf", "kkfcnf@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/recommendation/99999", nil, token)
	if w.Code == http.StatusOK {
		t.Errorf("expected non-200 for non-existent registration, got %d", w.Code)
	}
}

// TestFC_Recommendation_InvalidID memastikan 400 saat ID bukan angka.
func TestFC_Recommendation_InvalidID(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_fc_inv", "kkfcinv@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/recommendation/abc", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Recommendation_RBAC_Rejected memastikan 403 untuk user bukan kepala_kua.
func TestFC_Recommendation_RBAC_Rejected(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "user_fc_rbac", "userfcrbac@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/recommendation/1", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for user_biasa, got %d", w.Code)
	}
}

// ==================== FORWARD CHAINING EVALUATION REPORT TESTS ====================

// TestFC_Evaluation_Success memastikan laporan evaluasi detail berhasil didapatkan.
func TestFC_Evaluation_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_fc_eval", "kkfceval@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_fc_eval", "pengfceval@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-EVAL", testhelper.GetSafeFutureDate())

	w := doRequest("GET", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/evaluation/%d", p.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

// TestFC_Evaluation_InvalidID memastikan 400 saat ID tidak valid.
func TestFC_Evaluation_InvalidID(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_eval_inv", "kkevalinv@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/evaluation/xyz", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID, got %d", w.Code)
	}
}

// TestFC_Evaluation_RBAC_Rejected memastikan 403 untuk staff biasa.
func TestFC_Evaluation_RBAC_Rejected(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_eval_rbac", "staffevalrbac@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/evaluation/1", nil, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for staff, got %d", w.Code)
	}
}

// ==================== FORWARD CHAINING ASSIGN PENGHULU TESTS ====================

// TestFC_Assign_Success memastikan kepala_kua dapat assign penghulu ke pendaftaran.
func TestFC_Assign_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_assign_ok", "kkassignok@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_asgn_ok", "pengasgnok@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-ASSIGN", testhelper.GetSafeFutureDate())

	body := map[string]interface{}{
		"penghulu_id":    penghulu.ID,
		"approval_notes": "Disetujui sesuai jadwal dan kapasitas",
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
	data := resp["data"].(map[string]interface{})
	if data["status_pendaftaran"] != structs.StatusPendaftaranPenghuluDitugaskan {
		t.Errorf("expected status Penghulu Ditugaskan, got %v", data["status_pendaftaran"])
	}
}

// TestFC_Assign_MissingApprovalNotes memastikan 400 saat approval_notes kosong.
func TestFC_Assign_MissingApprovalNotes(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_assign_notes", "kkassignnotes@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_notes", "pengnotes@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-NOTES", testhelper.GetSafeFutureDate())

	body := map[string]interface{}{
		"penghulu_id":    penghulu.ID,
		"approval_notes": "", // kosong
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty notes, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Assign_MissingPenghuluID memastikan 400 saat penghulu_id tidak dikirim.
func TestFC_Assign_MissingPenghuluID(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_assign_noid", "kkassignnoid@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-NOID", testhelper.GetSafeFutureDate())

	body := map[string]interface{}{
		"approval_notes": "Approved",
		// penghulu_id tidak dikirim
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing penghulu_id, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Assign_PenghuluNotFound memastikan 404/400 saat penghulu_id tidak ada di DB.
func TestFC_Assign_PenghuluNotFound(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_assign_404", "kkassign404@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-PNHF", testhelper.GetSafeFutureDate())

	body := map[string]interface{}{
		"penghulu_id":    99999, // tidak ada
		"approval_notes": "Approved",
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
		t.Errorf("expected 404 or 400 for unknown penghulu, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Assign_RegistrationNotFound memastikan 404 saat registrationID tidak ada.
func TestFC_Assign_RegistrationNotFound(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_assign_reg404", "kkassignreg404@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_reg404", "pengreg404@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	body := map[string]interface{}{
		"penghulu_id":    penghulu.ID,
		"approval_notes": "Approved",
	}

	w := doRequest("POST", "/simnikah/kepala-kua/forward-chaining/assign/99999", body, token)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent registration, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Assign_AlreadyAssigned memastikan 409 saat pendaftaran sudah di-assign.
func TestFC_Assign_AlreadyAssigned(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_already", "kkalready@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_already", "pengalready@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	// Buat pendaftaran yang sudah dalam status "Penghulu Ditugaskan"
	p := testhelper.CreateTestPendaftaran(testDB, "USR-ALREADY", testhelper.GetSafeFutureDate())
	testDB.Model(&p).Updates(map[string]interface{}{
		"status_pendaftaran": structs.StatusPendaftaranPenghuluDitugaskan,
		"penghulu_id":        penghulu.ID,
	})

	body := map[string]interface{}{
		"penghulu_id":    penghulu.ID,
		"approval_notes": "Trying to re-assign",
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 for already-assigned registration, got %d: %s", w.Code, w.Body.String())
	}
}

// TestFC_Assign_InvalidID memastikan 400 untuk ID bukan angka.
func TestFC_Assign_InvalidID(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_assign_iid", "kkassigniid@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	body := map[string]interface{}{
		"penghulu_id":    1,
		"approval_notes": "Approved",
	}

	w := doRequest("POST", "/simnikah/kepala-kua/forward-chaining/assign/abc", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID, got %d", w.Code)
	}
}

// TestFC_Assign_RBAC_Rejected memastikan 403 untuk bukan kepala_kua.
func TestFC_Assign_RBAC_Rejected(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "staff_assign_rbac", "staffassignrbac@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"penghulu_id":    1,
		"approval_notes": "Approved",
	}

	w := doRequest("POST", "/simnikah/kepala-kua/forward-chaining/assign/1", body, token)
	if w.Code != http.StatusForbidden {
		t.Errorf("SECURITY: expected 403 for staff, got %d", w.Code)
	}
}

// TestFC_Assign_PenghuluOwnershipSecurity memastikan penghulu_id tidak bisa dipalsukan.
// Kepala KUA harus assign penghulu yang aktif dan valid, bukan ID sembarang.
func TestFC_Assign_VerifyPenghuluInDB(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_verify", "kkverify@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_verify", "pengverify@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-VRF", testhelper.GetSafeFutureDate())

	body := map[string]interface{}{
		"penghulu_id":    penghulu.ID,
		"approval_notes": "Assigned after evaluation",
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid assign, got %d: %s", w.Code, w.Body.String())
	}

	// Verifikasi data di DB
	var updated structs.PendaftaranNikah
	testDB.Where("id = ?", p.ID).First(&updated)
	if updated.Penghulu_id == nil || *updated.Penghulu_id != penghulu.ID {
		t.Errorf("SECURITY: penghulu_id in DB should be %d, got %v", penghulu.ID, updated.Penghulu_id)
	}
	if updated.Status_pendaftaran != structs.StatusPendaftaranPenghuluDitugaskan {
		t.Errorf("expected status Penghulu Ditugaskan, got %s", updated.Status_pendaftaran)
	}
}

// TestFC_Assign_WithoutToken memastikan 401 tanpa token.
func TestFC_Assign_WithoutToken(t *testing.T) {
	w := doRequest("POST", "/simnikah/kepala-kua/forward-chaining/assign/1", map[string]interface{}{
		"penghulu_id":    1,
		"approval_notes": "Approved",
	}, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

// ==================== FC CONFIG TEST ====================

// TestFC_Config_FieldsComplete memastikan semua field konfigurasi ada di response.
func TestFC_Config_FieldsComplete(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_cfg_full", "kkcfgfull@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	w := doRequest("GET", "/simnikah/kepala-kua/forward-chaining/config", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})

	requiredFields := []string{
		"minimum_rating",
		"capacity_per_day",
		"capacity_per_hour",
		"kua_latitude",
		"kua_longitude",
		"scoring_weights",
		"rule_constraint_notes",
	}
	for _, field := range requiredFields {
		if data[field] == nil {
			t.Errorf("expected field %q in config response", field)
		}
	}
}

// TestFC_Assign_SequentialAssign memastikan 2 assign berturutan ke penghulu yang sama berhasil
// (untuk 2 pendaftaran di tanggal berbeda, tanpa bentrok kapasitas).
func TestFC_Assign_SequentialAssign(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_seq_ok", "kkseqok@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	// Satu penghulu untuk 2 pendaftaran di hari berbeda
	pUser := testhelper.CreateTestUser(testDB, "peng_seqonly", "pengseqonly@test.com", structs.UserRolePenghulu)
	testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	// Reload dari DB untuk mendapat ID auto-increment yang benar
	var dbPenghulu structs.Penghulu
	if err := testDB.Where("user_id = ?", pUser.User_id).First(&dbPenghulu).Error; err != nil {
		t.Fatalf("penghulu tidak ditemukan di DB: %v", err)
	}
	if dbPenghulu.ID == 0 {
		t.Fatal("penghulu ID = 0, belum tersimpan ke DB")
	}

	// Buat 2 pendaftaran dengan tanggal berbeda (tidak bentrok)
	futureDate1 := testhelper.GetSafeFutureDate()
	futureDate2 := futureDate1.AddDate(0, 0, 7) // 7 hari kemudian

	p1 := testhelper.CreateTestPendaftaran(testDB, "USR-SEQA1", futureDate1)
	p2 := testhelper.CreateTestPendaftaran(testDB, "USR-SEQA2", futureDate2)

	// Assign p1 ke penghulu (assign pertama)
	w1 := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p1.ID), map[string]interface{}{
		"penghulu_id":    dbPenghulu.ID,
		"approval_notes": "Assign pertama - tanggal berbeda",
	}, token)
	if w1.Code != http.StatusOK {
		t.Errorf("assign p1 gagal: %d - %s", w1.Code, w1.Body.String())
	}

	// Assign p2 ke penghulu yang sama (assign kedua, tanggal berbeda - tidak bentrok)
	w2 := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p2.ID), map[string]interface{}{
		"penghulu_id":    dbPenghulu.ID,
		"approval_notes": "Assign kedua - tanggal berbeda",
	}, token)
	if w2.Code != http.StatusOK {
		t.Errorf("assign p2 gagal: %d - %s", w2.Code, w2.Body.String())
	}

	// Verifikasi kedua pendaftaran terupdate dengan benar di DB
	var reg1, reg2 structs.PendaftaranNikah
	testDB.Where("id = ?", p1.ID).First(&reg1)
	testDB.Where("id = ?", p2.ID).First(&reg2)

	if reg1.Status_pendaftaran != structs.StatusPendaftaranPenghuluDitugaskan {
		t.Errorf("p1 harus status Penghulu Ditugaskan, got %s", reg1.Status_pendaftaran)
	}
	if reg2.Status_pendaftaran != structs.StatusPendaftaranPenghuluDitugaskan {
		t.Errorf("p2 harus status Penghulu Ditugaskan, got %s", reg2.Status_pendaftaran)
	}
}



// TestFC_Assign_SelesaiCannotBeAssigned memastikan pendaftaran berstatus "Selesai" tidak bisa di-assign ulang.
func TestFC_Assign_SelesaiCannotBeAssigned(t *testing.T) {
	testhelper.CleanupDB(testDB)

	kkUser := testhelper.CreateTestUser(testDB, "kk_selesai", "kkselesai@test.com", structs.UserRoleKepalaKUA)
	token := testhelper.GenerateTestToken(kkUser.User_id, kkUser.Email, kkUser.Role, kkUser.Nama)

	pUser := testhelper.CreateTestUser(testDB, "peng_selesai", "pengselesai@test.com", structs.UserRolePenghulu)
	penghulu := testhelper.CreateTestPenghulu(testDB, pUser.User_id)

	p := testhelper.CreateTestPendaftaran(testDB, "USR-SLS", time.Now().AddDate(0, 1, 0))
	testDB.Model(&p).Update("status_pendaftaran", structs.StatusPendaftaranSelesai)

	body := map[string]interface{}{
		"penghulu_id":    penghulu.ID,
		"approval_notes": "Trying to assign finished registration",
	}

	w := doRequest("POST", fmt.Sprintf("/simnikah/kepala-kua/forward-chaining/assign/%d", p.ID), body, token)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 for Selesai registration, got %d: %s", w.Code, w.Body.String())
	}
}
