package security

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestSM4EncryptDecrypt(t *testing.T) {
	key := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}
	plaintext := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}

	encrypted := SM4Encrypt(key, plaintext)
	if len(encrypted) != 16 {
		t.Errorf("len=%d", len(encrypted))
	}

	decrypted := SM4Decrypt(key, encrypted)
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("got %x, want %x", decrypted, plaintext)
	}
}

func TestSM4KnownVector(t *testing.T) {
	// GB/T 32907-2016 test vector
	key := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}
	plaintext := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}
	expected := []byte{0x68, 0x1E, 0xDF, 0x34, 0xD2, 0x06, 0x96, 0x5E, 0x86, 0xB3, 0xE9, 0x4F, 0x53, 0x6E, 0x42, 0x46}

	encrypted := SM4Encrypt(key, plaintext)
	if !bytes.Equal(encrypted, expected) {
		t.Errorf("got %x, want %x", encrypted, expected)
	}
}

func TestSM4CBCRoundtrip(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	plaintext := []byte("Hello, SM4 CBC mode!")

	encrypted := SM4EncryptCBC(key, iv, plaintext)
	decrypted, err := SM4DecryptCBC(key, iv, encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("got %s, want %s", string(decrypted), string(plaintext))
	}
}

func TestSM4CBC_BlockAligned(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	plaintext := make([]byte, 32) // exactly 2 blocks

	encrypted := SM4EncryptCBC(key, iv, plaintext)
	decrypted, err := SM4DecryptCBC(key, iv, encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Error("roundtrip failed")
	}
}

func TestSM4CBC_InvalidIV(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid IV")
		}
	}()
	SM4EncryptCBC(make([]byte, 16), []byte{1}, []byte("test"))
}

func TestSM4DecryptCBC_InvalidCiphertext(t *testing.T) {
	_, err := SM4DecryptCBC(make([]byte, 16), make([]byte, 16), []byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for invalid ciphertext length")
	}
}

func TestSM4DecryptCBC_InvalidIV(t *testing.T) {
	_, err := SM4DecryptCBC(make([]byte, 16), []byte{1}, make([]byte, 16))
	if err == nil {
		t.Error("expected error for invalid IV")
	}
}

func TestSM4ECBRoundtrip(t *testing.T) {
	key := make([]byte, 16)
	plaintext := []byte("ECB mode test")

	encrypted := SM4EncryptECB(key, plaintext)
	decrypted, err := SM4DecryptECB(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Error("roundtrip failed")
	}
}

func TestSM4DecryptECB_InvalidLength(t *testing.T) {
	_, err := SM4DecryptECB(make([]byte, 16), []byte{1, 2, 3})
	if err == nil {
		t.Error("expected error")
	}
}

func TestSM4Encrypt_PanicsOnShortKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	SM4Encrypt([]byte{1, 2}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
}

func TestSecurityControlField(t *testing.T) {
	tests := []struct {
		sc SecurityControlField
		b  byte
	}{
		{SecurityControlField{0, false, false, false, false}, 0x00},
		{SecurityControlField{1, true, false, false, false}, 0x11},
		{SecurityControlField{2, true, true, false, false}, 0x32},
		{SecurityControlField{3, true, true, true, false}, 0x73},
		{SecurityControlField{5, true, true, false, true}, 0xB5},
	}
	for _, tc := range tests {
		b := tc.sc.ToBytes()
		if b != tc.b {
			t.Errorf("ToBytes(%+v) = 0x%02x, want 0x%02x", tc.sc, b, tc.b)
		}
		var parsed SecurityControlField
		parsed.FromBytes(b)
		if parsed != tc.sc {
			t.Errorf("FromBytes(0x%02x) = %+v, want %+v", b, parsed, tc.sc)
		}
	}
}

