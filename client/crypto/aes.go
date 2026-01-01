package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func Encrypt(data, key []byte) ([]byte, []byte) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nonce
}

func Decrypt(ciphertext, nonce, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	return gcm.Open(nil, nonce, ciphertext, nil)
}

