package crypto

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"github.com/go-jose/go-jose/v4"
)

// GetJWKThumbprint calculates the thumbprint of the JWK using SHA-256
func GetJWKThumbprint(jwk jose.JSONWebKey) (string, error) {
	thumbprint, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed to calculate JWK thumbprint: %v", err)
	}

	// Encode the thumbprint as Base64 (URL encoding, no padding)
	thumbprintBase64 := base64.RawURLEncoding.EncodeToString(thumbprint)

	// Replace all occurrences of "=" with an empty string ""
	return thumbprintBase64, nil
}
