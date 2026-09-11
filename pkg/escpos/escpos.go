package escpos

import (
	"bytes"
	"image"
	"image/color"
	"strings"
)

const (
	InitPrinter  = "\x1B\x40"             // ESC @
	LineFeed     = "\x0A"                 // LF
	BoldOn       = "\x1B\x45\x01"         // ESC E 1
	BoldOff      = "\x1B\x45\x00"         // ESC E 0
	UnderlineOn  = "\x1B\x2D\x01"         // ESC - 1
	UnderlineOff = "\x1B\x2D\x00"         // ESC - 0
	ItalicOn     = "\x1B\x34"             // ESC 4 (vendor-specific)
	ItalicOff    = "\x1B\x35"             // ESC 5 (vendor-specific)
	ReverseOn    = "\x1D\x42\x01"         // GS B 1
	ReverseOff   = "\x1D\x42\x00"         // GS B 0
	NormalMode   = "\x1B\x21\x00"         // ESC ! 0
	SizeNormal   = "\x1D\x21\x00"         // GS ! 0 (1x1)
	AlignLeft    = "\x1B\x61\x00"         // ESC a 0
	AlignCenter  = "\x1B\x61\x01"         // ESC a 1
	AlignRight   = "\x1B\x61\x02"         // ESC a 2
	Cut          = "\x1D\x56\x41\x10"     // GS V 65 16 (partial cut)
	CashDrawer   = "\x1B\x70\x00\x19\xFA" // ESC p 0 25 250
)

type Receipt struct {
	buf *bytes.Buffer
}

func NewReceipt() *Receipt {
	r := &Receipt{
		buf: new(bytes.Buffer),
	}
	r.buf.WriteString(InitPrinter)
	return r
}

func (r *Receipt) Bytes() []byte {
	return r.buf.Bytes()
}

func (r *Receipt) Init() *Receipt {
	r.buf.WriteString(InitPrinter)
	return r
}

func (r *Receipt) Text(t string) *Receipt {
	r.buf.WriteString(t)
	return r
}

func (r *Receipt) LineFeed() *Receipt {
	r.buf.WriteString(LineFeed)
	return r
}

func (r *Receipt) Bold() *Receipt {
	r.buf.WriteString(BoldOn)
	return r
}

func (r *Receipt) BoldOff() *Receipt {
	r.buf.WriteString(BoldOff)
	return r
}

func (r *Receipt) Underline() *Receipt {
	r.buf.WriteString(UnderlineOn)
	return r
}

func (r *Receipt) UnderlineOff() *Receipt {
	r.buf.WriteString(UnderlineOff)
	return r
}

func (r *Receipt) Italic() *Receipt {
	r.buf.WriteString(ItalicOn)
	return r
}

func (r *Receipt) ItalicOff() *Receipt {
	r.buf.WriteString(ItalicOff)
	return r
}

func (r *Receipt) Reverse() *Receipt {
	r.buf.WriteString(ReverseOn)
	return r
}

func (r *Receipt) ReverseOff() *Receipt {
	r.buf.WriteString(ReverseOff)
	return r
}

func (r *Receipt) Center() *Receipt {
	r.buf.WriteString(AlignCenter)
	return r
}

func (r *Receipt) Left() *Receipt {
	r.buf.WriteString(AlignLeft)
	return r
}

func (r *Receipt) Right() *Receipt {
	r.buf.WriteString(AlignRight)
	return r
}

func (r *Receipt) SetSize(width, height uint8) *Receipt {
	if width < 1 {
		width = 1
	}
	if width > 8 {
		width = 8
	}
	if height < 1 {
		height = 1
	}
	if height > 8 {
		height = 8
	}
	n := (width-1)<<4 | (height - 1)
	r.buf.WriteString("\x1D\x21")
	r.buf.WriteByte(n)
	return r
}

