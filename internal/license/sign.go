package license

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

func Sign(priv ed25519.PrivateKey, lic License) ([]byte, error) {
	payload, err := json.Marshal(lic)
	if err != nil {
		return nil, err
	}

	sig := ed25519.Sign(priv, payload)

	env := envelope{
		Payload:   base64.StdEncoding.EncodeToString(payload),
		Signature: base64.StdEncoding.EncodeToString(sig),
	}

	return json.MarshalIndent(env, "", "  ")
}

func GenerateKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return pub, priv, nil
}
