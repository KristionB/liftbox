package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ComputeHMAC computes HMAC-SHA256 of the data using the provided key
func ComputeHMAC(key []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// VerifyHMAC verifies that the provided HMAC matches the computed HMAC
func VerifyHMAC(key []byte, data []byte, providedMAC []byte) bool {
	expectedMAC := ComputeHMAC(key, data)
	return hmac.Equal(expectedMAC, providedMAC)
}

// HMACToHex converts HMAC bytes to hex string
func HMACToHex(mac []byte) string {
	return hex.EncodeToString(mac)
}

// HMACFromHex converts hex string to HMAC bytes
func HMACFromHex(hexMAC string) ([]byte, error) {
	return hex.DecodeString(hexMAC)
}

