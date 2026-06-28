package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"simnikah/tests/testhelper"

	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
	router http.Handler
)

// TestMain sets up the test database and router once for all tests
func TestMain(m *testing.M) {
	testDB = testhelper.SetupTestDB()
	router = testhelper.SetupRouter(testDB)
	code := m.Run()
	// Save captured API data for report generation
	testhelper.SaveReport("test_api_data.json")
	os.Exit(code)
}

func doRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", testhelper.AuthHeader(token))
	}

	w := httptest.NewRecorder()
	start := time.Now()
	router.ServeHTTP(w, req)
	dur := time.Since(start)

	// Auto-capture for report
	testhelper.RecordAPI(method, path, body, w, dur)

	return w
}

func parseJSON(w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

// ==================== REGISTER TESTS ====================

func TestRegister_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)

	body := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@test.com",
		"password": "password123",
		"nama":     "New User",
	}

	w := doRequest("POST", "/register", body, "")

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["success"] != true {
		t.Error("expected success=true")
	}

	data := resp["data"].(map[string]interface{})
	if data["role"] != "user_biasa" {
		t.Errorf("expected role user_biasa, got %v", data["role"])
	}
	if data["username"] != "newuser" {
		t.Errorf("expected username newuser, got %v", data["username"])
	}
}

func TestRegister_MissingFields(t *testing.T) {
	testhelper.CleanupDB(testDB)

	body := map[string]interface{}{
		"username": "testuser",
		// missing email, password, nama
	}

	w := doRequest("POST", "/register", body, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	testhelper.CleanupDB(testDB)

	body := map[string]interface{}{
		"username": "testuser",
		"email":    "test@test.com",
		"password": "123", // too short (min 6)
		"nama":     "Test",
	}

	w := doRequest("POST", "/register", body, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	testhelper.CleanupDB(testDB)

	body := map[string]interface{}{
		"username": "duplicate",
		"email":    "first@test.com",
		"password": "password123",
		"nama":     "First User",
	}
	doRequest("POST", "/register", body, "")

	// Second registration with same username
	body2 := map[string]interface{}{
		"username": "duplicate",
		"email":    "second@test.com",
		"password": "password123",
		"nama":     "Second User",
	}
	w := doRequest("POST", "/register", body2, "")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 (duplicate), got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	testhelper.CleanupDB(testDB)

	body := map[string]interface{}{
		"username": "user1",
		"email":    "same@test.com",
		"password": "password123",
		"nama":     "User 1",
	}
	doRequest("POST", "/register", body, "")

	body2 := map[string]interface{}{
		"username": "user2",
		"email":    "same@test.com",
		"password": "password123",
		"nama":     "User 2",
	}
	w := doRequest("POST", "/register", body2, "")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 (duplicate email), got %d", w.Code)
	}
}

func TestRegister_RoleAlwaysUserBiasa(t *testing.T) {
	testhelper.CleanupDB(testDB)

	// Even if user tries to send a role, it should be ignored
	body := map[string]interface{}{
		"username": "wannabe_admin",
		"email":    "admin@test.com",
		"password": "password123",
		"nama":     "Wannabe Admin",
		"role":     "kepala_kua", // should be ignored
	}

	w := doRequest("POST", "/register", body, "")

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["role"] != "user_biasa" {
		t.Errorf("SECURITY: role should always be user_biasa, got %v", data["role"])
	}
}

// ==================== LOGIN TESTS ====================

func TestLogin_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.CreateTestUser(testDB, "logintest", "login@test.com", "user_biasa")

	body := map[string]interface{}{
		"username": "logintest",
		"password": "password123",
	}

	w := doRequest("POST", "/login", body, "")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["token"] == nil || resp["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	testhelper.CleanupDB(testDB)
	testhelper.CreateTestUser(testDB, "wrongpw", "wrongpw@test.com", "user_biasa")

	body := map[string]interface{}{
		"username": "wrongpw",
		"password": "wrongpassword",
	}

	w := doRequest("POST", "/login", body, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	testhelper.CleanupDB(testDB)

	body := map[string]interface{}{
		"username": "ghost",
		"password": "password123",
	}

	w := doRequest("POST", "/login", body, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	testhelper.CleanupDB(testDB)
	user := testhelper.CreateTestUser(testDB, "inactive", "inactive@test.com", "user_biasa")

	// Deactivate user
	testDB.Model(&user).Update("status", "Nonaktif")

	body := map[string]interface{}{
		"username": "inactive",
		"password": "password123",
	}

	w := doRequest("POST", "/login", body, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for inactive user, got %d", w.Code)
	}
}

// ==================== PROFILE TESTS ====================

func TestGetProfile_Success(t *testing.T) {
	testhelper.CleanupDB(testDB)
	user := testhelper.CreateTestUser(testDB, "profiletest", "profile@test.com", "user_biasa")
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/profile", nil, token)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	data := resp["data"].(map[string]interface{})
	if data["username"] != "profiletest" {
		t.Errorf("expected username profiletest, got %v", data["username"])
	}
}

func TestGetProfile_WithoutToken(t *testing.T) {
	w := doRequest("GET", "/profile", nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

// ==================== HEALTH CHECK ====================

func TestHealthCheck(t *testing.T) {
	w := doRequest("GET", "/health", nil, "")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	if resp["status"] != "healthy" {
		t.Errorf("expected status healthy, got %v", resp["status"])
	}
}
