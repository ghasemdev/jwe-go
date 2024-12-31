package routes

import (
	"bytes"
	"encoding/base64"
	"jwe-go/model"
	"jwe-go/packages/json"
	"jwe-go/packages/pool"
	"jwe-go/packages/schema"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ghasemdev/go-jose/v4"
)

func DecryptEndpoint(context *gin.Context) {
	var decryption model.DecryptRequest

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
	if err := json.StrictUnmarshal(buf.Bytes(), &decryption); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or unknown field"})
		return
	}

	// Manually validate the struct using the validator
	if err := schema.Validate.Struct(decryption); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Decrypt the message using the private key
	decryptedObject, err := jose.ParseEncryptedCompact(
		decryption.Ciphertext,
		[]jose.KeyAlgorithm{jose.DIRECT},
		[]jose.ContentEncryption{jose.A256GCM},
	)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var decrypted []byte
	var secKeyError error

	switch {
	case len(decryption.SecretKey) > 0 && len(decryption.SecretKeyBase64) > 0:
		context.JSON(http.StatusBadRequest, gin.H{"error": "both SecretKey and SecretKeyBase64 cannot be set at the same time"})
		return
	case len(decryption.SecretKey) == 32:
		decrypted, secKeyError = decryptedObject.Decrypt(decryption.SecretKey)
	case len(decryption.SecretKeyBase64) > 0:
		secretKey, err := base64.RawURLEncoding.DecodeString(decryption.SecretKeyBase64)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		decrypted, secKeyError = decryptedObject.Decrypt(secretKey)
	default:
		context.JSON(http.StatusBadRequest, gin.H{"error": "No valid secret key provided"})
		return
	}

	if secKeyError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": secKeyError.Error()})
		return
	}

	context.String(http.StatusOK, string(decrypted))
}
