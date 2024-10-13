package model

type DecryptRequest struct {
	Ciphertext string `json:"ciphertext" validate:"required"`
	SecretKey  string `json:"secretKey" validate:"required,len=32"`
}
