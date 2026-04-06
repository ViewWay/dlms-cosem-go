package security

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Certificate represents an X.509 v3 certificate
type Certificate struct {
	Version       uint8
	SerialNumber  [16]byte
	Issuer        []byte
	Subject       []byte
	Validity      Validity
	PublicKey     PublicKeyInfo
	Signature     CertificateSignature
}

// Validity represents the validity period of a certificate
type Validity struct {
	NotBefore time.Time
	NotAfter  time.Time
}

// PublicKeyInfo represents public key information
type PublicKeyInfo struct {
	Algorithm  []byte
	PublicKey  []byte
}

// CertificateSignature represents a certificate signature
type CertificateSignature struct {
	Algorithm  []byte
	Signature  []byte
}

// CertificateStore manages multiple certificates
type CertificateStore struct {
	certificates []Certificate
}

// NewCertificate creates a new certificate template
func NewCertificate(subject []byte, publicKey PublicKeyInfo, validity Validity) Certificate {
	return Certificate{
		Version:      3,
		SerialNumber: [16]byte{},
		Issuer:       subject,
		Subject:      subject,
		Validity:     validity,
		PublicKey:    publicKey,
		Signature:    CertificateSignature{},
	}
}

// ToDER serializes the certificate to DER format (simplified)
func (c *Certificate) ToDER() ([]byte, error) {
	der := make([]byte, 0, 256)

	// Version (3) with tag
	der = append(der, 0xA0, 0x03, 0x02, 0x01, 0x02)

	// Serial number
	der = append(der, 0x02, 0x10) // 16 bytes
	der = append(der, c.SerialNumber[:]...)

	// Issuer (simplified)
	der = append(der, 0x04, byte(len(c.Issuer)))
	der = append(der, c.Issuer...)

	// Validity
	der = append(der, 0x30, 0x1E)
	der = append(der, c.encodeTime(c.Validity.NotBefore)...)
	der = append(der, c.encodeTime(c.Validity.NotAfter)...)

	// Subject
	der = append(der, 0x04, byte(len(c.Subject)))
	der = append(der, c.Subject...)

	// Public key info
	der = append(der, 0x30, 0x59)
	der = append(der, 0x04, byte(len(c.PublicKey.Algorithm)))
	der = append(der, c.PublicKey.Algorithm...)
	der = append(der, 0x04, byte(len(c.PublicKey.PublicKey)))
	der = append(der, c.PublicKey.PublicKey...)

	return der, nil
}

// Sign signs the certificate with the CA private key
func (c *Certificate) Sign(caPrivateKey SM2PrivateKey, issuer []byte) error {
	c.Issuer = issuer

	// Compute hash of certificate data
	certData, err := c.ToDER()
	if err != nil {
		return err
	}
	digest := sha256.Sum256(certData)

	// Sign the digest
	signature, err := SM2Sign(caPrivateKey, digest[:])
	if err != nil {
		return err
	}

	// Set signature
	c.Signature.Algorithm = []byte{0x06, 0x08} // OID for SM2
	c.Signature.Signature = signature.ToBytes()

	return nil
}

// Verify verifies the certificate signature
func (c *Certificate) Verify(issuerPublicKey SM2PublicKey) error {
	// Check validity period
	now := time.Now()
	if now.Before(c.Validity.NotBefore) {
		return fmt.Errorf("certificate not yet valid")
	}
	if now.After(c.Validity.NotAfter) {
		return fmt.Errorf("certificate expired")
	}

	// Verify signature
	certData, err := c.ToDER()
	if err != nil {
		return fmt.Errorf("certificate format error: %v", err)
	}
	digest := sha256.Sum256(certData)

	signature, err := SM2SignatureFromBytes(c.Signature.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature format: %v", err)
	}

	err = SM2Verify(issuerPublicKey, digest[:], signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	return nil
}

// encodeTime encodes a time to UTCTime format (simplified)
func (c *Certificate) encodeTime(t time.Time) []byte {
	result := make([]byte, 15)
	result[0] = 0x17 // UTCTime tag
	result[1] = 13  // Length

	// Simplified encoding
	secs := t.Unix()
	for i := 0; i < 13; i++ {
		result[i+2] = byte((secs >> uint(i*5)) & 0xFF)
	}

	return result
}

// NewCertificateStore creates a new certificate store
func NewCertificateStore() *CertificateStore {
	return &CertificateStore{
		certificates: make([]Certificate, 0),
	}
}

// Add adds a certificate to the store
func (s *CertificateStore) Add(cert Certificate) {
	s.certificates = append(s.certificates, cert)
}

// FindBySubject finds a certificate by subject
func (s *CertificateStore) FindBySubject(subject []byte) *Certificate {
	for i := range s.certificates {
		if string(s.certificates[i].Subject) == string(subject) {
			return &s.certificates[i]
		}
	}
	return nil
}

// FindByIssuer finds a certificate by issuer
func (s *CertificateStore) FindByIssuer(issuer []byte) *Certificate {
	for i := range s.certificates {
		if string(s.certificates[i].Issuer) == string(issuer) {
			return &s.certificates[i]
		}
	}
	return nil
}

// Remove removes a certificate by subject
func (s *CertificateStore) Remove(subject []byte) bool {
	for i, cert := range s.certificates {
		if string(cert.Subject) == string(subject) {
			s.certificates = append(s.certificates[:i], s.certificates[i+1:]...)
			return true
		}
	}
	return false
}

// VerifyChain verifies a certificate chain
// The chain should start with the end-entity certificate and end with a trusted root CA
func (s *CertificateStore) VerifyChain(chain []Certificate) error {
	if len(chain) == 0 {
		return fmt.Errorf("empty certificate chain")
	}

	// Verify each certificate in the chain
	for i := 0; i < len(chain); i++ {
		cert := &chain[i]

		if i == 0 {
			// First certificate: issuer should be a CA in the store
			issuerCert := s.FindByIssuer(cert.Issuer)
			if issuerCert == nil {
				return fmt.Errorf("unknown issuer")
			}

			issuerPubKey, err := SM2PublicKeyFromBytes(issuerCert.PublicKey.PublicKey)
			if err != nil {
				return fmt.Errorf("invalid issuer public key: %v", err)
			}

			if err := cert.Verify(issuerPubKey); err != nil {
				return fmt.Errorf("certificate verification failed: %v", err)
			}
		} else if i < len(chain)-1 {
			// Intermediate certificate: issuer is the previous certificate
			issuerCert := &chain[i-1]
			issuerPubKey, err := SM2PublicKeyFromBytes(issuerCert.PublicKey.PublicKey)
			if err != nil {
				return fmt.Errorf("invalid issuer public key: %v", err)
			}

			if err := cert.Verify(issuerPubKey); err != nil {
				return fmt.Errorf("certificate verification failed: %v", err)
			}
		} else {
			// Root certificate: should be self-signed
			if string(cert.Issuer) != string(cert.Subject) {
				return fmt.Errorf("root certificate not self-signed")
			}
		}
	}

	return nil
}

// All returns all certificates
func (s *CertificateStore) All() []Certificate {
	return s.certificates
}

// Count returns the number of certificates
func (s *CertificateStore) Count() int {
	return len(s.certificates)
}
