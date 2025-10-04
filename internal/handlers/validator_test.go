package handlers

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.id", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%s) = %v, expected %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"password123", true},
		{"123456", true},
		{"12345", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			result := ValidatePassword(tt.password)
			if result != tt.expected {
				t.Errorf("ValidatePassword(%s) = %v, expected %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		min      int
		max      int
		expected bool
	}{
		{"Valid length", "Hello", 3, 10, true},
		{"Too short", "Hi", 3, 10, false},
		{"Too long", "This is a very long string", 3, 10, false},
		{"Exact min", "ABC", 3, 10, true},
		{"Exact max", "1234567890", 3, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStringLength(tt.str, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("ValidateStringLength(%s, %d, %d) = %v, expected %v",
					tt.str, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}
