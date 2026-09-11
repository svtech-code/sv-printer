package usb

import (
	"context"
	"strings"

	"sv-print/internal/domain/printer"
)

type Discoverer struct{}

func NewDiscoverer() *Discoverer { return &Discoverer{} }

func (d *Discoverer) Discover(ctx context.Context) ([]printer.Printer, error) {
	return listUSB(ctx)
}

type usbInfo struct {
	vendorID     string
	productID    string
	manufacturer string
	product      string
	class        string
}

func toPrinter(info usbInfo) printer.Printer {
	id := "usb-" + info.vendorID + "-" + info.productID
	name := info.product
	if name == "" {
		name = info.manufacturer
	}
	if name == "" {
		name = id
	}
	return printer.Printer{
		ID:           id,
		Name:         name,
		Manufacturer: info.manufacturer,
		Model:        info.product,
		Connection:   printer.ConnectionUSB,
		Protocol:     printer.ProtocolEscpos,
		Status:       printer.StatusUnknown,
	}
}

func isPrinterClass(class string) bool {
	c := strings.TrimSpace(class)
	return c == "07" || c == "ff"
}
