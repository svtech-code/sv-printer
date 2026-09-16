//go:build unix || darwin

package system

import (
	"context"
	"os/exec"
	"strings"

	"sv-printer/internal/domain/printer"
)

type Discoverer struct{}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(ctx context.Context) ([]printer.Printer, error) {
	// Execute lpstat -p to get the list of printers
	cmd := exec.CommandContext(ctx, "lpstat", "-p")
	out, err := cmd.Output()
	if err != nil {
		// CUPS might not be running or installed
		return nil, nil
	}

	var printers []printer.Printer
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		// Output format: "printer PRINTER_NAME is idle.  enabled since..."
		if strings.HasPrefix(line, "printer ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[1]
				printers = append(printers, printer.Printer{
					ID:           "system-" + name,
					Name:         name,
					Manufacturer: "System",
					Connection:   printer.ConnectionSystem,
					Address:      name,
					Status:       printer.StatusUnknown,
					Protocol:     printer.ProtocolEscpos,
				})
			}
		}
	}

	return printers, nil
}
