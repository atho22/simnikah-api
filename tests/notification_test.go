package tests

import (
	"fmt"
	"net/http"
	"testing"

	"simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== NOTIFICATION CRUD TESTS ====================

func TestNotification_CreateByStaff(t *testing.T) {
	testhelper.CleanupDB(testDB)

	targetUser := testhelper.CreateTestUser(testDB, "notif_target", "target@test.com", structs.UserRoleUserBiasa)
	staffUser := testhelper.CreateTestUser(testDB, "notif_staff", "staff@test.com", structs.UserRoleStaff)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"user_id": targetUser.User_id,
		"judul":   "Test Notif",
		"pesan":   "Test message",
		"tipe":    "Info",
	}

	w := doRequest("POST", "/simnikah/notifikasi", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotification_GetMyNotifications(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "notif_user", "notifuser@test.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")

	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/notifikasi/user/me", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	notifs := resp["notifications"].([]interface{})
	if len(notifs) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(notifs))
	}
}

func TestNotification_GetByID_OwnerOnly(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "notif_owner", "owner@test.com", structs.UserRoleUserBiasa)
	notif := testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", fmt.Sprintf("/simnikah/notifikasi/%d", notif.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for owner, got %d", w.Code)
	}
}

func TestNotification_UpdateStatus(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "notif_update", "update@test.com", structs.UserRoleUserBiasa)
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

func TestNotification_MarkAllRead(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "notif_markall", "markall@test.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("PUT", "/simnikah/notifikasi/mark-all-read", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	if resp["updated_count"].(float64) != 2 {
		t.Errorf("expected 2 updated, got %v", resp["updated_count"])
	}
}

func TestNotification_Delete(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "notif_del", "del@test.com", structs.UserRoleUserBiasa)
	notif := testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("DELETE", fmt.Sprintf("/simnikah/notifikasi/%d", notif.ID), nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestNotification_Stats(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "notif_stats", "stats@test.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")
	testhelper.CreateTestNotification(testDB, user.User_id, "Error")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/notifikasi/stats", nil, token)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	stats := resp["stats"].(map[string]interface{})
	if stats["total"].(float64) != 3 {
		t.Errorf("expected total 3, got %v", stats["total"])
	}
	if stats["unread"].(float64) != 3 {
		t.Errorf("expected unread 3, got %v", stats["unread"])
	}
}

func TestNotification_InvalidType(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "notif_inv", "notifinv@test.com", structs.UserRoleStaff)
	targetUser := testhelper.CreateTestUser(testDB, "notif_tgt2", "tgt2@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"user_id": targetUser.User_id,
		"judul":   "Test",
		"pesan":   "Test",
		"tipe":    "InvalidType",
	}

	w := doRequest("POST", "/simnikah/notifikasi", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid type, got %d", w.Code)
	}
}

func TestNotification_SendToRole(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "notif_sender", "sender@test.com", structs.UserRoleStaff)
	// Create users with penghulu role AND active status
	recv1 := testhelper.CreateTestUser(testDB, "notif_recv1", "recv1@test.com", structs.UserRolePenghulu)
	recv2 := testhelper.CreateTestUser(testDB, "notif_recv2", "recv2@test.com", structs.UserRolePenghulu)
	// Ensure active status
	testDB.Model(&recv1).Update("status", structs.UserStatusAktif)
	testDB.Model(&recv2).Update("status", structs.UserStatusAktif)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	body := map[string]interface{}{
		"role":  "penghulu",
		"judul": "Announcement",
		"pesan": "Meeting at 14:00",
		"tipe":  "Info",
	}

	w := doRequest("POST", "/simnikah/notifikasi/send-to-role", body, token)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
		return
	}

	resp := parseJSON(w)
	if count, ok := resp["recipient_count"].(float64); ok {
		if count < 2 {
			t.Errorf("expected at least 2 recipients, got %v", count)
		}
	} else {
		t.Error("missing recipient_count in response")
	}
}

// ==================== IDOR SECURITY TESTS ====================

