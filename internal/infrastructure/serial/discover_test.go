package serial

import "testing"

func TestFromPath(t *testing.T) {
	p := fromPath("/dev/ttyUSB0")

	if p.ID != "serial-ttyUSB0" {
		t.Errorf("ID = %q, want serial-ttyUSB0", p.ID)
	}
	if p.Connection != "serial" {
		t.Errorf("Connection = %q, want serial", p.Connection)
	}
	if p.Protocol != "escpos" {
		t.Errorf("Protocol = %q, want escpos", p.Protocol)
	}
}
