package license

import "time"

const ProductName = "sv-print"

type Tier string

const (
	TierTrial Tier = "trial"
	TierBeta  Tier = "beta"
	TierFull  Tier = "full"
)

const (
	FeatureRawPrint  = "raw_print"
	FeatureWebsocket = "websocket"
	FeatureSerial    = "serial"
)

type License struct {
	LicenseID   string   `json:"license_id"`
	Product     string   `json:"product"`
	Customer    string   `json:"customer"`
	Tier        Tier     `json:"tier"`
	Features    []string `json:"features,omitempty"`
	Expiry      *string  `json:"expiry,omitempty"`
	Fingerprint *string  `json:"fingerprint,omitempty"`
}

func (l License) HasFeature(f string) bool {
	for _, feat := range l.Features {
		if feat == f {
			return true
		}
	}
	return false
}

func (l License) ExpiresAt() (time.Time, bool) {
	if l.Expiry == nil {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, *l.Expiry)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func (l License) Expired(now time.Time) bool {
	t, ok := l.ExpiresAt()
	if !ok {
		return false
	}
	return now.After(t)
}
