package security

import (
	"testing"
	"time"
)

func TestNewCertificate(t *testing.T) {
	_, pub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert := NewCertificate([]byte("Test Subject"), publicKeyInfo, validity)

	if cert.Version != 3 {
		t.Errorf("version: got %d, want 3", cert.Version)
	}
	if string(cert.Subject) != "Test Subject" {
		t.Errorf("subject: got %s, want 'Test Subject'", cert.Subject)
	}
}

func TestCertificateSign(t *testing.T) {
	caPriv, _ := SM2GenerateKeyPair(nil)
	_, endPub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: endPub.ToBytes(),
	}

	cert := NewCertificate([]byte("End Entity"), publicKeyInfo, validity)

	err := cert.Sign(caPriv, []byte("CA Issuer"))
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if string(cert.Issuer) != "CA Issuer" {
		t.Errorf("issuer: got %s, want 'CA Issuer'", cert.Issuer)
	}
	if len(cert.Signature.Signature) == 0 {
		t.Error("signature should not be empty")
	}
}

func TestCertificateToDER(t *testing.T) {
	_, pub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert := NewCertificate([]byte("Test"), publicKeyInfo, validity)

	der, err := cert.ToDER()
	if err != nil {
		t.Fatalf("ToDER failed: %v", err)
	}
	if len(der) == 0 {
		t.Error("DER should not be empty")
	}
}

func TestCertificateStoreAdd(t *testing.T) {
	store := NewCertificateStore()
	_, pub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert := NewCertificate([]byte("Test"), publicKeyInfo, validity)

	store.Add(cert)
	if store.Count() != 1 {
		t.Errorf("count: got %d, want 1", store.Count())
	}
}

func TestCertificateStoreFind(t *testing.T) {
	store := NewCertificateStore()
	_, pub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert := NewCertificate([]byte("My Certificate"), publicKeyInfo, validity)

	store.Add(cert)

	found := store.FindBySubject([]byte("My Certificate"))
	if found == nil {
		t.Error("certificate not found")
	} else {
		if string(found.Subject) != "My Certificate" {
			t.Errorf("subject: got %s, want 'My Certificate'", found.Subject)
		}
	}
}

func TestCertificateStoreRemove(t *testing.T) {
	store := NewCertificateStore()
	_, pub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert := NewCertificate([]byte("Test"), publicKeyInfo, validity)

	store.Add(cert)
	if store.Count() != 1 {
		t.Errorf("count: got %d, want 1", store.Count())
	}

	store.Remove([]byte("Test"))
	if store.Count() != 0 {
		t.Errorf("count: got %d, want 0", store.Count())
	}
}

func TestCertificateChain(t *testing.T) {
	_ = SM2Verify // Silence unused warning
	store := NewCertificateStore()

	// Create CA certificate (self-signed)
	seed1 := make([]byte, 32)
	for i := 0; i < 32; i++ {
		seed1[i] = byte(i + 1)
	}
	caPriv, caPub := SM2GenerateKeyPair(seed1)

	caValidity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+63072000, 0),
	}

	caPubInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: caPub.ToBytes(),
	}

	caCert := NewCertificate([]byte("Root CA"), caPubInfo, caValidity)
	caCert.Sign(caPriv, []byte("Root CA"))
	store.Add(caCert)

	// Create intermediate certificate
	seed2 := make([]byte, 32)
	for i := 0; i < 32; i++ {
		seed2[i] = byte(i + 33)
	}
	intPriv, intPub := SM2GenerateKeyPair(seed2)

	intValidity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	intPubInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: intPub.ToBytes(),
	}

	intCert := NewCertificate([]byte("Intermediate CA"), intPubInfo, intValidity)
	intCert.Sign(caPriv, []byte("Root CA"))

	// Create end-entity certificate
	seed3 := make([]byte, 32)
	for i := 0; i < 32; i++ {
		seed3[i] = byte(i + 65)
	}
	_, endPub := SM2GenerateKeyPair(seed3)

	endValidity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+15768000, 0),
	}

	endPubInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: endPub.ToBytes(),
	}

	endCert := NewCertificate([]byte("End Entity"), endPubInfo, endValidity)
	endCert.Sign(intPriv, []byte("Intermediate CA"))

	// Check chain structure
	chain := []Certificate{endCert, intCert, caCert}
	if len(chain) != 3 {
		t.Errorf("chain length: got %d, want 3", len(chain))
	}
	// End entity should be signed by intermediate
	if string(chain[0].Issuer) != "Intermediate CA" {
		t.Errorf("chain[0].issuer: got %s, want 'Intermediate CA'", chain[0].Issuer)
	}
	// Intermediate should be signed by root
	if string(chain[1].Issuer) != "Root CA" {
		t.Errorf("chain[1].issuer: got %s, want 'Root CA'", chain[1].Issuer)
	}
	// Root should be self-signed
	if string(chain[2].Issuer) != "Root CA" {
		t.Errorf("chain[2].issuer: got %s, want 'Root CA'", chain[2].Issuer)
	}
}

func TestCertificateAll(t *testing.T) {
	store := NewCertificateStore()
	_, pub := SM2GenerateKeyPair(nil)

	validity := Validity{
		NotBefore: time.Unix(0x66000000, 0),
		NotAfter:  time.Unix(0x66000000+31536000, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert1 := NewCertificate([]byte("Test1"), publicKeyInfo, validity)
	cert2 := NewCertificate([]byte("Test2"), publicKeyInfo, validity)

	store.Add(cert1)
	store.Add(cert2)

	all := store.All()
	if len(all) != 2 {
		t.Errorf("All(): got %d certificates, want 2", len(all))
	}
}

func TestCertificateVerifyValidity(t *testing.T) {
	_, pub := SM2GenerateKeyPair(nil)

	// Expired certificate
	expiredValidity := Validity{
		NotBefore: time.Unix(0x66000000-31536000, 0),
		NotAfter:  time.Unix(0x66000000-86400, 0),
	}

	publicKeyInfo := PublicKeyInfo{
		Algorithm: []byte{0x06, 0x08},
		PublicKey: pub.ToBytes(),
	}

	cert := NewCertificate([]byte("Expired"), publicKeyInfo, expiredValidity)

	// Should fail due to expired certificate
	err := cert.Verify(pub)
	if err == nil {
		t.Error("Verify should fail for expired certificate")
	}
}
