package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the bcrypt cost (10 = 2^10 rounds).
const DefaultCost = 12

// HashPassword returns a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ComparePassword returns nil if password matches hash.
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
