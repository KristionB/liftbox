package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 32
	KeySize    = 32
	Iterations = 100000
)

// DeriveKey derives a cryptographic key from a password using PBKDF2
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)
}

// GenerateSalt generates a random salt for key derivation
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// SaltFromHex converts a hex-encoded salt string to bytes
func SaltFromHex(hexSalt string) ([]byte, error) {
	return hex.DecodeString(hexSalt)
}

// SaltToHex converts a salt byte array to a hex-encoded string
func SaltToHex(salt []byte) string {
	return hex.EncodeToString(salt)
}

