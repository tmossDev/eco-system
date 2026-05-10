package utils_test

import (
	"testing"

	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

func TestGenerateJwt(t *testing.T) {
	secretKey := "secret"
	issuer := "1"
	token, _, err := utils.GenerateJwt(issuer, secretKey)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Ensure the generated token is not empty
	if token == "" {
		t.Error("Generated token is empty")
	}

}

func TestParseJwt(t *testing.T) {
	secretKey := "secret"
	issuer := "me"
	// Generate a token for testing
	testToken, _, err := utils.GenerateJwt(issuer, secretKey)
	if err != nil {
		t.Fatalf("Error generating test JWT: %v", err)
	}

	// Test parsing the generated token
	parsedIssuer, err := utils.ParseJwt(testToken, secretKey)
	if err != nil {
		t.Fatalf("Error parsing JWT: %v", err)
	}

	// Ensure the parsed issuer matches the expected issuer
	if parsedIssuer != issuer {
		t.Errorf("Expected issuer %s, but got %s", issuer, parsedIssuer)
	}
}
