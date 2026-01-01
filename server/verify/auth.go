package verify

import (
	"crypto/hmac"
	"crypto/ed25519"
	"encoding/hex"

	"github.com/KristionB/secure-file-sync/client/crypto"
)

// VerifySignature verifies the signature of an upload request
func VerifySignature(publicKeyHex, fileName string, fileData []byte, hmacHex, signatureHex string) bool {
	// Decode public key
	publicKey, err := crypto.PublicKeyFromHex(publicKeyHex)
	if err != nil {
		return false
	}

	// Decode signature
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	// Reconstruct signature payload (same as client)
	signaturePayload := append([]byte(fileName), fileData...)
	signaturePayload = append(signaturePayload, []byte(hmacHex)...)

	// Verify signature
	return crypto.VerifySignature(publicKey, signaturePayload, signature)
}

// VerifyHMAC verifies the HMAC of file data
func VerifyHMAC(key []byte, data []byte, hmacHex string) bool {
	expectedHMAC, err := hex.DecodeString(hmacHex)
	if err != nil {
		return false
	}
	computedHMAC := crypto.ComputeHMAC(data, key)
	return hmac.Equal(expectedHMAC, computedHMAC)
}

// GetPublicKeyFromHex is a helper to get public key from hex string
func GetPublicKeyFromHex(hexKey string) (ed25519.PublicKey, error) {
	return crypto.PublicKeyFromHex(hexKey)
}

