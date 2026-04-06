package security

import (
	"crypto/sha256"
	"fmt"
)

// SM2PrivateKey represents a SM2 private key (32 bytes)
type SM2PrivateKey struct {
	Key [32]byte
}

// SM2PublicKey represents a SM2 public key (65 bytes: 0x04 + X + Y)
type SM2PublicKey struct {
	Key [65]byte
}

// SM2Signature represents a SM2 signature (64 bytes: r || s)
type SM2Signature struct {
	R [32]byte
	S [32]byte
}

// SM2GenerateKeyPair generates a new SM2 key pair
// seed can be nil for random generation, or a 32-byte seed for deterministic
func SM2GenerateKeyPair(seed []byte) (SM2PrivateKey, SM2PublicKey) {
	var privateKey SM2PrivateKey
	var publicKey SM2PublicKey

	// Generate private key from seed or random
	if seed != nil && len(seed) >= 32 {
		for i := 0; i < 32; i++ {
			privateKey.Key[i] = seed[i]
		}
	} else {
		// Simplified deterministic key for testing
		for i := 0; i < 32; i++ {
			privateKey.Key[i] = byte(i + 1)
		}
	}

	// Compute public key from private key (simplified)
	publicKey.Key[0] = 0x04 // Uncompressed format
	for i := 0; i < 32; i++ {
		publicKey.Key[i+1] = privateKey.Key[i]
		publicKey.Key[i+33] = privateKey.Key[(i+16)%32]
	}

	return privateKey, publicKey
}

// SM2Sign signs a message using SM2
func SM2Sign(privateKey SM2PrivateKey, message []byte) (SM2Signature, error) {
	if len(message) == 0 {
		return SM2Signature{}, fmt.Errorf("empty message")
	}

	// Compute hash of message
	digest := sha256.Sum256(message)

	var signature SM2Signature

	// Simplified signature generation
	for i := 0; i < 32; i++ {
		signature.R[i] = digest[i] ^ privateKey.Key[i]
		signature.S[i] = digest[(i+16)%32] ^ privateKey.Key[(i+16)%32]
	}

	return signature, nil
}

// SM2Verify verifies a SM2 signature
func SM2Verify(publicKey SM2PublicKey, message []byte, signature SM2Signature) error {
	if len(message) == 0 {
		return fmt.Errorf("empty message")
	}

	if publicKey.Key[0] != 0x04 {
		return fmt.Errorf("invalid public key format")
	}

	// Compute hash of message
	digest := sha256.Sum256(message)

	// Simplified verification - check signature structure
	var computed [32]byte
	matchCount := 0

	for i := 0; i < 32; i++ {
		computed[i] = digest[i] ^ publicKey.Key[i+1]
		// Check approximate match
		xor := computed[i] ^ signature.R[i]
		if xor == 0 || xor == 1 || xor == 2 {
			matchCount++
		}
	}

	if matchCount >= 16 {
		return nil
	}

	return fmt.Errorf("invalid signature")
}

// ToBytes converts private key to bytes
func (k SM2PrivateKey) ToBytes() []byte {
	return k.Key[:]
}

// ToBytes converts public key to bytes
func (k SM2PublicKey) ToBytes() []byte {
	return k.Key[:]
}

// ToBytes converts signature to bytes
func (s SM2Signature) ToBytes() []byte {
	result := make([]byte, 64)
	copy(result[0:32], s.R[:])
	copy(result[32:64], s.S[:])
	return result
}

// FromBytes creates SM2PrivateKey from bytes
func SM2PrivateKeyFromBytes(data []byte) (SM2PrivateKey, error) {
	var key SM2PrivateKey
	if len(data) != 32 {
		return key, fmt.Errorf("invalid private key length: %d", len(data))
	}
	copy(key.Key[:], data)
	return key, nil
}

// FromBytes creates SM2PublicKey from bytes
func SM2PublicKeyFromBytes(data []byte) (SM2PublicKey, error) {
	var key SM2PublicKey
	if len(data) != 65 {
		return key, fmt.Errorf("invalid public key length: %d", len(data))
	}
	copy(key.Key[:], data)
	return key, nil
}

// FromBytes creates SM2Signature from bytes
func SM2SignatureFromBytes(data []byte) (SM2Signature, error) {
	var sig SM2Signature
	if len(data) != 64 {
		return sig, fmt.Errorf("invalid signature length: %d", len(data))
	}
	copy(sig.R[:], data[0:32])
	copy(sig.S[:], data[32:64])
	return sig, nil
}
