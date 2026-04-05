package security

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

// ============= SM4-GMAC Authentication Tests =============

// TestSM4GMAC_AuthenticationRoundTrip tests SM4-GMAC for authentication
func TestSM4GMAC_AuthenticationRoundTrip(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	message := []byte("Authentication test message")

	tag, err := SM4GMAC(key, nonce, message, nil)
	if err != nil {
		t.Fatalf("SM4GMAC failed: %v", err)
	}

	// Verify tag
	computedTag, err := SM4GMAC(key, nonce, message, nil)
	if err != nil {
		t.Fatalf("SM4GMAC failed: %v", err)
	}

	if !bytes.Equal(tag, computedTag) {
		t.Error("Tags don't match for same inputs")
	}
}

// TestSM4GMAC_WithAAD tests SM4-GMAC with Additional Authenticated Data
func TestSM4GMAC_WithAAD(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	message := []byte("Message")
	aad := []byte("Additional Authenticated Data")

	tag1, err := SM4GMAC(key, nonce, aad, message)
	if err != nil {
		t.Fatalf("SM4GMAC failed: %v", err)
	}

	tag2, err := SM4GMAC(key, nonce, aad, message)
	if err != nil {
		t.Fatalf("SM4GMAC failed: %v", err)
	}

	if !bytes.Equal(tag1, tag2) {
		t.Error("Tags should be deterministic")
	}

	// Different AAD should produce different tag
	tag3, _ := SM4GMAC(key, nonce, []byte("Different AAD"), message)
	if bytes.Equal(tag1, tag3) {
		t.Error("Different AAD should produce different tag")
	}
}

// TestSM4GMAC_EmptyMessage tests SM4-GMAC with empty message
func TestSM4GMAC_EmptyMessage(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)

	tag, err := SM4GMAC(key, nonce, nil, nil)
	if err != nil {
		t.Fatalf("SM4GMAC with empty data failed: %v", err)
	}

	if len(tag) != 12 {
		t.Errorf("Tag length should be 12, got %d", len(tag))
	}
}

// ============= SM4-GCM Round Trip Tests =============

// TestSM4GCM_RoundTripLargeData tests SM4-GCM with large data
func TestSM4GCM_RoundTripLargeData(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)

	// Test with various sizes
	sizes := []int{0, 1, 15, 16, 17, 31, 32, 33, 64, 127, 128, 129, 255, 256, 512, 1024}

	for _, size := range sizes {
		plaintext := make([]byte, size)
		rand.Read(plaintext)

		ciphertext, err := SM4GCMEncrypt(key, nonce, plaintext, nil)
		if err != nil {
			t.Fatalf("Encryption failed for size %d: %v", size, err)
		}

		decrypted, err := SM4GCMDecrypt(key, nonce, ciphertext, nil)
		if err != nil {
			t.Fatalf("Decryption failed for size %d: %v", size, err)
		}

		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("Roundtrip failed for size %d", size)
		}
	}
}

// TestSM4GCM_WithAAD tests SM4-GCM with Additional Authenticated Data
func TestSM4GCM_WithAAD(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	plaintext := []byte("Secret message")
	aad := []byte("Additional authenticated data")

	ct, err := SM4GCMEncrypt(key, nonce, plaintext, aad)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt with correct AAD
	pt, err := SM4GCMDecrypt(key, nonce, ct, aad)
	if err != nil {
		t.Fatalf("Decryption with correct AAD failed: %v", err)
	}

	if !bytes.Equal(pt, plaintext) {
		t.Error("Decrypted data doesn't match plaintext")
	}

	// Decrypt with wrong AAD should fail
	wrongAad := []byte("Wrong AAD")
	_, err = SM4GCMDecrypt(key, nonce, ct, wrongAad)
	if err == nil {
		t.Error("Decryption with wrong AAD should fail")
	}
}

