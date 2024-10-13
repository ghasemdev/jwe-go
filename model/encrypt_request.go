package model

type EncryptRequest struct {
	Plaintext    string `json:"plaintext" validate:"required"`
	PublicKeyPem string `json:"publicKeyPem"`
	CertificatePem string `json:"certificatePem"`
}
