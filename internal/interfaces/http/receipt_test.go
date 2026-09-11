package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrintReceiptHandler(t *testing.T) {
	api := newTestAPI()
	body := `{"printer_id":"net-192.168.1.100:9100","cut":true,"lines":[{"text":"SV TECH","style":{"bold":true,"align":"center"}}]}`
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
	if !strings.HasPrefix(resp.JobID, "job_") {
		t.Errorf("job_id = %q, want prefix job_", resp.JobID)
	}
	if resp.Status != "queued" {
		t.Errorf("status = %q, want queued", resp.Status)
	}
}

func TestPrintReceiptHandlerValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
		want int
		code string
	}{
		{"missing printer_id", `{"lines":[{"text":"x"}]}`, http.StatusBadRequest, "INVALID_PAYLOAD"},
		{"empty lines", `{"printer_id":"net-192.168.1.100:9100"}`, http.StatusBadRequest, "INVALID_PAYLOAD"},
		{"invalid align", `{"printer_id":"net-192.168.1.100:9100","lines":[{"text":"x","style":{"align":"diagonal"}}]}`, http.StatusBadRequest, "INVALID_PAYLOAD"},
		{"invalid size", `{"printer_id":"net-192.168.1.100:9100","lines":[{"text":"x","style":{"size":[99]}}]}`, http.StatusBadRequest, "INVALID_PAYLOAD"},
		{"unknown printer", `{"printer_id":"nope","lines":[{"text":"x"}]}`, http.StatusNotFound, "PRINTER_NOT_FOUND"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newTestAPI()
			req := httptest.NewRequest("POST", "/api/v1/print/receipt", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			api.PrintReceiptHandler(rr, req)
			if rr.Code != tt.want {
				t.Fatalf("status = %d, want %d", rr.Code, tt.want)
			}
			assertErrorCode(t, rr, tt.code)
		})
	}
}
