package handlers

import (
	"regexp"
)

// ValidateEmail - Validasi format email
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword - Validasi password minimal 6 karakter
func ValidatePassword(password string) bool {
	return len(password) >= 6
}

// ValidateStringLength - Validasi panjang string
func ValidateStringLength(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

// ValidateRequired - Validasi field required
func ValidateRequired(fields map[string]string) (bool, string) {
	for fieldName, fieldValue := range fields {
		if fieldValue == "" {
			return false, fieldName + " is required"
		}
	}
	return true, ""
}
