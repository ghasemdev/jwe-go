package model

type DecryptRequest struct {
	Ciphertext      string `json:"ciphertext" validate:"required"`
	SecretKey       string `json:"secretKey"`
	SecretKeyBase64 string `json:"secretKeyBase64"`
}
