package printing

import (
	"context"
	"testing"
	"time"

	"sv-printer/internal/domain/job"
	"sv-printer/internal/infrastructure/transport"
)

type mockTransport struct {
	openErr  error
	writeErr error
	closeErr error
	written  []byte
}

func (m *mockTransport) Open(ctx context.Context) error { return m.openErr }
func (m *mockTransport) Write(ctx context.Context, data []byte) error {
	m.written = data
	return m.writeErr
}
func (m *mockTransport) Close() error { return m.closeErr }

func TestWorker(t *testing.T) {
	q := NewInMemoryQueue()
	j := &job.PrintJob{
		ID:        "job1",
		PrinterID: "printer1",
		Payload:   []byte("hello"),
		CreatedAt: time.Now(),
	}
	q.Enqueue(j)

	tf := func(printerID string) (transport.PrinterTransport, error) {
		return &mockTransport{}, nil
	}

	worker := NewWorker(q, tf)
	ctx, cancel := context.WithCancel(context.Background())

	go worker.Start(ctx)

	// Wait for job to process
	time.Sleep(500 * time.Millisecond)
	cancel() // Stop the worker

	retrieved, _ := q.GetJob("job1")
	if retrieved.Status != job.StatusCompleted {
		t.Errorf("Expected job to be completed, got %v", retrieved.Status)
	}
}