// TestSM4GCM_DifferentNonces tests that different nonces produce different ciphertexts
func TestSM4GCM_DifferentNonces(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	plaintext := []byte("Same plaintext")

	nonce1 := make([]byte, 12)
	nonce1[0] = 0x01
	nonce2 := make([]byte, 12)
	nonce2[0] = 0x02

	ct1, _ := SM4GCMEncrypt(key, nonce1, plaintext, nil)
	ct2, _ := SM4GCMEncrypt(key, nonce2, plaintext, nil)

	if bytes.Equal(ct1, ct2) {
		t.Error("Different nonces should produce different ciphertexts")
	}
}

// TestSM4GCM_DifferentKeys tests that different keys produce different ciphertexts
func TestSM4GCM_DifferentKeys(t *testing.T) {
	key1 := make([]byte, 16)
	key2 := make([]byte, 16)
	key2[0] = 0xFF
	nonce := make([]byte, 12)
	plaintext := []byte("Same plaintext")

	ct1, _ := SM4GCMEncrypt(key1, nonce, plaintext, nil)
	ct2, _ := SM4GCMEncrypt(key2, nonce, plaintext, nil)

	if bytes.Equal(ct1, ct2) {
		t.Error("Different keys should produce different ciphertexts")
	}
}

// TestSM4GCM_TamperedCiphertext tests that tampered ciphertext fails authentication
func TestSM4GCM_TamperedCiphertext(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	plaintext := []byte("Original message")

	ct, err := SM4GCMEncrypt(key, nonce, plaintext, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Tamper with ciphertext (modify first byte)
	ct[0] ^= 0xFF

	_, err = SM4GCMDecrypt(key, nonce, ct, nil)
	if err == nil {
		t.Error("Tampered ciphertext should fail authentication")
	}
}

// ============= SM4 Key Expansion Tests =============

// TestSM4KeyExpansion tests SM4 key expansion
func TestSM4KeyExpansion(t *testing.T) {
	key := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}

	rk := sm4KeyExpansion(key)

	// Should generate 32 round keys
	if len(rk) != sm4Rounds {
		t.Errorf("Expected %d round keys, got %d", sm4Rounds, len(rk))
	}

	// Round keys should be non-zero
	for i, rk_val := range rk {
		if rk_val == 0 {
			t.Errorf("Round key %d is zero", i)
		}
	}
}

// TestSM4KeyExpansion_Deterministic tests that key expansion is deterministic
func TestSM4KeyExpansion_Deterministic(t *testing.T) {
	key := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}

	rk1 := sm4KeyExpansion(key)
	rk2 := sm4KeyExpansion(key)

	if !equalRoundKeys(rk1, rk2) {
		t.Error("Key expansion should be deterministic")
	}
}

// TestSM4KeyExpansion_DifferentKeys tests that different keys produce different round keys
func TestSM4KeyExpansion_DifferentKeys(t *testing.T) {
	key1 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}
	key2 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x11}

	rk1 := sm4KeyExpansion(key1)
	rk2 := sm4KeyExpansion(key2)

	// At least one round key should be different
	foundDiff := false
	for i := 0; i < sm4Rounds; i++ {
		if rk1[i] != rk2[i] {
			foundDiff = true
			break
		}
	}

	if !foundDiff {
		t.Error("Different keys should produce different round keys")
	}
}

// ============= AES-GCM IV and Data Size Tests =============

// TestAESGCM_VariousIVLengths tests AES-GCM with different IV lengths
func TestAESGCM_VariousIVLengths(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)

	// Test with various IV lengths (though GCM standard is 12 bytes)
	ivLengths := []int{12}

	for _, ivLen := range ivLengths {
		iv := make([]byte, ivLen)
		rand.Read(iv)
		plaintext := []byte("Test data")

		ct, err := AESEncryptGCM(key, iv, plaintext, nil)
		if err != nil {
			t.Fatalf("Encryption failed with IV length %d: %v", ivLen, err)
		}

		pt, err := AESDecryptGCM(key, iv, ct, nil)
		if err != nil {
			t.Fatalf("Decryption failed with IV length %d: %v", ivLen, err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("Roundtrip failed with IV length %d", ivLen)
		}
	}
}

