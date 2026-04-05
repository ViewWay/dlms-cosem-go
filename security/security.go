package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"fmt"
)

// Security suite definitions
const (
	SecuritySuite0 = 0 // No security
	SecuritySuite1 = 1 // AES-128-GCM
	SecuritySuite2 = 2 // AES-256-GCM
	SecuritySuite3 = 3 // SM4-GCM (Chinese standard)
	SecuritySuite4 = 4 // ECDSA + AES-128
	SecuritySuite5 = 5 // ECDSA + SM4
)

// SecurityControlField represents the 8-bit security control field.
type SecurityControlField struct {
	SecuritySuite int
	Authenticated bool
	Encrypted     bool
	BroadcastKey  bool
	Compressed    bool
}

// FromBytes parses security control from a byte.
func (s *SecurityControlField) FromBytes(b byte) {
	s.SecuritySuite = int(b & 0x0F)
	s.Authenticated = (b & 0x10) != 0
	s.Encrypted = (b & 0x20) != 0
	s.BroadcastKey = (b & 0x40) != 0
	s.Compressed = (b & 0x80) != 0
}

// ToBytes encodes security control to a byte.
func (s *SecurityControlField) ToBytes() byte {
	var b byte
	b |= byte(s.SecuritySuite & 0x0F)
	if s.Authenticated {
		b |= 0x10
	}
	if s.Encrypted {
		b |= 0x20
	}
	if s.BroadcastKey {
		b |= 0x40
	}
	if s.Compressed {
		b |= 0x80
	}
	return b
}

// SecurityProcessor provides cryptographic operations for DLMS/COSEM.
type SecurityProcessor struct {
	suite             int
	encryptionKey     []byte
	authenticationKey []byte
	systemTitle       []byte
}

// NewSecurityProcessor creates a new security processor.
func NewSecurityProcessor(suite int, encryptionKey, authenticationKey, systemTitle []byte) (*SecurityProcessor, error) {
	if err := validateSuite(suite); err != nil {
		return nil, err
	}
	if err := validateKey(suite, encryptionKey, "encryption"); err != nil {
		return nil, err
	}
	if authenticationKey != nil {
		if err := validateKey(suite, authenticationKey, "authentication"); err != nil {
			return nil, err
		}
	}
	if len(systemTitle) != 8 {
		return nil, fmt.Errorf("security: system title must be 8 bytes, got %d", len(systemTitle))
	}

	return &SecurityProcessor{
		suite:             suite,
		encryptionKey:     encryptionKey,
		authenticationKey: authenticationKey,
		systemTitle:       systemTitle,
	}, nil
}

func validateSuite(suite int) error {
	if suite < 0 || suite > 5 {
		return fmt.Errorf("security: invalid suite %d", suite)
	}
	return nil
}

func validateKey(suite int, key []byte, name string) error {
	// Suite0 has no encryption/authentication, nil keys are allowed
	if suite == SecuritySuite0 {
		return nil
	}
	if key == nil {
		return fmt.Errorf("security: %s key is nil", name)
	}
	var expectedLen int
	switch suite {
	case SecuritySuite1, SecuritySuite3, SecuritySuite4, SecuritySuite5:
		expectedLen = 16
	case SecuritySuite2:
		expectedLen = 32
	case SecuritySuite0:
		return nil
	}
	if len(key) != expectedLen {
		return fmt.Errorf("security: %s key must be %d bytes, got %d", name, expectedLen, len(key))
	}
	return nil
}

// MakeIV creates a 12-byte initialization vector from system title and invocation counter.
func (sp *SecurityProcessor) MakeIV(invocationCounter uint32) []byte {
	iv := make([]byte, 12)
	copy(iv[:8], sp.systemTitle)
	iv[8] = byte(invocationCounter >> 24)
	iv[9] = byte(invocationCounter >> 16)
	iv[10] = byte(invocationCounter >> 8)
	iv[11] = byte(invocationCounter)
	return iv
}

// AESEncryptGCM performs AES-GCM encryption.
func AESEncryptGCM(key, nonce, plaintext, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Seal(nil, nonce, plaintext, aad), nil
}

