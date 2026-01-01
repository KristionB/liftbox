package handlers

import (
	"encoding/json"
	"net/http"

	"secure-file-sync/server/storage"
)

// DownloadHandler handles file download requests
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Retrieve file from storage
	fileData, hmac, publicKey, err := storage.RetrieveFile(fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Return file data with metadata
	response := map[string]interface{}{
		"file_data":  fileData,
		"hmac":       hmac,
		"public_key": publicKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

