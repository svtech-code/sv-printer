package license

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func signForTest(t *testing.T, lic License) ([]byte, ed25519.PublicKey) {
	t.Helper()
	pub, priv, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	data, err := Sign(priv, lic)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	return data, pub
}

func TestSignAndVerify(t *testing.T) {
	lic := License{
		LicenseID: "lic_123",
		Product:   ProductName,
		Customer:  "ACME",
		Tier:      TierFull,
		Features:  []string{FeatureRawPrint, FeatureWebsocket},
	}

	data, pub := signForTest(t, lic)

	got, err := Verify(data, pub)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if got.LicenseID != "lic_123" {
		t.Errorf("LicenseID = %q", got.LicenseID)
	}
	if got.Customer != "ACME" {
		t.Errorf("Customer = %q", got.Customer)
	}
	if got.Tier != TierFull {
		t.Errorf("Tier = %q", got.Tier)
	}
	if !got.HasFeature(FeatureRawPrint) {
		t.Error("expected raw_print feature")
	}
	if got.HasFeature(FeatureSerial) {
		t.Error("did not expect serial feature")
	}
}

func TestVerifyTampered(t *testing.T) {
	lic := License{Product: ProductName, Customer: "ACME", Tier: TierBeta}
	data, pub := signForTest(t, lic)

	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	payload, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}

	// tamper with the payload: change "ACME" to "EVIL"
	tamperedPayload := bytes.Replace(payload, []byte("ACME"), []byte("EVIL"), 1)
	env.Payload = base64.StdEncoding.EncodeToString(tamperedPayload)

	tampered, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal tampered envelope: %v", err)
	}

	if _, err := Verify(tampered, pub); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("Verify() error = %v, want ErrInvalidSignature", err)
	}
}

func TestVerifyWrongKey(t *testing.T) {
	lic := License{Product: ProductName, Customer: "ACME", Tier: TierBeta}
	data, _ := signForTest(t, lic)

	otherPub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	if _, err := Verify(data, otherPub); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("Verify() error = %v, want ErrInvalidSignature with wrong key", err)
	}
}

func TestVerifyWrongProduct(t *testing.T) {
	lic := License{Product: "other-product", Customer: "ACME", Tier: TierFull}
	data, pub := signForTest(t, lic)

	if _, err := Verify(data, pub); !errors.Is(err, ErrWrongProduct) {
		t.Errorf("Verify() error = %v, want ErrWrongProduct", err)
	}
}

func TestExpiry(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	lic := License{Product: ProductName, Customer: "ACME", Tier: TierBeta, Expiry: &exp}

	if lic.Expired(time.Now()) {
		t.Error("license should not be expired")
	}

	past := time.Now().Add(-time.Hour).Format(time.RFC3339)
	lic2 := License{Product: ProductName, Expiry: &past}
	if !lic2.Expired(time.Now()) {
		t.Error("license should be expired")
	}

	lic3 := License{Product: ProductName}
	if lic3.Expired(time.Now()) {
		t.Error("license without expiry should never expire")
	}
}

func TestPublicKeyEmbedded(t *testing.T) {
	pub, err := PublicKey()
	if err != nil {
		t.Fatalf("PublicKey: %v", err)
	}
	if len(pub) != ed25519.PublicKeySize {
		t.Errorf("embedded public key length = %d, want %d", len(pub), ed25519.PublicKeySize)
	}
}
