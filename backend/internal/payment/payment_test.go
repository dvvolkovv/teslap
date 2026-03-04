package payment

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestValidateIBAN_Valid(t *testing.T) {
	validIBANs := []string{
		"DE89370400440532013000",
		"GB29NWBK60161331926819",
		"FR7630006000011234567890189",
		"NL91ABNA0417164300",
	}
	for _, iban := range validIBANs {
		if err := validateIBAN(iban); err != nil {
			t.Errorf("expected valid IBAN %s to pass, got error: %v", iban, err)
		}
	}
}

func TestValidateIBAN_Invalid(t *testing.T) {
	invalidIBANs := []string{
		"INVALID",
		"12ABCD1234567890",  // starts with digits
		"DE",                // too short
		"AB12",              // too short
		"AB1X123456789012345", // non-digit check digits
	}
	for _, iban := range invalidIBANs {
		if err := validateIBAN(iban); err == nil {
			t.Errorf("expected invalid IBAN %s to fail validation, got nil", iban)
		}
	}
}

func TestGetFXQuote_ValidPair(t *testing.T) {
	svc := &Service{}

	quote, err := svc.GetFXQuote(context.Background(), "EUR", "USD", "100")
	if err != nil {
		t.Fatalf("expected valid FX quote, got error: %v", err)
	}
	if quote == nil {
		t.Fatal("expected non-nil quote")
	}
	if quote.FromCurrency != "EUR" || quote.ToCurrency != "USD" {
		t.Errorf("expected EUR/USD, got %s/%s", quote.FromCurrency, quote.ToCurrency)
	}
	expectedRate := decimal.NewFromFloat(1.08)
	if !quote.Rate.Equal(expectedRate) {
		t.Errorf("expected rate 1.08, got %s", quote.Rate.String())
	}
	expectedConverted := decimal.NewFromFloat(108.00)
	if !quote.ConvertedAmount.Equal(expectedConverted) {
		t.Errorf("expected converted amount 108.00, got %s", quote.ConvertedAmount.String())
	}
}

func TestGetFXQuote_UnsupportedCurrency(t *testing.T) {
	svc := &Service{}

	_, err := svc.GetFXQuote(context.Background(), "XYZ", "USD", "100")
	if err == nil {
		t.Error("expected error for unsupported currency, got nil")
	}
}

func TestGetFXQuote_NegativeAmount(t *testing.T) {
	svc := &Service{}

	_, err := svc.GetFXQuote(context.Background(), "EUR", "USD", "-50")
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestGetFXQuote_ZeroAmount(t *testing.T) {
	svc := &Service{}

	_, err := svc.GetFXQuote(context.Background(), "EUR", "USD", "0")
	if err == nil {
		t.Error("expected error for zero amount, got nil")
	}
}

func TestGetFXQuote_ExpiryIs30Seconds(t *testing.T) {
	svc := &Service{}

	before := time.Now().UTC()
	quote, err := svc.GetFXQuote(context.Background(), "EUR", "USD", "100")
	after := time.Now().UTC()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	minExpiry := before.Add(29 * time.Second)
	maxExpiry := after.Add(31 * time.Second)

	if quote.ExpiresAt.Before(minExpiry) || quote.ExpiresAt.After(maxExpiry) {
		t.Errorf("quote expiry %v not within expected range [%v, %v]",
			quote.ExpiresAt, minExpiry, maxExpiry)
	}
}

func TestStringOrNil_NonEmpty(t *testing.T) {
	s := "hello"
	result := stringOrNil(s)
	if result == nil {
		t.Error("expected non-nil for non-empty string")
	}
	if *result != s {
		t.Errorf("expected %q, got %q", s, *result)
	}
}

func TestStringOrNil_Empty(t *testing.T) {
	result := stringOrNil("")
	if result != nil {
		t.Error("expected nil for empty string")
	}
}

func TestComputeNextExecution_Daily(t *testing.T) {
	now := time.Now().UTC()
	next := computeNextExecution(ScheduleTypeDaily, now)
	expected := now.Add(24 * time.Hour)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestComputeNextExecution_Weekly(t *testing.T) {
	now := time.Now().UTC()
	next := computeNextExecution(ScheduleTypeWeekly, now)
	expected := now.Add(7 * 24 * time.Hour)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestComputeNextExecution_Monthly(t *testing.T) {
	now := time.Now().UTC()
	next := computeNextExecution(ScheduleTypeMonthly, now)
	expected := now.AddDate(0, 1, 0)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestInternalPayment_SameAccountValidation(t *testing.T) {
	svc := &Service{repo: nil, producer: nil, logger: nil}

	req := &InternalPaymentRequest{
		SenderAccountID:    "550e8400-e29b-41d4-a716-446655440000",
		RecipientAccountID: "550e8400-e29b-41d4-a716-446655440000", // same!
		Amount:             "100.00",
		Currency:           "EUR",
	}

	_, err := svc.CreateInternalPayment(context.Background(), req)
	if err == nil {
		t.Error("expected error for same sender/recipient accounts, got nil")
	}
}

func TestInternalPayment_NegativeAmount(t *testing.T) {
	svc := &Service{}

	req := &InternalPaymentRequest{
		SenderAccountID:    "550e8400-e29b-41d4-a716-446655440000",
		RecipientAccountID: "550e8400-e29b-41d4-a716-446655440001",
		Amount:             "-100.00",
		Currency:           "EUR",
	}

	_, err := svc.CreateInternalPayment(context.Background(), req)
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestSEPAPayment_InvalidIBAN(t *testing.T) {
	svc := &Service{}

	req := &SEPAPaymentRequest{
		SenderAccountID: "550e8400-e29b-41d4-a716-446655440000",
		RecipientIBAN:   "INVALID_IBAN",
		RecipientName:   "Test Person",
		Amount:          "50.00",
		Currency:        "EUR",
	}

	_, err := svc.CreateSEPAPayment(context.Background(), req)
	if err == nil {
		t.Error("expected error for invalid IBAN, got nil")
	}
}
