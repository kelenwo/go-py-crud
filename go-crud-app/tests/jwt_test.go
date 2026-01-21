package tests

import (
	"testing"
	"time"

	"go-crud-app/internal/utils"
)

func TestGenerateToken(t *testing.T) {
	config := utils.JWTConfig{
		SecretKey:       "test-secret-key",
		ExpirationHours: 24,
	}

	token, err := utils.GenerateToken(1, "testuser", "test@example.com", config)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated, but got empty string")
	}
}

func TestValidateToken(t *testing.T) {
	config := utils.JWTConfig{
		SecretKey:       "test-secret-key",
		ExpirationHours: 24,
	}

	// Generate a valid token
	token, err := utils.GenerateToken(1, "testuser", "test@example.com", config)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		secretKey   string
		shouldError bool
	}{
		{
			name:        "Valid token",
			token:       token,
			secretKey:   config.SecretKey,
			shouldError: false,
		},
		{
			name:        "Invalid secret key",
			token:       token,
			secretKey:   "wrong-secret-key",
			shouldError: true,
		},
		{
			name:        "Invalid token format",
			token:       "invalid.token.format",
			secretKey:   config.SecretKey,
			shouldError: true,
		},
		{
			name:        "Empty token",
			token:       "",
			secretKey:   config.SecretKey,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ValidateToken(tt.token, tt.secretKey)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if claims == nil {
					t.Error("Expected claims to be returned, but got nil")
				}
				if claims.UserID != 1 {
					t.Errorf("Expected UserID to be 1, but got %d", claims.UserID)
				}
				if claims.Username != "testuser" {
					t.Errorf("Expected Username to be 'testuser', but got %s", claims.Username)
				}
				if claims.Email != "test@example.com" {
					t.Errorf("Expected Email to be 'test@example.com', but got %s", claims.Email)
				}
			}
		})
	}
}

func TestExpiredToken(t *testing.T) {
	config := utils.JWTConfig{
		SecretKey:       "test-secret-key",
		ExpirationHours: -1, // Expired token
	}

	token, err := utils.GenerateToken(1, "testuser", "test@example.com", config)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait a moment to ensure token is expired
	time.Sleep(100 * time.Millisecond)

	_, err = utils.ValidateToken(token, config.SecretKey)
	if err == nil {
		t.Error("Expected error for expired token, but got none")
	}
}
