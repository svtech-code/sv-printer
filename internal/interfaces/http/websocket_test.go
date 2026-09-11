package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"sv-print/internal/application/discovery"
	"sv-print/internal/application/events"
	"sv-print/internal/application/printing"
)

func TestEventsHandler(t *testing.T) {
	bus := events.NewBus()
	api := NewAPI(printing.NewInMemoryQueue(), discovery.NewRegistry(), "1.0.0", bus)

	srv := httptest.NewServer(http.HandlerFunc(api.EventsHandler))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	bus.Publish(events.NewEvent(events.PrintStarted, map[string]any{"job_id": "job1", "printer_id": "p1"}))

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !strings.Contains(string(data), `"event":"print.started"`) {
		t.Errorf("received %q, want print.started event", data)
	}
	if !strings.Contains(string(data), `"job_id":"job1"`) {
		t.Errorf("received %q, want job_id field", data)
	}
}

func TestEventsHandlerWithoutBus(t *testing.T) {
	api := NewAPI(printing.NewInMemoryQueue(), discovery.NewRegistry(), "1.0.0")

	req := httptest.NewRequest("GET", "/api/v1/events", nil)
	rr := httptest.NewRecorder()
	api.EventsHandler(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
}
