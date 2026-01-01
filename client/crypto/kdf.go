package crypto

import (
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key(
		[]byte(password),
		salt,
		100_000,
		32,
		sha256.New,
	)
}

