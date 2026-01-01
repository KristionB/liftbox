package sync

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"secure-file-sync/client/crypto"
)

// UploadRequest represents the request payload for file upload
type UploadRequest struct {
	FileName   string `json:"file_name"`
	FileData   []byte `json:"file_data"`
	HMAC       string `json:"hmac"`
	Signature  string `json:"signature"`
	PublicKey  string `json:"public_key"`
}

// UploadFile uploads an encrypted file to the server
func UploadFile(
	serverURL string,
	filePath string,
	encryptionKey []byte,
	privateKey ed25519.PrivateKey,
	publicKey ed25519.PublicKey,
) error {
	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Encrypt file data
	encryptedData, _, err := crypto.EncryptAES(encryptionKey, fileData)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Compute HMAC of encrypted data
	hmac := crypto.ComputeHMAC(encryptionKey, encryptedData)
	hmacHex := crypto.HMACToHex(hmac)

	// Create signature payload (file name + encrypted data + HMAC)
	signaturePayload := append([]byte(filepath.Base(filePath)), encryptedData...)
	signaturePayload = append(signaturePayload, []byte(hmacHex)...)

	// Sign the payload
	signature := crypto.SignData(privateKey, signaturePayload)
	signatureHex := crypto.HMACToHex(signature)

	// Create upload request
	req := UploadRequest{
		FileName:  filepath.Base(filePath),
		FileData:  encryptedData,
		HMAC:      hmacHex,
		Signature: signatureHex,
		PublicKey: crypto.PublicKeyToHex(publicKey),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send HTTP POST request
	resp, err := http.Post(serverURL+"/upload", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

