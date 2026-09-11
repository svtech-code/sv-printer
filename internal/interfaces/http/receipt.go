package http

import (
	"net/http"

	domainErrors "sv-print/internal/domain/errors"
	"sv-print/internal/receipt"
)

type receiptRequest struct {
	PrinterID string `json:"printer_id"`
	receipt.Document
}

func (a *API) PrintReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}

	var req receiptRequest
	if !decodeBody(w, r, &req) {
		return
	}

	if req.PrinterID == "" || len(req.Lines) == 0 {
		s, c, m := mapDomainError(domainErrors.ErrInvalidPayload)
		writeError(w, s, c, m)
		return
	}

	if err := req.Validate(); err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}

	if _, err := a.registry.Find(r.Context(), req.PrinterID); err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}

	wm, ok := a.enforcePrint(w)
	if !ok {
		return
	}

	a.enqueue(w, req.PrinterID, req.Build(wm).Bytes())
}
