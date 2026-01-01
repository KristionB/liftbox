package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
)

// GenerateKeyPair generates a new Ed25519 key pair for signing
func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(rand.Reader)
}

// SignData signs data using Ed25519 private key
func SignData(privateKey ed25519.PrivateKey, data []byte) []byte {
	return ed25519.Sign(privateKey, data)
}

// VerifySignature verifies a signature using Ed25519 public key
func VerifySignature(publicKey ed25519.PublicKey, data []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}

// PublicKeyToHex converts public key to hex string
func PublicKeyToHex(publicKey ed25519.PublicKey) string {
	return hex.EncodeToString(publicKey)
}

// PublicKeyFromHex converts hex string to public key
func PublicKeyFromHex(hexKey string) (ed25519.PublicKey, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(keyBytes), nil
}

// PrivateKeyToHex converts private key to hex string
func PrivateKeyToHex(privateKey ed25519.PrivateKey) string {
	return hex.EncodeToString(privateKey)
}

// PrivateKeyFromHex converts hex string to private key
func PrivateKeyFromHex(hexKey string) (ed25519.PrivateKey, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(keyBytes), nil
}

