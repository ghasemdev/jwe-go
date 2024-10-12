package model

type EncryptRequest struct {
	Plaintext string `json:"plaintext" validate:"required"`
	PublicKey string `json:"publicKey" validate:"required"`
}
