package escpos

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestReceiptBuilder(t *testing.T) {
	r := NewReceipt()
	r.Center().Bold().Text("SV TECH").LineFeed().BoldOff().Left().Text("Total").Cut()

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString(AlignCenter)
	expected.WriteString(BoldOn)
	expected.WriteString("SV TECH")
	expected.WriteString(LineFeed)
	expected.WriteString(BoldOff)
	expected.WriteString(AlignLeft)
	expected.WriteString("Total")
	expected.WriteString(Cut)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestReceiptResetStyle(t *testing.T) {
	r := NewReceipt()
	r.ResetStyle()

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString(NormalMode)
	expected.WriteString(AlignLeft)
	expected.WriteString(SizeNormal)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestResetStyleDoesNotReinit(t *testing.T) {
	r := NewReceipt()
	r.Bold().Center().ResetStyle().Text("x")

	got := r.Bytes()
	if bytes.Contains(got, []byte(InitPrinter+InitPrinter)) {
		t.Error("ResetStyle should not emit a second ESC @ (InitPrinter)")
	}
}

func TestUnderline(t *testing.T) {
	r := NewReceipt()
	r.Underline().Text("u").UnderlineOff()

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString(UnderlineOn)
	expected.WriteString("u")
	expected.WriteString(UnderlineOff)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestItalic(t *testing.T) {
	r := NewReceipt()
	r.Italic().Text("i").ItalicOff()

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString(ItalicOn)
	expected.WriteString("i")
	expected.WriteString(ItalicOff)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestReverse(t *testing.T) {
	r := NewReceipt()
	r.Reverse().Text("r").ReverseOff()

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString(ReverseOn)
	expected.WriteString("r")
	expected.WriteString(ReverseOff)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestSetSize(t *testing.T) {
	r := NewReceipt()
	r.SetSize(2, 2)

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x21")
	expected.WriteByte(0x11)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestSetSizeClamped(t *testing.T) {
	r := NewReceipt()
	r.SetSize(99, 0)

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x21")
	expected.WriteByte((8-1)<<4 | 0) // width 8, height 1

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestCharacterSpacing(t *testing.T) {
	r := NewReceipt()
	r.CharacterSpacing(2)

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1B\x20")
	expected.WriteByte(2)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestLineSpacing(t *testing.T) {
	r := NewReceipt()
	r.LineSpacing(30)

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1B\x33")
	expected.WriteByte(30)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestCashDrawer(t *testing.T) {
	r := NewReceipt()
	r.CashDrawer()

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString(CashDrawer)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestBarcodeCode39(t *testing.T) {
	r := NewReceipt()
	r.BarcodeCode39("ABC")

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x68\xA2")
	expected.WriteString("\x1D\x77\x03")
	expected.WriteString("\x1D\x6B")
	expected.WriteByte(0x04)
	expected.WriteString("ABC")
	expected.WriteByte(0x00)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestBarcodeEAN13(t *testing.T) {
	r := NewReceipt()
	r.BarcodeEAN13("123456789012")

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x68\xA2")
	expected.WriteString("\x1D\x77\x03")
	expected.WriteString("\x1D\x6B")
	expected.WriteByte(0x02)
	expected.WriteString("123456789012")
	expected.WriteByte(0x00)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestBarcodeCode128(t *testing.T) {
	r := NewReceipt()
	r.BarcodeCode128("ABC123")

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x68\xA2")
	expected.WriteString("\x1D\x77\x03")
	expected.WriteString("\x1D\x6B")
	expected.WriteByte(0x49)
	expected.WriteByte(6) // length
	expected.WriteString("ABC123")

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestQRCode(t *testing.T) {
	r := NewReceipt()
	r.QRCode("https://example.com")

	data := "https://example.com"
	n := len(data) + 3

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x28\x6B")
	expected.WriteByte(byte(n & 0xFF))
	expected.WriteByte(byte(n >> 8))
	expected.WriteByte(0x31)
	expected.WriteByte(0x50)
	expected.WriteByte(0x32)
	expected.WriteString(data)
	expected.WriteString("\x1D\x28\x6B\x03\x00\x31\x51\x30")

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.Black)
		}
	}

	r := NewReceipt()
	r.Image(img)

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("\x1D\x76\x30\x00")
	expected.WriteByte(1) // xL (widthBytes = 1)
	expected.WriteByte(0) // xH
	expected.WriteByte(8) // yL (height = 8)
	expected.WriteByte(0) // yH
	for i := 0; i < 8; i++ {
		expected.WriteByte(0xFF)
	}

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func TestTable(t *testing.T) {
	r := NewReceipt()
	r.Table([][]string{
		{"Item", "Price"},
		{"Cafe", "$2.00"},
	})

	expected := new(bytes.Buffer)
	expected.WriteString(InitPrinter)
	expected.WriteString("Item  Price")
	expected.WriteString(LineFeed)
	expected.WriteString("Cafe  $2.00")
	expected.WriteString(LineFeed)

	assertBytes(t, r.Bytes(), expected.Bytes())
}

func assertBytes(t *testing.T, got, want []byte) {
	t.Helper()
	if !bytes.Equal(got, want) {
		t.Errorf("bytes mismatch.\nGot:  %q\nWant: %q", got, want)
	}
}
