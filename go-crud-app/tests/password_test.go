package tests

import (
	"testing"

	"go-crud-app/internal/utils"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{
			name:        "Valid strong password",
			password:    "StrongPass123",
			shouldError: false,
		},
		{
			name:        "Weak password - too short",
			password:    "Short1",
			shouldError: true,
		},
		{
			name:        "Weak password - no uppercase",
			password:    "weakpass123",
			shouldError: true,
		},
		{
			name:        "Weak password - no lowercase",
			password:    "WEAKPASS123",
			shouldError: true,
		},
		{
			name:        "Weak password - no number",
			password:    "WeakPassword",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := utils.HashPassword(tt.password)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for password %s, but got none", tt.password)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for password %s: %v", tt.password, err)
				}
				if hash == "" {
					t.Error("Expected hash to be generated, but got empty string")
				}
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "TestPassword123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			expected: true,
		},
		{
			name:     "Incorrect password",
			password: "WrongPassword123",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CheckPassword(tt.password, tt.hash)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{
			name:        "Valid password",
			password:    "ValidPass123",
			shouldError: false,
		},
		{
			name:        "Too short",
			password:    "Short1",
			shouldError: true,
		},
		{
			name:        "No uppercase",
			password:    "lowercase123",
			shouldError: true,
		},
		{
			name:        "No lowercase",
			password:    "UPPERCASE123",
			shouldError: true,
		},
		{
			name:        "No number",
			password:    "NoNumberPass",
			shouldError: true,
		},
		{
			name:        "Complex valid password",
			password:    "C0mpl3xP@ssw0rd!",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidatePassword(tt.password)

			if tt.shouldError && err == nil {
				t.Errorf("Expected error for password %s, but got none", tt.password)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for password %s: %v", tt.password, err)
			}
		})
	}
}
