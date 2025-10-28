package catin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"simnikah/helper"
	"simnikah/structs"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== LOCATION & MAP ENDPOINTS ====================
// Endpoints untuk mendukung integrasi peta menggunakan OpenStreetMap (100% GRATIS)
// Menggunakan Nominatim API untuk geocoding dan reverse geocoding

// GetCoordinatesFromAddressEndpoint - Mendapatkan koordinat dari alamat (Geocoding)
// Digunakan saat user mengetik alamat dan butuh koordinat untuk map
func (h *InDB) GetCoordinatesFromAddressEndpoint(c *gin.Context) {
	var input struct {
		Alamat string `json:"alamat" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Alamat diperlukan",
			"error":   err.Error(),
		})
		return
	}

	// Validasi alamat tidak kosong
	if input.Alamat == "" || len(input.Alamat) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Alamat terlalu pendek",
			"error":   "Masukkan alamat lengkap minimal 10 karakter",
		})
		return
	}

	// Dapatkan koordinat menggunakan OpenStreetMap Nominatim API (dengan caching)
	lat, lon, err := helper.GetCoordinatesFromAddressCached(input.Alamat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mendapatkan koordinat",
			"error":   fmt.Sprintf("Tidak dapat menemukan lokasi untuk alamat: %s. Pastikan alamat lengkap dan benar.", input.Alamat),
			"details": err.Error(),
		})
		return
	}

	// Validasi koordinat valid
	if lat == 0 && lon == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Lokasi tidak ditemukan",
			"error":   "Alamat tidak dapat ditemukan di peta. Silakan periksa kembali alamat yang dimasukkan.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Koordinat berhasil ditemukan",
		"data": gin.H{
			"alamat":    input.Alamat,
			"latitude":  lat,
			"longitude": lon,
			"map_url":   fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, lon),
			"osm_url":   fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", lat, lon),
		},
	})
}

// GetAddressFromCoordinates - Mendapatkan alamat dari koordinat (Reverse Geocoding)
// Digunakan saat user memilih lokasi dari map dan butuh alamat
func (h *InDB) GetAddressFromCoordinates(c *gin.Context) {
	var input struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Latitude dan longitude diperlukan",
			"error":   err.Error(),
		})
		return
	}

	// Validasi koordinat
	if input.Latitude < -90 || input.Latitude > 90 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Latitude tidak valid",
			"error":   "Latitude harus antara -90 dan 90",
		})
		return
	}

	if input.Longitude < -180 || input.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Longitude tidak valid",
			"error":   "Longitude harus antara -180 dan 180",
		})
		return
	}

	// Reverse geocoding menggunakan Nominatim API
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=18&addressdetails=1", input.Latitude, input.Longitude)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat request",
			"error":   err.Error(),
		})
		return
	}

	// Set User-Agent (diperlukan oleh Nominatim)
	req.Header.Set("User-Agent", "SimNikah-KUA-Banjarmasin/1.0")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengirim request",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membaca response",
			"error":   err.Error(),
		})
		return
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal parse response",
			"error":   err.Error(),
		})
		return
	}

	// Cek apakah ada error dari API
	if errorMsg, exists := result["error"]; exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Lokasi tidak ditemukan",
			"error":   errorMsg,
		})
		return
	}

	// Extract alamat
	displayName, _ := result["display_name"].(string)
	address, _ := result["address"].(map[string]interface{})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alamat berhasil ditemukan",
		"data": gin.H{
			"latitude":  input.Latitude,
			"longitude": input.Longitude,
			"alamat":    displayName,
			"detail":    address,
			"map_url":   fmt.Sprintf("https://www.google.com/maps?q=%f,%f", input.Latitude, input.Longitude),
			"osm_url":   fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", input.Latitude, input.Longitude),
		},
	})
}

// UpdateWeddingLocationWithCoordinates - Update alamat nikah beserta koordinatnya
func (h *InDB) UpdateWeddingLocationWithCoordinates(c *gin.Context) {
	registrationID := c.Param("id")

	var input struct {
		AlamatAkad string   `json:"alamat_akad" binding:"required"`
		Latitude   *float64 `json:"latitude"`  // Optional, bisa didapat dari geocoding
		Longitude  *float64 `json:"longitude"` // Optional, bisa didapat dari geocoding
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	// Check if registration exists and belongs to this user
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ? AND pendaftar_id = ?", registrationID, userID.(string)).First(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pendaftaran tidak ditemukan",
		})
		return
	}

	// Check if wedding location is outside KUA
	if pendaftaran.Tempat_nikah != "Di Luar KUA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Alamat hanya dapat diubah untuk nikah di luar KUA",
		})
		return
	}

	// Jika koordinat tidak disediakan, dapatkan dari alamat (geocoding)
	var latitude, longitude *float64
	if input.Latitude != nil && input.Longitude != nil {
		latitude = input.Latitude
		longitude = input.Longitude
	} else {
		// Auto-geocoding (dengan caching)
		lat, lon, err := helper.GetCoordinatesFromAddressCached(input.AlamatAkad)
		if err != nil {
			// Log warning tapi tetap lanjut simpan alamat
			fmt.Printf("Warning: Gagal mendapatkan koordinat untuk alamat '%s': %v\n", input.AlamatAkad, err)
		} else {
			latitude = &lat
			longitude = &lon
		}
	}

	// Update alamat dan koordinat
	updates := map[string]interface{}{
		"alamat_akad": input.AlamatAkad,
		"updated_at":  time.Now(),
	}

	if latitude != nil {
		updates["latitude"] = latitude
	}
	if longitude != nil {
		updates["longitude"] = longitude
	}

	if err := h.DB.Model(&pendaftaran).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate lokasi",
			"error":   err.Error(),
		})
		return
	}

	// Get updated data
	h.DB.Where("id = ?", registrationID).First(&pendaftaran)

	response := gin.H{
		"pendaftaran_id":    pendaftaran.ID,
		"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
		"alamat_akad":       pendaftaran.Alamat_akad,
		"tempat_nikah":      pendaftaran.Tempat_nikah,
		"updated_at":        pendaftaran.Updated_at,
	}

	// Tambahkan koordinat dan map URL jika tersedia
	if pendaftaran.Latitude != nil && pendaftaran.Longitude != nil {
		response["latitude"] = *pendaftaran.Latitude
		response["longitude"] = *pendaftaran.Longitude
		response["map_url"] = fmt.Sprintf("https://www.google.com/maps?q=%f,%f", *pendaftaran.Latitude, *pendaftaran.Longitude)
		response["osm_url"] = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", *pendaftaran.Latitude, *pendaftaran.Longitude)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lokasi nikah berhasil diupdate",
		"data":    response,
	})
}

// GetWeddingLocationDetail - Mendapatkan detail lokasi nikah dengan koordinat untuk penghulu
func (h *InDB) GetWeddingLocationDetail(c *gin.Context) {
	registrationID := c.Param("id")

	// Get registration
	var pendaftaran structs.PendaftaranNikah
	if err := h.DB.Where("id = ?", registrationID).First(&pendaftaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pendaftaran tidak ditemukan",
		})
		return
	}

	response := gin.H{
		"pendaftaran_id":    pendaftaran.ID,
		"nomor_pendaftaran": pendaftaran.Nomor_pendaftaran,
		"tanggal_nikah":     pendaftaran.Tanggal_nikah,
		"waktu_nikah":       pendaftaran.Waktu_nikah,
		"tempat_nikah":      pendaftaran.Tempat_nikah,
		"alamat_akad":       pendaftaran.Alamat_akad,
	}

	// Jika ada koordinat, tambahkan info map
	if pendaftaran.Latitude != nil && pendaftaran.Longitude != nil {
		response["latitude"] = *pendaftaran.Latitude
		response["longitude"] = *pendaftaran.Longitude
		response["has_coordinates"] = true
		response["map_url"] = fmt.Sprintf("https://www.google.com/maps?q=%f,%f", *pendaftaran.Latitude, *pendaftaran.Longitude)
		response["google_maps_url"] = fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", *pendaftaran.Latitude, *pendaftaran.Longitude)
		response["google_maps_directions_url"] = fmt.Sprintf("https://www.google.com/maps/dir/?api=1&destination=%f,%f", *pendaftaran.Latitude, *pendaftaran.Longitude)
		response["osm_url"] = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", *pendaftaran.Latitude, *pendaftaran.Longitude)
		response["waze_url"] = fmt.Sprintf("https://www.waze.com/ul?ll=%f,%f&navigate=yes", *pendaftaran.Latitude, *pendaftaran.Longitude)
	} else {
		response["has_coordinates"] = false
		response["message"] = "Koordinat belum tersedia untuk lokasi ini"
	}

	// Jika nikah di luar KUA, tambahkan info bahwa ini nikah di luar
	if pendaftaran.Tempat_nikah == "Di Luar KUA" {
		response["is_outside_kua"] = true
		response["note"] = "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
	} else {
		response["is_outside_kua"] = false
		response["note"] = "Pernikahan dilaksanakan di KUA"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail lokasi nikah berhasil diambil",
		"data":    response,
	})
}

// SearchAddress - Search alamat menggunakan Nominatim API (untuk autocomplete)
func (h *InDB) SearchAddress(c *gin.Context) {
	query := c.Query("q")
	if query == "" || len(query) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Query terlalu pendek",
			"error":   "Masukkan minimal 3 karakter untuk pencarian",
		})
		return
	}

	// URL encode query
	encodedQuery := url.QueryEscape(query)

	// Nominatim API dengan filter Indonesia
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&countrycodes=id&limit=5&addressdetails=1", encodedQuery)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat request",
		})
		return
	}

	req.Header.Set("User-Agent", "SimNikah-KUA-Banjarmasin/1.0")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengirim request",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membaca response",
		})
		return
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal parse response",
		})
		return
	}

	// Format results
	var locations []gin.H
	for _, result := range results {
		lat, _ := result["lat"].(string)
		lon, _ := result["lon"].(string)
		displayName, _ := result["display_name"].(string)

		locations = append(locations, gin.H{
			"display_name": displayName,
			"latitude":     lat,
			"longitude":    lon,
			"address":      result["address"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hasil pencarian alamat",
		"data": gin.H{
			"query":   query,
			"results": locations,
			"count":   len(locations),
		},
	})
}
