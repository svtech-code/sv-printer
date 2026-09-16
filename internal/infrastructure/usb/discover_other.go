//go:build !linux

package usb

import (
	"context"

	"sv-printer/internal/domain/printer"
)

func listUSB(ctx context.Context) ([]printer.Printer, error) {
	return nil, nil
}
