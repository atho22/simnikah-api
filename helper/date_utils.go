package helper

import "time"

// CalculateWorkingDays calculates working days between two dates (excluding weekends)
func CalculateWorkingDays(startDate, endDate time.Time) int {
	if endDate.Before(startDate) {
		return 0
	}

	workingDays := 0
	current := startDate

	for current.Before(endDate) || current.Equal(endDate) {
		// Check if current day is not Saturday (6) or Sunday (0)
		weekday := current.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			workingDays++
		}
		current = current.AddDate(0, 0, 1)
	}

	return workingDays
}

// CalculateAge calculates age in years
func CalculateAge(birthDate, currentDate time.Time) int {
	age := currentDate.Year() - birthDate.Year()

	// Adjust if birthday hasn't occurred this year
	if currentDate.Month() < birthDate.Month() ||
		(currentDate.Month() == birthDate.Month() && currentDate.Day() < birthDate.Day()) {
		age--
	}

	return age
}
