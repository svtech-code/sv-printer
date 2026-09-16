package printing

import (
	"context"
	"testing"

	"sv-printer/internal/application/events"
	"sv-printer/internal/domain/job"
	"sv-printer/internal/infrastructure/transport"
)

func TestWorkerPublishesEvents(t *testing.T) {
	q := NewInMemoryQueue()
	_ = q.Enqueue(&job.PrintJob{ID: "job1", PrinterID: "p1", Payload: []byte("x")})

	bus := events.NewBus()
	ch, cancel := bus.Subscribe()
	defer cancel()

	tf := func(string) (transport.PrinterTransport, error) { return &mockTransport{}, nil }
	w := NewWorker(q, tf, bus)

	w.processNext(context.Background())

	if ev := <-ch; ev.Type != events.PrintStarted {
		t.Errorf("first event = %q, want %q", ev.Type, events.PrintStarted)
	}
	if ev := <-ch; ev.Type != events.PrintCompleted {
		t.Errorf("second event = %q, want %q", ev.Type, events.PrintCompleted)
	}
}

func TestWorkerPublishesFailedEvent(t *testing.T) {
	q := NewInMemoryQueue()
	_ = q.Enqueue(&job.PrintJob{ID: "job1", PrinterID: "p1", Payload: []byte("x")})

	bus := events.NewBus()
	ch, cancel := bus.Subscribe()
	defer cancel()

	tf := func(string) (transport.PrinterTransport, error) {
		return &mockTransport{openErr: context.DeadlineExceeded}, nil
	}
	w := NewWorker(q, tf, bus)

	w.processNext(context.Background())

	if ev := <-ch; ev.Type != events.PrintStarted {
		t.Errorf("first event = %q, want %q", ev.Type, events.PrintStarted)
	}
	if ev := <-ch; ev.Type != events.PrintFailed {
		t.Errorf("second event = %q, want %q", ev.Type, events.PrintFailed)
	}
}
