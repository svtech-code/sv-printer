package discovery

import (
	"context"

	"sv-print/internal/domain/printer"
)

type StaticDiscoverer struct {
	printers []printer.Printer
}

func NewStaticDiscoverer(printers []printer.Printer) *StaticDiscoverer {
	return &StaticDiscoverer{printers: printers}
}

func (s *StaticDiscoverer) Discover(ctx context.Context) ([]printer.Printer, error) {
	return s.printers, nil
}