func TestNewSecurityProcessor(t *testing.T) {
	key := make([]byte, 16)
	authKey := make([]byte, 16)
	title := make([]byte, 8)

	sp, err := NewSecurityProcessor(SecuritySuite1, key, authKey, title)
	if err != nil {
		t.Fatal(err)
	}
	if sp.suite != SecuritySuite1 {
		t.Error("wrong suite")
	}
}

func TestNewSecurityProcessor_InvalidSuite(t *testing.T) {
	_, err := NewSecurityProcessor(99, make([]byte, 16), nil, make([]byte, 8))
	if err == nil {
		t.Error("expected error")
	}
}

func TestNewSecurityProcessor_InvalidKeyLen(t *testing.T) {
	_, err := NewSecurityProcessor(SecuritySuite1, make([]byte, 8), nil, make([]byte, 8))
	if err == nil {
		t.Error("expected error")
	}
}

func TestNewSecurityProcessor_InvalidSystemTitle(t *testing.T) {
	_, err := NewSecurityProcessor(SecuritySuite1, make([]byte, 16), nil, make([]byte, 6))
	if err == nil {
		t.Error("expected error")
	}
}

func TestNewSecurityProcessor_Suite0(t *testing.T) {
	sp, err := NewSecurityProcessor(SecuritySuite0, nil, nil, make([]byte, 8))
	if err != nil {
		t.Fatal(err)
	}
	if sp.suite != 0 {
		t.Error("wrong suite")
	}
}

