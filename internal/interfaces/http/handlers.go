package http

import (
	"net/http"
	"runtime"
	"time"

	"sv-print/internal/application/discovery"
	"sv-print/internal/application/events"
	"sv-print/internal/application/printing"
	domainErrors "sv-print/internal/domain/errors"
	"sv-print/internal/domain/job"
	"sv-print/internal/domain/printer"
	"sv-print/internal/license"
	"sv-print/internal/licensing"
	"sv-print/internal/receipt"
)

type API struct {
	queue    printing.PrintQueue
	registry *discovery.Registry
	version  string
	bus      *events.Bus
	lic      *licensing.State
}

func NewAPI(q printing.PrintQueue, r *discovery.Registry, version string, bus ...*events.Bus) *API {
	a := &API{
		queue:    q,
		registry: r,
		version:  version,
	}
	if len(bus) > 0 {
		a.bus = bus[0]
	}
	return a
}

func (a *API) WithLicense(s *licensing.State) *API {
	a.lic = s
	return a
}

func (a *API) isTrial() bool {
	return a.lic != nil && a.lic.IsTrial()
}

func (a *API) hasFeature(f string) bool {
	return a.lic == nil || a.lic.HasFeature(f)
}

func (a *API) enforcePrint(w http.ResponseWriter) (string, bool) {
	if a.lic == nil {
		return "", true
	}
	if !a.lic.AllowPrint() {
		s, c, m := mapDomainError(domainErrors.ErrLicenseQuotaExceeded)
		writeError(w, s, c, m)
		return "", false
	}
	return a.lic.Watermark(), true
}

func (a *API) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": a.version,
	})
}

func (a *API) InfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}
	resp := map[string]any{
		"name":         "SV Print Agent",
		"version":      a.version,
		"platform":     runtime.GOOS,
		"architecture": runtime.GOARCH,
	}
	if a.lic != nil {
		resp["tier"] = string(a.lic.Tier())
		resp["licensed"] = a.lic.IsActive()
	}
	writeJSON(w, http.StatusOK, resp)
}

type printersResponse struct {
	Printers []printer.Printer `json:"printers"`
}

func (a *API) PrintersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}
	printers, err := a.registry.DiscoverAll(r.Context())
	if err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}
	if printers == nil {
		printers = []printer.Printer{}
	}
	writeJSON(w, http.StatusOK, printersResponse{Printers: printers})
}

func (a *API) DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}
	printers, err := a.registry.DiscoverAll(r.Context())
	if err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}
	if printers == nil {
		printers = []printer.Printer{}
	}
	writeJSON(w, http.StatusOK, printersResponse{Printers: printers})
}

func (a *API) GetPrinterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}
	p, err := a.registry.Find(r.Context(), r.PathValue("id"))
	if err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (a *API) TestPrinterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}

	id := r.PathValue("id")
	if _, err := a.registry.Find(r.Context(), id); err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}

	wm, ok := a.enforcePrint(w)
	if !ok {
		return
	}

	doc := receipt.Document{
		Cut: true,
		Lines: []receipt.Line{
			{Text: "SV Print Agent", Style: receipt.Style{Bold: true, Align: "center"}},
			{Text: "Test receipt"},
		},
	}

	a.enqueue(w, id, doc.Build(wm).Bytes())
}

type printRequest struct {
	PrinterID string `json:"printer_id"`
	Format    string `json:"format"`
	Payload   string `json:"payload"`
}

type printResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

func (a *API) PrintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}

	if !a.hasFeature(license.FeatureRawPrint) {
		s, c, m := mapDomainError(domainErrors.ErrLicenseRequired)
		writeError(w, s, c, m)
		return
	}

	var req printRequest
	if !decodeBody(w, r, &req) {
		return
	}

	if req.PrinterID == "" || req.Payload == "" {
		s, c, m := mapDomainError(domainErrors.ErrInvalidPayload)
		writeError(w, s, c, m)
		return
	}
	if req.Format != "" && req.Format != "escpos" {
		s, c, m := mapDomainError(domainErrors.ErrUnsupportedProtocol)
		writeError(w, s, c, m)
		return
	}

	if _, err := a.registry.Find(r.Context(), req.PrinterID); err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}

	a.enqueue(w, req.PrinterID, []byte(req.Payload))
}

func (a *API) enqueue(w http.ResponseWriter, printerID string, payload []byte) {
	newJob := &job.PrintJob{
		ID:        job.NewID(),
		PrinterID: printerID,
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := a.queue.Enqueue(newJob); err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}

	writeJSON(w, http.StatusCreated, printResponse{JobID: newJob.ID, Status: string(job.StatusQueued)})
}

func (a *API) GetJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed.")
		return
	}
	j, err := a.queue.GetJob(r.PathValue("id"))
	if err != nil {
		s, c, m := mapDomainError(err)
		writeError(w, s, c, m)
		return
	}
	writeJSON(w, http.StatusOK, j)
}
