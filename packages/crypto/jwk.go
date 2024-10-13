package crypto

import (
	"crypto/rsa"
	"github.com/go-jose/go-jose/v4"
)

// ConvertRSAPublicKeyToJWK converts the RSA public key to a JWK (JSON Web Key)
func ConvertRSAPublicKeyToJWK(pubKey *rsa.PublicKey) (jose.JSONWebKey, error) {
	jwk := jose.JSONWebKey{
		Key:       pubKey,                    // Optionally set a key ID
		Algorithm: string(jose.RSA1_5), // You can specify the algorithm
	}

	return jwk, nil
}
