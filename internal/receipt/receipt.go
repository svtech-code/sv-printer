package receipt

import (
	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/pkg/escpos"
)

type Document struct {
	Cut   bool   `json:"cut"`
	Lines []Line `json:"lines"`
}

type Line struct {
	Text  string `json:"text"`
	Style Style  `json:"style,omitempty"`
}

type Style struct {
	Bold      bool   `json:"bold,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Italic    bool   `json:"italic,omitempty"`
	Reverse   bool   `json:"reverse,omitempty"`
	Align     string `json:"align,omitempty"`
	Size      []int  `json:"size,omitempty"`
}

func (d Document) Validate() error {
	for _, l := range d.Lines {
		switch l.Style.Align {
		case "", "left", "center", "right":
		default:
			return domainErrors.ErrInvalidPayload
		}
		if len(l.Style.Size) != 0 && len(l.Style.Size) != 2 {
			return domainErrors.ErrInvalidPayload
		}
		for _, s := range l.Style.Size {
			if s < 1 || s > 8 {
				return domainErrors.ErrInvalidPayload
			}
		}
	}
	return nil
}

func (d Document) Build(watermark ...string) *escpos.Receipt {
	r := escpos.NewReceipt()

	wm := ""
	if len(watermark) > 0 {
		wm = watermark[0]
	}

	for _, l := range d.Lines {
		st := l.Style
		switch st.Align {
		case "center":
			r.Center()
		case "right":
			r.Right()
		default:
			r.Left()
		}
		if st.Bold {
			r.Bold()
		}
		if st.Underline {
			r.Underline()
		}
		if st.Italic {
			r.Italic()
		}
		if st.Reverse {
			r.Reverse()
		}
		if len(st.Size) == 2 {
			r.SetSize(uint8(st.Size[0]), uint8(st.Size[1]))
		}
		r.Text(l.Text)
		r.ResetStyle()
		r.LineFeed()
	}

	if wm != "" {
		r.Center().Bold().Text(wm).ResetStyle().LineFeed()
	}

	if d.Cut {
		r.Cut()
	}
	return r
}