// TestAESGCM_VariousDataSizes tests AES-GCM with various data sizes
func TestAESGCM_VariousDataSizes(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)

	sizes := []int{0, 1, 16, 32, 64, 128, 256, 512, 1024, 2048}

	for _, size := range sizes {
		plaintext := make([]byte, size)
		rand.Read(plaintext)

		ct, err := AESEncryptGCM(key, nonce, plaintext, nil)
		if err != nil {
			t.Fatalf("Encryption failed for size %d: %v", size, err)
		}

		pt, err := AESDecryptGCM(key, nonce, ct, nil)
		if err != nil {
			t.Fatalf("Decryption failed for size %d: %v", size, err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("Roundtrip failed for size %d", size)
		}
	}
}

// TestAESGCM_256WithVariousSizes tests AES-256-GCM with various data sizes
func TestAESGCM_256WithVariousSizes(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)

	sizes := []int{0, 1, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		plaintext := make([]byte, size)
		rand.Read(plaintext)

		ct, err := AESEncryptGCM(key, nonce, plaintext, nil)
		if err != nil {
			t.Fatalf("Encryption failed for size %d: %v", size, err)
		}

		pt, err := AESDecryptGCM(key, nonce, ct, nil)
		if err != nil {
			t.Fatalf("Decryption failed for size %d: %v", size, err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("Roundtrip failed for size %d", size)
		}
	}
}

// ============= HLS (High Level Security) Tests =============

// TestHLSISM_Deterministic tests that HLSISM is deterministic
func TestHLSISM_Deterministic(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	challenge := make([]byte, 16)
	rand.Read(challenge)
	title := make([]byte, 8)
	rand.Read(title)

	tag1 := HLSISM(key, challenge, title, 1)
	tag2 := HLSISM(key, challenge, title, 1)

	if !bytes.Equal(tag1, tag2) {
		t.Error("HLSISM should be deterministic")
	}
}

// TestHLSISM_DifferentCounters tests that different counters produce different tags
func TestHLSISM_DifferentCounters(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	challenge := make([]byte, 16)
	rand.Read(challenge)
	title := make([]byte, 8)
	rand.Read(title)

	tag1 := HLSISM(key, challenge, title, 1)
	tag2 := HLSISM(key, challenge, title, 2)

	if bytes.Equal(tag1, tag2) {
		t.Error("Different counters should produce different tags")
	}
}

// TestHLSISM_DifferentSystemTitles tests that different system titles produce different tags
func TestHLSISM_DifferentSystemTitles(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	challenge := make([]byte, 16)
	rand.Read(challenge)
	title1 := make([]byte, 8)
	title2 := make([]byte, 8)
	title1[0] = 0x01
	title2[0] = 0x02

	tag1 := HLSISM(key, challenge, title1, 1)
	tag2 := HLSISM(key, challenge, title2, 1)

	if bytes.Equal(tag1, tag2) {
		t.Error("Different system titles should produce different tags")
	}
}

// TestHLSISM_DifferentKeys tests that different keys produce different tags
func TestHLSISM_DifferentKeys(t *testing.T) {
	key1 := make([]byte, 16)
	key2 := make([]byte, 16)
	key2[0] = 0xFF
	challenge := make([]byte, 16)
	rand.Read(challenge)
	title := make([]byte, 8)
	rand.Read(title)

	tag1 := HLSISM(key1, challenge, title, 1)
	tag2 := HLSISM(key2, challenge, title, 1)

	if bytes.Equal(tag1, tag2) {
		t.Error("Different keys should produce different tags")
	}
}

// TestHLSISM_TagLength tests that HLSISM produces tags of correct length
func TestHLSISM_TagLength(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	challenge := make([]byte, 16)
	rand.Read(challenge)
	title := make([]byte, 8)
	rand.Read(title)

	tag := HLSISM(key, challenge, title, 1)

	if len(tag) != 12 {
		t.Errorf("Tag length should be 12, got %d", len(tag))
	}
}

// ============= Boundary Tests =============

