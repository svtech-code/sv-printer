package receipt

import (
	"bytes"
	"testing"

	"sv-printer/pkg/escpos"
)

func TestBuild(t *testing.T) {
	doc := Document{
		Cut: true,
		Lines: []Line{
			{Text: "SV TECH", Style: Style{Bold: true, Align: "center"}},
			{Text: "Total $25.000"},
		},
	}

	expected := escpos.NewReceipt().
		Center().
		Bold().
		Text("SV TECH").
		ResetStyle().
		LineFeed().
		Left().
		Text("Total $25.000").
		ResetStyle().
		LineFeed().
		Cut()

	if !bytes.Equal(doc.Build().Bytes(), expected.Bytes()) {
		t.Errorf("Build mismatch.\nGot:  %q\nWant: %q", doc.Build().Bytes(), expected.Bytes())
	}
}

func TestBuildStyleOptions(t *testing.T) {
	doc := Document{
		Lines: []Line{
			{Text: "x", Style: Style{Underline: true, Italic: true, Reverse: true, Align: "right", Size: []int{2, 2}}},
		},
	}

	expected := escpos.NewReceipt().
		Right().
		Underline().
		Italic().
		Reverse().
		SetSize(2, 2).
		Text("x").
		ResetStyle().
		LineFeed()

	if !bytes.Equal(doc.Build().Bytes(), expected.Bytes()) {
		t.Errorf("Build style mismatch.\nGot:  %q\nWant: %q", doc.Build().Bytes(), expected.Bytes())
	}
}

func TestBuildWithWatermark(t *testing.T) {
	doc := Document{
		Cut:   true,
		Lines: []Line{{Text: "hello"}},
	}

	wm := "*** TRIAL ***"

	expected := escpos.NewReceipt().
		Left().Text("hello").ResetStyle().LineFeed().
		Center().Bold().Text(wm).ResetStyle().LineFeed().
		Cut()

	if !bytes.Equal(doc.Build(wm).Bytes(), expected.Bytes()) {
		t.Errorf("Build with watermark mismatch.\nGot:  %q\nWant: %q", doc.Build(wm).Bytes(), expected.Bytes())
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		doc     Document
		wantErr bool
	}{
		{"valid empty lines", Document{}, false},
		{"valid align", Document{Lines: []Line{{Text: "x", Style: Style{Align: "center"}}}}, false},
		{"valid size", Document{Lines: []Line{{Text: "x", Style: Style{Size: []int{2, 2}}}}}, false},
		{"invalid align", Document{Lines: []Line{{Text: "x", Style: Style{Align: "diagonal"}}}}, true},
		{"invalid size len", Document{Lines: []Line{{Text: "x", Style: Style{Size: []int{2}}}}}, true},
		{"invalid size range", Document{Lines: []Line{{Text: "x", Style: Style{Size: []int{9, 1}}}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
