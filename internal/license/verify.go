package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidFormat    = errors.New("invalid license format")
	ErrInvalidSignature = errors.New("invalid license signature")
	ErrExpired          = errors.New("license expired")
	ErrWrongProduct     = errors.New("license product mismatch")
)

type envelope struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

func Verify(data []byte, pub ed25519.PublicKey) (License, error) {
	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return License{}, ErrInvalidFormat
	}

	payload, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return License{}, ErrInvalidFormat
	}

	sig, err := base64.StdEncoding.DecodeString(env.Signature)
	if err != nil {
		return License{}, ErrInvalidFormat
	}

	if !ed25519.Verify(pub, payload, sig) {
		return License{}, ErrInvalidSignature
	}

	var lic License
	if err := json.Unmarshal(payload, &lic); err != nil {
		return License{}, ErrInvalidFormat
	}

	if lic.Product != ProductName {
		return License{}, ErrWrongProduct
	}

	return lic, nil
}
