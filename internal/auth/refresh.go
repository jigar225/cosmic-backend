package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateRefreshToken returns a cryptographically random string suitable for refresh tokens.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashRefreshToken returns SHA256 hash of the token (hex-encoded) for storage.
func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// ValidateAndHashRefreshToken returns the hash of the raw token, or error if empty.
func ValidateAndHashRefreshToken(raw string) (hash string, err error) {
	if raw == "" {
		return "", fmt.Errorf("refresh token is required")
	}
	return HashRefreshToken(raw), nil
}
