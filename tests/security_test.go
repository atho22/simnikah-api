package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"simnikah/internal/models"
	"simnikah/tests/testhelper"

	"github.com/golang-jwt/jwt/v5"
)

// ==================== JWT SECURITY TESTS ====================

func TestJWT_ExpiredToken(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Create a token that's already expired
	claims := jwt.MapClaims{
		"user_id": "USR-EXPIRED",
		"email":   "expired@test.com",
		"role":    "user_biasa",
		"nama":    "Expired User",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
		"nbf":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret-key-for-unit-testing-only"))

	w := doRequest("GET", "/profile", nil, tokenStr)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("SECURITY: expected 401 for expired token, got %d", w.Code)
	}
}

func TestJWT_InvalidSignature(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Token signed with different key
	claims := jwt.MapClaims{
		"user_id": "USR-FAKE",
		"email":   "fake@test.com",
		"role":    "kepala_kua",
		"nama":    "Fake Admin",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("completely-different-secret-key"))

	w := doRequest("GET", "/profile", nil, tokenStr)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("SECURITY: expected 401 for wrong signature, got %d", w.Code)
	}
}

func TestJWT_NoneAlgorithm(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Attempt alg:none attack (no signature)
	claims := jwt.MapClaims{
		"user_id": "USR-ATTACKER",
		"email":   "attacker@test.com",
		"role":    "kepala_kua",
		"nama":    "Attacker",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	w := doRequest("GET", "/simnikah/kepala-kua/available-penghulu", nil, tokenStr)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("SECURITY: alg:none attack should fail, got %d", w.Code)
	}
}

func TestJWT_MalformedToken(t *testing.T) {
	testhelper.CleanupDB(testDB)

	malformedTokens := []string{
		"not.a.valid.jwt",
		"abc123",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9", // only header
	}

	for _, tokenStr := range malformedTokens {
		t.Run(fmt.Sprintf("malformed_%d", len(tokenStr)), func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/profile", nil)
			req.Header.Set("Authorization", "Bearer "+tokenStr)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("SECURITY: expected 401 for malformed token, got %d", w.Code)
			}
		})
	}

	// Test completely missing Authorization header
	t.Run("missing_auth_header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("SECURITY: expected 401 without auth header, got %d", w.Code)
		}
	})
}

func TestJWT_RoleEscalation_ForgedRole(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Token with forged kepala_kua role but invalid signature
	claims := jwt.MapClaims{
		"user_id": "USR-ESCALATE",
		"email":   "escalate@test.com",
		"role":    "kepala_kua", // attempting privilege escalation
		"nama":    "Escalated User",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("wrong-key"))

	w := doRequest("GET", "/simnikah/kepala-kua/available-penghulu", nil, tokenStr)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("SECURITY: role escalation with wrong key should fail, got %d", w.Code)
	}
}

func TestJWT_InvalidRole(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "inv_role", "invrole@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, "superadmin", user.Nama)

	w := doRequest("GET", "/profile", nil, token)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("SECURITY: invalid role should be rejected, got %d", w.Code)
	}
}

// ==================== INPUT INJECTION TESTS ====================

func TestInjection_SQLInjection_RegistrationID(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "inj_test", "inj@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	// Try SQL injection in registration ID
	maliciousIDs := []string{
		"1 OR 1=1",
		"1; DROP TABLE users;--",
		"' UNION SELECT * FROM users--",
		"1' AND '1'='1",
	}

	for _, id := range maliciousIDs {
		t.Run(id, func(t *testing.T) {
			w := doRequest("GET", fmt.Sprintf("/simnikah/pendaftaran/%s/location", id), nil, token)
			// Should get 400, 404, or 403, but never 200 with leaked data
			if w.Code == http.StatusOK {
				t.Errorf("SECURITY: SQL injection in ID should not return 200, got %d", w.Code)
			}
		})
	}
}

func TestInjection_XSSInNotification(t *testing.T) {
	testhelper.CleanupDB(testDB)

	staffUser := testhelper.CreateTestUser(testDB, "xss_staff", "xssstaff@test.com", structs.UserRoleStaff)
	targetUser := testhelper.CreateTestUser(testDB, "xss_target", "xsstarget@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(staffUser.User_id, staffUser.Email, staffUser.Role, staffUser.Nama)

	// Create notification with XSS payload
	body := map[string]interface{}{
		"user_id": targetUser.User_id,
		"judul":   "<script>alert('xss')</script>",
		"pesan":   "<img src=x onerror=alert('xss')>",
		"tipe":    "Info",
	}

	w := doRequest("POST", "/simnikah/notifikasi", body, token)
	// The notification should be created (server-side doesn't sanitize)
	// but the frontend must sanitize on render
	if w.Code != http.StatusCreated {
		t.Logf("Note: XSS payload notification returned %d (OK if validated)", w.Code)
	}

	// Verify content-type is application/json (prevents content-type sniffing XSS)
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("SECURITY: Content-Type should be application/json, got %s", contentType)
	}
}

// ==================== HTTP METHOD TESTS ====================

func TestHTTPMethods_NoUnexpectedMethods(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "method_test", "method@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	// Try DELETE on endpoints that should only accept GET
	wrongMethods := []struct {
		method string
		path   string
	}{
		{"DELETE", "/simnikah/kepala-kua/available-penghulu"},
		{"DELETE", "/simnikah/notifikasi/user/me"},
		{"POST", "/simnikah/penghulu/jadwal-penugasan"},
		{"PUT", "/health"},
	}

	for _, ep := range wrongMethods {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(ep.method, ep.path, nil, token)
			if w.Code == http.StatusOK {
				t.Errorf("SECURITY: %s should not work on %s", ep.method, ep.path)
			}
		})
	}
}

// ==================== RATE LIMITING ====================

func Test404_UnknownEndpoint(t *testing.T) {
	w := doRequest("GET", "/nonexistent-endpoint", nil, "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// ==================== EDGE CASES ====================

func TestRegistration_InvalidID(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "edge_id", "edgeid@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/pendaftaran/abc/location", nil, token)
	// Should fail gracefully
	if w.Code == http.StatusOK {
		t.Error("expected non-200 for invalid ID format")
	}
}

func TestNotification_PaginationLimits(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "page_test", "page@test.com", structs.UserRoleUserBiasa)
	for i := 0; i < 15; i++ {
		testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	}
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	// Request page 1 with limit 5
	w := doRequest("GET", "/simnikah/notifikasi/user/me?page=1&limit=5", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	pagination := resp["pagination"].(map[string]interface{})
	if pagination["total"].(float64) != 15 {
		t.Errorf("expected total 15, got %v", pagination["total"])
	}

	notifs := resp["notifications"].([]interface{})
	if len(notifs) != 5 {
		t.Errorf("expected 5 items per page, got %d", len(notifs))
	}
}

func TestNotification_FilterByStatus(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "filter_test", "filter@test.com", structs.UserRoleUserBiasa)
	testhelper.CreateTestNotification(testDB, user.User_id, "Info")
	testhelper.CreateTestNotification(testDB, user.User_id, "Warning")

	// Mark one as read
	var notif structs.Notifikasi
	testDB.Where("user_id = ? AND tipe = ?", user.User_id, "Info").First(&notif)
	testDB.Model(&notif).Update("status_baca", "Sudah Dibaca")

	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/notifikasi/user/me?status=Belum+Dibaca", nil, token)
	resp := parseJSON(w)
	notifs := resp["notifications"].([]interface{})
	if len(notifs) != 1 {
		t.Errorf("expected 1 unread notification, got %d", len(notifs))
	}
}
