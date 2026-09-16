package discovery

import (
	"context"
	"testing"

	"sv-printer/internal/domain/printer"
)

func TestRegistryFind(t *testing.T) {
	r := NewRegistry()
	r.Register(NewStaticDiscoverer([]printer.Printer{
		{ID: "net-192.168.1.100:9100", Name: "Caja 1"},
	}))

	p, err := r.Find(context.Background(), "net-192.168.1.100:9100")
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if p.ID != "net-192.168.1.100:9100" {
		t.Errorf("Find() ID = %q", p.ID)
	}

	if _, err := r.Find(context.Background(), "missing"); err == nil {
		t.Fatal("Find() expected error for missing printer")
	}
}
