package helper

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// IsValidTimeFormat validates if the time string is in HH:MM format
func IsValidTimeFormat(timeStr string) bool {
	// Check if the string matches HH:MM format
	matched, err := regexp.MatchString(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`, timeStr)
	if err != nil || !matched {
		return false
	}

	// Try to parse the time to ensure it's valid
	_, err = time.Parse("15:04", timeStr)
	return err == nil
}

// validateNIKExternal simulates external NIK validation (SIMKAH/Dukcapil)
// In a real implementation, this would call an external API
func validateNIKExternal(nik string) (map[string]interface{}, error) {
	// Basic NIK validation
	if len(nik) != 16 {
		return nil, errors.New("NIK harus 16 digit")
	}

	// Check if all characters are digits
	for _, char := range nik {
		if char < '0' || char > '9' {
			return nil, errors.New("NIK harus berupa angka")
		}
	}

	// Extract NIK components for validation
	provinceCode := nik[0:2]
	regencyCode := nik[2:4]
	districtCode := nik[4:6]
	birthDate := nik[6:12]
	genderCode := nik[15:16]

	// Validate province code (01-94 for Indonesia)
	provinceNum, _ := strconv.Atoi(provinceCode)
	if provinceNum < 1 || provinceNum > 94 {
		return nil, errors.New("kode provinsi tidak valid")
	}

	// Validate birth date format
	day := birthDate[0:2]
	month := birthDate[2:4]
	year := birthDate[4:6]

	dayNum, _ := strconv.Atoi(day)
	monthNum, _ := strconv.Atoi(month)
	yearNum, _ := strconv.Atoi(year)

	// Adjust day for female (add 40)
	genderNum, _ := strconv.Atoi(genderCode)
	if genderNum >= 4 {
		dayNum -= 40
	}

	// Basic date validation
	if dayNum < 1 || dayNum > 31 {
		return nil, errors.New("tanggal lahir tidak valid")
	}
	if monthNum < 1 || monthNum > 12 {
		return nil, errors.New("bulan lahir tidak valid")
	}
	if yearNum < 0 || yearNum > 99 {
		return nil, errors.New("tahun lahir tidak valid")
	}

	// Determine gender
	gender := "Laki-laki"
	if genderNum >= 4 {
		gender = "Perempuan"
	}

	// Simulate external API response
	nikData := map[string]interface{}{
		"nik":           nik,
		"provinsi":      "Provinsi " + provinceCode,
		"kabupaten":     "Kabupaten " + regencyCode,
		"kecamatan":     "Kecamatan " + districtCode,
		"tanggal_lahir": day + "/" + month + "/" + year,
		"jenis_kelamin": gender,
		"status":        "Valid",
		"message":       "NIK valid dan terdaftar di sistem",
	}

	return nikData, nil
}

// ValidateNIK validates NIK format and calls external validation
func ValidateNIK(nik string) (map[string]interface{}, error) {
	return validateNIKExternal(nik)
}
