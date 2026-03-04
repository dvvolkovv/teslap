package account

import (
	"strings"
	"testing"
)

func TestGenerateIBAN_Format(t *testing.T) {
	iban := GenerateIBAN("123456789012")
	if !strings.HasPrefix(iban, "LT") {
		t.Errorf("IBAN should start with LT, got %s", iban)
	}
	if len(iban) != 20 {
		t.Errorf("Lithuanian IBAN should be 20 characters, got %d: %s", len(iban), iban)
	}
}

func TestGenerateIBAN_ValidCheckDigits(t *testing.T) {
	iban := GenerateIBAN("123456789012")
	if !ValidateIBAN(iban) {
		t.Errorf("Generated IBAN should pass validation: %s", iban)
	}
}

func TestValidateIBAN_Valid(t *testing.T) {
	// Generate and validate.
	iban := GenerateIBAN("000000000001")
	if !ValidateIBAN(iban) {
		t.Errorf("Expected valid IBAN, got invalid: %s", iban)
	}
}

func TestValidateIBAN_Invalid(t *testing.T) {
	if ValidateIBAN("LT00TESL000000000001") {
		t.Error("Expected invalid IBAN with check digits 00 to fail validation")
	}
}

func TestValidateIBAN_TooShort(t *testing.T) {
	if ValidateIBAN("LT") {
		t.Error("Expected short string to fail validation")
	}
}

func TestGenerateExternalID_Format(t *testing.T) {
	id := GenerateExternalID()
	if !strings.HasPrefix(id, "TP-") {
		t.Errorf("External ID should start with 'TP-', got %s", id)
	}
	if len(id) != 11 {
		t.Errorf("External ID should be 11 characters (TP-XXXXXXXX), got %d: %s", len(id), id)
	}
}

func TestGenerateIBAN_DifferentAccounts(t *testing.T) {
	iban1 := GenerateIBAN("000000000001")
	iban2 := GenerateIBAN("000000000002")
	if iban1 == iban2 {
		t.Error("Different account numbers should produce different IBANs")
	}
}

func TestConvertToNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A", "10"},
		{"B", "11"},
		{"Z", "35"},
		{"LT", "2129"},
		{"12AB", "121011"},
	}

	for _, tc := range tests {
		result := convertToNumeric(tc.input)
		if result != tc.expected {
			t.Errorf("convertToNumeric(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}
