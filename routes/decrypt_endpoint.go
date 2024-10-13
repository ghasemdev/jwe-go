package routes

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"jwe-go/model"
	"jwe-go/packages/json"
	"jwe-go/packages/pool"
	"jwe-go/packages/schema"
	"net/http"
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

	decrypted, err := decryptedObject.Decrypt(decryption.SecretKey)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.String(http.StatusOK, string(decrypted))
}
