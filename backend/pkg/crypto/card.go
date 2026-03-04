package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// GenerateCardNumber generates a valid 16-digit card number with the given BIN prefix,
// filling random digits and computing the Luhn check digit.
func GenerateCardNumber(bin string) (string, error) {
	// Total digits = 16, last digit is Luhn check digit.
	totalLen := 16
	// We need (totalLen - 1 - len(bin)) random digits.
	randomLen := totalLen - 1 - len(bin)
	if randomLen < 0 {
		return "", fmt.Errorf("BIN too long")
	}
	digits := bin
	for i := 0; i < randomLen; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate digit: %w", err)
		}
		digits += n.String()
	}
	// Compute Luhn check digit.
	check := luhnCheckDigit(digits)
	return fmt.Sprintf("%s%d", digits, check), nil
}

// luhnCheckDigit computes the Luhn check digit for a partial card number (without check digit).
func luhnCheckDigit(partial string) int {
	sum := 0
	// The partial number will have an even number of digits after adding check digit.
	// We process from right to left; the rightmost existing digit is at position 2
	// (1-indexed from the right) when the check digit is added.
	nDigits := len(partial)
	for i, ch := range partial {
		d := int(ch - '0')
		// Position from right (1-indexed after adding check digit).
		pos := nDigits - i
		if pos%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return (10 - (sum % 10)) % 10
}

// ValidateLuhn validates a complete card number using the Luhn algorithm.
func ValidateLuhn(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	if len(number) < 13 || len(number) > 19 {
		return false
	}
	sum := 0
	nDigits := len(number)
	parity := nDigits % 2
	for i, ch := range number {
		d := int(ch - '0')
		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return sum%10 == 0
}

// MaskCardNumber returns a display-friendly masked card number.
// Input: "5425001234561234"
// Output: "5425 **** **** 1234"
// Returns the original string unchanged if it is not exactly 16 characters.
func MaskCardNumber(number string) string {
	if len(number) != 16 {
		return number
	}
	return fmt.Sprintf("%s **** **** %s", number[:4], number[12:])
}

// GenerateCVV returns a cryptographically random 3-digit CVV string (e.g. "042", "891").
// The range is [0, 999], zero-padded to 3 digits.
func GenerateCVV() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "", fmt.Errorf("generate CVV: %w", err)
	}
	return fmt.Sprintf("%03d", n.Int64()), nil
}

// HashCVV returns a SHA-256 hash of the CVV string, hex-encoded.
// This hash is stored in the database — the plaintext CVV is never stored.
func HashCVV(cvv string) string {
	h := sha256.Sum256([]byte(cvv))
	return hex.EncodeToString(h[:])
}

// EncryptCardNumber encrypts the full card number for at-rest storage using AES-256-GCM.
// A random 12-byte nonce is generated per encryption; the output is hex-encoded nonce+ciphertext.
// If key is nil or empty, the number is hex-encoded without encryption (dev mode only).
// In production, supply a 32-byte AES-256 key managed by a KMS.
func EncryptCardNumber(number string, key []byte) (string, error) {
	if len(key) == 0 {
		// Dev mode: just hex-encode without encryption.
		return hex.EncodeToString([]byte(number)), nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	// Seal appends ciphertext (with GCM tag) to nonce, producing nonce+ciphertext+tag.
	ciphertext := gcm.Seal(nonce, nonce, []byte(number), nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptCardNumber decrypts a card number previously encrypted with EncryptCardNumber.
// If key is nil or empty, the hex-encoded string is decoded and returned as-is (dev mode).
func DecryptCardNumber(encrypted string, key []byte) (string, error) {
	if len(key) == 0 {
		// Dev mode: reverse the hex-encoding.
		b, err := hex.DecodeString(encrypted)
		if err != nil {
			return "", fmt.Errorf("decode: %w", err)
		}
		return string(b), nil
	}
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("decode hex: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return string(plaintext), nil
}
