//go:build windows

package system

import (
	"context"

	"sv-printer/internal/domain/printer"
)

type Discoverer struct{}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(ctx context.Context) ([]printer.Printer, error) {
	// TODO: Implement Phase 2 Windows Spooler discovery (via wmic or PowerShell)
	return nil, nil
}
