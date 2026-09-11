package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sv-print/internal/application/discovery"
	"sv-print/internal/application/printing"
	"sv-print/internal/domain/job"
	"sv-print/internal/domain/printer"
)

func newTestAPI() *API {
	q := printing.NewInMemoryQueue()
	r := discovery.NewRegistry()
	r.Register(discovery.NewStaticDiscoverer([]printer.Printer{
		{ID: "net-192.168.1.100:9100", Name: "Caja 1", Connection: printer.ConnectionNetwork, Protocol: printer.ProtocolEscpos},
	}))
	return NewAPI(q, r, "1.0.0")
}

func TestInfoHandlerPlatform(t *testing.T) {
	api := newTestAPI()
	req := httptest.NewRequest("GET", "/api/v1/info", nil)
	rr := httptest.NewRecorder()
	api.InfoHandler(rr, req)

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["platform"] == "" || resp["platform"] == "unknown" {
		t.Errorf("platform = %q, want real GOOS", resp["platform"])
	}
	if resp["architecture"] == "" || resp["architecture"] == "unknown" {
		t.Errorf("architecture = %q, want real GOARCH", resp["architecture"])
	}
}

func TestGetPrinterNotFound(t *testing.T) {
	api := newTestAPI()
	req := httptest.NewRequest("GET", "/api/v1/printers/does-not-exist", nil)
	req.SetPathValue("id", "does-not-exist")
	rr := httptest.NewRecorder()
	api.GetPrinterHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	assertErrorCode(t, rr, "PRINTER_NOT_FOUND")
}

func TestGetPrinterFound(t *testing.T) {
	api := newTestAPI()
	req := httptest.NewRequest("GET", "/api/v1/printers/net-192.168.1.100:9100", nil)
	req.SetPathValue("id", "net-192.168.1.100:9100")
	rr := httptest.NewRecorder()
	api.GetPrinterHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var p printer.Printer
	if err := json.NewDecoder(rr.Body).Decode(&p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.ID != "net-192.168.1.100:9100" {
		t.Errorf("id = %q", p.ID)
	}
	if p.Name != "Caja 1" {
		t.Errorf("name = %q", p.Name)
	}
}

func TestPrintHandlerValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
		want int
		code string
	}{
		{"empty printer_id", `{"payload":"x"}`, http.StatusBadRequest, "INVALID_PAYLOAD"},
		{"empty payload", `{"printer_id":"net-192.168.1.100:9100"}`, http.StatusBadRequest, "INVALID_PAYLOAD"},
		{"unsupported format", `{"printer_id":"net-192.168.1.100:9100","payload":"x","format":"pdf"}`, http.StatusBadRequest, "UNSUPPORTED_PROTOCOL"},
		{"printer not found", `{"printer_id":"nope","payload":"x"}`, http.StatusNotFound, "PRINTER_NOT_FOUND"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newTestAPI()
			req := httptest.NewRequest("POST", "/api/v1/print", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			api.PrintHandler(rr, req)
			if rr.Code != tt.want {
				t.Fatalf("status = %d, want %d", rr.Code, tt.want)
			}
			assertErrorCode(t, rr, tt.code)
		})
	}
}

func TestPrintHandlerSuccess(t *testing.T) {
	api := newTestAPI()
	req := httptest.NewRequest("POST", "/api/v1/print", strings.NewReader(`{"printer_id":"net-192.168.1.100:9100","payload":"x"}`))
	rr := httptest.NewRecorder()
	api.PrintHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}
	var resp printResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !strings.HasPrefix(resp.JobID, "job_") {
		t.Errorf("job_id = %q, want prefix job_", resp.JobID)
	}
	if resp.Status != string(job.StatusQueued) {
		t.Errorf("status = %q, want queued", resp.Status)
	}
}

func TestGetJobHandler(t *testing.T) {
	api := newTestAPI()
	_ = api.queue.Enqueue(&job.PrintJob{ID: "job_abc123", PrinterID: "net-192.168.1.100:9100", Payload: []byte("secret")})

	req := httptest.NewRequest("GET", "/api/v1/jobs/job_abc123", nil)
	req.SetPathValue("id", "job_abc123")
	rr := httptest.NewRecorder()
	api.GetJobHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var resp job.PrintJob
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != "job_abc123" {
		t.Errorf("id = %q", resp.ID)
	}
	if resp.Status != job.StatusQueued {
		t.Errorf("status = %q, want queued", resp.Status)
	}
	if string(resp.Payload) != "" {
		t.Errorf("payload should not be exposed, got %q", resp.Payload)
	}
}

func TestGetJobNotFound(t *testing.T) {
	api := newTestAPI()
	req := httptest.NewRequest("GET", "/api/v1/jobs/nope", nil)
	req.SetPathValue("id", "nope")
	rr := httptest.NewRecorder()
	api.GetJobHandler(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	assertErrorCode(t, rr, "JOB_NOT_FOUND")
}

func TestPrintersJSONKeysLowercase(t *testing.T) {
	api := newTestAPI()
	req := httptest.NewRequest("GET", "/api/v1/printers", nil)
	rr := httptest.NewRecorder()
	api.PrintersHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	body := rr.Body.String()
	if strings.Contains(body, `"ID"`) || strings.Contains(body, `"Name"`) {
		t.Errorf("response should use lowercase keys, got %s", body)
	}
	if !strings.Contains(body, `"id"`) || !strings.Contains(body, `"name"`) {
		t.Errorf("response missing lowercase keys, got %s", body)
	}
}

func assertErrorCode(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	var resp errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Error.Code != want {
		t.Errorf("error code = %q, want %q", resp.Error.Code, want)
	}
}
