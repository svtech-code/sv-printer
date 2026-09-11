package discovery

import (
	"context"
	"testing"

	"sv-print/internal/domain/printer"
)

type mockDiscovery struct {
	printers []printer.Printer
}

func (m *mockDiscovery) Discover(ctx context.Context) ([]printer.Printer, error) {
	return m.printers, nil
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	m1 := &mockDiscovery{
		printers: []printer.Printer{
			{ID: "p1", Name: "Printer 1"},
		},
	}
	m2 := &mockDiscovery{
		printers: []printer.Printer{
			{ID: "p2", Name: "Printer 2"},
		},
	}

	r.Register(m1)
	r.Register(m2)

	ctx := context.Background()
	printers, err := r.DiscoverAll(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(printers) != 2 {
		t.Errorf("Expected 2 printers, got %d", len(printers))
	}

	if printers[0].ID != "p1" || printers[1].ID != "p2" {
		t.Errorf("Printer IDs mismatch")
	}
}
