package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GenerateUserID(prefix string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return prefix + string(b)
}

// Struct untuk response Nominatim API
type NominatimResponse struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// GetCoordinatesFromAddress mendapatkan koordinat dari alamat menggunakan OpenStreetMap Nominatim API (GRATIS)
func GetCoordinatesFromAddress(address string) (float64, float64, error) {
	// URL encode alamat
	encodedAddress := url.QueryEscape(address)

	// URL Nominatim API dengan filter Indonesia
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&countrycodes=id&limit=1", encodedAddress)

	// Buat HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal membuat request: %v", err)
	}

	// Set User-Agent (diperlukan oleh Nominatim)
	req.Header.Set("User-Agent", "SimNikah-KUA-Banjarmasin/1.0")

	// Kirim request
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal membaca response: %v", err)
	}

	// Parse JSON response
	var results []NominatimResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return 0, 0, fmt.Errorf("gagal parse JSON: %v", err)
	}

	// Cek apakah ada hasil
	if len(results) == 0 {
		return 0, 0, fmt.Errorf("alamat tidak ditemukan")
	}

	// Konversi string ke float64
	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal parse latitude: %v", err)
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal parse longitude: %v", err)
	}

	return lat, lon, nil
}
