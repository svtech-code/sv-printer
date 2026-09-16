package discovery

import (
	"context"
	"testing"

	"sv-printer/internal/domain/printer"
)

func TestStaticDiscoverer(t *testing.T) {
	expected := []printer.Printer{
		{ID: "net-192.168.1.100:9100", Name: "Caja 1", Connection: printer.ConnectionNetwork},
	}

	d := NewStaticDiscoverer(expected)

	got, err := d.Discover(context.Background())
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("Discover() returned %d printers, want 1", len(got))
	}
	if got[0].ID != expected[0].ID {
		t.Errorf("Discover() ID = %q, want %q", got[0].ID, expected[0].ID)
	}
}
