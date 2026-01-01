package tests

import (
	"testing"

	"secure-file-sync/client/crypto"
)

func TestKDF(t *testing.T) {
	password := "test-password"
	salt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	key1 := crypto.DeriveKey(password, salt)
	key2 := crypto.DeriveKey(password, salt)

	if len(key1) != crypto.KeySize {
		t.Errorf("Expected key size %d, got %d", crypto.KeySize, len(key1))
	}

	if len(key1) != len(key2) {
		t.Error("Key sizes don't match")
	}

	for i := range key1 {
		if key1[i] != key2[i] {
			t.Error("Derived keys don't match")
		}
	}

	// Test with different salt
	salt2, _ := crypto.GenerateSalt()
	key3 := crypto.DeriveKey(password, salt2)
	if len(key3) == len(key1) {
		// Keys should be different with different salts
		same := true
		for i := range key1 {
			if key1[i] != key3[i] {
				same = false
				break
			}
		}
		if same {
			t.Error("Keys should be different with different salts")
		}
	}
}

func TestAES(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := []byte("Hello, World! This is a test message.")

	encrypted, _, err := crypto.EncryptAES(key, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if len(encrypted) == 0 {
		t.Error("Encrypted data is empty")
	}

	decrypted, err := crypto.DecryptAES(key, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted text doesn't match. Expected: %s, Got: %s", plaintext, decrypted)
	}
}

func TestHMAC(t *testing.T) {
	key := []byte("test-key")
	data := []byte("test data")

	hmac1 := crypto.ComputeHMAC(key, data)
	hmac2 := crypto.ComputeHMAC(key, data)

	if len(hmac1) == 0 {
		t.Error("HMAC is empty")
	}

	if !crypto.VerifyHMAC(key, data, hmac1) {
		t.Error("HMAC verification failed")
	}

	if !crypto.VerifyHMAC(key, data, hmac2) {
		t.Error("HMAC verification failed for second computation")
	}

	// Test with wrong key
	wrongKey := []byte("wrong-key")
	if crypto.VerifyHMAC(wrongKey, data, hmac1) {
		t.Error("HMAC verification should fail with wrong key")
	}

	// Test with wrong data
	wrongData := []byte("wrong data")
	if crypto.VerifyHMAC(key, wrongData, hmac1) {
		t.Error("HMAC verification should fail with wrong data")
	}
}

func TestSigning(t *testing.T) {
	publicKey, privateKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	data := []byte("test data to sign")

	signature := crypto.SignData(privateKey, data)

	if len(signature) == 0 {
		t.Error("Signature is empty")
	}

	if !crypto.VerifySignature(publicKey, data, signature) {
		t.Error("Signature verification failed")
	}

	// Test with wrong data
	wrongData := []byte("wrong data")
	if crypto.VerifySignature(publicKey, wrongData, signature) {
		t.Error("Signature verification should fail with wrong data")
	}

	// Test hex encoding/decoding
	publicKeyHex := crypto.PublicKeyToHex(publicKey)
	privateKeyHex := crypto.PrivateKeyToHex(privateKey)

	decodedPublic, err := crypto.PublicKeyFromHex(publicKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	decodedPrivate, err := crypto.PrivateKeyFromHex(privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	if len(decodedPublic) != len(publicKey) {
		t.Error("Decoded public key length doesn't match")
	}

	if len(decodedPrivate) != len(privateKey) {
		t.Error("Decoded private key length doesn't match")
	}

	// Verify signature with decoded keys
	newSignature := crypto.SignData(decodedPrivate, data)
	if !crypto.VerifySignature(decodedPublic, data, newSignature) {
		t.Error("Signature verification failed with decoded keys")
	}
}

