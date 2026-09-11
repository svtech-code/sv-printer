//go:build linux

package usb

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"sv-print/internal/domain/printer"
)

func listUSB(ctx context.Context) ([]printer.Printer, error) {
	root := "/sys/bus/usb/devices"
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var printers []printer.Printer
	for _, e := range entries {
		base := filepath.Join(root, e.Name())
		info := usbInfo{
			vendorID:     readAttr(base, "idVendor"),
			productID:    readAttr(base, "idProduct"),
			manufacturer: readAttr(base, "manufacturer"),
			product:      readAttr(base, "product"),
			class:        readAttr(base, "bDeviceClass"),
		}
		if info.vendorID == "" || info.productID == "" {
			continue
		}
		if !isPrinterClass(info.class) {
			continue
		}
		printers = append(printers, toPrinter(info))
	}
	return printers, nil
}

func readAttr(base, name string) string {
	data, err := os.ReadFile(filepath.Join(base, name))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
