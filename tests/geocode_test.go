package tests

import (
	"fmt"
	"net/http"
	"testing"

	structs "simnikah/internal/models"
	"simnikah/tests/testhelper"
)

// ==================== GEOCODE INPUT VALIDATION TESTS ====================
// Catatan: endpoint geocode dan reverse-geocode memanggil API eksternal (Nominatim).
// Test di sini memfokuskan pada validasi INPUT (server-side) bukan hasil geocoding
// karena hasil geocoding bergantung pada jaringan dan API eksternal.

// TestGeocode_MissingAlamat memastikan 400 jika field 'alamat' tidak dikirim.
func TestGeocode_MissingAlamat(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "geo_miss", "geomiss@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		// 'alamat' tidak ada
	}

	w := doRequest("POST", "/simnikah/location/geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing alamat, got %d: %s", w.Code, w.Body.String())
	}
}

// TestGeocode_AlamatTerlalu Pendek memastikan 400 jika alamat < 10 karakter.
func TestGeocode_AlamatTerlaluPendek(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "geo_short", "geoshort@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"alamat": "Jl. A", // kurang dari 10 karakter
	}

	w := doRequest("POST", "/simnikah/location/geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for short address (<10 chars), got %d: %s", w.Code, w.Body.String())
	}
}

// TestGeocode_AlamatCukupPanjang memastikan validasi panjang lolos saat alamat >= 10 karakter.
// Catatan: mungkin mendapat 500 jika API Nominatim tidak bisa diakses, ini OK di environment test.
func TestGeocode_AlamatCukupPanjang(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "geo_long", "geolong@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"alamat": "Jl. A. Yani KM 5 Banjarmasin", // >= 10 karakter
	}

	w := doRequest("POST", "/simnikah/location/geocode", body, token)
	// Tidak boleh 400 (validasi lolos), bisa 200 (berhasil) atau 500 (network error)
	if w.Code == http.StatusBadRequest {
		t.Errorf("should not get 400 for valid-length address, got %d: %s", w.Code, w.Body.String())
	}
}

// TestGeocode_WithoutToken memastikan 401 tanpa token.
func TestGeocode_WithoutToken(t *testing.T) {
	body := map[string]interface{}{
		"alamat": "Jl. A. Yani KM 5 Banjarmasin",
	}

	w := doRequest("POST", "/simnikah/location/geocode", body, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

// ==================== REVERSE GEOCODE INPUT VALIDATION TESTS ====================

// TestReverseGeocode_MissingLatitude memastikan 400 jika latitude tidak dikirim.
func TestReverseGeocode_MissingLatitude(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_miss_lat", "revmisslat@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"longitude": 114.5881,
		// latitude tidak ada
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing latitude, got %d", w.Code)
	}
}

// TestReverseGeocode_MissingLongitude memastikan 400 jika longitude tidak dikirim.
func TestReverseGeocode_MissingLongitude(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_miss_lon", "revmisslon@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude": -3.2913,
		// longitude tidak ada
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing longitude, got %d", w.Code)
	}
}

// TestReverseGeocode_LatitudeTerlaluBesar memastikan 400 jika latitude > 90.
func TestReverseGeocode_LatitudeTerlaluBesar(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_lat_big", "revlatbig@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude":  95.0, // tidak valid, max 90
		"longitude": 114.5881,
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for latitude > 90, got %d: %s", w.Code, w.Body.String())
	}
}

// TestReverseGeocode_LatitudeTerlaluKecil memastikan 400 jika latitude < -90.
func TestReverseGeocode_LatitudeTerlaluKecil(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_lat_sml", "revlatsml@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude":  -95.0, // tidak valid, min -90
		"longitude": 114.5881,
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for latitude < -90, got %d: %s", w.Code, w.Body.String())
	}
}

// TestReverseGeocode_LongitudeTerlaluBesar memastikan 400 jika longitude > 180.
func TestReverseGeocode_LongitudeTerlaluBesar(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_lon_big", "revlonbig@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude":  -3.2913,
		"longitude": 190.0, // tidak valid, max 180
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for longitude > 180, got %d: %s", w.Code, w.Body.String())
	}
}

