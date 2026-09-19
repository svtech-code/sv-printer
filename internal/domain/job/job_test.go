package job

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID()
	if !strings.HasPrefix(id, "job_") {
		t.Errorf("NewID() = %q, want prefix job_", id)
	}
	if len(id) <= len("job_") {
		t.Errorf("NewID() = %q, want a non-empty suffix", id)
	}
}

func TestNewIDUnique(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 10000; i++ {
		id := NewID()
		if seen[id] {
			t.Fatalf("duplicate ID %q", id)
		}
		seen[id] = true
	}
}
