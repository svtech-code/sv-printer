package usb

import "testing"

func TestToPrinter(t *testing.T) {
	p := toPrinter(usbInfo{
		vendorID:     "1234",
		productID:    "5678",
		manufacturer: "XPrinter",
		product:      "XP-Q200",
	})

	if p.ID != "usb-1234-5678" {
		t.Errorf("ID = %q, want usb-1234-5678", p.ID)
	}
	if p.Name != "XP-Q200" {
		t.Errorf("Name = %q, want XP-Q200", p.Name)
	}
	if p.Manufacturer != "XPrinter" {
		t.Errorf("Manufacturer = %q, want XPrinter", p.Manufacturer)
	}
	if p.Connection != "usb" {
		t.Errorf("Connection = %q, want usb", p.Connection)
	}
}

func TestToPrinterFallsBackToManufacturer(t *testing.T) {
	p := toPrinter(usbInfo{vendorID: "1234", productID: "5678", manufacturer: "XPrinter"})

	if p.Name != "XPrinter" {
		t.Errorf("Name = %q, want XPrinter (fallback)", p.Name)
	}
}

func TestIsPrinterClass(t *testing.T) {
	if !isPrinterClass("07") {
		t.Error("class 07 should be a printer class")
	}
	if !isPrinterClass("ff") {
		t.Error("class ff should be a printer class (vendor-specific)")
	}
	if isPrinterClass("09") {
		t.Error("class 09 (hub) should not be a printer class")
	}
}