// AESDecryptGCM performs AES-GCM decryption.
func AESDecryptGCM(key, nonce, ciphertext, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, nonce, ciphertext, aad)
}

// SM4GCMEncrypt performs SM4-GCM encryption (GCM mode with SM4).
func SM4GCMEncrypt(key, nonce, plaintext, aad []byte) ([]byte, error) {
	// Use SM4 in CTR mode for GCM
	// For DLMS/COSEM, GMAC (authentication only) is common
	// Full GCM requires CTR + GHASH
	// Simplified: use SM4-ECB to generate keystream, XOR with plaintext
	// Then compute GMAC for authentication

	if len(nonce) != 12 {
		return nil, fmt.Errorf("sm4-gcm: nonce must be 12 bytes")
	}

	// Generate keystream using SM4-ECB
	keystream := make([]byte, len(plaintext))
	counter := make([]byte, 16)
	copy(counter, nonce)
	counter[15] = 1

	for i := 0; i < len(plaintext); i += 16 {
		encrypted := SM4Encrypt(key, counter)
		end := i + 16
		if end > len(plaintext) {
			end = len(plaintext)
		}
		copy(keystream[i:end], encrypted[:end-i])
		// Increment counter
		for j := 15; j >= 12; j-- {
			counter[j]++
			if counter[j] != 0 {
				break
			}
		}
	}

	// XOR plaintext with keystream
	ciphertext := make([]byte, len(plaintext))
	for i := range plaintext {
		ciphertext[i] = plaintext[i] ^ keystream[i]
	}

	// Compute GMAC tag
	tag, err := SM4GMAC(key, nonce, aad, ciphertext)
	if err != nil {
		return nil, err
	}

	return append(ciphertext, tag...), nil
}

// SM4GCMDecrypt performs SM4-GCM decryption.
func SM4GCMDecrypt(key, nonce, ciphertext, aad []byte) ([]byte, error) {
	if len(nonce) != 12 {
		return nil, fmt.Errorf("sm4-gcm: nonce must be 12 bytes")
	}
	if len(ciphertext) < 12 {
		return nil, fmt.Errorf("sm4-gcm: ciphertext too short")
	}

	// Split ciphertext and tag
	ct := ciphertext[:len(ciphertext)-12]
	tag := ciphertext[len(ciphertext)-12:]

	// Verify tag
	computedTag, err := SM4GMAC(key, nonce, aad, ct)
	if err != nil {
		return nil, err
	}
	if subtle.ConstantTimeCompare(tag, computedTag) != 1 {
		return nil, fmt.Errorf("sm4-gcm: authentication failed")
	}

	// Decrypt using CTR mode
	keystream := make([]byte, len(ct))
	counter := make([]byte, 16)
	copy(counter, nonce)
	counter[15] = 1

	for i := 0; i < len(ct); i += 16 {
		encrypted := SM4Encrypt(key, counter)
		end := i + 16
		if end > len(ct) {
			end = len(ct)
		}
		copy(keystream[i:end], encrypted[:end-i])
		for j := 15; j >= 12; j-- {
			counter[j]++
			if counter[j] != 0 {
				break
			}
		}
	}

	plaintext := make([]byte, len(ct))
	for i := range ct {
		plaintext[i] = ct[i] ^ keystream[i]
	}

	return plaintext, nil
}