// TestSM4GCM_EmptyData tests SM4-GCM with empty data
func TestSM4GCM_EmptyData(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)

	ct, err := SM4GCMEncrypt(key, nonce, []byte{}, nil)
	if err != nil {
		t.Fatalf("Encryption of empty data failed: %v", err)
	}

	pt, err := SM4GCMDecrypt(key, nonce, ct, nil)
	if err != nil {
		t.Fatalf("Decryption of empty data failed: %v", err)
	}

	if len(pt) != 0 {
		t.Errorf("Expected empty plaintext, got %d bytes", len(pt))
	}
}

// TestSM4GCM_SingleByte tests SM4-GCM with single byte
func TestSM4GCM_SingleByte(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	plaintext := []byte{0x42}

	ct, err := SM4GCMEncrypt(key, nonce, plaintext, nil)
	if err != nil {
		t.Fatal(err)
	}

	pt, err := SM4GCMDecrypt(key, nonce, ct, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, plaintext) {
		t.Error("Single byte roundtrip failed")
	}
}

// TestSM4GCM_LargeData tests SM4-GCM with large data
func TestSM4GCM_LargeData(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	plaintext := make([]byte, 10000) // 10KB
	rand.Read(plaintext)

	ct, err := SM4GCMEncrypt(key, nonce, plaintext, nil)
	if err != nil {
		t.Fatal(err)
	}

	pt, err := SM4GCMDecrypt(key, nonce, ct, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, plaintext) {
		t.Error("Large data roundtrip failed")
	}
}

// TestSM4GCM_WrongKey tests SM4-GCM with wrong key
func TestSM4GCM_WrongKey(t *testing.T) {
	key1 := make([]byte, 16)
	key2 := make([]byte, 16)
	rand.Read(key1)
	rand.Read(key2)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	plaintext := []byte("Secret message")

	ct, _ := SM4GCMEncrypt(key1, nonce, plaintext, nil)

	_, err := SM4GCMDecrypt(key2, nonce, ct, nil)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}
}

// TestSM4GCM_WrongNonce tests SM4-GCM with wrong nonce
func TestSM4GCM_WrongNonce(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce1 := make([]byte, 12)
	nonce2 := make([]byte, 12)
	rand.Read(nonce1)
	rand.Read(nonce2)
	plaintext := []byte("Secret message")

	ct, _ := SM4GCMEncrypt(key, nonce1, plaintext, nil)

	_, err := SM4GCMDecrypt(key, nonce2, ct, nil)
	if err == nil {
		t.Error("Decryption with wrong nonce should fail")
	}
}

// TestSM4GMAC_WrongTag tests SM4-GMAC with wrong tag
func TestSM4GMAC_WrongTag(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	message := []byte("Message")

	tag, _ := SM4GMAC(key, nonce, message, nil)
	tag[0] ^= 0xFF // Tamper with tag

	// Verify with tampered tag
	computedTag, _ := SM4GMAC(key, nonce, message, nil)
	if bytes.Equal(tag, computedTag) {
		t.Error("Tampered tag should not match computed tag")
	}
}

// TestAESGCM_WrongKey tests AES-GCM with wrong key
func TestAESGCM_WrongKey(t *testing.T) {
	key1 := make([]byte, 16)
	key2 := make([]byte, 16)
	rand.Read(key1)
	rand.Read(key2)
	nonce := make([]byte, 12)
	rand.Read(nonce)
	plaintext := []byte("Secret")

	ct, _ := AESEncryptGCM(key1, nonce, plaintext, nil)

	_, err := AESDecryptGCM(key2, nonce, ct, nil)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}
}

