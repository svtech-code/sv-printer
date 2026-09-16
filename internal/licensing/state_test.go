package licensing

import (
	"testing"
	"time"

	"sv-printer/internal/license"
)

func TestTrialState(t *testing.T) {
	s := New(nil)

	if !s.IsTrial() {
		t.Error("nil license should be trial")
	}
	if s.Tier() != license.TierTrial {
		t.Errorf("Tier = %q, want trial", s.Tier())
	}
	if s.Watermark() == "" {
		t.Error("trial should have a watermark")
	}
	if s.HasFeature(license.FeatureRawPrint) {
		t.Error("trial should not have features")
	}
}

func TestActiveState(t *testing.T) {
	lic := license.License{
		Product:  license.ProductName,
		Customer: "ACME",
		Tier:     license.TierFull,
		Features: []string{license.FeatureRawPrint},
	}
	s := New(&lic)

	if s.IsTrial() {
		t.Error("active license should not be trial")
	}
	if s.Tier() != license.TierFull {
		t.Errorf("Tier = %q, want full", s.Tier())
	}
	if s.Watermark() != "" {
		t.Error("active license should not watermark")
	}
	if !s.HasFeature(license.FeatureRawPrint) {
		t.Error("expected raw_print feature")
	}
	if s.HasFeature(license.FeatureWebsocket) {
		t.Error("did not expect websocket feature")
	}
}

func TestExpiredLicenseIsTrial(t *testing.T) {
	past := time.Now().Add(-time.Hour).Format(time.RFC3339)
	lic := license.License{Product: license.ProductName, Tier: license.TierBeta, Expiry: &past}
	s := New(&lic)

	if !s.IsTrial() {
		t.Error("expired license should be treated as trial")
	}
}

func TestQuota(t *testing.T) {
	s := New(nil)

	for i := 0; i < TrialDailyQuota; i++ {
		if !s.AllowPrint() {
			t.Fatalf("print %d should be allowed", i)
		}
	}
	if s.AllowPrint() {
		t.Error("print beyond quota should be denied")
	}
}

func TestQuotaUnlimitedForActive(t *testing.T) {
	lic := license.License{Product: license.ProductName, Tier: license.TierFull}
	s := New(&lic)

	for i := 0; i < 1000; i++ {
		if !s.AllowPrint() {
			t.Fatalf("active license print %d should be allowed", i)
		}
	}
}

func TestFromFileMissingIsTrial(t *testing.T) {
	s := FromFile(t.TempDir() + "/nope.key")
	if !s.IsTrial() {
		t.Error("missing license file should be trial")
	}
}

func TestFingerprintMatch(t *testing.T) {
	fp := "device-123"
	lic := license.License{Product: license.ProductName, Tier: license.TierFull, Fingerprint: &fp}
	s := NewBound(&lic, fp)

	if !s.IsActive() {
		t.Error("matching fingerprint should be active")
	}
}

func TestFingerprintMismatch(t *testing.T) {
	fp := "device-123"
	lic := license.License{Product: license.ProductName, Tier: license.TierFull, Fingerprint: &fp}
	s := NewBound(&lic, "other-device")

	if s.IsActive() {
		t.Error("mismatched fingerprint should be trial")
	}
	if !s.IsTrial() {
		t.Error("mismatched fingerprint should be treated as trial")
	}
}

func TestNoFingerprintIgnoresDevice(t *testing.T) {
	lic := license.License{Product: license.ProductName, Tier: license.TierFull}
	s := NewBound(&lic, "any-device")

	if !s.IsActive() {
		t.Error("license without fingerprint should ignore device id")
	}
}
