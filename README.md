# Secure File Sync

A secure file synchronization system with client-server architecture, featuring end-to-end encryption, digital signatures, and HMAC verification.

## Features

- **End-to-End Encryption**: Files are encrypted using AES-256-GCM before upload
- **Key Derivation**: PBKDF2 with SHA-256 for secure key derivation from passwords
- **Digital Signatures**: Ed25519 signatures for authentication and integrity
- **HMAC Verification**: HMAC-SHA256 for data integrity verification
- **Secure Storage**: Server-side storage with metadata management

## Architecture

```
client/          - Client application
  cmd/           - Main client entry point
  crypto/        - Cryptographic utilities (KDF, AES, HMAC, Signing)
  sync/          - File synchronization logic

server/          - Server application
  handlers/      - HTTP request handlers (upload, download)
  verify/        - Signature and HMAC verification
  storage/       - File storage backend (currently filesystem, S3-ready)

tests/           - Unit tests
```

## Building

```bash
# Build server
go build -o bin/server ./server/main.go

# Build client
go build -o bin/client ./client/cmd/main.go
```

## Usage

### Server

Start the server:

```bash
./bin/server -port 8080
```

Or using Docker:

```bash
docker build -t secure-file-sync .
docker run -p 8080:8080 secure-file-sync
```

### Client

Upload a file:

```bash
./bin/client \
  -server http://localhost:8080 \
  -file /path/to/file.txt \
  -password "your-secure-password" \
  -salt "hex-encoded-salt" \
  -private-key "hex-encoded-private-key" \
  -public-key "hex-encoded-public-key"
```

If you don't provide salt or keys, they will be generated automatically. Save the generated values for future use.

### Download

Download a file (using curl or similar):

```bash
curl "http://localhost:8080/download?file=filename.txt"
```

## Security Features

1. **Encryption**: Files are encrypted with AES-256-GCM before transmission
2. **Key Derivation**: PBKDF2 with 100,000 iterations and SHA-256
3. **Signatures**: Ed25519 signatures ensure file authenticity
4. **HMAC**: HMAC-SHA256 verifies data integrity
5. **Salt**: Random salt for each key derivation prevents rainbow table attacks

## Development

Run tests:

```bash
go test ./tests/...
```

## Storage Backend

The current implementation uses filesystem storage. To use S3, implement the S3 functions in `server/storage/s3.go` using the AWS SDK.

## License

MIT

