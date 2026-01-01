package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// In-memory storage (can be replaced with actual S3 implementation)
type fileStorage struct {
	files map[string]*storedFile
	mu    sync.RWMutex
}

type storedFile struct {
	Data      []byte
	HMAC      string
	PublicKey string
}

var storage = &fileStorage{
	files: make(map[string]*storedFile),
}

// StoreFile stores a file in storage
func StoreFile(fileName string, fileData []byte, hmac string, publicKey string) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.files[fileName] = &storedFile{
		Data:      fileData,
		HMAC:      hmac,
		PublicKey: publicKey,
	}

	// Optionally persist to disk
	return persistToDisk(fileName, fileData, hmac, publicKey)
}

// RetrieveFile retrieves a file from storage
func RetrieveFile(fileName string) ([]byte, string, string, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	file, exists := storage.files[fileName]
	if !exists {
		// Try to load from disk
		return loadFromDisk(fileName)
	}

	return file.Data, file.HMAC, file.PublicKey, nil
}

// persistToDisk persists file to local filesystem (for development)
func persistToDisk(fileName string, fileData []byte, hmac string, publicKey string) error {
	storageDir := "storage"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return err
	}

	// Store metadata
	metaPath := filepath.Join(storageDir, fileName+".meta")
	meta := fmt.Sprintf("hmac:%s\npublic_key:%s\n", hmac, publicKey)
	return os.WriteFile(metaPath, []byte(meta), 0644)
}

// loadFromDisk loads file from local filesystem
func loadFromDisk(fileName string) ([]byte, string, string, error) {
	storageDir := "storage"
	filePath := filepath.Join(storageDir, fileName)
	metaPath := filepath.Join(storageDir, fileName+".meta")

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", "", err
	}

	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, "", "", err
	}

	// Parse metadata (simple format)
	var hmac, publicKey string
	lines := strings.Split(string(metaData), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "hmac:") {
			hmac = strings.TrimPrefix(line, "hmac:")
		} else if strings.HasPrefix(line, "public_key:") {
			publicKey = strings.TrimPrefix(line, "public_key:")
		}
	}

	// Store in memory for future access
	storage.mu.Lock()
	storage.files[fileName] = &storedFile{
		Data:      fileData,
		HMAC:      hmac,
		PublicKey: publicKey,
	}
	storage.mu.Unlock()

	return fileData, hmac, publicKey, nil
}

// TODO: Replace with actual S3 implementation
// Example S3 implementation would use AWS SDK:
// func StoreFileS3(fileName string, fileData []byte, hmac string, publicKey string) error {
//     sess := session.Must(session.NewSession())
//     svc := s3.New(sess)
//     _, err := svc.PutObject(&s3.PutObjectInput{
//         Bucket: aws.String("your-bucket"),
//         Key:    aws.String(fileName),
//         Body:   bytes.NewReader(fileData),
//     })
//     return err
// }