// TestReverseGeocode_LongitudeTerlaluKecil memastikan 400 jika longitude < -180.
func TestReverseGeocode_LongitudeTerlaluKecil(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_lon_sml", "revlonsml@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude":  -3.2913,
		"longitude": -190.0, // tidak valid, min -180
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for longitude < -180, got %d: %s", w.Code, w.Body.String())
	}
}

// TestReverseGeocode_ValidCoordinates memastikan koordinat valid lolos validasi input.
// Catatan: mungkin mendapat 200 atau 500 (network), keduanya OK.
func TestReverseGeocode_ValidCoordinates(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "rev_valid", "revvalid@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	body := map[string]interface{}{
		"latitude":  -3.2913, // koordinat Banjarmasin
		"longitude": 114.5881,
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, token)
	// Tidak boleh error validasi (400)
	if w.Code == http.StatusBadRequest {
		t.Errorf("should not get 400 for valid coordinates: %s", w.Body.String())
	}
}

// TestReverseGeocode_WithoutToken memastikan 401 tanpa token.
func TestReverseGeocode_WithoutToken(t *testing.T) {
	body := map[string]interface{}{
		"latitude":  -3.2913,
		"longitude": 114.5881,
	}

	w := doRequest("POST", "/simnikah/location/reverse-geocode", body, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

// ==================== SEARCH ADDRESS INPUT VALIDATION TESTS ====================

// TestSearchAddress_QueryTerlaluPendek memastikan 400 jika query < 3 karakter.
func TestSearchAddress_QueryTerlaluPendek(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "srch_short", "srchshort@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/location/search?q=ab", nil, token) // 2 karakter, terlalu pendek
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for query < 3 chars, got %d: %s", w.Code, w.Body.String())
	}
}

// TestSearchAddress_QueryKosong memastikan 400 jika query tidak dikirim sama sekali.
func TestSearchAddress_QueryKosong(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "srch_empty", "srchempty@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/location/search", nil, token) // tanpa query param
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty query, got %d: %s", w.Code, w.Body.String())
	}
}

// TestSearchAddress_QueryValid memastikan query >= 3 karakter lolos validasi.
// Catatan: hasil mungkin 200 (berhasil) atau 500 (network issue), keduanya OK.
func TestSearchAddress_QueryValid(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "srch_valid", "srchvalid@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/location/search?q=Banjarmasin", nil, token)
	// Tidak boleh error validasi
	if w.Code == http.StatusBadRequest {
		t.Errorf("should not get 400 for valid query: %s", w.Body.String())
	}
}

// TestSearchAddress_WithoutToken memastikan 401 tanpa token.
func TestSearchAddress_WithoutToken(t *testing.T) {
	w := doRequest("GET", "/simnikah/location/search?q=Banjarmasin", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", w.Code)
	}
}

// TestSearchAddress_QueryTepat3Karakter memastikan tepat 3 karakter lolos validasi.
func TestSearchAddress_QueryTepat3Karakter(t *testing.T) {
	testhelper.CleanupDB(testDB)

	user := testhelper.CreateTestUser(testDB, "srch_3chr", "srch3chr@test.com", structs.UserRoleUserBiasa)
	token := testhelper.GenerateTestToken(user.User_id, user.Email, user.Role, user.Nama)

	w := doRequest("GET", "/simnikah/location/search?q=KUA", nil, token) // tepat 3 karakter
	// Tidak boleh 400
	if w.Code == http.StatusBadRequest {
		t.Errorf("should not get 400 for exactly 3-char query: %s", w.Body.String())
	}
}

// ==================== COMBINED LOCATION SECURITY TESTS ====================

// TestLocation_AllEndpoints_RequireToken memastikan semua endpoint location butuh token.
func TestLocation_AllEndpoints_RequireToken(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"POST", "/simnikah/location/geocode", map[string]interface{}{"alamat": "Jl. Test No. 1 Banjarmasin"}},
		{"POST", "/simnikah/location/reverse-geocode", map[string]interface{}{"latitude": -3.29, "longitude": 114.58}},
		{"GET", "/simnikah/location/search?q=Banjarmasin", nil},
	}

	for _, ep := range endpoints {
		t.Run(fmt.Sprintf("%s %s", ep.method, ep.path), func(t *testing.T) {
			w := doRequest(ep.method, ep.path, ep.body, "")
			if w.Code != http.StatusUnauthorized {
				t.Errorf("SECURITY: expected 401 without token on %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}
