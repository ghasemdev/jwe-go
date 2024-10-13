package routes

import (
	"bytes"
	"crypto/rsa"
	"jwe-go/model"
	"jwe-go/packages/crypto"
	"jwe-go/packages/json"
	"jwe-go/packages/pool"
	"jwe-go/packages/schema"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
)

func EncryptEndpoint(context *gin.Context) {
	var encryption model.EncryptRequest

	// Get a buffer from the pool
	buf := pool.BufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer pool.BufPool.Put(buf)

	// Read request body
	if _, err := buf.ReadFrom(context.Request.Body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Use strict unmarshaling
	if err := json.StrictUnmarshal(buf.Bytes(), &encryption); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or unknown field"})
		return
	}

	// Manually validate the struct using the validator
	if err := schema.Validate.Struct(encryption); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var publicKey *rsa.PublicKey
	var err error

	// Try to import the RSA public key from either the Public Key PEM or Certificate PEM.
	switch {
	case len(encryption.PublicKeyPem) > 0:
		publicKey, err = crypto.ImportRSAPublicKeyFromPEM(encryption.PublicKeyPem)
	case len(encryption.CertificatePem) > 0:
		publicKey, err = crypto.ImportRSAPublicKeyFromCertificatePEM(encryption.CertificatePem)
	default:
		context.JSON(http.StatusBadRequest, gin.H{"error": "No valid PEM provided"})
		return
	}

	// Handle any error from the import functions.
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the RSA public key to a JWK
	jwk, err := crypto.ConvertRSAPublicKeyToJWK(publicKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compute thumbprint of the public key (SHA256)
	thumbprint, err := crypto.GetJWKThumbprint(jwk)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute public key thumbprint"})
		return
	}

	// Create JWE Encrypter with RSA-OAEP-256 and AES-GCM
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM, // Content encryption algorithm
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256, // Key encryption algorithm
			Key:       publicKey,         // Recipient's public key
		},
		(&jose.EncrypterOptions{}).WithHeader("server_kid", thumbprint), // Add custom header (server_kid)
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Encrypt the data
	jwe, err := encrypter.Encrypt([]byte(encryption.Plaintext))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serialize JWE to a compact format
	serialized, err := jwe.CompactSerialize()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.String(http.StatusOK, serialized)
}
