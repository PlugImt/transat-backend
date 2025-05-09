package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// SecureRandom fills the provided byte slice with cryptographically secure random bytes
func SecureRandom(b []byte) (int, error) {
	return rand.Read(b)
}

// GenerateSecureRandomNumber generates a cryptographically secure random number in the range [min, max]
func GenerateSecureRandomNumber(min, max int64) (int64, error) {
	if min >= max {
		return 0, fmt.Errorf("min must be less than max")
	}

	// Create a range
	range_ := big.NewInt(max - min + 1)

	// Generate random number within range
	n, err := rand.Int(rand.Reader, range_)
	if err != nil {
		return 0, err
	}

	// Add min to shift to the desired range
	return n.Int64() + min, nil
}
