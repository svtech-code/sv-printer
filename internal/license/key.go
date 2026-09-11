package license

import (
	"crypto/ed25519"
	_ "embed"
	"encoding/hex"
	"os"
	"strings"
)

//go:embed public.key
var publicKeyHex string

func PublicKey() (ed25519.PublicKey, error) {
	h := strings.TrimSpace(publicKeyHex)
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(b), nil
}

func Load(path string) (License, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return License{}, err
	}
	pub, err := PublicKey()
	if err != nil {
		return License{}, err
	}
	return Verify(data, pub)
}
