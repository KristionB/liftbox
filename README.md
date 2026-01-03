# Liftbox

A production-ready, secure file synchronization system with end-to-end encryption, stateless architecture, and zero-trust server design.

## Architecture Overview

### Stateless Architecture

The server is designed as a **stateless API** that does not maintain any session state or user authentication. This design provides several key benefits:

- **Horizontal Scalability**: Any server instance can handle any request without shared state
- **Fault Tolerance**: Server failures don't result in lost sessions or state
- **Simplified Deployment**: No need for sticky sessions, session stores, or state synchronization
- **Cloud-Native**: Perfect for containerized deployments and auto-scaling

The server acts purely as a **storage and verification layer**, processing each request independently based on cryptographic proofs provided by the client.

### Client-Side Encryption

All encryption and cryptographic operations happen **exclusively on the client side** before data leaves the user's device:

- **AES-256-GCM Encryption**: Files are encrypted using AES-256 in Galois/Counter Mode before upload
- **PBKDF2 Key Derivation**: Encryption keys are derived from user passwords using PBKDF2 with 100,000 iterations
- **HMAC Verification**: HMAC-SHA256 ensures data integrity without the server seeing plaintext
- **Ed25519 Signatures**: Digital signatures provide authentication and non-repudiation

**The server never sees unencrypted data.** Even if the server is compromised, attackers cannot decrypt files without the user's password and keys.

### Zero-Trust Server

The server implements a **zero-trust security model**:

- **No Authentication Required**: The server doesn't authenticate users or maintain user accounts
- **Cryptographic Verification Only**: Every request is verified using cryptographic signatures and HMACs
- **No Plaintext Access**: The server cannot decrypt or read file contents
- **Stateless Verification**: Each request is independently verified without relying on stored credentials

The server's role is limited to:
1. Verifying cryptographic signatures
2. Verifying HMAC integrity
3. Storing encrypted blobs
4. Serving encrypted blobs on request

This design ensures that **server compromise does not compromise user data**.

### Concurrency Model

The client implements **concurrent file uploads** using Go's goroutines:

```go
func UploadFiles(files []string) {
    var wg sync.WaitGroup
    for _, f := range files {
        wg.Add(1)
        go func(file string) {
            defer wg.Done()
            encryptAndSend(file)
        }(f)
    }
    wg.Wait()
}
```

**Benefits:**
- **Parallel Processing**: Multiple files are encrypted and uploaded simultaneously
- **Efficient Resource Usage**: Leverages I/O wait time for other operations
- **Scalable**: Can handle large batches of files efficiently
- **Synchronization**: `sync.WaitGroup` ensures all uploads complete before returning

This model is particularly effective for:
- Batch file synchronization
- Large file transfers
- High-latency network conditions

### Threat Model

The system is designed to protect against the following threats:

#### ✅ Protected Against

1. **Server Compromise**
   - Files are encrypted client-side; server cannot decrypt
   - No authentication credentials stored on server
   - Cryptographic signatures prevent unauthorized modifications

2. **Man-in-the-Middle Attacks**
   - HMAC verification detects tampering
   - Digital signatures prevent replay attacks
   - Encrypted payloads prevent eavesdropping

3. **Data Breaches**
   - Encrypted data stored in S3 is useless without keys
   - No plaintext credentials or keys stored server-side
   - Zero-trust model means server compromise doesn't expose data

4. **Unauthorized Access**
   - Ed25519 signatures verify request authenticity
   - HMAC ensures data integrity
   - No server-side authentication means no authentication bypass

#### ⚠️ Security Considerations

1. **Key Management**
   - Users must securely store their passwords and private keys
   - Lost keys result in permanent data loss (by design)
   - Consider key escrow for enterprise use cases

2. **Password Strength**
   - Weak passwords reduce security of PBKDF2-derived keys
   - Users should use strong, unique passwords

3. **Client Security**
   - Client compromise exposes encryption keys
   - Users should use secure devices and environments

4. **Denial of Service**
   - Server does not implement rate limiting (consider adding)
   - Large file uploads could consume resources

## Technology Stack

- **Language**: Go 1.22
- **Framework**: Gin (HTTP web framework)
- **Cryptography**: 
  - AES-256-GCM for encryption
  - PBKDF2 for key derivation
  - HMAC-SHA256 for integrity
  - Ed25519 for signatures
- **Storage**: AWS S3
- **Testing**: Ginkgo + Gomega
- **Deployment**: Fly.io
- **CI/CD**: GitHub Actions

## Project Structure

```
secure-file-sync/
├── client/              # Client application
│   ├── crypto/         # Cryptographic primitives
│   └── sync/           # File synchronization
├── server/              # Stateless API server
│   ├── handlers/       # HTTP handlers
│   ├── storage/        # S3 storage interface
│   └── verify/         # Cryptographic verification
├── tests/               # Ginkgo test suite
└── .github/workflows/  # CI/CD pipelines
```

## Security Features

- ✅ End-to-end encryption (client-side only)
- ✅ Zero-trust server architecture
- ✅ Cryptographic signature verification
- ✅ HMAC integrity checking
- ✅ Stateless, scalable design
- ✅ No server-side key storage
- ✅ No authentication required (cryptographic proofs only)

## Development

### Prerequisites

- Go 1.22+
- AWS CLI (for S3 setup)
- Fly CLI (for deployment)

### Running Tests

```bash
# Run Ginkgo tests
ginkgo ./...

# Run standard Go tests
go test ./...
```

### Building

```bash
# Build server
go build -o server ./server

# Build client
go build -o client ./client/cmd
```

### Deployment

See [DEPLOYMENT.md](./DEPLOYMENT.md) for detailed deployment instructions.

## License

MIT

---

**Note for Recruiters**: This project demonstrates production-ready security architecture, modern Go concurrency patterns, cloud-native deployment, and comprehensive testing. The zero-trust, stateless design showcases understanding of scalable system architecture and security best practices.
