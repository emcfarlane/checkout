package checkout

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parsePAN tries to clean a PAN string to all numeric digits and reports if
// it successfully does so.
func parsePAN(cc string) (string, error) {
	var rs []rune
	for i, r := range cc {
		switch r {
		case ' ', '-', '_', ',':
			continue
		}
		if r < '0' || r > '9' {
			return "", fmt.Errorf("unexpected char %q", cc[i:i+1])
		}

		// Append all digits
		rs = append(rs, r)
	}

	// Check length of PAN is between 12-19 digits.
	if l := len(rs); l < 12 || l > 19 {
		return "", fmt.Errorf("invalid length %d", l)
	}

	return string(rs), nil
}

// checkLuhn validates a PAN string conforms to the luhn algorithm.
// Copies implementation from Wikipedia.
func checkLuhn(purportedCC string) bool {
	if len(purportedCC) == 0 {
		return false
	}
	nDigits := len(purportedCC)
	sum, err := strconv.Atoi(purportedCC[nDigits-1:])
	if err != nil {
		return false
	}
	parity := nDigits % 2
	for i := 0; i < nDigits-1; i++ {
		digit, err := strconv.Atoi(purportedCC[i : i+1])
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit = digit * 2
		}
		if digit > 9 {
			digit -= 9
		}
		sum += digit
	}
	return (sum % 10) == 0
}

// currencies is a subset of supported ISO 4217 currency codes.
var currencies = map[string]bool{
	"gbp": true,
	"usd": true,
	"eur": true,
}

// parseCurrency formats a currency string to lowercase and reports
// whether it is a supported currency.
func parseCurrency(currency string) (string, error) {
	c := strings.ToLower(currency)
	if ok := currencies[c]; !ok {
		return "", fmt.Errorf("invalid currency code %q", currency)
	}
	return c, nil
}

// checkCVV is a simple numeric check.
func checkCVV(cvv string) bool {
	if len(cvv) == 0 {
		return false
	}
	for _, r := range cvv {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// parseExpiry as a time.Time.
func parseExpiry(expYear, expMonth uint32) (time.Time, error) {
	if expYear > 3000 {
		return time.Time{}, fmt.Errorf("invalid expiry year %v", expYear)
	}
	if expMonth > 12 || expMonth == 0 {
		return time.Time{}, fmt.Errorf("invalid expiry month %v", expMonth)
	}
	m := time.Month(expMonth)
	return time.Date(int(expYear), m, 0, 0, 0, 0, 0, time.UTC), nil
}
