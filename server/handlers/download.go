package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/KristionB/secure-file-sync/server/storage"
)

// Download handles file download requests
func Download(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	// Retrieve file from storage
	fileData, hmac, publicKey, err := storage.RetrieveFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Return file data with metadata
	c.JSON(http.StatusOK, gin.H{
		"file_data":  fileData,
		"hmac":       hmac,
		"public_key": publicKey,
	})
}

