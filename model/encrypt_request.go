package model

type EncryptRequest struct {
	Plaintext       string `json:"plaintext" validate:"required"`
	SecretKeyBase64 string `json:"secretKeyBase64"  validate:"required"`
	PublicKeyPem    string `json:"publicKeyPem"`
	CertificatePem  string `json:"certificatePem"`
}
