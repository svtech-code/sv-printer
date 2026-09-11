package serial

import (
	"context"
	"path/filepath"

	"sv-print/internal/domain/printer"
)

type Discoverer struct{}

func NewDiscoverer() *Discoverer { return &Discoverer{} }

func (d *Discoverer) Discover(ctx context.Context) ([]printer.Printer, error) {
	var printers []printer.Printer
	for _, pattern := range patterns() {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, m := range matches {
			printers = append(printers, fromPath(m))
		}
	}
	return printers, nil
}

func fromPath(path string) printer.Printer {
	base := filepath.Base(path)
	return printer.Printer{
		ID:         "serial-" + base,
		Name:       base,
		Connection: printer.ConnectionSerial,
		Protocol:   printer.ProtocolEscpos,
		Status:     printer.StatusUnknown,
	}
}
