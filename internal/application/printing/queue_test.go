package printing

import (
	"testing"
	"time"

	"sv-printer/internal/domain/job"
)

func TestInMemoryQueue(t *testing.T) {
	q := NewInMemoryQueue()

	j := &job.PrintJob{
		ID:        "job1",
		PrinterID: "printer1",
		Payload:   []byte("test"),
		CreatedAt: time.Now(),
	}

	if err := q.Enqueue(j); err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	// Test GetJob
	retrieved, err := q.GetJob("job1")
	if err != nil {
		t.Fatalf("Failed to get job: %v", err)
	}
	if retrieved.Status != job.StatusQueued {
		t.Errorf("Expected status queued, got %v", retrieved.Status)
	}

	// Test UpdateStatus
	if err := q.UpdateStatus("job1", job.StatusPrinting); err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}
	retrieved, _ = q.GetJob("job1")
	if retrieved.Status != job.StatusPrinting {
		t.Errorf("Expected status printing, got %v", retrieved.Status)
	}

	// Test Dequeue
	dequeued := q.Dequeue()
	if dequeued == nil || dequeued.ID != "job1" {
		t.Fatalf("Dequeue failed")
	}

	// Queue should now be empty
	if empty := q.Dequeue(); empty != nil {
		t.Fatalf("Expected empty queue")
	}
}
