package licensing

import (
	"sync"
	"time"

	"sv-printer/internal/deviceid"
	"sv-printer/internal/license"
)

const (
	TrialWatermark  = "*** SV PRINTER — LICENCIA DE PRUEBA ***\nwww.svtech.cl"
	TrialDailyQuota = 50
)

type State struct {
	lic      *license.License
	deviceID string
	mu       sync.Mutex
	day      string
	count    int
}

func New(lic *license.License) *State {
	return &State{lic: lic}
}

func NewBound(lic *license.License, deviceID string) *State {
	return &State{lic: lic, deviceID: deviceID}
}

func FromFile(path string) *State {
	lic, err := license.Load(path)
	if err != nil {
		return New(nil)
	}
	id, _ := deviceid.ID()
	return NewBound(&lic, id)
}

func (s *State) IsActive() bool {
	if s.lic == nil || s.lic.Expired(time.Now()) {
		return false
	}
	if s.lic.Fingerprint != nil && *s.lic.Fingerprint != s.deviceID {
		return false
	}
	return true
}

func (s *State) IsTrial() bool {
	return !s.IsActive()
}

func (s *State) Tier() license.Tier {
	if s.IsActive() {
		return s.lic.Tier
	}
	return license.TierTrial
}

func (s *State) HasFeature(f string) bool {
	return s.IsActive() && s.lic.HasFeature(f)
}

func (s *State) Watermark() string {
	if s.IsTrial() {
		return TrialWatermark
	}
	return ""
}

func (s *State) AllowPrint() bool {
	if s.IsActive() {
		return true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if s.day != today {
		s.day = today
		s.count = 0
	}
	if s.count >= TrialDailyQuota {
		return false
	}
	s.count++
	return true
}
