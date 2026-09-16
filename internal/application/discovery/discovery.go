package discovery

import (
	"context"

	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/domain/printer"
)

type PrinterDiscovery interface {
	Discover(ctx context.Context) ([]printer.Printer, error)
}

type Registry struct {
	discoverers []PrinterDiscovery
}

func NewRegistry() *Registry {
	return &Registry{
		discoverers: make([]PrinterDiscovery, 0),
	}
}

func (r *Registry) Register(d PrinterDiscovery) {
	r.discoverers = append(r.discoverers, d)
}

func (r *Registry) DiscoverAll(ctx context.Context) ([]printer.Printer, error) {
	var allPrinters []printer.Printer

	for _, d := range r.discoverers {
		printers, err := d.Discover(ctx)
		if err != nil {
			// In a real scenario we might log this and continue instead of failing entirely
			continue
		}
		allPrinters = append(allPrinters, printers...)
	}

	return allPrinters, nil
}

func (r *Registry) Find(ctx context.Context, id string) (printer.Printer, error) {
	printers, err := r.DiscoverAll(ctx)
	if err != nil {
		return printer.Printer{}, err
	}

	for _, p := range printers {
		if p.ID == id {
			return p, nil
		}
	}

	return printer.Printer{}, domainErrors.ErrPrinterNotFound
}
