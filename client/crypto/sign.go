package crypto

import "crypto/ed25519"

func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	return publicKey, privateKey
}

