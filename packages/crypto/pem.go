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

// Import the RSA public key from PEM format
func ImportRSAPublicKeyFromPEM(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(addPEMHeaders(pemString)))
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

func addPEMHeaders(pem string) string {
	if !strings.Contains(pem, "-----BEGIN PUBLIC KEY-----") {
		pem = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s", pem)
	}
	if !strings.Contains(pem, "-----END PUBLIC KEY-----") {
		pem = fmt.Sprintf("%s\n-----END PUBLIC KEY-----", pem)
	}
	return pem
}
