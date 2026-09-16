package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sv-printer/internal/license"
	"sv-printer/internal/licensing"
)

func TestPrintHandlerTrialBlocked(t *testing.T) {
	api := newTestAPI().WithLicense(licensing.New(nil))
	req := httptest.NewRequest("POST", "/api/v1/print", strings.NewReader(`{"printer_id":"net-192.168.1.100:9100","payload":"x"}`))
	rr := httptest.NewRecorder()
	api.PrintHandler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusForbidden)
	}
	assertErrorCode(t, rr, "LICENSE_REQUIRED")
}

func TestEventsHandlerTrialBlocked(t *testing.T) {
	api := newTestAPI().WithLicense(licensing.New(nil))
	req := httptest.NewRequest("GET", "/api/v1/events", nil)
	rr := httptest.NewRecorder()
	api.EventsHandler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusForbidden)
	}
	assertErrorCode(t, rr, "LICENSE_REQUIRED")
}

func TestPrintReceiptTrialWatermark(t *testing.T) {
	api := newTestAPI().WithLicense(licensing.New(nil))
	body := `{"printer_id":"net-192.168.1.100:9100","lines":[{"text":"hello"}]}`
	req := httptest.NewRequest("POST", "/api/v1/print/receipt", strings.NewReader(body))
	rr := httptest.NewRecorder()
	api.PrintReceiptHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var resp printResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	j, err := api.queue.GetJob(resp.JobID)
	if err != nil {
		t.Fatalf("GetJob: %v", err)
	}
	if !strings.Contains(string(j.Payload), licensing.TrialWatermark) {
		t.Error("trial receipt should contain the watermark")
	}
}

func TestPrintReceiptQuotaExceeded(t *testing.T) {
	api := newTestAPI().WithLicense(licensing.New(nil))

	for i := 0; i < licensing.TrialDailyQuota; i++ {
		if !api.lic.AllowPrint() {
			t.Fatalf("print %d should be allowed", i)
		}
	}

	body := `{"printer_id":"net-192.168.1.100:9100","lines":[{"text":"x"}]}`
	req := httptest.NewRequest("POST", "/api/v1/print/receipt", strings.NewReader(body))
	rr := httptest.NewRecorder()
	api.PrintReceiptHandler(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusTooManyRequests)
	}
	assertErrorCode(t, rr, "LICENSE_QUOTA_EXCEEDED")
}

func TestLicensedNoWatermark(t *testing.T) {
	lic := &license.License{
		Product:  "sv-printer",
		Customer: "ACME",
		Tier:     "full",
		Features: []string{"raw_print", "websocket"},
	}
	api := newTestAPI().WithLicense(licensing.New(lic))

	body := `{"printer_id":"net-192.168.1.100:9100","lines":[{"text":"hello"}]}`
	req := httptest.NewRequest("POST", "/api/v1/print/receipt", strings.NewReader(body))
	rr := httptest.NewRecorder()
	api.PrintReceiptHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var resp printResponse
	_ = json.NewDecoder(rr.Body).Decode(&resp)

	j, _ := api.queue.GetJob(resp.JobID)
	if strings.Contains(string(j.Payload), licensing.TrialWatermark) {
		t.Error("licensed receipt should not contain the watermark")
	}
}
