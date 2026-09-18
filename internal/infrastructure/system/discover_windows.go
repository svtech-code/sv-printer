//go:build windows

package system

import (
	"context"
	"os/exec"
	"strings"
	"syscall"

	"sv-printer/internal/domain/printer"
)

type Discoverer struct{}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(ctx context.Context) ([]printer.Printer, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", "Get-Printer | Select-Object -ExpandProperty Name")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		// If powershell fails or isn't available, return nil instead of failing the whole discovery
		return nil, nil
	}

	var printers []printer.Printer
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		printers = append(printers, printer.Printer{
			ID:         "system-" + name,
			Name:       name,
			Address:    name,
			Protocol:   printer.ProtocolEscpos,
			Connection: printer.ConnectionSystem,
			Status:     printer.StatusUnknown,
		})
	}
	return printers, nil
}