func TestNewSecurityProcessor_Suite2(t *testing.T) {
	key := make([]byte, 32)
	_, err := NewSecurityProcessor(SecuritySuite2, key, nil, make([]byte, 8))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewSecurityProcessor_Suite2_WrongKeyLen(t *testing.T) {
	_, err := NewSecurityProcessor(SecuritySuite2, make([]byte, 16), nil, make([]byte, 8))
	if err == nil {
		t.Error("expected error")
	}
}

func TestMakeIV(t *testing.T) {
	title := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	sp, _ := NewSecurityProcessor(SecuritySuite1, make([]byte, 16), nil, title)
	iv := sp.MakeIV(42)
	if len(iv) != 12 {
		t.Errorf("len=%d", len(iv))
	}
	if !bytes.Equal(iv[:8], title) {
		t.Error("system title mismatch")
	}
}

func TestAESEncryptDecryptGCM(t *testing.T) {
	key := make([]byte, 16)
	nonce := make([]byte, 12)
	plaintext := []byte("test data for AES-GCM")

	ct, err := AESEncryptGCM(key, nonce, plaintext, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(ct) <= len(plaintext) {
		t.Error("ciphertext should include tag")
	}

	pt, err := AESDecryptGCM(key, nonce, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, plaintext) {
		t.Errorf("got %s, want %s", pt, plaintext)
	}
}

func TestAESEncryptGCM_WithAAD(t *testing.T) {
	key := make([]byte, 16)
	nonce := make([]byte, 12)
	aad := []byte("additional data")

	ct, err := AESEncryptGCM(key, nonce, []byte("secret"), aad)
	if err != nil {
		t.Fatal(err)
	}

	// Wrong AAD should fail
	_, err = AESDecryptGCM(key, nonce, ct, []byte("wrong aad"))
	if err == nil {
		t.Error("should fail with wrong AAD")
	}

	// Correct AAD should succeed
	pt, err := AESDecryptGCM(key, nonce, ct, aad)
	if err != nil {
		t.Fatal(err)
	}
	if string(pt) != "secret" {
		t.Errorf("got %s", pt)
	}
}

func TestAESDecryptGCM_WrongKey(t *testing.T) {
	key1 := make([]byte, 16)
	key2 := make([]byte, 16)
	key2[0] = 0xFF
	nonce := make([]byte, 12)

	ct, _ := AESEncryptGCM(key1, nonce, []byte("test"), nil)
	_, err := AESDecryptGCM(key2, nonce, ct, nil)
	if err == nil {
		t.Error("should fail with wrong key")
	}
}

func TestAESGCM_256(t *testing.T) {
	key := make([]byte, 32)
	nonce := make([]byte, 12)

	ct, err := AESEncryptGCM(key, nonce, []byte("AES-256-GCM"), nil)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := AESDecryptGCM(key, nonce, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(pt) != "AES-256-GCM" {
		t.Errorf("got %s", pt)
	}
}

func TestSM4GMAC(t *testing.T) {
	key := make([]byte, 16)
	nonce := make([]byte, 12)

	tag, err := SM4GMAC(key, nonce, nil, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tag) != 12 {
		t.Errorf("tag len=%d", len(tag))
	}

	// Same input should produce same tag
	tag2, _ := SM4GMAC(key, nonce, nil, []byte("test"))
	if !bytes.Equal(tag, tag2) {
		t.Error("deterministic MAC expected")
	}
}

func TestSM4GMAC_InvalidKey(t *testing.T) {
	_, err := SM4GMAC([]byte{1}, make([]byte, 12), nil, []byte("test"))
	if err == nil {
		t.Error("expected error")
	}
}

func TestSM4GMAC_InvalidNonce(t *testing.T) {
	_, err := SM4GMAC(make([]byte, 16), []byte{1}, nil, []byte("test"))
	if err == nil {
		t.Error("expected error")
	}
}

func TestSM4GCMEncryptDecrypt(t *testing.T) {
	key := make([]byte, 16)
	nonce := make([]byte, 12)
	plaintext := []byte("SM4-GCM test data!")

	ct, err := SM4GCMEncrypt(key, nonce, plaintext, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(ct) < len(plaintext) {
		t.Error("ciphertext too short")
	}

	pt, err := SM4GCMDecrypt(key, nonce, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, plaintext) {
		t.Errorf("got %s, want %s", pt, plaintext)
	}
}

func TestSM4GCMDecrypt_WrongKey(t *testing.T) {
	key1 := make([]byte, 16)
	key2 := make([]byte, 16)
	key2[0] = 0xFF
	nonce := make([]byte, 12)

	ct, _ := SM4GCMEncrypt(key1, nonce, []byte("test"), nil)
	_, err := SM4GCMDecrypt(key2, nonce, ct, nil)
	if err == nil {
		t.Error("should fail")
	}
}

func TestSM4GCMEncrypt_InvalidNonce(t *testing.T) {
	_, err := SM4GCMEncrypt(make([]byte, 16), []byte{1}, []byte("test"), nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSM4GCMDecrypt_InvalidNonce(t *testing.T) {
	_, err := SM4GCMDecrypt(make([]byte, 16), []byte{1}, make([]byte, 16), nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSM4GCMDecrypt_TooShort(t *testing.T) {
	_, err := SM4GCMDecrypt(make([]byte, 16), make([]byte, 12), []byte{1}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSecurityProcessor_EncryptDecrypt_AES(t *testing.T) {
	key := make([]byte, 16)
	authKey := make([]byte, 16)
	title := make([]byte, 8)
	sp, _ := NewSecurityProcessor(SecuritySuite1, key, authKey, title)

	ct, err := sp.Encrypt(1, []byte("test"), nil)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := sp.Decrypt(1, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, []byte("test")) {
		t.Errorf("got %s", pt)
	}
}

func TestSecurityProcessor_EncryptDecrypt_SM4(t *testing.T) {
	key := make([]byte, 16)
	authKey := make([]byte, 16)
	title := make([]byte, 8)
	sp, _ := NewSecurityProcessor(SecuritySuite3, key, authKey, title)

	ct, err := sp.Encrypt(1, []byte("test sm4"), nil)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := sp.Decrypt(1, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, []byte("test sm4")) {
		t.Errorf("got %s", pt)
	}
}

func TestSecurityProcessor_Encrypt_Suite0(t *testing.T) {
	title := make([]byte, 8)
	sp, _ := NewSecurityProcessor(SecuritySuite0, nil, nil, title)
	_, err := sp.Encrypt(1, []byte("test"), nil)
	if err == nil {
		t.Error("expected error for suite 0")
	}
}

func TestSecurityProcessor_Authenticate(t *testing.T) {
	key := make([]byte, 16)
	authKey := make([]byte, 16)
	title := make([]byte, 8)
	sp, _ := NewSecurityProcessor(SecuritySuite1, key, authKey, title)

	tag, err := sp.Authenticate(1, []byte("data"))
	if err != nil {
		t.Fatal(err)
	}
	// Suite1 uses AES-GCM, returns ciphertext+tag (16 bytes for nil plaintext)
	if len(tag) != 16 {
		t.Errorf("tag len=%d, want 16", len(tag))
	}
}

func TestHLSISM(t *testing.T) {
	key := make([]byte, 16)
	challenge := make([]byte, 16)
	title := make([]byte, 8)
	tag := HLSISM(key, challenge, title, 1)
	if len(tag) != 12 {
		t.Errorf("tag len=%d", len(tag))
	}
}

func TestSM4EncryptCBC_Empty(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	ct := SM4EncryptCBC(key, iv, []byte{})
	if len(ct) != 16 {
		t.Errorf("len=%d", len(ct))
	}
}

func TestSM4EncryptCBC_SingleBlock(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	ct := SM4EncryptCBC(key, iv, make([]byte, 1))
	// 1 byte + 15 bytes PKCS7 padding = 16 bytes = 1 block
	if len(ct) != 16 {
		t.Errorf("len=%d, want 16", len(ct))
	}
}

func TestRotl32(t *testing.T) {
	if rotl32(0x80000000, 1) != 1 {
		t.Error("rotl32")
	}
	if rotl32(1, 31) != 0x80000000 {
		t.Error("rotl32")
	}
}

func TestSM4SboxTransform(t *testing.T) {
	if sm4SboxTransform(0x00) != 0xD6 {
		t.Errorf("got 0x%02x", sm4SboxTransform(0x00))
	}
	if sm4SboxTransform(0xFF) != 0x48 {
		t.Errorf("got 0x%02x", sm4SboxTransform(0xFF))
	}
}

func TestValidateSuite(t *testing.T) {
	tests := []int{-1, 6, 100}
	for _, s := range tests {
		if validateSuite(s) == nil {
			t.Errorf("suite %d should be invalid", s)
		}
	}
}

func TestValidateKey_Suite0(t *testing.T) {
	if validateKey(0, nil, "test") != nil {
		t.Error("suite 0 should accept nil key")
	}
}

func TestNewSecurityProcessor_AllSuites(t *testing.T) {
	for _, suite := range []int{0, 1, 2, 3, 4, 5} {
		var key, authKey []byte
		switch suite {
		case 0:
		case 2:
			key = make([]byte, 32)
			authKey = make([]byte, 32)
		default:
			key = make([]byte, 16)
			authKey = make([]byte, 16)
		}
		_, err := NewSecurityProcessor(suite, key, authKey, make([]byte, 8))
		if err != nil {
			t.Errorf("suite %d: %v", suite, err)
		}
	}
}

func TestSecurityProcessor_RandomData(t *testing.T) {
	key := make([]byte, 16)
	authKey := make([]byte, 16)
	title := make([]byte, 8)
	sp, _ := NewSecurityProcessor(SecuritySuite1, key, authKey, title)

	for i := 0; i < 10; i++ {
		plaintext := make([]byte, 100+i*13)
		rand.Read(plaintext)

		ct, err := sp.Encrypt(uint32(i), plaintext, nil)
		if err != nil {
			t.Fatal(err)
		}

		pt, err := sp.Decrypt(uint32(i), ct, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(pt, plaintext) {
			t.Errorf("roundtrip failed for size %d", len(plaintext))
		}
	}
}

func TestSecurityProcessor_AuthKeyNil(t *testing.T) {
	key := make([]byte, 16)
	title := make([]byte, 8)
	sp, err := NewSecurityProcessor(SecuritySuite1, key, nil, title)
	if err != nil {
		t.Fatal(err)
	}
	_ = sp
}