func (r *Receipt) CharacterSpacing(n uint8) *Receipt {
	r.buf.WriteString("\x1B\x20")
	r.buf.WriteByte(n)
	return r
}

func (r *Receipt) LineSpacing(n uint8) *Receipt {
	r.buf.WriteString("\x1B\x33")
	r.buf.WriteByte(n)
	return r
}

func (r *Receipt) ResetStyle() *Receipt {
	r.buf.WriteString(NormalMode)
	r.buf.WriteString(AlignLeft)
	r.buf.WriteString(SizeNormal)
	return r
}

func (r *Receipt) Cut() *Receipt {
	r.buf.WriteString(Cut)
	return r
}

func (r *Receipt) CashDrawer() *Receipt {
	r.buf.WriteString(CashDrawer)
	return r
}

func (r *Receipt) BarcodeCode39(text string) *Receipt {
	return r.barcode(0x04, text, true)
}

func (r *Receipt) BarcodeEAN13(text string) *Receipt {
	return r.barcode(0x02, text, true)
}

func (r *Receipt) BarcodeCode128(text string) *Receipt {
	return r.barcode(0x49, text, false)
}

func (r *Receipt) barcode(mode byte, data string, nulTerminated bool) *Receipt {
	r.buf.WriteString("\x1D\x68\xA2") // GS h 162 (height)
	r.buf.WriteString("\x1D\x77\x03") // GS w 3 (width)
	r.buf.WriteString("\x1D\x6B")     // GS k
	r.buf.WriteByte(mode)
	if nulTerminated {
		r.buf.WriteString(data)
		r.buf.WriteByte(0x00)
	} else {
		r.buf.WriteByte(byte(len(data)))
		r.buf.WriteString(data)
	}
	return r
}

func (r *Receipt) QRCode(data string) *Receipt {
	// Store: GS ( k pL pH 49 80 50 data
	n := len(data) + 3
	r.buf.WriteString("\x1D\x28\x6B")
	r.buf.WriteByte(byte(n & 0xFF))
	r.buf.WriteByte(byte(n >> 8))
	r.buf.WriteByte(0x31) // cn = 49 (QR)
	r.buf.WriteByte(0x50) // fn = 80 (store)
	r.buf.WriteByte(0x32) // m = 50 (model 2)
	r.buf.WriteString(data)
	// Print: GS ( k 03 00 49 81 48
	r.buf.WriteString("\x1D\x28\x6B\x03\x00\x31\x51\x30")
	return r
}

func (r *Receipt) Image(img image.Image) *Receipt {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()
	if w <= 0 || h <= 0 {
		return r
	}

	widthBytes := (w + 7) / 8
	data := make([]byte, widthBytes*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if luminance(img.At(b.Min.X+x, b.Min.Y+y)) < 128 {
				data[y*widthBytes+x/8] |= 1 << (7 - (x % 8))
			}
		}
	}

	r.buf.WriteString("\x1D\x76\x30\x00") // GS v 0 0
	r.buf.WriteByte(byte(widthBytes & 0xFF))
	r.buf.WriteByte(byte(widthBytes >> 8))
	r.buf.WriteByte(byte(h & 0xFF))
	r.buf.WriteByte(byte(h >> 8))
	r.buf.Write(data)
	return r
}

func (r *Receipt) Table(rows [][]string) *Receipt {
	if len(rows) == 0 {
		return r
	}

	cols := 0
	for _, row := range rows {
		if len(row) > cols {
			cols = len(row)
		}
	}

	widths := make([]int, cols)
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	for _, row := range rows {
		var line strings.Builder
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			line.WriteString(cell)
			if i < cols-1 {
				line.WriteString(strings.Repeat(" ", widths[i]-len(cell)+2))
			}
		}
		r.Text(strings.TrimRight(line.String(), " ")).LineFeed()
	}
	return r
}

func luminance(c color.Color) uint8 {
	g := color.GrayModel.Convert(c).(color.Gray)
	return g.Y
}
