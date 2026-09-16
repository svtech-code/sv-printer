package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/infrastructure/transport"
)

type SystemTransport struct {
	printerName string
}

func NewSystemTransport(printerName string) transport.PrinterTransport {
	return &SystemTransport{
		printerName: printerName,
	}
}

func (t *SystemTransport) Open(ctx context.Context) error {
	// For system (CUPS) printing via lp, connection is stateless.
	// We check if lp command exists.
	_, err := exec.LookPath("lp")
	if err != nil {
		return fmt.Errorf("%w: lp command not found (CUPS not available)", domainErrors.ErrPrinterConnectionFailed)
	}
	return nil
}

func (t *SystemTransport) Write(ctx context.Context, data []byte) error {
	// Create a temporary file for the payload
	tmpFile, err := os.CreateTemp("", "sv-printer-cups-*.bin")
	if err != nil {
		return fmt.Errorf("%w: failed to create temp file: %v", domainErrors.ErrPrintFailed, err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("%w: failed to write to temp file: %v", domainErrors.ErrPrintFailed, err)
	}
	tmpFile.Close()

	// Execute lp command
	// lp -d <printerName> -o raw <tempFile>
	cmd := exec.CommandContext(ctx, "lp", "-d", t.printerName, "-o", "raw", tmpFile.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: lp command failed: %v, output: %s", domainErrors.ErrPrintFailed, err, string(out))
	}

	return nil
}

func (t *SystemTransport) Close() error {
	return nil
}