// SM4GMAC computes SM4-based GMAC (Galois Message Authentication Code).
func SM4GMAC(key, nonce, aad, ciphertext []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("sm4-gmac: key must be 16 bytes")
	}
	if len(nonce) != 12 {
		return nil, fmt.Errorf("sm4-gmac: nonce must be 12 bytes")
	}

	// Simplified GMAC: SM4-CBC-MAC over aad || ciphertext
	// Build the data to authenticate
	var data []byte
	if aad != nil {
		data = append(data, aad...)
	}
	data = append(data, ciphertext...)

	if len(data) == 0 {
		// Return MAC of empty
		block := SM4Encrypt(key, make([]byte, 16))
		return block[:12], nil
	}

	// Pad to block size
	padLen := 16 - (len(data) % 16)
	if padLen == 16 {
		padLen = 0
	}
	padded := make([]byte, len(data)+padLen)
	copy(padded, data)

	// CBC-MAC
	prev := make([]byte, 16)
	copy(prev, nonce)
	// Pad nonce to 16 bytes if needed
	nonceBlock := make([]byte, 16)
	copy(nonceBlock, nonce)

	for i := 0; i < len(padded); i += 16 {
		xored := make([]byte, 16)
		for j := 0; j < 16; j++ {
			xored[j] = padded[i+j] ^ prev[j]
		}
		prev = SM4Encrypt(key, xored)
	}

	return prev[:12], nil
}

// Encrypt performs encryption based on the security suite.
func (sp *SecurityProcessor) Encrypt(invocationCounter uint32, plaintext, aad []byte) ([]byte, error) {
	switch sp.suite {
	case SecuritySuite0:
		return nil, fmt.Errorf("security: suite 0 does not support encryption")
	}
	iv := sp.MakeIV(invocationCounter)
	switch sp.suite {
	case SecuritySuite1, SecuritySuite4:
		return AESEncryptGCM(sp.encryptionKey, iv, plaintext, aad)
	case SecuritySuite2:
		return AESEncryptGCM(sp.encryptionKey, iv, plaintext, aad)
	case SecuritySuite3, SecuritySuite5:
		return SM4GCMEncrypt(sp.encryptionKey, iv, plaintext, aad)
	default:
		return nil, fmt.Errorf("security: suite %d does not support encryption", sp.suite)
	}
}

// Decrypt performs decryption based on the security suite.
func (sp *SecurityProcessor) Decrypt(invocationCounter uint32, ciphertext, aad []byte) ([]byte, error) {
	iv := sp.MakeIV(invocationCounter)
	switch sp.suite {
	case SecuritySuite1, SecuritySuite4:
		return AESDecryptGCM(sp.encryptionKey, iv, ciphertext, aad)
	case SecuritySuite2:
		return AESDecryptGCM(sp.encryptionKey, iv, ciphertext, aad)
	case SecuritySuite3, SecuritySuite5:
		return SM4GCMDecrypt(sp.encryptionKey, iv, ciphertext, aad)
	default:
		return nil, fmt.Errorf("security: suite %d does not support decryption", sp.suite)
	}
}

// Authenticate computes authentication tag.
func (sp *SecurityProcessor) Authenticate(invocationCounter uint32, data []byte) ([]byte, error) {
	iv := sp.MakeIV(invocationCounter)
	switch sp.suite {
	case SecuritySuite1, SecuritySuite4:
		return AESEncryptGCM(sp.authenticationKey, iv, nil, data)
	case SecuritySuite2:
		return AESEncryptGCM(sp.authenticationKey, iv, nil, data)
	case SecuritySuite3, SecuritySuite5:
		return SM4GMAC(sp.authenticationKey, iv, data, nil)
	default:
		return nil, fmt.Errorf("security: suite %d does not support authentication", sp.suite)
	}
}

// HLSISM performs High Level Security Initialization with SM-MAC.
func HLSISM(authKey []byte, challenge, systemTitle []byte, fCount int) []byte {
	// Calculate HMAC-like value using SM4
	data := make([]byte, 0, len(challenge)+len(systemTitle)+8)
	data = append(data, challenge...)
	data = append(data, systemTitle...)
	// Add frame counter (4 bytes) + direction (4 bytes)
	data = append(data, byte(fCount>>24), byte(fCount>>16), byte(fCount>>8), byte(fCount))
	data = append(data, 0, 0, 0, 1) // client to server

	// Build 12-byte IV from system title
	iv := make([]byte, 12)
	copy(iv, systemTitle)
	iv[8] = byte(fCount >> 24)
	iv[9] = byte(fCount >> 16)
	iv[10] = byte(fCount >> 8)
	iv[11] = byte(fCount)

	tag, _ := SM4GMAC(authKey, iv, nil, data)
	return tag
}
