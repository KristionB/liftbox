package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"secure-file-sync/server/storage"
	"secure-file-sync/server/verify"
)

// UploadHandler handles file upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON request
	var req struct {
		FileName  string `json:"file_name"`
		FileData  []byte `json:"file_data"`
		HMAC      string `json:"hmac"`
		Signature string `json:"signature"`
		PublicKey string `json:"public_key"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Verify signature
	if !verify.VerifySignature(req.PublicKey, req.FileName, req.FileData, req.HMAC, req.Signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Store file
	if err := storage.StoreFile(req.FileName, req.FileData, req.HMAC, req.PublicKey); err != nil {
		http.Error(w, "Failed to store file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "File uploaded successfully",
	})
}

