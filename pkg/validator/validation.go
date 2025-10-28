package validator

import (
	"simnikah/pkg/utils"
)

import (
	"regexp"
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
