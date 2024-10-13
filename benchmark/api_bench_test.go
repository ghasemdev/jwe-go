package benchmark

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"jwe-go/model"
	"jwe-go/packages/schema"
	"jwe-go/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	schema.Validate = validator.New(validator.WithRequiredStructEnabled())
	router := gin.Default()
	router.POST("/encrypt", routes.EncryptEndpoint)
	router.POST("/decrypt", routes.DecryptEndpoint)
	return router
}

func BenchmarkEncryptEndpoint(b *testing.B) {
	router := setupRouter()
	reqBody := model.EncryptRequest{
		Plaintext:    "Hello, World!",
		PublicKeyPem: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt5TICyK59sJggCB8YbGp0uTYMTr3V4fJIvaZgujqEAEtGB6QDCS6IOnqtZzetVspDVH1tIV/wlOuFzgga3kKhMawb2Q/zKLTtK+QnNngaPcND8PClnY/ro1BBy9sxjO3FgCHKrRkRAnzif3qGLQHgvGNk1MWJ/qvdg8F2rCqAbcmCxdROUcLNEjbeW1pReSFEVOJRvrQDmDGvJRZArSx8CCCPRJPqzjByV3pSqqHCqIQ2P9aeXW8L1lvzOuwFCGpFupjoc5v3G8M8hthxGfueVjGz6iw0ka6+V/Zem6XkEJFXHWTnvmemYziMDswFE0GxpeizuVaY/jLZ30gCG/0CQIDAQAB",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}

func BenchmarkDecryptEndpoint(b *testing.B) {
	router := setupRouter()
	reqBody := model.DecryptRequest{
		Ciphertext: "eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..Drq7NpAeq6PVWlyT.Z95Y6gDDWHpqpgWzptrghpNiIst2S3qFSdgmqJT8yKM4M2B58-r3sFKAHPx4wtVWnihWgW6ez9ttf8V0CIPJEstJqXOudGxzlzpKosBFwHpnbouaLNOnWQjPQGthufc.p4JT_Z5GxdABNhld0YpYLQ",
		SecretKey:  "Understandably-Daring-Return-857",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/decrypt", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}