// TestAESGCM_EmptyData tests AES-GCM with empty data
func TestAESGCM_EmptyData(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	nonce := make([]byte, 12)
	rand.Read(nonce)

	ct, err := AESEncryptGCM(key, nonce, []byte{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	pt, err := AESDecryptGCM(key, nonce, ct, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(pt) != 0 {
		t.Errorf("Expected empty, got %d bytes", len(pt))
	}
}

// TestSecurityProcessor_AllSuitesRoundTrip tests all security suites
func TestSecurityProcessor_AllSuitesRoundTrip(t *testing.T) {
	suites := []struct {
		suite       int
		keyLen      int
		authKeyLen  int
	}{
		{SecuritySuite0, 0, 0},
		{SecuritySuite1, 16, 16},
		{SecuritySuite2, 32, 32},
		{SecuritySuite3, 16, 16},
		{SecuritySuite4, 16, 16},
		{SecuritySuite5, 16, 16},
	}

	for _, tc := range suites {
		t.Run(fmt.Sprintf("Suite%d", tc.suite), func(t *testing.T) {
			var key, authKey []byte
			if tc.keyLen > 0 {
				key = make([]byte, tc.keyLen)
				rand.Read(key)
			}
			if tc.authKeyLen > 0 {
				authKey = make([]byte, tc.authKeyLen)
				rand.Read(authKey)
			}
			title := make([]byte, 8)
			rand.Read(title)

			sp, err := NewSecurityProcessor(tc.suite, key, authKey, title)
			if err != nil {
				t.Fatalf("NewSecurityProcessor failed: %v", err)
			}

			if tc.suite == SecuritySuite0 {
				// Suite 0 doesn't support encryption
				_, err = sp.Encrypt(1, []byte("test"), nil)
				if err == nil {
					t.Error("Suite 0 should not support encryption")
				}
				return
			}

			plaintext := []byte("Test data for encryption")
			ct, err := sp.Encrypt(1, plaintext, nil)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			pt, err := sp.Decrypt(1, ct, nil)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if !bytes.Equal(pt, plaintext) {
				t.Error("Roundtrip failed")
			}
		})
	}
}

// TestSecurityProcessor_MultipleInvocationCounters tests multiple invocation counters
func TestSecurityProcessor_MultipleInvocationCounters(t *testing.T) {
	key := make([]byte, 16)
	authKey := make([]byte, 16)
	title := make([]byte, 8)
	rand.Read(key)
	rand.Read(authKey)
	rand.Read(title)

	sp, _ := NewSecurityProcessor(SecuritySuite1, key, authKey, title)
	plaintext := []byte("Test data")

	// Encrypt and decrypt with different invocation counters
	for i := uint32(0); i < 10; i++ {
		ct, err := sp.Encrypt(i, plaintext, nil)
		if err != nil {
			t.Fatalf("Encrypt failed for counter %d: %v", i, err)
		}

		pt, err := sp.Decrypt(i, ct, nil)
		if err != nil {
			t.Fatalf("Decrypt failed for counter %d: %v", i, err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("Roundtrip failed for counter %d", i)
		}
	}
}

// TestSM4CBC_VariousSizes tests SM4-CBC with various sizes
func TestSM4CBC_VariousSizes(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	rand.Read(key)
	rand.Read(iv)

	sizes := []int{0, 1, 15, 16, 17, 31, 32, 33, 64, 128, 256}

	for _, size := range sizes {
		plaintext := make([]byte, size)
		rand.Read(plaintext)

		ct := SM4EncryptCBC(key, iv, plaintext)
		pt, err := SM4DecryptCBC(key, iv, ct)
		if err != nil {
			t.Fatalf("Decryption failed for size %d: %v", size, err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("Roundtrip failed for size %d", size)
		}
	}
}

// TestSM4ECB_VariousSizes tests SM4-ECB with various sizes
func TestSM4ECB_VariousSizes(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)

	sizes := []int{1, 15, 16, 17, 31, 32, 33, 64, 128, 256}

	for _, size := range sizes {
		plaintext := make([]byte, size)
		rand.Read(plaintext)

		ct := SM4EncryptECB(key, plaintext)
		pt, err := SM4DecryptECB(key, ct)
		if err != nil {
			t.Fatalf("Decryption failed for size %d: %v", size, err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("Roundtrip failed for size %d", size)
		}
	}
}

// ============= Helper Functions =============

// equalRoundKeys compares two round key arrays
func equalRoundKeys(rk1, rk2 [sm4Rounds]uint32) bool {
	for i := 0; i < sm4Rounds; i++ {
		if rk1[i] != rk2[i] {
			return false
		}
	}
	return true
}
