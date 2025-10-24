package helper

import "strings"

func StartsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// IsValidEmail validates email format (basic validation)
func IsValidEmail(email string) bool {
	return len(email) > 0 && strings.Contains(email, "@") && strings.Contains(email, ".")
}
