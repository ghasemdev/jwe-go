package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

// Export the RSA public key to PEM format
func ExportRSAPublicKeyAsPEM(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("error marshalling public key to ASN.1: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	return string(pubPEM), nil
}

// ImportRSAPublicKeyFromCertificatePEM extracts the RSA public key from a PEM-encoded certificate.
func ImportRSAPublicKeyFromCertificatePEM(certificatePEM string) (*rsa.PublicKey, error) {
	// Decode the PEM block
	block, _ := pem.Decode([]byte(addPEMHeaders(certificatePEM, "CERTIFICATE")))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the certificate from the decoded block
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Extract the public key from the certificate and assert it as RSA
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate does not contain an RSA public key")
	}

	return pubKey, nil
}

// Import the RSA public key from PEM format
func ImportRSAPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(addPEMHeaders(publicKeyPEM, "PUBLIC KEY")))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Verify that the key is an RSA public key
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPub, nil
}

func addPEMHeaders(pem, pemType string) string {
	beginHeader := fmt.Sprintf("-----BEGIN %s-----", pemType)
	endHeader := fmt.Sprintf("-----END %s-----", pemType)

	if !strings.Contains(pem, beginHeader) {
		pem = fmt.Sprintf("%s\n%s", beginHeader, pem)
	}
	if !strings.Contains(pem, endHeader) {
		pem = fmt.Sprintf("%s\n%s", pem, endHeader)
	}

	return pem
}
