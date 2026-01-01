package main

import (
	"crypto/ed25519"
	"flag"
	"fmt"
	"os"

	"secure-file-sync/client/crypto"
	"secure-file-sync/client/sync"
)

func main() {
	var (
		serverURL      = flag.String("server", "http://localhost:8080", "Server URL")
		filePath     = flag.String("file", "", "Path to file to upload")
		password     = flag.String("password", "", "Password for key derivation")
		saltHex      = flag.String("salt", "", "Salt for key derivation (hex encoded)")
		privateKeyHex = flag.String("private-key", "", "Private key for signing (hex encoded)")
		publicKeyHex  = flag.String("public-key", "", "Public key for signing (hex encoded)")
	)
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintf(os.Stderr, "Error: -file is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *password == "" {
		fmt.Fprintf(os.Stderr, "Error: -password is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Derive encryption key
	var salt []byte
	var err error
	if *saltHex != "" {
		salt, err = crypto.SaltFromHex(*saltHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid salt: %v\n", err)
			os.Exit(1)
		}
	} else {
		salt, err = crypto.GenerateSalt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to generate salt: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated salt: %s\n", crypto.SaltToHex(salt))
	}

	encryptionKey := crypto.DeriveKey(*password, salt)

	// Get or generate signing keys
	var privateKey ed25519.PrivateKey
	var publicKey ed25519.PublicKey

	if *privateKeyHex != "" && *publicKeyHex != "" {
		privateKey, err = crypto.PrivateKeyFromHex(*privateKeyHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid private key: %v\n", err)
			os.Exit(1)
		}
		publicKey, err = crypto.PublicKeyFromHex(*publicKeyHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid public key: %v\n", err)
			os.Exit(1)
		}
	} else {
		publicKey, privateKey, err = crypto.GenerateKeyPair()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to generate key pair: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated public key: %s\n", crypto.PublicKeyToHex(publicKey))
		fmt.Printf("Generated private key: %s\n", crypto.PrivateKeyToHex(privateKey))
	}

	// Upload file
	fmt.Printf("Uploading file: %s\n", *filePath)
	err = sync.UploadFile(*serverURL, *filePath, encryptionKey, privateKey, publicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: upload failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File uploaded successfully!")
}

