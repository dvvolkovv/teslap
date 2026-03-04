package ledger

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestValidateDoubleEntry_BalancedEntries(t *testing.T) {
	svc := &Service{}

	entries := []EntryLine{
		{
			AccountID: uuid.New(),
			Side:      Debit,
			Amount:    decimal.NewFromFloat(100.00),
			Currency:  "EUR",
		},
		{
			AccountID: uuid.New(),
			Side:      Credit,
			Amount:    decimal.NewFromFloat(100.00),
			Currency:  "EUR",
		},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err != nil {
		t.Errorf("expected balanced entries to be valid, got error: %v", err)
	}
}

func TestValidateDoubleEntry_UnbalancedEntries(t *testing.T) {
	svc := &Service{}

	entries := []EntryLine{
		{
			AccountID: uuid.New(),
			Side:      Debit,
			Amount:    decimal.NewFromFloat(100.00),
			Currency:  "EUR",
		},
		{
			AccountID: uuid.New(),
			Side:      Credit,
			Amount:    decimal.NewFromFloat(99.99),
			Currency:  "EUR",
		},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err == nil {
		t.Error("expected unbalanced entries to return error, got nil")
	}
}

func TestValidateDoubleEntry_MultiCurrency(t *testing.T) {
	svc := &Service{}

	// Multi-leg FX transaction: debit EUR, credit USD, both balanced per currency.
	entries := []EntryLine{
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.NewFromFloat(100.00), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.NewFromFloat(100.00), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.NewFromFloat(108.50), Currency: "USD"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.NewFromFloat(108.50), Currency: "USD"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err != nil {
		t.Errorf("expected multi-currency balanced entries to be valid, got error: %v", err)
	}
}

func TestValidateDoubleEntry_NegativeAmount(t *testing.T) {
	svc := &Service{}

	entries := []EntryLine{
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.NewFromFloat(-100.00), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.NewFromFloat(-100.00), Currency: "EUR"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err == nil {
		t.Error("expected negative amount to return error, got nil")
	}
}

func TestValidateDoubleEntry_ZeroAmount(t *testing.T) {
	svc := &Service{}

	entries := []EntryLine{
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.Zero, Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.Zero, Currency: "EUR"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err == nil {
		t.Error("expected zero amount to return error, got nil")
	}
}

func TestValidateDoubleEntry_SingleEntry(t *testing.T) {
	svc := &Service{}

	entries := []EntryLine{
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.NewFromFloat(100.00), Currency: "EUR"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err == nil {
		t.Error("expected single entry to return error, got nil")
	}
}

func TestValidateDoubleEntry_InvalidSide(t *testing.T) {
	svc := &Service{}

	entries := []EntryLine{
		{AccountID: uuid.New(), Side: "invalid", Amount: decimal.NewFromFloat(100.00), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.NewFromFloat(100.00), Currency: "EUR"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err == nil {
		t.Error("expected invalid side to return error, got nil")
	}
}

func TestValidateDoubleEntry_PrecisionHandling(t *testing.T) {
	svc := &Service{}

	// Test that decimal precision is maintained correctly.
	entries := []EntryLine{
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.RequireFromString("0.0001"), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.RequireFromString("0.0001"), Currency: "EUR"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err != nil {
		t.Errorf("expected precise decimal entries to be valid, got error: %v", err)
	}
}

func TestValidateDoubleEntry_ThreeWaySplit(t *testing.T) {
	svc := &Service{}

	// Payment with fee: debit customer 100.20, credit beneficiary 100.00, credit fee revenue 0.20.
	entries := []EntryLine{
		{AccountID: uuid.New(), Side: Debit, Amount: decimal.RequireFromString("100.20"), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.RequireFromString("100.00"), Currency: "EUR"},
		{AccountID: uuid.New(), Side: Credit, Amount: decimal.RequireFromString("0.20"), Currency: "EUR"},
	}

	err := svc.ValidateDoubleEntry(entries)
	if err != nil {
		t.Errorf("expected three-way split to be valid, got error: %v", err)
	}
}

func TestComputeEventChecksum(t *testing.T) {
	data := map[string]any{"type": "debit.posted", "amount": "100.00"}
	checksum1 := ComputeEventChecksum(data, "")
	if checksum1 == "" {
		t.Error("expected non-empty checksum")
	}

	// Same data with same previous checksum should produce same result.
	checksum2 := ComputeEventChecksum(data, "")
	if checksum1 != checksum2 {
		t.Error("expected deterministic checksum")
	}

	// Different previous checksum should produce different result.
	checksum3 := ComputeEventChecksum(data, "abc123")
	if checksum1 == checksum3 {
		t.Error("expected different checksum with different chain input")
	}
}