func TestIDOR_CannotReadOtherUsersNotification(t *testing.T) {
	testhelper.CleanupDB(testDB)

	userA := testhelper.CreateTestUser(testDB, "idor_userA", "userA@test.com", structs.UserRoleUserBiasa)
	userB := testhelper.CreateTestUser(testDB, "idor_userB", "userB@test.com", structs.UserRoleUserBiasa)

	notifA := testhelper.CreateTestNotification(testDB, userA.User_id, "Info")
	tokenB := testhelper.GenerateTestToken(userB.User_id, userB.Email, userB.Role, userB.Nama)

	// User B tries to read User A's notification
	w := doRequest("GET", fmt.Sprintf("/simnikah/notifikasi/%d", notifA.ID), nil, tokenB)
	if w.Code != http.StatusNotFound {
		t.Errorf("IDOR: User B should NOT see User A's notification. Got %d", w.Code)
	}
}

func TestIDOR_CannotUpdateOtherUsersNotification(t *testing.T) {
	testhelper.CleanupDB(testDB)

	userA := testhelper.CreateTestUser(testDB, "idor_updA", "updA@test.com", structs.UserRoleUserBiasa)
	userB := testhelper.CreateTestUser(testDB, "idor_updB", "updB@test.com", structs.UserRoleUserBiasa)

	notifA := testhelper.CreateTestNotification(testDB, userA.User_id, "Info")
	tokenB := testhelper.GenerateTestToken(userB.User_id, userB.Email, userB.Role, userB.Nama)

	body := map[string]interface{}{"status_baca": "Sudah Dibaca"}
	w := doRequest("PUT", fmt.Sprintf("/simnikah/notifikasi/%d/status", notifA.ID), body, tokenB)
	if w.Code != http.StatusNotFound {
		t.Errorf("IDOR: User B should NOT update User A's notification. Got %d", w.Code)
	}
}

func TestIDOR_CannotDeleteOtherUsersNotification(t *testing.T) {
	testhelper.CleanupDB(testDB)

	userA := testhelper.CreateTestUser(testDB, "idor_delA", "delA@test.com", structs.UserRoleUserBiasa)
	userB := testhelper.CreateTestUser(testDB, "idor_delB", "delB@test.com", structs.UserRoleUserBiasa)

	notifA := testhelper.CreateTestNotification(testDB, userA.User_id, "Info")
	tokenB := testhelper.GenerateTestToken(userB.User_id, userB.Email, userB.Role, userB.Nama)

	w := doRequest("DELETE", fmt.Sprintf("/simnikah/notifikasi/%d", notifA.ID), nil, tokenB)
	if w.Code != http.StatusNotFound {
		t.Errorf("IDOR: User B should NOT delete User A's notification. Got %d", w.Code)
	}
}

func TestIDOR_CannotSeeOtherUsersNotifications(t *testing.T) {
	testhelper.CleanupDB(testDB)

	userA := testhelper.CreateTestUser(testDB, "idor_listA", "listA@test.com", structs.UserRoleUserBiasa)
	userB := testhelper.CreateTestUser(testDB, "idor_listB", "listB@test.com", structs.UserRoleUserBiasa)

	testhelper.CreateTestNotification(testDB, userA.User_id, "Info")
	testhelper.CreateTestNotification(testDB, userA.User_id, "Warning")
	tokenB := testhelper.GenerateTestToken(userB.User_id, userB.Email, userB.Role, userB.Nama)

	// User B's notification list should be empty (only sees own)
	w := doRequest("GET", "/simnikah/notifikasi/user/me", nil, tokenB)
	resp := parseJSON(w)

	var count float64
	if resp["notifications"] != nil {
		count = float64(len(resp["notifications"].([]interface{})))
	}
	if count != 0 {
		t.Errorf("IDOR: User B should see 0 notifications, got %.0f", count)
	}
}

func TestIDOR_MarkAllRead_OnlyOwnNotifications(t *testing.T) {
	testhelper.CleanupDB(testDB)

	userA := testhelper.CreateTestUser(testDB, "idor_markA", "markA@test.com", structs.UserRoleUserBiasa)
	userB := testhelper.CreateTestUser(testDB, "idor_markB", "markB@test.com", structs.UserRoleUserBiasa)

	testhelper.CreateTestNotification(testDB, userA.User_id, "Info")
	testhelper.CreateTestNotification(testDB, userB.User_id, "Info")
	tokenB := testhelper.GenerateTestToken(userB.User_id, userB.Email, userB.Role, userB.Nama)

	w := doRequest("PUT", "/simnikah/notifikasi/mark-all-read", nil, tokenB)
	resp := parseJSON(w)

	// Only User B's 1 notification should be marked as read
	if resp["updated_count"].(float64) != 1 {
		t.Errorf("expected only 1 notification marked, got %v", resp["updated_count"])
	}
}
