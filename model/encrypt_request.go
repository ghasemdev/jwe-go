package model

type EncryptRequest struct {
	Plaintext string `json:"plaintext" validate:"required"`
	PublicKeyPem string `json:"publicKeyPem" validate:"required"`
}
